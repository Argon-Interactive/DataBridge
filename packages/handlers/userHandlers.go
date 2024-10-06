package hlrs

import (
	"DataBridge/packages/database"
	"DataBridge/packages/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	authToken, refreshToken, err := dbmgr.Login(username, password)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"authtoken":authToken, "refreshtoken":refreshToken})
}

func RegisterHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	err := dbmgr.CreateUser(username, password)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"message":"User created successfully"})
}

func Refresh(c echo.Context) error {
	refToken := c.FormValue("refreshtoken")
	authToken, refreshToken, err := mdlwr.RefreshJWT(refToken)
	if err != nil { return err.JSON(c) }
	return c.JSON(http.StatusOK, map[string]string{"authtoken":authToken, "refreshtoken":refreshToken})
}
