package health

import (
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	liveness  = "liveness"
	readiness = "readiness"
)

var server *health.Server

func Initialize(srv *health.Server) {
	server = srv
	Live()
	NotReady()
}

func Live() {
	server.SetServingStatus(liveness, grpc_health_v1.HealthCheckResponse_SERVING)
}

func NotLive() {
	server.SetServingStatus(liveness, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func Ready() {
	server.SetServingStatus(readiness, grpc_health_v1.HealthCheckResponse_SERVING)
}

func NotReady() {
	server.SetServingStatus(readiness, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}
