# NetScaler Exporter

Prometheus exporter for Citrix NetScaler (ADC) and Citrix ADM (MPS) metrics via the Nitro API.

## Features

- Support for both ADC (NetScaler) and MPS (Citrix ADM) targets
- Flexible custom labels for metric identification
- Topology metrics for service graph visualization
- Static LB request Host rewrites for routing correlation
- Configurable module disabling for unsupported collectors

## Quick Start

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NETSCALER_URL` | NetScaler URL (e.g., `https://netscaler.example.com`) | Yes (or use `-url` flag) |
| `NETSCALER_USERNAME` | API username | Yes |
| `NETSCALER_PASSWORD` | API password | Yes |
| `NETSCALER_TYPE` | Target type: `adc` or `mps` | No (default: `adc`) |
| `NETSCALER_IGNORE_CERT` | Skip TLS verification (`true` or `1`) | No |
| `NETSCALER_CA_FILE` | Path to custom CA certificate file | No |
| `NETSCALER_LABELS` | Base labels (format: `key1=val1,key2=val2`), merged with `-labels` flag | No |
| `NETSCALER_DISABLED_MODULES` | Base disabled modules (comma-separated), merged with `-disabled-modules` flag | No |

### Binary

```bash
export NETSCALER_URL=https://netscaler.example.com
export NETSCALER_USERNAME=nsroot
export NETSCALER_PASSWORD=secret
./netscaler-exporter
```

With CLI flags:

```bash
./netscaler-exporter \
  -url https://netscaler.example.com \
  -type adc \
  -labels "env=prod,dc=us-east" \
  -disabled-modules "ns_capacity,ssl_certs"
```

### Docker

```bash
docker run -p 9280:9280 \
  -e NETSCALER_URL=https://netscaler.example.com \
  -e NETSCALER_USERNAME=nsroot \
  -e NETSCALER_PASSWORD=secret \
  ghcr.io/elohmeier/netscaler-exporter
```

For MPS (Citrix ADM):

```bash
docker run -p 9280:9280 \
  -e NETSCALER_URL=https://adm.example.com \
  -e NETSCALER_USERNAME=admin \
  -e NETSCALER_PASSWORD=secret \
  -e NETSCALER_TYPE=mps \
  ghcr.io/elohmeier/netscaler-exporter
```

## Configuration

All configuration is via CLI flags and environment variables. CLI flags override environment variables.

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-url` | NetScaler URL (overrides `NETSCALER_URL`) | |
| `-type` | Target type: `adc` or `mps` (overrides `NETSCALER_TYPE`) | `adc` |
| `-labels` | Custom labels (format: `key1=val1,key2=val2`) | |
| `-disabled-modules` | Modules to disable (comma-separated) | |
| `-bind-port` | HTTP server port | 9280 |
| `-parallelism` | Maximum concurrent API requests | 5 |
| `-debug` | Enable debug logging | false |
| `-version` | Display application version | |

### Disabling Modules

Use `-disabled-modules` to skip collectors that aren't supported by your device:

| Module | Description |
|--------|-------------|
| `ns_stats` | System stats (CPU, memory, network) |
| `ns_license` | License/model info |
| `ns_capacity` | Bandwidth capacity stats |
| `interfaces` | Network interface metrics |
| `virtual_servers` | LB virtual servers |
| `services` | Backend services |
| `service_groups` | Service groups |
| `gslb_services` | GSLB services |
| `gslb_vservers` | GSLB virtual servers |
| `cs_vservers` | Content switching virtual servers |
| `vpn_vservers` | VPN virtual servers |
| `aaa_stats` | Authentication stats |
| `topology` | Topology relationships |
| `rewrite_policies` | Static request-side Host rewrites on LB virtual servers |
| `protocol_http` | HTTP protocol stats |
| `protocol_tcp` | TCP protocol stats |
| `protocol_ip` | IP protocol stats |
| `ssl_stats` | SSL global stats |
| `ssl_certs` | SSL certificates |
| `ssl_vservers` | SSL virtual servers |
| `system_cpu` | Per-core CPU stats |
| `ha_stats` | High availability node stats |
| `mps_health` | MPS health stats (MPS targets only) |

## Endpoints

| Path | Description |
|------|-------------|
| `/metrics` | Prometheus metrics |
| `/health` | Health check (returns 200 OK) |

## Metrics

Target metrics include any custom labels defined via `-labels`. Exporter self-metrics use only the labels shown below:

| Metric | Description |
|--------|-------------|
| `netscaler_exporter_scrape_success` | Whether every enabled collector succeeded (`1`) or any failed (`0`) |
| `netscaler_exporter_collector_success{collector="..."}` | Whether an enabled collector succeeded (`1`) or failed (`0`); disabled collectors are absent |
| `netscaler_exporter_scrape_duration_seconds` | Total duration of the exporter scrape |
| `netscaler_exporter_build_info{version,revision,target_type}` | Exporter build and target information (always `1`) |

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
- **Rewrite policies**: Static request-side Host replacements bound to LB virtual servers

### HTTP Host Rewrite Mapping

The `rewrite_policies` module resolves LB vServer rewrite bindings through their policies and actions. It exports only request-side `replace` actions whose target is the HTTP `Host` header and whose replacement is a static hostname:

```text
netscaler_lb_vserver_http_host_rewrite_info{
  virtual_server="public-apps-vserver",
  policy="catalog-host-rewrite",
  priority="100",
  rule="HTTP.REQ.HEADER(\"Host\").EQ(\"legacy.example.com\") && HTTP.REQ.URL.STARTSWITH(\"/catalog\")",
  action="set-catalog-host",
  host="catalog.apps.example.com"
} 1
```

The normalized `host` label can be joined with an OpenShift route information metric when both are available in the same Prometheus-compatible datasource:

```promql
netscaler_lb_vserver_http_host_rewrite_info
  * on (host) group_left(namespace, route)
    max by (host, namespace, route) (openshift_route_info)
```

Dynamic rewrite expressions and rewrites of headers other than `Host` are intentionally omitted. Disable the collector with `-disabled-modules rewrite_policies` if the ADC version or API user does not support the rewrite configuration endpoints.

### MPS (Citrix ADM) Metrics

- **Health**: CPU usage, memory usage/free/total, disk usage/free/total/used

MPS requests use stateless `X-NITRO-USER`/`X-NITRO-PASS` authentication. This
avoids consuming the limited per-user ADM login sessions and remains usable when
the interactive session limit has already been reached.

## License

MIT
