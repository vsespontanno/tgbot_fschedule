package types

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

type League struct {
	Name           string
	CollectionName string
	Code           string
}
