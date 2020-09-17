package otelstats

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
)

func readMetricFromPath(metricPath string) int64 {
	dataBytes, err := ioutil.ReadFile(metricPath)
	if err != nil {
		log.Printf("Failed to read metric at %v: %v\n", metricPath, err)
	}

	data, err := strconv.Atoi(strings.TrimSuffix(string(dataBytes), "\n"))
	if err != nil {
		log.Printf("Failed to convert metric value at %v to int: %v\n", metricPath, err)
	}

	value := int64(data)
	return value
}

func createMetric(metricName string, metricInfo []MetricInfo) {
	meter := global.MeterProvider().Meter("otel-stats")
	metric.Must(meter).NewInt64ValueObserver(metricName,
		func(_ context.Context, result metric.Int64ObserverResult) {
			for _, info := range metricInfo {
				result.Observe(
					readMetricFromPath(info.Path),
					kv.String("device", info.Label),
				)
			}
		},
		metric.WithDescription(metricName),
	)
}

// CreateOtelMetricsForStatsfs creates a otel metric for every
// metric found in the given statsfsPath
func CreateOtelMetricsForStatsfs(statsfsPath string) error {
	m, err := CreateStatsfsMetrics(statsfsPath)
	if err != nil {
		return fmt.Errorf("Failed to create statsfs metrics for %v: %v", statsfsPath, err)
	}
	m.Print()

	for _, subsysMetrics := range m.Metrics {
		for metricName, metricInfo := range subsysMetrics.Metrics {
			createMetric(metricName, metricInfo)
		}
	}
	return nil
}
