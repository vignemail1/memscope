package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	
	"github.com/vignemail1/memscope/internal/model"
	"github.com/vignemail1/memscope/internal/recommend"
)

// Exporter handles data export operations
type Exporter struct {
	options ExportOptions
}

// ExportOptions configures export behavior
type ExportOptions struct {
	PrettyJSON    bool   `json:"pretty_json"`
	CSVDelimiter  string `json:"csv_delimiter"`
	IncludeHeader bool   `json:"include_header"`
	TimeFormat    string `json:"time_format"`
}

// NewExporter creates a new exporter with default options
func NewExporter() *Exporter {
	return &Exporter{
		options: ExportOptions{
			PrettyJSON:    true,
			CSVDelimiter:  ",",
			IncludeHeader: true,
			TimeFormat:    "2006-01-02T15:04:05Z07:00",
		},
	}
}

// NewExporterWithOptions creates an exporter with custom options
func NewExporterWithOptions(options ExportOptions) *Exporter {
	return &Exporter{options: options}
}

// ExportJSON exports a snapshot to JSON format
func (e *Exporter) ExportJSON(snapshot *model.Snapshot, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	if e.options.PrettyJSON {
		encoder.SetIndent("", "  ")
	}
	
	return encoder.Encode(snapshot)
}

// ExportMemoryCSV exports memory analysis data to CSV format
func (e *Exporter) ExportMemoryCSV(memoryData *MemoryAnalysis, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	if e.options.CSVDelimiter != "," {
		csvWriter.Comma = rune(e.options.CSVDelimiter[0])
	}
	defer csvWriter.Flush()
	
	// Write header
	if e.options.IncludeHeader {
		header := []string{
			"Slot", "Manufacturer", "PartNumber", "SerialNumber", "Capacity", 
			"Speed", "Voltage", "MemoryType", "FormFactor",
			"CL", "TRCD", "TRP", "TRAS", "TRC", "TRFC",
		}
		if err := csvWriter.Write(header); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
	}
	
	// Write memory module data
	for _, module := range memoryData.Modules {
		record := []string{
			module.Slot,
			module.Manufacturer,
			module.PartNumber,
			module.SerialNumber,
			module.Capacity,
			strconv.FormatUint(uint64(module.Speed), 10),
			strconv.FormatFloat(float64(module.Voltage), 'f', 2, 32),
			module.MemoryType,
			module.FormFactor,
		}
		
		// Add timing values
		timingFields := []string{"CL", "TRCD", "TRP", "TRAS", "TRC", "TRFC"}
		for _, field := range timingFields {
			if value, exists := module.Timings[field]; exists {
				record = append(record, strconv.FormatUint(uint64(value), 10))
			} else {
				record = append(record, "")
			}
		}
		
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}
	
	return nil
}

// ExportRecommendationsCSV exports recommendations to CSV format
func (e *Exporter) ExportRecommendationsCSV(recommendations []*recommend.Recommendation, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	if e.options.CSVDelimiter != "," {
		csvWriter.Comma = rune(e.options.CSVDelimiter[0])
	}
	defer csvWriter.Flush()
	
	// Write header
	if e.options.IncludeHeader {
		header := []string{
			"RuleID", "RuleVersion", "Category", "Status", "EvidenceCount",
			"HighConfidenceEvidence", "WarningCount", "ValidationSteps", "Timestamp",
		}
		if err := csvWriter.Write(header); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
	}
	
	// Write recommendation data
	for _, rec := range recommendations {
		highConfidenceEvidence := 0
		for _, evidence := range rec.Evidence {
			if evidence.Confidence >= 0.8 {
				highConfidenceEvidence++
			}
		}
		
		validationSteps := 0
		if rec.ValidationPlan != nil {
			validationSteps = len(rec.ValidationPlan.Steps)
		}
		
		record := []string{
			rec.RuleID,
			rec.RuleVersion,
			rec.Category,
			rec.Status,
			strconv.Itoa(len(rec.Evidence)),
			strconv.Itoa(highConfidenceEvidence),
			strconv.Itoa(len(rec.Warnings)),
			strconv.Itoa(validationSteps),
			rec.Timestamp.Format(e.options.TimeFormat),
		}
		
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}
	
	return nil
}

// ExportToFile exports data to a file with automatic format detection
func (e *Exporter) ExportToFile(data interface{}, filename string, format string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	// Open file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()
	
	// Export based on format
	switch strings.ToLower(format) {
	case "json":
		// Handle different data types for JSON export
		switch typedData := data.(type) {
		case *model.Snapshot:
			return e.ExportJSON(typedData, file)
		case *MemoryAnalysis, []*recommend.Recommendation, *RuntimeData:
			// Use generic JSON encoder for other types
			encoder := json.NewEncoder(file)
			if e.options.PrettyJSON {
				encoder.SetIndent("", "  ")
			}
			return encoder.Encode(data)
		default:
			return fmt.Errorf("unsupported data type for JSON export: %T", data)
		}
	case "csv":
		switch typedData := data.(type) {
		case *MemoryAnalysis:
			return e.ExportMemoryCSV(typedData, file)
		case []*recommend.Recommendation:
			return e.ExportRecommendationsCSV(typedData, file)
		default:
			return fmt.Errorf("unsupported data type for CSV export: %T", data)
		}
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// ConvertSnapshotToMemoryAnalysis converts a snapshot to memory analysis format
func (e *Exporter) ConvertSnapshotToMemoryAnalysis(snapshot *model.Snapshot) *MemoryAnalysis {
	analysis := &MemoryAnalysis{
		Timestamp: snapshot.CollectionStartedAt,
		Modules:   []*MemoryModule{},
	}
	
	// Convert system info
	analysis.SystemInfo = &SystemInfo{
		OS:          snapshot.Platform.OS,
		Architecture: snapshot.Platform.Arch,
	}
	
	// Extract additional system info from observations
	for _, obs := range snapshot.Observations {
		if obs.Value != nil && obs.Value.Text != nil {
			switch obs.Parameter {
			case "name":
				if strings.Contains(obs.DeviceID, "cpu") {
					analysis.SystemInfo.CPU = *obs.Value.Text
				}
			case "model":
				if strings.Contains(obs.DeviceID, "system") {
					analysis.SystemInfo.Motherboard = *obs.Value.Text
				}
			case "version":
				if strings.Contains(obs.DeviceID, "bios") {
					analysis.SystemInfo.BIOSVersion = *obs.Value.Text
				}
			}
		}
	}
	
	// Convert memory devices
	for _, device := range snapshot.Devices {
		if device.Kind == "memory" {
			module := &MemoryModule{
				Slot:    device.ID,
				Timings: make(map[string]uint16),
			}
			
			// Extract memory module properties from observations
			for _, obs := range snapshot.Observations {
				if obs.DeviceID == device.ID && obs.Value != nil {
					switch obs.Parameter {
					case "manufacturer":
						if obs.Value.Text != nil {
							module.Manufacturer = *obs.Value.Text
						}
					case "part_number":
						if obs.Value.Text != nil {
							module.PartNumber = *obs.Value.Text
						}
					case "serial_number":
						if obs.Value.Text != nil {
							module.SerialNumber = *obs.Value.Text
						}
					case "capacity":
						if obs.Value.Unsigned != nil {
							// Convert bytes to GB for display
							capacityGB := *obs.Value.Unsigned / (1024 * 1024 * 1024)
							module.Capacity = fmt.Sprintf("%dGB", capacityGB)
						}
					case "speed":
						if obs.Value.Unsigned != nil {
							module.Speed = uint32(*obs.Value.Unsigned)
						}
					case "device_locator":
						if obs.Value.Text != nil {
							module.Slot = *obs.Value.Text
						}
					case "memory_type":
						if obs.Value.Text != nil {
							module.MemoryType = *obs.Value.Text
						}
					case "form_factor":
						if obs.Value.Text != nil {
							module.FormFactor = *obs.Value.Text
						}
					}
				}
			}
			
			// Extract SPD profiles if available
			for _, profile := range snapshot.Profiles {
				if profile.DeviceID == device.ID && profile.Type == "spd" {
					spdProfile := &SPDProfile{
						Name:     profile.ID,
						Type:     profile.Type,
						Timings:  make(map[string]uint16),
						Supported: true,
					}
					
					// Extract profile observations
					for _, obs := range profile.Observations {
						if obs.Value != nil && obs.Value.Text != nil {
							switch obs.Parameter {
							case "manufacturer":
								// Profile manufacturer might override device manufacturer
								if module.Manufacturer == "" {
									module.Manufacturer = *obs.Value.Text
								}
							}
						}
					}
					
					module.SPDProfiles = append(module.SPDProfiles, spdProfile)
				}
			}
			
			// Set default voltage if not specified
			if module.Voltage == 0 {
				module.Voltage = 1.2 // Default DDR4 voltage
			}
			
			analysis.Modules = append(analysis.Modules, module)
		}
	}
	
	return analysis
}

// GenerateExportSummary creates a summary of export results
func (e *Exporter) GenerateExportSummary(result *ExportResult) string {
	var summary strings.Builder
	
	fmt.Fprintf(&summary, "Export Summary:\n")
	fmt.Fprintf(&summary, "  Success: %v\n", result.Success)
	fmt.Fprintf(&summary, "  Output: %s\n", result.OutputPath)
	fmt.Fprintf(&summary, "  Format: %s\n", result.Format)
	fmt.Fprintf(&summary, "  Records: %d\n", result.RecordCount)
	fmt.Fprintf(&summary, "  File Size: %d bytes\n", result.FileSize)
	fmt.Fprintf(&summary, "  Duration: %v\n", result.Duration)
	
	if result.Error != "" {
		fmt.Fprintf(&summary, "  Error: %s\n", result.Error)
	}
	
	return summary.String()
}