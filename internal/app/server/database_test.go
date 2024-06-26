//go:build unit

package server

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabase_ConfigMissing_MongoDBUri(t *testing.T) {
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	_, err := connectToMongoDB(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), MongoDBUriEnv)
}

func TestConnectToDatabase_ConfigMissing_MongoDBDatabase(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "mongodb://rootuser:rootpassword@mongodb:27017")

	_, err := connectToMongoDB(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), MongoDbDatabaseNameEnv)
}

func TestConnectToDatabase_UriMalformed(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "invalidUri")
	t.Setenv(MongoDbDatabaseNameEnv, "testdatabase")

	_, err := connectToMongoDB(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
}

func TestInitializeDatabase_AnyError_AppCrash(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		// will crash because of missing environment variables
		InitializeDatabase(context.Background(), make(chan struct{}), make(chan struct{}))
		return
	}

	runningApp := exec.Command(os.Args[0], "-test.run=TestInitializeDatabase_AnyError_AppCrash")
	runningApp.Env = append(os.Environ(), "EXECUTE=1")
	err := runningApp.Run()
	require.NotNil(t, err)
	exitCode := err.(*exec.ExitError).ExitCode()
	assert.Equal(t, 51, exitCode)
}

func TestConnectToDatabase_ConfigIsComplete_PingToDatabaseFailedBecauseItIsNotStarted(t *testing.T) {
	t.Setenv(MongoDBUriEnv, "mongodb://ihavenoaccess:iamalittlehack@localhost:666")
	t.Setenv(MongoDbDatabaseNameEnv, "notexisting")
	t.Setenv(MongoDbTimeoutEnv, "100ms")

	databaseConnected := make(chan struct{})
	databaseStopped := make(chan struct{})

	_, err := connectToMongoDB(context.Background(), databaseConnected, databaseStopped)
	// assert the database is connected is ensured over the ping of the client

	require.Error(t, err)
	assert.Contains(t, err.Error(), "server selection error: context deadline exceeded")
}
