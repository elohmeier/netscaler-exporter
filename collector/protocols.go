package collector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elohmeier/netscaler-exporter/config"
	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

// collectProtocolHTTPStats collects protocol HTTP statistics
func (e *Exporter) collectProtocolHTTPStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetProtocolHTTPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol HTTP stats", "target", target.URL, "err", err)
		return
	}

	baseLabels := e.buildLabelValues(target)
	http := stats.ProtocolHTTPStats

	// Counters
	e.sendMetric(ch, e.httpTotalRequests, http.TotalRequests, baseLabels)
	e.sendMetric(ch, e.httpTotalResponses, http.TotalResponses, baseLabels)
	e.sendMetric(ch, e.httpTotalPosts, http.TotalPosts, baseLabels)
	e.sendMetric(ch, e.httpTotalGets, http.TotalGets, baseLabels)
	e.sendMetric(ch, e.httpTotalOthers, http.TotalOthers, baseLabels)
	e.sendMetric(ch, e.httpTotalRxRequestBytes, http.TotalRxRequestBytes, baseLabels)
	e.sendMetric(ch, e.httpTotalRxResponseBytes, http.TotalRxResponseBytes, baseLabels)
	e.sendMetric(ch, e.httpTotalTxRequestBytes, http.TotalTxRequestBytes, baseLabels)
	e.sendMetric(ch, e.httpTotal10Requests, http.Total10Requests, baseLabels)
	e.sendMetric(ch, e.httpTotal11Requests, http.Total11Requests, baseLabels)
	e.sendMetric(ch, e.httpTotal10Responses, http.Total10Responses, baseLabels)
	e.sendMetric(ch, e.httpTotal11Responses, http.Total11Responses, baseLabels)
	e.sendMetric(ch, e.httpTotalChunkedRequests, http.TotalChunkedRequests, baseLabels)
	e.sendMetric(ch, e.httpTotalChunkedResponses, http.TotalChunkedResponses, baseLabels)
	e.sendMetric(ch, e.httpTotalSPDYStreams, http.TotalSPDYStreams, baseLabels)
	e.sendMetric(ch, e.httpTotalSPDYv2Streams, http.TotalSPDYv2Streams, baseLabels)
	e.sendMetric(ch, e.httpTotalSPDYv3Streams, http.TotalSPDYv3Streams, baseLabels)
	e.sendMetric(ch, e.httpErrNoReuseMultipart, http.ErrNoReuseMultipart, baseLabels)
	e.sendMetric(ch, e.httpErrIncompleteHeaders, http.ErrIncompleteHeaders, baseLabels)
	e.sendMetric(ch, e.httpErrIncompleteRequests, http.ErrIncompleteRequests, baseLabels)
	e.sendMetric(ch, e.httpErrIncompleteResponses, http.ErrIncompleteResponses, baseLabels)
	e.sendMetric(ch, e.httpErrServerBusy, http.ErrServerBusy, baseLabels)
	e.sendMetric(ch, e.httpErrLargeContent, http.ErrLargeContent, baseLabels)
	e.sendMetric(ch, e.httpErrLargeChunk, http.ErrLargeChunk, baseLabels)
	e.sendMetric(ch, e.httpErrLargeCtlen, http.ErrLargeCtlen, baseLabels)

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
}

// collectProtocolTCPStats collects protocol TCP statistics
func (e *Exporter) collectProtocolTCPStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetProtocolTCPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol TCP stats", "target", target.URL, "err", err)
		return
	}

	baseLabels := e.buildLabelValues(target)
	tcp := stats.ProtocolTCPStats

	// Counters
	e.sendMetric(ch, e.tcpTotalRxPackets, tcp.TotalRxPackets, baseLabels)
	e.sendMetric(ch, e.tcpTotalRxBytes, tcp.TotalRxBytes, baseLabels)
	e.sendMetric(ch, e.tcpTotalTxBytes, tcp.TotalTxBytes, baseLabels)
	e.sendMetric(ch, e.tcpTotalTxPackets, tcp.TotalTxPackets, baseLabels)
	e.sendMetric(ch, e.tcpTotalClientConnOpened, tcp.TotalClientConnOpened, baseLabels)
	e.sendMetric(ch, e.tcpTotalServerConnOpened, tcp.TotalServerConnOpened, baseLabels)
	e.sendMetric(ch, e.tcpTotalSyn, tcp.TotalSyn, baseLabels)
	e.sendMetric(ch, e.tcpTotalSynProbe, tcp.TotalSynProbe, baseLabels)
	e.sendMetric(ch, e.tcpTotalServerFin, tcp.TotalServerFin, baseLabels)
	e.sendMetric(ch, e.tcpTotalClientFin, tcp.TotalClientFin, baseLabels)

	// Gauges
	e.sendMetric(ch, e.tcpActiveServerConn, tcp.ActiveServerConn, baseLabels)
	e.sendMetric(ch, e.tcpCurClientConnEstablished, tcp.CurClientConnEstablished, baseLabels)
	e.sendMetric(ch, e.tcpCurServerConnEstablished, tcp.CurServerConnEstablished, baseLabels)
	e.sendMetric(ch, e.tcpRxPacketsRate, tcp.RxPacketsRate, baseLabels)
	e.sendMetric(ch, e.tcpRxBytesRate, tcp.RxBytesRate, baseLabels)
	e.sendMetric(ch, e.tcpTxPacketsRate, tcp.TxPacketsRate, baseLabels)
	e.sendMetric(ch, e.tcpTxBytesRate, tcp.TxBytesRate, baseLabels)
	e.sendMetric(ch, e.tcpClientConnOpenedRate, tcp.ClientConnOpenedRate, baseLabels)
	e.sendMetric(ch, e.tcpErrBadChecksum, tcp.ErrBadChecksum, baseLabels)
	e.sendMetric(ch, e.tcpErrBadChecksumRate, tcp.ErrBadChecksumRate, baseLabels)
	e.sendMetric(ch, e.tcpErrAnyPortFail, tcp.ErrAnyPortFail, baseLabels)
	e.sendMetric(ch, e.tcpErrIPPortFail, tcp.ErrIPPortFail, baseLabels)
	e.sendMetric(ch, e.tcpErrBadStateConn, tcp.ErrBadStateConn, baseLabels)
	e.sendMetric(ch, e.tcpErrRstThreshold, tcp.ErrRstThreshold, baseLabels)
	e.sendMetric(ch, e.tcpSynRate, tcp.SynRate, baseLabels)
	e.sendMetric(ch, e.tcpSynProbeRate, tcp.SynProbeRate, baseLabels)
}

// collectProtocolIPStats collects protocol IP statistics
func (e *Exporter) collectProtocolIPStats(ctx context.Context, nsClient *netscaler.NitroClient, target config.Target, ch chan<- prometheus.Metric) {
	stats, err := netscaler.GetProtocolIPStats(ctx, nsClient, "")
	if err != nil {
		e.logger.Error("failed to get protocol IP stats", "target", target.URL, "err", err)
		return
	}

	baseLabels := e.buildLabelValues(target)
	ip := stats.ProtocolIPStats

	// Counters
	e.sendMetric(ch, e.ipTotalRxPackets, ip.TotalRxPackets, baseLabels)
	e.sendMetric(ch, e.ipTotalRxBytes, ip.TotalRxBytes, baseLabels)
	e.sendMetric(ch, e.ipTotalTxPackets, ip.TotalTxPackets, baseLabels)
	e.sendMetric(ch, e.ipTotalTxBytes, ip.TotalTxBytes, baseLabels)
	e.sendMetric(ch, e.ipTotalRxMbits, ip.TotalRxMbits, baseLabels)
	e.sendMetric(ch, e.ipTotalTxMbits, ip.TotalTxMbits, baseLabels)
	e.sendMetric(ch, e.ipTotalRoutedPackets, ip.TotalRoutedPackets, baseLabels)
	e.sendMetric(ch, e.ipTotalRoutedMbits, ip.TotalRoutedMbits, baseLabels)
	e.sendMetric(ch, e.ipTotalFragments, ip.TotalFragments, baseLabels)
	e.sendMetric(ch, e.ipTotalSuccReassembly, ip.TotalSuccReassembly, baseLabels)
	e.sendMetric(ch, e.ipTotalAddrLookup, ip.TotalAddrLookup, baseLabels)
	e.sendMetric(ch, e.ipTotalAddrLookupFail, ip.TotalAddrLookupFail, baseLabels)
	e.sendMetric(ch, e.ipTotalUDPFragmentsFwd, ip.TotalUDPFragmentsFwd, baseLabels)
	e.sendMetric(ch, e.ipTotalTCPFragmentsFwd, ip.TotalTCPFragmentsFwd, baseLabels)
	e.sendMetric(ch, e.ipTotalBadChecksums, ip.TotalBadChecksums, baseLabels)
	e.sendMetric(ch, e.ipTotalUnsuccReassembly, ip.TotalUnsuccReassembly, baseLabels)
	e.sendMetric(ch, e.ipTotalTooBig, ip.TotalTooBig, baseLabels)
	e.sendMetric(ch, e.ipTotalDupFragments, ip.TotalDupFragments, baseLabels)
	e.sendMetric(ch, e.ipTotalOutOfOrderFrag, ip.TotalOutOfOrderFrag, baseLabels)
	e.sendMetric(ch, e.ipTotalVIPDown, ip.TotalVIPDown, baseLabels)
	e.sendMetric(ch, e.ipTotalTTLExpired, ip.TotalTTLExpired, baseLabels)
	e.sendMetric(ch, e.ipTotalMaxClients, ip.TotalMaxClients, baseLabels)
	e.sendMetric(ch, e.ipTotalUnknownSvcs, ip.TotalUnknownSvcs, baseLabels)
	e.sendMetric(ch, e.ipTotalInvalidHeaderSz, ip.TotalInvalidHeaderSz, baseLabels)
	e.sendMetric(ch, e.ipTotalInvalidPacketSize, ip.TotalInvalidPacketSize, baseLabels)
	e.sendMetric(ch, e.ipTotalTruncatedPackets, ip.TotalTruncatedPackets, baseLabels)
	e.sendMetric(ch, e.ipNonIPTotalTruncatedPkts, ip.NonIPTotalTruncatedPkts, baseLabels)
	e.sendMetric(ch, e.ipTotalBadMacAddrs, ip.TotalBadMacAddrs, baseLabels)

	// Gauges (rates)
	e.sendMetric(ch, e.ipRxPacketsRate, ip.RxPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipRxBytesRate, ip.RxBytesRate, baseLabels)
	e.sendMetric(ch, e.ipTxPacketsRate, ip.TxPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipTxBytesRate, ip.TxBytesRate, baseLabels)
	e.sendMetric(ch, e.ipRxMbitsRate, ip.RxMbitsRate, baseLabels)
	e.sendMetric(ch, e.ipTxMbitsRate, ip.TxMbitsRate, baseLabels)
	e.sendMetric(ch, e.ipRoutedPacketsRate, ip.RoutedPacketsRate, baseLabels)
	e.sendMetric(ch, e.ipRoutedMbitsRate, ip.RoutedMbitsRate, baseLabels)
}

// sendMetric is a helper to parse and send a metric value.
// Accepts string or FlexString (or any type that can be converted to string via fmt.Sprint).
func (e *Exporter) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value any, labels []string) {
	val, _ := strconv.ParseFloat(fmt.Sprint(value), 64)
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, labels...)
}
