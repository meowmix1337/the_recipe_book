package controller

import "github.com/meowmix1337/the_recipe_book/internal/config"

const (
	V1 = "v1"
)

type BaseController struct {
	Config config.Config
}
