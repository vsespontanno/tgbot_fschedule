package types

var (
	CLstage = map[string]float64{
		"PLAYOFFS":       0.25, // 1/16
		"LAST_16":        0.5,  // 1/8
		"QUARTER_FINALS": 0.75, // 1/4
		"SEMI_FINALS":    0.9,  // 1/2
		"FINAL":          1.0,  // Final
	}

	leagueNorm = map[string]float64{
		"Champions League": 1.0,
		"PremierLeague":    0.9,
		"LaLiga":           0.8,
		"SerieA":           0.8,
		"Bundesliga":       0.75,
		"Ligue1":           0.7,
	}

	teamsInLeague = map[string]int{
		"PremierLeague": 20,
		"LaLiga":        20,
		"Bundesliga":    18,
		"SerieA":        18,
		"Ligue1":        20,
	}

	derbys = map[[2]string]float64{
		// Англия (PremierLeague)
		{"Manchester United", "Manchester City"}: 0.27, // Манчестерское дерби
		{"Liverpool", "Everton"}:                 0.16, // Мерсисайдское дерби
		{"Arsenal", "Tottenham"}:                 0.25, // Северолондонское дерби
		{"Chelsea", "Arsenal"}:                   0.25,
		{"Chelsea", "Tottenham"}:                 0.25,
		{"Manchester United", "Liverpool"}:       0.26,
		{"Manchester United", "Leeds United"}:    0.15,
		{"Newcastle", "Sunderland"}:              0.14, // Тайн-Уир

		// Испания (LaLiga)
		{"Real Madrid", "Barcelona"}: 0.35, // Увеличено с 0.3 для Эль Класико
		{"Atletico", "Real Madrid"}:  0.26, // Мадридское дерби
		{"Sevilla", "Real Betis"}:    0.2,  // Севильское дерби
		{"Barcelona", "Espanyol"}:    0.18, // Барселонское дерби
		{"Valencia", "Levante"}:      0.14, // Валенсийское дерби

		// Германия (Bundesliga)
		{"Borussia D.", "Bayern"}:           0.28, // Дер Классикер
		{"Schalke 04", "Borussia Dortmund"}: 0.16, // Рурское дерби
		{"Hamburger SV", "Werder Bremen"}:   0.15, // Северное дерби
		{"Bayern", "1860 Munich"}:           0.14, // Мюнхенское дерби
		{"Cologne", "Borussia M."}:          0.14,

		// Италия (SerieA)
		{"Inter", "Milan"}:     0.29, // Миланское дерби
		{"Roma", "Lazio"}:      0.28, // Римское дерби
		{"Juventus", "Torino"}: 0.2,  // Дерби делла Моле
		{"Genoa", "Sampdoria"}: 0.18, // Дерби делла Лантерна
		{"Napoli", "Roma"}:     0.15,

		// Франция (Ligue1)
		{"PSG", "Marseille"}:                0.23, // Ле Классик
		{"Olympique Lyon", "Saint-Etienne"}: 0.18, // Ронское дерби
		{"Nice", "Monaco"}:                  0.14, // Лазурное дерби
		{"Lille", "RC Lens"}:                0.14, // Северное дерби
	}
)
