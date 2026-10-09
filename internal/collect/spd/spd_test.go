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

func TestParseSPDDataIntegration(t *testing.T) {
	// Create comprehensive mock DDR4 SPD data
	mockSPDData := make([]byte, 512)
	
	// Basic info
	mockSPDData[2] = 0x0C  // DDR4
	mockSPDData[3] = 0x01  // UDIMM
	
	// Capacity calculation data
	mockSPDData[4] = 0x03  // SDRAM capacity
	mockSPDData[12] = 0x03 // Module organization
	mockSPDData[13] = 0x03 // Bus width
	
	// JEDEC timing data
	mockSPDData[17] = 0x01 // MTB
	mockSPDData[18] = 0x0F // tCK min
	mockSPDData[24] = 0x6E // tAA min
	
	// Manufacturer data
	mockSPDData[320] = 0x80 // Manufacturer ID low
	mockSPDData[321] = 0x2C // Manufacturer ID high (Micron)
	// Part number with bounds checking
	partNumberBytes := []byte("TestModule123       ")
	if len(mockSPDData) >= 329+len(partNumberBytes) {
		copy(mockSPDData[329:], partNumberBytes)
	}
	
	profile, err := ParseSPDData(mockSPDData)
	if err != nil {
		t.Fatalf("ParseSPDData failed: %v", err)
	}
	
	// Verify all components are parsed
	if profile.DeviceType != "DDR4" {
		t.Errorf("Expected DDR4, got %s", profile.DeviceType)
	}
	
	if profile.Capacity == 0 {
		t.Error("Capacity should not be zero")
	}
	
	if profile.Manufacturer == "" {
		t.Error("Manufacturer should not be empty")
	}
	
	if profile.JEDECProfile == nil {
		t.Error("JEDEC profile should not be nil")
	}
	
	if profile.JEDECProfile.Frequency == 0 {
		t.Error("JEDEC frequency should not be zero")
	}
}

func TestParseSPDDataDDR5Integration(t *testing.T) {
	// Create comprehensive mock DDR5 SPD data
	mockSPDData := make([]byte, 512)
	
	// Basic info
	mockSPDData[2] = 0x12  // DDR5
	mockSPDData[3] = 0x02  // SODIMM
	
	// Capacity calculation data
	mockSPDData[4] = 0x02  // SDRAM capacity
	mockSPDData[12] = 0x02 // Module organization
	mockSPDData[13] = 0x03 // Bus width
	
	// Manufacturer data
	mockSPDData[320] = 0x80 // Manufacturer ID low
	mockSPDData[321] = 0xAD // Manufacturer ID high (SK Hynix)
	// Part number with bounds checking
	partNumberBytes := []byte("DDR5TestMod         ")
	if len(mockSPDData) >= 329+len(partNumberBytes) {
		copy(mockSPDData[329:], partNumberBytes)
	}
	
	profile, err := ParseSPDData(mockSPDData)
	if err != nil {
		t.Fatalf("ParseSPDData failed: %v", err)
	}
	
	// Verify DDR5-specific components
	if profile.DeviceType != "DDR5" {
		t.Errorf("Expected DDR5, got %s", profile.DeviceType)
	}
	
	if profile.ModuleType != "SODIMM" {
		t.Errorf("Expected SODIMM, got %s", profile.ModuleType)
	}
	
	if profile.Capacity == 0 {
		t.Error("Capacity should not be zero")
	}
	
	if profile.JEDECProfile == nil {
		t.Error("JEDEC profile should not be nil")
	}
}

func TestCapacityCalculation(t *testing.T) {
	tests := []struct {
		name         string
		setupData    func([]byte)
		expectedMin  uint64 // Minimum expected capacity
	}{
		{
			name: "DDR4 Basic",
			setupData: func(data []byte) {
				data[2] = 0x0C  // DDR4
				data[4] = 0x01  // 512MB per die
				data[12] = 0x00 // 1 rank, x4 width
				data[13] = 0x03 // 64-bit bus
			},
			expectedMin: 1024 * 1024 * 1024, // At least 1GB
		},
		{
			name: "DDR5 Basic",
			setupData: func(data []byte) {
				data[2] = 0x12  // DDR5
				data[4] = 0x01  // 1GB per die
				data[12] = 0x00 // 1 rank, x4 width
				data[13] = 0x03 // 64-bit bus
			},
			expectedMin: 2 * 1024 * 1024 * 1024, // At least 2GB
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, 512)
			tt.setupData(data)
			
			capacity := parseCapacity(data)
			if capacity < tt.expectedMin {
				t.Errorf("Capacity %d is less than expected minimum %d", capacity, tt.expectedMin)
			}
		})
	}
}

func TestParseSPDDataWithXMPIntegration(t *testing.T) {
	// Create mock DDR4 SPD data with XMP profile
	mockSPDData := make([]byte, 512)
	
	// Basic info
	mockSPDData[2] = 0x0C  // DDR4
	mockSPDData[3] = 0x01  // UDIMM
	
	// Capacity data  
	mockSPDData[4] = 0x02  // SDRAM capacity
	mockSPDData[12] = 0x02 // Module organization
	mockSPDData[13] = 0x03 // Bus width
	
	// Mock XMP signature at offset 384
	xmpOffset := 384
	mockSPDData[xmpOffset] = 0x0C
	mockSPDData[xmpOffset+1] = 0x0F
	mockSPDData[xmpOffset+2] = 0x0C
	mockSPDData[xmpOffset+3] = 0x4A
	mockSPDData[xmpOffset+13] = 0x01 // Organization byte - Profile 1 enabled
	
	// Manufacturer data
	mockSPDData[320] = 0x04 // Manufacturer ID low
	mockSPDData[321] = 0xCD // Manufacturer ID high (G.Skill)
	// Part number with bounds checking
	partNumberBytes := []byte("TestXMPModule       ")
	if len(mockSPDData) >= 329+len(partNumberBytes) {
		copy(mockSPDData[329:], partNumberBytes)
	}
	
	profile, err := ParseSPDData(mockSPDData)
	if err != nil {
		t.Fatalf("ParseSPDData failed: %v", err)
	}
	
	// Verify JEDEC profile exists
	if profile.JEDECProfile == nil {
		t.Error("JEDEC profile should not be nil")
	}
	
	// Verify XMP profiles were parsed
	if len(profile.XMPProfiles) == 0 {
		t.Error("Expected at least one XMP profile")
	}
	
	if profile.XMPProfiles[0].Type != "XMP" {
		t.Errorf("Expected XMP profile type, got %s", profile.XMPProfiles[0].Type)
	}
}

func TestParseSPDDataWithEXPOIntegration(t *testing.T) {
	// Create mock DDR5 SPD data with EXPO profile
	mockSPDData := make([]byte, 512)
	
	// Basic info
	mockSPDData[2] = 0x12  // DDR5
	mockSPDData[3] = 0x02  // SODIMM
	
	// Capacity data
	mockSPDData[4] = 0x03  // SDRAM capacity
	mockSPDData[12] = 0x01 // Module organization  
	mockSPDData[13] = 0x03 // Bus width
	
	// Mock EXPO signature at offset 480
	expoOffset := 480
	copy(mockSPDData[expoOffset:expoOffset+4], []byte("EXPO"))
	
	// Create a valid profile 1
	profileStart := expoOffset + 4
	mockSPDData[profileStart] = 0x01    // Profile present
	mockSPDData[profileStart+1] = 0x50  // Frequency indicator
	mockSPDData[profileStart+5] = 0x25  // Voltage byte
	
	// Add timing values
	mockSPDData[profileStart+8] = 32   // CL
	mockSPDData[profileStart+9] = 24   // TRCD
	mockSPDData[profileStart+10] = 24  // TRP
	mockSPDData[profileStart+11] = 60  // TRAS (must be > TRCD+TRP)
	
	// Manufacturer data
	mockSPDData[320] = 0x80 // Manufacturer ID low
	mockSPDData[321] = 0xCE // Manufacturer ID high (Samsung)
	// Part number with bounds checking
	partNumberBytes := []byte("TestEXPOModule      ")
	if len(mockSPDData) >= 329+len(partNumberBytes) {
		copy(mockSPDData[329:], partNumberBytes)
	}
	
	profile, err := ParseSPDData(mockSPDData)
	if err != nil {
		t.Fatalf("ParseSPDData failed: %v", err)
	}
	
	// Verify JEDEC profile exists
	if profile.JEDECProfile == nil {
		t.Error("JEDEC profile should not be nil")
	}
	
	// Verify EXPO profiles were parsed
	if len(profile.EXPOProfiles) == 0 {
		t.Error("Expected at least one EXPO profile")
	}
	
	if profile.EXPOProfiles[0].Type != "EXPO" {
		t.Errorf("Expected EXPO profile type, got %s", profile.EXPOProfiles[0].Type)
	}
}