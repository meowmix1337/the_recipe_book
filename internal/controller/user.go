package controller

import (
	"errors"
	"net/http"

	"github.com/meowmix1337/the_recipe_book/internal/controller/validation"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/meowmix1337/the_recipe_book/internal/service"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) AddUnprotectedRoutes(e *echo.Echo) {
	e.POST("/signup", uc.signup)
	e.POST("/login", uc.login)
}

func (uc *UserController) signup(c echo.Context) error {
	var req endpoint.UserSignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	validationErrors := make(map[string]interface{})
	if err := c.Validate(&req); err != nil {
		validationErrors = validation.FormatValidationError(err)
	}

	passwordErrors := validation.ValidatePassword(req.Password)
	if len(passwordErrors) > 0 {
		validationErrors["password"] = passwordErrors
	}

	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, &endpoint.UserSignupError{
			Message: "Validation errors",
			Errors:  validationErrors,
		})
	}

	err := uc.UserService.SignUp(req.ToDomain())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User created successfully"})
}

func (uc *UserController) login(c echo.Context) error {
	var req endpoint.UserCredentialsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	if err := c.Validate(&req); err != nil {
		validationErrors := validation.FormatValidationError(err)
		return c.JSON(http.StatusBadRequest, &endpoint.UserSignupError{
			Message: "Validation errors",
			Errors:  validationErrors,
		})
	}

	token, err := uc.UserService.Login(req.ToDomain())
	if err != nil {
		// we want to mask the actual error to the user
		if uc.isUnauthorizedErr(err) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal Server Error"})
	}

	// return JWT token to be stored in client's local storage
	return c.JSON(http.StatusOK, &endpoint.JWTResponse{
		Token: token,
	})
}

func (uc *UserController) isUnauthorizedErr(err error) bool {
	return errors.Is(err, domain.ErrInvalidCredentials) ||
		errors.Is(err, domain.ErrNoCredentialsProvided) ||
		errors.Is(err, domain.ErrUserNotFound)
}
