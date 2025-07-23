package main

import (
	"context"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi/v5"

	httpapi "github.com/bmstu-itstech/scriptum-back/internal/api/http"
	worker "github.com/bmstu-itstech/scriptum-back/internal/api/worker"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs"
	"github.com/bmstu-itstech/scriptum-back/pkg/server"
)

func main() {
	l := logs.DefaultLogger()
	context := context.Background()

	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	dispatcher, err := service.NewLauncher(pubsub)
	if err != nil {
		return
	}
	
	emailNotifier, err := service.NewEmailNotifier()
	if err != nil {
		return
	}

	jobRepo, err := service.NewJobRepo(context)
	if err != nil {
		return
	}

	resRepo, err := service.NewResRepo(context)
	if err != nil {
		return
	}

	pythonLauncher, err := service.NewPythonLauncher("", pubsub)
	if err != nil {
		return
	}

	fileManager, err := service.NewFileManager()
	if err != nil {
		return
	}

	scriptRepo, err := service.NewScriptRepo(context)
	if err != nil {
		return
	}

	userRepo, err := service.NewMockUserRepository()
	if err != nil {
		return
	}

	handler, err := worker.NewLaunchHandler(
		app.NewJobRunUC(jobRepo, pythonLauncher, emailNotifier, l),
		pubsub,
		logger,
	)
	if err != nil {
		return
	}

	handler.Listen(context)

	a := app.Application{
		Commands: app.Commands{
			CreateScript: app.NewScriptCreateUC(scriptRepo, l, userRepo, fileManager),
			CreateUser:   app.NewUserCreateUC(userRepo, l),
			DeleteScript: app.NewScriptDeleteUC(scriptRepo, userRepo, l, fileManager),
			DeleteUser:   app.NewUserDeleteUC(userRepo, l),
			StartJob:     app.NewJobStartUC(scriptRepo, jobRepo, dispatcher, userRepo, l),
			UpdateScript: app.NewScriptUpdateUC(scriptRepo, userRepo, l),
			UpdateUser:   app.NewUserUpdateUC(userRepo, l),
		},
		Queries: app.Queries{
			GetResults:          app.NewGetJobsUC(resRepo, l),
			GetScriptByID:       app.NewGetScript(scriptRepo, l),
			GetScripts:          app.NewGetScriptsUС(scriptRepo, userRepo, l),
			GetUser:             app.NewGetUserUC(userRepo, l),
			GetUsers:            app.NewGetUsersUC(userRepo, l),
			SearchResultsSubstr: app.NewSearchResultsSubstrUC(jobRepo, resRepo, userRepo, scriptRepo, l),
			SearchResultsID:     app.NewSearchResultIDUC(jobRepo, resRepo, userRepo, scriptRepo, l),
			SearchScripts:       app.NewScriptSearchUC(scriptRepo, userRepo, l),
		},
	}

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpapi.HandlerFromMux(httpapi.NewHTTPServer(&a), router)
	})
}
