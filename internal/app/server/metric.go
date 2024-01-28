package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	internalMetric "github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/metric"

	"go.opentelemetry.io/otel/sdk/metric"
)

const ServiceNameEnv = "OTEL_SERVICE_NAME"
const OtelMetricEndpointEnv = "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"

func InitializeMetrics(ctx context.Context) *metric.MeterProvider {
	provider, err := initializeMetrics(ctx)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to initialize metrics")
		os.Exit(40)
	}
	return provider
}

func initializeMetrics(ctx context.Context) (*metric.MeterProvider, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	provider, err := internalMetric.InitializeMetrics(ctx, config)
	return provider, err
}

func loadConfig() (internalMetric.OtelConfig, error) {
	endpoint, found := os.LookupEnv(OtelMetricEndpointEnv)
	if !found {
		return internalMetric.OtelConfig{}, fmt.Errorf("otel metric endpoint is not configured. Expected environment variable %v", OtelMetricEndpointEnv)
	}

	serviceName, found := os.LookupEnv(ServiceNameEnv)
	if !found {
		return internalMetric.OtelConfig{}, fmt.Errorf("otel service name is not configured. Expected environment variable %v", ServiceNameEnv)
	}

	return internalMetric.OtelConfig{
		OtelMetricEndpoint: endpoint,
		OtelServiceName:    serviceName,
	}, nil
}
