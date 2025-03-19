package types

type Match struct {
	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`
	Date     string `json:"date"`
}

type MatchResponse struct {
	Matches []Match `json:"matches"`
}

type Club struct {
	Name           string `json:"name"`
	Ligue          string `json:"ligue"`
	FoundationDate string `json:"foundation_date"`
	Stadium        string `json:"stadium"`
}
