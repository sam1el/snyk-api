package commands

import (
	"fmt"

	"github.com/sam1el/snyk-api/pkg/config"
	"github.com/spf13/cobra"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage default org/group/project values",
	}

	cmd.AddCommand(
		newContextShowCmd(),
		newContextSetCmd(),
		newContextClearCmd(),
	)

	return cmd
}

func newContextShowCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show resolved context (org/group/project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}
			env := envMap()
			res := config.Resolve(cfg, config.FlagOverrides{Profile: profile}, env)

			fmt.Printf("Profile: %s\n", res.ProfileName)
			fmt.Printf("Org ID: %s\n", res.OrgID)
			fmt.Printf("Group ID: %s\n", res.GroupID)
			fmt.Printf("Project ID: %s\n", res.ProjectID)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile to resolve (defaults to current)")
	return cmd
}

func newContextSetCmd() *cobra.Command {
	var profile string
	var orgID string
	var groupID string
	var projectID string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set default org/group/project for a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}

			targetProfile := profile
			if targetProfile == "" {
				targetProfile = cfg.Current
			}
			if targetProfile == "" {
				targetProfile = "default"
			}
			if cfg.Profiles == nil {
				cfg.Profiles = map[string]config.Profile{}
			}
			p := cfg.Profiles[targetProfile]

			if orgID != "" {
				p.OrgID = orgID
			}
			if groupID != "" {
				p.GroupID = groupID
			}
			if projectID != "" {
				p.ProjectID = projectID
			}

			cfg.Profiles[targetProfile] = p
			if cfg.Current == "" {
				cfg.Current = targetProfile
			}

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Context updated for profile %s\n", targetProfile)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Default organization ID")
	cmd.Flags().StringVar(&groupID, "group-id", "", "Default group ID")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Default project ID")
	return cmd
}

func newContextClearCmd() *cobra.Command {
	var profile string
	var orgID bool
	var groupID bool
	var projectID bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear default org/group/project for a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}

			targetProfile := profile
			if targetProfile == "" {
				targetProfile = cfg.Current
			}
			if targetProfile == "" {
				targetProfile = "default"
			}

			p, ok := cfg.Profiles[targetProfile]
			if !ok {
				return fmt.Errorf("profile %s not found", targetProfile)
			}

			if orgID {
				p.OrgID = ""
			}
			if groupID {
				p.GroupID = ""
			}
			if projectID {
				p.ProjectID = ""
			}

			cfg.Profiles[targetProfile] = p

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Context cleared for profile %s\n", targetProfile)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	cmd.Flags().BoolVar(&orgID, "org-id", false, "Clear organization ID")
	cmd.Flags().BoolVar(&groupID, "group-id", false, "Clear group ID")
	cmd.Flags().BoolVar(&projectID, "project-id", false, "Clear project ID")
	return cmd
}
