// Package monitor provides a client for Snyk v1 Monitor API.
package monitor

import (
	"context"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	v1projects "github.com/sam1el/snyk-api/pkg/apiclients/v1/projects"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Monitor API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Monitor client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// MonitorRequest represents a monitor request.
type MonitorRequest struct {
	DepGraph               v1projects.DepGraph `json:"depGraph"`
	Policy                 string              `json:"policy,omitempty"`
	TargetFile             string              `json:"targetFile,omitempty"`
	TargetFileRelativePath string              `json:"targetFileRelativePath,omitempty"`
	TargetReference        string              `json:"targetReference,omitempty"` // e.g., branch name
	RemoteRepoURL          string              `json:"remoteRepoUrl,omitempty"`
	Meta                   *MonitorMeta        `json:"meta,omitempty"`
}

// MonitorMeta represents metadata for monitoring.
type MonitorMeta struct {
	Method            string            `json:"method,omitempty"` // e.g., "cli"
	Hostname          string            `json:"hostname,omitempty"`
	ID                string            `json:"id,omitempty"`
	CI                bool              `json:"ci,omitempty"`
	PID               int               `json:"pid,omitempty"`
	NodeVersion       string            `json:"nodeVersion,omitempty"`
	Name              string            `json:"name,omitempty"`
	Version           string            `json:"version,omitempty"`
	Org               string            `json:"org,omitempty"`
	PluginName        string            `json:"pluginName,omitempty"`
	PluginRuntime     string            `json:"pluginRuntime,omitempty"`
	DockerImageID     string            `json:"dockerImageId,omitempty"`
	ProjectName       string            `json:"projectName,omitempty"`
	GradleProjectName string            `json:"gradleProjectName,omitempty"`
	PruneRepeated     bool              `json:"prune,omitempty"`
	Policy            string            `json:"policy,omitempty"`
	IsDocker          bool              `json:"isDocker,omitempty"`
	Platform          string            `json:"platform,omitempty"`
	MissingDeps       bool              `json:"missingDeps,omitempty"`
	Args              map[string]string `json:"args,omitempty"`
}

// MonitorResponse represents a monitor response.
type MonitorResponse struct {
	ID             string      `json:"id"`
	URI            string      `json:"uri"`
	IsMonitored    bool        `json:"isMonitored"`
	TrialStarted   bool        `json:"trialStarted,omitempty"`
	LicensesPolicy interface{} `json:"licensesPolicy,omitempty"`
	Path           string      `json:"path,omitempty"`
	ProjectName    string      `json:"projectName,omitempty"`
}

// ============================================================================
// Monitor
// ============================================================================

// MonitorDepGraph monitors a dependency graph.
// POST /monitor/dep-graph
func (c *Client) MonitorDepGraph(ctx context.Context, orgID string, req *MonitorRequest) (*MonitorResponse, error) {
	path := "/monitor/dep-graph"
	if orgID != "" {
		path = path + "?org=" + orgID
	}

	var result MonitorResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
