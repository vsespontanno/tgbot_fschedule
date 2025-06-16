package types

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

// структура для хранения информации о командах
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}
