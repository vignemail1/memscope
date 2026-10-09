package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegressionSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping regression tests in short mode")
	}
	
	binaryPath, cleanup, err := buildMemscope(t)
	if err != nil {
		t.Fatalf("Failed to build memscope: %v", err)
	}
	defer cleanup()
	
	tempDir, err := os.MkdirTemp("", "memscope_regression")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	t.Run("Command_Interface_Stability", func(t *testing.T) {
		testCommandInterfaceStability(t, binaryPath)
	})
	
	t.Run("Data_Format_Consistency", func(t *testing.T) {
		testDataFormatConsistency(t, binaryPath, tempDir)
	})
	
	t.Run("Backward_Compatibility", func(t *testing.T) {
		testBackwardCompatibility(t, binaryPath, tempDir)
	})
	
	t.Run("Error_Message_Quality", func(t *testing.T) {
		testErrorMessageQuality(t, binaryPath)
	})
}

func testCommandInterfaceStability(t *testing.T, binaryPath string) {
	// Verify all expected commands and flags exist
	
	expectedCommands := map[string][]string{
		"inspect":  {},
		"memory":   {"current", "profiles"},
		"recommend": {},
		"export":   {},
		"snapshot": {"create", "list", "compare", "delete"},
		"doctor":   {},
	}
	
	for command, subcommands := range expectedCommands {
		// Test main command help
		cmd := exec.Command(binaryPath, command, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Command %s help failed: %v", command, err)
			continue
		}
		
		helpText := string(output)
		
		// Verify subcommands are listed
		for _, sub := range subcommands {
			if !strings.Contains(helpText, sub) {
				t.Errorf("Command %s missing subcommand %s in help", command, sub)
			}
		}
		
		// Test subcommands
		for _, sub := range subcommands {
			cmd := exec.Command(binaryPath, command, sub, "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Subcommand %s %s help failed: %v", command, sub, err)
			}
		}
	}
}

func testDataFormatConsistency(t *testing.T, binaryPath, tempDir string) {
	// Test that data formats remain consistent
	
	// Create test snapshot
	testSnapshot := createTestSnapshot()
	snapshotFile := filepath.Join(tempDir, "format_test.json")
	if err := saveTestSnapshot(testSnapshot, snapshotFile); err != nil {
		t.Fatalf("Failed to create test snapshot: %v", err)
	}
	
	// Test JSON export format
	jsonExport := filepath.Join(tempDir, "format_test_export.json")
	cmd := exec.Command(binaryPath, "export", "--input", snapshotFile, "--output", jsonExport, "--format", "json")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("JSON export failed: %v", err)
		return
	}
	
	// Verify JSON structure
	data, err := os.ReadFile(jsonExport)
	if err != nil {
		t.Errorf("Failed to read JSON export: %v", err)
		return
	}
	
	var exportData map[string]interface{}
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Errorf("JSON export is not valid JSON: %v", err)
		return
	}
	
	// Verify expected fields exist
	expectedFields := []string{"id", "schema_version", "devices", "observations"}
	for _, field := range expectedFields {
		if _, exists := exportData[field]; !exists {
			t.Errorf("JSON export missing expected field: %s", field)
		}
	}
	
	// Test CSV export format
	csvExport := filepath.Join(tempDir, "format_test_export.csv")
	cmd = exec.Command(binaryPath, "export", "--input", snapshotFile, "--output", csvExport, "--format", "csv", "--type", "memory")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("CSV export failed: %v", err)
		return
	}
	
	// Verify CSV structure
	csvData, err := os.ReadFile(csvExport)
	if err != nil {
		t.Errorf("Failed to read CSV export: %v", err)
		return
	}
	
	csvContent := string(csvData)
	lines := strings.Split(strings.TrimSpace(csvContent), "\n")
	if len(lines) < 2 {
		t.Error("CSV export should have header + data rows")
		return
	}
	
	// Verify CSV headers
	header := lines[0]
	expectedHeaders := []string{"Manufacturer", "Capacity", "Speed"}
	for _, expectedHeader := range expectedHeaders {
		if !strings.Contains(header, expectedHeader) {
			t.Errorf("CSV header missing expected column: %s", expectedHeader)
		}
	}
}

func testBackwardCompatibility(t *testing.T, binaryPath, tempDir string) {
	// Test that older data formats are still supported
	
	// Create snapshot in "older" format (current format is the baseline)
	oldSnapshot := createTestSnapshot()
	oldSnapshotFile := filepath.Join(tempDir, "old_format.json")
	if err := saveTestSnapshot(oldSnapshot, oldSnapshotFile); err != nil {
		t.Fatalf("Failed to create old format snapshot: %v", err)
	}
	
	// Test that all commands can still process this format
	commands := [][]string{
		{"doctor", "--input", oldSnapshotFile},
		{"recommend", "--input", oldSnapshotFile},
		{"export", "--input", oldSnapshotFile, "--output", filepath.Join(tempDir, "old_export.json")},
	}
	
	for _, cmd := range commands {
		execCmd := exec.Command(binaryPath, cmd...)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			t.Errorf("Command %v failed with old format: %v, output: %s", cmd, err, output)
		}
	}
}

func testErrorMessageQuality(t *testing.T, binaryPath string) {
	// Test that error messages are helpful and actionable
	
	errorTests := []struct {
		name     string
		args     []string
		expected []string // Keywords that should appear in error message
	}{
		{
			name:     "invalid_command",
			args:     []string{"nonexistent-command"},
			expected: []string{"unknown", "command", "help"},
		},
		{
			name:     "missing_required_flag",
			args:     []string{"export"},
			expected: []string{"required", "output"},
		},
		{
			name:     "invalid_file",
			args:     []string{"export", "--input", "/nonexistent/file.json", "--output", "/tmp/test.json"},
			expected: []string{"file", "no such"},
		},
	}
	
	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, test.args...)
			output, err := cmd.CombinedOutput()
			
			if err == nil {
				t.Errorf("Expected command to fail but it succeeded")
				return
			}
			
			outputStr := strings.ToLower(string(output))
			
			for _, keyword := range test.expected {
				if !strings.Contains(outputStr, strings.ToLower(keyword)) {
					t.Errorf("Error message should contain '%s'. Got: %s", keyword, outputStr)
				}
			}
			
			// Error message should not contain debug info or stack traces
			if strings.Contains(outputStr, "goroutine") || strings.Contains(outputStr, "panic") {
				t.Errorf("Error message should not contain debug info: %s", outputStr)
			}
		})
	}
}