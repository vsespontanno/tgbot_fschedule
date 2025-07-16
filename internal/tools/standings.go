package tools

import "github.com/vsespontanno/tgbot_fschedule/internal/types"

func StandingsFilter(standings []types.Standing) {
	for i := range standings {
		switch standings[i].Team.Name {
		case "Wolverhampton Wanderers FC":
			standings[i].Team.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			standings[i].Team.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			standings[i].Team.Name = "Inter"
		case "Club Atlético de Madrid":
			standings[i].Team.Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			standings[i].Team.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			standings[i].Team.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			standings[i].Team.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			standings[i].Team.Name = "Real Sociedad"
		}
	}

	for i := range standings {
		switch standings[i].Team.ShortName {
		case "Leverkusen":
			standings[i].Team.ShortName = "Bayer"
		case "Dortmund":
			standings[i].Team.ShortName = "Borussia D."
		case "M'gladbach":
			standings[i].Team.ShortName = "Borussia M."
		case "Atleti":
			standings[i].Team.ShortName = "Atletico"
		case "Barça":
			standings[i].Team.ShortName = "Barcelona"
		case "Leganés":
			standings[i].Team.ShortName = "Leganes"
		case "Man United":
			standings[i].Team.ShortName = "Manchester United"
		case "Man City":
			standings[i].Team.ShortName = "Manchester City"

		}
	}

}
