package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"

	"go.mongodb.org/mongo-driver/mongo"
)

const connectTimeout = 10 * time.Second

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

	timeout, cancel := context.WithTimeout(ctx, connectTimeout)
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

	if err := initializeFileMetadataRepository(ctx, client); err != nil {
		return err
	}

	close(databaseConnected)
	return nil
}

func createClient(ctx context.Context) (*mongo.Client, error) {
	uri := "mongodb://rootuser:rootpassword@mongodb:27017" // TODO: get from env
	client, err := persistence.CreateMongoDBClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	/* TODO start database before connecting
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	*/

	return client, nil
}

func initializeFileMetadataRepository(ctx context.Context, client *mongo.Client) error {
	databaseName := "restaurant-file-store-db" // TODO: get from env

	var err error
	file.FileMetadataRepositoryInstance, err = persistence.NewMongoDBRepository(ctx, client, databaseName, "files") // TODO: get from env or extract

	return err
}
