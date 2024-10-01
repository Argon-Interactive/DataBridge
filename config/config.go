package cfg

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type config struct {
	ServerPort	string	
	JWTKey		string
}

var cfg config

func InitConfig(filepath string) {
	err := godotenv.Load(filepath)
	if err != nil { log.Panicf("ERROR: Couldn't load .env file: %s", err) }
	cfg = config {
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTKey: os.Getenv("JWT_KEY"),
	}
}

func GetConfig() config {
	return cfg
}
