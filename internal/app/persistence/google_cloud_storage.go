package persistence

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

type GoogleCloundStorageConfig struct {
}

func ConnectToGoogleCloudStorage(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}, config GoogleCloundStorageConfig) (*GoogleCloudStorageFileRepository, error) {
	client, err := storage.NewGRPCClient(ctx)
	if err != nil {
		close(storageDisconnected)
		return &GoogleCloudStorageFileRepository{}, err
	}
	go storageClientlistenToGracefulShutdown(client, storageDisconnected)

	close(storageConnected)

	return &GoogleCloudStorageFileRepository{
		Client: client,
	}, nil
}

func storageClientlistenToGracefulShutdown(client *storage.Client, storageDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	client.Close()

	close(storageDisconnected)
}
