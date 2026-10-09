// internal/collect/runtime/provider.go
package runtime

import (
	"fmt"
	"runtime"
)

// Provider interface for runtime memory parameter collection
type Provider interface {
	CollectMemoryParameters() (*MemoryParameters, error)
}

// NewProvider creates a platform-appropriate runtime provider
func NewProvider() Provider {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsProvider()
	default:
		return &UnsupportedProvider{}
	}
}

// UnsupportedProvider for non-Windows platforms
type UnsupportedProvider struct{}

func (p *UnsupportedProvider) CollectMemoryParameters() (*MemoryParameters, error) {
	return nil, fmt.Errorf("runtime memory parameter collection not supported on %s", runtime.GOOS)
}