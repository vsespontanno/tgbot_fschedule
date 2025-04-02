package handlers

import (
	"context"
	"fmt"
	"football_tgbot/db"
	"football_tgbot/rating"
	"football_tgbot/types"
	"sort"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleTopMatches обрабатывает запрос на получение топовых матчей
func HandleTopMatches(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, store db.MatchesStore, ratingService *rating.Service) error {
	ctx := context.Background()
	matches, err := store.GetMatches(ctx, "matches")
	if err != nil {
		return fmt.Errorf("failed to get matches: %w", err)
	}

	// Создаем структуру для хранения матчей с их рейтингами
	type MatchWithRating struct {
		Match  types.Match
		Rating float64
	}

	// Получаем рейтинги для всех матчей
	var matchesWithRatings []MatchWithRating
	for _, match := range matches {
		rating, err := ratingService.GetMatchRating(ctx, "team_ratings", match)
		if err != nil {
			continue // Пропускаем матчи с ошибками при получении рейтинга
		}
		matchesWithRatings = append(matchesWithRatings, MatchWithRating{
			Match:  match,
			Rating: rating,
		})
	}

	// Сортируем матчи по рейтингу (по убыванию)
	sort.Slice(matchesWithRatings, func(i, j int) bool {
		return matchesWithRatings[i].Rating > matchesWithRatings[j].Rating
	})

	// Формируем сообщение с топ-5 матчами
	response := "🏆 Топ матчи:\n\n"
	for i, mwr := range matchesWithRatings {
		if i >= 5 {
			break
		}
		match := mwr.Match
		ratingLevel := getRatingLevel(mwr.Rating)
		response += fmt.Sprintf("%d. %s vs %s\n📅 %s\n⭐ Рейтинг: %s\n\n",
			i+1,
			match.HomeTeam.Name,
			match.AwayTeam.Name,
			formatDate(match.UTCDate),
			ratingLevel,
		)
	}

	msgConfig := tgbotapi.NewMessage(msg.Chat.ID, response)
	_, err = bot.Send(msgConfig)
	return err
}

// getRatingLevel возвращает текстовое описание рейтинга
func getRatingLevel(rating float64) string {
	switch {
	case rating >= 0.8:
		return "⭐⭐⭐⭐⭐"
	case rating >= 0.6:
		return "⭐⭐⭐⭐"
	case rating >= 0.4:
		return "⭐⭐⭐"
	case rating >= 0.2:
		return "⭐⭐"
	default:
		return "⭐"
	}
}

// formatDate форматирует дату матча
func formatDate(utcDate string) string {
	t, err := time.Parse(time.RFC3339, utcDate)
	if err != nil {
		return utcDate
	}

	// Форматируем дату в удобный вид
	return t.Format("02.01.2006 15:04")
}
