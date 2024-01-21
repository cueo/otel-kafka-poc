package vmreceiver

import (
	"errors"
	"fmt"

	"github.com/mmynk/otel-kafka-poc/receivers/vmreceiver/internal/metadata"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	Delay int `mapstructure:",omitempty"` // Delay is the delay between consecutive `vmstat` calls.
	Count int `mapstructure:",omitempty"` // Count is the number of `vmstat` calls to make.

	metadata.MetricsBuilderConfig           `mapstructure:",squash"` // MetricsBuilderConfig to enable/disable specific metrics (default: all enabled)
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"` // ScraperControllerSettings to configure scraping interval (default: scrape every second)
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
