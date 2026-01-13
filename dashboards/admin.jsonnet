local g = import 'grafonnet-latest/main.libsonnet';

local var = g.dashboard.variable;
local row = g.panel.row;
local stat = g.panel.stat;
local gauge = g.panel.gauge;
local table = g.panel.table;
local timeSeries = g.panel.timeSeries;
local pieChart = g.panel.pieChart;
local barGauge = g.panel.barGauge;

// Variables
local datasource =
  var.datasource.new('datasource', 'prometheus')
  + var.datasource.generalOptions.withLabel('Data source');

local netscaler =
  var.query.new('netscaler')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('netscaler', 'netscaler_ha_cur_state')
  + var.query.withRefresh('time')
  + var.query.selectionOptions.withMulti(true)
  + var.query.selectionOptions.withIncludeAll(true)
  + var.query.generalOptions.withLabel('NetScaler');

// Helper for prometheus queries
local promQuery(expr, legend='') =
  g.query.prometheus.new('$datasource', expr)
  + g.query.prometheus.withLegendFormat(legend);

// State value mapping (0=DOWN, 1=UP)
local stateMapping = [
  g.panel.stat.standardOptions.mapping.ValueMap.withType()
  + g.panel.stat.standardOptions.mapping.ValueMap.withOptions({
    '0': { text: 'DOWN', color: 'red' },
    '1': { text: 'UP', color: 'green' },
  }),
];

// Thresholds for state (0=red, 1=green)
local stateThresholds =
  g.panel.stat.standardOptions.thresholds.withMode('absolute')
  + g.panel.stat.standardOptions.thresholds.withSteps([
    { color: 'red', value: null },
    { color: 'green', value: 1 },
  ]);

// Thresholds for percent (0-80 green, 80-90 yellow, 90+ red)
local percentThresholds =
  g.panel.gauge.standardOptions.thresholds.withMode('absolute')
  + g.panel.gauge.standardOptions.thresholds.withSteps([
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
// HA Status Row
// ============================================================================
local haClusterState =
  stat.new('HA Cluster State')
  + stat.queryOptions.withTargets([promQuery('netscaler_ha_cur_state{netscaler=~"$netscaler"}', '{{netscaler}}')])
  + stat.standardOptions.withMappings(stateMapping)
  + stateThresholds
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + { gridPos: { w: 3, h: 4, x: 0, y: 1 } };

local haNodeTable =
  table.new('HA Node States')
  + table.queryOptions.withTargets([
    promQuery('netscaler_ha_node_state{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('netscaler_ha_node_status{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'B' },
    promQuery('netscaler_ha_node_sync_state{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'C' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'node_id', mode: 'outer' },
    },
    {
      id: 'organize',
      options: {
        excludeByName: {
          Time: true,
          'Time 1': true,
          'Time 2': true,
          'Time 3': true,
          __name__: true,
          '__name__ 1': true,
          '__name__ 2': true,
          '__name__ 3': true,
          instance: true,
          'instance 1': true,
          'instance 2': true,
          'instance 3': true,
          job: true,
          'job 1': true,
          'job 2': true,
          'job 3': true,
          netscaler: true,
          'netscaler 1': true,
          'netscaler 2': true,
          'netscaler 3': true,
          'node_ip 1': true,
          'node_ip 2': true,
          'node_name 1': true,
          'node_name 2': true,
        },
        indexByName: {
          node_id: 0,
          node_ip: 1,
          'node_ip 3': 1,
          node_name: 2,
          'node_name 3': 2,
          'Value #A': 3,
          'Value #B': 4,
          'Value #C': 5,
        },
        renameByName: {
          node_id: 'Node ID',
          node_ip: 'Node IP',
          'node_ip 3': 'Node IP',
          node_name: 'Node Name',
          'node_name 3': 'Node Name',
          'Value #A': 'Primary',
          'Value #B': 'UP',
          'Value #C': 'Sync OK',
        },
      },
    },
  ])
  + { gridPos: { w: 15, h: 4, x: 9, y: 1 } };

local haMasterDuration =
  stat.new('Master State Duration')
  + stat.queryOptions.withTargets([promQuery('netscaler_ha_node_master_state_seconds{netscaler=~"$netscaler",node_id="0"}', '{{node_name}}')])
  + stat.standardOptions.withUnit('s')
  + stat.options.withColorMode('none')
  + { gridPos: { w: 3, h: 4, x: 3, y: 1 } };

local haSyncFailures =
  stat.new('Sync Failures')
  + stat.queryOptions.withTargets([promQuery('netscaler_ha_sync_failures_total{netscaler=~"$netscaler"}', '{{netscaler}}')])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 1 },
    { color: 'red', value: 5 },
  ])
  + stat.options.withColorMode('background')
  + { gridPos: { w: 3, h: 4, x: 6, y: 1 } };

local haRow =
  row.new('High Availability')
  + { gridPos: { h: 1, w: 24, x: 0, y: 0 } };

// ============================================================================
// SSL Certificates Row
// ============================================================================
local certExpiringSoon =
  stat.new('Expiring < 30 days')
  + stat.queryOptions.withTargets([promQuery('count(netscaler_ssl_cert_days_to_expire{netscaler=~"$netscaler"} < 30) or vector(0)', '')])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + { gridPos: { w: 4, h: 4, x: 0, y: 6 } };

local certCritical =
  stat.new('Critical < 7 days')
  + stat.queryOptions.withTargets([promQuery('count(netscaler_ssl_cert_days_to_expire{netscaler=~"$netscaler"} < 7) or vector(0)', '')])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + { gridPos: { w: 4, h: 4, x: 0, y: 10 } };

local certTable =
  table.new('SSL Certificate Expiry')
  + table.queryOptions.withTargets([
    promQuery('sort(netscaler_ssl_cert_days_to_expire{netscaler=~"$netscaler"})', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'organize',
      options: {
        excludeByName: { Time: true, __name__: true, instance: true, job: true, netscaler: true },
        renameByName: { certkey: 'Certificate', Value: 'Days to Expire' },
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
  + { gridPos: { w: 20, h: 8, x: 4, y: 6 } };

local sslRow =
  row.new('SSL Certificates')
  + { gridPos: { h: 1, w: 24, x: 0, y: 5 } };

// ============================================================================
// System Health Row
// ============================================================================
local mgmtCpu =
  gauge.new('Management CPU')
  + gauge.queryOptions.withTargets([promQuery('netscaler_mgmt_cpu_usage{netscaler=~"$netscaler"}', '{{netscaler}}')])
  + gauge.standardOptions.withUnit('percent')
  + gauge.standardOptions.withMin(0) + gauge.standardOptions.withMax(100)
  + percentThresholds
  + { gridPos: { w: 4, h: 6, x: 0, y: 15 } };

local pktCpu =
  gauge.new('Packet CPU')
  + gauge.queryOptions.withTargets([promQuery('netscaler_pkt_cpu_usage{netscaler=~"$netscaler"}', '{{netscaler}}')])
  + gauge.standardOptions.withUnit('percent')
  + gauge.standardOptions.withMin(0) + gauge.standardOptions.withMax(100)
  + percentThresholds
  + { gridPos: { w: 4, h: 6, x: 4, y: 15 } };

local memUsage =
  gauge.new('Memory Usage')
  + gauge.queryOptions.withTargets([promQuery('netscaler_mem_usage{netscaler=~"$netscaler"}', '{{netscaler}}')])
  + gauge.standardOptions.withUnit('percent')
  + gauge.standardOptions.withMin(0) + gauge.standardOptions.withMax(100)
  + percentThresholds
  + { gridPos: { w: 4, h: 6, x: 8, y: 15 } };

local cpuCores =
  timeSeries.new('CPU Usage per Core')
  + timeSeries.queryOptions.withTargets([promQuery('netscaler_cpu_core_usage_percent{netscaler=~"$netscaler"}', 'Core {{core_id}}')])
  + timeSeries.standardOptions.withUnit('percent')
  + timeSeries.standardOptions.withMin(0) + timeSeries.standardOptions.withMax(100)
  + { gridPos: { w: 12, h: 6, x: 12, y: 15 } };

local systemRow =
  row.new('System Health')
  + { gridPos: { h: 1, w: 24, x: 0, y: 14 } };

// ============================================================================
// Network Interfaces Row
// ============================================================================
local interfaceTraffic =
  timeSeries.new('Interface Traffic')
  + timeSeries.queryOptions.withTargets([
    promQuery('rate(netscaler_interfaces_received_bytes{netscaler=~"$netscaler"}[$__rate_interval])', '{{interface}} RX'),
    promQuery('rate(netscaler_interfaces_transmitted_bytes{netscaler=~"$netscaler"}[$__rate_interval])', '{{interface}} TX'),
  ])
  + timeSeries.standardOptions.withUnit('Bps')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(6);

local interfaceErrors =
  timeSeries.new('Interface Errors')
  + timeSeries.queryOptions.withTargets([promQuery('rate(netscaler_interfaces_error_packets_received{netscaler=~"$netscaler"}[$__rate_interval])', '{{interface}}')])
  + timeSeries.standardOptions.withUnit('pps')
  + timeSeries.gridPos.withW(6) + timeSeries.gridPos.withH(6);

local interfaceJumbo =
  timeSeries.new('Jumbo Packets')
  + timeSeries.queryOptions.withTargets([
    promQuery('rate(netscaler_interfaces_jumbo_packets_received{netscaler=~"$netscaler"}[$__rate_interval])', '{{interface}} RX'),
    promQuery('rate(netscaler_interfaces_jumbo_packets_transmitted{netscaler=~"$netscaler"}[$__rate_interval])', '{{interface}} TX'),
  ])
  + timeSeries.standardOptions.withUnit('pps')
  + timeSeries.gridPos.withW(6) + timeSeries.gridPos.withH(6);

local interfacesRow =
  row.new('Network Interfaces')
  + row.withCollapsed(true)
  + row.withPanels([interfaceTraffic, interfaceErrors, interfaceJumbo]);

// ============================================================================
// Traffic Row
// ============================================================================
local ipTraffic =
  timeSeries.new('IP Traffic')
  + timeSeries.queryOptions.withTargets([
    promQuery('netscaler_ip_rx_mbits_rate{netscaler=~"$netscaler"}', 'RX Mbps'),
    promQuery('netscaler_ip_tx_mbits_rate{netscaler=~"$netscaler"}', 'TX Mbps'),
  ])
  + timeSeries.standardOptions.withUnit('Mbits')
  + timeSeries.standardOptions.withMin(0)
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local httpRequests =
  timeSeries.new('HTTP Requests/Responses')
  + timeSeries.queryOptions.withTargets([
    promQuery('netscaler_http_requests_rate{netscaler=~"$netscaler"}', 'Requests'),
    promQuery('netscaler_http_responses_rate{netscaler=~"$netscaler"}', 'Responses'),
  ])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.standardOptions.withMin(0)
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local httpMethods =
  pieChart.new('HTTP Methods')
  + pieChart.queryOptions.withTargets([
    promQuery('netscaler_http_gets_total{netscaler=~"$netscaler"}', 'GET'),
    promQuery('netscaler_http_posts_total{netscaler=~"$netscaler"}', 'POST'),
    promQuery('netscaler_http_others_total{netscaler=~"$netscaler"}', 'Other'),
  ])
  + pieChart.options.withDisplayLabels(['name', 'percent'])
  + pieChart.options.reduceOptions.withCalcs(['lastNotNull'])
  + pieChart.options.reduceOptions.withFields('')
  + pieChart.options.reduceOptions.withValues(false)
  + pieChart.options.withPieType('pie')
  + pieChart.options.legend.withPlacement('right')
  + pieChart.options.legend.withShowLegend(true)
  + pieChart.gridPos.withW(8) + pieChart.gridPos.withH(6);

local trafficRow =
  row.new('Aggregate Traffic')
  + row.withPanels([ipTraffic, httpRequests, httpMethods]);

// ============================================================================
// Errors Row
// ============================================================================
local httpErrors =
  timeSeries.new('HTTP Errors')
  + timeSeries.queryOptions.withTargets([
    promQuery('rate(netscaler_http_err_server_busy_total{netscaler=~"$netscaler"}[$__rate_interval])', 'Server Busy'),
    promQuery('rate(netscaler_http_err_incomplete_requests_total{netscaler=~"$netscaler"}[$__rate_interval])', 'Incomplete Requests'),
    promQuery('rate(netscaler_http_err_incomplete_responses_total{netscaler=~"$netscaler"}[$__rate_interval])', 'Incomplete Responses'),
  ])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local ipErrors =
  timeSeries.new('IP Errors')
  + timeSeries.queryOptions.withTargets([
    promQuery('rate(netscaler_ip_bad_checksums_total{netscaler=~"$netscaler"}[$__rate_interval])', 'Bad Checksums'),
    promQuery('rate(netscaler_ip_ttl_expired_total{netscaler=~"$netscaler"}[$__rate_interval])', 'TTL Expired'),
    promQuery('rate(netscaler_ip_truncated_packets_total{netscaler=~"$netscaler"}[$__rate_interval])', 'Truncated'),
  ])
  + timeSeries.standardOptions.withUnit('pps')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local tcpErrors =
  timeSeries.new('TCP Errors')
  + timeSeries.queryOptions.withTargets([promQuery('netscaler_tcp_err_ip_port_fail{netscaler=~"$netscaler"}', 'IP Port Fail')])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local errorsRow =
  row.new('Protocol Errors')
  + row.withCollapsed(true)
  + row.withPanels([httpErrors, ipErrors, tcpErrors]);

// ============================================================================
// CS Virtual Servers Row
// ============================================================================
local csStatesTable =
  table.new('CS Virtual Server States')
  + table.queryOptions.withTargets([
    promQuery('netscaler_cs_virtual_servers_state{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('netscaler_cs_virtual_servers_total_hits{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'B' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'virtual_server', mode: 'outer' },
    },
    {
      id: 'organize',
      options: {
        excludeByName: {
          Time: true,
          'Time 1': true,
          'Time 2': true,
          __name__: true,
          '__name__ 1': true,
          '__name__ 2': true,
          instance: true,
          'instance 1': true,
          'instance 2': true,
          job: true,
          'job 1': true,
          'job 2': true,
          netscaler: true,
          'netscaler 1': true,
          'netscaler 2': true,
        },
        indexByName: {
          virtual_server: 0,
          'Value #A': 1,
          'Value #B': 2,
        },
        renameByName: {
          virtual_server: 'CS vServer',
          'Value #A': 'State',
          'Value #B': 'Total Hits',
        },
      },
    },
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Total Hits', desc: true }] },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'Total Hits' },
          properties: [
            { id: 'unit', value: 'short' },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local csTopHits =
  barGauge.new('Top CS vServers by Hits')
  + barGauge.queryOptions.withTargets([promQuery('topk(10, netscaler_cs_virtual_servers_total_hits{netscaler=~"$netscaler"})', '{{virtual_server}}')])
  + barGauge.queryOptions.withTransformations([
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Value', desc: true }] },
    },
  ])
  + barGauge.standardOptions.withUnit('short')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.gridPos.withW(12) + barGauge.gridPos.withH(8);

local csRow =
  row.new('CS Virtual Servers')
  + row.withCollapsed(true)
  + row.withPanels([csStatesTable, csTopHits]);

// ============================================================================
// LB Virtual Servers Row
// ============================================================================
local lbHealthTable =
  table.new('LB Virtual Server Health')
  + table.queryOptions.withTargets([
    promQuery('netscaler_virtual_servers_state{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('netscaler_virtual_servers_health{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'B' },
    promQuery('netscaler_virtual_servers_active_services{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'C' },
    promQuery('netscaler_virtual_servers_inactive_services{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true, refId: 'D' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'virtual_server', mode: 'outer' },
    },
    {
      id: 'organize',
      options: {
        excludeByName: {
          Time: true,
          'Time 1': true,
          'Time 2': true,
          'Time 3': true,
          'Time 4': true,
          __name__: true,
          '__name__ 1': true,
          '__name__ 2': true,
          '__name__ 3': true,
          '__name__ 4': true,
          instance: true,
          'instance 1': true,
          'instance 2': true,
          'instance 3': true,
          'instance 4': true,
          job: true,
          'job 1': true,
          'job 2': true,
          'job 3': true,
          'job 4': true,
          netscaler: true,
          'netscaler 1': true,
          'netscaler 2': true,
          'netscaler 3': true,
          'netscaler 4': true,
        },
        indexByName: {
          virtual_server: 0,
          'Value #A': 1,
          'Value #B': 2,
          'Value #C': 3,
          'Value #D': 4,
        },
        renameByName: {
          virtual_server: 'LB vServer',
          'Value #A': 'State',
          'Value #B': 'Health %',
          'Value #C': 'Active',
          'Value #D': 'Inactive',
        },
      },
    },
  ])
  + table.gridPos.withW(14) + table.gridPos.withH(8);

local lbInactive =
  barGauge.new('LB vServers with Inactive Services')
  + barGauge.queryOptions.withTargets([promQuery('netscaler_virtual_servers_inactive_services{netscaler=~"$netscaler"} > 0', '{{virtual_server}}')])
  + barGauge.queryOptions.withTransformations([
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Value', desc: true }] },
    },
  ])
  + barGauge.standardOptions.withUnit('short')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'yellow', value: null },
    { color: 'red', value: 2 },
  ])
  + barGauge.gridPos.withW(10) + barGauge.gridPos.withH(8);

local lbRow =
  row.new('LB Virtual Servers')
  + row.withCollapsed(true)
  + row.withPanels([lbHealthTable, lbInactive]);

// ============================================================================
// Service Groups Row
// ============================================================================
local sgStatesTable =
  table.new('Service Group Member States')
  + table.queryOptions.withTargets([
    promQuery('netscaler_servicegroup_state{netscaler=~"$netscaler"}', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'organize',
      options: {
        excludeByName: { Time: true, __name__: true, instance: true, job: true, netscaler: true },
        renameByName: {
          servicegroup: 'Service Group',
          member: 'Member',
          port: 'Port',
          Value: 'State',
        },
      },
    },
  ])
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local sgDown =
  table.new('Down Service Group Members')
  + table.queryOptions.withTargets([
    promQuery('netscaler_servicegroup_state{netscaler=~"$netscaler"} == 0', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'organize',
      options: {
        excludeByName: { Time: true, __name__: true, instance: true, job: true, netscaler: true, Value: true },
        renameByName: {
          servicegroup: 'Service Group',
          member: 'Member',
          port: 'Port',
        },
      },
    },
  ])
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local sgTtfb =
  barGauge.new('Top TTFB by Service Group')
  + barGauge.queryOptions.withTargets([promQuery('topk(10, avg by (servicegroup) (netscaler_servicegroup_average_time_to_first_byte{netscaler=~"$netscaler"}))', '{{servicegroup}}')])
  + barGauge.standardOptions.withUnit('ms')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.gridPos.withW(24) + barGauge.gridPos.withH(6);

local sgRow =
  row.new('Service Groups')
  + row.withCollapsed(true)
  + row.withPanels([sgStatesTable, sgDown, sgTtfb]);

// ============================================================================
// Dashboard
// ============================================================================

// HA panels with manual positioning (row + 4 panels in one line)
local haPanels = [
  haRow,
  haClusterState,
  haMasterDuration,
  haSyncFailures,
  haNodeTable,
];

// SSL and System panels with manual positioning
local sslAndSystemPanels = [
  sslRow,
  certExpiringSoon,
  certCritical,
  certTable,
  systemRow,
  mgmtCpu,
  pktCpu,
  memUsage,
  cpuCores,
];

// Remaining rows use grid layout starting at y=21
local otherRows = g.util.grid.makeGrid([
  interfacesRow,
  trafficRow,
  errorsRow,
  csRow,
  lbRow,
  sgRow,
], panelWidth=24, panelHeight=8, startY=21);

g.dashboard.new('NetScaler Admin Dashboard')
+ g.dashboard.withUid('netscaler-admin')
+ g.dashboard.withDescription('Infrastructure overview for NetScaler ADC operators')
+ g.dashboard.withRefresh('30s')
+ g.dashboard.withTimezone('browser')
+ g.dashboard.time.withFrom('now-1h')
+ g.dashboard.time.withTo('now')
+ g.dashboard.graphTooltip.withSharedCrosshair()
+ g.dashboard.withVariables([datasource, netscaler])
+ g.dashboard.withPanels(
  g.util.panel.setPanelIDs(haPanels + sslAndSystemPanels + otherRows)
)
