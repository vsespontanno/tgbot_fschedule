# Football Schedule Telegram Bot

Telegram бот для просмотра расписания матчей, турнирных таблиц и информации о командах в различных футбольных лигах.

## Поддерживаемые лиги

- Premier League (Англия)
- La Liga (Испания)
- Bundesliga (Германия)
- Serie A (Италия)
- Ligue 1 (Франция)
- Champions League (Лига чемпионов)

## Функциональность

- Просмотр списка команд в каждой лиге
- Просмотр турнирных таблиц с визуализацией
- Просмотр расписания матчей на ближайшие дни
- Автоматическое обновление данных через API football-data.org

## Требования

- Go 1.23.5 или выше
- MongoDB
- API ключ от [football-data.org](https://www.football-data.org/)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/tgbot_fschedule.git
cd tgbot_fschedule
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` в корневой директории проекта со следующим содержимым:
```env
TELEGRAM_BOT_API_KEY=your_telegram_bot_token
FOOTBALL_DATA_API_KEY=your_football_data_api_key
MONGODB_URI=your_mongodb_connection_string
```

## Инициализация базы данных

Для заполнения базы данных данными выполните следующие скрипты в указанном порядке:

```bash
# Очистка существующих коллекций
go run scripts/drop_coll/drop_coll.go

# Заполнение данными о командах
go run scripts/seed_teams/seed_teams.go

# Заполнение турнирных таблиц
go run scripts/seed_standings/seed_standings.go

# Заполнение расписания матчей
go run scripts/seed_matches/seed_matches.go
```

## Запуск бота

```bash
go run main.go
```

## Структура проекта

```
tgbot_fschedule/
├── bot/
│   ├── handlers/      # Обработчики сообщений и callback-запросов
│   ├── keyboards/     # Определение клавиатур и кнопок
│   └── response/      # Функции для отправки ответов
├── db/               # Работа с базой данных
├── scripts/          # Скрипты для инициализации данных
│   ├── drop_coll/    # Очистка коллекций
│   ├── seed_teams/   # Заполнение данными о командах
│   ├── seed_standings/ # Заполнение турнирных таблиц
│   └── seed_matches/  # Заполнение расписания матчей
├── types/            # Определения типов данных
└── main.go          # Точка входа в приложение
```

## Использование

После запуска бота, отправьте команду `/start` в Telegram для начала работы. Бот предоставит вам следующие возможности:

- Просмотр списка команд в выбранной лиге
- Просмотр турнирной таблицы с визуализацией
- Просмотр расписания матчей на ближайшие дни

## Лицензия

MIT
