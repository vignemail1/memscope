// Package collect defines acquisition interfaces without hardware dependencies.
package collect

import (
	"context"
	"errors"

	"github.com/vignemail1/memscope/internal/model"
)

var (
	ErrUnsupported      = errors.New("unsupported capability")
	ErrNotImplemented   = errors.New("collector not implemented")
	ErrPermissionDenied = errors.New("permission denied")
)

type Descriptor struct {
	ID                 string
	Version            string
	SupportedPlatforms []model.Platform
	Privileges         []string
}

// Provider describes and probes a collector without automatic elevation.
type Provider interface {
	Descriptor() Descriptor
	Capabilities(context.Context) ([]model.Capability, error)
}

type InventoryResult struct {
	Devices      []model.Device
	Observations []model.Observation
	Diagnostics  []model.Diagnostic
}

type InventoryProvider interface {
	Provider
	Inventory(context.Context) (InventoryResult, error)
}

type SPDResult struct {
	Profiles     []model.Profile
	Observations []model.Observation
	Diagnostics  []model.Diagnostic
}

type SpdProvider interface {
	Provider
	ReadSPD(context.Context) (SPDResult, error)
}

type RuntimeResult struct {
	Observations []model.Observation
	Diagnostics  []model.Diagnostic
}

type RuntimeMemoryProvider interface {
	Provider
	ReadRuntime(context.Context) (RuntimeResult, error)
}

// Providers may be nil when unavailable; callers must inspect them explicitly.
type Providers struct {
	Inventory InventoryProvider
	SPD       SpdProvider
	Runtime   RuntimeMemoryProvider
}
