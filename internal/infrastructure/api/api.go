package api

import (
	"encoding/json"
	"fmt"
	"football_tgbot/internal/types"
	"io"
	"net/http"
)

type FootballAPIClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewFootballAPIClient(httpClient *http.Client, apiKey string) *FootballAPIClient {
	return &FootballAPIClient{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

func (m *FootballAPIClient) GetMatchesSchedule(httpclient *http.Client, apiKey, from, to string) ([]types.Match, error) {

	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", from, to)
	// req, err := http.NewRequest("GET", today)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)

	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	var MatchesResponse types.MatchesResponse
	err = json.Unmarshal(body, &MatchesResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	leaguesSet := make(map[string]struct{})
	for _, match := range MatchesResponse.Matches {
		leaguesSet[match.Competition.Name] = struct{}{}
	}
	for i := range MatchesResponse.Matches {
		switch MatchesResponse.Matches[i].HomeTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse.Matches[i].HomeTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse.Matches[i].HomeTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse.Matches[i].HomeTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse.Matches[i].HomeTeam.Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse.Matches[i].HomeTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse.Matches[i].HomeTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse.Matches[i].HomeTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse.Matches[i].HomeTeam.Name = "Real Sociedad"
		}
		// Corrected loop for AwayTeam
		switch MatchesResponse.Matches[i].AwayTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse.Matches[i].AwayTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse.Matches[i].AwayTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse.Matches[i].AwayTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "AtLetico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse.Matches[i].AwayTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Sociedad"
		}

		switch MatchesResponse.Matches[i].Competition.Name {
		case "UEFA Champions League":
			MatchesResponse.Matches[i].Competition.Name = "UCL"
		case "UEFA Europa League":
			MatchesResponse.Matches[i].Competition.Name = "UEL"
		case "Primera Division":
			MatchesResponse.Matches[i].Competition.Name = "LaLiga"
		case "Primeira Liga":
			MatchesResponse.Matches[i].Competition.Name = "Primeira"
		case "Premier League":
			MatchesResponse.Matches[i].Competition.Name = "EPL"
		case "Serie A":
			MatchesResponse.Matches[i].Competition.Name = "SerieA"
		case "Bundesliga":
			MatchesResponse.Matches[i].Competition.Name = "Bundesliga"
		case "Ligue 1":
			MatchesResponse.Matches[i].Competition.Name = "Ligue1"
		case "Eredivisie":
			MatchesResponse.Matches[i].Competition.Name = "Eredivisie"
		}
	}

	leaguesSet2 := make(map[string]struct{})
	for _, match := range MatchesResponse.Matches {
		leaguesSet2[match.Competition.Name] = struct{}{}
	}

	// Фильтруем матчи только нужных лиг
	var filteredMatches []types.Match
	allowedLeagues := map[string]bool{
		"LaLiga":     true,
		"EPL":        true,
		"Bundesliga": true,
		"SerieA":     true,
		"Ligue1":     true,
		"UCL":        true,
	}
	for _, match := range MatchesResponse.Matches {
		if allowedLeagues[match.Competition.Name] {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches, nil
}
