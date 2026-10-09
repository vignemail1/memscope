package export

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestJSONExporter(t *testing.T) {
	exporter := NewExporter()
	
	// Create sample snapshot
	snapshot := &model.Snapshot{
		SchemaVersion:        "1.0",
		ToolVersion:         "0.1.0",
		ID:                  "test-snapshot",
		CollectionStartedAt:  time.Now(),
		CollectionFinishedAt: time.Now().Add(time.Second),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:       "memory-0",
				Kind:     "memory",
				ParentID: "",
			},
		},
	}
	
	var buffer bytes.Buffer
	err := exporter.ExportJSON(snapshot, &buffer)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}
	
	// Validate JSON structure
	var exported model.Snapshot
	if err := json.Unmarshal(buffer.Bytes(), &exported); err != nil {
		t.Fatalf("Exported JSON is invalid: %v", err)
	}
	
	if exported.Platform.OS != snapshot.Platform.OS {
		t.Errorf("OS mismatch: got %s, want %s", exported.Platform.OS, snapshot.Platform.OS)
	}
}

func TestCSVExporter(t *testing.T) {
	exporter := NewExporter()
	
	// Create sample memory analysis data
	memoryData := &MemoryAnalysis{
		Modules: []*MemoryModule{
			{
				Slot:         "DIMM_A1",
				Manufacturer: "Corsair",
				PartNumber:   "CMK16GX4M2B3200C16",
				Capacity:     "8GB",
				Speed:        3200,
				Voltage:      1.35,
				Timings: map[string]uint16{
					"CL":   16,
					"TRCD": 18,
					"TRP":  18,
					"TRAS": 38,
				},
			},
		},
	}
	
	var buffer bytes.Buffer
	err := exporter.ExportMemoryCSV(memoryData, &buffer)
	if err != nil {
		t.Fatalf("ExportMemoryCSV failed: %v", err)
	}
	
	// Validate CSV structure
	csvContent := buffer.String()
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")
	
	if len(lines) < 2 {
		t.Errorf("CSV should have header + at least 1 data row, got %d lines", len(lines))
	}
	
	// Check header
	expectedHeaders := []string{"Slot", "Manufacturer", "PartNumber", "Capacity", "Speed", "Voltage"}
	headerLine := lines[0]
	for _, header := range expectedHeaders {
		if !strings.Contains(headerLine, header) {
			t.Errorf("CSV header missing %s", header)
		}
	}
}

func TestFileExport(t *testing.T) {
	exporter := NewExporter()
	
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "memscope_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	snapshot := &model.Snapshot{
		SchemaVersion:        "1.0",
		ToolVersion:         "0.1.0",
		ID:                  "test-snapshot",
		CollectionStartedAt:  time.Now(),
		CollectionFinishedAt: time.Now().Add(time.Second),
		Platform: model.Platform{OS: "windows", Arch: "amd64"},
	}
	
	// Test JSON file export
	jsonFile := filepath.Join(tempDir, "test.json")
	err = exporter.ExportToFile(snapshot, jsonFile, "json")
	if err != nil {
		t.Fatalf("ExportToFile JSON failed: %v", err)
	}
	
	// Verify file exists and contains valid JSON
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read exported JSON: %v", err)
	}
	
	var exported model.Snapshot
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("Exported JSON file is invalid: %v", err)
	}
}

func TestExporterWithCustomOptions(t *testing.T) {
	options := ExportOptions{
		PrettyJSON:    false,
		CSVDelimiter:  ";",
		IncludeHeader: false,
		TimeFormat:    "2006-01-02",
	}
	exporter := NewExporterWithOptions(options)
	
	// Test CSV with custom delimiter and no header
	memoryData := &MemoryAnalysis{
		Modules: []*MemoryModule{
			{
				Slot:         "DIMM_A1",
				Manufacturer: "Test",
				PartNumber:   "TEST123",
				Capacity:     "8GB",
				Speed:        2400,
				Voltage:      1.2,
				Timings:      map[string]uint16{"CL": 14},
			},
		},
	}
	
	var buffer bytes.Buffer
	err := exporter.ExportMemoryCSV(memoryData, &buffer)
	if err != nil {
		t.Fatalf("ExportMemoryCSV failed: %v", err)
	}
	
	csvContent := buffer.String()
	if strings.Contains(csvContent, ",") {
		t.Error("CSV should use semicolon delimiter, found comma")
	}
	if strings.Contains(csvContent, ";") {
		t.Log("Custom delimiter working correctly")
	}
}

func TestConvertSnapshotToMemoryAnalysis(t *testing.T) {
	exporter := NewExporter()
	
	snapshot := &model.Snapshot{
		CollectionStartedAt: time.Now(),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{ID: "memory-0", Kind: "memory"},
			{ID: "cpu-0", Kind: "cpu"}, // Should be filtered out
		},
	}
	
	analysis := exporter.ConvertSnapshotToMemoryAnalysis(snapshot)
	
	if analysis.SystemInfo.OS != "windows" {
		t.Errorf("Expected OS windows, got %s", analysis.SystemInfo.OS)
	}
	
	if len(analysis.Modules) != 1 {
		t.Errorf("Expected 1 memory module, got %d", len(analysis.Modules))
	}
	
	if analysis.Modules[0].Slot != "memory-0" {
		t.Errorf("Expected slot memory-0, got %s", analysis.Modules[0].Slot)
	}
}