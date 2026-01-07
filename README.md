# NetScaler Exporter

Prometheus exporter for Citrix NetScaler (ADC) metrics via the Nitro API.

## Features

- Multi-target support with concurrent metric collection
- Flexible authentication (password, environment variable, or file)
- SSL certificate validation toggle per target
- Metrics for system, virtual servers, services, GSLB, content switching, AAA, and interfaces

## Quick Start

### Binary

```bash
./netscaler-exporter -targets-file config.yml
```

### Docker

```bash
docker run -p 9280:9280 -v ./config.yml:/config.yml \
  ghcr.io/elohmeier/netscaler-exporter -targets-file /config.yml
```

## Configuration

```yaml
# Global defaults (optional)
username: nsroot
passwordEnv: NETSCALER_PASSWORD
ignoreCert: true

# Targets
targets:
  - name: prod-ns1
    url: https://netscaler1.example.com/nitro/v1
  - name: prod-ns2
    url: https://netscaler2.example.com/nitro/v1
    username: admin
    passwordFile: /run/secrets/ns2-password
    ignoreCert: false
```

### Configuration Options

| Option | Description |
|--------|-------------|
| `url` | NetScaler Nitro API endpoint (required) |
| `name` | Target identifier (required) |
| `username` | API username |
| `password` | API password |
| `passwordEnv` | Environment variable containing password |
| `passwordFile` | File containing password |
| `ignoreCert` | Skip SSL certificate validation |
| `collectTopology` | Enable topology data collection |

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-targets-file` | Path to YAML/JSON config file | |
| `-targets` | Inline YAML/JSON configuration | |
| `-bind-port` | HTTP server port | 9280 |
| `-debug` | Enable debug logging | false |

## Endpoints

| Path | Description |
|------|-------------|
| `/metrics` | Prometheus metrics |
| `/health` | Health check (returns 200 OK) |

## Metrics

Collected metrics include:

- **System**: CPU, memory, storage, network throughput
- **Virtual Servers**: State, health, requests, connections, traffic
- **Services**: State, throughput, connections, response times
- **Service Groups**: Member state and traffic
- **GSLB**: Global server load balancing metrics
- **Content Switching**: CS virtual server statistics
- **AAA**: Authentication metrics
- **Interfaces**: Per-interface traffic statistics

## License

MIT
