// Package integrations provides a client for Snyk v1 Integrations API.
package integrations

import (
	"context"
	"fmt"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Integrations API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Integrations client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// Integration represents a Snyk integration.
type Integration struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
}

// IntegrationSettings represents integration settings.
type IntegrationSettings struct {
	AutoDepUpgradeEnabled              *bool    `json:"autoDepUpgradeEnabled,omitempty"`
	AutoDepUpgradeIgnoredDependencies  []string `json:"autoDepUpgradeIgnoredDependencies,omitempty"`
	AutoDepUpgradeMinAge               *int     `json:"autoDepUpgradeMinAge,omitempty"`
	AutoDepUpgradeLimit                *int     `json:"autoDepUpgradeLimit,omitempty"`
	PullRequestFailOnAnyVulns          *bool    `json:"pullRequestFailOnAnyVulns,omitempty"`
	PullRequestFailOnlyForHighSeverity *bool    `json:"pullRequestFailOnlyForHighSeverity,omitempty"`
	PullRequestTestEnabled             *bool    `json:"pullRequestTestEnabled,omitempty"`
	AutoPrEnabled                      *bool    `json:"autoPrEnabled,omitempty"`
}

// ============================================================================
// List Integrations
// ============================================================================

// Integrations represents a map of integration type to integration.
type Integrations map[string]Integration

// List lists all integrations for an organization.
// GET /org/{orgId}/integrations
func (c *Client) List(ctx context.Context, orgID string) (Integrations, error) {
	path := fmt.Sprintf("/org/%s/integrations", orgID)
	var result Integrations
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ============================================================================
// Get/Update Integration
// ============================================================================

// Get retrieves an integration by ID.
// GET /org/{orgId}/integrations/{integrationId}
func (c *Client) Get(ctx context.Context, orgID, integrationID string) (*Integration, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s", orgID, integrationID)
	var result Integration
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByType retrieves an integration by type.
// GET /org/{orgId}/integrations/{type}
func (c *Client) GetByType(ctx context.Context, orgID, integrationType string) (*Integration, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s", orgID, integrationType)
	var result Integration
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Settings
// ============================================================================

// GetSettings retrieves integration settings.
// GET /org/{orgId}/integrations/{integrationId}/settings
func (c *Client) GetSettings(ctx context.Context, orgID, integrationID string) (*IntegrationSettings, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/settings", orgID, integrationID)
	var result IntegrationSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates integration settings.
// PUT /org/{orgId}/integrations/{integrationId}/settings
func (c *Client) UpdateSettings(ctx context.Context, orgID, integrationID string, settings *IntegrationSettings) (*IntegrationSettings, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/settings", orgID, integrationID)
	var result IntegrationSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Authentication
// ============================================================================

// Credentials represents integration credentials.
type Credentials struct {
	Token        string `json:"token,omitempty"`
	URL          string `json:"url,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Region       string `json:"region,omitempty"`
	RegistryBase string `json:"registryBase,omitempty"`
}

// UpdateAuthentication updates integration authentication.
// PUT /org/{orgId}/integrations/{integrationId}/authentication
func (c *Client) UpdateAuthentication(ctx context.Context, orgID, integrationID string, creds *Credentials) error {
	path := fmt.Sprintf("/org/%s/integrations/%s/authentication", orgID, integrationID)
	return c.base.Put(ctx, path, creds, nil)
}

// ProvisionBrokerToken provisions a broker token.
// POST /org/{orgId}/integrations/{integrationId}/authentication/provision-token
func (c *Client) ProvisionBrokerToken(ctx context.Context, orgID, integrationID string) (*BrokerToken, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/authentication/provision-token", orgID, integrationID)
	var result BrokerToken
	if err := c.base.Post(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BrokerToken represents a broker token.
type BrokerToken struct {
	Token string `json:"token"`
}

// SwitchBrokerToken switches to a new broker token.
// POST /org/{orgId}/integrations/{integrationId}/authentication/switch-token
func (c *Client) SwitchBrokerToken(ctx context.Context, orgID, integrationID string) (*BrokerToken, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/authentication/switch-token", orgID, integrationID)
	var result BrokerToken
	if err := c.base.Post(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Clone
// ============================================================================

// CloneRequest represents a clone integration request.
type CloneRequest struct {
	SourceOrgID string `json:"sourceOrgId"`
}

// Clone clones an integration from another organization.
// POST /org/{orgId}/integrations/{integrationId}/clone
func (c *Client) Clone(ctx context.Context, targetOrgID, integrationID string, req *CloneRequest) error {
	path := fmt.Sprintf("/org/%s/integrations/%s/clone", targetOrgID, integrationID)
	return c.base.Post(ctx, path, req, nil)
}

// ============================================================================
// Import
// ============================================================================

// ImportRequest represents an import request.
type ImportRequest struct {
	Target ImportTarget `json:"target"`
	Files  *ImportFiles `json:"files,omitempty"`
}

// ImportTarget represents an import target.
type ImportTarget struct {
	Owner      string `json:"owner,omitempty"`
	Name       string `json:"name,omitempty"`
	Branch     string `json:"branch,omitempty"`
	ID         int    `json:"id,omitempty"`         // For GitLab
	ProjectKey string `json:"projectKey,omitempty"` // For Bitbucket Server
	RepoSlug   string `json:"repoSlug,omitempty"`   // For Bitbucket Server
}

// ImportFiles represents files to import.
type ImportFiles struct {
	Manifest string `json:"manifest,omitempty"`
}

// ImportResponse represents an import response.
type ImportResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Import imports a project from an integration.
// POST /org/{orgId}/integrations/{integrationId}/import
func (c *Client) Import(ctx context.Context, orgID, integrationID string, req *ImportRequest) (*ImportResponse, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/import", orgID, integrationID)
	var result ImportResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ImportJob represents an import job.
type ImportJob struct {
	ID      string      `json:"id"`
	Status  string      `json:"status"`
	Created string      `json:"created"`
	Logs    []ImportLog `json:"logs,omitempty"`
}

// ImportLog represents an import log entry.
type ImportLog struct {
	Name     string            `json:"name"`
	Created  string            `json:"created"`
	Status   string            `json:"status"`
	Projects []ImportedProject `json:"projects,omitempty"`
}

// ImportedProject represents an imported project.
type ImportedProject struct {
	ProjectURL string `json:"projectUrl"`
	Success    bool   `json:"success"`
	TargetFile string `json:"targetFile,omitempty"`
}

// GetImportJob retrieves an import job status.
// GET /org/{orgId}/integrations/{integrationId}/import/{jobId}
func (c *Client) GetImportJob(ctx context.Context, orgID, integrationID, jobID string) (*ImportJob, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/import/%s", orgID, integrationID, jobID)
	var result ImportJob
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
