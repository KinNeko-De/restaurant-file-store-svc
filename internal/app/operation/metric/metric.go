package metric

import (
	"context"
	"fmt"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/operation/logger"
	"go.opentelemetry.io/otel"

	"github.com/kinneko-de/restaurant-file-store-svc/build"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/go-logr/zerologr"
)

type OtelConfig struct {
	OtelMetricEndpoint string // is used by the otel sdk to identify the endpoint to send metrics to. According to document it Will be set implicitly by the otel sdk. But it does not work. I set it explicitly.
	OtelServiceName    string // is used by the otel sdk to identify the service name. I found no way to set it explicitly by the otel sdk. According to the specification setting an attribute with name "service.name" should work, but it does not.
}

var (
	version  = "0.2.0"
	provider *metric.MeterProvider
	meter    api.Meter
)

func InitializeMetrics(ctx context.Context, config OtelConfig) (*metric.MeterProvider, error) {
	metricLogger := zerologr.New(&logger.Logger)
	otel.SetLogger(metricLogger)

	provider, err := initializeOpenTelemetry(ctx, config)
	return provider, err
}

func initializeOpenTelemetry(ctx context.Context, config OtelConfig) (*metric.MeterProvider, error) {
	ressource, err := createRessource(config)
	if err != nil {
		return nil, err
	}

	readers, err := createReader(ctx, config)
	if err != nil {
		return nil, err
	}

	views := createViews()
	provider := createProvider(ressource, readers, views)
	metricError := createMetrics(provider, config)
	return provider, metricError
}

func createRessource(config OtelConfig) (*resource.Resource, error) {
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.OtelServiceName),
			semconv.ServiceVersionKey.String(build.Version),
		))
	if err != nil {
		return nil, fmt.Errorf("failed to create ressource for metric reader: %w", err)
	}

	return res, nil
}

func createViews() []metric.View {
	return []metric.View{}
}

func createReader(ctx context.Context, config OtelConfig) ([]metric.Reader, error) {
	otelGrpcExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(config.OtelMetricEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metric reader to otel collector: %w", err)
	}
	otelReader := metric.NewPeriodicReader(otelGrpcExporter)

	consoleExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metric reader to console: %w", err)
	}
	consoleReader := metric.NewPeriodicReader(consoleExporter)

	return []metric.Reader{otelReader, consoleReader}, nil
}

func createProvider(ressource *resource.Resource, readers []metric.Reader, views []metric.View) *metric.MeterProvider {
	options := []metric.Option{
		metric.WithResource(ressource),
		metric.WithView(views...),
	}
	for _, reader := range readers {
		options = append(options, metric.WithReader(reader))
	}

	provider = metric.NewMeterProvider(
		options...,
	)
	otel.SetMeterProvider(provider)
	return provider
}

// https://opentelemetry.io/docs/specs/otel/metrics/semantic_conventions/
func createMetrics(provider *metric.MeterProvider, config OtelConfig) error {
	// I decided to use the service name here as scope because this service is a microservice. one sccope per service approach.
	meter = provider.Meter(config.OtelServiceName, api.WithInstrumentationVersion(version))

	// var err error
	// errorTemplate := "failed to initialize metric '%v' %w"
	// define metrics

	return nil
}

func ForceFlush(ctx context.Context) {
	provider.ForceFlush(ctx)
}
