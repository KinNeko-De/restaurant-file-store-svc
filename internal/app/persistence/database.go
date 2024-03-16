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

func ConnectToMongoDB(ctx context.Context, databaseConnected chan struct{}, databaseDisconnected chan struct{}, config MongoDBConfig) (file.FileMetadataRepository, error) {
	logger.Logger.Debug().Msg("connecting to mongodb")
	mongoDBRepository, err := initializeMongoDbFileMetadataRepository(ctx, config)
	if err != nil {
		close(databaseDisconnected)
		return nil, err
	}
	go listenToGracefulShutdown(ctx, mongoDBRepository.client, databaseDisconnected)

	close(databaseConnected)
	return mongoDBRepository, nil
}

func initializeMongoDbFileMetadataRepository(ctx context.Context, config MongoDBConfig) (*MongoDBRepository, error) {
	client, err := createClient(ctx, config.HostUri)
	if err != nil {
		return nil, err
	}

	fileMetadataRepository, err := NewMongoDBRepository(ctx, client, config.DatabaseName, config.FileMetadataCollection)
	return fileMetadataRepository, err
}

func listenToGracefulShutdown(ctx context.Context, client *mongo.Client, databaseDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	client.Disconnect(ctx)
	close(databaseDisconnected)
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
