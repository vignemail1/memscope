package recommend

import (
	"time"
)

// Recommendation represents a single BIOS memory configuration recommendation
type Recommendation struct {
	RuleID           string                 `json:"rule_id"`
	RuleVersion      string                 `json:"rule_version"`
	Category         string                 `json:"category"`        // manufacturer_profile, experimental_candidate, tested_configuration
	Status           string                 `json:"status"`          // candidate, insufficient_data, unsupported, rejected
	Prerequisites    []string               `json:"prerequisites"`
	Evidence         []Evidence             `json:"evidence"`
	ProposedSettings map[string]interface{} `json:"proposed_settings"`
	MissingData      []string               `json:"missing_data,omitempty"`
	Warnings         []string               `json:"warnings,omitempty"`
	ValidationPlan   *ValidationPlan        `json:"validation_plan,omitempty"`
	Timestamp        time.Time              `json:"timestamp"`
}

// Evidence represents supporting evidence for a recommendation
type Evidence struct {
	Type        string      `json:"type"`        // profile_data, qvl_entry, test_result, manufacturer_spec
	Source      string      `json:"source"`      // URL, document reference, test identifier
	Description string      `json:"description"`
	Data        interface{} `json:"data,omitempty"`
	Confidence  float32     `json:"confidence"`  // 0.0 to 1.0
}

// ValidationPlan represents testing steps for validating a recommendation
type ValidationPlan struct {
	Steps       []ValidationStep `json:"steps"`
	Tools       []string         `json:"recommended_tools"`
	Duration    string           `json:"estimated_duration"`
	Precautions []string         `json:"precautions"`
}

// ValidationStep represents a single validation step
type ValidationStep struct {
	Order       int    `json:"order"`
	Action      string `json:"action"`
	Expected    string `json:"expected_result"`
	FailAction  string `json:"fail_action"`
	Required    bool   `json:"required"`
}

// RecommendationRequest represents input for generating recommendations
type RecommendationRequest struct {
	Snapshot      interface{}            `json:"snapshot"`
	Preferences   *UserPreferences       `json:"preferences,omitempty"`
	ExcludeRules  []string              `json:"exclude_rules,omitempty"`
	IncludeExperimental bool           `json:"include_experimental"`
}

// UserPreferences represents user-configurable recommendation preferences
type UserPreferences struct {
	PrioritizeStability bool    `json:"prioritize_stability"`
	MaxVoltageIncrease  float32 `json:"max_voltage_increase_v"`
	AllowExperimental   bool    `json:"allow_experimental"`
	PreferredManufacturer string `json:"preferred_manufacturer,omitempty"`
}