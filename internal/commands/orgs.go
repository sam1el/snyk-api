package commands

import (
	"context"
	"fmt"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/orgs"
	"github.com/sam1el/snyk-api/pkg/client"
)

func newOrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "orgs",
		Aliases: []string{"org", "organizations"},
		Short:   "Manage Snyk organizations",
		Long: `Manage Snyk organizations.

Organizations are the top-level entities in Snyk that contain projects,
members, and settings. Use these commands to list and view organizations
you have access to.`,
	}

	cmd.AddCommand(newOrgsListCmd())
	cmd.AddCommand(newOrgsGetCmd())

	return cmd
}

func newOrgsListCmd() *cobra.Command {
	var limit int
	var startingAfter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		Long: `List all organizations you have access to.

Examples:
  # List first 10 organizations
  snyk-api orgs list

  # List 50 organizations
  snyk-api orgs list --limit 50

  # List with table output
  snyk-api orgs list --output table

  # Pagination
  snyk-api orgs list --starting-after <cursor>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Create base client
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }() //nolint:errcheck // Best effort cleanup

			// Create orgs client
			orgsClient := orgs.NewOrgsClient(baseClient)

			// Build parameters
			params := &orgs.ListOrganizationsParams{
				Limit: &limit,
			}
			if startingAfter != "" {
				params.StartingAfter = &startingAfter
			}

			// List organizations
			result, err := orgsClient.ListOrganizations(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to list organizations: %w", err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().StringVar(&startingAfter, "starting-after", "", "Cursor for pagination")

	return cmd
}

func newOrgsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <org-id>",
		Short: "Get organization by ID",
		Long: `Get detailed information about a specific organization.

Examples:
  # Get organization
  snyk-api orgs get 00000000-0000-0000-0000-000000000000

  # Get with YAML output
  snyk-api orgs get 00000000-0000-0000-0000-000000000000 --output yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			orgID := args[0]

			// Create base client
			baseClient, err := createClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer func() { _ = baseClient.Close() }() //nolint:errcheck // Best effort cleanup

			// Create orgs client
			orgsClient := orgs.NewOrgsClient(baseClient)

			// Parse org ID as UUID
			orgUUID, err := parseUUID(orgID)
			if err != nil {
				return fmt.Errorf("invalid organization ID: %w", err)
			}

			// Get organization
			org, err := orgsClient.GetOrganization(ctx, orgUUID)
			if err != nil {
				return fmt.Errorf("failed to get organization: %w", err)
			}

			// Format output
			formatter := output.New(getOutputFormat())
			return formatter.Print(org)
		},
	}

	return cmd
}

// createClient creates a configured API client.
func createClient(ctx context.Context) (*client.Client, error) {
	opts := []client.Option{}

	// Add debug logging if enabled
	// TODO: Add logger configuration when debug flag is set
	_ = isDebug()

	// Override API URL if specified
	if apiURL := getAPIURL(); apiURL != "" {
		opts = append(opts, client.WithBaseURL(apiURL, apiURL+"/rest"))
	}

	// Override API version if specified
	if apiVersion := getAPIVersion(); apiVersion != "" {
		opts = append(opts, client.WithVersion(client.APIVersion(apiVersion)))
	}

	return client.New(ctx, opts...)
}

// parseUUID parses a string as a UUID.
func parseUUID(s string) (openapi_types.UUID, error) {
	var uuid openapi_types.UUID
	err := uuid.UnmarshalText([]byte(s))
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return uuid, nil
}
