package otelstats

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitOtelPipeline initializes an OpenTelemetry pipeline
// that crawls /sys/kernel/stats and exports all the available
// stats to a backend of choice (gcp, stdout, prometheus)
func InitOtelPipeline() {
	fmt.Println("In otel-metrics!")
	InitExporter()
}

func walkDirHelper(path string, info os.FileInfo, err error) error {
	handleErr(err, fmt.Sprintf("Failed to walk to file %v", path))
	if !info.IsDir() {
		fmt.Println(path)
	}
	return err
}

// WalkDir walks through all files in the given directory dir
// and performs action defined by the given WalkFunc
func WalkDir(dir string) {
	err := filepath.Walk(dir, walkDirHelper)
	handleErr(err, fmt.Sprintf("Failed to walk dir %v", dir))
}
