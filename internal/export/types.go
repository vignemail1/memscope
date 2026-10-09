package export

import (
	"time"
	
	"github.com/vignemail1/memscope/internal/recommend"
)

// ExportRequest represents a request to export data
type ExportRequest struct {
	Type       string                 `json:"type"`        // snapshot, memory, recommendations, runtime
	Format     string                 `json:"format"`      // json, csv
	OutputPath string                 `json:"output_path"`
	Options    map[string]interface{} `json:"options,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

// ExportResult represents the result of an export operation
type ExportResult struct {
	Success     bool      `json:"success"`
	OutputPath  string    `json:"output_path"`
	Format      string    `json:"format"`
	RecordCount int       `json:"record_count"`
	FileSize    int64     `json:"file_size_bytes"`
	Duration    time.Duration `json:"duration"`
	Error       string    `json:"error,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// MemoryAnalysis represents structured memory data for export
type MemoryAnalysis struct {
	Timestamp    time.Time       `json:"timestamp"`
	SystemInfo   *SystemInfo     `json:"system_info"`
	Modules      []*MemoryModule `json:"modules"`
	RuntimeData  *RuntimeData    `json:"runtime_data,omitempty"`
}

// SystemInfo represents system information for export
type SystemInfo struct {
	CPU         string `json:"cpu"`
	Motherboard string `json:"motherboard"`
	BIOSVersion string `json:"bios_version"`
	OS          string `json:"operating_system"`
	Architecture string `json:"architecture"`
}

// MemoryModule represents a single memory module for export
type MemoryModule struct {
	Slot         string             `json:"slot"`
	Manufacturer string             `json:"manufacturer"`
	PartNumber   string             `json:"part_number"`
	SerialNumber string             `json:"serial_number,omitempty"`
	Capacity     string             `json:"capacity"`
	Speed        uint32             `json:"speed_mhz"`
	Voltage      float32            `json:"voltage_v"`
	MemoryType   string             `json:"memory_type"`
	FormFactor   string             `json:"form_factor"`
	Timings      map[string]uint16  `json:"timings"`
	SPDProfiles  []*SPDProfile      `json:"spd_profiles,omitempty"`
}

// SPDProfile represents an SPD profile for export
type SPDProfile struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Frequency   uint32             `json:"frequency_mhz"`
	Voltage     float32            `json:"voltage_v"`
	Timings     map[string]uint16  `json:"timings"`
	Supported   bool               `json:"supported"`
}

// RuntimeData represents current runtime memory parameters
type RuntimeData struct {
	CurrentFrequency uint32                    `json:"current_frequency_mhz"`
	ConfiguredSpeed  uint32                    `json:"configured_speed_mhz"`
	ActiveProfile    string                    `json:"active_profile,omitempty"`
	Timings          map[string]uint16         `json:"current_timings"`
	Voltages         map[string]float32        `json:"current_voltages"`
	Temperature      map[string]float32        `json:"temperatures,omitempty"`
}

// RecommendationExport represents recommendation data for export
type RecommendationExport struct {
	Timestamp       time.Time                    `json:"timestamp"`
	SystemInfo      *SystemInfo                  `json:"system_info"`
	Recommendations []*recommend.Recommendation  `json:"recommendations"`
	Summary         *RecommendationSummary       `json:"summary"`
}

// RecommendationSummary provides aggregated recommendation statistics
type RecommendationSummary struct {
	TotalRecommendations int                    `json:"total_recommendations"`
	ByCategory          map[string]int         `json:"by_category"`
	ByStatus            map[string]int         `json:"by_status"`
	HighConfidence      int                    `json:"high_confidence_count"`
	RequiresValidation  int                    `json:"requires_validation_count"`
	SafetyWarnings      int                    `json:"safety_warnings_count"`
}