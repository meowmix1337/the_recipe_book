package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
)

func VerifyJWT(tokenString string, secretKey string) (*domain.JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, echo.ErrUnauthorized
	}

	// Check if the token has expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// check if the token is blacklisted

	return claims, nil
}

// JWTMiddleware verifies the JWT token on each request.
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.ErrUnauthorized
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := VerifyJWT(tokenString, secretKey)
			if err != nil {
				return err
			}

			c.Set("claims", claims)
			return next(c)
		}
	}
}
