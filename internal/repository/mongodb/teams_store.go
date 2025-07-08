package mongodb

import (
	"context"
	"fmt"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamsStore interface {
	GetAllTeams(ctx context.Context, collectionName string) ([]types.Team, error)
	GetTeamLeague(ctx context.Context, collectionName string, id int) (string, error)
	GetTeamsShortName(ctx context.Context, collectionName string, fullName string) (string, error)
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

	var standing types.Standing
	filter := bson.M{"team.id": id}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		return "Wrong League", nil
	}
	return collectionName, nil
}

func (m *MongoDBTeamsStore) GetTeamsShortName(ctx context.Context, collectionName string, fullName string) (string, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName)
	var team types.Team
	filter := bson.M{"name": fullName}
	err := collection.FindOne(ctx, filter).Decode(&team)
	if err != nil {
		return "", fmt.Errorf("error finding team with name %s: %w", fullName, err)
	}
	return team.ShortName, nil
}
