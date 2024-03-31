package server

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	googleHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/health"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server/shutdown"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/file"
)

func StartGrpcServer(grpcServerStarted chan struct{}, grpcServerStopped chan struct{}) {
	port := ":3110" // todo load from env, move os.exit up to here and refactor tests

	startGrpcServer(grpcServerStarted, grpcServerStopped, port)
}

func startGrpcServer(grpcServerStarted chan struct{}, grpcServerStopped chan struct{}, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Logger.Error().Err(err).Msgf("Failed to listen on port %v", port)
		os.Exit(50)
	}

	grpcServer := configureGrpcServer()
	healthServer := googleHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	health.Initialize(healthServer)

	var gracefulStop = shutdown.CreateGracefulStop()
	logger.Logger.Debug().Msg("starting grpc server")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Logger.Error().Err(err).Msg("failed to serve grpc server")
			os.Exit(51)
		}
	}()
	close(grpcServerStarted)

	stop := <-gracefulStop
	healthServer.Shutdown()
	grpcServer.GracefulStop()
	logger.Logger.Debug().Msgf("http server stopped. received signal %s", stop)
	close(grpcServerStopped)
}

func configureGrpcServer() *grpc.Server {
	// Handling of panic to prevent crash from example nil pointer exceptions
	logPanic := func(p any) (err error) {
		log.Error().Any("method", p).Err(err).Msg("Recovered from panic.")
		return status.Errorf(codes.Internal, "Internal server error occured.")
	}

	opts := []recovery.Option{
		recovery.WithRecoveryHandler(logPanic),
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(opts...),
		),
		grpc.StreamInterceptor(
			recovery.StreamServerInterceptor(opts...),
		),
	)
	RegisterAllGrpcServices(grpcServer)
	return grpcServer
}

func RegisterAllGrpcServices(grpcServer *grpc.Server) {
	apiRestaurantFile.RegisterFileServiceServer(grpcServer, &file.FileServiceServer{})
}
