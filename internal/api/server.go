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

	"github.com/meowmix1337/the_recipe_book/internal/api/middleware"
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/controller"
	"github.com/meowmix1337/the_recipe_book/internal/repo"
	"github.com/meowmix1337/the_recipe_book/internal/service"

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
		api := s.setUpAPI(echoRouter)

		// Initialize repositories
		userRepo := repo.NewUserRepository()

		// Initialize services
		userService := service.NewUserService(s.Config, userRepo)
		recipeService := service.NewRecipeService(s.Config)

		// Initialize controllers
		userController := controller.NewUserController(userService)
		userController.AddUnprotectedRoutes(echoRouter)

		recipeController := controller.NewRecipeController(recipeService)
		recipeController.AddRoutes(api)

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
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware(s.GetJWTSecret()))

	return api
}
