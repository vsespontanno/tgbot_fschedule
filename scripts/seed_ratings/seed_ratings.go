package main

import (
	"context"
	"football_tgbot/db"
	"football_tgbot/types"
	"log"
	"time"
)

func main() {
	// Подключаемся к MongoDB
	client, err := db.ConnectToMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	store := db.NewMongoDBMatchesStore(client, "football")

	// Получаем все таблицы
	leagues := []string{"laliga", "epl", "bundesliga", "seriea"}
	teamRatings := make(map[int]types.TeamRating)

	// Собираем информацию о командах из таблиц
	for _, league := range leagues {
		standings, err := store.GetStandings(context.Background(), league)
		if err != nil {
			log.Printf("Failed to get standings for %s: %v", league, err)
			continue
		}

		for _, standing := range standings {
			rating := types.TeamRating{
				TeamID:           standing.Team.ID,
				TeamName:         standing.Team.Name,
				Position:         standing.Position,
				Points:           standing.Points,
				Form:             0.5, // Базовое значение формы
				GoalDiff:         standing.GoalDifference,
				TournamentWeight: getTournamentWeight(league),
				LastUpdated:      time.Now().Format(time.RFC3339),
			}
			teamRatings[standing.Team.ID] = rating
		}
	}

	// Обновляем рейтинги в базе данных
	for _, rating := range teamRatings {
		if err := store.UpdateTeamRating(context.Background(), "team_ratings", rating); err != nil {
			log.Printf("Failed to update rating for team %s: %v", rating.TeamName, err)
			continue
		}
		log.Printf("Updated rating for team: %s", rating.TeamName)
	}

	log.Println("Successfully seeded team ratings")
}

func getTournamentWeight(league string) float64 {
	switch league {
	case "ucl":
		return 1.0
	case "uel":
		return 0.9
	case "epl":
		return 0.8
	case "laliga":
		return 0.8
	case "bundesliga":
		return 0.7
	case "seriea":
		return 0.7
	default:
		return 0.5
	}
}
