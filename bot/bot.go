package bot

import (
	"context"
	"fmt"
	"football_tgbot/db"
	"football_tgbot/types"
	"log"
	"os"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Создаем клавиатуру с кнопками лиг для команд
var keyboardLeagues = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("APL", "league_APL"),
		tgbotapi.NewInlineKeyboardButtonData("La Liga", "league_LaLiga"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "league_Bundesliga"),
		tgbotapi.NewInlineKeyboardButtonData("Serie A", "league_SerieA"),
	),
)

// Создаем клавиатуру с кнопками лиг для таблиц
var keyboardStandings = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("APL", "standings_APL"),
		tgbotapi.NewInlineKeyboardButtonData("La Liga", "standings_LaLiga"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Bundesliga", "standings_Bundesliga"),
		tgbotapi.NewInlineKeyboardButtonData("Serie A", "standings_SerieA"),
	),
)

func Start() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_API_KEY")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_API_KEY is not set")
	}

	// Создаем новый экземпляр бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %v", err)
	}

	// Включаем режим отладки (логирование запросов)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI is not set")
	}

	client, err := connectToMongoDB(mongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Инициализируем хранилище данных
	dbName := "football"
	store := db.NewMongoDBMatchesStore(client, dbName)

	// Создаем канал для получения обновлений от Telegram
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// Обрабатываем входящие обновления
	for update := range updates {
		// Если обновление содержит сообщение
		if update.Message != nil {
			// Обрабатываем команды
			switch update.Message.Command() {
			case "start":
				handleStartCommand(bot, update.Message)
			case "help":
				handleHelpCommand(bot, update.Message)
			case "leagues":
				handleLeaguesCommand(bot, update.Message)
			case "schedule":
				handleScheduleCommand(bot, update.Message, store)
			case "standings":
				handleStandingsCommand(bot, update.Message, store)
			default:
				handleUnknownCommand(bot, update.Message)
			}
		}

		// Если обновление содержит callback (нажатие на кнопку)
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery, store)
		}
	}

	return nil
}

// connectToMongoDB подключается к MongoDB
func connectToMongoDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}

// handleScheduleCommand обрабатывает команду /schedule
func handleScheduleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, store db.MatchesStore) error {
	matches, err := store.GetMatches(context.Background(), "matches")
	if err != nil {
		return fmt.Errorf("failed to get matches: %v", err)
	}

	// Формируем ответ
	response := "Расписание матчей на ближайшие 10 дней:\n"
	if len(matches) == 0 {
		response = "На сегодня матчей не запланировано.\n"
	} else {
		for _, match := range matches {
			response += fmt.Sprintf("- %s vs %s (%s)\n", match.HomeTeam.Name, match.AwayTeam.Name, match.UTCDate[0:10])
		}
	}

	// Отправляем ответ
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func handleStandingsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, store db.MatchesStore) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу для просмотра таблицы: ")
	msg.ReplyMarkup = keyboardStandings // Используем новую клавиатуру
	bot.Send(msg)
}

// handleStartCommand обрабатывает команду /start
func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я бот для футбольной статистики. Используй /help, чтобы узнать доступные команды.")
	bot.Send(msg)
}

// handleHelpCommand обрабатывает команду /help
func handleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpText := `Доступные команды:
	/start - Начать работу с ботом
	/help - Получить список команд
	/leagues - Показать список футбольных лиг
	/schedule - Показать расписание матчей
	/standings - Показать таблицы лиг`
	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	bot.Send(msg)
}

// handleLeaguesCommand обрабатывает команду /leagues
func handleLeaguesCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	// кидает движняк с клавой
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите лигу:")
	msg.ReplyMarkup = keyboardLeagues // Используем новую клавиатуру
	bot.Send(msg)
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используй /help, чтобы узнать доступные команды.")
	bot.Send(msg)
}

// handleCallbackQuery обрабатывает нажатие на кнопку
func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, store db.MatchesStore) {
	// Получаем данные из callback (название лиги)
	league := callbackQuery.Data

	// Определяем коллекцию в MongoDB на основе выбранной лиги
	collectionName := ""
	isStandings := false
	switch league {
	case "league_APL":
		collectionName = "PremierLeague"
	case "league_LaLiga":
		collectionName = "LaLiga"
	case "league_Bundesliga":
		collectionName = "Bundesliga"
	case "league_SerieA":
		collectionName = "SerieA"
	case "standings_APL":
		collectionName = "PremierLeague_standings"
		isStandings = true
	case "standings_LaLiga":
		collectionName = "LaLiga_standings"
		isStandings = true
	case "standings_Bundesliga":
		collectionName = "Bundesliga_standings"
		isStandings = true
	case "standings_SerieA":
		collectionName = "SerieA_standings"
		isStandings = true
	default:
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Неизвестная лига.")
		bot.Send(msg)
		return
	}

	if !isStandings {
		teams, err := store.GetTeams(context.Background(), collectionName)
		if err != nil {
			log.Printf("Error getting teams: %v", err)
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Произошла ошибка при получении данных.")
			bot.Send(msg)
			return
		}

		// Формируем ответ с командами
		response := fmt.Sprintf("Команды %s:\n", collectionName)
		for _, team := range teams {
			response += fmt.Sprintf("- %s\n", team.Name)
		}

		// Отправляем ответ
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, response)
		bot.Send(msg)
	} else {
		// Обработка запроса таблицы
		standings, err := getStandingsFromDB(store, collectionName)
		if err != nil {
			log.Printf("Error getting standings: %v", err)
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Произошла ошибка при получении данных.")
			bot.Send(msg)
			return
		}

		// Формируем ответ с таблицей
		response := formatStandings(standings, collectionName)

		// Отправляем ответ
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, response)
		bot.Send(msg)
	}

	// Подтверждаем обработку callback, т.к прекращаем тот противный белый какой-то пал вокруг кнопки
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Send(callback)
}

// getStandingsFromDB получает таблицу лиги из MongoDB.
func getStandingsFromDB(store db.MatchesStore, collectionName string) ([]types.Standing, error) {
	// Получаем все команды из коллекции
	var standings []types.Standing
	var err error
	standings, err = store.GetStandings(context.Background(), collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get standings: %v", err)
	}
	return standings, nil
}

// formatStandings форматирует таблицу лиги для отправки пользователю.
func formatStandings(standings []types.Standing, collectionName string) string {
	underscore := strings.Index(collectionName, "_")
	if underscore != -1 {
		collectionName = collectionName[:underscore]
	}
	// Проверяем, что таблица не пустаs
	if len(standings) == 0 {
		return fmt.Sprintf("Таблица %s пуста.", collectionName)
	}

	// Сортируем таблицу по позиции
	sort.Slice(standings, func(i, j int) bool {
		return standings[i].Position < standings[j].Position
	})

	response := fmt.Sprintf("Таблица %s:\n", collectionName)
	response += fmt.Sprintf("%-4s %-25s %-5s %-5s %-5s %-5s %-5s %-5s %-5s %-5s\n", "#", "Команда", "И", "В", "Н", "П", "ГЗ", "ГП", "РГ", "О")
	for _, standing := range standings {
		response += fmt.Sprintf("%-4d %-25s %-5d %-5d %-5d %-5d %-5d %-5d %-5d %-5d\n",
			standing.Position, standing.Team.Name, standing.PlayedGames, standing.Won, standing.Draw, standing.Lost, standing.GoalsFor, standing.GoalsAgainst, standing.GoalDifference, standing.Points)
	}
	return response
}
