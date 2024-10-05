package mdlwr

import (
	"DataBridge/packages/config"
	"fmt"
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

func ValidateJWT(tokenStr string) TokenState {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 { return nil, &emptyError{} }
		return []byte(cfg.GetConfig().JWTKey), nil
	})
	if err != nil { return Invalid }
	if !token.Valid { return Invalid }

	var claims jwt.MapClaims
	var exp float64
	var ok bool

	if claims, ok = token.Claims.(jwt.MapClaims); !ok { return Invalid }
	if exp, ok = claims["exp"].(float64); !ok { return Invalid }
	if time.Now().After(time.Unix(int64(exp), 0)) { return Expired }

	fmt.Println("token is not expired")
	return Valid
} 

func GenerateJWT(name string) (string, error) {
	claim := jwt.MapClaims{}
	claim["username"] = name
	claim["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(cfg.GetConfig().JWTKey))
}
