package vmreceiver

import (
	"context"
	"time"

	"github.com/mmynk/otel-kafka-poc/receivers/vmreceiver/internal/metadata"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

const (
	procs  = "proc"
	memory = "memory"
	swap   = "swap"
	io     = "io"
	cpu    = "cpu"
)

type scraper struct {
	logger         *zap.Logger              // Logger to log events
	metricsBuilder *metadata.MetricsBuilder // MetricsBuilder to build metrics
	reader         *vmStatReader            // vmStatReader to read vmstat output
}

func newScraper(cfg *Config, metricsBuilder *metadata.MetricsBuilder, logger *zap.Logger) *scraper {
	return &scraper{
		logger:         logger,
		metricsBuilder: metricsBuilder,
		reader:         newVmStatReader(cfg, logger),
	}
}

func (s *scraper) scrape(_ context.Context) (pmetric.Metrics, error) {
	s.logger.Info("Scraping vm stats")
	vmStat, err := s.reader.Read()
	if err != nil {
		return pmetric.Metrics{}, err
	}
	attr := newAttributeReader(s.logger).getAttributes()
	s.recordVmStats(vmStat, attr)
	return s.metricsBuilder.Emit(), nil
}

func (s *scraper) recordVmStats(stat *vmStat, attr *attributes) {
	now := pcommon.NewTimestampFromTime(time.Now())

	s.metricsBuilder.RecordRunnableProcsDataPoint(now, stat.RunnableProcs, attr.host, attr.os, attr.arch, procs)
	s.metricsBuilder.RecordTotalProcsDataPoint(now, stat.TotalProcs, attr.host, attr.os, attr.arch, procs)
	s.metricsBuilder.RecordSwappedDataPoint(now, stat.Swapped, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordFreeDataPoint(now, stat.Free, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordBufferedDataPoint(now, stat.Buffered, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordCachedDataPoint(now, stat.Cached, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordInactiveDataPoint(now, stat.Inactive, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordActiveDataPoint(now, stat.Active, attr.host, attr.os, attr.arch, memory)
	s.metricsBuilder.RecordSwapInDataPoint(now, stat.SwapIn, attr.host, attr.os, attr.arch, swap)
	s.metricsBuilder.RecordSwapOutDataPoint(now, stat.SwapOut, attr.host, attr.os, attr.arch, swap)
	s.metricsBuilder.RecordBlocksReceivedDataPoint(now, stat.BlocksReceived, attr.host, attr.os, attr.arch, io)
	s.metricsBuilder.RecordBlocksSentDataPoint(now, stat.BlocksSent, attr.host, attr.os, attr.arch, io)
	s.metricsBuilder.RecordInterruptsDataPoint(now, stat.Interrupts, attr.host, attr.os, attr.arch, io)
	s.metricsBuilder.RecordContextSwitchesDataPoint(now, stat.ContextSwitches, attr.host, attr.os, attr.arch, io)
	s.metricsBuilder.RecordUserTimeDataPoint(now, stat.UserTime, attr.host, attr.os, attr.arch, cpu)
	s.metricsBuilder.RecordSystemTimeDataPoint(now, stat.SystemTime, attr.host, attr.os, attr.arch, cpu)
	s.metricsBuilder.RecordIdleTimeDataPoint(now, stat.IdleTime, attr.host, attr.os, attr.arch, cpu)
	s.metricsBuilder.RecordIoWaitTimeDataPoint(now, stat.IoWaitTime, attr.host, attr.os, attr.arch, cpu)
	s.metricsBuilder.RecordStolenTimeDataPoint(now, stat.StolenTime, attr.host, attr.os, attr.arch, cpu)
}
