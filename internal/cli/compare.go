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
		Args:  cobra.ExactArgs(2), // Restore this validation from original spec
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Comparing %s vs %s (placeholder)\n", args[0], args[1])
			return nil
		},
	}
}