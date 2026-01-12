package netscaler

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// MPSClient represents the client used to connect to the Citrix ADM (MPS) Nitro v2 API.
type MPSClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewMPSClient creates a new client for interacting with the Citrix ADM (MPS) Nitro v2 API.
// Uses stateless Basic Auth for each request.
func NewMPSClient(url string, username string, password string, ignoreCert bool, caFile string) (*MPSClient, error) {
	transport := &http.Transport{
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

	return &MPSClient{
		url:      strings.Trim(url, " /") + "/nitro/v2/",
		username: username,
		password: password,
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}, nil
}

// CloseIdleConnections closes idle connections in the transport pool.
func (c *MPSClient) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

// get performs a GET request to the MPS Nitro v2 API with Basic Auth.
func (c *MPSClient) get(ctx context.Context, path string, querystring string) ([]byte, error) {
	url := c.url + path
	if querystring != "" {
		url = url + "?" + querystring
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
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

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("request failed: %s (%s)", resp.Status, string(body))
	}
	return body, nil
}

// GetStats sends a request to the MPS Nitro v2 API and retrieves stats for the given type.
func (c *MPSClient) GetStats(ctx context.Context, statsType string, querystring string) ([]byte, error) {
	return c.get(ctx, "stat/"+statsType, querystring)
}
