package types

// структура для хранения информации о матче
type Match struct {
	ID          int `json:"id"`
	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"competition"`
	HomeTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homeTeam"`
	AwayTeam struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"awayTeam"`
	UTCDate string `json:"utcDate"`
	Status  string `json:"status"`
	Score   struct {
		FullTime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fullTime"`
	} `json:"score"`
}

// структура для хранения информации о командах
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

// структура для хранения информации о матчах
type MatchesResponse struct {
	Matches []Match `json:"matches"`
}

// структура для хранения информации о площадке
type Area struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// структура для хранения информации о команде
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Tla       string `json:"tla"`
	CrestURL  string `json:"crestUrl"`
	Area      Area   `json:"area"`
	Founded   int    `json:"founded"`
}

// структура для хранения информации о таблице
type Standing struct {
	Position       int  `json:"position"`
	Team           Team `json:"team"`
	PlayedGames    int  `json:"playedGames"`
	Won            int  `json:"won"`
	Draw           int  `json:"draw"`
	Lost           int  `json:"lost"`
	Points         int  `json:"points"`
	GoalsFor       int  `json:"goalsFor"`
	GoalsAgainst   int  `json:"goalsAgainst"`
	GoalDifference int  `json:"goalDifference"`
}

// структура для хранения информации о таблице
type StandingsResponse struct {
	Standings []struct {
		Table []Standing `json:"table"`
	} `json:"standings"`
}

// структура для хранения рейтинга команды
type TeamRating struct {
	TeamID           int     `json:"teamId" bson:"teamId"`
	TeamName         string  `json:"teamName" bson:"teamName"`
	Position         int     `json:"position" bson:"position"`
	Points           int     `json:"points" bson:"points"`
	Form             float64 `json:"form" bson:"form"`                         // Форма команды (0-1)
	GoalDiff         int     `json:"goalDiff" bson:"goalDiff"`                 // Разница забитых и пропущенных
	TournamentWeight float64 `json:"tournamentWeight" bson:"tournamentWeight"` // Вес турнира (0-1)
	LastUpdated      string  `json:"lastUpdated" bson:"lastUpdated"`
}

// структура для хранения рейтингов команд
type TeamRatingsResponse struct {
	Ratings []TeamRating `json:"ratings"`
}

// CalculateRating вычисляет общий рейтинг команды
func (tr *TeamRating) CalculateRating() float64 {
	// Базовый рейтинг на основе позиции (чем выше позиция, тем выше рейтинг)
	positionScore := 1.0 - (float64(tr.Position) / 20.0) // Предполагаем максимум 20 команд

	// Рейтинг на основе очков (нормализованный)
	pointsScore := float64(tr.Points) / 100.0 // Предполагаем максимум 100 очков

	// Рейтинг на основе формы
	formScore := tr.Form

	// Рейтинг на основе разницы мячей (нормализованный)
	goalDiffScore := float64(tr.GoalDiff) / 50.0 // Предполагаем максимум разницы в 50 мячей

	// Веса для разных компонентов
	weights := struct {
		position   float64
		points     float64
		form       float64
		goalDiff   float64
		tournament float64
	}{
		position:   0.3,
		points:     0.3,
		form:       0.2,
		goalDiff:   0.1,
		tournament: 0.1,
	}

	// Вычисляем общий рейтинг
	totalRating := positionScore*weights.position +
		pointsScore*weights.points +
		formScore*weights.form +
		goalDiffScore*weights.goalDiff +
		tr.TournamentWeight*weights.tournament

	// Нормализуем рейтинг до диапазона 0-1
	if totalRating > 1.0 {
		totalRating = 1.0
	} else if totalRating < 0.0 {
		totalRating = 0.0
	}

	return totalRating
}

// GetRatingLevel возвращает уровень рейтинга команды
func (tr *TeamRating) GetRatingLevel() string {
	rating := tr.CalculateRating()
	switch {
	case rating >= 0.8:
		return "Высокий"
	case rating >= 0.5:
		return "Средний"
	default:
		return "Низкий"
	}
}
