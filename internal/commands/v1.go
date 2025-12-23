package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/v1/groups"
	"github.com/sam1el/snyk-api/pkg/apiclients/v1/integrations"
	v1orgs "github.com/sam1el/snyk-api/pkg/apiclients/v1/orgs"
	v1projects "github.com/sam1el/snyk-api/pkg/apiclients/v1/projects"
	"github.com/sam1el/snyk-api/pkg/apiclients/v1/reporting"
	v1testing "github.com/sam1el/snyk-api/pkg/apiclients/v1/testing"
	"github.com/sam1el/snyk-api/pkg/apiclients/v1/users"
	"github.com/sam1el/snyk-api/pkg/apiclients/v1/webhooks"
	"github.com/sam1el/snyk-api/pkg/client"
	"github.com/spf13/cobra"
)

// newV1Cmd creates the v1 command group.
func newV1Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "v1",
		Aliases: []string{"api-v1"},
		Short:   "Access Snyk v1 API endpoints",
		Long: `Access Snyk v1 API endpoints directly.

The v1 API provides access to endpoints that are not yet available in the REST API,
including aggregated issues, ignores, dependency graphs, and more.`,
	}

	cmd.AddCommand(
		newV1ProjectsCmd(),
		newV1OrgsCmd(),
		newV1UsersCmd(),
		newV1GroupsCmd(),
		newV1TestingCmd(),
		newV1ReportingCmd(),
		newV1IntegrationsCmd(),
		newV1WebhooksCmd(),
	)

	return cmd
}

// ============================================================================
// v1 Projects Commands
// ============================================================================

func newV1ProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "v1 Projects API operations",
	}

	cmd.AddCommand(
		newV1ProjectsAggregatedIssuesCmd(),
		newV1ProjectsIssuePathsCmd(),
		newV1ProjectsIgnoresCmd(),
		newV1ProjectsDepGraphCmd(),
		newV1ProjectsHistoryCmd(),
		newV1ProjectsSettingsCmd(),
	)

	return cmd
}

func newV1ProjectsAggregatedIssuesCmd() *cobra.Command {
	var orgID, projectID string
	var includeDeps, includeDesc bool
	var severities []string

	cmd := &cobra.Command{
		Use:   "aggregated-issues",
		Short: "Get aggregated issues for a project",
		Long:  "Retrieves all issues for a project, aggregated by vulnerability.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)

			req := &v1projects.AggregatedIssuesRequest{
				IncludeDescription:       includeDesc,
				IncludeIntroducedThrough: includeDeps,
			}
			if len(severities) > 0 {
				req.Filters = &v1projects.AggregatedFilters{
					Severities: severities,
				}
			}

			resp, err := c.GetAggregatedIssues(ctx, orgID, projectID, req)
			if err != nil {
				return fmt.Errorf("failed to get aggregated issues: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	cmd.Flags().BoolVar(&includeDeps, "include-deps", false, "Include introduced through info")
	cmd.Flags().BoolVar(&includeDesc, "include-desc", false, "Include descriptions")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "Filter by severity (critical,high,medium,low)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

func newV1ProjectsIssuePathsCmd() *cobra.Command {
	var orgID, projectID, issueID string

	cmd := &cobra.Command{
		Use:   "issue-paths",
		Short: "Get dependency paths for an issue",
		Long:  "Retrieves the dependency paths that introduce a specific vulnerability.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			resp, err := c.GetIssuePaths(ctx, orgID, projectID, issueID, nil)
			if err != nil {
				return fmt.Errorf("failed to get issue paths: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	cmd.Flags().StringVar(&issueID, "issue-id", "", "Issue ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("issue-id")

	return cmd
}

func newV1ProjectsIgnoresCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ignores",
		Short: "Manage project ignores",
	}

	cmd.AddCommand(
		newV1ProjectsIgnoresListCmd(),
		newV1ProjectsIgnoresAddCmd(),
		newV1ProjectsIgnoresDeleteCmd(),
	)

	return cmd
}

func newV1ProjectsIgnoresListCmd() *cobra.Command {
	var orgID, projectID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all ignores for a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			resp, err := c.ListIgnores(ctx, orgID, projectID)
			if err != nil {
				return fmt.Errorf("failed to list ignores: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

func newV1ProjectsIgnoresAddCmd() *cobra.Command {
	var orgID, projectID, issueID, reason, reasonType string
	var disregardIfFixable bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an ignore for an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			req := &v1projects.AddIgnoreRequest{
				Reason:             reason,
				ReasonType:         reasonType,
				DisregardIfFixable: disregardIfFixable,
				IgnorePath:         "*", // Ignore all paths
			}

			if err := c.AddIgnore(ctx, orgID, projectID, issueID, req); err != nil {
				return fmt.Errorf("failed to add ignore: %w", err)
			}

			fmt.Printf("Successfully added ignore for issue %s\n", issueID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	cmd.Flags().StringVar(&issueID, "issue-id", "", "Issue ID to ignore (required)")
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for ignoring")
	cmd.Flags().StringVar(&reasonType, "reason-type", "temporary-ignore", "Reason type (not-vulnerable, wont-fix, temporary-ignore)")
	cmd.Flags().BoolVar(&disregardIfFixable, "disregard-if-fixable", false, "Remove ignore if fix becomes available")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("issue-id")

	return cmd
}

func newV1ProjectsIgnoresDeleteCmd() *cobra.Command {
	var orgID, projectID, issueID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an ignore for an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			if err := c.DeleteIgnore(ctx, orgID, projectID, issueID); err != nil {
				return fmt.Errorf("failed to delete ignore: %w", err)
			}

			fmt.Printf("Successfully deleted ignore for issue %s\n", issueID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	cmd.Flags().StringVar(&issueID, "issue-id", "", "Issue ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("issue-id")

	return cmd
}

func newV1ProjectsDepGraphCmd() *cobra.Command {
	var orgID, projectID string

	cmd := &cobra.Command{
		Use:   "dep-graph",
		Short: "Get dependency graph for a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			resp, err := c.GetDepGraph(ctx, orgID, projectID)
			if err != nil {
				return fmt.Errorf("failed to get dep-graph: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

func newV1ProjectsHistoryCmd() *cobra.Command {
	var orgID, projectID string

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Get test history for a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			resp, err := c.GetHistory(ctx, orgID, projectID, nil)
			if err != nil {
				return fmt.Errorf("failed to get history: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

func newV1ProjectsSettingsCmd() *cobra.Command {
	var orgID, projectID string

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Get project settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1projects.New(baseClient)
			resp, err := c.GetSettings(ctx, orgID, projectID)
			if err != nil {
				return fmt.Errorf("failed to get settings: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

// ============================================================================
// v1 Orgs Commands
// ============================================================================

func newV1OrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "v1 Organizations API operations",
	}

	cmd.AddCommand(
		newV1OrgsListCmd(),
		newV1OrgsMembersCmd(),
		newV1OrgsDependenciesCmd(),
		newV1OrgsLicensesCmd(),
		newV1OrgsEntitlementsCmd(),
	)

	return cmd
}

func newV1OrgsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all organizations (v1)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1orgs.New(baseClient)
			resp, err := c.List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list orgs: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

func newV1OrgsMembersCmd() *cobra.Command {
	var orgID string
	var includeGroupAdmins bool

	cmd := &cobra.Command{
		Use:   "members",
		Short: "List organization members",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1orgs.New(baseClient)
			resp, err := c.ListMembers(ctx, orgID, includeGroupAdmins)
			if err != nil {
				return fmt.Errorf("failed to list members: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().BoolVar(&includeGroupAdmins, "include-group-admins", false, "Include group admins")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newV1OrgsDependenciesCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "dependencies",
		Short: "List organization dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1orgs.New(baseClient)
			resp, err := c.ListDependencies(ctx, orgID, nil)
			if err != nil {
				return fmt.Errorf("failed to list dependencies: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newV1OrgsLicensesCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "licenses",
		Short: "List organization licenses",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1orgs.New(baseClient)
			resp, err := c.ListLicenses(ctx, orgID, nil)
			if err != nil {
				return fmt.Errorf("failed to list licenses: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newV1OrgsEntitlementsCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "entitlements",
		Short: "Get organization entitlements",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1orgs.New(baseClient)
			resp, err := c.GetEntitlements(ctx, orgID)
			if err != nil {
				return fmt.Errorf("failed to get entitlements: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

// ============================================================================
// v1 Users Commands
// ============================================================================

func newV1UsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "v1 Users API operations",
	}

	cmd.AddCommand(
		newV1UsersMeCmd(),
		newV1UsersGetCmd(),
	)

	return cmd
}

func newV1UsersMeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Get current user details",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := users.New(baseClient)
			resp, err := c.GetMe(ctx)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

func newV1UsersGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [user-id]",
		Short: "Get user by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := users.New(baseClient)
			resp, err := c.Get(ctx, args[0])
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

// ============================================================================
// v1 Groups Commands
// ============================================================================

func newV1GroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "v1 Groups API operations",
	}

	cmd.AddCommand(
		newV1GroupsMembersCmd(),
		newV1GroupsOrgsCmd(),
		newV1GroupsRolesCmd(),
		newV1GroupsTagsCmd(),
	)

	return cmd
}

func newV1GroupsMembersCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "members",
		Short: "List group members",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := groups.New(baseClient)
			resp, err := c.ListMembers(ctx, groupID)
			if err != nil {
				return fmt.Errorf("failed to list members: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newV1GroupsOrgsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "List organizations in a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := groups.New(baseClient)
			resp, err := c.ListOrgs(ctx, groupID, 100, 1)
			if err != nil {
				return fmt.Errorf("failed to list orgs: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newV1GroupsRolesCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List group roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := groups.New(baseClient)
			resp, err := c.ListRoles(ctx, groupID)
			if err != nil {
				return fmt.Errorf("failed to list roles: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newV1GroupsTagsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List group tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := groups.New(baseClient)
			resp, err := c.ListTags(ctx, groupID, 100, 1)
			if err != nil {
				return fmt.Errorf("failed to list tags: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

// ============================================================================
// v1 Testing Commands
// ============================================================================

func newV1TestingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "v1 Testing API operations",
		Long:  "Test packages and dependencies for vulnerabilities using the v1 API.",
	}

	cmd.AddCommand(
		newV1TestNpmCmd(),
		newV1TestMavenCmd(),
		newV1TestPipCmd(),
		newV1TestDepGraphCmd(),
	)

	return cmd
}

func newV1TestNpmCmd() *cobra.Command {
	var orgID, packageName, version string

	cmd := &cobra.Command{
		Use:   "npm",
		Short: "Test an npm package",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1testing.New(baseClient)
			resp, err := c.TestNpmByName(ctx, orgID, packageName, version)
			if err != nil {
				return fmt.Errorf("failed to test npm package: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&packageName, "package", "", "Package name (required)")
	cmd.Flags().StringVar(&version, "version", "", "Package version (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("package")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("version")

	return cmd
}

func newV1TestMavenCmd() *cobra.Command {
	var orgID, groupID, artifactID, version string

	cmd := &cobra.Command{
		Use:   "maven",
		Short: "Test a Maven package",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1testing.New(baseClient)
			resp, err := c.TestMavenByCoords(ctx, orgID, groupID, artifactID, version)
			if err != nil {
				return fmt.Errorf("failed to test Maven package: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&groupID, "group", "", "Group ID (required)")
	cmd.Flags().StringVar(&artifactID, "artifact", "", "Artifact ID (required)")
	cmd.Flags().StringVar(&version, "version", "", "Version (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("artifact")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("version")

	return cmd
}

func newV1TestPipCmd() *cobra.Command {
	var orgID, packageName, version string

	cmd := &cobra.Command{
		Use:   "pip",
		Short: "Test a pip package",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := v1testing.New(baseClient)
			resp, err := c.TestPipByName(ctx, orgID, packageName, version)
			if err != nil {
				return fmt.Errorf("failed to test pip package: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&packageName, "package", "", "Package name (required)")
	cmd.Flags().StringVar(&version, "version", "", "Package version (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("package")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("version")

	return cmd
}

func newV1TestDepGraphCmd() *cobra.Command {
	var orgID, filePath string

	cmd := &cobra.Command{
		Use:   "dep-graph",
		Short: "Test a dependency graph",
		Long:  "Test a dependency graph JSON file for vulnerabilities.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			// Read dep-graph from file
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read dep-graph file: %w", err)
			}

			var depGraph v1projects.DepGraph
			if err := json.Unmarshal(data, &depGraph); err != nil {
				return fmt.Errorf("failed to parse dep-graph: %w", err)
			}

			c := v1testing.New(baseClient)
			resp, err := c.TestDepGraph(ctx, orgID, &depGraph)
			if err != nil {
				return fmt.Errorf("failed to test dep-graph: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to dep-graph JSON file (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

// ============================================================================
// v1 Reporting Commands
// ============================================================================

func newV1ReportingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reporting",
		Short: "v1 Reporting API operations",
	}

	cmd.AddCommand(
		newV1ReportingIssuesCmd(),
		newV1ReportingCountsCmd(),
	)

	return cmd
}

func newV1ReportingIssuesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issues",
		Short: "Get latest issues report",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := reporting.New(baseClient)
			resp, err := c.GetLatestIssues(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to get issues report: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

func newV1ReportingCountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "counts",
		Short: "Get issue counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := reporting.New(baseClient)
			resp, err := c.GetLatestIssueCounts(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to get issue counts: %w", err)
			}

			// Table output for counts
			if getOutputFormat() == "table" {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "DAY\tCRITICAL\tHIGH\tMEDIUM\tLOW\tTOTAL")
				for _, r := range resp.Results {
					fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\n",
						r.Day,
						r.Severity.Critical,
						r.Severity.High,
						r.Severity.Medium,
						r.Severity.Low,
						r.Count,
					)
				}
				return w.Flush()
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

// ============================================================================
// v1 Integrations Commands
// ============================================================================

func newV1IntegrationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "integrations",
		Short: "v1 Integrations API operations",
	}

	cmd.AddCommand(
		newV1IntegrationsListCmd(),
		newV1IntegrationsImportCmd(),
	)

	return cmd
}

func newV1IntegrationsListCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List integrations for an org",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := integrations.New(baseClient)
			resp, err := c.List(ctx, orgID)
			if err != nil {
				return fmt.Errorf("failed to list integrations: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newV1IntegrationsImportCmd() *cobra.Command {
	var orgID, integrationID, owner, name, branch string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import a project from an integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := integrations.New(baseClient)
			req := &integrations.ImportRequest{
				Target: integrations.ImportTarget{
					Owner:  owner,
					Name:   name,
					Branch: branch,
				},
			}

			resp, err := c.Import(ctx, orgID, integrationID, req)
			if err != nil {
				return fmt.Errorf("failed to import: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&integrationID, "integration-id", "", "Integration ID (required)")
	cmd.Flags().StringVar(&owner, "owner", "", "Repository owner (required)")
	cmd.Flags().StringVar(&name, "name", "", "Repository name (required)")
	cmd.Flags().StringVar(&branch, "branch", "", "Branch name")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("integration-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("owner")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

// ============================================================================
// v1 Webhooks Commands
// ============================================================================

func newV1WebhooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhooks",
		Short: "v1 Webhooks API operations",
	}

	cmd.AddCommand(
		newV1WebhooksListCmd(),
		newV1WebhooksCreateCmd(),
		newV1WebhooksDeleteCmd(),
		newV1WebhooksPingCmd(),
	)

	return cmd
}

func newV1WebhooksListCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List webhooks for an org",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := webhooks.New(baseClient)
			resp, err := c.List(ctx, orgID)
			if err != nil {
				return fmt.Errorf("failed to list webhooks: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newV1WebhooksCreateCmd() *cobra.Command {
	var orgID, url, secret string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := webhooks.New(baseClient)
			resp, err := c.Create(ctx, orgID, &webhooks.CreateWebhookRequest{
				URL:    url,
				Secret: secret,
			})
			if err != nil {
				return fmt.Errorf("failed to create webhook: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&url, "url", "", "Webhook URL (required)")
	cmd.Flags().StringVar(&secret, "secret", "", "Webhook secret")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func newV1WebhooksDeleteCmd() *cobra.Command {
	var orgID, webhookID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := webhooks.New(baseClient)
			if err := c.Delete(ctx, orgID, webhookID); err != nil {
				return fmt.Errorf("failed to delete webhook: %w", err)
			}

			fmt.Printf("Successfully deleted webhook %s\n", webhookID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&webhookID, "webhook-id", "", "Webhook ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("webhook-id")

	return cmd
}

func newV1WebhooksPingCmd() *cobra.Command {
	var orgID, webhookID string

	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping a webhook (test)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := webhooks.New(baseClient)
			if err := c.Ping(ctx, orgID, webhookID); err != nil {
				return fmt.Errorf("failed to ping webhook: %w", err)
			}

			fmt.Printf("Successfully pinged webhook %s\n", webhookID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&webhookID, "webhook-id", "", "Webhook ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("webhook-id")

	return cmd
}
