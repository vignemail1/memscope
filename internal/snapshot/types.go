package snapshot

import (
	"time"
)

// SnapshotMetadata contains information about a stored snapshot
type SnapshotMetadata struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	Timestamp   time.Time `json:"timestamp"`
	Size        int64     `json:"size_bytes"`
	Checksum    string    `json:"checksum,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

// ComparisonResult represents the result of comparing two snapshots
type ComparisonResult struct {
	Snapshot1    *SnapshotMetadata `json:"snapshot1"`
	Snapshot2    *SnapshotMetadata `json:"snapshot2"`
	HasChanges   bool              `json:"has_changes"`
	Changes      []*Change         `json:"changes"`
	Summary      *ComparisonSummary `json:"summary"`
	GeneratedAt  time.Time         `json:"generated_at"`
}

// Change represents a single change detected between snapshots
type Change struct {
	Category    string      `json:"category"`    // system, memory, device, observation
	Type        string      `json:"type"`        // added, removed, modified
	Path        string      `json:"path"`        // dot notation path to changed field
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	Severity    string      `json:"severity"`    // low, medium, high, critical
	Description string      `json:"description"`
}

// ComparisonSummary provides aggregated change statistics
type ComparisonSummary struct {
	TotalChanges     int            `json:"total_changes"`
	ByCategory       map[string]int `json:"by_category"`
	BySeverity       map[string]int `json:"by_severity"`
	ByType           map[string]int `json:"by_type"`
	SystemChanges    int            `json:"system_changes"`
	MemoryChanges    int            `json:"memory_changes"`
	DeviceChanges    int            `json:"device_changes"`
	ObservationChanges int          `json:"observation_changes"`
}

// SnapshotConfig contains configuration for snapshot management
type SnapshotConfig struct {
	StorageDir       string        `json:"storage_dir"`
	MaxSnapshots     int           `json:"max_snapshots"`
	RetentionPeriod  time.Duration `json:"retention_period"`
	AutoCleanup      bool          `json:"auto_cleanup"`
	CompressionLevel int           `json:"compression_level"`
	IncludeRuntime   bool          `json:"include_runtime"`
}

// ValidationResult contains snapshot validation results
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// ComparisonOptions configures comparison behavior
type ComparisonOptions struct {
	IgnoreTimestamps     bool     `json:"ignore_timestamps"`
	IgnoreCategories     []string `json:"ignore_categories,omitempty"`
	ShowOnlyChanges      bool     `json:"show_only_changes"`
	IncludeRuntimeData   bool     `json:"include_runtime_data"`
	SeverityFilter       []string `json:"severity_filter,omitempty"`
}