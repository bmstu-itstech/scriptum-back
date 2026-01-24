package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5/request"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context, token value.Token) (value.UserID, error)
}

type Middleware struct {
	verifier  TokenVerifier
	extractor request.Extractor
}

func NewMiddleware(verifier TokenVerifier) *Middleware {
	return &Middleware{
		verifier: verifier,
		extractor: request.MultiExtractor{
			request.AuthorizationHeaderExtractor,
			request.ArgumentExtractor{"jwtToken"},
		},
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := m.extractor.ExtractToken(r)
		if errors.Is(err, request.ErrNoTokenInRequest) {
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{"message": err.Error()})
			return
		}

		uid, err := m.verifier.VerifyToken(r.Context(), value.Token(tokenString))
		if errors.Is(err, ports.ErrTokenInvalid) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{"message": fmt.Sprintf("token is invalid: %s", err.Error())})
			return
		} else if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{"message": "internal server error"})
		}

		r = r.WithContext(toContext(r.Context(), string(uid)))

		next.ServeHTTP(w, r)
	})
}
