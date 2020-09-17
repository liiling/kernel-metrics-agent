package main

import (
	"flag"
	"log"
	"otelstats"
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

	err = otelstats.CreateOtelMetricsForStatsfs(*statsfsPath)
	if err != nil {
		log.Panic(err)
	}

	select { }
}
