package otelstats_test

import (
	"github.com/liiling/kernel-metrics-agent/statsfs/otelstats"

	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateStatsfsMetrics(t *testing.T) {
	statsfspath := "testsys/kernel/stats"
	expected := &otelstats.StatsfsMetrics{
		StatsfsPath: statsfspath,
		Metrics: map[string]otelstats.SubsysMetrics{
			"subsys0": {
				StatsfsPath:   statsfspath,
				SubSystemName: "subsys0",
				SubSystemPath: "testsys/kernel/stats/subsys0",
				Metrics: map[string]otelstats.MetricInfo{
					"subsys0/m0": otelstats.MetricInfo{
						Name: "m0",
						Flag: "GAUGE",
						Type: "FLOAT",
						Desc: "metric m0",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys0/dev0/m0": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys0"},
								{Key: "device", Value: "dev0"},
							},
							"testsys/kernel/stats/subsys0/dev1/m0": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys0"},
								{Key: "device", Value: "dev1"},
							},
						},
					},
					"subsys0/m1": otelstats.MetricInfo{
						Name: "m1",
						Flag: "GAUGE",
						Type: "FLOAT",
						Desc: "metric m1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys0/dev0/m1": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys0"},
								{Key: "device", Value: "dev0"},
							},
						},
					},
				},
			},
			"subsys1": {
				StatsfsPath:   statsfspath,
				SubSystemName: "subsys1",
				SubSystemPath: "testsys/kernel/stats/subsys1",
				Metrics: map[string]otelstats.MetricInfo{
					"subsys1/in_all_m": {
						Name: "in_all_m",
						Flag: "CUMULATIVE",
						Type: "INT",
						Desc: "a metric found in all devices under subsystem subsys1 and directly under subsys1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys1/dev0/in_all_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev0"},
							},
							"testsys/kernel/stats/subsys1/dev1/in_all_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev1"},
							},
							"testsys/kernel/stats/subsys1/in_all_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
							},
						},
					},
					"subsys1/in_both_devs_m": {
						Name: "in_both_devs_m",
						Flag: "CUMULATIVE",
						Type: "INT",
						Desc: "a metric found in both devices dev0 and dev1 under subsystem subsys1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys1/dev0/in_both_devs_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev0"},
							},
							"testsys/kernel/stats/subsys1/dev1/in_both_devs_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev1"},
							},
						},
					},
					"subsys1/in_top_and_dev0_m": {
						Name: "in_top_and_dev0_m",
						Flag: "CUMULATIVE",
						Type: "INT",
						Desc: "a metric found directly under subsystem subsys1 and device 0 in subsys1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys1/dev0/in_top_and_dev0_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev0"},
							},
							"testsys/kernel/stats/subsys1/in_top_and_dev0_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
							},
						},
					},
					"subsys1/only_in_dev0_m": {
						Name: "only_in_dev0_m",
						Flag: "CUMULATIVE",
						Type: "INT",
						Desc: "a metric found in device dev0 under subsystem subsys1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys1/dev0/only_in_dev0_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
								{Key: "device", Value: "dev0"},
							},
						},
					},
					"subsys1/top_level_m": {
						Name: "top_level_m",
						Flag: "CUMULATIVE",
						Type: "FLOAT",
						Desc: "a metric found directly under subsystem subsys1",
						PathToLabel: map[string][]otelstats.MetricLabel{
							"testsys/kernel/stats/subsys1/top_level_m": []otelstats.MetricLabel{
								{Key: "subsystem", Value: "subsys1"},
							},
						},
					},
				},
			},
		},
	}
	actual, err := otelstats.NewStatsfsMetrics(statsfspath)
	if err != nil {
		t.Errorf("NewStatsfsMetrics error: %v", err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("NewStatsfsMetrics mismatch (-expected +actual):\n%s", diff)
	}
}
