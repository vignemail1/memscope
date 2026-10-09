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
	Text     *string  `json:"text,omitempty"`
	Unsigned *uint64  `json:"unsigned,omitempty"`
	Decimal  *float64 `json:"decimal,omitempty"`
	Boolean  *bool    `json:"boolean,omitempty"`
}

type Device struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	ParentID string `json:"parent_id,omitempty"`
}

type Diagnostic struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Source    Source `json:"source"`
	DeviceID  string `json:"device_id,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type Observation struct {
	DeviceID      string       `json:"device_id"`
	Scope         string       `json:"scope"`
	Parameter     string       `json:"parameter"`
	Value         *Value       `json:"value,omitempty"`
	Unit          string       `json:"unit,omitempty"`
	Source        Source       `json:"source"`
	SourceVersion string       `json:"source_version,omitempty"`
	Status        Status       `json:"status"`
	CapturedAt    time.Time    `json:"captured_at"`
	Diagnostics   []Diagnostic `json:"diagnostics,omitempty"`
}

type Profile struct {
	ID           string        `json:"id"`
	DeviceID     string        `json:"device_id"`
	Type         string        `json:"type"`
	Version      string        `json:"version"`
	Observations []Observation `json:"observations,omitempty"`
	Diagnostics  []Diagnostic  `json:"diagnostics,omitempty"`
}

type Capability struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type Snapshot struct {
	SchemaVersion        string       `json:"schema_version"`
	ToolVersion          string       `json:"tool_version"`
	ID                   string       `json:"id"`
	CollectionStartedAt  time.Time    `json:"collection_started_at"`
	CollectionFinishedAt time.Time    `json:"collection_finished_at"`
	Platform             Platform     `json:"platform"`
	Devices              []Device     `json:"devices"`
	Observations         []Observation `json:"observations"`
	Profiles             []Profile    `json:"profiles"`
	Diagnostics          []Diagnostic `json:"diagnostics"`
	Capabilities         []Capability `json:"capabilities"`
}
