package ui

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// BudgetEnforcer provides automatic budget enforcement
type BudgetEnforcer struct {
	monitor       *PerformanceMonitor
	budgets       *PerformanceBudget
	enforcementMode EnforcementMode
	actionHandler  func(EnforcementAction)
	mu            sync.RWMutex
	enabled       bool
}

// EnforcementMode defines how budgets are enforced
type EnforcementMode int

const (
	// ModeMonitor - Only monitor, no enforcement
	ModeMonitor EnforcementMode = iota
	// ModeWarning - Warn when budgets are exceeded
	ModeWarning
	// ModeFail - Fail when budgets are exceeded
	ModeFail
	// ModeStrict - Strict enforcement with immediate action
	ModeStrict
)

// EnforcementAction represents an enforcement action
type EnforcementAction struct {
	Operation     string
	Action        string
	Reason        string
	Severity      AlertSeverity
	Timestamp     time.Time
	Measurement   PerformanceMeasurement
	Recommendations []string
}

// NewBudgetEnforcer creates a new budget enforcer
func NewBudgetEnforcer(monitor *PerformanceMonitor, budgets *PerformanceBudget) *BudgetEnforcer {
	return &BudgetEnforcer{
		monitor:        monitor,
		budgets:        budgets,
		enforcementMode: ModeWarning,
		enabled:        true,
	}
}

// SetEnforcementMode sets the enforcement mode
func (be *BudgetEnforcer) SetEnforcementMode(mode EnforcementMode) {
	be.mu.Lock()
	defer be.mu.Unlock()
	
	be.enforcementMode = mode
}

// SetActionHandler sets the action handler function
func (be *BudgetEnforcer) SetActionHandler(handler func(EnforcementAction)) {
	be.mu.Lock()
	defer be.mu.Unlock()
	
	be.actionHandler = handler
}

// Enable enables budget enforcement
func (be *BudgetEnforcer) Enable() {
	be.mu.Lock()
	defer be.mu.Unlock()
	
	be.enabled = true
}

// Disable disables budget enforcement
func (be *BudgetEnforcer) Disable() {
	be.mu.Lock()
	defer be.mu.Unlock()
	
	be.enabled = false
}

// IsEnabled returns whether budget enforcement is enabled
func (be *BudgetEnforcer) IsEnabled() bool {
	be.mu.RLock()
	defer be.mu.RUnlock()
	
	return be.enabled
}

// EnforceBudget enforces a budget for an operation
func (be *BudgetEnforcer) EnforceBudget(name string, duration time.Duration, metadata map[string]interface{}) EnforcementAction {
	if !be.IsEnabled() {
		return EnforcementAction{}
	}
	
	measurement := PerformanceMeasurement{
		Operation: name,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	
	// Check if within budget
	withinBudget, warning := be.checkBudget(name, duration)
	measurement.WithinBudget = withinBudget
	measurement.Warning = warning
	
	action := EnforcementAction{
		Operation:   name,
		Timestamp:  time.Now(),
		Measurement: measurement,
	}
	
	if !withinBudget {
		action.Severity = be.getSeverityForWarning(warning)
		action.Reason = fmt.Sprintf("Operation %s exceeded budget: %s (budget %s, %.1fx)", 
			name, duration, warning.Budget, warning.Ratio)
		
		// Get recommendations
		action.Recommendations = be.getRecommendations(name, duration)
		
		// Take action based on enforcement mode
		action.Action = be.determineAction(name, duration, warning)
		
		// Log the enforcement action
		be.logEnforcementAction(action)
		
		// Call action handler if set
		if be.actionHandler != nil {
			go be.actionHandler(action)
		}
	}
	
	return action
}

// checkBudget checks if an operation is within budget
func (be *BudgetEnforcer) checkBudget(name string, duration time.Duration) (bool, *PerformanceWarning) {
	budget := be.getBudgetForOperation(name)
	if budget == 0 {
		return true, nil
	}
	
	withinBudget := duration <= budget
	if !withinBudget {
		warning := NewPerformanceWarning(name, duration, budget)
		return false, warning
	}
	return true, nil
}

// getBudgetForOperation gets the budget for a specific operation
func (be *BudgetEnforcer) getBudgetForOperation(name string) time.Duration {
	switch name {
	case "first_paint":
		return be.budgets.FirstPaint
	case "frame_render":
		return be.budgets.FrameRender
	case "keystroke":
		return be.budgets.Keystroke
	case "file_op":
		return be.budgets.FileOp
	case "lsp_startup":
		return be.budgets.LSPStartup
	case "doctor_run":
		return be.budgets.DoctorRun
	case "dashboard_load":
		return be.budgets.DashboardLoad
	default:
		return 0
	}
}

// getSeverityForWarning returns the severity for a warning
func (be *BudgetEnforcer) getSeverityForWarning(warning *PerformanceWarning) AlertSeverity {
	if warning == nil {
		return SeverityInfo
	}
	
	if warning.IsCritical() {
		return SeverityCritical
	}
	
	if warning.Ratio > 1.5 {
		return SeverityError
	}
	
	return SeverityWarning
}

// determineAction determines the action to take based on enforcement mode
func (be *BudgetEnforcer) determineAction(name string, duration time.Duration, warning *PerformanceWarning) string {
	switch be.enforcementMode {
	case ModeMonitor:
		return "monitor"
	case ModeWarning:
		return "warn"
	case ModeFail:
		if warning.IsCritical() {
			return "fail"
		}
		return "warn"
	case ModeStrict:
		if warning.IsCritical() {
			return "fail"
		} else if warning.Ratio > 1.5 {
			return "throttle"
		}
		return "warn"
	default:
		return "monitor"
	}
}

// getRecommendations returns optimization recommendations for an operation
func (be *BudgetEnforcer) getRecommendations(name string, duration time.Duration) []string {
	recommendations := []string{}
	
	switch name {
	case "first_paint":
		recommendations = append(recommendations,
			"Use Skeleton widgets for async data loading",
			"Pre-compute styles at startup",
			"Implement lazy loading of non-critical components")
	case "frame_render":
		recommendations = append(recommendations,
			"Use virtualized lists for large datasets",
			"Cache computed layouts",
			"Reduce re-renders by using memoization",
			"Optimize border and style calculations")
	case "keystroke":
		recommendations = append(recommendations,
			"Use direct dispatch for simple actions",
			"Avoid async round-trips for typing",
			"Implement input debouncing")
	case "file_op":
		recommendations = append(recommendations,
			"Use file caching",
			"Implement batch operations",
			"Use parallel processing for large file sets")
	case "lsp_startup":
		recommendations = append(recommendations,
			"Implement LSP server pooling",
			"Use lazy initialization",
			"Cache project metadata")
	case "doctor_run":
		recommendations = append(recommendations,
			"Implement incremental analysis",
			"Use parallel processing",
			"Cache analysis results")
	case "dashboard_load":
		recommendations = append(recommendations,
			"Implement data pagination",
			"Use background data loading",
			"Cache dashboard state")
	}
	
	return recommendations
}

// logEnforcementAction logs the enforcement action
func (be *BudgetEnforcer) logEnforcementAction(action EnforcementAction) {
	log.Printf("Enforcement Action: %s - %s", action.Action, action.Reason)
	
	if len(action.Recommendations) > 0 {
		log.Printf("Recommendations:")
		for _, rec := range action.Recommendations {
			log.Printf("  - %s", rec)
		}
	}
}

// MonitorAndEnforce combines monitoring and enforcement
func (be *BudgetEnforcer) MonitorAndEnforce(name string, fn func() interface{}) interface{} {
	startTime := time.Now()
	
	// Execute the operation
	result := fn()
	duration := time.Since(startTime)
	
	// Enforce budget
	be.EnforceBudget(name, duration, map[string]interface{}{
		"result": result,
	})
	
	return result
}

// MonitorAndEnforceWithDuration combines monitoring and enforcement with duration
func (be *BudgetEnforcer) MonitorAndEnforceWithDuration(name string, fn func() interface{}) (interface{}, time.Duration, EnforcementAction) {
	startTime := time.Now()
	
	// Execute the operation
	result := fn()
	duration := time.Since(startTime)
	
	// Enforce budget
	action := be.EnforceBudget(name, duration, map[string]interface{}{
		"result": result,
	})
	
	return result, duration, action
}

// GetEnforcementStats returns enforcement statistics
func (be *BudgetEnforcer) GetEnforcementStats() EnforcementStats {
	be.mu.RLock()
	defer be.mu.RUnlock()
	
	stats := EnforcementStats{
		TotalOperations: 0,
		FailedOperations: 0,
		WarningOperations: 0,
		CriticalOperations: 0,
		AverageExcess: 0,
	}
	
	// Get metrics from monitor
	metrics := be.monitor.GetMetrics()
	
	for op, measurements := range metrics {
		stats.TotalOperations += len(measurements)
		
		for _, m := range measurements {
			if !m.WithinBudget {
				stats.FailedOperations++
				
				if m.Warning != nil {
					if m.Warning.IsCritical() {
						stats.CriticalOperations++
					} else {
						stats.WarningOperations++
					}
					
					// Calculate average excess
					excess := float64(m.Duration-m.Warning.Budget) / float64(m.Warning.Budget)
					stats.AverageExcess += excess
				}
			}
		}
	}
	
	if stats.FailedOperations > 0 {
		stats.AverageExcess /= float64(stats.FailedOperations)
	}
	
	return stats
}

// EnforcementStats represents enforcement statistics
type EnforcementStats struct {
	TotalOperations    int
	FailedOperations   int
	WarningOperations  int
	CriticalOperations int
	AverageExcess      float64
	FailureRate        float64
}

// CalculateFailureRate calculates the failure rate
func (es *EnforcementStats) CalculateFailureRate() {
	if es.TotalOperations > 0 {
		es.FailureRate = float64(es.FailedOperations) / float64(es.TotalOperations)
	}
}

// GetEnforcementReport returns a comprehensive enforcement report
func (be *BudgetEnforcer) GetEnforcementReport() string {
	be.mu.RLock()
	defer be.mu.RUnlock()
	
	report := "=== Budget Enforcement Report ===\n\n"
	
	// Enforcement mode
	report += fmt.Sprintf("Enforcement Mode: %s\n", be.enforcementMode)
	report += fmt.Sprintf("Enabled: %t\n", be.enabled)
	
	// Statistics
	stats := be.GetEnforcementStats()
	stats.CalculateFailureRate()
	
	report += "\nStatistics:\n"
	report += fmt.Sprintf("  Total Operations: %d\n", stats.TotalOperations)
	report += fmt.Sprintf("  Failed Operations: %d\n", stats.FailedOperations)
	report += fmt.Sprintf("  Warning Operations: %d\n", stats.WarningOperations)
	report += fmt.Sprintf("  Critical Operations: %d\n", stats.CriticalOperations)
	report += fmt.Sprintf("  Failure Rate: %.2f%%\n", stats.FailureRate*100)
	
	if stats.FailedOperations > 0 {
		report += fmt.Sprintf("  Average Excess: %.2fx\n", stats.AverageExcess)
	}
	
	// Recent enforcement actions
	report += "\nRecent Enforcement Actions:\n"
	metrics := be.monitor.GetMetrics()
	
	for op, measurements := range metrics {
		if len(measurements) == 0 {
			continue
		}
		
		lastMeasurement := measurements[len(measurements)-1]
		if !lastMeasurement.WithinBudget && lastMeasurement.Warning != nil {
			report += fmt.Sprintf("  %s: %s\n", op, lastMeasurement.Warning.String())
		}
	}
	
	return report
}