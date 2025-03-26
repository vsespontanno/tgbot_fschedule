package handlers

import (
	"context"
	"fmt"
	"football_tgbot/bot/keyboards"
	resp "football_tgbot/bot/response"
	"football_tgbot/db"
	"football_tgbot/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// обработка сообщений и команд
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, store db.MatchesStore) error {
	switch message.Command() {
	case "start":
		return resp.SendMessage(bot, message.Chat.ID, "Привет! Я бот для футбольной статистики. Используй /help, чтобы узнать доступные команды.")
	case "help":
		return resp.SendMessage(bot, message.Chat.ID, types.HelpText)
	case "leagues":
		return resp.SendMessageWithKeyboard(bot, message.Chat.ID, "Выберите лигу:", keyboards.KeyboardLeagues)
	case "schedule":
		return HandleScheduleCommand(bot, message, store)
	case "standings":
		return resp.SendMessageWithKeyboard(bot, message.Chat.ID, "Выберите лигу для просмотра таблицы:", keyboards.KeyboardStandings)
	default:
		return resp.SendMessage(bot, message.Chat.ID, "Неизвестная команда. Используй /help, чтобы узнать доступные команды.")
	}
}

// обработка команды расписания
func HandleScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, store db.MatchesStore) error {
	matches, err := store.GetMatches(context.Background(), "matches")
	if err != nil {
		return fmt.Errorf("failed to get matches: %w", err)
	}

	msgText := "Расписание матчей на ближайшие 10 дней:\n"
	if len(matches) == 0 {
		msgText = "На сегодня матчей не запланировано.\n"
	} else {
		for _, match := range matches {
			msgText += fmt.Sprintf("- %s vs %s (%s)\n", match.HomeTeam.Name, match.AwayTeam.Name, match.UTCDate[0:10])
		}
	}

	return resp.SendMessage(bot, message.Chat.ID, msgText)
}
