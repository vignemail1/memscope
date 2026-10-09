// internal/collect/runtime/windows.go
package runtime

import (
	"fmt"
	"time"
	
	"github.com/vignemail1/memscope/internal/platform/windows"
)

// WindowsProvider implements runtime collection for Windows
type WindowsProvider struct{}

// NewWindowsProvider creates a new Windows runtime provider
func NewWindowsProvider() *WindowsProvider {
	return &WindowsProvider{}
}

// CollectMemoryParameters gathers current memory parameters from Windows
func (p *WindowsProvider) CollectMemoryParameters() (*MemoryParameters, error) {
	timestamp := time.Now()
	
	// Collect memory controller information
	controllerInfo, err := windows.GetMemoryControllerInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory controller info: %w", err)
	}
	
	// Collect current memory timings
	timings, err := windows.GetCurrentMemoryTimings()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory timings: %w", err)
	}
	
	// Collect voltage readings
	voltages, err := windows.GetMemoryVoltages()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory voltages: %w", err)
	}
	
	// Convert to runtime parameters format
	params := &MemoryParameters{
		Timestamp:        timestamp,
		CurrentFrequency: controllerInfo.CurrentFrequency,
		ConfiguredSpeed:  controllerInfo.ConfiguredSpeed,
		Timings:          make(map[string]TimingValue),
		Voltages:         make(map[string]VoltageValue),
	}
	
	// Convert timings
	for name, timing := range timings {
		params.Timings[name] = TimingValue{
			Name:      name,
			Value:     timing.Value,
			Unit:      timing.Unit,
			Source:    "memory_controller",
			Timestamp: timestamp,
		}
	}
	
	// Convert voltages
	for name, voltage := range voltages {
		params.Voltages[name] = VoltageValue{
			Name:      name,
			Value:     voltage.Value,
			Unit:      voltage.Unit,
			Source:    "sensor",
			Timestamp: timestamp,
		}
	}
	
	// Detect active profile if possible
	if profile, err := windows.DetectActiveProfile(controllerInfo, timings); err == nil {
		params.Profile = &ActiveProfile{
			Name:      profile.Name,
			Type:      profile.Type,
			Frequency: profile.Frequency,
			Voltage:   profile.Voltage,
			Source:    "detection",
			Timestamp: timestamp,
		}
	}
	
	return params, nil
}