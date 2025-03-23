// /home/matthew/tgbot_fschedule/scripts/seed_teams.go (или новый файл get_standings.go)

package main

import (
	"encoding/json"
	"fmt"
	"football_tgbot/handlers"
	"football_tgbot/types"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getLeagueStandings(apiKey, leagueCode string) ([]types.Standing, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings", leagueCode)
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

	var standingsResponse types.StandingsResponse
	err = json.Unmarshal(body, &standingsResponse)
	if err != nil {
		return nil, err
	}

	// Возвращаем таблицу из первого элемента массива Standings
	if len(standingsResponse.Standings) > 0 {
		return standingsResponse.Standings[0].Table, nil
	}
	return nil, fmt.Errorf("no standings found for league code: %s", leagueCode)
}

func main() {
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
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}

	// Подключение к MongoDB
	client, err := connectToMongoDB(mongoURI)

	// Для каждой лиги получаем команды и сохраняем в MongoDB
	for leagueName, leagueCode := range leagues {
		// ... (существующий код для получения и сохранения команд) ...

		// Получаем таблицу лиги
		standings, err := getLeagueStandings(apiKey, leagueCode)
		if err != nil {
			log.Printf("Error getting standings for %s: %v\n", leagueName, err)
			continue
		}
		if len(standings) == 0 {
			log.Printf("No standings found for %s\n", leagueName)
			continue
		}

		// Сохраняем таблицу лиги в отдельную коллекцию (опционально)
		standingsCollection := client.Database("football").Collection(leagueName + "_standings")
		err = handlers.SaveStandingsToMongoDB(standingsCollection, standings)
		if err != nil {
			log.Printf("Error saving standings for %s: %v\n", leagueName, err)
			continue
		}

		fmt.Printf("Successfully saved standings for %s\n", leagueName)
	}
}
