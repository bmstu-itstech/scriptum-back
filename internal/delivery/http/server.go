package scriptumapi

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const MaxFileSize = 10 << 20

type Server struct {
	app *app.Application
}

func NewServer(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) GetJobs(w http.ResponseWriter, r *http.Request) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	results, err := s.app.GetJobs.Jobs(r.Context(), uint32(userID))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	resultsToSend := make([]Result, len(results))
	for i, res := range results {
		resultsToSend[i] = DTOToJobHttp(res)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resultsToSend)
}

func (s *Server) GetJobsSearch(w http.ResponseWriter, r *http.Request, params GetJobsSearchParams) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	resultsToSend := make([]Result, len(results))
	for i, res := range results {
		resultsToSend[i] = DTOToJobHttp(res)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resultsToSend)
}

func (s *Server) GetJobsIdResult(w http.ResponseWriter, r *http.Request, id JobId) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	req := PostScriptsJSONRequestBody{}
	if err := render.Decode(r, &req); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	script := req
	fileID := script.FileId

	in := make([]app.FieldDTO, len(script.InFields))
	for i, field := range script.InFields {
		in[i] = FieldToDTOHttp(field)
	}
	out := make([]app.FieldDTO, len(script.OutFields))
	for i, field := range script.OutFields {
		out[i] = FieldToDTOHttp(field)
	}
	reqDto := app.ScriptCreateDTO{
		OwnerID:           userID,
		ScriptName:        script.ScriptName,
		ScriptDescription: script.ScriptDescription,
		FileID:            fileID,
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

func (s *Server) PostScriptsUpload(w http.ResponseWriter, r *http.Request) {
	_, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		httpError(w, r, fmt.Errorf("failed to read file: %v", err),
			http.StatusInternalServerError)
		return
	}

	reqDto := app.FileDTO{
		Name:   uuid.New().String(),
		Reader: bytes.NewReader(fileBytes),
	}

	fileID, err := s.app.CreateFile.CreateFile(r.Context(), reqDto)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = file.Close()
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusCreated)
	mes := struct {
		FileId  int64  `json:"file_id"`
		Message string `json:"message"`
	}{
		FileId:  int64(fileID),
		Message: "File uploaded successfully",
	}
	render.JSON(w, r, mes)
}

func (s *Server) GetScriptsSearch(w http.ResponseWriter, r *http.Request, params GetScriptsSearchParams) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	in := PutScriptsIdJSONRequestBody{}
	if err := render.Decode(r, &in); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	log.Println("ok1")

	err = s.app.UpdateScript.UpdateScript(r.Context(), int64(userID), id, ScriptDataToDTOHttp(in))
	log.Println("ok2")
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
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
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

func DTOToJobHttp(job app.JobDTO) Result {
	expected := make([]Field, len(job.Expected))
	for i, field := range job.Expected {
		expected[i] = DTOToFieldHttp(field)
	}
	in := make([]Value, len(job.Input))
	for i, val := range job.Input {
		in[i] = DTOToValueHttp(val)
	}
	j := Job{
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
	if job.JobResult == nil {
		return Result{
			Job: &j,
		}
	}

	out := make([]Value, len(job.JobResult.Output))
	for i, val := range job.JobResult.Output {
		out[i] = DTOToValueHttp(val)
	}
	c := int(job.JobResult.Code)
	return Result{
		Out:          &out,
		Code:         &c,
		ErrorMessage: job.JobResult.ErrMsg,
		Job:          &j,
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
		FileId:            &script.FileID,
		ScriptDescription: &script.Desc,
		ScriptId:          &id,
		ScriptName:        &script.Name,
		Visibility:        (*Visibility)(&script.Visibility),
	}
}

func ScriptToDTOHttp(script Script) app.ScriptDTO {
	in := make([]app.FieldDTO, 0, len(*script.InFields))
	for _, field := range *script.InFields {
		in = append(in, FieldToDTOHttp(field))
	}
	out := make([]app.FieldDTO, 0, len(*script.OutFields))
	for _, field := range *script.OutFields {
		out = append(out, FieldToDTOHttp(field))
	}

	return app.ScriptDTO{
		ID:         int32(*script.ScriptId),
		Name:       *script.ScriptName,
		Desc:       *script.ScriptDescription,
		FileID:     *script.FileId,
		Visibility: string(*script.Visibility),
		Input:      in,
		Output:     out,
		OwnerID:    int64(script.Owner),
		CreatedAt:  *script.CreatedAt,
	}
}

func ScriptDataToDTOHttp(data ScriptUpdateData) app.ScriptUpdateDTO {
	var in, out []app.FieldDTO
	if data.InFields != nil {
		in = make([]app.FieldDTO, 0, len(*data.InFields))
		for _, field := range *data.InFields {
			in = append(in, FieldToDTOHttp(field))
		}
	} else {
		in = nil
	}
	if data.OutFields != nil {
		out = make([]app.FieldDTO, 0, len(*data.OutFields))
		for _, field := range *data.OutFields {
			out = append(out, FieldToDTOHttp(field))
		}
	} else {
		out = nil
	}

	sName := ""
	if data.ScriptName != nil {
		sName = *data.ScriptName
	}
	sDesc := ""
	if data.ScriptDescription != nil {
		sDesc = *data.ScriptDescription
	}
	return app.ScriptUpdateDTO{
		InFields:          in,
		OutFields:         out,
		ScriptName:        sName,
		ScriptDescription: sDesc,
	}
}

func renderCSVJobResult(w http.ResponseWriter, job app.JobDTO) error {
	if job.JobResult == nil {
		return fmt.Errorf("no results to render")
	}

	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	row := make([]string, len(job.JobResult.Output))
	for i, val := range job.JobResult.Output {
		row[i] = val.Data
	}
	if err := csvWriter.Write(row); err != nil {
		return err
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
