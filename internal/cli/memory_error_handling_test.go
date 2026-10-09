package cli

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func TestMemoryCurrentCommandErrorHandling(t *testing.T) {
	// Test the CLI error handling for memory current command
	cmd := newMemoryCurrentCmd()
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	
	// Run the command
	err := cmd.Execute()
	
	if runtime.GOOS != "windows" {
		// On non-Windows platforms, should show platform message
		if err != nil {
			t.Errorf("Command should not return error on non-Windows platforms, got: %v", err)
		}
		
		output := stdout.String()
		if !strings.Contains(output, "only supported on Windows") {
			t.Errorf("Expected platform message, got: %s", output)
		}
	}
	// Note: We can't test the Windows-specific error handling paths without 
	// being on Windows and having mock implementations
}

func TestMemoryCommandStructure(t *testing.T) {
	// Test that the memory command is properly structured
	cmd := newMemoryCmd()
	
	// Should have subcommands
	if !cmd.HasSubCommands() {
		t.Error("Memory command should have subcommands")
	}
	
	// Check for expected subcommands
	subCommands := cmd.Commands()
	hasCurrentCmd := false
	hasProfilesCmd := false
	
	for _, subCmd := range subCommands {
		switch subCmd.Use {
		case "current":
			hasCurrentCmd = true
		case "profiles":
			hasProfilesCmd = true
		}
	}
	
	if !hasCurrentCmd {
		t.Error("Memory command should have 'current' subcommand")
	}
	
	if !hasProfilesCmd {
		t.Error("Memory command should have 'profiles' subcommand")
	}
}