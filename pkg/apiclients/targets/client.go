package targets

import (
	"context"
	"fmt"
	"io"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/sam1el/snyk-api/pkg/client"
)

// TargetsClient wraps the generated OpenAPI client for the Targets API.
type TargetsClient struct {
	apiClient  *ClientWithResponses
	baseClient *client.Client
}

// NewTargetsClient creates a new TargetsClient.
func NewTargetsClient(baseClient *client.Client) *TargetsClient {
	// The generated client needs an http.Client that uses our baseClient's Execute method.
	httpClient := &http.Client{
		Transport: &roundTripperFunc{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return baseClient.Execute(req.Context(), req)
			},
		},
	}

	// Create the generated client
	apiClient, err := NewClientWithResponses(baseClient.RestBaseURL(), WithHTTPClient(httpClient))
	if err != nil {
		baseClient.GetLogger().Fatal().Err(err).Msg("Failed to create generated targets API client")
	}

	return &TargetsClient{
		apiClient:  apiClient,
		baseClient: baseClient,
	}
}

// ListTargets lists all targets for a given organization.
func (c *TargetsClient) ListTargets(ctx context.Context, orgID string, params *ListTargetsParams) (*TargetList, error) {
	uuidOrgID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	resp, err := c.apiClient.ListTargetsWithResponse(ctx, uuidOrgID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list targets for organization %s: %w", orgID, err)
	}
	if resp.ApplicationvndApiJSON200 == nil {
		if resp.Body != nil {
			defer func() { _ = resp.HTTPResponse.Body.Close() }()
			bodyBytes, err := io.ReadAll(resp.HTTPResponse.Body)
			if err != nil {
				return nil, fmt.Errorf("unexpected response from API: status %s, failed to read body: %w", resp.Status(), err)
			}
			return nil, fmt.Errorf("unexpected response from API: status %s, body: %s", resp.Status(), string(bodyBytes))
		}
		return nil, fmt.Errorf("unexpected response from API: status %s", resp.Status())
	}
	return resp.ApplicationvndApiJSON200, nil
}

// GetTarget retrieves a single target by ID.
func (c *TargetsClient) GetTarget(ctx context.Context, orgID, targetID string) (*TargetResponse, error) {
	uuidOrgID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}
	uuidTargetID, err := parseUUID(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target ID format: %w", err)
	}

	resp, err := c.apiClient.GetTargetWithResponse(ctx, uuidOrgID, uuidTargetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target %s for organization %s: %w", targetID, orgID, err)
	}
	if resp.ApplicationvndApiJSON200 == nil {
		if resp.Body != nil {
			defer func() { _ = resp.HTTPResponse.Body.Close() }()
			bodyBytes, err := io.ReadAll(resp.HTTPResponse.Body)
			if err != nil {
				return nil, fmt.Errorf("unexpected response from API: status %s, failed to read body: %w", resp.Status(), err)
			}
			return nil, fmt.Errorf("unexpected response from API: status %s, body: %s", resp.Status(), string(bodyBytes))
		}
		return nil, fmt.Errorf("unexpected response from API: status %s", resp.Status())
	}
	return resp.ApplicationvndApiJSON200, nil
}

// DeleteTarget deletes a target by ID.
func (c *TargetsClient) DeleteTarget(ctx context.Context, orgID, targetID string) error {
	uuidOrgID, err := parseUUID(orgID)
	if err != nil {
		return fmt.Errorf("invalid organization ID format: %w", err)
	}
	uuidTargetID, err := parseUUID(targetID)
	if err != nil {
		return fmt.Errorf("invalid target ID format: %w", err)
	}

	// Create the required request body
	requestBody := DeleteTargetRequest{
		Data: struct {
			Type DeleteTargetRequestDataType `json:"type"`
		}{
			Type: "target",
		},
	}

	resp, err := c.apiClient.DeleteTargetWithApplicationVndAPIPlusJSONBodyWithResponse(ctx, uuidOrgID, uuidTargetID, requestBody)
	if err != nil {
		return fmt.Errorf("failed to delete target %s for organization %s: %w", targetID, orgID, err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusNoContent {
		if resp.Body != nil {
			defer func() { _ = resp.HTTPResponse.Body.Close() }()
			bodyBytes, err := io.ReadAll(resp.HTTPResponse.Body)
			if err != nil {
				return fmt.Errorf("unexpected response from API: status %s, failed to read body: %w", resp.Status(), err)
			}
			return fmt.Errorf("unexpected response from API: status %s, body: %s", resp.Status(), string(bodyBytes))
		}
		return fmt.Errorf("unexpected response from API: status %s", resp.Status())
	}
	return nil
}

// parseUUID parses a string as a UUID.
func parseUUID(s string) (openapi_types.UUID, error) {
	var uuid openapi_types.UUID
	err := uuid.UnmarshalText([]byte(s))
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return uuid, nil
}

// roundTripperFunc is a helper to allow a function to implement http.RoundTripper.
type roundTripperFunc struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (rt *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.roundTrip(req)
}
