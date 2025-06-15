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
   ```

## Инициализация базы данных

Для заполнения MongoDB данными выполните скрипты (по порядку):

```bash
make drop         # Очистка коллекций
make seedteams    # Заполнение командами
make seedstandings # Заполнение турнирных таблиц
make seedmatches  # Заполнение расписания матчей
```
## Запуск бота

```bash
make run
```

## Структура проекта

## Использование

После запуска бота отправьте команду `/start` в Telegram. Доступные возможности:

- Просмотр списка команд в выбранной лиге
- Просмотр турнирной таблицы с визуализацией
- Просмотр расписания матчей на ближайшие дни
- Просмотр рейтингов команд и топовых матчей

## Лицензия

MIT
