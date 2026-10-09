package cli

import (
	"bytes"
	"testing"
)

func TestInspectCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"inspect"})
	
	err := cmd.Execute()
	
	// On Windows, command should succeed
	// On other platforms, it should fail with expected message
	if err != nil {
		errorMsg := err.Error()
		if errorMsg != "failed to collect hardware inventory: inventory collection not supported on darwin" &&
		   errorMsg != "failed to collect hardware inventory: inventory collection not supported on linux" {
			t.Fatalf("inspect command failed with unexpected error: %v", err)
		}
		// Test passed - got expected error on unsupported platform
		return
	}
	
	// If no error, we should have output (Windows case)
	output := buf.String()
	if output == "" {
		t.Fatal("inspect command succeeded but produced no output")
	}
}

func TestMemoryCurrentCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"memory", "current"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("memory current command failed: %v", err)
	}
}

func TestMemoryProfilesCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"memory", "profiles"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("memory profiles command failed: %v", err)
	}
}

func TestRecommendCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"recommend"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("recommend command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("recommend command produced no output")
	}
}

func TestDoctorCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"doctor"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("doctor command produced no output")
	}
}

func TestExportCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"export"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("export command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("export command produced no output")
	}
}

func TestSnapshotCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"snapshot"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("snapshot command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("snapshot command produced no output")
	}
}

func TestCompareCommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"compare", "snapshot1.json", "snapshot2.json"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("compare command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("compare command produced no output")
	}
}