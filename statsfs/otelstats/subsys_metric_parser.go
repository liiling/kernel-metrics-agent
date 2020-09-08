package otelstats

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Given a path to a statsfs file, return the metricName.
// path = {m.StatsfsPath}/{m.SubSystemName}/{device paths(if any)}/{metricFileName}
// metricName = {m.SubSystemName}/{metricFilename}
// Example:
//	Input:
//		path = /sys/kernel/stats/net/eth0/latency
//			where
//				m.StatsfsPath = /sys/kernel/stats
//				m.SubsystemName = net
//				device paths = eth0
//				metricFileName = latency
//	Output:
//		metricName = net/latency
func (m SubsysMetrics) getMetricName(path string) string {
	segs := strings.Split(path, "/")
	metricFileName := segs[len(segs)-1]
	metricName := strings.Join([]string{m.SubSystemName, metricFileName}, "/")
	return metricName
}

// Given a path to a statsfs file, return the label of the metric.
// path = {m.SubSystemPath}/{device paths(if any)}/{metricFileName}
// label = {device paths (if any)}
// Example:
//	Input:
//		path = /sys/kernel/stats/net/eth0/latency
//			where
//				m.StatsfsPath = /sys/kernel/stats
//				m.SubSystemName = net
//				m.SubSystemPath = /sys/kernel/stats/net
//				device paths = /eth0/sub0
//				metricFileName = latency
//	Output:
// 		label = /eth0/sub0
func (m SubsysMetrics) getMetricLabel(path string) string {
	metricStr := strings.Split(path, m.SubSystemPath)[1]
	labelSeg := strings.Split(metricStr, "/")
	label := strings.Join(labelSeg[:len(labelSeg)-1], "/")
	return label
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
//			Label: /eth0/sub0
//			Path: /sys/kernel/stats/net/eth0/sub0/latency
//		}
func (m SubsysMetrics) getMetricInfo(path string) MetricInfo {
	label := m.getMetricLabel(path)
	metricInfo := MetricInfo{Label: label, Path: path}
	return metricInfo
}

// Given a path to a statsfs file, update the metricMap of the corresponding
// SubsysMetric stuct
// Example:
//	Input:
//		m.SubSystemPath = /sys/kernel/stats/net
// 		m.SubSystemName = net
//		m.StatsfsPath = /sys/kernel/stats
//		path = /sys/kernel/stats/net/eth0/latency
//	Output:
//		metricName = net/latency
//		label = /eth0
//		new entry in m.Metrics[net/latency]:
//			MetricInfo{Label: /eth0, Path: /sys/kernel/stats/net/eth0/latency}
func (m SubsysMetrics) updateMetricMapOneEntry(path string) {
	metricName := m.getMetricName(path)
	metricInfo := m.getMetricInfo(path)
	m.Metrics[metricName] = append(m.Metrics[metricName], metricInfo)
}

func (m SubsysMetrics) updateMetricMap(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Failed to walk to file %v\n", path)
	}

	if !info.IsDir() {
		m.updateMetricMapOneEntry(path)
	}
	return err
}

func (m SubsysMetrics) constructMetricMap() {
	err := filepath.Walk(m.SubSystemPath, m.updateMetricMap)
	if err != nil {
		log.Printf("Failed to parse metrics for subsystem %v at %v", m.SubSystemName, m.SubSystemPath)
	}
}

func (m SubsysMetrics) print() {
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

// Print prints StatsfsMetrics struct
func (m StatsfsMetrics) Print() {
	fmt.Print("\n####################################\n")
	fmt.Printf("StatsfsPath: %v\n\n", m.StatsfsPath)
	for subsysName, subsysMetrics := range m.Metrics {
		fmt.Println("------------------")
		fmt.Printf("Statsfs metrics for subsystem %v:\n\n", subsysName)
		subsysMetrics.print()
	}
	fmt.Printf("####################################\n\n")
}

// MetricInfo contains a Label used to identify the specific device
// and a Path from where the metric for this device could be retrieved
type MetricInfo struct {
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

// StatsfsMetrics is a struct that represents metrics available in a statsfs
// filesystem found at StatsfsPath
// Each subsystem metrics is represented with a SubsysMetrics struct
type StatsfsMetrics struct {
	StatsfsPath string
	Metrics     map[string]SubsysMetrics
}

func initSubsysMetricStruct(statsfsPath, subsystemName string) SubsysMetrics {
	subsysMetrics := SubsysMetrics{
		StatsfsPath:   statsfsPath,
		SubSystemName: subsystemName,
		SubSystemPath: strings.Join([]string{statsfsPath, subsystemName}, "/"),
		Metrics:       make(map[string][]MetricInfo),
	}
	return subsysMetrics
}

// createSubsysMetrics creates a SubsysMetric struct given the mounting
// point of statsfs filesystem (statsfsPath) and the subsystemName
func createSubsysMetrics(statsfsPath, subsystemName string) SubsysMetrics {
	subsysMetrics := initSubsysMetricStruct(statsfsPath, subsystemName)
	subsysMetrics.constructMetricMap()
	return subsysMetrics
}

// CreateStatsfsMetrics creates a StatsfsMetrics struct given the mounting
// point of statsfs filesystem (statsfsPath)
func CreateStatsfsMetrics(statsfsPath string) (*StatsfsMetrics, error) {
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
		metrics.Metrics[subsystemName] = createSubsysMetrics(statsfsPath, subsystemName)
	}
	return &metrics, nil
}
