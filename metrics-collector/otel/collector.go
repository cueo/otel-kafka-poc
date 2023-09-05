package otel

import (
	"context"
	"mmynk/metrics-collector/below"
	"regexp"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
	// "go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Metric struct {
	Name   string
	Value  float64
	Unit   string
	Type   string
	Labels map[string]string
}

type Metrics []Metric

func parseOpenMetrics(text string) Metrics {
	metricsLines := strings.Split(text, "\n")
	var metrics Metrics

	metric := Metric{}
	for _, metricLine := range metricsLines {
		if len(metricLine) == 0 {
			continue
		}
		metricLine = strings.TrimSpace(metricLine)

		if strings.HasPrefix(metricLine, "# TYPE") {
			// parse type
			metric.Type = strings.Split(metricLine, " ")[3]
		} else if strings.HasPrefix(metricLine, "# UNIT") {
			// parse unit
			metric.Unit = strings.Split(metricLine, " ")[3]
		} else if !strings.HasPrefix(metricLine, "#") {
			// parse metric
			segments := strings.Split(metricLine, " ")
			name := segments[0]

			// parse attributes
			labels := make(map[string]string)
			re := regexp.MustCompile(`(.*)\{(.*)\}`)
			if re.MatchString(name) {
				attributesStr := re.FindStringSubmatch(name)[2]
				attributesArr := strings.Split(attributesStr, ",")
				for _, a := range attributesArr {
					segments := strings.Split(a, "=")
					labels[segments[0]] = removeQuotes(segments[1])
				}
				name = re.FindStringSubmatch(name)[1]
			}

			metric.Name = name
			metric.Value, _ = strconv.ParseFloat(segments[1], 64)
			metric.Labels = labels
			metrics = append(metrics, metric)
		} else if strings.HasPrefix(metricLine, "# EOF") {
			// end of file
			break
		}
	}
	logger.Debug("Parsed metrics", zap.Any("metrics", metrics))
	return metrics
}

func removeQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func gauge(om Metric, m pmetric.Metric) pmetric.Gauge {
	g := m.SetEmptyGauge()
	d := g.DataPoints().AppendEmpty()
	d.SetDoubleValue(om.Value)
	return g
}

func counter(om Metric, m pmetric.Metric) pmetric.Sum {
	ctr := m.Sum()
	ctr.SetIsMonotonic(true)
	d := ctr.DataPoints().AppendEmpty()
	d.SetIntValue(int64(om.Value))
	return ctr
}

// TODO: add labels
// func attributes(labels map[string]string) []attribute.KeyValue {
// 	var attrs []attribute.KeyValue
// 	for k, v := range labels {
// 		attrs = append(attrs, attribute.String(k, v))
// 	}
// 	return attrs
// }

func collectBelowMetrics() (pmetric.ResourceMetrics, error) {
	logger.Info("Collecting below metrics")
	var metrics Metrics
	bms := below.ReadMetrics()
	for _, bm := range bms {
		metrics = append(metrics, parseOpenMetrics(bm)...)
	}

	// mp := meterProvider()
	// meter := mp.Meter("below")

	rm := pmetric.NewResourceMetrics()
	sm := rm.ScopeMetrics().AppendEmpty()
	for _, om := range metrics {
		// create metric
		m := sm.Metrics().AppendEmpty()
		m.SetName(om.Name)
		m.SetUnit(om.Unit)

		switch om.Type {
		case "gauge":
			// gauge, err := meter.Float64ObservableGauge(om.Name, metric.WithUnit(om.Unit))
			// _, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
			// 	o.ObserveFloat64(gauge, om.Value, observeOption(om.Labels))
			// 	return nil
			// }, g)
			gauge(om, m)
		case "counter":
		default:
			// meter.Int64ObservableCounter(om.Name, metric.WithUnit(om.Unit))
			// ctr, err := meter.Int64Counter(om.Name, metric.WithUnit(om.Unit))
			// if err != nil {
			// 	logger.Error("Failed to create counter", zap.String("metric", om.Name), zap.Error(err))
			// } else {
			// 	ctr.Add(ctx, int64(om.Value), addOption(om.Labels))
			// }
			counter(om, m)
		}
	}

	return rm, nil
}

func CollectMetrics(ctx context.Context) error {
	logger.Info("Start metric collection")
	me, err := newExporter(ctx)
	if err != nil {
		logger.Error("Failed to create exporter", zap.Error(err))
		return err
	}

	err = me.Start(ctx, nil)
	if err != nil {
		logger.Error("Failed to start exporter", zap.Error(err))
		return err
	}
	defer me.Shutdown(ctx)

	ms := pmetric.NewMetrics()
	rm, err := collectBelowMetrics()
	if err != nil {
		logger.Error("Failed to collect below metrics", zap.Error(err))
	}
	r := ms.ResourceMetrics().AppendEmpty()
	rm.CopyTo(r)

	logger.Info("Publishing metrics", zap.Any("metrics", ms))
	err = me.ConsumeMetrics(ctx, ms)
	if err != nil {
		logger.Error("Failed to consume metrics", zap.Error(err))
		return err
	}

	// combine errors for different collectors once we add additional collectors
	logger.Info("Finished metric collection")
	return err
}
