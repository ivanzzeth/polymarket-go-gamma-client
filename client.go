package polymarketgamma

import (
	"fmt"
	"io"
	"net/http"
)

// Client is a client for the Gamma API (events and markets metadata)
type Client struct {
	host       string
	httpClient *http.Client
}

// NewClient creates a new Gamma API client for querying events and market metadata
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		host:       GAMMA_API_URL,
		httpClient: httpClient,
	}
}

// doRequest performs an HTTP request to the Gamma API
func (c *Client) doRequest(method, path string) ([]byte, error) {
	fullURL := c.host + path

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
