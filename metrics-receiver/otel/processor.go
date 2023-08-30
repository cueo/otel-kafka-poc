package otel

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func processor(ctx context.Context, ld pmetric.Metrics) error {
	logger.Info("Received metrics")
	metrics := ld.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	logger.Info("Metrics length", zap.Int("length", metrics.Len()))
	for i := 0; i < metrics.Len(); i++ {
		m := metrics.At(i)
		var val float64
		switch t := m.Type(); t {
		case pmetric.MetricTypeGauge:
			val = m.Gauge().DataPoints().At(0).DoubleValue()
		case pmetric.MetricTypeSum:
			val = m.Sum().DataPoints().At(0).DoubleValue()
		}

		logger.Info("Fetched metric",
			zap.String("metric", m.Name()),
			zap.String("unit", m.Unit()),
			zap.Float64("value", val),
		)
	}
	return nil
}

func ProcessMetrics(ctx context.Context) {
	r, err := newReceiver(ctx)
	if err != nil {
		logger.Error("error creating receiver", zap.Error(err))
	}
	defer r.Shutdown(ctx)

	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)
    go func() {
        <-c
        logger.Info("Shutting down...")
        os.Exit(1)
    }()

	r.Start(ctx, nil)
	<-done
}
