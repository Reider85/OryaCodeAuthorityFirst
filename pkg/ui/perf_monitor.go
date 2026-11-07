package ui

import (
	"sync"
	"time"
)

// PerformanceMonitor provides real-time performance monitoring
type PerformanceMonitor struct {
	collector    *PerformanceMetricsCollector
	budgets      *PerformanceBudget
	enabled      bool
	mu           sync.RWMutex
	alertHandler func(*PerformanceAlert)
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(budgets *PerformanceBudget) *PerformanceMonitor {
	return &PerformanceMonitor{
		collector: NewPerformanceMetricsCollector(budgets, 100), // Keep last 100 samples
		budgets:   budgets,
		enabled:   true,
	}
}

// SetAlertHandler sets the alert handler function
func (pm *PerformanceMonitor) SetAlertHandler(handler func(*PerformanceAlert)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.alertHandler = handler
}

// StartOperation starts monitoring an operation
func (pm *PerformanceMonitor) StartOperation(name string) OperationContext {
	if !pm.IsEnabled() {
		return OperationContext{}
	}
	
	return OperationContext{
		name:     name,
		startTime: time.Now(),
		monitor:  pm,
	}
}

// OperationContext represents an operation being monitored
type OperationContext struct {
	name      string
	startTime time.Time
	monitor   *PerformanceMonitor
}

// End ends the operation monitoring
func (oc OperationContext) End(metadata map[string]interface{}) {
	if oc.monitor == nil {
		return
	}
	
	oc.monitor.EndOperation(oc.name, oc.startTime, metadata)
}

// IsEnabled returns whether performance monitoring is enabled
func (pm *PerformanceMonitor) IsEnabled() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	return pm.enabled
}

// Enable enables performance monitoring
func (pm *PerformanceMonitor) Enable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.enabled = true
}

// Disable disables performance monitoring
func (pm *PerformanceMonitor) Disable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.enabled = false
}

// GetCollector returns the metrics collector
func (pm *PerformanceMonitor) GetCollector() *PerformanceMetricsCollector {
	return pm.collector
}

// GetPerformanceReport returns a performance report
func (pm *PerformanceMonitor) GetPerformanceReport() string {
	return pm.collector.GetPerformanceReport()
}

// GetMetrics returns all collected metrics
func (pm *PerformanceMonitor) GetMetrics() map[string][]PerformanceMeasurement {
	return pm.collector.GetMetrics()
}

// GetAverageDuration returns the average duration for an operation
func (pm *PerformanceMonitor) GetAverageDuration(operation string) (time.Duration, bool) {
	return pm.collector.GetAverageDuration(operation)
}

// GetBudgetCompliance returns the budget compliance rate for an operation
func (pm *PerformanceMonitor) GetBudgetCompliance(operation string) (float64, bool) {
	return pm.collector.GetBudgetCompliance(operation)
}

// ClearMetrics clears all collected metrics
func (pm *PerformanceMonitor) ClearMetrics() {
	pm.collector.ClearMetrics()
}

// EndOperation ends an operation and records the measurement
func (pm *PerformanceMonitor) EndOperation(name string, startTime time.Time, metadata map[string]interface{}) {
	if !pm.IsEnabled() {
		return
	}
	
	measurement := pm.collector.EndOperation(name, startTime, metadata)
	
	// Handle alerts if an alert handler is set
	if pm.alertHandler != nil && measurement.Warning != nil {
		alert := pm.collector.generateAlert(*measurement, metadata)
		if alert != nil {
			pm.alertHandler(alert)
		}
	}
}

// PerformanceMiddleware provides a middleware for performance monitoring
type PerformanceMiddleware struct {
	monitor *PerformanceMonitor
}

// NewPerformanceMiddleware creates a new performance middleware
func NewPerformanceMiddleware(monitor *PerformanceMonitor) *PerformanceMiddleware {
	return &PerformanceMiddleware{
		monitor: monitor,
	}
}

// MeasureOperation measures an operation and returns the result
func (pm *PerformanceMiddleware) MeasureOperation(name string, fn func() interface{}) interface{} {
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	ctx.End(nil)
	return result
}

// MeasureOperationWithDuration measures an operation and returns both result and duration
func (pm *PerformanceMiddleware) MeasureOperationWithDuration(name string, fn func() interface{}) (interface{}, time.Duration) {
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	duration := time.Since(ctx.startTime)
	ctx.End(nil)
	return result, duration
}

// MeasureAsyncOperation measures an async operation
func (pm *PerformanceMiddleware) MeasureAsyncOperation(name string, fn func()) {
	go func() {
		ctx := pm.monitor.StartOperation(name)
		fn()
		ctx.End(nil)
	}()
}

// TeaCmd creates a tea.Cmd for measuring performance
func (pm *PerformanceMiddleware) TeaCmd(name string, fn func() tea.Msg) tea.Cmd {
	return func() tea.Msg {
		ctx := pm.monitor.StartOperation(name)
		msg := fn()
		ctx.End(nil)
		return msg
	}
}

// PerformanceConfig holds configuration for performance monitoring
type PerformanceConfig struct {
	Enabled       bool
	MaxSamples    int
	AlertSeverity AlertSeverity
	EnableAlerts bool
}

// DefaultPerformanceConfig returns the default performance configuration
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		Enabled:       true,
		MaxSamples:    100,
		AlertSeverity: SeverityWarning,
		EnableAlerts:  true,
	}
}

// GlobalPerformanceMonitor provides a global performance monitor instance
var GlobalPerformanceMonitor *PerformanceMonitor

// InitializeGlobalMonitor initializes the global performance monitor
func InitializeGlobalMonitor(budgets *PerformanceBudget) {
	GlobalPerformanceMonitor = NewPerformanceMonitor(budgets)
}

// GetGlobalMonitor returns the global performance monitor
func GetGlobalMonitor() *PerformanceMonitor {
	return GlobalPerformanceMonitor
}

// MonitorOperation is a convenience function to monitor an operation using the global monitor
func MonitorOperation(name string) OperationContext {
	if GlobalPerformanceMonitor == nil {
		return OperationContext{}
	}
	return GlobalPerformanceMonitor.StartOperation(name)
}

// EndMonitoring is a convenience function to end monitoring using the global monitor
func EndMonitoring(name string, metadata map[string]interface{}) {
	if GlobalPerformanceMonitor == nil {
		return
	}
	
	ctx := GlobalPerformanceMonitor.StartOperation(name)
	ctx.End(metadata)
}