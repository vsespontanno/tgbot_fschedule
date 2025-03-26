package keyboards

import (
	"football_tgbot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// создание мапы для выбора лиги
	KeyboardsLeagues = map[string]types.League{
		"league_APL":        types.Leagues["PremierLeague"],
		"league_LaLiga":     types.Leagues["LaLiga"],
		"league_Bundesliga": types.Leagues["Bundesliga"],
		"league_SerieA":     types.Leagues["SerieA"],
		"league_Ligue1":     types.Leagues["Ligue1"],
		"league_CL":         types.Leagues["ChampionsLeague"],
	}

	// создание мапы для выбора таблицы
	KeyboardsStandings = map[string]types.League{
		"standings_APL":        types.Leagues["PremierLeague"],
		"standings_LaLiga":     types.Leagues["LaLiga"],
		"standings_Bundesliga": types.Leagues["Bundesliga"],
		"standings_SerieA":     types.Leagues["SerieA"],
		"standings_Ligue1":     types.Leagues["Ligue1"],
		"standings_CL":         types.Leagues["ChampionsLeague"],
	}

	// создание клавиатуры для выбора лиги
	KeyboardLeagues = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("APL", "league_APL"),
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "league_LaLiga"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "league_Bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "league_SerieA"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ligue 1", "league_Ligue1"),
			tgbotapi.NewInlineKeyboardButtonData("Champions League", "league_CL"),
		),
	)

	// создание клавиатуры для выбора таблицы
	KeyboardStandings = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("APL", "standings_APL"),
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "standings_LaLiga"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "standings_Bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "standings_SerieA"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ligue 1", "standings_Ligue1"),
			tgbotapi.NewInlineKeyboardButtonData("Champions League", "standings_CL"),
		),
	)
)
