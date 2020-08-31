package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initPrometheus() {
	port := 2112
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Printf("Prometheus server running on :%d\n", port)
}

func createIntSumCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A counter that increments every second",
	},
		[]string{"label1"})
}

func incrementCounterMetrics(counter prometheus.Counter) {
	for {
		counter.Inc()
		time.Sleep(time.Second)
	}
}

func main() {
	initPrometheus()

	counter := createIntSumCounter().With(prometheus.Labels{"label1": "value1"})

	fmt.Println("Start incrementing counter...")
	incrementCounterMetrics(counter)
}
