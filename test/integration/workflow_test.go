package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testCLIHelp(t *testing.T, binaryPath string) {
	// Test main help
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Main help failed: %v", err)
	}
	
	helpOutput := string(output)
	requiredCommands := []string{"inspect", "memory", "recommend", "export", "snapshot", "doctor"}
	for _, command := range requiredCommands {
		if !strings.Contains(helpOutput, command) {
			t.Errorf("Help output missing command: %s", command)
		}
	}
	
	// Test subcommand help
	subcommands := map[string][]string{
		"memory":   {"current", "profiles"},
		"snapshot": {"create", "list", "compare", "delete"},
		"export":   {},
	}
	
	for command, subs := range subcommands {
		cmd := exec.Command(binaryPath, command, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Help for %s failed: %v", command, err)
			continue
		}
		
		helpOutput := string(output)
		for _, sub := range subs {
			if !strings.Contains(helpOutput, sub) {
				t.Errorf("Help for %s missing subcommand: %s", command, sub)
			}
		}
	}
}

func testCommandIntegration(t *testing.T, binaryPath, tempDir string) {
	// Test command chaining: inspect → export → snapshot compare
	
	// 1. Test inspect command
	t.Log("Testing inspect command")
	cmd := exec.Command(binaryPath, "inspect")
	output, err := cmd.CombinedOutput()
	
	// Should work or gracefully fail on non-Windows
	if runtime.GOOS != "windows" {
		if err == nil || !strings.Contains(string(output), "Windows") {
			t.Log("Non-Windows platform detected, skipping Windows-specific tests")
			return
		}
	}
	
	// 2. Test export command with mock data
	exportFile := filepath.Join(tempDir, "test_export.json")
	t.Log("Testing export command")
	
	// Create a test snapshot for export
	testSnapshot := createTestSnapshot()
	snapshotFile := filepath.Join(tempDir, "test_snapshot.json")
	if err := saveTestSnapshot(testSnapshot, snapshotFile); err != nil {
		t.Fatalf("Failed to create test snapshot: %v", err)
	}
	
	cmd = exec.Command(binaryPath, "export", "--input", snapshotFile, "--output", exportFile, "--type", "snapshot")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Export command failed: %v, output: %s", err, output)
	}
	
	// Verify export file exists and is valid JSON
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	} else {
		// Validate JSON
		data, err := os.ReadFile(exportFile)
		if err != nil {
			t.Errorf("Failed to read export file: %v", err)
		} else {
			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				t.Errorf("Export file is not valid JSON: %v", err)
			}
		}
	}
	
	// 3. Test recommendation command
	t.Log("Testing recommend command")
	cmd = exec.Command(binaryPath, "recommend", "--input", snapshotFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Recommend command failed: %v, output: %s", err, output)
	}
	
	// 4. Test doctor command
	t.Log("Testing doctor command")
	doctorReportFile := filepath.Join(tempDir, "doctor_report.json")
	cmd = exec.Command(binaryPath, "doctor", "--input", snapshotFile, "--format", "json", "--output", doctorReportFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Doctor command failed: %v, output: %s", err, output)
	}
	
	// Verify doctor report
	if _, err := os.Stat(doctorReportFile); os.IsNotExist(err) {
		t.Error("Doctor report file was not created")
	}
}

func testDataFlowValidation(t *testing.T, binaryPath, tempDir string) {
	// Test data consistency across commands
	
	// Create test snapshot
	testSnapshot := createTestSnapshot()
	snapshotFile := filepath.Join(tempDir, "data_flow_test.json")
	if err := saveTestSnapshot(testSnapshot, snapshotFile); err != nil {
		t.Fatalf("Failed to create test snapshot: %v", err)
	}
	
	// Export in different formats and verify consistency
	jsonExport := filepath.Join(tempDir, "export.json")
	csvExport := filepath.Join(tempDir, "export.csv")
	
	// Export as JSON
	cmd := exec.Command(binaryPath, "export", "--input", snapshotFile, "--output", jsonExport, "--type", "memory", "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("JSON export failed: %v, output: %s", err, output)
	}
	
	// Export as CSV
	cmd = exec.Command(binaryPath, "export", "--input", snapshotFile, "--output", csvExport, "--type", "memory", "--format", "csv")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("CSV export failed: %v, output: %s", err, output)
	}
	
	// Verify both files exist and have reasonable content
	jsonData, err := os.ReadFile(jsonExport)
	if err != nil {
		t.Errorf("Failed to read JSON export: %v", err)
	} else if len(jsonData) < 10 {
		t.Error("JSON export appears to be empty or too small")
	}
	
	csvData, err := os.ReadFile(csvExport)
	if err != nil {
		t.Errorf("Failed to read CSV export: %v", err)
	} else if len(csvData) < 10 {
		t.Error("CSV export appears to be empty or too small")
	}
	
	// Test snapshot comparison
	snapshot2File := filepath.Join(tempDir, "data_flow_test2.json")
	testSnapshot2 := createModifiedTestSnapshot()
	if err := saveTestSnapshot(testSnapshot2, snapshot2File); err != nil {
		t.Fatalf("Failed to create second test snapshot: %v", err)
	}
	
	cmd = exec.Command(binaryPath, "snapshot", "compare", snapshotFile, snapshot2File)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Snapshot compare failed: %v, output: %s", err, output)
	}
	
	// Should detect differences
	outputStr := string(output)
	if !strings.Contains(outputStr, "changes") && !strings.Contains(outputStr, "differences") {
		t.Error("Snapshot comparison should detect changes between different snapshots")
	}
}

func testErrorHandling(t *testing.T, binaryPath, tempDir string) {
	// Test various error conditions
	
	// 1. Invalid command
	cmd := exec.Command(binaryPath, "invalid-command")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Invalid command should return error")
	}
	if !strings.Contains(string(output), "unknown command") && !strings.Contains(string(output), "invalid") {
		t.Error("Error message should indicate unknown command")
	}
	
	// 2. Missing required flags
	cmd = exec.Command(binaryPath, "export")
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Export without required flags should fail")
	}
	
	// 3. Invalid file paths
	cmd = exec.Command(binaryPath, "export", "--input", "/nonexistent/file.json", "--output", "/tmp/test.json")
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Export with nonexistent input file should fail")
	}
	
	// 4. Invalid JSON input
	invalidJsonFile := filepath.Join(tempDir, "invalid.json")
	if err := os.WriteFile(invalidJsonFile, []byte("invalid json content"), 0644); err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}
	
	cmd = exec.Command(binaryPath, "doctor", "--input", invalidJsonFile)
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Error("Doctor with invalid JSON should fail")
	}
}

// Helper functions for creating test data

func createTestSnapshot() map[string]interface{} {
	return map[string]interface{}{
		"id":                     "test-snapshot-integration",
		"schema_version":         "0.1.0",
		"tool_version":           "0.1.0",
		"collection_started_at":  time.Now().Format(time.RFC3339),
		"collection_finished_at": time.Now().Add(5 * time.Second).Format(time.RFC3339),
		"platform": map[string]interface{}{
			"os":   "windows",
			"arch": "amd64",
		},
		"devices": []map[string]interface{}{
			{
				"id":   "system-0",
				"kind": "system",
			},
			{
				"id":        "memory-0",
				"kind":      "memory",
				"parent_id": "system-0",
			},
			{
				"id":        "memory-1",
				"kind":      "memory",
				"parent_id": "system-0",
			},
		},
		"observations": []map[string]interface{}{
			{
				"device_id":   "memory-0",
				"scope":       "hardware",
				"parameter":   "manufacturer",
				"value":       map[string]interface{}{"text": "Corsair"},
				"source":      "wmi",
				"status":      "observed",
				"captured_at": time.Now().Format(time.RFC3339),
			},
			{
				"device_id":   "memory-0",
				"scope":       "hardware",
				"parameter":   "speed",
				"value":       map[string]interface{}{"unsigned": 3200},
				"unit":        "MT/s",
				"source":      "wmi",
				"status":      "observed",
				"captured_at": time.Now().Format(time.RFC3339),
			},
			{
				"device_id":   "memory-0",
				"scope":       "spd",
				"parameter":   "voltage_level",
				"value":       map[string]interface{}{"decimal": 1.35},
				"source":      "spd",
				"status":      "observed",
				"captured_at": time.Now().Format(time.RFC3339),
			},
			{
				"device_id":   "memory-1",
				"scope":       "hardware",
				"parameter":   "manufacturer",
				"value":       map[string]interface{}{"text": "Kingston"},
				"source":      "wmi",
				"status":      "observed",
				"captured_at": time.Now().Format(time.RFC3339),
			},
		},
		"profiles":     []map[string]interface{}{},
		"diagnostics":  []map[string]interface{}{},
		"capabilities": []map[string]interface{}{},
	}
}

func createModifiedTestSnapshot() map[string]interface{} {
	snapshot := createTestSnapshot()
	snapshot["id"] = "test-snapshot-integration-modified"
	
	// Modify observations to create differences
	observations := snapshot["observations"].([]map[string]interface{})
	if len(observations) >= 2 {
		// Change speed from 3200 to 3600
		observations[1]["value"] = map[string]interface{}{"unsigned": 3600}
		// Change voltage from 1.35 to 1.4
		observations[2]["value"] = map[string]interface{}{"decimal": 1.4}
	}
	
	return snapshot
}

func saveTestSnapshot(snapshot map[string]interface{}, filename string) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}