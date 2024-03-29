package persistence

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

type PersistentVolumeConfig struct {
	Path string
}

func ConnectToPersistentVolume(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}, config PersistentVolumeConfig) (*PersistentVolumeFileRepository, error) {
	err := EnsurePathExists(config)
	if err != nil {
		return nil, err
	}

	permissingErr := EnsureDirectoryIsWritable(config)
	if permissingErr != nil {
		return nil, permissingErr
	}

	close(storageConnected)

	go PersistentVolumeListenToGracefulShutdown(storageDisconnected)

	return &PersistentVolumeFileRepository{StoragePath: config.Path}, nil
}

func EnsurePathExists(config PersistentVolumeConfig) error {
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return fmt.Errorf("path '%v for persistent volume was not found. Please check the configuration of the mounted volume. Inner errrr: %w", config.Path, err)
	}
	return nil
}

func EnsureDirectoryIsWritable(config PersistentVolumeConfig) error {
	randomFileName := uuid.New().String()
	configTestFilePath := path.Join(config.Path, "."+randomFileName+"testpath")
	file, err := os.Create(configTestFilePath)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(configTestFilePath)
	return nil
}

func PersistentVolumeListenToGracefulShutdown(storageDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	close(storageDisconnected)
}
