package handlers

import (
	"football_tgbot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot   *tgbotapi.BotAPI
	store db.MatchesStore
}

func NewHandler(bot *tgbotapi.BotAPI, store db.MatchesStore) *Handler {
	return &Handler{
		bot:   bot,
		store: store,
	}
}
