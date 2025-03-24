package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Leagues = map[string]string{
	"Ligue1":          "FL1",
	"LaLiga":          "PD",
	"PremierLeague":   "PL",
	"Bundesliga":      "BL1",
	"SerieA":          "SA",
	"ChampionsLeague": "CL",
}

func ConnectToMongoDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	return client, nil
}
