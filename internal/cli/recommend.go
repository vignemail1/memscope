package cli

import (
	"encoding/json"
	"fmt"
	"runtime"
	
	"github.com/spf13/cobra"
	
	"github.com/vignemail1/memscope/internal/collect/inventory"
	"github.com/vignemail1/memscope/internal/model"
	"github.com/vignemail1/memscope/internal/recommend"
)

func newRecommendCmd() *cobra.Command {
	var (
		formatFlag       string
		includeExperimental bool
		snapshotFile     string
	)
	
	cmd := &cobra.Command{
		Use:   "recommend",
		Short: "Generate BIOS recommendations based on memory analysis",
		Long:  "Analyzes memory configuration and provides evidence-based BIOS recommendations",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Analyzing system for memory optimization recommendations...")
			
			var snapshot *model.Snapshot
			var err error
			
			// Load snapshot from file or collect live data
			if snapshotFile != "" {
				snapshot, err = loadSnapshotFromFile(snapshotFile)
				if err != nil {
					return fmt.Errorf("failed to load snapshot from file: %w", err)
				}
			} else {
				// Collect live system data
				if runtime.GOOS != "windows" {
					fmt.Fprintln(cmd.OutOrStdout(), "Live analysis requires Windows. Use --input with snapshot file for offline analysis.")
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
			
			// Generate recommendations
			engine := recommend.NewEngine()
			recommendations, err := engine.GenerateRecommendations(snapshot)
			if err != nil {
				return fmt.Errorf("failed to generate recommendations: %w", err)
			}
			
			// Display results
			return displayRecommendations(cmd, recommendations, formatFlag)
		},
	}
	
	cmd.Flags().StringVar(&formatFlag, "format", "text", "Output format (text, json)")
	cmd.Flags().BoolVar(&includeExperimental, "experimental", false, "Include experimental recommendations")
	cmd.Flags().StringVar(&snapshotFile, "input", "", "Use snapshot file for offline analysis")
	
	return cmd
}

func displayRecommendations(cmd *cobra.Command, recommendations []*recommend.Recommendation, format string) error {
	if format == "json" {
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(recommendations)
	}
	
	// Text format output
	fmt.Fprintln(cmd.OutOrStdout(), "\n=== BIOS Memory Recommendations ===")
	
	for i, rec := range recommendations {
		fmt.Fprintf(cmd.OutOrStdout(), "\n--- Recommendation %d ---\n", i+1)
		fmt.Fprintf(cmd.OutOrStdout(), "Rule: %s (v%s)\n", rec.RuleID, rec.RuleVersion)
		fmt.Fprintf(cmd.OutOrStdout(), "Category: %s\n", rec.Category)
		fmt.Fprintf(cmd.OutOrStdout(), "Status: %s\n", rec.Status)
		
		if len(rec.Evidence) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "\nEvidence:")
			for _, evidence := range rec.Evidence {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%.1f%% confidence): %s\n", 
					evidence.Type, evidence.Confidence*100, evidence.Description)
			}
		}
		
		if len(rec.ProposedSettings) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "\nProposed Settings:")
			for key, value := range rec.ProposedSettings {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s: %v\n", key, value)
			}
		}
		
		if len(rec.Warnings) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "\nWarnings:")
			for _, warning := range rec.Warnings {
				fmt.Fprintf(cmd.OutOrStdout(), "  ! %s\n", warning)
			}
		}
		
		if len(rec.MissingData) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "\nMissing Data:")
			for _, missing := range rec.MissingData {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", missing)
			}
		}
		
		if rec.ValidationPlan != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "\nValidation Plan:")
			fmt.Fprintf(cmd.OutOrStdout(), "  Duration: %s\n", rec.ValidationPlan.Duration)
			if len(rec.ValidationPlan.Steps) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  Steps:")
				for _, step := range rec.ValidationPlan.Steps {
					required := ""
					if step.Required {
						required = " (required)"
					}
					fmt.Fprintf(cmd.OutOrStdout(), "    %d. %s%s\n", step.Order, step.Action, required)
				}
			}
			if len(rec.ValidationPlan.Tools) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "  Recommended Tools:")
				for _, tool := range rec.ValidationPlan.Tools {
					fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", tool)
				}
			}
		}
	}
	
	fmt.Fprintln(cmd.OutOrStdout(), "\nIMPORTANT: Test all changes thoroughly. Keep BIOS recovery procedures ready.")
	fmt.Fprintln(cmd.OutOrStdout(), "No guarantees provided - recommendations are suggestions for testing only.")
	
	return nil
}

func loadSnapshotFromFile(filename string) (*model.Snapshot, error) {
	// This would be implemented with actual file loading logic
	return nil, fmt.Errorf("snapshot file loading not implemented in this task")
}