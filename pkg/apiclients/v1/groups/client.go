// Package groups provides a client for Snyk v1 Groups API.
package groups

import (
	"context"
	"fmt"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Groups API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Groups client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// Group represents a Snyk group.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// GroupSettings represents group settings.
type GroupSettings struct {
	SessionLength        *int  `json:"sessionLength,omitempty"`
	IgnorePublicExploits *bool `json:"ignorePublicExploits,omitempty"`
}

// Member represents a group member.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	OrgName  string `json:"orgName,omitempty"` // For org members
}

// MemberList represents a list of members.
type MemberList []Member

// Tag represents a group tag.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TagList represents a list of tags.
type TagList struct {
	Tags []Tag `json:"tags"`
}

// Organization represents an organization in a group.
type Organization struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug,omitempty"`
	URL   string `json:"url,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// OrganizationList represents a list of organizations.
type OrganizationList struct {
	Orgs []Organization `json:"orgs"`
}

// Role represents a group role.
type Role struct {
	PublicID    string `json:"publicId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

// RoleList represents a list of roles.
type RoleList []Role

// ============================================================================
// Settings
// ============================================================================

// GetSettings retrieves group settings.
// GET /group/{groupId}/settings
func (c *Client) GetSettings(ctx context.Context, groupID string) (*GroupSettings, error) {
	path := fmt.Sprintf("/group/%s/settings", groupID)
	var result GroupSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates group settings.
// PUT /group/{groupId}/settings
func (c *Client) UpdateSettings(ctx context.Context, groupID string, settings *GroupSettings) (*GroupSettings, error) {
	path := fmt.Sprintf("/group/%s/settings", groupID)
	var result GroupSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Members
// ============================================================================

// ListMembers lists all members of a group.
// GET /group/{groupId}/members
func (c *Client) ListMembers(ctx context.Context, groupID string) (MemberList, error) {
	path := fmt.Sprintf("/group/%s/members", groupID)
	var result MemberList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ListOrgMembers lists all members of an org within a group.
// GET /group/{groupId}/org/{orgId}/members
func (c *Client) ListOrgMembers(ctx context.Context, groupID, orgID string) (MemberList, error) {
	path := fmt.Sprintf("/group/%s/org/%s/members", groupID, orgID)
	var result MemberList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ============================================================================
// Tags
// ============================================================================

// ListTags lists all tags for a group.
// GET /group/{groupId}/tags
func (c *Client) ListTags(ctx context.Context, groupID string, perPage, page int) (*TagList, error) {
	path := fmt.Sprintf("/group/%s/tags", groupID)
	if perPage > 0 || page > 0 {
		path = fmt.Sprintf("%s?perPage=%d&page=%d", path, perPage, page)
	}
	var result TagList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteTagsRequest represents a request to delete tags.
type DeleteTagsRequest struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"` // If empty, deletes all values for the key
	Force bool   `json:"force,omitempty"` // Force delete even if projects use the tag
}

// DeleteTags deletes tags from a group.
// POST /group/{groupId}/tags/delete
func (c *Client) DeleteTags(ctx context.Context, groupID string, req *DeleteTagsRequest) error {
	path := fmt.Sprintf("/group/%s/tags/delete", groupID)
	return c.base.Post(ctx, path, req, nil)
}

// ============================================================================
// Organizations
// ============================================================================

// ListOrgs lists all organizations in a group.
// GET /group/{groupId}/orgs
func (c *Client) ListOrgs(ctx context.Context, groupID string, perPage, page int) (*OrganizationList, error) {
	path := fmt.Sprintf("/group/%s/orgs", groupID)
	if perPage > 0 || page > 0 {
		path = fmt.Sprintf("%s?perPage=%d&page=%d", path, perPage, page)
	}
	var result OrganizationList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Roles
// ============================================================================

// ListRoles lists all roles in a group.
// GET /group/{groupId}/roles
func (c *Client) ListRoles(ctx context.Context, groupID string) (RoleList, error) {
	path := fmt.Sprintf("/group/%s/roles", groupID)
	var result RoleList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}
