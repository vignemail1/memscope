package windows

import (
	"testing"
)

// TestTypeConversionSafety tests the type conversion logic without requiring Windows
func TestTypeConversionSafety(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]interface{}
		expected MemoryModule
		wantErr  bool
	}{
		{
			name: "normal uint64 capacity",
			input: map[string]interface{}{
				"Capacity": uint64(8589934592), // 8GB
				"Speed":    uint32(2400),
			},
			expected: MemoryModule{
				Capacity: 8589934592,
				Speed:    2400,
			},
			wantErr: false,
		},
		{
			name: "int64 capacity conversion",
			input: map[string]interface{}{
				"Capacity": int64(4294967296), // 4GB
				"Speed":    int32(1600),
			},
			expected: MemoryModule{
				Capacity: 4294967296,
				Speed:    1600,
			},
			wantErr: false,
		},
		{
			name: "int32 capacity conversion",
			input: map[string]interface{}{
				"Capacity": int32(1073741824), // 1GB
				"Speed":    uint16(1333),
			},
			expected: MemoryModule{
				Capacity: 1073741824,
				Speed:    1333,
			},
			wantErr: false,
		},
		{
			name: "uint32 capacity conversion",
			input: map[string]interface{}{
				"Capacity": uint32(2147483648), // 2GB
			},
			expected: MemoryModule{
				Capacity: 2147483648,
			},
			wantErr: false,
		},
		{
			name: "negative int64 capacity should error",
			input: map[string]interface{}{
				"Capacity": int64(-1),
			},
			wantErr: true,
		},
		{
			name: "negative int32 capacity should error",
			input: map[string]interface{}{
				"Capacity": int32(-1),
			},
			wantErr: true,
		},
		{
			name: "string fields with correct types",
			input: map[string]interface{}{
				"DeviceLocator": "DIMM1",
				"BankLabel":     "Bank 0",
				"Manufacturer":  "Samsung",
				"PartNumber":    "M471A1K43CB1-CTD",
				"SerialNumber":  "12345678",
			},
			expected: MemoryModule{
				DeviceLocator: "DIMM1",
				BankLabel:     "Bank 0",
				Manufacturer:  "Samsung",
				PartNumber:    "M471A1K43CB1-CTD",
				SerialNumber:  "12345678",
			},
			wantErr: false,
		},
		{
			name: "wrong type for string field should error",
			input: map[string]interface{}{
				"DeviceLocator": 123,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			module, err := convertMemoryModuleData(tc.input)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if module.Capacity != tc.expected.Capacity {
				t.Errorf("Capacity: expected %d, got %d", tc.expected.Capacity, module.Capacity)
			}
			if module.Speed != tc.expected.Speed {
				t.Errorf("Speed: expected %d, got %d", tc.expected.Speed, module.Speed)
			}
			if module.DeviceLocator != tc.expected.DeviceLocator {
				t.Errorf("DeviceLocator: expected %s, got %s", tc.expected.DeviceLocator, module.DeviceLocator)
			}
			if module.BankLabel != tc.expected.BankLabel {
				t.Errorf("BankLabel: expected %s, got %s", tc.expected.BankLabel, module.BankLabel)
			}
			if module.Manufacturer != tc.expected.Manufacturer {
				t.Errorf("Manufacturer: expected %s, got %s", tc.expected.Manufacturer, module.Manufacturer)
			}
			if module.PartNumber != tc.expected.PartNumber {
				t.Errorf("PartNumber: expected %s, got %s", tc.expected.PartNumber, module.PartNumber)
			}
			if module.SerialNumber != tc.expected.SerialNumber {
				t.Errorf("SerialNumber: expected %s, got %s", tc.expected.SerialNumber, module.SerialNumber)
			}
		})
	}
}