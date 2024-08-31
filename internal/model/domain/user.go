package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoCredentialsProvided = errors.New("no credentials provided")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrJWTGeneration         = errors.New("error generating jwt token")
	ErrUnauthorized          = errors.Join(ErrInvalidCredentials, ErrNoCredentialsProvided, ErrUserNotFound)
)

const (
	JWTExpiration = time.Hour * 72
)

type JWTCustomClaims struct {
	Email string `json:"email"`
	UUID  string `json:"uuid"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type UserSignup struct {
	Email    string
	Password string
}

type UserCredentials struct {
	Email    string
	Password string
}

type User struct {
	ID       uint
	UUID     string
	Email    string
	Password string
}
