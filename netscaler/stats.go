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

// GetServiceGroupMemberStats queries the Nitro API for service group member stats
func GetServiceGroupMemberStats(ctx context.Context, c *NitroClient, querystring string) (NSAPIResponse, error) {
	return getStats(ctx, c, "servicegroupmember", querystring)
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

// GetLBVServerServiceBindings retrieves all LB virtual server to service bindings.
func GetLBVServerServiceBindings(ctx context.Context, c *NitroClient) ([]LBVServerServiceBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_service_binding", "")
	if err != nil {
		return nil, fmt.Errorf("error getting lbvserver_service_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling lbvserver_service_binding response: %w", err)
	}

	return response.LBVServerServiceBindings, nil
}

// GetLBVServerServiceGroupBindings retrieves all LB virtual server to service group bindings.
func GetLBVServerServiceGroupBindings(ctx context.Context, c *NitroClient) ([]LBVServerServiceGroupBinding, error) {
	body, err := c.GetConfig(ctx, "lbvserver_servicegroup_binding", "")
	if err != nil {
		return nil, fmt.Errorf("error getting lbvserver_servicegroup_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling lbvserver_servicegroup_binding response: %w", err)
	}

	return response.LBVServerServiceGroupBindings, nil
}

// GetCSVServerLBVServerBindings retrieves all CS virtual server to LB virtual server bindings.
func GetCSVServerLBVServerBindings(ctx context.Context, c *NitroClient) ([]CSVServerLBVServerBinding, error) {
	body, err := c.GetConfig(ctx, "csvserver_lbvserver_binding", "")
	if err != nil {
		return nil, fmt.Errorf("error getting csvserver_lbvserver_binding: %w", err)
	}

	var response BindingsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshalling csvserver_lbvserver_binding response: %w", err)
	}

	return response.CSVServerLBVServerBindings, nil
}
