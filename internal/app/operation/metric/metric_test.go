package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeMetrics_ConfigMissing_ServiceName(t *testing.T) {
	expectedOtelMetricEndpoint := "otel-collector:4317"
	t.Setenv(OtelMetricEndpointEnv, expectedOtelMetricEndpoint)

	createdProvider, err := InitializeMetrics()

	require.Error(t, err)
	assert.Contains(t, err.Error(), ServiceNameEnv)
	assert.Nil(t, createdProvider)
}

func TestInitializeMetrics_ConfigMissing_OtelMetricEndpoint(t *testing.T) {
	expectedServiceName := "expectedServiceName"
	t.Setenv(ServiceNameEnv, expectedServiceName)

	createdProvider, err := InitializeMetrics()

	require.Error(t, err)
	assert.Contains(t, err.Error(), OtelMetricEndpointEnv)
	assert.Nil(t, createdProvider)
}

func TestInitializeMetrics_ConfigIsComplete(t *testing.T) {
	expectedServiceName := "expectedServiceName"
	expectedOtelMetricEndpoint := "otel-collector:4317"
	t.Setenv(ServiceNameEnv, expectedServiceName)
	t.Setenv(OtelMetricEndpointEnv, expectedOtelMetricEndpoint)

	createdProvider, err := InitializeMetrics()

	assert.NoError(t, err)
	assert.Equal(t, expectedServiceName, config.OtelServiceName)
	assert.Equal(t, expectedOtelMetricEndpoint, config.OtelMetricEndpoint)
	assert.NotNil(t, provider)
	assert.NotNil(t, createdProvider)
	assert.NotNil(t, meter)
	assert.NotNil(t, previewRequested)
	assert.NotNil(t, previewDelivered)
	assert.NotNil(t, documentGenerateSuccessful)
	assert.NotNil(t, documentGenerateFailed)
	assert.NotNil(t, documentGenerateDuration)
}
