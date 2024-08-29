package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) AddRoutes(e *echo.Echo) {
	e.POST("/signup", uc.signup)
}

func (uc *UserController) signup(c echo.Context) error {
	var req endpoint.UserSignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	// err := uc.UserService.Signup(req)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	// }

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}
