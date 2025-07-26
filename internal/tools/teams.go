package tools

import "github.com/vsespontanno/tgbot_fschedule/internal/types"

// TeamsFilter фильтрует команды, заменяя длинные названия на короткие
func TeamsFilter(teams types.TeamsResponse) {
	for i := range teams.Teams {
		switch teams.Teams[i].Name {
		case "Sevilla FC":
			teams.Teams[i].ShortName = "Sevilla"
		case "Wolverhampton Wanderers FC":
			teams.Teams[i].Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			teams.Teams[i].Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			teams.Teams[i].Name = "Inter"
		case "Club Atlético de Madrid":
			teams.Teams[i].Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			teams.Teams[i].Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			teams.Teams[i].Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			teams.Teams[i].Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			teams.Teams[i].Name = "Real Sociedad"
		}
	}

	for i := range teams.Teams {
		switch teams.Teams[i].ShortName {
		case "Sevilla FC":
			teams.Teams[i].ShortName = "Sevilla"
		case "Leverkusen":
			teams.Teams[i].ShortName = "Bayer"
		case "Dortmund":
			teams.Teams[i].ShortName = "Borussia D."
		case "M'gladbach":
			teams.Teams[i].ShortName = "Borussia M."
		case "Atleti":
			teams.Teams[i].ShortName = "Atletico"
		case "Barça":
			teams.Teams[i].ShortName = "Barcelona"
		case "Leganés":
			teams.Teams[i].ShortName = "Leganes"
		case "Man United":
			teams.Teams[i].ShortName = "Manchester United"
		case "Man City":
			teams.Teams[i].ShortName = "Manchester City"

		}
	}
}
