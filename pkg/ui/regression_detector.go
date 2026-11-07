package ui

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// RegressionDetector detects performance regressions
type RegressionDetector struct {
	monitor      *PerformanceMonitor
	budgets      *PerformanceBudget
	detectionMode DetectionMode
	alertHandler func(*RegressionAlert)
	mu           sync.RWMutex
	enabled      bool
	baseline     map[string]PerformanceBaseline
}

// DetectionMode defines how regressions are detected
type DetectionMode int

const (
	// ModeBasic - Basic regression detection
	ModeBasic DetectionMode = iota
	// ModeAdvanced - Advanced regression detection with trend analysis
	ModeAdvanced
	// ModeStrict - Strict regression detection with immediate alerts
	ModeStrict
)

// RegressionAlert represents a regression alert
type RegressionAlert struct {
	Operation     string
	Message       string
	Severity      AlertSeverity
	Timestamp     time.Time
	Current       PerformanceMeasurement
	Previous      PerformanceMeasurement
	RegressionType RegressionType
	Ratio         float64
	Trend         TrendDirection
}

// RegressionType defines the type of regression
type RegressionType int

const (
	// RegressionTypeSudden - Sudden performance degradation
	RegressionTypeSudden RegressionType = iota
	// RegressionTypeGradual - Gradual performance degradation
	RegressionTypeGradual
	// RegressionTypeSpike - Performance spike
	RegressionTypeSpike
	// RegressionTypeTrend - Consistent trend
	RegressionTypeTrend
)

// TrendDirection defines the trend direction
type TrendDirection int

const (
	// TrendUp - Performance improving
	TrendUp TrendDirection = iota
	// TrendDown - Performance degrading
	TrendDown
	// TrendStable - Performance stable
	TrendStable
)

// PerformanceBaseline represents a performance baseline
type PerformanceBaseline struct {
	Average      time.Duration
	Min          time.Duration
	Max          time.Duration
	StdDev       time.Duration
	Samples      int
	LastUpdated  time.Time
	Confidence   float64
}

// NewRegressionDetector creates a new regression detector
func NewRegressionDetector(monitor *PerformanceMonitor, budgets *PerformanceBudget) *RegressionDetector {
	return &RegressionDetector{
		monitor:      monitor,
		budgets:      budgets,
		detectionMode: ModeAdvanced,
		enabled:      true,
		baseline:     make(map[string]PerformanceBaseline),
	}
}

// SetDetectionMode sets the detection mode
func (rd *RegressionDetector) SetDetectionMode(mode DetectionMode) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	
	rd.detectionMode = mode
}

// SetAlertHandler sets the alert handler function
func (rd *RegressionDetector) SetAlertHandler(handler func(*RegressionAlert)) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	
	rd.alertHandler = handler
}

// Enable enables regression detection
func (rd *RegressionDetector) Enable() {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	
	rd.enabled = true
}

// Disable disables regression detection
func (rd *RegressionDetector) Disable() {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	
	rd.enabled = false
}

// IsEnabled returns whether regression detection is enabled
func (rd *RegressionDetector) IsEnabled() bool {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	
	return rd.enabled
}

// UpdateBaseline updates the performance baseline for an operation
func (rd *RegressionDetector) UpdateBaseline(operation string, measurements []PerformanceMeasurement) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	
	if len(measurements) == 0 {
		return
	}
	
	// Calculate statistics
	var total time.Duration
	min := measurements[0].Duration
	max := measurements[0].Duration
	
	for _, m := range measurements {
		total += m.Duration
		if m.Duration < min {
			min = m.Duration
		}
		if m.Duration > max {
			max = m.Duration
		}
	}
	
	avg := total / time.Duration(len(measurements))
	
	// Calculate standard deviation
	var sumSquares float64
	for _, m := range measurements {
		diff := float64(m.Duration - avg)
		sumSquares += diff * diff
	}
	stdDev := time.Duration(math.Sqrt(sumSquares / float64(len(measurements))))
	
	// Calculate confidence (based on sample size and variance)
	confidence := calculateConfidence(len(measurements), stdDev, avg)
	
	rd.baseline[operation] = PerformanceBaseline{
		Average:     avg,
		Min:         min,
		Max:         max,
		StdDev:      stdDev,
		Samples:     len(measurements),
		LastUpdated: time.Now(),
		Confidence:  confidence,
	}
}

// DetectRegressions detects regressions in recent measurements
func (rd *RegressionDetector) DetectRegressions() []*RegressionAlert {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	
	if !rd.enabled {
		return nil
	}
	
	alerts := []*RegressionAlert{}
	
	// Get recent measurements
	metrics := rd.monitor.GetMetrics()
	
	for operation, measurements := range metrics {
		if len(measurements) < 2 {
			continue
		}
		
		// Get baseline
		baseline, hasBaseline := rd.baseline[operation]
		if !hasBaseline {
			continue
		}
		
		// Analyze recent measurements
		recentAlerts := rd.analyzeMeasurements(operation, measurements, baseline)
		alerts = append(alerts, recentAlerts...)
	}
	
	return alerts
}

// analyzeMeasurements analyzes measurements for regressions
func (rd *RegressionDetector) analyzeMeasurements(operation string, measurements []PerformanceMeasurement, baseline PerformanceBaseline) []*RegressionAlert {
	alerts := []*RegressionAlert{}
	
	// Get the most recent measurement
	current := measurements[len(measurements)-1]
	
	// Compare with previous measurement
	if len(measurements) >= 2 {
		previous := measurements[len(measurements)-2]
		alert := rd.compareMeasurements(operation, current, previous, baseline)
		if alert != nil {
			alerts = append(alerts, alert)
		}
	}
	
	// Check for trend
	if rd.detectionMode == ModeAdvanced || rd.detectionMode == ModeStrict {
		trendAlert := rd.analyzeTrend(operation, measurements, baseline)
		if trendAlert != nil {
			alerts = append(alerts, trendAlert)
		}
	}
	
	// Check for spikes
	if rd.detectionMode == ModeStrict {
		spikeAlert := rd.detectSpike(operation, measurements, baseline)
		if spikeAlert != nil {
			alerts = append(alerts, spikeAlert)
		}
	}
	
	return alerts
}

// compareMeasurements compares two measurements for regression
func (rd *RegressionDetector) compareMeasurements(operation string, current, previous PerformanceMeasurement, baseline PerformanceBaseline) *RegressionAlert {
	// Calculate regression ratio
	ratio := float64(current.Duration) / float64(previous.Duration)
	
	// Check for significant regression
	if ratio > 1.2 { // 20% degradation
		alert := &RegressionAlert{
			Operation:     operation,
			Timestamp:     current.Timestamp,
			Current:       current,
			Previous:      previous,
			Ratio:         ratio,
			RegressionType: RegressionTypeSudden,
		}
		
		// Determine severity
		if ratio > 2.0 {
			alert.Severity = SeverityCritical
			alert.Message = fmt.Sprintf("CRITICAL: %s performance degraded by %.1fx (from %s to %s)", 
				operation, ratio, previous.Duration, current.Duration)
		} else {
			alert.Severity = SeverityError
			alert.Message = fmt.Sprintf("ERROR: %s performance degraded by %.1fx (from %s to %s)", 
				operation, ratio, previous.Duration, current.Duration)
		}
		
		// Determine trend
		alert.Trend = rd.getTrendDirection(current.Duration, baseline.Average)
		
		return alert
	}
	
	return nil
}

// analyzeTrend analyzes performance trend
func (rd *RegressionDetector) analyzeTrend(operation string, measurements []PerformanceMeasurement, baseline PerformanceBaseline) *RegressionAlert {
	if len(measurements) < 5 {
		return nil
	}
	
	// Get recent measurements (last 5)
	recent := measurements[len(measurements)-5:]
	
	// Calculate trend
	trend := calculateTrend(recent)
	
	// Check for consistent degrading trend
	if trend < -0.1 { // Consistent degradation
		alert := &RegressionAlert{
			Operation:     operation,
			Timestamp:     time.Now(),
			Current:       measurements[len(measurements)-1],
			RegressionType: RegressionTypeTrend,
			Trend:         TrendDown,
			Ratio:         float64(measurements[len(measurements)-1].Duration) / float64(baseline.Average),
		}
		
		alert.Severity = SeverityWarning
		alert.Message = fmt.Sprintf("WARNING: %s shows consistent degrading trend (current: %s, baseline: %s)", 
			operation, measurements[len(measurements)-1].Duration, baseline.Average)
		
		return alert
	}
	
	return nil
}

// detectSpike detects performance spikes
func (rd *RegressionDetector) detectSpike(operation string, measurements []PerformanceMeasurement, baseline PerformanceBaseline) *RegressionAlert {
	if len(measurements) < 3 {
		return nil
	}
	
	current := measurements[len(measurements)-1]
	
	// Check if current measurement is a spike
	if current.Duration > baseline.Max*2 { // More than 2x the historical max
		alert := &RegressionAlert{
			Operation:     operation,
			Timestamp:     current.Timestamp,
			Current:       current,
			RegressionType: RegressionTypeSpike,
			Ratio:         float64(current.Duration) / float64(baseline.Average),
		}
		
		alert.Severity = SeverityError
		alert.Message = fmt.Sprintf("ERROR: %s performance spike detected (current: %s, baseline: %s)", 
			operation, current.Duration, baseline.Average)
		
		return alert
	}
	
	return nil
}

// getTrendDirection determines the trend direction
func (rd *RegressionDetector) getTrendDirection(current, baseline time.Duration) TrendDirection {
	ratio := float64(current) / float64(baseline)
	
	if ratio > 1.1 {
		return TrendDown
	} else if ratio < 0.9 {
		return TrendUp
	}
	
	return TrendStable
}

// calculateTrend calculates the trend direction
func calculateTrend(measurements []PerformanceMeasurement) float64 {
	if len(measurements) < 2 {
		return 0
	}
	
	// Simple linear regression
	n := float64(len(measurements))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	
	for i, m := range measurements {
		x := float64(i)
		y := float64(m.Duration)
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	
	return slope
}

// calculateConfidence calculates confidence level
func calculateConfidence(samples int, stdDev, avg time.Duration) float64 {
	if samples == 0 {
		return 0
	}
	
	// Simple confidence calculation based on sample size and variance
	// More samples and lower variance = higher confidence
	confidence := float64(samples) / (float64(samples) + float64(stdDev)/float64(avg))
	
	return math.Min(1.0, math.Max(0.0, confidence))
}

// GetBaseline returns the baseline for an operation
func (rd *RegressionDetector) GetBaseline(operation string) (PerformanceBaseline, bool) {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	
	baseline, exists := rd.baseline[operation]
	return baseline, exists
}

// GetRegressionReport returns a regression detection report
func (rd *RegressionDetector) GetRegressionReport() string {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	
	report := "=== Regression Detection Report ===\n\n"
	
	// Detection mode
	report += fmt.Sprintf("Detection Mode: %s\n", rd.detectionMode)
	report += fmt.Sprintf("Enabled: %t\n", rd.enabled)
	
	// Baselines
	report += "\nPerformance Baselines:\n"
	for operation, baseline := range rd.baseline {
		report += fmt.Sprintf("  %s:\n", operation)
		report += fmt.Sprintf("    Average: %s\n", baseline.Average)
		report += fmt.Sprintf("    Min: %s\n", baseline.Min)
		report += fmt.Sprintf("    Max: %s\n", baseline.Max)
		report += fmt.Sprintf("    StdDev: %s\n", baseline.StdDev)
		report += fmt.Sprintf("    Samples: %d\n", baseline.Samples)
		report += fmt.Sprintf("    Confidence: %.2f\n", baseline.Confidence)
		report += fmt.Sprintf("    Last Updated: %s\n", baseline.LastUpdated.Format("2006-01-02 15:04:05"))
	}
	
	// Recent regressions
	alerts := rd.DetectRegressions()
	if len(alerts) > 0 {
		report += "\nRecent Regressions:\n"
		for _, alert := range alerts {
			report += fmt.Sprintf("  [%s] %s: %s\n", 
				alert.Timestamp.Format("15:04:05"), 
				alert.Operation, 
				alert.Message)
		}
	} else {
		report += "\nNo recent regressions detected.\n"
	}
	
	return report
}

// MonitorAndDetect combines monitoring and regression detection
func (rd *RegressionDetector) MonitorAndDetect(name string, fn func() interface{}) interface{} {
	startTime := time.Now()
	
	// Execute the operation
	result := fn()
	duration := time.Since(startTime)
	
	// Get recent measurements
	measurements := rd.monitor.GetRecentMetrics(name, 10)
	
	// Update baseline if we have enough samples
	if len(measurements) >= 5 {
		rd.UpdateBaseline(name, measurements)
	}
	
	// Detect regressions
	alerts := rd.DetectRegressions()
	
	// Handle alerts
	for _, alert := range alerts {
		if rd.alertHandler != nil {
			go rd.alertHandler(alert)
		}
	}
	
	return result
}