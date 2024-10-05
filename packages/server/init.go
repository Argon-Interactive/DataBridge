package server

import (
	"DataBridge/packages/handlers"
	"DataBridge/packages/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func initServer() {
	echoCtx.HideBanner = true
	echoCtx.HidePort = true
	echoCtx.POST("/register", hlrs.RegisterHandler)
	echoCtx.POST("/login", hlrs.LoginHandler)
	echoCtx.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello from Web") })
	echoCtx.GET("/private", func(c echo.Context) error { 
		token := c.QueryParam("token")
		if tokenState := mdlwr.ValidateJWT(token); tokenState != mdlwr.Valid {
			return c.String(http.StatusOK, "Acces Denied") 
		}
		return c.String(http.StatusOK, "Hello from secured connection") 
	})
}
