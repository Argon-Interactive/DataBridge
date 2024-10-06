package hlrs

import (
	"DataBridge/packages/database"
	"DataBridge/packages/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	authToken, refreshToken, err := dbmgr.Login(username, password)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"authtoken":authToken, "refreshtoken":refreshToken})
}

func Logout(c echo.Context) error {
	token := c.FormValue("token")

	name, err := mdlwr.ValidateAuthJWT(token)
	if err != nil { return err.JSON(c) }

	err = dbmgr.Logout(name)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"message":"User logout successfully"})
}

func Register(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	err := dbmgr.CreateUser(username, password)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"message":"User created successfully"})
}

func Refresh(c echo.Context) error {
	token := c.FormValue("token")

	name, err := mdlwr.ValidateRefreshJWT(token)
	if err != nil { return err.JSON(c) }

	authToken, refreshToken, err := dbmgr.Refresh(name)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"authtoken":authToken, "refreshtoken":refreshToken})
}

func ChangeName(c echo.Context) error {
	token := c.FormValue("token")
	newName := c.FormValue("name")

	name, err := mdlwr.ValidateAuthJWT(token)
	if err != nil { return err.JSON(c) }

	authToken, refreshToken, err := dbmgr.ChangeUsername(name, newName)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"authtoken":authToken, "refreshtoken":refreshToken})
}

func ChangePassword(c echo.Context) error {
	token := c.FormValue("token")
	password := c.FormValue("password")

	name, err := mdlwr.ValidateAuthJWT(token)
	if err != nil { return err.JSON(c) }

	err = dbmgr.ChangePassword(name, password)
	if err != nil { return err.JSON(c) }

	return c.JSON(http.StatusOK, map[string]string{"message":"Password changed successfully"})
}
