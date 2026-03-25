package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sam1el/snyk-api/pkg/config"
	"github.com/spf13/cobra"
)

// newAuthCmd exposes auth helpers similar to gh.
func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate and manage profiles",
	}

	cmd.AddCommand(
		newAuthLoginCmd(),
		newAuthLogoutCmd(),
		newAuthStatusCmd(),
	)

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var profile string
	var token string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Store a token for a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				fmt.Print("Enter SNYK token: ")
				reader := bufio.NewReader(os.Stdin)
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("read token: %w", err)
				}
				token = strings.TrimSpace(input)
			}
			if token == "" {
				return fmt.Errorf("token is required")
			}

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
			p.Token = token
			cfg.Profiles[targetProfile] = p
			if cfg.Current == "" {
				cfg.Current = targetProfile
			}

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Token stored for profile %s\n", targetProfile)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	cmd.Flags().StringVar(&token, "token", "", "Token value (optional; otherwise prompted)")
	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored token for a profile",
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

			p.Token = ""
			cfg.Profiles[targetProfile] = p

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Token removed for profile %s\n", targetProfile)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	return cmd
}

func newAuthStatusCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show active profile and token presence",
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
			env := envMap()
			resolved := config.Resolve(cfg, config.FlagOverrides{Profile: targetProfile}, env)

			if resolved.ProfileName == "" {
				resolved.ProfileName = "default"
			}

			fmt.Printf("Profile: %s\n", resolved.ProfileName)
			if resolved.Token != "" {
				fmt.Println("Token: set")
			} else {
				fmt.Println("Token: not set")
			}
			fmt.Printf("API URL: %s\n", resolved.APIURL)
			fmt.Printf("REST API URL: %s\n", resolved.RestAPIURL)
			fmt.Printf("API Version: %s\n", resolved.APIVersion)
			if resolved.OrgID != "" {
				fmt.Printf("Org ID: %s\n", resolved.OrgID)
			}
			if resolved.GroupID != "" {
				fmt.Printf("Group ID: %s\n", resolved.GroupID)
			}
			fmt.Printf("Output: %s\n", resolved.Output)
			fmt.Printf("Page Size: %d\n", resolved.PageSize)
			fmt.Printf("Debug: %t\n", resolved.Debug)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	return cmd
}
