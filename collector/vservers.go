package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
)

// LB Virtual Server collectors
func (e *Exporter) collectVirtualServerState(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersState.Reset()
	for _, vs := range ns.VirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersState.WithLabelValues(labels...).Set(state)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersWaitingRequests.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersWaitingRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersHealth.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.Health, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersHealth.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersInactiveServices.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersInactiveServices.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersActiveServices.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersActiveServices.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersTotalHits.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersTotalHits.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersTotalRequests.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersTotalResponses.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersCurrentClientConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.virtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.VirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.virtualServersCurrentServerConnections.WithLabelValues(labels...).Set(val)
	}
}

// GSLB Virtual Server collectors
func (e *Exporter) collectGSLBVirtualServerState(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersState.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersState.WithLabelValues(labels...).Set(state)
	}
}

func (e *Exporter) collectGSLBVirtualServerHealth(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersHealth.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.Health, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersHealth.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerInactiveServices(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersInactiveServices.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersInactiveServices.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerActiveServices(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersActiveServices.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersActiveServices.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersTotalHits.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersTotalHits.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersTotalRequests.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersTotalResponses.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersCurrentClientConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbVirtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.GSLBVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.gslbVirtualServersCurrentServerConnections.WithLabelValues(labels...).Set(val)
	}
}

// CS Virtual Server collectors
func (e *Exporter) collectCSVirtualServerState(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersState.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersState.WithLabelValues(labels...).Set(state)
	}
}

func (e *Exporter) collectCSVirtualServerTotalHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalHits.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalHits, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalHits.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalRequests.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalResponses.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersCurrentClientConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersCurrentClientConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersCurrentServerConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersCurrentServerConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerEstablishedConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersEstablishedConnections.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.EstablishedConnections, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersEstablishedConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsReceived(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalPacketsReceived.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalPacketsReceived, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalPacketsReceived.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsSent(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalPacketsSent.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalPacketsSent, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalPacketsSent.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalSpillovers(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalSpillovers.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalSpillovers, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalSpillovers.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerDeferredRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersDeferredRequests.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.DeferredRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersDeferredRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponse(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersNumberInvalidRequestResponse.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InvalidRequestResponse, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersNumberInvalidRequestResponse.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponseDropped(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersNumberInvalidRequestResponseDropped.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.InvalidRequestResponseDropped, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersNumberInvalidRequestResponseDropped.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerTotalVServerDownBackupHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersTotalVServerDownBackupHits.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalVServerDownBackupHits, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersTotalVServerDownBackupHits.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSessions(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersCurrentMultipathSessions.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentMultipathSessions, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersCurrentMultipathSessions.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSubflows(ns netscaler.NSAPIResponse, target config.Target) {
	e.csVirtualServersCurrentMultipathSubflows.Reset()
	for _, vs := range ns.CSVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.CurrentMultipathSubflows, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.csVirtualServersCurrentMultipathSubflows.WithLabelValues(labels...).Set(val)
	}
}

// VPN Virtual Server collectors
func (e *Exporter) collectVPNVirtualServerTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.vpnVirtualServersTotalRequests.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.vpnVirtualServersTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.vpnVirtualServersTotalResponses.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.vpnVirtualServersTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.vpnVirtualServersTotalRequestBytes.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.vpnVirtualServersTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.vpnVirtualServersTotalResponseBytes.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		val, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, vs.Name)
		e.vpnVirtualServersTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectVPNVirtualServerState(ns netscaler.NSAPIResponse, target config.Target) {
	e.vpnVirtualServersState.Reset()
	for _, vs := range ns.VPNVirtualServerStats {
		state := 0.0
		if vs.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, vs.Name)
		e.vpnVirtualServersState.WithLabelValues(labels...).Set(state)
	}
}
