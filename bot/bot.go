package bot

import (
	"context"
	"fmt"
	"football_tgbot/bot/handlers"
	"football_tgbot/db"
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

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(context.TODO())

	store := db.NewMongoDBMatchesStore(client, "football")
	return handleUpdates(bot, store)
}

func handleUpdates(bot *tgbotapi.BotAPI, store db.MatchesStore) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if err := handlers.HandleMessage(bot, update.Message, store); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}

		if update.CallbackQuery != nil {
			if err := handlers.HandleCallbackQuery(bot, update.CallbackQuery, store); err != nil {
				log.Printf("Error handling callback query: %v", err)
			}
		}
	}
	return nil
}
