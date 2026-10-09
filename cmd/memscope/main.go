package main

import (
	"os"

	"github.com/vignemail1/memscope/internal/cli"
)

// Build information set by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Pass version information to CLI
	cmd := cli.NewRootCmd()
	cmd.Version = version
	
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
