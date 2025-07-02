package jobs

import (
	"context"
	"net/http"
	"time"

	"football_tgbot/internal/adapters"
	"football_tgbot/internal/cache"
	"football_tgbot/internal/config"
	"football_tgbot/internal/infrastructure/api"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/service"

	"football_tgbot/internal/domain"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateMatchesDataWithCollection(ctx context.Context, mongoClient *mongo.Client, redisClient *cache.RedisClient) {
	cfg := config.LoadConfig("../../.env")
	footballAPI := cfg.FootballDataAPIKey

	logrus.Info("Updating data")

	from := time.Now().Format("2006-01-02")
	to := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	httpclient := &http.Client{}

	// Инфраструктурные клиенты и репозитории
	footballDataClient := api.NewFootballAPIClient(httpclient, footballAPI)
	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(mongoClient, "football")
	standingsStore := mongoRepo.NewMongoDBStandingsStore(mongoClient, "football")

	// Сервисы
	matchesService := service.NewMatchesService(matchesStore, footballDataClient)
	teamsService := service.NewTeamsService(teamsStore)
	standingsService := service.NewStandingService(standingsStore)

	// Адаптер домена, реализующий domain.Calculator через наши сервисы
	calculator := adapters.NewCalculatorAdapter(teamsService, standingsService, matchesService)

	logrus.Infof("Fetching matches from %s to %s…", from, to)
	matches, err := matchesService.HandleReqMatches(ctx, from, to)
	if err != nil {
		logrus.Errorf("Error fetching matches: %v", err)
		return
	}
	if len(matches) == 0 {
		logrus.Info("No matches to upsert")
		return
	}

	for _, match := range matches {
		if err := matchesService.HandleUpsertMatch(ctx, match); err != nil {
			logrus.Errorf("Error upserting match %d: %v", match.ID, err)
			continue
		}

		// Вычисляем рейтинг через домен
		rating, err := domain.CalculateRatingOfMatch(ctx, match, calculator)
		if err != nil {
			logrus.Warnf("Failed to calculate rating for match %d: %v", match.ID, err)
		} else {
			logrus.Infof("Match %d rated: %.2f", match.ID, rating)
		}
	}

	// Инвалидация кэша
	if err := redisClient.DeleteByPattern(context.Background(), "top_matches_image"); err != nil {
		logrus.Errorf("Failed to invalidate top_matches_image cache: %v", err)
	}
	if err := redisClient.DeleteByPattern(context.Background(), "all_matches_image*"); err != nil {
		logrus.Errorf("Failed to invalidate all_matches_image cache: %v", err)
	}
	if err := redisClient.DeleteByPattern(context.Background(), "table_image:*"); err != nil {
		logrus.Errorf("Failed to invalidate table_image cache: %v", err)
	}

	// Обработка паники
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Panic in updater: %v", r)
		}
	}()
}

func Start(mongoClient *mongo.Client, redisClient *cache.RedisClient) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			UpdateMatchesDataWithCollection(ctx, mongoClient, redisClient)
			cancel()
		}
	}()
}
