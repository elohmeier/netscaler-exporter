package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// NetScaler system metrics (descriptors)
var (
	modelID = prometheus.NewDesc(
		"model_id",
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		[]string{"ns_instance"},
		nil,
	)

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		[]string{"ns_instance"},
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		[]string{"ns_instance"},
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		[]string{"ns_instance"},
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		[]string{"ns_instance"},
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		[]string{"ns_instance"},
		nil,
	)

	totRxMB = prometheus.NewDesc(
		"total_received_mb",
		"Total number of Megabytes received by the NetScaler appliance",
		[]string{"ns_instance"},
		nil,
	)

	totTxMB = prometheus.NewDesc(
		"total_transmit_mb",
		"Total number of Megabytes transmitted by the NetScaler appliance",
		[]string{"ns_instance"},
		nil,
	)

	httpRequests = prometheus.NewDesc(
		"http_requests",
		"Total number of HTTP requests received",
		[]string{"ns_instance"},
		nil,
	)

	httpResponses = prometheus.NewDesc(
		"http_responses",
		"Total number of HTTP responses sent",
		[]string{"ns_instance"},
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		"tcp_current_client_connections",
		"Client connections, including connections in the Opening, Established, and Closing state.",
		[]string{"ns_instance"},
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_client_connections_established",
		"Current client connections in the Established state.",
		[]string{"ns_instance"},
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		"tcp_current_server_connections",
		"Server connections, including connections in the Opening, Established, and Closing state.",
		[]string{"ns_instance"},
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_server_connections_established",
		"Current server connections in the Established state.",
		[]string{"ns_instance"},
		nil,
	)
)

// Interface metrics
var (
	interfacesRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes",
			Help: "Number of bytes received by specific interfaces.",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes",
			Help: "Number of bytes transmitted by specific interfaces.",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets",
			Help: "Number of packets received by specific interfaces",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets",
			Help: "Number of packets transmitted by specific interfaces",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesJumboPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received",
			Help: "Number of bytes received by specific interfaces",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesJumboPacketsTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted",
			Help: "Number of jumbo packets transmitted by specific interfaces",
		},
		[]string{"ns_instance", "interface", "alias"},
	)

	interfacesErrorPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received",
			Help: "Number of error packets received by specific interfaces",
		},
		[]string{"ns_instance", "interface", "alias"},
	)
)

// AAA metrics
var (
	aaaAuthSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_auth_success",
			Help: "Count of authentication successes",
		},
		[]string{"ns_instance"},
	)

	aaaAuthFail = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_auth_fail",
			Help: "Count of authentication failures",
		},
		[]string{"ns_instance"},
	)

	aaaAuthOnlyHTTPSuccess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_auth_only_http_success",
			Help: "Count of HTTP connections that succeeded authorisation",
		},
		[]string{"ns_instance"},
	)

	aaaAuthOnlyHTTPFail = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_auth_only_http_fail",
			Help: "Count of HTTP connections that failed authorisation",
		},
		[]string{"ns_instance"},
	)

	aaaCurIcaSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_current_ica_sessions",
			Help: "Count of current Basic ICA only sessions",
		},
		[]string{"ns_instance"},
	)

	aaaCurIcaOnlyConn = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aaa_current_ica_only_connections",
			Help: "Count of current Basic ICA only connections",
		},
		[]string{"ns_instance"},
	)
)

// Interface collectors
func (e *Exporter) collectInterfacesRxBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesRxBytes.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedBytes, 64)
		e.interfacesRxBytes.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxBytes(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesTxBytes.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitBytes, 64)
		e.interfacesTxBytes.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesRxPackets(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesRxPackets.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedPackets, 64)
		e.interfacesRxPackets.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxPackets(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesTxPackets.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitPackets, 64)
		e.interfacesTxPackets.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRx(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesJumboPacketsRx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsReceived, 64)
		e.interfacesJumboPacketsRx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTx(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesJumboPacketsTx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsTransmitted, 64)
		e.interfacesJumboPacketsTx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRx(ns netscaler.NSAPIResponse, nsInstance string) {
	e.interfacesErrorPacketsRx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.ErrorPacketsReceived, 64)
		e.interfacesErrorPacketsRx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

// AAA collectors
func (e *Exporter) collectAaaAuthSuccess(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaAuthSuccess.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthSuccess, 64)
	e.aaaAuthSuccess.WithLabelValues(nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthFail(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaAuthFail.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthFail, 64)
	e.aaaAuthFail.WithLabelValues(nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPSuccess(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaAuthOnlyHTTPSuccess.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPSuccess, 64)
	e.aaaAuthOnlyHTTPSuccess.WithLabelValues(nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPFail(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaAuthOnlyHTTPFail.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPFail, 64)
	e.aaaAuthOnlyHTTPFail.WithLabelValues(nsInstance).Set(val)
}

func (e *Exporter) collectAaaCurIcaSessions(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaCurIcaSessions.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaSessions, 64)
	e.aaaCurIcaSessions.WithLabelValues(nsInstance).Set(val)
}

func (e *Exporter) collectAaaCurIcaOnlyConn(ns netscaler.NSAPIResponse, nsInstance string) {
	e.aaaCurIcaOnlyConn.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaOnlyConnections, 64)
	e.aaaCurIcaOnlyConn.WithLabelValues(nsInstance).Set(val)
}
