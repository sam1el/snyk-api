// Package users provides a client for Snyk v1 Users API.
package users

import (
	"context"
	"fmt"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Users API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Users client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// User represents a Snyk user.
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Orgs     []UserOrg `json:"orgs,omitempty"`
}

// UserOrg represents an organization a user belongs to.
type UserOrg struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug,omitempty"`
	URL   string `json:"url,omitempty"`
	Group *Group `json:"group,omitempty"`
}

// Group represents a group the org belongs to.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NotificationSettings represents notification settings.
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

// ============================================================================
// User Endpoints
// ============================================================================

// GetMe retrieves the current user's details.
// GET /user/me
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var result User
	if err := c.base.Get(ctx, "/user/me", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a user by ID.
// GET /user/{userId}
func (c *Client) Get(ctx context.Context, userID string) (*User, error) {
	path := fmt.Sprintf("/user/%s", userID)
	var result User
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Notification Settings
// ============================================================================

// GetOrgNotificationSettings retrieves notification settings for an org.
// GET /user/me/notification-settings/org/{orgId}
func (c *Client) GetOrgNotificationSettings(ctx context.Context, orgID string) (*NotificationSettings, error) {
	path := fmt.Sprintf("/user/me/notification-settings/org/%s", orgID)
	var result NotificationSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrgNotificationSettings updates notification settings for an org.
// PUT /user/me/notification-settings/org/{orgId}
func (c *Client) UpdateOrgNotificationSettings(ctx context.Context, orgID string, settings *NotificationSettings) (*NotificationSettings, error) {
	path := fmt.Sprintf("/user/me/notification-settings/org/%s", orgID)
	var result NotificationSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProjectNotificationSettings retrieves notification settings for a project.
// GET /user/me/notification-settings/org/{orgId}/project/{projectId}
func (c *Client) GetProjectNotificationSettings(ctx context.Context, orgID, projectID string) (*NotificationSettings, error) {
	path := fmt.Sprintf("/user/me/notification-settings/org/%s/project/%s", orgID, projectID)
	var result NotificationSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProjectNotificationSettings updates notification settings for a project.
// PUT /user/me/notification-settings/org/{orgId}/project/{projectId}
func (c *Client) UpdateProjectNotificationSettings(ctx context.Context, orgID, projectID string, settings *NotificationSettings) (*NotificationSettings, error) {
	path := fmt.Sprintf("/user/me/notification-settings/org/%s/project/%s", orgID, projectID)
	var result NotificationSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
