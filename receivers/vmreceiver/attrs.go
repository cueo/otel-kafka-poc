package vmreceiver

import (
	"os/exec"

	"go.uber.org/zap"
)

type attributes struct {
	host string
	os   string
	arch string
}

type attributeReader struct {
	logger *zap.Logger
}

func newAttributeReader(logger *zap.Logger) *attributeReader {
	return &attributeReader{
		logger: logger,
	}
}

func (a *attributeReader) getAttributes() *attributes {
	var h, o, ar string

	cmd := exec.Command("hostname")
	bytes, err := cmd.Output()
	if err != nil {
		a.logger.Error("failed to execute hostname", zap.Error(err))
		h = "unknown"
	}
	h = string(bytes)

	cmd = exec.Command("uname", "-s")
	bytes, err = cmd.Output()
	if err != nil {
		a.logger.Error("failed to execute uname -s", zap.Error(err))
		o = "unknown"
	}
	o = string(bytes)

	cmd = exec.Command("uname", "-m")
	bytes, err = cmd.Output()
	if err != nil {
		a.logger.Error("failed to execute uname -m", zap.Error(err))
		ar = "unknown"
	}
	ar = string(bytes)

	return &attributes{
		host: h,
		os:   o,
		arch: ar,
	}
}
