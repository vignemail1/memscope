//go:build !windows

package windows

import "errors"

type RegistryReader struct{}

func NewRegistryReader() (*RegistryReader, error) {
	return nil, errors.New("registry access not available on this platform")
}

func (r *RegistryReader) GetSystemInfo() (*SystemInfo, error) {
	return nil, errors.New("registry access not available on this platform")
}

func (r *RegistryReader) GetMemoryInfo() ([]MemoryModule, error) {
	return nil, errors.New("registry access not available on this platform")
}

func (r *RegistryReader) GetBIOSInfo() (*BIOSInfo, error) {
	return nil, errors.New("registry access not available on this platform")
}

func (r *RegistryReader) Close() {}