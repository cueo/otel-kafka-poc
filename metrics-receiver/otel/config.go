package otel

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func newReceiver(ctx context.Context) (receiver.Metrics, error) {
	f := kafkareceiver.NewFactory(kafkareceiver.WithMetricsUnmarshalers())

	cfg := f.CreateDefaultConfig().(*kafkareceiver.Config)
	cfg.Topic = "metrics"

	ts := component.TelemetrySettings{
		Logger:          logger,
		MeterProvider:   meterProvider(),
		TracerProvider:  tracerProvider(),
		MetricsLevel:    configtelemetry.LevelNormal,
	}
	set := receiver.CreateSettings{
		ID:               component.NewID("metrics"),
		TelemetrySettings: ts,
		BuildInfo:        component.NewDefaultBuildInfo(),
	}

	consumer, _ := consumer.NewMetrics(processor)

	logger.Info("Creating metrics receiver...")
	return f.CreateMetricsReceiver(ctx, set, cfg, consumer)
}

func meterProvider() *metric.MeterProvider {
	return metric.NewMeterProvider()
}

func tracerProvider() *trace.TracerProvider {
	return trace.NewTracerProvider()
}
