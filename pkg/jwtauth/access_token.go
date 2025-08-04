package jwtauth

import (
	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

type AccessTokenPayload struct {
	UserID int64
}
