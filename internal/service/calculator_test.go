package service

import (
	"context"
	"fmt"
	db "football_tgbot/internal/db"
	mongoRepo "football_tgbot/internal/repository/mongodb"
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
	matchesService := NewMatchesService(matchesStore, nil)
	// standingsService := service.NewStandingService(standingsStore)
	teamsService := NewTeamsService(teamsStore)

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
		t.Errorf("wanted %s and %s, got %s and %s", "LaLiga", "LaLiga", homeLeague, awayLeague)
	}

}

func TestMatchRatings(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")

	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB!")

	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")
	standingsStore := mongoRepo.NewMongoDBStandingsStore(mongoClient, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(mongoClient, "football")

	standingsService := NewStandingService(standingsStore)
	matchesService := NewMatchesService(matchesStore, nil)
	teamsService := NewTeamsService(teamsStore)

	ctx := context.Background()

	// Получаем все матчи
	homeTeamName := "Athletic Club"
	awayTeamName := "Deportivo Alavés"
	match1, err := FindMatchByTeamNames(t, matchesService, homeTeamName, awayTeamName)
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}
	homeTeamName1 := "Bayer 04 Leverkusen"
	awayTeamName1 := "Borussia Dortmund"
	match2, err := FindMatchByTeamNames(t, matchesService, homeTeamName1, awayTeamName1)
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}
	//Данные для первого матча
	// 1) Сила команд по позициям
	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, teamsService, standingsService, match1)
	if err != nil {
		t.Errorf("error calculating team strengths: %v", err)
	}
	fmt.Printf("Сила команд: Athletic Club - %f; Deportivo Alavés - %f\n", homeStrength, awayStrength)
	// 2) Лиги и вес
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamsService, match1.HomeTeam.ID, match1.AwayTeam.ID)
	if err != nil || homeLeague == "" || awayLeague == "" {
		t.Errorf("Матч %s - %s пропущен: проблема с лигами\n", match1.HomeTeam.Name, match1.AwayTeam.Name)
	}
	avgLeagueWeight := (leagueNorm[homeLeague] + leagueNorm[awayLeague]) / 2.0
	fmt.Printf("Сила ЛаЛига: %F\n", avgLeagueWeight)
	// 3) Форма команд
	recentMatchesHome, err := matchesStore.GetRecentMatches(ctx, match1.HomeTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for home team %d: %v", match1.HomeTeam.ID, err)
	}
	recentMatchesAway, err := matchesStore.GetRecentMatches(ctx, match1.AwayTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for away team %d: %v", match1.AwayTeam.ID, err)
	}
	homeForm := calculateForm(recentMatchesHome)
	awayForm := calculateForm(recentMatchesAway)
	formFactor := (homeForm + awayForm) / 2.0
	fmt.Printf("Формы каждой команды: Athletic Club - %f; Deportivo Alavés - %f\nОбщая форма - %f\n", homeForm, awayForm, formFactor)
	// 4) Бонусы
	derbyBonus := GetDerbyBonus(ctx, teamsService, match1)
	stageBonus := 0.0
	if homeLeague == "Champions League" && match1.Stage != "" {
		stageBonus = CLstage[match1.Stage]
	}
	crossLeagueBonus := 0.0
	if homeLeague != awayLeague {
		crossLeagueBonus = 0.15 // Увеличен с 0.1 для большего эффекта
	}
	fmt.Printf("Бонусы: %f, %f, %f\n", derbyBonus, stageBonus, crossLeagueBonus)
	// 5) Финальный рейтинг
	baseRating := (homeStrength+awayStrength)/2.0*0.15 + // Уменьшено влияние позиций
		avgLeagueWeight*0.35 + // Снижено с 0.4
		formFactor*0.15 // Добавлено влияние формы
	rating := baseRating * (1 + derbyBonus + stageBonus + crossLeagueBonus)
	fmt.Printf("Финальный рейтинг: %f\n", rating)
	fmt.Println("-----------------------------------")
	//Данные для второго матча
	// 1) Сила команд по позициям
	homeStrength1, awayStrength1, err := CalculatePositionOfTeams(ctx, teamsService, standingsService, match2)
	if err != nil {
		t.Errorf("error calculating team strengths: %v", err)
	}
	fmt.Printf("Сила команд: Bayer 04 Leverkusen - %f; Borussia Dortmund - %f\n", homeStrength1, awayStrength1)
	// 2) Лиги и вес
	homeLeague1, awayLeague1, err := getLeaguesForTeams(ctx, teamsService, match2.HomeTeam.ID, match2.AwayTeam.ID)
	if err != nil || homeLeague1 == "" || awayLeague1 == "" {
		t.Errorf("Матч %s - %s пропущен: проблема с лигами\n", match2.HomeTeam.Name, match2.AwayTeam.Name)
	}
	avgLeagueWeight1 := (leagueNorm[homeLeague1] + leagueNorm[awayLeague1]) / 2.0
	fmt.Printf("Сила Бундеслиги: %F\n", avgLeagueWeight1)
	// 3) Форма команд
	recentMatchesHome1, err := matchesStore.GetRecentMatches(ctx, match2.HomeTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for home team %d: %v", match2.HomeTeam.ID, err)
	}
	recentMatchesAway1, err := matchesStore.GetRecentMatches(ctx, match2.AwayTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for away team %d: %v", match2.AwayTeam.ID, err)
	}
	homeForm1 := calculateForm(recentMatchesHome1)
	awayForm1 := calculateForm(recentMatchesAway1)
	formFactor1 := (homeForm1 + awayForm1) / 2.0
	fmt.Printf("Формы каждой команды: Bayer 04 Leverkusen - %f; Borussia Dortmund - %f\nОбщая форма - %f\n", homeForm1, awayForm1, formFactor1)
	// 4) Бонусы
	derbyBonus1 := GetDerbyBonus(ctx, teamsService, match1)
	stageBonus1 := 0.0
	if homeLeague1 == "Champions League" && match2.Stage != "" {
		stageBonus1 = CLstage[match2.Stage]
	}
	crossLeagueBonus1 := 0.0
	if homeLeague1 != awayLeague1 {
		crossLeagueBonus1 = 0.15 // Увеличен с 0.1 для большего эффекта
	}
	fmt.Printf("Бонусы: %f, %f, %f\n", derbyBonus, stageBonus, crossLeagueBonus)
	// 5) Финальный рейтинг
	baseRating1 := (homeStrength1+awayStrength1)/2.0*0.15 + // Уменьшено влияние позиций
		avgLeagueWeight1*0.35 + // Снижено с 0.4
		formFactor1*0.15 // Добавлено влияние формы
	rating1 := baseRating1 * (1 + derbyBonus1 + stageBonus1 + crossLeagueBonus1)
	fmt.Printf("Финальный рейтинг: %f\n", rating1)
}

// TestFindMatchesByTeamNames - тестовая функция для поиска матчей по названиям команд
func FindMatchByTeamNames(t *testing.T, matchesService *MatchesService, homeTeamName, awayTeamName string) (types.Match, error) {
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

	// Получаем все матчи
	matches, err := matchesService.HandleGetMatches(context.Background())
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}

	// Названия команд для поиска (вы можете изменить их на нужные)

	var foundMatches []types.Match

	// Ищем матчи где одна из команд играет дома или в гостях
	for _, match := range matches {
		if match.HomeTeam.Name == homeTeamName && match.AwayTeam.Name == awayTeamName {
			foundMatches = append(foundMatches, match)
		}
	}

	fmt.Printf("Found %d matches for teams %s and %s:\n", len(foundMatches), homeTeamName, awayTeamName)

	// Выводим первые 2 найденных матча
	for i, match := range foundMatches {
		fmt.Printf("Match %d: %s vs %s (Date: %s, Status: %s)\n",
			i+1,
			match.HomeTeam.Name,
			match.AwayTeam.Name,
			match.UTCDate,
			match.Status)
	}

	// Проверяем, что нашли хотя бы один матч
	if len(foundMatches) == 0 {
		t.Fatalf("No matches found for teams %s and %s", homeTeamName, awayTeamName)
	}

	// Проверяем, что нашли не менее 2 матчей
	if len(foundMatches) < 2 {
		t.Logf("Found only %d matches, expected at least 2", len(foundMatches))
	} else {
		fmt.Printf("Successfully found %d matches (showing first 2)\n", len(foundMatches))
	}

	return foundMatches[0], nil
}
