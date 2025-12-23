// Package projects provides a client for Snyk v1 Projects API.
// This includes high-value endpoints like aggregated-issues, ignores, and dep-graph.
package projects

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Projects API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Projects client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// Project represents a Snyk project in v1 API.
type Project struct {
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	Created               time.Time              `json:"created"`
	Origin                string                 `json:"origin"`
	Type                  string                 `json:"type"`
	ReadOnly              bool                   `json:"readOnly"`
	TestFrequency         string                 `json:"testFrequency"`
	TotalDependencies     int                    `json:"totalDependencies"`
	IssueCountsBySeverity IssueCountsBySeverity  `json:"issueCountsBySeverity"`
	ImageTag              string                 `json:"imageTag,omitempty"`
	ImageID               string                 `json:"imageId,omitempty"`
	LastTestedDate        *time.Time             `json:"lastTestedDate,omitempty"`
	BrowseURL             string                 `json:"browseUrl"`
	ImportingUser         *User                  `json:"importingUser,omitempty"`
	Owner                 *User                  `json:"owner,omitempty"`
	Tags                  []Tag                  `json:"tags,omitempty"`
	Attributes            map[string]interface{} `json:"attributes,omitempty"`
	Branch                string                 `json:"branch,omitempty"`
}

// IssueCountsBySeverity contains issue counts by severity level.
type IssueCountsBySeverity struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

// User represents a Snyk user.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Tag represents a project tag.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ============================================================================
// Aggregated Issues (High Value)
// ============================================================================

// AggregatedIssuesRequest represents the request body for aggregated issues.
type AggregatedIssuesRequest struct {
	IncludeDescription       bool               `json:"includeDescription,omitempty"`
	IncludeIntroducedThrough bool               `json:"includeIntroducedThrough,omitempty"`
	Filters                  *AggregatedFilters `json:"filters,omitempty"`
}

// AggregatedFilters represents filters for aggregated issues.
type AggregatedFilters struct {
	Severities      []string        `json:"severities,omitempty"`
	ExploitMaturity []string        `json:"exploitMaturity,omitempty"`
	Types           []string        `json:"types,omitempty"`
	Ignored         bool            `json:"ignored,omitempty"`
	Patched         bool            `json:"patched,omitempty"`
	Priority        *PriorityFilter `json:"priority,omitempty"`
}

// PriorityFilter represents priority score filtering.
type PriorityFilter struct {
	Score *ScoreFilter `json:"score,omitempty"`
}

// ScoreFilter represents a score range filter.
type ScoreFilter struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

// AggregatedIssuesResponse represents the aggregated issues response.
type AggregatedIssuesResponse struct {
	Issues []AggregatedIssue `json:"issues"`
}

// AggregatedIssue represents an aggregated issue.
type AggregatedIssue struct {
	ID                string              `json:"id"`
	IssueType         string              `json:"issueType"`
	PkgName           string              `json:"pkgName"`
	PkgVersions       []string            `json:"pkgVersions"`
	IssueData         IssueData           `json:"issueData"`
	IsPatched         bool                `json:"isPatched"`
	IsIgnored         bool                `json:"isIgnored"`
	IgnoreReasons     []IgnoreReason      `json:"ignoreReasons,omitempty"`
	FixInfo           FixInfo             `json:"fixInfo"`
	Priority          *Priority           `json:"priority,omitempty"`
	Links             map[string]string   `json:"links,omitempty"`
	IntroducedThrough []IntroducedThrough `json:"introducedThrough,omitempty"`
}

// IssueData contains detailed issue information.
type IssueData struct {
	ID                 string       `json:"id"`
	Title              string       `json:"title"`
	Severity           string       `json:"severity"`
	OriginalSeverity   string       `json:"originalSeverity,omitempty"`
	URL                string       `json:"url"`
	Description        string       `json:"description,omitempty"`
	Identifiers        *Identifiers `json:"identifiers,omitempty"`
	Credit             []string     `json:"credit,omitempty"`
	ExploitMaturity    string       `json:"exploitMaturity,omitempty"`
	SemverVulnerable   []string     `json:"semver,omitempty"`
	PublicationTime    string       `json:"publicationTime,omitempty"`
	DisclosureTime     string       `json:"disclosureTime,omitempty"`
	CVSSv3             string       `json:"CVSSv3,omitempty"`
	CvssScore          float64      `json:"cvssScore,omitempty"`
	Language           string       `json:"language,omitempty"`
	Patches            []Patch      `json:"patches,omitempty"`
	NearestFixedIn     []string     `json:"nearestFixedInVersion,omitempty"`
	IsMaliciousPackage bool         `json:"isMaliciousPackage,omitempty"`
}

// Identifiers contains CVE, CWE, and other identifiers.
type Identifiers struct {
	CVE  []string `json:"CVE,omitempty"`
	CWE  []string `json:"CWE,omitempty"`
	GHSA []string `json:"GHSA,omitempty"`
}

// Patch represents a Snyk patch.
type Patch struct {
	ID               string   `json:"id"`
	URLs             []string `json:"urls"`
	Version          string   `json:"version"`
	ModificationTime string   `json:"modificationTime"`
	Comments         []string `json:"comments,omitempty"`
}

// IgnoreReason represents why an issue was ignored.
type IgnoreReason struct {
	Reason             string     `json:"reason"`
	Expires            *time.Time `json:"expires,omitempty"`
	ReasonType         string     `json:"reasonType,omitempty"`
	DisregardIfFixable bool       `json:"disregardIfFixable,omitempty"`
	Source             string     `json:"source,omitempty"`
	Created            *time.Time `json:"created,omitempty"`
}

// FixInfo contains fix information.
type FixInfo struct {
	IsUpgradable          bool     `json:"isUpgradable"`
	IsPinnable            bool     `json:"isPinnable"`
	IsPatchable           bool     `json:"isPatchable"`
	IsFixable             bool     `json:"isFixable"`
	IsPartiallyFixable    bool     `json:"isPartiallyFixable"`
	NearestFixedInVersion string   `json:"nearestFixedInVersion,omitempty"`
	FixedIn               []string `json:"fixedIn,omitempty"`
}

// Priority contains priority score information.
type Priority struct {
	Score   int      `json:"score"`
	Factors []Factor `json:"factors,omitempty"`
}

// Factor represents a priority factor.
type Factor struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// IntroducedThrough represents how a vulnerability was introduced.
type IntroducedThrough struct {
	Kind    string      `json:"kind"`
	Version string      `json:"version,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetAggregatedIssues retrieves aggregated issues for a project.
// POST /org/{orgId}/project/{projectId}/aggregated-issues
func (c *Client) GetAggregatedIssues(ctx context.Context, orgID, projectID string, req *AggregatedIssuesRequest) (*AggregatedIssuesResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/aggregated-issues", orgID, projectID)

	if req == nil {
		req = &AggregatedIssuesRequest{}
	}

	var result AggregatedIssuesResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Issue Paths
// ============================================================================

// IssuePathsRequest represents the request for issue paths.
type IssuePathsRequest struct {
	// Pagination parameters
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
}

// IssuePathsResponse represents the response for issue paths.
type IssuePathsResponse struct {
	SnapshotID string       `json:"snapshotId"`
	Paths      [][]PathNode `json:"paths"`
	Total      int          `json:"total"`
}

// PathNode represents a node in a vulnerability path.
type PathNode struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// GetIssuePaths retrieves the dependency paths for a specific issue.
// POST /org/{orgId}/project/{projectId}/issue/{issueId}/paths
func (c *Client) GetIssuePaths(ctx context.Context, orgID, projectID, issueID string, req *IssuePathsRequest) (*IssuePathsResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/issue/%s/paths", orgID, projectID, issueID)

	if req == nil {
		req = &IssuePathsRequest{}
	}

	var result IssuePathsResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Ignores
// ============================================================================

// Ignore represents an ignore rule.
type Ignore struct {
	ID                 string     `json:"*,omitempty"` // Key is the path
	Reason             string     `json:"reason,omitempty"`
	ReasonType         string     `json:"reasonType,omitempty"`
	Expires            *time.Time `json:"expires,omitempty"`
	Created            *time.Time `json:"created,omitempty"`
	IgnoredBy          *User      `json:"ignoredBy,omitempty"`
	DisregardIfFixable bool       `json:"disregardIfFixable,omitempty"`
	Source             string     `json:"source,omitempty"`
}

// IgnoresResponse represents the response for listing ignores.
// The response is a map of issue ID to list of ignores.
type IgnoresResponse map[string][]Ignore

// ListIgnores retrieves all ignore rules for a project.
// GET /org/{orgId}/project/{projectId}/ignores
func (c *Client) ListIgnores(ctx context.Context, orgID, projectID string) (IgnoresResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/ignores", orgID, projectID)

	var result IgnoresResponse
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddIgnoreRequest represents the request to add an ignore.
type AddIgnoreRequest struct {
	Reason             string     `json:"reason,omitempty"`
	ReasonType         string     `json:"reasonType,omitempty"` // "not-vulnerable", "wont-fix", "temporary-ignore"
	Expires            *time.Time `json:"expires,omitempty"`
	DisregardIfFixable bool       `json:"disregardIfFixable,omitempty"`
	IgnorePath         string     `json:"ignorePath,omitempty"` // "*" for all paths
}

// AddIgnore adds an ignore rule for an issue.
// POST /org/{orgId}/project/{projectId}/ignore/{issueId}
func (c *Client) AddIgnore(ctx context.Context, orgID, projectID, issueID string, req *AddIgnoreRequest) error {
	path := fmt.Sprintf("/org/%s/project/%s/ignore/%s", orgID, projectID, issueID)
	return c.base.Post(ctx, path, req, nil)
}

// UpdateIgnore updates an ignore rule for an issue.
// PUT /org/{orgId}/project/{projectId}/ignore/{issueId}
func (c *Client) UpdateIgnore(ctx context.Context, orgID, projectID, issueID string, req *AddIgnoreRequest) error {
	path := fmt.Sprintf("/org/%s/project/%s/ignore/%s", orgID, projectID, issueID)
	return c.base.Put(ctx, path, req, nil)
}

// DeleteIgnore removes an ignore rule for an issue.
// DELETE /org/{orgId}/project/{projectId}/ignore/{issueId}
func (c *Client) DeleteIgnore(ctx context.Context, orgID, projectID, issueID string) error {
	path := fmt.Sprintf("/org/%s/project/%s/ignore/%s", orgID, projectID, issueID)
	return c.base.Delete(ctx, path)
}

// ============================================================================
// Project Lifecycle
// ============================================================================

// Deactivate deactivates a project.
// POST /org/{orgId}/project/{projectId}/deactivate
func (c *Client) Deactivate(ctx context.Context, orgID, projectID string) error {
	path := fmt.Sprintf("/org/%s/project/%s/deactivate", orgID, projectID)
	return c.base.Post(ctx, path, nil, nil)
}

// Activate activates a project.
// POST /org/{orgId}/project/{projectId}/activate
func (c *Client) Activate(ctx context.Context, orgID, projectID string) error {
	path := fmt.Sprintf("/org/%s/project/%s/activate", orgID, projectID)
	return c.base.Post(ctx, path, nil, nil)
}

// ============================================================================
// Dependency Graph
// ============================================================================

// DepGraph represents a project's dependency graph.
type DepGraph struct {
	SchemaVersion string     `json:"schemaVersion"`
	PkgManager    PkgManager `json:"pkgManager"`
	Pkgs          []Package  `json:"pkgs"`
	Graph         Graph      `json:"graph"`
}

// PkgManager represents the package manager.
type PkgManager struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// Package represents a package in the dep graph.
type Package struct {
	ID   string      `json:"id"`
	Info PackageInfo `json:"info"`
}

// PackageInfo contains package information.
type PackageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Graph represents the dependency graph structure.
type Graph struct {
	RootNodeID string      `json:"rootNodeId"`
	Nodes      []GraphNode `json:"nodes"`
}

// GraphNode represents a node in the dependency graph.
type GraphNode struct {
	NodeID string `json:"nodeId"`
	PkgID  string `json:"pkgId"`
	Deps   []Dep  `json:"deps,omitempty"`
}

// Dep represents a dependency.
type Dep struct {
	NodeID string `json:"nodeId"`
}

// DepGraphResponse represents the dep graph response.
type DepGraphResponse struct {
	DepGraph DepGraph `json:"depGraph"`
}

// GetDepGraph retrieves the dependency graph for a project.
// GET /org/{orgId}/project/{projectId}/dep-graph
func (c *Client) GetDepGraph(ctx context.Context, orgID, projectID string) (*DepGraphResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/dep-graph", orgID, projectID)

	var result DepGraphResponse
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// History / Snapshots
// ============================================================================

// HistoryRequest represents the request for project history.
type HistoryRequest struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
}

// HistoryResponse represents the project history response.
type HistoryResponse struct {
	Snapshots []Snapshot `json:"snapshots"`
	Total     int        `json:"total"`
}

// Snapshot represents a project snapshot.
type Snapshot struct {
	ID                    string                `json:"id"`
	Created               time.Time             `json:"created"`
	TotalDependencies     int                   `json:"totalDependencies"`
	IssueCountsBySeverity IssueCountsBySeverity `json:"issueCountsBySeverity"`
}

// GetHistory retrieves project test history.
// POST /org/{orgId}/project/{projectId}/history
func (c *Client) GetHistory(ctx context.Context, orgID, projectID string, req *HistoryRequest) (*HistoryResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/history", orgID, projectID)

	if req == nil {
		req = &HistoryRequest{}
	}

	var result HistoryResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetHistoricalAggregatedIssues retrieves aggregated issues for a historical snapshot.
// POST /org/{orgId}/project/{projectId}/history/{snapshotId}/aggregated-issues
func (c *Client) GetHistoricalAggregatedIssues(ctx context.Context, orgID, projectID, snapshotID string, req *AggregatedIssuesRequest) (*AggregatedIssuesResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/history/%s/aggregated-issues", orgID, projectID, snapshotID)

	if req == nil {
		req = &AggregatedIssuesRequest{}
	}

	var result AggregatedIssuesResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetHistoricalIssuePaths retrieves issue paths for a historical snapshot.
// POST /org/{orgId}/project/{projectId}/history/{snapshotId}/issue/{issueId}/paths
func (c *Client) GetHistoricalIssuePaths(ctx context.Context, orgID, projectID, snapshotID, issueID string, req *IssuePathsRequest) (*IssuePathsResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/history/%s/issue/%s/paths", orgID, projectID, snapshotID, issueID)

	if req == nil {
		req = &IssuePathsRequest{}
	}

	var result IssuePathsResponse
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Jira Integration
// ============================================================================

// JiraIssue represents a linked Jira issue.
type JiraIssue struct {
	JiraIssue JiraIssueDetails `json:"jiraIssue"`
}

// JiraIssueDetails contains Jira issue details.
type JiraIssueDetails struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// JiraIssuesResponse represents the Jira issues response.
type JiraIssuesResponse map[string][]JiraIssue

// ListJiraIssues lists all Jira issues linked to a project.
// GET /org/{orgId}/project/{projectId}/jira-issues
func (c *Client) ListJiraIssues(ctx context.Context, orgID, projectID string) (JiraIssuesResponse, error) {
	path := fmt.Sprintf("/org/%s/project/%s/jira-issues", orgID, projectID)

	var result JiraIssuesResponse
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateJiraIssueRequest represents the request to create a Jira issue.
type CreateJiraIssueRequest struct {
	Fields map[string]interface{} `json:"fields"`
}

// CreateJiraIssue creates a Jira issue for a Snyk issue.
// POST /org/{orgId}/project/{projectId}/issue/{issueId}/jira-issue
func (c *Client) CreateJiraIssue(ctx context.Context, orgID, projectID, issueID string, req *CreateJiraIssueRequest) (*JiraIssue, error) {
	path := fmt.Sprintf("/org/%s/project/%s/issue/%s/jira-issue", orgID, projectID, issueID)

	var result JiraIssue
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Project Settings
// ============================================================================

// ProjectSettings represents project settings.
type ProjectSettings struct {
	AutoDepUpgradeEnabled              *bool                  `json:"autoDepUpgradeEnabled,omitempty"`
	AutoDepUpgradeIgnoredDependencies  []string               `json:"autoDepUpgradeIgnoredDependencies,omitempty"`
	AutoDepUpgradeMinAge               *int                   `json:"autoDepUpgradeMinAge,omitempty"`
	AutoDepUpgradeLimit                *int                   `json:"autoDepUpgradeLimit,omitempty"`
	PullRequestFailOnAnyVulns          *bool                  `json:"pullRequestFailOnAnyVulns,omitempty"`
	PullRequestFailOnlyForHighSeverity *bool                  `json:"pullRequestFailOnlyForHighSeverity,omitempty"`
	PullRequestTestEnabled             *bool                  `json:"pullRequestTestEnabled,omitempty"`
	PullRequestAssignment              *PullRequestAssignment `json:"pullRequestAssignment,omitempty"`
}

// PullRequestAssignment represents PR assignment settings.
type PullRequestAssignment struct {
	Enabled bool     `json:"enabled"`
	Type    string   `json:"type"` // "auto" or "manual"
	Users   []string `json:"users,omitempty"`
}

// GetSettings retrieves project settings.
// GET /org/{orgId}/project/{projectId}/settings
func (c *Client) GetSettings(ctx context.Context, orgID, projectID string) (*ProjectSettings, error) {
	path := fmt.Sprintf("/org/%s/project/%s/settings", orgID, projectID)

	var result ProjectSettings
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates project settings.
// PUT /org/{orgId}/project/{projectId}/settings
func (c *Client) UpdateSettings(ctx context.Context, orgID, projectID string, settings *ProjectSettings) (*ProjectSettings, error) {
	path := fmt.Sprintf("/org/%s/project/%s/settings", orgID, projectID)

	var result ProjectSettings
	if err := c.base.Put(ctx, path, settings, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Tags
// ============================================================================

// AddTagRequest represents the request to add tags.
type AddTagRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AddTag adds a tag to a project.
// POST /org/{orgId}/project/{projectId}/tags
func (c *Client) AddTag(ctx context.Context, orgID, projectID string, req *AddTagRequest) error {
	path := fmt.Sprintf("/org/%s/project/%s/tags", orgID, projectID)
	return c.base.Post(ctx, path, req, nil)
}

// RemoveTagRequest represents the request to remove tags.
type RemoveTagRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RemoveTag removes a tag from a project.
// POST /org/{orgId}/project/{projectId}/tags/remove
func (c *Client) RemoveTag(ctx context.Context, orgID, projectID string, req *RemoveTagRequest) error {
	path := fmt.Sprintf("/org/%s/project/%s/tags/remove", orgID, projectID)
	return c.base.Post(ctx, path, req, nil)
}

// ============================================================================
// Move Project
// ============================================================================

// MoveProjectRequest represents the request to move a project.
type MoveProjectRequest struct {
	TargetOrgID string `json:"targetOrgId"`
}

// Move moves a project to another organization.
// PUT /org/{orgId}/project/{projectId}/move
func (c *Client) Move(ctx context.Context, orgID, projectID string, req *MoveProjectRequest) error {
	path := fmt.Sprintf("/org/%s/project/%s/move", orgID, projectID)
	return c.base.Put(ctx, path, req, nil)
}

// ============================================================================
// Attributes
// ============================================================================

// Attributes represents project attributes.
type Attributes struct {
	Criticality []string `json:"criticality,omitempty"`
	Lifecycle   []string `json:"lifecycle,omitempty"`
	Environment []string `json:"environment,omitempty"`
}

// UpdateAttributes updates project attributes.
// POST /org/{orgId}/project/{projectId}/attributes
func (c *Client) UpdateAttributes(ctx context.Context, orgID, projectID string, attrs *Attributes) error {
	path := fmt.Sprintf("/org/%s/project/%s/attributes", orgID, projectID)
	return c.base.Post(ctx, path, attrs, nil)
}

// Get retrieves a single project.
// GET /org/{orgId}/project/{projectId}
func (c *Client) Get(ctx context.Context, orgID, projectID string) (*Project, error) {
	path := fmt.Sprintf("/org/%s/project/%s", orgID, projectID)

	var result Project
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a project.
// DELETE /org/{orgId}/project/{projectId}
func (c *Client) Delete(ctx context.Context, orgID, projectID string) error {
	path := fmt.Sprintf("/org/%s/project/%s", orgID, projectID)
	return c.base.Delete(ctx, path)
}
