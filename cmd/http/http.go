package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	httpapi "github.com/bmstu-itstech/scriptum-back/internal/delivery/http"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func main() {
	l := logs.NewLogger("prod")
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		panic(err)
	}

	emailNotifier, err := service.NewEmailNotifier(
		os.Getenv("EMAIL_TEMPLATE_PATH"),
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		port,
	)
	if err != nil {
		panic(err)
	}

	cfg := Config{
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	jobRepo := service.NewJobRepository(db)
	scriptRepo := service.NewScriptRepository(db)
	systemManager, err := service.NewSystemManager(".")
	if err != nil {
		panic(err)
	}

	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	dispatcher, err := service.NewLauncher(pubsub)
	if err != nil {
		panic(err)
	}

	handler, err := service.NewLaunchHandler(pubsub, logger)
	if err != nil {
		panic(err)
	}

	userProv, err := service.NewMockUserProvider()
	if err != nil {
		panic(err)
	}

	pythonLauncher, err := service.NewPythonLauncher("")
	if err != nil {
		panic(err)
	}

	usecase := app.NewJobRunUC(scriptRepo, jobRepo, pythonLauncher, emailNotifier, userProv, l)
	handler.Listen(ctx, usecase.Run)

	application := app.Application{
		CreateScript:  app.NewScriptCreateUC(scriptRepo, userProv, systemManager, l),
		DeleteScript:  app.NewScriptDeleteUC(scriptRepo, userProv, l, systemManager),
		UpdateScript:  app.NewScriptUpdateUC(scriptRepo, l),
		SearchScript:  app.NewSearchScriptsUC(scriptRepo, userProv, l),
		StartJob:      app.NewJobStartUC(scriptRepo, jobRepo, dispatcher, l),
		GetJob:        app.NewGetJobUC(jobRepo, userProv, l),
		GetJobs:       app.NewGetJobsUC(jobRepo, userProv, l),
		GetScriptByID: app.NewGetScript(scriptRepo, l),
		GetScripts:    app.NewGetScriptsUС(scriptRepo, userProv, l),
		SearchJob:     app.NewSearchJobsUC(jobRepo, userProv, l),
	}

	log.Println("Starting server")
	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpapi.HandlerFromMux(httpapi.NewServer(&application), router)
	})
}
