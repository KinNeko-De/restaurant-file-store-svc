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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	expected := &file.FileMetadata{
		Id: uuid.New(),
	}

	err = sut.StoreFileMetadata(ctx, expected)
	require.Nil(t, err)

	// Now let's try to find the inserted document
	var actual *file.FileMetadata
	err = sut.collection.FindOne(ctx, bson.M{"_id": expected.Id.String()}).Decode(&actual)
	require.Nil(t, err)
	// Check if the inserted document is the same as the original one
	assert.Equal(t, expected, actual)
}

func tearDown(t *testing.T, collection *mongo.Collection) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Failed to drop collection: %v", err)
	}
}
