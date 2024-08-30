package entity

import "github.com/meowmix1337/the_recipe_book/internal/model/domain"

type User struct {
	ID       uint   `db:"id"`
	UUID     string `db:"uuid"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:    u.ID,
		UUID:  u.UUID,
		Email: u.Email,
	}
}
