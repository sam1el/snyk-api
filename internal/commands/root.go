// Package commands provides the CLI command structure for snyk-api.
package commands

import (
	"github.com/spf13/cobra"
)

var (
	// Global flags
	debug        bool
	apiURL       string
	version      string
	outputFormat string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "snyk-api",
	Short: "Comprehensive Snyk API management tool",
	Long: `snyk-api is a comprehensive tool for interacting with Snyk APIs.

Built on top of Snyk's go-application-framework, it provides:
  • Type-safe API clients generated from OpenAPI specifications
  • Rate limiting and retry logic with exponential backoff
  • Support for both REST and v1 APIs
  • Multiple output formats (JSON, YAML, table)
  • Integration with official Snyk CLI

Authentication:
  Set SNYK_TOKEN environment variable with your Snyk API token.

Examples:
  # List organizations
  snyk-api orgs list

  # Get a specific organization
  snyk-api orgs get <org-id>

  # List with custom output format
  snyk-api orgs list --output table

For more information, visit: https://github.com/sam1el/snyk-api`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "Override Snyk API URL")
	rootCmd.PersistentFlags().StringVar(&version, "api-version", "", "Snyk API version (e.g., 2025-11-05)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format: json, yaml, table")

	// Add subcommands
	rootCmd.AddCommand(newOrgsCmd())
	rootCmd.AddCommand(newProjectsCmd())
	rootCmd.AddCommand(newTargetsCmd())
	rootCmd.AddCommand(newIssuesCmd())
	rootCmd.AddCommand(newVersionCmd())
}

// getOutputFormat returns the configured output format.
func getOutputFormat() string {
	return outputFormat
}

// isDebug returns whether debug mode is enabled.
func isDebug() bool {
	return debug
}

// getAPIURL returns the configured API URL.
func getAPIURL() string {
	return apiURL
}

// getAPIVersion returns the configured API version.
func getAPIVersion() string {
	return version
}
