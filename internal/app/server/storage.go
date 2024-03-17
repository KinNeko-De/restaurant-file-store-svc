package server

import (
	"context"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

type Storage int

const (
	Unspecified Storage = iota
	PersistentVolume
	StorageGoogleCloud
)

func InitializeStorage(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) {
	storage := Unspecified

	var fileRepository file.FileRepository
	var err error
	switch storage {
	case Unspecified, PersistentVolume:
		fileRepository, err = connectToPersistentVolume(ctx, storageConnected, storageDisconnected)
	case StorageGoogleCloud:
		fileRepository, err = connectToGoogleCloundStorage(ctx, storageConnected, storageDisconnected)
	}

	if storage == StorageGoogleCloud {
		fileRepository, err = connectToGoogleCloundStorage(ctx, storageConnected, storageDisconnected)
	}
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to storage")
		os.Exit(52)
	}
	injectFileRepository(fileRepository)
}

func connectToPersistentVolume(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) (file.FileRepository, error) {
	config := persistence.PersistentVolumeConfig{} // TODO load path
	repository, err := persistence.ConnectToPersistentVolume(ctx, storageConnected, storageDisconnected, config)
	return repository, err
}

func connectToGoogleCloundStorage(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) (file.FileRepository, error) {
	config, err := loadGoogleCloudStorageConfig()
	if err != nil {
		return nil, err
	}
	repository, err := persistence.ConnectToGoogleCloudStorage(ctx, storageConnected, storageDisconnected, config)
	return repository, err
}

func loadGoogleCloudStorageConfig() (persistence.GoogleCloundStorageConfig, error) {
	config := persistence.GoogleCloundStorageConfig{}
	return config, nil
}

func injectFileRepository(fileRepository file.FileRepository) {
	file.FileRepositoryInstance = fileRepository
}
