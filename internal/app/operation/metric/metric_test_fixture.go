package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func MockMetric(config OtelConfig) (*metric.ManualReader, *metric.MeterProvider) {

	reader := metric.NewManualReader()

	ressource, _ := createRessource(config)
	views := createViews()
	provider := createProvider(ressource, []metric.Reader{reader}, views)
	createMetrics(provider, config)
	return reader, provider
}

func ActualMetrics(t *testing.T, reader *metric.ManualReader) metricdata.ResourceMetrics {
	rm := metricdata.ResourceMetrics{}
	err := reader.Collect(context.Background(), &rm)
	if err != nil {
		assert.NoError(t, err)
	}
	return rm
}
