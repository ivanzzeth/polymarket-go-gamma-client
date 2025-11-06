package polymarketgamma

import (
	"context"
	"encoding/json"
	"fmt"
)

// HealthCheck checks if the Gamma API is healthy and responding
// Returns the status string (typically "OK") if successful
func (c *Client) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	respBody, err := c.doRequest(ctx, "GET", "/")
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	var health HealthResponse
	if err := json.Unmarshal(respBody, &health); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &health, nil
}
