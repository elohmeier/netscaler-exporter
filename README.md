# NetScaler Exporter

Prometheus exporter for Citrix NetScaler (ADC) metrics via the Nitro API.

## Features

- Multi-target support with concurrent metric collection
- Flexible labels for target identification
- Topology metrics for service graph visualization

## Quick Start

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NETSCALER_USERNAME` | API username | Yes |
| `NETSCALER_PASSWORD` | API password | Yes |
| `NETSCALER_IGNORE_CERT` | Skip TLS verification (`true` or `1`) | No |

### Binary

```bash
export NETSCALER_USERNAME=nsroot
export NETSCALER_PASSWORD=secret
./netscaler-exporter -config config.yaml
```

### Docker

```bash
docker run -p 9280:9280 \
  -e NETSCALER_USERNAME=nsroot \
  -e NETSCALER_PASSWORD=secret \
  -v ./config.yaml:/config.yaml \
  ghcr.io/elohmeier/netscaler-exporter -config /config.yaml
```

## Configuration

The config file only contains targets and optional labels. All operational settings (credentials, TLS) are via environment variables.

```yaml
# Global labels (applied to all targets)
labels:
  environment: production
  datacenter: us-east

# Targets to monitor
targets:
  - url: https://netscaler1.example.com/nitro/v1
    labels:
      cluster: primary

  - url: https://netscaler2.example.com/nitro/v1
    labels:
      cluster: secondary
```

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Path to YAML/JSON config file | |
| `-config-inline` | Inline YAML/JSON configuration | |
| `-bind-port` | HTTP server port | 9280 |
| `-debug` | Enable debug logging | false |

## Endpoints

| Path | Description |
|------|-------------|
| `/metrics` | Prometheus metrics |
| `/health` | Health check (returns 200 OK) |

## Metrics

All metrics include the `ns_instance` label (target URL) plus any custom labels defined in config.

### Categories

- **System**: CPU, memory, storage, network throughput
- **Virtual Servers**: State, health, requests, connections, traffic
- **Services**: State, throughput, connections, response times
- **Service Groups**: Member state and traffic
- **GSLB**: Global server load balancing metrics
- **Content Switching**: CS virtual server statistics
- **VPN**: VPN virtual server metrics
- **AAA**: Authentication metrics
- **Interfaces**: Per-interface traffic statistics
- **Topology**: Node and edge metrics for service graph visualization

## License

MIT
