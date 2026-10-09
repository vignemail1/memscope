package model

import (
	"fmt"
	"strings"
)

func required(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: %s is required", ErrInvalid, field)
	}
	return nil
}

func (s Status) Validate() error {
	switch s {
	case Observed, Reported, Derived, Unavailable, Invalid:
		return nil
	default:
		return fmt.Errorf("%w: unknown status %q", ErrInvalid, s)
	}
}

func (s Source) Validate() error {
	switch s {
	case SMBIOS, WMI, CPUID, SPD, Controller, Telemetry, Calculation, Import, User:
		return nil
	default:
		return fmt.Errorf("%w: unknown source %q", ErrInvalid, s)
	}
}

// Validate checks generic acquisition invariants, not parameter-specific ranges.
// Invalid observations contain no normalized value; raw evidence belongs elsewhere.
func (o Observation) Validate() error {
	for _, field := range []struct{ name, value string }{
		{"device_id", o.DeviceID}, {"scope", o.Scope}, {"parameter", o.Parameter},
	} {
		if err := required(field.name, field.value); err != nil {
			return err
		}
	}
	if err := o.Status.Validate(); err != nil {
		return err
	}
	if err := o.Source.Validate(); err != nil {
		return err
	}
	if o.CapturedAt.IsZero() {
		return fmt.Errorf("%w: captured_at is required", ErrInvalid)
	}
	if o.Status == Unavailable || o.Status == Invalid {
		if o.Value != nil {
			return fmt.Errorf("%w: %s observation must not contain a value", ErrInvalid, o.Status)
		}
	} else if err := o.Value.Validate(); err != nil {
		return err
	}
	if o.Source == Calculation && o.Status != Derived && o.Status != Unavailable && o.Status != Invalid {
		return fmt.Errorf("%w: calculation cannot be observed or reported", ErrInvalid)
	}
	return nil
}

// Validate checks the graph and acquisition interval. It accepts partial snapshots.
// Persisted schema versions and parameter-specific rules belong to later layers.
func (s Snapshot) Validate() error {
	for _, field := range []struct{ name, value string }{
		{"schema_version", s.SchemaVersion}, {"tool_version", s.ToolVersion},
		{"snapshot_id", s.ID}, {"platform.os", s.Platform.OS}, {"platform.arch", s.Platform.Arch},
	} {
		if err := required(field.name, field.value); err != nil {
			return err
		}
	}
	if s.CollectionStartedAt.IsZero() || s.CollectionFinishedAt.IsZero() || s.CollectionFinishedAt.Before(s.CollectionStartedAt) {
		return fmt.Errorf("%w: invalid collection interval", ErrInvalid)
	}
	devices := make(map[string]Device, len(s.Devices))
	for i, d := range s.Devices {
		if err := required("device.id", d.ID); err != nil {
			return fmt.Errorf("device[%d]: %w", i, err)
		}
		if err := required("device.kind", d.Kind); err != nil {
			return fmt.Errorf("device[%d]: %w", i, err)
		}
		if _, exists := devices[d.ID]; exists {
			return fmt.Errorf("%w: duplicate device %q", ErrInvalid, d.ID)
		}
		devices[d.ID] = d
	}
	for _, d := range s.Devices {
		if d.ParentID != "" {
			if _, exists := devices[d.ParentID]; !exists {
				return fmt.Errorf("%w: unknown parent %q", ErrInvalid, d.ParentID)
			}
		}
	}
	// Iterative graph traversal avoids recursion on imported device trees.
	finished := make(map[string]bool, len(devices))
	for _, d := range s.Devices {
		path := make(map[string]bool)
		for id := d.ID; id != "" && !finished[id]; id = devices[id].ParentID {
			if path[id] {
				return fmt.Errorf("%w: device parent cycle at %q", ErrInvalid, id)
			}
			path[id] = true
		}
		for id := range path {
			finished[id] = true
		}
	}
	check := func(o Observation) error {
		if err := o.Validate(); err != nil {
			return err
		}
		if _, exists := devices[o.DeviceID]; !exists {
			return fmt.Errorf("%w: unknown observation device %q", ErrInvalid, o.DeviceID)
		}
		if o.CapturedAt.Before(s.CollectionStartedAt) || o.CapturedAt.After(s.CollectionFinishedAt) {
			return fmt.Errorf("%w: observation outside collection interval", ErrInvalid)
		}
		return nil
	}
	for i, o := range s.Observations {
		if err := check(o); err != nil {
			return fmt.Errorf("observation[%d]: %w", i, err)
		}
	}
	profiles := make(map[[2]string]bool, len(s.Profiles))
	for i, p := range s.Profiles {
		if err := required("profile.id", p.ID); err != nil {
			return fmt.Errorf("profile[%d]: %w", i, err)
		}
		if err := required("profile.type", p.Type); err != nil {
			return fmt.Errorf("profile[%d]: %w", i, err)
		}
		if _, exists := devices[p.DeviceID]; !exists {
			return fmt.Errorf("%w: unknown profile device %q", ErrInvalid, p.DeviceID)
		}
		key := [2]string{p.DeviceID, p.ID}
		if profiles[key] {
			return fmt.Errorf("%w: duplicate profile %q on %q", ErrInvalid, p.ID, p.DeviceID)
		}
		profiles[key] = true
		for j, o := range p.Observations {
			if o.DeviceID != p.DeviceID || o.Scope != "profile" {
				return fmt.Errorf("%w: profile[%d] observation[%d] has inconsistent device or scope", ErrInvalid, i, j)
			}
			if err := check(o); err != nil {
				return fmt.Errorf("profile[%d] observation[%d]: %w", i, j, err)
			}
		}
	}
	return nil
}
