package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
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
	exporter, err := stdout.NewExportPipeline(stdout.Config{
		PrettyPrint: false},
		push.WithPeriod(time.Second*10),
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
func initExporter(exporterName string) (exporter *push.Controller) {
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

func createIntValueObserver(name string, desc string) metric.Int64ValueObserver {
	meter := global.Meter("otel-switch-backend")
	observer := metric.Must(meter).NewInt64ValueObserver(name,
		updateMetric,
		metric.WithDescription(desc),
	)
	return observer
}

func getVisitCounter() int64 {
	resp, err := http.Get("http://localhost:8090/getVisitCounter")
	handleErr(err, "Failed to issue GET request to localhost:8090/getVisitCounter")
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	handleErr(err, "Failed to read response body")

	numVisits, err := strconv.Atoi(string(respBody))
	handleErr(err, "Failed to convert response to int")

	return int64(numVisits)
}

func updateMetric(_ context.Context, result metric.Int64ObserverResult) {
	numVisits := getVisitCounter()
	fmt.Printf("Updated number of visits: %v\n", numVisits)
	result.Observe(numVisits, kv.String("label", "test"))
}

func main() {
	exporterName := parseExporterName()
	if exporterName == "prometheus" {
		initPrometheusExporter()
	} else {
		exporter := initExporter(exporterName)
		defer exporter.Stop()
	}

	createIntValueObserver("visit-observer",
		"A value observer representing number of times a website is visited.")

	for {
	}
}
