## OpenTelemetry Collector Demo

A simple demo testing the usage of [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector). A custom incrementing int counter is instrumented with [Prometheus](https://prometheus.io/). The host machine metrics is instrumented with [Prometheus/node_exporter](https://github.com/prometheus/node_exporter).
An OpenTelemetry Collector is then configured to receive the Prometheus info, then push it to another backend (currently this backend is StackDriver (Google Cloud Monitoring)).

## Running the Demo

From the directory kernel-metrics-agent/otel-collector:
```
docker-compose up
```
The command sets up the following three containers:

- otel-collector_prometheus_1
- otel-collector_otel-agent_1
- otel-collector_go-app_1
- otel-collector_node-exporter_1

Details of the containers can be found using `docker ps`.

The counter metric is exposed by go-app on port 2112, while the node_exporter metrics are exposed by node_exporter on port 9090. Both are scraped by the Prometheus receiver in otel-agent, which then export the metrics to port 8889.
Finally, it is exported to Prometheus and StackDriver. The metrics show up in the graph on port 9090 in the host machine, as well as in the Google Cloud Monitoring page of the corresponding GCP project.

As such, there are multiple ways to access the exposed metrics (each with some delays due to scraping interval configs):
1. Use `docker ps` to inspect the containers. For the custom counter metric, inspect `otel-collector_go-app_1`, note down the host ephemeral port mapped to the container port 2112 and access the metrics by visiting the ephemeral port. For instance, visit `localhost:32790/metrics` if 32790 is the port mapped to port 2112 in the container. For node_exporter metrics, visit `localhost:9100/metrics`.
2. Visit `localhost:8889/metrics`, where the Prometheus exporter in otel-agent exports the metrics data to.
3. Visit `localhost:9090`, where the Prometheus server is running.
4. Metrics on the status of the OpenTelemetry collector are found in `localhost:8888/metrics`.