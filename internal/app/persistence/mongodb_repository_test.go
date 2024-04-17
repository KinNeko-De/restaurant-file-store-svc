//go:build component

package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/test/testing/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStoreFileMetadata(t *testing.T) {
	input := file.FileMetadata{
		Id: uuid.New(),

		Revisions: []file.Revision{
			{
				Id:        uuid.New(),
				Extension: ".txt",
				MediaType: "text/plain",
				Size:      1024,
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	expectedFileMetadata := fileMetadata{
		Id: input.Id.String(),
		Revisions: []revision{
			{
				Id:        input.Revisions[0].Id.String(),
				Extension: input.Revisions[0].Extension,
				MediaType: input.Revisions[0].MediaType,
				Size:      input.Revisions[0].Size,
				CreatedAt: input.Revisions[0].CreatedAt,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb.MongoDbServer))
	require.Nil(t, err)
	defer disconnectClient(ctx, client)
	sut, err := NewMongoDBRepository(ctx, client, uuid.NewString(), t.Name()+uuid.NewString())
	require.Nil(t, err)
	defer tearDown(t, sut.collection)

	err = sut.StoreFileMetadata(ctx, input)
	require.Nil(t, err)

	var actualFileMetadata fileMetadata
	err = sut.collection.FindOne(ctx, bson.M{"_id": input.Id.String()}).Decode(&actualFileMetadata)
	require.Nil(t, err)
	assertDatamodeEqual(t, expectedFileMetadata, actualFileMetadata)
}

func TestFetchFileMetadata(t *testing.T) {
	fileId := uuid.New()

	expectedFileMetadata := file.FileMetadata{
		Id: fileId,
		Revisions: []file.Revision{
			{
				Id:        uuid.New(),
				Extension: ".txt",
				MediaType: "text/plain; charset=utf-8",
				Size:      1024,
				CreatedAt: time.Now().UTC().Round(time.Millisecond),
			},
		},
	}

	existingFile := fileMetadata{
		Id: fileId.String(),
		Revisions: []revision{
			{
				Id:        expectedFileMetadata.Revisions[0].Id.String(),
				Extension: expectedFileMetadata.Revisions[0].Extension,
				MediaType: expectedFileMetadata.Revisions[0].MediaType,
				Size:      expectedFileMetadata.Revisions[0].Size,
				CreatedAt: expectedFileMetadata.Revisions[0].CreatedAt,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb.MongoDbServer))
	require.Nil(t, err)
	defer disconnectClient(ctx, client)

	sut, err := NewMongoDBRepository(ctx, client, uuid.NewString(), t.Name()+uuid.NewString())
	require.Nil(t, err)
	defer tearDown(t, sut.collection)
	_, err = sut.collection.InsertOne(ctx, existingFile)
	require.Nil(t, err)

	actualFileMetadata, err := sut.FetchFileMetadata(ctx, fileId)
	require.Nil(t, err)

	assert.Equal(t, expectedFileMetadata, actualFileMetadata)
}

func TestStoreRevision_ARevisionIsAddedToAFile(t *testing.T) {
	fileId := uuid.New()
	storedRevision := file.Revision{
		Id:        uuid.New(),
		Extension: ".txt",
		MediaType: "text/plain",
		Size:      1024,
		CreatedAt: time.Now().UTC().Add(-time.Hour).Round(time.Millisecond),
	}
	storedFileMetadata := file.FileMetadata{
		Id: fileId,
		Revisions: []file.Revision{
			storedRevision,
		},
	}

	newRevsion := file.Revision{
		Id:        uuid.New(),
		Extension: ".md",
		MediaType: "text/plain; charset=utf-8",
		Size:      1069,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
	}

	expectedFileMetadata := file.FileMetadata{
		Id: fileId,
		Revisions: []file.Revision{
			storedRevision,
			newRevsion,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb.MongoDbServer))
	require.Nil(t, err)
	sut, err := NewMongoDBRepository(ctx, client, uuid.NewString(), t.Name()+uuid.NewString())
	require.Nil(t, err)
	defer tearDown(t, sut.collection)

	err = sut.StoreFileMetadata(ctx, storedFileMetadata)
	require.Nil(t, err)
	err = sut.StoreRevision(ctx, fileId, newRevsion)
	require.Nil(t, err)
	actualFileMetadata, err := sut.FetchFileMetadata(ctx, fileId)
	require.Nil(t, err)
	assert.Equal(t, expectedFileMetadata, actualFileMetadata)
}

func TestFetchFileMetadata_FileDoesNotExists(t *testing.T) {
	fileId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb.MongoDbServer))
	require.Nil(t, err)
	defer disconnectClient(ctx, client)

	sut, err := NewMongoDBRepository(ctx, client, uuid.NewString(), t.Name()+uuid.NewString())
	require.Nil(t, err)
	defer tearDown(t, sut.collection)

	_, actualError := sut.FetchFileMetadata(ctx, fileId)
	require.NotNil(t, actualError)
	assert.True(t, errors.Is(actualError, mongo.ErrNoDocuments))
}

func assertDatamodeEqual(t *testing.T, expectedFileMetadata fileMetadata, actualFileMetadata fileMetadata) {
	assert.Equal(t, expectedFileMetadata.Id, actualFileMetadata.Id)
	for i := range expectedFileMetadata.Revisions {
		assertRevisionsEqual(t, expectedFileMetadata.Revisions[i], actualFileMetadata.Revisions[i])
	}
}

func assertRevisionsEqual(t *testing.T, expectedRevision revision, actualRevision revision) {
	assert.Equal(t, expectedRevision.Id, actualRevision.Id)
	assert.Equal(t, expectedRevision.Extension, actualRevision.Extension)
	assert.Equal(t, expectedRevision.MediaType, actualRevision.MediaType)
	assert.Equal(t, expectedRevision.Size, actualRevision.Size)
	assert.WithinDuration(t, expectedRevision.CreatedAt, actualRevision.CreatedAt, time.Millisecond)
}

func tearDown(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Failed to drop collection: %v", err)
	}
}

func disconnectClient(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
