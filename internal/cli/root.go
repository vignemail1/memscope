package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memscope",
		Short: "Hardware memory diagnostics and analysis tool",
		Long: `memscope is a comprehensive tool for memory diagnostics and analysis.
It provides detailed information about memory modules, profiles, and system configuration
without making any modifications to BIOS or firmware settings.`,
	}

	// Add subcommands
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInspectCmd())
	cmd.AddCommand(newMemoryCmd())
	cmd.AddCommand(newRecommendCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newSnapshotCmd())
	cmd.AddCommand(newCompareCmd())

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "memscope development version")
			return nil
		},
	}
}