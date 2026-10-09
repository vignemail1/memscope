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