package recommend

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

// Engine implements the recommendation generation logic
type Engine struct {
	rules []Rule
}

// Rule represents a recommendation rule
type Rule struct {
	ID          string
	Version     string
	Name        string
	Category    string
	Selector    PlatformSelector
	Logic       RuleLogic
	Evidence    []Evidence
	Settings    map[string]interface{}
	Warnings    []string
	Validation  *ValidationPlan
}

// PlatformSelector defines what platforms a rule applies to
type PlatformSelector struct {
	CPUFamily        []string `json:"cpu_family,omitempty"`
	CPUModel         []string `json:"cpu_model,omitempty"`
	ChipsetFamily    []string `json:"chipset_family,omitempty"`
	MotherboardModel []string `json:"motherboard_model,omitempty"`
	BIOSVersion      []string `json:"bios_version,omitempty"`
	MemoryType       []string `json:"memory_type,omitempty"`
}

// RuleLogic defines the logic for applying a rule
type RuleLogic struct {
	RequiredData    []string               `json:"required_data"`
	Conditions      []Condition            `json:"conditions"`
	Actions         []Action               `json:"actions"`
	Transformations map[string]interface{} `json:"transformations,omitempty"`
}

// Condition represents a rule condition
type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Action represents a rule action
type Action struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

// NewEngine creates a new recommendation engine with default rules
func NewEngine() *Engine {
	engine := &Engine{
		rules: []Rule{},
	}
	
	// Load default rules
	engine.loadDefaultRules()
	
	return engine
}

// GenerateRecommendations analyzes a snapshot and generates recommendations
func (e *Engine) GenerateRecommendations(snapshot *model.Snapshot) ([]*Recommendation, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot cannot be nil")
	}
	
	var recommendations []*Recommendation
	timestamp := time.Now()
	
	// Analyze system compatibility
	systemInfo := e.analyzeSystem(snapshot)
	
	// Check each rule for applicability
	for _, rule := range e.rules {
		rec := e.evaluateRule(rule, snapshot, systemInfo, timestamp)
		if rec != nil {
			recommendations = append(recommendations, rec)
		}
	}
	
	// If no specific recommendations, generate insufficient data response
	if len(recommendations) == 0 {
		recommendations = append(recommendations, &Recommendation{
			RuleID:      "default_insufficient_data",
			RuleVersion: "1.0",
			Category:    "system_analysis",
			Status:      "insufficient_data",
			MissingData: e.identifyMissingData(snapshot),
			Timestamp:   timestamp,
		})
	}
	
	return recommendations, nil
}

// analyzeSystem extracts key system information for rule matching
func (e *Engine) analyzeSystem(snapshot *model.Snapshot) map[string]interface{} {
	info := make(map[string]interface{})
	
	// Extract CPU information
	cpuName := e.findObservationValue(snapshot, "cpu", "name")
	if cpuName != "" {
		info["cpu"] = cpuName
	}
	
	// Extract memory information
	var memoryDevices []model.Device
	for _, device := range snapshot.Devices {
		if device.Kind == "memory" {
			memoryDevices = append(memoryDevices, device)
		}
	}
	info["memory_devices"] = memoryDevices
	info["memory_count"] = len(memoryDevices)
	
	// Extract BIOS information
	biosVersion := e.findObservationValue(snapshot, "bios", "version")
	if biosVersion != "" {
		info["bios_version"] = biosVersion
	}
	
	// Extract motherboard information from system device
	systemModel := e.findObservationValue(snapshot, "system", "model")
	if systemModel != "" {
		info["motherboard"] = systemModel
	}
	
	return info
}

// findObservationValue finds an observation value for a specific device kind and parameter
func (e *Engine) findObservationValue(snapshot *model.Snapshot, deviceKind, parameter string) string {
	// Find device of the specified kind
	var deviceID string
	for _, device := range snapshot.Devices {
		if device.Kind == deviceKind {
			deviceID = device.ID
			break
		}
	}
	
	if deviceID == "" {
		return ""
	}
	
	// Find observation for this device and parameter
	for _, obs := range snapshot.Observations {
		if obs.DeviceID == deviceID && obs.Parameter == parameter {
			if obs.Value != nil && obs.Value.Text != nil {
				return *obs.Value.Text
			}
		}
	}
	
	return ""
}

// evaluateRule checks if a rule applies and generates a recommendation
func (e *Engine) evaluateRule(rule Rule, snapshot *model.Snapshot, systemInfo map[string]interface{}, timestamp time.Time) *Recommendation {
	// Check platform compatibility
	if !e.matchesPlatform(rule.Selector, systemInfo) {
		return nil
	}
	
	// Check if we have required data
	missingData := e.checkRequiredData(rule.Logic.RequiredData, systemInfo)
	if len(missingData) > 0 {
		return &Recommendation{
			RuleID:       rule.ID,
			RuleVersion:  rule.Version,
			Category:     rule.Category,
			Status:       "insufficient_data",
			MissingData:  missingData,
			Timestamp:    timestamp,
		}
	}
	
	// Evaluate conditions
	if !e.evaluateConditions(rule.Logic.Conditions, systemInfo) {
		return &Recommendation{
			RuleID:      rule.ID,
			RuleVersion: rule.Version,
			Category:    rule.Category,
			Status:      "rejected",
			Evidence:    []Evidence{{Type: "condition_check", Description: "Rule conditions not met"}},
			Timestamp:   timestamp,
		}
	}
	
	// Generate candidate recommendation
	return &Recommendation{
		RuleID:           rule.ID,
		RuleVersion:      rule.Version,
		Category:         rule.Category,
		Status:           "candidate",
		Evidence:         rule.Evidence,
		ProposedSettings: rule.Settings,
		Warnings:         rule.Warnings,
		ValidationPlan:   rule.Validation,
		Timestamp:        timestamp,
	}
}

// loadDefaultRules loads the initial set of recommendation rules
func (e *Engine) loadDefaultRules() {
	// Add basic XMP profile detection rule
	xmpRule := Rule{
		ID:       "xmp_profile_detection",
		Version:  "1.0",
		Name:     "XMP Profile Analysis",
		Category: "manufacturer_profile",
		Selector: PlatformSelector{
			CPUFamily: []string{"Intel", "AMD"},
		},
		Logic: RuleLogic{
			RequiredData: []string{"memory_devices"},
			Conditions: []Condition{
				{Field: "memory_count", Operator: "gt", Value: 0},
			},
		},
		Evidence: []Evidence{
			{
				Type:        "profile_data",
				Source:      "SPD_EEPROM",
				Description: "XMP profile data from memory module SPD",
				Confidence:  0.9,
			},
		},
		Settings: map[string]interface{}{
			"enable_xmp": true,
			"profile":    "XMP_1",
		},
		Warnings: []string{
			"XMP profiles may require BIOS support verification",
			"Some motherboards may need manual timing adjustment",
		},
		Validation: &ValidationPlan{
			Steps: []ValidationStep{
				{Order: 1, Action: "Enable XMP in BIOS", Expected: "System boots successfully", Required: true},
				{Order: 2, Action: "Run memory stress test", Expected: "No errors for 30 minutes", Required: true},
				{Order: 3, Action: "Monitor temperatures", Expected: "Within safe operating range", Required: false},
			},
			Tools:    []string{"MemTest86", "Prime95", "HWiNFO64"},
			Duration: "2-4 hours",
			Precautions: []string{
				"Keep BIOS recovery procedures ready",
				"Monitor system stability during testing",
			},
		},
	}
	
	e.rules = append(e.rules, xmpRule)
	
	// Add basic stability rule
	stabilityRule := Rule{
		ID:       "jedec_stability",
		Version:  "1.0",
		Name:     "JEDEC Stability Baseline",
		Category: "manufacturer_profile",
		Selector: PlatformSelector{
			CPUFamily: []string{"Intel", "AMD"},
		},
		Logic: RuleLogic{
			RequiredData: []string{"memory_devices"},
			Conditions: []Condition{
				{Field: "memory_count", Operator: "gt", Value: 0},
			},
		},
		Evidence: []Evidence{
			{
				Type:        "manufacturer_spec",
				Source:      "JEDEC_STANDARD",
				Description: "JEDEC specification compliance",
				Confidence:  1.0,
			},
		},
		Settings: map[string]interface{}{
			"use_jedec_timings": true,
			"voltage":           1.2,
		},
		Warnings: []string{
			"JEDEC settings prioritize stability over performance",
		},
		Validation: &ValidationPlan{
			Steps: []ValidationStep{
				{Order: 1, Action: "Apply JEDEC timings in BIOS", Expected: "System boots successfully", Required: true},
				{Order: 2, Action: "Run stability test", Expected: "No crashes or errors", Required: true},
			},
			Tools:    []string{"MemTest86", "Windows Memory Diagnostic"},
			Duration: "1-2 hours",
			Precautions: []string{
				"JEDEC timings should be universally compatible",
			},
		},
	}
	
	e.rules = append(e.rules, stabilityRule)
}

// Helper methods for rule evaluation
func (e *Engine) matchesPlatform(selector PlatformSelector, systemInfo map[string]interface{}) bool {
	// Basic platform matching logic
	if len(selector.CPUFamily) > 0 {
		cpu, ok := systemInfo["cpu"].(string)
		if !ok {
			return false
		}
		
		match := false
		for _, family := range selector.CPUFamily {
			if containsIgnoreCase(cpu, family) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	
	return true
}

func (e *Engine) checkRequiredData(required []string, systemInfo map[string]interface{}) []string {
	var missing []string
	
	for _, req := range required {
		if _, exists := systemInfo[req]; !exists {
			missing = append(missing, req)
		}
	}
	
	return missing
}

func (e *Engine) evaluateConditions(conditions []Condition, systemInfo map[string]interface{}) bool {
	for _, condition := range conditions {
		if !e.evaluateCondition(condition, systemInfo) {
			return false
		}
	}
	return true
}

func (e *Engine) evaluateCondition(condition Condition, systemInfo map[string]interface{}) bool {
	value, exists := systemInfo[condition.Field]
	if !exists {
		return false
	}
	
	switch condition.Operator {
	case "gt":
		if intVal, ok := value.(int); ok {
			if targetVal, ok := condition.Value.(int); ok {
				return intVal > targetVal
			}
		}
	case "eq":
		return value == condition.Value
	}
	
	return false
}

func (e *Engine) identifyMissingData(snapshot *model.Snapshot) []string {
	var missing []string
	
	// Check for CPU information
	if e.findObservationValue(snapshot, "cpu", "name") == "" {
		missing = append(missing, "cpu_model")
	}
	
	// Check for motherboard/system information
	if e.findObservationValue(snapshot, "system", "model") == "" {
		missing = append(missing, "motherboard_model")
	}
	
	// Check for BIOS information
	if e.findObservationValue(snapshot, "bios", "version") == "" {
		missing = append(missing, "bios_version")
	}
	
	// Check for memory modules
	hasMemory := false
	for _, device := range snapshot.Devices {
		if device.Kind == "memory" {
			hasMemory = true
			break
		}
	}
	
	if !hasMemory {
		missing = append(missing, "memory_modules")
	}
	
	return missing
}

// Utility functions
func containsIgnoreCase(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}