package server

import (
	"DataBridge/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func initServer() {
	echoCtx.HideBanner = true
	echoCtx.HidePort = true
	echoCtx.POST("/register", hlrs.RegisterHandler)
	echoCtx.POST("/login", hlrs.LoginHandler)
	echoCtx.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello from Web") })
}
