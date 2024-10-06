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
		password CHAR[60] NOT NULL,
		isLoggedIn BOOL,
		expirationDate TIMESTAMP,
		key CHAR[36]
		)
		`)
	return err
}

func CreateUser(username string, password string) *cerr.ServerError {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Failed to encrypt") }
	_, err = db.Exec("INSERT INTO users (name, password, isLoggedIn) VALUES (?, ?, false)", username, hashedPasswd)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Failed to update database") }
	return nil
}

func Logout(username string) *cerr.ServerError {
	_, err := db.Exec("UPDATE users SET isLoggedIn = false, expirationDate = DEFAULT, key = '' WHERE name = ?", username)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }
	return nil
}

func Login(username string, password string) (string, string, *cerr.ServerError) {
	if !verifyLogin(username, password) { return "", "", cerr.NewServerError(http.StatusUnauthorized, "Invalid credentials") }
	var isLoggedIn bool
	err := db.QueryRow("SELECT isLoggedIn FROM users WHERE name = ?", username).Scan(&isLoggedIn)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to read from database" + err.Error()) }
	if isLoggedIn { return "", "", cerr.NewServerError(http.StatusBadRequest, "User is already logged in") }

	uniqKey := uuid.New().String()
	authToken, err := generateJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to generate token") }

	refreshToken, err := generateRefreshJWT(username, uniqKey)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to generate token") }
	expirationDate := time.Now().Add(1 * time.Hour).Unix()

	_, err = db.Exec("UPDATE users SET isLoggedIn = true, expirationDate = ?, key = ? WHERE name = ?", expirationDate, uniqKey, username)
	if err != nil { return "", "", cerr.NewServerError(http.StatusInternalServerError, "Failed to update database") }


	
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

func ChangeUsername(oldUsername string, newUsername string) *cerr.ServerError {
	_, err := db.Exec("UPDATE users SET name = ? WHERE name = ?", newUsername, oldUsername)
	if err != nil { return cerr.NewServerError(http.StatusInternalServerError, "Internal server error") }
	return nil
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


