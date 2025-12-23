// Package orgs provides a manual client for Snyk REST Organizations API.
package orgs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/sam1el/snyk-api/pkg/apiclients/rest"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to REST Organizations API endpoints.
type Client struct {
	base *rest.BaseClient
}

// New creates a new REST Organizations client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: rest.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Common Types
// ============================================================================

// Organization represents a Snyk organization.
type Organization struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes OrganizationAttrs `json:"attributes"`
}

// OrganizationAttrs contains organization attributes.
type OrganizationAttrs struct {
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	GroupID    string     `json:"group_id,omitempty"`
	Created    *time.Time `json:"created_at,omitempty"`
	IsPersonal bool       `json:"is_personal,omitempty"`
}

// OrganizationResponse wraps a single organization.
type OrganizationResponse struct {
	Data  Organization `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// OrganizationListResponse wraps a list of organizations.
type OrganizationListResponse struct {
	Data  []Organization `json:"data"`
	Links rest.Links     `json:"links,omitempty"`
	Meta  rest.Meta      `json:"meta,omitempty"`
}

// ============================================================================
// Organizations CRUD
// ============================================================================

// List lists all organizations.
// GET /orgs
func (c *Client) List(ctx context.Context, params *rest.PaginationParams) (*OrganizationListResponse, error) {
	var result OrganizationListResponse
	if err := c.base.Get(ctx, "/orgs", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves an organization by ID.
// GET /orgs/{org_id}
func (c *Client) Get(ctx context.Context, orgID string) (*OrganizationResponse, error) {
	path := fmt.Sprintf("/orgs/%s", orgID)
	var result OrganizationResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRequest represents a request to update an organization.
type UpdateRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name string `json:"name,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// Update updates an organization.
// PATCH /orgs/{org_id}
func (c *Client) Update(ctx context.Context, orgID string, name string) (*OrganizationResponse, error) {
	path := fmt.Sprintf("/orgs/%s", orgID)
	req := UpdateRequest{}
	req.Data.Type = "org"
	req.Data.Attributes.Name = name
	var result OrganizationResponse
	if err := c.base.Patch(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Memberships
// ============================================================================

// Membership represents an organization membership.
type Membership struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes MembershipAttrs `json:"attributes"`
}

// MembershipAttrs contains membership attributes.
type MembershipAttrs struct {
	Email    string     `json:"email,omitempty"`
	Name     string     `json:"name,omitempty"`
	Username string     `json:"username,omitempty"`
	Role     string     `json:"role,omitempty"`
	Created  *time.Time `json:"created_at,omitempty"`
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

// ListMemberships lists all memberships for an organization.
// GET /orgs/{org_id}/memberships
func (c *Client) ListMemberships(ctx context.Context, orgID string, params *rest.PaginationParams) (*MembershipListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/memberships", orgID)
	var result MembershipListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMembership retrieves a membership by ID.
// GET /orgs/{org_id}/memberships/{membership_id}
func (c *Client) GetMembership(ctx context.Context, orgID, membershipID string) (*MembershipResponse, error) {
	path := fmt.Sprintf("/orgs/%s/memberships/%s", orgID, membershipID)
	var result MembershipResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateMembershipRequest represents a request to update a membership.
type UpdateMembershipRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Role string `json:"role"`
		} `json:"attributes"`
	} `json:"data"`
}

// UpdateMembership updates a membership's role.
// PATCH /orgs/{org_id}/memberships/{membership_id}
func (c *Client) UpdateMembership(ctx context.Context, orgID, membershipID, role string) (*MembershipResponse, error) {
	path := fmt.Sprintf("/orgs/%s/memberships/%s", orgID, membershipID)
	req := UpdateMembershipRequest{}
	req.Data.Type = "org_membership"
	req.Data.Attributes.Role = role
	var result MembershipResponse
	if err := c.base.Patch(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMembership removes a membership.
// DELETE /orgs/{org_id}/memberships/{membership_id}
func (c *Client) DeleteMembership(ctx context.Context, orgID, membershipID string) error {
	path := fmt.Sprintf("/orgs/%s/memberships/%s", orgID, membershipID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Invites
// ============================================================================

// Invite represents an organization invite.
type Invite struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes InviteAttrs `json:"attributes"`
}

// InviteAttrs contains invite attributes.
type InviteAttrs struct {
	Email     string     `json:"email"`
	Role      string     `json:"role,omitempty"`
	Created   *time.Time `json:"created_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// InviteListResponse wraps a list of invites.
type InviteListResponse struct {
	Data  []Invite   `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// InviteResponse wraps a single invite.
type InviteResponse struct {
	Data Invite `json:"data"`
}

// ListInvites lists all pending invites for an organization.
// GET /orgs/{org_id}/invites
func (c *Client) ListInvites(ctx context.Context, orgID string, params *rest.PaginationParams) (*InviteListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/invites", orgID)
	var result InviteListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateInviteRequest represents a request to create an invite.
type CreateInviteRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Email string `json:"email"`
			Role  string `json:"role,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateInvite creates a new invite.
// POST /orgs/{org_id}/invites
func (c *Client) CreateInvite(ctx context.Context, orgID, email, role string) (*InviteResponse, error) {
	path := fmt.Sprintf("/orgs/%s/invites", orgID)
	req := CreateInviteRequest{}
	req.Data.Type = "org_invitation"
	req.Data.Attributes.Email = email
	req.Data.Attributes.Role = role
	var result InviteResponse
	if err := c.base.Post(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteInvite cancels an invite.
// DELETE /orgs/{org_id}/invites/{invite_id}
func (c *Client) DeleteInvite(ctx context.Context, orgID, inviteID string) error {
	path := fmt.Sprintf("/orgs/%s/invites/%s", orgID, inviteID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Service Accounts
// ============================================================================

// ServiceAccount represents a service account.
type ServiceAccount struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes ServiceAccountAttrs `json:"attributes"`
}

// ServiceAccountAttrs contains service account attributes.
type ServiceAccountAttrs struct {
	Name           string     `json:"name"`
	AuthType       string     `json:"auth_type,omitempty"`
	RolePublicID   string     `json:"role_public_id,omitempty"`
	JwksURL        string     `json:"jwks_url,omitempty"`
	AccessTokenTTL int        `json:"access_token_ttl_seconds,omitempty"`
	Created        *time.Time `json:"created_at,omitempty"`
}

// ServiceAccountListResponse wraps a list of service accounts.
type ServiceAccountListResponse struct {
	Data  []ServiceAccount `json:"data"`
	Links rest.Links       `json:"links,omitempty"`
}

// ServiceAccountResponse wraps a single service account.
type ServiceAccountResponse struct {
	Data ServiceAccount `json:"data"`
}

// ListServiceAccounts lists all service accounts for an organization.
// GET /orgs/{org_id}/service_accounts
func (c *Client) ListServiceAccounts(ctx context.Context, orgID string, params *rest.PaginationParams) (*ServiceAccountListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/service_accounts", orgID)
	var result ServiceAccountListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServiceAccount retrieves a service account by ID.
// GET /orgs/{org_id}/service_accounts/{serviceaccount_id}
func (c *Client) GetServiceAccount(ctx context.Context, orgID, saID string) (*ServiceAccountResponse, error) {
	path := fmt.Sprintf("/orgs/%s/service_accounts/%s", orgID, saID)
	var result ServiceAccountResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateServiceAccountRequest represents a request to create a service account.
type CreateServiceAccountRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name           string `json:"name"`
			AuthType       string `json:"auth_type"`
			RolePublicID   string `json:"role_public_id,omitempty"`
			JwksURL        string `json:"jwks_url,omitempty"`
			AccessTokenTTL int    `json:"access_token_ttl_seconds,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// ServiceAccountWithSecret includes the secret (only returned on creation).
type ServiceAccountWithSecret struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name         string     `json:"name"`
		AuthType     string     `json:"auth_type"`
		APIKey       string     `json:"api_key,omitempty"`
		ClientID     string     `json:"client_id,omitempty"`
		ClientSecret string     `json:"client_secret,omitempty"`
		Created      *time.Time `json:"created_at,omitempty"`
	} `json:"attributes"`
}

// ServiceAccountCreateResponse wraps a created service account with secret.
type ServiceAccountCreateResponse struct {
	Data ServiceAccountWithSecret `json:"data"`
}

// CreateServiceAccount creates a new service account.
// POST /orgs/{org_id}/service_accounts
func (c *Client) CreateServiceAccount(ctx context.Context, orgID string, req *CreateServiceAccountRequest) (*ServiceAccountCreateResponse, error) {
	path := fmt.Sprintf("/orgs/%s/service_accounts", orgID)
	req.Data.Type = "service_account"
	var result ServiceAccountCreateResponse
	if err := c.base.Post(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateServiceAccountRequest represents a request to update a service account.
type UpdateServiceAccountRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name string `json:"name,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// UpdateServiceAccount updates a service account.
// PATCH /orgs/{org_id}/service_accounts/{serviceaccount_id}
func (c *Client) UpdateServiceAccount(ctx context.Context, orgID, saID, name string) (*ServiceAccountResponse, error) {
	path := fmt.Sprintf("/orgs/%s/service_accounts/%s", orgID, saID)
	req := UpdateServiceAccountRequest{}
	req.Data.Type = "service_account"
	req.Data.Attributes.Name = name
	var result ServiceAccountResponse
	if err := c.base.Patch(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteServiceAccount deletes a service account.
// DELETE /orgs/{org_id}/service_accounts/{serviceaccount_id}
func (c *Client) DeleteServiceAccount(ctx context.Context, orgID, saID string) error {
	path := fmt.Sprintf("/orgs/%s/service_accounts/%s", orgID, saID)
	return c.base.Delete(ctx, path, nil)
}

// ServiceAccountSecret represents a service account secret.
type ServiceAccountSecret struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		ClientSecret string     `json:"client_secret,omitempty"`
		APIKey       string     `json:"api_key,omitempty"`
		Created      *time.Time `json:"created_at,omitempty"`
	} `json:"attributes"`
}

// ServiceAccountSecretResponse wraps a service account secret.
type ServiceAccountSecretResponse struct {
	Data ServiceAccountSecret `json:"data"`
}

// RotateServiceAccountSecret rotates the secret for a service account.
// POST /orgs/{org_id}/service_accounts/{serviceaccount_id}/secrets
func (c *Client) RotateServiceAccountSecret(ctx context.Context, orgID, saID string) (*ServiceAccountSecretResponse, error) {
	path := fmt.Sprintf("/orgs/%s/service_accounts/%s/secrets", orgID, saID)
	var result ServiceAccountSecretResponse
	if err := c.base.Post(ctx, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Policies
// ============================================================================

// Policy represents an organization policy.
type Policy struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes PolicyAttrs `json:"attributes"`
}

// PolicyAttrs contains policy attributes.
type PolicyAttrs struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Severity    string     `json:"severity,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
}

// PolicyListResponse wraps a list of policies.
type PolicyListResponse struct {
	Data  []Policy   `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// PolicyResponse wraps a single policy.
type PolicyResponse struct {
	Data Policy `json:"data"`
}

// ListPolicies lists all policies for an organization.
// GET /orgs/{org_id}/policies
func (c *Client) ListPolicies(ctx context.Context, orgID string, params *rest.PaginationParams) (*PolicyListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/policies", orgID)
	var result PolicyListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPolicy retrieves a policy by ID.
// GET /orgs/{org_id}/policies/{policy_id}
func (c *Client) GetPolicy(ctx context.Context, orgID, policyID string) (*PolicyResponse, error) {
	path := fmt.Sprintf("/orgs/%s/policies/%s", orgID, policyID)
	var result PolicyResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePolicy deletes a policy.
// DELETE /orgs/{org_id}/policies/{policy_id}
func (c *Client) DeletePolicy(ctx context.Context, orgID, policyID string) error {
	path := fmt.Sprintf("/orgs/%s/policies/%s", orgID, policyID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Collections
// ============================================================================

// Collection represents a collection.
type Collection struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes CollectionAttrs `json:"attributes"`
}

// CollectionAttrs contains collection attributes.
type CollectionAttrs struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Created     *time.Time `json:"created_at,omitempty"`
	Updated     *time.Time `json:"updated_at,omitempty"`
}

// CollectionListResponse wraps a list of collections.
type CollectionListResponse struct {
	Data  []Collection `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// CollectionResponse wraps a single collection.
type CollectionResponse struct {
	Data Collection `json:"data"`
}

// ListCollections lists all collections for an organization.
// GET /orgs/{org_id}/collections
func (c *Client) ListCollections(ctx context.Context, orgID string, params *rest.PaginationParams) (*CollectionListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/collections", orgID)
	var result CollectionListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCollection retrieves a collection by ID.
// GET /orgs/{org_id}/collections/{collection_id}
func (c *Client) GetCollection(ctx context.Context, orgID, collectionID string) (*CollectionResponse, error) {
	path := fmt.Sprintf("/orgs/%s/collections/%s", orgID, collectionID)
	var result CollectionResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCollectionRequest represents a request to create a collection.
type CreateCollectionRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateCollection creates a new collection.
// POST /orgs/{org_id}/collections
func (c *Client) CreateCollection(ctx context.Context, orgID, name, description string) (*CollectionResponse, error) {
	path := fmt.Sprintf("/orgs/%s/collections", orgID)
	req := CreateCollectionRequest{}
	req.Data.Type = "resource"
	req.Data.Attributes.Name = name
	req.Data.Attributes.Description = description
	var result CollectionResponse
	if err := c.base.Post(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCollection deletes a collection.
// DELETE /orgs/{org_id}/collections/{collection_id}
func (c *Client) DeleteCollection(ctx context.Context, orgID, collectionID string) error {
	path := fmt.Sprintf("/orgs/%s/collections/%s", orgID, collectionID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Settings
// ============================================================================

// IaCSettings represents IaC settings.
type IaCSettings struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes IaCSettingsAttrs `json:"attributes"`
}

// IaCSettingsAttrs contains IaC settings attributes.
type IaCSettingsAttrs struct {
	CustomRulesEnabled bool `json:"custom_rules_enabled,omitempty"`
}

// IaCSettingsResponse wraps IaC settings.
type IaCSettingsResponse struct {
	Data IaCSettings `json:"data"`
}

// GetIaCSettings retrieves IaC settings for an organization.
// GET /orgs/{org_id}/settings/iac
func (c *Client) GetIaCSettings(ctx context.Context, orgID string) (*IaCSettingsResponse, error) {
	path := fmt.Sprintf("/orgs/%s/settings/iac", orgID)
	var result IaCSettingsResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SASTSettings represents SAST settings.
type SASTSettings struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes SASTSettingsAttrs `json:"attributes"`
}

// SASTSettingsAttrs contains SAST settings attributes.
type SASTSettingsAttrs struct {
	SASTEnabled bool `json:"sast_enabled,omitempty"`
}

// SASTSettingsResponse wraps SAST settings.
type SASTSettingsResponse struct {
	Data SASTSettings `json:"data"`
}

// GetSASTSettings retrieves SAST settings for an organization.
// GET /orgs/{org_id}/settings/sast
func (c *Client) GetSASTSettings(ctx context.Context, orgID string) (*SASTSettingsResponse, error) {
	path := fmt.Sprintf("/orgs/%s/settings/sast", orgID)
	var result SASTSettingsResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenSourceSettings represents open source settings.
type OpenSourceSettings struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes OpenSourceSettingsAttrs `json:"attributes"`
}

// OpenSourceSettingsAttrs contains open source settings attributes.
type OpenSourceSettingsAttrs struct {
	PackageRepositoryIntegrationEnabled bool `json:"package_repository_integration_enabled,omitempty"`
}

// OpenSourceSettingsResponse wraps open source settings.
type OpenSourceSettingsResponse struct {
	Data OpenSourceSettings `json:"data"`
}

// GetOpenSourceSettings retrieves open source settings for an organization.
// GET /orgs/{org_id}/settings/opensource
func (c *Client) GetOpenSourceSettings(ctx context.Context, orgID string) (*OpenSourceSettingsResponse, error) {
	path := fmt.Sprintf("/orgs/%s/settings/opensource", orgID)
	var result OpenSourceSettingsResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Audit Logs
// ============================================================================

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes AuditLogAttrs `json:"attributes"`
}

// AuditLogAttrs contains audit log attributes.
type AuditLogAttrs struct {
	Created *time.Time             `json:"created,omitempty"`
	UserID  string                 `json:"user_id,omitempty"`
	Event   string                 `json:"event,omitempty"`
	Content map[string]interface{} `json:"content,omitempty"`
}

// AuditLogListResponse wraps a list of audit logs.
type AuditLogListResponse struct {
	Data  []AuditLog `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// SearchAuditLogs searches audit logs for an organization.
// GET /orgs/{org_id}/audit_logs/search
func (c *Client) SearchAuditLogs(ctx context.Context, orgID string, params *rest.PaginationParams) (*AuditLogListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/audit_logs/search", orgID)
	var result AuditLogListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Projects (additional endpoints)
// ============================================================================

// Project represents a project.
type Project struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Attributes ProjectAttrs `json:"attributes"`
}

// ProjectAttrs contains project attributes.
type ProjectAttrs struct {
	Name                string     `json:"name"`
	Created             *time.Time `json:"created,omitempty"`
	Origin              string     `json:"origin,omitempty"`
	Type                string     `json:"type,omitempty"`
	TargetReference     string     `json:"target_reference,omitempty"`
	Status              string     `json:"status,omitempty"`
	BusinessCriticality []string   `json:"business_criticality,omitempty"`
	Environment         []string   `json:"environment,omitempty"`
	Lifecycle           []string   `json:"lifecycle,omitempty"`
}

// ProjectListResponse wraps a list of projects.
type ProjectListResponse struct {
	Data  []Project  `json:"data"`
	Links rest.Links `json:"links,omitempty"`
	Meta  rest.Meta  `json:"meta,omitempty"`
}

// ProjectResponse wraps a single project.
type ProjectResponse struct {
	Data Project `json:"data"`
}

// ListProjects lists all projects for an organization.
// GET /orgs/{org_id}/projects
func (c *Client) ListProjects(ctx context.Context, orgID string, params *rest.PaginationParams) (*ProjectListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/projects", orgID)
	var result ProjectListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProject retrieves a project by ID.
// GET /orgs/{org_id}/projects/{project_id}
func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*ProjectResponse, error) {
	path := fmt.Sprintf("/orgs/%s/projects/%s", orgID, projectID)
	var result ProjectResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProject deletes a project.
// DELETE /orgs/{org_id}/projects/{project_id}
func (c *Client) DeleteProject(ctx context.Context, orgID, projectID string) error {
	path := fmt.Sprintf("/orgs/%s/projects/%s", orgID, projectID)
	return c.base.Delete(ctx, path, nil)
}

// GetProjectSBOM retrieves the SBOM for a project.
// GET /orgs/{org_id}/projects/{project_id}/sbom
func (c *Client) GetProjectSBOM(ctx context.Context, orgID, projectID, format string) ([]byte, error) {
	path := fmt.Sprintf("/orgs/%s/projects/%s/sbom", orgID, projectID)
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	var result json.RawMessage
	if err := c.base.Get(ctx, path, q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ============================================================================
// Targets
// ============================================================================

// Target represents a target.
type Target struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes TargetAttrs `json:"attributes"`
}

// TargetAttrs contains target attributes.
type TargetAttrs struct {
	DisplayName string     `json:"display_name"`
	URL         string     `json:"url,omitempty"`
	Origin      string     `json:"origin,omitempty"`
	Created     *time.Time `json:"created_at,omitempty"`
}

// TargetListResponse wraps a list of targets.
type TargetListResponse struct {
	Data  []Target   `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// TargetResponse wraps a single target.
type TargetResponse struct {
	Data Target `json:"data"`
}

// ListTargets lists all targets for an organization.
// GET /orgs/{org_id}/targets
func (c *Client) ListTargets(ctx context.Context, orgID string, params *rest.PaginationParams) (*TargetListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/targets", orgID)
	var result TargetListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTarget retrieves a target by ID.
// GET /orgs/{org_id}/targets/{target_id}
func (c *Client) GetTarget(ctx context.Context, orgID, targetID string) (*TargetResponse, error) {
	path := fmt.Sprintf("/orgs/%s/targets/%s", orgID, targetID)
	var result TargetResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTarget deletes a target.
// DELETE /orgs/{org_id}/targets/{target_id}
func (c *Client) DeleteTarget(ctx context.Context, orgID, targetID string) error {
	path := fmt.Sprintf("/orgs/%s/targets/%s", orgID, targetID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Issues
// ============================================================================

// Issue represents an issue.
type Issue struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes IssueAttrs `json:"attributes"`
}

// IssueAttrs contains issue attributes.
type IssueAttrs struct {
	Title       string     `json:"title"`
	Severity    string     `json:"severity,omitempty"`
	Type        string     `json:"type,omitempty"`
	Status      string     `json:"status,omitempty"`
	Description string     `json:"description,omitempty"`
	Created     *time.Time `json:"created_at,omitempty"`
	Updated     *time.Time `json:"updated_at,omitempty"`
}

// IssueListResponse wraps a list of issues.
type IssueListResponse struct {
	Data  []Issue    `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// IssueResponse wraps a single issue.
type IssueResponse struct {
	Data Issue `json:"data"`
}

// ListIssuesParams represents parameters for listing issues.
type ListIssuesParams struct {
	rest.PaginationParams
	Severity string
	Type     string
	Status   string
}

// ToQuery converts issue params to URL query values.
func (p *ListIssuesParams) ToQuery() url.Values {
	q := p.PaginationParams.ToQuery()
	if p.Severity != "" {
		q.Set("severity", p.Severity)
	}
	if p.Type != "" {
		q.Set("type", p.Type)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	return q
}

// ListIssues lists all issues for an organization.
// GET /orgs/{org_id}/issues
func (c *Client) ListIssues(ctx context.Context, orgID string, params *ListIssuesParams) (*IssueListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/issues", orgID)
	var q url.Values
	if params != nil {
		q = params.ToQuery()
	}
	var result IssueListResponse
	if err := c.base.Get(ctx, path, q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssue retrieves an issue by ID.
// GET /orgs/{org_id}/issues/{issue_id}
func (c *Client) GetIssue(ctx context.Context, orgID, issueID string) (*IssueResponse, error) {
	path := fmt.Sprintf("/orgs/%s/issues/%s", orgID, issueID)
	var result IssueResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// SBOM Testing
// ============================================================================

// SBOMTestJob represents an SBOM test job.
type SBOMTestJob struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes SBOMTestJobAttrs `json:"attributes"`
}

// SBOMTestJobAttrs contains SBOM test job attributes.
type SBOMTestJobAttrs struct {
	Status    string     `json:"status"`
	Created   *time.Time `json:"created_at,omitempty"`
	Completed *time.Time `json:"completed_at,omitempty"`
}

// SBOMTestJobResponse wraps an SBOM test job.
type SBOMTestJobResponse struct {
	Data SBOMTestJob `json:"data"`
}

// CreateSBOMTestRequest represents a request to create an SBOM test.
type CreateSBOMTestRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			SBOM interface{} `json:"sbom"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateSBOMTest creates a new SBOM test.
// POST /orgs/{org_id}/sbom_tests
func (c *Client) CreateSBOMTest(ctx context.Context, orgID string, sbom interface{}) (*SBOMTestJobResponse, error) {
	path := fmt.Sprintf("/orgs/%s/sbom_tests", orgID)
	req := CreateSBOMTestRequest{}
	req.Data.Type = "sbom_test"
	req.Data.Attributes.SBOM = sbom
	var result SBOMTestJobResponse
	if err := c.base.Post(ctx, path, nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSBOMTestJob retrieves an SBOM test job status.
// GET /orgs/{org_id}/sbom_tests/{job_id}
func (c *Client) GetSBOMTestJob(ctx context.Context, orgID, jobID string) (*SBOMTestJobResponse, error) {
	path := fmt.Sprintf("/orgs/%s/sbom_tests/%s", orgID, jobID)
	var result SBOMTestJobResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SBOMTestResult represents SBOM test results.
type SBOMTestResult struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Summary SBOMTestSummary `json:"summary"`
			Issues  []Issue         `json:"issues,omitempty"`
		} `json:"attributes"`
	} `json:"data"`
}

// SBOMTestSummary contains SBOM test summary.
type SBOMTestSummary struct {
	TotalIssues int `json:"total_issues"`
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
}

// GetSBOMTestResults retrieves SBOM test results.
// GET /orgs/{org_id}/sbom_tests/{job_id}/results
func (c *Client) GetSBOMTestResults(ctx context.Context, orgID, jobID string) (*SBOMTestResult, error) {
	path := fmt.Sprintf("/orgs/%s/sbom_tests/%s/results", orgID, jobID)
	var result SBOMTestResult
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Container Images
// ============================================================================

// ContainerImage represents a container image.
type ContainerImage struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes ContainerImageAttrs `json:"attributes"`
}

// ContainerImageAttrs contains container image attributes.
type ContainerImageAttrs struct {
	Names    []string   `json:"names,omitempty"`
	Platform string     `json:"platform,omitempty"`
	Created  *time.Time `json:"created_at,omitempty"`
}

// ContainerImageListResponse wraps a list of container images.
type ContainerImageListResponse struct {
	Data  []ContainerImage `json:"data"`
	Links rest.Links       `json:"links,omitempty"`
}

// ContainerImageResponse wraps a single container image.
type ContainerImageResponse struct {
	Data ContainerImage `json:"data"`
}

// ListContainerImages lists all container images for an organization.
// GET /orgs/{org_id}/container_images
func (c *Client) ListContainerImages(ctx context.Context, orgID string, params *rest.PaginationParams) (*ContainerImageListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/container_images", orgID)
	var result ContainerImageListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetContainerImage retrieves a container image by ID.
// GET /orgs/{org_id}/container_images/{image_id}
func (c *Client) GetContainerImage(ctx context.Context, orgID, imageID string) (*ContainerImageResponse, error) {
	path := fmt.Sprintf("/orgs/%s/container_images/%s", orgID, imageID)
	var result ContainerImageResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Apps
// ============================================================================

// App represents an app.
type App struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Attributes AppAttrs `json:"attributes"`
}

// AppAttrs contains app attributes.
type AppAttrs struct {
	Name     string     `json:"name"`
	ClientID string     `json:"client_id,omitempty"`
	Created  *time.Time `json:"created_at,omitempty"`
}

// AppListResponse wraps a list of apps.
type AppListResponse struct {
	Data  []App      `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// AppResponse wraps a single app.
type AppResponse struct {
	Data App `json:"data"`
}

// ListApps lists all apps for an organization.
// GET /orgs/{org_id}/apps
func (c *Client) ListApps(ctx context.Context, orgID string, params *rest.PaginationParams) (*AppListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/apps", orgID)
	var result AppListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAppCreations lists all app creations for an organization.
// GET /orgs/{org_id}/apps/creations
func (c *Client) ListAppCreations(ctx context.Context, orgID string, params *rest.PaginationParams) (*AppListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/apps/creations", orgID)
	var result AppListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AppInstall represents an app installation.
type AppInstall struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes AppInstallAttrs `json:"attributes"`
}

// AppInstallAttrs contains app installation attributes.
type AppInstallAttrs struct {
	AppID   string     `json:"app_id"`
	Created *time.Time `json:"created_at,omitempty"`
}

// AppInstallListResponse wraps a list of app installations.
type AppInstallListResponse struct {
	Data  []AppInstall `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// ListAppInstalls lists all app installations for an organization.
// GET /orgs/{org_id}/apps/installs
func (c *Client) ListAppInstalls(ctx context.Context, orgID string, params *rest.PaginationParams) (*AppInstallListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/apps/installs", orgID)
	var result AppInstallListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Cloud
// ============================================================================

// CloudEnvironment represents a cloud environment.
type CloudEnvironment struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes CloudEnvironmentAttrs `json:"attributes"`
}

// CloudEnvironmentAttrs contains cloud environment attributes.
type CloudEnvironmentAttrs struct {
	Name     string     `json:"name"`
	Provider string     `json:"provider,omitempty"`
	NativeID string     `json:"native_id,omitempty"`
	Status   string     `json:"status,omitempty"`
	Created  *time.Time `json:"created_at,omitempty"`
}

// CloudEnvironmentListResponse wraps a list of cloud environments.
type CloudEnvironmentListResponse struct {
	Data  []CloudEnvironment `json:"data"`
	Links rest.Links         `json:"links,omitempty"`
}

// ListCloudEnvironments lists all cloud environments for an organization.
// GET /orgs/{org_id}/cloud/environments
func (c *Client) ListCloudEnvironments(ctx context.Context, orgID string, params *rest.PaginationParams) (*CloudEnvironmentListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/cloud/environments", orgID)
	var result CloudEnvironmentListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CloudScan represents a cloud scan.
type CloudScan struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes CloudScanAttrs `json:"attributes"`
}

// CloudScanAttrs contains cloud scan attributes.
type CloudScanAttrs struct {
	Status    string     `json:"status"`
	Created   *time.Time `json:"created_at,omitempty"`
	Completed *time.Time `json:"completed_at,omitempty"`
}

// CloudScanListResponse wraps a list of cloud scans.
type CloudScanListResponse struct {
	Data  []CloudScan `json:"data"`
	Links rest.Links  `json:"links,omitempty"`
}

// ListCloudScans lists all cloud scans for an organization.
// GET /orgs/{org_id}/cloud/scans
func (c *Client) ListCloudScans(ctx context.Context, orgID string, params *rest.PaginationParams) (*CloudScanListResponse, error) {
	path := fmt.Sprintf("/orgs/%s/cloud/scans", orgID)
	var result CloudScanListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Export
// ============================================================================

// ExportJob represents an export job.
type ExportJob struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes ExportJobAttrs `json:"attributes"`
}

// ExportJobAttrs contains export job attributes.
type ExportJobAttrs struct {
	Status      string     `json:"status"`
	Format      string     `json:"format,omitempty"`
	Created     *time.Time `json:"created_at,omitempty"`
	Completed   *time.Time `json:"completed_at,omitempty"`
	DownloadURL string     `json:"download_url,omitempty"`
}

// ExportJobResponse wraps an export job.
type ExportJobResponse struct {
	Data ExportJob `json:"data"`
}

// CreateExport creates a new export job.
// POST /orgs/{org_id}/export
func (c *Client) CreateExport(ctx context.Context, orgID string) (*ExportJobResponse, error) {
	path := fmt.Sprintf("/orgs/%s/export", orgID)
	var result ExportJobResponse
	if err := c.base.Post(ctx, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExport retrieves an export job status.
// GET /orgs/{org_id}/export/{export_id}
func (c *Client) GetExport(ctx context.Context, orgID, exportID string) (*ExportJobResponse, error) {
	path := fmt.Sprintf("/orgs/%s/export/%s", orgID, exportID)
	var result ExportJobResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
