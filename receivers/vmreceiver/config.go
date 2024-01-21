package vmreceiver

import (
	"errors"
	"fmt"

	"github.com/mmynk/otel-kafka-poc/receivers/vmreceiver/internal/metadata"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	// Delay is the delay between `vmstat` calls
	Delay int `mapstructure:"delay"`
	// Count is the number of `vmstat` calls to make
	Count int `mapstructure:"count"`

	// MetricsBuilderConfig to enable/disable specific metrics (default: all enabled)
	metadata.MetricsBuilderConfig `mapstructure:",squash"`
	// ScraperControllerSettings to configure scraping interval (default: scrape every second)
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
}

func (cfg *Config) Validate() error {
	var err error
	if cfg.Delay < 0 {
		err = errors.Join(err, fmt.Errorf("delay cannot be negative"))
	}
	if cfg.Count < 0 {
		err = errors.Join(err, fmt.Errorf("count cannot be negative"))
	}
	return err
}
