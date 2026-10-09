package recommend

import (
	"testing"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

func TestRecommendationEngine(t *testing.T) {
	engine := NewEngine()
	
	// Create sample snapshot with memory data
	snapshot := &model.Snapshot{
		SchemaVersion: "1.0.0",
		ID:           "test-snapshot",
		CollectionStartedAt: time.Now(),
		CollectionFinishedAt: time.Now(),
		Devices: []model.Device{
			{
				ID:   "cpu-0",
				Kind: "cpu",
			},
			{
				ID:   "memory-0",
				Kind: "memory",
			},
		},
		Observations: []model.Observation{
			{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "name",
				Value:      &model.Value{Text: stringPtr("AMD Ryzen 5 3600")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "manufacturer",
				Value:      &model.Value{Text: stringPtr("Corsair")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "part_number",
				Value:      &model.Value{Text: stringPtr("CMK16GX4M2B3200C16")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "capacity",
				Value:      &model.Value{Unsigned: uint64Ptr(17179869184)}, // 16GB
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
			{
				DeviceID:   "memory-0",
				Scope:      "hardware",
				Parameter:  "speed",
				Value:      &model.Value{Unsigned: uint64Ptr(3200)},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
		},
	}
	
	recommendations, err := engine.GenerateRecommendations(snapshot)
	if err != nil {
		t.Fatalf("GenerateRecommendations failed: %v", err)
	}
	
	if len(recommendations) == 0 {
		t.Error("Should generate at least one recommendation")
	}
	
	// Validate recommendation structure
	rec := recommendations[0]
	if rec.RuleID == "" {
		t.Error("Recommendation should have rule ID")
	}
	if rec.Category == "" {
		t.Error("Recommendation should have category")
	}
	if rec.Status == "" {
		t.Error("Recommendation should have status")
	}
}

func TestInsufficientDataHandling(t *testing.T) {
	engine := NewEngine()
	
	// Create incomplete snapshot
	snapshot := &model.Snapshot{
		SchemaVersion: "1.0.0",
		ID:           "incomplete-snapshot",
		CollectionStartedAt: time.Now(),
		CollectionFinishedAt: time.Now(),
		Observations: []model.Observation{
			{
				DeviceID:   "cpu-0",
				Scope:      "hardware",
				Parameter:  "name",
				Value:      &model.Value{Text: stringPtr("Unknown CPU")},
				Source:     model.WMI,
				Status:     model.Observed,
				CapturedAt: time.Now(),
			},
		},
	}
	
	recommendations, err := engine.GenerateRecommendations(snapshot)
	if err != nil {
		t.Fatalf("Should handle insufficient data gracefully: %v", err)
	}
	
	// Should return insufficient_data status
	found := false
	for _, rec := range recommendations {
		if rec.Status == "insufficient_data" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Should return insufficient_data status for incomplete snapshot")
	}
}

func TestRecommendationPolicyCompliance(t *testing.T) {
	engine := NewEngine()
	
	// Create complete snapshot
	snapshot := &model.Snapshot{
		SchemaVersion: "1.0.0",
		ID:           "policy-test-snapshot",
		CollectionStartedAt: time.Now(),
		CollectionFinishedAt: time.Now(),
		Devices: []model.Device{
			{ID: "cpu-0", Kind: "cpu"},
			{ID: "system-0", Kind: "system"},
			{ID: "bios-0", Kind: "bios"},
			{ID: "memory-0", Kind: "memory"},
		},
		Observations: []model.Observation{
			{
				DeviceID: "cpu-0", Scope: "hardware", Parameter: "name",
				Value: &model.Value{Text: stringPtr("AMD Ryzen 5 3600")},
				Source: model.WMI, Status: model.Observed, CapturedAt: time.Now(),
			},
			{
				DeviceID: "system-0", Scope: "hardware", Parameter: "model",
				Value: &model.Value{Text: stringPtr("ASUS B450-F")},
				Source: model.WMI, Status: model.Observed, CapturedAt: time.Now(),
			},
			{
				DeviceID: "bios-0", Scope: "firmware", Parameter: "version",
				Value: &model.Value{Text: stringPtr("4021")},
				Source: model.WMI, Status: model.Observed, CapturedAt: time.Now(),
			},
			{
				DeviceID: "memory-0", Scope: "hardware", Parameter: "manufacturer",
				Value: &model.Value{Text: stringPtr("Corsair")},
				Source: model.WMI, Status: model.Observed, CapturedAt: time.Now(),
			},
		},
	}
	
	recommendations, err := engine.GenerateRecommendations(snapshot)
	if err != nil {
		t.Fatalf("GenerateRecommendations failed: %v", err)
	}
	
	// Verify policy compliance for each recommendation
	for _, rec := range recommendations {
		// Every recommendation must have required fields per policy
		if rec.RuleID == "" {
			t.Error("Missing rule_id - required by recommendation policy")
		}
		if rec.RuleVersion == "" {
			t.Error("Missing rule_version - required by recommendation policy")
		}
		if rec.Category == "" {
			t.Error("Missing category - required by recommendation policy")
		}
		if rec.Status == "" {
			t.Error("Missing status - required by recommendation policy")
		}
		
		// Category must be one of the allowed values
		validCategories := []string{"manufacturer_profile", "experimental_candidate", "tested_configuration", "system_analysis"}
		validCategory := false
		for _, cat := range validCategories {
			if rec.Category == cat {
				validCategory = true
				break
			}
		}
		if !validCategory {
			t.Errorf("Invalid category %s - must be one of: %v", rec.Category, validCategories)
		}
		
		// Status must be one of the allowed values
		validStatuses := []string{"candidate", "insufficient_data", "unsupported", "rejected"}
		validStatus := false
		for _, status := range validStatuses {
			if rec.Status == status {
				validStatus = true
				break
			}
		}
		if !validStatus {
			t.Errorf("Invalid status %s - must be one of: %v", rec.Status, validStatuses)
		}
		
		// Candidate recommendations should have evidence and settings
		if rec.Status == "candidate" {
			if len(rec.Evidence) == 0 {
				t.Error("Candidate recommendations must have evidence")
			}
			if len(rec.ProposedSettings) == 0 {
				t.Error("Candidate recommendations must have proposed settings")
			}
			if rec.ValidationPlan == nil {
				t.Error("Candidate recommendations must have validation plan")
			}
		}
		
		// Insufficient data recommendations should list missing data
		if rec.Status == "insufficient_data" {
			if len(rec.MissingData) == 0 {
				t.Error("Insufficient data recommendations must specify missing data")
			}
		}
		
		// Evidence should have proper confidence values
		for _, evidence := range rec.Evidence {
			if evidence.Confidence < 0.0 || evidence.Confidence > 1.0 {
				t.Errorf("Evidence confidence must be between 0.0 and 1.0, got %f", evidence.Confidence)
			}
			if evidence.Type == "" {
				t.Error("Evidence must have type")
			}
			if evidence.Description == "" {
				t.Error("Evidence must have description")
			}
		}
		
		// Validation plans should be structured
		if rec.ValidationPlan != nil {
			if len(rec.ValidationPlan.Steps) == 0 {
				t.Error("Validation plan must have steps")
			}
			if rec.ValidationPlan.Duration == "" {
				t.Error("Validation plan must have duration estimate")
			}
			
			for _, step := range rec.ValidationPlan.Steps {
				if step.Action == "" {
					t.Error("Validation step must have action")
				}
				if step.Expected == "" {
					t.Error("Validation step must have expected result")
				}
			}
		}
	}
}

func TestAdvancedConditions(t *testing.T) {
	engine := NewEngine()
	
	testCases := []struct {
		name      string
		condition Condition
		systemInfo map[string]interface{}
		expected   bool
	}{
		{
			name: "contains operator",
			condition: Condition{Field: "cpu", Operator: "contains", Value: "ryzen"},
			systemInfo: map[string]interface{}{"cpu": "AMD Ryzen 5 3600"},
			expected: true,
		},
		{
			name: "in operator",
			condition: Condition{Field: "memory_type", Operator: "in", Value: []interface{}{"DDR4", "DDR5"}},
			systemInfo: map[string]interface{}{"memory_type": "DDR4"},
			expected: true,
		},
		{
			name: "gte operator",
			condition: Condition{Field: "memory_count", Operator: "gte", Value: 2},
			systemInfo: map[string]interface{}{"memory_count": 4},
			expected: true,
		},
		{
			name: "lte operator",
			condition: Condition{Field: "memory_count", Operator: "lte", Value: 4},
			systemInfo: map[string]interface{}{"memory_count": 4},
			expected: true,
		},
		{
			name: "ne operator",
			condition: Condition{Field: "status", Operator: "ne", Value: "failed"},
			systemInfo: map[string]interface{}{"status": "success"},
			expected: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := engine.evaluateCondition(tc.condition, tc.systemInfo)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestExpandedRuleSet(t *testing.T) {
	engine := NewEngine()
	
	// Should have more than just the 2 basic rules
	if len(engine.rules) < 4 {
		t.Errorf("Expected at least 4 rules, got %d", len(engine.rules))
	}
	
	// Check for specific rules
	ruleIDs := make(map[string]bool)
	for _, rule := range engine.rules {
		ruleIDs[rule.ID] = true
	}
	
	expectedRules := []string{
		"xmp_profile_detection",
		"jedec_stability",
		"memory_compatibility_check",
		"overclocking_stability",
		"thermal_management",
	}
	
	for _, expectedRule := range expectedRules {
		if !ruleIDs[expectedRule] {
			t.Errorf("Missing expected rule: %s", expectedRule)
		}
	}
}

func TestEnhancedPlatformMatching(t *testing.T) {
	engine := NewEngine()
	
	testCases := []struct {
		name       string
		selector   PlatformSelector
		systemInfo map[string]interface{}
		expected   bool
	}{
		{
			name: "CPU family match",
			selector: PlatformSelector{CPUFamily: []string{"AMD"}},
			systemInfo: map[string]interface{}{"cpu": "AMD Ryzen 5 3600"},
			expected: true,
		},
		{
			name: "CPU model match", 
			selector: PlatformSelector{CPUModel: []string{"3600"}},
			systemInfo: map[string]interface{}{"cpu": "AMD Ryzen 5 3600"},
			expected: true,
		},
		{
			name: "Motherboard match",
			selector: PlatformSelector{MotherboardModel: []string{"B450"}},
			systemInfo: map[string]interface{}{"motherboard": "ASUS B450-F Gaming"},
			expected: true,
		},
		{
			name: "BIOS version match",
			selector: PlatformSelector{BIOSVersion: []string{"4021"}},
			systemInfo: map[string]interface{}{"bios_version": "4021"},
			expected: true,
		},
		{
			name: "No match when missing data",
			selector: PlatformSelector{CPUFamily: []string{"Intel"}},
			systemInfo: map[string]interface{}{"cpu": "AMD Ryzen 5 3600"},
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := engine.matchesPlatform(tc.selector, tc.systemInfo)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// Helper functions for test data creation
func stringPtr(s string) *string {
	return &s
}

func uint64Ptr(u uint64) *uint64 {
	return &u
}