package otelstats

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
)

// InitOtelPipeline initializes an OpenTelemetry pipeline
// that crawls a user defined path and exports all the available
// stats to a backend of choice (gcp, stdout, prometheus)
func InitOtelPipeline(path string) {
	fmt.Println("In otel-metrics!")
	exporter := InitExporter()
	if exporter != nil {
		defer exporter.Stop()
	}
	WalkDir(path)
	for {
	}
}

func walkDirHelper(path string, info os.FileInfo, err error) error {
	handleErr(err, fmt.Sprintf("Failed to walk to file %v", path))
	if !info.IsDir() {
		fmt.Println(path)
		createMetric(path, path)
	} else {
		fmt.Printf("Directory: %v\n", path)
	}
	return err
}

// WalkDir walks through all files in the given directory dir
// and performs action defined by the given WalkFunc
func WalkDir(dir string) {
	err := filepath.Walk(dir, walkDirHelper)
	handleErr(err, fmt.Sprintf("Failed to walk dir %v", dir))
}

func createMetric(path string, desc string) {
	meter := global.MeterProvider().Meter("otel-stats")
	devName, statsName := filepath.Split(path)
	label := []kv.KeyValue{kv.String("device", devName), kv.String("stats", statsName)}

	metric.Must(meter).NewInt64ValueObserver(path,
		func(_ context.Context, result metric.Int64ObserverResult) {
			data, err := ioutil.ReadFile(path)
			handleErr(err, fmt.Sprintf("Failed to read file %v for metric update", path))

			dataNum, err := strconv.Atoi(strings.TrimSuffix(string(data), "\n"))
			handleErr(err, fmt.Sprintf("Failed to convert metric %v to int", path))

			result.Observe(int64(dataNum), label...)
		},
		metric.WithDescription(desc),
	)
}
