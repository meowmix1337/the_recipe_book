package service

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	SignUp(userSignup *domain.UserSignup) error
}

type userService struct {
	users map[uint]*domain.User
}

func NewUserService() *userService {
	return &userService{
		users: make(map[uint]*domain.User),
	}
}

// check UserService interface implementation on compile time.
var _ UserService = (*userService)(nil)

func (u *userService) SignUp(userSignup *domain.UserSignup) error {
	if userSignup == nil {
		return fmt.Errorf("no user sign up details provided")
	}

	magic := 8
	buf := make([]byte, magic)

	//nolint:errcheck // just testing
	rand.Read(buf) // Always succeeds, no need to check error
	id := binary.LittleEndian.Uint64(buf)

	newUser := &domain.User{
		ID:    uint(id),
		Email: userSignup.Email,
	}

	u.users[uint(id)] = newUser

	log.Debug().
		Uint("user_id", uint(id)).
		Str("email", newUser.Email).
		Msg("new user created")

	return nil
}
