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
	RawData       []byte // Original SPD data
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
	if len(data) < 128 {
		return nil, fmt.Errorf("SPD data too short: %d bytes", len(data))
	}
	
	profile := &SPDProfile{
		RawData: make([]byte, len(data)),
	}
	copy(profile.RawData, data)
	
	// Parse basic device information (existing)
	if err := parseBasicInfo(data, profile); err != nil {
		return nil, fmt.Errorf("failed to parse basic info: %w", err)
	}
	
	// Parse capacity (call the existing function)
	profile.Capacity = parseCapacity(data)
	
	// Parse manufacturer, part number, serial number (call existing function)
	if err := parseManufacturerInfo(profile, data); err != nil {
		// Don't fail if manufacturer info is missing, just log
		// This is non-critical information
	}
	
	// Parse JEDEC standard profile (call the existing function)
	jedec, err := ParseJEDECProfile(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JEDEC profile: %w", err)
	}
	profile.JEDECProfile = jedec
	
	// Parse XMP profiles if present (call the existing function)
	xmpProfiles, err := ParseXMPProfiles(data)
	if err == nil {
		profile.XMPProfiles = xmpProfiles
	}
	
	// Parse EXPO profiles if present (call the existing function)
	expoProfiles, err := ParseEXPOProfiles(data)
	if err == nil {
		profile.EXPOProfiles = expoProfiles
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
func parseCapacity(data []byte) uint64 {
	if len(data) < 14 {
		return 0
	}
	
	deviceType := data[2]
	
	switch deviceType {
	case 0x0C: // DDR4
		return calculateDDR4Capacity(data)
	case 0x12: // DDR5
		return calculateDDR5Capacity(data)
	case 0x0B: // DDR3
		return calculateDDR3Capacity(data)
	default:
		return 0
	}
}

// calculateDDR4Capacity calculates capacity for DDR4 modules
func calculateDDR4Capacity(data []byte) uint64 {
	if len(data) < 14 {
		return 0
	}
	
	// DDR4 capacity calculation from SPD bytes 4-6
	sdramCapacity := uint64(256) << ((data[4] & 0x0F))  // MB per die
	primaryBusWidth := uint32(8) << ((data[13] & 0x07)) // bits
	sdramWidth := uint32(4) << ((data[12] & 0x07))      // bits
	logicalRanks := uint32(((data[12] >> 3) & 0x07) + 1)
	
	if primaryBusWidth == 0 || sdramWidth == 0 {
		return 0
	}
	
	// Capacity = (sdramCapacity / 8) * (primaryBusWidth / sdramWidth) * logicalRanks
	capacity := (sdramCapacity / 8) * uint64(primaryBusWidth/sdramWidth) * uint64(logicalRanks)
	return capacity * 1024 * 1024 // Convert to bytes
}

// calculateDDR5Capacity calculates capacity for DDR5 modules
func calculateDDR5Capacity(data []byte) uint64 {
	if len(data) < 14 {
		return 0
	}
	
	// DDR5 has similar structure but different encoding
	sdramCapacity := uint64(512) << ((data[4] & 0x1F))  // MB per die (5 bits for DDR5)
	primaryBusWidth := uint32(8) << ((data[13] & 0x07)) // bits
	sdramWidth := uint32(4) << ((data[12] & 0x07))      // bits  
	logicalRanks := uint32(((data[12] >> 3) & 0x07) + 1)
	
	if primaryBusWidth == 0 || sdramWidth == 0 {
		return 0
	}
	
	capacity := (sdramCapacity / 8) * uint64(primaryBusWidth/sdramWidth) * uint64(logicalRanks)
	return capacity * 1024 * 1024 // Convert to bytes
}

// calculateDDR3Capacity calculates capacity for DDR3 modules  
func calculateDDR3Capacity(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	
	// DDR3 capacity calculation (simplified)
	sdramCapacity := uint64(256) << ((data[4] & 0x0F))  // MB per die
	primaryBusWidth := uint32(64)                        // DDR3 typically 64-bit
	sdramWidth := uint32(8) << ((data[7] & 0x07))       // bits
	logicalRanks := uint32(((data[7] >> 3) & 0x07) + 1)
	
	if sdramWidth == 0 {
		return 0
	}
	
	capacity := (sdramCapacity / 8) * uint64(primaryBusWidth/sdramWidth) * uint64(logicalRanks)
	return capacity * 1024 * 1024 // Convert to bytes
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

// parseManufacturerInfo extracts manufacturer, part number, and serial from SPD data
func parseManufacturerInfo(profile *SPDProfile, data []byte) error {
	// Extract manufacturer and part number (bytes 320-383 for DDR4/DDR5)
	if len(data) >= 384 {
		// Manufacturer (bytes 320-321)
		manufacturerID := uint16(data[320]) | (uint16(data[321]) << 8)
		profile.Manufacturer = getManufacturer(manufacturerID)
		
		// Part number (bytes 329-348, ASCII)
		partBytes := data[329:349]
		for i, b := range partBytes {
			if b == 0 {
				profile.PartNumber = string(partBytes[:i])
				break
			}
		}
		if profile.PartNumber == "" {
			profile.PartNumber = string(partBytes)
		}
		
		// Serial number (bytes 325-328)
		serial := uint32(data[325]) | (uint32(data[326]) << 8) | 
		         (uint32(data[327]) << 16) | (uint32(data[328]) << 24)
		profile.SerialNumber = fmt.Sprintf("%08X", serial)
	}
	
	return nil
}

// getManufacturer returns manufacturer name from JEDEC ID
func getManufacturer(jedecID uint16) string {
	if name, exists := manufacturerNames[jedecID]; exists {
		return name
	}
	return fmt.Sprintf("Unknown (0x%04X)", jedecID)
}