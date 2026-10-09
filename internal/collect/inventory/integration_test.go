package inventory

import (
	"runtime"
	"testing"
)

func TestCollectFullInventory(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Integration test requires Windows")
	}
	
	provider := NewProvider()
	inventory, err := provider.CollectInventory()
	if err != nil {
		t.Fatalf("CollectInventory failed: %v", err)
	}
	
	// Validate system info
	if inventory.System.Manufacturer == "" {
		t.Error("System manufacturer should not be empty")
	}
	
	if inventory.System.Model == "" {
		t.Error("System model should not be empty")
	}
	
	// Validate BIOS info
	if inventory.BIOS.Vendor == "" {
		t.Error("BIOS vendor should not be empty")
	}
	
	// Validate CPU info  
	if inventory.CPU.Name == "" {
		t.Error("CPU name should not be empty")
	}

	if inventory.CPU.Cores == 0 {
		t.Error("CPU cores should not be zero")
	}

	if inventory.CPU.Threads == 0 {
		t.Error("CPU threads should not be zero")
	}
	
	// Validate memory info
	if len(inventory.Memory) == 0 {
		t.Error("Should detect at least one memory module")
	}
	
	// Validate memory module data
	for i, module := range inventory.Memory {
		if module.Capacity == 0 {
			t.Errorf("Memory module %d should have non-zero capacity", i)
		}
		
		if module.Speed == 0 {
			t.Errorf("Memory module %d should have non-zero speed", i)
		}
	}
}

func TestInventoryToSnapshot(t *testing.T) {
	// Create mock inventory
	inv := &SystemInventory{
		System: SystemInfo{
			Manufacturer: "Test Manufacturer",
			Model:        "Test Model",
		},
		CPU: CPUInfo{
			Name:         "Test CPU",
			Manufacturer: "Test CPU Vendor",
			Architecture: "x64",
			Cores:        4,
			Threads:      8,
			MaxClockMHz:  3200,
		},
		Memory: []MemoryModuleInfo{
			{
				DeviceLocator: "DIMM_A1",
				Capacity:      8 * 1024 * 1024 * 1024, // 8GB
				Speed:         3200,
				Manufacturer:  "Test RAM",
			},
		},
	}
	
	snapshot, err := inv.ToSnapshot()
	if err != nil {
		t.Fatalf("ToSnapshot failed: %v", err)
	}
	
	if len(snapshot.Devices) == 0 {
		t.Error("Snapshot should contain devices")
	}
	
	if len(snapshot.Observations) == 0 {
		t.Error("Snapshot should contain observations")
	}

	// Validate CPU is in snapshot
	cpuDeviceFound := false
	cpuObservationsFound := 0
	for _, device := range snapshot.Devices {
		if device.ID == "cpu-0" {
			cpuDeviceFound = true
			break
		}
	}
	if !cpuDeviceFound {
		t.Error("Snapshot should contain CPU device")
	}

	for _, obs := range snapshot.Observations {
		if obs.DeviceID == "cpu-0" {
			cpuObservationsFound++
		}
	}
	if cpuObservationsFound == 0 {
		t.Error("Snapshot should contain CPU observations")
	}
}