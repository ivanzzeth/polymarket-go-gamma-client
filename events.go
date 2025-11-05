package polymarketgamma

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetEvents fetches all events with optional filtering
func (c *Client) GetEvents(params *GetEventsParams) ([]Event, error) {
	path := "/events?"

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
		if params.TagID != nil {
			urlParams.Add("tag_id", fmt.Sprintf("%d", *params.TagID))
		}
		for _, tagID := range params.ExcludeTagID {
			urlParams.Add("exclude_tag_id", fmt.Sprintf("%d", tagID))
		}
		if params.RelatedTags != nil {
			urlParams.Add("related_tags", fmt.Sprintf("%t", *params.RelatedTags))
		}
		if params.Featured != nil {
			urlParams.Add("featured", fmt.Sprintf("%t", *params.Featured))
		}
		if params.CYOM != nil {
			urlParams.Add("cyom", fmt.Sprintf("%t", *params.CYOM))
		}
		if params.IncludeChat != nil {
			urlParams.Add("include_chat", fmt.Sprintf("%t", *params.IncludeChat))
		}
		if params.IncludeTemplate != nil {
			urlParams.Add("include_template", fmt.Sprintf("%t", *params.IncludeTemplate))
		}
		if params.Recurrence != "" {
			urlParams.Add("recurrence", params.Recurrence)
		}
		if params.Closed != nil {
			urlParams.Add("closed", fmt.Sprintf("%t", *params.Closed))
		}
		if params.StartDateMin != "" {
			urlParams.Add("start_date_min", params.StartDateMin)
		}
		if params.StartDateMax != "" {
			urlParams.Add("start_date_max", params.StartDateMax)
		}
		if params.EndDateMin != "" {
			urlParams.Add("end_date_min", params.EndDateMin)
		}
		if params.EndDateMax != "" {
			urlParams.Add("end_date_max", params.EndDateMax)
		}
	}

	path += urlParams.Encode()

	respBody, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := json.Unmarshal(respBody, &events); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return events, nil
}

// GetEventBySlug fetches a specific event by its slug using /events/slug/{slug} with optional query parameters
func (c *Client) GetEventBySlug(slug string, params *GetEventBySlugQueryParams) (*Event, error) {
	path := fmt.Sprintf("/events/slug/%s", url.PathEscape(slug))

	// Add query parameters if provided
	if params != nil {
		urlParams := url.Values{}
		if params.IncludeChat != nil {
			urlParams.Add("include_chat", fmt.Sprintf("%t", *params.IncludeChat))
		}
		if params.IncludeTemplate != nil {
			urlParams.Add("include_template", fmt.Sprintf("%t", *params.IncludeTemplate))
		}
		if len(urlParams) > 0 {
			path += "?" + urlParams.Encode()
		}
	}

	respBody, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var event Event
	if err := json.Unmarshal(respBody, &event); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &event, nil
}

// GetEventByID fetches a specific event by its ID with optional query parameters
func (c *Client) GetEventByID(eventID string, params *GetEventByIDQueryParams) (*Event, error) {
	path := fmt.Sprintf("/events/%s", url.PathEscape(eventID))

	// Add query parameters if provided
	if params != nil {
		urlParams := url.Values{}
		if params.IncludeChat != nil {
			urlParams.Add("include_chat", fmt.Sprintf("%t", *params.IncludeChat))
		}
		if params.IncludeTemplate != nil {
			urlParams.Add("include_template", fmt.Sprintf("%t", *params.IncludeTemplate))
		}
		if len(urlParams) > 0 {
			path += "?" + urlParams.Encode()
		}
	}

	respBody, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var event Event
	if err := json.Unmarshal(respBody, &event); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &event, nil
}

// GetEventTags fetches all tags associated with a specific event
func (c *Client) GetEventTags(eventID string) ([]Tag, error) {
	path := fmt.Sprintf("/events/%s/tags", url.PathEscape(eventID))

	respBody, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	if err := json.Unmarshal(respBody, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tags, nil
}
