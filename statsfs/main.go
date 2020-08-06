package main

import (
	"otelstats"
)

func main() {
	otelstats.WalkDir("./stubsys")
}
