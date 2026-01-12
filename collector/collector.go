package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.targetType == "mps" {
		e.scrapeMPS(ch)
	} else {
		e.scrapeADC(ch)
	}
}

// scrapeADC scrapes the NetScaler ADC instance
func (e *Exporter) scrapeADC(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nsClient, err := netscaler.NewNitroClient(e.url, e.username, e.password, e.ignoreCert, e.caFile)
	if err != nil {
		e.logger.Error("failed to create Nitro client", "url", e.url, "err", err)
		return
	}
	defer nsClient.CloseIdleConnections()

	var wg sync.WaitGroup
	// Semaphore to limit concurrent requests to avoid overloading the NetScaler
	sem := make(chan struct{}, e.parallelism)

	// Helper to run a scrape function concurrently
	run := func(name string, scrapeFn func()) {
		if e.config.IsModuleDisabled(name) {
			return // Skip disabled modules
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}: // Acquire token
				defer func() { <-sem }() // Release token
				scrapeFn()
			case <-ctx.Done():
				e.logger.Warn("context cancelled, skipping scrape", "url", e.url, "name", name)
			}
		}()
	}

	// Build base label values
	baseLabels := e.buildLabelValues()

	// Collect topology metrics FIRST (synchronously) to populate chainMembership
	// This must complete before service_groups runs so it can use chain labels
	if !e.config.IsModuleDisabled("topology") {
		e.collectTopologyMetrics(ctx, nsClient, ch)
	}

	// 1. NS Stats
	run("ns_stats", func() {
		ns, err := netscaler.GetNSStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get NS stats", "url", e.url, "err", err)
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
			e.logger.Error("failed to get NS license", "url", e.url, "err", err)
			return
		}
		fltModelID, _ := strconv.ParseFloat(nslicense.NSLicense.ModelID, 64)
		ch <- prometheus.MustNewConstMetric(e.modelID, prometheus.GaugeValue, fltModelID, baseLabels...)
	})

	// 3. Interfaces
	run("interfaces", func() {
		interfaces, err := netscaler.GetInterfaceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get interface stats", "url", e.url, "err", err)
			return
		}
		e.collectInterfacesRxBytes(interfaces)
		e.interfacesRxBytes.Collect(ch)
		e.collectInterfacesTxBytes(interfaces)
		e.interfacesTxBytes.Collect(ch)
		e.collectInterfacesRxPackets(interfaces)
		e.interfacesRxPackets.Collect(ch)
		e.collectInterfacesTxPackets(interfaces)
		e.interfacesTxPackets.Collect(ch)
		e.collectInterfacesJumboPacketsRx(interfaces)
		e.interfacesJumboPacketsRx.Collect(ch)
		e.collectInterfacesJumboPacketsTx(interfaces)
		e.interfacesJumboPacketsTx.Collect(ch)
		e.collectInterfacesErrorPacketsRx(interfaces)
		e.interfacesErrorPacketsRx.Collect(ch)
	})

	// 4. Virtual Servers
	run("virtual_servers", func() {
		virtualServers, err := netscaler.GetVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get virtual server stats", "url", e.url, "err", err)
			return
		}
		e.collectVirtualServerState(virtualServers)
		e.virtualServersState.Collect(ch)
		e.collectVirtualServerWaitingRequests(virtualServers)
		e.virtualServersWaitingRequests.Collect(ch)
		e.collectVirtualServerHealth(virtualServers)
		e.virtualServersHealth.Collect(ch)
		e.collectVirtualServerInactiveServices(virtualServers)
		e.virtualServersInactiveServices.Collect(ch)
		e.collectVirtualServerActiveServices(virtualServers)
		e.virtualServersActiveServices.Collect(ch)
		e.collectVirtualServerTotalHits(virtualServers)
		e.virtualServersTotalHits.Collect(ch)
		e.collectVirtualServerTotalRequests(virtualServers)
		e.virtualServersTotalRequests.Collect(ch)
		e.collectVirtualServerTotalResponses(virtualServers)
		e.virtualServersTotalResponses.Collect(ch)
		e.collectVirtualServerTotalRequestBytes(virtualServers)
		e.virtualServersTotalRequestBytes.Collect(ch)
		e.collectVirtualServerTotalResponseBytes(virtualServers)
		e.virtualServersTotalResponseBytes.Collect(ch)
		e.collectVirtualServerCurrentClientConnections(virtualServers)
		e.virtualServersCurrentClientConnections.Collect(ch)
		e.collectVirtualServerCurrentServerConnections(virtualServers)
		e.virtualServersCurrentServerConnections.Collect(ch)
	})

	// 5. Services
	run("services", func() {
		services, err := netscaler.GetServiceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get service stats", "url", e.url, "err", err)
			return
		}
		e.collectServicesThroughput(services)
		e.servicesThroughput.Collect(ch)
		e.collectServicesAvgTTFB(services)
		e.servicesAvgTTFB.Collect(ch)
		e.collectServicesState(services)
		e.servicesState.Collect(ch)
		e.collectServicesTotalRequests(services)
		e.servicesTotalRequests.Collect(ch)
		e.collectServicesTotalResponses(services)
		e.servicesTotalResponses.Collect(ch)
		e.collectServicesTotalRequestBytes(services)
		e.servicesTotalRequestBytes.Collect(ch)
		e.collectServicesTotalResponseBytes(services)
		e.servicesTotalResponseBytes.Collect(ch)
		e.collectServicesCurrentClientConns(services)
		e.servicesCurrentClientConns.Collect(ch)
		e.collectServicesSurgeCount(services)
		e.servicesSurgeCount.Collect(ch)
		e.collectServicesCurrentServerConns(services)
		e.servicesCurrentServerConns.Collect(ch)
		e.collectServicesServerEstablishedConnections(services)
		e.servicesServerEstablishedConnections.Collect(ch)
		e.collectServicesCurrentReusePool(services)
		e.servicesCurrentReusePool.Collect(ch)
		e.collectServicesMaxClients(services)
		e.servicesMaxClients.Collect(ch)
		e.collectServicesCurrentLoad(services)
		e.servicesCurrentLoad.Collect(ch)
		e.collectServicesVirtualServerServiceHits(services)
		e.servicesVirtualServerServiceHits.Collect(ch)
		e.collectServicesActiveTransactions(services)
		e.servicesActiveTransactions.Collect(ch)
	})

	// 6. GSLB Services
	run("gslb_services", func() {
		gslbServices, err := netscaler.GetGSLBServiceStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get GSLB service stats", "url", e.url, "err", err)
			return
		}
		e.collectGSLBServicesState(gslbServices)
		e.gslbServicesState.Collect(ch)
		e.collectGSLBServicesTotalRequests(gslbServices)
		e.gslbServicesTotalRequests.Collect(ch)
		e.collectGSLBServicesTotalResponses(gslbServices)
		e.gslbServicesTotalResponses.Collect(ch)
		e.collectGSLBServicesTotalRequestBytes(gslbServices)
		e.gslbServicesTotalRequestBytes.Collect(ch)
		e.collectGSLBServicesTotalResponseBytes(gslbServices)
		e.gslbServicesTotalResponseBytes.Collect(ch)
		e.collectGSLBServicesCurrentClientConns(gslbServices)
		e.gslbServicesCurrentClientConns.Collect(ch)
		e.collectGSLBServicesCurrentServerConns(gslbServices)
		e.gslbServicesCurrentServerConns.Collect(ch)
		e.collectGSLBServicesEstablishedConnections(gslbServices)
		e.gslbServicesEstablishedConnections.Collect(ch)
		e.collectGSLBServicesCurrentLoad(gslbServices)
		e.gslbServicesCurrentLoad.Collect(ch)
		e.collectGSLBServicesVirtualServerServiceHits(gslbServices)
		e.gslbServicesVirtualServerServiceHits.Collect(ch)
	})

	// 7. GSLB Virtual Servers
	run("gslb_vservers", func() {
		gslbVirtualServers, err := netscaler.GetGSLBVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get GSLB virtual server stats", "url", e.url, "err", err)
			return
		}
		e.collectGSLBVirtualServerState(gslbVirtualServers)
		e.gslbVirtualServersState.Collect(ch)
		e.collectGSLBVirtualServerHealth(gslbVirtualServers)
		e.gslbVirtualServersHealth.Collect(ch)
		e.collectGSLBVirtualServerInactiveServices(gslbVirtualServers)
		e.gslbVirtualServersInactiveServices.Collect(ch)
		e.collectGSLBVirtualServerActiveServices(gslbVirtualServers)
		e.gslbVirtualServersActiveServices.Collect(ch)
		e.collectGSLBVirtualServerTotalHits(gslbVirtualServers)
		e.gslbVirtualServersTotalHits.Collect(ch)
		e.collectGSLBVirtualServerTotalRequests(gslbVirtualServers)
		e.gslbVirtualServersTotalRequests.Collect(ch)
		e.collectGSLBVirtualServerTotalResponses(gslbVirtualServers)
		e.gslbVirtualServersTotalResponses.Collect(ch)
		e.collectGSLBVirtualServerTotalRequestBytes(gslbVirtualServers)
		e.gslbVirtualServersTotalRequestBytes.Collect(ch)
		e.collectGSLBVirtualServerTotalResponseBytes(gslbVirtualServers)
		e.gslbVirtualServersTotalResponseBytes.Collect(ch)
		e.collectGSLBVirtualServerCurrentClientConnections(gslbVirtualServers)
		e.gslbVirtualServersCurrentClientConnections.Collect(ch)
		e.collectGSLBVirtualServerCurrentServerConnections(gslbVirtualServers)
		e.gslbVirtualServersCurrentServerConnections.Collect(ch)
	})

	// 8. CS Virtual Servers
	run("cs_vservers", func() {
		csVirtualServers, err := netscaler.GetCSVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get CS virtual server stats", "url", e.url, "err", err)
			return
		}
		e.collectCSVirtualServerState(csVirtualServers)
		e.csVirtualServersState.Collect(ch)
		e.collectCSVirtualServerTotalHits(csVirtualServers)
		e.csVirtualServersTotalHits.Collect(ch)
		e.collectCSVirtualServerTotalRequests(csVirtualServers)
		e.csVirtualServersTotalRequests.Collect(ch)
		e.collectCSVirtualServerTotalResponses(csVirtualServers)
		e.csVirtualServersTotalResponses.Collect(ch)
		e.collectCSVirtualServerTotalRequestBytes(csVirtualServers)
		e.csVirtualServersTotalRequestBytes.Collect(ch)
		e.collectCSVirtualServerTotalResponseBytes(csVirtualServers)
		e.csVirtualServersTotalResponseBytes.Collect(ch)
		e.collectCSVirtualServerCurrentClientConnections(csVirtualServers)
		e.csVirtualServersCurrentClientConnections.Collect(ch)
		e.collectCSVirtualServerCurrentServerConnections(csVirtualServers)
		e.csVirtualServersCurrentServerConnections.Collect(ch)
		e.collectCSVirtualServerEstablishedConnections(csVirtualServers)
		e.csVirtualServersEstablishedConnections.Collect(ch)
		e.collectCSVirtualServerTotalPacketsReceived(csVirtualServers)
		e.csVirtualServersTotalPacketsReceived.Collect(ch)
		e.collectCSVirtualServerTotalPacketsSent(csVirtualServers)
		e.csVirtualServersTotalPacketsSent.Collect(ch)
		e.collectCSVirtualServerTotalSpillovers(csVirtualServers)
		e.csVirtualServersTotalSpillovers.Collect(ch)
		e.collectCSVirtualServerDeferredRequests(csVirtualServers)
		e.csVirtualServersDeferredRequests.Collect(ch)
		e.collectCSVirtualServerNumberInvalidRequestResponse(csVirtualServers)
		e.csVirtualServersNumberInvalidRequestResponse.Collect(ch)
		e.collectCSVirtualServerNumberInvalidRequestResponseDropped(csVirtualServers)
		e.csVirtualServersNumberInvalidRequestResponseDropped.Collect(ch)
		e.collectCSVirtualServerTotalVServerDownBackupHits(csVirtualServers)
		e.csVirtualServersTotalVServerDownBackupHits.Collect(ch)
		e.collectCSVirtualServerCurrentMultipathSessions(csVirtualServers)
		e.csVirtualServersCurrentMultipathSessions.Collect(ch)
		e.collectCSVirtualServerCurrentMultipathSubflows(csVirtualServers)
		e.csVirtualServersCurrentMultipathSubflows.Collect(ch)
	})

	// 9. VPN Virtual Servers
	run("vpn_vservers", func() {
		vpnVirtualServers, err := netscaler.GetVPNVirtualServerStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get VPN virtual server stats", "url", e.url, "err", err)
			return
		}
		e.collectVPNVirtualServerTotalRequests(vpnVirtualServers)
		e.vpnVirtualServersTotalRequests.Collect(ch)
		e.collectVPNVirtualServerTotalResponses(vpnVirtualServers)
		e.vpnVirtualServersTotalResponses.Collect(ch)
		e.collectVPNVirtualServerTotalRequestBytes(vpnVirtualServers)
		e.vpnVirtualServersTotalRequestBytes.Collect(ch)
		e.collectVPNVirtualServerTotalResponseBytes(vpnVirtualServers)
		e.vpnVirtualServersTotalResponseBytes.Collect(ch)
		e.collectVPNVirtualServerState(vpnVirtualServers)
		e.vpnVirtualServersState.Collect(ch)
	})

	// 10. AAA Stats
	run("aaa_stats", func() {
		aaa, err := netscaler.GetAAAStats(ctx, nsClient, "")
		if err != nil {
			e.logger.Error("failed to get AAA stats", "url", e.url, "err", err)
			return
		}
		e.collectAaaAuthSuccess(aaa)
		e.aaaAuthSuccess.Collect(ch)
		e.collectAaaAuthFail(aaa)
		e.aaaAuthFail.Collect(ch)
		e.collectAaaAuthOnlyHTTPSuccess(aaa)
		e.aaaAuthOnlyHTTPSuccess.Collect(ch)
		e.collectAaaAuthOnlyHTTPFail(aaa)
		e.aaaAuthOnlyHTTPFail.Collect(ch)
		e.collectAaaCurIcaSessions(aaa)
		e.aaaCurIcaSessions.Collect(ch)
		e.collectAaaCurIcaOnlyConn(aaa)
		e.aaaCurIcaOnlyConn.Collect(ch)
	})

	// 11. Service Groups (Nested parallelization)
	run("service_groups", func() {
		servicegroups, err := netscaler.GetServiceGroups(ctx, nsClient, "attrs=servicegroupname")
		if err != nil {
			e.logger.Error("failed to get service groups", "url", e.url, "err", err)
			return
		}

		// Reset all servicegroup metrics once before processing
		e.serviceGroupsState.Reset()
		e.serviceGroupsAvgTTFB.Reset()
		e.serviceGroupsTotalRequests.Reset()
		e.serviceGroupsTotalResponses.Reset()
		e.serviceGroupsTotalRequestBytes.Reset()
		e.serviceGroupsTotalResponseBytes.Reset()
		e.serviceGroupsCurrentClientConnections.Reset()
		e.serviceGroupsSurgeCount.Reset()
		e.serviceGroupsCurrentServerConnections.Reset()
		e.serviceGroupsServerEstablishedConnections.Reset()
		e.serviceGroupsCurrentReusePool.Reset()
		e.serviceGroupsMaxClients.Reset()

		// Use a separate WaitGroup for service group goroutines
		var sgWg sync.WaitGroup

		// Track seen members globally across all goroutines to prevent duplicates
		var seenMu sync.Mutex
		seenMembers := make(map[string]bool)

		// Deduplicate service groups (API may return duplicates)
		seenServiceGroups := make(map[string]bool)
		for _, sg := range servicegroups.ServiceGroups {
			sgName := sg.Name // Capture for closure
			if seenServiceGroups[sgName] {
				continue
			}
			seenServiceGroups[sgName] = true
			sgWg.Add(1)
			go func() {
				defer sgWg.Done()
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-ctx.Done():
					return
				}

				// Create servicegroup topology node (if topology enabled)
				var sgChain string
				if !e.config.IsModuleDisabled("topology") {
					nodeID := "servicegroup:" + sgName
					sgChain = e.chainMembership[nodeID]
					nodeLabels := e.buildLabelValues(nodeID, sgName, "servicegroup", "UP", sgChain)
					e.topologyNode.WithLabelValues(nodeLabels...).Set(1.0)
				}

				stats, err2 := netscaler.GetServiceGroupMemberStats(ctx, nsClient, sgName)
				if err2 != nil {
					e.logger.Error("failed to get service group member stats", "service_group", sgName, "url", e.url, "err", err2)
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

					// Deduplicate members globally (API may return duplicates)
					key := fmt.Sprintf("%s:%s:%d", sgName, memberName, s.PrimaryPort)
					seenMu.Lock()
					if seenMembers[key] {
						seenMu.Unlock()
						continue
					}
					seenMembers[key] = true
					seenMu.Unlock()

					// Set metric values (no Reset, no Collect - done once after all goroutines)
					port := strconv.Itoa(s.PrimaryPort)
					labels := e.buildLabelValues(sgName, memberName, port)

					state := 0.0
					if s.State == "UP" {
						state = 1.0
					}
					e.serviceGroupsState.WithLabelValues(labels...).Set(state)

					if val, err := strconv.ParseFloat(s.AvgTimeToFirstByte, 64); err == nil {
						e.serviceGroupsAvgTTFB.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.TotalRequests, 64); err == nil {
						e.serviceGroupsTotalRequests.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.TotalResponses, 64); err == nil {
						e.serviceGroupsTotalResponses.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.TotalRequestBytes, 64); err == nil {
						e.serviceGroupsTotalRequestBytes.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.TotalResponseBytes, 64); err == nil {
						e.serviceGroupsTotalResponseBytes.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.CurrentClientConnections, 64); err == nil {
						e.serviceGroupsCurrentClientConnections.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.SurgeCount, 64); err == nil {
						e.serviceGroupsSurgeCount.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.CurrentServerConnections, 64); err == nil {
						e.serviceGroupsCurrentServerConnections.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.ServerEstablishedConnections, 64); err == nil {
						e.serviceGroupsServerEstablishedConnections.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.CurrentReusePool, 64); err == nil {
						e.serviceGroupsCurrentReusePool.WithLabelValues(labels...).Set(val)
					}
					if val, err := strconv.ParseFloat(s.MaxClients, 64); err == nil {
						e.serviceGroupsMaxClients.WithLabelValues(labels...).Set(val)
					}

					// Create topology server node and edge (reusing already-fetched data)
					if !e.config.IsModuleDisabled("topology") {
						serverID := fmt.Sprintf("server:%s:%d", s.PrimaryIPAddress, s.PrimaryPort)
						// Use server name (memberName) for title if available, otherwise fall back to IP
						serverTitle := fmt.Sprintf("%s:%d", memberName, s.PrimaryPort)
						topoState := "DOWN"
						value := 0.0
						if s.State == "UP" {
							topoState = "UP"
							value = 1.0
						}
						// Server inherits chain from its parent servicegroup
						nodeLabels := e.buildLabelValues(serverID, serverTitle, "server", topoState, sgChain)
						e.topologyNode.WithLabelValues(nodeLabels...).Set(value)

						edgeID := fmt.Sprintf("servicegroup:%s->server:%s:%d", sgName, s.PrimaryIPAddress, s.PrimaryPort)
						sourceID := "servicegroup:" + sgName
						edgeLabels := e.buildLabelValues(edgeID, sourceID, serverID, "1", "", sgChain)
						e.topologyEdge.WithLabelValues(edgeLabels...).Set(1)
					}
				}
			}()
		}

		// Wait for all service group goroutines to complete
		sgWg.Wait()

		// Collect all servicegroup metrics once after all data is set
		e.serviceGroupsState.Collect(ch)
		e.serviceGroupsAvgTTFB.Collect(ch)
		e.serviceGroupsTotalRequests.Collect(ch)
		e.serviceGroupsTotalResponses.Collect(ch)
		e.serviceGroupsTotalRequestBytes.Collect(ch)
		e.serviceGroupsTotalResponseBytes.Collect(ch)
		e.serviceGroupsCurrentClientConnections.Collect(ch)
		e.serviceGroupsSurgeCount.Collect(ch)
		e.serviceGroupsCurrentServerConnections.Collect(ch)
		e.serviceGroupsServerEstablishedConnections.Collect(ch)
		e.serviceGroupsCurrentReusePool.Collect(ch)
		e.serviceGroupsMaxClients.Collect(ch)
	})

	// 12. Protocol HTTP Stats
	run("protocol_http", func() {
		e.collectProtocolHTTPStats(ctx, nsClient, ch)
	})

	// 14. Protocol TCP Stats
	run("protocol_tcp", func() {
		e.collectProtocolTCPStats(ctx, nsClient, ch)
	})

	// 15. Protocol IP Stats
	run("protocol_ip", func() {
		e.collectProtocolIPStats(ctx, nsClient, ch)
	})

	// 16. SSL Stats
	run("ssl_stats", func() {
		e.collectSSLStats(ctx, nsClient, ch)
	})

	// 17. SSL Cert Keys
	run("ssl_certs", func() {
		e.collectSSLCertKeys(ctx, nsClient, ch)
	})

	// 18. SSL VServer Stats
	run("ssl_vservers", func() {
		e.collectSSLVServerStats(ctx, nsClient, ch)
	})

	// 19. System CPU per-core Stats
	run("system_cpu", func() {
		e.collectSystemCPUStats(ctx, nsClient, ch)
	})

	// 20. Bandwidth Capacity Stats
	run("ns_capacity", func() {
		e.collectNSCapacityStats(ctx, nsClient, ch)
	})

	wg.Wait()

	// Collect topology metrics after all modules have added their nodes/edges
	if !e.config.IsModuleDisabled("topology") {
		e.topologyNode.Collect(ch)
		e.topologyEdge.Collect(ch)
	}
}

// scrapeMPS scrapes the Citrix ADM (MPS) instance
func (e *Exporter) scrapeMPS(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mpsClient, err := netscaler.NewMPSClient(e.url, e.username, e.password, e.ignoreCert, e.caFile)
	if err != nil {
		e.logger.Error("failed to create MPS client", "url", e.url, "err", err)
		return
	}
	defer mpsClient.CloseIdleConnections()

	// MPS Health stats
	mpsHealth, err := netscaler.GetMPSHealth(ctx, mpsClient)
	if err != nil {
		e.logger.Error("failed to get MPS health stats", "url", e.url, "err", err)
		return
	}

	e.collectMPSHealth(mpsHealth)
	e.mpsHealthCPUUsage.Collect(ch)
	e.mpsHealthDiskUsage.Collect(ch)
	e.mpsHealthDiskFree.Collect(ch)
	e.mpsHealthDiskTotal.Collect(ch)
	e.mpsHealthDiskUsed.Collect(ch)
	e.mpsHealthMemoryUsage.Collect(ch)
	e.mpsHealthMemoryFree.Collect(ch)
	e.mpsHealthMemoryTotal.Collect(ch)
}
