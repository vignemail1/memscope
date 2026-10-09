//go:build !windows

package windows

import "errors"

func GetSystemInfo() (*SystemInfo, error) {
	return nil, errors.New("Windows hardware access not available on this platform")
}

func GetMemoryInfo() ([]MemoryModule, error) {
	return nil, errors.New("Windows hardware access not available on this platform")
}

type SystemInfo struct {
	Manufacturer string
	Model        string
	SerialNumber string
	BIOSVendor   string
	BIOSVersion  string
	BIOSDate     string
}

type MemoryModule struct {
	DeviceLocator  string
	BankLabel      string
	Capacity       uint64
	Speed          uint32
	Manufacturer   string
	PartNumber     string
	SerialNumber   string
	DataWidth      uint16
	TotalWidth     uint16
	FormFactor     uint16
	MemoryType     uint16
}