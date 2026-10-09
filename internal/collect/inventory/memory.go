//go:build !windows

package inventory

import "fmt"

// WindowsProvider stub for non-Windows platforms
type WindowsProvider struct{}

// NewWindowsProvider creates a stub provider that returns errors
func NewWindowsProvider() *WindowsProvider {
	return &WindowsProvider{}
}

// CollectInventory returns an error on non-Windows platforms
func (p *WindowsProvider) CollectInventory() (*SystemInventory, error) {
	return nil, fmt.Errorf("Windows inventory provider not available on this platform")
}

// CollectMemoryWithSPD returns an error on non-Windows platforms
func (p *WindowsProvider) CollectMemoryWithSPD() ([]MemoryModuleInfo, error) {
	return nil, fmt.Errorf("Windows memory with SPD collection not available on this platform")
}