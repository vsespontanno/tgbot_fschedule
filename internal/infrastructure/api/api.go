package api

import (
	"context"
	"encoding/json"
	"fmt"
	"football_tgbot/internal/types"
	"io"
	"net/http"
)

type FootballDataClient interface {
	// GetMatches возвращает список матчей за указанный период
	GetMatches(ctx context.Context, from, to string) ([]types.Match, error)
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

func (m *FootballAPIClient) GetMatches(ctx context.Context, from, to string) ([]types.Match, error) {

	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", from, to)
	// req, err := http.NewRequest("GET", today)
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

	Matches := Mapper(MatchesResponse)

	return Matches, nil
}
