package persistence

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBConfig struct {
	HostUri                string
	DatabaseName           string
	FileMetadataCollection string
}

func ConnectToDatabase(ctx context.Context, databaseConnected chan struct{}, databaseStopped chan struct{}, config MongoDBConfig) error {
	var client *mongo.Client
	go listenToGracefulShutdown(ctx, client, databaseStopped)
	logger.Logger.Debug().Msg("connecting to database")

	err := initializePersistence(ctx, config)
	if err != nil {
		return err
	}

	close(databaseConnected)
	return nil
}

func initializePersistence(ctx context.Context, config MongoDBConfig) error {
	fileMetadataRepository, err := initializeMongoDbFileMetadataRepository(ctx, config)
	if err != nil {
		return err
	}

	file.FileMetadataRepositoryInstance = fileMetadataRepository
	return nil
}

func initializeMongoDbFileMetadataRepository(ctx context.Context, config MongoDBConfig) (file.FileMetadataRepository, error) {
	client, err := createClient(ctx, config.HostUri)
	if err != nil {
		return nil, err
	}

	fileMetadataRepository, err := NewMongoDBRepository(ctx, client, config.DatabaseName, config.FileMetadataCollection)
	return fileMetadataRepository, err
}

func listenToGracefulShutdown(ctx context.Context, client *mongo.Client, databaseStopped chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	if client != nil {
		client.Disconnect(ctx)
	}
	close(databaseStopped)
}

func createClient(ctx context.Context, hostUri string) (*mongo.Client, error) {
	gracefulAbort, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	client, err := CreateMongoDBClient(gracefulAbort, hostUri)
	if err != nil {
		return nil, err
	}

	err = client.Ping(gracefulAbort, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
