package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/projects"
)

func newProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"project", "proj"},
		Short:   "Manage Snyk projects",
		Long: `Manage Snyk projects.

Projects represent monitored code repositories or packages in Snyk.
Each project contains vulnerability findings and configuration.`,
	}

	cmd.AddCommand(newProjectsListCmd())
	cmd.AddCommand(newProjectsGetCmd())
	cmd.AddCommand(newProjectsDeleteCmd())

	return cmd
}

func newProjectsListCmd() *cobra.Command {
	var orgID string
	var limit int
	var startingAfter string
	var targetID string
	var origin string
	var projectType string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects in an organization",
		Long: `List all projects in a Snyk organization.

Examples:
  # List first 10 projects
  snyk-api projects list --org-id=<org-id>

  # List 50 projects
  snyk-api projects list --org-id=<org-id> --limit 50

  # Filter by origin
  snyk-api projects list --org-id=<org-id> --origin github

  # Filter by type
  snyk-api projects list --org-id=<org-id> --type npm

  # List with table output
  snyk-api projects list --org-id=<org-id> --output table`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("--org-id is required")
			}

			ctx := context.Background()

			// Create base client
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }() //nolint:errcheck // Best effort cleanup

			// Create projects client
			projectsClient := projects.NewProjectsClient(baseClient)

			// Build parameters
			params := &projects.ListProjectsParams{
				Limit: &limit,
			}
			if startingAfter != "" {
				params.StartingAfter = &startingAfter
			}
			if targetID != "" {
				targetUUID, err := parseUUID(targetID)
				if err != nil {
					return fmt.Errorf("invalid target ID: %w", err)
				}
				params.TargetId = &targetUUID
			}
			if origin != "" {
				params.Origin = &origin
			}
			if projectType != "" {
				params.Type = &projectType
			}

			// List projects
			result, err := projectsClient.ListProjects(ctx, orgID, params)
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(result)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().StringVar(&startingAfter, "starting-after", "", "Cursor for pagination")
	cmd.Flags().StringVar(&targetID, "target-id", "", "Filter by target ID")
	cmd.Flags().StringVar(&origin, "origin", "", "Filter by origin (e.g., github, cli)")
	cmd.Flags().StringVar(&projectType, "type", "", "Filter by project type (e.g., npm, maven)")

	return cmd
}

func newProjectsGetCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "get <project-id>",
		Short: "Get project by ID",
		Long: `Get detailed information about a specific project.

Examples:
  # Get project
  snyk-api projects get <project-id> --org-id=<org-id>

  # Get with YAML output
  snyk-api projects get <project-id> --org-id=<org-id> --output yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("--org-id is required")
			}

			ctx := context.Background()
			projectID := args[0]

			// Create base client
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }() //nolint:errcheck // Best effort cleanup

			// Create projects client
			projectsClient := projects.NewProjectsClient(baseClient)

			// Get project
			project, err := projectsClient.GetProject(ctx, orgID, projectID)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(project)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")

	return cmd
}

func newProjectsDeleteCmd() *cobra.Command {
	var orgID string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <project-id>",
		Short: "Delete a project",
		Long: `Delete a project from Snyk.

WARNING: This action cannot be undone!

Examples:
  # Delete project (with confirmation prompt)
  snyk-api projects delete <project-id> --org-id=<org-id>

  # Delete project (skip confirmation)
  snyk-api projects delete <project-id> --org-id=<org-id> --confirm`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if orgID == "" {
				return fmt.Errorf("--org-id is required")
			}

			projectID := args[0]

			// Confirmation check
			if !confirm {
				fmt.Printf("Are you sure you want to delete project %s? (yes/no): ", projectID)
				var response string
				_, err := fmt.Scanln(&response)
				if err != nil || (response != "yes" && response != "y") {
					fmt.Println("Deletion cancelled.")
					return nil
				}
			}

			ctx := context.Background()

			// Create base client
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }() //nolint:errcheck // Best effort cleanup

			// Create projects client
			projectsClient := projects.NewProjectsClient(baseClient)

			// Delete project
			err = projectsClient.DeleteProject(ctx, orgID, projectID)
			if err != nil {
				return fmt.Errorf("failed to delete project: %w", err)
			}

			fmt.Printf("Project %s deleted successfully.\n", projectID)
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Skip confirmation prompt")

	return cmd
}
