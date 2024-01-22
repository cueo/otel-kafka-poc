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

	fields := strings.Fields(data[1])

	for _, line := range data[2:] {
		if line == "" {
			continue
		}
		values := strings.Fields(line)
		if len(values) != len(fields) {
			return nil, fmt.Errorf("invalid vmstat data to parse")
		}
		for i, v := range values {
			switch fields[i] {
			case "r":
				vmStat.RunnableProcs += parseInt(v)
			case "b":
				vmStat.TotalProcs += parseInt(v)
			case "swpd":
				vmStat.Swapped += parseInt(v)
			case "free":
				vmStat.Free += parseInt(v)
			case "buff":
				vmStat.Buffered += parseInt(v)
			case "cache":
				vmStat.Cached += parseInt(v)
			case "inact":
				vmStat.Inactive += parseInt(v)
			case "active":
				vmStat.Active += parseInt(v)
			case "si":
				vmStat.SwapIn += parseInt(v)
			case "so":
				vmStat.SwapOut += parseInt(v)
			case "bi":
				vmStat.BlocksReceived += parseInt(v)
			case "bo":
				vmStat.BlocksSent += parseInt(v)
			case "in":
				vmStat.Interrupts += parseInt(v)
			case "cs":
				vmStat.ContextSwitches += parseInt(v)
			case "us":
				vmStat.UserTime += parseInt(v)
			case "sy":
				vmStat.SystemTime += parseInt(v)
			case "id":
				vmStat.IdleTime += parseInt(v)
			case "wa":
				vmStat.IoWaitTime += parseInt(v)
			case "st":
				vmStat.StolenTime += parseInt(v)
			}
		}
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
