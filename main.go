package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/elohmeier/netscaler-exporter/collector"
	"github.com/elohmeier/netscaler-exporter/config"
)

var (
	app     = "Citrix-NetScaler-Exporter"
	version string
	build   string
)

func main() {
	var (
		url             string
		targetType      string
		labelsStr       string
		disabledModules string
		bindPort        int
		parallelism     int
		showVersion     bool
		debug           bool
	)

	flag.StringVar(&url, "url", "", "NetScaler URL (e.g., https://netscaler.example.com)")
	flag.StringVar(&targetType, "type", "", "Target type: adc or mps (default: adc)")
	flag.StringVar(&labelsStr, "labels", "", "Custom labels in key=value format, comma-separated (e.g., env=prod,dc=us-east)")
	flag.StringVar(&disabledModules, "disabled-modules", "", "Comma-separated list of modules to disable")
	flag.IntVar(&bindPort, "bind-port", 9280, "Port to bind the exporter endpoint to")
	flag.IntVar(&parallelism, "parallelism", 5, "Maximum concurrent API requests")
	flag.BoolVar(&showVersion, "version", false, "Display application version")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s v%s build %s\n", app, version, build)
		os.Exit(0)
	}

	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	})).With("app", app, "version", "v"+version, "build", build)

	// URL: CLI flag takes precedence over env var
	if url == "" {
		url = config.GetURL()
	}
	if url == "" {
		logger.Error("URL is required (use -url flag or NETSCALER_URL env var)")
		flag.Usage()
		os.Exit(1)
	}

	// Type: CLI flag takes precedence over env var, default to "adc"
	if targetType == "" {
		targetType = config.GetType()
	}
	if targetType == "" {
		targetType = "adc"
	}
	if targetType != "adc" && targetType != "mps" {
		logger.Error("invalid target type (must be 'adc' or 'mps')", "type", targetType)
		os.Exit(1)
	}

	// Parse labels: env var provides base, CLI flag extends/overrides
	envLabels := config.ParseLabels(config.GetLabels())
	cliLabels := config.ParseLabels(labelsStr)
	labels := envLabels
	for k, v := range cliLabels {
		labels[k] = v
	}

	// Parse disabled modules: env var provides base, CLI flag extends
	envDisabled := config.ParseDisabledModules(config.GetDisabledModules())
	cliDisabled := config.ParseDisabledModules(disabledModules)
	disabled := append(envDisabled, cliDisabled...)

	cfg := &config.Config{
		Labels:          labels,
		DisabledModules: disabled,
	}

	// Get credentials from environment
	username, password, err := config.GetCredentials()
	if err != nil {
		logger.Error("failed to get credentials", "err", err)
		os.Exit(1)
	}

	ignoreCert := config.GetIgnoreCert()
	if ignoreCert {
		logger.Info("TLS certificate verification disabled")
	}

	caFile := config.GetCAFile()
	if caFile != "" {
		logger.Info("using custom CA file", "path", caFile)
	}

	logger.Info("starting exporter", "url", url, "type", targetType, "labels", len(labels), "disabled_modules", len(disabled))

	// Create exporter
	exporter, err := collector.NewExporter(cfg, url, targetType, username, password, ignoreCert, caFile, parallelism, logger)
	if err != nil {
		logger.Error("failed to create exporter", "err", err)
		os.Exit(1)
	}

	// Register with default Prometheus registry
	prometheus.MustRegister(exporter)

	// Setup HTTP handlers
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app + " - /metrics for Prometheus metrics"))
	})

	listenAddr := ":" + strconv.Itoa(bindPort)
	logger.Info("starting server", "addr", listenAddr)

	srv := &http.Server{
		Addr:              listenAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}
