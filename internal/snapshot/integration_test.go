package snapshot

import (
	"os"
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestIntegration_SnapshotWorkflow(t *testing.T) {
	// Create temporary directory for test snapshots
	tempDir, err := os.MkdirTemp("", "memscope_integration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	manager := NewManager(tempDir)
	
	// Create baseline snapshot
	baselineSnapshot := createTestSnapshot("baseline", time.Now().Add(-time.Hour))
	
	baselinePath, err := manager.CreateSnapshot(baselineSnapshot, "baseline", "Initial system state")
	if err != nil {
		t.Fatalf("Failed to create baseline snapshot: %v", err)
	}
	
	// Create modified snapshot  
	modifiedSnapshot := createTestSnapshot("modified", time.Now())
	modifiedSnapshot.Devices = append(modifiedSnapshot.Devices, model.Device{
		ID:       "memory-1",
		Kind:     "memory",
		ParentID: "system-0",
	})
	
	modifiedPath, err := manager.CreateSnapshot(modifiedSnapshot, "modified", "System with changes")
	if err != nil {
		t.Fatalf("Failed to create modified snapshot: %v", err)
	}
	
	// List snapshots
	snapshots, err := manager.ListSnapshots()
	if err != nil {
		t.Fatalf("Failed to list snapshots: %v", err)
	}
	
	if len(snapshots) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(snapshots))
	}
	
	// Load and compare snapshots
	loadedBaseline, err := manager.LoadSnapshot(baselinePath)
	if err != nil {
		t.Fatalf("Failed to load baseline snapshot: %v", err)
	}
	
	loadedModified, err := manager.LoadSnapshot(modifiedPath)
	if err != nil {
		t.Fatalf("Failed to load modified snapshot: %v", err)
	}
	
	// Compare snapshots
	comparer := NewComparer()
	comparison, err := comparer.CompareSnapshots(loadedBaseline, loadedModified)
	if err != nil {
		t.Fatalf("Failed to compare snapshots: %v", err)
	}
	
	if !comparison.HasChanges {
		t.Error("Comparison should detect changes between snapshots")
	}
	
	if comparison.Summary.DeviceChanges == 0 {
		t.Error("Comparison should detect device changes")
	}
	
	// Validate snapshots
	validator := NewValidator()
	
	baselineValidation := validator.ValidateSnapshotDetailed(loadedBaseline)
	if !baselineValidation.Valid {
		t.Errorf("Baseline snapshot should be valid: %v", baselineValidation.Errors)
	}
	
	modifiedValidation := validator.ValidateSnapshotDetailed(loadedModified)
	if !modifiedValidation.Valid {
		t.Errorf("Modified snapshot should be valid: %v", modifiedValidation.Errors)
	}
	
	// Test cleanup functionality
	cleanupManager := NewManagerWithConfig(&SnapshotConfig{
		StorageDir:      tempDir,
		MaxSnapshots:    1, // Force cleanup
		RetentionPeriod: 24 * time.Hour,
		AutoCleanup:     true,
	})
	
	// Create one more snapshot to trigger cleanup
	cleanupSnapshot := createTestSnapshot("cleanup", time.Now().Add(time.Minute))
	_, err = cleanupManager.CreateSnapshot(cleanupSnapshot, "cleanup", "Triggers cleanup")
	if err != nil {
		t.Fatalf("Failed to create cleanup snapshot: %v", err)
	}
	
	// Verify cleanup happened
	snapshotsAfterCleanup, err := manager.ListSnapshots()
	if err != nil {
		t.Fatalf("Failed to list snapshots after cleanup: %v", err)
	}
	
	if len(snapshotsAfterCleanup) > 1 {
		t.Errorf("Expected cleanup to limit snapshots to 1, got %d", len(snapshotsAfterCleanup))
	}
	
	// Test delete functionality
	if len(snapshotsAfterCleanup) > 0 {
		err = manager.DeleteSnapshot(snapshotsAfterCleanup[0].Name)
		if err != nil {
			t.Fatalf("Failed to delete snapshot: %v", err)
		}
		
		snapshotsAfterDelete, err := manager.ListSnapshots()
		if err != nil {
			t.Fatalf("Failed to list snapshots after delete: %v", err)
		}
		
		if len(snapshotsAfterDelete) > 0 {
			t.Errorf("Expected 0 snapshots after delete, got %d", len(snapshotsAfterDelete))
		}
	}
}

func TestIntegration_ComparisonFormats(t *testing.T) {
	// Test different comparison scenarios and options
	
	// Create test snapshots with specific differences
	snapshot1 := createTestSnapshot("test1", time.Now().Add(-time.Hour))
	snapshot2 := createTestSnapshot("test2", time.Now())
	
	// Add platform difference
	snapshot2.Platform.OS = "linux"
	
	// Test with different comparison options
	testCases := []struct {
		name    string
		options *ComparisonOptions
	}{
		{
			name: "default_options",
			options: &ComparisonOptions{
				IgnoreTimestamps:   true,
				ShowOnlyChanges:    false,
				IncludeRuntimeData: true,
			},
		},
		{
			name: "include_timestamps",
			options: &ComparisonOptions{
				IgnoreTimestamps:   false,
				ShowOnlyChanges:    false,
				IncludeRuntimeData: true,
			},
		},
		{
			name: "changes_only",
			options: &ComparisonOptions{
				IgnoreTimestamps:   true,
				ShowOnlyChanges:    true,
				IncludeRuntimeData: true,
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comparer := NewComparerWithOptions(tc.options)
			
			comparison, err := comparer.CompareSnapshots(snapshot1, snapshot2)
			if err != nil {
				t.Fatalf("Comparison failed: %v", err)
			}
			
			// Should always detect platform changes
			if !comparison.HasChanges {
				t.Error("Should detect platform changes")
			}
			
			// Check that system changes are detected
			foundSystemChanges := false
			for _, change := range comparison.Changes {
				if change.Category == "system" {
					foundSystemChanges = true
					break
				}
			}
			
			if !foundSystemChanges {
				t.Error("Should detect system/platform changes")
			}
		})
	}
}

func TestIntegration_ValidationScenarios(t *testing.T) {
	validator := NewValidator()
	
	// Test various validation scenarios
	testCases := []struct {
		name      string
		snapshot  *model.Snapshot
		shouldBeValid bool
		expectedWarnings int
	}{
		{
			name:         "valid_snapshot",
			snapshot:     createTestSnapshot("valid", time.Now()),
			shouldBeValid: true,
			expectedWarnings: 1, // device memory-0 has no observations
		},
		{
			name: "future_timestamp",
			snapshot: createTestSnapshot("future", time.Now().Add(2*time.Hour)),
			shouldBeValid: true,
			expectedWarnings: 3, // Both start and end timestamps are in future + device has no observations
		},
		{
			name: "no_devices",
			snapshot: &model.Snapshot{
				ID:                   "no-devices",
				SchemaVersion:        "1.0",
				ToolVersion:          "test",
				CollectionStartedAt:  time.Now().Add(-time.Minute),
				CollectionFinishedAt: time.Now(),
				Platform:             model.Platform{OS: "windows", Arch: "amd64"},
				Devices:              []model.Device{},
				Observations:         []model.Observation{},
			},
			shouldBeValid: true,
			expectedWarnings: 2, // Warning about no devices + no observations
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateSnapshotDetailed(tc.snapshot)
			
			if result.Valid != tc.shouldBeValid {
				t.Errorf("Expected valid=%v, got %v. Errors: %v", tc.shouldBeValid, result.Valid, result.Errors)
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tc.expectedWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

// Helper function to create a test snapshot
func createTestSnapshot(id string, timestamp time.Time) *model.Snapshot {
	return &model.Snapshot{
		ID:                   id,
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  timestamp,
		CollectionFinishedAt: timestamp.Add(time.Second),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "system-0",
				Kind: "system",
			},
			{
				ID:       "memory-0",
				Kind:     "memory",
				ParentID: "system-0",
			},
		},
		Observations: []model.Observation{
			{
				DeviceID:   "system-0",
				Scope:      "hardware",
				Parameter:  "manufacturer",
				Value:      &model.Value{Text: stringPtr("Test System")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: timestamp,
			},
		},
		Profiles:     []model.Profile{},
		Diagnostics:  []model.Diagnostic{},
		Capabilities: []model.Capability{},
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}