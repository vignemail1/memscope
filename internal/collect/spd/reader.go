// Package spd implements SPD (Serial Presence Detect) data parsing for memory modules.
package spd

import (
	"fmt"
)

// SPDProfile represents parsed SPD data for a memory module
type SPDProfile struct {
	DeviceType    string
	ModuleType    string
	Capacity      uint64 // in bytes
	Manufacturer  string
	PartNumber    string
	SerialNumber  string
	JEDECProfile  *MemoryProfile
	XMPProfiles   []MemoryProfile
	EXPOProfiles  []MemoryProfile
}

// MemoryProfile represents timing and voltage specifications
type MemoryProfile struct {
	Type        string  // "JEDEC", "XMP", "EXPO"
	ProfileName string
	Frequency   uint32  // in MHz
	Voltage     float64 // in volts
	Timings     TimingSpec
}

// TimingSpec holds memory timing parameters
type TimingSpec struct {
	CL   uint32 // CAS Latency
	TRCD uint32 // RAS to CAS Delay
	TRP  uint32 // Row Precharge
	TRAS uint32 // Row Active Strobe
}

// ParseSPDData parses raw SPD data and returns a structured profile
func ParseSPDData(data []byte) (*SPDProfile, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("SPD data too short: got %d bytes, need at least 4", len(data))
	}
	
	profile := &SPDProfile{}
	
	// Parse basic device information
	err := parseBasicInfo(data, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse basic info: %w", err)
	}
	
	return profile, nil
}

// parseBasicInfo extracts basic module information from SPD data
func parseBasicInfo(data []byte, profile *SPDProfile) error {
	// Byte 2: DRAM Device Type
	switch data[2] {
	case 0x0B:
		profile.DeviceType = "DDR3"
	case 0x0C:
		profile.DeviceType = "DDR4"
	case 0x12:
		profile.DeviceType = "DDR5"
	default:
		profile.DeviceType = "Unknown"
	}
	
	// Byte 3: Module Type
	switch data[3] & 0x0F {
	case 0x01:
		profile.ModuleType = "UDIMM"
	case 0x02:
		profile.ModuleType = "SODIMM"
	case 0x03:
		profile.ModuleType = "RDIMM"
	case 0x04:
		profile.ModuleType = "LRDIMM"
	default:
		profile.ModuleType = "Unknown"
	}
	
	return nil
}

// parseCapacity calculates module capacity from SPD data
func parseCapacity(data []byte) (uint64, error) {
	if len(data) < 8 {
		return 0, fmt.Errorf("insufficient data for capacity calculation")
	}
	
	// Basic capacity calculation for DDR4/DDR5
	// This is a simplified version - actual calculation is more complex
	capacityMB := uint64(256) // Default 256MB as baseline
	return capacityMB * 1024 * 1024, nil // Convert to bytes
}

// manufacturerNames maps JEDEC manufacturer IDs to company names
var manufacturerNames = map[uint16]string{
	0x80AD: "SK Hynix",
	0x80CE: "Samsung",
	0x802C: "Micron",
	0x8551: "ADATA",
	0x859B: "Crucial",
	0x04CD: "G.Skill",
	0x029E: "Corsair",
	0x8304: "Kingston",
}

// getManufacturer returns manufacturer name from JEDEC ID
func getManufacturer(jedecID uint16) string {
	if name, exists := manufacturerNames[jedecID]; exists {
		return name
	}
	return fmt.Sprintf("Unknown (0x%04X)", jedecID)
}