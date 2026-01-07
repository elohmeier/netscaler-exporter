package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/elohmeier/netscaler-exporter/config"
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	config     *config.Config
	username   string
	password   string
	ignoreCert bool
	labelKeys  []string
	logger     *slog.Logger

	// System metrics (descriptors)
	modelID                                *prometheus.Desc
	mgmtCPUUsage                           *prometheus.Desc
	memUsage                               *prometheus.Desc
	pktCPUUsage                            *prometheus.Desc
	flashPartitionUsage                    *prometheus.Desc
	varPartitionUsage                      *prometheus.Desc
	totRxMB                                *prometheus.Desc
	totTxMB                                *prometheus.Desc
	httpRequests                           *prometheus.Desc
	httpResponses                          *prometheus.Desc
	tcpCurrentClientConnections            *prometheus.Desc
	tcpCurrentClientConnectionsEstablished *prometheus.Desc
	tcpCurrentServerConnections            *prometheus.Desc
	tcpCurrentServerConnectionsEstablished *prometheus.Desc

	// Interface metrics
	interfacesRxBytes        *prometheus.GaugeVec
	interfacesTxBytes        *prometheus.GaugeVec
	interfacesRxPackets      *prometheus.GaugeVec
	interfacesTxPackets      *prometheus.GaugeVec
	interfacesJumboPacketsRx *prometheus.GaugeVec
	interfacesJumboPacketsTx *prometheus.GaugeVec
	interfacesErrorPacketsRx *prometheus.GaugeVec

	// Virtual Server metrics
	virtualServersState                    *prometheus.GaugeVec
	virtualServersWaitingRequests          *prometheus.GaugeVec
	virtualServersHealth                   *prometheus.GaugeVec
	virtualServersInactiveServices         *prometheus.GaugeVec
	virtualServersActiveServices           *prometheus.GaugeVec
	virtualServersTotalHits                *prometheus.GaugeVec
	virtualServersTotalRequests            *prometheus.GaugeVec
	virtualServersTotalResponses           *prometheus.GaugeVec
	virtualServersTotalRequestBytes        *prometheus.GaugeVec
	virtualServersTotalResponseBytes       *prometheus.GaugeVec
	virtualServersCurrentClientConnections *prometheus.GaugeVec
	virtualServersCurrentServerConnections *prometheus.GaugeVec

	// Service metrics
	servicesThroughput                   *prometheus.GaugeVec
	servicesAvgTTFB                      *prometheus.GaugeVec
	servicesState                        *prometheus.GaugeVec
	servicesTotalRequests                *prometheus.GaugeVec
	servicesTotalResponses               *prometheus.GaugeVec
	servicesTotalRequestBytes            *prometheus.GaugeVec
	servicesTotalResponseBytes           *prometheus.GaugeVec
	servicesCurrentClientConns           *prometheus.GaugeVec
	servicesSurgeCount                   *prometheus.GaugeVec
	servicesCurrentServerConns           *prometheus.GaugeVec
	servicesServerEstablishedConnections *prometheus.GaugeVec
	servicesCurrentReusePool             *prometheus.GaugeVec
	servicesMaxClients                   *prometheus.GaugeVec
	servicesCurrentLoad                  *prometheus.GaugeVec
	servicesVirtualServerServiceHits     *prometheus.GaugeVec
	servicesActiveTransactions           *prometheus.GaugeVec

	// Service Group metrics
	serviceGroupsState                        *prometheus.GaugeVec
	serviceGroupsAvgTTFB                      *prometheus.GaugeVec
	serviceGroupsTotalRequests                *prometheus.GaugeVec
	serviceGroupsTotalResponses               *prometheus.GaugeVec
	serviceGroupsTotalRequestBytes            *prometheus.GaugeVec
	serviceGroupsTotalResponseBytes           *prometheus.GaugeVec
	serviceGroupsCurrentClientConnections     *prometheus.GaugeVec
	serviceGroupsSurgeCount                   *prometheus.GaugeVec
	serviceGroupsCurrentServerConnections     *prometheus.GaugeVec
	serviceGroupsServerEstablishedConnections *prometheus.GaugeVec
	serviceGroupsCurrentReusePool             *prometheus.GaugeVec
	serviceGroupsMaxClients                   *prometheus.GaugeVec

	// GSLB Service metrics
	gslbServicesState                    *prometheus.GaugeVec
	gslbServicesTotalRequests            *prometheus.GaugeVec
	gslbServicesTotalResponses           *prometheus.GaugeVec
	gslbServicesTotalRequestBytes        *prometheus.GaugeVec
	gslbServicesTotalResponseBytes       *prometheus.GaugeVec
	gslbServicesCurrentClientConns       *prometheus.GaugeVec
	gslbServicesCurrentServerConns       *prometheus.GaugeVec
	gslbServicesCurrentLoad              *prometheus.GaugeVec
	gslbServicesVirtualServerServiceHits *prometheus.GaugeVec
	gslbServicesEstablishedConnections   *prometheus.GaugeVec

	// GSLB Virtual Server metrics
	gslbVirtualServersState                    *prometheus.GaugeVec
	gslbVirtualServersHealth                   *prometheus.GaugeVec
	gslbVirtualServersInactiveServices         *prometheus.GaugeVec
	gslbVirtualServersActiveServices           *prometheus.GaugeVec
	gslbVirtualServersTotalHits                *prometheus.GaugeVec
	gslbVirtualServersTotalRequests            *prometheus.GaugeVec
	gslbVirtualServersTotalResponses           *prometheus.GaugeVec
	gslbVirtualServersTotalRequestBytes        *prometheus.GaugeVec
	gslbVirtualServersTotalResponseBytes       *prometheus.GaugeVec
	gslbVirtualServersCurrentClientConnections *prometheus.GaugeVec
	gslbVirtualServersCurrentServerConnections *prometheus.GaugeVec

	// CS Virtual Server metrics
	csVirtualServersState                              *prometheus.GaugeVec
	csVirtualServersTotalHits                          *prometheus.GaugeVec
	csVirtualServersTotalRequests                      *prometheus.GaugeVec
	csVirtualServersTotalResponses                     *prometheus.GaugeVec
	csVirtualServersTotalRequestBytes                  *prometheus.GaugeVec
	csVirtualServersTotalResponseBytes                 *prometheus.GaugeVec
	csVirtualServersCurrentClientConnections           *prometheus.GaugeVec
	csVirtualServersCurrentServerConnections           *prometheus.GaugeVec
	csVirtualServersEstablishedConnections             *prometheus.GaugeVec
	csVirtualServersTotalPacketsReceived               *prometheus.GaugeVec
	csVirtualServersTotalPacketsSent                   *prometheus.GaugeVec
	csVirtualServersTotalSpillovers                    *prometheus.GaugeVec
	csVirtualServersDeferredRequests                   *prometheus.GaugeVec
	csVirtualServersNumberInvalidRequestResponse       *prometheus.GaugeVec
	csVirtualServersNumberInvalidRequestResponseDropped *prometheus.GaugeVec
	csVirtualServersTotalVServerDownBackupHits         *prometheus.GaugeVec
	csVirtualServersCurrentMultipathSessions           *prometheus.GaugeVec
	csVirtualServersCurrentMultipathSubflows           *prometheus.GaugeVec

	// VPN Virtual Server metrics
	vpnVirtualServersTotalRequests      *prometheus.GaugeVec
	vpnVirtualServersTotalResponses     *prometheus.GaugeVec
	vpnVirtualServersTotalRequestBytes  *prometheus.GaugeVec
	vpnVirtualServersTotalResponseBytes *prometheus.GaugeVec
	vpnVirtualServersState              *prometheus.GaugeVec

	// AAA metrics
	aaaAuthSuccess         *prometheus.GaugeVec
	aaaAuthFail            *prometheus.GaugeVec
	aaaAuthOnlyHTTPSuccess *prometheus.GaugeVec
	aaaAuthOnlyHTTPFail    *prometheus.GaugeVec
	aaaCurIcaSessions      *prometheus.GaugeVec
	aaaCurIcaOnlyConn      *prometheus.GaugeVec

	// Topology metrics
	topologyNode *prometheus.GaugeVec
	topologyEdge *prometheus.GaugeVec
}

// NewExporter initialises the exporter with the given configuration
func NewExporter(cfg *config.Config, username, password string, ignoreCert bool, logger *slog.Logger) (*Exporter, error) {
	labelKeys := cfg.LabelKeys()

	// Build base label names for different metric types
	baseLabels := append([]string{"ns_instance"}, labelKeys...)
	vsLabels := append(baseLabels, "virtual_server")
	svcLabels := append(baseLabels, "service")
	sgLabels := append(baseLabels, "servicegroup", "member", "port")
	ifLabels := append(baseLabels, "interface", "alias")
	vpnVsLabels := append(baseLabels, "vpn_virtual_server")
	topoNodeLabels := append(baseLabels, "id", "title", "node_type", "state")
	topoEdgeLabels := append(baseLabels, "id", "source", "target", "weight", "priority")

	e := &Exporter{
		config:     cfg,
		username:   username,
		password:   password,
		ignoreCert: ignoreCert,
		labelKeys:  labelKeys,
		logger:     logger,

		// System metrics (descriptors)
		modelID:             prometheus.NewDesc("model_id", "NetScaler model - reflects the bandwidth available", baseLabels, nil),
		mgmtCPUUsage:        prometheus.NewDesc("mgmt_cpu_usage", "Current CPU utilisation for management", baseLabels, nil),
		memUsage:            prometheus.NewDesc("mem_usage", "Current memory utilisation", baseLabels, nil),
		pktCPUUsage:         prometheus.NewDesc("pkt_cpu_usage", "Current CPU utilisation for packet engines", baseLabels, nil),
		flashPartitionUsage: prometheus.NewDesc("flash_partition_usage", "Used space in /flash partition", baseLabels, nil),
		varPartitionUsage:   prometheus.NewDesc("var_partition_usage", "Used space in /var partition", baseLabels, nil),
		totRxMB:             prometheus.NewDesc("total_received_mb", "Total Megabytes received", baseLabels, nil),
		totTxMB:             prometheus.NewDesc("total_transmit_mb", "Total Megabytes transmitted", baseLabels, nil),
		httpRequests:        prometheus.NewDesc("http_requests", "Total HTTP requests received", baseLabels, nil),
		httpResponses:       prometheus.NewDesc("http_responses", "Total HTTP responses sent", baseLabels, nil),
		tcpCurrentClientConnections:            prometheus.NewDesc("tcp_current_client_connections", "Current client connections", baseLabels, nil),
		tcpCurrentClientConnectionsEstablished: prometheus.NewDesc("tcp_current_client_connections_established", "Current established client connections", baseLabels, nil),
		tcpCurrentServerConnections:            prometheus.NewDesc("tcp_current_server_connections", "Current server connections", baseLabels, nil),
		tcpCurrentServerConnectionsEstablished: prometheus.NewDesc("tcp_current_server_connections_established", "Current established server connections", baseLabels, nil),

		// Interface metrics
		interfacesRxBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_received_bytes", Help: "Bytes received by interface"}, ifLabels),
		interfacesTxBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_transmitted_bytes", Help: "Bytes transmitted by interface"}, ifLabels),
		interfacesRxPackets:      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_received_packets", Help: "Packets received by interface"}, ifLabels),
		interfacesTxPackets:      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_transmitted_packets", Help: "Packets transmitted by interface"}, ifLabels),
		interfacesJumboPacketsRx: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_jumbo_packets_received", Help: "Jumbo packets received by interface"}, ifLabels),
		interfacesJumboPacketsTx: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_jumbo_packets_transmitted", Help: "Jumbo packets transmitted by interface"}, ifLabels),
		interfacesErrorPacketsRx: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "interfaces_error_packets_received", Help: "Error packets received by interface"}, ifLabels),

		// Virtual Server metrics
		virtualServersState:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_state", Help: "Current state of the server"}, vsLabels),
		virtualServersWaitingRequests:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_waiting_requests", Help: "Number of waiting requests"}, vsLabels),
		virtualServersHealth:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_health", Help: "Percentage of UP services"}, vsLabels),
		virtualServersInactiveServices:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_inactive_services", Help: "Number of inactive services"}, vsLabels),
		virtualServersActiveServices:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_active_services", Help: "Number of active services"}, vsLabels),
		virtualServersTotalHits:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_total_hits", Help: "Total hits"}, vsLabels),
		virtualServersTotalRequests:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_total_requests", Help: "Total requests"}, vsLabels),
		virtualServersTotalResponses:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_total_responses", Help: "Total responses"}, vsLabels),
		virtualServersTotalRequestBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_total_request_bytes", Help: "Total request bytes"}, vsLabels),
		virtualServersTotalResponseBytes:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_total_response_bytes", Help: "Total response bytes"}, vsLabels),
		virtualServersCurrentClientConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_current_client_connections", Help: "Current client connections"}, vsLabels),
		virtualServersCurrentServerConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "virtual_servers_current_server_connections", Help: "Current server connections"}, vsLabels),

		// Service metrics
		servicesThroughput:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_throughput", Help: "Throughput in Mbps"}, svcLabels),
		servicesAvgTTFB:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_average_time_to_first_byte", Help: "Average TTFB"}, svcLabels),
		servicesState:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_state", Help: "Current state"}, svcLabels),
		servicesTotalRequests:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_total_requests", Help: "Total requests"}, svcLabels),
		servicesTotalResponses:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_total_responses", Help: "Total responses"}, svcLabels),
		servicesTotalRequestBytes:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_total_request_bytes", Help: "Total request bytes"}, svcLabels),
		servicesTotalResponseBytes:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_total_response_bytes", Help: "Total response bytes"}, svcLabels),
		servicesCurrentClientConns:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_current_client_connections", Help: "Current client connections"}, svcLabels),
		servicesSurgeCount:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_surge_count", Help: "Requests in surge queue"}, svcLabels),
		servicesCurrentServerConns:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_current_server_connections", Help: "Current server connections"}, svcLabels),
		servicesServerEstablishedConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_server_established_connections", Help: "Established server connections"}, svcLabels),
		servicesCurrentReusePool:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_current_reuse_pool", Help: "Requests in reuse pool"}, svcLabels),
		servicesMaxClients:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_max_clients", Help: "Max open connections"}, svcLabels),
		servicesCurrentLoad:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_current_load", Help: "Current load"}, svcLabels),
		servicesVirtualServerServiceHits:     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_virtual_server_service_hits", Help: "Service hits"}, svcLabels),
		servicesActiveTransactions:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "service_active_transactions", Help: "Active transactions"}, svcLabels),

		// Service Group metrics
		serviceGroupsState:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_state", Help: "Current state"}, sgLabels),
		serviceGroupsAvgTTFB:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_average_time_to_first_byte", Help: "Average TTFB"}, sgLabels),
		serviceGroupsTotalRequests:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_total_requests", Help: "Total requests"}, sgLabels),
		serviceGroupsTotalResponses:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_total_responses", Help: "Total responses"}, sgLabels),
		serviceGroupsTotalRequestBytes:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_total_request_bytes", Help: "Total request bytes"}, sgLabels),
		serviceGroupsTotalResponseBytes:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_total_response_bytes", Help: "Total response bytes"}, sgLabels),
		serviceGroupsCurrentClientConnections:     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_current_client_connections", Help: "Current client connections"}, sgLabels),
		serviceGroupsSurgeCount:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_surge_count", Help: "Requests in surge queue"}, sgLabels),
		serviceGroupsCurrentServerConnections:     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_current_server_connections", Help: "Current server connections"}, sgLabels),
		serviceGroupsServerEstablishedConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_server_established_connections", Help: "Established server connections"}, sgLabels),
		serviceGroupsCurrentReusePool:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_current_reuse_pool", Help: "Requests in reuse pool"}, sgLabels),
		serviceGroupsMaxClients:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "servicegroup_max_clients", Help: "Max open connections"}, sgLabels),

		// GSLB Service metrics
		gslbServicesState:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_state", Help: "Current state"}, svcLabels),
		gslbServicesTotalRequests:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_total_requests", Help: "Total requests"}, svcLabels),
		gslbServicesTotalResponses:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_total_responses", Help: "Total responses"}, svcLabels),
		gslbServicesTotalRequestBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_total_request_bytes", Help: "Total request bytes"}, svcLabels),
		gslbServicesTotalResponseBytes:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_total_response_bytes", Help: "Total response bytes"}, svcLabels),
		gslbServicesCurrentClientConns:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_current_client_connections", Help: "Current client connections"}, svcLabels),
		gslbServicesCurrentServerConns:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_current_server_connections", Help: "Current server connections"}, svcLabels),
		gslbServicesCurrentLoad:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_current_load", Help: "Current load"}, svcLabels),
		gslbServicesVirtualServerServiceHits: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_virtual_server_service_hits", Help: "Service hits"}, svcLabels),
		gslbServicesEstablishedConnections:   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_service_established_connections", Help: "Established connections"}, svcLabels),

		// GSLB Virtual Server metrics
		gslbVirtualServersState:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_state", Help: "Current state"}, vsLabels),
		gslbVirtualServersHealth:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_health", Help: "Percentage of UP services"}, vsLabels),
		gslbVirtualServersInactiveServices:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_inactive_services", Help: "Inactive services"}, vsLabels),
		gslbVirtualServersActiveServices:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_active_services", Help: "Active services"}, vsLabels),
		gslbVirtualServersTotalHits:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_total_hits", Help: "Total hits"}, vsLabels),
		gslbVirtualServersTotalRequests:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_total_requests", Help: "Total requests"}, vsLabels),
		gslbVirtualServersTotalResponses:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_total_responses", Help: "Total responses"}, vsLabels),
		gslbVirtualServersTotalRequestBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_total_request_bytes", Help: "Total request bytes"}, vsLabels),
		gslbVirtualServersTotalResponseBytes:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_total_response_bytes", Help: "Total response bytes"}, vsLabels),
		gslbVirtualServersCurrentClientConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_current_client_connections", Help: "Current client connections"}, vsLabels),
		gslbVirtualServersCurrentServerConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "gslb_virtual_servers_current_server_connections", Help: "Current server connections"}, vsLabels),

		// CS Virtual Server metrics
		csVirtualServersState:                               prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_state", Help: "Current state"}, vsLabels),
		csVirtualServersTotalHits:                           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_hits", Help: "Total hits"}, vsLabels),
		csVirtualServersTotalRequests:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_requests", Help: "Total requests"}, vsLabels),
		csVirtualServersTotalResponses:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_responses", Help: "Total responses"}, vsLabels),
		csVirtualServersTotalRequestBytes:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_request_bytes", Help: "Total request bytes"}, vsLabels),
		csVirtualServersTotalResponseBytes:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_response_bytes", Help: "Total response bytes"}, vsLabels),
		csVirtualServersCurrentClientConnections:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_current_client_connections", Help: "Current client connections"}, vsLabels),
		csVirtualServersCurrentServerConnections:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_current_server_connections", Help: "Current server connections"}, vsLabels),
		csVirtualServersEstablishedConnections:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_established_connections", Help: "Established connections"}, vsLabels),
		csVirtualServersTotalPacketsReceived:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_packets_received", Help: "Total packets received"}, vsLabels),
		csVirtualServersTotalPacketsSent:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_packets_sent", Help: "Total packets sent"}, vsLabels),
		csVirtualServersTotalSpillovers:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_spillovers", Help: "Total spillovers"}, vsLabels),
		csVirtualServersDeferredRequests:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_deferred_requests", Help: "Deferred requests"}, vsLabels),
		csVirtualServersNumberInvalidRequestResponse:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_number_invalid_request_response", Help: "Invalid request/responses"}, vsLabels),
		csVirtualServersNumberInvalidRequestResponseDropped: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_number_invalid_request_response_dropped", Help: "Invalid request/responses dropped"}, vsLabels),
		csVirtualServersTotalVServerDownBackupHits:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_total_vserver_down_backup_hits", Help: "Backup hits when vserver down"}, vsLabels),
		csVirtualServersCurrentMultipathSessions:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_current_multipath_sessions", Help: "Current multipath TCP sessions"}, vsLabels),
		csVirtualServersCurrentMultipathSubflows:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "cs_virtual_servers_current_multipath_subflows", Help: "Current multipath TCP subflows"}, vsLabels),

		// VPN Virtual Server metrics
		vpnVirtualServersTotalRequests:      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "vpn_virtual_servers_total_requests", Help: "Total requests"}, vpnVsLabels),
		vpnVirtualServersTotalResponses:     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "vpn_virtual_servers_total_responses", Help: "Total responses"}, vpnVsLabels),
		vpnVirtualServersTotalRequestBytes:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "vpn_virtual_servers_total_request_bytes", Help: "Total request bytes"}, vpnVsLabels),
		vpnVirtualServersTotalResponseBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "vpn_virtual_servers_total_response_bytes", Help: "Total response bytes"}, vpnVsLabels),
		vpnVirtualServersState:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "vpn_virtual_servers_state", Help: "Current state"}, vpnVsLabels),

		// AAA metrics
		aaaAuthSuccess:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_auth_success", Help: "Authentication successes"}, baseLabels),
		aaaAuthFail:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_auth_fail", Help: "Authentication failures"}, baseLabels),
		aaaAuthOnlyHTTPSuccess: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_auth_only_http_success", Help: "HTTP auth successes"}, baseLabels),
		aaaAuthOnlyHTTPFail:    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_auth_only_http_fail", Help: "HTTP auth failures"}, baseLabels),
		aaaCurIcaSessions:      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_current_ica_sessions", Help: "Current ICA sessions"}, baseLabels),
		aaaCurIcaOnlyConn:      prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "aaa_current_ica_only_connections", Help: "Current ICA connections"}, baseLabels),

		// Topology metrics
		topologyNode: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_topology_node", Help: "Node for topology visualization"}, topoNodeLabels),
		topologyEdge: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_topology_edge", Help: "Edge between frontend and backend"}, topoEdgeLabels),
	}

	return e, nil
}

// buildLabelValues creates the label values slice for a target
func (e *Exporter) buildLabelValues(target config.Target, extraLabels ...string) []string {
	labels := target.MergedLabels(e.config.Labels)
	values := make([]string, 0, 1+len(e.labelKeys)+len(extraLabels))
	values = append(values, target.URL)
	for _, k := range e.labelKeys {
		values = append(values, labels[k])
	}
	values = append(values, extraLabels...)
	return values
}

// Describe implements Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.modelID
	ch <- e.mgmtCPUUsage
	ch <- e.memUsage
	ch <- e.pktCPUUsage
	ch <- e.flashPartitionUsage
	ch <- e.varPartitionUsage
	ch <- e.totRxMB
	ch <- e.totTxMB
	ch <- e.httpResponses
	ch <- e.httpRequests
	ch <- e.tcpCurrentClientConnections
	ch <- e.tcpCurrentClientConnectionsEstablished
	ch <- e.tcpCurrentServerConnections
	ch <- e.tcpCurrentServerConnectionsEstablished

	e.interfacesRxBytes.Describe(ch)
	e.interfacesTxBytes.Describe(ch)
	e.interfacesRxPackets.Describe(ch)
	e.interfacesTxPackets.Describe(ch)
	e.interfacesJumboPacketsRx.Describe(ch)
	e.interfacesJumboPacketsTx.Describe(ch)
	e.interfacesErrorPacketsRx.Describe(ch)

	e.virtualServersState.Describe(ch)
	e.virtualServersWaitingRequests.Describe(ch)
	e.virtualServersHealth.Describe(ch)
	e.virtualServersInactiveServices.Describe(ch)
	e.virtualServersActiveServices.Describe(ch)
	e.virtualServersTotalHits.Describe(ch)
	e.virtualServersTotalRequests.Describe(ch)
	e.virtualServersTotalResponses.Describe(ch)
	e.virtualServersTotalRequestBytes.Describe(ch)
	e.virtualServersTotalResponseBytes.Describe(ch)
	e.virtualServersCurrentClientConnections.Describe(ch)
	e.virtualServersCurrentServerConnections.Describe(ch)

	e.servicesThroughput.Describe(ch)
	e.servicesAvgTTFB.Describe(ch)
	e.servicesState.Describe(ch)
	e.servicesTotalRequests.Describe(ch)
	e.servicesTotalResponses.Describe(ch)
	e.servicesTotalRequestBytes.Describe(ch)
	e.servicesTotalResponseBytes.Describe(ch)
	e.servicesCurrentClientConns.Describe(ch)
	e.servicesSurgeCount.Describe(ch)
	e.servicesCurrentServerConns.Describe(ch)
	e.servicesServerEstablishedConnections.Describe(ch)
	e.servicesCurrentReusePool.Describe(ch)
	e.servicesMaxClients.Describe(ch)
	e.servicesCurrentLoad.Describe(ch)
	e.servicesVirtualServerServiceHits.Describe(ch)
	e.servicesActiveTransactions.Describe(ch)

	e.serviceGroupsState.Describe(ch)
	e.serviceGroupsAvgTTFB.Describe(ch)
	e.serviceGroupsTotalRequests.Describe(ch)
	e.serviceGroupsTotalResponses.Describe(ch)
	e.serviceGroupsTotalRequestBytes.Describe(ch)
	e.serviceGroupsTotalResponseBytes.Describe(ch)
	e.serviceGroupsCurrentClientConnections.Describe(ch)
	e.serviceGroupsSurgeCount.Describe(ch)
	e.serviceGroupsCurrentServerConnections.Describe(ch)
	e.serviceGroupsServerEstablishedConnections.Describe(ch)
	e.serviceGroupsCurrentReusePool.Describe(ch)
	e.serviceGroupsMaxClients.Describe(ch)

	e.gslbServicesState.Describe(ch)
	e.gslbServicesTotalRequests.Describe(ch)
	e.gslbServicesTotalResponses.Describe(ch)
	e.gslbServicesTotalRequestBytes.Describe(ch)
	e.gslbServicesTotalResponseBytes.Describe(ch)
	e.gslbServicesCurrentClientConns.Describe(ch)
	e.gslbServicesCurrentServerConns.Describe(ch)
	e.gslbServicesCurrentLoad.Describe(ch)
	e.gslbServicesVirtualServerServiceHits.Describe(ch)
	e.gslbServicesEstablishedConnections.Describe(ch)

	e.gslbVirtualServersState.Describe(ch)
	e.gslbVirtualServersHealth.Describe(ch)
	e.gslbVirtualServersInactiveServices.Describe(ch)
	e.gslbVirtualServersActiveServices.Describe(ch)
	e.gslbVirtualServersTotalHits.Describe(ch)
	e.gslbVirtualServersTotalRequests.Describe(ch)
	e.gslbVirtualServersTotalResponses.Describe(ch)
	e.gslbVirtualServersTotalRequestBytes.Describe(ch)
	e.gslbVirtualServersTotalResponseBytes.Describe(ch)
	e.gslbVirtualServersCurrentClientConnections.Describe(ch)
	e.gslbVirtualServersCurrentServerConnections.Describe(ch)

	e.csVirtualServersState.Describe(ch)
	e.csVirtualServersTotalHits.Describe(ch)
	e.csVirtualServersTotalRequests.Describe(ch)
	e.csVirtualServersTotalResponses.Describe(ch)
	e.csVirtualServersTotalRequestBytes.Describe(ch)
	e.csVirtualServersTotalResponseBytes.Describe(ch)
	e.csVirtualServersCurrentClientConnections.Describe(ch)
	e.csVirtualServersCurrentServerConnections.Describe(ch)
	e.csVirtualServersEstablishedConnections.Describe(ch)
	e.csVirtualServersTotalPacketsReceived.Describe(ch)
	e.csVirtualServersTotalPacketsSent.Describe(ch)
	e.csVirtualServersTotalSpillovers.Describe(ch)
	e.csVirtualServersDeferredRequests.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponse.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponseDropped.Describe(ch)
	e.csVirtualServersTotalVServerDownBackupHits.Describe(ch)
	e.csVirtualServersCurrentMultipathSessions.Describe(ch)
	e.csVirtualServersCurrentMultipathSubflows.Describe(ch)

	e.vpnVirtualServersTotalRequests.Describe(ch)
	e.vpnVirtualServersTotalResponses.Describe(ch)
	e.vpnVirtualServersTotalRequestBytes.Describe(ch)
	e.vpnVirtualServersTotalResponseBytes.Describe(ch)
	e.vpnVirtualServersState.Describe(ch)

	e.aaaAuthSuccess.Describe(ch)
	e.aaaAuthFail.Describe(ch)
	e.aaaAuthOnlyHTTPSuccess.Describe(ch)
	e.aaaAuthOnlyHTTPFail.Describe(ch)
	e.aaaCurIcaSessions.Describe(ch)
	e.aaaCurIcaOnlyConn.Describe(ch)

	e.topologyNode.Describe(ch)
	e.topologyEdge.Describe(ch)
}
