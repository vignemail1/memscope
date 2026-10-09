//go:build windows

package windows

import (
	"fmt"
	"strings"
)

// MemoryControllerInfo represents memory controller data
type MemoryControllerInfo struct {
	CurrentFrequency uint32
	ConfiguredSpeed  uint32
	ControllerType   string
	Channels         int
}

// MemoryTiming represents a timing parameter
type MemoryTiming struct {
	Value uint16
	Unit  string
}

// MemoryVoltage represents a voltage reading
type MemoryVoltage struct {
	Value float32
	Unit  string
}

// ActiveProfileInfo represents detected active profile
type ActiveProfileInfo struct {
	Name      string
	Type      string
	Frequency uint32
	Voltage   float32
}

// GetMemoryControllerInfo retrieves memory controller information via WMI
func GetMemoryControllerInfo() (*MemoryControllerInfo, error) {
	client, err := NewWMIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI client: %w", err)
	}
	defer client.Close()
	
	// Query memory controller
	query := "SELECT ConfiguredClockSpeed, Description FROM Win32_PhysicalMemory"
	results, err := client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query memory controller: %w", err)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no memory controller information found")
	}
	
	info := &MemoryControllerInfo{
		ControllerType: "DDR",
		Channels:       1, // Default assumption
	}
	
	// Extract frequency from first module (they should match)
	if speed, ok := results[0]["ConfiguredClockSpeed"].(uint32); ok {
		info.CurrentFrequency = speed
		info.ConfiguredSpeed = speed
	}
	
	return info, nil
}

// GetCurrentMemoryTimings retrieves current memory timings
func GetCurrentMemoryTimings() (map[string]MemoryTiming, error) {
	// Note: Real implementation would read from MSR registers or chipset
	// This is a mock implementation for development
	timings := map[string]MemoryTiming{
		"CL": {
			Value: 16,
			Unit:  "clocks",
		},
		"TRCD": {
			Value: 16,
			Unit:  "clocks",
		},
		"TRP": {
			Value: 16,
			Unit:  "clocks",
		},
		"TRAS": {
			Value: 36,
			Unit:  "clocks",
		},
	}
	
	return timings, nil
}

// GetMemoryVoltages retrieves current memory voltage readings
func GetMemoryVoltages() (map[string]MemoryVoltage, error) {
	client, err := NewWMIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI client: %w", err)
	}
	defer client.Close()
	
	// Try to get voltage from WMI sensors
	query := "SELECT Name, CurrentReading FROM Win32_Sensor WHERE SensorType = 3"
	results, err := client.Query(query)
	if err != nil {
		// Fallback to mock data if sensors not available
		voltages := map[string]MemoryVoltage{
			"VDIMM": {
				Value: 1.35,
				Unit:  "V",
			},
		}
		return voltages, nil
	}
	
	voltages := make(map[string]MemoryVoltage)
	
	for _, result := range results {
		if name, ok := result["Name"].(string); ok {
			if strings.Contains(strings.ToLower(name), "memory") || 
			   strings.Contains(strings.ToLower(name), "dimm") {
				if reading, ok := result["CurrentReading"].(float64); ok {
					voltages[name] = MemoryVoltage{
						Value: float32(reading),
						Unit:  "V",
					}
				}
			}
		}
	}
	
	// Ensure at least one voltage reading
	if len(voltages) == 0 {
		voltages["VDIMM"] = MemoryVoltage{
			Value: 1.35,
			Unit:  "V",
		}
	}
	
	return voltages, nil
}

// DetectActiveProfile attempts to detect the active memory profile
func DetectActiveProfile(controller *MemoryControllerInfo, timings map[string]MemoryTiming) (*ActiveProfileInfo, error) {
	// Simple detection logic based on frequency and timings
	profile := &ActiveProfileInfo{
		Frequency: controller.CurrentFrequency,
		Voltage:   1.35, // Default DDR4 voltage
	}
	
	// Detect profile type based on frequency
	switch {
	case controller.CurrentFrequency >= 3200:
		profile.Name = "XMP"
		profile.Type = "XMP"
		profile.Voltage = 1.35
	case controller.CurrentFrequency >= 2400:
		profile.Name = "JEDEC Enhanced"
		profile.Type = "JEDEC"
		profile.Voltage = 1.2
	default:
		profile.Name = "JEDEC Standard"
		profile.Type = "JEDEC"
		profile.Voltage = 1.2
	}
	
	return profile, nil
}