package spd

import (
	"fmt"
)

// EXPO constants for AMD Extended Profiles for Overclocking
const (
	EXPOSignature   = "EXPO"     // ASCII signature for EXPO
	EXPOOffset      = 480        // Byte offset for EXPO data in SPD
	EXPOProfileSize = 16         // Size of each EXPO profile in bytes
)

// ParseEXPOProfiles extracts AMD EXPO profiles from SPD data
func ParseEXPOProfiles(data []byte) ([]MemoryProfile, error) {
	if len(data) < 512 {
		return nil, fmt.Errorf("insufficient data for EXPO parsing: need 512 bytes, got %d", len(data))
	}

	// EXPO is primarily for DDR5
	if data[2] != 0x12 {
		return nil, fmt.Errorf("EXPO profiles only supported on DDR5, got DDR type: 0x%02X", data[2])
	}

	if len(data) < EXPOOffset+32 {
		return nil, fmt.Errorf("insufficient data for EXPO parsing at offset %d", EXPOOffset)
	}

	expoData := data[EXPOOffset:]

	// Verify EXPO signature
	if !verifyEXPOSignature(expoData) {
		return nil, fmt.Errorf("EXPO signature not found")
	}

	var profiles []MemoryProfile

	// Check which profiles are present
	// EXPO typically has up to 2 profiles
	for i := 0; i < 2; i++ {
		profile, err := parseEXPOProfile(expoData, i+1)
		if err == nil && profile != nil {
			profiles = append(profiles, *profile)
		}
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no valid EXPO profiles found")
	}

	return profiles, nil
}

// parseEXPOProfile parses a specific EXPO profile
func parseEXPOProfile(expoData []byte, profileNum int) (*MemoryProfile, error) {
	// Calculate offset for the specific profile
	profileOffset := 4 + (profileNum-1)*EXPOProfileSize // Skip 4-byte signature

	if len(expoData) < profileOffset+EXPOProfileSize {
		return nil, fmt.Errorf("insufficient data for EXPO profile %d", profileNum)
	}

	profileData := expoData[profileOffset:]

	// Check if profile is valid (first byte should not be 0x00)
	if profileData[0] == 0x00 {
		return nil, fmt.Errorf("EXPO profile %d is empty", profileNum)
	}

	profile := &MemoryProfile{
		Type:        "EXPO",
		ProfileName: fmt.Sprintf("EXPO Profile %d", profileNum),
	}

	// Extract frequency
	profile.Frequency = extractEXPOFrequency(profileData)

	// Extract voltage
	profile.Voltage = extractEXPOVoltage(profileData)

	// Extract timings
	profile.Timings = extractEXPOTimings(profileData)

	// Validate profile makes sense
	if err := validateEXPOProfile(profile); err != nil {
		return nil, fmt.Errorf("invalid EXPO profile %d: %w", profileNum, err)
	}

	return profile, nil
}

// verifyEXPOSignature checks if EXPO signature is present
func verifyEXPOSignature(expoData []byte) bool {
	if len(expoData) < 4 {
		return false
	}

	// Check for ASCII "EXPO" signature
	signature := string(expoData[0:4])
	return signature == EXPOSignature
}

// extractEXPOFrequency calculates frequency from EXPO profile data
func extractEXPOFrequency(profileData []byte) uint32 {
	if len(profileData) < 4 {
		return 5200 // Default DDR5 EXPO frequency
	}

	// EXPO frequency encoding (simplified)
	// Real implementation would decode the specific frequency bytes
	// For now, return common EXPO frequencies based on profile pattern
	freqByte := profileData[1]
	
	switch {
	case freqByte >= 0x50:
		return 6000 // DDR5-6000
	case freqByte >= 0x40:
		return 5600 // DDR5-5600
	case freqByte >= 0x30:
		return 5200 // DDR5-5200
	default:
		return 4800 // DDR5-4800 fallback
	}
}

// extractEXPOVoltage extracts voltage from EXPO profile data
func extractEXPOVoltage(profileData []byte) float64 {
	if len(profileData) < 6 {
		return 1.25 // Default DDR5 EXPO voltage
	}

	// EXPO voltage encoding (byte 5 typically holds voltage info)
	voltageByte := profileData[5]
	
	// Voltage calculation for EXPO (simplified)
	// Base voltage 1.10V + increments
	voltage := 1.10 + float64(voltageByte)*0.00625
	
	// Clamp to reasonable DDR5 range
	if voltage < 1.05 || voltage > 1.50 {
		return 1.25 // Safe fallback
	}

	return voltage
}

// extractEXPOTimings extracts timing parameters from EXPO profile data
func extractEXPOTimings(profileData []byte) TimingSpec {
	if len(profileData) < 12 {
		// Default tight timings for high-performance EXPO
		return TimingSpec{
			CL:   30,
			TRCD: 36,
			TRP:  36,
			TRAS: 52,
		}
	}

	// Extract timings from profile data (simplified)
	// Real implementation would decode specific timing bytes
	return TimingSpec{
		CL:   uint32(profileData[8]),  // CAS Latency
		TRCD: uint32(profileData[9]),  // RAS to CAS Delay
		TRP:  uint32(profileData[10]), // Row Precharge
		TRAS: uint32(profileData[11]), // Row Active Strobe
	}
}

// validateEXPOProfile performs basic validation on extracted EXPO profile
func validateEXPOProfile(profile *MemoryProfile) error {
	// Check frequency is in reasonable DDR5 range
	if profile.Frequency < 4800 || profile.Frequency > 8000 {
		return fmt.Errorf("frequency %d MHz out of valid range", profile.Frequency)
	}

	// Check voltage is reasonable for DDR5
	if profile.Voltage < 1.05 || profile.Voltage > 1.60 {
		return fmt.Errorf("voltage %.3fV out of valid range", profile.Voltage)
	}

	// Check timings make sense
	timings := profile.Timings
	if timings.CL == 0 || timings.CL > 100 {
		return fmt.Errorf("invalid CAS latency: %d", timings.CL)
	}

	if timings.TRAS <= timings.TRCD+timings.TRP {
		return fmt.Errorf("invalid timing relationship: TRAS(%d) <= TRCD(%d)+TRP(%d)", 
			timings.TRAS, timings.TRCD, timings.TRP)
	}

	return nil
}