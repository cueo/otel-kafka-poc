package main

import (
	"mmynk/metrics-receiver/otel"
	"context"

	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	logger.Info("Starting metrics receiver...")
	ctx := context.Background()
	otel.ProcessMetrics(ctx)
}
