package otel

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()

	attempts = 5
	sleep    = 2 * time.Second

	// Kafka
	broker = "kafka:9092"
	topic  = "avx-metrics"
)

func retryCreateReceiver(
	factory receiver.Factory,
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics) (receiver.Metrics, error) {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		mr, err := factory.CreateMetricsReceiver(ctx, set, cfg, nextConsumer)
		if err == nil {
			return mr, nil
		}
	}
	return nil, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func newReceiver(ctx context.Context) (receiver.Metrics, error) {
	f := kafkareceiver.NewFactory(kafkareceiver.WithMetricsUnmarshalers())

	cfg := f.CreateDefaultConfig().(*kafkareceiver.Config)
	cfg.Brokers = []string{broker}
	cfg.Topic = topic

	ts := component.TelemetrySettings{
		Logger:         logger,
		MeterProvider:  meterProvider(),
		TracerProvider: tracerProvider(),
		MetricsLevel:   configtelemetry.LevelNormal,
	}
	set := receiver.CreateSettings{
		ID:                component.NewID("metrics"),
		TelemetrySettings: ts,
		BuildInfo:         component.NewDefaultBuildInfo(),
	}

	nc, _ := consumer.NewMetrics(processor)

	logger.Info("Creating metrics receiver...")
	return retryCreateReceiver(f, ctx, set, cfg, nc)
}

func meterProvider() *metric.MeterProvider {
	return metric.NewMeterProvider()
}

func tracerProvider() *trace.TracerProvider {
	return trace.NewTracerProvider()
}
