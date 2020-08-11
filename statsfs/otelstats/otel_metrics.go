package otelstats

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
)

// InitOtelPipeline initializes an OpenTelemetry pipeline
// that crawls a user defined path and exports all the available
// stats to a backend of choice (gcp, stdout, prometheus)
func InitOtelPipeline(statsfsPath string) {
	fmt.Println("In otel-metrics!")
	exporter := InitExporter()
	if exporter != nil {
		defer exporter.Stop()
	}
	createOtelMetricsForStatsfs(statsfsPath)

	for {
	}
}

func readMetricFromPath(metricPath string) (value int64) {
	dataBytes, err := ioutil.ReadFile(metricPath)
	handleErr(err, fmt.Sprintf("Failed to read metric at %v", metricPath))

	data, err := strconv.Atoi(strings.TrimSuffix(string(dataBytes), "\n"))
	handleErr(err, fmt.Sprintf("Failed to convert metric value at %v to int", metricPath))

	value = int64(data)
	return
}

func createMetric(metricName string, metricInfo []MetricInfo) {
	meter := global.MeterProvider().Meter("otel-stats")
	metric.Must(meter).NewInt64UpDownSumObserver(metricName,
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

func createOtelMetricsForStatsfs(statsfsPath string) {
	m := CreateStatsfsMetrics(statsfsPath)

	for _, subsysMetrics := range m.Metrics {
		for metricName, metricInfo := range subsysMetrics.Metrics {
			createMetric(metricName, metricInfo)
		}
	}
}
