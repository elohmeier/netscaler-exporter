package netscaler

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NitroClient represents the client used to connect to the API
type NitroClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewNitroClient creates a new client used to interact with the Nitro API.
// Uses stateless Basic Auth for each request.
func NewNitroClient(url string, username string, password string, ignoreCert bool) *NitroClient {
	transport := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     30 * time.Second,
	}

	if ignoreCert {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
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
	}
}

// CloseIdleConnections closes idle connections in the transport pool.
func (c *NitroClient) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

// get performs a GET request to the Nitro API with Basic Auth.
func (c *NitroClient) get(ctx context.Context, path string, querystring string) ([]byte, error) {
	url := c.url + path
	if querystring != "" {
		url = url + "?" + querystring
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
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
