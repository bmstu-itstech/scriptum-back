package delivery

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"

	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
)

type Server struct {
	app *app.Application
}

func (s Server) GetJobs(w http.ResponseWriter, r *http.Request) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	results, err := s.app.Queries.GetResults.JobResults(r.Context(), uint32(userID))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVResults(w, results)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) GetJobsSearch(w http.ResponseWriter, r *http.Request, params GetJobsSearchParams) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	results, err := s.app.Queries.SearchResultsSubstr.SearchResultBySubstr(r.Context(), uint32(userID), params.ScriptName)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVResults(w, results)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) GetJobsIdResult(w http.ResponseWriter, r *http.Request, id JobId) {
	_, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	result, err := s.app.Queries.SearchResultsID.SearchResultByID(r.Context(), scripts.JobID(id))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVResults(w, []app.ResultDTO{result})
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s Server) GetScripts(w http.ResponseWriter, r *http.Request) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	scripts, err := s.app.Queries.GetScripts.Scripts(r.Context(), uint32(userID))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVScripts(w, scripts)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) GetScriptsId(w http.ResponseWriter, r *http.Request, id int) {
	_, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	script, err := s.app.Queries.GetScriptByID.Script(r.Context(), id)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVScripts(w, []app.ScriptDTO{script})
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) GetScriptsSearch(w http.ResponseWriter, r *http.Request, params GetScriptsSearchParams) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	script, err := s.app.Queries.SearchScripts.Search(r.Context(), uint32(userID), params.ScriptName)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVScripts(w, script)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) GetUsersId(w http.ResponseWriter, r *http.Request, id UserId) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	user, err := s.app.Queries.GetUser.GetUser(r.Context(), uint32(userID), uint32(id))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVUsers(w, []app.UserDTO{user})
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	users, err := s.app.Queries.GetUsers.GetUsers(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = renderCSVUsers(w, users)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) DeleteScriptsId(w http.ResponseWriter, r *http.Request, id ScriptId) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = s.app.Commands.DeleteScript.DeleteScript(r.Context(), uint32(userID), uint32(id))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) PostScripts(w http.ResponseWriter, r *http.Request) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	PostScript := PostScriptsJSONBody{}
	if err := render.Decode(r, &PostScript); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	script, err := postScriptsJSONBodyToDTO(PostScript)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	_, err = s.app.Commands.CreateScript.CreateScript(r.Context(), uint32(userID), script)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) PutScriptsId(w http.ResponseWriter, r *http.Request, id ScriptId) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	PostScript := Script{}
	if err := render.Decode(r, &PostScript); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	script := app.ScriptDTO{
		ID:          uint32(*PostScript.ScriptId),
		Name:        *PostScript.ScriptName,
		Description: *PostScript.ScriptDesc,
		InFields:    convertFields(PostScript.InFields),
		OutFields:   convertFields(PostScript.OutFields),
		Owner:       uint32(PostScript.Owner),
		Visibility:  string(*PostScript.Visibility),
	}

	err = s.app.Commands.UpdateScript.Update(r.Context(), uint32(userID), script)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) PostScriptsIdStart(w http.ResponseWriter, r *http.Request, id ScriptId) {
	_, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}

	params := PostScriptsIdStartJSONBody{}
	if err := render.Decode(r, &params); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	input := app.ScriptRunDTO{
		ScriptID:     uint32(id),
		InParams:     convertValues(params.InParams),
		NeedToNotify: *params.NotifyByEmail,
	}

	err = s.app.Commands.StartJob.StartJob(r.Context(), input)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) DeleteUsersId(w http.ResponseWriter, r *http.Request, id UserId) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	err = s.app.Commands.DeleteUser.DeleteUser(r.Context(), uint32(userID), uint32(id))
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) PostUsers(w http.ResponseWriter, r *http.Request) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	PostUser := User{}
	if err := render.Decode(r, &PostUser); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	user := app.UserDTO{
		ID:       uint32(*PostUser.UserId),
		FullName: PostUser.FullName,
		Email:    string(PostUser.Email),
		IsAdmin:  PostUser.IsAdmin,
	}

	_, err = s.app.Commands.CreateUser.CreateUser(r.Context(), uint32(userID), user)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func (s Server) PutUsersId(w http.ResponseWriter, r *http.Request, id UserId) {
	userUUID, err := jwtauth.UserUUIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseUint(userUUID, 10, 32)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	PostUser := User{}
	if err := render.Decode(r, &PostUser); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	user := app.UserDTO{
		ID:       uint32(*PostUser.UserId),
		FullName: PostUser.FullName,
		Email:    string(PostUser.Email),
		IsAdmin:  PostUser.IsAdmin,
	}

	err = s.app.Commands.UpdateUser.UpdateUser(r.Context(), uint32(userID), user)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

}

func NewHTTPServer(app *app.Application) *Server {
	return &Server{app: app}
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	msg := err.Error()
	render.Status(r, code)
	render.JSON(w, r, Error{Message: &msg})
}

func renderCSVResults(w http.ResponseWriter, results []app.ResultDTO) error {
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	header := []string{
		"JobID",
		"UserID",
		"In",
		"Command",
		"StartedAt",
		"Code",
		"Out",
		"ErrorMes",
		"ClosedAt",
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, r := range results {
		inJSON, err := json.Marshal(r.Job.In)
		if err != nil {
			return err
		}
		outJSON, err := json.Marshal(r.Out)
		if err != nil {
			return err
		}

		errorMessage := ""
		if r.ErrorMes != nil {
			errorMessage = *r.ErrorMes
		}

		row := []string{
			strconv.FormatUint(uint64(r.Job.JobID), 10),
			strconv.FormatUint(uint64(r.Job.UserID), 10),
			string(inJSON),
			r.Job.Command,
			r.Job.StartedAt.Format(time.RFC3339),
			strconv.Itoa(r.Code),
			string(outJSON),
			errorMessage,
			r.ClosedAt.Format(time.RFC3339),
		}

		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}

func renderCSVScripts(w http.ResponseWriter, answers []app.ScriptDTO) error {
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	header := []string{
		"ID",
		"Name",
		"Description",
		"InFields",
		"OutFields",
		"Path",
		"Owner",
		"Visibility",
		"CreatedAt",
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, s := range answers {
		inFieldsJSON, err := json.Marshal(s.InFields)
		if err != nil {
			return err
		}
		outFieldsJSON, err := json.Marshal(s.OutFields)
		if err != nil {
			return err
		}

		row := []string{
			strconv.FormatUint(uint64(s.ID), 10),
			s.Name,
			s.Description,
			string(inFieldsJSON),
			string(outFieldsJSON),
			s.Path,
			strconv.FormatUint(uint64(s.Owner), 10),
			string(s.Visibility),
			s.CreatedAt.Format(time.RFC3339),
		}

		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}

func renderCSVUsers(w http.ResponseWriter, users []app.UserDTO) error {
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	header := []string{"ID", "FullName", "Email", "IsAdmin"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, u := range users {
		row := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.FullName,
			u.Email,
			strconv.FormatBool(u.IsAdmin),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return csvWriter.Error()
}
func convertFields(fields *[]Field) []app.FieldDTO {
	if fields == nil {
		return nil
	}
	result := make([]app.FieldDTO, 0, len(*fields))
	for _, f := range *fields {
		result = append(result, app.FieldDTO{
			Type:        *f.Type,
			Name:        *f.Name,
			Description: *f.Description,
			Unit:        *f.Unit,
		})
	}
	return result
}

func postScriptsJSONBodyToDTO(body PostScriptsJSONBody) (app.ScriptCreateDTO, error) {
	dto := app.ScriptCreateDTO{}

	if body.Script != nil {
		if body.Script.ScriptName != nil {
			dto.ScriptName = *body.Script.ScriptName
		}
		if body.Script.ScriptDesc != nil {
			dto.ScriptDescription = *body.Script.ScriptDesc
		}
		dto.InFields = convertFields(body.Script.InFields)
		dto.OutFields = convertFields(body.Script.OutFields)
	}

	if body.File != nil {
		bytes, err := body.File.Content.Bytes()
		if err != nil {
			return dto, err
		}

		dto.File = app.FileDTO{
			Name:     *body.File.Name,
			FileType: *body.File.FileType,
			Content:  string(bytes),
		}
	}
	return dto, nil
}

func convertValues(values []Value) []app.ValueDTO {
	result := make([]app.ValueDTO, 0, len(values))
	for _, v := range values {
		var dto app.ValueDTO
		if v.Type != nil {
			dto.Type = *v.Type
		}
		if v.Data != nil {
			dto.Data = *v.Data
		}
		result = append(result, dto)
	}
	return result
}
