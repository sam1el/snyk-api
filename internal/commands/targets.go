package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/targets"
)

// newTargetsCmd creates the targets subcommand.
func newTargetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "targets",
		Aliases: []string{"target"},
		Short:   "Manage Snyk targets.",
		Long: `Manage Snyk targets (repositories and container registries).

Targets represent source control repositories or container registries
that have been imported into Snyk for monitoring.`,
	}

	cmd.AddCommand(newTargetsListCmd())
	cmd.AddCommand(newTargetsGetCmd())
	cmd.AddCommand(newTargetsDeleteCmd())

	return cmd
}

func newTargetsListCmd() *cobra.Command {
	var orgID string
	var limit int
	var startingAfter string
	var origin string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List targets in an organization",
		Long: `List all targets within a specified Snyk organization.

Examples:
  # List first 10 targets in an organization
  snyk-api targets list --org-id=<org-id>

  # List targets with filters
  snyk-api targets list --org-id=<org-id> --origin github

  # List with table output
  snyk-api targets list --org-id=<org-id> --output table

  # Pagination
  snyk-api targets list --org-id=<org-id> --starting-after <cursor>
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

			targetsClient := targets.NewTargetsClient(baseClient)

			params := &targets.ListTargetsParams{
				Limit:         &limit,
				StartingAfter: &startingAfter,
			}

			if origin != "" {
				params.Origin = &origin
			}

			targetList, err := targetsClient.ListTargets(ctx, orgID, params)
			if err != nil {
				return fmt.Errorf("failed to list targets for organization %s: %w", orgID, err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(targetList)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().StringVar(&startingAfter, "starting-after", "", "Cursor for pagination")
	cmd.Flags().StringVar(&origin, "origin", "", "Filter by origin (e.g., github, gitlab, cli)")

	return cmd
}

func newTargetsGetCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "get <target-id>",
		Short: "Get target by ID",
		Args:  cobra.ExactArgs(1),
		Long: `Get a single target by its ID within a specified Snyk organization.

Examples:
  # Get target details
  snyk-api targets get <target-id> --org-id=<org-id>
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

			targetsClient := targets.NewTargetsClient(baseClient)
			targetID := args[0]

			target, err := targetsClient.GetTarget(ctx, orgID, targetID)
			if err != nil {
				return fmt.Errorf("failed to get target %s for organization %s: %w", targetID, orgID, err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(target)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")

	return cmd
}

func newTargetsDeleteCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "delete <target-id>",
		Short: "Delete a target",
		Args:  cobra.ExactArgs(1),
		Long: `Delete a target by its ID within a specified Snyk organization.

Examples:
  # Delete a target
  snyk-api targets delete <target-id> --org-id=<org-id>
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

			targetsClient := targets.NewTargetsClient(baseClient)
			targetID := args[0]

			err = targetsClient.DeleteTarget(ctx, orgID, targetID)
			if err != nil {
				return fmt.Errorf("failed to delete target %s for organization %s: %w", targetID, orgID, err)
			}

			fmt.Printf("Target %s deleted successfully from organization %s.\n", targetID, orgID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")

	return cmd
}
