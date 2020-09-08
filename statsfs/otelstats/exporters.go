package otelstats

import (
	"fmt"
	"net/http"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func initPrometheusExporter() error {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{}, pull.WithCachePeriod(time.Second*10))
	if err != nil {
		return fmt.Errorf("Failed to initialize Prometheus metric exporter: %v", err)
	}

	port := 2112
	http.HandleFunc("/metrics", exporter.ServeHTTP)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Printf("Prometheus server running on :%d\n", port)
	global.SetMeterProvider(exporter.Provider())
	return nil
}

func initStdoutExporter() (*push.Controller, error) {
	exportOpts := []stdout.Option{stdout.WithPrettyPrint()}
	pushOpts := []push.Option{push.WithPeriod(time.Second * 10)}
	_, exporter, err := stdout.NewExportPipeline(
		exportOpts,
		pushOpts,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize stdout metric exporter: %v", err)
	}
	global.SetMeterProvider(exporter.Provider())
	return exporter, nil
}

func initGCPExporter() (*push.Controller, error) {
	opts := []mexporter.Option{}
	// Minimum interval for GCP exporter is 10s
	exporter, err := mexporter.NewExportPipeline(opts, push.WithPeriod(time.Second*10))
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize GCP metric exporter: %v", err)
	}

	global.SetMeterProvider(exporter.Provider())
	return exporter, nil
}

// InitExporter initialise gcp, stdout or prometheus exporter based on cmd flag
func InitExporter(exporterName string) (*push.Controller, error) {
	switch exporterName {
	case "prometheus":
		err := initPrometheusExporter()
		return nil, err
	case "stdout":
		return initStdoutExporter()
	case "gcp":
		return initGCPExporter()
	default:
		return nil, fmt.Errorf("invalid exporter name %v", exporterName)
	}
}
