package endpoint

import "github.com/meowmix1337/the_recipe_book/internal/model/domain"

type UserSignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *UserSignupRequest) ToDomain() *domain.UserSignup {
	return &domain.UserSignup{
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserSignupError struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type UserCredentialsRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *UserCredentialsRequest) ToDomain() *domain.UserCredentials {
	return &domain.UserCredentials{
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
