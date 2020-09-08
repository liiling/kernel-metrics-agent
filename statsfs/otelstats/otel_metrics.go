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

// InitOtelPipeline initializes an OpenTelemetry pipeline
// that crawls a user defined path and exports all the available
// stats to a backend of choice (gcp, stdout, prometheus)
func InitOtelPipeline(exporterName, statsfsPath string) {
	exporter, err := InitExporter(exporterName)
	if err != nil {
		log.Panicf("Failed to initialize exporter %v: %v\n", exporterName, err)
	}
	if exporter != nil {
		defer exporter.Stop()
	}

	err = createOtelMetricsForStatsfs(statsfsPath)
	if err != nil {
		log.Panic(err)
	}

	for {
	}
}

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

func createOtelMetricsForStatsfs(statsfsPath string) error {
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
