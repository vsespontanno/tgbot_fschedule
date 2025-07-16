package tools

import "github.com/vsespontanno/tgbot_fschedule/internal/types"

func MatchFilter(MatchesResponse []types.Match) []types.Match {
	for i := range MatchesResponse {
		switch MatchesResponse[i].HomeTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse[i].HomeTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse[i].HomeTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse[i].HomeTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse[i].HomeTeam.Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse[i].HomeTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse[i].HomeTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse[i].HomeTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse[i].HomeTeam.Name = "Real Sociedad"
		}
		// Corrected loop for AwayTeam
		switch MatchesResponse[i].AwayTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse[i].AwayTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse[i].AwayTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse[i].AwayTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse[i].AwayTeam.Name = "AtLetico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse[i].AwayTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse[i].AwayTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse[i].AwayTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse[i].AwayTeam.Name = "Real Sociedad"
		}

		switch MatchesResponse[i].Competition.Name {
		case "UEFA Champions League":
			MatchesResponse[i].Competition.Name = "UCL"
		case "UEFA Europa League":
			MatchesResponse[i].Competition.Name = "UEL"
		case "Primera Division":
			MatchesResponse[i].Competition.Name = "LaLiga"
		case "Primeira Liga":
			MatchesResponse[i].Competition.Name = "Primeira"
		case "Premier League":
			MatchesResponse[i].Competition.Name = "EPL"
		case "Serie A":
			MatchesResponse[i].Competition.Name = "SerieA"
		case "Bundesliga":
			MatchesResponse[i].Competition.Name = "Bundesliga"
		case "Ligue 1":
			MatchesResponse[i].Competition.Name = "Ligue1"
		case "Eredivisie":
			MatchesResponse[i].Competition.Name = "Eredivisie"
		}
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
	for _, match := range MatchesResponse {
		if allowedLeagues[match.Competition.Name] {
			filteredMatches = append(filteredMatches, match)
		}
	}
	return filteredMatches
}
