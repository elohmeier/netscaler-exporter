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

local netscalerCluster =
  var.query.new('netscaler_cluster')
  + var.query.withDatasourceFromVariable(datasource)
  + var.query.queryTypes.withLabelValues('netscaler_cluster', 'netscaler_ha_cur_state{deployment_environment_name=~"$environment"}')
  + var.query.withRefresh('time')
  + var.query.selectionOptions.withMulti(true)
  + var.query.selectionOptions.withIncludeAll(true)
  + var.query.generalOptions.withLabel('Cluster');

// Helper for prometheus queries
local promQuery(expr, legend='') =
  g.query.prometheus.new('$datasource', expr)
  + g.query.prometheus.withLegendFormat(legend)
  + g.query.prometheus.withInterval('1m');

// ============================================================================
// Problem Summary Row
// ============================================================================
local downVservers =
  stat.new('DOWN vServers')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (virtual_server, netscaler_cluster) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(0) + stat.gridPos.withY(1);

local downCsVservers =
  stat.new('DOWN CS vServers')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (virtual_server, netscaler_cluster) (netscaler_cs_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(4) + stat.gridPos.withY(1);

local degradedVservers =
  stat.new('Degraded vServers')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (virtual_server, netscaler_cluster) (netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0 and netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 100)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 1 },
    { color: 'red', value: 5 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(8) + stat.gridPos.withY(1);

local downServiceGroupMembers =
  stat.new('DOWN SG Members')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (servicegroup, member, port, netscaler_cluster) (netscaler_servicegroup_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(12) + stat.gridPos.withY(1);

local certsCritical =
  stat.new('Certs < 30 days')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (certkey, netscaler_cluster) (netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 30 and netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'red', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(16) + stat.gridPos.withY(1);

local certsWarning =
  stat.new('Certs 30-90 days')
  + stat.queryOptions.withTargets([
    promQuery('count(max by (certkey, netscaler_cluster) (netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} >= 30 and netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 90)) or vector(0)', ''),
  ])
  + stat.standardOptions.thresholds.withMode('absolute')
  + stat.standardOptions.thresholds.withSteps([
    { color: 'green', value: null },
    { color: 'yellow', value: 1 },
  ])
  + stat.options.withColorMode('background')
  + stat.options.withGraphMode('none')
  + stat.gridPos.withW(4) + stat.gridPos.withH(4) + stat.gridPos.withX(20) + stat.gridPos.withY(1);

local summaryRow =
  row.new('Problem Summary')
  + row.gridPos.withY(0);

// ============================================================================
// DOWN vServers Row
// ============================================================================
local downLbVserversTable =
  table.new('DOWN LB Virtual Servers')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(virtual_server|netscaler_cluster|deployment_environment_name)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          virtual_server: 2,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          virtual_server: 'LB vServer',
        },
      },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'LB vServer' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View Chain',
                  url: '/d/netscaler-chain?var-environment=${__data.fields.Environment}&var-chain=${__value.raw}',
                },
              ],
            },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local downCsVserversTable =
  table.new('DOWN CS Virtual Servers')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_cs_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(virtual_server|netscaler_cluster|deployment_environment_name)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          virtual_server: 2,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          virtual_server: 'CS vServer',
        },
      },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'CS vServer' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View Chain',
                  url: '/d/netscaler-chain?var-environment=${__data.fields.Environment}&var-chain=${__value.raw}',
                },
              ],
            },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local downVserversRow =
  row.new('DOWN Virtual Servers')
  + row.withPanels([downLbVserversTable, downCsVserversTable]);

// ============================================================================
// Degraded Health Row
// ============================================================================
local degradedHealthTable =
  table.new('vServers with Degraded Health')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0 and netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 100)', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_active_services{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} and on(virtual_server, netscaler_cluster) (netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0 and netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 100))', '')
    + { format: 'table', instant: true, refId: 'B' },
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_inactive_services{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} and on(virtual_server, netscaler_cluster) (netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0 and netscaler_virtual_servers_health{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 100))', '')
    + { format: 'table', instant: true, refId: 'C' },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'joinByField',
      options: { byField: 'virtual_server', mode: 'outer' },
    },
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(virtual_server|netscaler_cluster|deployment_environment_name|Value #[ABC])$' },
      },
    },
    {
      id: 'organize',
      options: {
        excludeByName: {
          'netscaler_cluster 1': true,
          'netscaler_cluster 2': true,
          'deployment_environment_name 1': true,
          'deployment_environment_name 2': true,
        },
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          virtual_server: 2,
          'Value #A': 3,
          'Value #B': 4,
          'Value #C': 5,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          virtual_server: 'vServer',
          'Value #A': 'Health %',
          'Value #B': 'Active',
          'Value #C': 'Inactive',
        },
      },
    },
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Health %', desc: false }] },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'Health %' },
          properties: [
            { id: 'unit', value: 'percent' },
            { id: 'thresholds', value: { mode: 'absolute', steps: [
              { color: 'red', value: null },
              { color: 'yellow', value: 50 },
              { color: 'green', value: 100 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
        {
          matcher: { id: 'byName', options: 'Inactive' },
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
          matcher: { id: 'byName', options: 'vServer' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View Chain',
                  url: '/d/netscaler-chain?var-environment=${__data.fields.Environment}&var-chain=${__value.raw}',
                },
              ],
            },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(24) + table.gridPos.withH(10);

local degradedHealthRow =
  row.new('Degraded Health')
  + row.withPanels([degradedHealthTable]);

// ============================================================================
// DOWN Service Group Members Row
// ============================================================================
local downSgMembersTable =
  table.new('DOWN Service Group Members')
  + table.queryOptions.withTargets([
    promQuery('max by (servicegroup, member, port, netscaler_cluster, deployment_environment_name) (netscaler_servicegroup_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(servicegroup|member|port|netscaler_cluster|deployment_environment_name)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          servicegroup: 2,
          member: 3,
          port: 4,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          servicegroup: 'Service Group',
          member: 'Member',
          port: 'Port',
        },
      },
    },
  ])
  + table.gridPos.withW(24) + table.gridPos.withH(10);

local downSgMembersRow =
  row.new('DOWN Service Group Members')
  + row.withPanels([downSgMembersTable]);

// ============================================================================
// SSL Certificates Row
// ============================================================================
local certsCriticalTable =
  table.new('SSL Certificates Expiring < 30 days')
  + table.queryOptions.withTargets([
    promQuery('max by (certkey, netscaler_cluster, deployment_environment_name) (netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 30 and netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 0)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(certkey|netscaler_cluster|deployment_environment_name|Value)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          certkey: 2,
          Value: 3,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
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
              { color: 'orange', value: 14 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local certsWarningTable =
  table.new('SSL Certificates Expiring 30-90 days')
  + table.queryOptions.withTargets([
    promQuery('max by (certkey, netscaler_cluster, deployment_environment_name) (netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} >= 30 and netscaler_ssl_cert_days_to_expire{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} < 90)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(certkey|netscaler_cluster|deployment_environment_name|Value)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          certkey: 2,
          Value: 3,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
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
              { color: 'yellow', value: null },
              { color: 'green', value: 60 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(12) + table.gridPos.withH(8);

local certsRow =
  row.new('SSL Certificate Expiry')
  + row.withPanels([certsCriticalTable, certsWarningTable]);

// ============================================================================
// Potentially Unused vServers Row (collapsed)
// ============================================================================
local zeroHitsTable =
  table.new('vServers with Zero Hits (potentially unused)')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_total_hits{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0) and on(virtual_server, netscaler_cluster) max by (virtual_server, netscaler_cluster) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 1)', '')
    + { format: 'table', instant: true },
  ])
  + table.queryOptions.withTransformations([
    {
      id: 'filterFieldsByName',
      options: {
        include: { pattern: '^(virtual_server|netscaler_cluster|deployment_environment_name)$' },
      },
    },
    {
      id: 'organize',
      options: {
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          virtual_server: 2,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          virtual_server: 'vServer (UP but zero hits)',
        },
      },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'vServer (UP but zero hits)' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View Chain',
                  url: '/d/netscaler-chain?var-environment=${__data.fields.Environment}&var-chain=${__value.raw}',
                },
              ],
            },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(24) + table.gridPos.withH(10);

local unusedRow =
  row.new('Potentially Unused Configuration')
  + row.withCollapsed(true)
  + row.withPanels([zeroHitsTable]);

// ============================================================================
// High Inactive Services Row (collapsed)
// ============================================================================
local highInactiveTable =
  table.new('vServers with Many Inactive Services')
  + table.queryOptions.withTargets([
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_inactive_services{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 2)', '')
    + { format: 'table', instant: true, refId: 'A' },
    promQuery('max by (virtual_server, netscaler_cluster, deployment_environment_name) (netscaler_virtual_servers_active_services{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} and on(virtual_server, netscaler_cluster) (netscaler_virtual_servers_inactive_services{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} > 2))', '')
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
        include: { pattern: '^(virtual_server|netscaler_cluster|deployment_environment_name|Value #[AB])$' },
      },
    },
    {
      id: 'organize',
      options: {
        excludeByName: {
          'netscaler_cluster 1': true,
          'deployment_environment_name 1': true,
        },
        indexByName: {
          deployment_environment_name: 0,
          netscaler_cluster: 1,
          virtual_server: 2,
          'Value #B': 3,
          'Value #A': 4,
        },
        renameByName: {
          deployment_environment_name: 'Environment',
          netscaler_cluster: 'Cluster',
          virtual_server: 'vServer',
          'Value #B': 'Active',
          'Value #A': 'Inactive',
        },
      },
    },
    {
      id: 'sortBy',
      options: { sort: [{ field: 'Inactive', desc: true }] },
    },
  ])
  + {
    fieldConfig: {
      defaults: {},
      overrides: [
        {
          matcher: { id: 'byName', options: 'vServer' },
          properties: [
            {
              id: 'links',
              value: [
                {
                  title: 'View Chain',
                  url: '/d/netscaler-chain?var-environment=${__data.fields.Environment}&var-chain=${__value.raw}',
                },
              ],
            },
          ],
        },
        {
          matcher: { id: 'byName', options: 'Inactive' },
          properties: [
            { id: 'thresholds', value: { mode: 'absolute', steps: [
              { color: 'yellow', value: null },
              { color: 'orange', value: 5 },
              { color: 'red', value: 10 },
            ] } },
            { id: 'custom.cellOptions', value: { type: 'color-background' } },
          ],
        },
      ],
    },
  }
  + table.gridPos.withW(24) + table.gridPos.withH(10);

local highInactiveRow =
  row.new('High Inactive Services')
  + row.withCollapsed(true)
  + row.withPanels([highInactiveTable]);

// ============================================================================
// Problem Trends Row (collapsed)
// ============================================================================
local downVserversTrend =
  timeSeries.new('DOWN vServers Over Time')
  + timeSeries.queryOptions.withTargets([
    promQuery('count(max by (virtual_server, netscaler_cluster) (netscaler_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', 'LB DOWN'),
    promQuery('count(max by (virtual_server, netscaler_cluster) (netscaler_cs_virtual_servers_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', 'CS DOWN'),
  ])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.standardOptions.color.withMode('palette-classic')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(8);

local downServicesTrend =
  timeSeries.new('DOWN Service Group Members Over Time')
  + timeSeries.queryOptions.withTargets([
    promQuery('count(max by (servicegroup, member, port, netscaler_cluster) (netscaler_servicegroup_state{deployment_environment_name=~"$environment",netscaler_cluster=~"$netscaler_cluster"} == 0)) or vector(0)', 'DOWN Members'),
  ])
  + timeSeries.standardOptions.withUnit('short')
  + timeSeries.standardOptions.color.withMode('palette-classic')
  + timeSeries.gridPos.withW(12) + timeSeries.gridPos.withH(8);

local trendsRow =
  row.new('Problem Trends')
  + row.withCollapsed(true)
  + row.withPanels([downVserversTrend, downServicesTrend]);

// ============================================================================
// Dashboard
// ============================================================================
// Manually positioned summary panels
local summaryPanels = [
  summaryRow,
  downVservers,
  downCsVservers,
  degradedVservers,
  downServiceGroupMembers,
  certsCritical,
  certsWarning,
];

// Grid-positioned rows (starting at y=5, after summary row header + stat panels)
local gridRows = g.util.grid.makeGrid([
  downVserversRow,
  degradedHealthRow,
  downSgMembersRow,
  certsRow,
  unusedRow,
  highInactiveRow,
  trendsRow,
], panelWidth=24, panelHeight=10, startY=5);

g.dashboard.new('NetScaler Problems')
+ g.dashboard.withUid('netscaler-problems')
+ g.dashboard.withDescription('Problem detection dashboard for identifying misconfigurations, down services, and issues requiring attention')
+ g.dashboard.withRefresh('1m')
+ g.dashboard.withTimezone('browser')
+ g.dashboard.time.withFrom('now-1h')
+ g.dashboard.time.withTo('now')
+ g.dashboard.graphTooltip.withSharedCrosshair()
+ g.dashboard.withVariables([datasource, environment, netscalerCluster])
+ g.dashboard.withPanels(
  g.util.panel.setPanelIDs(summaryPanels + gridRows)
)
