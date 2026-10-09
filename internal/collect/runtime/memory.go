// internal/collect/runtime/memory.go
package runtime

import (
	"time"
)

// MemoryParameters represents current active memory configuration
type MemoryParameters struct {
	Timestamp        time.Time                `json:"timestamp"`
	CurrentFrequency uint32                   `json:"current_frequency_mhz"`
	ConfiguredSpeed  uint32                   `json:"configured_speed_mhz"`
	Timings          map[string]TimingValue   `json:"timings"`
	Voltages         map[string]VoltageValue  `json:"voltages"`
	Profile          *ActiveProfile           `json:"active_profile,omitempty"`
}

// TimingValue represents a memory timing parameter
type TimingValue struct {
	Name        string    `json:"name"`
	Value       uint16    `json:"value"`
	Unit        string    `json:"unit"`
	Description string    `json:"description,omitempty"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
}

// VoltageValue represents a voltage reading
type VoltageValue struct {
	Name        string    `json:"name"`
	Value       float32   `json:"value"`
	Unit        string    `json:"unit"`
	Description string    `json:"description,omitempty"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
}

// ActiveProfile represents the currently active memory profile
type ActiveProfile struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"` // JEDEC, XMP, EXPO, Custom
	Frequency   uint32    `json:"frequency_mhz"`
	Voltage     float32   `json:"voltage_v"`
	Timings     []string  `json:"timings"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
}