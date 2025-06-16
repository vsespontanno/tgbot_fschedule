package rating

import (
	"context"
	"fmt"
	"football_tgbot/internal/config"
	db "football_tgbot/internal/db"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/service"
	"testing"
)

func TestCalculatePositionOfTeams(t *testing.T) {
	cfg := config.LoadConfig()
	mongoURI := cfg.MongoURI

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB!")

	matchesStore := mongoRepo.NewMongoDBMatchesStore(client, "football")
	standingsStore := mongoRepo.NewMongoDBStandingsStore(client, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(client, "football")
	matchesService := service.NewMatchesService(matchesStore)
	standingsService := service.NewStandingService(standingsStore)
	teamsService := service.NewTeamsService(teamsStore)

	matches, err := matchesService.HandleGetMatches(context.Background())
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}

	match := matches[1]
	want1 := []int{1, 113}
	want2 := []int{14, 104}
	want := [][]int{want1, want2}
	got, err := CalculatePositionOfTeams(context.Background(), teamsService, standingsService, match, 1)
	if err != nil {
		t.Fatalf("Failed to calculate position of teams: %v", err)
	}

	if len(got) != len(want) {
		t.Errorf("voobshe kapec: got %v want %v", got, want)
	} else if got[0][0] != want[0][0] || got[1][0] != want[1][0] {
		t.Errorf("wrong val: got %v want %v", got, want)

	}
}
