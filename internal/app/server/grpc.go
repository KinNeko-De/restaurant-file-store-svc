package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	googleHealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/health"
	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/mylogger"
)

func StartGrpcServer(grpcServerStopped chan struct{}, grpcServerStarted chan struct{}) {
	port := ":3110" // todo load from env, move os.exit up to here and refactor tests

	startGrpcServer(grpcServerStopped, grpcServerStarted, port)
}

func startGrpcServer(grpcServerStopped chan struct{}, grpcServerStarted chan struct{}, port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		mylogger.Logger.Error().Err(err).Msgf("Failed to listen on port %v", port)
		os.Exit(50)
	}

	grpcServer := configureGrpcServer()
	healthServer := googleHealth.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	health.Initialize(healthServer)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	mylogger.Logger.Debug().Msg("starting grpc server")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			mylogger.Logger.Error().Err(err).Msg("failed to serve grpc server")
			os.Exit(51)
		}
	}()
	close(grpcServerStarted)

	stop := <-gracefulStop
	healthServer.Shutdown()
	grpcServer.GracefulStop()
	mylogger.Logger.Debug().Msgf("http server stopped. received signal %s", stop)
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
}
