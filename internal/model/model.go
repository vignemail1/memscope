// Package model defines platform-independent acquisition contracts.
// These types are internal; the persisted JSON schema is not implemented yet.
package model

import "time"

type Status string

const (
	Observed    Status = "observed"
	Reported    Status = "reported"
	Derived     Status = "derived"
	Unavailable Status = "unavailable"
	Invalid     Status = "invalid"
)

type Source string

const (
	SMBIOS      Source = "smbios"
	WMI         Source = "wmi"
	CPUID       Source = "cpuid"
	SPD         Source = "spd"
	Controller  Source = "controller"
	Telemetry   Source = "telemetry"
	Calculation Source = "calculation"
	Import      Source = "import"
	User        Source = "user"
)

// Value holds exactly one typed value when present. Validation is a future task.
// A nil observation Value represents missing data, never a measured zero.
type Value struct {
	Text     *string
	Unsigned *uint64
	Decimal  *float64
	Boolean  *bool
}

type Device struct {
	ID       string
	Kind     string
	ParentID string
}

type Diagnostic struct {
	Code      string
	Severity  string
	Message   string
	Source    Source
	DeviceID  string
	Parameter string
}

type Observation struct {
	DeviceID      string
	Scope         string
	Parameter     string
	Value         *Value
	Unit          string
	Source        Source
	SourceVersion string
	Status        Status
	CapturedAt    time.Time
	Diagnostics   []Diagnostic
}

type Profile struct {
	ID           string
	DeviceID     string
	Type         string
	Version      string
	Observations []Observation
	Diagnostics  []Diagnostic
}

type Capability struct {
	Name      string
	Available bool
	Reason    string
}

type Platform struct {
	OS   string
	Arch string
}

type Snapshot struct {
	SchemaVersion        string
	ToolVersion          string
	ID                   string
	CollectionStartedAt  time.Time
	CollectionFinishedAt time.Time
	Platform             Platform
	Devices              []Device
	Observations         []Observation
	Profiles             []Profile
	Diagnostics          []Diagnostic
	Capabilities         []Capability
}
