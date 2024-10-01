package main

import (
	"fmt"
	"database/sql"
	"net/http"
	"github.com/labstack/echo/v4"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./server.db")
	if err != nil { return nil, err }

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id INTIGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
		)
		`)
	if err != nil { return nil, err }
	return db, nil
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Web")
}

func main() {
	fmt.Println("This compiles at least")
	e := echo.New()

	db, err := initDB()
	if err != nil {
		fmt.Println("Error initializing database: ", err)
		return
	}
	defer db.Close()

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":8000"))
}
