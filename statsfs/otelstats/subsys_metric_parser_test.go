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

func TestGetMetricName(t *testing.T) {
	subsysMetric := SubsysMetrics{
		SubSystemName: "subsystem",
	}

	paths := []string{
		"/sys/kernel/stats/subsystem/metrics",
		"/sys/kernel/stats/subsystem/device/metrics",
		"/sys/kernel/stats/subsystem/device/subdevice/metrics",
	}
	expectedMetricName := "subsystem/metrics"

	for _, path := range paths {
		actualMetricName := subsysMetric.getMetricName(path)
		if diff := cmp.Diff(expectedMetricName, actualMetricName); diff != "" {
			t.Errorf("getMetricName mismatch on input path = %s,(-expected +actual):\n%s", path, diff)
		}
	}
}

func TestGetMetricLabel(t *testing.T) {
	subsysMetric := SubsysMetrics{
		SubSystemPath: "/sys/kernel/stats/subsystem",
	}

	paths := []string{
		"/sys/kernel/stats/subsystem/metrics",
		"/sys/kernel/stats/subsystem/device/metrics",
		"/sys/kernel/stats/subsystem/device/subdevice/metrics",
	}
	expectedLabels := []string{
		"",
		"/device",
		"/device/subdevice",
	}

	for i, path := range paths {
		actualLabel := subsysMetric.getMetricLabel(path)
		if diff := cmp.Diff(expectedLabels[i], actualLabel); diff != "" {
			t.Errorf("getMetricLabel mismatch on input path = %s,(-expected +actual):\n%s", path, diff)
		}
	}
}

func TestGetMetricInfo(t *testing.T) {
	subsysMetric := SubsysMetrics{
		SubSystemPath: "/sys/kernel/stats/subsystem",
	}

	paths := []string{
		"/sys/kernel/stats/subsystem/metrics",
		"/sys/kernel/stats/subsystem/device/metrics",
		"/sys/kernel/stats/subsystem/device/subdevice/metrics",
	}
	expectedMetricInfo := []MetricInfo{
		MetricInfo{Label: "", Path: paths[0]},
		MetricInfo{Label: "/device", Path: paths[1]},
		MetricInfo{Label: "/device/subdevice", Path: paths[2]},
	}

	for i, path := range paths {
		actualMetricInfo := subsysMetric.getMetricInfo(path)
		if diff := cmp.Diff(expectedMetricInfo[i], actualMetricInfo); diff != "" {
			t.Errorf("getMetricInfo mismatch on input path = %s,(-expected +actual):\n%s", path, diff)
		}
	}
}
