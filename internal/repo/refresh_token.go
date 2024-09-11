package repo

import (
	"context"
	"time"

	"github.com/meowmix1337/go-core/db"
)

type RefreshTokenRepo interface {
	CreateRefreshToken(ctx context.Context, refreshToken string, userID uint) error
	DeleteRefreshToken(ctx context.Context, userID uint) error
}

type refreshTokenRepo struct {
	DB db.DB
}

func NewRefreshTokenRepo(db db.DB) *refreshTokenRepo {
	return &refreshTokenRepo{
		DB: db,
	}
}

var _ RefreshTokenRepo = (*refreshTokenRepo)(nil)

const (
	expiresAt = time.Hour * 24

	deleteTokenQuery = `UPDATE refresh_tokens SET deleted_at = $1 WHERE user_id = $2 AND deleted_at IS NULL`
)

func (r *refreshTokenRepo) CreateRefreshToken(ctx context.Context, refreshToken string, userID uint) error {
	err := r.DB.Transaction(ctx, func(ctx context.Context, tx db.Tx) error {
		_, err := tx.Exec(ctx, deleteTokenQuery, time.Now().UTC(), userID)
		if err != nil {
			return err
		}

		defaultExpireAt := time.Now().Add(expiresAt).UTC()
		query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ctx, query, userID, refreshToken, defaultExpireAt)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *refreshTokenRepo) DeleteRefreshToken(ctx context.Context, userID uint) error {
	_, err := r.DB.Exec(ctx, deleteTokenQuery, time.Now().UTC(), userID)
	if err != nil {
		return err
	}

	return err
}
