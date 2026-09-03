package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestAbsoluteCounterVecExposesAbsoluteValuesAsCounter(t *testing.T) {
	counter := newAbsoluteCounterVec(
		prometheus.CounterOpts{
			Name: "test_external_events_total",
			Help: "Events reported by an external system.",
		},
		[]string{"source"},
	)
	registry := prometheus.NewRegistry()
	registry.MustRegister(counter)

	counter.WithLabelValues("adc").Set(42)
	assertCounterFamilyValue(t, gatherMetricFamilies(t, registry), "test_external_events_total", 42)

	// A lower absolute value represents a reset in the external system. It must
	// replace the exported sample rather than be added to exporter-local state.
	counter.Reset()
	counter.WithLabelValues("adc").Set(3)
	assertCounterFamilyValue(t, gatherMetricFamilies(t, registry), "test_external_events_total", 3)
}

func TestAbsoluteCounterVecRequiresTotalSuffix(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("newAbsoluteCounterVec() did not reject a counter name without _total")
		}
	}()

	newAbsoluteCounterVec(prometheus.CounterOpts{Name: "test_external_events"}, nil)
}

func TestMetricHelpersExposeExplicitTypes(t *testing.T) {
	exporter := &Exporter{}
	tests := []struct {
		name     string
		send     func(chan<- prometheus.Metric, *prometheus.Desc, any, []string)
		wantType dto.MetricType
	}{
		{
			name:     "gauge",
			send:     exporter.sendMetric,
			wantType: dto.MetricType_GAUGE,
		},
		{
			name:     "counter",
			send:     exporter.sendCounterMetric,
			wantType: dto.MetricType_COUNTER,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desc := prometheus.NewDesc("test_"+test.name, "Test metric.", nil, nil)
			ch := make(chan prometheus.Metric, 1)
			test.send(ch, desc, "12", nil)

			metric := <-ch
			var encoded dto.Metric
			if err := metric.Write(&encoded); err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			switch test.wantType {
			case dto.MetricType_COUNTER:
				if encoded.Counter == nil || encoded.Counter.GetValue() != 12 {
					t.Fatalf("counter = %v, want value 12", encoded.Counter)
				}
			case dto.MetricType_GAUGE:
				if encoded.Gauge == nil || encoded.Gauge.GetValue() != 12 {
					t.Fatalf("gauge = %v, want value 12", encoded.Gauge)
				}
			}
		})
	}
}

func assertCounterFamilyValue(t *testing.T, families map[string]*dto.MetricFamily, name string, want float64) {
	t.Helper()
	family := families[name]
	if family == nil {
		t.Fatalf("metric family %q is absent", name)
	}
	if family.GetType() != dto.MetricType_COUNTER {
		t.Fatalf("%s type = %s, want COUNTER", name, family.GetType())
	}
	metric := singleMetric(t, families, name)
	if got := metric.GetCounter().GetValue(); got != want {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}
