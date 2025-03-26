package handlers

import (
	"context"
	"fmt"
	"football_tgbot/bot/keyboards"
	resp "football_tgbot/bot/response"
	"football_tgbot/db"
	"football_tgbot/types"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// обработка callback запросов таких как выбор лиги, выбор команды, выбор таблицы через кнопки
func HandleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore) error {
	if league, ok := keyboards.KeyboardsLeagues[query.Data]; ok {
		return HandleLeagueCallback(bot, query, store, league)
	}

	if league, ok := keyboards.KeyboardsStandings[query.Data]; ok {
		return HandleStandingsCallback(bot, query, store, league)
	}

	if query.Data == "schedule" {
		return HandleScheduleCallback(bot, query, store)
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, "Неизвестная команда.")
}

// обработка лиги по кнопке и вывод команд принадлежащих этой лиге
func HandleLeagueCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore, league types.League) error {
	teams, err := store.GetTeams(context.Background(), league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting teams: %w", err)
	}

	response := fmt.Sprintf("Команды %s:\n", league.Name)
	for _, team := range teams {
		response += fmt.Sprintf("- %s\n", team.Name)
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, response)
}

// обработка таблицы и вывод изображения
func HandleStandingsCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore, league types.League) error {
	standings, err := store.GetStandings(context.Background(), league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting standings: %w", err)
	}

	fmt.Printf("Retrieved %d standings for %s\n", len(standings), league.CollectionName)
	for _, s := range standings {
		fmt.Printf("Position: %d, Team: %s, Points: %d\n", s.Position, s.Team.Name, s.Points)
	}

	imagePath := fmt.Sprintf("%s.png", league.CollectionName)
	defer os.Remove(imagePath)

	if err := GenerateTableImage(standings, imagePath); err != nil {
		return fmt.Errorf("error generating image: %w", err)
	}

	if err := resp.SendPhoto(bot, query.Message.Chat.ID, imagePath); err != nil {
		return err
	}

	// Отправляем ответ на callback query
	callback := tgbotapi.NewCallback(query.ID, "")
	_, err = bot.Request(callback)
	return err
}

// обработка расписания и вывод матчей
func HandleScheduleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore) error {
	matches, err := store.GetMatches(context.Background(), "matches")
	if err != nil {
		return fmt.Errorf("error getting matches: %w", err)
	}

	response := "Расписание матчей на ближайшие 10 дней:\n"
	if len(matches) == 0 {
		response = "На сегодня матчей не запланировано.\n"
	} else {
		for _, match := range matches {
			response += fmt.Sprintf("- %s vs %s (%s)\n", match.HomeTeam.Name, match.AwayTeam.Name, match.UTCDate[0:10])
		}
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, response)
}
