package controller

import (
	"errors"
	"net/http"
	"rewrite/internal/user/dto"
	"rewrite/internal/user/service"
	"rewrite/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ErrBadRequestBody     = errors.New("bad request body")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNoUserFound        = errors.New("no user found")
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (u *UserController) InitRoutes(e *echo.Echo) {
	// Routes with authentication
	secure := e.Group("")
	secure.Use(middleware.JWT([]byte(config.JWT_SECRET)))

	secure.GET("/users", u.GetAllUser)

	// Public routes
	e.POST("/users", u.CreateUser)
	e.POST("/login", u.Login)
}

func (u *UserController) GetAllUser(c echo.Context) error {
	users, err := u.userService.FindAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(users) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, ErrNoUserFound.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success getting users",
		"data":    users,
	})
}

func (u *UserController) CreateUser(c echo.Context) error {
	var user dto.UserRequest
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrBadRequestBody.Error())
	}

	err = u.userService.CreateUser(user, c.Request().Context())
	if err != nil {
		if err == service.ErrUserExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Success creating user",
	})
}

func (u *UserController) Login(c echo.Context) error {
	var user dto.UserRequest
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrBadRequestBody.Error())
	}

	token, err := u.userService.Login(user, c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidCredentials.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login success",
		"token":   token,
	})
}
