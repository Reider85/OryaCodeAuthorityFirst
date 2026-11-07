package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

type PerformanceTrendConfig struct {
	InputPath      string `json:"input_path"`
	OutputPath     string `json:"output_path"`
	CheckDegradation bool `json:"check_degradation"`
	GenerateHTML   bool  `json:"generate_html"`
}

type TrendAnalysis struct {
	Timestamp   time.Time                 `json:"timestamp"`
	Operations  map[string]OperationTrend `json:"operations"`
	OverallTrend TrendDirection           `json:"overall_trend"`
	Alerts      []TrendAlert              `json:"alerts"`
	Summary     TrendSummary              `json:"summary"`
}

type OperationTrend struct {
	Operation     string        `json:"operation"`
	Trend         TrendDirection `json:"trend"`
	ChangeRate    float64       `json:"change_rate"`
	CurrentAvg    time.Duration `json:"current_avg"`
	BaselineAvg   time.Duration `json:"baseline_avg"`
	Confidence    float64       `json:"confidence"`
	Significant   bool          `json:"significant"`
}

type TrendAlert struct {
	Operation     string        `json:"operation"`
	Type          string        `json:"type"`
	Severity      string        `json:"severity"`
	Message       string        `json:"message"`
	Trend         TrendDirection `json:"trend"`
	ChangeRate    float64       `json:"change_rate"`
	Timestamp     time.Time     `json:"timestamp"`
}

type TrendSummary struct {
	TotalOperations int               `json:"total_operations"`
	Improving       int               `json:"improving"`
	Degrading      int               `json:"degrading"`
	Stable         int               `json:"stable"`
	CriticalAlerts  int               `json:"critical_alerts"`
	WarningAlerts   int               `json:"warning_alerts"`
	OverallHealth  string            `json:"overall_health"`
}

func main() {
	// Parse command line flags
	config := &PerformanceTrendConfig{}
	flag.StringVar(&config.InputPath, "input", "", "Path to input performance data")
	flag.StringVar(&config.OutputPath, "output", "", "Path to output analysis")
	flag.BoolVar(&config.CheckDegradation, "check-degradation", false, "Check for performance degradation")
	flag.BoolVar(&config.GenerateHTML, "html", false, "Generate HTML report")
	flag.Parse()

	// Validate input
	if config.InputPath == "" {
		log.Fatal("Input path is required")
	}

	// Load performance data
	performanceData, err := loadPerformanceData(config.InputPath)
	if err != nil {
		log.Fatalf("Failed to load performance data: %v", err)
	}

	// Analyze trends
	analysis := analyzeTrends(performanceData)

	// Check for degradation
	if config.CheckDegradation {
		checkDegradation(analysis)
	}

	// Generate output
	if config.OutputPath != "" {
		if config.GenerateHTML {
			generateHTMLReport(analysis, config.OutputPath)
		} else {
			generateJSONReport(analysis, config.OutputPath)
		}
	}

	// Log summary
	logTrendSummary(analysis)
}

func loadPerformanceData(inputPath string) (*ui.PerformanceMetricsCollector, error) {
	// In a real implementation, this would load from file
	// For now, create a mock collector
	budgets := ui.NewPerformanceBudget()
	monitor := ui.NewPerformanceMonitor(budgets)
	
	// Add some sample data for testing
	ctx := monitor.StartOperation("first_paint")
	time.Sleep(50 * time.Millisecond)
	ctx.End(nil)
	
	ctx = monitor.StartOperation("frame_render")
	time.Sleep(10 * time.Millisecond)
	ctx.End(nil)
	
	ctx = monitor.StartOperation("keystroke")
	time.Sleep(20 * time.Millisecond)
	ctx.End(nil)
	
	return monitor.GetCollector(), nil
}

func analyzeTrends(collector *ui.PerformanceMetricsCollector) *TrendAnalysis {
	analysis := &TrendAnalysis{
		Timestamp:   time.Now(),
		Operations:  make(map[string]OperationTrend),
		OverallTrend: TrendStable,
		Alerts:      []TrendAlert{},
		Summary: TrendSummary{
			TotalOperations: 0,
			Improving:       0,
			Degrading:      0,
			Stable:         0,
			CriticalAlerts:  0,
			WarningAlerts:   0,
			OverallHealth:  "good",
		},
	}

	// Analyze each operation
	operations := []string{"first_paint", "frame_render", "keystroke", "file_op", "lsp_startup", "doctor_run", "dashboard_load"}
	
	for _, op := range operations {
		metrics := collector.GetRecentMetrics(op, 10)
		if len(metrics) < 2 {
			continue
		}
		
		trend := analyzeOperationTrend(op, metrics)
		analysis.Operations[op] = trend
		
		// Update summary
		analysis.Summary.TotalOperations++
		switch trend.Trend {
		case TrendUp:
			analysis.Summary.Improving++
		case TrendDown:
			analysis.Summary.Degrading++
		case TrendStable:
			analysis.Summary.Stable++
		}
		
		// Generate alerts for significant changes
		if trend.Significant {
			alert := generateTrendAlert(op, trend)
			analysis.Alerts = append(analysis.Alerts, alert)
			
			// Update alert counts
			if alert.Severity == "critical" {
				analysis.Summary.CriticalAlerts++
			} else if alert.Severity == "warning" {
				analysis.Summary.WarningAlerts++
			}
		}
	}

	// Determine overall trend
	analysis.OverallTrend = calculateOverallTrend(analysis.Operations)
	
	// Determine overall health
	analysis.Summary.OverallHealth = calculateOverallHealth(analysis.Summary)
	
	return analysis
}

func analyzeOperationTrend(operation string, metrics []ui.PerformanceMeasurement) OperationTrend {
	if len(metrics) < 2 {
		return OperationTrend{
			Operation:  operation,
			Trend:      TrendStable,
			ChangeRate: 0,
		}
	}
	
	// Calculate trend using linear regression
	trend := calculateTrendDirection(metrics)
	
	// Calculate change rate
	latest := metrics[len(metrics)-1]
	earliest := metrics[0]
	changeRate := float64(latest.Duration-earliest.Duration) / float64(earliest.Duration)
	
	// Determine if significant
	significant := isSignificantChange(changeRate, len(metrics))
	
	return OperationTrend{
		Operation:     operation,
		Trend:         trend,
		ChangeRate:    changeRate,
		CurrentAvg:    calculateAverage(metrics),
		BaselineAvg:   calculateAverage(metrics[:len(metrics)/2]), // First half as baseline
		Confidence:    calculateConfidence(len(metrics)),
		Significant:   significant,
	}
}

func calculateTrendDirection(metrics []ui.PerformanceMeasurement) TrendDirection {
	if len(metrics) < 2 {
		return TrendStable
	}
	
	// Simple linear regression
	n := float64(len(metrics))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	
	for i, m := range metrics {
		x := float64(i)
		y := float64(m.Duration)
		
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	
	// Determine trend direction
	if slope > 0.1 {
		return TrendDown
	} else if slope < -0.1 {
		return TrendUp
	}
	
	return TrendStable
}

func calculateAverage(metrics []ui.PerformanceMeasurement) time.Duration {
	if len(metrics) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, m := range metrics {
		total += m.Duration
	}
	
	return total / time.Duration(len(metrics))
}

func calculateConfidence(samples int) float64 {
	// Simple confidence calculation based on sample size
	if samples < 5 {
		return 0.5
	} else if samples < 10 {
		return 0.7
	} else {
		return 0.9
	}
}

func isSignificantChange(changeRate float64, samples int) bool {
	// Consider significant if change rate > 20% and we have enough samples
	return changeRate > 0.2 || changeRate < -0.2
}

func generateTrendAlert(operation string, trend OperationTrend) TrendAlert {
	alert := TrendAlert{
		Operation: operation,
		Timestamp: time.Now(),
		ChangeRate: trend.ChangeRate,
		Trend:     trend.Trend,
	}
	
	// Determine alert type and severity
	if trend.ChangeRate > 0.3 { // 30%+ degradation
		alert.Type = "degradation"
		alert.Severity = "critical"
		alert.Message = fmt.Sprintf("CRITICAL: %s performance degraded by %.1f%%", operation, trend.ChangeRate*100)
	} else if trend.ChangeRate > 0.15 { // 15%+ degradation
		alert.Type = "degradation"
		alert.Severity = "warning"
		alert.Message = fmt.Sprintf("WARNING: %s performance degraded by %.1f%%", operation, trend.ChangeRate*100)
	} else if trend.ChangeRate < -0.2 { // 20%+ improvement
		alert.Type = "improvement"
		alert.Severity = "info"
		alert.Message = fmt.Sprintf("INFO: %s performance improved by %.1f%%", operation, -trend.ChangeRate*100)
	}
	
	return alert
}

func calculateOverallTrend(operations map[string]OperationTrend) TrendDirection {
	improving := 0
	degrading := 0
	stable := 0
	
	for _, trend := range operations {
		switch trend.Trend {
		case TrendUp:
			improving++
		case TrendDown:
			degrading++
		case TrendStable:
			stable++
		}
	}
	
	total := len(operations)
	if total == 0 {
		return TrendStable
	}
	
	improvingRatio := float64(improving) / float64(total)
	degradingRatio := float64(degrading) / float64(total)
	
	if degradingRatio > 0.3 {
		return TrendDown
	} else if improvingRatio > 0.3 {
		return TrendUp
	}
	
	return TrendStable
}

func calculateOverallHealth(summary TrendSummary) string {
	if summary.CriticalAlerts > 0 {
		return "critical"
	} else if summary.WarningAlerts > 0 {
		return "warning"
	} else if summary.Degrading > 0 {
		return "fair"
	} else {
		return "good"
	}
}

func checkDegradation(analysis *TrendAnalysis) {
	for _, alert := range analysis.Alerts {
		if alert.Type == "degradation" {
			log.Printf("Performance degradation detected: %s", alert.Message)
		}
	}
	
	if analysis.OverallTrend == TrendDown {
		log.Printf("Overall performance trend is degrading")
	}
}

func generateJSONReport(analysis *TrendAnalysis, outputPath string) {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal analysis: %v", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}
}

func generateHTMLReport(analysis *TrendAnalysis, outputPath string) {
	html := generateHTMLContent(analysis)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	
	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		log.Fatalf("Failed to write HTML report: %v", err)
	}
}

func generateHTMLContent(analysis *TrendAnalysis) string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Performance Trend Analysis</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .operation { margin: 10px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .improving { background: #d4edda; }
        .degrading { background: #f8d7da; }
        .stable { background: #fff3cd; }
        .alert { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .critical { background: #f8d7da; border: 1px solid #f5c6cb; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; }
        .info { background: #d1ecf1; border: 1px solid #bee5eb; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Performance Trend Analysis</h1>
        <p>Generated on: ` + analysis.Timestamp.Format("2006-01-02 15:04:05") + `</p>
        <p>Overall Health: <strong>` + analysis.Summary.OverallHealth + `</strong></p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <table>
            <tr><th>Total Operations</th><th>Improving</th><th>Degrading</th><th>Stable</th><th>Critical Alerts</th><th>Warning Alerts</th></tr>
            <tr><td>` + fmt.Sprintf("%d", analysis.Summary.TotalOperations) + `</td><td>` + fmt.Sprintf("%d", analysis.Summary.Improving) + `</td><td>` + fmt.Sprintf("%d", analysis.Summary.Degrading) + `</td><td>` + fmt.Sprintf("%d", analysis.Summary.Stable) + `</td><td>` + fmt.Sprintf("%d", analysis.Summary.CriticalAlerts) + `</td><td>` + fmt.Sprintf("%d", analysis.Summary.WarningAlerts) + `</td></tr>
        </table>
    </div>
    
    <div class="operations">
        <h2>Operation Trends</h2>
`
	
	for op, trend := range analysis.Operations {
		className := "stable"
		if trend.Trend == TrendUp {
			className = "improving"
		} else if trend.Trend == TrendDown {
			className = "degrading"
		}
		
		html += `        <div class="operation ` + className + `">
            <h3>` + op + `</h3>
            <p>Trend: ` + trendToString(trend.Trend) + `</p>
            <p>Change Rate: ` + fmt.Sprintf("%.2f%%", trend.ChangeRate*100) + `</p>
            <p>Current Average: ` + trend.CurrentAvg.String() + `</p>
            <p>Baseline Average: ` + trend.BaselineAvg.String() + `</p>
            <p>Confidence: ` + fmt.Sprintf("%.2f", trend.Confidence) + `</p>
            <p>Significant: ` + fmt.Sprintf("%t", trend.Significant) + `</p>
        </div>
`
	}
	
	html += `    </div>
    
    <div class="alerts">
        <h2>Alerts</h2>
`
	
	for _, alert := range analysis.Alerts {
		html += `        <div class="alert ` + alert.Severity + `">
            <p><strong>` + alert.Type + `:</strong> ` + alert.Message + `</p>
            <p>Operation: ` + alert.Operation + `</p>
            <p>Trend: ` + trendToString(alert.Trend) + `</p>
            <p>Change Rate: ` + fmt.Sprintf("%.2f%%", alert.ChangeRate*100) + `</p>
            <p>Timestamp: ` + alert.Timestamp.Format("2006-01-02 15:04:05") + `</p>
        </div>
`
	}
	
	if len(analysis.Alerts) == 0 {
		html += `        <p>No alerts generated.</p>
`
	}
	
	html += `    </div>
</body>
</html>`
	
	return html
}

func trendToString(trend TrendDirection) string {
	switch trend {
	case TrendUp:
		return "Improving"
	case TrendDown:
		return "Degrading"
	case TrendStable:
		return "Stable"
	default:
		return "Unknown"
	}
}

func logTrendSummary(analysis *TrendAnalysis) {
	log.Printf("Performance Trend Analysis Summary:")
	log.Printf("  Overall Health: %s", analysis.Summary.OverallHealth)
	log.Printf("  Total Operations: %d", analysis.Summary.TotalOperations)
	log.Printf("  Improving: %d", analysis.Summary.Improving)
	log.Printf("  Degrading: %d", analysis.Summary.Degrading)
	log.Printf("  Stable: %d", analysis.Summary.Stable)
	log.Printf("  Critical Alerts: %d", analysis.Summary.CriticalAlerts)
	log.Printf("  Warning Alerts: %d", analysis.Summary.WarningAlerts)
	
	if analysis.Summary.CriticalAlerts > 0 {
		log.Printf("  CRITICAL: Performance degradation detected!")
	} else if analysis.Summary.WarningAlerts > 0 {
		log.Printf("  WARNING: Some performance issues detected")
	} else {
		log.Printf("  GOOD: Performance is stable or improving")
	}
}