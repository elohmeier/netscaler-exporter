package collector

import (
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestExporterSelfMetricsReportCollectorAndScrapeStatus(t *testing.T) {
	var fail atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nitro/v1/config/lbvserver_rewritepolicy_binding" {
			http.NotFound(w, r)
			return
		}
		if fail.Load() {
			http.Error(w, "temporary failure", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"lbvserver_rewritepolicy_binding":[],"errorcode":0}`)
	}))
	t.Cleanup(server.Close)

	exporter := newRewriteTestExporter(t, server.URL, nil)
	exporter.buildInfo = BuildInfo{Version: "1.2.3", Revision: "abc123"}
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families := gatherMetricFamilies(t, registry)
	assertSingleGaugeValue(t, families, "netscaler_exporter_scrape_success", 1)
	assertSingleGaugeValue(t, families, "netscaler_exporter_collector_success", 1)

	collectorMetric := families["netscaler_exporter_collector_success"].Metric[0]
	if got := labelsByName(collectorMetric)["collector"]; got != "rewrite_policies" {
		t.Fatalf("collector label = %q, want %q", got, "rewrite_policies")
	}

	duration := singleMetric(t, families, "netscaler_exporter_scrape_duration_seconds").GetGauge().GetValue()
	if duration < 0 || math.IsNaN(duration) || math.IsInf(duration, 0) {
		t.Fatalf("scrape duration = %v, want a finite non-negative value", duration)
	}

	buildMetric := singleMetric(t, families, "netscaler_exporter_build_info")
	if got := buildMetric.GetGauge().GetValue(); got != 1 {
		t.Fatalf("build info value = %v, want 1", got)
	}
	wantBuildLabels := map[string]string{
		"version":     "1.2.3",
		"revision":    "abc123",
		"target_type": "adc",
	}
	buildLabels := labelsByName(buildMetric)
	for name, want := range wantBuildLabels {
		if got := buildLabels[name]; got != want {
			t.Errorf("build label %q = %q, want %q", name, got, want)
		}
	}

	fail.Store(true)
	families = gatherMetricFamilies(t, registry)
	assertSingleGaugeValue(t, families, "netscaler_exporter_scrape_success", 0)
	assertSingleGaugeValue(t, families, "netscaler_exporter_collector_success", 0)
}

func TestExporterSelfMetricsOmitDisabledCollectors(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calls.Add(1)
	}))
	t.Cleanup(server.Close)

	exporter := newRewriteTestExporter(t, server.URL, []string{"rewrite_policies"})
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families := gatherMetricFamilies(t, registry)
	if family := families["netscaler_exporter_collector_success"]; family != nil {
		t.Fatalf("disabled collectors emitted %d collector_success metrics, want none", len(family.Metric))
	}
	assertSingleGaugeValue(t, families, "netscaler_exporter_scrape_success", 1)
	if got := calls.Load(); got != 0 {
		t.Fatalf("disabled collector made %d HTTP requests, want none", got)
	}
}

func TestExporterSelfMetricsUseMPSCollectorAndTargetType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nitro/v2/stat/mps_health" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"mps_health":[],"errorcode":0}`)
	}))
	t.Cleanup(server.Close)

	exporter, err := NewExporter(
		&config.Config{Labels: map[string]string{}},
		server.URL,
		"mps",
		"",
		"",
		false,
		"",
		1,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		BuildInfo{Version: "1.2.3", Revision: "abc123"},
	)
	if err != nil {
		t.Fatalf("NewExporter() error = %v", err)
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families := gatherMetricFamilies(t, registry)
	assertSingleGaugeValue(t, families, "netscaler_exporter_collector_success", 1)
	collectorMetric := singleMetric(t, families, "netscaler_exporter_collector_success")
	if got := labelsByName(collectorMetric)["collector"]; got != "mps_health" {
		t.Fatalf("collector label = %q, want %q", got, "mps_health")
	}
	buildMetric := singleMetric(t, families, "netscaler_exporter_build_info")
	if got := labelsByName(buildMetric)["target_type"]; got != "mps" {
		t.Fatalf("target_type label = %q, want %q", got, "mps")
	}
}

func gatherMetricFamilies(t *testing.T, registry *prometheus.Registry) map[string]*dto.MetricFamily {
	t.Helper()
	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Gather() error = %v", err)
	}
	result := make(map[string]*dto.MetricFamily, len(families))
	for _, family := range families {
		result[family.GetName()] = family
	}
	return result
}

func assertSingleGaugeValue(t *testing.T, families map[string]*dto.MetricFamily, name string, want float64) {
	t.Helper()
	metric := singleMetric(t, families, name)
	if got := metric.GetGauge().GetValue(); got != want {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

func singleMetric(t *testing.T, families map[string]*dto.MetricFamily, name string) *dto.Metric {
	t.Helper()
	family := families[name]
	if family == nil {
		t.Fatalf("metric family %q is absent", name)
	}
	if len(family.Metric) != 1 {
		t.Fatalf("metric family %q has %d metrics, want 1", name, len(family.Metric))
	}
	return family.Metric[0]
}
