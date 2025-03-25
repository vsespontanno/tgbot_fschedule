package db

import "football_tgbot/bot/models"

var Leagues = map[string]models.League{
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
}
