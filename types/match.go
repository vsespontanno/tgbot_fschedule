package types

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

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type MatchesResponse struct {
	Matches []Match `json:"matches"`
}

type Club struct {
	Name           string `json:"name"`
	Ligue          string `json:"ligue"`
	FoundationDate string `json:"foundation_date"`
	Stadium        string `json:"stadium"`
}

type Area struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Tla       string `json:"tla"`
	CrestURL  string `json:"crestUrl"`
	Area      Area   `json:"area"`
	Founded   int    `json:"founded"`
}
