package controller

import (
	"net/http"

	"github.com/meowmix1337/the_recipe_book/internal/service"

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
	e.GET("/"+V1+"/recipes", rc.index)
}

func (rc *RecipeController) index(c echo.Context) error {
	userID, err := rc.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	paginationParams, err := rc.GetPaginationParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	filterParams, err := rc.RecipeService.ParseAllParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	recipes, err := rc.RecipeService.All(c.Request().Context(), userID, paginationParams, filterParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// potentially a pagination struct wrapper here to encapsulate the recipes

	return c.JSON(http.StatusOK, recipes)
}
