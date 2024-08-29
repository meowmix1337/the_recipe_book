package api

import (
	"fmt"

	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/controller"
	"github.com/rs/zerolog/log"
)

func Start(cfg *config.Config) {
	log.Info().Msg(fmt.Sprintf("Starting server on port: %v and environment: %v", cfg.Port, cfg.Environment))
	echoRouter := newRouter()

	// Public routes
	// e.POST("/signup", userController.Signup)
	// e.POST("/login", userController.Login)

	// // Private routes
	// api := e.Group("/api")
	// api.Use(middleware.JWTAuth())
	// api.GET("/protected", userController.Protected)

	// // Recipe routes
	// api.POST("/recipes", recipeController.CreateRecipe)
	// api.GET("/recipes/:id", recipeController.GetRecipe)

	// Initialize repositories
	//  userRepo := repository.NewUserRepository(db)

	// Initialize services
	//  userService := service.NewUserService(userRepo, []byte(cfg.JWTSecret))

	// Initialize controllers
	userController := controller.NewUserController()
	userController.AddRoutes(echoRouter)

	echoRouter.Logger.Fatal(echoRouter.Start(fmt.Sprintf(":%v", cfg.Port)))
}
