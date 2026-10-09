//go:build windows

package inventory

import (
	"fmt"
	"time"

	"github.com/vignemail1/memscope/internal/platform/windows"
)

// WindowsProvider implements inventory collection for Windows systems
type WindowsProvider struct{}

// NewWindowsProvider creates a new Windows inventory provider
func NewWindowsProvider() *WindowsProvider {
	return &WindowsProvider{}
}

// CollectInventory gathers complete hardware inventory information
func (p *WindowsProvider) CollectInventory() (*SystemInventory, error) {
	// Get system information
	systemInfo, err := windows.GetSystemInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}
	
	// Get memory information  
	memoryModules, err := windows.GetMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}
	
	// Get CPU information
	cpuInfo, err := windows.GetCPUInfo()
	if err != nil {
		// CPU collection is not critical, log and continue with placeholder
		cpuInfo = &windows.CPUInfo{
			Name: "CPU info collection failed",
		}
	}
	
	// Create inventory structure
	inv := &SystemInventory{
		Timestamp: time.Now(),
		System: SystemInfo{
			Manufacturer: systemInfo.Manufacturer,
			Model:        systemInfo.Model,
			SerialNumber: systemInfo.SerialNumber,
		},
		CPU: CPUInfo{
			Name:         cpuInfo.Name,
			Manufacturer: cpuInfo.Manufacturer,
			Architecture: cpuInfo.Architecture,
			Cores:        cpuInfo.Cores,
			Threads:      cpuInfo.Threads,
			MaxClockMHz:  cpuInfo.MaxClockMHz,
		},
		BIOS: BIOSInfo{
			Vendor:  systemInfo.BIOSVendor,
			Version: systemInfo.BIOSVersion,
			Date:    systemInfo.BIOSDate,
		},
		Memory: make([]MemoryModuleInfo, len(memoryModules)),
	}
	
	// Convert memory modules
	for i, module := range memoryModules {
		inv.Memory[i] = MemoryModuleInfo{
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

// CollectMemoryWithSPD gathers memory information including SPD data
func (p *WindowsProvider) CollectMemoryWithSPD() ([]MemoryModuleInfo, error) {
	// Start with basic memory information
	memoryModules, err := windows.GetMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	// Convert to our format
	modules := make([]MemoryModuleInfo, len(memoryModules))
	for i, module := range memoryModules {
		modules[i] = MemoryModuleInfo{
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
			// SPD data collection will be added in future tasks
			SPD: nil,
		}
	}

	return modules, nil
}