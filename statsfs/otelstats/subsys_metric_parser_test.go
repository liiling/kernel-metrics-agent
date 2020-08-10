package otelstats

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInitSubsysMetricStruct(t *testing.T) {
	statsfsPath := "/sys/kernel/stats"
	subsystemName := "subsystem"
	actual := initSubsysMetricStruct(statsfsPath, subsystemName)

	expected := SubsysMetrics{
		StatsfsPath:   "/sys/kernel/stats",
		SubSystemName: "subsystem",
		SubSystemPath: "/sys/kernel/stats/subsystem",
		Metrics:       make(map[string][]MetricInfo),
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("initSubsysMetricStruct mismatch (-expected +actual):\n%s", diff)
	}
}

func TestGetMetricPath(t *testing.T) {
	subsysMetric := SubsysMetrics{
		SubSystemPath: "/sys/kernel/stats/subsystem",
	}
	inputPaths := []string{
		"/sys/kernel/stats/subsystem/metrics",
		"/sys/kernel/stats/subsystem/device/metrics",
		"/sys/kernel/stats/subsystem/device/subdevice/metrics",
	}
	expected := []string{
		"/metrics",
		"/device/metrics",
		"/device/subdevice/metrics",
	}

	for i, inputPath := range inputPaths {
		actual := subsysMetric.getMetricPath(inputPath)
		if diff := cmp.Diff(expected[i], actual); diff != "" {
			t.Errorf("getMetricPath mismatch \ninput path = %s,(-expected +actual):\n%s", inputPath, diff)
		}
	}
}

func TestGetMetricNameAndLabel(t *testing.T) {
	subsysMetric := SubsysMetrics{
		SubSystemName: "subsystem",
	}
	metricPaths := []string{
		"/metrics",
		"/device/metrics",
		"/device/subdevice/metrics",
	}

	expectedMetricName := "subsystem/metrics"
	expectedLabels := []string{
		"",
		"/device",
		"/device/subdevice",
	}

	for i, path := range metricPaths {
		actualMetricName, actualLabel := subsysMetric.getMetricNameAndLabel(path)
		if metricNameDiff := cmp.Diff(expectedMetricName, actualMetricName); metricNameDiff != "" {
			t.Errorf("getMetricNameAndLabel mismatch on metric name\ninput path = %s,(-expected +actual):\n%s", path, metricNameDiff)
		}
		if labelDiff := cmp.Diff(expectedLabels[i], actualLabel); labelDiff != "" {
			t.Errorf("getMetricNameAndLabel mismatch on label\ninput path = %s,(-expected +actual):\n%s", path, labelDiff)
		}
	}
}
