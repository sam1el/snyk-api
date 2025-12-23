// Package self provides a manual client for Snyk REST Self API.
// These endpoints relate to the current authenticated user.
package self

import (
	"context"
	"fmt"
	"time"

	"github.com/sam1el/snyk-api/pkg/apiclients/rest"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to REST Self API endpoints.
type Client struct {
	base *rest.BaseClient
}

// New creates a new REST Self client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: rest.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Self
// ============================================================================

// User represents the current user.
type User struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Attributes UserAttrs `json:"attributes"`
}

// UserAttrs contains user attributes.
type UserAttrs struct {
	Name      string     `json:"name"`
	Email     string     `json:"email,omitempty"`
	Username  string     `json:"username,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Created   *time.Time `json:"created_at,omitempty"`
}

// UserResponse wraps the current user.
type UserResponse struct {
	Data User `json:"data"`
}

// Get retrieves the current user.
// GET /self
func (c *Client) Get(ctx context.Context) (*UserResponse, error) {
	var result UserResponse
	if err := c.base.Get(ctx, "/self", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Access Requests
// ============================================================================

// AccessRequest represents an access request.
type AccessRequest struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes AccessRequestAttrs `json:"attributes"`
}

// AccessRequestAttrs contains access request attributes.
type AccessRequestAttrs struct {
	OrgID   string     `json:"org_id"`
	OrgName string     `json:"org_name,omitempty"`
	Status  string     `json:"status,omitempty"`
	Created *time.Time `json:"created_at,omitempty"`
}

// AccessRequestListResponse wraps a list of access requests.
type AccessRequestListResponse struct {
	Data  []AccessRequest `json:"data"`
	Links rest.Links      `json:"links,omitempty"`
}

// ListAccessRequests lists all access requests for the current user.
// GET /self/access_requests
func (c *Client) ListAccessRequests(ctx context.Context, params *rest.PaginationParams) (*AccessRequestListResponse, error) {
	var result AccessRequestListResponse
	if err := c.base.Get(ctx, "/self/access_requests", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Apps
// ============================================================================

// App represents a user app.
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

// ListApps lists all apps for the current user.
// GET /self/apps
func (c *Client) ListApps(ctx context.Context, params *rest.PaginationParams) (*AppListResponse, error) {
	var result AppListResponse
	if err := c.base.Get(ctx, "/self/apps", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetApp retrieves an app by ID.
// GET /self/apps/{app_id}
func (c *Client) GetApp(ctx context.Context, appID string) (*AppResponse, error) {
	path := fmt.Sprintf("/self/apps/%s", appID)
	var result AppResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteApp deletes an app.
// DELETE /self/apps/{app_id}
func (c *Client) DeleteApp(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/self/apps/%s", appID)
	return c.base.Delete(ctx, path, nil)
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

// AppInstallResponse wraps a single app installation.
type AppInstallResponse struct {
	Data AppInstall `json:"data"`
}

// ListAppInstalls lists all app installations for the current user.
// GET /self/apps/installs
func (c *Client) ListAppInstalls(ctx context.Context, params *rest.PaginationParams) (*AppInstallListResponse, error) {
	var result AppInstallListResponse
	if err := c.base.Get(ctx, "/self/apps/installs", params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAppInstall retrieves an app installation by ID.
// GET /self/apps/installs/{install_id}
func (c *Client) GetAppInstall(ctx context.Context, installID string) (*AppInstallResponse, error) {
	path := fmt.Sprintf("/self/apps/installs/%s", installID)
	var result AppInstallResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeAppInstall revokes an app installation.
// DELETE /self/apps/installs/{install_id}
func (c *Client) RevokeAppInstall(ctx context.Context, installID string) error {
	path := fmt.Sprintf("/self/apps/installs/%s", installID)
	return c.base.Delete(ctx, path, nil)
}

// ============================================================================
// App Sessions
// ============================================================================

// AppSession represents an app session.
type AppSession struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes AppSessionAttrs `json:"attributes"`
}

// AppSessionAttrs contains app session attributes.
type AppSessionAttrs struct {
	LastActive *time.Time `json:"last_active,omitempty"`
	Created    *time.Time `json:"created_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// AppSessionListResponse wraps a list of app sessions.
type AppSessionListResponse struct {
	Data  []AppSession `json:"data"`
	Links rest.Links   `json:"links,omitempty"`
}

// AppSessionResponse wraps a single app session.
type AppSessionResponse struct {
	Data AppSession `json:"data"`
}

// ListAppSessions lists all sessions for an app.
// GET /self/apps/{app_id}/sessions
func (c *Client) ListAppSessions(ctx context.Context, appID string, params *rest.PaginationParams) (*AppSessionListResponse, error) {
	path := fmt.Sprintf("/self/apps/%s/sessions", appID)
	var result AppSessionListResponse
	if err := c.base.Get(ctx, path, params.ToQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAppSession retrieves an app session by ID.
// GET /self/apps/{app_id}/sessions/{session_id}
func (c *Client) GetAppSession(ctx context.Context, appID, sessionID string) (*AppSessionResponse, error) {
	path := fmt.Sprintf("/self/apps/%s/sessions/%s", appID, sessionID)
	var result AppSessionResponse
	if err := c.base.Get(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RevokeAppSession revokes an app session.
// DELETE /self/apps/{app_id}/sessions/{session_id}
func (c *Client) RevokeAppSession(ctx context.Context, appID, sessionID string) error {
	path := fmt.Sprintf("/self/apps/%s/sessions/%s", appID, sessionID)
	return c.base.Delete(ctx, path, nil)
}
