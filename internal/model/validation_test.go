package model

import (
	"errors"
	"testing"
	"time"
)

func validSnapshot() Snapshot {
	now := time.Date(2026, 10, 9, 6, 0, 0, 0, time.UTC)
	return Snapshot{
		SchemaVersion: "1", ToolVersion: "test", ID: "test",
		CollectionStartedAt: now, CollectionFinishedAt: now,
		Platform: Platform{OS: "windows", Arch: "amd64"},
		Devices: []Device{{ID: "imc0", Kind: "controller"}},
		Observations: []Observation{{
			DeviceID: "imc0", Scope: "controller", Parameter: "tCL",
			Value: NewUnsigned(46), Unit: "cycles", Source: Controller,
			Status: Observed, CapturedAt: now,
		}},
	}
}

func TestSnapshotValidate(t *testing.T) {
	if err := validSnapshot().Validate(); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name   string
		mutate func(*Snapshot)
	}{
		{"missing id", func(s *Snapshot) { s.ID = " " }},
		{"reversed interval", func(s *Snapshot) { s.CollectionFinishedAt = s.CollectionStartedAt.Add(-time.Second) }},
		{"duplicate device", func(s *Snapshot) { s.Devices = append(s.Devices, s.Devices[0]) }},
		{"unknown parent", func(s *Snapshot) { s.Devices[0].ParentID = "missing" }},
		{"self cycle", func(s *Snapshot) { s.Devices[0].ParentID = "imc0" }},
		{"two node cycle", func(s *Snapshot) {
			s.Devices[0].ParentID = "other"
			s.Devices = append(s.Devices, Device{ID: "other", Kind: "controller", ParentID: "imc0"})
		}},
		{"unknown device", func(s *Snapshot) { s.Observations[0].DeviceID = "missing" }},
		{"unknown status", func(s *Snapshot) { s.Observations[0].Status = "unknown" }},
		{"unknown source", func(s *Snapshot) { s.Observations[0].Source = "unknown" }},
		{"missing value", func(s *Snapshot) { s.Observations[0].Value = nil }},
		{"unavailable with value", func(s *Snapshot) { s.Observations[0].Status = Unavailable }},
		{"calculation observed", func(s *Snapshot) { s.Observations[0].Source = Calculation }},
		{"outside interval", func(s *Snapshot) { s.Observations[0].CapturedAt = s.CollectionStartedAt.Add(-time.Second) }},
		{"missing timestamp", func(s *Snapshot) { s.Observations[0].CapturedAt = time.Time{} }},
		{"profile unknown device", func(s *Snapshot) { s.Profiles = []Profile{{ID: "p", Type: "XMP", DeviceID: "missing"}} }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := validSnapshot()
			tt.mutate(&s)
			if err := s.Validate(); !errors.Is(err, ErrInvalid) {
				t.Fatalf("got %v, want ErrInvalid", err)
			}
		})
	}
}

func TestPartialSnapshot(t *testing.T) {
	s := validSnapshot()
	s.Observations[0].Status = Unavailable
	s.Observations[0].Value = nil
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Observations = nil
	s.Devices = nil
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestProfileValidation(t *testing.T) {
	s := validSnapshot()
	o := s.Observations[0]
	o.Scope, o.Source, o.Status = "profile", SPD, Reported
	s.Profiles = []Profile{{ID: "p1", DeviceID: "imc0", Type: "XMP", Observations: []Observation{o}}}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Profiles = append(s.Profiles, s.Profiles[0])
	if !errors.Is(s.Validate(), ErrInvalid) {
		t.Fatal("duplicate profile accepted")
	}
	s.Profiles = s.Profiles[:1]
	s.Profiles[0].Observations[0].Scope = "controller"
	if !errors.Is(s.Validate(), ErrInvalid) {
		t.Fatal("wrong profile scope accepted")
	}
}
