package main

import (
	"context"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/health"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server"
)

func main() {
	logger.SetLogLevel(logger.LogLevel)
	logger.Logger.Info().Msg("Starting application.")

	ctx, cancel := context.WithCancel(context.Background())

	provider := server.InitializeMetrics(ctx)

	grpcServerStarted := make(chan struct{})
	grpcServerStopped := make(chan struct{})
	storageConnected := make(chan struct{})
	storageDisconnected := make(chan struct{})
	databaseConnected := make(chan struct{})
	databaseDisconnected := make(chan struct{})
	go server.StartGrpcServer(grpcServerStarted, grpcServerStopped)
	go server.InitializeStorage(ctx, storageConnected, storageDisconnected)
	go server.InitializeDatabase(ctx, databaseConnected, databaseDisconnected)

	go func() {
		<-databaseConnected
		<-storageConnected
		<-grpcServerStarted
		logger.Logger.Info().Msg("Application started.")
		health.Ready()
	}()

	<-grpcServerStopped
	<-storageDisconnected
	<-databaseDisconnected
	provider.Shutdown(ctx)
	cancel()
	logger.Logger.Info().Msg("Application stopped.")
	os.Exit(0)
}
