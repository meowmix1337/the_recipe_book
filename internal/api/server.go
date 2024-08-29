package api

import (
	"fmt"

	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/rs/zerolog/log"
)

func Start(cfg *config.Config) {
	log.Info().Msg(fmt.Sprintf("Starting server on port: %v and environment: %v", cfg.Port, cfg.Environment))
	echoRouter := newRouter()
	echoRouter.Logger.Fatal(echoRouter.Start(fmt.Sprintf(":%v", cfg.Port)))
}
