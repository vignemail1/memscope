package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vignemail1/memscope/internal/collect/inventory"
)

func newInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Display hardware overview and system information",
		Long: `Inspect provides a comprehensive overview of system hardware including:
- CPU information and capabilities
- Motherboard and BIOS details
- Memory modules and configuration
- Platform-specific details`,
		RunE: runInspect,
	}
	
	return cmd
}

func runInspect(cmd *cobra.Command, args []string) error {
	provider := inventory.NewProvider()
	
	inv, err := provider.CollectInventory()
	if err != nil {
		return fmt.Errorf("failed to collect hardware inventory: %w", err)
	}
	
	fmt.Fprintln(cmd.OutOrStdout(), "=== Hardware Inventory ===")
	fmt.Fprintln(cmd.OutOrStdout())
	
	// System Information
	fmt.Fprintln(cmd.OutOrStdout(), "System Information:")
	fmt.Fprintf(cmd.OutOrStdout(), "  Manufacturer: %s\n", inv.System.Manufacturer)
	fmt.Fprintf(cmd.OutOrStdout(), "  Model:        %s\n", inv.System.Model)
	if inv.System.SerialNumber != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  Serial:       %s\n", inv.System.SerialNumber)
	}
	fmt.Fprintln(cmd.OutOrStdout())
	
	// BIOS Information
	fmt.Fprintln(cmd.OutOrStdout(), "BIOS Information:")
	fmt.Fprintf(cmd.OutOrStdout(), "  Vendor:       %s\n", inv.BIOS.Vendor)
	fmt.Fprintf(cmd.OutOrStdout(), "  Version:      %s\n", inv.BIOS.Version)
	if inv.BIOS.Date != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  Date:         %s\n", inv.BIOS.Date)
	}
	fmt.Fprintln(cmd.OutOrStdout())
	
	// Memory Information
	fmt.Fprintf(cmd.OutOrStdout(), "Memory Modules (%d detected):\n", len(inv.Memory))
	for i, module := range inv.Memory {
		fmt.Fprintf(cmd.OutOrStdout(), "  Module %d:\n", i+1)
		if module.DeviceLocator != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Location:     %s\n", module.DeviceLocator)
		}
		if module.BankLabel != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Bank:         %s\n", module.BankLabel)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "    Capacity:     %s\n", formatCapacity(module.Capacity))
		if module.Speed > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "    Speed:        %d MT/s\n", module.Speed)
		}
		if module.Manufacturer != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Manufacturer: %s\n", module.Manufacturer)
		}
		if module.PartNumber != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Part Number:  %s\n", module.PartNumber)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}
	
	return nil
}

// formatCapacity converts bytes to human-readable format
func formatCapacity(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}