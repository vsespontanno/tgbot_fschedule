package db

import (
	"context"
	"fmt"
	"football_tgbot/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// интерфейс для взаимодействия с данными матчей и команд
type MatchesStore interface {
	GetMatches(ctx context.Context, collectionName string) ([]types.Match, error)
	GetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error)
}

// структура для взаимодействия с данными матчей и команд
type MongoDBMatchesStore struct {
	dbName string
	client *mongo.Client
}

// функция для создания новой структуры для взаимодействия с данными матчей и команд
// client - клиент MongoDB
// dbName - имя базы данных
// возвращает *MongoDBMatchesStore

func NewMongoDBMatchesStore(client *mongo.Client, dbName string) *MongoDBMatchesStore {
	return &MongoDBMatchesStore{
		client: client,
		dbName: dbName,
	}
}

// функция для поиска документов в MONGODB
func (m *MongoDBMatchesStore) findDocuments(ctx context.Context, collectionName string, result interface{}) error {
	collection := m.client.Database(m.dbName).Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error finding documents in collection %s: %w", collectionName, err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return fmt.Errorf("error decoding documents in collection %s: %w", collectionName, err)
	}

	return nil
}

// функция для получения матчей из MONGODB
func (m *MongoDBMatchesStore) GetMatches(ctx context.Context, collectionName string) ([]types.Match, error) {
	var matches []types.Match
	if err := m.findDocuments(ctx, collectionName, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (m *MongoDBMatchesStore) GetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error) {
	collection := m.client.Database(m.dbName).Collection("matches")
	filter := bson.M{
		"$or": []bson.M{
			{"homeTeam.id": teamID},
			{"awayTeam.id": teamID},
		},
	}
	opts := options.Find().SetSort(bson.M{"date": -1}).SetLimit(int64(lastN))
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find recent matches: %w", err)
	}
	defer cursor.Close(ctx)

	var matches []types.Match
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, fmt.Errorf("failed to decode recent matches: %w", err)
	}
	return matches, err
}
