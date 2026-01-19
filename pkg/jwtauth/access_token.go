package jwtauth

import (
	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims

	UID string `json:"uid"`
}
