package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func (repository *MongoDBRepository) StoreFileMetadata(ctx context.Context, fileMetadata file.FileMetadata) error {
	dataModel := fileMetadataToDataModel(fileMetadata)

	_, err := repository.collection.InsertOne(ctx, dataModel)
	if err != nil {
		return fmt.Errorf("failed to insert file metadata: %v", err)
	}

	return nil
}

func fileMetadataToDataModel(domainModel file.FileMetadata) fileMetadata {
	return fileMetadata{
		Id:        domainModel.Id.String(),
		Revisions: revisionsToDataModel(domainModel.Revisions),
	}
}

func revisionsToDataModel(domainModel []file.Revision) []revision {
	var dataModel []revision
	for _, revision := range domainModel {
		dataModel = append(dataModel, revisionToDataModel(revision))
	}
	return dataModel
}

func revisionToDataModel(domainModel file.Revision) revision {
	return revision{
		Id:        domainModel.Id.String(),
		Extension: domainModel.Extension,
		MediaType: domainModel.MediaType,
		Size:      domainModel.Size,
		CreatedAt: domainModel.CreatedAt,
	}
}

type fileMetadata struct {
	Id        string `bson:"_id"`
	Revisions []revision
}

type revision struct {
	Id        string `bson:"_id"`
	Extension string
	MediaType string
	Size      uint64
	CreatedAt time.Time
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
