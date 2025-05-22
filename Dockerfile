# Используем официальный образ Go для сборки
FROM golang:1.23 as builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

# Финальный минимальный образ
FROM ubuntu:25.10

# Устанавливаем необходимые зависимости для работы шрифтов и изображений
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    fonts-noto \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник и необходимые файлы
COPY --from=builder /bot /app/bot
COPY .env /app/.env
COPY --from=builder /app/internal/bot/handlers/fonts /usr/share/fonts/noto  

# Открываем порт (если нужен)
EXPOSE 8080

# Запуск бота
ENTRYPOINT ["/app/bot"]