package handlers

import (
	"context"
	"fmt"
	"football_tgbot/bot/keyboards"
	"football_tgbot/bot/models"
	"football_tgbot/db"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore) error {
	if league, ok := keyboards.Leagues[query.Data]; ok {
		return HandleLeagueCallback(bot, query, store, league)
	}

	if league, ok := keyboards.Standings[query.Data]; ok {
		return HandleStandingsCallback(bot, query, store, league)
	}

	return SendMessage(bot, query.Message.Chat.ID, "Неизвестная лига.")
}

func HandleLeagueCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore, league models.League) error {
	teams, err := store.GetTeams(context.Background(), league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting teams: %w", err)
	}

	response := fmt.Sprintf("Команды %s:\n", league.Name)
	for _, team := range teams {
		response += fmt.Sprintf("- %s\n", team.Name)
	}

	return SendMessage(bot, query.Message.Chat.ID, response)
}

func HandleStandingsCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, store db.MatchesStore, league models.League) error {
	standings, err := GetStandingsFromDB(store, league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting standings: %w", err)
	}

	imagePath := fmt.Sprintf("%s.png", league.CollectionName)
	defer os.Remove(imagePath)

	if err := GenerateTableImage(standings, imagePath); err != nil {
		return fmt.Errorf("error generating image: %w", err)
	}

	if err := SendPhoto(bot, query.Message.Chat.ID, imagePath); err != nil {
		return err
	}

	return SendCallbackResponse(bot, query.ID)
}

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

	return SendMessage(bot, query.Message.Chat.ID, response)
}
