package tools

import "github.com/vsespontanno/tgbot_fschedule/internal/types"

func TeamsFilter(teams []types.Team, leagueName string) {
	for i := range teams {
		teams[i].League = leagueName
		switch teams[i].Name {
		case "Sevilla FC":
			teams[i].ShortName = "Sevilla"
		case "Wolverhampton Wanderers FC":
			teams[i].Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			teams[i].Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			teams[i].Name = "Inter"
		case "Club Atlético de Madrid":
			teams[i].Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			teams[i].Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			teams[i].Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			teams[i].Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			teams[i].Name = "Real Sociedad"
		}
	}

	for i := range teams {
		switch teams[i].ShortName {
		case "Sevilla FC":
			teams[i].ShortName = "Sevilla"
		case "Leverkusen":
			teams[i].ShortName = "Bayer"
		case "Dortmund":
			teams[i].ShortName = "Borussia D."
		case "M'gladbach":
			teams[i].ShortName = "Borussia M."
		case "Atleti":
			teams[i].ShortName = "Atletico"
		case "Barça":
			teams[i].ShortName = "Barcelona"
		case "Leganés":
			teams[i].ShortName = "Leganes"
		case "Man United":
			teams[i].ShortName = "Manchester United"
		case "Man City":
			teams[i].ShortName = "Manchester City"

		}
	}
}
