package mongodb

import (
	"context"
	"fmt"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Интерфейс для взаимодействия с данными матчей
type MatchesStore interface {
	GetMatchesInPeriod(ctx context.Context, league, from, to string) ([]types.Match, error)
	SaveMatchesToMongoDB(matches []types.Match, from, to string) error
	UpdateMatchRatingInMongoDB(match types.Match, rating float64) error
	UpsertMatch(ctx context.Context, match types.Match) error
}

// Интерфейс для взаимодействия с данными матчей в контексте калькуляции рейтинга матчей
type MatchCalcStore interface {
	GetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error)
	GetMatchesInPeriod(ctx context.Context, league string, from, to string) ([]types.Match, error)
}

// Структура для взаимодействия с данными матчей и команд
type MongoDBMatchesStore struct {
	dbName   string
	client   *mongo.Client
	collName string
}

// Конструктор структуры для взаимодействия с данными матчей и команд
func NewMongoDBMatchesStore(client *mongo.Client, dbName string, collName string) *MongoDBMatchesStore {
	return &MongoDBMatchesStore{
		client:   client,
		dbName:   dbName,
		collName: collName,
	}
}

// Общий метод для поиска документов в MONGODB
func (m *MongoDBMatchesStore) findDocuments(ctx context.Context, result interface{}) error {
	collection := m.client.Database(m.dbName).Collection(m.collName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error finding documents in collection %s: %w", m.collName, err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return fmt.Errorf("error decoding documents in collection %s: %w", m.collName, err)
	}

	return nil
}

// Метод для получения всех матчей из MONGODB
func (m *MongoDBMatchesStore) GetMatches(ctx context.Context) ([]types.Match, error) {
	var matches []types.Match
	if err := m.findDocuments(ctx, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

// Метод для получения последнтх матчей той или иной команды из MONGODB
func (m *MongoDBMatchesStore) GetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error) {
	collection := m.client.Database(m.dbName).Collection(m.collName)
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

// Метод для сохранения матчей в базу MongoDb
func (m *MongoDBMatchesStore) SaveMatchesToMongoDB(matches []types.Match, from, to string) error {
	if len(matches) == 0 {
		logrus.Infof("No matches found from %s to %s\n", from, to)
		return nil
	}

	collection := m.client.Database(m.dbName).Collection(m.collName)

	var documents []interface{}
	for _, match := range matches {
		documents = append(documents, match)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		return fmt.Errorf("error inserting matches: %v", err)
	}

	return nil
}

// Метод для обновления рейтинга матча в базе MongoDB
func (m *MongoDBMatchesStore) UpdateMatchRatingInMongoDB(match types.Match, rating float64) error {
	collection := m.client.Database(m.dbName).Collection(m.collName)

	filter := bson.M{"id": match.ID}
	update := bson.M{"$set": bson.M{"rating": rating}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating match rating for ID %d: %w", match.ID, err)
	}

	return nil
}

// Метод для обновления матчей в базе MongoDB, дабы при обновлении в фоне не было дубликатов
func (m *MongoDBMatchesStore) UpsertMatch(ctx context.Context, match types.Match) error {
	collection := m.client.Database(m.dbName).Collection(m.collName)
	filter := bson.M{"id": match.ID}
	update := bson.M{"$set": match}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Метод для получения всех матчей за тот или иной период
func (m *MongoDBMatchesStore) GetMatchesInPeriod(ctx context.Context, league, from, to string) ([]types.Match, error) {
	coll := m.client.Database(m.dbName).Collection(m.collName)
	var (
		matches []types.Match
		filter  bson.M
	)
	if league == "" {
		filter = bson.M{
			"utcdate": bson.M{
				"$gte": from + "T00:00:00Z",
				"$lte": to + "T23:59:59Z",
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{"utcdate": bson.M{
					"$gte": from + "T00:00:00Z",
					"$lte": to + "T23:59:59Z",
				}},
				{"competition.name": league},
			},
		}
	}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding matches in period for league %s: %w", league, err)
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &matches); err != nil {
		return nil, fmt.Errorf("error decoding matches in period for league %s: %w", league, err)
	}
	return matches, nil
}
