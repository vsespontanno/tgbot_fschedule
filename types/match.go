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
}

// структура для хранения информации о командах
type TeamsResponse struct {
	Teams []Team `json:"teams"`
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

// структура для хранения информации о команде
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Tla       string `json:"tla"`
	CrestURL  string `json:"crestUrl"`
	Area      Area   `json:"area"`
	Founded   int    `json:"founded"`
}

// структура для хранения информации о таблице
type Standing struct {
	Position       int  `json:"position"`
	Team           Team `json:"team"`
	PlayedGames    int  `json:"playedGames"`
	Won            int  `json:"won"`
	Draw           int  `json:"draw"`
	Lost           int  `json:"lost"`
	Points         int  `json:"points"`
	GoalsFor       int  `json:"goalsFor"`
	GoalsAgainst   int  `json:"goalsAgainst"`
	GoalDifference int  `json:"goalDifference"`
}

// структура для хранения информации о таблице
type StandingsResponse struct {
	Standings []struct {
		Table []Standing `json:"table"`
	} `json:"standings"`
}
