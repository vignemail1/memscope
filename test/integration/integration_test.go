package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMemscopeEndToEndWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Build memscope binary
	binaryPath, cleanup, err := buildMemscope(t)
	if err != nil {
		t.Fatalf("Failed to build memscope: %v", err)
	}
	defer cleanup()
	
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "memscope_integration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	t.Run("CLI_Help_System", func(t *testing.T) {
		testCLIHelp(t, binaryPath)
	})
	
	t.Run("Command_Integration", func(t *testing.T) {
		testCommandIntegration(t, binaryPath, tempDir)
	})
	
	t.Run("Data_Flow_Validation", func(t *testing.T) {
		testDataFlowValidation(t, binaryPath, tempDir)
	})
	
	t.Run("Error_Handling", func(t *testing.T) {
		testErrorHandling(t, binaryPath, tempDir)
	})
}

func TestPerformanceValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	binaryPath, cleanup, err := buildMemscope(t)
	if err != nil {
		t.Fatalf("Failed to build memscope: %v", err)
	}
	defer cleanup()
	
	t.Run("Memory_Usage", func(t *testing.T) {
		testMemoryUsage(t, binaryPath)
	})
	
	t.Run("Response_Times", func(t *testing.T) {
		testResponseTimes(t, binaryPath)
	})
	
	t.Run("Large_Dataset_Handling", func(t *testing.T) {
		testLargeDatasetHandling(t, binaryPath)
	})
}

func TestCrossPlatformValidation(t *testing.T) {
	binaryPath, cleanup, err := buildMemscope(t)
	if err != nil {
		t.Fatalf("Failed to build memscope: %v", err)
	}
	defer cleanup()
	
	t.Run("Platform_Detection", func(t *testing.T) {
		testPlatformDetection(t, binaryPath)
	})
	
	t.Run("Graceful_Fallbacks", func(t *testing.T) {
		testGracefulFallbacks(t, binaryPath)
	})
	
	t.Run("File_Operations", func(t *testing.T) {
		testFileOperations(t, binaryPath)
	})
}

func buildMemscope(t *testing.T) (string, func(), error) {
	// Use existing memscope binary if it exists, otherwise build it
	projectRoot := filepath.Join("..", "..")
	existingBinary := filepath.Join(projectRoot, "memscope")
	if os.Getenv("GOOS") == "windows" {
		existingBinary += ".exe"
	}
	
	// Check if existing binary is usable
	if _, err := os.Stat(existingBinary); err == nil {
		// Test if it works
		cmd := exec.Command(existingBinary, "version")
		if err := cmd.Run(); err == nil {
			// Return existing binary with no-op cleanup
			return existingBinary, func() {}, nil
		}
	}
	
	// Build memscope binary for testing
	binaryPath := filepath.Join(os.TempDir(), "memscope_test")
	if os.Getenv("GOOS") == "windows" {
		binaryPath += ".exe"
	}
	
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/memscope")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", func() {}, fmt.Errorf("build failed: %v, output: %s", err, output)
	}
	
	// Return with cleanup function
	cleanup := func() {
		os.Remove(binaryPath)
	}
	return binaryPath, cleanup, nil
}