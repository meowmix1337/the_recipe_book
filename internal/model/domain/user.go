package domain

import (
	"errors"
	"time"
)

var (
	ErrNoCredentialsProvided = errors.New("no credentials provided")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrJWTGeneration         = errors.New("error generating jwt token")
	ErrUnauthorized          = errors.Join(ErrInvalidCredentials, ErrNoCredentialsProvided, ErrUserNotFound)
)

type UserSignup struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type UserCredentials struct {
	Email    string
	Password string
}

type User struct {
	ID        uint
	UUID      string
	Email     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	DeletedAt time.Time
}
