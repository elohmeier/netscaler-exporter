package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// LB Virtual Server metrics
var (
	virtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersWaitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_waiting_requests",
			Help: "Number of requests waiting on a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersTotalHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	virtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{"ns_instance", "virtual_server"},
	)
)

// GSLB Virtual Server metrics
var (
	gslbVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersTotalHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	gslbVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{"ns_instance", "virtual_server"},
	)
)

// CS Virtual Server metrics
var (
	csVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_established_connections",
			Help: "Number of client connections in ESTABLISHED state.",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalPacketsReceived = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_packets_received",
			Help: "Total number of packets received",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalPacketsSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_packets_sent",
			Help: "Total number of packets sent.",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalSpillovers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_spillovers",
			Help: "Number of times vserver experienced spill over.",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersDeferredRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_deferred_requests",
			Help: "Number of deferred request on this vserver",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersNumberInvalidRequestResponse = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_number_invalid_request_response",
			Help: "Number invalid requests/responses on this vserver",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersNumberInvalidRequestResponseDropped = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_number_invalid_request_response_dropped",
			Help: "Number invalid requests/responses dropped on this vserver",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersTotalVServerDownBackupHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_total_vserver_down_backup_hits",
			Help: "Number of times traffic was diverted to backup vserver since primary vserver was DOWN.",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersCurrentMultipathSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_sessions",
			Help: "Current Multipath TCP sessions",
		},
		[]string{"ns_instance", "virtual_server"},
	)

	csVirtualServersCurrentMultipathSubflows = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_subflows",
			Help: "Current Multipath TCP subflows",
		},
		[]string{"ns_instance", "virtual_server"},
	)
)

// VPN Virtual Server metrics
var (
	vpnVirtualServersTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_total_requests",
			Help: "Total VPN virtual server requests",
		},
		[]string{"ns_instance", "vpn_virtual_server"},
	)

	vpnVirtualServersTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_total_responses",
			Help: "Total VPN virtual server responses",
		},
		[]string{"ns_instance", "vpn_virtual_server"},
	)

	vpnVirtualServersTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_total_request_bytes",
			Help: "Total VPN virtual server request bytes",
		},
		[]string{"ns_instance", "vpn_virtual_server"},
	)

	vpnVirtualServersTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_total_response_bytes",
			Help: "Total VPN virtual server response bytes",
		},
		[]string{"ns_instance", "vpn_virtual_server"},
	)

	vpnVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_state",
			Help: "Current state of the VPN virtual server",
		},
		[]string{"ns_instance", "vpn_virtual_server"},
	)
)

// LB Virtual Server collectors
func (e *Exporter) collectVirtualServerState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersState.Reset()
	for _, vs := range ns.VirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		e.virtualServersState.WithLabelValues(nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersWaitingRequests.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersHealth.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersInactiveServices.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersActiveServices.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersTotalHits.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersTotalRequests.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersTotalResponses.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.virtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

// GSLB Virtual Server collectors
func (e *Exporter) collectGSLBVirtualServerState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersState.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		e.gslbVirtualServersState.WithLabelValues(nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBVirtualServerHealth(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersHealth.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.Health, 64)
		e.gslbVirtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerInactiveServices(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersInactiveServices.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.gslbVirtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerActiveServices(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersActiveServices.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.gslbVirtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersTotalHits.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.gslbVirtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersTotalRequests.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.gslbVirtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersTotalResponses.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.gslbVirtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.gslbVirtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.gslbVirtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.gslbVirtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbVirtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.gslbVirtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

// CS Virtual Server collectors
func (e *Exporter) collectCSVirtualServerState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersState.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		e.csVirtualServersState.WithLabelValues(nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectCSVirtualServerTotalHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalHits.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.csVirtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalRequests.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.csVirtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalResponses.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.csVirtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.csVirtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.csVirtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.csVirtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.csVirtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerEstablishedConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersEstablishedConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.EstablishedConnections, 64)
		e.csVirtualServersEstablishedConnections.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsReceived(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalPacketsReceived.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalPacketsReceived, 64)
		e.csVirtualServersTotalPacketsReceived.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsSent(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalPacketsSent.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalPacketsSent, 64)
		e.csVirtualServersTotalPacketsSent.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalSpillovers(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalSpillovers.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalSpillovers, 64)
		e.csVirtualServersTotalSpillovers.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerDeferredRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersDeferredRequests.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.DeferredRequests, 64)
		e.csVirtualServersDeferredRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponse(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersNumberInvalidRequestResponse.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InvalidRequestResponse, 64)
		e.csVirtualServersNumberInvalidRequestResponse.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponseDropped(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersNumberInvalidRequestResponseDropped.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InvalidRequestResponseDropped, 64)
		e.csVirtualServersNumberInvalidRequestResponseDropped.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalVServerDownBackupHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersTotalVServerDownBackupHits.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalVServerDownBackupHits, 64)
		e.csVirtualServersTotalVServerDownBackupHits.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSessions(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersCurrentMultipathSessions.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentMultipathSessions, 64)
		e.csVirtualServersCurrentMultipathSessions.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSubflows(ns netscaler.NSAPIResponse, nsInstance string) {
	e.csVirtualServersCurrentMultipathSubflows.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentMultipathSubflows, 64)
		e.csVirtualServersCurrentMultipathSubflows.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

// VPN Virtual Server collectors
func (e *Exporter) collectVPNVirtualServerTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.vpnVirtualServersTotalRequests.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.vpnVirtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.vpnVirtualServersTotalResponses.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.vpnVirtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.vpnVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.vpnVirtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.vpnVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.vpnVirtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.vpnVirtualServersState.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		e.vpnVirtualServersState.WithLabelValues(nsInstance, vs.Name).Set(state)
	}
}
