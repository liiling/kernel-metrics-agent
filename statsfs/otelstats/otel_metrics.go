package otelstats

import "fmt"

// InitOtelPipeline initializes an OpenTelemetry pipeline
// that crawls /sys/kernel/stats and exports all the available
// stats to a backend of choice (gcp, stdout, prometheus)
func InitOtelPipeline() {
	fmt.Println("In otel-metrics!")
	InitExporter()
}
