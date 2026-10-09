//go:build !windows

package windows

import "errors"

// MockSPDReader provides mock SPD access for non-Windows platforms
type MockSPDReader struct{}

// NewMockSPDReader creates a mock SPD reader
func NewMockSPDReader() *MockSPDReader {
	return &MockSPDReader{}
}

// ReadSPD returns an error on non-Windows platforms
func (m *MockSPDReader) ReadSPD(slot int) ([]byte, error) {
	return nil, errors.New("SPD reading not available on this platform")
}

// GetSlotCount returns an error on non-Windows platforms  
func (m *MockSPDReader) GetSlotCount() (int, error) {
	return 0, errors.New("SPD reading not available on this platform")
}