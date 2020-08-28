package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
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

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func createIntSumCounter() metric.Int64Counter {
	meter := global.Meter("test")
	counter := metric.Must(meter).NewInt64Counter("counter",
		metric.WithDescription("A counter that increments every 10 seconds"))
	return counter
}

func incrementCounterMetrics(ctx context.Context, counter metric.BoundInt64Counter) {
	for i := 1; ; {
		counter.Add(ctx, int64(i))
		time.Sleep(time.Second)
	}
}

func main() {
	initPrometheusExporter()

	counter := createIntSumCounter().Bind(label.String("label1", "value1"))
	defer counter.Unbind()

	ctx := context.Background()
	fmt.Println("Start incrementing counter...")
	incrementCounterMetrics(ctx, counter)
}
