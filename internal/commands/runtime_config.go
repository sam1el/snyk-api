package commands

import (
	"os"
	"strconv"

	"github.com/sam1el/snyk-api/internal/output"
	cfgpkg "github.com/sam1el/snyk-api/pkg/config"
	"github.com/spf13/cobra"
)

func applyConfigPreRun(cmd *cobra.Command, _ []string) error {
	res, err := resolveActiveConfig(cmd)
	if err != nil {
		return err
	}

	// Propagate CLI overrides via env so lower layers (client) can see them.
	if cmd.Flags().Changed("profile") && res.ProfileName != "" {
		_ = os.Setenv("SNYK_API_PROFILE", res.ProfileName)
	}
	if cmd.Flags().Changed("api-url") && res.APIURL != "" {
		_ = os.Setenv("SNYK_API", res.APIURL)
	}
	if cmd.Flags().Changed("api-version") && res.APIVersion != "" {
		_ = os.Setenv("SNYK_API_VERSION", res.APIVersion)
	}
	if cmd.Flags().Changed("output") && res.Output != "" {
		_ = os.Setenv("SNYK_OUTPUT", res.Output)
	}
	if cmd.Flags().Changed("debug") {
		_ = os.Setenv("SNYK_DEBUG", strconv.FormatBool(res.Debug))
	}

	if !cmd.Flags().Changed("output") && res.Output != "" {
		outputFormat = res.Output
	}

	if getTemplate() != "" {
		output.SetDefaultTemplate(getTemplate())
	}
	if getJQ() != "" {
		output.SetDefaultJQ(getJQ())
	}

	return nil
}

func resolveActiveConfig(cmd *cobra.Command) (cfgpkg.Resolved, error) {
	path, err := cfgpkg.DefaultPath()
	if err != nil {
		return cfgpkg.Resolved{}, err
	}
	fileCfg, err := cfgpkg.LoadFile(path)
	if err != nil {
		return cfgpkg.Resolved{}, err
	}

	env := envMap()
	overrides := cfgpkg.FlagOverrides{}

	if cmd.Flags().Changed("profile") && getProfile() != "" {
		overrides.Profile = getProfile()
	}
	if cmd.Flags().Changed("api-url") && getAPIURL() != "" {
		overrides.APIURL = getAPIURL()
	}
	if cmd.Flags().Changed("api-version") && getAPIVersion() != "" {
		overrides.APIVersion = getAPIVersion()
	}
	if cmd.Flags().Changed("output") && getOutputFormat() != "" {
		overrides.Output = getOutputFormat()
	}
	if cmd.Flags().Changed("debug") {
		dbg := isDebug()
		overrides.Debug = &dbg
	}

	res := cfgpkg.Resolve(fileCfg, overrides, env)
	return res, nil
}
