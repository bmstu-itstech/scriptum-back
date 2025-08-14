package app

type Application struct {
	CreateScript  ScriptCreateUC
	DeleteScript  ScriptDeleteUC
	SearchScript  SearchScriptsUC
	StartJob      JobStartUC
	GetJob        GetJobUC
	GetJobs       GetJobsUC
	GetScriptByID GetScriptUC
	GetScripts    GetScriptsUC
	SearchJob     SearchJobsUC
	CreateFile    FileCreateUC
}
