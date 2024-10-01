package hlrs

import (
	cfg "DataBridge/config"
	"DataBridge/database"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")

	res := dbmgr.VerifyLogin(username, password)
	if !res { return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"}) }

	claim := jwt.MapClaims{}
	claim["username"] = username
	claim["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(cfg.GetConfig().JWTKey))

	if err != nil { return c.JSON(http.StatusInternalServerError, map[string]string{"message":"Failed to create a token for JWT"}) }
	return c.JSON(http.StatusOK, map[string]string{"token":signedToken})
}

func RegisterHandler(c echo.Context) error {
	username := c.QueryParam("name")
	password := c.QueryParam("password")

	err := dbmgr.CreateUser(username, password)
	if err != nil { return c.JSON(http.StatusBadRequest, map[string]string{"message":err.Error()}) }

	return c.JSON(http.StatusOK, map[string]string{"message":"User created successfully"})
}
