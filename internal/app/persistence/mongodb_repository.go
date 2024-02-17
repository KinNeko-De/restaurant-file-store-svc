package persistence

import (
	"context"
	"fmt"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (repository *MongoDBRepository) CreateFileMetadata(ctx context.Context, fileMetadata file.FileMetadata) error {
	_, err := repository.collection.InsertOne(ctx, fileMetadata)
	if err != nil {
		return fmt.Errorf("failed to insert file metadata: %v", err)
	}

	return nil
}

func CreateMongoDBClient(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, err
}

func NewMongoDBRepository(ctx context.Context, client *mongo.Client, databaseName string, collectionName string) (*MongoDBRepository, error) {
	db := client.Database(databaseName)
	col := db.Collection(collectionName)

	return &MongoDBRepository{
		client:     client,
		database:   db,
		collection: col,
	}, nil
}
