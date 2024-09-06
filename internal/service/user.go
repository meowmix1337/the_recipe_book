package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/repo"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, userSignup *domain.UserSignup) error
	Login(ctx context.Context, userCredentials *domain.UserCredentials) (string, error)
	Logout(ctx context.Context, userID uint) error

	ByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userService struct {
	*BaseService

	userRepo repo.UserRepo
}

func NewUserService(cfg config.Config, userRepo repo.UserRepo) *userService {
	return &userService{
		BaseService: &BaseService{Config: cfg},
		userRepo:    userRepo,
	}
}

// check UserService interface implementation on compile time.
var _ UserService = (*userService)(nil)

func (u *userService) SignUp(ctx context.Context, userSignup *domain.UserSignup) error {
	if userSignup == nil {
		return fmt.Errorf("no user sign up details provided")
	}

	// check if email exists already
	user, err := u.ByEmail(ctx, userSignup.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	if user != nil {
		return domain.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userSignup.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Err(err).Msg("error generating hash password")
		return err
	}

	// generate uuid
	uuid := u.GenerateUUIDHash("user")

	err = u.userRepo.Create(ctx, uuid, userSignup, string(hashedPassword))
	if err != nil {
		log.Err(err).Msg("error creating user")
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (u *userService) Login(ctx context.Context, userCredentials *domain.UserCredentials) (string, error) {
	if userCredentials == nil {
		log.Err(domain.ErrNoCredentialsProvided).Msg("no credentials were provided")
		return "", fmt.Errorf("no user login credentials provided: %w", domain.ErrNoCredentialsProvided)
	}

	user, err := u.ByEmailWithPassword(ctx, userCredentials.Email)
	if err != nil {
		return "", err
	}

	// Compare the stored hash with the provided password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentials.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Err(domain.ErrInvalidCredentials).Msg("invalid credentials")
			return "", fmt.Errorf("invalid credentials: %w", domain.ErrInvalidCredentials)
		}
		log.Err(err).Msg("error comparing password")
		return "", err
	}

	claims := &domain.JWTCustomClaims{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UUID:      user.UUID,
		Admin:     false,
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "Recipe App",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(domain.JWTExpiration)),
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(u.Config.GetJWTSecret()))
	if err != nil {
		log.Err(err).Msg("error generating JWT token")
		return "", domain.ErrJWTGeneration
	}

	// TODO: generate refresh token

	return tokenString, nil
}

func (u *userService) Logout(ctx context.Context, userID uint) error {
	// TODO blacklist the token via redis

	// TODO revoke refresh token

	return nil
}

func (u *userService) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Err(domain.ErrUserNotFound).Msg("user not found")
			return nil, fmt.Errorf("user not found: %w", domain.ErrUserNotFound)
		}
		log.Err(err).Msg("error retreiving user by email")
		return nil, err
	}
	return user, nil
}

func (u *userService) ByEmailWithPassword(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.ByEmailWithPassword(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Err(domain.ErrUserNotFound).Msg("user not found")
			return nil, fmt.Errorf("user not found: %w", domain.ErrUserNotFound)
		}
		log.Err(err).Msg("error retreiving user by email")
		return nil, err
	}
	return user, nil
}
