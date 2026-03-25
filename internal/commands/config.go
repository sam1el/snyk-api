package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sam1el/snyk-api/pkg/config"
	"github.com/spf13/cobra"
)

// newConfigCmd exposes config management commands.
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage snyk-api configuration profiles",
	}

	cmd.AddCommand(
		newConfigGetCmd(),
		newConfigSetCmd(),
		newConfigListCmd(),
		newConfigUseProfileCmd(),
		newConfigCurrentCmd(),
	)

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get a config value for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(args[0])

			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}
			if profile == "" {
				profile = cfg.Current
			}
			if profile == "" {
				profile = "default"
			}
			p := cfg.Profiles[profile]

			val, ok := getProfileValue(p, key)
			if !ok {
				return fmt.Errorf("unknown key: %s", key)
			}
			fmt.Println(val)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	return cmd
}

func newConfigSetCmd() *cobra.Command {
	var profile string
	var allowSecret bool

	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a config value for a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(args[0])
			val := args[1]

			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}

			if profile == "" {
				profile = cfg.Current
			}
			if profile == "" {
				profile = "default"
			}
			if cfg.Profiles == nil {
				cfg.Profiles = map[string]config.Profile{}
			}
			p := cfg.Profiles[profile]

			if err := setProfileValue(&p, key, val, allowSecret); err != nil {
				return err
			}

			cfg.Profiles[profile] = p
			if cfg.Current == "" {
				cfg.Current = profile
			}

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Set %s for profile %s\n", key, profile)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name (defaults to current)")
	cmd.Flags().BoolVar(&allowSecret, "with-secret", false, "Allow setting token values")
	return cmd
}

func newConfigListCmd() *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List profiles or values within a profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}

			if profile == "" {
				for name := range cfg.Profiles {
					currentMark := ""
					if cfg.Current == name {
						currentMark = " (current)"
					}
					fmt.Printf("%s%s\n", name, currentMark)
				}
				if len(cfg.Profiles) == 0 {
					fmt.Println("no profiles set")
				}
				return nil
			}

			p, ok := cfg.Profiles[profile]
			if !ok {
				return fmt.Errorf("profile %s not found", profile)
			}

			printProfileValues(profile, p)
			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Profile name to inspect")
	return cmd
}

func newConfigUseProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-profile [name]",
		Short: "Set the current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}
			if cfg.Profiles == nil {
				cfg.Profiles = map[string]config.Profile{}
			}
			if _, ok := cfg.Profiles[name]; !ok {
				cfg.Profiles[name] = config.Profile{}
			}
			cfg.Current = name

			if err := config.SaveFile(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Switched to profile %s\n", name)
			return nil
		},
	}

	return cmd
}

func newConfigCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the current profile name",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.LoadFile(path)
			if err != nil {
				return err
			}
			profile := cfg.Current
			if profile == "" {
				profile = "default"
			}
			fmt.Println(profile)
			return nil
		},
	}

	return cmd
}

func getProfileValue(p config.Profile, key string) (string, bool) {
	switch key {
	case "token":
		return p.Token, true
	case "api_url":
		return p.APIURL, true
	case "rest_api_url":
		return p.RestAPIURL, true
	case "api_version":
		return p.APIVersion, true
	case "org_id":
		return p.OrgID, true
	case "group_id":
		return p.GroupID, true
	case "project_id":
		return p.ProjectID, true
	case "output":
		return p.Output, true
	case "page_size":
		if p.PageSize == 0 {
			return "", true
		}
		return strconv.Itoa(p.PageSize), true
	case "debug":
		return strconv.FormatBool(p.Debug), true
	default:
		return "", false
	}
}

func setProfileValue(p *config.Profile, key, val string, allowSecret bool) error {
	switch key {
	case "token":
		if !allowSecret {
			return fmt.Errorf("refusing to set token without --with-secret")
		}
		p.Token = val
	case "api_url":
		p.APIURL = val
	case "rest_api_url":
		p.RestAPIURL = val
	case "api_version":
		p.APIVersion = val
	case "org_id":
		p.OrgID = val
	case "group_id":
		p.GroupID = val
	case "project_id":
		p.ProjectID = val
	case "output":
		p.Output = val
	case "page_size":
		n, err := strconv.Atoi(val)
		if err != nil || n <= 0 {
			return fmt.Errorf("invalid page_size: %s", val)
		}
		p.PageSize = n
	case "debug":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid debug value: %s", val)
		}
		p.Debug = b
	default:
		return fmt.Errorf("unknown key: %s", key)
	}
	return nil
}

func printProfileValues(name string, p config.Profile) {
	fmt.Printf("Profile: %s\n", name)
	fmt.Printf("  api_url: %s\n", p.APIURL)
	fmt.Printf("  rest_api_url: %s\n", p.RestAPIURL)
	fmt.Printf("  api_version: %s\n", p.APIVersion)
	fmt.Printf("  org_id: %s\n", p.OrgID)
	fmt.Printf("  group_id: %s\n", p.GroupID)
	fmt.Printf("  output: %s\n", p.Output)
	if p.PageSize > 0 {
		fmt.Printf("  page_size: %d\n", p.PageSize)
	} else {
		fmt.Printf("  page_size: \n")
	}
	fmt.Printf("  debug: %t\n", p.Debug)
	if p.Token != "" {
		fmt.Printf("  token: [set]\n")
	} else {
		fmt.Printf("  token: \n")
	}
}
