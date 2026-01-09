package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Config holds the exporter configuration.
type Config struct {
	Labels          map[string]string
	DisabledModules []string
}

// IsModuleDisabled returns true if the given module name is in the disabled list.
func (c *Config) IsModuleDisabled(name string) bool {
	for _, m := range c.DisabledModules {
		if m == name {
			return true
		}
	}
	return false
}

// LabelKeys returns the sorted list of label keys.
func (c *Config) LabelKeys() []string {
	keys := make([]string, 0, len(c.Labels))
	for k := range c.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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

// GetURL reads the URL from environment variable.
func GetURL() string {
	return os.Getenv("NETSCALER_URL")
}

// GetType reads the target type from environment variable.
func GetType() string {
	return os.Getenv("NETSCALER_TYPE")
}

// GetLabels reads labels from environment variable.
func GetLabels() string {
	return os.Getenv("NETSCALER_LABELS")
}

// GetDisabledModules reads disabled modules from environment variable.
func GetDisabledModules() string {
	return os.Getenv("NETSCALER_DISABLED_MODULES")
}

// ParseLabels parses a comma-separated key=value string into a map.
func ParseLabels(labelsStr string) map[string]string {
	labels := make(map[string]string)
	if labelsStr == "" {
		return labels
	}

	pairs := strings.Split(labelsStr, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				labels[key] = value
			}
		}
	}
	return labels
}

// ParseDisabledModules parses a comma-separated list of module names.
func ParseDisabledModules(modulesStr string) []string {
	if modulesStr == "" {
		return nil
	}

	parts := strings.Split(modulesStr, ",")
	modules := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			modules = append(modules, part)
		}
	}
	return modules
}
