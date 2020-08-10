package otelstats

import (
	"fmt"
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
	for {
	}
}
