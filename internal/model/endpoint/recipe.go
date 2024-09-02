package endpoint

import "github.com/meowmix1337/the_recipe_book/internal/model/domain"

type Recipe struct {
	ID    uint   `json:"id"`
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

func NewRecipe(recipe *domain.Recipe) *Recipe {
	return &Recipe{
		ID:    recipe.ID,
		UUID:  recipe.UUID,
		Title: recipe.Title,
	}
}
