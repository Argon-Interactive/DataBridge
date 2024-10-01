package main

import (
	"DataBridge/database"
	"DataBridge/config"

	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
)


func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Web")
}

func main() {
	e := echo.New()
	cfg, _ := cfg.GetConfig(".env")

	db, err := dbmgr.InitDB()
	if err != nil {
		fmt.Println("Error initializing database: ", err)
		return
	}
	defer db.Close()

	e.GET("/test", hello)
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
