package netscaler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const sessionExpiredBody = `{ "errorcode": 444, "message": "Session expired or killed. Please login again", "severity": "ERROR" }`

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestResponseRequiresSessionRefresh(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		body       string
		want       bool
	}{
		{name: "HTTP 401", statusCode: http.StatusUnauthorized, body: `{}`, want: true},
		{name: "permission denied with HTTP 403", statusCode: http.StatusForbidden, body: `{}`, want: false},
		{name: "session expired with HTTP 403", statusCode: http.StatusForbidden, body: sessionExpiredBody, want: true},
		{name: "session expired with HTTP 200", statusCode: http.StatusOK, body: sessionExpiredBody, want: true},
		{name: "auth timeout with HTTP 400", statusCode: http.StatusBadRequest, body: `{"errorcode":1027}`, want: true},
		{name: "unrelated API error", statusCode: http.StatusBadRequest, body: `{"errorcode":123}`, want: false},
		{name: "invalid JSON", statusCode: http.StatusInternalServerError, body: `not JSON`, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := responseRequiresSessionRefresh(tt.statusCode, []byte(tt.body)); got != tt.want {
				t.Fatalf("responseRequiresSessionRefresh(%d, %q) = %v, want %v", tt.statusCode, tt.body, got, tt.want)
			}
		})
	}
}

func TestNitroClientRetriesSessionExpiredWithHTTP401(t *testing.T) {
	t.Parallel()

	var loginRequests atomic.Int32
	var statsRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/nitro/v1/config/login":
			loginNumber := loginRequests.Add(1)
			fmt.Fprintf(w, `{"sessionid":"session-%d","errorcode":0}`, loginNumber)
		case r.Method == http.MethodGet && r.URL.Path == "/nitro/v1/stat/ns":
			statsRequests.Add(1)
			switch cookie, _ := r.Cookie("sessionid"); {
			case cookie != nil && cookie.Value == "session-1":
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, sessionExpiredBody)
			case cookie != nil && cookie.Value == "session-2":
				fmt.Fprint(w, `{"ns":{"mgmtcpuusagepcnt":1},"errorcode":0}`)
			default:
				t.Errorf("unexpected session cookie: %v", cookie)
				w.WriteHeader(http.StatusUnauthorized)
			}
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewNitroClient(server.URL, "monitor", "secret", false, "", discardLogger())
	if err != nil {
		t.Fatalf("NewNitroClient() error = %v", err)
	}

	body, err := client.GetStats(context.Background(), "ns", "")
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}
	if !strings.Contains(string(body), `"errorcode":0`) {
		t.Fatalf("GetStats() body = %s", body)
	}
	if got := loginRequests.Load(); got != 2 {
		t.Fatalf("login requests = %d, want 2", got)
	}
	if got := statsRequests.Load(); got != 2 {
		t.Fatalf("stats requests = %d, want 2", got)
	}
	if got := client.currentSession(); got != "session-2" {
		t.Fatalf("current session = %q, want session-2", got)
	}
}

func TestNitroClientRefreshesRejectedSessionOnlyOnceConcurrently(t *testing.T) {
	t.Parallel()

	const requestCount = 8
	var loginRequests atomic.Int32
	var staleRequests atomic.Int32
	var freshRequests atomic.Int32
	var releaseOnce sync.Once
	allStaleRequestsStarted := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/nitro/v1/config/login":
			loginRequests.Add(1)
			fmt.Fprint(w, `{"sessionid":"fresh","errorcode":0}`)
		case r.Method == http.MethodGet && r.URL.Path == "/nitro/v1/stat/ns":
			cookie, _ := r.Cookie("sessionid")
			if cookie != nil && cookie.Value == "stale" {
				if staleRequests.Add(1) == requestCount {
					releaseOnce.Do(func() { close(allStaleRequestsStarted) })
				}
				select {
				case <-allStaleRequestsStarted:
				case <-r.Context().Done():
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, sessionExpiredBody)
				return
			}
			if cookie != nil && cookie.Value == "fresh" {
				freshRequests.Add(1)
				fmt.Fprint(w, `{"errorcode":0}`)
				return
			}
			t.Errorf("unexpected session cookie: %v", cookie)
			w.WriteHeader(http.StatusUnauthorized)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewNitroClient(server.URL, "monitor", "secret", false, "", discardLogger())
	if err != nil {
		t.Fatalf("NewNitroClient() error = %v", err)
	}
	client.sessionID = "stale"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	errs := make(chan error, requestCount)
	for range requestCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetStats(ctx, "ns", "")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Errorf("GetStats() error = %v", err)
		}
	}

	if got := loginRequests.Load(); got != 1 {
		t.Fatalf("login requests = %d, want 1", got)
	}
	if got := staleRequests.Load(); got != requestCount {
		t.Fatalf("stale requests = %d, want %d", got, requestCount)
	}
	if got := freshRequests.Load(); got != requestCount {
		t.Fatalf("fresh requests = %d, want %d", got, requestCount)
	}
}

func TestNitroClientClearsSessionWhenRetryIsRejected(t *testing.T) {
	t.Parallel()

	var loginRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/nitro/v1/config/login" {
			loginNumber := loginRequests.Add(1)
			fmt.Fprintf(w, `{"sessionid":"session-%d","errorcode":0}`, loginNumber)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/nitro/v1/stat/ns" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, sessionExpiredBody)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client, err := NewNitroClient(server.URL, "monitor", "secret", false, "", discardLogger())
	if err != nil {
		t.Fatalf("NewNitroClient() error = %v", err)
	}

	_, err = client.GetStats(context.Background(), "ns", "")
	if err == nil || !strings.Contains(err.Error(), "request failed after session refresh") {
		t.Fatalf("GetStats() error = %v, want rejected retry error", err)
	}
	if client.HasSession() {
		t.Fatal("client retained the rejected session")
	}
	if got := loginRequests.Load(); got != 2 {
		t.Fatalf("login requests = %d, want 2", got)
	}
}

func TestMPSClientUsesStatelessAuthWhenSessionLimitIsExhausted(t *testing.T) {
	t.Parallel()

	var loginRequests atomic.Int32
	var statsRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/nitro/v2/config/login":
			loginRequests.Add(1)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"errorcode":10020,"message":"Login failed. The user session count has exceeded the maximum limit."}`)
		case r.Method == http.MethodGet && r.URL.Path == "/nitro/v2/stat/mps_health":
			statsRequests.Add(1)
			if got := r.Header.Get("X-NITRO-USER"); got != "monitor" {
				t.Errorf("X-NITRO-USER = %q, want monitor", got)
			}
			if got := r.Header.Get("X-NITRO-PASS"); got != "secret" {
				t.Errorf("X-NITRO-PASS = %q, want secret", got)
			}
			if cookie := r.Header.Get("Cookie"); cookie != "" {
				t.Errorf("Cookie = %q, want no session cookie", cookie)
			}
			fmt.Fprint(w, `{"mps_health":[],"errorcode":0}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client, err := NewMPSClient(server.URL, "monitor", "secret", false, "", discardLogger())
	if err != nil {
		t.Fatalf("NewMPSClient() error = %v", err)
	}

	for range 2 {
		if _, err := client.GetStats(context.Background(), "mps_health", ""); err != nil {
			t.Fatalf("GetStats() error = %v", err)
		}
	}
	if got := loginRequests.Load(); got != 0 {
		t.Fatalf("login requests = %d, want 0", got)
	}
	if got := statsRequests.Load(); got != 2 {
		t.Fatalf("stats requests = %d, want 2", got)
	}
}
