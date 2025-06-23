package service

import (
	"context"
	"football_tgbot/internal/cache"
	"football_tgbot/internal/config"
	"football_tgbot/internal/infrastructure/api"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateMatchesDataWithCollection(ctx context.Context, mongoClient *mongo.Client, redisClient *cache.RedisClient) {
	cfg := config.LoadConfig("../../.env")
	footballAPI := cfg.FootballDataAPIKey

	logrus.Info("Updating data")

	from := time.Now().Format("2006-01-02")
	to := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	defer mongoClient.Disconnect(context.TODO())
	logrus.Info("Connected to MongoDB")

	httpclient := &http.Client{}

	FootballData := api.NewFootballAPIClient(httpclient, footballAPI)
	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(mongoClient, "football")

	standingsStore := mongoRepo.NewMongoDBStandingsStore(mongoClient, "football")
	matchesService := NewMatchesService(matchesStore, FootballData)
	teamsService := NewTeamsService(teamsStore)
	standingsService := NewStandingService(standingsStore)

	logrus.Infof("Fetching matches from %s to %s…", from, to)
	matches, err := matchesService.HandleReqMatches(ctx, from, to)
	if err != nil {
		logrus.Errorf("Error fetching matches: %v", err)
		return
	}

	var newMatches []types.Match
	for _, match := range matches {
		existingMatch, err := matchesService.HandleGetMatchByID(ctx, match.ID)
		if err != nil || existingMatch.ID == 0 {
			newMatches = append(newMatches, match)
		}
	}

	if len(matches) == 0 {
		logrus.Info("No new matches to save")
		return
	}

	err = matchesService.HandleSaveMatches(newMatches, from, to)
	if err != nil {
		logrus.Errorf("Error saving matches: %v", err)
		return
	}
	logrus.Infof("Successfully saved %d matches", len(matches))

	for _, match := range matches {
		_, err := CalculateRatingOfMatch(ctx, match, teamsService, standingsService, matchesStore)
		if err != nil {
			logrus.Warnf("Failed to calculate rating for match %d: %v", match.ID, err)
		}
	}

	// Инвалидация кэша
	redisClient.DeleteByPattern(context.Background(), "schedule_image:*")
	redisClient.DeleteByPattern(context.Background(), "table_image:*")

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
