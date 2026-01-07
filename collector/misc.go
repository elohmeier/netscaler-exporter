package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
)

// Interface collectors
func (e *Exporter) collectInterfacesRxBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesRxBytes.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedBytes, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesRxBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxBytes(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesTxBytes.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitBytes, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesTxBytes.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesRxPackets(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesRxPackets.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedPackets, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesRxPackets.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxPackets(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesTxPackets.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitPackets, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesTxPackets.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRx(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesJumboPacketsRx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsReceived, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesJumboPacketsRx.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTx(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesJumboPacketsTx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsTransmitted, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesJumboPacketsTx.WithLabelValues(labels...).Set(val)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRx(ns netscaler.NSAPIResponse, target config.Target) {
	e.interfacesErrorPacketsRx.Reset()
	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.ErrorPacketsReceived, 64)
		labels := e.buildLabelValues(target, iface.ID, iface.Alias)
		e.interfacesErrorPacketsRx.WithLabelValues(labels...).Set(val)
	}
}

// AAA collectors
func (e *Exporter) collectAaaAuthSuccess(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaAuthSuccess.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthSuccess, 64)
	labels := e.buildLabelValues(target)
	e.aaaAuthSuccess.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectAaaAuthFail(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaAuthFail.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthFail, 64)
	labels := e.buildLabelValues(target)
	e.aaaAuthFail.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPSuccess(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaAuthOnlyHTTPSuccess.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPSuccess, 64)
	labels := e.buildLabelValues(target)
	e.aaaAuthOnlyHTTPSuccess.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPFail(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaAuthOnlyHTTPFail.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPFail, 64)
	labels := e.buildLabelValues(target)
	e.aaaAuthOnlyHTTPFail.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectAaaCurIcaSessions(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaCurIcaSessions.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaSessions, 64)
	labels := e.buildLabelValues(target)
	e.aaaCurIcaSessions.WithLabelValues(labels...).Set(val)
}

func (e *Exporter) collectAaaCurIcaOnlyConn(ns netscaler.NSAPIResponse, target config.Target) {
	e.aaaCurIcaOnlyConn.Reset()
	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaOnlyConnections, 64)
	labels := e.buildLabelValues(target)
	e.aaaCurIcaOnlyConn.WithLabelValues(labels...).Set(val)
}
