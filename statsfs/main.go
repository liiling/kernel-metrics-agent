package main

import (
	"otelstats"
)

func main() {
	otelstats.InitOtelPipeline("otelstats/testsys/kernel/stats")
}
