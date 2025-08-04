package app

type Application struct {
	CreateScript  ScriptCreateUC
	DeleteScript  ScriptDeleteUC
	UpdateScript  ScriptUpdateUC
	SearchScript  SearchScriptsUC
	StartJob      JobStartUC
	GetJob        GetJobUC
	GetJobs       GetJobsUC
	GetScriptByID GetScriptUC
	GetScripts    GetScriptsUC
	SearchJob     SearchJobsUC
}
