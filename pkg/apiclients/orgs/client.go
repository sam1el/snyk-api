// Package orgs provides a high-level client for the Snyk Organizations API.
package orgs

import (
	"context"
	"encoding/json"
	"fmt"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sam1el/snyk-api/pkg/client"
)

// OrgsClient wraps the generated API client with our base client for
// authentication, rate limiting, and error handling.
type OrgsClient struct {
	baseClient *client.Client
	apiClient  *Client
}

// NewOrgsClient creates a new organizations API client.
func NewOrgsClient(baseClient *client.Client) *OrgsClient {
	// Create the generated client using our base client's HTTP client
	apiClient, err := NewClient(
		baseClient.RestBaseURL(),
		WithHTTPClient(baseClient.GetNetworkAccess().GetHttpClient()),
	)
	if err != nil {
		// This should never happen with our setup
		panic(fmt.Sprintf("failed to create orgs API client: %v", err))
	}

	return &OrgsClient{
		baseClient: baseClient,
		apiClient:  apiClient,
	}
}

// ListOrganizations retrieves all organizations the authenticated user has access to.
func (c *OrgsClient) ListOrganizations(ctx context.Context, params *ListOrganizationsParams) (*OrganizationList, error) {
	resp, err := c.apiClient.ListOrganizations(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // Best effort cleanup

	if resp.StatusCode != 200 {
		return nil, client.NewAPIError(resp.StatusCode, resp.Status, "", resp.Header)
	}

	var result OrganizationList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetOrganization retrieves a specific organization by ID.
func (c *OrgsClient) GetOrganization(ctx context.Context, orgID openapi_types.UUID) (*Organization, error) {
	resp, err := c.apiClient.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // Best effort cleanup

	if resp.StatusCode != 200 {
		return nil, client.NewAPIError(resp.StatusCode, resp.Status, "", resp.Header)
	}

	var result struct {
		Data Organization `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}
