package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	httpapi "github.com/bmstu-itstech/scriptum-back/internal/delivery/http"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/sl"
	"github.com/bmstu-itstech/scriptum-back/pkg/server"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	l := logs.NewLogger("prod")
	ctx := context.Background()

	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatalf("failed get port: %s", err.Error())
	}

	emailNotifier, err := service.NewEmailNotifier(
		os.Getenv("EMAIL_TEMPLATE_PATH"),
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		port,
	)
	if err != nil {
		log.Fatalf("failed get email notifier: %s", err.Error())
	}

	if os.Getenv("DATABASE_URI") == "" {
		log.Fatalf("DATABASE_URI is empty")
	}
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Fatalf("failed connect postgres: %s", err.Error())
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	jobRepo := service.NewJobRepository(db)
	scriptRepo := service.NewScriptRepository(db)
	fileRepo := service.NewFileRepository(db)
	systemManager, err := service.NewSystemManager(".", 1024*1024*7)
	if err != nil {
		log.Fatalf("failed get system manager: %s", err.Error())
	}

	logger := sl.NewWatermillLoggerAdapter(l)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	dispatcher, err := service.NewLauncher(pubsub)
	if err != nil {
		log.Fatalf("failed get launcher: %s", err.Error())
	}

	handler, err := service.NewLaunchHandler(pubsub, logger)
	if err != nil {
		log.Fatalf("failed get handler: %s", err.Error())
	}

	userProv, err := service.NewMockUserProvider()
	if err != nil {
		log.Fatalf("failed get user provider: %s", err.Error())
	}

	pythonLauncher, err := service.NewPythonLauncher(os.Getenv("PYTHON_INTERPRETER"))
	if err != nil {
		log.Fatalf("failed get python launcher: %s", err.Error())
	}

	usecase := app.NewJobRunUC(scriptRepo, jobRepo, pythonLauncher, emailNotifier, userProv, l)
	handler.Listen(ctx, usecase.Run)

	application := app.Application{
		CreateScript:  app.NewScriptCreateUC(scriptRepo, userProv, fileRepo, systemManager, l),
		DeleteScript:  app.NewScriptDeleteUC(scriptRepo, userProv, systemManager, fileRepo, l),
		UpdateScript:  app.NewScriptUpdateUC(scriptRepo, l),
		SearchScript:  app.NewSearchScriptsUC(scriptRepo, userProv, l),
		StartJob:      app.NewJobStartUC(scriptRepo, fileRepo, jobRepo, dispatcher, l),
		GetJob:        app.NewGetJobUC(jobRepo, userProv, l),
		GetJobs:       app.NewGetJobsUC(jobRepo, userProv, l),
		GetScriptByID: app.NewGetScript(scriptRepo, l),
		GetScripts:    app.NewGetScriptsUС(scriptRepo, userProv, l),
		SearchJob:     app.NewSearchJobsUC(jobRepo, userProv, l),
		CreateFile:    app.NewFileCreateUC(userProv, fileRepo, systemManager, l),
	}

	l.Info("Starting server")
	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpapi.HandlerFromMux(httpapi.NewServer(&application), router)
	})
}
