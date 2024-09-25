package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/rs/zerolog/log"
)

type RecipeService interface {
	ParseAllParams(c echo.Context) (*domain.RecipeAllParams, error)
	All(ctx context.Context, userID uint, paginationParams *endpoint.PagniationParams, filterParams *domain.RecipeAllParams) ([]*endpoint.Recipe, error)
}

type recipeService struct {
	*BaseService

	recipes map[uint]*domain.Recipe
}

func NewRecipeService(base *BaseService) *recipeService {
	return &recipeService{
		BaseService: base,
		recipes:     make(map[uint]*domain.Recipe),
	}
}

// check UserService interface implementation on compile time.
var _ RecipeService = (*recipeService)(nil)

func (s *recipeService) ParseAllParams(c echo.Context) (*domain.RecipeAllParams, error) {
	params := new(domain.RecipeAllParams)
	if err := c.Bind(params); err != nil {
		return nil, err
	}

	if len(params.Title) > domain.MaxTitleLength {
		return nil, domain.ErrTitleParamInvalid
	}

	if params.Rating < domain.MinRating || params.Rating > domain.MaxRating {
		return nil, domain.ErrRatingParamInvalid
	}

	for _, tag := range params.Tags {
		if len(tag) > domain.MaxTagLength {
			return nil, domain.ErrTagsParamInvalid
		}
	}

	return params, nil
}

func (s *recipeService) All(ctx context.Context, userID uint, paginationParams *endpoint.PagniationParams, filterParams *domain.RecipeAllParams) ([]*endpoint.Recipe, error) {
	log.Info().
		Str("recipe_title", filterParams.Title).
		Int("rating", filterParams.Rating).
		Interface("tags", filterParams.Tags).
		Int("limit", paginationParams.Limit).
		Str("order", paginationParams.Order).
		Str("cursor", paginationParams.Cursor).
		Msg("Retrieving recipes")

	return []*endpoint.Recipe{}, nil
}
