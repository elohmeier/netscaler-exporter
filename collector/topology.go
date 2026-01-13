package collector

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

// CSToLBMapping represents a resolved CS vserver → LB vserver relationship.
type CSToLBMapping struct {
	CSVServer   string
	LBVServer   string
	Priority    string
	PolicyName  string // For policy-based routing
}

func (e *Exporter) collectTopologyMetrics(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) {
	e.topologyNode.Reset()
	e.topologyEdge.Reset()
	e.topologyNodeState.Reset()
	e.topologyNodeHealth.Reset()
	e.topologyNodeRequestsTotal.Reset()
	e.topologyNodeConnections.Reset()
	e.topologyNodeTTFBMs.Reset()

	// Fetch all bindings in parallel using bulk APIs
	var allSvcBindings []netscaler.LBVServerServiceBinding
	var allSgBindings []netscaler.LBVServerServiceGroupBinding
	var allCsLbBindings []netscaler.CSVServerLBVServerBinding
	var allCsPolicyBindings []netscaler.CSVServerCSPolicyBinding
	var allCsPolicies []netscaler.CSPolicy
	var allCsActions []netscaler.CSAction

	var bindingsWg sync.WaitGroup
	bindingsWg.Add(6)

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
			e.logger.Debug("error getting bulk csvserver_lbvserver bindings", "url", e.url, "err", err)
			return
		}
		allCsLbBindings = bindings
	}()

	go func() {
		defer bindingsWg.Done()
		bindings, err := netscaler.GetAllCSVServerCSPolicyBindings(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting bulk csvserver_cspolicy bindings", "url", e.url, "err", err)
			return
		}
		allCsPolicyBindings = bindings
	}()

	go func() {
		defer bindingsWg.Done()
		policies, err := netscaler.GetAllCSPolicies(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting cspolicies", "url", e.url, "err", err)
			return
		}
		allCsPolicies = policies
	}()

	go func() {
		defer bindingsWg.Done()
		actions, err := netscaler.GetAllCSActions(ctx, nsClient)
		if err != nil {
			e.logger.Debug("error getting csactions", "url", e.url, "err", err)
			return
		}
		allCsActions = actions
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

	// Build policy → action lookup
	policyToAction := make(map[string]string)
	for _, p := range allCsPolicies {
		if p.Action != "" {
			policyToAction[p.PolicyName] = p.Action
		}
	}

	// Build action → targetlbvserver lookup
	actionToLB := make(map[string]string)
	for _, a := range allCsActions {
		if a.TargetLBVServer != "" {
			actionToLB[a.Name] = a.TargetLBVServer
		}
	}

	// Resolve all CS → LB mappings from multiple sources
	csToLBMappings := e.resolveCSToLBMappings(allCsLbBindings, allCsPolicyBindings, policyToAction, actionToLB)

	// Build csBindingsByVS for edge creation and chain membership
	csBindingsByVS := make(map[string][]CSToLBMapping)
	for _, m := range csToLBMappings {
		csBindingsByVS[m.CSVServer] = append(csBindingsByVS[m.CSVServer], m)
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

			// mainStat = health %, secondaryStat = connections
			mainStat := vs.Health
			secondaryStat := vs.CurrentClientConnections

			labels := e.buildLabelValues(nodeID, vs.Name, "lbvserver", state, chain, mainStat, secondaryStat)
			e.topologyNode.WithLabelValues(labels...).Set(value)

			// Emit topology node stats
			statsLabels := e.buildLabelValues(nodeID, "lbvserver", chain)
			e.topologyNodeState.WithLabelValues(statsLabels...).Set(value)

			if health, err := strconv.ParseFloat(vs.Health, 64); err == nil {
				e.topologyNodeHealth.WithLabelValues(statsLabels...).Set(health)
			}
			if requests, err := strconv.ParseFloat(vs.TotalRequests, 64); err == nil {
				e.topologyNodeRequestsTotal.WithLabelValues(statsLabels...).Set(requests)
			}
			if conns, err := strconv.ParseFloat(vs.CurrentClientConnections, 64); err == nil {
				e.topologyNodeConnections.WithLabelValues(statsLabels...).Set(conns)
			}
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

			// CS vservers: mainStat = "" (no health), secondaryStat = connections
			mainStat := ""
			secondaryStat := vs.CurrentClientConnections

			labels := e.buildLabelValues(nodeID, vs.Name, "csvserver", state, chain, mainStat, secondaryStat)
			e.topologyNode.WithLabelValues(labels...).Set(value)

			// Emit topology node stats (CS vservers use total_hits as main stat)
			statsLabels := e.buildLabelValues(nodeID, "csvserver", chain)
			e.topologyNodeState.WithLabelValues(statsLabels...).Set(value)

			if hits, err := strconv.ParseFloat(vs.TotalHits, 64); err == nil {
				e.topologyNodeRequestsTotal.WithLabelValues(statsLabels...).Set(hits)
			}
			if conns, err := strconv.ParseFloat(vs.CurrentClientConnections, 64); err == nil {
				e.topologyNodeConnections.WithLabelValues(statsLabels...).Set(conns)
			}
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

			// Services: mainStat = "", secondaryStat = ""
			labels := e.buildLabelValues(nodeID, svc.Name, "service", state, chain, "", "")
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

	// Collect CS VServer -> LB VServer edges using resolved mappings
	if len(csVServers.CSVirtualServerStats) > 0 {
		for _, vs := range csVServers.CSVirtualServerStats {
			sourceID := "csvserver:" + vs.Name
			sourceChain := e.chainMembership[sourceID]

			for _, m := range csBindingsByVS[vs.Name] {
				edgeID := "csvserver:" + m.CSVServer + "->lbvserver:" + m.LBVServer
				targetID := "lbvserver:" + m.LBVServer
				priority := m.Priority
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

// resolveCSToLBMappings resolves all CS vserver → LB vserver relationships from multiple sources:
// 1. Direct csvserver_lbvserver_binding entries
// 2. Policy bindings with direct targetlbvserver
// 3. Policy bindings resolved via policy → action → targetlbvserver
func (e *Exporter) resolveCSToLBMappings(
	directBindings []netscaler.CSVServerLBVServerBinding,
	policyBindings []netscaler.CSVServerCSPolicyBinding,
	policyToAction map[string]string,
	actionToLB map[string]string,
) []CSToLBMapping {
	var mappings []CSToLBMapping

	// Track seen mappings to avoid duplicates (CS+LB pair)
	seen := make(map[string]bool)
	addMapping := func(csv, lb, priority, policy string) {
		key := csv + ":" + lb
		if seen[key] {
			return
		}
		seen[key] = true
		mappings = append(mappings, CSToLBMapping{
			CSVServer:  csv,
			LBVServer:  lb,
			Priority:   priority,
			PolicyName: policy,
		})
	}

	// 1. Direct bindings (csvserver_lbvserver_binding)
	for _, b := range directBindings {
		addMapping(b.Name, b.LBVServer, b.Priority, "")
	}

	// 2. Policy bindings
	for _, pb := range policyBindings {
		var targetLB string

		// Check if binding has direct targetlbvserver
		if pb.TargetLBVServer != "" {
			targetLB = pb.TargetLBVServer
		} else {
			// Resolve via policy → action → targetlbvserver
			actionName := policyToAction[pb.PolicyName]
			if actionName != "" {
				targetLB = actionToLB[actionName]
			}
		}

		if targetLB != "" {
			addMapping(pb.Name, targetLB, pb.Priority, pb.PolicyName)
		}
	}

	return mappings
}

// buildChainMembership computes which chain(s) each node belongs to.
// A chain is identified by its top-level frontend: csvserver name, or lbvserver name if standalone.
// Returns map[nodeID] → comma-separated chain names (sorted alphabetically).
func (e *Exporter) buildChainMembership(
	csBindingsByVS map[string][]CSToLBMapping,
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
	for csvName, mappings := range csBindingsByVS {
		for _, m := range mappings {
			lbToCsMap[m.LBVServer] = append(lbToCsMap[m.LBVServer], csvName)
		}
	}

	// Process csvservers: each csvserver is its own chain
	for csvName, mappings := range csBindingsByVS {
		chain := csvName
		csvNodeID := "csvserver:" + csvName
		addChain(csvNodeID, chain)

		// Traverse to lbvservers
		for _, m := range mappings {
			lbNodeID := "lbvserver:" + m.LBVServer
			addChain(lbNodeID, chain)

			// Traverse from lbvserver to services
			for _, svcB := range svcBindingsByVS[m.LBVServer] {
				svcNodeID := "service:" + svcB.ServiceName
				addChain(svcNodeID, chain)
			}

			// Traverse from lbvserver to servicegroups
			for _, sgB := range sgBindingsByVS[m.LBVServer] {
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
