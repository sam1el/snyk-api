package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/issues"
)

// newIssuesCmd creates the issues subcommand.
func newIssuesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issues",
		Aliases: []string{"issue"},
		Short:   "Manage Snyk issues.",
		Long: `Manage Snyk issues (vulnerabilities and license issues).

Issues represent security vulnerabilities and license problems found in your projects.`,
	}

	cmd.AddCommand(newIssuesListCmd())
	cmd.AddCommand(newIssuesGetCmd())

	return cmd
}

func newIssuesListCmd() *cobra.Command {
	var orgID string
	var limit int
	var startingAfter string
	var severity string
	var issueType string
	var status string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues in an organization",
		Long: `List all issues across projects within a specified Snyk organization.

Examples:
  # List first 10 issues in an organization
  snyk-api issues list --org-id=<org-id>

  # List critical vulnerabilities
  snyk-api issues list --org-id=<org-id> --severity critical --type vuln

  # List open issues
  snyk-api issues list --org-id=<org-id> --status open

  # List with table output
  snyk-api issues list --org-id=<org-id> --output table

  # Pagination
  snyk-api issues list --org-id=<org-id> --starting-after <cursor>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("organization ID is required. Use --org-id flag")
			}

			ctx := cmd.Context()
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }()

			issuesClient := issues.NewIssuesClient(baseClient)

			params := &issues.ListIssuesForOrgParams{
				Limit:         &limit,
				StartingAfter: &startingAfter,
			}

			if severity != "" {
				sevParam := issues.ListIssuesForOrgParamsSeverity(severity)
				params.Severity = &sevParam
			}
			if issueType != "" {
				typeParam := issues.ListIssuesForOrgParamsType(issueType)
				params.Type = &typeParam
			}
			if status != "" {
				statusParam := issues.ListIssuesForOrgParamsStatus(status)
				params.Status = &statusParam
			}

			issueList, err := issuesClient.ListIssuesForOrg(ctx, orgID, params)
			if err != nil {
				return fmt.Errorf("failed to list issues for organization %s: %w", orgID, err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(issueList)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().StringVar(&startingAfter, "starting-after", "", "Cursor for pagination")
	cmd.Flags().StringVar(&severity, "severity", "", "Filter by severity (critical, high, medium, low)")
	cmd.Flags().StringVar(&issueType, "type", "", "Filter by issue type (vuln, license)")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (open, resolved, ignored)")

	return cmd
}

func newIssuesGetCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "get <issue-id>",
		Short: "Get issue by ID",
		Args:  cobra.ExactArgs(1),
		Long: `Get a single issue by its ID within a specified Snyk organization.

Examples:
  # Get issue details
  snyk-api issues get <issue-id> --org-id=<org-id>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("organization ID is required. Use --org-id flag")
			}

			ctx := cmd.Context()
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }()

			issuesClient := issues.NewIssuesClient(baseClient)
			issueID := args[0]

			issue, err := issuesClient.GetIssue(ctx, orgID, issueID)
			if err != nil {
				return fmt.Errorf("failed to get issue %s for organization %s: %w", issueID, orgID, err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(issue)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")

	return cmd
}

