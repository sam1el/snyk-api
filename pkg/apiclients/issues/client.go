package issues

import (
	"context"
	"fmt"
	"io"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/sam1el/snyk-api/pkg/client"
)

// IssuesClient wraps the generated OpenAPI client for the Issues API.
type IssuesClient struct {
	apiClient  *ClientWithResponses
	baseClient *client.Client
}

// NewIssuesClient creates a new IssuesClient.
func NewIssuesClient(baseClient *client.Client) *IssuesClient {
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
		baseClient.GetLogger().Fatal().Err(err).Msg("Failed to create generated issues API client")
	}

	return &IssuesClient{
		apiClient:  apiClient,
		baseClient: baseClient,
	}
}

// ListIssuesForOrg lists all issues for a given organization.
func (c *IssuesClient) ListIssuesForOrg(ctx context.Context, orgID string, params *ListIssuesForOrgParams) (*IssueList, error) {
	uuidOrgID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}

	resp, err := c.apiClient.ListIssuesForOrgWithResponse(ctx, uuidOrgID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues for organization %s: %w", orgID, err)
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

// GetIssue retrieves a single issue by ID.
func (c *IssuesClient) GetIssue(ctx context.Context, orgID, issueID string) (*IssueResponse, error) {
	uuidOrgID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID format: %w", err)
	}
	uuidIssueID, err := parseUUID(issueID)
	if err != nil {
		return nil, fmt.Errorf("invalid issue ID format: %w", err)
	}

	resp, err := c.apiClient.GetIssueWithResponse(ctx, uuidOrgID, uuidIssueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %s for organization %s: %w", issueID, orgID, err)
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
