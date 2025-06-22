package types

// структура для хранения информации о команде
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Tla       string `json:"tla"`
	Founded   int    `json:"founded"`
}

// структура для хранения информации о командах
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type TeamForm struct {
	Wins   int `bson:"wins" json:"wins"`
	Losses int `bson:"losses" json:"losses"`
	Draws  int `bson:"draws" json:"draws"`
}
