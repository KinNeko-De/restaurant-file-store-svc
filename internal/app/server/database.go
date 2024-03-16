package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

const MongoDBUriEnv = "MONGODB_URI"
const MongoDbDatabaseNameEnv = "MONGODB_DATABASE"

func InitializeDatabase(ctx context.Context, databaseConnected chan struct{}, databaseDisconnected chan struct{}) {
	fileMetadataRepository, err := connectToMongoDB(ctx, databaseConnected, databaseDisconnected)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to database")
		os.Exit(51)
	}

	file.FileMetadataRepositoryInstance = fileMetadataRepository
}

func injectFileMetadaRepository(fileMetadataRepository file.FileMetadataRepository) {
	file.FileMetadataRepositoryInstance = fileMetadataRepository
}

func connectToMongoDB(ctx context.Context, databaseConnected chan struct{}, databaseDisconnected chan struct{}) (file.FileMetadataRepository, error) {
	config, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	repository, err := persistence.ConnectToMongoDB(ctx, databaseConnected, databaseDisconnected, config)
	return repository, err
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
