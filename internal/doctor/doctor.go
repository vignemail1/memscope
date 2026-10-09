package doctor

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/vignemail1/memscope/internal/model"
)

// Doctor performs comprehensive system diagnostics
type Doctor struct {
	options *DiagnosticOptions
}

// NewDoctor creates a new system doctor with default options
func NewDoctor() *Doctor {
	return &Doctor{
		options: &DiagnosticOptions{
			IncludePerformanceAnalysis: true,
			IncludeCompatibilityCheck:  true,
			SkipSlowChecks:            false,
			MinSeverity:               "low",
		},
	}
}

// NewDoctorWithOptions creates a doctor with custom options
func NewDoctorWithOptions(options *DiagnosticOptions) *Doctor {
	return &Doctor{options: options}
}

// RunDiagnostics performs comprehensive system analysis
func (d *Doctor) RunDiagnostics(snapshot *model.Snapshot) (*DiagnosticResult, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot cannot be nil")
	}
	
	timestamp := time.Now()
	result := &DiagnosticResult{
		Timestamp: timestamp,
		Checks:    []*DiagnosticCheck{},
		SystemInfo: d.extractSystemInfo(snapshot),
	}
	
	// Run all diagnostic checks
	systemChecks := d.runSystemChecks(snapshot)
	result.Checks = append(result.Checks, systemChecks...)
	
	memoryChecks := d.runMemoryChecks(snapshot)
	result.Checks = append(result.Checks, memoryChecks...)
	
	if d.options.IncludePerformanceAnalysis {
		performanceChecks := d.runPerformanceChecks(snapshot)
		result.Checks = append(result.Checks, performanceChecks...)
	}
	
	if d.options.IncludeCompatibilityCheck {
		compatibilityChecks := d.runCompatibilityChecks(snapshot)
		result.Checks = append(result.Checks, compatibilityChecks...)
	}
	
	// Generate summary and score
	result.Summary = d.generateSummary(result.Checks)
	result.OverallScore = d.calculateOverallScore(result.Summary)
	result.HealthStatus = d.determineHealthStatus(result.OverallScore)
	
	// Generate recommendations
	result.Recommendations = d.generateRecommendations(result.Checks, snapshot)
	
	return result, nil
}

// runSystemChecks performs basic system health checks
func (d *Doctor) runSystemChecks(snapshot *model.Snapshot) []*DiagnosticCheck {
	var checks []*DiagnosticCheck
	timestamp := time.Now()
	
	// Check if we have any system devices
	hasSystemDevice := false
	for _, device := range snapshot.Devices {
		if device.Kind == "system" {
			hasSystemDevice = true
			break
		}
	}
	
	if !hasSystemDevice {
		checks = append(checks, &DiagnosticCheck{
			Name:      "System Information",
			Category:  "system",
			Status:    "error",
			Severity:  "high",
			Message:   "System information is missing",
			Details:   "Cannot perform system-level diagnostics without basic system information",
			Timestamp: timestamp,
		})
		return checks
	}
	
	// CPU information check
	cpuName := d.getCPUName(snapshot)
	if cpuName == "" {
		checks = append(checks, &DiagnosticCheck{
			Name:      "CPU Information",
			Category:  "system",
			Status:    "warning",
			Severity:  "medium",
			Message:   "CPU information is missing",
			Suggestion: "Verify system detection is working correctly",
			Timestamp: timestamp,
		})
	} else {
		checks = append(checks, &DiagnosticCheck{
			Name:      "CPU Information",
			Category:  "system", 
			Status:    "pass",
			Severity:  "low",
			Message:   fmt.Sprintf("CPU detected: %s", cpuName),
			Timestamp: timestamp,
		})
	}
	
	// BIOS version check
	biosVersion := d.getBIOSVersion(snapshot)
	if biosVersion == "" {
		checks = append(checks, &DiagnosticCheck{
			Name:      "BIOS Information",
			Category:  "system",
			Status:    "warning",
			Severity:  "low",
			Message:   "BIOS version information is missing",
			Suggestion: "Consider updating system firmware detection",
			Timestamp: timestamp,
		})
	} else {
		checks = append(checks, &DiagnosticCheck{
			Name:      "BIOS Information",
			Category:  "system",
			Status:    "pass", 
			Severity:  "low",
			Message:   fmt.Sprintf("BIOS version: %s", biosVersion),
			Timestamp: timestamp,
		})
	}
	
	return checks
}

// runMemoryChecks performs memory-specific diagnostics
func (d *Doctor) runMemoryChecks(snapshot *model.Snapshot) []*DiagnosticCheck {
	var checks []*DiagnosticCheck
	timestamp := time.Now()
	
	memoryDevices := d.getMemoryDevices(snapshot.Devices)
	
	// Memory presence check
	if len(memoryDevices) == 0 {
		checks = append(checks, &DiagnosticCheck{
			Name:      "Memory Detection",
			Category:  "memory",
			Status:    "error",
			Severity:  "critical",
			Message:   "No memory devices detected",
			Details:   "System cannot function without memory modules",
			Timestamp: timestamp,
		})
		return checks
	}
	
	checks = append(checks, &DiagnosticCheck{
		Name:      "Memory Detection",
		Category:  "memory",
		Status:    "pass",
		Severity:  "low",
		Message:   fmt.Sprintf("Detected %d memory module(s)", len(memoryDevices)),
		Timestamp: timestamp,
	})
	
	// Check for voltage issues in SPD data
	for i, device := range memoryDevices {
		voltage := d.getMemoryVoltage(snapshot, device.ID)
		if voltage > 0 {
			voltageCheck := d.checkMemoryVoltage(voltage, i, timestamp)
			checks = append(checks, voltageCheck)
		}
	}
	
	return checks
}

// runPerformanceChecks analyzes system performance characteristics
func (d *Doctor) runPerformanceChecks(snapshot *model.Snapshot) []*DiagnosticCheck {
	var checks []*DiagnosticCheck
	timestamp := time.Now()
	
	memoryDevices := d.getMemoryDevices(snapshot.Devices)
	if len(memoryDevices) == 0 {
		return checks // Cannot analyze performance without memory info
	}
	
	// Memory speed analysis
	cpuName := d.getCPUName(snapshot)
	if cpuName != "" {
		speedCheck := d.analyzeMemorySpeed(snapshot, cpuName, memoryDevices, timestamp)
		if speedCheck != nil {
			checks = append(checks, speedCheck)
		}
	}
	
	return checks
}

// runCompatibilityChecks performs hardware compatibility analysis
func (d *Doctor) runCompatibilityChecks(snapshot *model.Snapshot) []*DiagnosticCheck {
	var checks []*DiagnosticCheck
	timestamp := time.Now()
	
	// Basic compatibility check placeholder
	checks = append(checks, &DiagnosticCheck{
		Name:      "Platform Compatibility",
		Category:  "compatibility",
		Status:    "pass",
		Severity:  "low",
		Message:   "Platform configuration appears compatible",
		Timestamp: timestamp,
	})
	
	return checks
}

// Helper methods

func (d *Doctor) getMemoryDevices(devices []model.Device) []model.Device {
	var memoryDevices []model.Device
	for _, device := range devices {
		if device.Kind == "memory" {
			memoryDevices = append(memoryDevices, device)
		}
	}
	return memoryDevices
}

func (d *Doctor) getCPUName(snapshot *model.Snapshot) string {
	for _, obs := range snapshot.Observations {
		if obs.Scope == "hardware" && obs.Parameter == "name" {
			// Find the device
			for _, device := range snapshot.Devices {
				if device.ID == obs.DeviceID && device.Kind == "cpu" {
					if obs.Value != nil && obs.Value.Text != nil {
						return *obs.Value.Text
					}
				}
			}
		}
	}
	return ""
}

func (d *Doctor) getBIOSVersion(snapshot *model.Snapshot) string {
	for _, obs := range snapshot.Observations {
		if obs.Scope == "firmware" && obs.Parameter == "version" {
			// Find the device
			for _, device := range snapshot.Devices {
				if device.ID == obs.DeviceID && device.Kind == "bios" {
					if obs.Value != nil && obs.Value.Text != nil {
						return *obs.Value.Text
					}
				}
			}
		}
	}
	return ""
}

func (d *Doctor) getMemoryVoltage(snapshot *model.Snapshot, deviceID string) float64 {
	for _, obs := range snapshot.Observations {
		if obs.DeviceID == deviceID && obs.Scope == "spd" && obs.Parameter == "voltage_level" {
			if obs.Value != nil && obs.Value.Decimal != nil {
				return *obs.Value.Decimal
			}
		}
	}
	return 0
}

func (d *Doctor) checkMemoryVoltage(voltage float64, index int, timestamp time.Time) *DiagnosticCheck {
	status := "pass"
	severity := "low"
	message := fmt.Sprintf("Memory module %d voltage: %.2fV", index+1, voltage)
	
	if voltage > 1.45 {
		status = "warning"
		severity = "high"
		message = fmt.Sprintf("Memory module %d has high voltage (%.2fV) - may reduce lifespan", index+1, voltage)
	} else if voltage < 1.0 || voltage > 1.4 {
		status = "warning"
		severity = "medium"
		message = fmt.Sprintf("Memory module %d voltage (%.2fV) is outside typical range", index+1, voltage)
	}
	
	return &DiagnosticCheck{
		Name:      fmt.Sprintf("Memory Module %d Voltage", index+1),
		Category:  "memory",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Timestamp: timestamp,
	}
}

func (d *Doctor) analyzeMemorySpeed(snapshot *model.Snapshot, cpuName string, memoryDevices []model.Device, timestamp time.Time) *DiagnosticCheck {
	// Get memory speed from observations
	var totalSpeed uint64
	validDevices := 0
	
	for _, device := range memoryDevices {
		for _, obs := range snapshot.Observations {
			if obs.DeviceID == device.ID && obs.Scope == "hardware" && obs.Parameter == "speed" {
				if obs.Value != nil && obs.Value.Unsigned != nil {
					totalSpeed += *obs.Value.Unsigned
					validDevices++
				}
			}
		}
	}
	
	if validDevices == 0 {
		return nil
	}
	
	avgSpeed := totalSpeed / uint64(validDevices)
	
	// Analyze performance based on CPU
	cpuUpper := strings.ToUpper(cpuName)
	expectedSpeed := d.getExpectedMemorySpeed(cpuUpper)
	
	status := "pass"
	severity := "low"
	message := fmt.Sprintf("Average memory speed: %d MT/s", avgSpeed)
	
	if expectedSpeed > 0 && avgSpeed < expectedSpeed {
		status = "info"
		severity = "medium"
		message = fmt.Sprintf("Memory speed (%d MT/s) is below optimal for %s (recommended: %d MT/s+)", 
			avgSpeed, cpuName, expectedSpeed)
	}
	
	return &DiagnosticCheck{
		Name:      "Memory Performance Analysis",
		Category:  "performance",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Timestamp: timestamp,
	}
}

func (d *Doctor) getExpectedMemorySpeed(cpu string) uint64 {
	// Basic CPU-based memory speed recommendations
	if strings.Contains(cpu, "RYZEN 7") || strings.Contains(cpu, "RYZEN 9") {
		return 3600 // AMD Ryzen high-end prefers faster memory
	} else if strings.Contains(cpu, "RYZEN 5") || strings.Contains(cpu, "RYZEN 3") {
		return 3200
	} else if strings.Contains(cpu, "INTEL") && (strings.Contains(cpu, "I7") || strings.Contains(cpu, "I9")) {
		return 3200
	} else if strings.Contains(cpu, "INTEL") && strings.Contains(cpu, "I5") {
		return 2666
	}
	
	return 2133 // Conservative default
}

func (d *Doctor) extractSystemInfo(snapshot *model.Snapshot) *SystemInfo {
	info := &SystemInfo{
		Platform: snapshot.Platform.OS + "/" + snapshot.Platform.Arch,
	}
	
	// Extract CPU name
	info.CPU = d.getCPUName(snapshot)
	
	// Extract BIOS version
	info.BIOSVersion = d.getBIOSVersion(snapshot)
	
	// Count memory devices
	memoryDevices := d.getMemoryDevices(snapshot.Devices)
	info.PopulatedSlots = len(memoryDevices)
	
	return info
}

func (d *Doctor) generateSummary(checks []*DiagnosticCheck) *DiagnosticSummary {
	summary := &DiagnosticSummary{
		TotalChecks: len(checks),
		ByCategory:  make(map[string]int),
		BySeverity:  make(map[string]int),
	}
	
	for _, check := range checks {
		summary.ByCategory[check.Category]++
		summary.BySeverity[check.Severity]++
		
		switch check.Status {
		case "pass":
			summary.PassedChecks++
		case "warning":
			summary.WarningChecks++
		case "error":
			summary.ErrorChecks++
		}
		
		if check.Severity == "critical" {
			summary.CriticalIssues++
		}
	}
	
	// Calculate category-specific scores
	summary.MemoryScore = d.calculateCategoryScore(checks, "memory")
	summary.PerformanceScore = d.calculateCategoryScore(checks, "performance") 
	summary.CompatibilityScore = d.calculateCategoryScore(checks, "compatibility")
	
	return summary
}

func (d *Doctor) calculateCategoryScore(checks []*DiagnosticCheck, category string) float32 {
	categoryChecks := 0
	totalScore := float32(0)
	
	for _, check := range checks {
		if check.Category == category {
			categoryChecks++
			switch check.Status {
			case "pass":
				totalScore += 100
			case "warning":
				totalScore += 50
			case "error":
				totalScore += 0
			case "info":
				totalScore += 80
			}
		}
	}
	
	if categoryChecks == 0 {
		return 100 // No issues found
	}
	
	return totalScore / float32(categoryChecks)
}

func (d *Doctor) calculateOverallScore(summary *DiagnosticSummary) float32 {
	if summary.TotalChecks == 0 {
		return 0
	}
	
	// Weight different factors
	passWeight := float32(100)
	warningWeight := float32(50)
	errorWeight := float32(0)
	
	totalScore := float32(summary.PassedChecks)*passWeight + 
		float32(summary.WarningChecks)*warningWeight + 
		float32(summary.ErrorChecks)*errorWeight
	
	// Apply critical issue penalty
	if summary.CriticalIssues > 0 {
		totalScore *= 0.5 // Reduce score by 50% for critical issues
	}
	
	return totalScore / float32(summary.TotalChecks)
}

func (d *Doctor) determineHealthStatus(score float32) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 75:
		return "good"
	case score >= 50:
		return "fair"
	case score >= 25:
		return "poor"
	default:
		return "critical"
	}
}

func (d *Doctor) generateRecommendations(checks []*DiagnosticCheck, snapshot *model.Snapshot) []*Recommendation {
	var recommendations []*Recommendation
	
	// Generate recommendations based on failed checks
	for _, check := range checks {
		if check.Status == "warning" || check.Status == "error" {
			rec := d.createRecommendationFromCheck(check)
			if rec != nil {
				recommendations = append(recommendations, rec)
			}
		}
	}
	
	return recommendations
}

func (d *Doctor) createRecommendationFromCheck(check *DiagnosticCheck) *Recommendation {
	if check.Suggestion == "" {
		return nil // No specific recommendation for this check
	}
	
	priority := "medium"
	if check.Severity == "critical" || check.Severity == "high" {
		priority = "high"
	} else if check.Severity == "low" {
		priority = "low"
	}
	
	return &Recommendation{
		ID:          fmt.Sprintf("rec_%s", strings.ReplaceAll(strings.ToLower(check.Name), " ", "_")),
		Title:       fmt.Sprintf("Address %s", check.Name),
		Description: check.Suggestion,
		Category:    check.Category,
		Priority:    priority,
		Confidence:  0.8, // Default confidence
		Actions:     []string{check.Suggestion},
		Timestamp:   time.Now(),
	}
}