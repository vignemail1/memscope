package main

import (
	"os"

	"github.com/vignemail1/memscope/internal/cli"
)

var version = "dev"

func main() {
	cmd := cli.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
