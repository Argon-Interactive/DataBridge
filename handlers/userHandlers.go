package hlrs

import (
	"DataBridge/database"
	mdlwr "DataBridge/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	res := dbmgr.VerifyLogin(username, password)
	if !res { return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"}) }

	signedToken, err := mdlwr.GenerateJWT(username)
	if err != nil { return c.JSON(http.StatusInternalServerError, map[string]string{"message":"Failed to create a token for JWT"}) }
	return c.JSON(http.StatusOK, map[string]string{"token":signedToken})
}

func RegisterHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	err := dbmgr.CreateUser(username, password)
	if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"message":err.Error()}) }

	return c.JSON(http.StatusOK, map[string]string{"message":"User created successfully"})
}
