package dbmgr

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil { log.Panicf("Failed to initialize the database: %s", err.Error()) }
	err = initUser()
	if err != nil { log.Panicf("Failed to initalize Users in the database: %s", err.Error()) }
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	db.Close()
}
