package netscaler

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetMPSHealth queries the Citrix ADM Nitro v2 API for mps_health stats.
func GetMPSHealth(ctx context.Context, c *MPSClient) (MPSAPIResponse, error) {
	data, err := c.GetStats(ctx, "mps_health", "")
	if err != nil {
		return MPSAPIResponse{}, err
	}

	var response MPSAPIResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return MPSAPIResponse{}, fmt.Errorf("error unmarshalling mps_health response: %w", err)
	}

	return response, nil
}
