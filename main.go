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
		targetsFile  string
		targets      string
		bindPort     int
		showVersion  bool
		debug        bool
	)

	flag.StringVar(&targetsFile, "targets-file", "", "Path to YAML/JSON file containing target configurations")
	flag.StringVar(&targets, "targets", "", "Inline YAML/JSON target configuration")
	flag.IntVar(&bindPort, "bind-port", 9280, "Port to bind the exporter endpoint to")
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

	if targetsFile != "" && targets != "" {
		logger.Error("cannot specify both -targets-file and -targets")
		os.Exit(1)
	}

	if targetsFile != "" {
		cfg, err = config.LoadFile(targetsFile)
	} else if targets != "" {
		cfg, err = config.Parse(targets)
	} else {
		logger.Error("must specify either -targets-file or -targets")
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		logger.Error("failed to load configuration", "err", err)
		os.Exit(1)
	}

	logger.Info("loaded configuration", "targets", len(cfg.Targets))

	// Create exporter with all targets
	exporter, err := collector.NewExporter(cfg.Targets, logger)
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
