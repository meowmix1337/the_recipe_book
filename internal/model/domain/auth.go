package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	JWTExpiration = time.Hour * 72
)

var (
	ErrUnableToVerifyClaim   = errors.New("unable to verify claims")
	ErrUnableToRetrieveToken = errors.New("unable to retrieve token")
	ErrRefreshTokenNotFound  = errors.New("refresh token does not exist")
)

type JWTCustomClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	UUID   string `json:"uuid"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        uint
	UserID    uint
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
