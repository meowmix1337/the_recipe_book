package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/controller"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/repo"
	"github.com/meowmix1337/the_recipe_book/internal/service"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog/log"
)

const shutdownTime = time.Second * 5

type Server struct {
	config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Start() {
	echoRouter := newRouter()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		// api := s.setUpAPI(echoRouter)
		// api.Use(echojwt.JWT([]byte(cfg.GetJWTSecret())))

		// // Recipe routes
		// api.POST("/recipes", recipeController.CreateRecipe)
		// api.GET("/recipes/:id", recipeController.GetRecipe)

		// Initialize repositories
		userRepo := repo.NewUserRepository()

		// Initialize services
		userService := service.NewUserService(s.Config, userRepo)

		// Initialize controllers
		userController := controller.NewUserController(userService)
		userController.AddUnprotectedRoutes(echoRouter)
		// userController.AddRoutes(api)

		log.Info().
			Msg(fmt.Sprintf("Starting server on port: %v and environment: %v", s.Config.GetPort(), s.Config.GetEnvironment()))
		if err := echoRouter.Start(fmt.Sprintf(":%v", s.Config.GetPort())); err != nil && errors.Is(err, http.ErrServerClosed) {
			// TODO do things before shutdown
			echoRouter.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()
	if err := echoRouter.Shutdown(ctx); err != nil {
		echoRouter.Logger.Fatal(err)
	}
}

func (s *Server) setUpAPI(e *echo.Echo) *echo.Group {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(domain.JWTCustomClaims)
		},
		SigningKey: []byte(s.Config.GetJWTSecret()),
	}

	api := e.Group("/api")
	api.Use(echojwt.WithConfig(config))

	return api
}
