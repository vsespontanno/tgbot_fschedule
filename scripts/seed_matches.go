package main

import (
	"context"
	"encoding/json"
	"fmt"
	"football_tgbot/types"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получение значений из .env
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}

	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	mongoURI := os.Getenv("MONGODB_URI")

	// Подключение к MongoDB
	client, err := connectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	httpclient := &http.Client{}
	matches, err := getMatchesSchedule(apiKey, today, tomorrow, httpclient)
	if err != nil {
		log.Fatal(err)
	}
	err = saveMatchesToMongoDB(client, matches, today)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully saved %d matches\n", len(matches))

	client.Disconnect(context.TODO())
}

func getMatchesSchedule(apiKey string, today string, tomorrow string, client *http.Client) ([]types.Match, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", today, tomorrow)
	// req, err := http.NewRequest("GET", today)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var MatchesResponse types.MatchesResponse
	err = json.Unmarshal(body, &MatchesResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return MatchesResponse.Matches, nil

}

func saveMatchesToMongoDB(client *mongo.Client, matches []types.Match, today string) error {
	if len(matches) == 0 {
		log.Printf("No matches found for today: %s\n", today)
		return nil
	}

	collection := client.Database("football").Collection("matches")

	var documents []interface{}
	for _, match := range matches {
		documents = append(documents, match)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		return fmt.Errorf("error inserting matches: %v", err)
	}

	return nil
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}
