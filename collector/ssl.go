package collector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// collectSSLStats collects SSL global statistics
func (e *Exporter) collectSSLStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetSSLStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get SSL stats", "target", target.URL, "err", err)
		return
	}

	baseLabels := e.buildLabelValues(target)
	ssl := stats.SSLStats

	// Counters
	e.sendMetric(ch, e.sslTotalTLSv11Sessions, ssl.TotalTLSv11Sessions, baseLabels)
	e.sendMetric(ch, e.sslTotalSSLv2Sessions, ssl.TotalSSLv2Sessions, baseLabels)
	e.sendMetric(ch, e.sslTotalSessions, ssl.TotalSessions, baseLabels)
	e.sendMetric(ch, e.sslTotalSSLv2Handshakes, ssl.TotalSSLv2Handshakes, baseLabels)
	e.sendMetric(ch, e.sslTotalEnc, ssl.TotalEnc, baseLabels)
	e.sendMetric(ch, e.sslCryptoUtilization, ssl.CryptoUtilizationStat, baseLabels)
	e.sendMetric(ch, e.sslTotalNewSessions, ssl.TotalNewSessions, baseLabels)

	// Gauges
	e.sendMetric(ch, e.sslSessionsRate, ssl.SessionsRate, baseLabels)
	e.sendMetric(ch, e.sslDecRate, ssl.DecRate, baseLabels)
	e.sendMetric(ch, e.sslEncRate, ssl.EncRate, baseLabels)
	e.sendMetric(ch, e.sslSSLv2HandshakesRate, ssl.SSLv2HandshakesRate, baseLabels)
	e.sendMetric(ch, e.sslNewSessionsRate, ssl.NewSessionsRate, baseLabels)
}

// collectSSLCertKeys collects SSL certificate expiration metrics
func (e *Exporter) collectSSLCertKeys(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetSSLCertKeys(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get SSL cert keys", "target", target.URL, "err", err)
		return
	}

	e.sslCertDaysToExpire.Reset()
	for _, cert := range stats.SSLCertKeys {
		val, _ := strconv.ParseFloat(fmt.Sprint(cert.DaysToExpiration), 64)
		labels := e.buildLabelValues(target, cert.CertKey)
		e.sslCertDaysToExpire.WithLabelValues(labels...).Set(val)
	}
	e.sslCertDaysToExpire.Collect(ch)
}

// collectSSLVServerStats collects SSL virtual server statistics
func (e *Exporter) collectSSLVServerStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetSSLVServerStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get SSL vserver stats", "target", target.URL, "err", err)
		return
	}

	// Reset all gauges
	e.sslVServerTotalDecBytes.Reset()
	e.sslVServerTotalEncBytes.Reset()
	e.sslVServerTotalHWDecBytes.Reset()
	e.sslVServerTotalHWEncBytes.Reset()
	e.sslVServerTotalSessionNew.Reset()
	e.sslVServerTotalSessionHits.Reset()
	e.sslVServerTotalClientAuthSuccess.Reset()
	e.sslVServerTotalClientAuthFailure.Reset()
	e.sslVServerHealth.Reset()
	e.sslVServerActiveServices.Reset()
	e.sslVServerClientAuthSuccessRate.Reset()
	e.sslVServerClientAuthFailureRate.Reset()
	e.sslVServerEncBytesRate.Reset()
	e.sslVServerDecBytesRate.Reset()
	e.sslVServerHWEncBytesRate.Reset()
	e.sslVServerHWDecBytesRate.Reset()
	e.sslVServerSessionNewRate.Reset()
	e.sslVServerSessionHitsRate.Reset()

	for _, vs := range stats.SSLVServerStats {
		labels := e.buildLabelValues(target, vs.VServerName, vs.Type, vs.PrimaryIPAddress)

		setGaugeVal(e.sslVServerTotalDecBytes, labels, vs.TotalDecBytes)
		setGaugeVal(e.sslVServerTotalEncBytes, labels, vs.TotalEncBytes)
		setGaugeVal(e.sslVServerTotalHWDecBytes, labels, vs.TotalHWDecBytes)
		setGaugeVal(e.sslVServerTotalHWEncBytes, labels, vs.TotalHWEncBytes)
		setGaugeVal(e.sslVServerTotalSessionNew, labels, vs.TotalSessionNew)
		setGaugeVal(e.sslVServerTotalSessionHits, labels, vs.TotalSessionHits)
		setGaugeVal(e.sslVServerTotalClientAuthSuccess, labels, vs.TotalClientAuthSuccess)
		setGaugeVal(e.sslVServerTotalClientAuthFailure, labels, vs.TotalClientAuthFailure)
		setGaugeVal(e.sslVServerHealth, labels, vs.Health)
		setGaugeVal(e.sslVServerActiveServices, labels, vs.ActiveServices)
		setGaugeVal(e.sslVServerClientAuthSuccessRate, labels, vs.ClientAuthSuccessRate)
		setGaugeVal(e.sslVServerClientAuthFailureRate, labels, vs.ClientAuthFailureRate)
		setGaugeVal(e.sslVServerEncBytesRate, labels, vs.EncBytesRate)
		setGaugeVal(e.sslVServerDecBytesRate, labels, vs.DecBytesRate)
		setGaugeVal(e.sslVServerHWEncBytesRate, labels, vs.HWEncBytesRate)
		setGaugeVal(e.sslVServerHWDecBytesRate, labels, vs.HWDecBytesRate)
		setGaugeVal(e.sslVServerSessionNewRate, labels, vs.SessionNewRate)
		setGaugeVal(e.sslVServerSessionHitsRate, labels, vs.SessionHitsRate)
	}

	// Collect all
	e.sslVServerTotalDecBytes.Collect(ch)
	e.sslVServerTotalEncBytes.Collect(ch)
	e.sslVServerTotalHWDecBytes.Collect(ch)
	e.sslVServerTotalHWEncBytes.Collect(ch)
	e.sslVServerTotalSessionNew.Collect(ch)
	e.sslVServerTotalSessionHits.Collect(ch)
	e.sslVServerTotalClientAuthSuccess.Collect(ch)
	e.sslVServerTotalClientAuthFailure.Collect(ch)
	e.sslVServerHealth.Collect(ch)
	e.sslVServerActiveServices.Collect(ch)
	e.sslVServerClientAuthSuccessRate.Collect(ch)
	e.sslVServerClientAuthFailureRate.Collect(ch)
	e.sslVServerEncBytesRate.Collect(ch)
	e.sslVServerDecBytesRate.Collect(ch)
	e.sslVServerHWEncBytesRate.Collect(ch)
	e.sslVServerHWDecBytesRate.Collect(ch)
	e.sslVServerSessionNewRate.Collect(ch)
	e.sslVServerSessionHitsRate.Collect(ch)
}

// collectSystemCPUStats collects per-core CPU statistics
func (e *Exporter) collectSystemCPUStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetSystemCPUStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get system CPU stats", "target", target.URL, "err", err)
		return
	}

	e.cpuCoreUsage.Reset()
	for _, cpu := range stats.SystemCPUStats {
		val, _ := strconv.ParseFloat(cpu.PerCPUUsage, 64)
		labels := e.buildLabelValues(target, cpu.ID)
		e.cpuCoreUsage.WithLabelValues(labels...).Set(val)
	}
	e.cpuCoreUsage.Collect(ch)
}

// collectNSCapacityStats collects bandwidth capacity statistics
func (e *Exporter) collectNSCapacityStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetNSCapacityStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get bandwidth capacity stats", "target", target.URL, "err", err)
		return
	}

	baseLabels := e.buildLabelValues(target)
	cap := stats.NSCapacityStats

	e.sendMetric(ch, e.capacityMaxBandwidth, cap.MaxBandwidth, baseLabels)
	e.sendMetric(ch, e.capacityMinBandwidth, cap.MinBandwidth, baseLabels)
	e.sendMetric(ch, e.capacityActualBandwidth, cap.ActualBandwidth, baseLabels)
	e.sendMetric(ch, e.capacityBandwidth, cap.Bandwidth, baseLabels)
}

// setGaugeVal is a helper to set a gauge value
func setGaugeVal(g *prometheus.GaugeVec, labels []string, value any) {
	val, _ := strconv.ParseFloat(fmt.Sprint(value), 64)
	g.WithLabelValues(labels...).Set(val)
}
