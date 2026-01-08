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
		configFile   string
		configInline string
		bindPort     int
		parallelism  int
		showVersion  bool
		debug        bool
	)

	flag.StringVar(&configFile, "config", "", "Path to YAML/JSON configuration file")
	flag.StringVar(&configInline, "config-inline", "", "Inline YAML/JSON configuration")
	flag.IntVar(&bindPort, "bind-port", 9280, "Port to bind the exporter endpoint to")
	flag.IntVar(&parallelism, "parallelism", 5, "Maximum concurrent API requests per target")
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

	// Load configuration
	var cfg *config.Config
	var err error

	if configFile != "" && configInline != "" {
		logger.Error("cannot specify both -config and -config-inline")
		os.Exit(1)
	}

	if configFile != "" {
		cfg, err = config.LoadFile(configFile)
	} else if configInline != "" {
		cfg, err = config.Parse(configInline)
	} else {
		logger.Error("must specify either -config or -config-inline")
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		logger.Error("failed to load configuration", "err", err)
		os.Exit(1)
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

	logger.Info("loaded configuration", "targets", len(cfg.Targets))

	// Create exporter with all targets
	exporter, err := collector.NewExporter(cfg, username, password, ignoreCert, caFile, parallelism, logger)
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
