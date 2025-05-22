# Football Schedule Telegram Bot

Телеграм-бот для просмотра расписания футбольных матчей, турнирных таблиц и информации о командах в популярных лигах.

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
- Расчет и отображение рейтингов команд и матчей

## Требования

- Go 1.23.5 или выше
- MongoDB
- (опционально) PostgreSQL для хранения пользователей и миграций
- API-ключ от [football-data.org](https://www.football-data.org/)
- API-ключ Telegram Bot

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

3. Создайте файл `.env` в корне проекта:
   ```env
   TELEGRAM_BOT_API_KEY=your_telegram_bot_token
   FOOTBALL_DATA_API_KEY=your_football_data_api_key
   MONGODB_URI=your_mongodb_connection_string
   POSTGRES_URI=your_postgres_connection_string # если используете PostgreSQL
   ```

## Инициализация базы данных

Для заполнения MongoDB данными выполните скрипты (по порядку):

```bash
make drop         # Очистка коллекций
make seedteams    # Заполнение командами
make seedstandings # Заполнение турнирных таблиц
make seedmatches  # Заполнение расписания матчей
```

## Миграции PostgreSQL

Если вы используете PostgreSQL для хранения пользователей или другой информации:

1. Установите [goose](https://github.com/pressly/goose):
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Примените миграции:
   ```bash
   make migrate-up
   ```

## Запуск бота

```bash
make run
```

## Структура проекта

```
tgbot_fschedule/
├── cmd/
│   └── bot/                # Точка входа (main.go)
├── internal/
│   ├── bot/                # Логика Telegram-бота
│   │   ├── handlers/       # Обработчики сообщений и callback-запросов
│   │   ├── keyboards/      # Клавиатуры и кнопки
│   │   └── response/       # Формирование ответов
│   ├── db/                 # Работа с базой данных
│   ├── rating/             # Расчет рейтингов
│   ├── types/              # Типы данных (Team, Match, User и др.)
│   └── scripts/            # Скрипты для инициализации данных
│       ├── drop_coll/
│       ├── seed_teams/
│       ├── seed_standings/
│       └── seed_matches/
├── migrations/             # Миграции для PostgreSQL
├── Makefile
├── go.mod
├── .env
└── README.md
```

## Использование

После запуска бота отправьте команду `/start` в Telegram. Доступные возможности:

- Просмотр списка команд в выбранной лиге
- Просмотр турнирной таблицы с визуализацией
- Просмотр расписания матчей на ближайшие дни
- Просмотр рейтингов команд и топовых матчей

## Лицензия

MIT
