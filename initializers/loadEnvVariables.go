package initializers

import "log"
import "github.com/joho/godotenv"

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
