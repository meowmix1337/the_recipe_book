package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/meowmix1337/the_recipe_book/internal/service"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) AddRoutes(e *echo.Echo) {
	e.POST("/signup", uc.signup)
}

func (uc *UserController) signup(c echo.Context) error {
	var req endpoint.UserSignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	err := uc.UserService.SignUp(req.ToDomain())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
