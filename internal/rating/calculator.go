package rating

import (
	"encoding/json"
	"fmt"
	"football_tgbot/internal/types"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/joho/godotenv"
)

// Calculator структура для расчета рейтингов
type Calculator struct {
	// Веса для разных компонентов рейтинга
	weights struct {
		position   float64
		points     float64
		form       float64
		goalDiff   float64
		tournament float64
	}
}

// NewCalculator создает новый калькулятор рейтингов
func NewCalculator() *Calculator {
	calc := &Calculator{}
	// Инициализация весов
	calc.weights.position = 0.3
	calc.weights.points = 0.3
	calc.weights.form = 0.2
	calc.weights.goalDiff = 0.1
	calc.weights.tournament = 0.1
	return calc
}

// CalculateTeamRating вычисляет рейтинг команды на основе всех доступных данных
func (c *Calculator) CalculateTeamRating(standing types.Standing, matches []types.Match, teamID int) types.TeamRating {
	// Получаем форму команды
	form := c.calculateForm(matches, teamID)

	// Получаем вес турнира
	tournamentWeight := c.getTournamentWeight(standing.Team.Area.Code)

	rating := types.TeamRating{
		TeamID:           teamID,
		TeamName:         standing.Team.Name,
		Position:         standing.Position,
		Points:           standing.Points,
		Form:             form,
		GoalDiff:         standing.GoalDifference,
		TournamentWeight: tournamentWeight,
		LastUpdated:      time.Now().Format(time.RFC3339),
	}

	return rating
}

// calculateForm вычисляет форму команды на основе последних матчей
func (c *Calculator) calculateForm(matches []types.Match, teamID int) float64 {
	// Берем последние 5 матчей
	lastMatches := c.getLastMatches(matches, teamID, 5)
	if len(lastMatches) == 0 {
		return 0.5 // Нейтральная форма, если нет матчей
	}

	var form float64
	for _, match := range lastMatches {
		// Определяем результат матча для команды
		var result float64
		if match.Status == "FINISHED" {
			if match.HomeTeam.ID == teamID {
				if match.Score.FullTime.Home > match.Score.FullTime.Away {
					result = 1.0 // Победа
				} else if match.Score.FullTime.Home == match.Score.FullTime.Away {
					result = 0.5 // Ничья
				} else {
					result = 0.0 // Поражение
				}
			} else if match.AwayTeam.ID == teamID {
				if match.Score.FullTime.Away > match.Score.FullTime.Home {
					result = 1.0 // Победа
				} else if match.Score.FullTime.Home == match.Score.FullTime.Away {
					result = 0.5 // Ничья
				} else {
					result = 0.0 // Поражение
				}
			}
		}
		form += result
	}

	return form / float64(len(lastMatches))
}

// getLastMatches возвращает последние N матчей команды
func (c *Calculator) getLastMatches(matches []types.Match, teamID int, n int) []types.Match {
	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Формируем URL для запроса матчей команды
	// Замените YOUR_API_KEY на реальный ключ API
	url := fmt.Sprintf("http://api.football-data.org/v4/teams/%d/matches?limit=%d&status=FINISHED", teamID, n)
	fmt.Println("---stage 1")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return matches // Возвращаем переданные матчи как fallback
	}
	fmt.Println("---stage 2")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Println("---stage 3")
	// Добавляем заголовок с API ключом
	req.Header.Add("X-Auth-Token", os.Getenv("FOOTBALL_API_KEY"))
	fmt.Println("---stage 4")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching matches: %v", err)
		return matches // Возвращаем переданные матчи как fallback
	}
	defer resp.Body.Close()
	fmt.Println("---stage 5")
	var apiResponse struct {
		Matches []types.Match `json:"matches"`
	}
	fmt.Printf("---stage 6, apiResponse: %v\n", apiResponse)
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return matches // Возвращаем переданные матчи как fallback
	}
	fmt.Printf("---stage 7, apiResponse: %v\n", apiResponse)
	// Сортируем по дате (новые первыми)
	sort.Slice(apiResponse.Matches, func(i, j int) bool {
		dateI, _ := time.Parse(time.RFC3339, apiResponse.Matches[i].UTCDate)
		dateJ, _ := time.Parse(time.RFC3339, apiResponse.Matches[j].UTCDate)
		return dateI.After(dateJ)
	})
	fmt.Printf("---stage 8, apiResponse: %v\n", apiResponse)
	return apiResponse.Matches
}

// getTournamentWeight возвращает вес турнира
func (c *Calculator) getTournamentWeight(areaCode string) float64 {
	// Определяем вес турнира на основе кода страны/региона
	switch areaCode {
	case "ENG": // Англия
		return 0.9
	case "ESP": // Испания
		return 0.9
	case "GER": // Германия
		return 0.8
	case "ITA": // Италия
		return 0.8
	case "FRA": // Франция
		return 0.8
	case "UEFA": // Лига чемпионов/Лига Европы
		return 1.0
	default:
		return 0.5 // Для остальных турниров
	}
}

// CalculateMatchRating вычисляет рейтинг матча на основе рейтингов команд и других факторов
func (c *Calculator) CalculateMatchRating(match types.Match, home types.TeamRating, away types.TeamRating) float64 {
	base := (home.CalculateRating() + away.CalculateRating()) / 2

	// Бонус за дерби (пример)
	if isDerby(match) {
		base += 0.1
	}

	// Бонус за матч лидеров
	if home.Position <= 2 && away.Position <= 2 {
		base += 0.1
	}

	// Ограничение диапазона
	if base > 1.0 {
		base = 1.0
	}
	return base
}

// isDerby определяет, является ли матч дерби (пример)
func isDerby(match types.Match) bool {
	derbies := [][]string{
		{"Real Madrid", "Atletico Madrid"},
		{"Manchester United", "Manchester City"},
		// TDOD : Другие дерби
	}
	for _, pair := range derbies {
		if (match.HomeTeam.Name == pair[0] && match.AwayTeam.Name == pair[1]) ||
			(match.HomeTeam.Name == pair[1] && match.AwayTeam.Name == pair[0]) {
			return true
		}
	}
	return false
}
