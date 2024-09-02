package service

import (
	"github.com/meowmix1337/the_recipe_book/internal/config"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
)

type RecipeService interface {
	All() ([]*endpoint.Recipe, error)
}

type recipeService struct {
	*BaseService

	recipes map[uint]*domain.Recipe
}

func NewRecipeService(cfg config.Config) *recipeService {
	return &recipeService{
		BaseService: &BaseService{Config: cfg},
		recipes:     make(map[uint]*domain.Recipe),
	}
}

// check UserService interface implementation on compile time.
var _ RecipeService = (*recipeService)(nil)

func (u *recipeService) All() ([]*endpoint.Recipe, error) {

	return []*endpoint.Recipe{}, nil
}
