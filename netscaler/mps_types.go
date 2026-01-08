package netscaler

// MPSHealthStats represents health statistics from a Citrix ADM (MPS) node.
type MPSHealthStats struct {
	NodeType    string `json:"node_type"`
	CPUUsage    string `json:"cpu_usage"`
	DiskUsage   string `json:"disk_usage"`
	DiskFree    string `json:"disk_free"`
	DiskTotal   string `json:"disk_total"`
	DiskUsed    string `json:"disk_used"`
	MemoryUsage string `json:"memory_usage"`
	MemoryFree  string `json:"memory_free"`
	MemoryTotal string `json:"memory_total"`
}

// MPSAPIResponse represents the response from the Citrix ADM Nitro v2 API.
type MPSAPIResponse struct {
	MPSHealth []MPSHealthStats `json:"mps_health"`
}
