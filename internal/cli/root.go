package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memscope",
		Short: "Hardware memory diagnostics and analysis tool",
		Long: `memscope is a comprehensive tool for memory diagnostics and analysis.
It provides detailed information about memory modules, profiles, and system configuration
without making any modifications to BIOS or firmware settings.

Examples:
  memscope inspect                    # System inspection and analysis
  memscope doctor                     # Comprehensive health check  
  memscope recommend                  # BIOS optimization recommendations
  memscope export --format json      # Export system data
  memscope snapshot create --name baseline    # Create system snapshot`,
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
	var showDetailed bool
	
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			version := cmd.Root().Version
			if version == "" {
				version = "dev"
			}
			
			if showDetailed {
				fmt.Fprintf(cmd.OutOrStdout(), "memscope %s\n", version)
				fmt.Fprintf(cmd.OutOrStdout(), "Built with: %s\n", runtime.Version())
				fmt.Fprintf(cmd.OutOrStdout(), "Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), version)
			}
			return nil
		},
	}
	
	cmd.Flags().BoolVarP(&showDetailed, "detailed", "d", false, "Show detailed version information")
	
	return cmd
}