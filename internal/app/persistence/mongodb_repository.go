package persistence

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func CreateMongoDBClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, err
}

func NewMongoDBRepository(ctx context.Context, client *mongo.Client, databaseName string, collection string) (*MongoDBRepository, error) {
	db := client.Database(databaseName)
	col := db.Collection(collection)

	return &MongoDBRepository{
		client:     client,
		database:   db,
		collection: col,
	}, nil
}
