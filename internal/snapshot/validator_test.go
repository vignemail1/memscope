package snapshot

import (
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestSnapshotValidator_ValidateSnapshot_ValidSnapshot(t *testing.T) {
	validator := NewValidator()
	
	validSnapshot := &model.Snapshot{
		ID:                   "valid-snapshot",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  time.Now().Add(-time.Minute),
		CollectionFinishedAt: time.Now(),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "memory-device-1",
				Kind: "memory",
			},
		},
	}
	
	if err := validator.ValidateSnapshot(validSnapshot); err != nil {
		t.Errorf("Valid snapshot should pass validation: %v", err)
	}
}

func TestSnapshotValidator_ValidateSnapshot_InvalidSnapshot(t *testing.T) {
	validator := NewValidator()
	
	// Test nil snapshot
	if err := validator.ValidateSnapshot(nil); err == nil {
		t.Error("Nil snapshot should fail validation")
	}
	
	// Test snapshot with missing required fields
	invalidSnapshot := &model.Snapshot{
		// Missing required fields
	}
	
	if err := validator.ValidateSnapshot(invalidSnapshot); err == nil {
		t.Error("Invalid snapshot should fail validation")
	}
}

func TestSnapshotValidator_ValidateSnapshotDetailed(t *testing.T) {
	validator := NewValidator()
	
	// Test valid snapshot
	validSnapshot := &model.Snapshot{
		ID:                   "valid-snapshot",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  time.Now().Add(-time.Minute),
		CollectionFinishedAt: time.Now(),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "memory-device-1",
				Kind: "memory",
			},
		},
	}
	
	result := validator.ValidateSnapshotDetailed(validSnapshot)
	if !result.Valid {
		t.Errorf("Valid snapshot should pass detailed validation")
	}
	
	if len(result.Errors) > 0 {
		t.Errorf("Valid snapshot should have no errors, got: %v", result.Errors)
	}
	
	// Test invalid snapshot
	invalidSnapshot := &model.Snapshot{
		ID: "invalid-snapshot",
		// Missing required fields
	}
	
	result = validator.ValidateSnapshotDetailed(invalidSnapshot)
	if result.Valid {
		t.Error("Invalid snapshot should fail detailed validation")
	}
	
	if len(result.Errors) == 0 {
		t.Error("Invalid snapshot should have errors")
	}
}

func TestSnapshotValidator_ValidateSnapshotDetailed_FutureTimestamp(t *testing.T) {
	validator := NewValidator()
	
	// Snapshot with future timestamp should generate warnings
	futureSnapshot := &model.Snapshot{
		ID:                   "future-snapshot",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  time.Now().Add(2 * time.Hour), // Future timestamp
		CollectionFinishedAt: time.Now().Add(2*time.Hour + time.Minute),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
	}
	
	result := validator.ValidateSnapshotDetailed(futureSnapshot)
	
	// Should still be valid but have warnings
	if !result.Valid {
		t.Error("Snapshot with future timestamp should be valid but have warnings")
	}
	
	if len(result.Warnings) == 0 {
		t.Error("Snapshot with future timestamp should have warnings")
	}
}