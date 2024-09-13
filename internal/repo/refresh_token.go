package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/meowmix1337/go-core/db"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/entity"
)

type RefreshTokenRepo interface {
	CreateRefreshToken(ctx context.Context, refreshToken string, userID uint) error
	DeleteRefreshToken(ctx context.Context, userID uint) error

	ByRefreshToken(ctx context.Context, userID uint, refreshToken string) (*domain.RefreshToken, error)
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

func (r *refreshTokenRepo) ByRefreshToken(ctx context.Context, userID uint, refreshToken string) (*domain.RefreshToken, error) {
	query := `
	SELECT * 
		FROM refresh_tokens 
	WHERE refresh_tokens.token = $1 
		AND refresh_tokens.user_id = $2 
		AND refresh_tokens.deleted_at IS NULL`

	var refreshTokenEntity entity.RefreshToken
	err := r.DB.Get(ctx, &refreshTokenEntity, query, refreshToken, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return refreshTokenEntity.ToDomain(), nil
}
