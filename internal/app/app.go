package app

import (
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/command"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/app/query"
)

type Commands struct {
	CreateBlueprint command.CreateBlueprintHandler
	DeleteBlueprint command.DeleteBlueprintHandler
	RunJob          command.RunJobHandler
	StartJob        command.StartJobHandler
	UploadFile      command.UploadFileHandler
}

type Queries struct {
	GetBlueprint     query.GetBlueprintHandler
	GetBlueprints    query.GetBlueprintsHandler
	GetJob           query.GetJobHandler
	GetJobs          query.GetJobsHandler
	SearchBlueprints query.SearchBlueprintsHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}

type Infra struct {
	BlueprintProvider   ports.BlueprintProvider
	BlueprintRepository ports.BlueprintRepository
	FileReader          ports.FileReader
	FileUploader        ports.FileUploader
	JobProvider         ports.JobProvider
	JobPublisher        ports.JobPublisher
	JobRepository       ports.JobRepository
	Runner              ports.Runner
	UserProvider        ports.UserProvider
}

func NewApp(infra Infra, l *slog.Logger) *App {
	return &App{
		Commands: Commands{
			CreateBlueprint: command.NewCreateBlueprintHandler(infra.BlueprintRepository, infra.UserProvider, l),
			DeleteBlueprint: command.NewDeleteBlueprintHandler(infra.BlueprintRepository, l),
			RunJob:          command.NewRunJobHandler(infra.Runner, infra.JobRepository, infra.FileReader, l),
			StartJob:        command.NewStartJobHandler(infra.BlueprintProvider, infra.JobRepository, infra.JobPublisher, l),
			UploadFile:      command.NewUploadFileHandler(infra.FileUploader, l),
		},
		Queries: Queries{
			GetBlueprint:     query.NewGetBlueprintHandler(infra.BlueprintProvider, l),
			GetBlueprints:    query.NewGetBlueprintsHandler(infra.BlueprintProvider, l),
			GetJob:           query.NewGetJobHandler(infra.JobProvider, l),
			GetJobs:          query.NewGetJobsHandler(infra.JobProvider, l),
			SearchBlueprints: query.NewSearchBlueprintsHandler(infra.BlueprintProvider, l),
		},
	}
}
