package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func initGCPExporter() *push.Controller {
	opts := []mexporter.Option{}
	// Minimum interval for GCP exporter is 10s
	exporter, err := mexporter.NewExportPipeline(opts, push.WithPeriod(time.Second*10))
	handleErr(err, "Failed to initialize metric exporter")

	global.SetMeterProvider(exporter.Provider())

	return exporter
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func updateMetric(_ context.Context, result metric.Int64ObserverResult) {
	fmt.Println("Set metric to 1")
	result.Observe(1, kv.String("label", "test"))
}

func parseAsyncMetricType() (bool, string) {
	sync := flag.Bool("sync", true, "Whether the synchronous metrics are used")
	metricType := flag.String("mtype", "sum", "sum, updown, value")
	flag.Parse()

	fmt.Printf("sync: %v, metric type used: %v\n", *sync, *metricType)
	return *sync, *metricType
}

func main() {
	sync, metricType := parseAsyncMetricType()

	exporter := initGCPExporter()
	defer exporter.Stop()

	meter := global.Meter("gcp-asynctest")

	if sync {
		if metricType == "sum" {
			fmt.Println("Creating an Int64Counter...")
			syncInstrument := metric.Must(meter).NewInt64Counter("Counter")
			syncInstrument.Add(context.Background(), 1)
		} else if metricType == "updown" {
			fmt.Println("Creating an Int64UpDownCounter...")
			syncInstrument := metric.Must(meter).NewInt64UpDownCounter("UpDownCounter")
			syncInstrument.Add(context.Background(), 1)
		} else if metricType == "value" {
			fmt.Println("Creating an Int64ValueRecorder...")
			syncInstrument := metric.Must(meter).NewInt64ValueRecorder("ValueRecorder")
			syncInstrument.Record(context.Background(), 1)
		}
	} else {
		if metricType == "sum" {
			fmt.Println("Creating an Int64SumObserver...")
			metric.Must(meter).NewInt64SumObserver("SumObserver", updateMetric)
		} else if metricType == "updown" {
			fmt.Println("Creating an Int64UpDownSumObserver...")
			metric.Must(meter).NewInt64UpDownSumObserver("UpDownSumObserver", updateMetric)
		} else if metricType == "value" {
			fmt.Println("Creating an Int64ValueObserver...")
			metric.Must(meter).NewInt64ValueObserver("ValueObserver", updateMetric)
		} else {
			fmt.Printf("Invalid metric type %v..\n", metricType)
		}
	}
}
