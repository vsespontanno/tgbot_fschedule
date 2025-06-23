// В файле internal/service/data_updater_test.go

package service

import (
	"context"
	"football_tgbot/internal/config"
	db "football_tgbot/internal/db"
	"football_tgbot/internal/infrastructure/api"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"net/http"
	"testing"
)

func TestUpdateMatchesDataToTestCollection(t *testing.T) {
	cfg := config.LoadConfig("../../.env")
	mongoURI := cfg.MongoURI
	footballAPI := cfg.FootballDataAPIKey

	httpclient := &http.Client{}
	ctx := context.Background()
	FootballData := api.NewFootballAPIClient(httpclient, footballAPI)

	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())

	// Используем коллекцию matches_test
	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")

	matchesService := NewMatchesService(matchesStore, FootballData)

	from := "2025-05-07"
	to := "2025-05-14"

	matches, err := matchesService.HandleReqMatches(ctx, from, to)
	if err != nil {
		t.Fatalf("Error fetching matches: %v", err)
	}

	// Сохраняем в matches_test
	err = matchesService.HandleSaveMatches(matches, from, to)
	if err != nil {
		t.Fatalf("Error saving matches to matches_test: %v", err)
	}
}
