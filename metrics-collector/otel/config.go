package otel

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/otel/sdk/metric"

	// "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	// semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

// func newResource() *resource.Resource {
// 	r, _ := resource.Merge(
// 		resource.Default(),
// 		resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceName("metrics_daemon"),
// 			semconv.ServiceVersion("0.0.1"),
// 		),
// 	)
// 	return r
// }

func newExporter(ctx context.Context) (exporter.Metrics, error) {
	f := kafkaexporter.NewFactory(kafkaexporter.WithMetricsMarshalers())
	cfg := f.CreateDefaultConfig().(*kafkaexporter.Config)
	cfg.Topic = "metrics"
	ts := component.TelemetrySettings{
		Logger:          logger,
		MeterProvider:   meterProvider(),
		TracerProvider:  tracerProvider(),
		MetricsLevel:    configtelemetry.LevelNormal,
	}
	cs := exporter.CreateSettings{
		ID:               component.NewID("metrics"),
		TelemetrySettings: ts,
		BuildInfo:        component.NewDefaultBuildInfo(),
	}

	return f.CreateMetricsExporter(ctx, cs, cfg)
}

func meterProvider() *metric.MeterProvider {
	return metric.NewMeterProvider()
}

func tracerProvider() *trace.TracerProvider {
	return trace.NewTracerProvider()
}
