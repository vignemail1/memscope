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
	
	// Add comprehensive voltage monitoring check
	voltageMonitoringCheck := d.checkVoltageMonitoring(snapshot, timestamp)
	if voltageMonitoringCheck != nil {
		checks = append(checks, voltageMonitoringCheck)
	}
	
	// Add temperature monitoring check
	temperatureCheck := d.checkTemperatureMonitoring(snapshot, timestamp)
	if temperatureCheck != nil {
		checks = append(checks, temperatureCheck)
	}
	
	// Add SPD profile validation
	spdCheck := d.validateSPDProfiles(memoryDevices, snapshot, timestamp)
	if spdCheck != nil {
		checks = append(checks, spdCheck)
	}
	
	// Add mixed memory configuration check
	mixedCheck := d.checkMixedMemoryConfiguration(memoryDevices, snapshot, timestamp)
	if mixedCheck != nil {
		checks = append(checks, mixedCheck)
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
	
	// Add memory timing analysis
	timingCheck := d.analyzeMemoryTimings(memoryDevices, snapshot, timestamp)
	if timingCheck != nil {
		checks = append(checks, timingCheck)
	}
	
	// Add bandwidth and latency analysis
	bandwidthCheck := d.analyzeBandwidthAndLatency(memoryDevices, snapshot, timestamp)
	if bandwidthCheck != nil {
		checks = append(checks, bandwidthCheck)
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

// New diagnostic methods

func (d *Doctor) checkVoltageMonitoring(snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	// Check if voltage monitoring is available through runtime observations
	voltageObservations := 0
	if snapshot.Observations != nil {
		for _, obs := range snapshot.Observations {
			if obs.Scope == "spd" && obs.Parameter == "voltage_level" {
				voltageObservations++
			}
		}
	}
	
	status := "warning"
	severity := "medium"
	message := "No voltage monitoring data available"
	
	if voltageObservations > 0 {
		status = "pass"
		severity = "low"
		message = fmt.Sprintf("Voltage monitoring active (%d sensors)", voltageObservations)
	}
	
	return &DiagnosticCheck{
		Name:      "Voltage Monitoring",
		Category:  "system",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Enable voltage monitoring for better system health tracking",
		Timestamp: timestamp,
	}
}

func (d *Doctor) checkTemperatureMonitoring(snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	temperatureObservations := 0
	if snapshot.Observations != nil {
		for _, obs := range snapshot.Observations {
			if obs.Scope == "thermal" && strings.Contains(obs.Parameter, "temperature") {
				temperatureObservations++
			}
		}
	}
	
	status := "warning"
	severity := "medium" 
	message := "No temperature monitoring data available"
	
	if temperatureObservations > 0 {
		status = "pass"
		severity = "low"
		message = fmt.Sprintf("Temperature monitoring active (%d sensors)", temperatureObservations)
	}
	
	return &DiagnosticCheck{
		Name:      "Temperature Monitoring",
		Category:  "system",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Enable temperature monitoring for thermal management",
		Timestamp: timestamp,
	}
}

func (d *Doctor) validateSPDProfiles(devices []model.Device, snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	validProfiles := 0
	totalDevices := len(devices)
	
	for _, device := range devices {
		// Check for SPD profile data in observations
		hasValidSPD := false
		for _, obs := range snapshot.Observations {
			if obs.DeviceID == device.ID && obs.Scope == "spd" {
				hasValidSPD = true
				break
			}
		}
		if hasValidSPD {
			validProfiles++
		}
	}
	
	status := "warning"
	severity := "medium"
	message := fmt.Sprintf("SPD profiles available for %d/%d memory modules", validProfiles, totalDevices)
	
	if validProfiles == totalDevices {
		status = "pass"
		severity = "low"
		message = "All memory modules have valid SPD profile data"
	} else if validProfiles == 0 {
		status = "error"
		severity = "high"
		message = "No SPD profile data available for any memory modules"
	}
	
	return &DiagnosticCheck{
		Name:      "SPD Profile Validation",
		Category:  "memory",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Verify memory module SPD EEPROM integrity and accessibility",
		Timestamp: timestamp,
	}
}

func (d *Doctor) analyzeMemoryTimings(devices []model.Device, snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	timingIssues := []string{}
	optimizationOpportunities := []string{}
	
	for i, device := range devices {
		// Look for timing data in observations
		cl := uint16(0)
		trcd := uint16(0)
		
		for _, obs := range snapshot.Observations {
			if obs.DeviceID == device.ID && obs.Scope == "spd" {
				if obs.Parameter == "cl_timing" && obs.Value != nil && obs.Value.Unsigned != nil {
					cl = uint16(*obs.Value.Unsigned)
				}
				if obs.Parameter == "trcd_timing" && obs.Value != nil && obs.Value.Unsigned != nil {
					trcd = uint16(*obs.Value.Unsigned)
				}
			}
		}
		
		// Analyze key timing parameters if available
		if cl > 0 {
			if cl > 20 {
				timingIssues = append(timingIssues, fmt.Sprintf("Module %d has loose CL timing (%d)", i+1, cl))
			} else if cl <= 14 {
				optimizationOpportunities = append(optimizationOpportunities, fmt.Sprintf("Module %d has tight CL timing (%d) - good for performance", i+1, cl))
			}
		}
		
		// Check TRCD timing
		if trcd > 22 {
			timingIssues = append(timingIssues, fmt.Sprintf("Module %d has loose TRCD timing (%d)", i+1, trcd))
		}
	}
	
	status := "pass"
	severity := "low"
	message := "Memory timing analysis complete"
	
	if len(timingIssues) > 0 {
		status = "warning"
		severity = "medium"
		message = fmt.Sprintf("Timing issues detected: %s", strings.Join(timingIssues, "; "))
	} else if len(optimizationOpportunities) > 0 {
		status = "pass"
		severity = "low"
		message = fmt.Sprintf("Good timing configuration: %s", strings.Join(optimizationOpportunities, "; "))
	}
	
	return &DiagnosticCheck{
		Name:      "Memory Timing Analysis",
		Category:  "performance",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Consider XMP/EXPO profiles for optimized timings",
		Timestamp: timestamp,
	}
}

func (d *Doctor) checkMixedMemoryConfiguration(devices []model.Device, snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	speeds := make(map[uint64][]int)
	manufacturers := make(map[string][]int)
	voltages := make(map[float64][]int)
	partNumbers := make(map[string][]int)
	
	for i, device := range devices {
		for _, obs := range snapshot.Observations {
			if obs.DeviceID == device.ID && obs.Value != nil {
				switch obs.Parameter {
				case "speed":
					if obs.Value.Unsigned != nil {
						speeds[*obs.Value.Unsigned] = append(speeds[*obs.Value.Unsigned], i+1)
					}
				case "manufacturer":
					if obs.Value.Text != nil {
						manufacturers[*obs.Value.Text] = append(manufacturers[*obs.Value.Text], i+1)
					}
				case "voltage_level":
					if obs.Value.Decimal != nil {
						voltages[*obs.Value.Decimal] = append(voltages[*obs.Value.Decimal], i+1)
					}
				case "part_number":
					if obs.Value.Text != nil {
						partNumbers[*obs.Value.Text] = append(partNumbers[*obs.Value.Text], i+1)
					}
				}
			}
		}
	}
	
	issues := []string{}
	warnings := []string{}
	
	if len(speeds) > 1 {
		issues = append(issues, "Mixed memory speeds may cause system to run at lowest speed")
		for speed, slots := range speeds {
			warnings = append(warnings, fmt.Sprintf("%d MT/s in slots %v", speed, slots))
		}
	}
	
	if len(manufacturers) > 1 {
		issues = append(issues, "Mixed manufacturers may cause compatibility issues")
	}
	
	if len(voltages) > 1 {
		issues = append(issues, "Mixed voltages may cause stability problems")
	}
	
	if len(partNumbers) == 1 && len(partNumbers) > 0 {
		// All modules are identical - best case
		return &DiagnosticCheck{
			Name:      "Mixed Memory Configuration",
			Category:  "memory",
			Status:    "pass",
			Severity:  "low",
			Message:   "All memory modules are identical (matched kit)",
			Timestamp: timestamp,
		}
	}
	
	status := "warning"
	severity := "medium"
	if len(issues) > 2 {
		severity = "high"
	}
	
	message := fmt.Sprintf("Mixed memory configuration detected: %s", strings.Join(issues, "; "))
	if len(warnings) > 0 {
		message += fmt.Sprintf(" Details: %s", strings.Join(warnings, "; "))
	}
	
	return &DiagnosticCheck{
		Name:      "Mixed Memory Configuration",
		Category:  "memory", 
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Use matched memory kits from same manufacturer with identical specifications",
		Timestamp: timestamp,
	}
}

func (d *Doctor) analyzeBandwidthAndLatency(devices []model.Device, snapshot *model.Snapshot, timestamp time.Time) *DiagnosticCheck {
	if len(devices) == 0 {
		return nil
	}
	
	// Calculate theoretical bandwidth and latency
	totalBandwidth := float64(0)
	avgLatency := float64(0)
	validDevices := 0
	
	for _, device := range devices {
		speed := uint64(0)
		cl := uint16(15) // Default CL
		
		for _, obs := range snapshot.Observations {
			if obs.DeviceID == device.ID && obs.Value != nil {
				if obs.Parameter == "speed" && obs.Value.Unsigned != nil {
					speed = *obs.Value.Unsigned
				}
				if obs.Parameter == "cl_timing" && obs.Value.Unsigned != nil {
					cl = uint16(*obs.Value.Unsigned)
				}
			}
		}
		
		if speed > 0 {
			// DDR memory: theoretical bandwidth = speed * bus_width * 2 / 8 (bytes/sec)
			// Assuming 64-bit bus width
			bandwidth := float64(speed) * 64 * 2 / 8 / 1000 // GB/s
			totalBandwidth += bandwidth
			
			// Latency in nanoseconds: (CL / (speed * 1000000)) * 1000000000
			deviceLatency := (float64(cl) / (float64(speed) * 1000000)) * 1000000000
			avgLatency += deviceLatency
			validDevices++
		}
	}
	
	if validDevices == 0 {
		return nil
	}
	
	avgLatency = avgLatency / float64(validDevices)
	
	status := "pass"
	severity := "low"
	message := fmt.Sprintf("Theoretical bandwidth: %.1f GB/s, Average latency: %.1f ns", totalBandwidth, avgLatency)
	
	if totalBandwidth < 25.0 {
		status = "warning"
		severity = "medium"
		message = fmt.Sprintf("Low memory bandwidth (%.1f GB/s) may impact performance", totalBandwidth)
	} else if avgLatency > 15.0 {
		status = "warning"
		severity = "medium"
		message = fmt.Sprintf("High memory latency (%.1f ns) may affect responsiveness", avgLatency)
	}
	
	return &DiagnosticCheck{
		Name:      "Memory Bandwidth and Latency",
		Category:  "performance",
		Status:    status,
		Severity:  severity,
		Message:   message,
		Suggestion: "Consider faster memory or tighter timings for better performance",
		Timestamp: timestamp,
	}
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