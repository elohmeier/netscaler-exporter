package collector

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestResolveHTTPHostRewrites(t *testing.T) {
	bindings := []netscaler.LBVServerRewritePolicyBinding{
		{Name: "public-apps-vserver", PolicyName: "catalog-host-rewrite", Priority: "100", BindPoint: " request "},
		{Name: "public-apps-vserver", PolicyName: "accounts-host-rewrite", Priority: "110", BindPoint: "REQUEST"},
		{Name: "public-apps-vserver", PolicyName: "response-header-rewrite", Priority: "1", BindPoint: "RESPONSE"},
		{Name: "public-apps-vserver", PolicyName: "dynamic-host-rewrite", Priority: "120", BindPoint: "REQUEST"},
		{Name: "public-apps-vserver", PolicyName: "insert-host-rewrite", Priority: "130", BindPoint: "REQUEST"},
		{Name: "public-apps-vserver", PolicyName: "missing-action-rewrite", Priority: "140", BindPoint: "REQUEST"},
		{Name: "public-apps-vserver", PolicyName: "catalog-host-rewrite", Priority: "100", BindPoint: "REQUEST"},
	}
	policies := []netscaler.RewritePolicy{
		{Name: "catalog-host-rewrite", Rule: `HTTP.REQ.URL.STARTSWITH("/catalog")`, Action: "set-catalog-host"},
		{Name: "accounts-host-rewrite", Rule: `HTTP.REQ.URL.STARTSWITH("/accounts")`, Action: "set-accounts-host"},
		{Name: "response-header-rewrite", Rule: "true", Action: "set-response-header"},
		{Name: "dynamic-host-rewrite", Rule: "true", Action: "set-dynamic-host"},
		{Name: "insert-host-rewrite", Rule: "true", Action: "insert-host"},
		{Name: "missing-action-rewrite", Rule: "true", Action: "missing-action"},
	}
	actions := []netscaler.RewriteAction{
		{Name: "set-catalog-host", Type: " RePlAcE ", Target: `http.req.header( "host" )`, StringBuilderExpr: `"CATALOG.APPS.EXAMPLE.COM:443"`},
		{Name: "set-accounts-host", Type: "replace", Target: `HTTP.REQ.HEADER("Host")`, StringBuilderExpr: `"ACCOUNTS.APPS.EXAMPLE.COM."`},
		{Name: "set-response-header", Type: "replace", Target: `HTTP.RES.HEADER("Host")`, StringBuilderExpr: `"ignored.example.com"`},
		{Name: "set-dynamic-host", Type: "replace", Target: `HTTP.REQ.HEADER("Host")`, StringBuilderExpr: `HTTP.REQ.URL`},
		{Name: "insert-host", Type: "insert_http_header", Target: `HTTP.REQ.HEADER("Host")`, StringBuilderExpr: `"ignored.example.com"`},
	}

	got := resolveHTTPHostRewrites(bindings, policies, actions)
	want := []httpHostRewriteMapping{
		{
			VirtualServer: "public-apps-vserver",
			Policy:        "catalog-host-rewrite",
			Priority:      "100",
			Rule:          `HTTP.REQ.URL.STARTSWITH("/catalog")`,
			Action:        "set-catalog-host",
			Host:          "catalog.apps.example.com",
		},
		{
			VirtualServer: "public-apps-vserver",
			Policy:        "accounts-host-rewrite",
			Priority:      "110",
			Rule:          `HTTP.REQ.URL.STARTSWITH("/accounts")`,
			Action:        "set-accounts-host",
			Host:          "accounts.apps.example.com",
		},
	}

	if len(got) != len(want) {
		t.Fatalf("resolveHTTPHostRewrites() returned %d mappings, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("mapping[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestNormalizeStaticHostExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       string
		wantOK     bool
	}{
		{name: "hostname", expression: `"Route.Example.COM"`, want: "route.example.com", wantOK: true},
		{name: "trailing dot", expression: `"route.example.com."`, want: "route.example.com", wantOK: true},
		{name: "port", expression: `"route.example.com:8443"`, want: "route.example.com", wantOK: true},
		{name: "dynamic", expression: `HTTP.REQ.URL`, wantOK: false},
		{name: "empty", expression: `""`, wantOK: false},
		{name: "path", expression: `"route.example.com/path"`, wantOK: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := normalizeStaticHostExpression(test.expression)
			if got != test.want || ok != test.wantOK {
				t.Errorf("normalizeStaticHostExpression(%q) = (%q, %t), want (%q, %t)", test.expression, got, ok, test.want, test.wantOK)
			}
		})
	}
}

func TestHTTPHostRewriteCollectorMetricAndReset(t *testing.T) {
	var failActions atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nitro/v1/config/lbvserver_rewritepolicy_binding":
			io.WriteString(w, `{"lbvserver_rewritepolicy_binding":[{"name":"public-apps-vserver","policyname":"catalog-host-rewrite","priority":"100","bindpoint":"REQUEST"}],"errorcode":0}`)
		case "/nitro/v1/config/rewritepolicy":
			io.WriteString(w, `{"rewritepolicy":[{"name":"catalog-host-rewrite","rule":"HTTP.REQ.URL.STARTSWITH(\"/catalog\")","action":"set-catalog-host"}],"errorcode":0}`)
		case "/nitro/v1/config/rewriteaction":
			if failActions.Load() {
				http.Error(w, "temporary failure", http.StatusInternalServerError)
				return
			}
			io.WriteString(w, `{"rewriteaction":[{"name":"set-catalog-host","type":"replace","target":"HTTP.REQ.HEADER(\"Host\")","stringbuilderexpr":"\"CATALOG.APPS.EXAMPLE.COM\""}],"errorcode":0}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	exporter := newRewriteTestExporter(t, server.URL, nil)
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	metrics := gatherMetricFamily(t, registry, "netscaler_lb_vserver_http_host_rewrite_info")
	if len(metrics) != 1 {
		t.Fatalf("first scrape returned %d rewrite metrics, want 1", len(metrics))
	}
	labels := labelsByName(metrics[0])
	wantLabels := map[string]string{
		"environment":    "production",
		"virtual_server": "public-apps-vserver",
		"policy":         "catalog-host-rewrite",
		"priority":       "100",
		"rule":           `HTTP.REQ.URL.STARTSWITH("/catalog")`,
		"action":         "set-catalog-host",
		"host":           "catalog.apps.example.com",
	}
	for name, want := range wantLabels {
		if got := labels[name]; got != want {
			t.Errorf("label %q = %q, want %q", name, got, want)
		}
	}
	if got := metrics[0].GetGauge().GetValue(); got != 1 {
		t.Errorf("metric value = %v, want 1", got)
	}

	failActions.Store(true)
	if metrics := gatherMetricFamily(t, registry, "netscaler_lb_vserver_http_host_rewrite_info"); len(metrics) != 0 {
		t.Fatalf("failed scrape retained %d stale rewrite metrics, want 0", len(metrics))
	}
}

func TestRewritePoliciesModuleCanBeDisabled(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		http.Error(w, "unexpected request", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	exporter := newRewriteTestExporter(t, server.URL, []string{"rewrite_policies"})
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)
	gatherMetricFamily(t, registry, "netscaler_lb_vserver_http_host_rewrite_info")
	if got := calls.Load(); got != 0 {
		t.Fatalf("disabled rewrite_policies module made %d API calls, want 0", got)
	}
}

func TestHTTPHostRewriteCollectorSkipsLookupsWithoutRequestBindings(t *testing.T) {
	var followupCalls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/nitro/v1/config/lbvserver_rewritepolicy_binding" {
			io.WriteString(w, `{"lbvserver_rewritepolicy_binding":[{"name":"public-apps-vserver","policyname":"security-response-header","priority":"1","bindpoint":"RESPONSE"}],"errorcode":0}`)
			return
		}
		followupCalls.Add(1)
		http.Error(w, "unexpected follow-up request", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	exporter := newRewriteTestExporter(t, server.URL, nil)
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)
	if metrics := gatherMetricFamily(t, registry, "netscaler_lb_vserver_http_host_rewrite_info"); len(metrics) != 0 {
		t.Fatalf("response-only bindings produced %d rewrite metrics, want 0", len(metrics))
	}
	if got := followupCalls.Load(); got != 0 {
		t.Fatalf("response-only bindings caused %d policy/action API calls, want 0", got)
	}
}

func newRewriteTestExporter(t *testing.T, url string, extraDisabled []string) *Exporter {
	t.Helper()
	disabled := []string{
		"topology", "ns_stats", "ns_license", "interfaces", "virtual_servers", "services",
		"gslb_services", "gslb_vservers", "cs_vservers", "vpn_vservers", "aaa_stats",
		"service_groups", "protocol_http", "protocol_tcp", "protocol_ip", "ssl_stats",
		"ssl_certs", "ssl_vservers", "system_cpu", "ns_capacity", "ha_stats",
	}
	disabled = append(disabled, extraDisabled...)
	exporter, err := NewExporter(
		&config.Config{
			Labels:          map[string]string{"environment": "production"},
			DisabledModules: disabled,
		},
		url,
		"adc",
		"",
		"",
		false,
		"",
		1,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	if err != nil {
		t.Fatalf("NewExporter() error = %v", err)
	}
	return exporter
}

func gatherMetricFamily(t *testing.T, registry *prometheus.Registry, name string) []*dto.Metric {
	t.Helper()
	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Gather() error = %v", err)
	}
	for _, family := range families {
		if family.GetName() == name {
			return family.Metric
		}
	}
	return nil
}

func labelsByName(metric *dto.Metric) map[string]string {
	labels := make(map[string]string, len(metric.Label))
	for _, label := range metric.Label {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}
