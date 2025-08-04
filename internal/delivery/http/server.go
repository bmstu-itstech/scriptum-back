package scriptumapi

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
	"github.com/go-chi/render"
)

type Server struct {
	app *app.Application
}

func NewServer(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) GetJobs(w http.ResponseWriter, r *http.Request) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	results, err := s.app.GetJobs.Jobs(r.Context(), uint32(userID))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	resultsToSend := make([]Job, len(results))
	for i, res := range results {
		resultsToSend[i] = DTOToJobHttp(res)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resultsToSend)
}

func (s *Server) GetJobsSearch(w http.ResponseWriter, r *http.Request, params GetJobsSearchParams) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	results, err := s.app.SearchJob.Search(r.Context(), uint32(userID), string(params.State))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		httpError(w, r, fmt.Errorf("no results found"), http.StatusNotFound)
		return
	}
	resultsToSend := make([]Job, len(results))
	for i, res := range results {
		resultsToSend[i] = DTOToJobHttp(res)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resultsToSend)
}

func (s *Server) GetJobsIdResult(w http.ResponseWriter, r *http.Request, id JobId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	job, err := s.app.GetJob.Job(r.Context(), uint32(userID), int64(id))
	if err != nil {
		if errors.Is(err, fmt.Errorf("permission denied")) {
			httpError(w, r, err, http.StatusForbidden)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, DTOToJobHttp(job))
}

func (s *Server) GetJobsIdResultDownload(w http.ResponseWriter, r *http.Request, id JobId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	job, err := s.app.GetJob.Job(r.Context(), uint32(userID), int64(id))
	if err != nil {
		if errors.Is(err, fmt.Errorf("permission denied")) {
			httpError(w, r, err, http.StatusForbidden)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=job_result_export.csv")
	err = renderCSVJobResult(w, job)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
}

func (s *Server) GetScripts(w http.ResponseWriter, r *http.Request) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	scrs, err := s.app.GetScripts.Scripts(r.Context(), uint32(userID))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	res := make([]Script, len(scrs))
	for i, scr := range scrs {
		res[i] = DTOToScriptHttp(scr)
	}
	render.JSON(w, r, res)
	render.Status(r, http.StatusOK)
}

func (s *Server) PostScripts(w http.ResponseWriter, r *http.Request) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	req := PostScriptsJSONRequestBody{}
	if err := render.Decode(r, &req); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	script, file := *req.Script, *req.File
	if userID != script.Owner {
		httpError(w, r, fmt.Errorf("permission denied"), http.StatusForbidden)
		return
	}

	in := make([]app.FieldDTO, len(*script.InFields))
	for i, field := range *script.InFields {
		in[i] = FieldToDTOHttp(field)
	}
	out := make([]app.FieldDTO, len(*script.OutFields))
	for i, field := range *script.OutFields {
		out[i] = FieldToDTOHttp(field)
	}
	reqDto := app.ScriptCreateDTO{
		OwnerID:           userID,
		ScriptName:        *script.ScriptName,
		ScriptDescription: *script.ScriptDescription,
		File:              FileToDTOHttp(file),
		InFields:          in,
		OutFields:         out,
	}

	scriptID, err := s.app.CreateScript.CreateScript(r.Context(), reqDto)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	mes := struct {
		ScriptId int64  `json:"script_id"`
		Message  string `json:"message"`
	}{
		ScriptId: int64(scriptID),
		Message:  "Script created successfully",
	}
	render.JSON(w, r, mes)
}

func (s *Server) GetScriptsSearch(w http.ResponseWriter, r *http.Request, params GetScriptsSearchParams) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	scrs, err := s.app.SearchScript.Search(r.Context(), uint32(userID), params.ScriptName)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	if len(scrs) == 0 {
		httpError(w, r, fmt.Errorf("script not found"), http.StatusNotFound)
		return
	}
	render.Status(r, http.StatusOK)
	res := make([]Script, len(scrs))
	for i, scr := range scrs {
		res[i] = DTOToScriptHttp(scr)
	}
	render.JSON(w, r, res)
}

func (s *Server) DeleteScriptsId(w http.ResponseWriter, r *http.Request, id ScriptId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	err = s.app.DeleteScript.DeleteScript(r.Context(), uint32(userID), uint32(id))
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		if errors.Is(err, fmt.Errorf("permission denied")) {
			httpError(w, r, err, http.StatusForbidden)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, struct {
		Message string `json:"message"`
	}{
		Message: "Script deleted successfully"},
	)

}

func (s *Server) GetScriptsId(w http.ResponseWriter, r *http.Request, id ScriptId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	script, err := s.app.GetScriptByID.Script(r.Context(), int64(userID), int32(id))
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, DTOToScriptHttp(script))
}

func (s *Server) PutScriptsId(w http.ResponseWriter, r *http.Request, id ScriptId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	scri := Script{}
	if err := render.Decode(r, &scri); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	err = s.app.UpdateScript.UpdateScript(r.Context(), int64(userID), ScriptToDTOHttp(scri))
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, struct {
		Message string `json:"message"`
	}{
		Message: "Script updated successfully"},
	)
}

func (s *Server) PostScriptsIdStart(w http.ResponseWriter, r *http.Request, id ScriptId) {
	gotUserID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(gotUserID, 10, 64)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	req := PostScriptsIdStartJSONRequestBody{}
	if err := render.Decode(r, &req); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	in := make([]app.ValueDTO, len(req.InParams))
	for i, val := range req.InParams {
		in[i] = ValueToDTOHttp(val)
	}
	reqDto := app.ScriptRunDTO{
		ScriptID:     uint32(id),
		InParams:     in,
		NeedToNotify: *req.NotifyByEmail,
	}
	err = s.app.StartJob.StartJob(r.Context(), int64(userID), reqDto)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user not found")) {
			httpError(w, r, err, http.StatusNotFound)
			return
		}
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, struct {
		Message string `json:"message"`
	}{
		Message: "Job started successfully"},
	)
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	msg := err.Error()
	render.Status(r, code)
	render.JSON(w, r, Error{Message: &msg})
}

func DTOToJobHttp(job app.JobDTO) Job {
	expected := make([]Field, len(job.Expected))
	for i, field := range job.Expected {
		expected[i] = DTOToFieldHttp(field)
	}
	in := make([]Value, len(job.Input))
	for i, val := range job.Input {
		in[i] = DTOToValueHttp(val)
	}
	return Job{
		CreatedAt:    &job.CreatedAt,
		Expected:     &expected,
		FinishedAt:   job.FinishedAt,
		In:           &in,
		JobId:        &job.JobID,
		NeedToNotify: &job.NeedToNotify,
		Path:         &job.Url,
		ScriptId:     &job.ScriptID,
		Status:       (*Status)(&job.State),
		UserId:       &job.OwnerID,
	}
}

func DTOToFieldHttp(field app.FieldDTO) Field {
	return Field{
		Description: &field.Desc,
		Name:        &field.Name,
		Type:        &field.Type,
		Unit:        &field.Unit,
	}
}

func FieldToDTOHttp(field Field) app.FieldDTO {
	return app.FieldDTO{
		Desc: *field.Description,
		Name: *field.Name,
		Type: *field.Type,
		Unit: *field.Unit,
	}
}

func DTOToValueHttp(val app.ValueDTO) Value {
	return Value{
		Data: &val.Data,
		Type: &val.Type,
	}
}

func ValueToDTOHttp(val Value) app.ValueDTO {
	return app.ValueDTO{
		Data: *val.Data,
		Type: *val.Type,
	}
}

func DTOToScriptHttp(script app.ScriptDTO) Script {
	in := make([]Field, len(script.Input))
	for i, field := range script.Input {
		in[i] = DTOToFieldHttp(field)
	}
	out := make([]Field, len(script.Output))
	for i, field := range script.Output {
		out[i] = DTOToFieldHttp(field)
	}
	id := (ScriptId)(script.ID)
	return Script{
		CreatedAt:         &script.CreatedAt,
		InFields:          &in,
		OutFields:         &out,
		Owner:             script.OwnerID,
		Path:              &script.URL,
		ScriptDescription: &script.Desc,
		ScriptId:          &id,
		ScriptName:        &script.Name,
		Visibility:        (*Visibility)(&script.Visibility),
	}
}

func ScriptToDTOHttp(script Script) app.ScriptDTO {
	in := make([]app.FieldDTO, len(*script.InFields))
	for i, field := range *script.InFields {
		in[i] = FieldToDTOHttp(field)
	}
	out := make([]app.FieldDTO, len(*script.OutFields))
	for i, field := range *script.OutFields {
		out[i] = FieldToDTOHttp(field)
	}

	return app.ScriptDTO{
		ID:         int32(*script.ScriptId),
		Name:       *script.ScriptName,
		Desc:       *script.ScriptDescription,
		URL:        *script.Path,
		Visibility: string(*script.Visibility),
		Input:      in,
		Output:     out,
		OwnerID:    int64(script.Owner),
		CreatedAt:  *script.CreatedAt,
	}
}

func DTOToFileHttp(file app.FileDTO) File {
	result := File{}

	if file.Name != "" {
		name := file.Name
		result.Name = &name
	}

	if len(file.Content) > 0 {
		result.Content.InitFromBytes(file.Content, file.Name)
	}

	return result
}

func FileToDTOHttp(file File) app.FileDTO {
	c, _ := file.Content.Bytes()
	return app.FileDTO{
		Name:    *file.Name,
		Content: c,
	}
}

func renderCSVJobResult(w http.ResponseWriter, job app.JobDTO) error {
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	header := []string{
		"JobID",
		"OwnerID",
		"ScriptID",
		"Input",
		"Expected",
		"Url",
		"State",
		"CreatedAt",
		"FinishedAt",
		"NeedToNotify",
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	inputJSON, err := json.Marshal(job.Input)
	if err != nil {
		return err
	}
	expectedJSON, err := json.Marshal(job.Expected)
	if err != nil {
		return err
	}
	row := []string{
		strconv.FormatUint(uint64(job.JobID), 10),
		strconv.FormatUint(uint64(job.OwnerID), 10),
		strconv.FormatUint(uint64(job.ScriptID), 10),
		string(inputJSON),
		string(expectedJSON),
		job.Url,
		job.State,
		job.CreatedAt.Format(time.RFC3339),
		job.FinishedAt.Format(time.RFC3339),
		strconv.FormatBool(job.NeedToNotify),
	}
	if err := csvWriter.Write(row); err != nil {
		return err
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
