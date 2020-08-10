package otelstats

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// path = m.SubSystemPath + metricPath
// Example:
//	Input:
//		path = /sys/kernel/stats/net/eth0/latency,
//		m.SubsystemPath = /sys/kernel/stats/net
//	Output:
//		metricPath = /eth0/latency
func (m SubsysMetrics) getMetricPath(path string) string {
	segs := strings.Split(path, m.SubSystemPath)
	return segs[1]
}

// parse metricPath into metricName and label
// metricPath = {label}/metricFileName
// metricName = m.SubSystemName/metricfileName (last part of metricPath)
// Example:
// 	Input:
//		m.SubSystemName = net
//		metricPath = /eth0/latency
//	Output:
//		metricName = net/latency
//		label = /eth0
func (m SubsysMetrics) getMetricNameAndLabel(metricPath string) (metricName, label string) {
	segs := strings.Split(metricPath, "/")
	lastIdx := len(segs) - 1

	metricName = strings.Join([]string{m.SubsystemName, segs[lastIdx]}, "/")
	label = strings.Join(segs[:lastIdx], "/")
	return
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
//		metricPath = /eth0/latency
//		metricName = net/latency
//		label = /eth0
//		new entry in m.Metrics[net/latency]:
//			MetricInfo{Label: /eth0, MetricPath: /sys/kernel/stats/net/eth0/latency}
func (m SubsysMetrics) updateMetricMapOneEntry(path string) {
	metricPath := m.getMetricPath(path)
	metricName, label := m.getMetricNameAndLabel(metricPath)
	metricInfo := MetricInfo{Label: label, MetricPath: path}
	m.Metrics[metricName] = append(m.Metrics[metricName], metricInfo)
}

func (m SubsysMetrics) updateMetricMap(path string, info os.FileInfo, err error) error {
	handleErr(err, fmt.Sprintf("Failed to walk to file %v", path))

	if !info.IsDir() {
		m.updateMetricMapOneEntry(path)
	}
	return err
}

func (m SubsysMetrics) constructMetricMap() {
	err := filepath.Walk(m.SubSystemPath, m.updateMetricMap)
	handleErr(err, fmt.Sprintf("Failed to parse metrics for subsystem %v at %v", m.SubsystemName, m.SubSystemPath))
}

func (m SubsysMetrics) print() {
	fmt.Println("------------------")
	fmt.Printf("StatsfsPath: %v\n", m.StatsfsPath)
	fmt.Printf("SubsystemName: %v\n", m.SubsystemName)
	fmt.Printf("SubsystemPath: %v\n", m.SubSystemPath)
	fmt.Println("Metrics:")
	for metricName, labels := range m.Metrics {
		fmt.Printf("\tmetricName: %v, info: %v\n", metricName, labels)
	}
	fmt.Println("------------------")
}

// MetricInfo contains a Label used to identify the specific device
// and a MetricPath from where the metric for this device could be retrieved
type MetricInfo struct {
	Label      string
	MetricPath string
}

// SubsysMetrics is a struct that represents metrics of a subsystem.
// Path gives the base path to the subsystem's stats in statsfs
// (usually /sys/kernel/stats/subsystemName), Metrics is a map with key
// being the metric name, and value being a list of labels (devices with
// the metric registered)
type SubsysMetrics struct {
	StatsfsPath   string
	SubsystemName string
	SubSystemPath string
	Metrics       map[string][]MetricInfo
}

// CreateSubsysMetrics creates a SubsysMetric struct given the mounting
// point of statsfs filesystem (statsfsPath) and the subsystemName
func CreateSubsysMetrics(statsfsPath string, subsystemName string) (subsysMetrics SubsysMetrics) {
	subsysMetrics = SubsysMetrics{
		StatsfsPath:   statsfsPath,
		SubsystemName: subsystemName,
		SubSystemPath: strings.Join([]string{statsfsPath, subsystemName}, "/"),
		Metrics:       make(map[string][]MetricInfo),
	}
	subsysMetrics.constructMetricMap()
	subsysMetrics.print()
	return
}
