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

	"github.com/meowmix1337/go-core/cache"
	"github.com/meowmix1337/go-core/db"
	"github.com/meowmix1337/the_recipe_book/internal/api/middleware"
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/controller"
	"github.com/meowmix1337/the_recipe_book/internal/repo"
	"github.com/meowmix1337/the_recipe_book/internal/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		db, err := s.initializeDB()
		if err != nil {
			echoRouter.Logger.Fatal("failed to initilize DB, shutting down: %w", err)
		}

		cache, err := s.initializeRedis()
		if err != nil {
			echoRouter.Logger.Fatal("failed to initilize Redis, shutting down: %w", err)
		}

		api := s.setUpAPI(echoRouter, cache)

		// Initialize repositories
		userRepo := repo.NewUserRepository(db)

		// Initialize services
		baseService := service.NewBaseService(s.Config, cache)
		userService := service.NewUserService(baseService, userRepo)
		recipeService := service.NewRecipeService(baseService)

		// Initialize controllers
		baseController := controller.NewBaseController(s.Config, cache)
		userController := controller.NewUserController(baseController, userService)
		userController.AddUnprotectedRoutes(echoRouter)

		recipeController := controller.NewRecipeController(baseController, recipeService)
		recipeController.AddRoutes(api)

		log.Info().
			Msg(fmt.Sprintf("Starting server on port: %v and environment: %v", s.Config.GetPort(), s.Config.GetEnvironment()))
		if err = echoRouter.Start(fmt.Sprintf(":%v", s.Config.GetPort())); err != nil && errors.Is(err, http.ErrServerClosed) {
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

func (s *Server) setUpAPI(e *echo.Echo, cache cache.Cache) *echo.Group {
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware(s.GetJWTSecret(), cache))

	return api
}

func (s *Server) initializeDB() (db.DB, error) {
	// TODO: add reader too
	dbDSN := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		s.Config.GetDBUser(),
		s.Config.GetDBPassword(),
		s.Config.GetDBHost(),
		s.Config.GetDBPort(),
		s.Config.GetDBName(),
	)
	db := db.NewPostgres(dbDSN, dbDSN)

	if err := s.runMigrations(dbDSN); err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}

	log.Info().Msg("database initilized")

	return db, nil
}

func (s *Server) runMigrations(writerDSN string) error {
	log.Info().Msg("Running migrations")

	// Create a new migrate instance
	m, err := migrate.New(
		fmt.Sprintf("file://%s", s.Config.GetMigrationPath()),
		writerDSN,
	)
	if err != nil {
		// debug mode, try ../migration
		m, err = migrate.New(
			"file://../migrations",
			writerDSN,
		)
		if err != nil {
			return fmt.Errorf("error creating migrate instance: %w", err)
		}
		log.Info().Msg("Running in debug mode, using ../migration path")
	}

	// Run migrations
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("no migrations to run, no changes detected")
			return nil
		}
		return fmt.Errorf("error running migration: %w", err)
	}

	return nil
}

func (s *Server) initializeRedis() (cache.Cache, error) {
	addr := fmt.Sprintf("%v:%v", s.Config.GetRedisHost(), s.Config.GetRedisPort())
	cache, err := cache.NewRedisCache(addr, s.Config.GetRedisPassword(), 0)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
