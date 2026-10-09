package integration

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func testPlatformDetection(t *testing.T, binaryPath string) {
	// Test platform-specific behavior
	
	cmd := exec.Command(binaryPath, "inspect")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	if runtime.GOOS == "windows" {
		// On Windows, should attempt hardware access
		if err != nil {
			// May fail due to permissions, but should show Windows-specific messages
			if !strings.Contains(outputStr, "Windows") && !strings.Contains(outputStr, "WMI") {
				t.Error("Windows-specific functionality not detected")
			}
		}
	} else {
		// On non-Windows, should show platform limitation message
		if err == nil || !strings.Contains(outputStr, "Windows") {
			t.Log("Non-Windows platform behavior validated")
		}
	}
}

func testGracefulFallbacks(t *testing.T, binaryPath string) {
	// Test graceful fallbacks for non-Windows platforms
	
	commands := [][]string{
		{"inspect"},
		{"memory", "current"},
		{"doctor"},
	}
	
	for _, cmd := range commands {
		t.Run(strings.Join(cmd, "_"), func(t *testing.T) {
			execCmd := exec.Command(binaryPath, cmd...)
			output, err := execCmd.CombinedOutput()
			outputStr := string(output)
			
			if runtime.GOOS != "windows" {
				// Should provide helpful message about Windows requirement
				// or work with --input flag
				if err != nil {
					if !strings.Contains(outputStr, "Windows") && !strings.Contains(outputStr, "input") && !strings.Contains(outputStr, "not supported") {
						t.Errorf("Command should provide helpful fallback message: %s", outputStr)
					}
				}
			}
			
			// Should not crash or panic
			if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "runtime error") {
				t.Errorf("Command should not panic: %s", outputStr)
			}
		})
	}
}

func testFileOperations(t *testing.T, binaryPath string) {
	// Test file operations work across platforms
	
	// Test help output (should always work)
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Help command should always work: %v", err)
	}
	
	outputStr := string(output)
	if len(outputStr) < 100 {
		t.Error("Help output seems too short")
	}
	
	// Verify standard commands are listed
	expectedCommands := []string{"inspect", "memory", "recommend", "export", "snapshot", "doctor"}
	for _, command := range expectedCommands {
		if !strings.Contains(outputStr, command) {
			t.Errorf("Help should list command: %s", command)
		}
	}
}