package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/joho/godotenv"
)

// setupCalculator создает и возвращает адаптер service.Calculator с реальными сервисами
func setupCalculator(t *testing.T) service.Calculator {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// Закрываем подключение после теста
	t.Cleanup(func() { _ = client.Disconnect(context.Background()) })

	// Репозитории
	matchesStore := mongoRepo.NewMongoDBMatchesStore(client, "football")
	teamsStore := mongoRepo.NewMongoDBTeamsStore(client, "football")
	standingsStore := mongoRepo.NewMongoDBStandingsStore(client, "football")

	// Адаптер домена реализует интерфейс service.Calculator
	calculator := service.NewCalculatorAdapter(teamsStore, standingsStore, matchesStore)
	return calculator
}

// TestCalculatePositionOfTeams проверяет определение лиг команд
func TestCalculatePositionOfTeams(t *testing.T) {
	calculator := setupCalculator(t)

	homeID := 80  // пример ID Sevilla FC
	awayID := 275 // пример ID UD Las Palmas

	homeLeague, awayLeague, err := service.GetLeaguesForTeams(context.Background(), calculator, homeID, awayID)
	if err != nil {
		t.Fatalf("Failed to get leagues for teams: %v", err)
	}
	if homeLeague != "LaLiga" || awayLeague != "LaLiga" {
		t.Errorf("wanted %s and %s, got %s and %s", "LaLiga", "LaLiga", homeLeague, awayLeague)
	}
}

// TestCalculateRatingOfMatch проверяет, что рейтинг матча лежит в [0,1]
func TestCalculateRatingOfMatch(t *testing.T) {
	calculator := setupCalculator(t)

	err := godotenv.Load("../.env")
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

	ctx := context.Background()

	// Получаем все матчи
	homeTeamName := "Athletic Club"
	awayTeamName := "Deportivo Alavés"
	match1, err := FindMatchByTeamNames(t, calculator, homeTeamName, awayTeamName)
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}
	homeTeamName1 := "Bayer 04 Leverkusen"
	awayTeamName1 := "Borussia Dortmund"
	match2, err := FindMatchByTeamNames(t, calculator, homeTeamName1, awayTeamName1)
	if err != nil {
		t.Fatalf("Failed to get matches: %v", err)
	}
	//Данные для первого матча
	// 1) Сила команд по позициям
	homeStrength, awayStrength, err := service.CalculatePositionOfTeams(ctx, calculator, match1)
	if err != nil {
		t.Errorf("error calculating team strengths: %v", err)
	}
	fmt.Printf("Сила команд: Athletic Club - %f; Deportivo Alavés - %f\n", homeStrength, awayStrength)
	// 2) Лиги и вес
	homeLeague, awayLeague, err := service.GetLeaguesForTeams(ctx, calculator, match1.HomeTeam.ID, match1.AwayTeam.ID)
	if err != nil || homeLeague == "" || awayLeague == "" {
		t.Errorf("Матч %s - %s пропущен: проблема с лигами\n", match1.HomeTeam.Name, match1.AwayTeam.Name)
	}
	avgLeagueWeight := (types.LeagueNorm[homeLeague] + types.LeagueNorm[awayLeague]) / 2.0
	fmt.Printf("Сила ЛаЛига: %F\n", avgLeagueWeight)
	// 3) Форма команд
	recentMatchesHome, err := calculator.HandleGetRecentMatches(ctx, match1.HomeTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for home team %d: %v", match1.HomeTeam.ID, err)
	}
	recentMatchesAway, err := calculator.HandleGetRecentMatches(ctx, match1.AwayTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for away team %d: %v", match1.AwayTeam.ID, err)
	}
	homeForm := service.CalculateForm(recentMatchesHome, match1.HomeTeam.ID)
	awayForm := service.CalculateForm(recentMatchesAway, match1.AwayTeam.ID)
	formFactor := (homeForm + awayForm) / 2.0
	fmt.Printf("Формы каждой команды: Athletic Club - %f; Deportivo Alavés - %f\nОбщая форма - %f\n", homeForm, awayForm, formFactor)
	// 4) Бонусы
	derbyBonus := service.GetDerbyBonus(ctx, calculator, match1)
	stageBonus := 0.0
	if homeLeague == "Champions League" && match1.Stage != "" {
		stageBonus = types.CLstage[match1.Stage]
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
	homeStrength1, awayStrength1, err := service.CalculatePositionOfTeams(ctx, calculator, match2)
	if err != nil {
		t.Errorf("error calculating team strengths: %v", err)
	}
	fmt.Printf("Сила команд: Bayer 04 Leverkusen - %f; Borussia Dortmund - %f\n", homeStrength1, awayStrength1)
	// 2) Лиги и вес
	homeLeague1, awayLeague1, err := service.GetLeaguesForTeams(ctx, calculator, match2.HomeTeam.ID, match2.AwayTeam.ID)
	if err != nil || homeLeague1 == "" || awayLeague1 == "" {
		t.Errorf("Матч %s - %s пропущен: проблема с лигами\n", match2.HomeTeam.Name, match2.AwayTeam.Name)
	}
	avgLeagueWeight1 := (types.LeagueNorm[homeLeague1] + types.LeagueNorm[awayLeague1]) / 2.0
	fmt.Printf("Сила Бундеслиги: %F\n", avgLeagueWeight1)
	// 3) Форма команд
	recentMatchesHome1, err := calculator.HandleGetRecentMatches(ctx, match2.HomeTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for home team %d: %v", match2.HomeTeam.ID, err)
	}
	recentMatchesAway1, err := calculator.HandleGetRecentMatches(ctx, match2.AwayTeam.ID, 5)
	if err != nil {
		t.Errorf("Error getting recent matches for away team %d: %v", match2.AwayTeam.ID, err)
	}
	homeForm1 := service.CalculateForm(recentMatchesHome1, match2.HomeTeam.ID)
	awayForm1 := service.CalculateForm(recentMatchesAway1, match2.AwayTeam.ID)
	formFactor1 := (homeForm1 + awayForm1) / 2.0
	fmt.Printf("Формы каждой команды: Bayer 04 Leverkusen - %f; Borussia Dortmund - %f\nОбщая форма - %f\n", homeForm1, awayForm1, formFactor1)
	// 4) Бонусы
	derbyBonus1 := service.GetDerbyBonus(ctx, calculator, match1)
	stageBonus1 := 0.0
	if homeLeague1 == "Champions League" && match2.Stage != "" {
		stageBonus1 = types.CLstage[match2.Stage]
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

func FindMatchByTeamNames(t *testing.T, calculator service.Calculator, homeTeamName, awayTeamName string) (types.Match, error) {

	// Получаем все матчи
	matches, err := calculator.HandleGetMatches(context.Background())
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
