package scripts

import "encoding/json"

type LaunchRequest struct {
	job          Job
	scriptFields []Field
	userEmail    Email
	needToNotify bool
}

func (r LaunchRequest) Job() Job {
	return r.job
}

func (r LaunchRequest) ScriptFields() []Field {
	return r.scriptFields
}

func (r LaunchRequest) UserEmail() Email {
	return r.userEmail
}

func (r LaunchRequest) NeedToNotify() bool {
	return r.needToNotify
}

func NewLaunchRequest(job Job, scriptFields []Field, userEmail Email, needToNotify bool) LaunchRequest {
	return LaunchRequest{
		job:          job,
		scriptFields: scriptFields,
		userEmail:    userEmail,
		needToNotify: needToNotify,
	}
}

func (r *LaunchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Job          Job     `json:"job"`
		ScriptFields []Field `json:"script_fields"`
		UserEmail    Email   `json:"user_email"`
		NeedToNotify bool    `json:"need_to_notify"`
	}{
		Job:          r.job,
		ScriptFields: r.scriptFields,
		UserEmail:    r.userEmail,
		NeedToNotify: r.needToNotify,
	})
}

func (r *LaunchRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		Job          Job     `json:"job"`
		ScriptFields []Field `json:"script_fields"`
		UserEmail    Email   `json:"user_email"`
		NeedToNotify bool    `json:"need_to_notify"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.job = aux.Job
	r.scriptFields = aux.ScriptFields
	r.userEmail = aux.UserEmail
	r.needToNotify = aux.NeedToNotify

	return nil
}
