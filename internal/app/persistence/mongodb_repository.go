package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
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

func (repository *MongoDBRepository) StoreFileMetadata(ctx context.Context, fileMetadata file.FileMetadata) error {
	dataModel := fileMetadataToDataModel(fileMetadata)

	_, err := repository.collection.InsertOne(ctx, dataModel)
	if err != nil {
		return fmt.Errorf("failed to insert file metadata: %w", err)
	}

	return nil
}

func (repository *MongoDBRepository) StoreRevision(ctx context.Context, fileId uuid.UUID, revision file.Revision) error {
	requestedId := fileId.String()
	dataModel := revisionToDataModel(revision)

	result, err := repository.collection.UpdateOne(ctx, bson.M{"_id": requestedId}, bson.M{"$push": bson.M{"revisions": dataModel}})
	if err != nil {
		return fmt.Errorf("failed to insert revision: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("file metadata not found")
	}

	return nil
}

func (repository *MongoDBRepository) FetchFileMetadata(ctx context.Context, fileId uuid.UUID) (file.FileMetadata, error) {
	requestedId := fileId.String()
	var dataModel fileMetadata
	err := repository.collection.FindOne(ctx, bson.M{"_id": requestedId}).Decode(&dataModel)
	if err != nil {
		return file.FileMetadata{}, fmt.Errorf("failed to fetch file metadata: %w", err)
	}

	return fileMetadataToDomainModel(dataModel), nil
}

func (repository *MongoDBRepository) NotFoundError() error {
	return mongo.ErrNoDocuments
}

func fileMetadataToDomainModel(dataModel fileMetadata) file.FileMetadata {
	return file.FileMetadata{
		Id:        uuid.MustParse(dataModel.Id),
		Revisions: revisionsToDomainModel(dataModel.Revisions),
	}
}

func revisionsToDomainModel(revision []revision) []file.Revision {
	var domainModel []file.Revision
	for _, revision := range revision {
		domainModel = append(domainModel, revisionToDomainModel(revision))
	}
	return domainModel
}

func revisionToDomainModel(revision revision) file.Revision {
	return file.Revision{
		Id:        uuid.MustParse(revision.Id),
		Extension: revision.Extension,
		MediaType: revision.MediaType,
		Size:      revision.Size,
		CreatedAt: revision.CreatedAt,
	}
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

func CreateMongoDBClient(ctx context.Context, config MongoDBConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.HostUri)

	if config.Timeout != 0 {
		clientOptions = clientOptions.SetTimeout(config.Timeout)
	}

	client, err := mongo.Connect(ctx, clientOptions)
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
