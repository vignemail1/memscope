//go:build windows

package windows

import (
	"testing"
)

// Test the CalculateTypicalTimings function
func TestCalculateTypicalTimings(t *testing.T) {
	tests := []struct {
		name     string
		speed    uint32
		expected map[string]MemoryTiming
	}{
		{
			name:  "DDR4-3200 should have typical XMP timings",
			speed: 3200,
			expected: map[string]MemoryTiming{
				"CL":   {Value: 16, Unit: "clocks"},
				"TRCD": {Value: 18, Unit: "clocks"},
				"TRP":  {Value: 18, Unit: "clocks"},
				"TRAS": {Value: 38, Unit: "clocks"},
				"TRC":  {Value: 56, Unit: "clocks"},
			},
		},
		{
			name:  "DDR4-2666 should have moderate timings",
			speed: 2666,
			expected: map[string]MemoryTiming{
				"CL":   {Value: 15, Unit: "clocks"},
				"TRCD": {Value: 17, Unit: "clocks"},
				"TRP":  {Value: 17, Unit: "clocks"},
				"TRAS": {Value: 35, Unit: "clocks"},
				"TRC":  {Value: 52, Unit: "clocks"},
			},
		},
		{
			name:  "DDR4-2133 should have JEDEC standard timings",
			speed: 2133,
			expected: map[string]MemoryTiming{
				"CL":   {Value: 15, Unit: "clocks"},
				"TRCD": {Value: 15, Unit: "clocks"},
				"TRP":  {Value: 15, Unit: "clocks"},
				"TRAS": {Value: 35, Unit: "clocks"},
				"TRC":  {Value: 50, Unit: "clocks"},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timings := CalculateTypicalTimings(tt.speed)
			
			for expectedName, expectedTiming := range tt.expected {
				actualTiming, exists := timings[expectedName]
				if !exists {
					t.Errorf("Missing timing %s", expectedName)
					continue
				}
				
				if actualTiming.Value != expectedTiming.Value {
					t.Errorf("Timing %s: expected value %d, got %d", 
						expectedName, expectedTiming.Value, actualTiming.Value)
				}
				
				if actualTiming.Unit != expectedTiming.Unit {
					t.Errorf("Timing %s: expected unit %s, got %s", 
						expectedName, expectedTiming.Unit, actualTiming.Unit)
				}
			}
		})
	}
}

// Test that getDefaultVoltages returns reasonable values
func TestGetDefaultVoltages(t *testing.T) {
	voltages := getDefaultVoltages()
	
	// Should have at least VDIMM
	vdimm, exists := voltages["VDIMM"]
	if !exists {
		t.Error("Should have VDIMM voltage")
	}
	
	// VDIMM should be in reasonable range
	if vdimm.Value < 1.0 || vdimm.Value > 2.0 {
		t.Errorf("VDIMM voltage %f is out of reasonable range", vdimm.Value)
	}
	
	// Should have VTT as well
	vtt, exists := voltages["VTT"]
	if !exists {
		t.Error("Should have VTT voltage")
	}
	
	// VTT should be half of VDIMM approximately
	if vtt.Value < 0.5 || vtt.Value > 1.0 {
		t.Errorf("VTT voltage %f is out of reasonable range", vtt.Value)
	}
}

// Test enhanced profile detection
func TestDetectActiveProfileWithTimings(t *testing.T) {
	controller := &MemoryControllerInfo{
		CurrentFrequency: 3200,
		ConfiguredSpeed:  3200,
		ControllerType:   "DDR4",
		Channels:         2,
	}
	
	// Test XMP-like timings (tight)
	xmpTimings := map[string]MemoryTiming{
		"CL":   {Value: 16, Unit: "clocks"},
		"TRCD": {Value: 18, Unit: "clocks"},
	}
	
	profile, err := DetectActiveProfile(controller, xmpTimings)
	if err != nil {
		t.Fatalf("DetectActiveProfile failed: %v", err)
	}
	
	if profile.Type != "XMP" {
		t.Errorf("Expected XMP profile, got %s", profile.Type)
	}
	
	// Test loose timings (custom)
	customTimings := map[string]MemoryTiming{
		"CL":   {Value: 20, Unit: "clocks"}, // Loose timing
		"TRCD": {Value: 22, Unit: "clocks"},
	}
	
	profile, err = DetectActiveProfile(controller, customTimings)
	if err != nil {
		t.Fatalf("DetectActiveProfile failed: %v", err)
	}
	
	if profile.Type != "Custom" {
		t.Errorf("Expected Custom profile, got %s", profile.Type)
	}
}