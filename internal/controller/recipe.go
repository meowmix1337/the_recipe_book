package controller

import (
	"net/http"

	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/service"
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type RecipeController struct {
	*BaseController
	RecipeService service.RecipeService
}

func NewRecipeController(base *BaseController, recipeService service.RecipeService) *RecipeController {
	return &RecipeController{
		BaseController: base,
		RecipeService:  recipeService,
	}
}

func (rc *RecipeController) AddRoutes(e *echo.Group) {
	e.GET("/"+V1+"/recipes", rc.all)
}

func (rc *RecipeController) all(c echo.Context) error {
	claims, ok := c.Get("claims").(*domain.JWTCustomClaims)
	if !ok {
		log.Error().Msg("Failed to assert claims")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": domain.ErrUnableToVerifyClaim})
	}

	// TODO hook up pagination

	_, err := rc.RecipeService.All()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": claims,
	})
}
