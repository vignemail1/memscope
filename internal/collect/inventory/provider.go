package inventory

import (
	"fmt"
	"runtime"
)

// Provider defines the interface for hardware inventory collection
type Provider interface {
	CollectInventory() (*SystemInventory, error)
	CollectMemoryWithSPD() ([]MemoryModuleInfo, error)
}

// NewProvider creates a platform-specific inventory provider
func NewProvider() Provider {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsProvider()
	default:
		return &UnsupportedProvider{platform: runtime.GOOS}
	}
}

// UnsupportedProvider handles non-Windows platforms
type UnsupportedProvider struct {
	platform string
}

// CollectInventory returns an error for unsupported platforms
func (p *UnsupportedProvider) CollectInventory() (*SystemInventory, error) {
	return nil, fmt.Errorf("inventory collection not supported on %s", p.platform)
}

// CollectMemoryWithSPD returns an error for unsupported platforms
func (p *UnsupportedProvider) CollectMemoryWithSPD() ([]MemoryModuleInfo, error) {
	return nil, fmt.Errorf("memory with SPD collection not supported on %s", p.platform)
}