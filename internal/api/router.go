package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func newRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("uri", v.URI).
				Str("method", c.Request().Method).
				Str("remote_ip", v.RemoteIP).
				Str("real_ip", c.RealIP()).
				Str("request_id", v.RequestID).
				Time("start_time", v.StartTime).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())

	return e
}
