package repo

import (
	"context"

	"github.com/meowmix1337/go-core/db"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/entity"
)

type UserRepo interface {
	Create(ctx context.Context, uuid, email, password string) error
	ByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepo struct {
	users map[uint]*entity.User
	DB    db.DB
}

func NewUserRepository(db db.DB) *userRepo {
	return &userRepo{
		users: make(map[uint]*entity.User),
		DB:    db,
	}
}

var _ UserRepo = (*userRepo)(nil)

func (u *userRepo) Create(ctx context.Context, uuid, email, password string) error {

	err := u.DB.Transaction(ctx, func(ctx context.Context, tx db.Tx) error {

		query := `INSERT INTO users (uuid, email) VALUES ($1, $2) RETURNING id`
		var userID int
		err := tx.Get(ctx, &userID, query, uuid, email)
		if err != nil {
			return err
		}

		query = `INSERT INTO user_passwords (user_id, password) VALUES ($1, $2)`
		_, err = tx.Exec(ctx, query, userID, password)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (u *userRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {

	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`

	var userEntity entity.User
	err := u.DB.Get_RO(ctx, &userEntity, query, email)
	if err != nil {
		return nil, err
	}

	return userEntity.ToDomain(), nil
}
