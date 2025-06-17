package handlers

import (
	"context"
	"fmt"
	"football_tgbot/internal/bot/keyboards"
	resp "football_tgbot/internal/bot/response"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
	"os"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// обработка callback запросов таких как выбор лиги, выбор команды, выбор таблицы через кнопки
func HandleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, matchService *service.MatchesService, standingsService *service.StandingsService) error {
	// Отправляем пустой ответ, чтобы убрать "часики" у кнопки
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := bot.Request(callback); err != nil {
		return err
	}

	switch query.Data {
	case "show_top_matches":
		return HandleTopMatches(bot, query, matchService)
	case "show_all_matches":
		return HandleDefaultScheduleCommand(bot, query.Message)
	}

	if league, ok := keyboards.KeyboardsStandings[query.Data]; ok {
		return HandleStandingsCallback(bot, query, standingsService, league)
	}

	if league, ok := keyboards.KeyboardsSchedule[query.Data]; ok {
		return HandleScheduleCallback(bot, query, matchService, league)
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, "Неизвестная команда.")
}

// обработка таблицы и вывод изображения
func HandleStandingsCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, standingsService *service.StandingsService, league types.League) error {
	standings, err := standingsService.HandleGetStandings(context.Background(), league.CollectionName)
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

// TODO: переписать как в standings
// HandleScheduleCallback обрабатывает callback запросы на расписание матчей
func HandleScheduleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, service *service.MatchesService, league types.League) error {
	// Отвечаем на callback запрос
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	leagueCode := strings.TrimPrefix(callback.Data, "schedule_")
	leagueName := getLeagueName(leagueCode)

	matches, err := service.HandleGetMatches(context.Background())
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при получении расписания матчей")
		bot.Send(msg)
		if _, err := bot.Request(callbackConfig); err != nil {
			return err
		}
		return err
	}

	// Фильтруем матчи только по лиге
	var leagueMatches []types.Match
	for _, match := range matches {
		if match.Competition.Name == leagueName {
			leagueMatches = append(leagueMatches, match)
		}
	}

	if len(leagueMatches) == 0 {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("В %s матчей не запланировано", leagueName))
		bot.Send(msg)
		if _, err := bot.Request(callbackConfig); err != nil {
			return err
		}
		return nil
	}

	// Генерируем изображение с расписанием
	buf, err := GenerateScheduleImage(leagueMatches)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при создании изображения с расписанием")
		bot.Send(msg)
		if _, err := bot.Request(callbackConfig); err != nil {
			return err
		}
		return err
	}

	// Отправляем изображение
	photo := tgbotapi.FileBytes{
		Name:  "schedule.png",
		Bytes: buf.Bytes(),
	}
	msg := tgbotapi.NewPhoto(callback.Message.Chat.ID, photo)
	_, err = bot.Send(msg)

	if _, err := bot.Request(callbackConfig); err != nil {
		return err
	}

	return err
}

// getLeagueName возвращает полное название лиги по её коду
func getLeagueName(code string) string {
	switch code {
	case "laliga":
		return "La Liga"
	case "epl":
		return "EPL"
	case "primeira":
		return "Primeira"
	case "eredivisie":
		return "Eredivisie"
	case "bundesliga":
		return "Bundesliga"
	case "seriea":
		return "Serie A"
	case "ucl":
		return "UCL"
	case "uel":
		return "UEL"
	default:
		return code
	}
}

func HandleTopMatches(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, service *service.MatchesService) error {
	ctx := context.Background()
	matches, err := service.HandleGetMatches(ctx)
	if err != nil {
		return err
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Rating > matches[j].Rating
	})

	if len(matches) > 13 {
		matches = matches[:13]
	}

	if len(matches) == 0 {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "На ближайшие дни нет топовых матчей.")
		bot.Send(msg)
		return nil
	}

	buf, err := GenerateScheduleImage(matches)
	if err != nil {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Произошла ошибка при создании изображения с расписанием топовых матчей")
		bot.Send(msg)
		return err
	}

	photo := tgbotapi.FileBytes{
		Name:  "top_matches_schedule.png",
		Bytes: buf.Bytes(),
	}
	msg := tgbotapi.NewPhoto(callback.Message.Chat.ID, photo)
	_, err = bot.Send(msg)

	return err

}
