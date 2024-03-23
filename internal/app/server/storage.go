package server

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

const StorageTypeEnv = "STORAGE_TYPE"
const PersistentVolumePathEnv = "STORAGE_PERSISTENT_VOLUME_PATH"

type Storage int

const (
	Unspecified Storage = iota
	PersistentVolume
	StorageGoogleCloud
)

func InitializeStorage(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) {
	fileRepository, err := createFileRepository(ctx, storageConnected, storageDisconnected)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to connect to storage")
		os.Exit(52)
	}
	injectFileRepository(fileRepository)
}

func createFileRepository(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) (file.FileRepository, error) {
	storage, err := loadStorageTypConfig()
	if err != nil {
		return nil, err
	}
	var fileRepository file.FileRepository
	switch storage {
	case Unspecified, PersistentVolume:
		fileRepository, err = connectToPersistentVolume(ctx, storageConnected, storageDisconnected)
	case StorageGoogleCloud:
		fileRepository, err = connectToGoogleCloundStorage(ctx, storageConnected, storageDisconnected)
	}
	return fileRepository, err
}

func connectToPersistentVolume(ctx context.Context, storageConnected chan struct{}, storageDisconnected chan struct{}) (file.FileRepository, error) {
	config, err := loadPersistentVolumeConfig()
	if err != nil {
		return nil, err
	}
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

func loadStorageTypConfig() (Storage, error) {
	storageConfig, found := os.LookupEnv(StorageTypeEnv)
	if !found {
		return Unspecified, nil
	}
	storage, convertErr := strconv.Atoi(storageConfig)
	if convertErr != nil {
		return Unspecified, convertErr
	}
	return Storage(storage), nil
}

func loadGoogleCloudStorageConfig() (persistence.GoogleCloundStorageConfig, error) {
	config := persistence.GoogleCloundStorageConfig{}
	return config, nil
}

func loadPersistentVolumeConfig() (persistence.PersistentVolumeConfig, error) {
	path, found := os.LookupEnv(PersistentVolumePathEnv)
	if !found {
		return persistence.PersistentVolumeConfig{}, fmt.Errorf("persistent volume path is not configured. Expected environment variable %v", PersistentVolumePathEnv)
	}

	return persistence.PersistentVolumeConfig{
		Path: path,
	}, nil
}

func injectFileRepository(fileRepository file.FileRepository) {
	file.FileRepositoryInstance = fileRepository
}
