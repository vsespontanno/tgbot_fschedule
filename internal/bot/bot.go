package bot

import (
	"context"
	"fmt"
	"football_tgbot/internal/bot/handlers"
	"football_tgbot/internal/cache"
	"football_tgbot/internal/config"
	"football_tgbot/internal/db"
	"football_tgbot/internal/infrastructure/api"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	pgRepo "football_tgbot/internal/repository/postgres"
	"football_tgbot/internal/service"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start() error {
	cfg := config.LoadConfig("")
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// MongoDB connection
	mongoClient, err := db.ConnectToMongoDB(cfg.MongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer mongoClient.Disconnect(context.TODO())

	// PostgreSQL connection
	pg, err := db.ConnectToPostgres(cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB, cfg.PostgresHost, cfg.PostgresPort)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	defer pg.Close()

	// Redis connection
	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	defer redisClient.Close()

	// Initialize stores and services
	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")
	standingsStore := mongoRepo.NewMongoDBStandingsStore(mongoClient, "football")
	userStore := pgRepo.NewPGUserStore(pg)

	footballData := api.NewFootballAPIClient(&http.Client{}, cfg.FootballDataAPIKey)

	standingsService := service.NewStandingService(standingsStore)
	matchesService := service.NewMatchesService(matchesStore, footballData)
	userService := service.NewUserService(userStore)

	service.Start(mongoClient, redisClient)

	return handleUpdates(bot, standingsService, matchesService, userService, redisClient)
}

func handleUpdates(bot *tgbotapi.BotAPI, standingsService *service.StandingsService, matchesService *service.MatchesService, userService *service.UserService, redisClient *cache.RedisClient) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if err := handlers.HandleMessage(bot, update.Message, userService); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}

		if update.CallbackQuery != nil {
			if err := handlers.HandleCallbackQuery(bot, update.CallbackQuery, matchesService, standingsService, redisClient); err != nil {
				log.Printf("Error handling callback query: %v", err)
			}
		}
	}
	return nil
}
