package app

import (
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/command"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/app/query"
)

type Commands struct {
	CreateBox  command.CreateBoxHandler
	DeleteBox  command.DeleteBoxHandler
	RunJob     command.RunJobHandler
	StartJob   command.StartJobHandler
	UploadFile command.UploadFileHandler
}

type Queries struct {
	GetBox      query.GetBoxHandler
	GetBoxes    query.GetBoxesHandler
	GetJob      query.GetJobHandler
	GetJobs     query.GetJobsHandler
	SearchBoxes query.SearchBoxesHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}

type Infra struct {
	BoxProvider     ports.BoxProvider
	BoxRepo         ports.BoxRepository
	FileReader      ports.FileReader
	FileUploader    ports.FileUploader
	IsAdminProvider ports.IsAdminProvider
	JobProvider     ports.JobProvider
	JobPublisher    ports.JobPublisher
	JobRepository   ports.JobRepository
	Runner          ports.Runner
}

func NewApp(infra Infra, l *slog.Logger) *App {
	return &App{
		Commands: Commands{
			CreateBox:  command.NewCreateBoxHandler(infra.BoxRepo, infra.IsAdminProvider, l),
			DeleteBox:  command.NewDeleteBoxHandler(infra.BoxRepo, l),
			RunJob:     command.NewRunJobHandler(infra.Runner, infra.JobRepository, infra.FileReader, l),
			StartJob:   command.NewStartJobHandler(infra.BoxProvider, infra.JobRepository, infra.JobPublisher, l),
			UploadFile: command.NewUploadFileHandler(infra.FileUploader, l),
		},
		Queries: Queries{
			GetBox:      query.NewGetBoxHandler(infra.BoxProvider, l),
			GetBoxes:    query.NewGetBoxesHandler(infra.BoxProvider, l),
			GetJob:      query.NewGetJobHandler(infra.JobProvider, l),
			GetJobs:     query.NewGetJobsHandler(infra.JobProvider, l),
			SearchBoxes: query.NewSearchBoxesHandler(infra.BoxProvider, l),
		},
	}
}
