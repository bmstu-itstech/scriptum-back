package app

type Application struct {
	CreateScript  ScriptCreateUC
	DeleteScript  ScriptDeleteUC
	UpdateScript  ScriptUpdateUC
	SearchScript  SearchScriptsUC
	StartJob      JobStartUC
	GetJobs       GetJobsUC
	GetScriptByID GetScriptUC
	GetScripts    GetScriptsUC
	SearchJob     JobDTO
}
