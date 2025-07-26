package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации приложения
// Содержит ключи API, параметры подключения к базам данных и другие настройки
// Используется для загрузки переменных окружения из .env файла
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

// LoadConfig функция для загрузки конфигурации из .env файла
// Использует пакет godotenv для чтения переменных окружения
// Возвращает указатель на Config с заполненными полями
// Если не удается загрузить .env файл, выводит ошибку в лог
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
