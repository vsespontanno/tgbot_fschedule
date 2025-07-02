# Используем официальный образ Go для сборки
FROM golang:1.23 as builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
# Загружаем зависимости
RUN go mod download
# Проверяем целостность модулей
RUN go mod verify

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bot ./cmd/bot

# Финальный минимальный образ
FROM ubuntu:22.04

# Устанавливаем зависимости (fontconfig нужен для работы со шрифтами)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    fontconfig \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник
COPY --from=builder /bot /app/bot

# Копируем папку со шрифтами
COPY fonts /app/fonts

# Обновляем кэш шрифтов
RUN fc-cache -fv

# Открываем порт
EXPOSE 8080

# Запуск бота
ENTRYPOINT ["/app/bot"]