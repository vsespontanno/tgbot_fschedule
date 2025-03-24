package handlers

import (
	"context"
	"fmt"
	"football_tgbot/db"
	"football_tgbot/types"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func SaveStandingsToMongoDB(collection *mongo.Collection, standings []types.Standing) error {
	var documents []interface{}
	for _, standing := range standings {
		documents = append(documents, standing)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func GetStandingsFromDB(store db.MatchesStore, collectionName string) ([]types.Standing, error) {
	// Получаем все команды из коллекции
	var standings []types.Standing
	var err error
	standings, err = store.GetStandings(context.Background(), collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get standings: %v", err)
	}
	return standings, nil
}

// formatStandings форматирует таблицу лиги для отправки пользователю.
func FormatStandings(standings []types.Standing, collectionName string) string {
	underscore := strings.Index(collectionName, "_")
	if underscore != -1 {
		collectionName = collectionName[:underscore]
	}
	// Проверяем, что таблица не пустаs
	if len(standings) == 0 {
		return fmt.Sprintf("Таблица %s пуста.", collectionName)
	}

	// Сортируем таблицу по позиции
	sort.Slice(standings, func(i, j int) bool {
		return standings[i].Position < standings[j].Position
	})

	response := fmt.Sprintf("Таблица %s:\n", collectionName)
	response += fmt.Sprintf("%-4s %-25s %-5s %-5s %-5s %-5s %-5s %-5s %-5s %-5s\n", "#", "Команда", "И", "В", "Н", "П", "ГЗ", "ГП", "РГ", "О")
	for _, standing := range standings {
		response += fmt.Sprintf("%-4d %-25s %-5d %-5d %-5d %-5d %-5d %-5d %-5d %-5d\n",
			standing.Position, standing.Team.Name, standing.PlayedGames, standing.Won, standing.Draw, standing.Lost, standing.GoalsFor, standing.GoalsAgainst, standing.GoalDifference, standing.Points)
	}
	return response
}

// func GetStandingsFromDB(store db.MatchesStore, collectionName string) ([]types.Standing, uint8, error) {
// 	var flag uint8
// 	// Получаем все команды из коллекции
// 	switch collectionName {
// 	case "PL":
// 		flag = 1
// 	case "BL1":
// 		flag = 2
// 	case "FL1":
// 		flag = 3
// 	case "SA":
// 		flag = 4
// 	case "PD":
// 		flag = 5
// 	}
// 	var standings []types.Standing
// 	var err error
// 	standings, err = store.GetStandings(context.Background(), collectionName)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get standings: %v", err)
// 	}
// 	return standings, flag, nil
// }
