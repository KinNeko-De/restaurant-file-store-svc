package persistence

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

type GoogleCloundStorageConfig struct {
}

type GoogleCloudStorageFileRepository struct {
	Client *storage.Client
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

func (g *GoogleCloudStorageFileRepository) CreateFile(ctx context.Context, fileId uuid.UUID, revisionId uuid.UUID, chunkSize int) (io.WriteCloser, error) {
	bucket := g.Client.Bucket("kinneko-de")
	objectname := fileId.String() + "/" + revisionId.String()
	object := bucket.Object(objectname)
	writer := object.NewWriter(ctx)
	writer.ChunkSize = chunkSize
	// TODO: checksum control
	return writer, nil
}
