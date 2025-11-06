package polymarketgamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetMarketByID fetches a specific market by its market ID (numeric ID)
func (c *Client) GetMarketByID(ctx context.Context, marketID string, params *GetMarketByIDQueryParams) (*Market, error) {
	path := fmt.Sprintf("/markets/%s", url.PathEscape(marketID))

	if params != nil && params.IncludeTag != nil {
		urlParams := url.Values{}
		urlParams.Add("include_tag", fmt.Sprintf("%t", *params.IncludeTag))
		path += "?" + urlParams.Encode()
	}

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var market Market
	if err := json.Unmarshal(respBody, &market); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &market, nil
}

// GetMarkets fetches markets with optional filtering and pagination
// Reference: https://docs.polymarket.com/api-reference/markets/list-markets
func (c *Client) GetMarkets(ctx context.Context, params *GetMarketsParams) ([]*Market, error) {
	path := "/markets?"

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
		for _, id := range params.ID {
			urlParams.Add("id", fmt.Sprintf("%d", id))
		}
		for _, slug := range params.Slug {
			urlParams.Add("slug", slug)
		}
		for _, tokenID := range params.ClobTokenIDs {
			urlParams.Add("clob_token_ids", tokenID)
		}
		for _, conditionID := range params.ConditionIDs {
			urlParams.Add("condition_ids", conditionID)
		}
		for _, addr := range params.MarketMakerAddress {
			urlParams.Add("market_maker_address", addr)
		}
		if params.LiquidityNumMin != nil {
			urlParams.Add("liquidity_num_min", fmt.Sprintf("%f", *params.LiquidityNumMin))
		}
		if params.LiquidityNumMax != nil {
			urlParams.Add("liquidity_num_max", fmt.Sprintf("%f", *params.LiquidityNumMax))
		}
		if params.VolumeNumMin != nil {
			urlParams.Add("volume_num_min", fmt.Sprintf("%f", *params.VolumeNumMin))
		}
		if params.VolumeNumMax != nil {
			urlParams.Add("volume_num_max", fmt.Sprintf("%f", *params.VolumeNumMax))
		}
		if params.StartDateMin != nil {
			urlParams.Add("start_date_min", params.StartDateMin.Time().Format("2006-01-02T15:04:05Z07:00"))
		}
		if params.StartDateMax != nil {
			urlParams.Add("start_date_max", params.StartDateMax.Time().Format("2006-01-02T15:04:05Z07:00"))
		}
		if params.EndDateMin != nil {
			urlParams.Add("end_date_min", params.EndDateMin.Time().Format("2006-01-02T15:04:05Z07:00"))
		}
		if params.EndDateMax != nil {
			urlParams.Add("end_date_max", params.EndDateMax.Time().Format("2006-01-02T15:04:05Z07:00"))
		}
		if params.TagID != nil {
			urlParams.Add("tag_id", fmt.Sprintf("%d", *params.TagID))
		}
		if params.RelatedTags != nil {
			urlParams.Add("related_tags", fmt.Sprintf("%t", *params.RelatedTags))
		}
		if params.CYOM != nil {
			urlParams.Add("cyom", fmt.Sprintf("%t", *params.CYOM))
		}
		if params.UMAResolutionStatus != "" {
			urlParams.Add("uma_resolution_status", params.UMAResolutionStatus)
		}
		if params.GameID != "" {
			urlParams.Add("game_id", params.GameID)
		}
		for _, smt := range params.SportsMarketTypes {
			urlParams.Add("sports_market_types", smt)
		}
		if params.RewardsMinSize != nil {
			urlParams.Add("rewards_min_size", fmt.Sprintf("%f", *params.RewardsMinSize))
		}
		for _, qid := range params.QuestionIDs {
			urlParams.Add("question_ids", qid)
		}
		if params.IncludeTag != nil {
			urlParams.Add("include_tag", fmt.Sprintf("%t", *params.IncludeTag))
		}
		if params.Closed != nil {
			urlParams.Add("closed", fmt.Sprintf("%t", *params.Closed))
		}
	}

	path += urlParams.Encode()

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var markets []*Market
	if err := json.Unmarshal(respBody, &markets); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return markets, nil
}

// GetMarketTags fetches all tags associated with a specific market
func (c *Client) GetMarketTags(ctx context.Context, marketID string) ([]Tag, error) {
	path := fmt.Sprintf("/markets/%s/tags", url.PathEscape(marketID))

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	if err := json.Unmarshal(respBody, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tags, nil
}

// GetMarketBySlug fetches a specific market by its slug
func (c *Client) GetMarketBySlug(ctx context.Context, slug string, params *GetMarketByIDQueryParams) (*Market, error) {
	path := fmt.Sprintf("/markets/slug/%s", url.PathEscape(slug))

	if params != nil && params.IncludeTag != nil {
		urlParams := url.Values{}
		urlParams.Add("include_tag", fmt.Sprintf("%t", *params.IncludeTag))
		path += "?" + urlParams.Encode()
	}

	respBody, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var market Market
	if err := json.Unmarshal(respBody, &market); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &market, nil
}
