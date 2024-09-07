package controller

import (
	"github.com/meowmix1337/go-core/cache"
	"github.com/meowmix1337/the_recipe_book/internal/config"
)

const (
	V1 = "v1"
)

type BaseController struct {
	Config config.Config
	Cache  cache.Cache
}

func NewBaseController(cfg config.Config, cache cache.Cache) *BaseController {
	return &BaseController{
		Config: cfg,
		Cache:  cache,
	}
}
