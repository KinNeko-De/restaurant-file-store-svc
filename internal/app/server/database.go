package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

const MongoDBUriEnv = "MONGODB_URI"
const MongoDbDatabaseNameEnv = "MONGODB_DATABASE"

func InitializeDatabase(ctx context.Context, databaseConnected chan struct{}, databaseStopped chan struct{}) {
	err := connectToDatabase(ctx, databaseConnected, databaseStopped)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to database")
		os.Exit(51)
	}
}

func connectToDatabase(ctx context.Context, databaseConnected chan struct{}, databaseStopped chan struct{}) error {
	config, err := loadDatabaseConfig()
	if err != nil {
		return err
	}
	err = persistence.ConnectToDatabase(ctx, databaseConnected, databaseStopped, config)
	return err
}

func loadDatabaseConfig() (persistence.MongoDBConfig, error) {
	uri, found := os.LookupEnv(MongoDBUriEnv)
	if !found {
		return persistence.MongoDBConfig{}, fmt.Errorf("mongodb uri is not configured. Expected environment variable %v", MongoDBUriEnv)
	}

	databaseName, found := os.LookupEnv(MongoDbDatabaseNameEnv)
	if !found {
		return persistence.MongoDBConfig{}, fmt.Errorf("mongodb database name is not configured. Expected environment variable %v", MongoDbDatabaseNameEnv)
	}

	config := persistence.MongoDBConfig{
		HostUri:                uri,
		DatabaseName:           databaseName,
		FileMetadataCollection: "files",
	}
	return config, nil
}
