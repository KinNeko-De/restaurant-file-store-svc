//go:build component

package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateFileMetadata(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://rootuser:rootpassword@localhost:27017"))
	require.Nil(t, err)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	sut, err := NewMongoDBRepository(ctx, client, "test", "test")
	require.Nil(t, err)
	defer tearDown(t, sut.collection)

	input := file.FileMetadata{
		Id:        uuid.New(),
		CreatedAt: time.Now().UTC(),
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
		Id:        input.Id.String(),
		CreatedAt: input.CreatedAt,
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

	err = sut.StoreFileMetadata(ctx, input)
	require.Nil(t, err)

	var actualFileMetadata fileMetadata
	err = sut.collection.FindOne(ctx, bson.M{"_id": input.Id.String()}).Decode(&actualFileMetadata)
	require.Nil(t, err)

	assert.Equal(t, expectedFileMetadata, actualFileMetadata)
}

func tearDown(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Failed to drop collection: %v", err)
	}
}
