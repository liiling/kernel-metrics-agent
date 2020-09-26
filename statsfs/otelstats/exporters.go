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

const (
	prometheusPort = 2112
	pullPeriod = 10 * time.Second
	pushPeriod = 10 * time.Second
)

func initPrometheusExporter() error {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{}, pull.WithCachePeriod(pullPeriod))
	if err != nil {
		return fmt.Errorf("Failed to initialize Prometheus metric exporter: %v", err)
	}

	http.HandleFunc("/metrics", exporter.ServeHTTP)
	go http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
	fmt.Printf("Prometheus server running on :%d\n", prometheusPort)
	global.SetMeterProvider(exporter.Provider())
	return nil
}

func initStdoutExporter() (*push.Controller, error) {
	exportOpts := []stdout.Option{stdout.WithPrettyPrint()}
	pushOpts := []push.Option{push.WithPeriod(pushPeriod)}
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
	exporter, err := mexporter.NewExportPipeline(opts, push.WithPeriod(pushPeriod))
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
		return nil, initPrometheusExporter()
	case "stdout":
		return initStdoutExporter()
	case "gcp":
		return initGCPExporter()
	default:
		return nil, fmt.Errorf("invalid exporter name %v", exporterName)
	}
}
