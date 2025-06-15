package db

import (
	"context"
	"database/sql"
	"fmt"

	"football_tgbot/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// функция для подключения к MongoDB
// uri - строка подключения к MongoDB
// возвращает *mongo.Client и ошибку

func ConnectToMongoDB(uri string) (*mongo.Client, error) {
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

type PG struct {
	Conn *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PG, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB, cfg.PostgresHost)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &PG{Conn: db}, nil
}
