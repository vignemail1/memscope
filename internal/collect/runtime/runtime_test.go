// internal/collect/runtime/runtime_test.go
package runtime

import (
	"runtime"
	"strings"
	"testing"
	
	"github.com/vignemail1/memscope/internal/platform/windows"
)

func TestCollectMemoryParameters(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Runtime collection requires Windows")
	}
	
	provider := NewProvider()
	params, err := provider.CollectMemoryParameters()
	if err != nil {
		t.Fatalf("CollectMemoryParameters failed: %v", err)
	}
	
	// Validate basic parameters
	if params.CurrentFrequency == 0 {
		t.Error("Current frequency should not be zero")
	}
	
	if len(params.Timings) == 0 {
		t.Error("Should have at least one timing parameter")
	}
	
	// Validate timing data
	for name, timing := range params.Timings {
		if timing.Value == 0 {
			t.Errorf("Timing %s should have non-zero value", name)
		}
	}
}

func TestMemoryVoltages(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Voltage collection requires Windows")
	}
	
	provider := NewProvider()
	params, err := provider.CollectMemoryParameters()
	if err != nil {
		t.Fatalf("CollectMemoryParameters failed: %v", err)
	}
	
	// Validate voltage readings
	if len(params.Voltages) == 0 {
		t.Error("Should have at least one voltage reading")
	}
	
	for name, voltage := range params.Voltages {
		if voltage.Value <= 0 || voltage.Value > 2.0 {
			t.Errorf("Voltage %s value %f is out of reasonable range", name, voltage.Value)
		}
	}
}

// Test that memory timings are realistic based on detected frequency
func TestMemoryTimingsRealistic(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Memory timing tests require Windows")
	}
	
	// Get memory controller info first
	controllerInfo, err := windows.GetMemoryControllerInfo()
	if err != nil {
		t.Fatalf("Failed to get controller info: %v", err)
	}
	
	timings, err := windows.GetCurrentMemoryTimings()
	if err != nil {
		t.Fatalf("Failed to get memory timings: %v", err)
	}
	
	// Verify we have essential timings
	essentialTimings := []string{"CL", "TRCD", "TRP", "TRAS"}
	for _, essential := range essentialTimings {
		if _, exists := timings[essential]; !exists {
			t.Errorf("Missing essential timing: %s", essential)
		}
	}
	
	// Verify timings are reasonable for the frequency
	if clTiming, exists := timings["CL"]; exists {
		// High frequency memory should have higher CL values
		if controllerInfo.CurrentFrequency >= 3200 && clTiming.Value < 14 {
			t.Errorf("CL timing %d seems too low for frequency %d MHz", clTiming.Value, controllerInfo.CurrentFrequency)
		}
		if controllerInfo.CurrentFrequency <= 2133 && clTiming.Value > 17 {
			t.Errorf("CL timing %d seems too high for frequency %d MHz", clTiming.Value, controllerInfo.CurrentFrequency)
		}
	}
}

// Test that voltage collection provides helpful error messages when sensors fail
func TestVoltageCollectionErrorHandling(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Voltage collection tests require Windows")
	}
	
	voltages, err := windows.GetMemoryVoltages()
	
	// Even if there's an error, we should get default voltages
	if len(voltages) == 0 {
		t.Error("Should always return at least default voltage values")
	}
	
	// If there's an error, it should be informative
	if err != nil {
		errorMsg := err.Error()
		if !strings.Contains(errorMsg, "typical") && !strings.Contains(errorMsg, "default") && !strings.Contains(errorMsg, "sensor") {
			t.Errorf("Error message should be informative about fallback behavior, got: %s", errorMsg)
		}
	}
	
	// Should always have VDIMM at minimum
	if _, exists := voltages["VDIMM"]; !exists {
		t.Error("Should always provide VDIMM voltage reading")
	}
}

// Test that profile detection considers timing values, not just frequency
func TestProfileDetectionUsesTimings(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Profile detection tests require Windows")
	}
	
	controllerInfo, err := windows.GetMemoryControllerInfo()
	if err != nil {
		t.Fatalf("Failed to get controller info: %v", err)
	}
	
	timings, err := windows.GetCurrentMemoryTimings()
	if err != nil {
		t.Fatalf("Failed to get timings: %v", err)
	}
	
	profile, err := windows.DetectActiveProfile(controllerInfo, timings)
	if err != nil {
		t.Fatalf("Failed to detect profile: %v", err)
	}
	
	// Profile detection should be more sophisticated than just frequency
	if profile.Type == "" {
		t.Error("Profile type should not be empty")
	}
	
	// For high frequency memory, profile should consider if timings are tight (XMP) or loose (custom)
	if controllerInfo.CurrentFrequency >= 3200 {
		if clTiming, exists := timings["CL"]; exists {
			if profile.Type == "XMP" && clTiming.Value > 18 {
				t.Errorf("Profile detected as XMP but CL %d is too loose for XMP at %d MHz", clTiming.Value, controllerInfo.CurrentFrequency)
			}
		}
	}
}

