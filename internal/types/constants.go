package types

// мапа для хранения информации о лигах
var Leagues = map[string]League{
	"PremierLeague": {
		Name:           "APL",
		CollectionName: "PremierLeague",
		Code:           "PL",
	},
	"LaLiga": {
		Name:           "La Liga",
		CollectionName: "LaLiga",
		Code:           "PD",
	},
	"Bundesliga": {
		Name:           "Bundesliga",
		CollectionName: "Bundesliga",
		Code:           "BL1",
	},
	"SerieA": {
		Name:           "Serie A",
		CollectionName: "SerieA",
		Code:           "SA",
	},
	"Ligue1": {
		Name:           "Ligue 1",
		CollectionName: "Ligue1",
		Code:           "FL1",
	},
	"ChampionsLeague": {
		Name:           "Champions League",
		CollectionName: "ChampionsLeague",
		Code:           "CL",
	},
}

// константа для ответа на команду help
const HelpText = `Доступные команды:
/start - Начать работу с ботом
/help - Получить список команд
/leagues - Показать список футбольных лиг
/schedule - Показать расписание матчей
/standings - Показать таблицы лиг`

// структура для хранения информации о лиге
type League struct {
	Name           string
	CollectionName string
	Code           string
}
