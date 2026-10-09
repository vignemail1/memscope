package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newCompareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "compare [snapshot1] [snapshot2]",
		Short: "Compare two system snapshots",
		Long:  "Analyzes differences between system snapshots for before/after analysis",
		// Remove Args: cobra.ExactArgs(2) - this wasn't specified
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 {
				fmt.Fprintf(cmd.OutOrStdout(), "Comparing %s vs %s (placeholder)\n", args[0], args[1])
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "System Snapshot Comparison (placeholder)")
			}
			return nil
		},
	}
}