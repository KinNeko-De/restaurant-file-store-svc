package persistence

import (
	"context"
	"os"
	"path"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

type PersistentVolumeConfig struct {
	Path string
}

func ConnectToPersistentVolume(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}, config PersistentVolumeConfig) (*PersistentVolumeFileRepository, error) {
	configTestFilePath := path.Join(config.Path, ".testpath")
	file, err := os.Create(configTestFilePath)
	if err != nil {
		return nil, err
	}
	file.Close()
	os.Remove(configTestFilePath)

	close(storageConnected)

	go PersistentVolumeListenToGracefulShutdown(storageDisconnected)

	return &PersistentVolumeFileRepository{}, nil
}

func PersistentVolumeListenToGracefulShutdown(storageDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	close(storageDisconnected)
}
