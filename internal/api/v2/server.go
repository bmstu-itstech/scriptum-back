package apiv2

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"

	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/pkg/jwtauth"
)

const maxMultipartFormDataSize = 32 << 20 // 32 Kb

var ErrAuthorizationRequired = errors.New("authorization required")

type Server struct {
	app *app.App
}

func NewServer(app *app.App) *Server {
	return &Server{app}
}

func (s *Server) StartJob(w http.ResponseWriter, r *http.Request, id string) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	req := StartJobRequest{}
	if err := render.Decode(r, &req); err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	jobID, err := s.app.Commands.StartJob.Handle(r.Context(), startJobRequestToDTO(req, uid, id))
	var iiErr domain.InvalidInputError
	if errors.As(err, &iiErr) {
		renderInvalidInputError(w, r, iiErr, http.StatusBadRequest)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := StartJobResponse{JobID: jobID}
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, res)
}

func (s *Server) UploadFile(w http.ResponseWriter, r *http.Request) {
	_, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(maxMultipartFormDataSize)
	if err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	f, header, err := r.FormFile("attachment")
	if err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	fileID, err := s.app.Commands.UploadFile.Handle(r.Context(), request.UploadFileRequest{
		Name:   header.Filename,
		Reader: f,
	})
	if err != nil {
		renderPlainError(w, r, err, http.StatusInternalServerError)
		return
	}

	// TODO: return file size
	res := UploadFileResponse{FileID: fileID, Size: 0}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}

func (s *Server) GetBlueprints(w http.ResponseWriter, r *http.Request) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	bs, err := s.app.Queries.GetBlueprints.Handle(r.Context(), request.GetBlueprints{ActorID: uid})
	if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := blueprintsToAPI(bs)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) CreateBlueprint(w http.ResponseWriter, r *http.Request) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	req := CreateBlueprintRequest{}
	if err := render.Decode(r, &req); err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := s.app.Commands.CreateBlueprint.Handle(r.Context(), createBlueprintToDTO(req, uid))
	var iiErr domain.InvalidInputError
	if errors.As(err, &iiErr) {
		renderInvalidInputError(w, r, iiErr, http.StatusBadRequest)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := CreateBlueprintResponse{BlueprintID: id}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}

func (s *Server) SearchBlueprints(w http.ResponseWriter, r *http.Request, params SearchBlueprintsParams) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	bs, err := s.app.Queries.SearchBlueprints.Handle(r.Context(), request.SearchBlueprints{ActorID: uid, Name: params.Name})
	if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := blueprintsToAPI(bs)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) DeleteBlueprint(w http.ResponseWriter, r *http.Request, id string) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	err := s.app.Commands.DeleteBlueprint.Handle(r.Context(), request.DeleteBlueprint{BlueprintID: id, ActorID: uid})
	if errors.Is(err, ports.ErrBlueprintNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	render.NoContent(w, r)
}

func (s *Server) GetBlueprint(w http.ResponseWriter, r *http.Request, id string) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	b, err := s.app.Queries.GetBlueprint.Handle(r.Context(), request.GetBlueprint{BlueprintID: id, ActorID: uid})
	if errors.Is(err, ports.ErrBlueprintNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := blueprintToAPI(b)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) GetJobs(w http.ResponseWriter, r *http.Request, params GetJobsParams) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	js, err := s.app.Queries.GetJobs.Handle(r.Context(), request.GetJobs{ActorID: uid, State: (*string)(params.State)})
	if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := jobsToAPI(js)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) GetJob(w http.ResponseWriter, r *http.Request, id string) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	j, err := s.app.Queries.GetJob.Handle(r.Context(), request.GetJob{ActorID: uid, JobID: id})
	if errors.Is(err, ports.ErrJobNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := jobToAPI(j)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	req := LoginRequest{}
	if err := render.Decode(r, &req); err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	token, err := s.app.Commands.Login.Handle(r.Context(), request.Login{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, domain.ErrInvalidCredentials) {
		renderPlainError(w, r, err, http.StatusUnauthorized)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := LoginResponse{AccessToken: token.AccessToken}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	user, err := s.app.Queries.GetUser.Handle(r.Context(), request.GetUser{ActorID: uid, UserID: id})
	if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if errors.Is(err, ports.ErrUserNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := userToAPI(user)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) GetUserMe(w http.ResponseWriter, r *http.Request) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	user, err := s.app.Queries.GetUser.Handle(r.Context(), request.GetUser{ActorID: uid, UserID: uid})
	if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if errors.Is(err, ports.ErrUserNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := userToAPI(user)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	uid, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	users, err := s.app.Queries.GetUsers.Handle(r.Context(), request.GetUsers{ActorID: uid})
	if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := usersToAPI(users)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	actorID, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	req := CreateUserRequest{}
	if err := render.Decode(r, &req); err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	appRes, err := s.app.Commands.CreateUser.Handle(r.Context(), createUserToDTO(req, actorID))
	var iiErr domain.InvalidInputError
	if errors.As(err, &iiErr) {
		renderInvalidInputError(w, r, iiErr, http.StatusBadRequest)
		return
	} else if errors.Is(err, ports.ErrUserAlreadyExists) {
		renderPlainError(w, r, err, http.StatusConflict)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := CreateUserResponse{UserID: appRes.UID}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	actorID, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	err := s.app.Commands.DeleteUser.Handle(r.Context(), request.DeleteUser{
		ActorID: actorID,
		UID:     id,
	})
	if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if errors.Is(err, ports.ErrUserNotFound) {
		renderPlainError(w, r, err, http.StatusNotFound)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (s *Server) PatchUser(w http.ResponseWriter, r *http.Request, id string) {
	actorID, ok := jwtauth.FromContext(r.Context())
	if !ok {
		renderPlainError(w, r, ErrAuthorizationRequired, http.StatusUnauthorized)
		return
	}

	req := PatchUserRequest{}
	if err := render.Decode(r, &req); err != nil {
		renderPlainError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := s.app.Commands.UpdateUser.Handle(r.Context(), patchUserToDTO(req, actorID, id))
	var iiErr domain.InvalidInputError
	if errors.As(err, &iiErr) {
		renderInvalidInputError(w, r, iiErr, http.StatusBadRequest)
		return
	} else if errors.Is(err, domain.ErrPermissionDenied) {
		renderPlainError(w, r, err, http.StatusForbidden)
		return
	} else if err != nil {
		renderInternalServerError(w, r)
		return
	}

	res := userToAPI(user)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}
