package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/client"
	"github.com/vsespontanno/tgbot_fschedule/internal/config"
	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	"github.com/vsespontanno/tgbot_fschedule/internal/jobs"
	"github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"

	"github.com/go-co-op/gocron"
)

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация подключений
	mongoClient, err := db.ConnectToMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	apiClient := client.NewFootballAPIClient(&http.Client{}, cfg.FootballDataAPIKey)

	// Инициализация сервисов
	matchesStore := mongodb.NewMongoDBMatchesStore(mongoClient, "football")
	standingsStore := mongodb.NewMongoDBStandingsStore(mongoClient, "football")
	teamsStore := mongodb.NewMongoDBTeamsStore(mongoClient, "football")

	matchesService := service.NewMatchesService(matchesStore, apiClient)
	standingsService := service.NewStandingService(standingsStore)
	teamsService := service.NewTeamsService(teamsStore)
	calculator := service.NewCalculatorAdapter(teamsStore, standingsStore, matchesStore)

	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Создаем планировщик
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SetMaxConcurrentJobs(1, gocron.RescheduleMode)

	// Регистрируем задачи
	jobs.RegisterStandingsJob(scheduler, standingsService, redisClient, apiClient)
	jobs.RegisterTeamsJob(scheduler, teamsService, apiClient)
	jobs.RegisterMatchesJob(scheduler, matchesService, redisClient, apiClient, calculator)

	// Запускаем планировщик
	scheduler.StartAsync()

	log.Println("Data updater service started")

	// Ожидание сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down data updater...")
	scheduler.Stop()
	cancel()
	log.Println("Service gracefully stopped")
}
