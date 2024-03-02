//go:build unit

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeMetrics_ConfigMissing_ServiceName(t *testing.T) {
	expectedOtelMetricEndpoint := "otel-collector:4317"
	t.Setenv(OtelMetricEndpointEnv, expectedOtelMetricEndpoint)

	createdProvider, err := initializeMetrics(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), ServiceNameEnv)
	assert.Nil(t, createdProvider)
}

func TestInitializeMetrics_ConfigMissing_OtelMetricEndpoint(t *testing.T) {
	expectedServiceName := "expectedServiceName"
	t.Setenv(ServiceNameEnv, expectedServiceName)

	createdProvider, err := initializeMetrics(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), OtelMetricEndpointEnv)
	assert.Nil(t, createdProvider)
}

func TestInitializeMetrics_ConfigIsComplete(t *testing.T) {
	expectedServiceName := "expectedServiceName"
	expectedOtelMetricEndpoint := "otel-collector:4317"
	t.Setenv(ServiceNameEnv, expectedServiceName)
	t.Setenv(OtelMetricEndpointEnv, expectedOtelMetricEndpoint)

	createdProvider, err := initializeMetrics(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, createdProvider)
}
