package snapshot

import (
	"fmt"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

// Validator handles snapshot validation
type Validator struct{}

// NewValidator creates a new snapshot validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateSnapshot validates a snapshot for completeness and consistency
func (v *Validator) ValidateSnapshot(snapshot *model.Snapshot) error {
	result := v.ValidateSnapshotDetailed(snapshot)
	if !result.Valid {
		if len(result.Errors) > 0 {
			return fmt.Errorf("snapshot validation failed: %s", result.Errors[0])
		}
		return fmt.Errorf("snapshot validation failed")
	}
	return nil
}

// ValidateSnapshotDetailed performs detailed validation and returns full results
func (v *Validator) ValidateSnapshotDetailed(snapshot *model.Snapshot) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}
	
	if snapshot == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "snapshot is nil")
		return result
	}
	
	// Use the existing model validation but catch any additional issues
	if err := snapshot.Validate(); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("model validation failed: %v", err))
	}
	
	// Additional validation beyond the model's validation
	
	// Check if timestamps are in the future (warning, not error)
	now := time.Now()
	if snapshot.CollectionStartedAt.After(now.Add(time.Hour)) {
		result.Warnings = append(result.Warnings, "collection start timestamp is more than 1 hour in the future")
	}
	
	if snapshot.CollectionFinishedAt.After(now.Add(time.Hour)) {
		result.Warnings = append(result.Warnings, "collection end timestamp is more than 1 hour in the future")
	}
	
	// Validate device distribution
	v.validateDeviceDistribution(snapshot, result)
	
	// Validate observation distribution
	v.validateObservationDistribution(snapshot, result)
	
	return result
}

// validateDeviceDistribution checks for reasonable device counts and types
func (v *Validator) validateDeviceDistribution(snapshot *model.Snapshot, result *ValidationResult) {
	if len(snapshot.Devices) == 0 {
		result.Warnings = append(result.Warnings, "no devices found in snapshot")
		return
	}
	
	deviceKinds := make(map[string]int)
	for _, device := range snapshot.Devices {
		deviceKinds[device.Kind]++
	}
	
	// Check for memory devices
	if deviceKinds["memory"] == 0 {
		result.Warnings = append(result.Warnings, "no memory devices found in snapshot")
	}
	
	// Warn about unusual device counts
	if deviceKinds["memory"] > 32 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("unusually high number of memory devices: %d", deviceKinds["memory"]))
	}
	
	// Check for device hierarchy issues
	v.validateDeviceHierarchy(snapshot.Devices, result)
}

// validateDeviceHierarchy checks for parent-child relationship consistency
func (v *Validator) validateDeviceHierarchy(devices []model.Device, result *ValidationResult) {
	deviceMap := make(map[string]model.Device)
	for _, device := range devices {
		deviceMap[device.ID] = device
	}
	
	// Check for orphaned devices (parent doesn't exist)
	for _, device := range devices {
		if device.ParentID != "" {
			if parent, exists := deviceMap[device.ParentID]; !exists {
				result.Warnings = append(result.Warnings, fmt.Sprintf("device %s references non-existent parent %s", device.ID, device.ParentID))
			} else {
				// Warn about unusual parent-child relationships
				if device.Kind == "memory" && parent.Kind != "motherboard" && parent.Kind != "system" {
					result.Warnings = append(result.Warnings, fmt.Sprintf("memory device %s has unusual parent kind: %s", device.ID, parent.Kind))
				}
			}
		}
	}
}

// validateObservationDistribution checks observation patterns
func (v *Validator) validateObservationDistribution(snapshot *model.Snapshot, result *ValidationResult) {
	if len(snapshot.Observations) == 0 {
		result.Warnings = append(result.Warnings, "no observations found in snapshot")
		return
	}
	
	// Count observations by device
	deviceObservations := make(map[string]int)
	observationScopes := make(map[string]int)
	observationStatuses := make(map[string]int)
	
	for _, obs := range snapshot.Observations {
		deviceObservations[obs.DeviceID]++
		observationScopes[obs.Scope]++
		observationStatuses[string(obs.Status)]++
	}
	
	// Warn about devices with no observations
	deviceMap := make(map[string]bool)
	for _, device := range snapshot.Devices {
		deviceMap[device.ID] = true
	}
	
	for deviceID := range deviceMap {
		if deviceObservations[deviceID] == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("device %s has no observations", deviceID))
		}
	}
	
	// Warn about unusual observation status distributions
	totalObservations := len(snapshot.Observations)
	if observationStatuses["unavailable"] > totalObservations/2 {
		result.Warnings = append(result.Warnings, "more than half of observations are unavailable")
	}
	
	if observationStatuses["invalid"] > totalObservations/10 {
		result.Warnings = append(result.Warnings, "more than 10% of observations are invalid")
	}
	
	// Check for temporal consistency
	v.validateObservationTimestamps(snapshot, result)
}

// validateObservationTimestamps checks that observation timestamps are within collection period
func (v *Validator) validateObservationTimestamps(snapshot *model.Snapshot, result *ValidationResult) {
	start := snapshot.CollectionStartedAt
	end := snapshot.CollectionFinishedAt
	
	for i, obs := range snapshot.Observations {
		if obs.CapturedAt.Before(start) || obs.CapturedAt.After(end) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("observation %d timestamp (%v) is outside collection period (%v - %v)", i, obs.CapturedAt, start, end))
		}
	}
}