package projects

import (
	"context"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sam1el/snyk-api/pkg/client"
)

// ProjectsClient wraps the generated OpenAPI client for the Projects API.
type ProjectsClient struct {
	apiClient  *ClientWithResponses
	baseClient *client.Client
}

// NewProjectsClient creates a new ProjectsClient.
func NewProjectsClient(baseClient *client.Client) *ProjectsClient {
	// The generated client needs an http.Client that uses our baseClient's Execute method.
	httpClient := &http.Client{
		Transport: &roundTripperFunc{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return baseClient.Execute(req.Context(), req)
			},
		},
	}

	// Create the generated client
	apiClient, err := NewClientWithResponses(baseClient.BaseURL(), WithHTTPClient(httpClient))
	if err != nil {
		// Fatal error - can't proceed without client
		panic(fmt.Sprintf("Failed to create generated projects API client: %v", err))
	}

	return &ProjectsClient{
		apiClient:  apiClient,
		baseClient: baseClient,
	}
}

// ListProjects lists all projects in an organization.
func (c *ProjectsClient) ListProjects(ctx context.Context, orgID string, params *ListProjectsParams) (*ProjectList, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	resp, err := c.apiClient.ListProjectsWithResponse(ctx, orgUUID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	if resp.ApplicationvndApiJSON200 == nil {
		return nil, fmt.Errorf("unexpected response from API: status %d, body: %s", resp.HTTPResponse.StatusCode, string(resp.Body))
	}
	return resp.ApplicationvndApiJSON200, nil
}

// GetProject retrieves a single project by ID.
func (c *ProjectsClient) GetProject(ctx context.Context, orgID, projectID string) (*Project, error) {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	projectUUID, err := parseUUID(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	resp, err := c.apiClient.GetProjectWithResponse(ctx, orgUUID, projectUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", projectID, err)
	}
	if resp.ApplicationvndApiJSON200 == nil {
		return nil, fmt.Errorf("unexpected response from API: status %d, body: %s", resp.HTTPResponse.StatusCode, string(resp.Body))
	}
	return &resp.ApplicationvndApiJSON200.Data, nil
}

// DeleteProject deletes a project.
func (c *ProjectsClient) DeleteProject(ctx context.Context, orgID, projectID string) error {
	orgUUID, err := parseUUID(orgID)
	if err != nil {
		return fmt.Errorf("invalid organization ID: %w", err)
	}

	projectUUID, err := parseUUID(projectID)
	if err != nil {
		return fmt.Errorf("invalid project ID: %w", err)
	}

	resp, err := c.apiClient.DeleteProjectWithResponse(ctx, orgUUID, projectUUID)
	if err != nil {
		return fmt.Errorf("failed to delete project %s: %w", projectID, err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response from API: status %d, body: %s", resp.HTTPResponse.StatusCode, string(resp.Body))
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
