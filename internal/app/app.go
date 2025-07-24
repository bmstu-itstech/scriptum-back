package app

type Application struct {
	CreateScript  ScriptCreateUC
	DeleteScript  ScriptDeleteUC
	StartJob      JobStartUC
	GetScriptByID GetScriptUC
	GetScripts    GetScriptsUC
}
