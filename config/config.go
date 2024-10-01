package cfg

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort	string	
	JWTKey		string
}

func GetConfig(filepath string) (Config, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Panicf("ERROR: Couldn't load .env file: %s", err)
	}
	return Config {
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTKey: os.Getenv("JWT_KEY"),
	}, nil
}
