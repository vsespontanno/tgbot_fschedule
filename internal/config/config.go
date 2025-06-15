package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken      string
	FootballDataAPIKey string
	MongoURI           string
	PostgresUser       string
	PostgresPass       string
	PostgresDB         string
	PostgresHost       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		TelegramToken:      os.Getenv("TELEGRAM_BOT_API_KEY"),
		FootballDataAPIKey: os.Getenv("FOOTBALL_DATA_API_KEY"),
		MongoURI:           os.Getenv("MONGODB_URI"),
		PostgresUser:       os.Getenv("POSTGRES_USER"),
		PostgresPass:       os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:         os.Getenv("POSTGRES_DB"),
		PostgresHost:       os.Getenv("POSTGRES_HOST"),
	}
}
