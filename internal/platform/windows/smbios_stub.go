//go:build !windows

package windows

import "errors"

type SMBIOSReader struct{}

func NewSMBIOSReader() (*SMBIOSReader, error) {
	return nil, errors.New("SMBIOS access not available on this platform")
}

func (s *SMBIOSReader) GetBIOSInfo() (*BIOSInfo, error) {
	return nil, errors.New("SMBIOS access not available on this platform")
}

func (s *SMBIOSReader) GetSystemInfo() (*SystemInfo, error) {
	return nil, errors.New("SMBIOS access not available on this platform")
}

func (s *SMBIOSReader) GetMemoryInfo() ([]MemoryModule, error) {
	return nil, errors.New("SMBIOS access not available on this platform")
}

func (s *SMBIOSReader) Close() {}

type BIOSInfo struct {
	Vendor      string
	Version     string
	ReleaseDate string
}