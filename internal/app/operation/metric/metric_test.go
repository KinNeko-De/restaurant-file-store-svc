package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeMetrics_ConfigIsComplete(t *testing.T) {
	expectedServiceName := "expectedServiceName"
	expectedOtelMetricEndpoint := "otel-collector:4317"

	config := OtelConfig{
		OtelServiceName:    expectedServiceName,
		OtelMetricEndpoint: expectedOtelMetricEndpoint,
	}

	createdProvider, err := InitializeMetrics(context.Background(), config)

	assert.NoError(t, err)
	assert.NotNil(t, createdProvider)
	assert.NotNil(t, provider)
	assert.NotNil(t, meter)
}
