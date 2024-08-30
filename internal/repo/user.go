package repo

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/entity"
)

type UserRepo interface {
	Create(email, password string) (*domain.User, error)
	ByEmail(email string) (*domain.User, error)
}

type userRepo struct {
	users map[uint]*entity.User
}

func NewUserRepository() *userRepo {
	return &userRepo{
		users: make(map[uint]*entity.User),
	}
}

var _ UserRepo = (*userRepo)(nil)

func (u *userRepo) Create(email, password string) (*domain.User, error) {
	magic := 8
	buf := make([]byte, magic)

	//nolint:errcheck // just testing
	rand.Read(buf) // Always succeeds, no need to check error
	id := binary.LittleEndian.Uint64(buf)

	newUser := &entity.User{
		ID:       uint(id),
		Email:    email,
		Password: password,
	}

	u.users[uint(id)] = newUser

	return newUser.ToDomain(), nil
}

func (u *userRepo) ByEmail(email string) (*domain.User, error) {
	for _, user := range u.users {
		if user.Email != email {
			continue
		}

		return user.ToDomain(), nil
	}

	return nil, sql.ErrNoRows
}
