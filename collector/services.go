package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
)

// Service collectors
func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesThroughput.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesThroughput.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesAvgTTFB.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesAvgTTFB.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesState.Reset()
	for _, service := range ns.ServiceStats {
		state := 0.0
		if service.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, service.Name)
		e.servicesState.WithLabelValues(labels...).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesTotalRequests.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesTotalResponses.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesTotalRequestBytes.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesTotalResponseBytes.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesCurrentClientConns.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesCurrentClientConns.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesSurgeCount.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesSurgeCount.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesCurrentServerConns.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesCurrentServerConns.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesServerEstablishedConnections.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesServerEstablishedConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesCurrentReusePool.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesCurrentReusePool.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesMaxClients.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesMaxClients.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesCurrentLoad.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesCurrentLoad.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesVirtualServerServiceHits.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesVirtualServerServiceHits.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse, target config.Target) {
	e.servicesActiveTransactions.Reset()
	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.servicesActiveTransactions.WithLabelValues(labels...).Set(val)
	}
}

// GSLB Service collectors
func (e *Exporter) collectGSLBServicesState(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesState.Reset()
	for _, service := range ns.GSLBServiceStats {
		state := 0.0
		if service.State == "UP" {
			state = 1.0
		}
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesState.WithLabelValues(labels...).Set(state)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequests(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesTotalRequests.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesTotalRequests.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponses(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesTotalResponses.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesTotalResponses.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequestBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesTotalRequestBytes.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesTotalRequestBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponseBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesTotalResponseBytes.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesTotalResponseBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentClientConns(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesCurrentClientConns.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesCurrentClientConns.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentServerConns(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesCurrentServerConns.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesCurrentServerConns.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesEstablishedConnections(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesEstablishedConnections.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.EstablishedConnections, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesEstablishedConnections.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentLoad(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesCurrentLoad.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesCurrentLoad.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse, target config.Target) {
	e.gslbServicesVirtualServerServiceHits.Reset()
	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		labels := e.buildLabelValues(target, service.Name)
		e.gslbServicesVirtualServerServiceHits.WithLabelValues(labels...).Set(val)
	}
}

// Service Group collectors
func (e *Exporter) collectServiceGroupsState(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsState.Reset()
	state := 0.0
	if sg.State == "UP" {
		state = 1.0
	}
	port := strconv.Itoa(sg.PrimaryPort)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsState.WithLabelValues(labels...).Set(state)
}

func (e *Exporter) collectServiceGroupsAvgTTFB(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsAvgTTFB.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsAvgTTFB.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequests(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsTotalRequests.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsTotalRequests.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponses(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsTotalResponses.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsTotalResponses.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsTotalRequestBytes.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsTotalRequestBytes.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsTotalResponseBytes.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsTotalResponseBytes.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsCurrentClientConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsCurrentClientConnections.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsSurgeCount(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsSurgeCount.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsSurgeCount.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsCurrentServerConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsCurrentServerConnections.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsServerEstablishedConnections.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsServerEstablishedConnections.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsCurrentReusePool.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsCurrentReusePool.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectServiceGroupsMaxClients(sg netscaler.ServiceGroupMemberStats, sgName string, servername string, target config.Target) {
	e.serviceGroupsMaxClients.Reset()
	port := strconv.Itoa(sg.PrimaryPort)
	val, _ := strconv.ParseFloat(sg.MaxClients, 64)
	labels := e.buildLabelValues(target, sgName, servername, port)
	e.serviceGroupsMaxClients.WithLabelValues(labels...).Set(val)
}
