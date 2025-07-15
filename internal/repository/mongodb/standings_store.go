package mongodb

import (
	"context"
	"fmt"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// интерфейс для взаимодействия с данными матчей и команд
type StandingsStore interface {
	GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error)
	SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error
}

type StandingsCalcStore interface {
	GetTeamStanding(ctx context.Context, collectionName string, id int) (int, error)
}

// структура для взаимодействия с данными матчей и команд
type MongoDBStandingsStore struct {
	dbName string
	client *mongo.Client
}

// функция для создания новой структуры для взаимодействия с данными матчей и команд
// client - клиент MongoDB
// dbName - имя базы данных
// возвращает *MongoDBMatchesStore

func NewMongoDBStandingsStore(client *mongo.Client, dbName string) *MongoDBStandingsStore {
	return &MongoDBStandingsStore{
		client: client,
		dbName: dbName,
	}
}

// функция для получения таблицы из MONGODB
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

// функция для сохранения таблицы в MONGODB
func (m *MongoDBStandingsStore) SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error {
	collection := m.client.Database(m.dbName).Collection(collectionName + "_standings")

	// Clear existing standings
	_, err := collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	// Insert new standings
	documents := make([]interface{}, len(standings))
	for i, standing := range standings {
		documents[i] = standing
	}

	_, err = collection.InsertMany(ctx, documents)
	return err
}

func (m *MongoDBStandingsStore) GetTeamStanding(ctx context.Context, collectionName string, id int) (int, error) {
	fullCollectionName := collectionName + "_standings"
	collection := m.client.Database(m.dbName).Collection(fullCollectionName)

	// Добавлено логирование для отладки

	var standing types.Standing
	filter := bson.M{"team.id": id}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Добавлено логирование при отсутствии документов
			fmt.Printf("No document found for team ID %d in %s.%s\n", id, m.dbName, fullCollectionName)
			return -1, nil
		}
		return 0, fmt.Errorf("error finding team standing for ID %d in %s: %w", id, fullCollectionName, err)
	}

	return standing.Position, nil
}
