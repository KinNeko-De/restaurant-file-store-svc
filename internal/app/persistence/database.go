package persistence

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/shutdown"
	"go.mongodb.org/mongo-driver/mongo"
)

var FileMetadataRepository file.FileMetadataRepository

func ConnectToDatabase(ctx context.Context, databaseStoped chan struct{}, databaseConnected chan struct{}) {
	err := connectToDatabase(ctx, databaseStoped, databaseConnected)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to database")
		os.Exit(50)
	}
}

func connectToDatabase(ctx context.Context, databaseStopped chan struct{}, databaseConnected chan struct{}) error {
	gracefulStop := shutdown.CreateGracefulStop()
	logger.Logger.Debug().Msg("connecting to database")

	timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	abort, stop := signal.NotifyContext(timeout, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	client, err := createClient(abort)
	if err != nil {
		select {
		case <-abort.Done():
			close(databaseStopped)
			return nil
		default:
			return err
		}
	}

	go func() {
		<-gracefulStop
		client.Disconnect(ctx)
		close(databaseStopped)
	}()

	if err := initializeMongoDBRepository(ctx, client); err != nil {
		return err
	}

	close(databaseConnected)
	return nil
}

func createClient(ctx context.Context) (*mongo.Client, error) {
	uri := "mongodb://rootuser:rootpassword@mongodb:27017" // TODO: get from env
	client, err := CreateClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func initializeMongoDBRepository(ctx context.Context, client *mongo.Client) error {
	databaseName := "restaurant-file-store-db" // TODO: get from env

	var err error
	FileMetadataRepository, err = NewMongoDBRepository(ctx, client, databaseName, "files")

	return err
}
