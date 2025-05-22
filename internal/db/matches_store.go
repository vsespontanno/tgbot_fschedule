package db

import (
	"context"
	"fmt"
	"football_tgbot/internal/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// интерфейс для взаимодействия с данными матчей и команд
type MatchesStore interface {
	GetTeams(ctx context.Context, collectionName string) ([]types.Team, error)
	GetMatches(ctx context.Context, collectionName string) ([]types.Match, error)
	GetStandings(ctx context.Context, collectionName string) ([]types.Standing, error)
	SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error
	// Новые методы для работы с рейтингами
	GetTeamRatings(ctx context.Context, collectionName string) ([]types.TeamRating, error)
	SaveTeamRatings(ctx context.Context, collectionName string, ratings []types.TeamRating) error
	UpdateTeamRating(ctx context.Context, collectionName string, rating types.TeamRating) error
	GetTeamRating(ctx context.Context, collectionName string, teamID int) (*types.TeamRating, error)
	GetTeamStanding(ctx context.Context, collectionName string, teamID int) (*types.Standing, error)
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

// функция для получения команд из MONGODB
func (m *MongoDBMatchesStore) GetTeams(ctx context.Context, collectionName string) ([]types.Team, error) {
	var teams []types.Team
	if err := m.findDocuments(ctx, collectionName, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// функция для получения матчей из MONGODB
func (m *MongoDBMatchesStore) GetMatches(ctx context.Context, collectionName string) ([]types.Match, error) {
	var matches []types.Match
	if err := m.findDocuments(ctx, collectionName, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

// функция для получения таблицы из MONGODB
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

// функция для сохранения таблицы в MONGODB
func (m *MongoDBMatchesStore) SaveStandings(ctx context.Context, collectionName string, standings []types.Standing) error {
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

// функция для получения рейтингов команд из MONGODB
func (m *MongoDBMatchesStore) GetTeamRatings(ctx context.Context, collectionName string) ([]types.TeamRating, error) {
	var ratings []types.TeamRating
	collection := m.client.Database(m.dbName).Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding team ratings: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &ratings); err != nil {
		return nil, fmt.Errorf("error decoding team ratings: %w", err)
	}

	return ratings, nil
}

// функция для сохранения рейтингов команд в MONGODB
func (m *MongoDBMatchesStore) SaveTeamRatings(ctx context.Context, collectionName string, ratings []types.TeamRating) error {
	collection := m.client.Database(m.dbName).Collection(collectionName)

	// Очищаем существующие рейтинги
	_, err := collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("error clearing existing ratings: %w", err)
	}

	// Вставляем новые рейтинги
	documents := make([]interface{}, len(ratings))
	for i, rating := range ratings {
		documents[i] = rating
	}

	_, err = collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting new ratings: %w", err)
	}

	return nil
}

// функция для обновления рейтинга конкретной команды
func (m *MongoDBMatchesStore) UpdateTeamRating(ctx context.Context, collectionName string, rating types.TeamRating) error {
	collection := m.client.Database(m.dbName).Collection(collectionName)

	filter := bson.M{"teamId": rating.TeamID}
	update := bson.M{"$set": rating}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error updating team rating: %w", err)
	}

	return nil
}

// функция для получения рейтинга конкретной команды
func (m *MongoDBMatchesStore) GetTeamRating(ctx context.Context, collectionName string, teamID int) (*types.TeamRating, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName)

	var standing types.Standing
	// Используем правильное имя поля для поиска команды
	filter := bson.M{"team.id": teamID}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding team standing: %w", err)
	}

	rating := &types.TeamRating{
		TeamID:      teamID,
		TeamName:    standing.Team.Name,
		Position:    standing.Position,
		Points:      standing.Points,
		GoalDiff:    standing.GoalDifference,
		Form:        0,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	return rating, nil
}

// GetTeamStanding получает данные о положении команды в турнирной таблице
func (m *MongoDBMatchesStore) GetTeamStanding(ctx context.Context, collectionName string, teamID int) (*types.Standing, error) {
	collection := m.client.Database(m.dbName).Collection(collectionName)

	var standing types.Standing
	filter := bson.M{"team.id": teamID}
	err := collection.FindOne(ctx, filter).Decode(&standing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding team standing: %w", err)
	}

	return &standing, nil
}
