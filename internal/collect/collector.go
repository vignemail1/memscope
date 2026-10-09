package collect

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/vignemail1/memscope/internal/model"
)

var ErrConfiguration = errors.New("invalid collector configuration")

// Config identifies a capture and bounds the entire collection operation.
// SnapshotID must be supplied by the caller; this package generates no identity.
type Config struct {
	SnapshotID    string
	ToolVersion   string
	SchemaVersion string
	Timeout       time.Duration
}

// Collector coordinates providers sequentially, without automatic elevation.
// Providers must honor context cancellation and must not return aliased mutable data.
// A provider interface must be nil or contain a non-nil implementation.
type Collector struct {
	providers Providers
	config    Config
}

func New(providers Providers, config Config) (*Collector, error) {
	for _, field := range []struct{ name, value string }{
		{"snapshot_id", config.SnapshotID}, {"tool_version", config.ToolVersion},
		{"schema_version", config.SchemaVersion},
	} {
		if strings.TrimSpace(field.value) == "" {
			return nil, fmt.Errorf("%w: %s is required", ErrConfiguration, field.name)
		}
	}
	if config.Timeout <= 0 {
		return nil, fmt.Errorf("%w: timeout must be positive", ErrConfiguration)
	}
	return &Collector{providers: providers, config: config}, nil
}

// Collect returns a snapshot even when providers return errors.
// Callers must inspect both results; a non-nil error must not be discarded.
// Missing providers alone are not errors. Context deadlines are cooperative.
func (c *Collector) Collect(parent context.Context) (model.Snapshot, error) {
	if c == nil || parent == nil {
		return model.Snapshot{}, fmt.Errorf("%w: collector and context must not be nil", ErrConfiguration)
	}
	ctx, cancel := context.WithTimeout(parent, c.config.Timeout)
	defer cancel()
	snapshot := model.Snapshot{
		SchemaVersion: c.config.SchemaVersion, ToolVersion: c.config.ToolVersion,
		ID: c.config.SnapshotID, CollectionStartedAt: time.Now().UTC(),
		Platform: model.Platform{OS: runtime.GOOS, Arch: runtime.GOARCH},
	}
	var failures []error
	record := func(stage string, err error) {
		if err == nil {
			return
		}
		failures = append(failures, fmt.Errorf("%s: %w", stage, err))
		snapshot.Diagnostics = append(snapshot.Diagnostics, model.Diagnostic{
			Code: "collect_error", Severity: "error", Message: stage + ": " + err.Error(),
		})
	}
	run := func(name string, provider Provider, acquire func()) {
		if provider == nil {
			snapshot.Capabilities = append(snapshot.Capabilities, model.Capability{
				Name: name, Available: false, Reason: "provider not configured",
			})
			return
		}
		if err := ctx.Err(); err != nil {
			snapshot.Capabilities = append(snapshot.Capabilities, model.Capability{
				Name: name, Available: false, Reason: err.Error(),
			})
			return
		}
		capabilities, err := provider.Capabilities(ctx)
		snapshot.Capabilities = append(snapshot.Capabilities, capabilities...)
		record(name+".capabilities", err)
		// A failed probe may still supply useful metadata, but acquisition is skipped.
		if err != nil || ctx.Err() != nil {
			return
		}
		acquire()
	}
	run("inventory", c.providers.Inventory, func() {
		result, err := c.providers.Inventory.Inventory(ctx)
		snapshot.Devices = append(snapshot.Devices, result.Devices...)
		snapshot.Observations = append(snapshot.Observations, result.Observations...)
		snapshot.Diagnostics = append(snapshot.Diagnostics, result.Diagnostics...)
		record("inventory", err)
	})
	run("spd", c.providers.SPD, func() {
		result, err := c.providers.SPD.ReadSPD(ctx)
		snapshot.Profiles = append(snapshot.Profiles, result.Profiles...)
		snapshot.Observations = append(snapshot.Observations, result.Observations...)
		snapshot.Diagnostics = append(snapshot.Diagnostics, result.Diagnostics...)
		record("spd", err)
	})
	run("runtime", c.providers.Runtime, func() {
		result, err := c.providers.Runtime.ReadRuntime(ctx)
		snapshot.Observations = append(snapshot.Observations, result.Observations...)
		snapshot.Diagnostics = append(snapshot.Diagnostics, result.Diagnostics...)
		record("runtime", err)
	})
	record("context", ctx.Err())
	snapshot.CollectionFinishedAt = time.Now().UTC()
	if err := snapshot.Validate(); err != nil {
		record("snapshot.validate", err)
	}
	return snapshot, errors.Join(failures...)
}
