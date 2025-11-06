package polymarketgamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetSeries fetches all series with optional filtering
func (c *Client) GetSeries(ctx context.Context, params *GetSeriesParams) ([]Series, error) {
	path := "/series?"

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
		for _, slug := range params.Slug {
			urlParams.Add("slug", slug)
		}
		for _, id := range params.CategoriesIDs {
			urlParams.Add("categories_ids", fmt.Sprintf("%d", id))
		}
		for _, label := range params.CategoriesLabels {
			urlParams.Add("categories_labels", label)
		}
		if params.Closed != nil {
			urlParams.Add("closed", fmt.Sprintf("%t", *params.Closed))
		}
		if params.IncludeChat != nil {
			urlParams.Add("include_chat", fmt.Sprintf("%t", *params.IncludeChat))
		}
		if params.Recurrence != "" {
			urlParams.Add("recurrence", params.Recurrence)
		}
	}

	path += urlParams.Encode()

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var series []Series
	if err := json.Unmarshal(respBody, &series); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return series, nil
}

// GetSeriesByID fetches a specific series by its ID
func (c *Client) GetSeriesByID(ctx context.Context, seriesID string, params *GetSeriesByIDQueryParams) (*Series, error) {
	path := fmt.Sprintf("/series/%s", url.PathEscape(seriesID))

	if params != nil && params.IncludeChat != nil {
		urlParams := url.Values{}
		urlParams.Add("include_chat", fmt.Sprintf("%t", *params.IncludeChat))
		path += "?" + urlParams.Encode()
	}

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var series Series
	if err := json.Unmarshal(respBody, &series); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &series, nil
}
