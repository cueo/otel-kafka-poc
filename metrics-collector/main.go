package main

import "go.uber.org/zap"

var logger, _ = zap.NewDevelopment()

func main() {
	logger.Info("Starting metrics daemon...")
	RunPublisher()
	// publishMetrics()
}
