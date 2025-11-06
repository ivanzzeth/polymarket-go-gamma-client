package polymarketgamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetTeams fetches teams with optional filtering and pagination
// Reference: https://gamma-api.polymarket.com/teams
func (c *Client) GetTeams(ctx context.Context, params *GetTeamsParams) ([]Team, error) {
	path := "/teams?"

	urlParams := url.Values{}

	if params != nil {
		if params.Limit > 0 {
			urlParams.Add("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			urlParams.Add("offset", fmt.Sprintf("%d", params.Offset))
		}
		if params.Order != "" {
			urlParams.Add("order", params.Order)
		}
		if params.Ascending != nil {
			urlParams.Add("ascending", fmt.Sprintf("%t", *params.Ascending))
		}
		for _, league := range params.League {
			urlParams.Add("league", league)
		}
		for _, name := range params.Name {
			urlParams.Add("name", name)
		}
		for _, abbreviation := range params.Abbreviation {
			urlParams.Add("abbreviation", abbreviation)
		}
	}

	path += urlParams.Encode()

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var teams []Team
	if err := json.Unmarshal(respBody, &teams); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return teams, nil
}

// GetSportsMetadata retrieves metadata for various sports including images, resolution sources,
// ordering preferences, tags, and series information
// Reference: https://gamma-api.polymarket.com/sports
func (c *Client) GetSportsMetadata(ctx context.Context) ([]SportMetadata, error) {
	path := "/sports"

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var metadata []SportMetadata
	if err := json.Unmarshal(respBody, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return metadata, nil
}
