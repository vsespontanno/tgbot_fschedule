package handlers

import (
	"context"
	"log"

	"github.com/vsespontanno/tgbot_fschedule/internal/bot/keyboards"
	resp "github.com/vsespontanno/tgbot_fschedule/internal/bot/response"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обрабатывает все входящие сообщения
func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, userService *service.UserService) error {
	if msg.Text == "" {
		return nil
	}
	switch msg.Text {
	case "/start":
		return handleStart(bot, msg, userService)
	case "/help":
		return handleHelp(bot, msg)
	case "/schedule":
		return handleScheduleCommand(bot, msg)
	case "/table":
		return handleTableCommand(bot, msg)
	default:
		return handleUnknownCommand(bot, msg)
	}
}

// Обрабатывает команду /start
func handleStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, userService *service.UserService) error {
	ctx := context.Background()
	user := &types.User{
		TelegramID: msg.Chat.ID,
		Username:   msg.Chat.UserName,
	}

	err := userService.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error saving user: %v", err)
		return err
	}
	response := "Привет! Я бот для отслеживания футбольных матчей. Доступные команды:\n" +
		"/schedule - показать расписание матчей\n" +
		"/table - показать турнирную таблицу\n" +
		"/help - показать справку"

	return resp.SendMessage(bot, msg.Chat.ID, response)
}

// Обрабатывает команду /help
func handleHelp(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	response := "Привет! Я бот для отслеживания футбольных матчей. Доступные команды:\n" +
		"/schedule - показать расписание всех матчей\n" +
		"/table - показать турнирную таблицу\n" +
		"/help - показать справку"
	return resp.SendMessage(bot, msg.Chat.ID, response)
}

// Обрабатывает команду /table
func handleTableCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	text := "Выберите лигу для просмотра турнирной таблицы:"
	return resp.SendMessageWithKeyboard(bot, message.Chat.ID, text, keyboards.KeyboardStandings)
}

// Обрабатывает команду /schedule
func handleScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	text := "Выберите тип расписания:"
	return resp.SendMessageWithKeyboard(bot, message.Chat.ID, text, keyboards.Keyboard_Schedule)
}

// Обрабатывает неизвестные команды
// Отправляет сообщение о том, что команда не распознана
func handleUnknownCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	response := "Неизвестная команда. Используйте /help для просмотра доступных команд."
	return resp.SendMessage(bot, msg.Chat.ID, response)
}

// Обрабатывает нажатие на кнопку для получения расписания матчей в лиге на 7 дней.
// Отправляет кнопки для выбора лиги и получает расписание матчей.
func HandleDefaultScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	msg := "Выберите лигу для просмотра расписания:"
	return resp.SendMessageWithKeyboard(bot, message.Chat.ID, msg, keyboards.KeyboardDefaultSchedule)
}
