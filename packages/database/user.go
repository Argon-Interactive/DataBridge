package dbmgr

import (
	"DataBridge/packages/error"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)


func initUser() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password CHAR(60) NOT NULL,
		expirationDate INTIGER DEFAULT 0,
		key CHAR(36) DEFAULT ''
		)
		`)
	return err
}

func GetLoggedAmount() int {
	var amount int 
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE expirationDate > ?", time.Now().Unix()).Scan(&amount)
	if err != nil { return 0 }
	return amount
}

func CreateUser(username string, password string) *cerr.ServerError {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Failed to encrypt") }
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", username, hashedPasswd)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Failed to update database" + err.Error()) }
	return nil
}

func Logout(username string) *cerr.ServerError {
	_, err := db.Exec("UPDATE users SET expirationDate = 0, key = '' WHERE name = ?", username)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Internal server error" + err.Error()) }
	return nil
}

func Login(username string, password string) (string, string, *cerr.ServerError) {
	if !verifyLogin(username, password) { return "", "", cerr.NewServerError(http.StatusUnauthorized, "Invalid credentials") }
	var expiration int64
	err := db.QueryRow("SELECT expirationDate FROM users WHERE name = ?", username).Scan(&expiration)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to read from database" + err.Error()) }
	if time.Now().Before(time.Unix(expiration, 0)) { return "", "", cerr.NewServerError(http.StatusBadRequest, "User is already logged in") }

	uniqKey := uuid.New().String()
	authToken, err := generateJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to generate token") }

	refreshToken, err := generateRefreshJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to generate token") }
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	_, err = db.Exec("UPDATE users SET expirationDate = ? , key = ? WHERE name = ?", expirationTime, uniqKey, username)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to update database" + err.Error()) }


	
	return authToken, refreshToken, nil
}

func Refresh(username string) (string, string, *cerr.ServerError) {
	var isLoggedIn bool
	err := db.QueryRow("SELECT isLoggedIn FROM users WHERE name = ?", username).Scan(&isLoggedIn)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }
	if !isLoggedIn { return "", "", cerr.NewServerError(http.StatusBadRequest, "User isn't logged in") }

	uniqKey := uuid.New().String()
	authToken, err := generateJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	refreshToken, err := generateRefreshJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	expirationDate := time.Now().Add(1 * time.Hour).Unix()

	_, err = db.Exec("UPDATE users SET expirationDate = ?, key = ? WHERE name = ?", expirationDate, uniqKey, username)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	return authToken, refreshToken, nil
}

func ChangeUsername(oldUsername string, newUsername string) (string, string, *cerr.ServerError) {
	if !IsUsernameFree(newUsername) { return "", "", cerr.NewServerError(http.StatusBadRequest, "Username is taken") }
	uniqKey := uuid.New().String()
	authToken, err := generateJWT(newUsername, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	refreshToken, err := generateRefreshJWT(newUsername, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	expirationDate := time.Now().Add(1 * time.Hour).Unix()

	_, err = db.Exec("UPDATE users SET expirationDate = ?, key = ?, name = ? WHERE name = ?", expirationDate, uniqKey, newUsername, oldUsername)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }

	return authToken, refreshToken, nil
}

func ChangePassword(username string, password string) *cerr.ServerError {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }
	_, err = db.Exec("UPDATE users SET password = ? WHERE name = ?", hashedPasswd, username)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }
	return nil
}

func IsUsernameValid(username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", username).Scan(&count)
	if err != nil { return false }
	if count != 0 { return true }
	return false
}
func IsUsernameFree(username string) bool { return !IsUsernameValid(username) }


