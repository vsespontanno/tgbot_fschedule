package types

// Мапа для хранения информации о лигах
// То, что закомментировано ниже, не используется в коде, но оставлено для возможного будущего использования
// Если нужно будет использовать эти лиги, раскомментируйте и добавьте их в мапу Leagues
// Если нет, то удалите эти комментарии и не используйте
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
	// "EuropaLeague": {
	// 	Name:           "Europa League",
	// 	CollectionName: "EuropaLeague",
	// 	Code:           "EC",
	// },

	// "DFB-Pokal": {
	// 	Name:           "DFB Pokal",
	// 	CollectionName: "DFB-Pokal",
	// 	Code:           "DFB",
	// },

	// "FA Cup": {
	// 	Name:           "FA Cup",
	// 	CollectionName: "FA Cup",
	// 	Code:           "FAC",
	// },

	// "Copa del Rey": {
	// 	Name:           "Copa del Rey",
	// 	CollectionName: "Copa del Rey",
	// 	Code:           "CDR",
	// },

	// "World Cup": {
	// 	Name:           "World Cup",
	// 	CollectionName: "World Cup",
	// 	Code:           "WC",
	// },
}

// Структура для хранения информации о лиге
type League struct {
	Name           string
	CollectionName string
	Code           string
}
