package entity

import (
	"database/sql"
	"time"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
)

type RefreshToken struct {
	ID        uint         `db:"id"`
	UserID    uint         `db:"user_id"`
	Token     string       `db:"token"`
	ExpiresAt time.Time    `db:"expires_at"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (r *RefreshToken) ToDomain() *domain.RefreshToken {
	rt := new(domain.RefreshToken)
	rt.ID = r.ID
	rt.UserID = r.UserID
	rt.Token = r.Token
	rt.ExpiresAt = r.ExpiresAt
	rt.CreatedAt = r.CreatedAt
	rt.UpdatedAt = r.UpdatedAt
	if r.DeletedAt.Valid {
		rt.DeletedAt = r.DeletedAt.Time
	}

	return rt
}
