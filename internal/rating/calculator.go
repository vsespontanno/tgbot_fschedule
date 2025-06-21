package rating

import (
	"context"
	"fmt"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
	"log"
)

var (
	CLstage = map[string]float64{
		"Round of 16":    0.25,
		"Quarter-finals": 0.50,
		"Semi-finals":    0.75,
		"Final":          1.0,
	}

	leagueNorm = map[string]float64{
		"Champions League": 1.0,
		"PremierLeague":    0.9,
		"LaLiga":           0.8,
		"SerieA":           0.8,
		"Bundesliga":       0.7,
		"Ligue1":           0.65,
	}

	teamsInLeague = map[string]int{
		"PremierLeague": 20,
		"LaLiga":        20,
		"Bundesliga":    18,
		"SerieA":        18,
		"Ligue1":        20,
	}

	derbys = map[[2]string]float64{
		// Англия (Premier League)
		{"Manchester United", "Manchester City"}: 0.27, // Манчестерское дерби
		{"Liverpool", "Everton"}:                 0.16, // Мерсисайдское дерби
		{"Arsenal", "Tottenham"}:                 0.25, // Северолондонское дерби
		{"Chelsea", "Arsenal"}:                   0.25,
		{"Chelsea", "Tottenham"}:                 0.25,
		{"Manchester United", "Liverpool"}:       0.26,
		{"Manchester United", "Leeds United"}:    0.15,
		{"Newcastle", "Sunderland"}:              0.14, // Тайн-Уир (если обе в лиге)

		// Испания (La Liga)
		{"Real Madrid", "Barcelona"}: 0.3,  // Эль Класико
		{"Atletico", "Real Madrid"}:  0.26, // Мадридское дерби
		{"Sevilla", "Real Betis"}:    0.2,  // Севильское дерби
		{"Barcelona", "Espanyol"}:    0.18, // Барселонское дерби
		{"Valencia", "Levante"}:      0.14, // Валенсийское дерби

		// Германия (Bundesliga)
		{"Borussia D.", "Bayern"}:           0.28, // Дер Классикер
		{"Schalke 04", "Borussia Dortmund"}: 0.16, // Рурское дерби
		{"Hamburger SV", "Werder Bremen"}:   0.15, // Северное дерби
		{"Bayern", "1860 Munich"}:           0.14, // Мюнхенское дерби
		{"Cologne", "Borussia M."}:          0.14,

		// Италия (Serie A)
		{"Inter", "Milan"}:     0.29, // Миланское дерби
		{"Roma", "Lazio"}:      0.28, // Римское дерби
		{"Juventus", "Torino"}: 0.2,  // Дерби делла Моле
		{"Genoa", "Sampdoria"}: 0.18, // Дерби делла Лантерна
		{"Napoli", "Roma"}:     0.15,

		// Франция (Ligue 1)
		{"PSG", "Marseille"}:                0.23, // Ле Классик
		{"Olympique Lyon", "Saint-Etienne"}: 0.18, // Ронское дерби
		{"Nice", "Monaco"}:                  0.14, // Лазурное дерби
		{"Lille", "RC Lens"}:                0.14, // Северное дерби
	}
)

// getLeaguesForTeams пытается определить лиги для домашней и гостевой команд.
func getLeaguesForTeams(ctx context.Context, teamsService *service.TeamsService, homeTeamID int, awayTeamID int) (homeLeague string, awayLeague string, err error) {
	// Найти лигу для домашней команды
	foundHomeLeague := false
	for leagueKey := range leagueNorm {
		// league, getLeagueErr := teamsService.HandleGetLeague()

		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, homeTeamID)

		if getLeagueErr != nil {
			// Можно логировать ошибку или обработать ее иначе, пока продолжаем поиск
			// fmt.Printf("Error checking league %s for home team %d: %v\n", leagueKey, homeTeamID, getLeagueErr)
			// Сохраняем первую возникшую ошибку, если она есть
			if err == nil {
				err = fmt.Errorf("error checking league %s for home team %d: %w", leagueKey, homeTeamID, getLeagueErr)
			}
			continue
		}
		if league != "Wrong League" {
			homeLeague = league
			foundHomeLeague = true
			err = nil // Сбрасываем ошибку, если лига найдена
			break
		}
	}
	if !foundHomeLeague {
		fmt.Printf("Команда с id=%d не найдена ни в одной лиге, матч будет пропущен\n", homeTeamID)
		return "", "", nil // или return "", "", fmt.Errorf("not found") и обработать это выше
	}

	// Найти лигу для гостевой команды
	var awayErr error // Отдельная переменная для ошибок гостевой команды
	foundAwayLeague := false
	for leagueKey := range leagueNorm {
		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, awayTeamID)
		if getLeagueErr != nil {
			// fmt.Printf("Error checking league %s for away team %d: %v\n", leagueKey, awayTeamID, getLeagueErr)
			if awayErr == nil {
				awayErr = fmt.Errorf("error checking league %s for away team %d: %w", leagueKey, awayTeamID, getLeagueErr)
			}
			continue
		}
		if league != "Wrong League" {
			awayLeague = league
			foundAwayLeague = true
			break
		}
	}
	if !foundAwayLeague {
		if awayErr == nil {
			return homeLeague, "", fmt.Errorf("could not determine league for away team ID %d", awayTeamID)
		}
		return homeLeague, "", fmt.Errorf("could not determine league for away team ID %d (last error: %w)", awayTeamID, awayErr)
	}
	return homeLeague, awayLeague, nil // Возвращаем найденные лиги и nil в качестве ошибки, если обе найдены
}

func CalculatePositionOfTeams(ctx context.Context, teamsService *service.TeamsService, standingsService *service.StandingsService, match types.Match) (homeTeam, awayTeam float64, err error) {
	var homeTeamRating float64
	var awayTeamRating float64

	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := getLeaguesForTeams(ctx, teamsService, HomeID, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	// Теперь HomeLeague и AwayLeague должны содержать корректные имена лиг
	posHome, err := standingsService.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting home team standing (league: %s, id: %d): %w", HomeLeague, HomeID, err)
	}

	posAway, err := standingsService.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting away team standing (league: %s, id: %d): %w", AwayLeague, AwayID, err)
	}

	homeTeamRating = (float64(teamsInLeague[HomeLeague] - posHome)) / float64(teamsInLeague[HomeLeague]-1)
	awayTeamRating = (float64(teamsInLeague[AwayLeague] - posAway)) / float64(teamsInLeague[AwayLeague]-1)
	return homeTeamRating, awayTeamRating, nil
}

func CalculateRatingOfMatch(ctx context.Context, match types.Match, teamService *service.TeamsService, standingsService *service.StandingsService) (float64, error) {
	// 1) Считаем силу по позиции в таблице (от 0.0 до 1.0)

	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, teamService, standingsService, match)
	if err != nil {
		return 0, fmt.Errorf("error calculating team strengths: %w", err)
	}

	// 2) Определяем лиги
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		return 0, fmt.Errorf("error getting leagues for teams: %w", err)
	}
	if homeLeague == "" || awayLeague == "" {
		fmt.Printf("Матч %s - %s пропущен: одна из команд не найдена в нужных лигах\n", match.HomeTeam.Name, match.AwayTeam.Name)
		return 0, nil
	}

	// 3) Средний вес лиги (раньше вы брали сумму, теперь — среднее)
	avgLeagueWeight := (leagueNorm[homeLeague] + leagueNorm[awayLeague]) / 2.0

	// 4) Бонусы (пока нулевые, добавите позже дерби и стадию ЛЧ)
	var derbyBonus float64 // = 0.15 если дерби
	var stageBonus float64 // = CLstage[stage] * stageWeight

	// derbyBonus = GetDerbyBonus(ctx, teamService, match)

	rating := ((homeStrength+awayStrength)/2.0)*0.2 +
		avgLeagueWeight*0.4 +
		derbyBonus +
		stageBonus

	// Гарантируем, что не выйдем за 1.0
	if rating > 1.0 {
		rating = 1.0
	}
	return rating, nil
}

func GetDerbyBonus(ctx context.Context, teamService *service.TeamsService, match types.Match) float64 {
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		log.Fatal(err)
	}
	homeShortName, err := teamService.HandleGetTeamShortName(ctx, homeLeague, match.HomeTeam.Name)
	if err != nil {
		log.Fatal(err)
	}
	awayShortName, err := teamService.HandleGetTeamShortName(ctx, awayLeague, match.AwayTeam.Name)
	if err != nil {
		log.Fatal(err)
	}
	key := [2]string{homeShortName, awayShortName}
	reverseKey := [2]string{awayShortName, homeShortName}

	if bonus, ok := derbys[key]; ok {
		return bonus
	}
	if bonus, ok := derbys[reverseKey]; ok {
		return bonus
	}
	return 0.0
}
