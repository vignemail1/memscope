package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	runtimecollect "github.com/vignemail1/memscope/internal/collect/runtime"
)

func newMemoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Memory analysis and profile management",
		Long: `Memory commands provide detailed analysis of memory configuration:
- Current active memory settings
- Available memory profiles (JEDEC, XMP, EXPO)
- Profile comparison and validation`,
	}
	
	cmd.AddCommand(newMemoryCurrentCmd())
	cmd.AddCommand(newMemoryProfilesCmd())
	
	return cmd
}

func newMemoryCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Display current active memory configuration",
		Long:  "Shows the currently active memory timings, frequencies, and voltages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS != "windows" {
				fmt.Fprintln(cmd.OutOrStdout(), "Memory parameter collection is only supported on Windows")
				return nil
			}
			
			fmt.Fprintln(cmd.OutOrStdout(), "Collecting current memory parameters...")
			
			provider := runtimecollect.NewProvider()
			params, err := provider.CollectMemoryParameters()
			if err != nil {
				return fmt.Errorf("failed to collect memory parameters: %w", err)
			}
			
			// Display current configuration
			fmt.Fprintln(cmd.OutOrStdout(), "\n=== Current Memory Configuration ===")
			fmt.Fprintf(cmd.OutOrStdout(), "Current Frequency: %d MHz\n", params.CurrentFrequency)
			fmt.Fprintf(cmd.OutOrStdout(), "Configured Speed: %d MHz\n", params.ConfiguredSpeed)
			
			// Display active profile
			if params.Profile != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Active Profile: %s (%s)\n", params.Profile.Name, params.Profile.Type)
				fmt.Fprintf(cmd.OutOrStdout(), "Profile Voltage: %.2f V\n", params.Profile.Voltage)
			}
			
			// Display timings
			if len(params.Timings) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\n=== Memory Timings ===")
				for name, timing := range params.Timings {
					fmt.Fprintf(cmd.OutOrStdout(), "%s: %d %s\n", name, timing.Value, timing.Unit)
				}
			}
			
			// Display voltages
			if len(params.Voltages) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\n=== Voltages ===")
				for name, voltage := range params.Voltages {
					fmt.Fprintf(cmd.OutOrStdout(), "%s: %.3f %s\n", name, voltage.Value, voltage.Unit)
				}
			}
			
			fmt.Fprintf(cmd.OutOrStdout(), "\nData collected at: %s\n", params.Timestamp.Format("2006-01-02 15:04:05"))
			
			return nil
		},
	}
}

func newMemoryProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profiles",
		Short: "Display available memory profiles",
		Long:  "Shows all available memory profiles from SPD data (JEDEC, XMP, EXPO)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Available Memory Profiles (placeholder)")
			fmt.Fprintln(cmd.OutOrStdout(), "JEDEC: [Not implemented]")
			fmt.Fprintln(cmd.OutOrStdout(), "XMP: [Not implemented]")
			return nil
		},
	}
}