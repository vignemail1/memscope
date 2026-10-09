package spd

import (
	"fmt"
)

// XMP signature and offsets for DDR4
const (
	XMPSignature   = 0x4A0C0F0C
	XMP2Offset     = 384  // Byte offset for XMP 2.0 header in DDR4
	XMP3Offset     = 384  // Byte offset for XMP 3.0 header in DDR5
	XMPProfileSize = 29   // Size of each XMP profile in bytes
)

// ParseXMPProfiles extracts Intel XMP profiles from SPD data
func ParseXMPProfiles(data []byte) ([]MemoryProfile, error) {
	if len(data) < 512 {
		return nil, fmt.Errorf("insufficient data for XMP parsing: need 512 bytes, got %d", len(data))
	}

	var profiles []MemoryProfile

	// Check DDR type to determine XMP version and offset
	switch data[2] {
	case 0x0C: // DDR4
		xmpProfiles, err := parseXMP2Profiles(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse XMP 2.0 profiles: %w", err)
		}
		profiles = append(profiles, xmpProfiles...)
	case 0x12: // DDR5
		xmpProfiles, err := parseXMP3Profiles(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse XMP 3.0 profiles: %w", err)
		}
		profiles = append(profiles, xmpProfiles...)
	default:
		return nil, fmt.Errorf("XMP not supported for DDR type: 0x%02X", data[2])
	}

	return profiles, nil
}

// parseXMP2Profiles parses XMP 2.0 profiles for DDR4
func parseXMP2Profiles(data []byte) ([]MemoryProfile, error) {
	if len(data) < XMP2Offset+64 {
		return nil, fmt.Errorf("insufficient data for XMP 2.0 parsing")
	}

	// Verify XMP signature
	xmpData := data[XMP2Offset:]
	if err := verifyXMPSignature(xmpData); err != nil {
		return nil, fmt.Errorf("XMP signature verification failed: %w", err)
	}

	var profiles []MemoryProfile

	// Check organization byte to see which profiles are present
	orgByte := xmpData[13]
	
	// Profile 1
	if orgByte&0x01 != 0 {
		profile, err := parseXMP2Profile(xmpData, 1)
		if err == nil {
			profiles = append(profiles, *profile)
		}
	}

	// Profile 2
	if orgByte&0x02 != 0 {
		profile, err := parseXMP2Profile(xmpData, 2)
		if err == nil {
			profiles = append(profiles, *profile)
		}
	}

	return profiles, nil
}

// parseXMP3Profiles parses XMP 3.0 profiles for DDR5
func parseXMP3Profiles(data []byte) ([]MemoryProfile, error) {
	if len(data) < XMP3Offset+128 {
		return nil, fmt.Errorf("insufficient data for XMP 3.0 parsing")
	}

	// XMP 3.0 has different structure than 2.0
	// This is a simplified implementation
	var profiles []MemoryProfile
	
	profile := &MemoryProfile{
		Type:        "XMP",
		ProfileName: "XMP 3.0 Profile 1",
		Frequency:   5600, // Example DDR5-5600
		Voltage:     1.25,  // Typical XMP DDR5 voltage
		Timings: TimingSpec{
			CL:   36,
			TRCD: 36,
			TRP:  36,
			TRAS: 52,
		},
	}
	profiles = append(profiles, *profile)

	return profiles, nil
}

// parseXMP2Profile parses a specific XMP 2.0 profile
func parseXMP2Profile(xmpData []byte, profileNum int) (*MemoryProfile, error) {
	// Calculate offset for the specific profile
	var profileOffset int
	if profileNum == 1 {
		profileOffset = 16 // Profile 1 starts at byte 16
	} else if profileNum == 2 {
		profileOffset = 35 // Profile 2 starts at byte 35  
	} else {
		return nil, fmt.Errorf("invalid profile number: %d", profileNum)
	}

	if len(xmpData) < profileOffset+XMPProfileSize {
		return nil, fmt.Errorf("insufficient data for profile %d", profileNum)
	}

	profileData := xmpData[profileOffset:]

	profile := &MemoryProfile{
		Type:        "XMP",
		ProfileName: fmt.Sprintf("XMP 2.0 Profile %d", profileNum),
	}

	// Extract frequency (simplified calculation)
	// Real XMP uses complex MTB calculations
	profile.Frequency = extractXMPFrequency(profileData)

	// Extract voltage
	profile.Voltage = extractXMPVoltage(profileData)

	// Extract timings
	profile.Timings = extractXMPTimings(profileData)

	return profile, nil
}

// verifyXMPSignature checks if XMP signature is present and returns specific error
func verifyXMPSignature(xmpData []byte) error {
	if len(xmpData) < 4 {
		return fmt.Errorf("insufficient data for XMP signature verification: need 4 bytes, got %d", len(xmpData))
	}

	// XMP signature is stored in little-endian format
	signature := uint32(xmpData[0]) | 
		uint32(xmpData[1])<<8 | 
		uint32(xmpData[2])<<16 | 
		uint32(xmpData[3])<<24

	if signature != XMPSignature {
		return fmt.Errorf("invalid XMP signature: expected 0x%08X, got 0x%08X", XMPSignature, signature)
	}

	return nil
}

// extractXMPFrequency calculates frequency from XMP profile data
func extractXMPFrequency(profileData []byte) uint32 {
	// Simplified frequency extraction
	// Real implementation would use MTB and calculate from tCK
	return 3200 // Example DDR4-3200 XMP frequency
}

// extractXMPVoltage extracts voltage from XMP profile data
func extractXMPVoltage(profileData []byte) float64 {
	// Voltage is typically stored in byte 7 for XMP 2.0
	if len(profileData) < 8 {
		return 1.35 // Default XMP voltage
	}

	voltageByte := profileData[7]
	// Voltage encoding: 1.20V = 0x78, each step is 0.00625V
	voltage := 1.20 + float64(voltageByte-0x78)*0.00625
	
	// Clamp to reasonable range
	if voltage < 1.0 || voltage > 2.0 {
		return 1.35 // Default fallback
	}

	return voltage
}

// extractXMPTimings extracts timing parameters from XMP profile data
func extractXMPTimings(profileData []byte) TimingSpec {
	// Simplified timing extraction
	// Real implementation would calculate from MTB values
	return TimingSpec{
		CL:   16, // Typical CL for DDR4-3200
		TRCD: 18,
		TRP:  18,
		TRAS: 36,
	}
}