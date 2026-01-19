package jwtauth

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

type Middleware struct {
	secret []byte
}

func NewMiddleware(secret string) (*Middleware, error) {
	if secret == "" {
		return nil, errors.New("secret key is required")
	}
	return &Middleware{
		secret: []byte(secret),
	}, nil
}

func MustNewMiddleware(secret string) *Middleware {
	m, err := NewMiddleware(secret)
	if err != nil {
		panic(err)
	}
	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims accessTokenClaims
		token, err := request.ParseFromRequest(
			r,
			request.MultiExtractor{
				request.AuthorizationHeaderExtractor,
				request.ArgumentExtractor{"jwtToken"},
			},
			func(_ *jwt.Token) (_ interface{}, _ error) {
				return m.secret, nil
			},
			request.WithClaims(&claims),
		)
		if errors.Is(err, request.ErrNoTokenInRequest) {
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{"message": err.Error()})
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{"message": "token is invalid"})
		}

		r = r.WithContext(toContext(r.Context(), claims.UID))

		next.ServeHTTP(w, r)
	})
}
