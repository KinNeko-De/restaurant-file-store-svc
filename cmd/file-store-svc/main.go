package main

import (
	"context"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/health"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/metric"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server"
)

func main() {
	logger.SetLogLevel(logger.LogLevel)
	logger.Logger.Info().Msg("Starting application.")
	ctx := context.Background()

	err := file.FileRepositoryInstance.Initialize()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to initialize storage")
		os.Exit(45)
	}

	provider, err := metric.InitializeMetrics()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to initialize metrics")
		os.Exit(40)
	}

	grpcServerStopped := make(chan struct{})
	grpcServerStarted := make(chan struct{})
	databaseStopped := make(chan struct{})
	databaseConnected := make(chan struct{})
	go server.StartGrpcServer(grpcServerStopped, grpcServerStarted, ":3110")
	go file.ConnectToDatabase(ctx, databaseStopped, databaseConnected)

	go func() {
		<-databaseConnected
		<-grpcServerStarted
		health.Ready()
	}()

	<-databaseStopped
	<-grpcServerStopped
	provider.Shutdown(ctx)
	logger.Logger.Info().Msg("Application stopped.")
	os.Exit(0)
}
