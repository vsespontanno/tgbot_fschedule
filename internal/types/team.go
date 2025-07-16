package types

// Cтруктура для хранения информации о команде
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Tla       string `json:"tla"`
	Founded   int    `json:"founded"`
	League    string `json:"league"`
}

// Cтруктура для декодинга Json-файла из API
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

// Cтруктура для хранения статистики команды для подсчёта рейтинга матча
type TeamForm struct {
	Wins   int `bson:"wins" json:"wins"`
	Losses int `bson:"losses" json:"losses"`
	Draws  int `bson:"draws" json:"draws"`
}
