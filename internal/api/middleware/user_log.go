package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/rs/zerolog/log"
)

type ContextID string

const (
	UserIDContextKey ContextID = "userID"
)

// Middleware to extract user ID and attach it to the context and logs.
func UserIDLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the user from the JWT
		claims, ok := c.Get("claims").(*domain.JWTCustomClaims)
		if !ok {
			log.Warn().Msg("no user_id available")
			return next(c)
		}

		// Add userID to context
		ctx := context.WithValue(c.Request().Context(), UserIDContextKey, claims.UserID)
		c.SetRequest(c.Request().WithContext(ctx))

		// Update zerolog global logger with userID context
		log.Logger = log.With().Uint("user_id", claims.UserID).Logger()

		// Continue to next handler
		return next(c)
	}
}
