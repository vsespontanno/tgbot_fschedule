package handlers

import (
	"football_tgbot/bot/keyboards"
	"football_tgbot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleMessage обрабатывает все входящие сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, store db.MatchesStore) error {
	if !message.IsCommand() {
		return nil
	}

	switch message.Command() {
	case "start":
		return handleStartCommand(bot, message)
	case "help":
		return handleHelpCommand(bot, message)
	case "table":
		return handleTableCommand(bot, message)
	case "schedule":
		return handleScheduleCommand(bot, message)
	default:
		return nil
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	text := `Привет! Я бот для просмотра футбольных матчей и турнирных таблиц.

	Доступные команды:
	/schedule - Посмотреть расписание матчей
	/table - Посмотреть турнирную таблицу
	/help - Показать справку`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := bot.Send(msg)
	return err
}

func handleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	text := `Список доступных команд:

	/schedule - Показывает расписание матчей на сегодня для выбранной лиги
	/table - Показывает турнирную таблицу для выбранной лиги
	/help - Показывает это сообщение

	При выборе команды /schedule или /table вам будет предложено выбрать лигу из списка.`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := bot.Send(msg)
	return err
}

func handleTableCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу для просмотра турнирной таблицы:")
	msg.ReplyMarkup = keyboards.KeyboardStandings
	_, err := bot.Send(msg)
	return err
}

func handleScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу для просмотра расписания матчей:")
	msg.ReplyMarkup = keyboards.KeyboardSchedule
	_, err := bot.Send(msg)
	return err
}
