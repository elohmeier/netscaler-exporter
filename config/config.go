package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// Config holds the full exporter configuration.
type Config struct {
	// Global defaults - apply to all targets unless overridden
	Username     string `yaml:"username,omitempty" json:"username,omitempty"`
	Password     string `yaml:"password,omitempty" json:"password,omitempty"`
	PasswordEnv  string `yaml:"passwordEnv,omitempty" json:"passwordEnv,omitempty"`
	PasswordFile string `yaml:"passwordFile,omitempty" json:"passwordFile,omitempty"`
	IgnoreCert   bool   `yaml:"ignoreCert,omitempty" json:"ignoreCert,omitempty"`

	Targets []Target `yaml:"targets" json:"targets"`
}

// Target represents a single NetScaler instance to scrape.
type Target struct {
	URL             string `yaml:"url" json:"url"`
	Name            string `yaml:"name" json:"name"`
	Username        string `yaml:"username,omitempty" json:"username,omitempty"`
	Password        string `yaml:"password,omitempty" json:"password,omitempty"`
	PasswordEnv     string `yaml:"passwordEnv,omitempty" json:"passwordEnv,omitempty"`
	PasswordFile    string `yaml:"passwordFile,omitempty" json:"passwordFile,omitempty"`
	IgnoreCert      *bool  `yaml:"ignoreCert,omitempty" json:"ignoreCert,omitempty"`
	CollectTopology bool   `yaml:"collectTopology,omitempty" json:"collectTopology,omitempty"`
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

	// Resolve global password first
	globalPassword, err := resolvePassword(cfg.Password, cfg.PasswordFile, cfg.PasswordEnv)
	if err != nil {
		return nil, fmt.Errorf("global config: %w", err)
	}

	for i := range cfg.Targets {
		cfg.Targets[i].applyDefaults(&cfg, globalPassword)
		if err := cfg.Targets[i].validate(); err != nil {
			return nil, fmt.Errorf("target %q: %w", cfg.Targets[i].Name, err)
		}
	}

	return &cfg, nil
}

// applyDefaults applies global defaults to target fields that aren't set.
func (t *Target) applyDefaults(cfg *Config, globalPassword string) {
	if t.Username == "" {
		t.Username = cfg.Username
	}
	if t.IgnoreCert == nil {
		t.IgnoreCert = &cfg.IgnoreCert
	}

	// Resolve target password, falling back to global
	password, _ := resolvePassword(t.Password, t.PasswordFile, t.PasswordEnv)
	if password == "" {
		password = globalPassword
	}
	t.Password = password
	t.PasswordFile = ""
	t.PasswordEnv = ""
}

// resolvePassword resolves a password from direct value, file, or env var.
func resolvePassword(password, passwordFile, passwordEnv string) (string, error) {
	if password != "" {
		return password, nil
	}
	if passwordFile != "" {
		data, err := os.ReadFile(passwordFile)
		if err != nil {
			return "", fmt.Errorf("reading password file: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	if passwordEnv != "" {
		return os.Getenv(passwordEnv), nil
	}
	return "", nil
}

func (t *Target) validate() error {
	if t.URL == "" {
		return fmt.Errorf("url is required")
	}
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if t.Username == "" {
		return fmt.Errorf("username is required")
	}
	if t.Password == "" {
		return fmt.Errorf("password is required (set password, passwordFile, or passwordEnv)")
	}
	return nil
}

// GetIgnoreCert returns the ignoreCert value, defaulting to false if not set.
func (t *Target) GetIgnoreCert() bool {
	if t.IgnoreCert == nil {
		return false
	}
	return *t.IgnoreCert
}
