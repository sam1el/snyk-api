package commands

import (
	"context"
	"fmt"

	"github.com/sam1el/snyk-api/internal/output"
	"github.com/sam1el/snyk-api/pkg/apiclients/rest"
	restgroups "github.com/sam1el/snyk-api/pkg/apiclients/rest/groups"
	restorgs "github.com/sam1el/snyk-api/pkg/apiclients/rest/orgs"
	restself "github.com/sam1el/snyk-api/pkg/apiclients/rest/self"
	resttenants "github.com/sam1el/snyk-api/pkg/apiclients/rest/tenants"
	"github.com/sam1el/snyk-api/pkg/client"
	"github.com/spf13/cobra"
)

// newRESTCmd creates the REST command group for full REST API access.
func newRESTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rest",
		Aliases: []string{"api"},
		Short:   "Full REST API access",
		Long: `Direct access to all Snyk REST API endpoints.

This provides complete coverage of Snyk's REST API including:
  - Organizations (memberships, invites, service accounts, policies, collections, etc.)
  - Groups (memberships, policies, SSO, assets, etc.)
  - Tenants (broker deployments, connections, credentials)
  - Self (current user, apps, sessions)`,
	}

	cmd.AddCommand(
		newRESTOrgsCmd(),
		newRESTGroupsCmd(),
		newRESTTenantsCmd(),
		newRESTSelfCmd(),
	)

	return cmd
}

// ============================================================================
// REST Organizations Commands
// ============================================================================

func newRESTOrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "REST Organizations API",
	}

	cmd.AddCommand(
		newRESTOrgsListCmd(),
		newRESTOrgsGetCmd(),
		newRESTOrgsMembershipsCmd(),
		newRESTOrgsInvitesCmd(),
		newRESTOrgsServiceAccountsCmd(),
		newRESTOrgsPoliciesCmd(),
		newRESTOrgsCollectionsCmd(),
		newRESTOrgsSettingsCmd(),
		newRESTOrgsAuditLogsCmd(),
		newRESTOrgsProjectsCmd(),
		newRESTOrgsTargetsCmd(),
		newRESTOrgsIssuesCmd(),
		newRESTOrgsSBOMCmd(),
		newRESTOrgsContainerImagesCmd(),
		newRESTOrgsAppsCmd(),
		newRESTOrgsCloudCmd(),
		newRESTOrgsExportCmd(),
	)

	return cmd
}

func newRESTOrgsListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			params := &rest.PaginationParams{Limit: limit}
			resp, err := c.List(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to list orgs: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 100, "Number of results")

	return cmd
}

func newRESTOrgsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [org-id]",
		Short: "Get organization by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.Get(ctx, args[0])
			if err != nil {
				return fmt.Errorf("failed to get org: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	return cmd
}

// Memberships
func newRESTOrgsMembershipsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memberships",
		Short: "Manage organization memberships",
	}

	cmd.AddCommand(
		newRESTOrgsMembershipsListCmd(),
		newRESTOrgsMembershipsGetCmd(),
		newRESTOrgsMembershipsUpdateCmd(),
		newRESTOrgsMembershipsDeleteCmd(),
	)

	return cmd
}

func newRESTOrgsMembershipsListCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organization memberships",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.ListMemberships(ctx, orgID, nil)
			if err != nil {
				return fmt.Errorf("failed to list memberships: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsMembershipsGetCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "get [membership-id]",
		Short: "Get membership by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.GetMembership(ctx, orgID, args[0])
			if err != nil {
				return fmt.Errorf("failed to get membership: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsMembershipsUpdateCmd() *cobra.Command {
	var orgID, role string

	cmd := &cobra.Command{
		Use:   "update [membership-id]",
		Short: "Update membership role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.UpdateMembership(ctx, orgID, args[0], role)
			if err != nil {
				return fmt.Errorf("failed to update membership: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&role, "role", "", "New role (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func newRESTOrgsMembershipsDeleteCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "delete [membership-id]",
		Short: "Remove membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			if err := c.DeleteMembership(ctx, orgID, args[0]); err != nil {
				return fmt.Errorf("failed to delete membership: %w", err)
			}

			fmt.Println("Membership deleted successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

// Invites
func newRESTOrgsInvitesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invites",
		Short: "Manage organization invites",
	}

	cmd.AddCommand(
		newRESTOrgsInvitesListCmd(),
		newRESTOrgsInvitesCreateCmd(),
		newRESTOrgsInvitesDeleteCmd(),
	)

	return cmd
}

func newRESTOrgsInvitesListCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pending invites",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.ListInvites(ctx, orgID, nil)
			if err != nil {
				return fmt.Errorf("failed to list invites: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsInvitesCreateCmd() *cobra.Command {
	var orgID, email, role string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.CreateInvite(ctx, orgID, email, role)
			if err != nil {
				return fmt.Errorf("failed to create invite: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	cmd.Flags().StringVar(&email, "email", "", "Email to invite (required)")
	cmd.Flags().StringVar(&role, "role", "collaborator", "Role for the invite")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("email")

	return cmd
}

func newRESTOrgsInvitesDeleteCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "delete [invite-id]",
		Short: "Cancel an invite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			if err := c.DeleteInvite(ctx, orgID, args[0]); err != nil {
				return fmt.Errorf("failed to delete invite: %w", err)
			}

			fmt.Println("Invite cancelled successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

// Service Accounts
func newRESTOrgsServiceAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-accounts",
		Short: "Manage service accounts",
	}

	cmd.AddCommand(
		newRESTOrgsServiceAccountsListCmd(),
		newRESTOrgsServiceAccountsGetCmd(),
		newRESTOrgsServiceAccountsDeleteCmd(),
		newRESTOrgsServiceAccountsRotateCmd(),
	)

	return cmd
}

func newRESTOrgsServiceAccountsListCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List service accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.ListServiceAccounts(ctx, orgID, nil)
			if err != nil {
				return fmt.Errorf("failed to list service accounts: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsServiceAccountsGetCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "get [service-account-id]",
		Short: "Get service account by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.GetServiceAccount(ctx, orgID, args[0])
			if err != nil {
				return fmt.Errorf("failed to get service account: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsServiceAccountsDeleteCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "delete [service-account-id]",
		Short: "Delete service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			if err := c.DeleteServiceAccount(ctx, orgID, args[0]); err != nil {
				return fmt.Errorf("failed to delete service account: %w", err)
			}

			fmt.Println("Service account deleted successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsServiceAccountsRotateCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "rotate [service-account-id]",
		Short: "Rotate service account secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			defer baseClient.Close()

			c := restorgs.New(baseClient)
			resp, err := c.RotateServiceAccountSecret(ctx, orgID, args[0])
			if err != nil {
				return fmt.Errorf("failed to rotate secret: %w", err)
			}

			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

// Simplified commands for other resources (patterns are similar)

func newRESTOrgsPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policies",
		Short: "Manage organization policies",
	}

	var orgID string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListPolicies(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	listCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = listCmd.MarkFlagRequired("org-id")

	cmd.AddCommand(listCmd)
	return cmd
}

func newRESTOrgsCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections",
		Short: "Manage collections",
	}

	var orgID string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListCollections(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	listCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = listCmd.MarkFlagRequired("org-id")

	cmd.AddCommand(listCmd)
	return cmd
}

func newRESTOrgsSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage organization settings",
	}

	var orgID string

	iacCmd := &cobra.Command{
		Use:   "iac",
		Short: "Get IaC settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.GetIaCSettings(ctx, orgID)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	iacCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = iacCmd.MarkFlagRequired("org-id")

	sastCmd := &cobra.Command{
		Use:   "sast",
		Short: "Get SAST settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.GetSASTSettings(ctx, orgID)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	sastCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = sastCmd.MarkFlagRequired("org-id")

	cmd.AddCommand(iacCmd, sastCmd)
	return cmd
}

func newRESTOrgsAuditLogsCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "audit-logs",
		Short: "Search audit logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.SearchAuditLogs(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsProjectsCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListProjects(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsTargetsCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "targets",
		Short: "List targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListTargets(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsIssuesCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "issues",
		Short: "List issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListIssues(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsSBOMCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "sbom",
		Short: "SBOM testing",
	}

	listJobsCmd := &cobra.Command{
		Use:   "get-job [job-id]",
		Short: "Get SBOM test job status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.GetSBOMTestJob(ctx, orgID, args[0])
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	listJobsCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = listJobsCmd.MarkFlagRequired("org-id")

	cmd.AddCommand(listJobsCmd)
	return cmd
}

func newRESTOrgsContainerImagesCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "container-images",
		Short: "List container images",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListContainerImages(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsAppsCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "apps",
		Short: "List apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListApps(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("org-id")

	return cmd
}

func newRESTOrgsCloudCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Cloud environments and scans",
	}

	envCmd := &cobra.Command{
		Use:   "environments",
		Short: "List cloud environments",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListCloudEnvironments(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	envCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = envCmd.MarkFlagRequired("org-id")

	scansCmd := &cobra.Command{
		Use:   "scans",
		Short: "List cloud scans",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.ListCloudScans(ctx, orgID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	scansCmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = scansCmd.MarkFlagRequired("org-id")

	cmd.AddCommand(envCmd, scansCmd)
	return cmd
}

func newRESTOrgsExportCmd() *cobra.Command {
	var orgID string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Create data export",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restorgs.New(baseClient)
			resp, err := c.CreateExport(ctx, orgID)
			if err != nil {
				return err
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
// REST Groups Commands
// ============================================================================

func newRESTGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "REST Groups API",
	}

	cmd.AddCommand(
		newRESTGroupsListCmd(),
		newRESTGroupsGetCmd(),
		newRESTGroupsMembershipsCmd(),
		newRESTGroupsOrgsCmd(),
		newRESTGroupsPoliciesCmd(),
		newRESTGroupsServiceAccountsCmd(),
		newRESTGroupsIssuesCmd(),
		newRESTGroupsAssetsCmd(),
		newRESTGroupsAuditLogsCmd(),
		newRESTGroupsSSOCmd(),
	)

	return cmd
}

func newRESTGroupsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.List(ctx, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTGroupsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [group-id]",
		Short: "Get group by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.Get(ctx, args[0])
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTGroupsMembershipsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "memberships",
		Short: "List group memberships",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListMemberships(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsOrgsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "List organizations in group",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListOrgs(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsPoliciesCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "policies",
		Short: "List group policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListPolicies(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsServiceAccountsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "service-accounts",
		Short: "List group service accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListServiceAccounts(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsIssuesCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "issues",
		Short: "List group issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListIssues(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsAssetsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Search group assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.SearchAssets(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsAuditLogsCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "audit-logs",
		Short: "Search group audit logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.SearchAuditLogs(ctx, groupID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("group-id")

	return cmd
}

func newRESTGroupsSSOCmd() *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "sso",
		Short: "List SSO connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restgroups.New(baseClient)
			resp, err := c.ListSSOConnections(ctx, groupID, nil)
			if err != nil {
				return err
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
// REST Tenants Commands
// ============================================================================

func newRESTTenantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenants",
		Short: "REST Tenants API",
	}

	cmd.AddCommand(
		newRESTTenantsListCmd(),
		newRESTTenantsGetCmd(),
		newRESTTenantsMembershipsCmd(),
		newRESTTenantsRolesCmd(),
		newRESTTenantsBrokerDeploymentsCmd(),
	)

	return cmd
}

func newRESTTenantsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tenants",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := resttenants.New(baseClient)
			resp, err := c.List(ctx, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTTenantsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [tenant-id]",
		Short: "Get tenant by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := resttenants.New(baseClient)
			resp, err := c.Get(ctx, args[0])
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTTenantsMembershipsCmd() *cobra.Command {
	var tenantID string

	cmd := &cobra.Command{
		Use:   "memberships",
		Short: "List tenant memberships",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := resttenants.New(baseClient)
			resp, err := c.ListMemberships(ctx, tenantID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Tenant ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("tenant-id")

	return cmd
}

func newRESTTenantsRolesCmd() *cobra.Command {
	var tenantID string

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List tenant roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := resttenants.New(baseClient)
			resp, err := c.ListRoles(ctx, tenantID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Tenant ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("tenant-id")

	return cmd
}

func newRESTTenantsBrokerDeploymentsCmd() *cobra.Command {
	var tenantID string

	cmd := &cobra.Command{
		Use:   "broker-deployments",
		Short: "List broker deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := resttenants.New(baseClient)
			resp, err := c.ListBrokerDeployments(ctx, tenantID, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}

	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Tenant ID (required)")
	//nolint:errcheck // Cobra handles this
	_ = cmd.MarkFlagRequired("tenant-id")

	return cmd
}

// ============================================================================
// REST Self Commands
// ============================================================================

func newRESTSelfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self",
		Short: "Current user endpoints",
	}

	cmd.AddCommand(
		newRESTSelfGetCmd(),
		newRESTSelfAppsCmd(),
		newRESTSelfAccessRequestsCmd(),
	)

	return cmd
}

func newRESTSelfGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restself.New(baseClient)
			resp, err := c.Get(ctx)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTSelfAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "List user apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restself.New(baseClient)
			resp, err := c.ListApps(ctx, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}

func newRESTSelfAccessRequestsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-requests",
		Short: "List access requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			baseClient, err := client.New(ctx)
			if err != nil {
				return err
			}
			defer baseClient.Close()
			c := restself.New(baseClient)
			resp, err := c.ListAccessRequests(ctx, nil)
			if err != nil {
				return err
			}
			return output.New(getOutputFormat()).Print(resp)
		},
	}
	return cmd
}
