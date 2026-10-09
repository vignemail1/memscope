package snapshot

import (
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestSnapshotComparer_CompareSnapshots_NoChanges(t *testing.T) {
	// Create two identical snapshots
	timestamp1 := time.Now().Add(-time.Hour)
	timestamp2 := time.Now()
	
	snapshot1 := &model.Snapshot{
		ID:                   "snapshot1",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  timestamp1,
		CollectionFinishedAt: timestamp1.Add(time.Second),
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
	
	snapshot2 := &model.Snapshot{
		ID:                   "snapshot2",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  timestamp2,
		CollectionFinishedAt: timestamp2.Add(time.Second),
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
	
	comparer := NewComparer()
	comparison, err := comparer.CompareSnapshots(snapshot1, snapshot2)
	if err != nil {
		t.Fatalf("CompareSnapshots failed: %v", err)
	}
	
	// Should detect no changes (ignoring timestamps)
	if comparison.HasChanges {
		t.Error("Should not detect changes between identical snapshots (ignoring timestamps)")
	}
	
	if comparison.Summary.TotalChanges != 0 {
		t.Errorf("Expected 0 total changes, got %d", comparison.Summary.TotalChanges)
	}
}

func TestSnapshotComparer_CompareSnapshots_WithDeviceChanges(t *testing.T) {
	timestamp1 := time.Now().Add(-time.Hour)
	timestamp2 := time.Now()
	
	snapshot1 := &model.Snapshot{
		ID:                   "snapshot1",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  timestamp1,
		CollectionFinishedAt: timestamp1.Add(time.Second),
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
	
	snapshot2 := &model.Snapshot{
		ID:                   "snapshot2", 
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  timestamp2,
		CollectionFinishedAt: timestamp2.Add(time.Second),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "memory-device-1",
				Kind: "memory",
			},
			{
				ID:   "memory-device-2", // Added device
				Kind: "memory",
			},
		},
	}
	
	comparer := NewComparer()
	comparison, err := comparer.CompareSnapshots(snapshot1, snapshot2)
	if err != nil {
		t.Fatalf("CompareSnapshots failed: %v", err)
	}
	
	// Should detect changes
	if !comparison.HasChanges {
		t.Error("Should detect changes between snapshots")
	}
	
	// Should find device changes
	foundDeviceChanges := false
	for _, change := range comparison.Changes {
		if change.Category == "device" {
			foundDeviceChanges = true
			break
		}
	}
	
	if !foundDeviceChanges {
		t.Error("Should detect device configuration changes")
	}
	
	if comparison.Summary.DeviceChanges == 0 {
		t.Error("Summary should show device changes")
	}
}

func TestSnapshotComparer_CompareSnapshots_WithPlatformChanges(t *testing.T) {
	timestamp1 := time.Now().Add(-time.Hour)
	timestamp2 := time.Now()
	
	snapshot1 := &model.Snapshot{
		ID:                   "snapshot1",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version-1",
		CollectionStartedAt:  timestamp1,
		CollectionFinishedAt: timestamp1.Add(time.Second),
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
	
	snapshot2 := &model.Snapshot{
		ID:                   "snapshot2",
		SchemaVersion:        "1.0", 
		ToolVersion:          "test-version-2", // Changed version
		CollectionStartedAt:  timestamp2,
		CollectionFinishedAt: timestamp2.Add(time.Second),
		Platform: model.Platform{
			OS:   "linux", // Changed OS
			Arch: "amd64",
		},
		Devices: []model.Device{
			{
				ID:   "memory-device-1",
				Kind: "memory",
			},
		},
	}
	
	comparer := NewComparer()
	comparison, err := comparer.CompareSnapshots(snapshot1, snapshot2)
	if err != nil {
		t.Fatalf("CompareSnapshots failed: %v", err)
	}
	
	// Should detect changes
	if !comparison.HasChanges {
		t.Error("Should detect changes between snapshots")
	}
	
	// Should find system changes
	foundSystemChanges := false
	for _, change := range comparison.Changes {
		if change.Category == "system" {
			foundSystemChanges = true
			break
		}
	}
	
	if !foundSystemChanges {
		t.Error("Should detect system platform changes")
	}
}

func TestSnapshotComparer_CompareSnapshots_NilInput(t *testing.T) {
	comparer := NewComparer()
	
	_, err := comparer.CompareSnapshots(nil, nil)
	if err == nil {
		t.Error("CompareSnapshots should fail with nil inputs")
	}
	
	snapshot := &model.Snapshot{
		ID:                   "test",
		SchemaVersion:        "1.0",
		ToolVersion:          "test-version",
		CollectionStartedAt:  time.Now(),
		CollectionFinishedAt: time.Now().Add(time.Second),
		Platform: model.Platform{
			OS:   "windows",
			Arch: "amd64",
		},
	}
	
	_, err = comparer.CompareSnapshots(snapshot, nil)
	if err == nil {
		t.Error("CompareSnapshots should fail with one nil input")
	}
	
	_, err = comparer.CompareSnapshots(nil, snapshot)
	if err == nil {
		t.Error("CompareSnapshots should fail with one nil input")
	}
}