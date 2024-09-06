package entity

import (
	"database/sql"
	"time"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
)

type User struct {
	ID        uint         `db:"id"`
	UUID      string       `db:"uuid"`
	Email     string       `db:"email"`
	FirstName string       `db:"first_name"`
	LastName  string       `db:"last_name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type UserWithPassword struct {
	ID        uint         `db:"id"`
	Password  string       `db:"password"`
	UUID      string       `db:"uuid"`
	Email     string       `db:"email"`
	FirstName string       `db:"first_name"`
	LastName  string       `db:"last_name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (u *UserWithPassword) ToDomain() *domain.User {
	user := new(domain.User)
	user.ID = u.ID
	user.Password = u.Password
	user.UUID = u.UUID
	user.Email = u.Email
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.CreatedAt = u.CreatedAt
	if u.DeletedAt.Valid {
		user.DeletedAt = u.DeletedAt.Time
	}

	return user
}

func (u *User) ToDomain() *domain.User {
	user := new(domain.User)
	user.ID = u.ID
	user.UUID = u.UUID
	user.Email = u.Email
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.CreatedAt = u.CreatedAt
	if u.DeletedAt.Valid {
		user.DeletedAt = u.DeletedAt.Time
	}

	return user
}
