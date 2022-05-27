package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ports/reader"
)

type MongoClient struct {
	C          *mongo.Client
	Collection *mongo.Collection
}

type PortDocument struct {
	Name        string    `bson:"name"`
	City        string    `bson:"city"`
	Country     string    `bson:"country"`
	Alias       []string  `bson:"alias"`
	Regions     []string  `bson:"regions"`
	Coordinates []float64 `bson:"coordinates"`
	Province    string    `bson:"province"`
	Timezone    string    `bson:"timezone"`
	Unlocs      []string  `bson:"unlocs"`
	Code        string    `bson:"code"`
}

func New(ctx context.Context, dbHost, user, password, database, collection string) (*MongoClient, error) {
	auth := options.Credential{
		Username: user,
		Password: password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbHost).SetAuth(auth))
	if err != nil {
		return nil, fmt.Errorf("error creating mongo client, %w", err)
	}

	coll := client.Database(database).Collection(collection)

	return &MongoClient{
		C:          client,
		Collection: coll,
	}, nil
}

func (mc *MongoClient) Write(ctx context.Context, result reader.ReadResult) error {
	d := mapToPortDocument(result)

	update := bson.D{{"$set", d}}
	options := options.Update().SetUpsert(true)

	_, err := mc.Collection.UpdateByID(ctx, result.Id, update, options)
	if err != nil {
		return fmt.Errorf("error inserting document for [%s], %w", result.Id, err)
	}

	return nil
}

func mapToPortDocument(result reader.ReadResult) *PortDocument {
	return &PortDocument{
		Name:        result.Port.Name,
		City:        result.Port.City,
		Country:     result.Port.Country,
		Alias:       result.Port.Alias,
		Regions:     result.Port.Regions,
		Coordinates: result.Port.Coordinates,
		Province:    result.Port.Province,
		Timezone:    result.Port.Timezone,
		Unlocs:      result.Port.Unlocs,
		Code:        result.Port.Code,
	}
}

func (mc *MongoClient) Stop(ctx context.Context) error {
	if err := mc.C.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
