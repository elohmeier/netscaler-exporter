package netscaler

import (
	"context"
	"encoding/json"
	"fmt"
)

// getStats is a helper that retrieves and unmarshals stats from the Nitro API.
func getStats(ctx context.Context, c *NitroClient, statsType string, querystring string) (NSAPIResponse, error) {
	data, err := c.GetStats(ctx, statsType, querystring)
	if err != nil {
		return NSAPIResponse{}, err
	}

	var response NSAPIResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return NSAPIResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return response, nil
}

// getConfig is a helper that retrieves and unmarshals config from the Nitro API.
func getConfig(ctx context.Context, c *NitroClient, configType string, querystring string) (NSAPIResponse, error) {
	data, err := c.GetConfig(ctx, configType, querystring)
	if err != nil {
		return NSAPIResponse{}, err
	}

	var response NSAPIResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return NSAPIResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return response, nil
}

// GetNSStats queries the Nitro API for ns stats
func GetNSStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "ns", querystring)
}

// GetInterfaceStats queries the Nitro API for interface stats
func GetInterfaceStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "Interface", querystring)
}

// GetVirtualServerStats queries the Nitro API for virtual server stats
func GetVirtualServerStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "lbvserver", querystring)
}

// GetServiceStats queries the Nitro API for service stats
func GetServiceStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "service", querystring)
}

// GetServiceGroupMemberStats queries the Nitro API for service group member stats.
// Uses the servicegroup/{name}?statbindings=yes endpoint which returns members inline.
func GetServiceGroupMemberStats(ctx context.Context, c *NitroClient, servicegroupName string) (NSAPIResponse, error) {
	return getStats(ctx, c, "servicegroup/"+servicegroupName, "statbindings=yes")
}

// GetGSLBServiceStats queries the Nitro API for GSLB service stats
func GetGSLBServiceStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "gslbservice", querystring)
}

// GetGSLBVirtualServerStats queries the Nitro API for GSLB virtual server stats
func GetGSLBVirtualServerStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "gslbvserver", querystring)
}

// GetCSVirtualServerStats queries the Nitro API for CS virtual server stats
func GetCSVirtualServerStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "csvserver", querystring)
}

// GetVPNVirtualServerStats queries the Nitro API for VPN virtual server stats
func GetVPNVirtualServerStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "vpnvserver", querystring)
}

// GetAAAStats queries the Nitro API for AAA stats
func GetAAAStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "aaa", querystring)
}

// GetNSLicense queries the Nitro API for license config
func GetNSLicense(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getConfig(ctx, c, "nslicense", querystring)
}

// GetServiceGroups queries the Nitro API for service group config
func GetServiceGroups(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getConfig(ctx, c, "servicegroup", querystring)
}

// GetLBVServerServiceBindings retrieves service bindings for a specific LB virtual server.
func GetLBVServerServiceBindings(ctx context.Context, c *NitroClient, lbvserverName string) ([]LBVServerServiceBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_service_binding/"+lbvserverName, "")
	if err != nil {
		return nil, fmt.Errorf("error getting lbvserver_service_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling lbvserver_service_binding response: %w", err)
	}

	return response.LBVServerServiceBindings, nil
}

// GetLBVServerServiceGroupBindings retrieves service group bindings for a specific LB virtual server.
func GetLBVServerServiceGroupBindings(ctx context.Context, c *NitroClient, lbvserverName string) ([]LBVServerServiceGroupBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_servicegroup_binding/"+lbvserverName, "")
	if err != nil {
		return nil, fmt.Errorf("error getting lbvserver_servicegroup_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling lbvserver_servicegroup_binding response: %w", err)
	}

	return response.LBVServerServiceGroupBindings, nil
}

// GetCSVServerLBVServerBindings retrieves LB vserver bindings for a specific CS virtual server.
func GetCSVServerLBVServerBindings(ctx context.Context, c *NitroClient, csvserverName string) ([]CSVServerLBVServerBinding, error) {
	body, err := c.GetConfig(ctx, "csvserver_lbvserver_binding/"+csvserverName, "")
	if err != nil {
		return nil, fmt.Errorf("error getting csvserver_lbvserver_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling csvserver_lbvserver_binding response: %w", err)
	}

	return response.CSVServerLBVServerBindings, nil
}

// Bulk binding functions using bulkbindings=yes (NS 11.1+)
// These fetch all bindings in a single API call instead of per-vserver queries.

// GetAllLBVServerServiceBindings retrieves all service bindings for all LB vservers in one call.
func GetAllLBVServerServiceBindings(ctx context.Context, c *NitroClient) ([]LBVServerServiceBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_service_binding", "bulkbindings=yes")
	if err != nil {
		return nil, fmt.Errorf("error getting bulk lbvserver_service_binding: %w", err)
	}

	var response BulkLBVServerServiceBindingResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling bulk lbvserver_service_binding: %w", err)
	}

	return response.LBVServerServiceBindings, nil
}

// GetAllLBVServerServiceGroupBindings retrieves all service group bindings for all LB vservers in one call.
func GetAllLBVServerServiceGroupBindings(ctx context.Context, c *NitroClient) ([]LBVServerServiceGroupBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_servicegroup_binding", "bulkbindings=yes")
	if err != nil {
		return nil, fmt.Errorf("error getting bulk lbvserver_servicegroup_binding: %w", err)
	}

	var response BulkLBVServerServiceGroupBindingResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling bulk lbvserver_servicegroup_binding: %w", err)
	}

	return response.LBVServerServiceGroupBindings, nil
}

// GetAllCSVServerLBVServerBindings retrieves all LB vserver bindings for all CS vservers in one call.
func GetAllCSVServerLBVServerBindings(ctx context.Context, c *NitroClient) ([]CSVServerLBVServerBinding, error) {
	body, err := c.GetConfig(ctx, "csvserver_lbvserver_binding", "bulkbindings=yes")
	if err != nil {
		return nil, fmt.Errorf("error getting bulk csvserver_lbvserver_binding: %w", err)
	}

	var response BulkCSVServerLBVServerBindingResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling bulk csvserver_lbvserver_binding: %w", err)
	}

	return response.CSVServerLBVServerBindings, nil
}

// GetProtocolHTTPStats queries the Nitro API for protocol HTTP stats
func GetProtocolHTTPStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "protocolhttp", querystring)
}

// GetProtocolTCPStats queries the Nitro API for protocol TCP stats
func GetProtocolTCPStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "protocoltcp", querystring)
}

// GetProtocolIPStats queries the Nitro API for protocol IP stats
func GetProtocolIPStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "protocolip", querystring)
}

// GetSSLStats queries the Nitro API for SSL stats
func GetSSLStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "ssl", querystring)
}

// GetSSLCertKeys queries the Nitro API for SSL certificate keys
func GetSSLCertKeys(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getConfig(ctx, c, "sslcertkey", querystring)
}

// GetSSLVServerStats queries the Nitro API for SSL virtual server stats
func GetSSLVServerStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "sslvserver", querystring)
}

// GetSystemCPUStats queries the Nitro API for system CPU stats
func GetSystemCPUStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "systemcpu", querystring)
}

// GetNSCapacityStats queries the Nitro API for bandwidth capacity stats
func GetNSCapacityStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "nscapacity", querystring)
}
