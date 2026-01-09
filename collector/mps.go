package collector

import (
	"strconv"

	"github.com/elohmeier/netscaler-exporter/netscaler"
)

// collectMPSHealth collects MPS health metrics (CPU, disk, memory usage).
func (e *Exporter) collectMPSHealth(mps netscaler.MPSAPIResponse) {
	e.mpsHealthCPUUsage.Reset()
	e.mpsHealthDiskUsage.Reset()
	e.mpsHealthDiskFree.Reset()
	e.mpsHealthDiskTotal.Reset()
	e.mpsHealthDiskUsed.Reset()
	e.mpsHealthMemoryUsage.Reset()
	e.mpsHealthMemoryFree.Reset()
	e.mpsHealthMemoryTotal.Reset()

	for _, health := range mps.MPSHealth {
		labels := e.buildLabelValues(health.NodeType)

		if cpuUsage, err := strconv.ParseFloat(health.CPUUsage, 64); err == nil {
			e.mpsHealthCPUUsage.WithLabelValues(labels...).Set(cpuUsage)
		}

		if diskUsage, err := strconv.ParseFloat(health.DiskUsage, 64); err == nil {
			e.mpsHealthDiskUsage.WithLabelValues(labels...).Set(diskUsage)
		}

		if diskFree, err := strconv.ParseFloat(health.DiskFree, 64); err == nil {
			e.mpsHealthDiskFree.WithLabelValues(labels...).Set(diskFree)
		}

		if diskTotal, err := strconv.ParseFloat(health.DiskTotal, 64); err == nil {
			e.mpsHealthDiskTotal.WithLabelValues(labels...).Set(diskTotal)
		}

		if diskUsed, err := strconv.ParseFloat(health.DiskUsed, 64); err == nil {
			e.mpsHealthDiskUsed.WithLabelValues(labels...).Set(diskUsed)
		}

		if memoryUsage, err := strconv.ParseFloat(health.MemoryUsage, 64); err == nil {
			e.mpsHealthMemoryUsage.WithLabelValues(labels...).Set(memoryUsage)
		}

		if memoryFree, err := strconv.ParseFloat(health.MemoryFree, 64); err == nil {
			e.mpsHealthMemoryFree.WithLabelValues(labels...).Set(memoryFree)
		}

		if memoryTotal, err := strconv.ParseFloat(health.MemoryTotal, 64); err == nil {
			e.mpsHealthMemoryTotal.WithLabelValues(labels...).Set(memoryTotal)
		}
	}
}
