package endpoint

import "github.com/meowmix1337/the_recipe_book/internal/model/domain"

type UserSignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSignupRequest) ToDomain() *domain.UserSignup {
	return &domain.UserSignup{
		Email:    u.Email,
		Password: u.Password,
	}
}
