package doctor

import (
	"time"
)

// DiagnosticResult represents the complete system diagnostic analysis
type DiagnosticResult struct {
	Timestamp     time.Time          `json:"timestamp"`
	OverallScore  float32            `json:"overall_score"`  // 0-100
	HealthStatus  string             `json:"health_status"`  // excellent, good, fair, poor, critical
	Checks        []*DiagnosticCheck `json:"checks"`
	Summary       *DiagnosticSummary `json:"summary"`
	Recommendations []*Recommendation `json:"recommendations,omitempty"`
	SystemInfo    *SystemInfo        `json:"system_info"`
}

// DiagnosticCheck represents a single system check
type DiagnosticCheck struct {
	Name        string    `json:"name"`
	Category    string    `json:"category"`    // system, memory, performance, compatibility
	Status      string    `json:"status"`      // pass, warning, error, info
	Severity    string    `json:"severity"`    // low, medium, high, critical
	Message     string    `json:"message"`
	Details     string    `json:"details,omitempty"`
	Evidence    []string  `json:"evidence,omitempty"`
	Suggestion  string    `json:"suggestion,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// DiagnosticSummary provides aggregated diagnostic statistics
type DiagnosticSummary struct {
	TotalChecks    int            `json:"total_checks"`
	PassedChecks   int            `json:"passed_checks"`
	WarningChecks  int            `json:"warning_checks"`
	ErrorChecks    int            `json:"error_checks"`
	ByCategory     map[string]int `json:"by_category"`
	BySeverity     map[string]int `json:"by_severity"`
	CriticalIssues int            `json:"critical_issues"`
	MemoryScore    float32        `json:"memory_score"`
	PerformanceScore float32      `json:"performance_score"`
	CompatibilityScore float32    `json:"compatibility_score"`
}

// Recommendation represents a suggested action for improvement
type Recommendation struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`    // high, medium, low
	Confidence  float32   `json:"confidence"`  // 0.0-1.0
	Actions     []string  `json:"actions"`
	References  []string  `json:"references,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// SystemInfo represents system information for diagnostic context
type SystemInfo struct {
	CPU           string    `json:"cpu"`
	Motherboard   string    `json:"motherboard"`
	BIOSVersion   string    `json:"bios_version"`
	MemorySlots   int       `json:"memory_slots"`
	PopulatedSlots int      `json:"populated_slots"`
	TotalCapacity string    `json:"total_capacity"`
	MemoryType    string    `json:"memory_type"`
	Platform      string    `json:"platform"`
}

// DiagnosticOptions configures diagnostic behavior
type DiagnosticOptions struct {
	IncludePerformanceAnalysis bool     `json:"include_performance_analysis"`
	IncludeCompatibilityCheck  bool     `json:"include_compatibility_check"`
	SkipSlowChecks            bool     `json:"skip_slow_checks"`
	CategoryFilter            []string `json:"category_filter,omitempty"`
	MinSeverity               string   `json:"min_severity"`
}

// CheckResult represents the result of an individual diagnostic check
type CheckResult struct {
	Passed    bool      `json:"passed"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Evidence  []string  `json:"evidence,omitempty"`
	Score     float32   `json:"score"`     // 0-100
	Timestamp time.Time `json:"timestamp"`
}