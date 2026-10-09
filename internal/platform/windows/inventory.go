package windows

import (
	"fmt"
	"time"

	"github.com/vignemail1/memscope/internal/collect/inventory"
)

// InventoryProvider implements inventory collection for Windows systems
type InventoryProvider struct{}

// NewInventoryProvider creates a new Windows inventory provider
func NewInventoryProvider() *InventoryProvider {
	return &InventoryProvider{}
}

// CollectInventory gathers hardware inventory information
func (p *InventoryProvider) CollectInventory() (*inventory.SystemInventory, error) {
	// Get system information
	systemInfo, err := GetSystemInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}
	
	// Get memory information  
	memoryModules, err := GetMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	
	// Create inventory structure
	inv := &inventory.SystemInventory{
		Timestamp: time.Now(),
		System: inventory.SystemInfo{
			Manufacturer: systemInfo.Manufacturer,
			Model:        systemInfo.Model,
			SerialNumber: systemInfo.SerialNumber,
		},
		BIOS: inventory.BIOSInfo{
			Vendor:  systemInfo.BIOSVendor,
			Version: systemInfo.BIOSVersion,
			Date:    systemInfo.BIOSDate,
		},
		Memory: make([]inventory.MemoryModuleInfo, len(memoryModules)),
	}
	
	// Convert memory modules
	for i, module := range memoryModules {
		inv.Memory[i] = inventory.MemoryModuleInfo{
			DeviceLocator: module.DeviceLocator,
			BankLabel:     module.BankLabel,
			Capacity:      module.Capacity,
			Speed:         module.Speed,
			Manufacturer:  module.Manufacturer,
			PartNumber:    module.PartNumber,
			SerialNumber:  module.SerialNumber,
			DataWidth:     module.DataWidth,
			TotalWidth:    module.TotalWidth,
			FormFactor:    module.FormFactor,
			MemoryType:    module.MemoryType,
		}
	}
	
	return inv, nil
}