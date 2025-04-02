package bot

import (
	"context"
	"fmt"
	"football_tgbot/bot/handlers"
	"football_tgbot/db"
	"football_tgbot/rating"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func Start() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_API_KEY")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_API_KEY is not set")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI is not set")
	}

	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		return fmt.Errorf("POSTGRES_URI is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// MongoDB connection
	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer mongoClient.Disconnect(context.TODO())

	mongoStore := db.NewMongoDBMatchesStore(mongoClient, "football")

	// Создаем сервис рейтингов
	ratingService := rating.NewService(mongoStore)

	// Запускаем обновление рейтингов в отдельной горутине
	ctx := context.Background()
	go ratingService.StartRatingUpdater(ctx, "football", 1*time.Hour)

	return handleUpdates(bot, mongoStore, ratingService)
}

func handleUpdates(bot *tgbotapi.BotAPI, mongoStore db.MatchesStore, ratingService *rating.Service) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if err := handlers.HandleMessage(bot, update.Message, mongoStore, ratingService); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}

		if update.CallbackQuery != nil {
			if err := handlers.HandleCallbackQuery(bot, update.CallbackQuery, mongoStore, ratingService); err != nil {
				log.Printf("Error handling callback query: %v", err)
			}
		}
	}
	return nil
}
