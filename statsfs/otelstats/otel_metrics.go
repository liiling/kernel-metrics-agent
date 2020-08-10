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
func InitOtelPipeline(statsfsPath string, subsystemName string) {
	fmt.Println("In otel-metrics!")
	exporter := InitExporter()
	if exporter != nil {
		defer exporter.Stop()
	}
	CreateOtelMetricsForSubsys(statsfsPath, subsystemName)

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
					readMetricFromPath(info.MetricPath),
					kv.String("device", info.Label),
				)
			}
		},
		metric.WithDescription(metricName),
	)
}

func CreateOtelMetricsForSubsys(statsfsPath string, subsystemName string) {
	m := CreateSubsysMetrics(statsfsPath, subsystemName)

	for metricName, metricInfo := range m.Metrics {
		createMetric(metricName, metricInfo)
	}
}
