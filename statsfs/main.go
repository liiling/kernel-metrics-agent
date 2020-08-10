package main

import (
	"otelstats"
)

func main() {
	// otelstats.InitOtelPipeline("./stubsys/kernel/stats")
	otelstats.CreateSubsysMetrics("stubsys/kernel/stats", "net")
}
