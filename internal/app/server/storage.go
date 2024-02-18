package server

import (
	"context"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/persistence"
)

func InitializeStorage(ctx context.Context, storageStoped chan struct{}, storageConnected chan struct{}) error {
	err := persistence.ConnectToStorage(ctx, storageStoped, storageConnected)
	return err
}
