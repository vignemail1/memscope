package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/spf13/cobra"
	
	"github.com/vignemail1/memscope/internal/collect/inventory"
	"github.com/vignemail1/memscope/internal/collect/runtime"
	"github.com/vignemail1/memscope/internal/model"
	"github.com/vignemail1/memscope/internal/snapshot"
)

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage system snapshots",
		Long:  "Create, list, compare, and manage system configuration snapshots",
	}
	
	cmd.AddCommand(
		newSnapshotCreateCmd(),
		newSnapshotListCmd(),
		newSnapshotCompareCmd(),
		newSnapshotDeleteCmd(),
	)
	
	return cmd
}

func newSnapshotCreateCmd() *cobra.Command {
	var (
		nameFlag        string
		descriptionFlag string
		outputDirFlag   string
		includeRuntime  bool
	)
	
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new system snapshot",
		Long:  "Capture current system configuration and save as a snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default snapshot name if not provided
			if nameFlag == "" {
				nameFlag = fmt.Sprintf("snapshot_%s", time.Now().Format("2006-01-02_15-04-05"))
			}
			
			// Default output directory
			if outputDirFlag == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				outputDirFlag = filepath.Join(homeDir, ".memscope", "snapshots")
			}
			
			fmt.Fprintln(cmd.OutOrStdout(), "Creating system snapshot...")
			
			// Create a simple snapshot using the legacy inventory system
			// TODO: Upgrade to use the full model collector system when it's fully implemented
			
			// Collect hardware inventory
			inventoryProvider := inventory.NewProvider()
			inv, err := inventoryProvider.CollectInventory()
			if err != nil {
				return fmt.Errorf("failed to collect hardware inventory: %w", err)
			}
			
			// Convert to model snapshot format using the existing ToSnapshot method
			snapshotPtr, err := inv.ToSnapshot()
			if err != nil {
				return fmt.Errorf("failed to convert inventory to snapshot: %w", err)
			}
			modelSnapshot := *snapshotPtr
			
			// Update snapshot metadata
			modelSnapshot.ID = nameFlag
			
			// Add runtime data if requested
			if includeRuntime {
				runtimeProvider := runtime.NewProvider()
				params, err := runtimeProvider.CollectMemoryParameters()
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "Warning: Failed to collect runtime parameters: %v\n", err)
				} else {
					// Add runtime observations
					addRuntimeObservations(&modelSnapshot, params)
				}
			}
			
			// Validate snapshot
			validator := snapshot.NewValidator()
			validationResult := validator.ValidateSnapshotDetailed(&modelSnapshot)
			if !validationResult.Valid {
				fmt.Fprintf(cmd.OutOrStderr(), "Warning: Snapshot validation failed:\n")
				for _, err := range validationResult.Errors {
					fmt.Fprintf(cmd.OutOrStderr(), "  - %s\n", err)
				}
			}
			if len(validationResult.Warnings) > 0 {
				fmt.Fprintf(cmd.OutOrStderr(), "Validation warnings:\n")
				for _, warning := range validationResult.Warnings {
					fmt.Fprintf(cmd.OutOrStderr(), "  - %s\n", warning)
				}
			}
			
			// Save snapshot
			manager := snapshot.NewManager(outputDirFlag)
			snapshotPath, err := manager.CreateSnapshot(&modelSnapshot, nameFlag, descriptionFlag)
			if err != nil {
				return fmt.Errorf("failed to create snapshot: %w", err)
			}
			
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot created successfully: %s\n", snapshotPath)
			return nil
		},
	}
	
	cmd.Flags().StringVar(&nameFlag, "name", "", "Snapshot name (default: auto-generated)")
	cmd.Flags().StringVar(&descriptionFlag, "description", "", "Snapshot description")
	cmd.Flags().StringVar(&outputDirFlag, "output-dir", "", "Output directory (default: ~/.memscope/snapshots)")
	cmd.Flags().BoolVar(&includeRuntime, "include-runtime", true, "Include runtime memory parameters")
	
	return cmd
}

func newSnapshotListCmd() *cobra.Command {
	var (
		formatFlag    string
		outputDirFlag string
	)
	
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available snapshots",
		Long:  "Display all stored snapshots with metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default output directory
			if outputDirFlag == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				outputDirFlag = filepath.Join(homeDir, ".memscope", "snapshots")
			}
			
			manager := snapshot.NewManager(outputDirFlag)
			snapshots, err := manager.ListSnapshots()
			if err != nil {
				return fmt.Errorf("failed to list snapshots: %w", err)
			}
			
			if len(snapshots) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
				return nil
			}
			
			// Display snapshots
			if formatFlag == "json" {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(snapshots)
			}
			
			// Text format
			fmt.Fprintf(cmd.OutOrStdout(), "Found %d snapshot(s):\n\n", len(snapshots))
			for _, snap := range snapshots {
				fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", snap.Name)
				fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", snap.Description)
				fmt.Fprintf(cmd.OutOrStdout(), "Created: %s\n", snap.Timestamp.Format("2006-01-02 15:04:05"))
				fmt.Fprintf(cmd.OutOrStdout(), "Size: %d bytes\n", snap.Size)
				fmt.Fprintf(cmd.OutOrStdout(), "Path: %s\n", snap.FilePath)
				if snap.Checksum != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "Checksum: %s\n", snap.Checksum)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "---")
			}
			
			return nil
		},
	}
	
	cmd.Flags().StringVar(&formatFlag, "format", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&outputDirFlag, "output-dir", "", "Snapshot directory (default: ~/.memscope/snapshots)")
	
	return cmd
}

func newSnapshotCompareCmd() *cobra.Command {
	var (
		formatFlag     string
		outputFlag     string
		ignoreTimestamps bool
		showOnlyChanges  bool
	)
	
	cmd := &cobra.Command{
		Use:   "compare <snapshot1> <snapshot2>",
		Short: "Compare two snapshots",
		Long:  "Compare two system snapshots and show differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshot1Path := args[0]
			snapshot2Path := args[1]
			
			// Load snapshots
			manager := snapshot.NewManager("")
			snapshot1, err := manager.LoadSnapshot(snapshot1Path)
			if err != nil {
				return fmt.Errorf("failed to load snapshot1: %w", err)
			}
			
			snapshot2, err := manager.LoadSnapshot(snapshot2Path)
			if err != nil {
				return fmt.Errorf("failed to load snapshot2: %w", err)
			}
			
			// Configure comparison options
			options := &snapshot.ComparisonOptions{
				IgnoreTimestamps:   ignoreTimestamps,
				ShowOnlyChanges:    showOnlyChanges,
				IncludeRuntimeData: true,
			}
			
			// Compare snapshots
			comparer := snapshot.NewComparerWithOptions(options)
			result, err := comparer.CompareSnapshots(snapshot1, snapshot2)
			if err != nil {
				return fmt.Errorf("comparison failed: %w", err)
			}
			
			// Output results
			if outputFlag != "" {
				file, err := os.Create(outputFlag)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()
				
				if formatFlag == "json" {
					encoder := json.NewEncoder(file)
					encoder.SetIndent("", "  ")
					return encoder.Encode(result)
				}
			}
			
			// Display results to stdout
			if formatFlag == "json" {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(result)
			}
			
			// Text format
			return displayComparisonResult(cmd, result)
		},
	}
	
	cmd.Flags().StringVar(&formatFlag, "format", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&outputFlag, "output", "", "Output file path")
	cmd.Flags().BoolVar(&ignoreTimestamps, "ignore-timestamps", false, "Ignore timestamp differences")
	cmd.Flags().BoolVar(&showOnlyChanges, "changes-only", false, "Show only items that changed")
	
	return cmd
}

func newSnapshotDeleteCmd() *cobra.Command {
	var (
		outputDirFlag string
		forceFlag     bool
	)
	
	cmd := &cobra.Command{
		Use:   "delete <snapshot-name>",
		Short: "Delete a snapshot",
		Long:  "Remove a snapshot and its metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotName := args[0]
			
			// Default output directory
			if outputDirFlag == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				outputDirFlag = filepath.Join(homeDir, ".memscope", "snapshots")
			}
			
			manager := snapshot.NewManager(outputDirFlag)
			
			// Confirm deletion unless forced
			if !forceFlag {
				fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to delete snapshot '%s'? (y/N): ", snapshotName)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
					fmt.Fprintln(cmd.OutOrStdout(), "Deletion cancelled.")
					return nil
				}
			}
			
			// Delete snapshot
			if err := manager.DeleteSnapshot(snapshotName); err != nil {
				return fmt.Errorf("failed to delete snapshot: %w", err)
			}
			
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot '%s' deleted successfully.\n", snapshotName)
			return nil
		},
	}
	
	cmd.Flags().StringVar(&outputDirFlag, "output-dir", "", "Snapshot directory (default: ~/.memscope/snapshots)")
	cmd.Flags().BoolVar(&forceFlag, "force", false, "Skip confirmation prompt")
	
	return cmd
}

func displayComparisonResult(cmd *cobra.Command, result *snapshot.ComparisonResult) error {
	fmt.Fprintln(cmd.OutOrStdout(), "=== Snapshot Comparison Results ===")
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot 1: %s (%s)\n", 
		result.Snapshot1.Name, 
		result.Snapshot1.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot 2: %s (%s)\n", 
		result.Snapshot2.Name, 
		result.Snapshot2.Timestamp.Format("2006-01-02 15:04:05"))
	
	if !result.HasChanges {
		fmt.Fprintln(cmd.OutOrStdout(), "\nNo differences found between snapshots.")
		return nil
	}
	
	// Display summary
	fmt.Fprintf(cmd.OutOrStdout(), "\n=== Summary ===\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Total changes: %d\n", result.Summary.TotalChanges)
	fmt.Fprintf(cmd.OutOrStdout(), "System changes: %d\n", result.Summary.SystemChanges)
	fmt.Fprintf(cmd.OutOrStdout(), "Memory changes: %d\n", result.Summary.MemoryChanges)
	fmt.Fprintf(cmd.OutOrStdout(), "Device changes: %d\n", result.Summary.DeviceChanges)
	fmt.Fprintf(cmd.OutOrStdout(), "Observation changes: %d\n", result.Summary.ObservationChanges)
	
	// Display changes by severity
	fmt.Fprintf(cmd.OutOrStdout(), "\nBy severity:\n")
	for severity, count := range result.Summary.BySeverity {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s: %d\n", severity, count)
	}
	
	// Display individual changes
	fmt.Fprintln(cmd.OutOrStdout(), "\n=== Changes ===")
	for i, change := range result.Changes {
		fmt.Fprintf(cmd.OutOrStdout(), "%d. [%s] %s\n", i+1, strings.ToUpper(change.Severity), change.Description)
		fmt.Fprintf(cmd.OutOrStdout(), "   Path: %s\n", change.Path)
		fmt.Fprintf(cmd.OutOrStdout(), "   Type: %s\n", change.Type)
		
		if change.OldValue != nil && change.NewValue != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "   Old: %v\n", change.OldValue)
			fmt.Fprintf(cmd.OutOrStdout(), "   New: %v\n", change.NewValue)
		} else if change.OldValue != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "   Removed: %v\n", change.OldValue)
		} else if change.NewValue != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "   Added: %v\n", change.NewValue)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "")
	}
	
	return nil
}



// addRuntimeObservations adds runtime parameter observations to the snapshot
func addRuntimeObservations(snapshot *model.Snapshot, params *runtime.MemoryParameters) {
	timestamp := time.Now()
	
	// Add timing observations
	for name, timing := range params.Timings {
		obs := model.Observation{
			DeviceID:   "memory-runtime",
			Scope:      "runtime",
			Parameter:  fmt.Sprintf("timing_%s", name),
			Value:      model.NewUnsigned(uint64(timing.Value)),
			Unit:       timing.Unit,
			Source:     model.Controller,
			Status:     model.Observed,
			CapturedAt: timestamp,
		}
		snapshot.Observations = append(snapshot.Observations, obs)
	}
	
	// Add voltage observations
	for name, voltage := range params.Voltages {
		voltageValue, err := model.NewDecimal(float64(voltage.Value))
		if err != nil {
			// Skip invalid voltage values
			continue
		}
		obs := model.Observation{
			DeviceID:   "memory-runtime",
			Scope:      "runtime",
			Parameter:  fmt.Sprintf("voltage_%s", name),
			Value:      voltageValue,
			Unit:       voltage.Unit,
			Source:     model.Controller,
			Status:     model.Observed,
			CapturedAt: timestamp,
		}
		snapshot.Observations = append(snapshot.Observations, obs)
	}
}