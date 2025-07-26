package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/bot/keyboards"
	resp "github.com/vsespontanno/tgbot_fschedule/internal/bot/response"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка callback запросов таких как выбор лиги, выбор команды, выбор таблицы через кнопки
func HandleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, matchService *service.MatchesService, standingsService *service.StandingsService, redisClient *cache.RedisClient) error {
	switch query.Data {
	case "show_top_matches":
		resp.SendCallbackResponse(bot, query.ID)
		return HandleTopMatches(bot, query, matchService, redisClient, "")
	case "show_all_matches":
		resp.SendCallbackResponse(bot, query.ID)
		return HandleDefaultScheduleCommand(bot, query.Message)

	}

	if league, ok := keyboards.KeyboardsStandings[query.Data]; ok {
		resp.SendCallbackResponse(bot, query.ID)
		return HandleStandingsCallback(bot, query, standingsService, redisClient, league)
	}

	if league, ok := keyboards.KeyboardsSchedule[query.Data]; ok {
		resp.SendCallbackResponse(bot, query.ID)
		return HandleScheduleCallback(bot, query, matchService, redisClient, league, query.Data)
	}

	return resp.SendMessage(bot, query.Message.Chat.ID, "Неизвестная команда.")
}

// Обработка колбэков для турнирной таблицы и расписания матчей
// Здесь мы получаем таблицу для выбранной лиги и отправляем ее пользователю в виде изображения
func HandleStandingsCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, standingsService *service.StandingsService, redisClient *cache.RedisClient, league types.League) error {
	standings, err := standingsService.HandleGetStandings(context.Background(), league.CollectionName)
	if err != nil {
		return fmt.Errorf("error getting standings: %w", err)
	}
	imagePath := fmt.Sprintf("%s.png", league.CollectionName)
	defer os.Remove(imagePath)

	if err := GenerateTableImage(standings, league.Code, imagePath, redisClient); err != nil {
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

// Обработка колбэков для расписания матчей
// Здесь мы получаем расписание матчей для выбранной лиги и отправляем его пользователю в виде изображения
func HandleScheduleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, service *service.MatchesService, redisClient *cache.RedisClient, league types.League, button string) error {
	var (
		leagueName = strings.TrimPrefix(query.Data, "schedule_")
		imagePath  = fmt.Sprintf("%s.png", leagueName)
		cacheKey   = "all_matches_image" + imagePath
		from       = time.Now()
		to         = from.AddDate(0, 0, 7)
		ctx        = context.Background()
	)

	// Проверяем кэш
	if _, err := redisClient.GetBytes(ctx, cacheKey); err == nil {
		logrus.WithField("cache_key", cacheKey).Info("Cache hit for schedule image")
		err = resp.SendPhoto(bot, query.Message.Chat.ID, imagePath)
		resp.SendCallbackResponse(bot, query.ID)
		if err != nil {
			resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с расписанием")
			return err
		}
		return nil

	} else if errors.Is(err, redis.Nil) {
		logrus.WithField("cache_key", cacheKey).Warn("Cache error: ", err)
	}

	matches, err := service.HandleGetMatchesForPeriod(context.Background(), leagueName, from.Format("2006-01-02"), to.Format("2006-01-02"))

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

	// Генерируем изображение с расписанием
	if err := GenerateScheduleImage(leagueMatches, imagePath, redisClient); err != nil {
		resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с матчами: \n")
		resp.SendCallbackResponse(bot, query.ID)
		return err
	}

	return err
}

// Обработка команды для получения расписания топовых матчей
// Здесь мы получаем топовые матчи за неделю и отправляем их пользователю в виде
// изображения.
func HandleTopMatches(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, service *service.MatchesService, redisClient *cache.RedisClient, button string) error {
	var (
		ctx       = context.Background()
		cacheKey  = "top_matches_image"
		imagePath = fmt.Sprintf("%s.png", "top_matches")
		from      = time.Now()
		to        = from.AddDate(0, 0, 7)
	)

	// Проверяем кэш
	if _, err := redisClient.GetBytes(ctx, cacheKey); err == nil {
		logrus.WithField("cache_key", cacheKey).Info("Cache hit for schedule image")
		err = resp.SendPhoto(bot, query.Message.Chat.ID, imagePath)
		if err != nil {
			resp.SendMessage(bot, query.Message.Chat.ID, "Произошла ошибка при отправке изображения с расписанием")
			return err
		}
		return resp.SendCallbackResponse(bot, query.ID)

	} else if errors.Is(err, redis.Nil) {
		logrus.WithField("cache_key", cacheKey).Warn("Cache error: ", err)
	}

	matches, err := service.HandleGetMatchesForPeriod(ctx, "", from.Format("2006-01-02"), to.Format("2006-01-02"))
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		err = resp.SendMessage(bot, query.Message.Chat.ID, "На ближайшие дни нет топовых матчей.")
		if err != nil {
			return err
		}

		return nil
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Rating > matches[j].Rating
	})

	if len(matches) > 13 {
		matches = matches[:13]
	}

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
