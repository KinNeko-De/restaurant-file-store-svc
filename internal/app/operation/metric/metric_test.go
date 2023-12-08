package metric

import (
	"context"
	"testing"
	"time"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
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

func TestPreviewRequested(t *testing.T) {
	expectedValue := int64(1)
	expectedMetric := metricdata.Metrics{
		Name:        MetricNameDocumentPreviewRequested,
		Description: MetricDescriptionDocumentPreviewRequested,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: expectedValue}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	reader, provider := MockMetric()
	defer provider.Shutdown(context.Background())

	PreviewRequested()

	actualMetrics := ActualMetrics(t, reader)
	require.Len(t, actualMetrics.ScopeMetrics, 1)
	require.Len(t, actualMetrics.ScopeMetrics[0].Metrics, 1)
	metricdatatest.AssertEqual(t, expectedMetric, actualMetrics.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
}

func TestPreviewDelivered(t *testing.T) {
	expectedValue := int64(1)
	expectedMetric := metricdata.Metrics{
		Name:        MetricNameDocumentPreviewDelivered,
		Description: MetricDescriptionDocumentPreviewDelivered,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: expectedValue}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	reader, provider := MockMetric()
	defer provider.Shutdown(context.Background())

	PreviewDelivered()

	actualMetrics := ActualMetrics(t, reader)
	require.Len(t, actualMetrics.ScopeMetrics, 1)
	require.Len(t, actualMetrics.ScopeMetrics[0].Metrics, 1)
	metricdatatest.AssertEqual(t, expectedMetric, actualMetrics.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
}

func TestPreviewRequested_DocumentGenerated_Successful(t *testing.T) {
	expectedValue := int64(1)
	expectedDocumentType := "Invoice"
	expectedAttriutes := attribute.NewSet(attribute.String(MetricAttributeDocumentType, expectedDocumentType))
	expectedDocumentMetric := metricdata.Metrics{
		Name:        MetricNameDocumentGenerateSuccessful,
		Description: MetricDescriptionDocumentGenerateSuccessful,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: expectedAttriutes, Value: expectedValue}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}
	duration := 5280127 * time.Microsecond
	expectedDuration := float64(duration.Milliseconds())
	expectedDurationMetric := metricdata.Metrics{
		Name:        MetricNameDocumentGenerateDuration,
		Description: MetricDescriptionDocumentGenerateDuration,
		Unit:        "ms",
		Data: metricdata.Histogram[float64]{
			DataPoints: []metricdata.HistogramDataPoint[float64]{{Attributes: expectedAttriutes, Sum: expectedDuration, Count: 1,
				Bounds:       []float64{1000, 4000, 7000, 10000, 20000},
				BucketCounts: []uint64{0, 0, 1, 0, 0, 0}}},
			Temporality: metricdata.CumulativeTemporality,
		},
	}

	reader, provider := MockMetric()
	defer provider.Shutdown(context.Background())

	DocumentGenerated(expectedDocumentType, duration, nil)

	actualMetrics := ActualMetrics(t, reader)
	require.Len(t, actualMetrics.ScopeMetrics, 1)
	require.Len(t, actualMetrics.ScopeMetrics[0].Metrics, 2)
	metricdatatest.AssertEqual(t, expectedDocumentMetric, actualMetrics.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
	metricdatatest.AssertEqual(t, expectedDurationMetric, actualMetrics.ScopeMetrics[0].Metrics[1], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreExemplars())
}

func TestPreviewRequested_DocumentGenerated_Failed(t *testing.T) {
	errorOccurred := errors.New("expected error")
	expectedValue := int64(1)
	expectedDocumentType := "Invoice"
	expectedAttriutes := attribute.NewSet(attribute.String(MetricAttributeDocumentType, expectedDocumentType))
	expectedDocumentMetric := metricdata.Metrics{
		Name:        MetricNameDocumentGenerateFailed,
		Description: MetricDescriptionDocumentGenerateFailed,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: expectedAttriutes, Value: expectedValue}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}
	duration := 5280127 * time.Microsecond

	reader, provider := MockMetric()
	defer provider.Shutdown(context.Background())

	DocumentGenerated(expectedDocumentType, duration, errorOccurred)

	actualMetrics := ActualMetrics(t, reader)
	require.Len(t, actualMetrics.ScopeMetrics, 1)
	require.Len(t, actualMetrics.ScopeMetrics[0].Metrics, 1)
	metricdatatest.AssertEqual(t, expectedDocumentMetric, actualMetrics.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
}
