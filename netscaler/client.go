package netscaler

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Nitro API error codes for session management
const (
	NSERR_SESSION_EXPIRED = 0x1BC // 444 - Session expired
	NSERR_AUTHTIMEOUT     = 0x403 // 1027 - Auth timeout
)

// NitroClient represents the client used to connect to the API.
// It uses session-based authentication with automatic re-login on session expiration.
type NitroClient struct {
	url       string
	username  string
	password  string
	client    *http.Client
	sessionID string
	sessionMu sync.Mutex
	logger    *slog.Logger
}

// NewNitroClient creates a new client used to interact with the Nitro API.
// Uses session-based authentication with automatic re-login on session expiration.
// If caFile is provided, it will be used for TLS verification.
// If ignoreCert is true, TLS verification is skipped entirely.
func NewNitroClient(url string, username string, password string, ignoreCert bool, caFile string, logger *slog.Logger) (*NitroClient, error) {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     30 * time.Second,
	}

	if ignoreCert {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		transport.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	return &NitroClient{
		url:      strings.Trim(url, " /") + "/nitro/v1/",
		username: username,
		password: password,
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		logger: logger,
	}, nil
}

// CloseIdleConnections closes idle connections in the transport pool.
func (c *NitroClient) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

// loginResponse represents the JSON response from the login endpoint.
type loginResponse struct {
	SessionID string `json:"sessionid"`
	ErrorCode int    `json:"errorcode"`
	Message   string `json:"message"`
}

// apiErrorResponse contains the fields used to identify authentication errors.
// NetScaler appliances may return these errors with either a successful or an
// HTTP error status, depending on the appliance version.
type apiErrorResponse struct {
	ErrorCode int `json:"errorcode"`
}

func responseRequiresSessionRefresh(statusCode int, body []byte) bool {
	// A bare 401 always means the session cannot be used. Do not treat every
	// 403 as a session error: Nitro also uses it for valid sessions that lack
	// permission for a resource.
	if statusCode == http.StatusUnauthorized {
		return true
	}

	var apiResp apiErrorResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false
	}

	return apiResp.ErrorCode == NSERR_SESSION_EXPIRED || apiResp.ErrorCode == NSERR_AUTHTIMEOUT
}

// Login authenticates with the Nitro API and stores the session ID.
// If already logged in, this is a no-op.
func (c *NitroClient) Login(ctx context.Context) error {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()
	return c.loginLocked(ctx)
}

// loginLocked authenticates while sessionMu is held.
func (c *NitroClient) loginLocked(ctx context.Context) error {

	// Already logged in
	if c.sessionID != "" {
		return nil
	}

	// No credentials - skip login (for unauthenticated access)
	if c.username == "" {
		return nil
	}

	payload := map[string]interface{}{
		"login": map[string]string{
			"username": c.username,
			"password": c.password,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %w", err)
	}

	url := c.url + "config/login"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	var loginResp loginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	if loginResp.ErrorCode != 0 {
		return fmt.Errorf("login failed: %s (errorcode: %d)", loginResp.Message, loginResp.ErrorCode)
	}

	c.sessionID = loginResp.SessionID
	if c.logger != nil {
		c.logger.Info("session login successful", "url", c.url)
	}
	return nil
}

// Logout clears the session state.
func (c *NitroClient) Logout() {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()
	c.sessionID = ""
}

// HasSession returns true if a session is active.
func (c *NitroClient) HasSession() bool {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()
	return c.sessionID != ""
}

func (c *NitroClient) currentSession() string {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()
	return c.sessionID
}

// refreshSession replaces staleSessionID with a newly authenticated session.
// If another concurrent request already refreshed the session, its result is
// reused instead of creating another login session.
func (c *NitroClient) refreshSession(ctx context.Context, staleSessionID string) error {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()

	if c.sessionID != staleSessionID {
		return nil
	}

	c.sessionID = ""
	return c.loginLocked(ctx)
}

// invalidateSession clears a rejected session without discarding a newer
// session installed concurrently by another request.
func (c *NitroClient) invalidateSession(rejectedSessionID string) {
	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()
	if c.sessionID == rejectedSessionID {
		c.sessionID = ""
	}
}

// get performs a GET request to the Nitro API using session-based auth.
// Automatically handles session expiration by re-logging in.
func (c *NitroClient) get(ctx context.Context, path string, querystring string) ([]byte, error) {
	// Ensure we have a session (or no auth needed)
	if c.username != "" {
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
	}

	return c.doGet(ctx, path, querystring, true)
}

// doGet performs the actual GET request. If retryOnSessionExpiry is true and
// the session has expired, it will re-login and retry once.
func (c *NitroClient) doGet(ctx context.Context, path string, querystring string, retryOnSessionExpiry bool) ([]byte, error) {
	url := c.url + path
	if querystring != "" {
		url = url + "?" + querystring
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Use session cookie if we have one, otherwise no auth. Keep the exact
	// session used by this request so concurrent refreshes cannot invalidate a
	// newer session when this response arrives.
	sessionID := c.currentSession()
	if sessionID != "" {
		req.Header.Set("Cookie", "sessionid="+sessionID)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check the Nitro error body regardless of HTTP status. In particular,
	// NetScaler returns errorcode 444 with HTTP 401 after an HA failover or
	// appliance restart.
	if responseRequiresSessionRefresh(resp.StatusCode, body) {
		if retryOnSessionExpiry {
			if c.logger != nil {
				c.logger.Info("session rejected, re-logging in", "url", c.url, "status", resp.StatusCode)
			}
			if err := c.refreshSession(ctx, sessionID); err != nil {
				return nil, fmt.Errorf("re-login failed: %w", err)
			}
			// Retry once without allowing further retries.
			return c.doGet(ctx, path, querystring, false)
		}

		// Do not retain a session rejected again after the retry. A later scrape
		// can then establish a fresh session instead of remaining wedged.
		c.invalidateSession(sessionID)
		return body, fmt.Errorf("request failed after session refresh: %s (%s)", resp.Status, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("request failed: %s (%s)", resp.Status, string(body))
	}
	return body, nil
}

// GetStats sends a request to the Nitro API and retrieves stats for the given type.
func (c *NitroClient) GetStats(ctx context.Context, statsType string, querystring string) ([]byte, error) {
	return c.get(ctx, "stat/"+statsType, querystring)
}

// GetConfig sends a request to the Nitro API and retrieves configuration for the given type.
func (c *NitroClient) GetConfig(ctx context.Context, configType string, querystring string) ([]byte, error) {
	return c.get(ctx, "config/"+configType, querystring)
}
