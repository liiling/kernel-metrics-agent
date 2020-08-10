package otelstats

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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

func initExporter(exporterName string) (exporter *push.Controller) {
	if exporterName == "stdout" {
		exporter = initStdoutExporter()
	} else {
		err := errors.New("Invalid exporter name")
		handleErr(err, fmt.Sprintf("Exporter name %v is not allowed.", exporterName))
	}
	return
}

func parseExporterName() string {
	exporterName := flag.String("exporter", "prometheus", "Exporter to use. Choose between prometheus (default), stdout and gcp")
	flag.Parse()

	fmt.Printf("Exporter: %v\n", *exporterName)
	return *exporterName
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

// InitExporter initialise gcp, stdout or prometheus exporter based on cmd flag
func InitExporter() *push.Controller {
	fmt.Println("In exporters.go!")
	exporterName := parseExporterName()
	if exporterName == "prometheus" {
		initPrometheusExporter()
		return nil
	}

	exporter := initExporter(exporterName)
	return exporter
}
