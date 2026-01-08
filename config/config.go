package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// Config holds the full exporter configuration.
type Config struct {
	Labels  map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Targets []Target          `yaml:"targets" json:"targets"`
}

// Target represents a single NetScaler instance to scrape.
type Target struct {
	URL    string            `yaml:"url" json:"url"`
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

// LoadFile loads configuration from a YAML or JSON file.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	return Parse(string(data))
}

// Parse parses configuration from a YAML or JSON string.
func Parse(data string) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("no targets configured")
	}

	for i, t := range cfg.Targets {
		if t.URL == "" {
			return nil, fmt.Errorf("target %d: url is required", i)
		}
	}

	return &cfg, nil
}

// GetCredentials reads credentials from environment variables.
func GetCredentials() (username, password string, err error) {
	username = os.Getenv("NETSCALER_USERNAME")
	if username == "" {
		return "", "", fmt.Errorf("NETSCALER_USERNAME environment variable is required")
	}

	password = os.Getenv("NETSCALER_PASSWORD")
	if password == "" {
		return "", "", fmt.Errorf("NETSCALER_PASSWORD environment variable is required")
	}

	return username, password, nil
}

// GetIgnoreCert reads the ignore cert setting from environment variable.
func GetIgnoreCert() bool {
	val := strings.ToLower(os.Getenv("NETSCALER_IGNORE_CERT"))
	return val == "true" || val == "1"
}

// GetCAFile reads the CA file path from environment variable.
func GetCAFile() string {
	return os.Getenv("NETSCALER_CA_FILE")
}

// MergedLabels returns the target's labels merged with global labels.
// Target labels override global labels with the same key.
func (t *Target) MergedLabels(global map[string]string) map[string]string {
	result := make(map[string]string, len(global)+len(t.Labels))
	for k, v := range global {
		result[k] = v
	}
	for k, v := range t.Labels {
		result[k] = v
	}
	return result
}

// LabelKeys returns the sorted list of all label keys from global and all targets.
func (c *Config) LabelKeys() []string {
	keys := make(map[string]struct{})
	for k := range c.Labels {
		keys[k] = struct{}{}
	}
	for _, t := range c.Targets {
		for k := range t.Labels {
			keys[k] = struct{}{}
		}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}

	// Sort for consistent ordering
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}
