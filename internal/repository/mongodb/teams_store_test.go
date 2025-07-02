package mongodb

import (
	"context"
	"football_tgbot/internal/db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetTeamShortName(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	store := NewMongoDBTeamsStore(client, "football")

	want := "Sevilla"
	got, err := store.GetTeamsShortName(context.Background(), "LaLiga", "Sevilla FC")
	if err != nil {
		t.Fatalf("GetTeamShortName returned an error: %v", err)
	}
	if got != want {
		t.Errorf("GetTeamShortName returned %q, want %q", got, want)
	}

	want1 := "PSG"
	got1, err := store.GetTeamsShortName(context.Background(), "Ligue1", "Paris Saint-Germain FC")
	if err != nil {
		t.Fatalf("GetTeamShortName returned an error: %v", err)
	}
	if got1 != want1 {
		t.Errorf("GetTeamShortName returned %q, want %q", got1, want1)
	}

	want2 := "Manchester City"
	got2, err := store.GetTeamsShortName(context.Background(), "PremierLeague", "Manchester City FC")
	if err != nil {
		t.Fatalf("GetTeamShortName returned an error: %v", err)
	}
	if got2 != want2 {
		t.Errorf("GetTeamShortName returned %q, want %q", got1, want1)
	}
}
