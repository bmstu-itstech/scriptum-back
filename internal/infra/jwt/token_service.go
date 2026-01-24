package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type claims struct {
	UID string `json:"uid"`

	jwt.RegisteredClaims
}

type TokenService struct {
	cfg config.JWT
}

func NewTokenService(cfg config.JWT) (*TokenService, error) {
	if cfg.Secret == "" {
		return nil, errors.New("secret key required")
	}
	if cfg.AccessTTL == 0 {
		return nil, errors.New("access TTL required")
	}
	return &TokenService{cfg}, nil
}

func MustNewTokenService(cfg config.JWT) *TokenService {
	service, err := NewTokenService(cfg)
	if err != nil {
		panic(err)
	}
	return service
}

func (s *TokenService) GenerateToken(_ context.Context, userID value.UserID) (value.Token, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTTL)

	c := &claims{
		UID: string(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   string(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return value.Token(tokenString), nil
}

func (s *TokenService) VerifyToken(_ context.Context, token value.Token) (value.UserID, error) {
	parsedToken, err := jwt.ParseWithClaims(string(token), &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s", ports.ErrTokenInvalid, err.Error())
	}

	c, ok := parsedToken.Claims.(*claims)
	if !ok || !parsedToken.Valid {
		return "", ports.ErrTokenInvalid
	}

	return value.UserID(c.UID), nil
}
