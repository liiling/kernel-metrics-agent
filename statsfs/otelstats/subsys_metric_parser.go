package otelstats

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// MetricInfo contains a Label used to identify the specific device
// and a Path from where the metric for this device could be retrieved
type MetricInfo struct {
	Name  string
	Label string
	Path  string
}

// SubsysMetrics is a struct that represents metrics of a subsystem.
// Path gives the base path to the subsystem's stats in statsfs
// (usually /sys/kernel/stats/subsystemName), Metrics is a map with key
// being the metric name, and value being a list of labels (devices with
// the metric registered)
type SubsysMetrics struct {
	StatsfsPath   string
	SubSystemName string
	SubSystemPath string
	Metrics       map[string][]MetricInfo
}

// newSubsysMetric creates a SubsysMetric struct given the mounting
// point of statsfs filesystem (statsfsPath) and the subsystemName
func newSubsysMetric(statsfsPath, subsystemName string) SubsysMetrics {
	subsysMetrics := SubsysMetrics{
		StatsfsPath:   statsfsPath,
		SubSystemName: subsystemName,
		SubSystemPath: strings.Join([]string{statsfsPath, subsystemName}, "/"),
		Metrics:       make(map[string][]MetricInfo),
	}
	subsysMetrics.constructMetricMap()
	return subsysMetrics
}

func (m *SubsysMetrics) constructMetricMap() {
	err := filepath.Walk(m.SubSystemPath, m.updateMetricMap)
	if err != nil {
		log.Printf("Failed to parse metrics for subsystem %v at %v", m.SubSystemName, m.SubSystemPath)
	}
}

func (m *SubsysMetrics) updateMetricMap(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Failed to walk to file %v\n", path)
	}

	if info.IsDir() {
		return nil
	}
	m.updateMetricMapOneEntry(path)
	return nil
}

func (m *SubsysMetrics) updateMetricMapOneEntry(path string) {
	metricInfo := m.getMetricInfo(path)
	m.Metrics[metricInfo.Name] = append(m.Metrics[metricInfo.Name], metricInfo)
}

// Given a path to a statsfs file, return a MetricInfo struct with label
// computed by getMetricLabel method and Path being the input path
// (the path from where the metric could be retrieved)
// Example:
//	Input:
//		m.SubSystemPath = /sys/kernel/stats/net
//		path = /sys/kernel/stats/net/eth0/sub0/latency
//	Output:
//		MetricInfo{
//			Name: net/latency
//			Label: /eth0/sub0
//			Path: /sys/kernel/stats/net/eth0/sub0/latency
//		}
func (m *SubsysMetrics) getMetricInfo(path string) MetricInfo {
	return MetricInfo{
		Name:  m.getMetricName(path),
		Label: m.getMetricLabel(path),
		Path:  path,
	}
}

func (m *SubsysMetrics) getMetricName(path string) string {
	segs := strings.Split(path, "/")
	metricFileName := segs[len(segs)-1]
	metricName := strings.Join([]string{m.SubSystemName, metricFileName}, "/")
	return metricName
}

func (m *SubsysMetrics) getMetricLabel(path string) string {
	metricStr := strings.Split(path, m.SubSystemPath)[1]
	labelSeg := strings.Split(metricStr, "/")
	label := strings.Join(labelSeg[:len(labelSeg)-1], "/")
	return label
}

func (m *SubsysMetrics) print() {
	fmt.Printf("StatsfsPath: %v\n", m.StatsfsPath)
	fmt.Printf("SubSystemName: %v\n", m.SubSystemName)
	fmt.Printf("SubSystemPath: %v\n", m.SubSystemPath)
	fmt.Println("Metrics:")
	for metricName, labels := range m.Metrics {
		fmt.Printf("\tmetricName: %v,\n\tinfo: \n", metricName)
		for _, label := range labels {
			fmt.Printf("\t\tLabel: %v, Path: %v\n", label.Label, label.Path)
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
		metrics.Metrics[subsystemName] = newSubsysMetric(statsfsPath, subsystemName)
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
