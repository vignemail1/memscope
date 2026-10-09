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
		Long:  "Export system snapshots, memory analysis, and recommendations to JSON or CSV formats",
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
	
	cmd.Flags().StringVar(&formatFlag, "format", "", "Export format (json, csv) - auto-detected from file extension")
	cmd.Flags().StringVar(&outputFlag, "output", "", "Output file path (required)")
	cmd.Flags().StringVar(&typeFlag, "type", "snapshot", "Data type to export (snapshot, memory, recommendations, runtime)")
	cmd.Flags().StringVar(&inputFlag, "input", "", "Input snapshot file for offline analysis")
	cmd.Flags().BoolVar(&prettyFlag, "pretty", true, "Pretty-print JSON output")
	cmd.Flags().StringVar(&delimiterFlag, "delimiter", ",", "CSV delimiter character")
	
	return cmd
}

// exportSnapshot exports system snapshot data
func exportSnapshot(exporter *export.Exporter, inputFile, outputFile, format string) (*export.ExportResult, error) {
	var snapshot *model.Snapshot
	var err error
	
	if inputFile != "" {
		// Load from existing snapshot file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		
		snapshot = &model.Snapshot{}
		if err := json.Unmarshal(data, snapshot); err != nil {
			return nil, fmt.Errorf("failed to parse input snapshot: %w", err)
		}
	} else {
		// Collect live data
		if runtime.GOOS != "windows" {
			return nil, fmt.Errorf("live data collection requires Windows, use --input with snapshot file")
		}
		
		provider := inventory.NewProvider()
		inv, err := provider.CollectInventory()
		if err != nil {
			return nil, fmt.Errorf("failed to collect inventory: %w", err)
		}
		
		snapshot, err = inv.ToSnapshot()
		if err != nil {
			return nil, fmt.Errorf("failed to convert inventory to snapshot: %w", err)
		}
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
	var snapshot *model.Snapshot
	var err error
	
	if inputFile != "" {
		// Load from existing snapshot file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		
		snapshot = &model.Snapshot{}
		if err := json.Unmarshal(data, snapshot); err != nil {
			return nil, fmt.Errorf("failed to parse input snapshot: %w", err)
		}
	} else {
		// Collect live data
		if runtime.GOOS != "windows" {
			return nil, fmt.Errorf("live data collection requires Windows, use --input with snapshot file")
		}
		
		provider := inventory.NewProvider()
		inv, err := provider.CollectInventory()
		if err != nil {
			return nil, fmt.Errorf("failed to collect inventory: %w", err)
		}
		
		snapshot, err = inv.ToSnapshot()
		if err != nil {
			return nil, fmt.Errorf("failed to convert inventory to snapshot: %w", err)
		}
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
	var snapshot *model.Snapshot
	var err error
	
	if inputFile != "" {
		// Load from existing snapshot file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		
		snapshot = &model.Snapshot{}
		if err := json.Unmarshal(data, snapshot); err != nil {
			return nil, fmt.Errorf("failed to parse input snapshot: %w", err)
		}
	} else {
		// Collect live data
		if runtime.GOOS != "windows" {
			return nil, fmt.Errorf("live data collection requires Windows, use --input with snapshot file")
		}
		
		provider := inventory.NewProvider()
		inv, err := provider.CollectInventory()
		if err != nil {
			return nil, fmt.Errorf("failed to collect inventory: %w", err)
		}
		
		snapshot, err = inv.ToSnapshot()
		if err != nil {
			return nil, fmt.Errorf("failed to convert inventory to snapshot: %w", err)
		}
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
	
	// For now, create a basic runtime data structure
	// In a real implementation, this would use the runtime provider
	runtimeData := &export.RuntimeData{
		CurrentFrequency: 2666,
		ConfiguredSpeed:  2666,
		Timings:         make(map[string]uint16),
		Voltages:        make(map[string]float32),
	}
	
	// Add some sample timing data
	runtimeData.Timings["CL"] = 19
	runtimeData.Timings["TRCD"] = 19
	runtimeData.Timings["TRP"] = 19
	runtimeData.Timings["TRAS"] = 43
	
	runtimeData.Voltages["VDIMM"] = 1.2
	runtimeData.Voltages["VCCSA"] = 1.0
	
	// For CSV export, create a simple analysis structure
	if format == "csv" {
		memoryAnalysis := &export.MemoryAnalysis{
			Timestamp:   time.Now(),
			RuntimeData: runtimeData,
			Modules: []*export.MemoryModule{
				{
					Slot:     "Runtime",
					Speed:    runtimeData.CurrentFrequency,
					Timings:  runtimeData.Timings,
				},
			},
		}
		
		err := exporter.ExportToFile(memoryAnalysis, outputFile, format)
		if err != nil {
			return &export.ExportResult{Success: false, Error: err.Error()}, err
		}
	} else {
		err := exporter.ExportToFile(runtimeData, outputFile, format)
		if err != nil {
			return &export.ExportResult{Success: false, Error: err.Error()}, err
		}
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
		RecordCount: 1,
		FileSize:    fileSize,
		Timestamp:   time.Now(),
	}, nil
}