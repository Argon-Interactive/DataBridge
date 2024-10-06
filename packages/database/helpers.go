package dbmgr

import (
	"DataBridge/packages/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)


func generateJWT(name string, key string) (string, error) {
	claim := jwt.MapClaims{}
	claim["name"] = name
	claim["use"] = "Authorization"
	claim["key"] = key
	claim["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(cfg.GetConfig().JWTKey))
}

func generateRefreshJWT(name string, key string) (string, error) {
	claim := jwt.MapClaims{}
	claim["name"] = name
	claim["use"] = "Refresh"
	claim["key"] = key
	claim["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(cfg.GetConfig().JWTKey))
}

func verifyLogin(username string, password string) bool {
	if len(username) == 0 || len(password) == 0 { return false }
	var hashedPasswd string
	err := db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&hashedPasswd)
	if err != nil { return false }
	err = bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password))
	return err == nil
}
