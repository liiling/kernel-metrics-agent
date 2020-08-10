package main

import (
	"otelstats"
)

func main() {
	otelstats.InitOtelPipeline("stubsys/kernel/stats", "net")
}
