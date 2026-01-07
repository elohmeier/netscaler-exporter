package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// Service metrics
var (
	servicesThroughput = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_throughput",
			Help: "Number of bytes received or sent by this service (Mbps)",
		},
		[]string{"ns_instance", "service"},
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.",
		},
		[]string{"ns_instance", "service"},
	)

	servicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_state",
			Help: "Current state of the service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{"ns_instance", "service"},
	)

	servicesSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_surge_count",
			Help: "Number of requests in the surge queue",
		},
		[]string{"ns_instance", "service"},
	)

	servicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{"ns_instance", "service"},
	)

	servicesServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{"ns_instance", "service"},
	)

	servicesCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{"ns_instance", "service"},
	)

	servicesMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_max_clients",
			Help: "Maximum open connections allowed on this service",
		},
		[]string{"ns_instance", "service"},
	)

	servicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{"ns_instance", "service"},
	)

	servicesVirtualServerServiceHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{"ns_instance", "service"},
	)

	servicesActiveTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_active_transactions",
			Help: "Number of active transactions handled by this service.",
		},
		[]string{"ns_instance", "service"},
	)
)

// GSLB Service metrics
var (
	gslbServicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_state",
			Help: "Current state of the service",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{"ns_instance", "service"},
	)

	gslbServicesVirtualServerServiceHits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{"ns_instance", "service"},
	)
)

// Service Group metrics
var (
	serviceGroupsState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_state",
			Help: "Current state of the server",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsTotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsTotalResponses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_total_responses",
			Help: "Number of responses received on this service.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsTotalRequestBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsTotalResponseBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_total_response_bytes",
			Help: "Number of response bytes received by this service",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_client_connections",
			Help: "Number of current client connections.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_surge_count",
			Help: "Number of requests in the surge queue.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_server_connections",
			Help: "Number of current connections to the actual servers behind the virtual server.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)

	serviceGroupsMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_max_clients",
			Help: "Maximum open connections allowed on this service.",
		},
		[]string{"ns_instance", "servicegroup", "member", "port"},
	)
)

// Service collectors
func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesThroughput.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		e.servicesThroughput.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesAvgTTFB.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		e.servicesAvgTTFB.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesState.Reset()
	for _, service := range ns.ServiceStats {
		state := 0.0
		if service.State == "UP" {
			state = 1.0
		}
		e.servicesState.WithLabelValues(nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesTotalRequests.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesTotalResponses.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesTotalRequestBytes.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesTotalResponseBytes.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesCurrentClientConns.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesSurgeCount.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesCurrentServerConns.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesServerEstablishedConnections.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesCurrentReusePool.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesMaxClients.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesCurrentLoad.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesVirtualServerServiceHits.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse, nsInstance string) {
	e.servicesActiveTransactions.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

// GSLB Service collectors
func (e *Exporter) collectGSLBServicesState(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesState.Reset()
	for _, service := range ns.GSLBServiceStats {
		state := 0.0
		if service.State == "UP" {
			state = 1.0
		}
		e.gslbServicesState.WithLabelValues(nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequests(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesTotalRequests.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.gslbServicesTotalRequests.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponses(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesTotalResponses.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.gslbServicesTotalResponses.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequestBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesTotalRequestBytes.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.gslbServicesTotalRequestBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponseBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesTotalResponseBytes.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.gslbServicesTotalResponseBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentClientConns(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesCurrentClientConns.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.gslbServicesCurrentClientConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentServerConns(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesCurrentServerConns.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.gslbServicesCurrentServerConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesEstablishedConnections(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesEstablishedConnections.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.EstablishedConnections, 64)
		e.gslbServicesEstablishedConnections.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentLoad(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesCurrentLoad.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.gslbServicesCurrentLoad.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse, nsInstance string) {
	e.gslbServicesVirtualServerServiceHits.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.gslbServicesVirtualServerServiceHits.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

// Service Group collectors
func (e *Exporter) collectServiceGroupsState(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsState.Reset()
	state := 0.0
	if sg.State == "UP" {
		state = 1.0
	}
	port := strconv.Itoa(sg.PrimaryPort)
	e.serviceGroupsState.WithLabelValues(nsInstance, sgName, servername, port).Set(state)
}

func (e *Exporter) collectServiceGroupsAvgTTFB(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsAvgTTFB.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
	e.serviceGroupsAvgTTFB.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequests(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsTotalRequests.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
	e.serviceGroupsTotalRequests.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponses(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsTotalResponses.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
	e.serviceGroupsTotalResponses.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsTotalRequestBytes.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
	e.serviceGroupsTotalRequestBytes.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsTotalResponseBytes.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
	e.serviceGroupsTotalResponseBytes.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsCurrentClientConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
	e.serviceGroupsCurrentClientConnections.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsSurgeCount(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsSurgeCount.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
	e.serviceGroupsSurgeCount.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsCurrentServerConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
	e.serviceGroupsCurrentServerConnections.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsServerEstablishedConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
	e.serviceGroupsServerEstablishedConnections.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsCurrentReusePool.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
	e.serviceGroupsCurrentReusePool.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsMaxClients(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, nsInstance string) {
	e.serviceGroupsMaxClients.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.MaxClients, 64)
	e.serviceGroupsMaxClients.WithLabelValues(nsInstance, sgName, servername, port).Set(val)
}
