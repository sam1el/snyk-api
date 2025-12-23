// Package orgs provides a client for Snyk v1 Organizations API.
package orgs

import (
	"context"
	"fmt"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Organizations API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Organizations client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// Organization represents a Snyk organization in v1 API.
type Organization struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug,omitempty"`
	URL   string `json:"url,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// Group represents a Snyk group.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrganizationList represents a list of organizations.
type OrganizationList struct {
	Orgs []Organization `json:"orgs"`
}

// ============================================================================
// List Organizations
// ============================================================================

// List lists all organizations the user has access to.
// GET /orgs
func (c *Client) List(ctx context.Context) (*OrganizationList, error) {
	var result OrganizationList
	if err := c.base.Get(ctx, "/orgs", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Organization Settings
// ============================================================================

// NotificationSettings represents organization notification settings.
type NotificationSettings struct {
	NewIssuesRemediations *NotificationSetting `json:"new-issues-remediations,omitempty"`
	ProjectImported       *NotificationSetting `json:"project-imported,omitempty"`
	TestLimit             *NotificationSetting `json:"test-limit,omitempty"`
	WeeklyReport          *NotificationSetting `json:"weekly-report,omitempty"`
}

// NotificationSetting represents a single notification setting.
type NotificationSetting struct {
	Enabled       bool   `json:"enabled"`
	IssueSeverity string `json:"issueSeverity,omitempty"`
	IssueType     string `json:"issueType,omitempty"`
	Inherited     bool   `json:"inherited,omitempty"`
}

// GetNotificationSettings retrieves organization notification settings.
// GET /org/{orgId}/notification-settings
func (c *Client) GetNotificationSettings(ctx context.Context, orgID string) (*NotificationSettings, error) {
	path := fmt.Sprintf("/org/%s/notification-settings", orgID)
	var result NotificationSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateNotificationSettings updates organization notification settings.
// PUT /org/{orgId}/notification-settings
func (c *Client) UpdateNotificationSettings(ctx context.Context, orgID string, settings *NotificationSettings) (*NotificationSettings, error) {
	path := fmt.Sprintf("/org/%s/notification-settings", orgID)
	var result NotificationSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Members
// ============================================================================

// Member represents an organization member.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// MemberList represents a list of members.
type MemberList []Member

// ListMembers lists all members of an organization.
// GET /org/{orgId}/members
func (c *Client) ListMembers(ctx context.Context, orgID string, includeGroupAdmins bool) (MemberList, error) {
	path := fmt.Sprintf("/org/%s/members", orgID)
	if includeGroupAdmins {
		path += "?includeGroupAdmins=true"
	}
	var result MemberList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateMemberRoleRequest represents the request to update a member's role.
type UpdateMemberRoleRequest struct {
	Role string `json:"role"` // "admin" or "collaborator"
}

// UpdateMemberRole updates a member's role in an organization.
// PUT /org/{orgId}/members/update/{userId}
func (c *Client) UpdateMemberRole(ctx context.Context, orgID, userID string, role string) error {
	path := fmt.Sprintf("/org/%s/members/update/%s", orgID, userID)
	req := &UpdateMemberRoleRequest{Role: role}
	return c.base.Put(ctx, path, req, nil)
}

// RemoveMember removes a member from an organization.
// DELETE /org/{orgId}/members/{userId}
func (c *Client) RemoveMember(ctx context.Context, orgID, userID string) error {
	path := fmt.Sprintf("/org/%s/members/%s", orgID, userID)
	return c.base.Delete(ctx, path)
}

// ============================================================================
// Invites
// ============================================================================

// InviteRequest represents an invitation request.
type InviteRequest struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin,omitempty"`
}

// Invite invites a user to an organization.
// POST /org/{orgId}/invite
func (c *Client) Invite(ctx context.Context, orgID string, req *InviteRequest) error {
	path := fmt.Sprintf("/org/%s/invite", orgID)
	return c.base.Post(ctx, path, req, nil)
}

// ============================================================================
// Organization Settings
// ============================================================================

// OrgSettings represents organization settings.
type OrgSettings struct {
	RequestAccess *RequestAccessSettings `json:"requestAccess,omitempty"`
}

// RequestAccessSettings represents request access settings.
type RequestAccessSettings struct {
	Enabled bool `json:"enabled"`
}

// GetSettings retrieves organization settings.
// GET /org/{orgId}/settings
func (c *Client) GetSettings(ctx context.Context, orgID string) (*OrgSettings, error) {
	path := fmt.Sprintf("/org/%s/settings", orgID)
	var result OrgSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates organization settings.
// PUT /org/{orgId}/settings
func (c *Client) UpdateSettings(ctx context.Context, orgID string, settings *OrgSettings) (*OrgSettings, error) {
	path := fmt.Sprintf("/org/%s/settings", orgID)
	var result OrgSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Provision
// ============================================================================

// ProvisionRequest represents a provision request.
type ProvisionRequest struct {
	Name string `json:"name"`
}

// ProvisionResponse represents the provision response.
type ProvisionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Provision provisions a new organization.
// POST /org/{orgId}/provision
func (c *Client) Provision(ctx context.Context, orgID string, req *ProvisionRequest) (*ProvisionResponse, error) {
	path := fmt.Sprintf("/org/%s/provision", orgID)
	var result ProvisionResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Dependencies
// ============================================================================

// DependenciesParams represents parameters for listing dependencies.
type DependenciesParams struct {
	SortBy  string `json:"sortBy,omitempty"` // "dependency" or "severity"
	Order   string `json:"order,omitempty"`  // "asc" or "desc"
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"perPage,omitempty"`
}

// DependenciesResponse represents the dependencies response.
type DependenciesResponse struct {
	Results []Dependency `json:"results"`
	Total   int          `json:"total"`
}

// Dependency represents a dependency.
type Dependency struct {
	ID                         string       `json:"id"`
	Name                       string       `json:"name"`
	Version                    string       `json:"version"`
	LatestVersion              string       `json:"latestVersion,omitempty"`
	LatestVersionPublishedDate string       `json:"latestVersionPublishedDate,omitempty"`
	FirstPublishedDate         string       `json:"firstPublishedDate,omitempty"`
	IsDeprecated               bool         `json:"isDeprecated"`
	DeprecatedVersions         []string     `json:"deprecatedVersions,omitempty"`
	DependenciesWithIssues     []string     `json:"dependenciesWithIssues,omitempty"`
	Type                       string       `json:"type"`
	Projects                   []ProjectRef `json:"projects,omitempty"`
	IssuesCritical             int          `json:"issuesCritical"`
	IssuesHigh                 int          `json:"issuesHigh"`
	IssuesMedium               int          `json:"issuesMedium"`
	IssuesLow                  int          `json:"issuesLow"`
	Licenses                   []License    `json:"licenses,omitempty"`
}

// ProjectRef is a reference to a project.
type ProjectRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// License represents a license.
type License struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity,omitempty"`
}

// ListDependencies lists all dependencies for an organization.
// POST /org/{orgId}/dependencies
func (c *Client) ListDependencies(ctx context.Context, orgID string, params *DependenciesParams) (*DependenciesResponse, error) {
	path := fmt.Sprintf("/org/%s/dependencies", orgID)
	if params == nil {
		params = &DependenciesParams{}
	}
	var result DependenciesResponse
	if err := c.base.Post(ctx, path, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Licenses
// ============================================================================

// LicensesParams represents parameters for listing licenses.
type LicensesParams struct {
	SortBy  string `json:"sortBy,omitempty"`
	Order   string `json:"order,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"perPage,omitempty"`
}

// LicensesResponse represents the licenses response.
type LicensesResponse struct {
	Results []LicenseResult `json:"results"`
	Total   int             `json:"total"`
}

// LicenseResult represents a license result.
type LicenseResult struct {
	ID           string       `json:"id"`
	Severity     string       `json:"severity"`
	Instructions string       `json:"instructions,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
	Projects     []ProjectRef `json:"projects"`
}

// ListLicenses lists all licenses for an organization.
// POST /org/{orgId}/licenses
func (c *Client) ListLicenses(ctx context.Context, orgID string, params *LicensesParams) (*LicensesResponse, error) {
	path := fmt.Sprintf("/org/%s/licenses", orgID)
	if params == nil {
		params = &LicensesParams{}
	}
	var result LicensesResponse
	if err := c.base.Post(ctx, path, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Entitlements
// ============================================================================

// Entitlements represents organization entitlements.
type Entitlements map[string]bool

// GetEntitlements retrieves organization entitlements.
// GET /org/{orgId}/entitlements
func (c *Client) GetEntitlements(ctx context.Context, orgID string) (Entitlements, error) {
	path := fmt.Sprintf("/org/%s/entitlements", orgID)
	var result Entitlements
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetEntitlement retrieves a specific entitlement.
// GET /org/{orgId}/entitlement/{entitlementKey}
func (c *Client) GetEntitlement(ctx context.Context, orgID, entitlementKey string) (bool, error) {
	path := fmt.Sprintf("/org/%s/entitlement/%s", orgID, entitlementKey)
	var result bool
	if err := c.base.Get(ctx, path, &result); err != nil {
		return false, err
	}
	return result, nil
}
