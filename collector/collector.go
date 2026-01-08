package collector

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	// Scrape all targets concurrently
	for _, target := range e.config.Targets {
		wg.Add(1)
		go func(t config.Target) {
			defer wg.Done()
			e.scrapeTarget(t, ch)
		}(target)
	}

	wg.Wait()
}

// scrapeTarget scrapes a single NetScaler target
func (e *Exporter) scrapeTarget(target config.Target, ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nsClient := netscaler.NewNitroClient(target.URL, e.username, e.password, e.ignoreCert)
	defer nsClient.CloseIdleConnections()

	var wg sync.WaitGroup
	// Semaphore to limit concurrent requests to avoid overloading the NetScaler
	sem := make(chan struct{}, e.parallelism)

	// Helper to run a scrape function concurrently
	run := func(name string, scrapeFn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}: // Acquire token
				defer func() { <-sem }() // Release token
				scrapeFn()
			case <-ctx.Done():
				e.logger.Warn("context cancelled, skipping scrape", "target", target.URL, "name", name)
			}
		}()
	}

	// Build base label values for this target
	baseLabels := e.buildLabelValues(target)

	// 1. NS Stats
	run("ns_stats", func() {
		ns, err := netscaler.GetNSStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get NS stats", "target", target.URL, "err", err)
			return
		}

		fltTotRxMB, _ := strconv.ParseFloat(ns.NSStats.TotalReceivedMB, 64)
		fltTotTxMB, _ := strconv.ParseFloat(ns.NSStats.TotalTransmitMB, 64)
		fltHTTPRequests, _ := strconv.ParseFloat(ns.NSStats.HTTPRequests, 64)
		fltHTTPResponses, _ := strconv.ParseFloat(ns.NSStats.HTTPResponses, 64)
		fltTCPCurrentClientConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnections, 64)
		fltTCPCurrentClientConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnectionsEstablished, 64)
		fltTCPCurrentServerConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnections, 64)
		fltTCPCurrentServerConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnectionsEstablished, 64)

		ch <- prometheus.MustNewConstMetric(e.mgmtCPUUsage, prometheus.GaugeValue, ns.NSStats.MgmtCPUUsagePcnt, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.memUsage, prometheus.GaugeValue, ns.NSStats.MemUsagePcnt, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.pktCPUUsage, prometheus.GaugeValue, ns.NSStats.PktCPUUsagePcnt, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.flashPartitionUsage, prometheus.GaugeValue, ns.NSStats.FlashPartitionUsage, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.varPartitionUsage, prometheus.GaugeValue, ns.NSStats.VarPartitionUsage, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.totRxMB, prometheus.GaugeValue, fltTotRxMB, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.totTxMB, prometheus.GaugeValue, fltTotTxMB, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.httpRequests, prometheus.GaugeValue, fltHTTPRequests, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.httpResponses, prometheus.GaugeValue, fltHTTPResponses, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.tcpCurrentClientConnections, prometheus.GaugeValue, fltTCPCurrentClientConnections, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.tcpCurrentClientConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentClientConnectionsEstablished, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.tcpCurrentServerConnections, prometheus.GaugeValue, fltTCPCurrentServerConnections, baseLabels...)
		ch <- prometheus.MustNewConstMetric(e.tcpCurrentServerConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentServerConnectionsEstablished, baseLabels...)
	})

	// 2. NS License
	run("ns_license", func() {
		nslicense, err := netscaler.GetNSLicense(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get NS license", "target", target.URL, "err", err)
			return
		}
		fltModelID, _ := strconv.ParseFloat(nslicense.NSLicense.ModelID, 64)
		ch <- prometheus.MustNewConstMetric(e.modelID, prometheus.GaugeValue, fltModelID, baseLabels...)
	})

	// 3. Interfaces
	run("interfaces", func() {
		interfaces, err := netscaler.GetInterfaceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get interface stats", "target", target.URL, "err", err)
			return
		}
		e.collectInterfacesRxBytes(interfaces, target)
		e.interfacesRxBytes.Collect(ch)
		e.collectInterfacesTxBytes(interfaces, target)
		e.interfacesTxBytes.Collect(ch)
		e.collectInterfacesRxPackets(interfaces, target)
		e.interfacesRxPackets.Collect(ch)
		e.collectInterfacesTxPackets(interfaces, target)
		e.interfacesTxPackets.Collect(ch)
		e.collectInterfacesJumboPacketsRx(interfaces, target)
		e.interfacesJumboPacketsRx.Collect(ch)
		e.collectInterfacesJumboPacketsTx(interfaces, target)
		e.interfacesJumboPacketsTx.Collect(ch)
		e.collectInterfacesErrorPacketsRx(interfaces, target)
		e.interfacesErrorPacketsRx.Collect(ch)
	})

	// 4. Virtual Servers
	run("virtual_servers", func() {
		virtualServers, err := netscaler.GetVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get virtual server stats", "target", target.URL, "err", err)
			return
		}
		e.collectVirtualServerState(virtualServers, target)
		e.virtualServersState.Collect(ch)
		e.collectVirtualServerWaitingRequests(virtualServers, target)
		e.virtualServersWaitingRequests.Collect(ch)
		e.collectVirtualServerHealth(virtualServers, target)
		e.virtualServersHealth.Collect(ch)
		e.collectVirtualServerInactiveServices(virtualServers, target)
		e.virtualServersInactiveServices.Collect(ch)
		e.collectVirtualServerActiveServices(virtualServers, target)
		e.virtualServersActiveServices.Collect(ch)
		e.collectVirtualServerTotalHits(virtualServers, target)
		e.virtualServersTotalHits.Collect(ch)
		e.collectVirtualServerTotalRequests(virtualServers, target)
		e.virtualServersTotalRequests.Collect(ch)
		e.collectVirtualServerTotalResponses(virtualServers, target)
		e.virtualServersTotalResponses.Collect(ch)
		e.collectVirtualServerTotalRequestBytes(virtualServers, target)
		e.virtualServersTotalRequestBytes.Collect(ch)
		e.collectVirtualServerTotalResponseBytes(virtualServers, target)
		e.virtualServersTotalResponseBytes.Collect(ch)
		e.collectVirtualServerCurrentClientConnections(virtualServers, target)
		e.virtualServersCurrentClientConnections.Collect(ch)
		e.collectVirtualServerCurrentServerConnections(virtualServers, target)
		e.virtualServersCurrentServerConnections.Collect(ch)
	})

	// 5. Services
	run("services", func() {
		services, err := netscaler.GetServiceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get service stats", "target", target.URL, "err", err)
			return
		}
		e.collectServicesThroughput(services, target)
		e.servicesThroughput.Collect(ch)
		e.collectServicesAvgTTFB(services, target)
		e.servicesAvgTTFB.Collect(ch)
		e.collectServicesState(services, target)
		e.servicesState.Collect(ch)
		e.collectServicesTotalRequests(services, target)
		e.servicesTotalRequests.Collect(ch)
		e.collectServicesTotalResponses(services, target)
		e.servicesTotalResponses.Collect(ch)
		e.collectServicesTotalRequestBytes(services, target)
		e.servicesTotalRequestBytes.Collect(ch)
		e.collectServicesTotalResponseBytes(services, target)
		e.servicesTotalResponseBytes.Collect(ch)
		e.collectServicesCurrentClientConns(services, target)
		e.servicesCurrentClientConns.Collect(ch)
		e.collectServicesSurgeCount(services, target)
		e.servicesSurgeCount.Collect(ch)
		e.collectServicesCurrentServerConns(services, target)
		e.servicesCurrentServerConns.Collect(ch)
		e.collectServicesServerEstablishedConnections(services, target)
		e.servicesServerEstablishedConnections.Collect(ch)
		e.collectServicesCurrentReusePool(services, target)
		e.servicesCurrentReusePool.Collect(ch)
		e.collectServicesMaxClients(services, target)
		e.servicesMaxClients.Collect(ch)
		e.collectServicesCurrentLoad(services, target)
		e.servicesCurrentLoad.Collect(ch)
		e.collectServicesVirtualServerServiceHits(services, target)
		e.servicesVirtualServerServiceHits.Collect(ch)
		e.collectServicesActiveTransactions(services, target)
		e.servicesActiveTransactions.Collect(ch)
	})

	// 6. GSLB Services
	run("gslb_services", func() {
		gslbServices, err := netscaler.GetGSLBServiceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get GSLB service stats", "target", target.URL, "err", err)
			return
		}
		e.collectGSLBServicesState(gslbServices, target)
		e.gslbServicesState.Collect(ch)
		e.collectGSLBServicesTotalRequests(gslbServices, target)
		e.gslbServicesTotalRequests.Collect(ch)
		e.collectGSLBServicesTotalResponses(gslbServices, target)
		e.gslbServicesTotalResponses.Collect(ch)
		e.collectGSLBServicesTotalRequestBytes(gslbServices, target)
		e.gslbServicesTotalRequestBytes.Collect(ch)
		e.collectGSLBServicesTotalResponseBytes(gslbServices, target)
		e.gslbServicesTotalResponseBytes.Collect(ch)
		e.collectGSLBServicesCurrentClientConns(gslbServices, target)
		e.gslbServicesCurrentClientConns.Collect(ch)
		e.collectGSLBServicesCurrentServerConns(gslbServices, target)
		e.gslbServicesCurrentServerConns.Collect(ch)
		e.collectGSLBServicesEstablishedConnections(gslbServices, target)
		e.gslbServicesEstablishedConnections.Collect(ch)
		e.collectGSLBServicesCurrentLoad(gslbServices, target)
		e.gslbServicesCurrentLoad.Collect(ch)
		e.collectGSLBServicesVirtualServerServiceHits(gslbServices, target)
		e.gslbServicesVirtualServerServiceHits.Collect(ch)
	})

	// 7. GSLB Virtual Servers
	run("gslb_vservers", func() {
		gslbVirtualServers, err := netscaler.GetGSLBVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get GSLB virtual server stats", "target", target.URL, "err", err)
			return
		}
		e.collectGSLBVirtualServerState(gslbVirtualServers, target)
		e.gslbVirtualServersState.Collect(ch)
		e.collectGSLBVirtualServerHealth(gslbVirtualServers, target)
		e.gslbVirtualServersHealth.Collect(ch)
		e.collectGSLBVirtualServerInactiveServices(gslbVirtualServers, target)
		e.gslbVirtualServersInactiveServices.Collect(ch)
		e.collectGSLBVirtualServerActiveServices(gslbVirtualServers, target)
		e.gslbVirtualServersActiveServices.Collect(ch)
		e.collectGSLBVirtualServerTotalHits(gslbVirtualServers, target)
		e.gslbVirtualServersTotalHits.Collect(ch)
		e.collectGSLBVirtualServerTotalRequests(gslbVirtualServers, target)
		e.gslbVirtualServersTotalRequests.Collect(ch)
		e.collectGSLBVirtualServerTotalResponses(gslbVirtualServers, target)
		e.gslbVirtualServersTotalResponses.Collect(ch)
		e.collectGSLBVirtualServerTotalRequestBytes(gslbVirtualServers, target)
		e.gslbVirtualServersTotalRequestBytes.Collect(ch)
		e.collectGSLBVirtualServerTotalResponseBytes(gslbVirtualServers, target)
		e.gslbVirtualServersTotalResponseBytes.Collect(ch)
		e.collectGSLBVirtualServerCurrentClientConnections(gslbVirtualServers, target)
		e.gslbVirtualServersCurrentClientConnections.Collect(ch)
		e.collectGSLBVirtualServerCurrentServerConnections(gslbVirtualServers, target)
		e.gslbVirtualServersCurrentServerConnections.Collect(ch)
	})

	// 8. CS Virtual Servers
	run("cs_vservers", func() {
		csVirtualServers, err := netscaler.GetCSVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get CS virtual server stats", "target", target.URL, "err", err)
			return
		}
		e.collectCSVirtualServerState(csVirtualServers, target)
		e.csVirtualServersState.Collect(ch)
		e.collectCSVirtualServerTotalHits(csVirtualServers, target)
		e.csVirtualServersTotalHits.Collect(ch)
		e.collectCSVirtualServerTotalRequests(csVirtualServers, target)
		e.csVirtualServersTotalRequests.Collect(ch)
		e.collectCSVirtualServerTotalResponses(csVirtualServers, target)
		e.csVirtualServersTotalResponses.Collect(ch)
		e.collectCSVirtualServerTotalRequestBytes(csVirtualServers, target)
		e.csVirtualServersTotalRequestBytes.Collect(ch)
		e.collectCSVirtualServerTotalResponseBytes(csVirtualServers, target)
		e.csVirtualServersTotalResponseBytes.Collect(ch)
		e.collectCSVirtualServerCurrentClientConnections(csVirtualServers, target)
		e.csVirtualServersCurrentClientConnections.Collect(ch)
		e.collectCSVirtualServerCurrentServerConnections(csVirtualServers, target)
		e.csVirtualServersCurrentServerConnections.Collect(ch)
		e.collectCSVirtualServerEstablishedConnections(csVirtualServers, target)
		e.csVirtualServersEstablishedConnections.Collect(ch)
		e.collectCSVirtualServerTotalPacketsReceived(csVirtualServers, target)
		e.csVirtualServersTotalPacketsReceived.Collect(ch)
		e.collectCSVirtualServerTotalPacketsSent(csVirtualServers, target)
		e.csVirtualServersTotalPacketsSent.Collect(ch)
		e.collectCSVirtualServerTotalSpillovers(csVirtualServers, target)
		e.csVirtualServersTotalSpillovers.Collect(ch)
		e.collectCSVirtualServerDeferredRequests(csVirtualServers, target)
		e.csVirtualServersDeferredRequests.Collect(ch)
		e.collectCSVirtualServerNumberInvalidRequestResponse(csVirtualServers, target)
		e.csVirtualServersNumberInvalidRequestResponse.Collect(ch)
		e.collectCSVirtualServerNumberInvalidRequestResponseDropped(csVirtualServers, target)
		e.csVirtualServersNumberInvalidRequestResponseDropped.Collect(ch)
		e.collectCSVirtualServerTotalVServerDownBackupHits(csVirtualServers, target)
		e.csVirtualServersTotalVServerDownBackupHits.Collect(ch)
		e.collectCSVirtualServerCurrentMultipathSessions(csVirtualServers, target)
		e.csVirtualServersCurrentMultipathSessions.Collect(ch)
		e.collectCSVirtualServerCurrentMultipathSubflows(csVirtualServers, target)
		e.csVirtualServersCurrentMultipathSubflows.Collect(ch)
	})

	// 9. VPN Virtual Servers
	run("vpn_vservers", func() {
		vpnVirtualServers, err := netscaler.GetVPNVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get VPN virtual server stats", "target", target.URL, "err", err)
			return
		}
		e.collectVPNVirtualServerTotalRequests(vpnVirtualServers, target)
		e.vpnVirtualServersTotalRequests.Collect(ch)
		e.collectVPNVirtualServerTotalResponses(vpnVirtualServers, target)
		e.vpnVirtualServersTotalResponses.Collect(ch)
		e.collectVPNVirtualServerTotalRequestBytes(vpnVirtualServers, target)
		e.vpnVirtualServersTotalRequestBytes.Collect(ch)
		e.collectVPNVirtualServerTotalResponseBytes(vpnVirtualServers, target)
		e.vpnVirtualServersTotalResponseBytes.Collect(ch)
		e.collectVPNVirtualServerState(vpnVirtualServers, target)
		e.vpnVirtualServersState.Collect(ch)
	})

	// 10. AAA Stats
	run("aaa_stats", func() {
		aaa, err := netscaler.GetAAAStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get AAA stats", "target", target.URL, "err", err)
			return
		}
		e.collectAaaAuthSuccess(aaa, target)
		e.aaaAuthSuccess.Collect(ch)
		e.collectAaaAuthFail(aaa, target)
		e.aaaAuthFail.Collect(ch)
		e.collectAaaAuthOnlyHTTPSuccess(aaa, target)
		e.aaaAuthOnlyHTTPSuccess.Collect(ch)
		e.collectAaaAuthOnlyHTTPFail(aaa, target)
		e.aaaAuthOnlyHTTPFail.Collect(ch)
		e.collectAaaCurIcaSessions(aaa, target)
		e.aaaCurIcaSessions.Collect(ch)
		e.collectAaaCurIcaOnlyConn(aaa, target)
		e.aaaCurIcaOnlyConn.Collect(ch)
	})

	// 11. Service Groups (Nested parallelization)
	run("service_groups", func() {
		servicegroups, err := netscaler.GetServiceGroups(ctx, nsClient, "attrs=servicegroupname")
		if err != nil {
			e.logger.Error("failed to get service groups", "target", target.URL, "err", err)
			return
		}

		for _, sg := range servicegroups.ServiceGroups {
			sgName := sg.Name // Capture for closure
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-ctx.Done():
					return
				}

				stats, err2 := netscaler.GetServiceGroupMemberStats(ctx, nsClient, sgName)
				if err2 != nil {
					e.logger.Error("failed to get service group member stats", "service_group", sgName, "target", target.URL, "err", err2)
					return
				}

				if len(stats.ServiceGroups) == 0 || len(stats.ServiceGroups[0].ServiceGroupMembers) == 0 {
					return
				}

				for _, s := range stats.ServiceGroups[0].ServiceGroupMembers {
					// Extract member name from ServiceGroupName (format: "sgname?servername" or use IP)
					memberName := s.PrimaryIPAddress
					if parts := strings.Split(s.ServiceGroupName, "?"); len(parts) > 1 {
						memberName = parts[1]
					}

					e.collectServiceGroupsState(s, sgName, memberName, target)
					e.serviceGroupsState.Collect(ch)
					e.collectServiceGroupsAvgTTFB(s, sgName, memberName, target)
					e.serviceGroupsAvgTTFB.Collect(ch)
					e.collectServiceGroupsTotalRequests(s, sgName, memberName, target)
					e.serviceGroupsTotalRequests.Collect(ch)
					e.collectServiceGroupsTotalResponses(s, sgName, memberName, target)
					e.serviceGroupsTotalResponses.Collect(ch)
					e.collectServiceGroupsTotalRequestBytes(s, sgName, memberName, target)
					e.serviceGroupsTotalRequestBytes.Collect(ch)
					e.collectServiceGroupsTotalResponseBytes(s, sgName, memberName, target)
					e.serviceGroupsTotalResponseBytes.Collect(ch)
					e.collectServiceGroupsCurrentClientConnections(s, sgName, memberName, target)
					e.serviceGroupsCurrentClientConnections.Collect(ch)
					e.collectServiceGroupsSurgeCount(s, sgName, memberName, target)
					e.serviceGroupsSurgeCount.Collect(ch)
					e.collectServiceGroupsCurrentServerConnections(s, sgName, memberName, target)
					e.serviceGroupsCurrentServerConnections.Collect(ch)
					e.collectServiceGroupsServerEstablishedConnections(s, sgName, memberName, target)
					e.serviceGroupsServerEstablishedConnections.Collect(ch)
					e.collectServiceGroupsCurrentReusePool(s, sgName, memberName, target)
					e.serviceGroupsCurrentReusePool.Collect(ch)
					e.collectServiceGroupsMaxClients(s, sgName, memberName, target)
					e.serviceGroupsMaxClients.Collect(ch)
				}
			}()
		}
	})

	// 12. Topology metrics (always collected)
	run("topology", func() {
		e.collectTopologyMetrics(ctx, nsClient, target, ch)
	})

	// 13. Protocol HTTP Stats
	run("protocol_http", func() {
		e.collectProtocolHTTPStats(ctx, nsClient, target, ch)
	})

	// 14. Protocol TCP Stats
	run("protocol_tcp", func() {
		e.collectProtocolTCPStats(ctx, nsClient, target, ch)
	})

	// 15. Protocol IP Stats
	run("protocol_ip", func() {
		e.collectProtocolIPStats(ctx, nsClient, target, ch)
	})

	// 16. SSL Stats
	run("ssl_stats", func() {
		e.collectSSLStats(ctx, nsClient, target, ch)
	})

	// 17. SSL Cert Keys
	run("ssl_certs", func() {
		e.collectSSLCertKeys(ctx, nsClient, target, ch)
	})

	// 18. SSL VServer Stats
	run("ssl_vservers", func() {
		e.collectSSLVServerStats(ctx, nsClient, target, ch)
	})

	// 19. System CPU per-core Stats
	run("system_cpu", func() {
		e.collectSystemCPUStats(ctx, nsClient, target, ch)
	})

	// 20. Bandwidth Capacity Stats
	run("ns_capacity", func() {
		e.collectNSCapacityStats(ctx, nsClient, target, ch)
	})

	wg.Wait()

	// Set probe success to 1 (we completed successfully)
	e.probeSuccess.Reset()
	e.probeSuccess.WithLabelValues(baseLabels...).Set(1)
	e.probeSuccess.Collect(ch)
}
