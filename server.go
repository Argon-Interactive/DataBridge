package main

import (
	"DataBridge/database"
	"DataBridge/config"
	"DataBridge/handlers"

	"net/http"
	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Web")
}

func main() {

	e := echo.New()
	cfg.InitConfig(".env")

	dbmgr.InitDB()
	defer dbmgr.CloseDB()

	e.POST("/register", hlrs.RegisterHandler)
	e.POST("/login", hlrs.LoginHandler)
	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":" + cfg.GetConfig().ServerPort))
}
