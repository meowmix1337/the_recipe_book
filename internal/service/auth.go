package service

import (
	"context"
	"fmt"
	"time"

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
	BlacklistToken(ctx context.Context, token string, userID uint, expiresAt time.Time) error

	ByRefreshToken(ctx context.Context, userID uint, refreshToken string) (*domain.RefreshToken, error)
}

type authService struct {
	*BaseService

	refreshTokenRepo repo.RefreshTokenRepo
}

func NewAuthService(base *BaseService, refreshTokenRepo repo.RefreshTokenRepo) *authService {
	return &authService{
		BaseService:      base,
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

func (s *authService) ByRefreshToken(ctx context.Context, userID uint, refreshToken string) (*domain.RefreshToken, error) {
	return s.refreshTokenRepo.ByRefreshToken(ctx, userID, refreshToken)
}

func (s *authService) BlacklistToken(ctx context.Context, token string, userID uint, expiresAt time.Time) error {
	key := fmt.Sprintf("%v_%v", userID, token)
	expirationTime := time.Unix(expiresAt.Unix(), 0)
	ttl := time.Until(expirationTime)

	err := s.Cache.Set(ctx, key, "", int(ttl))
	if err != nil {
		log.Err(err).Msg("error blacklisting token")
		return err
	}

	return nil
}
