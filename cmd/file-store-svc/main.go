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

	grpcServerStopped := make(chan struct{})
	grpcServerStarted := make(chan struct{})
	databaseDisconnected := make(chan struct{})
	databaseConnected := make(chan struct{})
	go server.StartGrpcServer(grpcServerStopped, grpcServerStarted)
	go server.InitializeDatabase(ctx, databaseDisconnected, databaseConnected)

	go func() {
		<-databaseConnected
		<-grpcServerStarted
		logger.Logger.Info().Msg("Application started.")
		health.Ready()
	}()

	<-grpcServerStopped
	<-databaseDisconnected
	provider.Shutdown(ctx)
	cancel()
	logger.Logger.Info().Msg("Application stopped.")
	os.Exit(0)
}
