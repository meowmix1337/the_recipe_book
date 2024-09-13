package controller

import (
	"errors"
	"net/http"

	"github.com/meowmix1337/the_recipe_book/internal/api/middleware"
	"github.com/meowmix1337/the_recipe_book/internal/controller/validation"
	"github.com/meowmix1337/the_recipe_book/internal/model/domain"
	"github.com/meowmix1337/the_recipe_book/internal/model/endpoint"
	"github.com/meowmix1337/the_recipe_book/internal/service"
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	*BaseController
	UserService service.UserService
}

func NewUserController(base *BaseController, userService service.UserService) *UserController {
	return &UserController{
		BaseController: base,
		UserService:    userService,
	}
}

func (uc *UserController) AddUnprotectedRoutes(e *echo.Echo) {
	e.POST("/signup", uc.signup)
	e.POST("/login", uc.login)

	// logout needs the middleware since we need to retrieve the JWT claims.
	e.POST("/logout", uc.logout, middleware.JWTMiddleware(uc.Config.GetJWTSecret(), uc.Cache))
	e.POST("/refresh", uc.refresh, middleware.JWTMiddleware(uc.Config.GetJWTSecret(), uc.Cache))

	// TODO: add refresh token route
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

	err := uc.UserService.SignUp(c.Request().Context(), req.ToDomain())
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal Server Error"})
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

	token, err := uc.UserService.Login(c.Request().Context(), req.ToDomain())
	if err != nil {
		// we want to mask the actual error to the user
		if uc.isUnauthorizedErr(err) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal Server Error"})
	}

	// TODO: add refresh token to the response

	// return JWT token to be stored in client's local storage
	return c.JSON(http.StatusOK, token)
}

func (uc *UserController) logout(c echo.Context) error {
	claims, ok := c.Get("claims").(*domain.JWTCustomClaims)
	if !ok {
		log.Error().Msg("Failed to assert claims")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": domain.ErrUnableToVerifyClaim.Error()})
	}

	token, ok := c.Get("jwt_token").(string)
	if !ok {
		log.Error().Msg("Failed to assert token")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": domain.ErrUnableToRetrieveToken.Error()})
	}

	err := uc.UserService.Logout(c.Request().Context(), token, claims)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully logged out",
	})
}

func (uc *UserController) isUnauthorizedErr(err error) bool {
	return errors.Is(err, domain.ErrInvalidCredentials) ||
		errors.Is(err, domain.ErrNoCredentialsProvided) ||
		errors.Is(err, domain.ErrUserNotFound)
}
