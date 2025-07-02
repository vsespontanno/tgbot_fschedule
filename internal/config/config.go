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
	PostgresPort       string
	RedisURL           string
}

func LoadConfig(opt string) *Config {
	if opt == "" {
		opt = ".env"
	}
	err := godotenv.Load(opt)
	if err != nil {
		log.Fatalf("Error loading .env file: %v: %v", err, opt)
	}

	return &Config{
		TelegramToken:      os.Getenv("TELEGRAM_BOT_API_KEY"),
		FootballDataAPIKey: os.Getenv("FOOTBALL_DATA_API_KEY"),
		MongoURI:           os.Getenv("MONGODB_URI"),
		PostgresUser:       os.Getenv("PG_USER"),
		PostgresPass:       os.Getenv("PG_PASSWORD"),
		PostgresDB:         os.Getenv("PG_DB"),
		PostgresHost:       os.Getenv("PG_HOST"),
		PostgresPort:       os.Getenv("PG_PORT"),
		RedisURL:           os.Getenv("REDIS_URL"),
	}
}
