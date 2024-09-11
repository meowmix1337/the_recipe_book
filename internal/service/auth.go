package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	GenerateToken(user *domain.User) (string, error)
}

type authService struct {
	config.Config
}

func NewAuthService(cfg config.Config) *authService {
	return &authService{
		Config: cfg,
	}
}

// check UserService interface implementation on compile time.
var _ AuthService = (*authService)(nil)

func (s *authService) GenerateToken(user *domain.User) (string, error) {
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
