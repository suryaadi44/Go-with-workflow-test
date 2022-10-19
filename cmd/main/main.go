package main

import (
	"rewrite/pkg/config"
	"rewrite/pkg/controller"
	"rewrite/pkg/database"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}

	err = database.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	controller.InitControllers(e, db)

	e.Logger.Fatal(e.Start(config.PORT))
}
