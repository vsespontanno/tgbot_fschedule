package types

// структура для хранения информации о матче
type Match struct {
	ID          int `json:"id"`
	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"competition"`
	HomeTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homeTeam"`
	AwayTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"awayTeam"`
	UTCDate string `json:"utcDate"`
	Status  string `json:"status"`
	Score   struct {
		FullTime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fullTime"`
	} `json:"score"`
	Rating float64 `json:"rating"`
}

// структура для хранения информации о матчах
type MatchesResponse struct {
	Matches []Match `json:"matches"`
}

// структура для хранения информации о площадке
type Area struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// структура для хранения информации о таблице
type Standing struct {
	Position       int  `json:"position" bson:"position"`
	Team           Team `json:"team" bson:"team"`
	PlayedGames    int  `json:"playedGames" bson:"playedgames"`
	Won            int  `json:"won" bson:"won"`
	Draw           int  `json:"draw" bson:"draw"`
	Lost           int  `json:"lost" bson:"lost"`
	Points         int  `json:"points" bson:"points"`
	GoalsFor       int  `json:"goalsFor" bson:"goalsfor"`
	GoalsAgainst   int  `json:"goalsAgainst" bson:"goalsagainst"`
	GoalDifference int  `json:"goalDifference" bson:"goaldifference"`
}

// структура для хранения информации о таблице
type StandingsResponse struct {
	Standings []struct {
		Table []Standing `json:"table"`
	} `json:"standings"`
}
