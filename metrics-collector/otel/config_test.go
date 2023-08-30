package otel

import (
	"context"
	"testing"
)

func TestCreateExporter(t *testing.T) {
	ctx := context.Background()
	_, err := newExporter(ctx)
	if err != nil {
		t.Errorf("Expected exporter to be created")
	}
}
