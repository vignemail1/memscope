package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

// Manager handles snapshot storage and retrieval
type Manager struct {
	config *SnapshotConfig
}

// NewManager creates a new snapshot manager with default configuration
func NewManager(storageDir string) *Manager {
	return &Manager{
		config: &SnapshotConfig{
			StorageDir:       storageDir,
			MaxSnapshots:     100,
			RetentionPeriod:  30 * 24 * time.Hour, // 30 days
			AutoCleanup:      true,
			CompressionLevel: 0, // No compression by default
			IncludeRuntime:   true,
		},
	}
}

// NewManagerWithConfig creates a snapshot manager with custom configuration
func NewManagerWithConfig(config *SnapshotConfig) *Manager {
	return &Manager{config: config}
}

// CreateSnapshot saves a snapshot to storage with metadata
func (m *Manager) CreateSnapshot(snapshot *model.Snapshot, name, description string) (string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(m.config.StorageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}
	
	// Generate filename with timestamp
	timestamp := time.Now()
	
	safeFilename := sanitizeFilename(name)
	filename := fmt.Sprintf("%s_%s.json", 
		timestamp.Format("2006-01-02_15-04-05"), 
		safeFilename)
	filePath := filepath.Join(m.config.StorageDir, filename)
	
	// Save snapshot to file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create snapshot file: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(snapshot); err != nil {
		return "", fmt.Errorf("failed to encode snapshot: %w", err)
	}
	
	// Calculate file checksum
	checksum, err := m.calculateChecksum(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}
	
	// Create and save metadata
	metadata := &SnapshotMetadata{
		Name:        name,
		Description: description,
		FilePath:    filePath,
		Timestamp:   timestamp,
		Checksum:    checksum,
	}
	
	// Get file size
	if fileInfo, err := os.Stat(filePath); err == nil {
		metadata.Size = fileInfo.Size()
	}
	
	// Save metadata
	if err := m.saveMetadata(metadata); err != nil {
		return "", fmt.Errorf("failed to save metadata: %w", err)
	}
	
	// Perform cleanup if enabled
	if m.config.AutoCleanup {
		if err := m.cleanup(); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: cleanup failed: %v\n", err)
		}
	}
	
	return filePath, nil
}

// ListSnapshots returns all available snapshots with metadata
func (m *Manager) ListSnapshots() ([]*SnapshotMetadata, error) {
	metadataDir := filepath.Join(m.config.StorageDir, ".metadata")
	
	// Check if metadata directory exists
	if _, err := os.Stat(metadataDir); os.IsNotExist(err) {
		return []*SnapshotMetadata{}, nil
	}
	
	// Read metadata files
	entries, err := os.ReadDir(metadataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata directory: %w", err)
	}
	
	var snapshots []*SnapshotMetadata
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".json") {
			metadataPath := filepath.Join(metadataDir, entry.Name())
			metadata, err := m.loadMetadata(metadataPath)
			if err != nil {
				continue // Skip invalid metadata files
			}
			snapshots = append(snapshots, metadata)
		}
	}
	
	// Sort by timestamp (newest first)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})
	
	return snapshots, nil
}

// LoadSnapshot loads a snapshot from file
func (m *Manager) LoadSnapshot(filePath string) (*model.Snapshot, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}
	
	var snapshot model.Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}
	
	return &snapshot, nil
}

// DeleteSnapshot removes a snapshot and its metadata
func (m *Manager) DeleteSnapshot(name string) error {
	snapshots, err := m.ListSnapshots()
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}
	
	var targetSnapshot *SnapshotMetadata
	for _, snapshot := range snapshots {
		if snapshot.Name == name {
			targetSnapshot = snapshot
			break
		}
	}
	
	if targetSnapshot == nil {
		return fmt.Errorf("snapshot '%s' not found", name)
	}
	
	// Remove snapshot file
	if err := os.Remove(targetSnapshot.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove snapshot file: %w", err)
	}
	
	// Remove metadata file
	metadataPath := m.getMetadataPath(targetSnapshot.Name, targetSnapshot.Timestamp)
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove metadata file: %w", err)
	}
	
	return nil
}

// Helper methods

func (m *Manager) saveMetadata(metadata *SnapshotMetadata) error {
	metadataDir := filepath.Join(m.config.StorageDir, ".metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}
	
	metadataPath := m.getMetadataPath(metadata.Name, metadata.Timestamp)
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(metadata)
}

func (m *Manager) loadMetadata(filePath string) (*SnapshotMetadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}
	
	var metadata SnapshotMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	
	return &metadata, nil
}

func (m *Manager) getMetadataPath(name string, timestamp time.Time) string {
	safeFilename := sanitizeFilename(name)
	filename := fmt.Sprintf("%s_%s.json",
		timestamp.Format("2006-01-02_15-04-05"),
		safeFilename)
	return filepath.Join(m.config.StorageDir, ".metadata", filename)
}

func (m *Manager) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (m *Manager) cleanup() error {
	snapshots, err := m.ListSnapshots()
	if err != nil {
		return err
	}
	
	// Remove snapshots exceeding max count
	if len(snapshots) > m.config.MaxSnapshots {
		excess := snapshots[m.config.MaxSnapshots:]
		for _, snapshot := range excess {
			if err := m.DeleteSnapshot(snapshot.Name); err != nil {
				return fmt.Errorf("failed to delete excess snapshot %s: %w", snapshot.Name, err)
			}
		}
	}
	
	// Remove snapshots older than retention period
	cutoff := time.Now().Add(-m.config.RetentionPeriod)
	for _, snapshot := range snapshots {
		if snapshot.Timestamp.Before(cutoff) {
			if err := m.DeleteSnapshot(snapshot.Name); err != nil {
				return fmt.Errorf("failed to delete old snapshot %s: %w", snapshot.Name, err)
			}
		}
	}
	
	return nil
}

func sanitizeFilename(name string) string {
	// Replace invalid filename characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}