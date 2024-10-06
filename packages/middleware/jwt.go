package mdlwr

import (
	"DataBridge/packages/config"
	"DataBridge/packages/database"
	cerr "DataBridge/packages/error"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type emptyError struct {}
func (e *emptyError) Error() string { return fmt.Sprintf("Loggin empty error") }


type TokenState int
const (
	Invalid TokenState = iota
	Expired
	Valid
)

func validateJWT(tokenStr string, tokenType string) (string, *cerr.ServerError) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 { return nil, &emptyError{} }
		return []byte(cfg.GetConfig().JWTKey), nil
	})
	if err != nil { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if !token.Valid { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }

	var claims jwt.MapClaims
	var exp float64
	var	key string 
	var	name string 
	var	use string 
	var ok bool

	if claims, ok = token.Claims.(jwt.MapClaims); !ok { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if exp, ok = claims["exp"].(float64); !ok { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if time.Now().After(time.Unix(int64(exp), 0)) { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is expired") }

	if key, ok = claims["key"].(string); !ok { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if name, ok = claims["name"].(string); !ok { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if use, ok = claims["use"].(string); !ok { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }

	if use != tokenType { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }

	var keyTemplate string
	err = dbmgr.GetDB().QueryRow("SELECT key FROM users WHERE name = ?", name).Scan(&keyTemplate)
	if err != nil { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }

	if keyTemplate != key || keyTemplate == "" { return "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }

	return name, nil
} 

func ValidateAuthJWT(tokenStr string) (string, *cerr.ServerError) { return validateJWT(tokenStr, "Authorization") }
func ValidateRefreshJWT(tokenStr string) (string, *cerr.ServerError) { return validateJWT(tokenStr, "Refresh") }
