package scripts

type LaunchRequest struct {
	Job          Job     `json:"job"`
	ScriptFields []Field `json:"script_fields"`
	UserEmail    Email   `json:"user_email"`
	NeedToNotify bool    `json:"need_to_notify"`
}

func (r LaunchRequest) GetJob() Job {
	return r.Job
}

func (r LaunchRequest) GetScriptFields() []Field {
	return r.ScriptFields
}

func (r LaunchRequest) GetUserEmail() Email {
	return r.UserEmail
}

func (r LaunchRequest) GetNeedToNotify() bool {
	return r.NeedToNotify
}

func NewLaunchRequest(job Job, scriptFields []Field, userEmail Email, needToNotify bool) LaunchRequest {
	return LaunchRequest{
		Job:          job,
		ScriptFields: scriptFields,
		UserEmail:    userEmail,
		NeedToNotify: needToNotify,
	}
}
