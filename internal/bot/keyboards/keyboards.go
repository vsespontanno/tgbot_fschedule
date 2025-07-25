package keyboards

import (
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (

	// Мапа для выбора лиги для турнирной таблицы на основе нажатой кнопки
	KeyboardsStandings = map[string]types.League{
		"standings_EPL":        types.Leagues["PremierLeague"],
		"standings_LaLiga":     types.Leagues["LaLiga"],
		"standings_Bundesliga": types.Leagues["Bundesliga"],
		"standings_SerieA":     types.Leagues["SerieA"],
		"standings_Ligue1":     types.Leagues["Ligue1"],
		"standings_CL":         types.Leagues["ChampionsLeague"],
	}

	// Мапа для выбора лиги для расписания матчей на основе нажатой кнопки
	KeyboardsSchedule = map[string]types.League{
		"schedule_LaLiga":     types.Leagues["LaLiga"],
		"schedule_EPL":        types.Leagues["PremierLeague"],
		"schedule_Primeira":   types.Leagues["Primeira"],
		"schedule_Eredivisie": types.Leagues["Eredivisie"],
		"schedule_Bundesliga": types.Leagues["Bundesliga"],
		"schedule_SerieA":     types.Leagues["SerieA"],
		"schedule_UCL":        types.Leagues["ChampionsLeague"],
		"schedule_UEL":        types.Leagues["EuropaLeague"],
	}

	// Инлайн-клавиатура для выбора лиг для турнирной таблицы
	KeyboardStandings = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("EPL", "standings_EPL"),
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

	// Инлайн-клавиатура для выбора лиг
	KeyboardDefaultSchedule = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("La Liga", "schedule_LaLiga"),
			tgbotapi.NewInlineKeyboardButtonData("EPL", "schedule_EPL"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Primeira", "schedule_Primeira"),
			tgbotapi.NewInlineKeyboardButtonData("Eredivisie", "schedule_Eredivisie"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "schedule_Bundesliga"),
			tgbotapi.NewInlineKeyboardButtonData("Serie A", "schedule_SerieA"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("UCL", "schedule_UCL"),
			tgbotapi.NewInlineKeyboardButtonData("UEL", "schedule_UEL"),
		),
	)
	Keyboard_Schedule = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все матчи", "show_all_matches"),
			tgbotapi.NewInlineKeyboardButtonData("Топ матчи", "show_top_matches"),
		),
	)
)
