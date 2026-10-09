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
	if err != nil {
		t.Fatalf("inspect command failed: %v", err)
	}
	
	output := buf.String()
	if output == "" {
		t.Fatal("inspect command produced no output")
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