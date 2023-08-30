package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mmynk/metrics-collector/otel"
)

func publishMetrics(ctx context.Context) {
	// metrics := below.ReadEthtoolMetrics()
	// kafka.PublishEthtoolMetrics(metrics)
	otel.CollectMetrics(ctx)
}

func RunPublisher() {
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	done := make(chan bool)

	// kill after 5 min
	go func() {
		// time.Sleep(5 * time.Minute)
		// logger.Info("Stopping metrics daemon...")
		// done <- true
		wait := make(chan os.Signal, 1)
		signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
		<-wait
		done <- true
	}()

	for {
		select {
		case <-ticker.C:
			publishMetrics(ctx)
		case <-done:
			logger.Info("Done publishing metrics")
			return
		}
	}
}
