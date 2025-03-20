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

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func saveTeamsToMongoDB(collection *mongo.Collection, teams []types.Team) error {
	var documents []interface{}
	for _, team := range teams {
		documents = append(documents, team)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func getTeamsFromAPI(apiKey, leagueCode string) ([]types.Team, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/teams", leagueCode)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var teamsResponse types.TeamsResponse
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return nil, err
	}

	return teamsResponse.Teams, nil
}

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

	mongoURI := os.Getenv("MONGODB_URI")

	// Подключение к MongoDB
	client, err := connectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO())

	// Лиги и их коды
	leagues := map[string]string{
		"Ligue1":          "FL1",
		"LaLiga":          "PD",
		"PremierLeague":   "PL",
		"Bundesliga":      "BL1",
		"SerieA":          "SA",
		"ChampionsLeague": "CL",
	} // я люблю Машу !!!!!!!!!!!!! (с) Виталик (хозяин тимура)

	// Для каждой лиги получаем команды и сохраняем в MongoDB
	for leagueName, leagueCode := range leagues {
		teams, err := getTeamsFromAPI(apiKey, leagueCode)
		if err != nil {
			log.Printf("Error getting teams for %s: %v\n", leagueName, err)
			continue
		}
		if len(teams) == 0 {
			log.Printf("No teams found for %s\n", leagueName)
			continue
		}
		collection := client.Database("football").Collection(leagueName)
		err = saveTeamsToMongoDB(collection, teams)
		if err != nil {
			log.Printf("Error saving teams for %s: %v\n", leagueName, err)
			continue
		}

		fmt.Printf("Successfully saved %d teams for %s\n", len(teams), leagueName)

	}
}
