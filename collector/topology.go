package collector

import (
	"context"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collectTopologyMetrics(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	e.topologyNode.Reset()
	e.topologyEdge.Reset()

	// Collect LB Virtual Server nodes
	lbVServers, err := netscaler.GetVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting LB vserver stats for topology", "target", target.URL, "err", err)
	} else {
		for _, vs := range lbVServers.VirtualServerStats {
			nodeID := "lbvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			labels := e.buildLabelValues(target, nodeID, vs.Name, "lbvserver", state)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect CS Virtual Server nodes
	csVServers, err := netscaler.GetCSVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting CS vserver stats for topology", "target", target.URL, "err", err)
	} else {
		for _, vs := range csVServers.CSVirtualServerStats {
			nodeID := "csvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			labels := e.buildLabelValues(target, nodeID, vs.Name, "csvserver", state)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect Service nodes
	services, err := netscaler.GetServiceStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting service stats for topology", "target", target.URL, "err", err)
	} else {
		for _, svc := range services.ServiceStats {
			nodeID := "service:" + svc.Name
			state := "DOWN"
			value := 0.0
			if svc.State == "UP" {
				state = "UP"
				value = 1.0
			}
			labels := e.buildLabelValues(target, nodeID, svc.Name, "service", state)
			e.topologyNode.WithLabelValues(labels...).Set(value)
		}
	}

	// Collect Service Group nodes
	serviceGroups, err := netscaler.GetServiceGroups(ctx, nsClient, "attrs=servicegroupname")
	if err != nil {
		e.logger.Error("error getting service groups for topology", "target", target.URL, "err", err)
	} else {
		for _, sg := range serviceGroups.ServiceGroups {
			nodeID := "servicegroup:" + sg.Name
			labels := e.buildLabelValues(target, nodeID, sg.Name, "servicegroup", "UP")
			e.topologyNode.WithLabelValues(labels...).Set(1.0)
		}
	}

	// Collect LB VServer -> Service and Service Group edges
	if len(lbVServers.VirtualServerStats) > 0 {
		for _, vs := range lbVServers.VirtualServerStats {
			// Service bindings
			lbSvcBindings, err := netscaler.GetLBVServerServiceBindings(ctx, nsClient, vs.Name)
			if err != nil {
				e.logger.Debug("error getting LB vserver service bindings", "lbvserver", vs.Name, "target", target.URL, "err", err)
			} else {
				for _, b := range lbSvcBindings {
					edgeID := "lbvserver:" + b.Name + "->service:" + b.ServiceName
					sourceID := "lbvserver:" + b.Name
					targetID := "service:" + b.ServiceName
					weight := b.Weight
					if weight == "" {
						weight = "1"
					}
					labels := e.buildLabelValues(target, edgeID, sourceID, targetID, weight, "")
					e.topologyEdge.WithLabelValues(labels...).Set(1)
				}
			}

			// Service group bindings
			lbSgBindings, err := netscaler.GetLBVServerServiceGroupBindings(ctx, nsClient, vs.Name)
			if err != nil {
				e.logger.Debug("error getting LB vserver service group bindings", "lbvserver", vs.Name, "target", target.URL, "err", err)
			} else {
				for _, b := range lbSgBindings {
					edgeID := "lbvserver:" + b.Name + "->servicegroup:" + b.ServiceGroupName
					sourceID := "lbvserver:" + b.Name
					targetID := "servicegroup:" + b.ServiceGroupName
					weight := b.Weight
					if weight == "" {
						weight = "1"
					}
					labels := e.buildLabelValues(target, edgeID, sourceID, targetID, weight, "")
					e.topologyEdge.WithLabelValues(labels...).Set(1)
				}
			}
		}
	}

	// Collect CS VServer -> LB VServer edges
	if len(csVServers.CSVirtualServerStats) > 0 {
		for _, vs := range csVServers.CSVirtualServerStats {
			csLbBindings, err := netscaler.GetCSVServerLBVServerBindings(ctx, nsClient, vs.Name)
			if err != nil {
				e.logger.Debug("error getting CS vserver LB vserver bindings", "csvserver", vs.Name, "target", target.URL, "err", err)
			} else {
				for _, b := range csLbBindings {
					edgeID := "csvserver:" + b.Name + "->lbvserver:" + b.LBVServer
					sourceID := "csvserver:" + b.Name
					targetID := "lbvserver:" + b.LBVServer
					priority := b.Priority
					if priority == "" {
						priority = "0"
					}
					labels := e.buildLabelValues(target, edgeID, sourceID, targetID, "", priority)
					e.topologyEdge.WithLabelValues(labels...).Set(1)
				}
			}
		}
	}

	e.topologyNode.Collect(ch)
	e.topologyEdge.Collect(ch)
}
