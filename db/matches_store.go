package db

import (
	"context"
	"fmt"
	"football_tgbot/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MatchesStore defines the interface for interacting with match and team data.
type MatchesStore interface {
	GetTeams(ctx context.Context, collectionName string) ([]types.Team, error)
	GetMatches(ctx context.Context, collectionName string) ([]types.Match, error)
	GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error) // Добавлен GetStandings
}

// MongoDBMatchesStore is a concrete implementation of MatchesStore for MongoDB.
type MongoDBMatchesStore struct {
	dbName string
	client *mongo.Client
}

// NewMongoDBMatchesStore creates a new MongoDBMatchesStore.
func NewMongoDBMatchesStore(client *mongo.Client, dbName string) *MongoDBMatchesStore {
	return &MongoDBMatchesStore{
		client: client,
		dbName: dbName,
	}
}

// findDocuments is a generic helper function to find documents in a collection.
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

func (m *MongoDBMatchesStore) GetTeams(ctx context.Context, collectionName string) ([]types.Team, error) {
	var teams []types.Team
	if err := m.findDocuments(ctx, collectionName, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (m *MongoDBMatchesStore) GetMatches(ctx context.Context, collectionName string) ([]types.Match, error) {
	var matches []types.Match
	if err := m.findDocuments(ctx, collectionName, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (m *MongoDBMatchesStore) GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error) {
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
