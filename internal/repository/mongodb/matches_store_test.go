package mongodb

import (
	"context"
	"fmt"
	"testing"

	"github.com/vsespontanno/tgbot_fschedule/internal/db"
)

func TestGetMatchesInPeriod(t *testing.T) {
	from := "2025-05-26"
	to := "2025-05-31"
	league1 := "UCL"

	client, err := db.ConnectToMongoDB("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	store := NewMongoDBMatchesStore(client, "football", "matches")

	matches, err := store.GetMatchesInPeriod(context.Background(), league1, from, to)
	if err != nil {
		t.Fatalf("Failed to get matches in period: %v", err)
	}
	fmt.Println(matches)

}
