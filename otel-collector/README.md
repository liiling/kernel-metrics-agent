## OpenTelemetry Collector Demo

A simple demo testing the usage of [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector). An incrementing int counter is instrumented with OpenTelemetry, then exported to [Prometheus](https://prometheus.io/) using the [Prometheus exporter](https://github.com/open-telemetry/opentelemetry-go/tree/master/exporters/metric/prometheus). An OpenTelemetry Collector is then configured to receive the Prometheus info, then push it to another backend (currently this backend is Prometheus as well...).

## Running the Demo

From the directory kernel-metrics-agent/otel-collector:
```
docker-compose up
```
The command sets up the following three containers:

- otel-collector_prometheus_1
- otel-collector_otel-agent_1
- otel-collector_go-app_1

Details of the containers can be found using `docker ps`.

The counter metric is first exposed by go-app on port 2112, collected by the receiver in otel-agent, and then further exported by the exporter in otel-agent on port 8889. Finally, it is scraped by Prometheus and shows up in the graph on port 9090 in the host machine.

As such, there are multiple ways to access the counter metrics (each with some delays due to scraping interval configs):
1. Use `docker ps` to inspect the container `otel-collector_go-app_1`, note down the host ephemeral port mapped to the container port 2112 and access the metrics by visiting the ephemeral port. For instance, visit `localhost:32790` if 32790 is the port mapped to port 2112 in the container. 
2. Visit `localhost:8889`, where the Prometheus exporter in otel-agent exports the metrics data to.
3. Visit `localhost:9090`, where the Prometheus server is running.
