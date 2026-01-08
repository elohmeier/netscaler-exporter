package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/elohmeier/netscaler-exporter/config"
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	config      *config.Config
	username    string
	password    string
	ignoreCert  bool
	parallelism int
	labelKeys   []string
	logger      *slog.Logger

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

	// Probe success metric
	probeSuccess *prometheus.GaugeVec

	// Protocol HTTP metrics
	httpTotalRequests              *prometheus.Desc
	httpTotalResponses             *prometheus.Desc
	httpTotalPosts                 *prometheus.Desc
	httpTotalGets                  *prometheus.Desc
	httpTotalOthers                *prometheus.Desc
	httpTotalRxRequestBytes        *prometheus.Desc
	httpTotalRxResponseBytes       *prometheus.Desc
	httpTotalTxRequestBytes        *prometheus.Desc
	httpTotal10Requests            *prometheus.Desc
	httpTotal11Requests            *prometheus.Desc
	httpTotal10Responses           *prometheus.Desc
	httpTotal11Responses           *prometheus.Desc
	httpTotalChunkedRequests       *prometheus.Desc
	httpTotalChunkedResponses      *prometheus.Desc
	httpTotalSPDYStreams           *prometheus.Desc
	httpTotalSPDYv2Streams         *prometheus.Desc
	httpTotalSPDYv3Streams         *prometheus.Desc
	httpErrNoReuseMultipart        *prometheus.Desc
	httpErrIncompleteHeaders       *prometheus.Desc
	httpErrIncompleteRequests      *prometheus.Desc
	httpErrIncompleteResponses     *prometheus.Desc
	httpErrServerBusy              *prometheus.Desc
	httpErrLargeContent            *prometheus.Desc
	httpErrLargeChunk              *prometheus.Desc
	httpErrLargeCtlen              *prometheus.Desc
	httpRequestsRate               *prometheus.Desc
	httpResponsesRate              *prometheus.Desc
	httpPostsRate                  *prometheus.Desc
	httpGetsRate                   *prometheus.Desc
	httpOthersRate                 *prometheus.Desc
	httpRxRequestBytesRate         *prometheus.Desc
	httpRxResponseBytesRate        *prometheus.Desc
	httpTxRequestBytesRate         *prometheus.Desc
	httpRequest10Rate              *prometheus.Desc
	httpRequest11Rate              *prometheus.Desc
	httpResponse10Rate             *prometheus.Desc
	httpResponse11Rate             *prometheus.Desc
	httpChunkedRequestsRate        *prometheus.Desc
	httpChunkedResponsesRate       *prometheus.Desc
	httpSPDYStreamsRate            *prometheus.Desc
	httpSPDYv2StreamsRate          *prometheus.Desc
	httpSPDYv3StreamsRate          *prometheus.Desc
	httpErrNoReuseMultipartRate    *prometheus.Desc
	httpErrIncompleteRequestsRate  *prometheus.Desc
	httpErrIncompleteResponsesRate *prometheus.Desc
	httpErrServerBusyRate          *prometheus.Desc

	// Protocol TCP metrics
	tcpTotalRxPackets           *prometheus.Desc
	tcpTotalRxBytes             *prometheus.Desc
	tcpTotalTxBytes             *prometheus.Desc
	tcpTotalTxPackets           *prometheus.Desc
	tcpTotalClientConnOpened    *prometheus.Desc
	tcpTotalServerConnOpened    *prometheus.Desc
	tcpTotalSyn                 *prometheus.Desc
	tcpTotalSynProbe            *prometheus.Desc
	tcpTotalServerFin           *prometheus.Desc
	tcpTotalClientFin           *prometheus.Desc
	tcpActiveServerConn         *prometheus.Desc
	tcpCurClientConnEstablished *prometheus.Desc
	tcpCurServerConnEstablished *prometheus.Desc
	tcpRxPacketsRate            *prometheus.Desc
	tcpRxBytesRate              *prometheus.Desc
	tcpTxPacketsRate            *prometheus.Desc
	tcpTxBytesRate              *prometheus.Desc
	tcpClientConnOpenedRate     *prometheus.Desc
	tcpErrBadChecksum           *prometheus.Desc
	tcpErrBadChecksumRate       *prometheus.Desc
	tcpErrAnyPortFail           *prometheus.Desc
	tcpErrIPPortFail            *prometheus.Desc
	tcpErrBadStateConn          *prometheus.Desc
	tcpErrRstThreshold          *prometheus.Desc
	tcpSynRate                  *prometheus.Desc
	tcpSynProbeRate             *prometheus.Desc

	// Protocol IP metrics
	ipTotalRxPackets          *prometheus.Desc
	ipTotalRxBytes            *prometheus.Desc
	ipTotalTxPackets          *prometheus.Desc
	ipTotalTxBytes            *prometheus.Desc
	ipTotalRxMbits            *prometheus.Desc
	ipTotalTxMbits            *prometheus.Desc
	ipTotalRoutedPackets      *prometheus.Desc
	ipTotalRoutedMbits        *prometheus.Desc
	ipTotalFragments          *prometheus.Desc
	ipTotalSuccReassembly     *prometheus.Desc
	ipTotalAddrLookup         *prometheus.Desc
	ipTotalAddrLookupFail     *prometheus.Desc
	ipTotalUDPFragmentsFwd    *prometheus.Desc
	ipTotalTCPFragmentsFwd    *prometheus.Desc
	ipTotalBadChecksums       *prometheus.Desc
	ipTotalUnsuccReassembly   *prometheus.Desc
	ipTotalTooBig             *prometheus.Desc
	ipTotalDupFragments       *prometheus.Desc
	ipTotalOutOfOrderFrag     *prometheus.Desc
	ipTotalVIPDown            *prometheus.Desc
	ipTotalTTLExpired         *prometheus.Desc
	ipTotalMaxClients         *prometheus.Desc
	ipTotalUnknownSvcs        *prometheus.Desc
	ipTotalInvalidHeaderSz    *prometheus.Desc
	ipTotalInvalidPacketSize  *prometheus.Desc
	ipTotalTruncatedPackets   *prometheus.Desc
	ipNonIPTotalTruncatedPkts *prometheus.Desc
	ipTotalBadMacAddrs        *prometheus.Desc
	ipRxPacketsRate           *prometheus.Desc
	ipRxBytesRate             *prometheus.Desc
	ipTxPacketsRate           *prometheus.Desc
	ipTxBytesRate             *prometheus.Desc
	ipRxMbitsRate             *prometheus.Desc
	ipTxMbitsRate             *prometheus.Desc
	ipRoutedPacketsRate       *prometheus.Desc
	ipRoutedMbitsRate         *prometheus.Desc

	// SSL global metrics
	sslTotalTLSv11Sessions  *prometheus.Desc
	sslTotalSSLv2Sessions   *prometheus.Desc
	sslTotalSessions        *prometheus.Desc
	sslTotalSSLv2Handshakes *prometheus.Desc
	sslTotalEnc             *prometheus.Desc
	sslCryptoUtilization    *prometheus.Desc
	sslTotalNewSessions     *prometheus.Desc
	sslSessionsRate         *prometheus.Desc
	sslDecRate              *prometheus.Desc
	sslEncRate              *prometheus.Desc
	sslSSLv2HandshakesRate  *prometheus.Desc
	sslNewSessionsRate      *prometheus.Desc

	// SSL certificate metrics
	sslCertDaysToExpire *prometheus.GaugeVec

	// SSL VServer metrics
	sslVServerTotalDecBytes          *prometheus.GaugeVec
	sslVServerTotalEncBytes          *prometheus.GaugeVec
	sslVServerTotalHWDecBytes        *prometheus.GaugeVec
	sslVServerTotalHWEncBytes        *prometheus.GaugeVec
	sslVServerTotalSessionNew        *prometheus.GaugeVec
	sslVServerTotalSessionHits       *prometheus.GaugeVec
	sslVServerTotalClientAuthSuccess *prometheus.GaugeVec
	sslVServerTotalClientAuthFailure *prometheus.GaugeVec
	sslVServerHealth                 *prometheus.GaugeVec
	sslVServerActiveServices         *prometheus.GaugeVec
	sslVServerClientAuthSuccessRate  *prometheus.GaugeVec
	sslVServerClientAuthFailureRate  *prometheus.GaugeVec
	sslVServerEncBytesRate           *prometheus.GaugeVec
	sslVServerDecBytesRate           *prometheus.GaugeVec
	sslVServerHWEncBytesRate         *prometheus.GaugeVec
	sslVServerHWDecBytesRate         *prometheus.GaugeVec
	sslVServerSessionNewRate         *prometheus.GaugeVec
	sslVServerSessionHitsRate        *prometheus.GaugeVec

	// System CPU per-core metrics
	cpuCoreUsage *prometheus.GaugeVec

	// Bandwidth capacity metrics
	capacityMaxBandwidth    *prometheus.Desc
	capacityMinBandwidth    *prometheus.Desc
	capacityActualBandwidth *prometheus.Desc
	capacityBandwidth       *prometheus.Desc
}

// NewExporter initialises the exporter with the given configuration
func NewExporter(cfg *config.Config, username, password string, ignoreCert bool, parallelism int, logger *slog.Logger) (*Exporter, error) {
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
	sslCertLabels := append(baseLabels, "certkey")
	sslVsLabels := append(baseLabels, "vserver", "type", "ip")
	cpuCoreLabels := append(baseLabels, "core_id")

	e := &Exporter{
		config:      cfg,
		username:    username,
		password:    password,
		ignoreCert:  ignoreCert,
		parallelism: parallelism,
		labelKeys:   labelKeys,
		logger:      logger,

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

		// Probe success metric
		probeSuccess: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_probe_success", Help: "Probe success (1 = success, 0 = failure)"}, baseLabels),

		// Protocol HTTP metrics
		httpTotalRequests:              prometheus.NewDesc("netscaler_http_requests_total", "Total HTTP requests", baseLabels, nil),
		httpTotalResponses:             prometheus.NewDesc("netscaler_http_responses_total", "Total HTTP responses", baseLabels, nil),
		httpTotalPosts:                 prometheus.NewDesc("netscaler_http_posts_total", "Total HTTP POST requests", baseLabels, nil),
		httpTotalGets:                  prometheus.NewDesc("netscaler_http_gets_total", "Total HTTP GET requests", baseLabels, nil),
		httpTotalOthers:                prometheus.NewDesc("netscaler_http_others_total", "Total other HTTP requests", baseLabels, nil),
		httpTotalRxRequestBytes:        prometheus.NewDesc("netscaler_http_rx_request_bytes_total", "Total HTTP request bytes received", baseLabels, nil),
		httpTotalRxResponseBytes:       prometheus.NewDesc("netscaler_http_rx_response_bytes_total", "Total HTTP response bytes received", baseLabels, nil),
		httpTotalTxRequestBytes:        prometheus.NewDesc("netscaler_http_tx_request_bytes_total", "Total HTTP request bytes transmitted", baseLabels, nil),
		httpTotal10Requests:            prometheus.NewDesc("netscaler_http_10_requests_total", "Total HTTP/1.0 requests", baseLabels, nil),
		httpTotal11Requests:            prometheus.NewDesc("netscaler_http_11_requests_total", "Total HTTP/1.1 requests", baseLabels, nil),
		httpTotal10Responses:           prometheus.NewDesc("netscaler_http_10_responses_total", "Total HTTP/1.0 responses", baseLabels, nil),
		httpTotal11Responses:           prometheus.NewDesc("netscaler_http_11_responses_total", "Total HTTP/1.1 responses", baseLabels, nil),
		httpTotalChunkedRequests:       prometheus.NewDesc("netscaler_http_chunked_requests_total", "Total chunked HTTP requests", baseLabels, nil),
		httpTotalChunkedResponses:      prometheus.NewDesc("netscaler_http_chunked_responses_total", "Total chunked HTTP responses", baseLabels, nil),
		httpTotalSPDYStreams:           prometheus.NewDesc("netscaler_http_spdy_streams_total", "Total SPDY streams", baseLabels, nil),
		httpTotalSPDYv2Streams:         prometheus.NewDesc("netscaler_http_spdy_v2_streams_total", "Total SPDY v2 streams", baseLabels, nil),
		httpTotalSPDYv3Streams:         prometheus.NewDesc("netscaler_http_spdy_v3_streams_total", "Total SPDY v3 streams", baseLabels, nil),
		httpErrNoReuseMultipart:        prometheus.NewDesc("netscaler_http_err_noreuse_multipart_total", "No-reuse multipart errors", baseLabels, nil),
		httpErrIncompleteHeaders:       prometheus.NewDesc("netscaler_http_err_incomplete_headers_total", "Incomplete header errors", baseLabels, nil),
		httpErrIncompleteRequests:      prometheus.NewDesc("netscaler_http_err_incomplete_requests_total", "Incomplete request errors", baseLabels, nil),
		httpErrIncompleteResponses:     prometheus.NewDesc("netscaler_http_err_incomplete_responses_total", "Incomplete response errors", baseLabels, nil),
		httpErrServerBusy:              prometheus.NewDesc("netscaler_http_err_server_busy_total", "Server busy errors", baseLabels, nil),
		httpErrLargeContent:            prometheus.NewDesc("netscaler_http_err_large_content_total", "Large content errors", baseLabels, nil),
		httpErrLargeChunk:              prometheus.NewDesc("netscaler_http_err_large_chunk_total", "Large chunk errors", baseLabels, nil),
		httpErrLargeCtlen:              prometheus.NewDesc("netscaler_http_err_large_ctlen_total", "Large content-length errors", baseLabels, nil),
		httpRequestsRate:               prometheus.NewDesc("netscaler_http_requests_rate", "HTTP requests rate", baseLabels, nil),
		httpResponsesRate:              prometheus.NewDesc("netscaler_http_responses_rate", "HTTP responses rate", baseLabels, nil),
		httpPostsRate:                  prometheus.NewDesc("netscaler_http_posts_rate", "HTTP POST rate", baseLabels, nil),
		httpGetsRate:                   prometheus.NewDesc("netscaler_http_gets_rate", "HTTP GET rate", baseLabels, nil),
		httpOthersRate:                 prometheus.NewDesc("netscaler_http_others_rate", "Other HTTP requests rate", baseLabels, nil),
		httpRxRequestBytesRate:         prometheus.NewDesc("netscaler_http_rx_request_bytes_rate", "HTTP request bytes received rate", baseLabels, nil),
		httpRxResponseBytesRate:        prometheus.NewDesc("netscaler_http_rx_response_bytes_rate", "HTTP response bytes received rate", baseLabels, nil),
		httpTxRequestBytesRate:         prometheus.NewDesc("netscaler_http_tx_request_bytes_rate", "HTTP request bytes transmitted rate", baseLabels, nil),
		httpRequest10Rate:              prometheus.NewDesc("netscaler_http_10_requests_rate", "HTTP/1.0 requests rate", baseLabels, nil),
		httpRequest11Rate:              prometheus.NewDesc("netscaler_http_11_requests_rate", "HTTP/1.1 requests rate", baseLabels, nil),
		httpResponse10Rate:             prometheus.NewDesc("netscaler_http_10_responses_rate", "HTTP/1.0 responses rate", baseLabels, nil),
		httpResponse11Rate:             prometheus.NewDesc("netscaler_http_11_responses_rate", "HTTP/1.1 responses rate", baseLabels, nil),
		httpChunkedRequestsRate:        prometheus.NewDesc("netscaler_http_chunked_requests_rate", "Chunked requests rate", baseLabels, nil),
		httpChunkedResponsesRate:       prometheus.NewDesc("netscaler_http_chunked_responses_rate", "Chunked responses rate", baseLabels, nil),
		httpSPDYStreamsRate:            prometheus.NewDesc("netscaler_http_spdy_streams_rate", "SPDY streams rate", baseLabels, nil),
		httpSPDYv2StreamsRate:          prometheus.NewDesc("netscaler_http_spdy_v2_streams_rate", "SPDY v2 streams rate", baseLabels, nil),
		httpSPDYv3StreamsRate:          prometheus.NewDesc("netscaler_http_spdy_v3_streams_rate", "SPDY v3 streams rate", baseLabels, nil),
		httpErrNoReuseMultipartRate:    prometheus.NewDesc("netscaler_http_err_noreuse_multipart_rate", "No-reuse multipart errors rate", baseLabels, nil),
		httpErrIncompleteRequestsRate:  prometheus.NewDesc("netscaler_http_err_incomplete_requests_rate", "Incomplete requests rate", baseLabels, nil),
		httpErrIncompleteResponsesRate: prometheus.NewDesc("netscaler_http_err_incomplete_responses_rate", "Incomplete responses rate", baseLabels, nil),
		httpErrServerBusyRate:          prometheus.NewDesc("netscaler_http_err_server_busy_rate", "Server busy errors rate", baseLabels, nil),

		// Protocol TCP metrics
		tcpTotalRxPackets:           prometheus.NewDesc("netscaler_tcp_rx_packets_total", "Total TCP packets received", baseLabels, nil),
		tcpTotalRxBytes:             prometheus.NewDesc("netscaler_tcp_rx_bytes_total", "Total TCP bytes received", baseLabels, nil),
		tcpTotalTxBytes:             prometheus.NewDesc("netscaler_tcp_tx_bytes_total", "Total TCP bytes transmitted", baseLabels, nil),
		tcpTotalTxPackets:           prometheus.NewDesc("netscaler_tcp_tx_packets_total", "Total TCP packets transmitted", baseLabels, nil),
		tcpTotalClientConnOpened:    prometheus.NewDesc("netscaler_tcp_client_connections_opened_total", "Total TCP client connections opened", baseLabels, nil),
		tcpTotalServerConnOpened:    prometheus.NewDesc("netscaler_tcp_server_connections_opened_total", "Total TCP server connections opened", baseLabels, nil),
		tcpTotalSyn:                 prometheus.NewDesc("netscaler_tcp_syn_total", "Total TCP SYN packets", baseLabels, nil),
		tcpTotalSynProbe:            prometheus.NewDesc("netscaler_tcp_syn_probe_total", "Total TCP SYN probe packets", baseLabels, nil),
		tcpTotalServerFin:           prometheus.NewDesc("netscaler_tcp_server_fin_total", "Total TCP server FIN packets", baseLabels, nil),
		tcpTotalClientFin:           prometheus.NewDesc("netscaler_tcp_client_fin_total", "Total TCP client FIN packets", baseLabels, nil),
		tcpActiveServerConn:         prometheus.NewDesc("netscaler_tcp_active_server_connections", "Active TCP server connections", baseLabels, nil),
		tcpCurClientConnEstablished: prometheus.NewDesc("netscaler_tcp_cur_client_connections_established", "Current established client connections", baseLabels, nil),
		tcpCurServerConnEstablished: prometheus.NewDesc("netscaler_tcp_cur_server_connections_established", "Current established server connections", baseLabels, nil),
		tcpRxPacketsRate:            prometheus.NewDesc("netscaler_tcp_rx_packets_rate", "TCP packets received rate", baseLabels, nil),
		tcpRxBytesRate:              prometheus.NewDesc("netscaler_tcp_rx_bytes_rate", "TCP bytes received rate", baseLabels, nil),
		tcpTxPacketsRate:            prometheus.NewDesc("netscaler_tcp_tx_packets_rate", "TCP packets transmitted rate", baseLabels, nil),
		tcpTxBytesRate:              prometheus.NewDesc("netscaler_tcp_tx_bytes_rate", "TCP bytes transmitted rate", baseLabels, nil),
		tcpClientConnOpenedRate:     prometheus.NewDesc("netscaler_tcp_client_connections_opened_rate", "TCP client connections opened rate", baseLabels, nil),
		tcpErrBadChecksum:           prometheus.NewDesc("netscaler_tcp_err_bad_checksum_total", "TCP bad checksum errors", baseLabels, nil),
		tcpErrBadChecksumRate:       prometheus.NewDesc("netscaler_tcp_err_bad_checksum_rate", "TCP bad checksum errors rate", baseLabels, nil),
		tcpErrAnyPortFail:           prometheus.NewDesc("netscaler_tcp_err_any_port_fail", "TCP any port fail errors", baseLabels, nil),
		tcpErrIPPortFail:            prometheus.NewDesc("netscaler_tcp_err_ip_port_fail", "TCP IP port fail errors", baseLabels, nil),
		tcpErrBadStateConn:          prometheus.NewDesc("netscaler_tcp_err_bad_state_conn", "TCP bad state connection errors", baseLabels, nil),
		tcpErrRstThreshold:          prometheus.NewDesc("netscaler_tcp_err_rst_threshold", "TCP RST threshold errors", baseLabels, nil),
		tcpSynRate:                  prometheus.NewDesc("netscaler_tcp_syn_rate", "TCP SYN rate", baseLabels, nil),
		tcpSynProbeRate:             prometheus.NewDesc("netscaler_tcp_syn_probe_rate", "TCP SYN probe rate", baseLabels, nil),

		// Protocol IP metrics
		ipTotalRxPackets:          prometheus.NewDesc("netscaler_ip_rx_packets_total", "Total IP packets received", baseLabels, nil),
		ipTotalRxBytes:            prometheus.NewDesc("netscaler_ip_rx_bytes_total", "Total IP bytes received", baseLabels, nil),
		ipTotalTxPackets:          prometheus.NewDesc("netscaler_ip_tx_packets_total", "Total IP packets transmitted", baseLabels, nil),
		ipTotalTxBytes:            prometheus.NewDesc("netscaler_ip_tx_bytes_total", "Total IP bytes transmitted", baseLabels, nil),
		ipTotalRxMbits:            prometheus.NewDesc("netscaler_ip_rx_mbits_total", "Total IP Mbits received", baseLabels, nil),
		ipTotalTxMbits:            prometheus.NewDesc("netscaler_ip_tx_mbits_total", "Total IP Mbits transmitted", baseLabels, nil),
		ipTotalRoutedPackets:      prometheus.NewDesc("netscaler_ip_routed_packets_total", "Total routed packets", baseLabels, nil),
		ipTotalRoutedMbits:        prometheus.NewDesc("netscaler_ip_routed_mbits_total", "Total routed Mbits", baseLabels, nil),
		ipTotalFragments:          prometheus.NewDesc("netscaler_ip_fragments_total", "Total IP fragments", baseLabels, nil),
		ipTotalSuccReassembly:     prometheus.NewDesc("netscaler_ip_successful_reassembly_total", "Total successful reassemblies", baseLabels, nil),
		ipTotalAddrLookup:         prometheus.NewDesc("netscaler_ip_address_lookup_total", "Total address lookups", baseLabels, nil),
		ipTotalAddrLookupFail:     prometheus.NewDesc("netscaler_ip_address_lookup_fail_total", "Total failed address lookups", baseLabels, nil),
		ipTotalUDPFragmentsFwd:    prometheus.NewDesc("netscaler_ip_udp_fragments_fwd_total", "Total UDP fragments forwarded", baseLabels, nil),
		ipTotalTCPFragmentsFwd:    prometheus.NewDesc("netscaler_ip_tcp_fragments_fwd_total", "Total TCP fragments forwarded", baseLabels, nil),
		ipTotalBadChecksums:       prometheus.NewDesc("netscaler_ip_bad_checksums_total", "Total bad checksums", baseLabels, nil),
		ipTotalUnsuccReassembly:   prometheus.NewDesc("netscaler_ip_unsuccessful_reassembly_total", "Total unsuccessful reassemblies", baseLabels, nil),
		ipTotalTooBig:             prometheus.NewDesc("netscaler_ip_too_big_total", "Total too big packets", baseLabels, nil),
		ipTotalDupFragments:       prometheus.NewDesc("netscaler_ip_duplicate_fragments_total", "Total duplicate fragments", baseLabels, nil),
		ipTotalOutOfOrderFrag:     prometheus.NewDesc("netscaler_ip_out_of_order_fragments_total", "Total out of order fragments", baseLabels, nil),
		ipTotalVIPDown:            prometheus.NewDesc("netscaler_ip_vip_down_total", "Total VIP down events", baseLabels, nil),
		ipTotalTTLExpired:         prometheus.NewDesc("netscaler_ip_ttl_expired_total", "Total TTL expired", baseLabels, nil),
		ipTotalMaxClients:         prometheus.NewDesc("netscaler_ip_max_clients_total", "Total max clients reached", baseLabels, nil),
		ipTotalUnknownSvcs:        prometheus.NewDesc("netscaler_ip_unknown_services_total", "Total unknown services", baseLabels, nil),
		ipTotalInvalidHeaderSz:    prometheus.NewDesc("netscaler_ip_invalid_header_size_total", "Total invalid header sizes", baseLabels, nil),
		ipTotalInvalidPacketSize:  prometheus.NewDesc("netscaler_ip_invalid_packet_size_total", "Total invalid packet sizes", baseLabels, nil),
		ipTotalTruncatedPackets:   prometheus.NewDesc("netscaler_ip_truncated_packets_total", "Total truncated packets", baseLabels, nil),
		ipNonIPTotalTruncatedPkts: prometheus.NewDesc("netscaler_ip_non_ip_truncated_packets_total", "Total non-IP truncated packets", baseLabels, nil),
		ipTotalBadMacAddrs:        prometheus.NewDesc("netscaler_ip_bad_mac_addresses_total", "Total bad MAC addresses", baseLabels, nil),
		ipRxPacketsRate:           prometheus.NewDesc("netscaler_ip_rx_packets_rate", "IP packets received rate", baseLabels, nil),
		ipRxBytesRate:             prometheus.NewDesc("netscaler_ip_rx_bytes_rate", "IP bytes received rate", baseLabels, nil),
		ipTxPacketsRate:           prometheus.NewDesc("netscaler_ip_tx_packets_rate", "IP packets transmitted rate", baseLabels, nil),
		ipTxBytesRate:             prometheus.NewDesc("netscaler_ip_tx_bytes_rate", "IP bytes transmitted rate", baseLabels, nil),
		ipRxMbitsRate:             prometheus.NewDesc("netscaler_ip_rx_mbits_rate", "IP Mbits received rate", baseLabels, nil),
		ipTxMbitsRate:             prometheus.NewDesc("netscaler_ip_tx_mbits_rate", "IP Mbits transmitted rate", baseLabels, nil),
		ipRoutedPacketsRate:       prometheus.NewDesc("netscaler_ip_routed_packets_rate", "Routed packets rate", baseLabels, nil),
		ipRoutedMbitsRate:         prometheus.NewDesc("netscaler_ip_routed_mbits_rate", "Routed Mbits rate", baseLabels, nil),

		// SSL global metrics
		sslTotalTLSv11Sessions:  prometheus.NewDesc("netscaler_ssl_tls11_sessions_total", "Total TLS v1.1 sessions", baseLabels, nil),
		sslTotalSSLv2Sessions:   prometheus.NewDesc("netscaler_ssl_v2_sessions_total", "Total SSL v2 sessions", baseLabels, nil),
		sslTotalSessions:        prometheus.NewDesc("netscaler_ssl_sessions_total", "Total SSL sessions", baseLabels, nil),
		sslTotalSSLv2Handshakes: prometheus.NewDesc("netscaler_ssl_v2_handshakes_total", "Total SSL v2 handshakes", baseLabels, nil),
		sslTotalEnc:             prometheus.NewDesc("netscaler_ssl_encode_total", "Total SSL encodes", baseLabels, nil),
		sslCryptoUtilization:    prometheus.NewDesc("netscaler_ssl_crypto_utilization", "SSL crypto utilization", baseLabels, nil),
		sslTotalNewSessions:     prometheus.NewDesc("netscaler_ssl_new_sessions_total", "Total new SSL sessions", baseLabels, nil),
		sslSessionsRate:         prometheus.NewDesc("netscaler_ssl_sessions_rate", "SSL sessions rate", baseLabels, nil),
		sslDecRate:              prometheus.NewDesc("netscaler_ssl_decode_rate", "SSL decode rate", baseLabels, nil),
		sslEncRate:              prometheus.NewDesc("netscaler_ssl_encode_rate", "SSL encode rate", baseLabels, nil),
		sslSSLv2HandshakesRate:  prometheus.NewDesc("netscaler_ssl_v2_handshakes_rate", "SSL v2 handshakes rate", baseLabels, nil),
		sslNewSessionsRate:      prometheus.NewDesc("netscaler_ssl_new_sessions_rate", "New SSL sessions rate", baseLabels, nil),

		// SSL certificate metrics
		sslCertDaysToExpire: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_ssl_cert_days_to_expire", Help: "Days until SSL certificate expires"}, sslCertLabels),

		// SSL VServer metrics
		sslVServerTotalDecBytes:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_decrypt_bytes_total", Help: "Total bytes decrypted"}, sslVsLabels),
		sslVServerTotalEncBytes:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_encrypt_bytes_total", Help: "Total bytes encrypted"}, sslVsLabels),
		sslVServerTotalHWDecBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_hw_decrypt_bytes_total", Help: "Total hardware decrypted bytes"}, sslVsLabels),
		sslVServerTotalHWEncBytes:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_hw_encrypt_bytes_total", Help: "Total hardware encrypted bytes"}, sslVsLabels),
		sslVServerTotalSessionNew:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_session_new_total", Help: "Total new sessions"}, sslVsLabels),
		sslVServerTotalSessionHits:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_session_hits_total", Help: "Total session hits"}, sslVsLabels),
		sslVServerTotalClientAuthSuccess: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_client_auth_success_total", Help: "Total client auth successes"}, sslVsLabels),
		sslVServerTotalClientAuthFailure: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_client_auth_failure_total", Help: "Total client auth failures"}, sslVsLabels),
		sslVServerHealth:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_health", Help: "SSL vserver health"}, sslVsLabels),
		sslVServerActiveServices:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_active_services", Help: "Active services"}, sslVsLabels),
		sslVServerClientAuthSuccessRate:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_client_auth_success_rate", Help: "Client auth success rate"}, sslVsLabels),
		sslVServerClientAuthFailureRate:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_client_auth_failure_rate", Help: "Client auth failure rate"}, sslVsLabels),
		sslVServerEncBytesRate:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_encrypt_bytes_rate", Help: "Encrypt bytes rate"}, sslVsLabels),
		sslVServerDecBytesRate:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_decrypt_bytes_rate", Help: "Decrypt bytes rate"}, sslVsLabels),
		sslVServerHWEncBytesRate:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_hw_encrypt_bytes_rate", Help: "HW encrypt bytes rate"}, sslVsLabels),
		sslVServerHWDecBytesRate:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_hw_decrypt_bytes_rate", Help: "HW decrypt bytes rate"}, sslVsLabels),
		sslVServerSessionNewRate:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_session_new_rate", Help: "New session rate"}, sslVsLabels),
		sslVServerSessionHitsRate:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_sslvserver_session_hits_rate", Help: "Session hits rate"}, sslVsLabels),

		// System CPU per-core metrics
		cpuCoreUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "netscaler_cpu_core_usage_percent", Help: "CPU usage per core"}, cpuCoreLabels),

		// Bandwidth capacity metrics
		capacityMaxBandwidth:    prometheus.NewDesc("netscaler_capacity_max_bandwidth", "Maximum licensed bandwidth in Mbps", baseLabels, nil),
		capacityMinBandwidth:    prometheus.NewDesc("netscaler_capacity_min_bandwidth", "Minimum licensed bandwidth in Mbps", baseLabels, nil),
		capacityActualBandwidth: prometheus.NewDesc("netscaler_capacity_actual_bandwidth", "Actual bandwidth in Mbps", baseLabels, nil),
		capacityBandwidth:       prometheus.NewDesc("netscaler_capacity_allocated_bandwidth", "Allocated licensed bandwidth in Mbps", baseLabels, nil),
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

	// Probe success
	e.probeSuccess.Describe(ch)

	// Protocol HTTP metrics
	ch <- e.httpTotalRequests
	ch <- e.httpTotalResponses
	ch <- e.httpTotalPosts
	ch <- e.httpTotalGets
	ch <- e.httpTotalOthers
	ch <- e.httpTotalRxRequestBytes
	ch <- e.httpTotalRxResponseBytes
	ch <- e.httpTotalTxRequestBytes
	ch <- e.httpTotal10Requests
	ch <- e.httpTotal11Requests
	ch <- e.httpTotal10Responses
	ch <- e.httpTotal11Responses
	ch <- e.httpTotalChunkedRequests
	ch <- e.httpTotalChunkedResponses
	ch <- e.httpTotalSPDYStreams
	ch <- e.httpTotalSPDYv2Streams
	ch <- e.httpTotalSPDYv3Streams
	ch <- e.httpErrNoReuseMultipart
	ch <- e.httpErrIncompleteHeaders
	ch <- e.httpErrIncompleteRequests
	ch <- e.httpErrIncompleteResponses
	ch <- e.httpErrServerBusy
	ch <- e.httpErrLargeContent
	ch <- e.httpErrLargeChunk
	ch <- e.httpErrLargeCtlen
	ch <- e.httpRequestsRate
	ch <- e.httpResponsesRate
	ch <- e.httpPostsRate
	ch <- e.httpGetsRate
	ch <- e.httpOthersRate
	ch <- e.httpRxRequestBytesRate
	ch <- e.httpRxResponseBytesRate
	ch <- e.httpTxRequestBytesRate
	ch <- e.httpRequest10Rate
	ch <- e.httpRequest11Rate
	ch <- e.httpResponse10Rate
	ch <- e.httpResponse11Rate
	ch <- e.httpChunkedRequestsRate
	ch <- e.httpChunkedResponsesRate
	ch <- e.httpSPDYStreamsRate
	ch <- e.httpSPDYv2StreamsRate
	ch <- e.httpSPDYv3StreamsRate
	ch <- e.httpErrNoReuseMultipartRate
	ch <- e.httpErrIncompleteRequestsRate
	ch <- e.httpErrIncompleteResponsesRate
	ch <- e.httpErrServerBusyRate

	// Protocol TCP metrics
	ch <- e.tcpTotalRxPackets
	ch <- e.tcpTotalRxBytes
	ch <- e.tcpTotalTxBytes
	ch <- e.tcpTotalTxPackets
	ch <- e.tcpTotalClientConnOpened
	ch <- e.tcpTotalServerConnOpened
	ch <- e.tcpTotalSyn
	ch <- e.tcpTotalSynProbe
	ch <- e.tcpTotalServerFin
	ch <- e.tcpTotalClientFin
	ch <- e.tcpActiveServerConn
	ch <- e.tcpCurClientConnEstablished
	ch <- e.tcpCurServerConnEstablished
	ch <- e.tcpRxPacketsRate
	ch <- e.tcpRxBytesRate
	ch <- e.tcpTxPacketsRate
	ch <- e.tcpTxBytesRate
	ch <- e.tcpClientConnOpenedRate
	ch <- e.tcpErrBadChecksum
	ch <- e.tcpErrBadChecksumRate
	ch <- e.tcpErrAnyPortFail
	ch <- e.tcpErrIPPortFail
	ch <- e.tcpErrBadStateConn
	ch <- e.tcpErrRstThreshold
	ch <- e.tcpSynRate
	ch <- e.tcpSynProbeRate

	// Protocol IP metrics
	ch <- e.ipTotalRxPackets
	ch <- e.ipTotalRxBytes
	ch <- e.ipTotalTxPackets
	ch <- e.ipTotalTxBytes
	ch <- e.ipTotalRxMbits
	ch <- e.ipTotalTxMbits
	ch <- e.ipTotalRoutedPackets
	ch <- e.ipTotalRoutedMbits
	ch <- e.ipTotalFragments
	ch <- e.ipTotalSuccReassembly
	ch <- e.ipTotalAddrLookup
	ch <- e.ipTotalAddrLookupFail
	ch <- e.ipTotalUDPFragmentsFwd
	ch <- e.ipTotalTCPFragmentsFwd
	ch <- e.ipTotalBadChecksums
	ch <- e.ipTotalUnsuccReassembly
	ch <- e.ipTotalTooBig
	ch <- e.ipTotalDupFragments
	ch <- e.ipTotalOutOfOrderFrag
	ch <- e.ipTotalVIPDown
	ch <- e.ipTotalTTLExpired
	ch <- e.ipTotalMaxClients
	ch <- e.ipTotalUnknownSvcs
	ch <- e.ipTotalInvalidHeaderSz
	ch <- e.ipTotalInvalidPacketSize
	ch <- e.ipTotalTruncatedPackets
	ch <- e.ipNonIPTotalTruncatedPkts
	ch <- e.ipTotalBadMacAddrs
	ch <- e.ipRxPacketsRate
	ch <- e.ipRxBytesRate
	ch <- e.ipTxPacketsRate
	ch <- e.ipTxBytesRate
	ch <- e.ipRxMbitsRate
	ch <- e.ipTxMbitsRate
	ch <- e.ipRoutedPacketsRate
	ch <- e.ipRoutedMbitsRate

	// SSL global metrics
	ch <- e.sslTotalTLSv11Sessions
	ch <- e.sslTotalSSLv2Sessions
	ch <- e.sslTotalSessions
	ch <- e.sslTotalSSLv2Handshakes
	ch <- e.sslTotalEnc
	ch <- e.sslCryptoUtilization
	ch <- e.sslTotalNewSessions
	ch <- e.sslSessionsRate
	ch <- e.sslDecRate
	ch <- e.sslEncRate
	ch <- e.sslSSLv2HandshakesRate
	ch <- e.sslNewSessionsRate

	// SSL cert metrics
	e.sslCertDaysToExpire.Describe(ch)

	// SSL VServer metrics
	e.sslVServerTotalDecBytes.Describe(ch)
	e.sslVServerTotalEncBytes.Describe(ch)
	e.sslVServerTotalHWDecBytes.Describe(ch)
	e.sslVServerTotalHWEncBytes.Describe(ch)
	e.sslVServerTotalSessionNew.Describe(ch)
	e.sslVServerTotalSessionHits.Describe(ch)
	e.sslVServerTotalClientAuthSuccess.Describe(ch)
	e.sslVServerTotalClientAuthFailure.Describe(ch)
	e.sslVServerHealth.Describe(ch)
	e.sslVServerActiveServices.Describe(ch)
	e.sslVServerClientAuthSuccessRate.Describe(ch)
	e.sslVServerClientAuthFailureRate.Describe(ch)
	e.sslVServerEncBytesRate.Describe(ch)
	e.sslVServerDecBytesRate.Describe(ch)
	e.sslVServerHWEncBytesRate.Describe(ch)
	e.sslVServerHWDecBytesRate.Describe(ch)
	e.sslVServerSessionNewRate.Describe(ch)
	e.sslVServerSessionHitsRate.Describe(ch)

	// CPU core metrics
	e.cpuCoreUsage.Describe(ch)

	// Bandwidth capacity metrics
	ch <- e.capacityMaxBandwidth
	ch <- e.capacityMinBandwidth
	ch <- e.capacityActualBandwidth
	ch <- e.capacityBandwidth
}
