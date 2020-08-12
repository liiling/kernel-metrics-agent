package otelstats

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func initPrometheusExporter() {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{}, pull.WithCachePeriod(time.Second*10))
	handleErr(err, "Failed to initialize Prometheus metric exporter")

	port := 2112
	http.HandleFunc("/metrics", exporter.ServeHTTP)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Printf("Prometheus server running on :%d\n", port)
	global.SetMeterProvider(exporter.Provider())
}

func initStdoutExporter() *push.Controller {
	exportOpts := []stdout.Option{stdout.WithPrettyPrint()}
	pushOpts := []push.Option{push.WithPeriod(time.Second * 10)}
	_, exporter, err := stdout.NewExportPipeline(
		exportOpts,
		pushOpts,
	)
	handleErr(err, "Failed to initialize Stdout metric exporter")
	global.SetMeterProvider(exporter.Provider())
	return exporter
}

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

// InitExporter initialise gcp, stdout or prometheus exporter based on cmd flag
func InitExporter(exporterName string) (exporter *push.Controller) {
	if exporterName == "prometheus" {
		initPrometheusExporter()
		return nil
	}
	if exporterName == "stdout" {
		exporter = initStdoutExporter()
	} else if exporterName == "gcp" {
		exporter = initGCPExporter()
	} else {
		err := errors.New("Invalid exporter name")
		handleErr(err, fmt.Sprintf("Exporter name %v is not allowed.", exporterName))
	}
	return
}