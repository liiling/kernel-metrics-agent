package otelstats

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const schemaFilename = ".schema"

// SubsysMetrics is a struct that represents metrics of a subsystem.
// SubsystemPath gives the base path to the subsystem's stats in statsfs
// (usually /sys/kernel/stats/subsystemName), Metrics is a map with key
// being the metric name, and value being the associated MetricInfo
type SubsysMetrics struct {
	StatsfsPath   string
	SubSystemName string
	SubSystemPath string
	Metrics       map[string]MetricInfo
}

// MetricInfo is a struct that represents a metric from a subsystem.
// A metric contains name, flag (CUMULATIVE or GAUGE), type (INT or FLAT),
// description and a map with key = path to the metric file, value =
// the associated list of metricLabels
type MetricInfo struct {
	Name        string
	Flag        string
	Type        string
	Desc        string
	PathToLabel map[string][]metricLabel
}

// newSubsysMetric creates a SubsysMetric struct given the mounting
// point of statsfs filesystem (statsfsPath) and the subsystemName
func newSubsysMetric(statsfsPath, subsystemName string) (*SubsysMetrics, error) {
	m := SubsysMetrics{
		StatsfsPath:   statsfsPath,
		SubSystemName: subsystemName,
		SubSystemPath: strings.Join([]string{statsfsPath, subsystemName}, "/"),
		Metrics:       make(map[string]MetricInfo),
	}

	if err := filepath.Walk(m.SubSystemPath, m.updateMetricMap); err != nil {
		return nil, fmt.Errorf("failed to parse metrics for subsystem %v at %v: %v", m.SubSystemName, m.SubSystemPath, err)
	}
	return &m, nil
}

func (m *SubsysMetrics) updateMetricMap(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("failed to walk to file %v", path)
	}

	if dirname, filename := filepath.Split(path); filename == schemaFilename {
		if metricSchemas, err := parseSchema(path); err != nil {
			return fmt.Errorf("failed to parse .schema file at %v", path)
		} else {
			fmt.Printf("schema: %v\n", metricSchemas)
			for _, metricSchema := range metricSchemas {
				metricPath := filepath.Join(dirname, metricSchema.mname)
				// update metric schema
				if schema, ok := m.Metrics[metricSchema.mname]; ok {
					fmt.Printf("labels: %v\n", metricSchema.mlabels)
					schema.PathToLabel[metricPath] = metricSchema.mlabels
				} else {
					m.Metrics[metricSchema.mname] = MetricInfo{
						Name:        metricSchema.mname,
						Flag:        metricSchema.mflag,
						Type:        metricSchema.mtype,
						Desc:        metricSchema.mdesc,
						PathToLabel: map[string][]metricLabel{metricPath: metricSchema.mlabels},
					}
				}
			}
		}
	}
	return nil
}

func (m *SubsysMetrics) print() {
	fmt.Printf("StatsfsPath: %v\n", m.StatsfsPath)
	fmt.Printf("SubSystemName: %v\n", m.SubSystemName)
	fmt.Printf("SubSystemPath: %v\n", m.SubSystemPath)
	fmt.Println("Metrics:")
	for metricName, info := range m.Metrics {
		fmt.Printf("\tmetric: %v\n", metricName)
		fmt.Printf("\t\tname: %v, flag: %v, type: %v, desc: %v\n", info.Name, info.Flag, info.Type, info.Desc)
		fmt.Printf("\t\tPath to labels:\n")
		for path, labels := range info.PathToLabel {
			fmt.Printf("\t\t\tpath: %v, labels: %p\t\n", path, labels)
		}
	}
}

// StatsfsMetrics is a struct that represents metrics available in a statsfs
// filesystem found at StatsfsPath
// Each subsystem metrics is represented with a SubsysMetrics struct
type StatsfsMetrics struct {
	StatsfsPath string
	Metrics     map[string]SubsysMetrics
}

// NewStatsfsMetrics creates a StatsfsMetrics struct given the mounting
// point of statsfs filesystem (statsfsPath)
func NewStatsfsMetrics(statsfsPath string) (*StatsfsMetrics, error) {
	metrics := StatsfsMetrics{
		StatsfsPath: statsfsPath,
		Metrics:     make(map[string]SubsysMetrics),
	}
	statsfsDir, err := os.Open(statsfsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open statsfs dir at %v: %v", statsfsPath, err)
	}
	defer statsfsDir.Close()

	subsystemNames, err := statsfsDir.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read dirnames from statsfs dir at %v: %v", statsfsPath, err)
	}

	for _, subsystemName := range subsystemNames {
		if subsysMetric, err := newSubsysMetric(statsfsPath, subsystemName); err != nil {
			log.Printf("failed to generate metrics for subsystem %v: %v\n", subsystemName, err)
		} else {
			metrics.Metrics[subsystemName] = *subsysMetric
		}
	}
	return &metrics, nil
}

// Print prints StatsfsMetrics struct
func (m *StatsfsMetrics) Print() {
	fmt.Print("\n####################################\n")
	fmt.Printf("StatsfsPath: %v\n\n", m.StatsfsPath)
	for subsysName, subsysMetrics := range m.Metrics {
		fmt.Println("------------------")
		fmt.Printf("Statsfs metrics for subsystem %v:\n\n", subsysName)
		subsysMetrics.print()
	}
	fmt.Printf("####################################\n\n")
}
