package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCodeCoverage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping coverage test in short mode")
	}
	
	// Run tests with coverage on just a few key packages to avoid timeout
	packages := []string{
		"./internal/cli",
		"./internal/export", 
		"./internal/model",
	}
	
	for _, pkg := range packages {
		cmd := exec.Command("go", "test", "-cover", pkg)
		cmd.Dir = "../.."
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Coverage test failed for %s: %v", pkg, err)
			continue
		}
		
		// Basic validation that we got some coverage info
		outputStr := string(output)
		if !strings.Contains(outputStr, "coverage:") {
			t.Errorf("No coverage info found for %s", pkg)
		}
	}
}

func TestDocumentationCompleteness(t *testing.T) {
	// Build binary for help testing
	binaryPath, cleanup, err := buildMemscope(t)
	if err != nil {
		t.Fatalf("Failed to build memscope: %v", err)
	}
	defer cleanup()
	
	// Test that all commands have help text
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Main help failed: %v", err)
	}
	
	helpText := string(output)
	
	// Check for essential documentation elements
	requiredElements := []string{
		"Usage:",
		"Available Commands:",
		"Flags:",
		"Use",
		"for more information",
	}
	
	for _, element := range requiredElements {
		if !strings.Contains(helpText, element) {
			t.Errorf("Help text missing required element: %s", element)
		}
	}
	
	// Test that help text is reasonably comprehensive
	if len(helpText) < 500 {
		t.Error("Help text seems too brief")
	}
	
	// Test subcommand help completeness
	commands := []string{"inspect", "memory", "recommend", "export", "snapshot", "doctor"}
	for _, command := range commands {
		cmd := exec.Command(binaryPath, command, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("Help for %s failed: %v", command, err)
			continue
		}
		
		cmdHelp := string(output)
		if len(cmdHelp) < 100 {
			t.Errorf("Help for %s seems too brief", command)
		}
		
		// Should have description
		if !strings.Contains(cmdHelp, "Usage:") {
			t.Errorf("Help for %s missing usage information", command)
		}
	}
}

func getPackageFromLine(line string) string {
	// Extract package name from coverage line
	if strings.Contains(line, "/") {
		parts := strings.Fields(line)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return "unknown"
}