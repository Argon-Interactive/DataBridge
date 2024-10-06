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

func validateJWT(tokenStr string, tokenType string) (TokenState, string) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 { return nil, &emptyError{} }
		return []byte(cfg.GetConfig().JWTKey), nil
	})
	if err != nil { return Invalid, "" }
	if !token.Valid { return Invalid, "" }

	var claims jwt.MapClaims
	var exp float64
	var	key string 
	var	name string 
	var	use string 
	var ok bool

	if claims, ok = token.Claims.(jwt.MapClaims); !ok { return Invalid, "" }
	if exp, ok = claims["exp"].(float64); !ok { return Invalid, "" }
	if time.Now().After(time.Unix(int64(exp), 0)) { return Expired, "" }

	if key, ok = claims["key"].(string); !ok { return Invalid, "" }
	if name, ok = claims["name"].(string); !ok { return Invalid, "" }
	if use, ok = claims["use"].(string); !ok { return Invalid, "" }

	if use != tokenType { return Invalid, "" }

	var keyTemplate string
	err = dbmgr.GetDB().QueryRow("SELECT key FROM users WHERE name = ?", name).Scan(&keyTemplate)
	if err != nil { return Invalid, "" }

	if keyTemplate != key { return Invalid, "" }

	return Valid, name
} 

func ValidateAuthJWT(tokenStr string) (TokenState, string) { return validateJWT(tokenStr, "Authorization") }
func ValidateRefreshJWT(tokenStr string) (TokenState, string) { return validateJWT(tokenStr, "Refresh") }

func RefreshJWT(tokenStr string) (string, string, *cerr.ServerError) {
	state, name := ValidateRefreshJWT(tokenStr)
	if state == Invalid { return "", "", cerr.NewServerError(http.StatusUnauthorized, "Token is invalid") }
	if state == Expired { return "", "", cerr.NewServerError(http.StatusUnauthorized, "Token is expired") }
	return dbmgr.Refresh(name)
	
}
