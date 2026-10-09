package windows

import (
	"testing"

	"github.com/vignemail1/memscope/internal/collect/inventory"
)

// TestInventoryProviderExists verifies that the inventory provider can be instantiated
// This test can run on any platform
func TestInventoryProviderExists(t *testing.T) {
	provider := inventory.NewProvider()
	if provider == nil {
		t.Error("NewProvider should not return nil")
	}
}

// TestInventoryProviderInterface verifies the provider implements the expected interface
func TestInventoryProviderInterface(t *testing.T) {
	provider := inventory.NewProvider()
	
	// This will fail at compile time if the interface is not satisfied
	var _ interface {
		CollectInventory() (*inventory.SystemInventory, error)
	} = provider
	
	// If we get here, the interface is correctly implemented
}