package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/bmstu-itstech/scriptum-back/internal/api"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/mock"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}
	cfg := config.MustLoad(cfgPath)

	l := logs.NewLogger(cfg.Logging)
	repos := postgres.MustNewRepository(cfg.Postgres, l)
	runner := docker.MustNewRunner(cfg.Docker, l)
	storage := local.MustNewStorage(cfg.Storage, l)
	mockIAP := mock.NewIsAdminProvider()

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	infra := app.Infra{
		BoxProvider:     repos,
		BoxRepo:         repos,
		FileReader:      storage,
		FileUploader:    storage,
		IsAdminProvider: mockIAP,
		JobProvider:     repos,
		JobPublisher:    jPub,
		JobRepository:   repos,
		Runner:          runner,
	}
	a := app.NewApp(infra, l)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	errCh := make(chan error, 1)

	go func() {
		err2 := jSub.Listen(ctx, func(ctx2 context.Context, jobID string) error {
			return a.Commands.RunJob.Handle(ctx2, request.RunJob{JobID: jobID})
		})
		errCh <- err2
	}()

	err := server.RunGRPCServerOnAddr(ctx, l, fmt.Sprintf(":%d", cfg.GRPC.Port), func(s *grpc.Server) {
		api.RegisterBoxService(s, a, l)
		api.RegisterFileService(s, a, l)
		api.RegisterJobService(s, a, l)
	})
	if err != nil {
		l.Error("failed to start grpc server", slog.String("error", err.Error()))
		errCh <- err
	}

	select {
	case <-ctx.Done():
		l.Info("received cancel signal, gracefully shutting down")
	case err = <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			l.Error("listen error", slog.String("error", err.Error()))
			cancel()
		}
	}
}
