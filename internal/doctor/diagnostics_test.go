package doctor

import (
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestSystemDiagnostics(t *testing.T) {
	doctor := NewDoctor()
	
	// Create sample system snapshot
	snapshot := &model.Snapshot{
		CollectionStartedAt: time.Now(),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "system-0",
				Kind: "system",
			},
			{
				ID:       "cpu-0",
				Kind:     "cpu",
				ParentID: "system-0",
			},
			{
				ID:       "memory-0",
				Kind:     "memory",
				ParentID: "system-0",
			},
		},
		Observations: []model.Observation{
			{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "name",
				Value:      &model.Value{Text: stringPtr("AMD Ryzen 5 3600")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "manufacturer",
				Value:      &model.Value{Text: stringPtr("Corsair")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "speed",
				Value:      &model.Value{Unsigned: uint64Ptr(3200)},
				Unit:       "MT/s",
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
		},
	}
	
	// Run diagnostics
	result, err := doctor.RunDiagnostics(snapshot)
	if err != nil {
		t.Fatalf("RunDiagnostics failed: %v", err)
	}
	
	// Validate diagnostic result structure
	if result == nil {
		t.Fatal("Diagnostic result should not be nil")
	}
	
	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Overall score should be 0-100, got %f", result.OverallScore)
	}
	
	if len(result.Checks) == 0 {
		t.Error("Should have at least one diagnostic check")
	}
	
	// Validate check structure
	for _, check := range result.Checks {
		if check.Name == "" {
			t.Error("Check should have a name")
		}
		if check.Status == "" {
			t.Error("Check should have a status")
		}
	}
}

func TestMemoryHealthCheck(t *testing.T) {
	doctor := NewDoctor()
	
	// Test with problematic memory configuration (high voltage)
	problemSnapshot := &model.Snapshot{
		CollectionStartedAt: time.Now(),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "memory-0",
				Kind: "memory",
			},
		},
		Observations: []model.Observation{
			{
				DeviceID:   "memory-0",
				Scope:      "spd",
				Parameter:  "voltage_level",
				Value:      &model.Value{Decimal: float64Ptr(1.5)}, // Too high voltage
				Source:     model.SPD,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
		},
	}
	
	result, err := doctor.RunDiagnostics(problemSnapshot)
	if err != nil {
		t.Fatalf("RunDiagnostics failed: %v", err)
	}
	
	// Should detect voltage issue
	foundVoltageIssue := false
	for _, check := range result.Checks {
		if check.Category == "memory" && (check.Status == "warning" || check.Status == "error") {
			foundVoltageIssue = true
			break
		}
	}
	
	if !foundVoltageIssue {
		t.Error("Should detect high voltage issue")
	}
}

func TestPerformanceAnalysis(t *testing.T) {
	doctor := NewDoctor()
	
	// Test with suboptimal configuration (high-end CPU with slow memory)
	snapshot := &model.Snapshot{
		CollectionStartedAt: time.Now(),
		Devices: []model.Device{
			{
				ID:   "cpu-0",
				Kind: "cpu",
			},
			{
				ID:   "memory-0",
				Kind: "memory",
			},
		},
		Observations: []model.Observation{
			{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "name",
				Value:      &model.Value{Text: stringPtr("AMD Ryzen 7 3700X")}, // High-end CPU
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "speed",
				Value:      &model.Value{Unsigned: uint64Ptr(2133)}, // Slow memory for this CPU
				Unit:       "MT/s",
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
		},
	}
	
	result, err := doctor.RunDiagnostics(snapshot)
	if err != nil {
		t.Fatalf("RunDiagnostics failed: %v", err)
	}
	
	// Should identify performance bottleneck
	foundPerformanceIssue := false
	for _, check := range result.Checks {
		if check.Category == "performance" {
			foundPerformanceIssue = true
			break
		}
	}
	
	if !foundPerformanceIssue {
		t.Error("Should identify performance optimization opportunity")
	}
}

// Helper functions for test data
func stringPtr(s string) *string {
	return &s
}

func uint64Ptr(u uint64) *uint64 {
	return &u
}

func float64Ptr(f float64) *float64 {
	return &f
}