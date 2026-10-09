package spd

import (
	"testing"
)

func TestParseSPDData(t *testing.T) {
	// Mock DDR4 SPD data (simplified)
	mockSPDData := make([]byte, 512)
	mockSPDData[0] = 0x24  // Number of bytes used/total bytes in SPD device
	mockSPDData[1] = 0x00  // SPD revision
	mockSPDData[2] = 0x0C  // DRAM device type (DDR4)
	mockSPDData[3] = 0x01  // Module type (UDIMM)
	
	profile, err := ParseSPDData(mockSPDData)
	if err != nil {
		t.Fatalf("ParseSPDData failed: %v", err)
	}
	
	if profile.DeviceType != "DDR4" {
		t.Errorf("Expected DDR4, got %s", profile.DeviceType)
	}
}

func TestJEDECProfile(t *testing.T) {
	mockSPDData := make([]byte, 512)
	mockSPDData[2] = 0x0C  // DDR4
	mockSPDData[18] = 0x08 // Timings table
	
	jedec, err := ParseJEDECProfile(mockSPDData)
	if err != nil {
		t.Fatalf("ParseJEDECProfile failed: %v", err)
	}
	
	if jedec == nil {
		t.Fatal("JEDEC profile should not be nil")
	}
}

func TestXMPProfiles(t *testing.T) {
	mockSPDData := make([]byte, 512)
	mockSPDData[2] = 0x0C  // DDR4
	
	// Mock XMP signature at offset 384
	xmpOffset := 384
	mockSPDData[xmpOffset] = 0x0C
	mockSPDData[xmpOffset+1] = 0x0F
	mockSPDData[xmpOffset+2] = 0x0C
	mockSPDData[xmpOffset+3] = 0x4A
	mockSPDData[xmpOffset+13] = 0x01 // Organization byte - Profile 1 enabled
	
	profiles, err := ParseXMPProfiles(mockSPDData)
	if err != nil {
		t.Fatalf("ParseXMPProfiles failed: %v", err)
	}
	
	if len(profiles) == 0 {
		t.Fatal("Expected at least one XMP profile")
	}
	
	if profiles[0].Type != "XMP" {
		t.Errorf("Expected XMP profile type, got %s", profiles[0].Type)
	}
}

func TestEXPOProfiles(t *testing.T) {
	mockSPDData := make([]byte, 512)
	mockSPDData[2] = 0x12  // DDR5
	
	// Mock EXPO signature at offset 480
	expoOffset := 480
	copy(mockSPDData[expoOffset:expoOffset+4], []byte("EXPO"))
	
	// Create a valid profile 1 (starts at offset+4)
	profileStart := expoOffset + 4
	mockSPDData[profileStart] = 0x01    // Profile present (non-zero first byte)
	mockSPDData[profileStart+1] = 0x40  // Frequency indicator
	mockSPDData[profileStart+5] = 0x20  // Voltage byte (at byte 5 of profile)
	
	// Add timing values (bytes 8-11 of profile data)
	// TRAS must be > TRCD + TRP for validation to pass
	mockSPDData[profileStart+8] = 30   // CL
	mockSPDData[profileStart+9] = 20   // TRCD
	mockSPDData[profileStart+10] = 20  // TRP
	mockSPDData[profileStart+11] = 50  // TRAS (must be > 20+20=40)
	
	profiles, err := ParseEXPOProfiles(mockSPDData)
	if err != nil {
		t.Fatalf("ParseEXPOProfiles failed: %v", err)
	}
	
	if len(profiles) == 0 {
		t.Fatal("Expected at least one EXPO profile")
	}
	
	if profiles[0].Type != "EXPO" {
		t.Errorf("Expected EXPO profile type, got %s", profiles[0].Type)
	}
}