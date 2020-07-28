## Description about OpenTelemetry

[OpenTelemetry](https://opentelemetry.io/) is a set tools (e.g. APIs, SDKs) for the collection and management of telemetry data such as traces and metrics. 
A [trace](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/overview.md#distributed-tracing) is a set of events caused by the same logical operation, with impact across multiple parts of an application. For instance, a trace could be started when sending a request to retrive a website.
A [metric](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/overview.md#metrics) is a raw measurement with a set of labels. Some examples of metrics are CPU usage, number of cache misses, etc. 
More details on metrics are provided in [OpenTelemetry metric API](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/metrics/api.md#metrics-api).
In this project, only metrics will be utilised.

Once the telemetry data is collected using OpenTelemtry, it can be exported to a [supported backend system](https://github.com/open-telemetry/opentelemetry-go/tree/master/exporters) or a [third-party telemetry system](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/exporters). 
As a result, changing the backend monitoring system requires a change of exporters, but no change in the actual instrumentation of the program.

OpenTelemetry also offers [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/about/), a binary that allows for vendor-agnostic collection and export of telemetry data to multiple monitoring systems.

## Description of the Example Pipeline

An example go application that runs a simple http server and exports a counter metric to different backend depending on the flag provided. Allowed exporter flags include `stdout`, `prometheus` and `gcp`.

## Runing the Example

1. Start the go server at `localhost:8090` with `go run server.go`
2. Start the exporter with `go run main.go --exporter x ` where x is one of [`stdout`, `prometheus`, `gcp`]. 
Exporter `stdout` requires no extra set up. For the other two exporters, follow the instructions below to set up [prometheus](#set-up-prometheus) and [Google Cloud Monitoring](#set-up-google-cloud-monitoring). 
3. The counter metric exposed represents the number of times the page `localhost:8090/visit` is visited. Visit `localhost:8090/visit` to increment the counter.

### Set Up Prometheus 
#### 1. [Install Prometheus](https://prometheus.io/docs/introduction/first_steps/#downloading-prometheus)

1. Download [latest Prometheus release](https://prometheus.io/download/).
2. 
    ```
    tar xvfz prometheus-*.tar.gz
    ```
3. Place the Prometheus folder to a desired location and add the location to PATH.

#### 2. Running Prometheus

1. Start Prometheus: `prometheus --config.file=prometheus-config.yml`
2. Access the text-formatted metrics at `localhost:2112/metrics`, or access the Prometheus client at `localhost:9090/graph`

### Set up Google Cloud Monitoring

See [OpenTelemetry Google CLoud Monitoring Exporter](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/master/exporter/metric/README.md) for more info.

1. [Create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects) on Google Cloud Platform. 
2. [Create a Workspace](https://cloud.google.com/monitoring/workspaces/create) inside the GCP project.
3. [Create a service account](https://cloud.google.com/docs/authentication/production#create_service_account) and download the JSON key file.
4. Set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` to the file path of the JSON key file.
5. On Google Cloud console, navigate to `Monitoring` -> `Metric Explorer` and search for `custom/opentelemetry/otel-test-counter` to see the metric update.