package mongodb

import (
	"context"
	"fmt"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Интерфейс для взаимодействия с данными турнирных таблиц
type StandingsStore interface {
	GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error)
	SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error
}

// Интерфейс для взаимодействия с данными турнирных таблиц в контексте калькуляции рейтинга матчей
type StandingsCalcStore interface {
	GetTeamStanding(ctx context.Context, collectionName string, id int) (int, error)
}

// Структура для взаимодействия с данными турнирных таблиц
type MongoDBStandingsStore struct {
	dbName string
	client *mongo.Client
}

// Конструктор структуры для взаимодействия с данными турнирных таблиц
// client - клиент MongoDB
// dbName - имя базы данных
func NewMongoDBStandingsStore(client *mongo.Client, dbName string) *MongoDBStandingsStore {
	return &MongoDBStandingsStore{
		client: client,
		dbName: dbName,
	}
}

// Метод для получения турнирных таблиц из MONGODB
func (m *MongoDBStandingsStore) GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error) {
	var standings []types.Standing
	collection := m.client.Database(m.dbName).Collection(collectionName + "_standings")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding standings: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &standings); err != nil {
		return nil, fmt.Errorf("error decoding standings: %w", err)
	}

	// Добавляем логирование
	fmt.Printf("Raw standings from MongoDB for %s:\n", collectionName)
	for i, s := range standings {
		fmt.Printf("Standing %d: %+v\n", i, s)
	}

	return standings, nil
}

// Метод для для сохранения турнирных таблиц в MONGODB
func (m *MongoDBStandingsStore) SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error {
	collection := m.client.Database(m.dbName).Collection(collectionName + "_standings")

	// Insert new standings
	documents := make([]interface{}, len(standings))
	for i, standing := range standings {
		documents[i] = standing
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}

// Метод для получения турнирной таблицы из MONGODB
func (m *MongoDBStandingsStore) GetTeamStanding(ctx context.Context, collectionName string, id int) (int, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName + "_standings")

	var standing types.Standing
	filter := bson.M{"team.id": id}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Добавлено логирование при отсутствии документов
			fmt.Printf("No document found for team ID %d in collection %s\n", id, collectionName)
			return -1, nil
		}
		return 0, fmt.Errorf("error finding team standing for ID %d in %s: %w", id, collectionName, err)
	}

	return standing.Position, nil
}
