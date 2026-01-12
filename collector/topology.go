package collector

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collectTopologyMetrics(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) {
	e.topologyNode.Reset()
	e.topologyEdge.Reset()

	// Fetch all bindings in parallel using bulk APIs (3 API calls instead of 152+)
	var allSvcBindings []netscaler.LBVServerServiceBinding
	var allSgBindings []netscaler.LBVServerServiceGroupBinding
	var allCsBindings []netscaler.CSVServerLBVServerBinding

	var bindingsWg sync.WaitGroup
	bindingsWg.Add(3)

	go func() {
		defer bindingsWg.Done()
		bindings, err := netscaler.GetAllLBVServerServiceBindings(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting bulk service bindings", "url", e.url, "err", err)
			return
		}
		allSvcBindings = bindings
	}()

	go func() {
		defer bindingsWg.Done()
		bindings, err := netscaler.GetAllLBVServerServiceGroupBindings(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting bulk servicegroup bindings", "url", e.url, "err", err)
			return
		}
		allSgBindings = bindings
	}()

	go func() {
		defer bindingsWg.Done()
		bindings, err := netscaler.GetAllCSVServerLBVServerBindings(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting bulk csvserver bindings", "url", e.url, "err", err)
			return
		}
		allCsBindings = bindings
	}()

	bindingsWg.Wait()

	// Build lookup maps by vserver name
	svcBindingsByVS := make(map[string][]netscaler.LBVServerServiceBinding)
	for _, b := range allSvcBindings {
		svcBindingsByVS[b.Name] = append(svcBindingsByVS[b.Name], b)
	}

	sgBindingsByVS := make(map[string][]netscaler.LBVServerServiceGroupBinding)
	for _, b := range allSgBindings {
		sgBindingsByVS[b.Name] = append(sgBindingsByVS[b.Name], b)
	}

	csBindingsByVS := make(map[string][]netscaler.CSVServerLBVServerBinding)
	for _, b := range allCsBindings {
		csBindingsByVS[b.Name] = append(csBindingsByVS[b.Name], b)
	}

	// Build chain membership map
	e.chainMembership = e.buildChainMembership(csBindingsByVS, svcBindingsByVS, sgBindingsByVS)

	// Collect LB Virtual Server nodes
	lbVServers, err := netscaler.GetVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting LB vserver stats for topology", "url", e.url, "err", err)
	} else {
		for _, vs := range lbVServers.VirtualServerStats {
			nodeID := "lbvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			chain := e.chainMembership[nodeID]
			labels := e.buildLabelValues(nodeID, vs.Name, "lbvserver", state, chain)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect CS Virtual Server nodes
	csVServers, err := netscaler.GetCSVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting CS vserver stats for topology", "url", e.url, "err", err)
	} else {
		for _, vs := range csVServers.CSVirtualServerStats {
			nodeID := "csvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			chain := e.chainMembership[nodeID]
			labels := e.buildLabelValues(nodeID, vs.Name, "csvserver", state, chain)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect Service nodes
	services, err := netscaler.GetServiceStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting service stats for topology", "url", e.url, "err", err)
	} else {
		for _, svc := range services.ServiceStats {
			nodeID := "service:" + svc.Name
			state := "DOWN"
			value := 0.0
			if svc.State == "UP" {
				state = "UP"
				value = 1.0
			}
			chain := e.chainMembership[nodeID]
			labels := e.buildLabelValues(nodeID, svc.Name, "service", state, chain)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Note: Service Group nodes are created in the service_groups collector
	// to avoid redundant API calls (collector.go already fetches servicegroups)

	// Collect LB VServer -> Service and Service Group edges using lookup maps
	if len(lbVServers.VirtualServerStats) > 0 {
		for _, vs := range lbVServers.VirtualServerStats {
			sourceID := "lbvserver:" + vs.Name
			sourceChain := e.chainMembership[sourceID]

			// Service bindings from lookup map
			for _, b := range svcBindingsByVS[vs.Name] {
				edgeID := "lbvserver:" + b.Name + "->service:" + b.ServiceName
				targetID := "service:" + b.ServiceName
				weight := b.Weight
				if weight == "" {
					weight = "1"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, weight, "", sourceChain)
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}

			// Service group bindings from lookup map
			for _, b := range sgBindingsByVS[vs.Name] {
				edgeID := "lbvserver:" + b.Name + "->servicegroup:" + b.ServiceGroupName
				targetID := "servicegroup:" + b.ServiceGroupName
				weight := b.Weight
				if weight == "" {
					weight = "1"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, weight, "", sourceChain)
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}
		}
	}

	// Collect CS VServer -> LB VServer edges using lookup map
	if len(csVServers.CSVirtualServerStats) > 0 {
		for _, vs := range csVServers.CSVirtualServerStats {
			sourceID := "csvserver:" + vs.Name
			sourceChain := e.chainMembership[sourceID]

			for _, b := range csBindingsByVS[vs.Name] {
				edgeID := "csvserver:" + b.Name + "->lbvserver:" + b.LBVServer
				targetID := "lbvserver:" + b.LBVServer
				priority := b.Priority
				if priority == "" {
					priority = "0"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, "", priority, sourceChain)
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}
		}
	}

	// Note: Collect() is called later in scrapeADC after service_groups has added its nodes/edges
}

// buildChainMembership computes which chain(s) each node belongs to.
// A chain is identified by its top-level frontend: csvserver name, or lbvserver name if standalone.
// Returns map[nodeID] → comma-separated chain names (sorted alphabetically).
func (e *Exporter) buildChainMembership(
	csBindingsByVS map[string][]netscaler.CSVServerLBVServerBinding,
	svcBindingsByVS map[string][]netscaler.LBVServerServiceBinding,
	sgBindingsByVS map[string][]netscaler.LBVServerServiceGroupBinding,
) map[string]string {
	// nodeChains collects all chains for each node
	nodeChains := make(map[string][]string)

	// addChain adds a chain to a node, avoiding duplicates
	addChain := func(nodeID, chain string) {
		for _, existing := range nodeChains[nodeID] {
			if existing == chain {
				return
			}
		}
		nodeChains[nodeID] = append(nodeChains[nodeID], chain)
	}

	// Build reverse lookup: lbvserver → csvservers that reference it
	lbToCsMap := make(map[string][]string)
	for csvName, bindings := range csBindingsByVS {
		for _, b := range bindings {
			lbToCsMap[b.LBVServer] = append(lbToCsMap[b.LBVServer], csvName)
		}
	}

	// Process csvservers: each csvserver is its own chain
	for csvName, bindings := range csBindingsByVS {
		chain := csvName
		csvNodeID := "csvserver:" + csvName
		addChain(csvNodeID, chain)

		// Traverse to lbvservers
		for _, b := range bindings {
			lbNodeID := "lbvserver:" + b.LBVServer
			addChain(lbNodeID, chain)

			// Traverse from lbvserver to services
			for _, svcB := range svcBindingsByVS[b.LBVServer] {
				svcNodeID := "service:" + svcB.ServiceName
				addChain(svcNodeID, chain)
			}

			// Traverse from lbvserver to servicegroups
			for _, sgB := range sgBindingsByVS[b.LBVServer] {
				sgNodeID := "servicegroup:" + sgB.ServiceGroupName
				addChain(sgNodeID, chain)
			}
		}
	}

	// Process standalone lbvservers (not behind any csvserver)
	for lbName := range svcBindingsByVS {
		if _, hasCsParent := lbToCsMap[lbName]; !hasCsParent {
			chain := lbName
			lbNodeID := "lbvserver:" + lbName
			addChain(lbNodeID, chain)

			// Traverse to services
			for _, svcB := range svcBindingsByVS[lbName] {
				svcNodeID := "service:" + svcB.ServiceName
				addChain(svcNodeID, chain)
			}
		}
	}
	for lbName := range sgBindingsByVS {
		if _, hasCsParent := lbToCsMap[lbName]; !hasCsParent {
			chain := lbName
			lbNodeID := "lbvserver:" + lbName
			addChain(lbNodeID, chain)

			// Traverse to servicegroups
			for _, sgB := range sgBindingsByVS[lbName] {
				sgNodeID := "servicegroup:" + sgB.ServiceGroupName
				addChain(sgNodeID, chain)
			}
		}
	}

	// Convert to comma-separated strings (sorted for consistency)
	result := make(map[string]string, len(nodeChains))
	for nodeID, chains := range nodeChains {
		sort.Strings(chains)
		result[nodeID] = strings.Join(chains, ",")
	}

	return result
}
