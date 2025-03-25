package keyboards

import (
	"football_tgbot/bot/models"
	"football_tgbot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Leagues = map[string]models.League{
		"league_APL":        db.Leagues["PremierLeague"],
		"league_LaLiga":     db.Leagues["LaLiga"],
		"league_Bundesliga": db.Leagues["Bundesliga"],
		"league_SerieA":     db.Leagues["SerieA"],
	}

	Standings = map[string]models.League{
		"standings_APL":        db.Leagues["PremierLeague"],
		"standings_LaLiga":     db.Leagues["LaLiga"],
		"standings_Bundesliga": db.Leagues["Bundesliga"],
		"standings_SerieA":     db.Leagues["SerieA"],
	}

	KeyboardLeagues = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("APL", "league_APL"),
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "league_LaLiga"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "league_Bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "league_SerieA"),
		),
	)

	KeyboardStandings = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("APL", "standings_APL"),
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "standings_LaLiga"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "standings_Bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "standings_SerieA"),
		),
	)
)
