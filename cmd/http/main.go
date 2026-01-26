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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	apiv2 "github.com/bmstu-itstech/scriptum-back/internal/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/bcrypt"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/jwt"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/postgres"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/watermill"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/sl"
)

const bcryptPasswordHasherCost = 12
const corsMaxAge = 300

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

	repos := postgres.MustNewRepository(cfg.Postgres)
	runner := docker.MustNewRunner(cfg.Docker, l)
	storage := local.MustNewStorage(cfg.Storage, l)
	hasher := bcrypt.NewPasswordHasher(bcryptPasswordHasherCost)
	tokenService := jwt.MustNewTokenService(cfg.JWT)

	jPub, jSub := watermill.NewJobPubSubGoChannels(l)

	infra := app.Infra{
		BlueprintProvider:   repos,
		BlueprintRepository: repos,
		FileReader:          storage,
		FileUploader:        storage,
		JobProvider:         repos,
		JobPublisher:        jPub,
		JobRepository:       repos,
		PasswordHasher:      hasher,
		Runner:              runner,
		TokenService:        tokenService,
		UserProvider:        repos,
		UserRepository:      repos,
	}
	a := app.NewApp(infra, l)

	root := chi.NewRouter()
	root.Use(middleware.RequestID)
	root.Use(middleware.RealIP)
	root.Use(sl.NewLoggerMiddleware(l))
	root.Use(middleware.Recoverer)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTP.CORSAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	})
	root.Use(corsMiddleware.Handler)
	root.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	root.Use(middleware.NoCache)
	root.Use(jwtauth.NewMiddleware(tokenService).Handler)
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: apiv2.HandlerFromMuxWithBaseURL(apiv2.NewServer(a), root, "/api/v2"),
	}

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
		l.Info("starting http server", slog.String("addr", s.Addr))
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
