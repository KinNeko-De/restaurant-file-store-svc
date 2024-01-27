package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"

	"go.mongodb.org/mongo-driver/mongo"
)

func ConnectToDatabase(ctx context.Context, databaseStoped chan struct{}, databaseConnected chan struct{}) {
	err := connectToDatabase(ctx, databaseStoped, databaseConnected)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to database")
		os.Exit(51)
	}
}

func connectToDatabase(ctx context.Context, databaseStopped chan struct{}, databaseConnected chan struct{}) error {
	gracefulStop := shutdown.CreateGracefulStop()
	logger.Logger.Debug().Msg("connecting to database")

	gracefulAbort, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	client, err := createClient(gracefulAbort)
	if err != nil {
		close(databaseStopped)
		return err
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

	/*
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
