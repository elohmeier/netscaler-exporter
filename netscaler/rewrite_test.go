package netscaler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRewriteConfigGetters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nitro/v1/config/lbvserver_rewritepolicy_binding":
			if got := r.URL.Query().Get("bulkbindings"); got != "yes" {
				t.Errorf("bulkbindings = %q, want yes", got)
			}
			io.WriteString(w, `{"lbvserver_rewritepolicy_binding":[{"name":"public-apps-vserver","policyname":"catalog-host-rewrite","priority":"100","bindpoint":"REQUEST"}],"errorcode":0}`)
		case "/nitro/v1/config/rewritepolicy":
			io.WriteString(w, `{"rewritepolicy":[{"name":"catalog-host-rewrite","rule":"HTTP.REQ.URL.STARTSWITH(\"/catalog\")","action":"set-catalog-host"}],"errorcode":0}`)
		case "/nitro/v1/config/rewriteaction":
			io.WriteString(w, `{"rewriteaction":[{"name":"set-catalog-host","type":"replace","target":"HTTP.REQ.HEADER(\"Host\")","stringbuilderexpr":"\"catalog.apps.example.com\""}],"errorcode":0}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewNitroClient(server.URL, "", "", false, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("NewNitroClient() error = %v", err)
	}

	bindings, err := GetAllLBVServerRewritePolicyBindings(context.Background(), client)
	if err != nil {
		t.Fatalf("GetAllLBVServerRewritePolicyBindings() error = %v", err)
	}
	if len(bindings) != 1 || bindings[0].Name != "public-apps-vserver" || bindings[0].BindPoint != "REQUEST" {
		t.Fatalf("unexpected bindings: %#v", bindings)
	}

	policies, err := GetAllRewritePolicies(context.Background(), client)
	if err != nil {
		t.Fatalf("GetAllRewritePolicies() error = %v", err)
	}
	if len(policies) != 1 || policies[0].Action != "set-catalog-host" {
		t.Fatalf("unexpected policies: %#v", policies)
	}

	actions, err := GetAllRewriteActions(context.Background(), client)
	if err != nil {
		t.Fatalf("GetAllRewriteActions() error = %v", err)
	}
	if len(actions) != 1 || actions[0].StringBuilderExpr != `"catalog.apps.example.com"` {
		t.Fatalf("unexpected actions: %#v", actions)
	}
}

func TestRewriteConfigGetterPropagatesHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not permitted", http.StatusForbidden)
	}))
	t.Cleanup(server.Close)

	client, err := NewNitroClient(server.URL, "", "", false, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("NewNitroClient() error = %v", err)
	}

	if _, err := GetAllRewritePolicies(context.Background(), client); err == nil {
		t.Fatal("GetAllRewritePolicies() error = nil, want HTTP error")
	}
}
