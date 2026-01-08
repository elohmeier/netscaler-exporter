package netscaler

import (
	"encoding/json"
	"fmt"
)

// FlexString handles JSON values that may be string or number.
// NetScaler API returns some fields as strings and others as numbers depending on context.
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}
	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexString(n.String())
		return nil
	}
	return fmt.Errorf("FlexString: cannot unmarshal %s", string(data))
}

// NSAPIResponse represents the main portion of the Nitro API response
type NSAPIResponse struct {
	Errorcode               int64                     `json:"errorcode"`
	Message                 string                    `json:"message"`
	Severity                string                    `json:"severity"`
	NSLicense               NSLicense                 `json:"nslicense"`
	NSStats                 NSStats                   `json:"ns"`
	InterfaceStats          []InterfaceStats          `json:"Interface"`
	VirtualServerStats      []VirtualServerStats      `json:"lbvserver"`
	ServiceStats            []ServiceStats            `json:"service"`
	ServiceGroups           []ServiceGroups           `json:"servicegroup"`
	ServiceGroupMemberStats []ServiceGroupMemberStats `json:"servicegroupmember"`
	GSLBServiceStats        []GSLBServiceStats        `json:"gslbservice"`
	GSLBVirtualServerStats  []GSLBVirtualServerStats  `json:"gslbvserver"`
	CSVirtualServerStats    []CSVirtualServerStats    `json:"csvserver"`
	VPNVirtualServerStats   []VPNVirtualServerStats   `json:"vpnvserver"`
	AAAStats                AAAStats                  `json:"aaa"`
	ProtocolHTTPStats       ProtocolHTTPStats         `json:"protocolhttp"`
	ProtocolTCPStats        ProtocolTCPStats          `json:"protocoltcp"`
	ProtocolIPStats         ProtocolIPStats           `json:"protocolip"`
	SSLStats                SSLStats                  `json:"ssl"`
	SSLCertKeys             []SSLCertKey              `json:"sslcertkey"`
	SSLVServerStats         []SSLVServerStats         `json:"sslvserver"`
	SystemCPUStats          []SystemCPUStats          `json:"systemcpu"`
	NSCapacityStats         NSCapacityStats           `json:"nscapacity"`
}

// NSLicense represents the data returned from the /config/nslicense Nitro API endpoint
type NSLicense struct {
	ModelID string `json:"modelid"`
}

// NSStats represents the data returned from the /stat/ns Nitro API endpoint
type NSStats struct {
	CPUUsagePcnt                           float64 `json:"cpuusagepcnt"`
	MemUsagePcnt                           float64 `json:"memusagepcnt"`
	MgmtCPUUsagePcnt                       float64 `json:"mgmtcpuusagepcnt"`
	PktCPUUsagePcnt                        float64 `json:"pktcpuusagepcnt"`
	FlashPartitionUsage                    float64 `json:"disk0perusage"`
	VarPartitionUsage                      float64 `json:"disk1perusage"`
	TotalReceivedMB                        string  `json:"totrxmbits"`
	TotalTransmitMB                        string  `json:"tottxmbits"`
	HTTPRequests                           string  `json:"httptotrequests"`
	HTTPResponses                          string  `json:"httptotresponses"`
	TCPCurrentClientConnections            string  `json:"tcpcurclientconn"`
	TCPCurrentClientConnectionsEstablished string  `json:"tcpcurclientconnestablished"`
	TCPCurrentServerConnections            string  `json:"tcpcurserverconn"`
	TCPCurrentServerConnectionsEstablished string  `json:"tcpcurserverconnestablished"`
}

// InterfaceStats represents the data returned from the /stat/interface Nitro API endpoint
type InterfaceStats struct {
	ID                      string `json:"id"`
	TotalReceivedBytes      string `json:"totrxbytes"`
	TotalTransmitBytes      string `json:"tottxbytes"`
	TotalReceivedPackets    string `json:"totrxpkts"`
	TotalTransmitPackets    string `json:"tottxpkts"`
	JumboPacketsReceived    string `json:"jumbopktsreceived"`
	JumboPacketsTransmitted string `json:"jumbopktstransmitted"`
	ErrorPacketsReceived    string `json:"errpktrx"`
	Alias                   string `json:"interfacealias"`
}

// VirtualServerStats represents the data returned from the /stat/lbvserver Nitro API endpoint
type VirtualServerStats struct {
	Name                     string `json:"name"`
	State                    string `json:"state"`
	WaitingRequests          string `json:"vsvrsurgecount"`
	Health                   string `json:"vslbhealth"`
	InactiveServices         string `json:"inactsvcs"`
	ActiveServices           string `json:"actsvcs"`
	TotalHits                string `json:"tothits"`
	TotalRequests            string `json:"totalrequests"`
	TotalResponses           string `json:"totalresponses"`
	TotalRequestBytes        string `json:"totalrequestbytes"`
	TotalResponseBytes       string `json:"totalresponsebytes"`
	CurrentClientConnections string `json:"curclntconnections"`
	CurrentServerConnections string `json:"cursrvrconnections"`
}

// ServiceStats represents the data returned from the /stat/service Nitro API endpoint
type ServiceStats struct {
	Name                         string `json:"name"`
	Throughput                   string `json:"throughput"`
	AvgTimeToFirstByte           string `json:"avgsvrttfb"`
	State                        string `json:"state"`
	TotalRequests                string `json:"totalrequests"`
	TotalResponses               string `json:"totalresponses"`
	TotalRequestBytes            string `json:"totalrequestbytes"`
	TotalResponseBytes           string `json:"totalresponsebytes"`
	CurrentClientConnections     string `json:"curclntconnections"`
	SurgeCount                   string `json:"surgecount"`
	CurrentServerConnections     string `json:"cursrvrconnections"`
	ServerEstablishedConnections string `json:"svrestablishedconn"`
	CurrentReusePool             string `json:"curreusepool"`
	MaxClients                   string `json:"maxclients"`
	CurrentLoad                  string `json:"curload"`
	ServiceHits                  string `json:"vsvrservicehits"`
	ActiveTransactions           string `json:"activetransactions"`
}

// ServiceGroups represents the data returned from the /config/servicegroup Nitro API endpoint
type ServiceGroups struct {
	Name                string                    `json:"servicegroupname"`
	ServiceGroupMembers []ServiceGroupMemberStats `json:"servicegroupmember"`
}

// ServiceGroupMemberStats represents the data returned from the /stat/servicegroupmember Nitro API endpoint
type ServiceGroupMemberStats struct {
	PrimaryPort                  int    `json:"primaryport"`
	State                        string `json:"state"`
	AvgTimeToFirstByte           string `json:"avgsvrttfb"`
	TotalRequests                string `json:"totalrequests"`
	TotalResponses               string `json:"totalresponses"`
	TotalRequestBytes            string `json:"totalrequestbytes"`
	TotalResponseBytes           string `json:"totalresponsebytes"`
	CurrentClientConnections     string `json:"curclntconnections"`
	SurgeCount                   string `json:"surgecount"`
	CurrentServerConnections     string `json:"cursrvrconnections"`
	ServerEstablishedConnections string `json:"svrestablishedconn"`
	CurrentReusePool             string `json:"curreusepool"`
	MaxClients                   string `json:"maxclients"`
	PrimaryIPAddress             string `json:"primaryipaddress"`
	ServiceGroupName             string `json:"servicegroupname"`
}

// GSLBServiceStats represents the data returned from the /stat/gslbservice Nitro API endpoint
type GSLBServiceStats struct {
	Name                     string `json:"servicename"`
	State                    string `json:"state"`
	TotalRequests            string `json:"totalrequests"`
	TotalResponses           string `json:"totalresponses"`
	TotalRequestBytes        string `json:"totalrequestbytes"`
	TotalResponseBytes       string `json:"totalresponsebytes"`
	CurrentClientConnections string `json:"curclntconnections"`
	CurrentServerConnections string `json:"cursrvrconnections"`
	EstablishedConnections   string `json:"establishedconn"`
	CurrentLoad              string `json:"curload"`
	ServiceHits              string `json:"vsvrservicehits"`
}

// GSLBVirtualServerStats represents the data returned from the /stat/gslbvserver Nitro API endpoint
type GSLBVirtualServerStats struct {
	Name                     string `json:"name"`
	State                    string `json:"state"`
	Health                   string `json:"vslbhealth"`
	InactiveServices         string `json:"inactsvcs"`
	ActiveServices           string `json:"actsvcs"`
	TotalHits                string `json:"tothits"`
	TotalRequests            string `json:"totalrequests"`
	TotalResponses           string `json:"totalresponses"`
	TotalRequestBytes        string `json:"totalrequestbytes"`
	TotalResponseBytes       string `json:"totalresponsebytes"`
	CurrentClientConnections string `json:"curclntconnections"`
	CurrentServerConnections string `json:"cursrvrconnections"`
}

// CSVirtualServerStats represents the data returned from the /stat/csvserver Nitro API endpoint
type CSVirtualServerStats struct {
	Name                          string `json:"name"`
	State                         string `json:"state"`
	TotalHits                     string `json:"tothits"`
	TotalRequests                 string `json:"totalrequests"`
	TotalResponses                string `json:"totalresponses"`
	TotalRequestBytes             string `json:"totalrequestbytes"`
	TotalResponseBytes            string `json:"totalresponsebytes"`
	CurrentClientConnections      string `json:"curclntconnections"`
	CurrentServerConnections      string `json:"cursrvrconnections"`
	EstablishedConnections        string `json:"establishedconn"`
	TotalPacketsReceived          string `json:"totalpktsrecvd"`
	TotalPacketsSent              string `json:"totalpktssent"`
	TotalSpillovers               string `json:"totspillovers"`
	DeferredRequests              string `json:"deferredreq"`
	InvalidRequestResponse        string `json:"invalidrequestresponse"`
	InvalidRequestResponseDropped string `json:"invalidrequestresponsedropped"`
	TotalVServerDownBackupHits    string `json:"totvserverdownbackuphits"`
	CurrentMultipathSessions      string `json:"curmptcpsessions"`
	CurrentMultipathSubflows      string `json:"cursubflowconn"`
}

// VPNVirtualServerStats represents the data returned from the /stat/vpnvserver Nitro API endpoint
type VPNVirtualServerStats struct {
	Name               string `json:"name"`
	TotalRequests      string `json:"totalrequests"`
	TotalResponses     string `json:"totalresponses"`
	TotalRequestBytes  string `json:"totalrequestbytes"`
	TotalResponseBytes string `json:"totalresponsebytes"`
	State              string `json:"state"`
}

// AAAStats represents the data returned from the /stat/aaa Nitro API endpoint
type AAAStats struct {
	AuthSuccess               string `json:"aaaauthsuccess"`
	AuthFail                  string `json:"aaaauthfail"`
	AuthOnlyHTTPSuccess       string `json:"aaaauthonlyhttpsuccess"`
	AuthOnlyHTTPFail          string `json:"aaaauthonlyhttpfail"`
	CurrentIcaSessions        string `json:"aaacuricasessions"`
	CurrentIcaOnlyConnections string `json:"aaacuricaonlyconn"`
}

// LBVServerServiceBinding represents a binding between an LB virtual server and a service.
type LBVServerServiceBinding struct {
	Name        string `json:"name"`
	ServiceName string `json:"servicename"`
	Weight      string `json:"weight"`
}

// LBVServerServiceGroupBinding represents a binding between an LB virtual server and a service group.
type LBVServerServiceGroupBinding struct {
	Name             string `json:"name"`
	ServiceGroupName string `json:"servicegroupname"`
	Weight           string `json:"weight"`
}

// CSVServerLBVServerBinding represents a binding between a CS virtual server and an LB virtual server.
type CSVServerLBVServerBinding struct {
	Name      string `json:"name"`
	LBVServer string `json:"lbvserver"`
	Priority  string `json:"priority"`
}

// BindingsResponse holds the response from binding API calls.
type BindingsResponse struct {
	LBVServerServiceBindings      []LBVServerServiceBinding      `json:"lbvserver_service_binding"`
	LBVServerServiceGroupBindings []LBVServerServiceGroupBinding `json:"lbvserver_servicegroup_binding"`
	CSVServerLBVServerBindings    []CSVServerLBVServerBinding    `json:"csvserver_lbvserver_binding"`
}

// ProtocolHTTPStats represents the data returned from the /stat/protocolhttp Nitro API endpoint
type ProtocolHTTPStats struct {
	// Counters
	TotalRequests                   string `json:"httptotrequests"`
	TotalResponses                  string `json:"httptotresponses"`
	TotalPosts                      string `json:"httptotposts"`
	TotalGets                       string `json:"httptotgets"`
	TotalOthers                     string `json:"httptotothers"`
	TotalRxRequestBytes             string `json:"httptotrxrequestbytes"`
	TotalRxResponseBytes            string `json:"httptotrxresponsebytes"`
	TotalTxRequestBytes             string `json:"httptottxrequestbytes"`
	Total10Requests                 string `json:"httptot10requests"`
	Total11Requests                 string `json:"httptot11requests"`
	Total10Responses                string `json:"httptot10responses"`
	Total11Responses                string `json:"httptot11responses"`
	TotalChunkedRequests            string `json:"httptotchunkedrequests"`
	TotalChunkedResponses           string `json:"httptotchunkedresponses"`
	TotalSPDYStreams                string `json:"spdytotstreams"`
	TotalSPDYv2Streams              string `json:"spdyv2totstreams"`
	TotalSPDYv3Streams              string `json:"spdyv3totstreams"`
	ErrNoReuseMultipart             string `json:"httperrnoreusemultipart"`
	ErrIncompleteHeaders            string `json:"httperrincompleteheaders"`
	ErrIncompleteRequests           string `json:"httperrincompleterequests"`
	ErrIncompleteResponses          string `json:"httperrincompleteresponses"`
	ErrServerBusy                   string `json:"httperrserverbusy"`
	ErrLargeContent                 string `json:"httperrlargecontent"`
	ErrLargeChunk                   string `json:"httperrlargechunk"`
	ErrLargeCtlen                   string `json:"httperrlargectlen"`
	// Gauges (rates) - use FlexString as NetScaler may return number or string
	RequestsRate                    FlexString `json:"httprequestsrate"`
	ResponsesRate                   FlexString `json:"httpresponsesrate"`
	PostsRate                       FlexString `json:"httppostsrate"`
	GetsRate                        FlexString `json:"httpgetsrate"`
	OthersRate                      FlexString `json:"httpothersrate"`
	RxRequestBytesRate              FlexString `json:"httprxrequestbytesrate"`
	RxResponseBytesRate             FlexString `json:"httprxresponsebytesrate"`
	TxRequestBytesRate              FlexString `json:"httptxrequestbytesrate"`
	Request10Rate                   FlexString `json:"http10requestsrate"`
	Request11Rate                   FlexString `json:"http11requestsrate"`
	Response10Rate                  FlexString `json:"http10responsesrate"`
	Response11Rate                  FlexString `json:"http11responsesrate"`
	ChunkedRequestsRate             FlexString `json:"httpchunkedrequestsrate"`
	ChunkedResponsesRate            FlexString `json:"httpchunkedresponsesrate"`
	SPDYStreamsRate                 FlexString `json:"spdystreamsrate"`
	SPDYv2StreamsRate               FlexString `json:"spdyv2streamsrate"`
	SPDYv3StreamsRate               FlexString `json:"spdyv3streamsrate"`
	ErrNoReuseMultipartRate         FlexString `json:"httperrnoreusemultipartrate"`
	ErrIncompleteRequestsRate       FlexString `json:"httperrincompleterequestsrate"`
	ErrIncompleteResponsesRate      FlexString `json:"httperrincompleteresponsesrate"`
	ErrServerBusyRate               FlexString `json:"httperrserverbusyrate"`
}

// ProtocolTCPStats represents the data returned from the /stat/protocoltcp Nitro API endpoint
type ProtocolTCPStats struct {
	// Counters
	TotalRxPackets           string `json:"tcptotrxpkts"`
	TotalRxBytes             string `json:"tcptotrxbytes"`
	TotalTxBytes             string `json:"tcptottxbytes"`
	TotalTxPackets           string `json:"tcptottxpkts"`
	TotalClientConnOpened    string `json:"tcptotclientconnopened"`
	TotalServerConnOpened    string `json:"tcptotserverconnopened"`
	TotalSyn                 string `json:"tcptotsyn"`
	TotalSynProbe            string `json:"tcptotsynprobe"`
	TotalServerFin           string `json:"tcptotsvrfin"`
	TotalClientFin           string `json:"tcptotcltfin"`
	// Gauges - use FlexString for rate fields as NetScaler may return number or string
	ActiveServerConn         string     `json:"tcpactiveserverconn"`
	CurClientConnEstablished string     `json:"tcpcurclientconnestablished"`
	CurServerConnEstablished string     `json:"tcpcurserverconnestablished"`
	RxPacketsRate            FlexString `json:"tcprxpktsrate"`
	RxBytesRate              FlexString `json:"tcprxbytesrate"`
	TxPacketsRate            FlexString `json:"tcptxpktsrate"`
	TxBytesRate              FlexString `json:"tcptxbytesrate"`
	ClientConnOpenedRate     FlexString `json:"tcpclientconnopenedrate"`
	ErrBadChecksum           string     `json:"tcperrbadchecksum"`
	ErrBadChecksumRate       FlexString `json:"tcperrbadchecksumrate"`
	ErrAnyPortFail           string     `json:"tcperranyportfail"`
	ErrIPPortFail            string     `json:"tcperripportfail"`
	ErrBadStateConn          string     `json:"tcperrbadstateconn"`
	ErrRstThreshold          string     `json:"tcperrrstthreshold"`
	SynRate                  FlexString `json:"tcpsynrate"`
	SynProbeRate             FlexString `json:"tcpsynproberate"`
}

// ProtocolIPStats represents the data returned from the /stat/protocolip Nitro API endpoint
type ProtocolIPStats struct {
	// Counters
	TotalRxPackets           string `json:"iptotrxpkts"`
	TotalRxBytes             string `json:"iptotrxbytes"`
	TotalTxPackets           string `json:"iptottxpkts"`
	TotalTxBytes             string `json:"iptottxbytes"`
	TotalRxMbits             string `json:"iptotrxmbits"`
	TotalTxMbits             string `json:"iptottxmbits"`
	TotalRoutedPackets       string `json:"iptotroutedpkts"`
	TotalRoutedMbits         string `json:"iptotroutedmbits"`
	TotalFragments           string `json:"iptotfragments"`
	TotalSuccReassembly      string `json:"iptotsuccreassembly"`
	TotalAddrLookup          string `json:"iptotaddrlookup"`
	TotalAddrLookupFail      string `json:"iptotaddrlookupfail"`
	TotalUDPFragmentsFwd     string `json:"iptotudpfragmentsfwd"`
	TotalTCPFragmentsFwd     string `json:"iptottcpfragmentsfwd"`
	TotalBadChecksums        string `json:"iptotbadchecksums"`
	TotalUnsuccReassembly    string `json:"iptotunsuccreassembly"`
	TotalTooBig              string `json:"iptottoobig"`
	TotalDupFragments        string `json:"iptotdupfragments"`
	TotalOutOfOrderFrag      string `json:"iptotoutoforderfrag"`
	TotalVIPDown             string `json:"iptotvipdown"`
	TotalTTLExpired          string `json:"iptotttlexpired"`
	TotalMaxClients          string `json:"iptotmaxclients"`
	TotalUnknownSvcs         string `json:"iptotunknownsvcs"`
	TotalInvalidHeaderSz     string `json:"iptotinvalidheadersz"`
	TotalInvalidPacketSize   string `json:"iptotinvalidpacketsize"`
	TotalTruncatedPackets    string `json:"iptottruncatedpackets"`
	NonIPTotalTruncatedPkts  string `json:"noniptottruncatedpackets"`
	TotalBadMacAddrs         string `json:"iptotbadmacaddrs"`
	// Gauges (rates) - use FlexString as NetScaler may return number or string
	RxPacketsRate            FlexString `json:"iprxpktsrate"`
	RxBytesRate              FlexString `json:"iprxbytesrate"`
	TxPacketsRate            FlexString `json:"iptxpktsrate"`
	TxBytesRate              FlexString `json:"iptxbytesrate"`
	RxMbitsRate              FlexString `json:"iprxmbitsrate"`
	TxMbitsRate              FlexString `json:"iptxmbitsrate"`
	RoutedPacketsRate        FlexString `json:"iproutedpktsrate"`
	RoutedMbitsRate          FlexString `json:"iproutedmbitsrate"`
}

// SSLStats represents the data returned from the /stat/ssl Nitro API endpoint
type SSLStats struct {
	// Counters
	TotalTLSv11Sessions      string     `json:"ssltottlsv11sessions"`
	TotalSSLv2Sessions       string     `json:"ssltotsslv2sessions"`
	TotalSessions            string     `json:"ssltotsessions"`
	TotalSSLv2Handshakes     string     `json:"ssltotsslv2handshakes"`
	TotalEnc                 string     `json:"ssltotenc"`
	CryptoUtilizationStat    FlexString `json:"sslcryptoutilizationstat"`
	TotalNewSessions         string     `json:"ssltotnewsessions"`
	// Gauges - use FlexString for rate fields as NetScaler may return number or string
	SessionsRate             FlexString `json:"sslsessionsrate"`
	DecRate                  FlexString `json:"ssldecrate"`
	EncRate                  FlexString `json:"sslencrate"`
	SSLv2HandshakesRate      FlexString `json:"sslsslv2handshakesrate"`
	NewSessionsRate          FlexString `json:"sslnewsessionsrate"`
}

// SSLCertKey represents the data returned from the /config/sslcertkey Nitro API endpoint
type SSLCertKey struct {
	CertKey          string     `json:"certkey"`
	DaysToExpiration FlexString `json:"daystoexpiration"`
}

// SSLVServerStats represents the data returned from the /stat/sslvserver Nitro API endpoint
type SSLVServerStats struct {
	VServerName              string `json:"vservername"`
	Type                     string `json:"type"`
	PrimaryIPAddress         string `json:"primaryipaddress"`
	State                    string `json:"state"`
	// Counters
	TotalDecBytes            string `json:"sslctxtotdecbytes"`
	TotalEncBytes            string `json:"sslctxtotencbytes"`
	TotalHWDecBytes          string `json:"sslctxtothwdec_bytes"`
	TotalHWEncBytes          string `json:"sslctxtothwencbytes"`
	TotalSessionNew          string `json:"sslctxtotsessionnew"`
	TotalSessionHits         string `json:"sslctxtotsessionhits"`
	TotalClientAuthSuccess   string `json:"ssltotclientauthsuccess"`
	TotalClientAuthFailure   string `json:"ssltotclientauthfailure"`
	// Gauges - use FlexString for rate fields as NetScaler may return number or string
	Health                   string     `json:"vslbhealth"`
	ActiveServices           string     `json:"actsvcs"`
	ClientAuthSuccessRate    FlexString `json:"sslclientauthsuccessrate"`
	ClientAuthFailureRate    FlexString `json:"sslclientauthfailurerate"`
	EncBytesRate             FlexString `json:"sslctxencbytesrate"`
	DecBytesRate             FlexString `json:"sslctxdecbytesrate"`
	HWEncBytesRate           FlexString `json:"sslctxhwencbytesrate"`
	HWDecBytesRate           FlexString `json:"sslctxhwdec_bytesrate"`
	SessionNewRate           FlexString `json:"sslctxsessionnewrate"`
	SessionHitsRate          FlexString `json:"sslctxsessionhitsrate"`
}

// SystemCPUStats represents the data returned from the /stat/systemcpu Nitro API endpoint
type SystemCPUStats struct {
	ID          string `json:"id"`
	PerCPUUsage string `json:"percpuuse"`
}

// NSCapacityStats represents the data returned from the /stat/nscapacity Nitro API endpoint
type NSCapacityStats struct {
	MaxBandwidth      string `json:"maxbandwidth"`
	MinBandwidth      string `json:"minbandwidth"`
	ActualBandwidth   string `json:"actualbandwidth"`
	Bandwidth         string `json:"bandwidth"`
}

// Bulk binding response types for bulkbindings=yes queries (NS 11.1+)

// BulkLBVServerServiceBindingResponse for lbvserver_service_binding?bulkbindings=yes response
type BulkLBVServerServiceBindingResponse struct {
	LBVServerServiceBindings []LBVServerServiceBinding `json:"lbvserver_service_binding,omitempty"`
}

// BulkLBVServerServiceGroupBindingResponse for lbvserver_servicegroup_binding?bulkbindings=yes response
type BulkLBVServerServiceGroupBindingResponse struct {
	LBVServerServiceGroupBindings []LBVServerServiceGroupBinding `json:"lbvserver_servicegroup_binding,omitempty"`
}

// BulkCSVServerLBVServerBindingResponse for csvserver_lbvserver_binding?bulkbindings=yes response
type BulkCSVServerLBVServerBindingResponse struct {
	CSVServerLBVServerBindings []CSVServerLBVServerBinding `json:"csvserver_lbvserver_binding,omitempty"`
}
