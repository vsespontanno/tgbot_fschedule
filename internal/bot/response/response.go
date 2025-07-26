package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция для отправки ответа на callback запрос
func SendCallbackResponse(bot *tgbotapi.BotAPI, queryID string) error {
	callback := tgbotapi.NewCallback(queryID, "")
	_, err := bot.Send(callback)
	return err
}

// Функция для отправки сообщения
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	return err
}

// Функция для отправки сообщения с клавиатурой
func SendMessageWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	return err
}

// Функция для отправки фото
func SendPhoto(bot *tgbotapi.BotAPI, chatID int64, path string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(path))
	_, err := bot.Send(photo)
	return err
}
