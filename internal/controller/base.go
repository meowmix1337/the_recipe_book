package controller

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/go-core/cache"
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/rs/zerolog/log"
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

func (bc *BaseController) GetPaginationParams(c echo.Context) (*endpoint.PagniationParams, error) {
	params := new(endpoint.PagniationParams)
	if err := c.Bind(params); err != nil {
		return nil, err
	}

	if params.Order == "" {
		return nil, endpoint.ErrMissingOrder
	}

	if params.Limit == 0 {
		return nil, endpoint.ErrMissingLimit
	}

	if strings.ToLower(params.Order) != "asc" && strings.ToLower(params.Order) != "desc" {
		return nil, endpoint.ErrUnsupportedOrder
	}

	if params.Limit > endpoint.MaxPaginationLimit || params.Limit < endpoint.MinPaginationLimit {
		return nil, endpoint.ErrPaginationLimit
	}

	return params, nil
}

func (bc *BaseController) GetUserID(c echo.Context) (uint, error) {
	claims, ok := c.Get("claims").(*domain.JWTCustomClaims)
	if !ok {
		log.Error().Msg("Failed to assert claims")
		return 0, domain.ErrUnableToVerifyClaim
	}

	return claims.UserID, nil
}
