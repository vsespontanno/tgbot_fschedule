package types

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
