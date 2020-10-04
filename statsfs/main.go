package main

import (
	"flag"
	"log"
	"github.com/liiling/kernel-metrics-agent/statsfs/otelstats"
)

var (
	exporterName = flag.String("exporter", "prometheus", "Exporter to use. Choose between prometheus (default), stdout and gcp")
	statsfsPath  = flag.String("statsfspath", "otelstats/testsys/kernel/stats", "Path to statsfs filesystem for metrics (default to the teset filesystem)")
)

func main() {
	flag.Parse()
	exporter, err := otelstats.InitExporter(*exporterName)
	if err != nil {
		log.Panicf("Failed to initialize exporter %v: %v\n", *exporterName, err)
	}
	if exporter != nil {
		defer exporter.Stop()
	}

	if err = otelstats.CreateOtelMetricsForStatsfs(*statsfsPath); err != nil {
		log.Panic(err)
	}

	// block forever since the exporter is pulling/pushing metrics periodically
	// in a separate goroutine
	select {}
}
