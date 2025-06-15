package bot

import (
	"context"
	"fmt"
	"football_tgbot/internal/bot/handlers"
	"football_tgbot/internal/db"
	"football_tgbot/internal/service"
	"log"
	"os"

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

	matchesStore := db.NewMongoDBMatchesStore(mongoClient, "football")
	standingsStore := db.NewMongoDBStandingsStore(mongoClient, "football")

	// ratingService := rating.NewService(mongoStore)

	standingsService := service.NewStandingService(standingsStore)
	matchesService := service.NewMatchesService(matchesStore)

	return handleUpdates(bot, standingsService, matchesService)
}

func handleUpdates(bot *tgbotapi.BotAPI, standingsService *service.StandingsService, matchesService *service.MatchesService) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if err := handlers.HandleMessage(bot, update.Message); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}

		if update.CallbackQuery != nil {
			if err := handlers.HandleCallbackQuery(bot, update.CallbackQuery, matchesService, standingsService); err != nil {
				log.Printf("Error handling callback query: %v", err)
			}
		}
	}
	return nil
}
