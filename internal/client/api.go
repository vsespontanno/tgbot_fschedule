package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vsespontanno/tgbot_fschedule/internal/tools"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

type MatchApiClient interface {
	FetchMatches(ctx context.Context, from, to string) ([]types.Match, error)
}

type StandingsApiClient interface {
	FetchStandings(ctx context.Context, leagueCode string) ([]types.Standing, error)
}

type TeamsApiClient interface {
	FetchTeams(ctx context.Context, leagueCode string) ([]types.Team, error)
}

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

func (m *FootballAPIClient) FetchMatches(ctx context.Context, from, to string) ([]types.Match, error) {

	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", from, to)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", m.apiKey)

	resp, err := m.httpClient.Do(req)
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

	filteredMatches := tools.MatchFilter(MatchesResponse)

	return filteredMatches, nil
}

func (m *FootballAPIClient) FetchStandings(ctx context.Context, leagueCode string) ([]types.Standing, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings?season=2024", leagueCode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Add("X-Auth-Token", m.apiKey)
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	var standingsResponse types.StandingsResponse
	err = json.Unmarshal(body, &standingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %s, body: %s", err, string(body))
	}

	if len(standingsResponse.Standings) > 0 && len(standingsResponse.Standings[0].Table) > 0 {
		tools.StandingsFilter(standingsResponse.Standings[0].Table)
		return standingsResponse.Standings[0].Table, nil
	}
	return nil, fmt.Errorf("no standings found for league code: %s", leagueCode)
}

func (m *FootballAPIClient) FetchTeams(ctx context.Context, leagueCode string) ([]types.Team, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/teams", leagueCode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", m.apiKey)
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var teamsResponse types.TeamsResponse
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return nil, err
	}
	tools.TeamsFilter(teamsResponse)

	return teamsResponse.Teams, nil
}
