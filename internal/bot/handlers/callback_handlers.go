package handlers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vsespontanno/tgbot_fschedule/internal/bot/keyboards"
	resp "github.com/vsespontanno/tgbot_fschedule/internal/bot/response"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// обработка callback запросов таких как выбор лиги, выбор команды, выбор таблицы через кнопки
func HandleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, matchService *service.MatchesService, standingsService *service.StandingsService, redisClient *cache.RedisClient) error {
	// Отправляем пустой ответ, чтобы убрать "часики" у кнопки
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := bot.Request(callback); err != nil {
		return err
	}

	switch query.Data {
	case "show_top_matches":
		return HandleTopMatches(bot, query, matchService, redisClient, "")
	case "show_all_matches":
		return HandleDefaultScheduleCommand(bot, query.Message)
	}

	if league, ok := keyboards.KeyboardsStandings[query.Data]; ok {
		return HandleStandingsCallback(bot, query, standingsService, redisClient, league)
	}

	if league, ok := keyboards.KeyboardsSchedule[query.Data]; ok {
		return HandleScheduleCallback(bot, query, matchService, redisClient, league, query.Data)
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, "Неизвестная команда.")
}

// обработка таблицы и вывод изображения
// todo: если нет стендингов, отправить что их нет.
func HandleStandingsCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, standingsService *service.StandingsService, redisClient *cache.RedisClient, league types.League) error {
	standings, err := standingsService.HandleGetStandings(context.Background(), league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting standings: %w", err)
	}
	imagePath := fmt.Sprintf("%s.png", league.CollectionName)
	defer os.Remove(imagePath)

	if err := GenerateTableImage(standings, imagePath, redisClient); err != nil {
		return fmt.Errorf("error generating image: %w", err)
	}

	err = resp.SendPhoto(bot, query.Message.Chat.ID, imagePath)
	if err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с таблицей")
		return fmt.Errorf("error sending image for table: %w", err)
	}

	// Отправляем ответ на callback query
	resp.SendCallbackResponse(bot, query.ID)
	return err
}

func HandleScheduleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, service *service.MatchesService, redisClient *cache.RedisClient, league types.League, button string) error {
	// Отвечаем на callback запрос
	leagueCode := strings.TrimPrefix(query.Data, "schedule_")
	leagueName := getLeagueName(leagueCode)

	matches, err := service.HandleGetMatches(context.Background())
	if err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при получении расписания матчей")
		resp.SendCallbackResponse(bot, query.ID)

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
		resp.SendMessage(bot, query.Message.Chat.ID, "На ближайшие дни нет матчей в лиге "+leagueName+".")
		resp.SendCallbackResponse(bot, query.ID)

		return nil
	}
	imagePath := fmt.Sprintf("%s.png", leagueName)

	// Генерируем изображение с расписанием
	if err := GenerateScheduleImage(leagueMatches, imagePath, redisClient); err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с матчами: \n")
		resp.SendCallbackResponse(bot, query.ID)

		return err
	}

	err = resp.SendPhoto(bot, query.Message.Chat.ID, imagePath)
	if err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с расписанием")
		return err
	}
	resp.SendCallbackResponse(bot, query.ID)

	return err
}

// getLeagueName возвращает полное название лиги по её коду
func getLeagueName(code string) string {
	switch code {
	case "laliga":
		return "LaLiga"
	case "epl":
		return "EPL"
	case "primeira":
		return "Primeira"
	case "eredivisie":
		return "Eredivisie"
	case "bundesliga":
		return "Bundesliga"
	case "seriea":
		return "SerieA"
	case "ucl":
		return "UCL"
	case "uel":
		return "UEL"
	default:
		return code
	}
}

func HandleTopMatches(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, service *service.MatchesService, redisClient *cache.RedisClient, button string) error {

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
		err = resp.SendMessage(bot, query.Message.Chat.ID, "На ближайшие дни нет топовых матчей.")
		if err != nil {
			return err
		}

		return nil
	}
	imagePath := fmt.Sprintf("%s.png", "top_matches")

	if err := GenerateScheduleImage(matches, imagePath, redisClient); err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при создании изображения с топ-матчами")
		return err
	}

	err = resp.SendPhoto(bot, query.Message.Chat.ID, imagePath)
	if err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с таблицей")
		return fmt.Errorf("error sending image for table: %w", err)
	}

	resp.SendCallbackResponse(bot, query.ID)

	return err

}
