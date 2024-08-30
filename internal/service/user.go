package service

import (
	"fmt"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/repo"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	SignUp(userSignup *domain.UserSignup) error
}

type userService struct {
	users    map[uint]*domain.User
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *userService {
	return &userService{
		users:    make(map[uint]*domain.User),
		userRepo: userRepo,
	}
}

// check UserService interface implementation on compile time.
var _ UserService = (*userService)(nil)

func (u *userService) SignUp(userSignup *domain.UserSignup) error {
	if userSignup == nil {
		return fmt.Errorf("no user sign up details provided")
	}

	user, err := u.userRepo.Create(userSignup.Email, userSignup.Password)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	log.Debug().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("new user created")

	return nil
}
