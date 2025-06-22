package types

// структура для хранения информации о матче
type Match struct {
	ID          int `json:"id"`
	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"competition"`
	Stage    string `json:"stage"`
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
		Winner   string `json:"winner"`
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
