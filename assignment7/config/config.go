package config

import (
	"log"
	"github.com/joho/godotenv"
)


func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File.env not found, system environment variables are used")
	}
}