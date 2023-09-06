# otel-kafka-poc

A POC for collecting, publishing and receiving metrics using OpenTelemety with Kafka.

## Running the services

All the services are Dockerized and use `docker-compose` to easily build and run all the services: collector / exporter, receiver and Kafka.

```shell
docker compose up
```

That's all, the `metrics-collector` service should start collecting the metrices every minute and exporting them to Kafka whereas the `metrics-receiver` service should start listening to the Kafka topic and logging the metrics.
