// Package cli implements the command-line entry point.
package cli

import (
	"fmt"
	"io"
)

const help = `memscope - hardware and memory diagnostics (development scaffold)

Usage: memscope --help | --version

Planned commands (not implemented):
  inspect
  memory current
  memory profiles
  recommend
  doctor
  export
  snapshot
  compare
`

// Run executes the scaffold CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "--help" || args[0] == "-h")) {
		if _, err := io.WriteString(stdout, help); err != nil {
			return 1
		}
		return 0
	}
	if len(args) == 1 && args[0] == "--version" {
		if _, err := fmt.Fprintf(stdout, "memscope %s\n", version); err != nil {
			return 1
		}
		return 0
	}
	switch args[0] {
	case "inspect", "memory", "recommend", "doctor", "export", "snapshot", "compare":
		if _, err := fmt.Fprintln(stderr, "requested command is not implemented; no hardware was accessed"); err != nil {
			return 1
		}
		return 3
	default:
		if _, err := fmt.Fprintf(stderr, "unknown command: %q\n", args[0]); err != nil {
			return 1
		}
		return 2
	}
}
