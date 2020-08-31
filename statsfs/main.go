package main

import (
	"flag"
	"fmt"
	"otelstats"
)

func parseFlags() (string, string) {
	exporterName := flag.String("exporter", "prometheus", "Exporter to use. Choose between prometheus (default), stdout and gcp")
	statsfsPath := flag.String("statsfspath", "otelstats/testsys/kernel/stats", "Path to statsfs filesystem for metrics (default to the teset filesystem)")
	flag.Parse()

	fmt.Printf("Exporter: %v\nStatsfs path: %v\n", *exporterName, *statsfsPath)
	return *exporterName, *statsfsPath
}

func main() {
	exporterName, statsfsPath := parseFlags()
	otelstats.InitOtelPipeline(exporterName, statsfsPath)
}
