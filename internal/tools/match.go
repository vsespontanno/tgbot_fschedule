package tools

import "github.com/vsespontanno/tgbot_fschedule/internal/types"

func MatchFilter(MatchesResponse types.MatchesResponse) []types.Match {
	for i, match := range MatchesResponse.Matches {
		switch match.HomeTeam.Name {
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
		switch match.AwayTeam.Name {
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

		switch match.Competition.Name {
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
	return filteredMatches
}
