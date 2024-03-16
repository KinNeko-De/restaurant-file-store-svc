package server

import (
	"context"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

func InitializeStorage(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) {
	fileRepository, err := connectToGoogleCloundStorage(ctx, storageConnected, storageDisconnected)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to storage")
		os.Exit(52)
	}
	injectFileRepository(fileRepository)
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
