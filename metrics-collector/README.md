# Metrics Collector

This module is responsible for collecting metrics from multiple sources and then exporting it with OpenTelemetry using [`kafkaexporter`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/kafkaexporter).

## Run the collector

The collector collects metrics from the sources every 1 min and publishes it to a Kafka topic.

```shell
go mod download
go run .
```
