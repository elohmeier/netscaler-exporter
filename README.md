# NetScaler Exporter

Prometheus exporter for Citrix NetScaler (ADC) and Citrix ADM (MPS) metrics via the Nitro API.

## Features

- Multi-target support with concurrent metric collection
- Support for both ADC (NetScaler) and MPS (Citrix ADM) targets
- Flexible labels for target identification
- Topology metrics for service graph visualization

## Quick Start

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NETSCALER_USERNAME` | API username | Yes |
| `NETSCALER_PASSWORD` | API password | Yes |
| `NETSCALER_IGNORE_CERT` | Skip TLS verification (`true` or `1`) | No |
| `NETSCALER_CA_FILE` | Path to custom CA certificate file | No |

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

The config file contains targets and optional labels. All operational settings (credentials, TLS) are via environment variables.

```yaml
# Global labels (applied to all targets)
labels:
  environment: production
  datacenter: us-east

# ADC (NetScaler) targets
adc_targets:
  - url: https://netscaler1.example.com/nitro/v1
    labels:
      cluster: primary

  - url: https://netscaler2.example.com/nitro/v1
    labels:
      cluster: secondary

# MPS (Citrix ADM) targets (optional)
mps_targets:
  - url: https://adm.example.com/nitro/v1
    labels:
      region: us-east
```

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Path to YAML/JSON config file | |
| `-config-inline` | Inline YAML/JSON configuration | |
| `-bind-port` | HTTP server port | 9280 |
| `-parallelism` | Maximum concurrent API requests per target | 5 |
| `-debug` | Enable debug logging | false |
| `-version` | Display application version | |

## Endpoints

| Path | Description |
|------|-------------|
| `/metrics` | Prometheus metrics |
| `/health` | Health check (returns 200 OK) |

## Metrics

All metrics include the `ns_instance` label (target URL) plus any custom labels defined in config.

### ADC Metrics

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

### MPS (Citrix ADM) Metrics

- **Health**: CPU usage, memory usage/free/total, disk usage/free/total/used

## License

MIT
