package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabase_ConfigMissing_MongoDBUri(t *testing.T) {
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	err := connectToDatabase(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), MongoDBUriEnv)
}

func TestConnectToDatabase_ConfigMissing_MongoDBDatabase(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "mongodb://rootuser:rootpassword@mongodb:27017")

	err := connectToDatabase(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), MongoDbDatabaseNameEnv)
}

func TestConnectToDatabase_ConfigIsComplete(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "mongodb://rootuser:rootpassword@mongodb:27017")
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	databaseConnected := make(chan struct{})
	databaseStopped := make(chan struct{})

	err := connectToDatabase(context.Background(), databaseConnected, databaseStopped)

	assert.NoError(t, err)
	// TODO assert that the database is connected
}
