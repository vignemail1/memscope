// internal/platform/windows/hardware_test.go
package windows

import (
	"runtime"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}
	
	info, err := GetSystemInfo()
	if err != nil {
		t.Fatalf("GetSystemInfo failed: %v", err)
	}
	
	if info.Manufacturer == "" {
		t.Error("Manufacturer should not be empty")
	}
	
	if info.Model == "" {
		t.Error("Model should not be empty")
	}
}

func TestGetMemoryInfo(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}
	
	memory, err := GetMemoryInfo()
	if err != nil {
		t.Fatalf("GetMemoryInfo failed: %v", err)
	}
	
	if len(memory) == 0 {
		t.Error("Should detect at least one memory module")
	}
}

func TestInventoryProviderIntegration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}
	
	provider := NewInventoryProvider()
	inventory, err := provider.CollectInventory()
	if err != nil {
		t.Fatalf("CollectInventory failed: %v", err)
	}
	
	// Validate inventory structure
	if inventory.System.Manufacturer == "" {
		t.Error("System manufacturer should not be empty")
	}
	
	if inventory.System.Model == "" {
		t.Error("System model should not be empty")
	}
	
	if len(inventory.Memory) == 0 {
		t.Error("Should detect at least one memory module")
	}
	
	// Validate memory module data
	for i, module := range inventory.Memory {
		if module.Capacity == 0 {
			t.Errorf("Memory module %d should have non-zero capacity", i)
		}
	}
}