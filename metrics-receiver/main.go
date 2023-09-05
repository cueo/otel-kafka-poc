package main

import (
	"context"
	"mmynk/metrics-receiver/otel"

	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	logger.Info("Starting metrics receiver...")
	ctx := context.Background()
	otel.ProcessMetrics(ctx)
}
