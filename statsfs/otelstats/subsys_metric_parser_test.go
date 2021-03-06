package otelstats_test

import (
	"otelstats"
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
				Metrics: map[string][]otelstats.MetricInfo{
					"subsys0/m0": []otelstats.MetricInfo{
						otelstats.MetricInfo{
							Name:  "subsys0/m0",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys0/dev0/m0",
						},
						otelstats.MetricInfo{
							Name:  "subsys0/m0",
							Label: "/dev1",
							Path:  "testsys/kernel/stats/subsys0/dev1/m0",
						},
					},
					"subsys0/m1": []otelstats.MetricInfo{
						otelstats.MetricInfo{
							Name:  "subsys0/m1",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys0/dev0/m1",
						},
					},
				},
			},
			"subsys1": otelstats.SubsysMetrics{
				StatsfsPath:   statsfspath,
				SubSystemName: "subsys1",
				SubSystemPath: "testsys/kernel/stats/subsys1",
				Metrics: map[string][]otelstats.MetricInfo{
					"subsys1/in_all_m": {
						otelstats.MetricInfo{
							Name:  "subsys1/in_all_m",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys1/dev0/in_all_m",
						},
						otelstats.MetricInfo{
							Name:  "subsys1/in_all_m",
							Label: "/dev1",
							Path:  "testsys/kernel/stats/subsys1/dev1/in_all_m",
						},
						otelstats.MetricInfo{
							Name:  "subsys1/in_all_m",
							Label: "",
							Path:  "testsys/kernel/stats/subsys1/in_all_m",
						},
					},
					"subsys1/in_both_devs_m": {
						otelstats.MetricInfo{
							Name:  "subsys1/in_both_devs_m",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys1/dev0/in_both_devs_m",
						},
						otelstats.MetricInfo{
							Name:  "subsys1/in_both_devs_m",
							Label: "/dev1",
							Path:  "testsys/kernel/stats/subsys1/dev1/in_both_devs_m",
						},
					},
					"subsys1/in_top_and_dev0_m": {
						otelstats.MetricInfo{
							Name:  "subsys1/in_top_and_dev0_m",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys1/dev0/in_top_and_dev0_m",
						},
						otelstats.MetricInfo{
							Name:  "subsys1/in_top_and_dev0_m",
							Label: "",
							Path:  "testsys/kernel/stats/subsys1/in_top_and_dev0_m",
						},
					},
					"subsys1/only_in_dev0_m": {
						otelstats.MetricInfo{
							Name:  "subsys1/only_in_dev0_m",
							Label: "/dev0",
							Path:  "testsys/kernel/stats/subsys1/dev0/only_in_dev0_m",
						},
					},
					"subsys1/top_level_m": {
						otelstats.MetricInfo{
							Name:  "subsys1/top_level_m",
							Label: "",
							Path:  "testsys/kernel/stats/subsys1/top_level_m",
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
