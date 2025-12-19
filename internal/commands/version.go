package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information (set by build flags)
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print detailed version information including git commit and build date.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("snyk-api version %s\n", Version)
			fmt.Printf("  Git commit: %s\n", GitCommit)
			fmt.Printf("  Built:      %s\n", BuildDate)
			fmt.Printf("  Module:     github.com/sam1el/snyk-api\n")
		},
	}
}
