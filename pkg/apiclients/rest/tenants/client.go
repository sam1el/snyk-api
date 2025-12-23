// Package tenants provides a manual client for Snyk REST Tenants API.
package tenants

import (
	"context"
	"fmt"
	"time"

	"github.com/sam1el/snyk-api/pkg/apiclients/rest"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to REST Tenants API endpoints.
type Client struct {
	base *rest.BaseClient
}

// New creates a new REST Tenants client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: rest.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Tenants
// ============================================================================

// Tenant represents a Snyk tenant.
type Tenant struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes TenantAttrs `json:"attributes"`
}

// TenantAttrs contains tenant attributes.
type TenantAttrs struct {
	Name    string     `json:"name"`
	Slug    string     `json:"slug,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// TenantResponse wraps a single tenant.
type TenantResponse struct {
	Data Tenant `json:"data"`
}

// TenantListResponse wraps a list of tenants.
type TenantListResponse struct {
	Data  []Tenant   `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// List lists all tenants.
// GET /tenants
func (c *Client) List(ctx context.Context, params *rest.PaginationParams) (*TenantListResponse, error) {
	var result TenantListResponse
	if err := c.base.Get(ctx, "/tenants", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a tenant by ID.
// GET /tenants/{tenant_id}
func (c *Client) Get(ctx context.Context, tenantID string) (*TenantResponse, error) {
	path := fmt.Sprintf("/tenants/%s", tenantID)
	var result TenantResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Memberships
// ============================================================================

// Membership represents a tenant membership.
type Membership struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes MembershipAttrs `json:"attributes"`
}

// MembershipAttrs contains membership attributes.
type MembershipAttrs struct {
	Email   string     `json:"email,omitempty"`
	Name    string     `json:"name,omitempty"`
	Role    string     `json:"role,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// MembershipListResponse wraps a list of memberships.
type MembershipListResponse struct {
	Data  []Membership `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// MembershipResponse wraps a single membership.
type MembershipResponse struct {
	Data Membership `json:"data"`
}

// ListMemberships lists all memberships for a tenant.
// GET /tenants/{tenant_id}/memberships
func (c *Client) ListMemberships(ctx context.Context, tenantID string, params *rest.PaginationParams) (*MembershipListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/memberships", tenantID)
	var result MembershipListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMembership retrieves a membership by ID.
// GET /tenants/{tenant_id}/memberships/{membership_id}
func (c *Client) GetMembership(ctx context.Context, tenantID, membershipID string) (*MembershipResponse, error) {
	path := fmt.Sprintf("/tenants/%s/memberships/%s", tenantID, membershipID)
	var result MembershipResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMembership removes a membership.
// DELETE /tenants/{tenant_id}/memberships/{membership_id}
func (c *Client) DeleteMembership(ctx context.Context, tenantID, membershipID string) error {
	path := fmt.Sprintf("/tenants/%s/memberships/%s", tenantID, membershipID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Roles
// ============================================================================

// Role represents a tenant role.
type Role struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Attributes RoleAttrs `json:"attributes"`
}

// RoleAttrs contains role attributes.
type RoleAttrs struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Created     *time.Time `json:"created_at,omitempty"`
}

// RoleListResponse wraps a list of roles.
type RoleListResponse struct {
	Data  []Role     `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// RoleResponse wraps a single role.
type RoleResponse struct {
	Data Role `json:"data"`
}

// ListRoles lists all roles for a tenant.
// GET /tenants/{tenant_id}/roles
func (c *Client) ListRoles(ctx context.Context, tenantID string, params *rest.PaginationParams) (*RoleListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/roles", tenantID)
	var result RoleListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRole retrieves a role by ID.
// GET /tenants/{tenant_id}/roles/{role_id}
func (c *Client) GetRole(ctx context.Context, tenantID, roleID string) (*RoleResponse, error) {
	path := fmt.Sprintf("/tenants/%s/roles/%s", tenantID, roleID)
	var result RoleResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Broker Deployments
// ============================================================================

// BrokerDeployment represents a broker deployment.
type BrokerDeployment struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes BrokerDeploymentAttrs `json:"attributes"`
}

// BrokerDeploymentAttrs contains broker deployment attributes.
type BrokerDeploymentAttrs struct {
	Name    string     `json:"name"`
	Status  string     `json:"status,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// BrokerDeploymentListResponse wraps a list of broker deployments.
type BrokerDeploymentListResponse struct {
	Data  []BrokerDeployment `json:"data"`
	Links rest.Links         `json:"links,omitempty"`
}

// BrokerDeploymentResponse wraps a single broker deployment.
type BrokerDeploymentResponse struct {
	Data BrokerDeployment `json:"data"`
}

// ListBrokerDeployments lists all broker deployments for a tenant.
// GET /tenants/{tenant_id}/brokers/deployments
func (c *Client) ListBrokerDeployments(ctx context.Context, tenantID string, params *rest.PaginationParams) (*BrokerDeploymentListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/deployments", tenantID)
	var result BrokerDeploymentListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBrokerDeployment retrieves a broker deployment by ID.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}
func (c *Client) GetBrokerDeployment(ctx context.Context, tenantID, installID, deploymentID string) (*BrokerDeploymentResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s", tenantID, installID, deploymentID)
	var result BrokerDeploymentResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInstallDeployments lists broker deployments for an install.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments
func (c *Client) ListInstallDeployments(ctx context.Context, tenantID, installID string, params *rest.PaginationParams) (*BrokerDeploymentListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments", tenantID, installID)
	var result BrokerDeploymentListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Broker Connections
// ============================================================================

// BrokerConnection represents a broker connection.
type BrokerConnection struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes BrokerConnectionAttrs `json:"attributes"`
}

// BrokerConnectionAttrs contains broker connection attributes.
type BrokerConnectionAttrs struct {
	Name    string     `json:"name"`
	Status  string     `json:"status,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// BrokerConnectionListResponse wraps a list of broker connections.
type BrokerConnectionListResponse struct {
	Data  []BrokerConnection `json:"data"`
	Links rest.Links         `json:"links,omitempty"`
}

// BrokerConnectionResponse wraps a single broker connection.
type BrokerConnectionResponse struct {
	Data BrokerConnection `json:"data"`
}

// ListBrokerConnections lists broker connections for a deployment.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/connections
func (c *Client) ListBrokerConnections(ctx context.Context, tenantID, installID, deploymentID string, params *rest.PaginationParams) (*BrokerConnectionListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/connections", tenantID, installID, deploymentID)
	var result BrokerConnectionListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBrokerConnection retrieves a broker connection by ID.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/connections/{connection_id}
func (c *Client) GetBrokerConnection(ctx context.Context, tenantID, installID, deploymentID, connectionID string) (*BrokerConnectionResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/connections/%s", tenantID, installID, deploymentID, connectionID)
	var result BrokerConnectionResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteBrokerConnection removes a broker connection.
// DELETE /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/connections/{connection_id}
func (c *Client) DeleteBrokerConnection(ctx context.Context, tenantID, installID, deploymentID, connectionID string) error {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/connections/%s", tenantID, installID, deploymentID, connectionID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Broker Credentials
// ============================================================================

// BrokerCredential represents a broker credential.
type BrokerCredential struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes BrokerCredentialAttrs `json:"attributes"`
}

// BrokerCredentialAttrs contains broker credential attributes.
type BrokerCredentialAttrs struct {
	Name    string     `json:"name"`
	Type    string     `json:"type,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// BrokerCredentialListResponse wraps a list of broker credentials.
type BrokerCredentialListResponse struct {
	Data  []BrokerCredential `json:"data"`
	Links rest.Links         `json:"links,omitempty"`
}

// BrokerCredentialResponse wraps a single broker credential.
type BrokerCredentialResponse struct {
	Data BrokerCredential `json:"data"`
}

// ListBrokerCredentials lists broker credentials for a deployment.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/credentials
func (c *Client) ListBrokerCredentials(ctx context.Context, tenantID, installID, deploymentID string, params *rest.PaginationParams) (*BrokerCredentialListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/credentials", tenantID, installID, deploymentID)
	var result BrokerCredentialListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBrokerCredential retrieves a broker credential by ID.
// GET /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/credentials/{credential_id}
func (c *Client) GetBrokerCredential(ctx context.Context, tenantID, installID, deploymentID, credentialID string) (*BrokerCredentialResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/credentials/%s", tenantID, installID, deploymentID, credentialID)
	var result BrokerCredentialResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteBrokerCredential removes a broker credential.
// DELETE /tenants/{tenant_id}/brokers/installs/{install_id}/deployments/{deployment_id}/credentials/{credential_id}
func (c *Client) DeleteBrokerCredential(ctx context.Context, tenantID, installID, deploymentID, credentialID string) error {
	path := fmt.Sprintf("/tenants/%s/brokers/installs/%s/deployments/%s/credentials/%s", tenantID, installID, deploymentID, credentialID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Broker Integrations
// ============================================================================

// BrokerIntegration represents a broker integration.
type BrokerIntegration struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes BrokerIntegrationAttrs `json:"attributes"`
}

// BrokerIntegrationAttrs contains broker integration attributes.
type BrokerIntegrationAttrs struct {
	Type    string     `json:"type"`
	Created *time.Time `json:"created_at,omitempty"`
}

// BrokerIntegrationListResponse wraps a list of broker integrations.
type BrokerIntegrationListResponse struct {
	Data  []BrokerIntegration `json:"data"`
	Links rest.Links          `json:"links,omitempty"`
}

// ListBrokerIntegrations lists broker integrations for a connection.
// GET /tenants/{tenant_id}/brokers/connections/{connection_id}/integrations
func (c *Client) ListBrokerIntegrations(ctx context.Context, tenantID, connectionID string, params *rest.PaginationParams) (*BrokerIntegrationListResponse, error) {
	path := fmt.Sprintf("/tenants/%s/brokers/connections/%s/integrations", tenantID, connectionID)
	var result BrokerIntegrationListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
