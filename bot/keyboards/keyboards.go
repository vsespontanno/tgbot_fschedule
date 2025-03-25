package keyboards

import (
	"football_tgbot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Leagues = map[string]types.League{
		"league_APL":        types.Leagues["PremierLeague"],
		"league_LaLiga":     types.Leagues["LaLiga"],
		"league_Bundesliga": types.Leagues["Bundesliga"],
		"league_SerieA":     types.Leagues["SerieA"],
	}

	Standings = map[string]types.League{
		"standings_APL":        types.Leagues["PremierLeague"],
		"standings_LaLiga":     types.Leagues["LaLiga"],
		"standings_Bundesliga": types.Leagues["Bundesliga"],
		"standings_SerieA":     types.Leagues["SerieA"],
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
