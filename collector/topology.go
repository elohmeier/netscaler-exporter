package collector

import (
	"context"

	"github.com/elohmeier/netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	topologyNode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "netscaler_topology_node",
			Help: "Node for topology visualization (1=UP, 0=DOWN)",
		},
		[]string{
			"ns_instance",
			"id",
			"title",
			"node_type",
			"state",
		},
	)

	topologyEdge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "netscaler_topology_edge",
			Help: "Edge representing binding between frontend and backend",
		},
		[]string{
			"ns_instance",
			"id",
			"source",
			"target",
			"weight",
			"priority",
		},
	)
)

func (e *Exporter) collectTopologyMetrics(ctx context.Context, nsClient *netscaler.NitroClient, nsInstance string, ch chan<- prometheus.Metric) {
	e.topologyNode.Reset()
	e.topologyEdge.Reset()

	// Collect LB Virtual Server nodes
	lbVServers, err := netscaler.GetVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting LB vserver stats for topology", "err", err)
	} else {
		for _, vs := range lbVServers.VirtualServerStats {
			nodeID := "lbvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			e.topologyNode.WithLabelValues(nsInstance, nodeID, vs.Name, "lbvserver", state).Set(value)
		}
	}

	// Collect CS Virtual Server nodes
	csVServers, err := netscaler.GetCSVirtualServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting CS vserver stats for topology", "err", err)
	} else {
		for _, vs := range csVServers.CSVirtualServerStats {
			nodeID := "csvserver:" + vs.Name
			state := "DOWN"
			value := 0.0
			if vs.State == "UP" {
				state = "UP"
				value = 1.0
			}
			e.topologyNode.WithLabelValues(nsInstance, nodeID, vs.Name, "csvserver", state).Set(value)
		}
	}

	// Collect Service nodes
	services, err := netscaler.GetServiceStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("error getting service stats for topology", "err", err)
	} else {
		for _, svc := range services.ServiceStats {
			nodeID := "service:" + svc.Name
			state := "DOWN"
			value := 0.0
			if svc.State == "UP" {
				state = "UP"
				value = 1.0
			}
			e.topologyNode.WithLabelValues(nsInstance, nodeID, svc.Name, "service", state).Set(value)
		}
	}

	// Collect Service Group nodes
	serviceGroups, err := netscaler.GetServiceGroups(ctx, nsClient, "attrs=servicegroupname")
	if err != nil {
		e.logger.Error("error getting service groups for topology", "err", err)
	} else {
		for _, sg := range serviceGroups.ServiceGroups {
			nodeID := "servicegroup:" + sg.Name
			e.topologyNode.WithLabelValues(nsInstance, nodeID, sg.Name, "servicegroup", "UP").Set(1.0)
		}
	}

	// Collect LB VServer -> Service edges
	lbSvcBindings, err := netscaler.GetLBVServerServiceBindings(ctx, nsClient)
	if err != nil {
		e.logger.Debug("error getting LB vserver service bindings", "err", err)
	} else {
		for _, b := range lbSvcBindings {
			edgeID := "lbvserver:" + b.Name + "->service:" + b.ServiceName
			sourceID := "lbvserver:" + b.Name
			targetID := "service:" + b.ServiceName
			weight := b.Weight
			if weight == "" {
				weight = "1"
			}
			e.topologyEdge.WithLabelValues(nsInstance, edgeID, sourceID, targetID, weight, "").Set(1)
		}
	}

	// Collect LB VServer -> Service Group edges
	lbSgBindings, err := netscaler.GetLBVServerServiceGroupBindings(ctx, nsClient)
	if err != nil {
		e.logger.Debug("error getting LB vserver service group bindings", "err", err)
	} else {
		for _, b := range lbSgBindings {
			edgeID := "lbvserver:" + b.Name + "->servicegroup:" + b.ServiceGroupName
			sourceID := "lbvserver:" + b.Name
			targetID := "servicegroup:" + b.ServiceGroupName
			weight := b.Weight
			if weight == "" {
				weight = "1"
			}
			e.topologyEdge.WithLabelValues(nsInstance, edgeID, sourceID, targetID, weight, "").Set(1)
		}
	}

	// Collect CS VServer -> LB VServer edges
	csLbBindings, err := netscaler.GetCSVServerLBVServerBindings(ctx, nsClient)
	if err != nil {
		e.logger.Debug("error getting CS vserver LB vserver bindings", "err", err)
	} else {
		for _, b := range csLbBindings {
			edgeID := "csvserver:" + b.Name + "->lbvserver:" + b.LBVServer
			sourceID := "csvserver:" + b.Name
			targetID := "lbvserver:" + b.LBVServer
			priority := b.Priority
			if priority == "" {
				priority = "0"
			}
			e.topologyEdge.WithLabelValues(nsInstance, edgeID, sourceID, targetID, "", priority).Set(1)
		}
	}

	e.topologyNode.Collect(ch)
	e.topologyEdge.Collect(ch)
}
