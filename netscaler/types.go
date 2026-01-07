package netscaler

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
