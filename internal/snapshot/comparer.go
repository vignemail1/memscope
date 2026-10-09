package snapshot

import (
	"fmt"
	"reflect"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

// Comparer handles snapshot comparison operations
type Comparer struct {
	options *ComparisonOptions
}

// NewComparer creates a new comparer with default options
func NewComparer() *Comparer {
	return &Comparer{
		options: &ComparisonOptions{
			IgnoreTimestamps:   true, // By default, ignore timestamp differences
			ShowOnlyChanges:    false,
			IncludeRuntimeData: true,
		},
	}
}

// NewComparerWithOptions creates a comparer with custom options
func NewComparerWithOptions(options *ComparisonOptions) *Comparer {
	return &Comparer{options: options}
}

// CompareSnapshots compares two snapshots and returns detailed differences
func (c *Comparer) CompareSnapshots(snapshot1, snapshot2 *model.Snapshot) (*ComparisonResult, error) {
	if snapshot1 == nil || snapshot2 == nil {
		return nil, fmt.Errorf("snapshots cannot be nil")
	}
	
	result := &ComparisonResult{
		Snapshot1: &SnapshotMetadata{
			Name:      "snapshot1",
			Timestamp: snapshot1.CollectionStartedAt,
		},
		Snapshot2: &SnapshotMetadata{
			Name:      "snapshot2", 
			Timestamp: snapshot2.CollectionStartedAt,
		},
		Changes:     []*Change{},
		GeneratedAt: time.Now(),
	}
	
	// Compare platform information
	platformChanges := c.comparePlatform(snapshot1.Platform, snapshot2.Platform)
	result.Changes = append(result.Changes, platformChanges...)
	
	// Compare tool version and schema version (system-level changes)
	systemChanges := c.compareSystemInfo(snapshot1, snapshot2)
	result.Changes = append(result.Changes, systemChanges...)
	
	// Compare devices
	deviceChanges := c.compareDevices(snapshot1.Devices, snapshot2.Devices)
	result.Changes = append(result.Changes, deviceChanges...)
	
	// Compare observations (runtime data)
	if c.options.IncludeRuntimeData {
		observationChanges := c.compareObservations(snapshot1.Observations, snapshot2.Observations)
		result.Changes = append(result.Changes, observationChanges...)
	}
	
	// Generate summary
	result.HasChanges = len(result.Changes) > 0
	result.Summary = c.generateSummary(result.Changes)
	
	return result, nil
}

// compareSystemInfo compares system-level information
func (c *Comparer) compareSystemInfo(snapshot1, snapshot2 *model.Snapshot) []*Change {
	var changes []*Change
	
	// Compare tool version
	if snapshot1.ToolVersion != snapshot2.ToolVersion {
		changes = append(changes, &Change{
			Category:    "system",
			Type:        "modified",
			Path:        "tool_version",
			OldValue:    snapshot1.ToolVersion,
			NewValue:    snapshot2.ToolVersion,
			Severity:    "low",
			Description: "Tool version changed",
		})
	}
	
	// Compare schema version
	if snapshot1.SchemaVersion != snapshot2.SchemaVersion {
		changes = append(changes, &Change{
			Category:    "system",
			Type:        "modified",
			Path:        "schema_version",
			OldValue:    snapshot1.SchemaVersion,
			NewValue:    snapshot2.SchemaVersion,
			Severity:    "medium",
			Description: "Schema version changed",
		})
	}
	
	return changes
}

// comparePlatform compares platform information between snapshots
func (c *Comparer) comparePlatform(platform1, platform2 model.Platform) []*Change {
	var changes []*Change
	
	// Compare OS
	if platform1.OS != platform2.OS {
		changes = append(changes, &Change{
			Category:    "system",
			Type:        "modified",
			Path:        "platform.os",
			OldValue:    platform1.OS,
			NewValue:    platform2.OS,
			Severity:    "high",
			Description: "Operating system changed",
		})
	}
	
	// Compare architecture
	if platform1.Arch != platform2.Arch {
		changes = append(changes, &Change{
			Category:    "system",
			Type:        "modified",
			Path:        "platform.arch",
			OldValue:    platform1.Arch,
			NewValue:    platform2.Arch,
			Severity:    "high",
			Description: "System architecture changed",
		})
	}
	
	return changes
}

// compareDevices compares device lists between snapshots
func (c *Comparer) compareDevices(devices1, devices2 []model.Device) []*Change {
	var changes []*Change
	
	// Create maps for easier comparison
	deviceMap1 := make(map[string]model.Device)
	deviceMap2 := make(map[string]model.Device)
	
	for _, device := range devices1 {
		deviceMap1[device.ID] = device
	}
	
	for _, device := range devices2 {
		deviceMap2[device.ID] = device
	}
	
	// Find added devices
	for id, device2 := range deviceMap2 {
		if _, exists := deviceMap1[id]; !exists {
			changes = append(changes, &Change{
				Category:    "device",
				Type:        "added",
				Path:        fmt.Sprintf("devices.%s", id),
				OldValue:    nil,
				NewValue:    device2,
				Severity:    c.getDeviceSeverity(device2, "added"),
				Description: fmt.Sprintf("%s device added: %s", device2.Kind, id),
			})
		}
	}
	
	// Find removed devices
	for id, device1 := range deviceMap1 {
		if _, exists := deviceMap2[id]; !exists {
			changes = append(changes, &Change{
				Category:    "device",
				Type:        "removed",
				Path:        fmt.Sprintf("devices.%s", id),
				OldValue:    device1,
				NewValue:    nil,
				Severity:    c.getDeviceSeverity(device1, "removed"),
				Description: fmt.Sprintf("%s device removed: %s", device1.Kind, id),
			})
		}
	}
	
	// Find modified devices
	for id, device1 := range deviceMap1 {
		if device2, exists := deviceMap2[id]; exists {
			deviceChanges := c.compareDeviceProperties(device1, device2)
			changes = append(changes, deviceChanges...)
		}
	}
	
	return changes
}

// compareDeviceProperties compares properties of two devices
func (c *Comparer) compareDeviceProperties(device1, device2 model.Device) []*Change {
	var changes []*Change
	
	// Compare device kind
	if device1.Kind != device2.Kind {
		changes = append(changes, &Change{
			Category:    "device",
			Type:        "modified",
			Path:        fmt.Sprintf("devices.%s.kind", device1.ID),
			OldValue:    device1.Kind,
			NewValue:    device2.Kind,
			Severity:    "high",
			Description: fmt.Sprintf("Device kind changed from %s to %s", device1.Kind, device2.Kind),
		})
	}
	
	// Compare parent ID
	if device1.ParentID != device2.ParentID {
		changes = append(changes, &Change{
			Category:    "device",
			Type:        "modified",
			Path:        fmt.Sprintf("devices.%s.parent_id", device1.ID),
			OldValue:    device1.ParentID,
			NewValue:    device2.ParentID,
			Severity:    "medium",
			Description: fmt.Sprintf("Device parent changed from %s to %s", device1.ParentID, device2.ParentID),
		})
	}
	
	return changes
}

// compareObservations compares observation data between snapshots
func (c *Comparer) compareObservations(obs1, obs2 []model.Observation) []*Change {
	var changes []*Change
	
	// Create maps for easier comparison
	obsMap1 := make(map[string]model.Observation)
	obsMap2 := make(map[string]model.Observation)
	
	for _, obs := range obs1 {
		key := fmt.Sprintf("%s_%s_%s", obs.DeviceID, obs.Scope, obs.Parameter)
		obsMap1[key] = obs
	}
	
	for _, obs := range obs2 {
		key := fmt.Sprintf("%s_%s_%s", obs.DeviceID, obs.Scope, obs.Parameter)
		obsMap2[key] = obs
	}
	
	// Find added observations
	for key, obs2 := range obsMap2 {
		if _, exists := obsMap1[key]; !exists {
			changes = append(changes, &Change{
				Category:    "observation",
				Type:        "added",
				Path:        fmt.Sprintf("observations.%s", key),
				OldValue:    nil,
				NewValue:    obs2,
				Severity:    c.getObservationSeverity(obs2.Parameter),
				Description: fmt.Sprintf("Observation added: %s", obs2.Parameter),
			})
		}
	}
	
	// Find removed observations
	for key, obs1 := range obsMap1 {
		if _, exists := obsMap2[key]; !exists {
			changes = append(changes, &Change{
				Category:    "observation",
				Type:        "removed",
				Path:        fmt.Sprintf("observations.%s", key),
				OldValue:    obs1,
				NewValue:    nil,
				Severity:    c.getObservationSeverity(obs1.Parameter),
				Description: fmt.Sprintf("Observation removed: %s", obs1.Parameter),
			})
		}
	}
	
	// Find modified observations
	for key, obs1 := range obsMap1 {
		if obs2, exists := obsMap2[key]; exists {
			// Compare values (ignoring timestamps if configured)
			if !c.options.IgnoreTimestamps && !obs1.CapturedAt.Equal(obs2.CapturedAt) {
				changes = append(changes, &Change{
					Category:    "observation",
					Type:        "modified",
					Path:        fmt.Sprintf("observations.%s.captured_at", key),
					OldValue:    obs1.CapturedAt,
					NewValue:    obs2.CapturedAt,
					Severity:    "low",
					Description: fmt.Sprintf("Observation timestamp changed: %s", obs1.Parameter),
				})
			}
			
			// Compare observation values
			if !reflect.DeepEqual(obs1.Value, obs2.Value) {
				changes = append(changes, &Change{
					Category:    "observation",
					Type:        "modified",
					Path:        fmt.Sprintf("observations.%s.value", key),
					OldValue:    obs1.Value,
					NewValue:    obs2.Value,
					Severity:    c.getObservationSeverity(obs1.Parameter),
					Description: fmt.Sprintf("Observation value changed: %s", obs1.Parameter),
				})
			}
			
			// Compare status
			if obs1.Status != obs2.Status {
				changes = append(changes, &Change{
					Category:    "observation",
					Type:        "modified",
					Path:        fmt.Sprintf("observations.%s.status", key),
					OldValue:    obs1.Status,
					NewValue:    obs2.Status,
					Severity:    "medium",
					Description: fmt.Sprintf("Observation status changed: %s", obs1.Parameter),
				})
			}
		}
	}
	
	return changes
}

// generateSummary creates a summary of all changes
func (c *Comparer) generateSummary(changes []*Change) *ComparisonSummary {
	summary := &ComparisonSummary{
		TotalChanges: len(changes),
		ByCategory:   make(map[string]int),
		BySeverity:   make(map[string]int),
		ByType:       make(map[string]int),
	}
	
	for _, change := range changes {
		summary.ByCategory[change.Category]++
		summary.BySeverity[change.Severity]++
		summary.ByType[change.Type]++
		
		// Count specific categories
		switch change.Category {
		case "system":
			summary.SystemChanges++
		case "memory":
			summary.MemoryChanges++
		case "device":
			summary.DeviceChanges++
		case "observation":
			summary.ObservationChanges++
		}
	}
	
	return summary
}

// Helper methods for categorization and severity assessment

func (c *Comparer) getDeviceSeverity(device model.Device, changeType string) string {
	if device.Kind == "memory" {
		return "high"
	}
	return "medium"
}

func (c *Comparer) getObservationSeverity(parameter string) string {
	// Memory-related observations are typically more critical
	memoryParams := []string{"frequency", "voltage", "timing", "latency", "bandwidth"}
	for _, param := range memoryParams {
		if parameter == param {
			return "high"
		}
	}
	return "medium"
}