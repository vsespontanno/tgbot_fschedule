# Используем официальный образ Go для сборки
FROM golang:1.23 as builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
# Сначала загружаем зависимости, чтобы использовать кэш Docker
RUN go mod download
# Проверяем целостность загруженных модулей (хорошая практика)
RUN go mod verify

# Копируем исходники. Для лучшего кэширования можно копировать выборочно,
# но для простоты оставим COPY . .
COPY . .

# Собираем бинарник
# Флаги -s -w убирают отладочную информацию и символы, уменьшая размер бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot ./cmd/bot

# Финальный минимальный образ
# Рекомендуется использовать LTS (Long-Term Support) версию Ubuntu, например, ubuntu:22.04.
# ubuntu:25.10 - это версия с краткосрочной поддержкой.
FROM ubuntu:22.04

# Устанавливаем необходимые зависимости
# ca-certificates - для HTTPS соединений
# fontconfig - для управления шрифтами и утилиты fc-cache
# fonts-noto - если вашему приложению нужны системные шрифты Noto.
# Если все необходимые шрифты находятся в /app/internal/bot/handlers/fonts вашего проекта,
# то установку fonts-noto можно пропустить, оставив только fontconfig.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    fontconfig \
    fonts-noto \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /bot /app/bot

# Копируем .env файл.
# ВАЖНО: Я понимаю, что для ля production-окружений крайне рекомендуется передавать конфигурацию
# через переменные окружения Docker (например, `docker run -e KEY=VALUE`)
# или использовать Docker secrets, а не копировать .env файл напрямую в образ.
# Но сейчас делается всё локально
COPY .env /app/.env

# Обновляем кэш шрифтов, чтобы система и приложения "увидели" новые шрифты
RUN fc-cache -fv

# Открываем порт (если ваш бот работает в режиме webhook или имеет какой-либо HTTP интерфейс)
EXPOSE 8080

# Запуск бота
ENTRYPOINT ["/app/bot"]
