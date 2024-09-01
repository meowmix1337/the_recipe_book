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
	"github.com/meowmix1337/the_recipe_book/internal/repo"
	"github.com/meowmix1337/the_recipe_book/internal/service"

	"github.com/rs/zerolog/log"
)

const shutdownTime = time.Second * 5

func Start(cfg config.Config) {
	echoRouter := newRouter()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		// Public routes
		// e.POST("/signup", userController.Signup)
		// e.POST("/login", userController.Login)

		// Private routes
		// api := echoRouter.Group("/api")
		// api.Use(middleware.JWTAuth())
		// api.GET("/protected", userController.Protected)

		// // Recipe routes
		// api.POST("/recipes", recipeController.CreateRecipe)
		// api.GET("/recipes/:id", recipeController.GetRecipe)

		// Initialize repositories
		userRepo := repo.NewUserRepository()

		// Initialize services
		userService := service.NewUserService(cfg, userRepo)

		// Initialize controllers
		userController := controller.NewUserController(userService)
		userController.AddUnprotectedRoutes(echoRouter)
		// userController.AddRoutes(api)

		log.Info().Msg(fmt.Sprintf("Starting server on port: %v and environment: %v", cfg.GetPort(), cfg.GetEnvironment()))
		if err := echoRouter.Start(fmt.Sprintf(":%v", cfg.GetPort())); err != nil && errors.Is(err, http.ErrServerClosed) {
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
