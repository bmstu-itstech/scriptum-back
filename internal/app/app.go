package app

import (
	"github.com/bmstu-itstech/scriptum-back/internal/app/command"
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
