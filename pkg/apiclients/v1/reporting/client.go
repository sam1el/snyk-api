// Package reporting provides a client for Snyk v1 Reporting API.
package reporting

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Reporting API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Reporting client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// ReportFilters represents filters for reports.
type ReportFilters struct {
	Orgs          []string `json:"orgs,omitempty"`
	Severity      []string `json:"severity,omitempty"`
	Types         []string `json:"types,omitempty"`
	Languages     []string `json:"languages,omitempty"`
	Projects      []string `json:"projects,omitempty"`
	Issues        []string `json:"issues,omitempty"`
	Identifier    string   `json:"identifier,omitempty"`
	Ignored       *bool    `json:"ignored,omitempty"`
	Patched       *bool    `json:"patched,omitempty"`
	Fixable       *bool    `json:"fixable,omitempty"`
	IsFixed       *bool    `json:"isFixed,omitempty"`
	IsUpgradeable *bool    `json:"isUpgradeable,omitempty"`
	IsPatchable   *bool    `json:"isPatchable,omitempty"`
	IsPinnable    *bool    `json:"isPinnable,omitempty"`
}

// ============================================================================
// Issues Reporting
// ============================================================================

// IssuesRequest represents a request for issues report.
type IssuesRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
	Page    int           `json:"page,omitempty"`
	PerPage int           `json:"perPage,omitempty"`
	SortBy  string        `json:"sortBy,omitempty"`
	Order   string        `json:"order,omitempty"`
}

// IssuesResponse represents the response for issues report.
type IssuesResponse struct {
	Results []ReportIssue `json:"results"`
	Total   int           `json:"total"`
}

// ReportIssue represents an issue in a report.
type ReportIssue struct {
	Issue          IssueDetails `json:"issue"`
	IsFixed        bool         `json:"isFixed"`
	IntroducedDate time.Time    `json:"introducedDate"`
	FixedDate      *time.Time   `json:"fixedDate,omitempty"`
	Project        ProjectInfo  `json:"project"`
}

// IssueDetails contains issue details.
type IssueDetails struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	Severity         string       `json:"severity"`
	OriginalSeverity string       `json:"originalSeverity,omitempty"`
	Type             string       `json:"type"`
	URL              string       `json:"url"`
	Package          string       `json:"package,omitempty"`
	Version          string       `json:"version,omitempty"`
	Language         string       `json:"language,omitempty"`
	PackageManager   string       `json:"packageManager,omitempty"`
	Ignored          bool         `json:"ignored"`
	Patched          bool         `json:"patched"`
	ExploitMaturity  string       `json:"exploitMaturity,omitempty"`
	CVSSv3           string       `json:"CVSSv3,omitempty"`
	CvssScore        float64      `json:"cvssScore,omitempty"`
	Identifiers      *Identifiers `json:"identifiers,omitempty"`
	Fixable          bool         `json:"fixable"`
	IsUpgradable     bool         `json:"isUpgradable"`
	IsPatchable      bool         `json:"isPatchable"`
	IsPinnable       bool         `json:"isPinnable"`
}

// Identifiers contains CVE and other identifiers.
type Identifiers struct {
	CVE []string `json:"CVE,omitempty"`
	CWE []string `json:"CWE,omitempty"`
}

// ProjectInfo contains project information.
type ProjectInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
	URL    string `json:"url,omitempty"`
}

// GetLatestIssues retrieves the latest issues report.
// POST /reporting/issues/latest
func (c *Client) GetLatestIssues(ctx context.Context, req *IssuesRequest) (*IssuesResponse, error) {
	if req == nil {
		req = &IssuesRequest{}
	}
	var result IssuesResponse
	if err := c.base.Post(ctx, "/reporting/issues/latest", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// IssuesHistoryRequest represents a request for issues history.
type IssuesHistoryRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
	From    time.Time     `json:"from"`
	To      time.Time     `json:"to"`
	Page    int           `json:"page,omitempty"`
	PerPage int           `json:"perPage,omitempty"`
	SortBy  string        `json:"sortBy,omitempty"`
	Order   string        `json:"order,omitempty"`
}

// GetIssuesHistory retrieves issues over a time period.
// POST /reporting/issues
func (c *Client) GetIssuesHistory(ctx context.Context, req *IssuesHistoryRequest) (*IssuesResponse, error) {
	if req == nil {
		req = &IssuesHistoryRequest{}
	}
	var result IssuesResponse
	if err := c.base.Post(ctx, "/reporting/issues", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Counts Reporting
// ============================================================================

// CountsRequest represents a request for counts.
type CountsRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
}

// CountsHistoryRequest represents a request for historical counts.
type CountsHistoryRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
	From    time.Time     `json:"from"`
	To      time.Time     `json:"to"`
}

// IssueCountsResponse represents issue counts response.
type IssueCountsResponse struct {
	Results []IssueCount `json:"results"`
}

// IssueCount represents issue counts.
type IssueCount struct {
	Day      string         `json:"day,omitempty"`
	Severity SeverityCounts `json:"severity"`
	Count    int            `json:"count"`
}

// SeverityCounts represents counts by severity.
type SeverityCounts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

// GetLatestIssueCounts retrieves latest issue counts.
// POST /reporting/counts/issues/latest
func (c *Client) GetLatestIssueCounts(ctx context.Context, req *CountsRequest) (*IssueCountsResponse, error) {
	if req == nil {
		req = &CountsRequest{}
	}
	var result IssueCountsResponse
	if err := c.base.Post(ctx, "/reporting/counts/issues/latest", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetIssueCountsHistory retrieves issue counts over time.
// POST /reporting/counts/issues
func (c *Client) GetIssueCountsHistory(ctx context.Context, req *CountsHistoryRequest) (*IssueCountsResponse, error) {
	if req == nil {
		req = &CountsHistoryRequest{}
	}
	var result IssueCountsResponse
	if err := c.base.Post(ctx, "/reporting/counts/issues", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ProjectCountsResponse represents project counts response.
type ProjectCountsResponse struct {
	Results []ProjectCount `json:"results"`
}

// ProjectCount represents project counts.
type ProjectCount struct {
	Day   string `json:"day,omitempty"`
	Count int    `json:"count"`
}

// GetLatestProjectCounts retrieves latest project counts.
// POST /reporting/counts/projects/latest
func (c *Client) GetLatestProjectCounts(ctx context.Context, req *CountsRequest) (*ProjectCountsResponse, error) {
	if req == nil {
		req = &CountsRequest{}
	}
	var result ProjectCountsResponse
	if err := c.base.Post(ctx, "/reporting/counts/projects/latest", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProjectCountsHistory retrieves project counts over time.
// POST /reporting/counts/projects
func (c *Client) GetProjectCountsHistory(ctx context.Context, req *CountsHistoryRequest) (*ProjectCountsResponse, error) {
	if req == nil {
		req = &CountsHistoryRequest{}
	}
	var result ProjectCountsResponse
	if err := c.base.Post(ctx, "/reporting/counts/projects", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TestCountsRequest represents a request for test counts.
type TestCountsRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
	From    time.Time     `json:"from"`
	To      time.Time     `json:"to"`
	GroupBy string        `json:"groupBy,omitempty"` // "day", "week", "month"
}

// TestCountsResponse represents test counts response.
type TestCountsResponse struct {
	Results []TestCount `json:"results"`
}

// TestCount represents test counts.
type TestCount struct {
	Day   string `json:"day,omitempty"`
	Count int    `json:"count"`
}

// GetTestCounts retrieves test counts over time.
// POST /reporting/counts/tests
func (c *Client) GetTestCounts(ctx context.Context, req *TestCountsRequest) (*TestCountsResponse, error) {
	if req == nil {
		now := time.Now()
		req = &TestCountsRequest{
			From: now.AddDate(0, -1, 0),
			To:   now,
		}
	}
	var result TestCountsResponse
	if err := c.base.Post(ctx, "/reporting/counts/tests", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExportRequest represents a report export request.
type ExportRequest struct {
	Filters ReportFilters `json:"filters,omitempty"`
	From    time.Time     `json:"from,omitempty"`
	To      time.Time     `json:"to,omitempty"`
	Format  string        `json:"format,omitempty"` // "csv", "json"
	Columns []string      `json:"columns,omitempty"`
}

// ExportJob represents an export job.
type ExportJob struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Created string `json:"created"`
	URL     string `json:"url,omitempty"` // Download URL when complete
}

// StartExport starts an export job.
func (c *Client) StartExport(ctx context.Context, orgID string, req *ExportRequest) (*ExportJob, error) {
	path := fmt.Sprintf("/org/%s/reporting/issues/export", orgID)
	if req == nil {
		req = &ExportRequest{}
	}
	var result ExportJob
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExportJob retrieves an export job status.
func (c *Client) GetExportJob(ctx context.Context, orgID, jobID string) (*ExportJob, error) {
	path := fmt.Sprintf("/org/%s/reporting/issues/export/%s", orgID, jobID)
	var result ExportJob
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
