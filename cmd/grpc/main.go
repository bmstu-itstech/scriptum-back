package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/command"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/query"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/mock"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
)

const RunnerTimeout = 15 * time.Minute
const LocalStorageBasePath = "/tmp/scriptum-back"

func connectDB() (*sqlx.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		return nil, errors.New("DATABASE_URI must be set")
	}
	return sqlx.Connect("postgres", uri)
}

func main() {
	l := logs.DefaultLogger()

	db, err := connectDB()
	if err != nil {
		l.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repos := postgres.NewRepository(db, l)
	runner := docker.MustNewRunner(l)
	storage := local.MustNewStorage(LocalStorageBasePath, l)
	mockIAP := mock.NewIsAdminProvider()

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	a := app.App{
		Commands: app.Commands{
			CreateBox:  command.NewCreateBoxHandler(repos, mockIAP, l),
			DeleteBox:  command.NewDeleteBoxHandler(repos, l),
			RunJob:     command.NewRunJobHandler(runner, repos, storage, l),
			StartJob:   command.NewStartJobHandler(repos, jPub, l),
			UploadFile: command.NewUploadFileHandler(storage, l),
		},
		Queries: app.Queries{
			GetBox:      query.NewGetBoxHandler(repos, l),
			GetBoxes:    query.NewGetBoxesHandler(repos, l),
			GetJob:      query.NewGetJobHandler(repos, l),
			GetJobs:     query.NewGetJobsHandler(repos, l),
			SearchBoxes: query.NewSearchBoxesHandler(repos, l),
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	errCh := make(chan error, 1)

	go func() {
		err2 := jSub.Listen(ctx, func(ctx2 context.Context, jobID string) error {
			ctx2, cancel2 := context.WithTimeout(ctx, RunnerTimeout)
			defer cancel2()
			return a.Commands.RunJob.Handle(ctx2, request.RunJob{JobID: jobID})
		})
		errCh <- err2
	}()

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
