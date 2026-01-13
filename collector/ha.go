package collector

import (
	"context"
	"strconv"
	"strings"

	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// collectHAStats collects HA (High Availability) metrics from both config and stat endpoints
func (e *Exporter) collectHAStats(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) {
	baseLabels := e.buildLabelValues()

	// Reset GaugeVec metrics
	e.haNodeState.Reset()
	e.haNodeStatus.Reset()
	e.haNodeSyncState.Reset()
	e.haNodeMasterStateSeconds.Reset()

	// Fetch HA node config (per-node info)
	haConfig, err := netscaler.GetHANodeConfig(ctx, nsClient)
	if err != nil {
		e.logger.Error("failed to get HA node config", "url", e.url, "err", err)
	} else {
		for _, node := range haConfig.HANodes {
			labels := e.buildLabelValues(node.ID, node.Name, node.IPAddress)

			// State: 1=Primary, 0=Secondary
			state := 0.0
			if strings.EqualFold(node.State, "Primary") {
				state = 1.0
			}
			e.haNodeState.WithLabelValues(labels...).Set(state)

			// Status: 1=UP, 0=DOWN
			status := 0.0
			if strings.EqualFold(node.HAStatus, "UP") {
				status = 1.0
			}
			e.haNodeStatus.WithLabelValues(labels...).Set(status)

			// Sync state: 1=SUCCESS or ENABLED, 0=other
			syncState := 0.0
			if strings.EqualFold(node.HASync, "SUCCESS") || strings.EqualFold(node.HASync, "ENABLED") {
				syncState = 1.0
			}
			e.haNodeSyncState.WithLabelValues(labels...).Set(syncState)

			// Master state time (seconds)
			e.haNodeMasterStateSeconds.WithLabelValues(labels...).Set(float64(node.MasterStateTime))
		}
	}

	// Collect GaugeVec metrics
	e.haNodeState.Collect(ch)
	e.haNodeStatus.Collect(ch)
	e.haNodeSyncState.Collect(ch)
	e.haNodeMasterStateSeconds.Collect(ch)

	// Fetch HA node stats (global stats)
	haStats, err := netscaler.GetHANodeStats(ctx, nsClient)
	if err != nil {
		e.logger.Error("failed to get HA node stats", "url", e.url, "err", err)
		return
	}

	// Current state: 1=UP, 0=DOWN
	curState := 0.0
	if strings.EqualFold(haStats.HANode.HACurState, "UP") {
		curState = 1.0
	}
	ch <- prometheus.MustNewConstMetric(e.haCurState, prometheus.GaugeValue, curState, baseLabels...)

	// Packets received total
	pktRx, _ := strconv.ParseFloat(haStats.HANode.HATotPktRx, 64)
	ch <- prometheus.MustNewConstMetric(e.haPacketsRxTotal, prometheus.CounterValue, pktRx, baseLabels...)

	// Packets transmitted total
	pktTx, _ := strconv.ParseFloat(haStats.HANode.HATotPktTx, 64)
	ch <- prometheus.MustNewConstMetric(e.haPacketsTxTotal, prometheus.CounterValue, pktTx, baseLabels...)

	// Sync failures total
	syncFailures, _ := strconv.ParseFloat(haStats.HANode.HAErrSyncFailure, 64)
	ch <- prometheus.MustNewConstMetric(e.haSyncFailuresTotal, prometheus.CounterValue, syncFailures, baseLabels...)

	// Propagation timeouts total
	propTimeouts, _ := strconv.ParseFloat(haStats.HANode.HAErrPropTimeout, 64)
	ch <- prometheus.MustNewConstMetric(e.haPropTimeoutsTotal, prometheus.CounterValue, propTimeouts, baseLabels...)
}
