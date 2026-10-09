package snapshot

import (
	"os"
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestSnapshotManager_CreateSnapshot(t *testing.T) {
	// Create temporary directory for test snapshots
	tempDir, err := os.MkdirTemp("", "memscope_snapshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	manager := NewManager(tempDir)
	
	// Test snapshot creation
	snapshot := &model.Snapshot{
		ID:                   "test-snapshot-id",
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
				ID:   "test-device",
				Kind: "memory",
			},
		},
	}
	
	snapshotPath, err := manager.CreateSnapshot(snapshot, "test-snapshot", "Test description")
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	
	if snapshotPath == "" {
		t.Error("CreateSnapshot should return non-empty path")
	}
	
	// Test snapshot listing
	snapshots, err := manager.ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	
	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}
	
	if snapshots[0].Name != "test-snapshot" {
		t.Errorf("Expected snapshot name 'test-snapshot', got '%s'", snapshots[0].Name)
	}
}

func TestSnapshotManager_ListSnapshots_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "memscope_snapshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	manager := NewManager(tempDir)
	
	snapshots, err := manager.ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	
	if len(snapshots) != 0 {
		t.Errorf("Expected 0 snapshots in empty directory, got %d", len(snapshots))
	}
}

func TestSnapshotManager_LoadSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "memscope_snapshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	manager := NewManager(tempDir)
	
	// Create a snapshot first
	originalSnapshot := &model.Snapshot{
		ID:                   "test-load-id",
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
				ID:   "test-device",
				Kind: "memory",
			},
		},
	}
	
	snapshotPath, err := manager.CreateSnapshot(originalSnapshot, "test-load", "Test load")
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	
	// Load the snapshot back
	loadedSnapshot, err := manager.LoadSnapshot(snapshotPath)
	if err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}
	
	if loadedSnapshot.ID != originalSnapshot.ID {
		t.Errorf("Expected snapshot ID '%s', got '%s'", originalSnapshot.ID, loadedSnapshot.ID)
	}
	
	if loadedSnapshot.SchemaVersion != originalSnapshot.SchemaVersion {
		t.Errorf("Expected schema version '%s', got '%s'", originalSnapshot.SchemaVersion, loadedSnapshot.SchemaVersion)
	}
}

func TestSnapshotManager_DeleteSnapshot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "memscope_snapshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	manager := NewManager(tempDir)
	
	// Create a snapshot first
	snapshot := &model.Snapshot{
		ID:                   "test-delete-id",
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
				ID:   "test-device",
				Kind: "memory",
			},
		},
	}
	
	_, err = manager.CreateSnapshot(snapshot, "test-delete", "Test delete")
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	
	// Verify it exists
	snapshots, err := manager.ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("Expected 1 snapshot, got %d", len(snapshots))
	}
	
	// Delete the snapshot
	err = manager.DeleteSnapshot("test-delete")
	if err != nil {
		t.Fatalf("DeleteSnapshot failed: %v", err)
	}
	
	// Verify it's gone
	snapshots, err = manager.ListSnapshots()
	if err != nil {
		t.Fatalf("ListSnapshots failed: %v", err)
	}
	if len(snapshots) != 0 {
		t.Errorf("Expected 0 snapshots after deletion, got %d", len(snapshots))
	}
}