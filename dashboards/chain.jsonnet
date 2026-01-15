local g = import 'grafonnet-latest/main.libsonnet';

local var = g.dashboard.variable;
local row = g.panel.row;
local stat = g.panel.stat;
local table = g.panel.table;
local timeSeries = g.panel.timeSeries;
local barGauge = g.panel.barGauge;
local nodeGraph = g.panel.nodeGraph;

// Variables
local datasource =
  var.datasource.new('datasource', 'prometheus')
  + var.datasource.generalOptions.withLabel('Data source');

local environment =
  var.query.new('environment')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('deployment_environment_name', 'netscaler_topology_node{chain!=""}')
  + var.query.withRefresh('time')
  + var.query.generalOptions.withLabel('Environment');

// Chain variable - gets all unique chain values from topology nodes
// Chains can start with either CS vservers or standalone LB vservers
local chain =
  var.query.new('chain')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('chain', 'netscaler_topology_node{deployment_environment_name=~"$environment",chain!=""}')
  + var.query.withRefresh('time')
  + var.query.selectionOptions.withMulti(true)
  + var.query.generalOptions.withLabel('Chain')
  + { allowCustomValue: true };

// LB vservers in the selected chain
local lbvserver =
  var.query.new('lbvserver')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('title', 'netscaler_topology_node{deployment_environment_name=~"$environment",chain=~".*$chain.*",node_type="lbvserver"}')
  + var.query.withRefresh('time')
  + var.query.selectionOptions.withMulti(true)
  + var.query.selectionOptions.withIncludeAll(true)
  + var.query.generalOptions.withLabel('LB vServer');

// Service groups in the selected chain
local servicegroup =
  var.query.new('servicegroup')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('title', 'netscaler_topology_node{deployment_environment_name=~"$environment",chain=~".*$chain.*",node_type="servicegroup"}')
  + var.query.withRefresh('time')
  + var.query.selectionOptions.withMulti(true)
  + var.query.selectionOptions.withIncludeAll(true)
  + var.query.generalOptions.withLabel('Service Group');

// Helper for prometheus queries
local promQuery(expr, legend='') =
  g.query.prometheus.new('$datasource', expr)
  + g.query.prometheus.withLegendFormat(legend)
  + g.query.prometheus.withInterval('1m');

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

// ============================================================================
// Topology Row
// ============================================================================

// The node graph panel uses:
// - 'id' for node identification
// - 'title' for display name
// - 'subtitle' for descriptive stats shown below the title
//
// Subtitle content by node type:
// - LB vserver: "Health: X%, Conns: Y"
// - CS vserver: "Conns: Y"
// - ServiceGroup: "Avg TTFB: Xms, Members: Y/Z"
// - Server: "TTFB: Xms, Conns: Y"
local topologyGraph =
  nodeGraph.new('Chain Topology')
  + nodeGraph.panelOptions.withDescription('Routing topology: CS vServer -> LB vServer -> Service Group -> Server')
  + nodeGraph.queryOptions.withTargets([
    // Nodes query
    promQuery('netscaler_topology_node{chain=~".*$chain.*"}', '')
    + { format: 'table', instant: true, refId: 'nodes' },
    // Edges query
    promQuery('netscaler_topology_edge{chain=~".*$chain.*"}', '')
    + { format: 'table', instant: true, refId: 'edges' },
  ])
  + {
    options: {
      edges: {},
      nodes: {},
      layoutAlgorithm: 'layered',
      zoomMode: 'cooperative',
    },
  }
  + { gridPos: { h: 24, w: 12, x: 0, y: 1 } };

local topologyRow =
  row.new('Topology')
  + { gridPos: { h: 1, w: 24, x: 0, y: 0 } };

// ============================================================================
// Chain Health Row
// ============================================================================
// Root state - tries CS first, falls back to LB (for chains that start with LB)
local chainRootState =
  stat.new('Chain Root State')
  + stat.queryOptions.withTargets([
    // Try CS vserver state first
    promQuery('netscaler_cs_virtual_servers_state{virtual_server=~"$chain"}', 'CS: {{virtual_server}}'),
    // Also show LB vserver if that's the root
    promQuery('netscaler_virtual_servers_state{virtual_server=~"$chain"}', 'LB: {{virtual_server}}'),
  ])
  + stat.standardOptions.withMappings(stateMapping)
  + stateThresholds
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + { gridPos: { h: 4, w: 12, x: 12, y: 1 } };

local chainHits =
  stat.new('Total Hits')
  + stat.queryOptions.withTargets([
    promQuery('netscaler_cs_virtual_servers_total_hits{virtual_server=~"$chain"}', 'CS'),
    promQuery('netscaler_virtual_servers_total_hits{virtual_server=~"$chain"}', 'LB'),
  ])
  + stat.standardOptions.withUnit('short')
  + stat.options.withColorMode('none')
  + stat.options.withGraphMode('area')
  + { gridPos: { h: 4, w: 6, x: 12, y: 5 } };

local chainConnections =
  stat.new('Active Connections')
  + stat.queryOptions.withTargets([
    promQuery('netscaler_cs_virtual_servers_current_client_connections{virtual_server=~"$chain"}', 'CS Client'),
    promQuery('netscaler_virtual_servers_current_client_connections{virtual_server=~"$chain"}', 'LB Client'),
  ])
  + stat.standardOptions.withUnit('short')
  + stat.options.withColorMode('none')
  + stat.options.withGraphMode('area')
  + { gridPos: { h: 4, w: 6, x: 18, y: 5 } };

local chainComponentsTable =
  table.new('Chain Components')
  + table.queryOptions.withTargets([
    promQuery('netscaler_topology_node{chain=~".*$chain.*"}', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(title|node_type|state)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          node_type: 0,
          title: 1,
          state: 2,
        },
        renameByName: {
          title: 'Component',
          node_type: 'Type',
          state: 'State',
        },
      },
    },
  ])
  + { gridPos: { h: 16, w: 12, x: 12, y: 9 } };

// ============================================================================
// LB Virtual Servers Row
// ============================================================================
local lbStatesTable =
  table.new('LB vServer States')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server) (netscaler_virtual_servers_state{virtual_server=~"$lbvserver"})', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('max by (virtual_server) (netscaler_virtual_servers_health{virtual_server=~"$lbvserver"})', '')
    + { format: 'table', instant: true, refId: 'B' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'virtual_server', mode: 'outer' },
    },
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(virtual_server|Value #[AB])$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          virtual_server: 0,
          'Value #A': 1,
          'Value #B': 2,
        },
        renameByName: {
          virtual_server: 'LB vServer',
          'Value #A': 'State',
          'Value #B': 'Health %',
        },
      },
    },
  ])
  + table.gridPos.withW(8) + table.gridPos.withH(6);

local lbHealthGauge =
  barGauge.new('LB vServer Health')
  + barGauge.queryOptions.withTargets([promQuery('max by (virtual_server) (netscaler_virtual_servers_health{virtual_server=~"$lbvserver"})', '{{virtual_server}}')])
  + barGauge.standardOptions.withUnit('percent')
  + barGauge.standardOptions.withMin(0) + barGauge.standardOptions.withMax(100)
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'red', value: null },
    { color: 'yellow', value: 50 },
    { color: 'green', value: 100 },
  ])
  + barGauge.gridPos.withW(8) + barGauge.gridPos.withH(6);

local lbActiveInactive =
  timeSeries.new('Active / Inactive Services')
  + timeSeries.queryOptions.withTargets([
    promQuery('max by (virtual_server) (netscaler_virtual_servers_active_services{virtual_server=~"$lbvserver"})', '{{virtual_server}} Active'),
    promQuery('max by (virtual_server) (netscaler_virtual_servers_inactive_services{virtual_server=~"$lbvserver"})', '{{virtual_server}} Inactive'),
  ])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.gridPos.withW(8) + timeSeries.gridPos.withH(6);

local lbRequests =
  timeSeries.new('LB Requests')
  + timeSeries.queryOptions.withTargets([promQuery('max by (virtual_server) (rate(netscaler_virtual_servers_total_requests{virtual_server=~"$lbvserver"}[$__rate_interval]))', '{{virtual_server}}')])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(6);

local lbTraffic =
  timeSeries.new('LB Traffic')
  + timeSeries.queryOptions.withTargets([
    promQuery('max by (virtual_server) (rate(netscaler_virtual_servers_total_request_bytes{virtual_server=~"$lbvserver"}[$__rate_interval]))', '{{virtual_server}} RX'),
    promQuery('max by (virtual_server) (rate(netscaler_virtual_servers_total_response_bytes{virtual_server=~"$lbvserver"}[$__rate_interval]))', '{{virtual_server}} TX'),
  ])
  + timeSeries.standardOptions.withUnit('Bps')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(6);

local lbRow =
  row.new('LB Virtual Servers')
  + row.withPanels([lbStatesTable, lbHealthGauge, lbActiveInactive, lbRequests, lbTraffic]);

// ============================================================================
// Service Groups Row
// ============================================================================
local sgMembersTable =
  table.new('Service Group Members')
  + table.queryOptions.withTargets([
    promQuery('max by (servicegroup, member, port) (netscaler_servicegroup_state{servicegroup=~"$servicegroup"})', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('max by (servicegroup, member, port) (netscaler_servicegroup_average_time_to_first_byte{servicegroup=~"$servicegroup"})', '')
    + { format: 'table', instant: true, refId: 'B' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'member', mode: 'outer' },
    },
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(servicegroup|member|port|Value #[AB])$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          servicegroup: 0,
          member: 1,
          port: 2,
          'Value #A': 3,
          'Value #B': 4,
        },
        renameByName: {
          servicegroup: 'Service Group',
          member: 'Member',
          port: 'Port',
          'Value #A': 'State',
          'Value #B': 'TTFB (ms)',
        },
      },
    },
  ])
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local sgTtfb =
  barGauge.new('TTFB per Member')
  + barGauge.queryOptions.withTargets([promQuery('max by (servicegroup, member, port) (netscaler_servicegroup_average_time_to_first_byte{servicegroup=~"$servicegroup"})', '{{member}}:{{port}}')])
  + barGauge.standardOptions.withUnit('ms')
  + barGauge.options.withDisplayMode('gradient')
  + barGauge.options.withOrientation('horizontal')
  + barGauge.standardOptions.thresholds.withMode('absolute')
  + barGauge.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 100 },
    { color: 'red', value: 500 },
  ])
  + barGauge.gridPos.withW(12) + barGauge.gridPos.withH(8);

local sgRequests =
  timeSeries.new('Member Requests')
  + timeSeries.queryOptions.withTargets([promQuery('max by (servicegroup, member, port) (rate(netscaler_servicegroup_total_requests{servicegroup=~"$servicegroup"}[$__rate_interval]))', '{{member}}:{{port}}')])
  + timeSeries.standardOptions.withUnit('reqps')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(6);

local sgConnections =
  timeSeries.new('Member Connections')
  + timeSeries.queryOptions.withTargets([promQuery('max by (servicegroup, member, port) (netscaler_servicegroup_current_server_connections{servicegroup=~"$servicegroup"})', '{{member}}:{{port}}')])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(6);

local sgRow =
  row.new('Service Groups')
  + row.withCollapsed(true)
  + row.withPanels([sgMembersTable, sgTtfb, sgRequests, sgConnections]);

// ============================================================================
// Backend Servers Row
// ============================================================================
local serverStatesTable =
  table.new('Backend Server States')
  + table.queryOptions.withTargets([
    promQuery('max by (title, state) (netscaler_topology_node{chain=~".*$chain.*",node_type="server"})', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(title|state)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          title: 0,
          state: 1,
        },
        renameByName: {
          title: 'Server',
          state: 'State',
        },
      },
    },
  ])
  + table.gridPos.withW(12) + table.gridPos.withH(6);

local serverHealthSummary =
  stat.new('Servers UP')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (title) (netscaler_topology_node{chain=~".*$chain.*",node_type="server",state="UP"})) or vector(0)', 'UP'),
    promQuery('count(max by (title) (netscaler_topology_node{chain=~".*$chain.*",node_type="server"})) or vector(0)', 'Total'),
  ])
  + stat.options.withColorMode('none')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(6) + stat.gridPos.withH(6);

local serverDown =
  table.new('Down Servers')
  + table.queryOptions.withTargets([
    promQuery('max by (title) (netscaler_topology_node{chain=~".*$chain.*",node_type="server",state="DOWN"})', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^title$' },
      },
    },
    {
      id: 'organize',
      options: {
        renameByName: { title: 'Server' },
      },
    },
  ])
  + table.standardOptions.color.withMode('fixed')
  + table.standardOptions.color.withFixedColor('red')
  + table.gridPos.withW(6) + table.gridPos.withH(6);

local serverRow =
  row.new('Backend Servers')
  + row.withCollapsed(true)
  + row.withPanels([serverStatesTable, serverHealthSummary, serverDown]);

// ============================================================================
// Dashboard
// ============================================================================

// Topology panels with manual positioning (node graph left, stats+table right)
local topologyPanels = [
  topologyRow,
  topologyGraph,
  chainRootState,
  chainHits,
  chainConnections,
  chainComponentsTable,
];

// Other rows use grid layout starting at y=25 (after topology section)
local otherRows = g.util.grid.makeGrid([
  lbRow,
  sgRow,
  serverRow,
], panelWidth=24, startY=25);

g.dashboard.new('NetScaler Chain')
+ g.dashboard.withUid('netscaler-chain')
+ g.dashboard.withDescription('Chain-focused view for NetScaler consumers - select chains to see routing topology and health')
+ g.dashboard.withRefresh('1m')
+ g.dashboard.withTimezone('browser')
+ g.dashboard.time.withFrom('now-1h')
+ g.dashboard.time.withTo('now')
+ g.dashboard.graphTooltip.withSharedCrosshair()
+ g.dashboard.withVariables([datasource, environment, chain, lbvserver, servicegroup])
+ g.dashboard.withPanels(
  g.util.panel.setPanelIDs(topologyPanels + otherRows)
)
