// Package webhooks provides a client for Snyk v1 Webhooks API.
package webhooks

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Webhooks API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Webhooks client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// Webhook represents a webhook.
type Webhook struct {
	ID      string    `json:"id"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

// WebhookList represents a list of webhooks.
type WebhookList struct {
	Results []Webhook `json:"results"`
	Total   int       `json:"total"`
}

// CreateWebhookRequest represents a request to create a webhook.
type CreateWebhookRequest struct {
	URL    string `json:"url"`
	Secret string `json:"secret,omitempty"`
}

// ============================================================================
// CRUD Operations
// ============================================================================

// List lists all webhooks for an organization.
// GET /org/{orgId}/webhooks
func (c *Client) List(ctx context.Context, orgID string) (*WebhookList, error) {
	path := fmt.Sprintf("/org/%s/webhooks", orgID)
	var result WebhookList
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new webhook.
// POST /org/{orgId}/webhooks
func (c *Client) Create(ctx context.Context, orgID string, req *CreateWebhookRequest) (*Webhook, error) {
	path := fmt.Sprintf("/org/%s/webhooks", orgID)
	var result Webhook
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a webhook by ID.
// GET /org/{orgId}/webhooks/{webhookId}
func (c *Client) Get(ctx context.Context, orgID, webhookID string) (*Webhook, error) {
	path := fmt.Sprintf("/org/%s/webhooks/%s", orgID, webhookID)
	var result Webhook
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a webhook.
// DELETE /org/{orgId}/webhooks/{webhookId}
func (c *Client) Delete(ctx context.Context, orgID, webhookID string) error {
	path := fmt.Sprintf("/org/%s/webhooks/%s", orgID, webhookID)
	return c.base.Delete(ctx, path)
}

// Ping triggers a test ping to a webhook.
// POST /org/{orgId}/webhooks/{webhookId}/ping
func (c *Client) Ping(ctx context.Context, orgID, webhookID string) error {
	path := fmt.Sprintf("/org/%s/webhooks/%s/ping", orgID, webhookID)
	return c.base.Post(ctx, path, nil, nil)
}
