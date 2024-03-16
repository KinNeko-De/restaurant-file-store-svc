//go:build component

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDatabase_ConfigIsComplete(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "mongodb://rootuser:rootpassword@mongodb:27017")
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	databaseConnected := make(chan struct{})
	databaseStopped := make(chan struct{})

	err := connectToDatabase(context.Background(), databaseConnected, databaseStopped)

	assert.NoError(t, err)
	// TODO assert that the database is connected
}
