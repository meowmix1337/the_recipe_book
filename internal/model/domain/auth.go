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
)

type JWTCustomClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	UUID   string `json:"uuid"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}
