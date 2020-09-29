package otelstats

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/api/metric"
)

const meterName = "otel-stats"

// CreateOtelMetricsForStatsfs creates a otel metric for every
// metric found in the given statsfsPath
func CreateOtelMetricsForStatsfs(statsfsPath string) error {
	m, err := NewStatsfsMetrics(statsfsPath)
	if err != nil {
		return fmt.Errorf("failed to create statsfs metrics for %v: %v", statsfsPath, err)
	}
	m.Print()

	for _, subsysMetrics := range m.Metrics {
		for metricName, metricInfo := range subsysMetrics.Metrics {
			if err := createMetric(metricName, metricInfo); err != nil {
				return err
			}
		}
	}
	return nil
}

func createMetric(metricName string, metricInfo MetricInfo) error {
	switch metricInfo.Type {
	case intType:
		return createIntMetric(metricName, metricInfo)
	case floatType:
		return createFloatMetric(metricName, metricInfo)
	}
	return nil
}

func createIntMetric(metricName string, metricInfo MetricInfo) error {
	meter := global.MeterProvider().Meter(meterName)

	intMetricCallback := func(_ context.Context, result metric.Int64ObserverResult) {
		for metricPath, metricLabels := range metricInfo.PathToLabel {
			if val, err := readIntMetricFromPath(metricPath); err != nil {
				log.Printf("Error reading metric at %v: %v\n", metricPath, err)
			} else {
				result.Observe(val, metricLabelToOtel(metricLabels)...)
			}
		}
	}

	switch metricInfo.Flag {
	case cumulative:
		_, err := meter.NewInt64SumObserver(metricName, intMetricCallback,
			metric.WithDescription(metricInfo.Desc),
		)
		if err != nil {
			return fmt.Errorf("failed to create metric %v: %v", metricName, err)
		}
	case gauge:
		_, err := meter.NewInt64ValueObserver(metricName, intMetricCallback,
			metric.WithDescription(metricInfo.Desc),
		)

		if err != nil {
			return fmt.Errorf("failed to create metric %v: %v", metricName, err)
		}
	default:
		return fmt.Errorf("unknown metric flag: %v", metricInfo.Flag)
	}
	return nil
}

func createFloatMetric(metricName string, metricInfo MetricInfo) error {
	meter := global.MeterProvider().Meter(meterName)

	floatMetricCallback := func(_ context.Context, result metric.Float64ObserverResult) {
		for metricPath, metricLabels := range metricInfo.PathToLabel {
			if val, err := readFloatMetricFromPath(metricPath); err != nil {
				log.Printf("Error reading metric at %v: %v\n", metricPath, err)
			} else {
				result.Observe(val, metricLabelToOtel(metricLabels)...)
			}
		}
	}

	switch metricInfo.Flag {
	case cumulative:
		_, err := meter.NewFloat64SumObserver(metricName, floatMetricCallback,
			metric.WithDescription(metricInfo.Desc),
		)
		if err != nil {
			return fmt.Errorf("failed to create metric %v: %v", metricName, err)
		}
	case gauge:
		_, err := meter.NewFloat64ValueObserver(metricName, floatMetricCallback,
			metric.WithDescription(metricInfo.Desc),
		)

		if err != nil {
			return fmt.Errorf("failed to create metric %v: %v", metricName, err)
		}
	default:
		return fmt.Errorf("unknown metric flag: %v", metricInfo.Flag)
	}
	return nil
}

func metricLabelToOtel(metricLabels []MetricLabel) []label.KeyValue {
	labels := []label.KeyValue{}
	for _, mLabel := range metricLabels {
		labels = append(labels, label.String(mLabel.Key, mLabel.Value))
	}
	return labels
}

func readIntMetricFromPath(metricPath string) (int64, error) {
	dataBytes, err := ioutil.ReadFile(metricPath)
	if err != nil {
		return -1, fmt.Errorf("failed to read metric at %v: %v", metricPath, err)
	}

	data, err := strconv.Atoi(strings.TrimSuffix(string(dataBytes), "\n"))
	if err != nil {
		return -1, fmt.Errorf("failed to convert metric value at %v to int: %v", metricPath, err)
	}

	return int64(data), nil
}

//TODO: actually read float...
func readFloatMetricFromPath(metricPath string) (float64, error) {
	dataBytes, err := ioutil.ReadFile(metricPath)
	if err != nil {
		return -1, fmt.Errorf("failed to read metric at %v: %v", metricPath, err)
	}

	data, err := strconv.Atoi(strings.TrimSuffix(string(dataBytes), "\n"))
	if err != nil {
		return -1, fmt.Errorf("failed to convert metric value at %v to int: %v", metricPath, err)
	}

	return float64(data), nil
}
