package rating

import (
	"context"
	"fmt"
	db "football_tgbot/internal/db"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestCalculatePositionOfTeams(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB!")

	matchesStore := mongoRepo.NewMongoDBMatchesStore(client, "football")
	// standingsStore := mongoRepo.NewMongoDBStandingsStore(client, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(client, "football")
	matchesService := service.NewMatchesService(matchesStore)
	// standingsService := service.NewStandingService(standingsStore)
	teamsService := service.NewTeamsService(teamsStore)

	matches, err := matchesService.HandleGetMatches(context.Background())
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}

	var match types.Match

	for _, show := range matches {
		if show.HomeTeam.Name == "Sevilla FC" {
			if show.AwayTeam.Name == "UD Las Palmas" {
				match = show
			}
		}
	}

	homeLeague, awayLeague, err := getLeaguesForTeams(context.Background(), teamsService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		t.Fatalf("Failed to get leagues for teams: %v", err)
	}
	if homeLeague != "LaLiga" || awayLeague != "LaLiga" {
		t.Errorf("pizdec wanted %s and %s, got %s and %s", "LaLiga", "LaLiga", homeLeague, awayLeague)
	}

}
