package vmreceiver

import (
	"context"
	"fmt"

	"github.com/mmynk/otel-kafka-poc/receivers/vmreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/zap"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(CreateVmStatReceiver, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Delay: 1,
		Count: 1,
	}
}

func CreateVmStatReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	logger := settings.Logger
	config, ok := cfg.(*Config)
	if !ok {
		em := "failed to cast to type Config"
		logger.Error(em)
		return nil, fmt.Errorf(em)
	}

	mb := metadata.NewMetricsBuilder(config.MetricsBuilderConfig, settings)

	ns := newScraper(config, mb, logger)
	scraper, err := scraperhelper.NewScraper(metadata.Type, ns.scrape)
	if err != nil {
		logger.Error("failed to create scraper", zap.Error(err))
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(
		&config.ScraperControllerSettings,
		settings,
		consumer,
		scraperhelper.AddScraper(scraper),
	)
}
