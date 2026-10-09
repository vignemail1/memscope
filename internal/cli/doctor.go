package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	
	"github.com/spf13/cobra"
	
	"github.com/vignemail1/memscope/internal/collect/inventory"
	"github.com/vignemail1/memscope/internal/doctor"
	"github.com/vignemail1/memscope/internal/model"
)

func newDoctorCmd() *cobra.Command {
	var (
		formatFlag         string
		outputFlag         string
		inputFlag          string
		skipPerformance    bool
		skipCompatibility  bool
		categoryFilter     []string
		minSeverity        string
		exportReport       bool
	)
	
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run comprehensive system diagnostics",
		Long:  "Analyze system configuration, detect issues, and provide optimization recommendations",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Running system diagnostics...")
			
			var snapshot *model.Snapshot
			var err error
			
			// Load snapshot from file or collect live data
			if inputFlag != "" {
				snapshot, err = loadSnapshotFromFile(inputFlag)
				if err != nil {
					return fmt.Errorf("failed to load snapshot: %w", err)
				}
			} else {
				// Collect live data
				if runtime.GOOS != "windows" {
					fmt.Fprintln(cmd.OutOrStdout(), "Live diagnostics require Windows. Use --input with snapshot file for offline analysis.")
					return nil
				}
				
				provider := inventory.NewProvider()
				inv, err := provider.CollectInventory()
				if err != nil {
					return fmt.Errorf("failed to collect system inventory: %w", err)
				}
				
				snapshot, err = inv.ToSnapshot()
				if err != nil {
					return fmt.Errorf("failed to convert inventory to snapshot: %w", err)
				}
			}
			
			// Configure diagnostic options
			options := &doctor.DiagnosticOptions{
				IncludePerformanceAnalysis: !skipPerformance,
				IncludeCompatibilityCheck:  !skipCompatibility,
				CategoryFilter:            categoryFilter,
				MinSeverity:               minSeverity,
			}
			
			// Run diagnostics
			doc := doctor.NewDoctorWithOptions(options)
			result, err := doc.RunDiagnostics(snapshot)
			if err != nil {
				return fmt.Errorf("diagnostic analysis failed: %w", err)
			}
			
			// Output results
			if outputFlag != "" {
				err = outputDiagnosticResult(result, outputFlag, formatFlag)
				if err != nil {
					return fmt.Errorf("failed to save diagnostic report: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Diagnostic report saved to: %s\n", outputFlag)
			}
			
			// Display results to stdout
			return displayDiagnosticResult(cmd, result, formatFlag)
		},
	}
	
	cmd.Flags().StringVar(&formatFlag, "format", "text", "Output format (text, json)")
	cmd.Flags().StringVar(&outputFlag, "output", "", "Save diagnostic report to file")
	cmd.Flags().StringVar(&inputFlag, "input", "", "Input snapshot file for offline analysis")
	cmd.Flags().BoolVar(&skipPerformance, "skip-performance", false, "Skip performance analysis")
	cmd.Flags().BoolVar(&skipCompatibility, "skip-compatibility", false, "Skip compatibility checks")
	cmd.Flags().StringSliceVar(&categoryFilter, "categories", nil, "Filter by categories (system,memory,performance,compatibility)")
	cmd.Flags().StringVar(&minSeverity, "min-severity", "low", "Minimum severity level (low,medium,high,critical)")
	cmd.Flags().BoolVar(&exportReport, "export", false, "Export detailed report for technical support")
	
	return cmd
}

func displayDiagnosticResult(cmd *cobra.Command, result *doctor.DiagnosticResult, format string) error {
	if format == "json" {
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	
	// Text format display
	fmt.Fprintln(cmd.OutOrStdout(), "\n=== System Diagnostic Report ===")
	fmt.Fprintf(cmd.OutOrStdout(), "Overall Health: %s (Score: %.1f/100)\n", 
		result.HealthStatus, result.OverallScore)
	fmt.Fprintf(cmd.OutOrStdout(), "Generated: %s\n", 
		result.Timestamp.Format("2006-01-02 15:04:05"))
	
	// System information
	if result.SystemInfo != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "\n=== System Information ===")
		if result.SystemInfo.CPU != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "CPU: %s\n", result.SystemInfo.CPU)
		}
		if result.SystemInfo.Motherboard != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Motherboard: %s\n", result.SystemInfo.Motherboard)
		}
		if result.SystemInfo.BIOSVersion != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "BIOS Version: %s\n", result.SystemInfo.BIOSVersion)
		}
		if result.SystemInfo.TotalCapacity != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Memory: %s (%d modules)\n", 
				result.SystemInfo.TotalCapacity, result.SystemInfo.PopulatedSlots)
		}
		if result.SystemInfo.Platform != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Platform: %s\n", result.SystemInfo.Platform)
		}
	}
	
	// Summary statistics
	if result.Summary != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "\n=== Summary ===")
		fmt.Fprintf(cmd.OutOrStdout(), "Total Checks: %d\n", result.Summary.TotalChecks)
		fmt.Fprintf(cmd.OutOrStdout(), "Passed: %d, Warnings: %d, Errors: %d\n",
			result.Summary.PassedChecks, result.Summary.WarningChecks, result.Summary.ErrorChecks)
		
		if result.Summary.CriticalIssues > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "⚠️  Critical Issues: %d\n", result.Summary.CriticalIssues)
		}
		
		// Category scores
		fmt.Fprintln(cmd.OutOrStdout(), "\nCategory Scores:")
		fmt.Fprintf(cmd.OutOrStdout(), "  Memory: %.1f/100\n", result.Summary.MemoryScore)
		fmt.Fprintf(cmd.OutOrStdout(), "  Performance: %.1f/100\n", result.Summary.PerformanceScore)
		fmt.Fprintf(cmd.OutOrStdout(), "  Compatibility: %.1f/100\n", result.Summary.CompatibilityScore)
	}
	
	// Detailed checks
	if len(result.Checks) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\n=== Detailed Checks ===")
		
		categories := map[string][]*doctor.DiagnosticCheck{}
		for _, check := range result.Checks {
			categories[check.Category] = append(categories[check.Category], check)
		}
		
		for category, checks := range categories {
			fmt.Fprintf(cmd.OutOrStdout(), "\n--- %s ---\n", strings.ToTitle(category))
			for _, check := range checks {
				icon := getStatusIcon(check.Status)
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s: %s\n", icon, check.Name, check.Message)
				
				if check.Details != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "   Details: %s\n", check.Details)
				}
				if check.Suggestion != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "   Suggestion: %s\n", check.Suggestion)
				}
			}
		}
	}
	
	// Recommendations
	if len(result.Recommendations) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\n=== Recommendations ===")
		for i, rec := range result.Recommendations {
			priority := strings.ToUpper(rec.Priority)
			fmt.Fprintf(cmd.OutOrStdout(), "%d. [%s] %s\n", i+1, priority, rec.Title)
			fmt.Fprintf(cmd.OutOrStdout(), "   %s\n", rec.Description)
			if len(rec.Actions) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "   Actions: %s\n", strings.Join(rec.Actions, "; "))
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}
	
	// Final summary
	if result.OverallScore < 75 {
		fmt.Fprintln(cmd.OutOrStdout(), "💡 Consider addressing the warnings and errors above to improve system health.")
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "✅ System appears to be in good health!")
	}
	
	return nil
}

func outputDiagnosticResult(result *doctor.DiagnosticResult, filename, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()
	
	if format == "json" || strings.HasSuffix(filename, ".json") {
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	
	// Text format output to file would go here
	return fmt.Errorf("text format file output not implemented")
}

func getStatusIcon(status string) string {
	switch status {
	case "pass":
		return "✅"
	case "warning":
		return "⚠️ "
	case "error":
		return "❌"
	case "info":
		return "ℹ️ "
	default:
		return "❓"
	}
}

