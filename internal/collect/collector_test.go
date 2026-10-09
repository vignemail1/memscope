package collect

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vignemail1/memscope/internal/model"
)

type fakeInventory struct {
	probeErr error
	acquire  func(context.Context) (InventoryResult, error)
}

func (f *fakeInventory) Descriptor() Descriptor { return Descriptor{ID: "fake", Version: "test"} }

func (f *fakeInventory) Capabilities(ctx context.Context) ([]model.Capability, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []model.Capability{{Name: "inventory", Available: f.probeErr == nil}}, f.probeErr
}

func (f *fakeInventory) Inventory(ctx context.Context) (InventoryResult, error) {
	return f.acquire(ctx)
}

func testConfig() Config {
	return Config{SnapshotID: "test", ToolVersion: "test", SchemaVersion: "1", Timeout: time.Second}
}

func TestConfiguration(t *testing.T) {
	for _, mutate := range []func(*Config){
		func(c *Config) { c.SnapshotID = " " },
		func(c *Config) { c.ToolVersion = "" },
		func(c *Config) { c.SchemaVersion = "" },
		func(c *Config) { c.Timeout = 0 },
		func(c *Config) { c.Timeout = -time.Second },
	} {
		config := testConfig()
		mutate(&config)
		if _, err := New(Providers{}, config); !errors.Is(err, ErrConfiguration) {
			t.Fatalf("got %v, want ErrConfiguration", err)
		}
	}
}

func TestAbsentProviders(t *testing.T) {
	c, err := New(Providers{}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Capabilities) != 3 || len(s.Observations) != 0 {
		t.Fatal("unexpected missing-provider snapshot")
	}
	for _, capability := range s.Capabilities {
		if capability.Available {
			t.Fatal("missing provider marked available")
		}
	}
}

func TestPartialResult(t *testing.T) {
	provider := &fakeInventory{acquire: func(context.Context) (InventoryResult, error) {
		return InventoryResult{
			Devices: []model.Device{{ID: "dimm0", Kind: "memory"}},
			Observations: []model.Observation{{
				DeviceID: "dimm0", Scope: "device", Parameter: "capacity",
				Value: model.NewUnsigned(0), Unit: "bytes", Source: model.SMBIOS,
				Status: model.Reported, CapturedAt: time.Now().UTC(),
			}},
		}, ErrPermissionDenied
	}}
	c, err := New(Providers{Inventory: provider}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.Collect(context.Background())
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("got %v, want ErrPermissionDenied", err)
	}
	if len(s.Devices) != 1 || len(s.Observations) != 1 || len(s.Diagnostics) == 0 {
		t.Fatal("partial result lost")
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProbeFailure(t *testing.T) {
	provider := &fakeInventory{probeErr: ErrUnsupported, acquire: func(context.Context) (InventoryResult, error) {
		t.Fatal("acquisition called after failed probe")
		return InventoryResult{}, nil
	}}
	c, err := New(Providers{Inventory: provider}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Collect(context.Background()); !errors.Is(err, ErrUnsupported) {
		t.Fatalf("got %v, want ErrUnsupported", err)
	}
}

func TestCanceledContext(t *testing.T) {
	provider := &fakeInventory{acquire: func(context.Context) (InventoryResult, error) {
		t.Fatal("acquisition called with canceled context")
		return InventoryResult{}, nil
	}}
	c, err := New(Providers{Inventory: provider}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Collect(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
}

func TestDeadline(t *testing.T) {
	provider := &fakeInventory{acquire: func(ctx context.Context) (InventoryResult, error) {
		<-ctx.Done()
		return InventoryResult{}, ctx.Err()
	}}
	config := testConfig()
	config.Timeout = 10 * time.Millisecond
	c, err := New(Providers{Inventory: provider}, config)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Collect(context.Background()); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want context.DeadlineExceeded", err)
	}
}

func TestInvalidSnapshot(t *testing.T) {
	provider := &fakeInventory{acquire: func(context.Context) (InventoryResult, error) {
		return InventoryResult{Devices: []model.Device{{ID: "duplicate", Kind: "memory"}, {ID: "duplicate", Kind: "memory"}}}, nil
	}}
	c, err := New(Providers{Inventory: provider}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.Collect(context.Background())
	if !errors.Is(err, model.ErrInvalid) || len(s.Devices) != 2 {
		t.Fatalf("invalid result not preserved and flagged: %v", err)
	}
}

func TestNilInputs(t *testing.T) {
	var c *Collector
	if _, err := c.Collect(context.Background()); !errors.Is(err, ErrConfiguration) {
		t.Fatal("nil collector accepted")
	}
	c, err := New(Providers{}, testConfig())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Collect(nil); !errors.Is(err, ErrConfiguration) {
		t.Fatal("nil context accepted")
	}
}
