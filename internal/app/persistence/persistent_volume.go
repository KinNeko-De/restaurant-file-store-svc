package persistence

import (
	"context"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

type PersistentVolumeConfig struct {
	Path string
}

func ConnectToPersistentVolume(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}, config PersistentVolumeConfig) (*PersistentVolumeFileRepository, error) {
	// TODO add check if path exists and is accessible
	close(storageConnected)

	go PersistentVolumeListenToGracefulShutdown(storageDisconnected)

	return &PersistentVolumeFileRepository{}, nil
}

func PersistentVolumeListenToGracefulShutdown(storageDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	close(storageDisconnected)
}
