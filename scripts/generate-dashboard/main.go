package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

type DashboardConfig struct {
	InputPath    string `json:"input_path"`
	OutputPath   string `json:"output_path"`
	Title        string `json:"title"`
	AutoRefresh  bool   `json:"auto_refresh"`
	Theme        string `json:"theme"`
}

type PerformanceDashboard struct {
	Title        string                 `json:"title"`
	GeneratedAt  time.Time              `json:"generated_at"`
	OverallScore float64                `json:"overall_score"`
	Status       string                 `json:"status"`
	Operations   []OperationDashboard   `json:"operations"`
	Alerts       []DashboardAlert       `json:"alerts"`
	Trends       []TrendDashboard        `json:"trends"`
	Enforcement  EnforcementDashboard   `json:"enforcement"`
	Optimizations OptimizationDashboard `json:"optimizations"`
}

type OperationDashboard struct {
	Name         string        `json:"name"`
	Status       string        `json:"status"`
	Score        float64       `json:"score"`
	AvgDuration  time.Duration `json:"avg_duration"`
	Budget       time.Duration `json:"budget"`
	Compliance   float64       `json:"compliance"`
	Trend        string        `json:"trend"`
	LastUpdated  time.Time     `json:"last_updated"`
	Impact       string        `json:"impact"`
}

type DashboardAlert struct {
	Type        string    `json:"type"`
	Operation   string    `json:"operation"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Resolved    bool      `json:"resolved"`
}

type TrendDashboard struct {
	Name        string        `json:"name"`
	Direction   string        `json:"direction"`
	ChangeRate  float64       `json:"change_rate"`
	CurrentAvg  time.Duration `json:"current_avg"`
	BaselineAvg time.Duration `json:"baseline_avg"`
	Confidence  float64       `json:"confidence"`
}

type EnforcementDashboard struct {
	TotalOperations    int     `json:"total_operations"`
	FailedOperations   int     `json:"failed_operations"`
	FailureRate        float64 `json:"failure_rate"`
	BudgetCompliance   float64 `json:"budget_compliance"`
	CriticalAlerts     int     `json:"critical_alerts"`
	WarningAlerts      int     `json:"warning_alerts"`
}

type OptimizationDashboard struct {
	TotalOptimizations   int     `json:"total_optimizations"`
	SuccessfulOptimizations int     `json:"successful_optimizations"`
	FailedOptimizations int     `json:"failed_optimizations"`
	AverageImprovement  float64 `json:"average_improvement"`
	ActiveOptimizations int     `json:"active_optimizations"`
}

func main() {
	// Parse command line flags
	config := &DashboardConfig{}
	flag.StringVar(&config.InputPath, "input", "", "Path to input performance data")
	flag.StringVar(&config.OutputPath, "output", "", "Path to output dashboard")
	flag.StringVar(&config.Title, "title", "Performance Dashboard", "Dashboard title")
	flag.BoolVar(&config.AutoRefresh, "auto-refresh", true, "Enable auto-refresh")
	flag.StringVar(&config.Theme, "theme", "light", "Dashboard theme (light/dark)")
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

	// Generate dashboard
	dashboard := generateDashboard(performanceData, config)

	// Generate HTML output
	if config.OutputPath != "" {
		generateHTMLDashboard(dashboard, config.OutputPath)
	}

	// Log summary
	logDashboardSummary(dashboard)
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
	
	ctx = monitor.StartOperation("file_op")
	time.Sleep(200 * time.Millisecond)
	ctx.End(nil)
	
	return monitor.GetCollector(), nil
}

func generateDashboard(collector *ui.PerformanceMetricsCollector, config *DashboardConfig) *PerformanceDashboard {
	dashboard := &PerformanceDashboard{
		Title:       config.Title,
		GeneratedAt: time.Now(),
	}

	// Generate operation data
	dashboard.Operations = generateOperationData(collector)
	
	// Generate alerts
	dashboard.Alerts = generateAlerts(collector)
	
	// Generate trends
	dashboard.Trends = generateTrends(collector)
	
	// Generate enforcement data
	dashboard.Enforcement = generateEnforcementData(collector)
	
	// Generate optimization data
	dashboard.Optimizations = generateOptimizationData(collector)
	
	// Calculate overall score and status
	dashboard.OverallScore, dashboard.Status = calculateOverallScore(dashboard)
	
	return dashboard
}

func generateOperationData(collector *ui.PerformanceMetricsCollector) []OperationDashboard {
	operations := []OperationDashboard{}
	
	// Define operations to track
	operationNames := []string{"first_paint", "frame_render", "keystroke", "file_op", "lsp_startup", "doctor_run", "dashboard_load"}
	
	budgets := ui.NewPerformanceBudget()
	
	for _, name := range operationNames {
		metrics := collector.GetRecentMetrics(name, 10)
		if len(metrics) == 0 {
			continue
		}
		
		// Calculate average duration
		avgDuration, hasAvg := collector.GetAverageDuration(name)
		if !hasAvg {
			continue
		}
		
		// Get budget
		budget := getBudgetForOperation(name, budgets)
		
		// Calculate compliance
		compliance, hasCompliance := collector.GetBudgetCompliance(name)
		if !hasCompliance {
			compliance = 0
		}
		
		// Determine status
		status := "good"
		if compliance < 0.5 {
			status = "critical"
		} else if compliance < 0.8 {
			status = "warning"
		}
		
		// Calculate score
		score := calculateOperationScore(avgDuration, budget, compliance)
		
		// Determine trend
		trend := determineTrend(metrics)
		
		// Determine impact
		impact := determineImpact(name, avgDuration)
		
		operation := OperationDashboard{
			Name:        name,
			Status:      status,
			Score:       score,
			AvgDuration: avgDuration,
			Budget:      budget,
			Compliance:  compliance,
			Trend:       trend,
			LastUpdated: metrics[len(metrics)-1].Timestamp,
			Impact:      impact,
		}
		
		operations = append(operations, operation)
	}
	
	return operations
}

func getBudgetForOperation(name string, budgets *ui.PerformanceBudget) time.Duration {
	switch name {
	case "first_paint":
		return budgets.FirstPaint
	case "frame_render":
		return budgets.FrameRender
	case "keystroke":
		return budgets.Keystroke
	case "file_op":
		return budgets.FileOp
	case "lsp_startup":
		return budgets.LSPStartup
	case "doctor_run":
		return budgets.DoctorRun
	case "dashboard_load":
		return budgets.DashboardLoad
	default:
		return 0
	}
}

func calculateOperationScore(duration, budget time.Duration, compliance float64) float64 {
	// Calculate score based on duration and compliance
	score := 100.0
	
	// Deduct points for exceeding budget
	if budget > 0 {
		ratio := float64(duration) / float64(budget)
		if ratio > 1 {
			score -= (ratio - 1) * 50
		}
	}
	
	// Deduct points for low compliance
	score -= (1 - compliance) * 30
	
	// Ensure score is between 0 and 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}
	
	return score
}

func determineTrend(metrics []ui.PerformanceMeasurement) string {
	if len(metrics) < 2 {
		return "stable"
	}
	
	latest := metrics[len(metrics)-1]
	earliest := metrics[0]
	
	ratio := float64(latest.Duration) / float64(earliest.Duration)
	
	if ratio > 1.1 {
		return "degrading"
	} else if ratio < 0.9 {
		return "improving"
	}
	
	return "stable"
}

func determineImpact(name string, duration time.Duration) string {
	switch name {
	case "first_paint":
		if duration > 100*time.Millisecond {
			return "high"
		}
	case "frame_render":
		if duration > 16*time.Millisecond {
			return "high"
		}
	case "keystroke":
		if duration > 50*time.Millisecond {
			return "high"
		}
	}
	
	return "medium"
}

func generateAlerts(collector *ui.PerformanceMetricsCollector) []DashboardAlert {
	alerts := []DashboardAlert{}
	
	// Check for budget violations
	operations := []string{"first_paint", "frame_render", "keystroke"}
	
	for _, name := range operations {
		metrics := collector.GetRecentMetrics(name, 5)
		if len(metrics) == 0 {
			continue
		}
		
		latest := metrics[len(metrics)-1]
		if !latest.WithinBudget && latest.Warning != nil {
			alert := DashboardAlert{
				Type:      "budget_violation",
				Operation: name,
				Severity:  "warning",
				Message:   fmt.Sprintf("%s exceeded budget: %s (budget %s)", name, latest.Duration, latest.Warning.Budget),
				Timestamp: latest.Timestamp,
				Resolved:  false,
			}
			alerts = append(alerts, alert)
		}
	}
	
	return alerts
}

func generateTrends(collector *ui.PerformanceMetricsCollector) []TrendDashboard {
	trends := []TrendDashboard{}
	
	// Analyze trends for key operations
	operations := []string{"first_paint", "frame_render", "keystroke"}
	
	for _, name := range operations {
		metrics := collector.GetRecentMetrics(name, 10)
		if len(metrics) < 5 {
			continue
		}
		
		latest := metrics[len(metrics)-1]
		baselineAvg := calculateAverage(metrics[:5])
		currentAvg := calculateAverage(metrics[len(metrics)-5:])
		
		changeRate := float64(currentAvg-baselineAvg) / float64(baselineAvg)
		direction := "stable"
		if changeRate > 0.1 {
			direction = "up"
		} else if changeRate < -0.1 {
			direction = "down"
		}
		
		trend := TrendDashboard{
			Name:         name,
			Direction:     direction,
			ChangeRate:   changeRate,
			CurrentAvg:   currentAvg,
			BaselineAvg:  baselineAvg,
			Confidence:   calculateConfidence(len(metrics)),
		}
		
		trends = append(trends, trend)
	}
	
	return trends
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
	if samples < 5 {
		return 0.5
	} else if samples < 10 {
		return 0.7
	} else {
		return 0.9
	}
}

func generateEnforcementData(collector *ui.PerformanceMetricsCollector) EnforcementDashboard {
	// Mock enforcement data
	return EnforcementDashboard{
		TotalOperations:    100,
		FailedOperations:   5,
		FailureRate:        0.05,
		BudgetCompliance:   0.95,
		CriticalAlerts:     1,
		WarningAlerts:      2,
	}
}

func generateOptimizationData(collector *ui.PerformanceMetricsCollector) OptimizationDashboard {
	// Mock optimization data
	return OptimizationDashboard{
		TotalOptimizations:     10,
		SuccessfulOptimizations: 8,
		FailedOptimizations:     2,
		AverageImprovement:      0.25,
		ActiveOptimizations:     5,
	}
}

func calculateOverallScore(dashboard *PerformanceDashboard) (float64, string) {
	// Calculate overall score based on all metrics
	score := 100.0
	
	// Deduct points for alerts
	score -= float64(len(dashboard.Alerts)) * 5
	
	// Deduct points for failed operations
	if dashboard.Enforcement.TotalOperations > 0 {
		score -= float64(dashboard.Enforcement.FailedOperations) * 2
	}
	
	// Add points for successful optimizations
	score += dashboard.Optimizations.AverageImprovement * 10
	
	// Ensure score is between 0 and 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}
	
	// Determine status
	status := "excellent"
	if score < 70 {
		status = "poor"
	} else if score < 85 {
		status = "fair"
	} else if score < 95 {
		status = "good"
	}
	
	return score, status
}

func generateHTMLDashboard(dashboard *PerformanceDashboard, outputPath string) {
	// Create HTML template
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .header { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .score { font-size: 2em; font-weight: bold; margin: 10px 0; }
        .status { font-size: 1.2em; font-weight: bold; margin: 10px 0; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .card h3 { margin-top: 0; color: #333; }
        .operation { padding: 10px; margin: 10px 0; border-radius: 4px; border-left: 4px solid #ddd; }
        .operation.good { border-left-color: #28a745; }
        .operation.warning { border-left-color: #ffc107; }
        .operation.critical { border-left-color: #dc3545; }
        .alert { padding: 10px; margin: 10px 0; border-radius: 4px; }
        .alert.warning { background: #fff3cd; border: 1px solid #ffeaa7; }
        .alert.critical { background: #f8d7da; border: 1px solid #f5c6cb; }
        .trend { padding: 10px; margin: 10px 0; border-radius: 4px; border: 1px solid #ddd; }
        .trend.up { border-left-color: #28a745; }
        .trend.down { border-left-color: #dc3545; }
        .trend.stable { border-left-color: #6c757d; }
        .metric { display: flex; justify-content: space-between; margin: 5px 0; }
        .metric-label { font-weight: bold; }
        .metric-value { color: #666; }
        .chart { height: 200px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 40px; color: #666; }
        .refresh { float: right; cursor: pointer; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Title}}</h1>
        <div class="score">{{.OverallScore}}%</div>
        <div class="status">Status: {{.Status}}</div>
        <div class="refresh" onclick="location.reload()">🔄 Refresh</div>
        <p>Generated on: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>
    </div>
    
    <div class="grid">
        <div class="card">
            <h3>Operations</h3>
            {{range .Operations}}
            <div class="operation {{.Status}}">
                <div class="metric">
                    <span class="metric-label">{{.Name}}</span>
                    <span class="metric-value">{{.Score}}%</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Average:</span>
                    <span class="metric-value">{{.AvgDuration}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Compliance:</span>
                    <span class="metric-value">{{printf "%.1f%%" .Compliance}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Trend:</span>
                    <span class="metric-value">{{.Trend}}</span>
                </div>
            </div>
            {{end}}
        </div>
        
        <div class="card">
            <h3>Alerts</h3>
            {{if .Alerts}}
            {{range .Alerts}}
            <div class="alert {{.Severity}}">
                <div class="metric">
                    <span class="metric-label">{{.Type}}</span>
                    <span class="metric-value">{{.Operation}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">{{.Message}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">{{.Timestamp.Format "15:04:05"}}</span>
                </div>
            </div>
            {{end}}
            {{else}}
            <p>No alerts</p>
            {{end}}
        </div>
        
        <div class="card">
            <h3>Trends</h3>
            {{range .Trends}}
            <div class="trend {{.Direction}}">
                <div class="metric">
                    <span class="metric-label">{{.Name}}</span>
                    <span class="metric-value">{{printf "%.1f%%" .ChangeRate}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Current:</span>
                    <span class="metric-value">{{.CurrentAvg}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Baseline:</span>
                    <span class="metric-value">{{.BaselineAvg}}</span>
                </div>
                <div class="metric">
                    <span class="metric-label">Confidence:</span>
                    <span class="metric-value">{{printf "%.1f" .Confidence}}</span>
                </div>
            </div>
            {{end}}
        </div>
        
        <div class="card">
            <h3>Enforcement</h3>
            <div class="metric">
                <span class="metric-label">Total Operations:</span>
                <span class="metric-value">{{.Enforcement.TotalOperations}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Failed Operations:</span>
                <span class="metric-value">{{.Enforcement.FailedOperations}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Failure Rate:</span>
                <span class="metric-value">{{printf "%.1f%%" .Enforcement.FailureRate}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Budget Compliance:</span>
                <span class="metric-value">{{printf "%.1f%%" .Enforcement.BudgetCompliance}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Critical Alerts:</span>
                <span class="metric-value">{{.Enforcement.CriticalAlerts}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Warning Alerts:</span>
                <span class="metric-value">{{.Enforcement.WarningAlerts}}</span>
            </div>
        </div>
        
        <div class="card">
            <h3>Optimizations</h3>
            <div class="metric">
                <span class="metric-label">Total Optimizations:</span>
                <span class="metric-value">{{.Optimizations.TotalOptimizations}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Successful:</span>
                <span class="metric-value">{{.Optimizations.SuccessfulOptimizations}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Failed:</span>
                <span class="metric-value">{{.Optimizations.FailedOptimizations}}</span>
            </div>
            <div class="metric">
                <span class="metric-label">Average Improvement:</span>
                <span class="metric-value">{{printf "%.1f%%" .Optimizations.AverageImprovement}}</span>
            </div>
            <div class="metric">
                <span class="metric-label"> Active Optimizations:</span>
                <span class="metric-value">{{.Optimizations.ActiveOptimizations}}</span>
            </div>
        </div>
    </div>
    
    <div class="footer">
        <p>Performance Dashboard - Auto-refresh every 30 seconds</p>
    </div>
    
    <script>
        // Auto-refresh functionality
        {{if .AutoRefresh}}
        setTimeout(function() {
            location.reload();
        }, 30000);
        {{end}}
    </script>
</body>
</html>`

	// Parse template
	tmpl, err := template.New("dashboard").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Create file
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, dashboard); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}
}

func logDashboardSummary(dashboard *PerformanceDashboard) {
	log.Printf("Performance Dashboard Summary:")
	log.Printf("  Title: %s", dashboard.Title)
	log.Printf("  Generated At: %s", dashboard.GeneratedAt.Format("2006-01-02 15:04:05"))
	log.Printf("  Overall Score: %.1f%%", dashboard.OverallScore)
	log.Printf("  Status: %s", dashboard.Status)
	log.Printf("  Operations: %d", len(dashboard.Operations))
	log.Printf("  Alerts: %d", len(dashboard.Alerts))
	log.Printf("  Trends: %d", len(dashboard.Trends))
	log.Printf("  Enforcement: Total=%d, Failed=%d, Rate=%.1f%%", 
		dashboard.Enforcement.TotalOperations, 
		dashboard.Enforcement.FailedOperations, 
		dashboard.Enforcement.FailureRate*100)
	log.Printf("  Optimizations: Total=%d, Success=%d, Improvement=%.1f%%", 
		dashboard.Optimizations.TotalOptimizations, 
		dashboard.Optimizations.SuccessfulOptimizations, 
		dashboard.Optimizations.AverageImprovement*100)
}