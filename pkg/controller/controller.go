package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	userControllerPkg "rewrite/internal/user/controller"
	userRepositoryPkg "rewrite/internal/user/repository"
	userServicePkg "rewrite/internal/user/service"
)

func InitControllers(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Recover())

	e.GET("/ping", Ping)

	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository)
	userController := userControllerPkg.NewUserController(userService)
	userController.InitRoutes(e)
}
