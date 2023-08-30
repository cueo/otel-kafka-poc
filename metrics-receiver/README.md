# Metrics Receiver

This module is responsible for receiving the exported metrics with OpenTelemetry using [`kafkareceiver`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver).

## Run the receiver

The receiver listens to the given topic and logs the metric fields.

```shell
go mod download
go run .
```
