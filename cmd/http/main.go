package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"

	apiv2 "github.com/bmstu-itstech/scriptum-back/internal/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()
	if cfgPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	cfg := config.MustLoad(cfgPath)

	l := logs.NewLogger(cfg.Logging)

	l.Debug(fmt.Sprintf("config: %+v", cfg))

	repos := postgres.MustNewRepository(cfg.Postgres, l)
	runner := docker.MustNewRunner(cfg.Docker, l)
	storage := local.MustNewStorage(cfg.Storage, l)

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	infra := app.Infra{
		BlueprintProvider:   repos,
		BlueprintRepository: repos,
		FileReader:          storage,
		FileUploader:        storage,
		JobProvider:         repos,
		JobPublisher:        jPub,
		JobRepository:       repos,
		Runner:              runner,
		UserProvider:        repos,
	}
	a := app.NewApp(infra, l)

	s := server.NewHTTPServer(cfg.HTTP, l, func(r chi.Router) http.Handler {
		return apiv2.HandlerFromMuxWithBaseURL(apiv2.NewServer(a), r, "/api/v2")
	})

	// start

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	errCh := make(chan error, 1)

	go func() {
		err := jSub.Listen(ctx, func(ctx2 context.Context, jobID string) error {
			return a.Commands.RunJob.Handle(ctx2, request.RunJob{JobID: jobID})
		})
		errCh <- err
	}()

	go func() {
		err := s.ListenAndServe()
		errCh <- err
	}()

	var err error
	select {
	case <-ctx.Done():
		l.Info("received cancel signal, gracefully shutting down")
		err = s.Shutdown(context.Background())
		if err != nil {
			l.Error("error shutting down http server", "error", err)
		}
	case err = <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			l.Error("listen error", slog.String("error", err.Error()))
			cancel()
		}
	}
}
