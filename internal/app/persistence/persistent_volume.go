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
	directoryErr := EnsurePathExists(config)
	if directoryErr != nil {
		return nil, directoryErr
	}

	fileErr := EnsureDirectoryIsWritable(config)
	if fileErr != nil {
		return nil, fileErr
	}

	close(storageConnected)

	go PersistentVolumeListenToGracefulShutdown(storageDisconnected)

	return &PersistentVolumeFileRepository{}, nil
}

func EnsureDirectoryIsWritable(config PersistentVolumeConfig) error {
	configTestFilePath := path.Join(config.Path, ".testpath")
	file, err := os.Create(configTestFilePath)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(configTestFilePath)
	return nil
}

func EnsurePathExists(config PersistentVolumeConfig) error {
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		directoryErr := os.MkdirAll(config.Path, os.ModePerm)
		if directoryErr != nil {
			return directoryErr
		}
	}
	return nil
}

func PersistentVolumeListenToGracefulShutdown(storageDisconnected chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	close(storageDisconnected)
}
