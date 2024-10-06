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
	echoCtx.POST("/register", hlrs.Register)
	echoCtx.POST("/login", hlrs.Login)
	echoCtx.POST("/logout", hlrs.Logout)
	echoCtx.POST("/refresh", hlrs.Refresh)
	echoCtx.POST("/changeusername", hlrs.ChangeName)
	echoCtx.POST("/changepassword", hlrs.ChangePassword)


	echoCtx.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello from Web") })
	echoCtx.GET("/private", func(c echo.Context) error { 
		token := c.QueryParam("token")
		if name, err := mdlwr.ValidateAuthJWT(token); err != nil {
			return c.String(http.StatusOK, "Acces Denied") 
		} else {
			return c.String(http.StatusOK, "Hello from secured connection, " + name) 
		}
	})
}
