package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ping(e echo.Context) error {
	return e.JSON(http.StatusOK, echo.Map{
		"message": "pong",
	})
}
