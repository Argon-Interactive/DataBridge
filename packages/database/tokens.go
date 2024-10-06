package dbmgr

import (

)

func initTokens() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tokensBL (
		key TEXT NOT NULL UNIQUE,
		)
		`)
	return err
}

func BlacklistToken(tokenStr string) {
	
}
