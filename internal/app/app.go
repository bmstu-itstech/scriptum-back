package app

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateScript ScriptCreateUC
	CreateUser   UserCreateUC
	DeleteScript ScriptDeleteUC
	DeleteUser   UserDeleteUC
	StartJob     JobStartUC
	UpdateScript ScriptUpdateUC
	UpdateUser   UserUpdateUC
}

type Queries struct {
	GetResults          GetJobResultsUC
	GetScriptByID       GetScriptUC
	GetScripts          GetScriptsUC
	GetUser             GetUserUC
	GetUsers            GetUsersUC
	SearchResultsSubstr SearchResultsSubstrUC
	SearchResultsID     SearchResultIDUC
	SearchScripts       ScriptSearchUC
}
