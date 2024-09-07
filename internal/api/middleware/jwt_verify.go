package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/go-core/cache"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
)

func VerifyJWT(ctx context.Context, cache cache.Cache, tokenString string, secretKey string) (*domain.JWTCustomClaims, error) {
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
		return nil, echo.ErrUnauthorized
	}

	// check if the token is blacklisted
	if isBlacklisted(ctx, cache, claims.UserID, token.Raw) {
		return nil, echo.ErrUnauthorized
	}

	return claims, nil
}

func isBlacklisted(ctx context.Context, cache cache.Cache, userID uint, token string) bool {
	key := fmt.Sprintf("%v_%v", userID, token)
	_, err := cache.Get(ctx, key)
	return err == nil
}

// JWTMiddleware verifies the JWT token on each request.
func JWTMiddleware(secretKey string, cache cache.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims, err := VerifyJWT(c.Request().Context(), cache, tokenString, secretKey)
			if err != nil {
				if errors.Is(err, echo.ErrUnauthorized) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}

			c.Set("claims", claims)
			c.Set("jwt_token", tokenString)
			return next(c)
		}
	}
}
