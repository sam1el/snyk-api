// Package v1 provides clients for Snyk's v1 API endpoints.
// The v1 API uses a different base URL (api.snyk.io/v1) and authentication
// pattern compared to the REST API.
package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sam1el/snyk-api/pkg/client"
)

// BaseClient provides common functionality for all v1 API clients.
type BaseClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// NewBaseClient creates a new v1 API base client.
func NewBaseClient(baseClient *client.Client) *BaseClient {
	// Get the base URL and adjust for v1
	baseURL := baseClient.BaseURL()
	if baseURL == "" {
		baseURL = "https://api.snyk.io"
	}
	// Ensure we use /v1 path
	baseURL = baseURL + "/v1"

	return &BaseClient{
		httpClient: baseClient.HTTPClient(),
		baseURL:    baseURL,
		token:      baseClient.Token(),
	}
}

// NewBaseClientWithConfig creates a v1 client with explicit configuration.
func NewBaseClientWithConfig(httpClient *http.Client, baseURL, token string) *BaseClient {
	if baseURL == "" {
		baseURL = "https://api.snyk.io"
	}
	return &BaseClient{
		httpClient: httpClient,
		baseURL:    baseURL + "/v1",
		token:      token,
	}
}

// Request performs an HTTP request to the v1 API.
func (c *BaseClient) Request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build URL
	reqURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - v1 API uses token auth
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr V1Error
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return &apiErr
	}

	// Parse result
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request.
func (c *BaseClient) Get(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request.
func (c *BaseClient) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request.
func (c *BaseClient) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, body, result)
}

// Delete performs a DELETE request.
func (c *BaseClient) Delete(ctx context.Context, path string) error {
	return c.Request(ctx, http.MethodDelete, path, nil, nil)
}

// V1Error represents a v1 API error response.
type V1Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error_  string `json:"error,omitempty"`
}

func (e *V1Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Error_ != "" {
		return e.Error_
	}
	return fmt.Sprintf("v1 API error (code: %d)", e.Code)
}
