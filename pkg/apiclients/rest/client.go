// Package rest provides manual clients for Snyk's REST API endpoints.
// Unlike the v1 API, REST API uses date-based versioning via query parameter.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sam1el/snyk-api/pkg/client"
)

// DefaultVersion is the default REST API version.
const DefaultVersion = "2025-11-05"

// BaseClient provides common functionality for all REST API clients.
type BaseClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
	version    string
}

// NewBaseClient creates a new REST API base client.
func NewBaseClient(baseClient *client.Client) *BaseClient {
	baseURL := baseClient.RestBaseURL()
	if baseURL == "" {
		baseURL = "https://api.snyk.io/rest"
	}

	version := DefaultVersion
	if v := baseClient.APIVersion(); v != "" {
		version = string(v)
	}

	return &BaseClient{
		httpClient: baseClient.HTTPClient(),
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		token:      baseClient.Token(),
		version:    version,
	}
}

// NewBaseClientWithConfig creates a REST client with explicit configuration.
func NewBaseClientWithConfig(httpClient *http.Client, baseURL, token, version string) *BaseClient {
	if baseURL == "" {
		baseURL = "https://api.snyk.io/rest"
	}
	if version == "" {
		version = DefaultVersion
	}
	return &BaseClient{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		token:      token,
		version:    version,
	}
}

// Request performs an HTTP request to the REST API.
func (c *BaseClient) Request(ctx context.Context, method, path string, query url.Values, body interface{}, result interface{}) error {
	// Build URL with version parameter
	reqURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters including version
	q := reqURL.Query()
	q.Set("version", c.version)
	for k, v := range query {
		for _, val := range v {
			q.Add(k, val)
		}
	}
	reqURL.RawQuery = q.Encode()

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
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - REST API uses token auth and JSON:API content type
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")

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
		var apiErr RESTError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		apiErr.StatusCode = resp.StatusCode
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
func (c *BaseClient) Get(ctx context.Context, path string, query url.Values, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, query, nil, result)
}

// Post performs a POST request.
func (c *BaseClient) Post(ctx context.Context, path string, query url.Values, body, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, query, body, result)
}

// Patch performs a PATCH request.
func (c *BaseClient) Patch(ctx context.Context, path string, query url.Values, body, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, path, query, body, result)
}

// Put performs a PUT request.
func (c *BaseClient) Put(ctx context.Context, path string, query url.Values, body, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, query, body, result)
}

// Delete performs a DELETE request.
func (c *BaseClient) Delete(ctx context.Context, path string, query url.Values) error {
	return c.Request(ctx, http.MethodDelete, path, query, nil, nil)
}

// RESTError represents a REST API error response.
type RESTError struct {
	StatusCode int           `json:"-"`
	Errors     []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail represents a single error in the REST API response.
type ErrorDetail struct {
	ID     string            `json:"id,omitempty"`
	Status string            `json:"status,omitempty"`
	Code   string            `json:"code,omitempty"`
	Title  string            `json:"title,omitempty"`
	Detail string            `json:"detail,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
}

func (e *RESTError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("REST API error: %s - %s", e.Errors[0].Title, e.Errors[0].Detail)
	}
	return fmt.Sprintf("REST API error (status %d)", e.StatusCode)
}

// ============================================================================
// Common Types for JSON:API responses
// ============================================================================

// Links represents JSON:API pagination links.
type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Self  string `json:"self,omitempty"`
}

// Meta represents JSON:API response metadata.
type Meta struct {
	Count int `json:"count,omitempty"`
}

// PaginationParams represents common pagination parameters.
type PaginationParams struct {
	Limit         int    `json:"limit,omitempty"`
	StartingAfter string `json:"starting_after,omitempty"`
	EndingBefore  string `json:"ending_before,omitempty"`
}

// ToQuery converts pagination params to URL query values.
func (p *PaginationParams) ToQuery() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.StartingAfter != "" {
		q.Set("starting_after", p.StartingAfter)
	}
	if p.EndingBefore != "" {
		q.Set("ending_before", p.EndingBefore)
	}
	return q
}
