package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testMemoryUsage(t *testing.T, binaryPath string) {
	// Test memory usage during various operations
	testCases := []struct {
		name string
		args []string
	}{
		{"help", []string{"--help"}},
		{"version", []string{"version"}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Monitor memory usage using system tools
			var memBefore, memAfter runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&memBefore)
			
			cmd := exec.Command(binaryPath, tc.args...)
			_, err := cmd.CombinedOutput()
			
			runtime.GC()
			runtime.ReadMemStats(&memAfter)
			
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				t.Errorf("Command failed: %v", err)
			}
			
			// Basic memory usage validation
			memUsed := memAfter.TotalAlloc - memBefore.TotalAlloc
			if memUsed > 100*1024*1024 { // 100MB threshold
				t.Errorf("Command used excessive memory: %d bytes", memUsed)
			}
		})
	}
}

func testResponseTimes(t *testing.T, binaryPath string) {
	// Test response times for various commands
	testCases := []struct {
		name     string
		args     []string
		maxTime  time.Duration
	}{
		{"help", []string{"--help"}, 2 * time.Second},
		{"version", []string{"version"}, 2 * time.Second},
		{"inspect_help", []string{"inspect", "--help"}, 2 * time.Second},
		{"export_help", []string{"export", "--help"}, 2 * time.Second},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			
			cmd := exec.Command(binaryPath, tc.args...)
			_, err := cmd.CombinedOutput()
			
			elapsed := time.Since(start)
			
			if err != nil && !strings.Contains(err.Error(), "exit status") {
				t.Errorf("Command failed: %v", err)
			}
			
			if elapsed > tc.maxTime {
				t.Errorf("Command took too long: %v (max: %v)", elapsed, tc.maxTime)
			}
			
			t.Logf("Command %s completed in %v", tc.name, elapsed)
		})
	}
}

func testLargeDatasetHandling(t *testing.T, binaryPath string) {
	// Create a large test dataset
	tempDir, err := os.MkdirTemp("", "memscope_large_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create snapshot with many memory modules
	largeSnapshot := createLargeTestSnapshot(16) // 16 memory modules
	snapshotFile := filepath.Join(tempDir, "large_snapshot.json")
	if err := saveTestSnapshot(largeSnapshot, snapshotFile); err != nil {
		t.Fatalf("Failed to create large test snapshot: %v", err)
	}
	
	// Test operations with large dataset
	testCases := []struct {
		name string
		args []string
	}{
		{"doctor_large", []string{"doctor", "--input", snapshotFile}},
		{"export_large", []string{"export", "--input", snapshotFile, "--output", filepath.Join(tempDir, "large_export.json")}},
		{"recommend_large", []string{"recommend", "--input", snapshotFile}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			
			cmd := exec.Command(binaryPath, tc.args...)
			output, err := cmd.CombinedOutput()
			
			elapsed := time.Since(start)
			
			if err != nil {
				t.Errorf("Command failed: %v, output: %s", err, output)
			}
			
			// Should handle large datasets within reasonable time (30 seconds)
			if elapsed > 30*time.Second {
				t.Errorf("Command took too long with large dataset: %v", elapsed)
			}
			
			t.Logf("Large dataset command %s completed in %v", tc.name, elapsed)
		})
	}
}

func createLargeTestSnapshot(numModules int) map[string]interface{} {
	snapshot := createTestSnapshot()
	snapshot["id"] = fmt.Sprintf("test-snapshot-large-%d", numModules)
	
	// Create many memory devices
	devices := []map[string]interface{}{
		{
			"id":   "system-0",
			"kind": "system",
		},
	}
	
	for i := 0; i < numModules; i++ {
		devices = append(devices, map[string]interface{}{
			"id":        fmt.Sprintf("memory-%d", i),
			"kind":      "memory",
			"parent_id": "system-0",
		})
	}
	
	snapshot["devices"] = devices
	
	// Add many observations (multiple per module)
	observations := []map[string]interface{}{}
	for i := 0; i < numModules; i++ {
		deviceID := fmt.Sprintf("memory-%d", i)
		
		// Manufacturer
		observations = append(observations, map[string]interface{}{
			"device_id":   deviceID,
			"scope":       "hardware",
			"parameter":   "manufacturer",
			"value":       map[string]interface{}{"text": fmt.Sprintf("Manufacturer_%d", i%4)},
			"source":      "wmi",
			"status":      "observed",
			"captured_at": time.Now().Format(time.RFC3339),
		})
		
		// Speed
		observations = append(observations, map[string]interface{}{
			"device_id":   deviceID,
			"scope":       "hardware",
			"parameter":   "speed",
			"value":       map[string]interface{}{"unsigned": 3200 + (i%3)*400},
			"unit":        "MT/s",
			"source":      "wmi",
			"status":      "observed",
			"captured_at": time.Now().Format(time.RFC3339),
		})
		
		// Voltage
		observations = append(observations, map[string]interface{}{
			"device_id":   deviceID,
			"scope":       "spd",
			"parameter":   "voltage_level",
			"value":       map[string]interface{}{"decimal": 1.35 + float64(i%2)*0.05},
			"source":      "spd",
			"status":      "observed",
			"captured_at": time.Now().Format(time.RFC3339),
		})
	}
	
	snapshot["observations"] = observations
	
	return snapshot
}