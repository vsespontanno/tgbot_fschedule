package db

import (
	"context"
	"football_tgbot/internal/config"
	"football_tgbot/internal/db"
	"testing"
)

var test1 int = 57
var collection = "PremierLeague"
var ctx context.Context = context.Background()

func TestGetTeamStanding(t *testing.T) {
	cfg := config.LoadConfig()
	mongoURI := cfg.MongoURI

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	store := NewMongoDBStandingsStore(client, "football")

	position, err := store.GetTeamStanding(ctx, collection, test1)
	if err != nil {
		t.Fatalf("Failed to get team standing: %v", err)
	}

	if position != 2 {
		t.Errorf("Expected position 1, got %d", position)
	}

}
