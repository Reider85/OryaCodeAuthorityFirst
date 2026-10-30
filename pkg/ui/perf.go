package ui

import (
	"fmt"
	"time"
	
	"github.com/charmbracelet/bubbletea"
)

// PerformanceBudget defines performance budgets for the TUI
// INV-UI-001: TUI first paint <100ms on pilot repos
// INV-UI-002: Frame render <16ms (60fps)
// INV-UI-003: Keystroke response <50ms for simple actions
type PerformanceBudget struct {
	FirstPaint     time.Duration // <100ms
	FrameRender    time.Duration // <16ms
	Keystroke     time.Duration // <50ms
	FileOp        time.Duration // <1s for 10k files
	LSPStartup    time.Duration // ~500ms
	DoctorRun     time.Duration // <5s
	DashboardLoad time.Duration // <500ms
}

// NewPerformanceBudget creates a new performance budget
func NewPerformanceBudget() *PerformanceBudget {
	return &PerformanceBudget{
		FirstPaint:     100 * time.Millisecond,
		FrameRender:    16 * time.Millisecond,
		Keystroke:      50 * time.Millisecond,
		FileOp:         1 * time.Second,
		LSPStartup:     500 * time.Millisecond,
		DoctorRun:      5 * time.Second,
		DashboardLoad:  500 * time.Millisecond,
	}
}

// PerfTracker tracks performance metrics
type PerfTracker struct {
	startTime time.Time
	measures  map[string]time.Duration
	budgets   *PerformanceBudget
}

// NewPerfTracker creates a new performance tracker
func NewPerfTracker(budgets *PerformanceBudget) *PerfTracker {
	return &PerfTracker{
		measures: make(map[string]time.Duration),
		budgets:  budgets,
	}
}

// Start starts a timer for a named operation
func (pt *PerfTracker) Start(name string) {
	pt.startTime = time.Now()
}

// End stops the timer and records the duration
func (pt *PerfTracker) End(name string) time.Duration {
	duration := time.Since(pt.startTime)
	pt.measures[name] = duration
	return duration
}

// GetDuration returns the duration for a named operation
func (pt *PerfTracker) GetDuration(name string) (time.Duration, bool) {
	duration, exists := pt.measures[name]
	return duration, exists
}

// CheckBudget checks if a named operation is within budget
func (pt *PerfTracker) CheckBudget(name string) (bool, time.Duration, time.Duration) {
	duration, exists := pt.measures[name]
	if !exists {
		return false, 0, 0
	}
	
	var budget time.Duration
	switch name {
	case "first_paint":
		budget = pt.budgets.FirstPaint
	case "frame_render":
		budget = pt.budgets.FrameRender
	case "keystroke":
		budget = pt.budgets.Keystroke
	case "file_op":
		budget = pt.budgets.FileOp
	case "lsp_startup":
		budget = pt.budgets.LSPStartup
	case "doctor_run":
		budget = pt.budgets.DoctorRun
	case "dashboard_load":
		budget = pt.budgets.DashboardLoad
	default:
		return false, 0, 0
	}
	
	withinBudget := duration <= budget
	return withinBudget, duration, budget
}

// GetReport returns a performance report
func (pt *PerfTracker) GetReport() string {
	report := "Performance Report:\n"
	
	// Check first paint
	if within, duration, budget := pt.CheckBudget("first_paint"); within {
		report += fmt.Sprintf("✓ First paint: %s (within %s)\n", duration, budget)
	} else {
		report += fmt.Sprintf("✗ First paint: %s (exceeds %s)\n", duration, budget)
	}
	
	// Check frame render
	if within, duration, budget := pt.CheckBudget("frame_render"); within {
		report += fmt.Sprintf("✓ Frame render: %s (within %s)\n", duration, budget)
	} else {
		report += fmt.Sprintf("✗ Frame render: %s (exceeds %s)\n", duration, budget)
	}
	
	// Check keystroke
	if within, duration, budget := pt.CheckBudget("keystroke"); within {
		report += fmt.Sprintf("✓ Keystroke: %s (within %s)\n", duration, budget)
	} else {
		report += fmt.Sprintf("✗ Keystroke: %s (exceeds %s)\n", duration, budget)
	}
	
	// Check other operations
	operations := []string{"file_op", "lsp_startup", "doctor_run", "dashboard_load"}
	for _, op := range operations {
		if within, duration, budget := pt.CheckBudget(op); within {
			report += fmt.Sprintf("✓ %s: %s (within %s)\n", op, duration, budget)
		} else {
			report += fmt.Sprintf("✗ %s: %s (exceeds %s)\n", op, duration, budget)
		}
	}
	
	return report
}

// PerfMsg represents a performance measurement message
type PerfMsg struct {
	Name      string
	Duration  time.Duration
	Timestamp time.Time
}

// MeasurePerformance is a tea.Cmd that measures performance
func MeasurePerformance(name string) tea.Cmd {
	return func() tea.Msg {
		return PerfMsg{
			Name:      name,
			Duration:  time.Since(time.Now()), // This would be tracked by the caller
			Timestamp: time.Now(),
		}
	}
}

// PerformanceWarning represents a performance warning
type PerformanceWarning struct {
	Operation string
	Actual    time.Duration
	Budget    time.Duration
	Ratio     float64
}

// NewPerformanceWarning creates a new performance warning
func NewPerformanceWarning(operation string, actual, budget time.Duration) *PerformanceWarning {
	ratio := float64(actual) / float64(budget)
	return &PerformanceWarning{
		Operation: operation,
		Actual:    actual,
		Budget:    budget,
		Ratio:     ratio,
	}
}

// IsCritical returns true if the performance issue is critical
func (pw *PerformanceWarning) IsCritical() bool {
	return pw.Ratio > 2.0 // >2x budget is critical
}

// String returns a string representation of the warning
func (pw *PerformanceWarning) String() string {
	if pw.IsCritical() {
		return fmt.Sprintf("CRITICAL: %s took %s (budget %s, %.1fx)", pw.Operation, pw.Actual, pw.Budget, pw.Ratio)
	}
	return fmt.Sprintf("WARNING: %s took %s (budget %s, %.1fx)", pw.Operation, pw.Actual, pw.Budget, pw.Ratio)
}

// EnforceBudget enforces a performance budget for an operation
func EnforceBudget(name string, actual time.Duration, budget time.Duration) (*PerformanceWarning, bool) {
	if actual <= budget {
		return nil, true
	}
	
	warning := NewPerformanceWarning(name, actual, budget)
	return warning, false
}

// Benchmark runs a benchmark for an operation
func Benchmark(name string, fn func() time.Duration) (time.Duration, *PerformanceWarning, bool) {
	duration := fn()
	
	budget := NewPerformanceBudget()
	var warning *PerformanceWarning
	var withinBudget bool
	
	switch name {
	case "first_paint":
		warning, withinBudget = EnforceBudget(name, duration, budget.FirstPaint)
	case "frame_render":
		warning, withinBudget = EnforceBudget(name, duration, budget.FrameRender)
	case "keystroke":
		warning, withinBudget = EnforceBudget(name, duration, budget.Keystroke)
	default:
		warning, withinBudget = EnforceBudget(name, duration, budget.FileOp)
	}
	
	return duration, warning, withinBudget
}

// OptimizeForPerformance returns optimization suggestions
func OptimizeForPerformance(perfTracker *PerfTracker) []string {
	suggestions := []string{}
	
	// Check first paint
	if within, _, _ := perfTracker.CheckBudget("first_paint"); !within {
		suggestions = append(suggestions, "Consider using Skeleton widgets for async data loading")
		suggestions = append(suggestions, "Pre-compute styles at startup")
	}
	
	// Check frame render
	if within, _, _ := perfTracker.CheckBudget("frame_render"); !within {
		suggestions = append(suggestions, "Use virtualized lists for large datasets")
		suggestions = append(suggestions, "Cache computed layouts")
		suggestions = append(suggestions, "Reduce re-renders by using memoization")
	}
	
	// Check keystroke
	if within, _, _ := perfTracker.CheckBudget("keystroke"); !within {
		suggestions = append(suggestions, "Use direct dispatch for simple actions")
		suggestions = append(suggestions, "Avoid async round-trips for typing")
	}
	
	return suggestions
}