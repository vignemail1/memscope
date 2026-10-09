package inventory

import "time"

// SystemInventory represents complete system hardware inventory
type SystemInventory struct {
	Timestamp time.Time
	System    SystemInfo
	BIOS      BIOSInfo
	Memory    []MemoryModuleInfo
}

// SystemInfo represents basic system information
type SystemInfo struct {
	Manufacturer string
	Model        string
	SerialNumber string
}

// BIOSInfo represents BIOS/UEFI information
type BIOSInfo struct {
	Vendor  string
	Version string
	Date    string
}

// MemoryModuleInfo represents a physical memory module
type MemoryModuleInfo struct {
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