package collector

import (
	"context"
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
			labels := e.buildLabelValues(nodeID, vs.Name, "lbvserver", state)
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
			labels := e.buildLabelValues(nodeID, vs.Name, "csvserver", state)
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
			labels := e.buildLabelValues(nodeID, svc.Name, "service", state)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect Service Group nodes
	serviceGroups, err := netscaler.GetServiceGroups(ctx, nsClient, "attrs=servicegroupname")
	if err != nil {
		e.logger.Error("error getting service groups for topology", "url", e.url, "err", err)
	} else {
		for _, sg := range serviceGroups.ServiceGroups {
			nodeID := "servicegroup:" + sg.Name
			labels := e.buildLabelValues(nodeID, sg.Name, "servicegroup", "UP")
			e.topologyNode.WithLabelValues(labels...).Set(1.0)
		}
	}

	// Collect LB VServer -> Service and Service Group edges using lookup maps
	if len(lbVServers.VirtualServerStats) > 0 {
		for _, vs := range lbVServers.VirtualServerStats {
			// Service bindings from lookup map
			for _, b := range svcBindingsByVS[vs.Name] {
				edgeID := "lbvserver:" + b.Name + "->service:" + b.ServiceName
				sourceID := "lbvserver:" + b.Name
				targetID := "service:" + b.ServiceName
				weight := b.Weight
				if weight == "" {
					weight = "1"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, weight, "")
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}

			// Service group bindings from lookup map
			for _, b := range sgBindingsByVS[vs.Name] {
				edgeID := "lbvserver:" + b.Name + "->servicegroup:" + b.ServiceGroupName
				sourceID := "lbvserver:" + b.Name
				targetID := "servicegroup:" + b.ServiceGroupName
				weight := b.Weight
				if weight == "" {
					weight = "1"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, weight, "")
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}
		}
	}

	// Collect CS VServer -> LB VServer edges using lookup map
	if len(csVServers.CSVirtualServerStats) > 0 {
		for _, vs := range csVServers.CSVirtualServerStats {
			for _, b := range csBindingsByVS[vs.Name] {
				edgeID := "csvserver:" + b.Name + "->lbvserver:" + b.LBVServer
				sourceID := "csvserver:" + b.Name
				targetID := "lbvserver:" + b.LBVServer
				priority := b.Priority
				if priority == "" {
					priority = "0"
				}
				labels := e.buildLabelValues(edgeID, sourceID, targetID, "", priority)
				e.topologyEdge.WithLabelValues(labels...).Set(1)
			}
		}
	}

	e.topologyNode.Collect(ch)
	e.topologyEdge.Collect(ch)
}
