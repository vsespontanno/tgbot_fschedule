package mongodb

import (
	"context"
	"fmt"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamsStore interface {
	SaveTeamsToMongoDB(ctx context.Context, collectionName string, teams []types.Team) error
	UpsertTeamToMongoDB(ctx context.Context, collectionName string, teams types.Team) error
}

type TeamsCalcStore interface {
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

func (m *MongoDBTeamsStore) GetTeamLeague(ctx context.Context, collectionName string, id int) (string, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName)
	var team types.Team
	filter := bson.M{"id": id}
	err := collection.FindOne(ctx, filter).Decode(&team)
	if err != nil {
		return "", fmt.Errorf("error finding team with ID %d: %s", id, err)
	}
	return team.League, nil
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

func (m *MongoDBTeamsStore) SaveTeamsToMongoDB(ctx context.Context, collectionName string, teams []types.Team) error {
	collection := m.client.Database(m.dbName).Collection(collectionName)
	var documents []interface{}
	for _, team := range teams {
		documents = append(documents, team)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func (m *MongoDBTeamsStore) UpsertTeamToMongoDB(ctx context.Context, collectionName string, team types.Team) error {
	collection := m.client.Database(m.dbName).Collection(collectionName)
	filter := bson.M{"id": team.ID}
	update := bson.M{"$set": team}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
