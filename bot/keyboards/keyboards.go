package keyboards

import (
	"football_tgbot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (

	// создание мапы для выбора таблицы
	KeyboardsStandings = map[string]types.League{
		"standings_APL":        types.Leagues["PremierLeague"],
		"standings_LaLiga":     types.Leagues["LaLiga"],
		"standings_Bundesliga": types.Leagues["Bundesliga"],
		"standings_SerieA":     types.Leagues["SerieA"],
		"standings_Ligue1":     types.Leagues["Ligue1"],
		"standings_CL":         types.Leagues["ChampionsLeague"],
	}
	KeyboardsSchedule = map[string]types.League{
		"schedule_laliga":     types.Leagues["LaLiga"],
		"schedule_epl":        types.Leagues["PremierLeague"],
		"schedule_primeira":   types.Leagues["Primeira"],
		"schedule_eredivisie": types.Leagues["Eredivisie"],
		"schedule_bundesliga": types.Leagues["Bundesliga"],
		"schedule_seriea":     types.Leagues["SerieA"],
		"schedule_ucl":        types.Leagues["ChampionsLeague"],
		"schedule_uel":        types.Leagues["EuropaLeague"],
	}

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
	KeyboardSchedule = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "schedule_laliga"),
			tgbotapi.NewInlineKeyboardButtonData("EPL", "schedule_epl"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Primeira", "schedule_primeira"),
			tgbotapi.NewInlineKeyboardButtonData("Eredivisie", "schedule_eredivisie"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "schedule_bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "schedule_seriea"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("UCL", "schedule_ucl"),
			tgbotapi.NewInlineKeyboardButtonData("UEL", "schedule_uel"),
		),
	)
)
