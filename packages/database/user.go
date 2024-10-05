package dbmgr

import (
	"golang.org/x/crypto/bcrypt"
)

func initUser() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password VARCHAR[255] NOT NULL
		)
		`)
	if err != nil { return err }
	return nil
}

func CreateUser(username string, password string) error {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return err }
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", username, hashedPasswd)
	return err
}

func VerifyLogin(username string, password string) bool {
	if len(username) == 0 || len(password) == 0 { return false }
	var hashedPasswd string
	err := db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&hashedPasswd)
	if err != nil { return false }
	err = bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password))
	return err == nil
}

