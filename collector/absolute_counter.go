package collector

import (
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// absoluteCounterVec exposes absolute cumulative values read from the
// NetScaler as Prometheus counters. Unlike CounterVec, values are replaced on
// each scrape instead of added to exporter-local state.
type absoluteCounterVec struct {
	desc       *prometheus.Desc
	labelCount int

	mu     sync.RWMutex
	values map[string]absoluteCounterSample
}

type absoluteCounterSample struct {
	value  float64
	labels []string
}

type absoluteCounter struct {
	vec    *absoluteCounterVec
	key    string
	labels []string
}

func newAbsoluteCounterVec(opts prometheus.CounterOpts, labelNames []string) *absoluteCounterVec {
	if !strings.HasSuffix(opts.Name, "_total") {
		panic("absolute counter metric name must end in _total")
	}

	return &absoluteCounterVec{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
			opts.Help,
			labelNames,
			opts.ConstLabels,
		),
		labelCount: len(labelNames),
		values:     make(map[string]absoluteCounterSample),
	}
}

func (v *absoluteCounterVec) Describe(ch chan<- *prometheus.Desc) {
	ch <- v.desc
}

func (v *absoluteCounterVec) Collect(ch chan<- prometheus.Metric) {
	v.mu.RLock()
	samples := make([]absoluteCounterSample, 0, len(v.values))
	for _, sample := range v.values {
		samples = append(samples, sample)
	}
	v.mu.RUnlock()

	for _, sample := range samples {
		ch <- prometheus.MustNewConstMetric(v.desc, prometheus.CounterValue, sample.value, sample.labels...)
	}
}

func (v *absoluteCounterVec) Reset() {
	v.mu.Lock()
	clear(v.values)
	v.mu.Unlock()
}

func (v *absoluteCounterVec) WithLabelValues(labels ...string) *absoluteCounter {
	if len(labels) != v.labelCount {
		panic("inconsistent label cardinality")
	}

	labels = append([]string(nil), labels...)
	return &absoluteCounter{
		vec:    v,
		key:    absoluteCounterLabelKey(labels),
		labels: labels,
	}
}

func (c *absoluteCounter) Set(value float64) {
	c.vec.mu.Lock()
	c.vec.values[c.key] = absoluteCounterSample{value: value, labels: c.labels}
	c.vec.mu.Unlock()
}

func absoluteCounterLabelKey(labels []string) string {
	var key strings.Builder
	for _, label := range labels {
		key.WriteString(strconv.Itoa(len(label)))
		key.WriteByte(':')
		key.WriteString(label)
	}
	return key.String()
}
