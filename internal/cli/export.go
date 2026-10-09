package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	"github.com/spf13/cobra"
	
	"github.com/vignemail1/memscope/internal/collect/inventory"
	runtimeCollect "github.com/vignemail1/memscope/internal/collect/runtime"
	"github.com/vignemail1/memscope/internal/export"
	"github.com/vignemail1/memscope/internal/model"
	"github.com/vignemail1/memscope/internal/recommend"
)

func newExportCmd() *cobra.Command {
	var (
		formatFlag   string
		outputFlag   string
		typeFlag     string
		inputFlag    string
		prettyFlag   bool
		delimiterFlag string
	)
	
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export system data and analysis results",
		Long: `Export system snapshots, memory analysis, and recommendations to JSON or CSV formats

Supported export types:
  - snapshot: Complete system snapshot (JSON only)
  - memory: Memory analysis with timing and SPD data (JSON, CSV)  
  - recommendations: BIOS optimization recommendations (JSON, CSV)
  - runtime: Current memory parameters (JSON, CSV)

CSV format is optimized for spreadsheet analysis and includes:
  - Memory modules with timing parameters
  - Recommendation summaries with evidence metrics
  - Runtime parameters with current configuration

Examples:
  memscope export --type snapshot --output data.json
  memscope export --type memory --format csv --output memory.csv
  memscope export --type recommendations --format csv --output recommendations.csv --delimiter ";"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()
			
			// Validate flags
			if outputFlag == "" {
				return fmt.Errorf("output file path is required (use --output)")
			}
			
			// Auto-detect format from file extension if not specified
			if formatFlag == "" {
				ext := strings.ToLower(filepath.Ext(outputFlag))
				switch ext {
				case ".json":
					formatFlag = "json"
				case ".csv":
					formatFlag = "csv"
				default:
					return fmt.Errorf("cannot auto-detect format from extension %s, please specify --format", ext)
				}
			}
			
			// Validate format compatibility
			if formatFlag == "csv" && typeFlag == "snapshot" {
				return fmt.Errorf("CSV format is not supported for snapshot export, use JSON instead")
			}
			
			// Create exporter with options
			exportOptions := export.ExportOptions{
				PrettyJSON:    prettyFlag,
				CSVDelimiter:  delimiterFlag,
				IncludeHeader: true,
				TimeFormat:    "2006-01-02T15:04:05Z07:00",
			}
			exporter := export.NewExporterWithOptions(exportOptions)
			
			// Collect or load data based on type and input
			var exportResult *export.ExportResult
			var err error
			
			switch typeFlag {
			case "snapshot":
				exportResult, err = exportSnapshot(exporter, inputFlag, outputFlag, formatFlag)
			case "memory":
				exportResult, err = exportMemoryAnalysis(exporter, inputFlag, outputFlag, formatFlag)
			case "recommendations":
				exportResult, err = exportRecommendations(exporter, inputFlag, outputFlag, formatFlag)
			case "runtime":
				exportResult, err = exportRuntimeData(exporter, outputFlag, formatFlag)
			default:
				return fmt.Errorf("unsupported export type: %s (supported: snapshot, memory, recommendations, runtime)", typeFlag)
			}
			
			if err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
			
			exportResult.Duration = time.Since(start)
			
			// Display results
			fmt.Fprintln(cmd.OutOrStdout(), exporter.GenerateExportSummary(exportResult))
			
			return nil
		},
	}
	
	cmd.Flags().StringVar(&formatFlag, "format", "", "Export format (json, csv) - auto-detected from file extension. Note: CSV only supports memory, recommendations, and runtime types")
	cmd.Flags().StringVar(&outputFlag, "output", "", "Output file path (required)")
	cmd.Flags().StringVar(&typeFlag, "type", "snapshot", "Data type to export: snapshot (JSON only), memory (JSON/CSV), recommendations (JSON/CSV), runtime (JSON/CSV)")
	cmd.Flags().StringVar(&inputFlag, "input", "", "Input snapshot file for offline analysis")
	cmd.Flags().BoolVar(&prettyFlag, "pretty", true, "Pretty-print JSON output")
	cmd.Flags().StringVar(&delimiterFlag, "delimiter", ",", "CSV delimiter character")
	
	return cmd
}

// loadSnapshotFromSource loads a snapshot from file or live collection
func loadSnapshotFromSource(inputFile string) (*model.Snapshot, error) {
	if inputFile != "" {
		// Load from existing snapshot file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		
		var snapshot model.Snapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return nil, fmt.Errorf("failed to parse input snapshot: %w", err)
		}
		
		return &snapshot, nil
	}
	
	// Collect live data
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("live data collection requires Windows, use --input with snapshot file")
	}
	
	provider := inventory.NewProvider()
	inv, err := provider.CollectInventory()
	if err != nil {
		return nil, fmt.Errorf("failed to collect inventory: %w", err)
	}
	
	snapshot, err := inv.ToSnapshot()
	if err != nil {
		return nil, fmt.Errorf("failed to convert inventory to snapshot: %w", err)
	}
	
	return snapshot, nil
}

// exportSnapshot exports system snapshot data
func exportSnapshot(exporter *export.Exporter, inputFile, outputFile, format string) (*export.ExportResult, error) {
	snapshot, err := loadSnapshotFromSource(inputFile)
	if err != nil {
		return nil, err
	}
	
	// Export to file
	err = exporter.ExportToFile(snapshot, outputFile, format)
	if err != nil {
		return &export.ExportResult{Success: false, Error: err.Error()}, err
	}
	
	// Get file stats
	fileInfo, _ := os.Stat(outputFile)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	return &export.ExportResult{
		Success:     true,
		OutputPath:  outputFile,
		Format:      format,
		RecordCount: 1, // One snapshot
		FileSize:    fileSize,
		Timestamp:   time.Now(),
	}, nil
}

// exportMemoryAnalysis exports memory-specific analysis data
func exportMemoryAnalysis(exporter *export.Exporter, inputFile, outputFile, format string) (*export.ExportResult, error) {
	snapshot, err := loadSnapshotFromSource(inputFile)
	if err != nil {
		return nil, err
	}
	
	// Convert to memory analysis format
	memoryAnalysis := exporter.ConvertSnapshotToMemoryAnalysis(snapshot)
	
	// Export to file
	err = exporter.ExportToFile(memoryAnalysis, outputFile, format)
	if err != nil {
		return &export.ExportResult{Success: false, Error: err.Error()}, err
	}
	
	// Get file stats
	fileInfo, _ := os.Stat(outputFile)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	return &export.ExportResult{
		Success:     true,
		OutputPath:  outputFile,
		Format:      format,
		RecordCount: len(memoryAnalysis.Modules),
		FileSize:    fileSize,
		Timestamp:   time.Now(),
	}, nil
}

// exportRecommendations exports recommendation analysis results
func exportRecommendations(exporter *export.Exporter, inputFile, outputFile, format string) (*export.ExportResult, error) {
	snapshot, err := loadSnapshotFromSource(inputFile)
	if err != nil {
		return nil, err
	}
	
	// Generate recommendations
	engine := recommend.NewEngine()
	recommendations, err := engine.GenerateRecommendations(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}
	
	// Export to file
	err = exporter.ExportToFile(recommendations, outputFile, format)
	if err != nil {
		return &export.ExportResult{Success: false, Error: err.Error()}, err
	}
	
	// Get file stats
	fileInfo, _ := os.Stat(outputFile)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	return &export.ExportResult{
		Success:     true,
		OutputPath:  outputFile,
		Format:      format,
		RecordCount: len(recommendations),
		FileSize:    fileSize,
		Timestamp:   time.Now(),
	}, nil
}

// exportRuntimeData exports current runtime memory parameters
func exportRuntimeData(exporter *export.Exporter, outputFile, format string) (*export.ExportResult, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("runtime data collection requires Windows")
	}
	
	// Collect actual runtime parameters
	runtimeProvider := runtimeCollect.NewProvider()
	params, err := runtimeProvider.CollectMemoryParameters()
	if err != nil {
		return nil, fmt.Errorf("failed to collect runtime parameters: %w", err)
	}
	
	// Convert to export format with validation
	runtimeData := &export.RuntimeData{
		CurrentFrequency: params.CurrentFrequency,
		ConfiguredSpeed:  params.ConfiguredSpeed,
		Timings:         make(map[string]uint16),
		Voltages:        make(map[string]float32),
		Temperature:     make(map[string]float32),
	}
	
	// Set active profile if available
	if params.Profile != nil {
		runtimeData.ActiveProfile = params.Profile.Name
	}
	
	// Convert timings with validation
	for name, timing := range params.Timings {
		if timing.Value > 0 {
			runtimeData.Timings[name] = timing.Value
		}
	}
	
	// Convert voltages with validation
	for name, voltage := range params.Voltages {
		if voltage.Value > 0 {
			runtimeData.Voltages[name] = voltage.Value
		}
	}
	
	var exportErr error
	
	// For CSV export, create a comprehensive analysis structure
	if format == "csv" {
		memoryAnalysis := &export.MemoryAnalysis{
			Timestamp:   time.Now(),
			RuntimeData: runtimeData,
			SystemInfo: &export.SystemInfo{
				CPU: "Runtime Collection",
			},
			Modules: []*export.MemoryModule{
				{
					Slot:     "Runtime Parameters",
					Speed:    runtimeData.CurrentFrequency,
					Voltage:  getAverageVoltage(runtimeData.Voltages),
					Timings:  runtimeData.Timings,
				},
			},
		}
		
		exportErr = exporter.ExportToFile(memoryAnalysis, outputFile, format)
	} else {
		exportErr = exporter.ExportToFile(runtimeData, outputFile, format)
	}
	
	if exportErr != nil {
		return &export.ExportResult{Success: false, Error: exportErr.Error()}, exportErr
	}
	
	// Get file stats
	fileInfo, _ := os.Stat(outputFile)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	
	return &export.ExportResult{
		Success:     true,
		OutputPath:  outputFile,
		Format:      format,
		RecordCount: len(runtimeData.Timings) + len(runtimeData.Voltages),
		FileSize:    fileSize,
		Timestamp:   time.Now(),
	}, nil
}

// Helper function to calculate average voltage
func getAverageVoltage(voltages map[string]float32) float32 {
	if len(voltages) == 0 {
		return 0
	}
	
	var sum float32
	for _, voltage := range voltages {
		sum += voltage
	}
	return sum / float32(len(voltages))
}