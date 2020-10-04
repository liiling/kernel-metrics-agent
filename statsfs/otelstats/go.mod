module github.com/liiling/kernel-metrics-agent/statsfs/otelstats

go 1.14

require (
	cloud.google.com/go v0.66.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.11.1-0.20200918040108-45880ec6f71f
	github.com/google/go-cmp v0.5.2
	github.com/prometheus/common v0.14.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.11.0
	go.opentelemetry.io/otel/exporters/stdout v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	golang.org/x/sys v0.0.0-20200929083018-4d22bbb62b3c // indirect
	google.golang.org/api v0.32.0 // indirect
	google.golang.org/genproto v0.0.0-20200925023002-c2d885f95484 // indirect
	google.golang.org/grpc v1.32.0 // indirect
)
