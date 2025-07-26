package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToMongoDB создает подключение к MongoDB и возвращает клиент
// Принимает строку подключения в формате MongoDB URI
// Возвращает указатель на mongo.Client и ошибку, если есть
func ConnectToMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5 секунд таймаут для пинга
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// СОздание индекса
	if err := createMatchesIndexes(client, "football"); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	logrus.Info("Connected to MongoDB")
	return client, nil
}

// createMatchesIndexes создает индекс на коллекции матчей в MongoDB
// Индекс включает поля homeTeam.id, awayTeam.id и date для ускорения запросов
// Принимает указатель на mongo.Client и имя базы данных
// Возвращает ошибку, если не удалось создать индекс
// Использует контекст с таймаутом 10 секунд для создания индекса
func createMatchesIndexes(client *mongo.Client, dbName string) error {
	collection := client.Database(dbName).Collection("matches")
	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "homeTeam.id", Value: 1},
			{Key: "awayTeam.id", Value: 1},
			{Key: "date", Value: -1},
		},
		Options: options.Index().SetName("homeTeam.id_1_awayTeam.id_1_date_-1"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return fmt.Errorf("failed to create index on matches: %w", err)
	}
	logrus.Info("Created index on matches collection")
	return nil
}

// ConnectToPostgres создает подключение к PostgreSQL и возвращает указатель на sql.DB
// Принимает параметры подключения: пользователь, пароль, имя базы данных, хост и порт
func ConnectToPostgres(user, password, dbname, host, port string) (*sql.DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)                 // Максимальное число открытых соединений
	db.SetMaxIdleConns(5)                  // Максимальное число бездействующих соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения
	db.SetConnMaxIdleTime(2 * time.Minute) // Максимальное время бездействия соединения\
	logrus.Info("Connected to Postgres")
	return db, nil
}
