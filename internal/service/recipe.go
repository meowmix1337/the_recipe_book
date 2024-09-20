package service

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/rs/zerolog/log"
)

type RecipeService interface {
	ParseAllParams(c echo.Context) (*domain.RecipeAllParams, error)
	All(ctx context.Context, userID uint, filterParams *domain.RecipeAllParams) ([]*endpoint.Recipe, error)
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
	recipeTitle := c.QueryParam("title")
	ratingParam := c.QueryParam("rating")
	tags := c.QueryParams()["tags"]

	if len(recipeTitle) > domain.MaxTitleLength {
		return nil, domain.ErrTitleParamInvalid
	}

	var rating int = -1
	var err error
	if ratingParam != "" {
		rating, err = strconv.Atoi(ratingParam)
		if err != nil || rating < domain.MinRating || rating > domain.MaxRating {
			return nil, domain.ErrRatingParamInvalid
		}
	}

	for _, tag := range tags {
		if len(tag) > domain.MaxTagLength {
			return nil, domain.ErrTagsParamInvalid
		}
	}

	return &domain.RecipeAllParams{
		Title:  recipeTitle,
		Rating: rating,
		Tags:   tags,
	}, nil

}

func (s *recipeService) All(ctx context.Context, userID uint, filterParams *domain.RecipeAllParams) ([]*endpoint.Recipe, error) {
	log.Info().
		Str("recipe_title", filterParams.Title).
		Int("rating", filterParams.Rating).
		Interface("tags", filterParams.Tags).
		Msg("Retrieving recipes")

	return []*endpoint.Recipe{}, nil
}
