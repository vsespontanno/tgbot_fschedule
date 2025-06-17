package db

import (
	"context"
	"fmt"
	"football_tgbot/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamsStore interface {
	GetAllTeams(ctx context.Context, collectionName string) ([]types.Team, error)
	GetTeamLeague(ctx context.Context, collectionName string, id int) (string, error)
}

type MongoDBTeamsStore struct {
	dbName string
	client *mongo.Client
}

func NewMongoDBTeamsStore(client *mongo.Client, dbName string) *MongoDBTeamsStore {
	return &MongoDBTeamsStore{
		client: client,
		dbName: dbName,
	}
}

func (m *MongoDBTeamsStore) GetAllTeams(ctx context.Context, collectionName string) ([]types.Team, error) {
	var teams []types.Team
	collection := m.client.Database(m.dbName).Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (m *MongoDBTeamsStore) GetTeamLeague(ctx context.Context, collectionName string, id int) (string, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName + "_standings")
	fmt.Println(collectionName)
	var standing types.Standing
	filter := bson.M{"team.id": id}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Wrong League", nil
		}
		return "", fmt.Errorf("error finding team standing: %w", err)
	}
	return collectionName, nil
}
