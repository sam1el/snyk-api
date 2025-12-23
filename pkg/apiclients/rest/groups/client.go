// Package groups provides a manual client for Snyk REST Groups API.
package groups

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/sam1el/snyk-api/pkg/apiclients/rest"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to REST Groups API endpoints.
type Client struct {
	base *rest.BaseClient
}

// New creates a new REST Groups client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: rest.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Groups
// ============================================================================

// Group represents a Snyk group.
type Group struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes GroupAttrs `json:"attributes"`
}

// GroupAttrs contains group attributes.
type GroupAttrs struct {
	Name    string     `json:"name"`
	Created *time.Time `json:"created_at,omitempty"`
}

// GroupResponse wraps a single group.
type GroupResponse struct {
	Data Group `json:"data"`
}

// GroupListResponse wraps a list of groups.
type GroupListResponse struct {
	Data  []Group    `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// List lists all groups.
// GET /groups
func (c *Client) List(ctx context.Context, params *rest.PaginationParams) (*GroupListResponse, error) {
	var result GroupListResponse
	if err := c.base.Get(ctx, "/groups", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a group by ID.
// GET /groups/{group_id}
func (c *Client) Get(ctx context.Context, groupID string) (*GroupResponse, error) {
	path := fmt.Sprintf("/groups/%s", groupID)
	var result GroupResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Memberships
// ============================================================================

// Membership represents a group membership.
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

// ListMemberships lists all memberships for a group.
// GET /groups/{group_id}/memberships
func (c *Client) ListMemberships(ctx context.Context, groupID string, params *rest.PaginationParams) (*MembershipListResponse, error) {
	path := fmt.Sprintf("/groups/%s/memberships", groupID)
	var result MembershipListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMembership retrieves a membership by ID.
// GET /groups/{group_id}/memberships/{membership_id}
func (c *Client) GetMembership(ctx context.Context, groupID, membershipID string) (*MembershipResponse, error) {
	path := fmt.Sprintf("/groups/%s/memberships/%s", groupID, membershipID)
	var result MembershipResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteMembership removes a membership.
// DELETE /groups/{group_id}/memberships/{membership_id}
func (c *Client) DeleteMembership(ctx context.Context, groupID, membershipID string) error {
	path := fmt.Sprintf("/groups/%s/memberships/%s", groupID, membershipID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Org Memberships
// ============================================================================

// OrgMembershipListResponse wraps a list of org memberships in a group.
type OrgMembershipListResponse struct {
	Data  []Membership `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// ListOrgMemberships lists all org memberships for a group.
// GET /groups/{group_id}/org_memberships
func (c *Client) ListOrgMemberships(ctx context.Context, groupID string, params *rest.PaginationParams) (*OrgMembershipListResponse, error) {
	path := fmt.Sprintf("/groups/%s/org_memberships", groupID)
	var result OrgMembershipListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Organizations
// ============================================================================

// Organization represents an organization in a group.
type Organization struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Attributes OrgAttrs `json:"attributes"`
}

// OrgAttrs contains organization attributes.
type OrgAttrs struct {
	Name    string     `json:"name"`
	Slug    string     `json:"slug,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// OrgListResponse wraps a list of organizations.
type OrgListResponse struct {
	Data  []Organization `json:"data"`
	Links rest.Links     `json:"links,omitempty"`
}

// ListOrgs lists all organizations in a group.
// GET /groups/{group_id}/orgs
func (c *Client) ListOrgs(ctx context.Context, groupID string, params *rest.PaginationParams) (*OrgListResponse, error) {
	path := fmt.Sprintf("/groups/%s/orgs", groupID)
	var result OrgListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Service Accounts
// ============================================================================

// ServiceAccount represents a group service account.
type ServiceAccount struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes ServiceAccountAttrs `json:"attributes"`
}

// ServiceAccountAttrs contains service account attributes.
type ServiceAccountAttrs struct {
	Name         string     `json:"name"`
	AuthType     string     `json:"auth_type,omitempty"`
	RolePublicID string     `json:"role_public_id,omitempty"`
	Created      *time.Time `json:"created_at,omitempty"`
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

// ListServiceAccounts lists all service accounts for a group.
// GET /groups/{group_id}/service_accounts
func (c *Client) ListServiceAccounts(ctx context.Context, groupID string, params *rest.PaginationParams) (*ServiceAccountListResponse, error) {
	path := fmt.Sprintf("/groups/%s/service_accounts", groupID)
	var result ServiceAccountListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServiceAccount retrieves a service account by ID.
// GET /groups/{group_id}/service_accounts/{serviceaccount_id}
func (c *Client) GetServiceAccount(ctx context.Context, groupID, saID string) (*ServiceAccountResponse, error) {
	path := fmt.Sprintf("/groups/%s/service_accounts/%s", groupID, saID)
	var result ServiceAccountResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteServiceAccount deletes a service account.
// DELETE /groups/{group_id}/service_accounts/{serviceaccount_id}
func (c *Client) DeleteServiceAccount(ctx context.Context, groupID, saID string) error {
	path := fmt.Sprintf("/groups/%s/service_accounts/%s", groupID, saID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Policies
// ============================================================================

// Policy represents a group policy.
type Policy struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes PolicyAttrs `json:"attributes"`
}

// PolicyAttrs contains policy attributes.
type PolicyAttrs struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
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

// ListPolicies lists all policies for a group.
// GET /groups/{group_id}/policies
func (c *Client) ListPolicies(ctx context.Context, groupID string, params *rest.PaginationParams) (*PolicyListResponse, error) {
	path := fmt.Sprintf("/groups/%s/policies", groupID)
	var result PolicyListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPolicy retrieves a policy by ID.
// GET /groups/{group_id}/policies/{policy_id}
func (c *Client) GetPolicy(ctx context.Context, groupID, policyID string) (*PolicyResponse, error) {
	path := fmt.Sprintf("/groups/%s/policies/%s", groupID, policyID)
	var result PolicyResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePolicy deletes a policy.
// DELETE /groups/{group_id}/policies/{policy_id}
func (c *Client) DeletePolicy(ctx context.Context, groupID, policyID string) error {
	path := fmt.Sprintf("/groups/%s/policies/%s", groupID, policyID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// Settings
// ============================================================================

// IaCSettings represents group IaC settings.
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

// GetIaCSettings retrieves IaC settings for a group.
// GET /groups/{group_id}/settings/iac
func (c *Client) GetIaCSettings(ctx context.Context, groupID string) (*IaCSettingsResponse, error) {
	path := fmt.Sprintf("/groups/%s/settings/iac", groupID)
	var result IaCSettingsResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PRTemplateSettings represents pull request template settings.
type PRTemplateSettings struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes PRTemplateAttrs `json:"attributes"`
}

// PRTemplateAttrs contains PR template attributes.
type PRTemplateAttrs struct {
	IsEnabled    bool   `json:"is_enabled"`
	TemplateBody string `json:"template_body,omitempty"`
}

// PRTemplateSettingsResponse wraps PR template settings.
type PRTemplateSettingsResponse struct {
	Data PRTemplateSettings `json:"data"`
}

// GetPRTemplateSettings retrieves PR template settings for a group.
// GET /groups/{group_id}/settings/pull_request_template
func (c *Client) GetPRTemplateSettings(ctx context.Context, groupID string) (*PRTemplateSettingsResponse, error) {
	path := fmt.Sprintf("/groups/%s/settings/pull_request_template", groupID)
	var result PRTemplateSettingsResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// SSO Connections
// ============================================================================

// SSOConnection represents an SSO connection.
type SSOConnection struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes SSOConnectionAttrs `json:"attributes"`
}

// SSOConnectionAttrs contains SSO connection attributes.
type SSOConnectionAttrs struct {
	Name     string     `json:"name"`
	Provider string     `json:"provider,omitempty"`
	Created  *time.Time `json:"created_at,omitempty"`
}

// SSOConnectionListResponse wraps a list of SSO connections.
type SSOConnectionListResponse struct {
	Data  []SSOConnection `json:"data"`
	Links rest.Links      `json:"links,omitempty"`
}

// ListSSOConnections lists all SSO connections for a group.
// GET /groups/{group_id}/sso_connections
func (c *Client) ListSSOConnections(ctx context.Context, groupID string, params *rest.PaginationParams) (*SSOConnectionListResponse, error) {
	path := fmt.Sprintf("/groups/%s/sso_connections", groupID)
	var result SSOConnectionListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SSOUser represents an SSO user.
type SSOUser struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Attributes SSOUserAttrs `json:"attributes"`
}

// SSOUserAttrs contains SSO user attributes.
type SSOUserAttrs struct {
	Email   string     `json:"email"`
	Name    string     `json:"name,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// SSOUserListResponse wraps a list of SSO users.
type SSOUserListResponse struct {
	Data  []SSOUser  `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// ListSSOUsers lists all users for an SSO connection.
// GET /groups/{group_id}/sso_connections/{sso_id}/users
func (c *Client) ListSSOUsers(ctx context.Context, groupID, ssoID string, params *rest.PaginationParams) (*SSOUserListResponse, error) {
	path := fmt.Sprintf("/groups/%s/sso_connections/%s/users", groupID, ssoID)
	var result SSOUserListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteSSOUser removes an SSO user.
// DELETE /groups/{group_id}/sso_connections/{sso_id}/users/{user_id}
func (c *Client) DeleteSSOUser(ctx context.Context, groupID, ssoID, userID string) error {
	path := fmt.Sprintf("/groups/%s/sso_connections/%s/users/%s", groupID, ssoID, userID)
	return c.base.Delete(ctx, path, nil)
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

// SearchAuditLogs searches audit logs for a group.
// GET /groups/{group_id}/audit_logs/search
func (c *Client) SearchAuditLogs(ctx context.Context, groupID string, params *rest.PaginationParams) (*AuditLogListResponse, error) {
	path := fmt.Sprintf("/groups/%s/audit_logs/search", groupID)
	var result AuditLogListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Issues
// ============================================================================

// Issue represents a group issue.
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

// ListIssues lists all issues for a group.
// GET /groups/{group_id}/issues
func (c *Client) ListIssues(ctx context.Context, groupID string, params *ListIssuesParams) (*IssueListResponse, error) {
	path := fmt.Sprintf("/groups/%s/issues", groupID)
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
// GET /groups/{group_id}/issues/{issue_id}
func (c *Client) GetIssue(ctx context.Context, groupID, issueID string) (*IssueResponse, error) {
	path := fmt.Sprintf("/groups/%s/issues/%s", groupID, issueID)
	var result IssueResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Assets
// ============================================================================

// Asset represents an asset.
type Asset struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes AssetAttrs `json:"attributes"`
}

// AssetAttrs contains asset attributes.
type AssetAttrs struct {
	Name      string     `json:"name"`
	AssetType string     `json:"asset_type,omitempty"`
	Created   *time.Time `json:"created_at,omitempty"`
}

// AssetListResponse wraps a list of assets.
type AssetListResponse struct {
	Data  []Asset    `json:"data"`
	Links rest.Links `json:"links,omitempty"`
}

// AssetResponse wraps a single asset.
type AssetResponse struct {
	Data Asset `json:"data"`
}

// SearchAssets searches assets for a group.
// GET /groups/{group_id}/assets/search
func (c *Client) SearchAssets(ctx context.Context, groupID string, params *rest.PaginationParams) (*AssetListResponse, error) {
	path := fmt.Sprintf("/groups/%s/assets/search", groupID)
	var result AssetListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAsset retrieves an asset by ID.
// GET /groups/{group_id}/assets/{asset_id}
func (c *Client) GetAsset(ctx context.Context, groupID, assetID string) (*AssetResponse, error) {
	path := fmt.Sprintf("/groups/%s/assets/%s", groupID, assetID)
	var result AssetResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// App Installs
// ============================================================================

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

// ListAppInstalls lists all app installations for a group.
// GET /groups/{group_id}/apps/installs
func (c *Client) ListAppInstalls(ctx context.Context, groupID string, params *rest.PaginationParams) (*AppInstallListResponse, error) {
	path := fmt.Sprintf("/groups/%s/apps/installs", groupID)
	var result AppInstallListResponse
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
	Created     *time.Time `json:"created_at,omitempty"`
	DownloadURL string     `json:"download_url,omitempty"`
}

// ExportJobResponse wraps an export job.
type ExportJobResponse struct {
	Data ExportJob `json:"data"`
}

// CreateExport creates a new export job.
// POST /groups/{group_id}/export
func (c *Client) CreateExport(ctx context.Context, groupID string) (*ExportJobResponse, error) {
	path := fmt.Sprintf("/groups/%s/export", groupID)
	var result ExportJobResponse
	if err := c.base.Post(ctx, path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExport retrieves an export job status.
// GET /groups/{group_id}/export/{export_id}
func (c *Client) GetExport(ctx context.Context, groupID, exportID string) (*ExportJobResponse, error) {
	path := fmt.Sprintf("/groups/%s/export/%s", groupID, exportID)
	var result ExportJobResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUser retrieves a user in a group.
// GET /groups/{group_id}/users/{id}
func (c *Client) GetUser(ctx context.Context, groupID, userID string) (*MembershipResponse, error) {
	path := fmt.Sprintf("/groups/%s/users/%s", groupID, userID)
	var result MembershipResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
