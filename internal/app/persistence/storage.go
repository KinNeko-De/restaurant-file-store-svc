package persistence

import (
	"context"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"
)

func ConnectToStorage(ctx context.Context, storageStopped chan struct{}, storageConnected chan struct{}) error {

	go listenToStorageGracefulShutdown(ctx, storageStopped)

	file.FileRepositoryInstance = &LocalStorageRepository{}

	close(storageConnected)
	return nil
}

func listenToStorageGracefulShutdown(ctx context.Context, storageStopped chan struct{}) {
	gracefulShutdown := shutdown.CreateGracefulStop()
	<-gracefulShutdown
	close(storageStopped)
}
