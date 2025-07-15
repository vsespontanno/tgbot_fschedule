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

// интерфейс для взаимодействия с данными матчей и команд
type MatchesStore interface {
	GetMatches(ctx context.Context, collectionName string) ([]types.Match, error)
	SaveMatchesToMongoDB(matches []types.Match, from, to string) error
	UpdateMatchRatingInMongoDB(match types.Match, rating float64) error
	UpsertMatch(ctx context.Context, match types.Match) error
}

type MatchCalcStore interface {
	GetMatches(ctx context.Context, collectionName string) ([]types.Match, error)
	GetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error)
}

// структура для взаимодействия с данными матчей и команд
type MongoDBMatchesStore struct {
	dbName string
	client *mongo.Client
}

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

func (m *MongoDBMatchesStore) SaveMatchesToMongoDB(matches []types.Match, from, to string) error {
	if len(matches) == 0 {
		logrus.Infof("No matches found from %s to %s\n", from, to)
		return nil
	}

	collection := m.client.Database("football").Collection("matches")

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

func (m *MongoDBMatchesStore) UpdateMatchRatingInMongoDB(match types.Match, rating float64) error {
	collection := m.client.Database("football").Collection("matches")

	filter := bson.M{"id": match.ID}
	update := bson.M{"$set": bson.M{"rating": rating}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating match rating for ID %d: %w", match.ID, err)
	}

	return nil
}

func (m *MongoDBMatchesStore) UpsertMatch(ctx context.Context, match types.Match) error {
	collection := m.client.Database("football").Collection("matches")
	filter := bson.M{"id": match.ID}
	update := bson.M{"$set": match}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
