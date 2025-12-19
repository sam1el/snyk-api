package main

import (
	"os"

	"github.com/sam1el/snyk-api/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
