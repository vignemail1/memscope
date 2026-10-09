//go:build !windows

package windows

import "errors"

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

// GetMemoryControllerInfo stub for non-Windows platforms
func GetMemoryControllerInfo() (*MemoryControllerInfo, error) {
	return nil, errors.New("memory controller access not available on this platform")
}

// GetCurrentMemoryTimings stub for non-Windows platforms
func GetCurrentMemoryTimings() (map[string]MemoryTiming, error) {
	return nil, errors.New("memory timing access not available on this platform")
}

// GetMemoryVoltages stub for non-Windows platforms
func GetMemoryVoltages() (map[string]MemoryVoltage, error) {
	return nil, errors.New("memory voltage access not available on this platform")
}

// DetectActiveProfile stub for non-Windows platforms
func DetectActiveProfile(controller *MemoryControllerInfo, timings map[string]MemoryTiming) (*ActiveProfileInfo, error) {
	return nil, errors.New("profile detection not available on this platform")
}