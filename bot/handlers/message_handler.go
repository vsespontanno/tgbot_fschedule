package handlers

import (
	"football_tgbot/bot/keyboards"
	"football_tgbot/db"
	"football_tgbot/rating"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleMessage обрабатывает все входящие сообщения
func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, store db.MatchesStore, ratingService *rating.Service) error {
	// Обрабатываем только текстовые сообщения
	if msg.Text == "" {
		return nil
	}

	// Обработка команд
	switch msg.Text {
	case "/start":
		return handleStart(bot, msg)
	case "/help":
		return handleHelp(bot, msg)
	case "/schedule":
		return handleScheduleCommand(bot, msg)
	case "/top":
		return HandleTopMatches(bot, msg, store, ratingService)
	case "/table":
		return handleTableCommand(bot, msg)
	default:
		return handleUnknownCommand(bot, msg)
	}
}

// handleStart обрабатывает команду /start
func handleStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	response := "Привет! Я бот для отслеживания футбольных матчей. Доступные команды:\n" +
		"/schedule - показать расписание всех матчей\n" +
		"/top - показать расписание топовых матчей\n" +
		"/table - показать турнирную таблицу\n" +
		"/help - показать справку"

	msgConfig := tgbotapi.NewMessage(msg.Chat.ID, response)
	_, err := bot.Send(msgConfig)
	return err
}

// handleHelp обрабатывает команду /help
func handleHelp(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	response := "Привет! Я бот для отслеживания футбольных матчей. Доступные команды:\n" +
		"/schedule - показать расписание всех матчей\n" +
		"/top - показать расписание топовых матчей\n" +
		"/table - показать турнирную таблицу\n" +
		"/help - показать справку"

	msgConfig := tgbotapi.NewMessage(msg.Chat.ID, response)
	_, err := bot.Send(msgConfig)
	return err
}

func handleTableCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу для просмотра турнирной таблицы:")
	msg.ReplyMarkup = keyboards.KeyboardStandings
	_, err := bot.Send(msg)
	return err
}

func handleDefaultScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу для просмотра расписания матчей:")
	msg.ReplyMarkup = keyboards.KeyboardSchedule
	_, err := bot.Send(msg)
	return err
}

func handleScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	// Создаем inline-клавиатуру с двумя опциями
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все матчи", "show_all_matches"),
			tgbotapi.NewInlineKeyboardButtonData("Топ матчи", "show_top_matches"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите тип расписания:")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	return err
}

// handleUnknownCommand обрабатывает неизвестные команды
func handleUnknownCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	response := "Неизвестная команда. Используйте /help для просмотра доступных команд."
	msgConfig := tgbotapi.NewMessage(msg.Chat.ID, response)
	_, err := bot.Send(msgConfig)
	return err
}
