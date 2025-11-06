package polymarketgamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Search searches for markets, events, and profiles
func (c *Client) Search(ctx context.Context, params *SearchParams) (*SearchResponse, error) {
	if params == nil || params.Q == "" {
		return nil, fmt.Errorf("search query (q) is required")
	}

	path := "/public-search?"
	urlParams := url.Values{}

	// Required parameter
	urlParams.Add("q", params.Q)

	// Optional parameters
	if params.Cache != nil {
		urlParams.Add("cache", fmt.Sprintf("%t", *params.Cache))
	}
	if params.EventsStatus != "" {
		urlParams.Add("events_status", params.EventsStatus)
	}
	if params.LimitPerType != nil {
		urlParams.Add("limit_per_type", fmt.Sprintf("%d", *params.LimitPerType))
	}
	if params.Page != nil {
		urlParams.Add("page", fmt.Sprintf("%d", *params.Page))
	}
	for _, tag := range params.EventsTag {
		urlParams.Add("events_tag", tag)
	}
	if params.KeepClosedMarkets != nil {
		urlParams.Add("keep_closed_markets", fmt.Sprintf("%d", *params.KeepClosedMarkets))
	}
	if params.Sort != "" {
		urlParams.Add("sort", params.Sort)
	}
	if params.Ascending != nil {
		urlParams.Add("ascending", fmt.Sprintf("%t", *params.Ascending))
	}
	if params.SearchTags != nil {
		urlParams.Add("search_tags", fmt.Sprintf("%t", *params.SearchTags))
	}
	if params.SearchProfiles != nil {
		urlParams.Add("search_profiles", fmt.Sprintf("%t", *params.SearchProfiles))
	}
	if params.Recurrence != "" {
		urlParams.Add("recurrence", params.Recurrence)
	}
	for _, tagID := range params.ExcludeTagID {
		urlParams.Add("exclude_tag_id", fmt.Sprintf("%d", tagID))
	}
	if params.Optimized != nil {
		urlParams.Add("optimized", fmt.Sprintf("%t", *params.Optimized))
	}

	path += urlParams.Encode()

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var response SearchResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
