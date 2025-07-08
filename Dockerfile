FROM golang:1.23-alpine

WORKDIR /app

RUN apk update && apk add --no-cache \
    ca-certificates \
    fontconfig \
    gcc \
    musl-dev

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Проверяем целостность модулей
RUN go mod verify

# Копируем весь проект в контейнер
COPY . .

# Устанавливаем переменные окружения для сборки с поддержкой C-библиотек
ENV CGO_ENABLED=1

# Открываем порт
EXPOSE 8080

# Сборка приложения
RUN go build -tags dynamic -o tgbot_fschedule ./cmd/bot

# Запуск приложения
CMD ["./tgbot_fschedule"]