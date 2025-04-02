package rating

import (
	"football_tgbot/types"
	"sort"
	"time"
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
	var teamMatches []types.Match

	// Фильтруем матчи команды
	for _, match := range matches {
		if match.HomeTeam.ID == teamID || match.AwayTeam.ID == teamID {
			teamMatches = append(teamMatches, match)
		}
	}

	// Сортируем по дате (новые первыми)
	sort.Slice(teamMatches, func(i, j int) bool {
		dateI, _ := time.Parse(time.RFC3339, teamMatches[i].UTCDate)
		dateJ, _ := time.Parse(time.RFC3339, teamMatches[j].UTCDate)
		return dateI.After(dateJ)
	})

	// Возвращаем последние N матчей
	if len(teamMatches) > n {
		return teamMatches[:n]
	}
	return teamMatches
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

// CalculateMatchRating вычисляет рейтинг матча на основе рейтингов команд
func (c *Calculator) CalculateMatchRating(match types.Match, homeRating, awayRating types.TeamRating) float64 {
	// Рейтинг матча - среднее значение рейтингов команд
	return (homeRating.CalculateRating() + awayRating.CalculateRating()) / 2
}
