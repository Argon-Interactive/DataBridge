package dbmgr

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func initUser(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTIGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password VARCHAR[255] NOT NULL
		)
		`)
	if err != nil { return err }
	return nil
}

func CreateUser(db *sql.DB, username string, password string) error {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return err }
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", username, hashedPasswd)
	return err
}

func VerifyLogin(db *sql.DB, username string, password string) bool {
	var hashedPasswd string
	err := db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&hashedPasswd)
	if err != nil { return false }
	err = bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(password))
	return err == nil
}

