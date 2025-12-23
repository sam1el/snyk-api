// Package testing provides a client for Snyk v1 Testing API.
// This includes endpoints for testing packages, dependencies, and dep-graphs.
package testing

import (
	"context"
	"fmt"

	v1 "github.com/sam1el/snyk-api/pkg/apiclients/v1"
	v1projects "github.com/sam1el/snyk-api/pkg/apiclients/v1/projects"
	"github.com/sam1el/snyk-api/pkg/client"
)

// Client provides access to v1 Testing API endpoints.
type Client struct {
	base *v1.BaseClient
}

// New creates a new v1 Testing client.
func New(baseClient *client.Client) *Client {
	return &Client{
		base: v1.NewBaseClient(baseClient),
	}
}

// ============================================================================
// Types
// ============================================================================

// TestResult represents the result of a test.
type TestResult struct {
	OK                    bool                             `json:"ok"`
	IssuesCount           int                              `json:"issuesCount,omitempty"`
	DependencyCount       int                              `json:"dependencyCount"`
	Issues                Issues                           `json:"issues,omitempty"`
	PackageManager        string                           `json:"packageManager"`
	Org                   TestOrg                          `json:"org,omitempty"`
	IssueCountsBySeverity v1projects.IssueCountsBySeverity `json:"issueSeverity,omitempty"`
	LicenseIssues         []LicenseIssue                   `json:"licensesPolicy,omitempty"`
	FilePath              string                           `json:"path,omitempty"`
	RemediationAdvice     *RemediationAdvice               `json:"remediation,omitempty"`
}

// Issues contains vulnerabilities and license issues.
type Issues struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities,omitempty"`
	Licenses        []LicenseIssue  `json:"licenses,omitempty"`
}

// Vulnerability represents a vulnerability finding.
type Vulnerability struct {
	ID                 string                  `json:"id"`
	Title              string                  `json:"title"`
	Severity           string                  `json:"severity"`
	URL                string                  `json:"url"`
	Description        string                  `json:"description,omitempty"`
	From               []string                `json:"from"`
	Package            string                  `json:"package"`
	Version            string                  `json:"version"`
	Name               string                  `json:"name,omitempty"`
	ModuleName         string                  `json:"moduleName,omitempty"`
	CVSSv3             string                  `json:"CVSSv3,omitempty"`
	CvssScore          float64                 `json:"cvssScore,omitempty"`
	Identifiers        *v1projects.Identifiers `json:"identifiers,omitempty"`
	ExploitMaturity    string                  `json:"exploitMaturity,omitempty"`
	IsUpgradable       bool                    `json:"isUpgradable"`
	IsPatchable        bool                    `json:"isPatchable"`
	IsPinnable         bool                    `json:"isPinnable"`
	UpgradePath        []interface{}           `json:"upgradePath,omitempty"`
	Patches            []v1projects.Patch      `json:"patches,omitempty"`
	NearestFixedIn     string                  `json:"nearestFixedInVersion,omitempty"`
	PublicationTime    string                  `json:"publicationTime,omitempty"`
	DisclosureTime     string                  `json:"disclosureTime,omitempty"`
	Credit             []string                `json:"credit,omitempty"`
	IsMaliciousPackage bool                    `json:"isMaliciousPackage,omitempty"`
}

// LicenseIssue represents a license issue.
type LicenseIssue struct {
	ID           string   `json:"id"`
	License      string   `json:"license"`
	Severity     string   `json:"severity"`
	Title        string   `json:"title,omitempty"`
	Instructions string   `json:"instructions,omitempty"`
	LegalContent *string  `json:"legalContent,omitempty"`
	Package      string   `json:"package,omitempty"`
	Version      string   `json:"version,omitempty"`
	From         []string `json:"from,omitempty"`
}

// TestOrg represents the organization context of a test.
type TestOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RemediationAdvice contains remediation advice for issues.
type RemediationAdvice struct {
	Unresolved []UnresolvedIssue        `json:"unresolved,omitempty"`
	Upgrade    map[string]UpgradeAdvice `json:"upgrade,omitempty"`
	Patch      map[string]PatchAdvice   `json:"patch,omitempty"`
	Pin        map[string]PinAdvice     `json:"pin,omitempty"`
}

// UnresolvedIssue represents an issue that cannot be automatically resolved.
type UnresolvedIssue struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Severity string   `json:"severity"`
	Path     []string `json:"path"`
}

// UpgradeAdvice represents upgrade advice for a package.
type UpgradeAdvice struct {
	UpgradeTo string   `json:"upgradeTo"`
	Upgrades  []string `json:"upgrades"`
	Vulns     []string `json:"vulns"`
}

// PatchAdvice represents patch advice for a vulnerability.
type PatchAdvice struct {
	Paths [][]string `json:"paths"`
}

// PinAdvice represents pin advice for a package.
type PinAdvice struct {
	IsTransitive bool   `json:"isTransitive"`
	UpgradeTo    string `json:"upgradeTo,omitempty"`
}

// ============================================================================
// Dependency Graph Testing
// ============================================================================

// TestDepGraphRequest represents the request to test a dep-graph.
type TestDepGraphRequest struct {
	DepGraph v1projects.DepGraph `json:"depGraph"`
}

// TestDepGraph tests a dependency graph for vulnerabilities.
// POST /test/dep-graph
func (c *Client) TestDepGraph(ctx context.Context, orgID string, depGraph *v1projects.DepGraph) (*TestResult, error) {
	path := "/test/dep-graph"
	if orgID != "" {
		path = fmt.Sprintf("/test/dep-graph?org=%s", orgID)
	}

	req := &TestDepGraphRequest{DepGraph: *depGraph}
	var result TestResult
	if err := c.base.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Maven Testing
// ============================================================================

// TestMavenByCoords tests a Maven package by coordinates.
// GET /test/maven/{groupId}/{artifactId}/{version}
func (c *Client) TestMavenByCoords(ctx context.Context, orgID, groupID, artifactID, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/maven/%s/%s/%s", groupID, artifactID, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MavenFile represents a Maven manifest file for testing.
type MavenFile struct {
	Encoding   string            `json:"encoding,omitempty"`
	Files      map[string]string `json:"files"`
	Additional map[string]string `json:"additional,omitempty"`
}

// TestMavenFile tests a Maven pom.xml file.
// POST /test/maven
func (c *Client) TestMavenFile(ctx context.Context, orgID string, file *MavenFile) (*TestResult, error) {
	path := "/test/maven"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// NPM Testing
// ============================================================================

// TestNpmByName tests an npm package by name and version.
// GET /test/npm/{packageName}/{version}
func (c *Client) TestNpmByName(ctx context.Context, orgID, packageName, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/npm/%s/%s", packageName, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NpmFile represents npm manifest files for testing.
type NpmFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // package.json, package-lock.json
}

// TestNpmFile tests npm package files.
// POST /test/npm
func (c *Client) TestNpmFile(ctx context.Context, orgID string, file *NpmFile) (*TestResult, error) {
	path := "/test/npm"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Yarn Testing
// ============================================================================

// YarnFile represents yarn manifest files for testing.
type YarnFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // package.json, yarn.lock
}

// TestYarnFile tests yarn package files.
// POST /test/yarn
func (c *Client) TestYarnFile(ctx context.Context, orgID string, file *YarnFile) (*TestResult, error) {
	path := "/test/yarn"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Pip Testing
// ============================================================================

// TestPipByName tests a pip package by name and version.
// GET /test/pip/{packageName}/{version}
func (c *Client) TestPipByName(ctx context.Context, orgID, packageName, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/pip/%s/%s", packageName, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PipFile represents pip manifest files for testing.
type PipFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // requirements.txt
}

// TestPipFile tests pip requirements files.
// POST /test/pip
func (c *Client) TestPipFile(ctx context.Context, orgID string, file *PipFile) (*TestResult, error) {
	path := "/test/pip"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// RubyGems Testing
// ============================================================================

// TestRubyGemsByName tests a Ruby gem by name and version.
// GET /test/rubygems/{gemName}/{version}
func (c *Client) TestRubyGemsByName(ctx context.Context, orgID, gemName, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/rubygems/%s/%s", gemName, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RubyGemsFile represents Ruby manifest files for testing.
type RubyGemsFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // Gemfile, Gemfile.lock
}

// TestRubyGemsFile tests Ruby Gemfile.
// POST /test/rubygems
func (c *Client) TestRubyGemsFile(ctx context.Context, orgID string, file *RubyGemsFile) (*TestResult, error) {
	path := "/test/rubygems"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Gradle Testing
// ============================================================================

// TestGradleByCoords tests a Gradle package by coordinates.
// GET /test/gradle/{group}/{name}/{version}
func (c *Client) TestGradleByCoords(ctx context.Context, orgID, group, name, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/gradle/%s/%s/%s", group, name, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GradleFile represents Gradle manifest files for testing.
type GradleFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // build.gradle
}

// TestGradleFile tests Gradle build files.
// POST /test/gradle
func (c *Client) TestGradleFile(ctx context.Context, orgID string, file *GradleFile) (*TestResult, error) {
	path := "/test/gradle"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// SBT Testing
// ============================================================================

// TestSbtByCoords tests an SBT package by coordinates.
// GET /test/sbt/{groupId}/{artifactId}/{version}
func (c *Client) TestSbtByCoords(ctx context.Context, orgID, groupID, artifactID, version string) (*TestResult, error) {
	path := fmt.Sprintf("/test/sbt/%s/%s/%s", groupID, artifactID, version)
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SbtFile represents SBT manifest files for testing.
type SbtFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // build.sbt
}

// TestSbtFile tests SBT build files.
// POST /test/sbt
func (c *Client) TestSbtFile(ctx context.Context, orgID string, file *SbtFile) (*TestResult, error) {
	path := "/test/sbt"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Composer Testing
// ============================================================================

// ComposerFile represents Composer manifest files for testing.
type ComposerFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // composer.json, composer.lock
}

// TestComposerFile tests Composer files.
// POST /test/composer
func (c *Client) TestComposerFile(ctx context.Context, orgID string, file *ComposerFile) (*TestResult, error) {
	path := "/test/composer"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================================
// Go Testing
// ============================================================================

// GoDepFile represents Go dep manifest files for testing.
type GoDepFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // Gopkg.toml, Gopkg.lock
}

// TestGoDepFile tests Go dep files.
// POST /test/golangdep
func (c *Client) TestGoDepFile(ctx context.Context, orgID string, file *GoDepFile) (*TestResult, error) {
	path := "/test/golangdep"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GoVendorFile represents Go vendor manifest files for testing.
type GoVendorFile struct {
	Encoding string            `json:"encoding,omitempty"`
	Files    map[string]string `json:"files"` // vendor.json
}

// TestGoVendorFile tests Go vendor files.
// POST /test/govendor
func (c *Client) TestGoVendorFile(ctx context.Context, orgID string, file *GoVendorFile) (*TestResult, error) {
	path := "/test/govendor"
	if orgID != "" {
		path = fmt.Sprintf("%s?org=%s", path, orgID)
	}

	var result TestResult
	if err := c.base.Post(ctx, path, file, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
