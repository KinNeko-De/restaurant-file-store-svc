//go:build component

package server

import (
	"context"
	"testing"

	"github.com/kinneko-de/restaurant-file-store-svc/test/testing/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestConnectToDatabase_ConfigIsComplete(t *testing.T) {
	t.Setenv(MongoDBUriEnv, mongodb.MongoDbServer)
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	databaseConnected := make(chan struct{})
	databaseStopped := make(chan struct{})

	err := connectToDatabase(context.Background(), databaseConnected, databaseStopped)

	assert.NoError(t, err)
	// assert the database is connected is ensured over the ping of the client
}
