package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/sl"
)

const corsMaxAge = 300

func NewHTTPServer(cfg config.HTTP, l *slog.Logger, createHandler func(router chi.Router) http.Handler) *http.Server {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter, l, cfg)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/", createHandler(apiRouter))

	addr := fmt.Sprintf(":%d", cfg.Port)
	l.Info("starting HTTP server", slog.String("addr", addr))

	return &http.Server{
		Addr:    addr,
		Handler: rootRouter,
	}
}

func setMiddlewares(router *chi.Mux, l *slog.Logger, cfg config.HTTP) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(sl.NewLoggerMiddleware(l))
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router, cfg)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
	router.Use(jwtauth.MustNewMiddleware(cfg.JWTSecret).Handler)
}

func addCorsMiddleware(router *chi.Mux, cfg config.HTTP) {
	if len(cfg.CORSAllowOrigins) == 0 {
		// Оставить по умолчанию
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	})
	router.Use(corsMiddleware.Handler)
}
