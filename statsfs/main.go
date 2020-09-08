package main

import (
	"flag"
	"otelstats"
)

var (
	exporterName = flag.String("exporter", "prometheus", "Exporter to use. Choose between prometheus (default), stdout and gcp")
	statsfsPath  = flag.String("statsfspath", "otelstats/testsys/kernel/stats", "Path to statsfs filesystem for metrics (default to the teset filesystem)")
)

func main() {
	flag.Parse()
	otelstats.InitOtelPipeline(*exporterName, *statsfsPath)
}
