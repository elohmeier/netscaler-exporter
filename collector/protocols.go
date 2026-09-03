package collector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// collectProtocolHTTPStats collects protocol HTTP statistics
func (e *Exporter) collectProtocolHTTPStats(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) bool {
	stats, err := netscaler.GetProtocolHTTPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol HTTP stats", "url", e.url, "err", err)
		return false
	}

	baseLabels := e.buildLabelValues()
	http := stats.ProtocolHTTPStats

	// Counters
	e.sendCounterMetric(ch, e.httpTotalRequests, http.TotalRequests, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalResponses, http.TotalResponses, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalPosts, http.TotalPosts, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalGets, http.TotalGets, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalOthers, http.TotalOthers, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalRxRequestBytes, http.TotalRxRequestBytes, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalRxResponseBytes, http.TotalRxResponseBytes, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalTxRequestBytes, http.TotalTxRequestBytes, baseLabels)
	e.sendCounterMetric(ch, e.httpTotal10Requests, http.Total10Requests, baseLabels)
	e.sendCounterMetric(ch, e.httpTotal11Requests, http.Total11Requests, baseLabels)
	e.sendCounterMetric(ch, e.httpTotal10Responses, http.Total10Responses, baseLabels)
	e.sendCounterMetric(ch, e.httpTotal11Responses, http.Total11Responses, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalChunkedRequests, http.TotalChunkedRequests, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalChunkedResponses, http.TotalChunkedResponses, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalSPDYStreams, http.TotalSPDYStreams, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalSPDYv2Streams, http.TotalSPDYv2Streams, baseLabels)
	e.sendCounterMetric(ch, e.httpTotalSPDYv3Streams, http.TotalSPDYv3Streams, baseLabels)
	e.sendCounterMetric(ch, e.httpErrNoReuseMultipart, http.ErrNoReuseMultipart, baseLabels)
	e.sendCounterMetric(ch, e.httpErrIncompleteHeaders, http.ErrIncompleteHeaders, baseLabels)
	e.sendCounterMetric(ch, e.httpErrIncompleteRequests, http.ErrIncompleteRequests, baseLabels)
	e.sendCounterMetric(ch, e.httpErrIncompleteResponses, http.ErrIncompleteResponses, baseLabels)
	e.sendCounterMetric(ch, e.httpErrServerBusy, http.ErrServerBusy, baseLabels)
	e.sendCounterMetric(ch, e.httpErrLargeContent, http.ErrLargeContent, baseLabels)
	e.sendCounterMetric(ch, e.httpErrLargeChunk, http.ErrLargeChunk, baseLabels)
	e.sendCounterMetric(ch, e.httpErrLargeCtlen, http.ErrLargeCtlen, baseLabels)

	// Gauges (rates)
	e.sendMetric(ch, e.httpRequestsRate, http.RequestsRate, baseLabels)
	e.sendMetric(ch, e.httpResponsesRate, http.ResponsesRate, baseLabels)
	e.sendMetric(ch, e.httpPostsRate, http.PostsRate, baseLabels)
	e.sendMetric(ch, e.httpGetsRate, http.GetsRate, baseLabels)
	e.sendMetric(ch, e.httpOthersRate, http.OthersRate, baseLabels)
	e.sendMetric(ch, e.httpRxRequestBytesRate, http.RxRequestBytesRate, baseLabels)
	e.sendMetric(ch, e.httpRxResponseBytesRate, http.RxResponseBytesRate, baseLabels)
	e.sendMetric(ch, e.httpTxRequestBytesRate, http.TxRequestBytesRate, baseLabels)
	e.sendMetric(ch, e.httpRequest10Rate, http.Request10Rate, baseLabels)
	e.sendMetric(ch, e.httpRequest11Rate, http.Request11Rate, baseLabels)
	e.sendMetric(ch, e.httpResponse10Rate, http.Response10Rate, baseLabels)
	e.sendMetric(ch, e.httpResponse11Rate, http.Response11Rate, baseLabels)
	e.sendMetric(ch, e.httpChunkedRequestsRate, http.ChunkedRequestsRate, baseLabels)
	e.sendMetric(ch, e.httpChunkedResponsesRate, http.ChunkedResponsesRate, baseLabels)
	e.sendMetric(ch, e.httpSPDYStreamsRate, http.SPDYStreamsRate, baseLabels)
	e.sendMetric(ch, e.httpSPDYv2StreamsRate, http.SPDYv2StreamsRate, baseLabels)
	e.sendMetric(ch, e.httpSPDYv3StreamsRate, http.SPDYv3StreamsRate, baseLabels)
	e.sendMetric(ch, e.httpErrNoReuseMultipartRate, http.ErrNoReuseMultipartRate, baseLabels)
	e.sendMetric(ch, e.httpErrIncompleteRequestsRate, http.ErrIncompleteRequestsRate, baseLabels)
	e.sendMetric(ch, e.httpErrIncompleteResponsesRate, http.ErrIncompleteResponsesRate, baseLabels)
	e.sendMetric(ch, e.httpErrServerBusyRate, http.ErrServerBusyRate, baseLabels)
	return true
}

// collectProtocolTCPStats collects protocol TCP statistics
func (e *Exporter) collectProtocolTCPStats(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) bool {
	stats, err := netscaler.GetProtocolTCPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol TCP stats", "url", e.url, "err", err)
		return false
	}

	baseLabels := e.buildLabelValues()
	tcp := stats.ProtocolTCPStats

	// Counters
	e.sendCounterMetric(ch, e.tcpTotalRxPackets, tcp.TotalRxPackets, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalRxBytes, tcp.TotalRxBytes, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalTxBytes, tcp.TotalTxBytes, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalTxPackets, tcp.TotalTxPackets, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalClientConnOpened, tcp.TotalClientConnOpened, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalServerConnOpened, tcp.TotalServerConnOpened, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalSyn, tcp.TotalSyn, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalSynProbe, tcp.TotalSynProbe, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalServerFin, tcp.TotalServerFin, baseLabels)
	e.sendCounterMetric(ch, e.tcpTotalClientFin, tcp.TotalClientFin, baseLabels)

	// Gauges (current values and rates)
	e.sendMetric(ch, e.tcpActiveServerConn, tcp.ActiveServerConn, baseLabels)
	e.sendMetric(ch, e.tcpCurClientConnEstablished, tcp.CurClientConnEstablished, baseLabels)
	e.sendMetric(ch, e.tcpCurServerConnEstablished, tcp.CurServerConnEstablished, baseLabels)
	e.sendMetric(ch, e.tcpRxPacketsRate, tcp.RxPacketsRate, baseLabels)
	e.sendMetric(ch, e.tcpRxBytesRate, tcp.RxBytesRate, baseLabels)
	e.sendMetric(ch, e.tcpTxPacketsRate, tcp.TxPacketsRate, baseLabels)
	e.sendMetric(ch, e.tcpTxBytesRate, tcp.TxBytesRate, baseLabels)
	e.sendMetric(ch, e.tcpClientConnOpenedRate, tcp.ClientConnOpenedRate, baseLabels)
	e.sendMetric(ch, e.tcpErrBadChecksumRate, tcp.ErrBadChecksumRate, baseLabels)
	e.sendMetric(ch, e.tcpSynRate, tcp.SynRate, baseLabels)
	e.sendMetric(ch, e.tcpSynProbeRate, tcp.SynProbeRate, baseLabels)

	// Error counters
	e.sendCounterMetric(ch, e.tcpErrBadChecksum, tcp.ErrBadChecksum, baseLabels)
	e.sendCounterMetric(ch, e.tcpErrAnyPortFail, tcp.ErrAnyPortFail, baseLabels)
	e.sendCounterMetric(ch, e.tcpErrIPPortFail, tcp.ErrIPPortFail, baseLabels)
	e.sendCounterMetric(ch, e.tcpErrBadStateConn, tcp.ErrBadStateConn, baseLabels)
	e.sendCounterMetric(ch, e.tcpErrRstThreshold, tcp.ErrRstThreshold, baseLabels)
	return true
}

// collectProtocolIPStats collects protocol IP statistics
func (e *Exporter) collectProtocolIPStats(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) bool {
	stats, err := netscaler.GetProtocolIPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol IP stats", "url", e.url, "err", err)
		return false
	}

	baseLabels := e.buildLabelValues()
	ip := stats.ProtocolIPStats

	// Counters
	e.sendCounterMetric(ch, e.ipTotalRxPackets, ip.TotalRxPackets, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalRxBytes, ip.TotalRxBytes, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTxPackets, ip.TotalTxPackets, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTxBytes, ip.TotalTxBytes, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalRxMbits, ip.TotalRxMbits, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTxMbits, ip.TotalTxMbits, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalRoutedPackets, ip.TotalRoutedPackets, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalRoutedMbits, ip.TotalRoutedMbits, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalFragments, ip.TotalFragments, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalSuccReassembly, ip.TotalSuccReassembly, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalAddrLookup, ip.TotalAddrLookup, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalAddrLookupFail, ip.TotalAddrLookupFail, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalUDPFragmentsFwd, ip.TotalUDPFragmentsFwd, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTCPFragmentsFwd, ip.TotalTCPFragmentsFwd, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalBadChecksums, ip.TotalBadChecksums, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalUnsuccReassembly, ip.TotalUnsuccReassembly, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTooBig, ip.TotalTooBig, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalDupFragments, ip.TotalDupFragments, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalOutOfOrderFrag, ip.TotalOutOfOrderFrag, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalVIPDown, ip.TotalVIPDown, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTTLExpired, ip.TotalTTLExpired, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalMaxClients, ip.TotalMaxClients, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalUnknownSvcs, ip.TotalUnknownSvcs, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalInvalidHeaderSz, ip.TotalInvalidHeaderSz, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalInvalidPacketSize, ip.TotalInvalidPacketSize, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalTruncatedPackets, ip.TotalTruncatedPackets, baseLabels)
	e.sendCounterMetric(ch, e.ipNonIPTotalTruncatedPkts, ip.NonIPTotalTruncatedPkts, baseLabels)
	e.sendCounterMetric(ch, e.ipTotalBadMacAddrs, ip.TotalBadMacAddrs, baseLabels)

	// Gauges (rates)
	e.sendMetric(ch, e.ipRxPacketsRate, ip.RxPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipRxBytesRate, ip.RxBytesRate, baseLabels)
	e.sendMetric(ch, e.ipTxPacketsRate, ip.TxPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipTxBytesRate, ip.TxBytesRate, baseLabels)
	e.sendMetric(ch, e.ipRxMbitsRate, ip.RxMbitsRate, baseLabels)
	e.sendMetric(ch, e.ipTxMbitsRate, ip.TxMbitsRate, baseLabels)
	e.sendMetric(ch, e.ipRoutedPacketsRate, ip.RoutedPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipRoutedMbitsRate, ip.RoutedMbitsRate, baseLabels)
	return true
}

// sendMetric is a helper to parse and send a metric value.
// Accepts string or FlexString (or any type that can be converted to string via fmt.Sprint).
func (e *Exporter) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value any, labels []string) {
	val, _ := strconv.ParseFloat(fmt.Sprint(value), 64)
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, labels...)
}

// sendCounterMetric parses and sends an absolute cumulative value reported by
// the NetScaler.
func (e *Exporter) sendCounterMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value any, labels []string) {
	val, _ := strconv.ParseFloat(fmt.Sprint(value), 64)
	ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, val, labels...)
}
