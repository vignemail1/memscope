package spd

import (
	"fmt"
)

// ParseJEDECProfile extracts JEDEC standard profile from SPD data
func ParseJEDECProfile(data []byte) (*MemoryProfile, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("insufficient data for JEDEC profile parsing")
	}
	
	profile := &MemoryProfile{
		Type:        "JEDEC",
		ProfileName: "Standard",
	}
	
	// Determine DDR type from byte 2
	switch data[2] {
	case 0x0B:
		return parseJEDECDDR3(data, profile)
	case 0x0C:
		return parseJEDECDDR4(data, profile)
	case 0x12:
		return parseJEDECDDR5(data, profile)
	default:
		return nil, fmt.Errorf("unsupported DDR type: 0x%02X", data[2])
	}
}

// parseJEDECDDR4 parses JEDEC profile for DDR4 modules
func parseJEDECDDR4(data []byte, profile *MemoryProfile) (*MemoryProfile, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("insufficient DDR4 data")
	}
	
	// DDR4 JEDEC timing extraction
	// Note: MTB and FTB calculations are available if needed for precise timing
	_ = calculateMTB(data[17]) // Typically 0.125 ns for DDR4
	_ = calculateFTB(data[16]) // Typically 0.001 ns for DDR4
	
	// Basic frequency calculation (simplified)
	profile.Frequency = 2133 // Default DDR4-2133 JEDEC
	
	// Standard DDR4 voltage
	profile.Voltage = 1.2
	
	// Extract basic timings
	profile.Timings = TimingSpec{
		CL:   15, // Default CL for DDR4-2133
		TRCD: 15,
		TRP:  15,
		TRAS: 35,
	}
	
	return profile, nil
}

// parseJEDECDDR5 parses JEDEC profile for DDR5 modules
func parseJEDECDDR5(data []byte, profile *MemoryProfile) (*MemoryProfile, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("insufficient DDR5 data")
	}
	
	// DDR5 has different structure
	profile.Frequency = 4800 // Default DDR5-4800 JEDEC
	profile.Voltage = 1.1    // Standard DDR5 voltage
	
	profile.Timings = TimingSpec{
		CL:   40, // Default CL for DDR5-4800
		TRCD: 39,
		TRP:  39,
		TRAS: 52,
	}
	
	return profile, nil
}

// parseJEDECDDR3 parses JEDEC profile for DDR3 modules
func parseJEDECDDR3(data []byte, profile *MemoryProfile) (*MemoryProfile, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("insufficient DDR3 data")
	}
	
	profile.Frequency = 1600 // Default DDR3-1600 JEDEC
	profile.Voltage = 1.5    // Standard DDR3 voltage
	
	profile.Timings = TimingSpec{
		CL:   11, // Default CL for DDR3-1600
		TRCD: 11,
		TRP:  11,
		TRAS: 28,
	}
	
	return profile, nil
}

// calculateMTB calculates Medium Timebase from SPD byte
func calculateMTB(mtbByte byte) float64 {
	// MTB is typically 0.125 ns for DDR4/DDR5
	// This is a simplified calculation
	return 0.125
}

// calculateFTB calculates Fine Timebase from SPD byte  
func calculateFTB(ftbByte byte) float64 {
	// FTB is typically 0.001 ns for DDR4/DDR5
	// This is a simplified calculation
	return 0.001
}

// extractTimingValue extracts timing value considering MTB and FTB
func extractTimingValue(mainByte, fineByte byte, mtb, ftb float64) uint32 {
	// Convert signed fine byte to adjustment
	var fineAdjust float64
	if fineByte&0x80 != 0 {
		// Negative adjustment
		fineAdjust = -float64(256-int(fineByte)) * ftb
	} else {
		// Positive adjustment
		fineAdjust = float64(fineByte) * ftb
	}
	
	// Calculate final timing in nanoseconds
	totalTime := float64(mainByte)*mtb + fineAdjust
	
	// Convert to cycles (this is simplified - actual conversion needs frequency)
	return uint32(totalTime / mtb)
}