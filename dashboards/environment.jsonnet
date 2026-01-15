local g = import 'grafonnet-latest/main.libsonnet';

local var = g.dashboard.variable;
local row = g.panel.row;
local stat = g.panel.stat;
local table = g.panel.table;
local timeSeries = g.panel.timeSeries;
local barGauge = g.panel.barGauge;

// Variables
local datasource =
  var.datasource.new('datasource', 'prometheus')
  + var.datasource.generalOptions.withLabel('Data source');

local environment =
  var.query.new('environment')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('deployment_environment_name', 'netscaler_ha_cur_state')
  + var.query.withRefresh('time')
  + var.query.generalOptions.withLabel('Environment');

// Helper for prometheus queries
local promQuery(expr, legend='') =
  g.query.prometheus.new('$datasource', expr)
  + g.query.prometheus.withLegendFormat(legend)
  + g.query.prometheus.withInterval('1m');

// Thresholds for state (0=red, 1=green)
local stateThresholds =
  g.panel.stat.standardOptions.thresholds.withMode('absolute')
  + g.panel.stat.standardOptions.thresholds.withSteps([
    { color: 'red', value: null },
    { color: 'green', value: 1 },
  ]);

// Thresholds for percent (0-80 green, 80-90 yellow, 90+ red)
local percentThresholds =
  g.panel.barGauge.standardOptions.thresholds.withMode('absolute')
  + g.panel.barGauge.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 80 },
    { color: 'red', value: 90 },
  ]);

// Thresholds for cert expiry (red < 7, yellow < 30, green >= 30)
local certThresholds =
  g.panel.stat.standardOptions.thresholds.withMode('absolute')
  + g.panel.stat.standardOptions.thresholds.withSteps([
    { color: 'red', value: null },
    { color: 'yellow', value: 7 },
    { color: 'green', value: 30 },
  ]);

// ============================================================================
// Fleet Overview Row
// ============================================================================
local totalClusters =
  stat.new('Total Clusters')
  + stat.queryOptions.withTargets([
    promQuery('count(count by (netscaler_cluster) (netscaler_ha_cur_state{deployment_environment_name=~"$environment"}))', ''),
  ])
  + stat.options.withColorMode('none')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(6) + stat.gridPos.withH(4) + stat.gridPos.withX(0) + stat.gridPos.withY(1);

local clustersUp =
  stat.new('Clusters UP')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (netscaler_cluster) (netscaler_ha_cur_state{deployment_environment_name=~"$environment"}) == 1)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(6) + stat.gridPos.withH(4) + stat.gridPos.withX(6) + stat.gridPos.withY(1);

local clustersWithIssues =
  stat.new('Clusters with Issues')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (netscaler_cluster) (netscaler_ha_cur_state{deployment_environment_name=~"$environment"}) == 0) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(6) + stat.gridPos.withH(4) + stat.gridPos.withX(12) + stat.gridPos.withY(1);

local certsExpiringSoon =
  stat.new('Certs Expiring < 30d')
  + stat.queryOptions.withTargets([
    promQuery('count(netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment"} < 30) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(6) + stat.gridPos.withH(4) + stat.gridPos.withX(18) + stat.gridPos.withY(1);

local fleetOverviewRow =
  row.new('Fleet Overview')
  + row.gridPos.withY(0);

// ============================================================================
// Cluster Health Table
// ============================================================================
local clusterHealthTable =
  table.new('Cluster Health')
  + table.queryOptions.withTargets([
    promQuery('max by (netscaler_cluster) (netscaler_ha_cur_state{deployment_environment_name=~"$environment"})', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('sum by (netscaler_cluster) (netscaler_ha_sync_failures_total{deployment_environment_name=~"$environment"})', '')
    + { format: 'table', instant: true, refId: 'B' },
    promQuery('count by (netscaler_cluster) (netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment"} < 30)', '')
    + { format: 'table', instant: true, refId: 'C' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'netscaler_cluster', mode: 'outer' },
    },
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(netscaler_cluster|Value #[ABC])$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          netscaler_cluster: 0,
          'Value #A': 1,
          'Value #B': 2,
          'Value #C': 3,
        },
        renameByName: {
          netscaler_cluster: 'Cluster',
          'Value #A': 'Status',
          'Value #B': 'Sync Failures',
          'Value #C': 'Certs < 30d',
        },
      },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'Cluster' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View HA Pair',
                  url: '/d/netscaler-ha-pair?var-environment=${environment}&var-netscaler_cluster=${__value.raw}',
                },
              ],
            },
          ],
        },
        {
          matcher: { id: 'byName', options: 'Status' },
          properties: [
            {
              id: 'mappings',
              value: [
                { type: 'value', options: { '0': { text: 'DOWN', color: 'red' }, '1': { text: 'UP', color: 'green' } } },
              ],
            },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
        {
          matcher: { id: 'byName', options: 'Sync Failures' },
          properties: [
            { id: 'thresholds', value: { mode: 'absolute', steps: [
              { color: 'green', value: null },
              { color: 'yellow', value: 1 },
              { color: 'red', value: 5 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
        {
          matcher: { id: 'byName', options: 'Certs < 30d' },
          properties: [
            { id: 'thresholds', value: { mode: 'absolute', steps: [
              { color: 'green', value: null },
              { color: 'yellow', value: 1 },
              { color: 'red', value: 3 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(24) + table.gridPos.withH(8);

local clusterHealthRow =
  row.new('Cluster Health')
  + row.withPanels([clusterHealthTable]);

// ============================================================================
// Traffic & Workload Row
// ============================================================================
local requestsByCluster =
  timeSeries.new('Requests by Cluster')
  + timeSeries.queryOptions.withTargets([
    promQuery('sum by (netscaler_cluster) (rate(netscaler_cs_virtual_servers_total_requests{deployment_environment_name=~"$environment"}[$__rate_interval])) + sum by (netscaler_cluster) (rate(netscaler_virtual_servers_total_requests{deployment_environment_name=~"$environment"}[$__rate_interval]))', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local throughputByCluster =
  timeSeries.new('Throughput by Cluster')
  + timeSeries.queryOptions.withTargets([
    promQuery('sum by (netscaler_cluster) (rate(netscaler_tcp_rx_bytes_total{deployment_environment_name=~"$environment"}[$__rate_interval]) + rate(netscaler_tcp_tx_bytes_total{deployment_environment_name=~"$environment"}[$__rate_interval]))', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('Bps')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local connectionsByCluster =
  timeSeries.new('Connections by Cluster')
  + timeSeries.queryOptions.withTargets([
    promQuery('sum by (netscaler_cluster) (netscaler_cs_virtual_servers_current_client_connections{deployment_environment_name=~"$environment"}) + sum by (netscaler_cluster) (netscaler_virtual_servers_current_client_connections{deployment_environment_name=~"$environment"})', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local topServices =
  timeSeries.new('Top Services')
  + timeSeries.queryOptions.withTargets([
    promQuery('topk(5, sum by (virtual_server) (rate(netscaler_cs_virtual_servers_total_requests{deployment_environment_name=~"$environment"}[$__rate_interval])))', '{{virtual_server}}'),
    promQuery('topk(5, sum by (virtual_server) (rate(netscaler_virtual_servers_total_requests{deployment_environment_name=~"$environment"}[$__rate_interval])))', '{{virtual_server}}'),
  ])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.standardOptions.withLinks([
    {
      title: 'View Chain',
      url: '/d/netscaler-chain?var-environment=${environment}&var-chain=${__field.labels.virtual_server}',
    },
  ])
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(8);

local httpErrors =
  timeSeries.new('HTTP Errors')
  + timeSeries.queryOptions.withTargets([
    promQuery('sum by (netscaler_cluster) (rate(netscaler_http_err_server_busy_total{deployment_environment_name=~"$environment"}[$__rate_interval]) + rate(netscaler_http_err_incomplete_requests_total{deployment_environment_name=~"$environment"}[$__rate_interval]))', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(8);

local trafficWorkloadRow =
  row.new('Traffic & Workload')
  + row.withPanels([requestsByCluster, throughputByCluster, connectionsByCluster, topServices, httpErrors]);

// ============================================================================
// Resource Utilization Row
// ============================================================================
local memoryUsage =
  timeSeries.new('Memory Usage')
  + timeSeries.queryOptions.withTargets([
    promQuery('max by (netscaler_cluster) (netscaler_mem_usage{deployment_environment_name=~"$environment"})', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('percent')
  + timeSeries.standardOptions.withMin(0) + timeSeries.standardOptions.withMax(100)
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local pktCpuUsage =
  timeSeries.new('Packet CPU')
  + timeSeries.queryOptions.withTargets([
    promQuery('max by (netscaler_cluster) (netscaler_pkt_cpu_usage{deployment_environment_name=~"$environment"})', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('percent')
  + timeSeries.standardOptions.withMin(0) + timeSeries.standardOptions.withMax(100)
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local mgmtCpuUsage =
  timeSeries.new('Management CPU')
  + timeSeries.queryOptions.withTargets([
    promQuery('max by (netscaler_cluster) (netscaler_mgmt_cpu_usage{deployment_environment_name=~"$environment"})', '{{netscaler_cluster}}'),
  ])
  + timeSeries.standardOptions.withUnit('percent')
  + timeSeries.standardOptions.withMin(0) + timeSeries.standardOptions.withMax(100)
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(8);

local resourceRow =
  row.new('Resource Utilization')
  + row.withPanels([memoryUsage, pktCpuUsage, mgmtCpuUsage]);

// ============================================================================
// Virtual Server Health Row
// ============================================================================
local lbVserversUp =
  barGauge.new('LB vServers UP')
  + barGauge.queryOptions.withTargets([
    promQuery('count by (netscaler_cluster) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment"} == 1)', '{{netscaler_cluster}}')
    + { instant: true },
  ])
  + barGauge.standardOptions.withUnit('short')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
  ])
  + barGauge.gridPos.withW(8) + barGauge.gridPos.withH(8);

local lbVserversDown =
  barGauge.new('LB vServers DOWN')
  + barGauge.queryOptions.withTargets([
    promQuery('count by (netscaler_cluster) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment"} == 0)', '{{netscaler_cluster}}')
    + { instant: true },
  ])
  + barGauge.standardOptions.withUnit('short')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + barGauge.gridPos.withW(8) + barGauge.gridPos.withH(8);

local csVserversUp =
  barGauge.new('CS vServers UP')
  + barGauge.queryOptions.withTargets([
    promQuery('count by (netscaler_cluster) (netscaler_cs_virtual_servers_state{deployment_environment_name=~"$environment"} == 1)', '{{netscaler_cluster}}')
    + { instant: true },
  ])
  + barGauge.standardOptions.withUnit('short')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
  ])
  + barGauge.gridPos.withW(8) + barGauge.gridPos.withH(8);

local vserverHealthRow =
  row.new('Virtual Server Health')
  + row.withPanels([lbVserversUp, lbVserversDown, csVserversUp]);

// ============================================================================
// Certificates Expiring Row (collapsed)
// ============================================================================
local certsExpiringTable =
  table.new('Certificates Expiring Soon')
  + table.queryOptions.withTargets([
    promQuery('netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment"} < 30', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(netscaler_cluster|certkey|Value)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          netscaler_cluster: 0,
          certkey: 1,
          Value: 2,
        },
        renameByName: {
          netscaler_cluster: 'Cluster',
          certkey: 'Certificate',
          Value: 'Days to Expire',
        },
      },
    },
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Days to Expire', desc: false }] },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'Days to Expire' },
          properties: [
            { id: 'unit', value: 'd' },
            { id: 'thresholds', value: { mode: 'absolute', steps: [
              { color: 'red', value: null },
              { color: 'yellow', value: 7 },
              { color: 'green', value: 30 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(24) + table.gridPos.withH(8);

local certsRow =
  row.new('Certificates Expiring')
  + row.withCollapsed(true)
  + row.withPanels([certsExpiringTable]);

// ============================================================================
// Dashboard
// ============================================================================
// Manually positioned Fleet Overview panels
local fleetPanels = [
  fleetOverviewRow,
  totalClusters,
  clustersUp,
  clustersWithIssues,
  certsExpiringSoon,
];

// Grid-positioned remaining rows (starting at y=5, after fleet row header + stat panels)
local gridRows = g.util.grid.makeGrid([
  clusterHealthRow,
  trafficWorkloadRow,
  resourceRow,
  vserverHealthRow,
  certsRow,
], panelWidth=24, panelHeight=8, startY=5);

g.dashboard.new('NetScaler Environment')
+ g.dashboard.withUid('netscaler-environment')
+ g.dashboard.withDescription('Fleet overview of all NetScaler HA clusters in an environment')
+ g.dashboard.withRefresh('1m')
+ g.dashboard.withTimezone('browser')
+ g.dashboard.time.withFrom('now-1h')
+ g.dashboard.time.withTo('now')
+ g.dashboard.graphTooltip.withSharedCrosshair()
+ g.dashboard.withVariables([datasource, environment])
+ g.dashboard.withPanels(
  g.util.panel.setPanelIDs(fleetPanels + gridRows)
)
