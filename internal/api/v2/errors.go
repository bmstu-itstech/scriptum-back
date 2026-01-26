package apiv2

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

func renderInvalidInputError(w http.ResponseWriter, r *http.Request, iiErr domain.InvalidInputError, code int) {
	e := InvalidInputError{
		Code:    iiErr.Code,
		Message: nilOnEmpty(iiErr.Message),
	}
	render.Status(r, code)
	render.JSON(w, r, e)
}

func renderPlainError(w http.ResponseWriter, r *http.Request, inner error, code int) {
	e := PlainError{Message: inner.Error()}
	render.Status(r, code)
	render.JSON(w, r, e)
}

func renderInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("internal server error"))
}
