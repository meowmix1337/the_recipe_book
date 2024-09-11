package service

import (
	"context"
	"time"

	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/repo"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	GenerateToken(ctx context.Context, user *domain.User) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uint) (string, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error
}

type authService struct {
	config.Config

	refreshTokenRepo repo.RefreshTokenRepo
}

func NewAuthService(cfg config.Config, refreshTokenRepo repo.RefreshTokenRepo) *authService {
	return &authService{
		Config:           cfg,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// check UserService interface implementation on compile time.
var _ AuthService = (*authService)(nil)

func (s *authService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	claims := &domain.JWTCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		UUID:   user.UUID,
		Admin:  false,
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "Recipe App",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(domain.JWTExpiration)),
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.Config.GetJWTSecret()))
	if err != nil {
		log.Err(err).Msg("error generating JWT token")
		return "", domain.ErrJWTGeneration
	}

	return tokenString, nil
}

func (s *authService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	uuid := uuid.NewString()

	err := s.refreshTokenRepo.CreateRefreshToken(ctx, uuid, userID)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (s *authService) DeleteRefreshToken(ctx context.Context, userID uint) error {
	return s.refreshTokenRepo.DeleteRefreshToken(ctx, userID)
}
