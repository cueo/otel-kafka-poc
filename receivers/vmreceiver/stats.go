package vmreceiver

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type vmStat struct {
	RunnableProcs   int64
	TotalProcs      int64
	Swapped         int64
	Free            int64
	Buffered        int64
	Cached          int64
	Inactive        int64
	Active          int64
	SwapIn          int64
	SwapOut         int64
	BlocksReceived  int64
	BlocksSent      int64
	Interrupts      int64
	ContextSwitches int64
	UserTime        int64
	SystemTime      int64
	IdleTime        int64
	IoWaitTime      int64
	StolenTime      int64
}

type vmStatReader struct {
	delay int
	count int

	logger *zap.Logger
}

func newVmStatReader(cfg *Config, logger *zap.Logger) *vmStatReader {
	return &vmStatReader{
		delay:  cfg.Delay,
		count:  cfg.Count,
		logger: logger,
	}
}

func (r *vmStatReader) Read() (*vmStat, error) {
	cmd := exec.Command("vmstat", fmt.Sprintf("%d", r.delay), fmt.Sprintf("%d", r.count))
	out, err := cmd.Output()
	if err != nil {
		r.logger.Error("failed to execute vmstat", zap.Error(err))
		return nil, err
	}
	return r.parse(out)
}

func (r *vmStatReader) parse(out []byte) (*vmStat, error) {
	o := string(out)
	var vmStat vmStat
	data := strings.Split(o, "\n")
	if len(data) < 3 {
		return nil, fmt.Errorf("no vmstat data to parse")
	}

	parseInt := func(s string) int64 {
		i, err := strconv.Atoi(s)
		if err != nil {
			r.logger.Error("failed to parse int, returning 0", zap.Error(err))
			return 0
		}
		return int64(i)
	}

	for _, line := range data[2:] {
		fields := strings.Fields(line)
		if len(fields) < 18 {
			continue
		}
		vmStat.RunnableProcs += parseInt(fields[0])
		vmStat.TotalProcs += parseInt(fields[1])
		vmStat.Swapped += parseInt(fields[2])
		vmStat.Free += parseInt(fields[3])
		vmStat.Buffered += parseInt(fields[4])
		vmStat.Cached += parseInt(fields[5])
		vmStat.Inactive += parseInt(fields[6])
		vmStat.Active += parseInt(fields[7])
		vmStat.SwapIn += parseInt(fields[8])
		vmStat.SwapOut += parseInt(fields[9])
		vmStat.BlocksReceived += parseInt(fields[10])
		vmStat.BlocksSent += parseInt(fields[11])
		vmStat.Interrupts += parseInt(fields[12])
		vmStat.ContextSwitches += parseInt(fields[13])
		vmStat.UserTime += parseInt(fields[14])
		vmStat.SystemTime += parseInt(fields[15])
		vmStat.IdleTime += parseInt(fields[16])
		vmStat.IoWaitTime += parseInt(fields[17])
		vmStat.StolenTime += parseInt(fields[18])
	}

	n := int64(len(data[2:]))
	if n > 1 {
		vmStat.RunnableProcs /= n
		vmStat.TotalProcs /= n
		vmStat.Swapped /= n
		vmStat.Free /= n
		vmStat.Buffered /= n
		vmStat.Cached /= n
		vmStat.Inactive /= n
		vmStat.Active /= n
		vmStat.SwapIn /= n
		vmStat.SwapOut /= n
		vmStat.BlocksReceived /= n
		vmStat.BlocksSent /= n
		vmStat.Interrupts /= n
		vmStat.ContextSwitches /= n
		vmStat.UserTime /= n
		vmStat.SystemTime /= n
		vmStat.IdleTime /= n
		vmStat.IoWaitTime /= n
		vmStat.StolenTime /= n
	}

	return &vmStat, nil
}
