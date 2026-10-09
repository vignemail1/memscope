// internal/collect/runtime/runtime_test.go
package runtime

import (
	"runtime"
	"testing"
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