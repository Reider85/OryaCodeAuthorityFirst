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

type PerformanceTestConfig struct {
	Threshold          float64 `json:"threshold"`
	MaxRegression      float64 `json:"max_regression"`
	MinBudgetCompliance float64 `json:"min_budget_compliance"`
	BaselinePath       string  `json:"baseline_path"`
	OutputPath         string  `json:"output_path"`
	Comprehensive      bool    `json:"comprehensive"`
	CheckRegression    bool    `json:"check_regression"`
	GenerateReport     bool    `json:"generate_report"`
	QualityGate        bool    `json:"quality_gate"`
}

type PerformanceReport struct {
	Timestamp     time.Time              `json:"timestamp"`
	Commit        string                 `json:"commit"`
	Branch        string                 `json:"branch"`
	Metrics       map[string][]ui.PerformanceMeasurement `json:"metrics"`
	Enforcement   ui.EnforcementStats    `json:"enforcement"`
	Optimizations ui.OptimizationStats  `json:"optimizations"`
	Regressions   []ui.RegressionAlert  `json:"regressions"`
	Score         float64                `json:"score"`
	Status        string                 `json:"status"`
	Issues        []PerformanceIssue    `json:"issues"`
}

type PerformanceIssue struct {
	Type        string `json:"type"`
	Operation   string `json:"operation"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Recommendation string `json:"recommendation"`
}

func main() {
	// Parse command line flags
	config := &PerformanceTestConfig{}
	flag.Float64Var(&config.Threshold, "threshold", 0.1, "Performance threshold")
	flag.Float64Var(&config.MaxRegression, "max-regression", 0.15, "Maximum allowed regression")
	flag.Float64Var(&config.MinBudgetCompliance, "min-budget-compliance", 0.9, "Minimum budget compliance")
	flag.StringVar(&config.BaselinePath, "baseline", "", "Path to baseline performance data")
	flag.StringVar(&config.OutputPath, "output", "", "Path to output report")
	flag.BoolVar(&config.Comprehensive, "comprehensive", false, "Run comprehensive performance tests")
	flag.BoolVar(&config.CheckRegression, "check-regression", false, "Check for performance regressions")
	flag.BoolVar(&config.GenerateReport, "report", false, "Generate performance report")
	flag.BoolVar(&config.QualityGate, "quality-gate", false, "Run performance quality gate")
	flag.Parse()

	// Initialize performance monitoring system
	budgets := ui.NewPerformanceBudget()
	monitor := ui.NewPerformanceMonitor(budgets)
	enforcer := ui.NewBudgetEnforcer(monitor, budgets)
	regression := ui.NewRegressionDetector(monitor, budgets)
	optimizer := ui.NewPerformanceOptimizer(monitor, enforcer, regression)

	// Set up global monitor
	ui.InitializeGlobalMonitor(budgets)

	// Run performance tests
	if config.Comprehensive {
		runComprehensiveTests(monitor, enforcer, regression, optimizer)
	}

	// Check for regressions
	if config.CheckRegression {
		checkRegressions(regression, config.BaselinePath)
	}

	// Generate report
	if config.GenerateReport {
		generateReport(monitor, enforcer, regression, optimizer, config)
	}

	// Run quality gate
	if config.QualityGate {
		runQualityGate(enforcer, regression, config)
	}
}

func runComprehensiveTests(monitor *ui.PerformanceMonitor, enforcer *ui.BudgetEnforcer, regression *ui.RegressionDetector, optimizer *ui.PerformanceOptimizer) {
	log.Println("Running comprehensive performance tests...")

	// Test first paint performance
	testFirstPaint(monitor)

	// Test frame render performance
	testFrameRender(monitor)

	// Test keystroke performance
	testKeystroke(monitor)

	// Test file operations
	testFileOperations(monitor)

	// Test LSP startup
	testLSPStartup(monitor)

	// Test doctor run
	testDoctorRun(monitor)

	// Test dashboard load
	testDashboardLoad(monitor)

	// Run optimization tests
	optimizer.AutoOptimize()

	log.Println("Comprehensive performance tests completed")
}

func testFirstPaint(monitor *ui.PerformanceMonitor) {
	log.Println("Testing first paint performance...")
	
	ctx := monitor.StartOperation("first_paint")
	// Simulate first paint
	time.Sleep(50 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testFrameRender(monitor *ui.PerformanceMonitor) {
	log.Println("Testing frame render performance...")
	
	ctx := monitor.StartOperation("frame_render")
	// Simulate frame render
	time.Sleep(10 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testKeystroke(monitor *ui.PerformanceMonitor) {
	log.Println("Testing keystroke performance...")
	
	ctx := monitor.StartOperation("keystroke")
	// Simulate keystroke handling
	time.Sleep(20 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testFileOperations(monitor *ui.PerformanceMonitor) {
	log.Println("Testing file operations...")
	
	ctx := monitor.StartOperation("file_op")
	// Simulate file operations
	time.Sleep(200 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testLSPStartup(monitor *ui.PerformanceMonitor) {
	log.Println("Testing LSP startup...")
	
	ctx := monitor.StartOperation("lsp_startup")
	// Simulate LSP startup
	time.Sleep(300 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testDoctorRun(monitor *ui.PerformanceMonitor) {
	log.Println("Testing doctor run...")
	
	ctx := monitor.StartOperation("doctor_run")
	// Simulate doctor run
	time.Sleep(2000 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func testDashboardLoad(monitor *ui.PerformanceMonitor) {
	log.Println("Testing dashboard load...")
	
	ctx := monitor.StartOperation("dashboard_load")
	// Simulate dashboard load
	time.Sleep(300 * time.Millisecond) // Simulate work
	ctx.End(nil)
}

func checkRegressions(regression *ui.RegressionDetector, baselinePath string) {
	log.Println("Checking for performance regressions...")
	
	// Detect regressions
	alerts := regression.DetectRegressions()
	
	if len(alerts) > 0 {
		log.Printf("Found %d regressions:", len(alerts))
		for _, alert := range alerts {
			log.Printf("  - %s: %s", alert.Operation, alert.Message)
		}
		
		// Check against baseline if provided
		if baselinePath != "" {
			checkAgainstBaseline(regression, baselinePath)
		}
		
		os.Exit(1) // Exit with error if regressions found
	} else {
		log.Println("No regressions detected")
	}
}

func checkAgainstBaseline(regression *ui.RegressionDetector, baselinePath string) {
	log.Println("Checking against baseline...")
	
	// Load baseline data
	// In a real implementation, this would load from file
	// For now, just log
	log.Println("Baseline check not implemented yet")
}

func generateReport(monitor *ui.PerformanceMonitor, enforcer *ui.BudgetEnforcer, regression *ui.RegressionDetector, optimizer *ui.PerformanceOptimizer, config *PerformanceTestConfig) {
	log.Println("Generating performance report...")
	
	// Collect data
	metrics := monitor.GetMetrics()
	enforcementStats := enforcer.GetEnforcementStats()
	optimizationStats := optimizer.GetOptimizationStats()
	regressionAlerts := regression.DetectRegressions()
	
	// Calculate performance score
	score := calculatePerformanceScore(enforcementStats, optimizationStats, len(regressionAlerts))
	
	// Generate issues
	issues := generateIssues(metrics, enforcementStats, regressionAlerts)
	
	// Create report
	report := PerformanceReport{
		Timestamp:     time.Now(),
		Commit:        getCommitHash(),
		Branch:        getBranchName(),
		Metrics:       metrics,
		Enforcement:   enforcementStats,
		Optimizations: optimizationStats,
		Regressions:   regressionAlerts,
		Score:         score,
		Status:        determineStatus(score, len(issues)),
		Issues:        issues,
	}
	
	// Write report to file
	if config.OutputPath != "" {
		writeReport(report, config.OutputPath)
	} else {
		writeReport(report, "performance-report.json")
	}
	
	log.Printf("Performance report generated with score: %.2f", score)
}

func calculatePerformanceScore(enforcementStats ui.EnforcementStats, optimizationStats ui.OptimizationStats, regressionCount int) float64 {
	// Calculate score based on various metrics
	score := 100.0
	
	// Deduct points for failed operations
	if enforcementStats.TotalOperations > 0 {
		failureRate := float64(enforcementStats.FailedOperations) / float64(enforcementStats.TotalOperations)
		score -= failureRate * 30
	}
	
	// Deduct points for regressions
	score -= float64(regressionCount) * 5
	
	// Add points for successful optimizations
	if optimizationStats.SuccessfulOptimizations > 0 {
		improvementBonus := optimizationStats.AverageImprovement * 20
		score += improvementBonus
	}
	
	// Ensure score is between 0 and 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}
	
	return score
}

func generateIssues(metrics map[string][]ui.PerformanceMeasurement, enforcementStats ui.EnforcementStats, regressionAlerts []ui.RegressionAlert) []PerformanceIssue {
	issues := []PerformanceIssue{}
	
	// Generate issues from regressions
	for _, alert := range regressionAlerts {
		issue := PerformanceIssue{
			Type:        "regression",
			Operation:   alert.Operation,
			Severity:    string(alert.Severity),
			Message:     alert.Message,
			Recommendation: getRecommendation(alert.Operation),
		}
		issues = append(issues, issue)
	}
	
	// Generate issues from enforcement statistics
	if enforcementStats.FailureRate > 0.1 { // More than 10% failures
		issue := PerformanceIssue{
			Type:        "enforcement",
			Severity:    "warning",
			Message:     fmt.Sprintf("High failure rate: %.2f%%", enforcementStats.FailureRate*100),
			Recommendation: "Review performance budgets and optimize critical operations",
		}
		issues = append(issues, issue)
	}
	
	return issues
}

func getRecommendation(operation string) string {
	recommendations := map[string]string{
		"first_paint":  "Consider using Skeleton widgets for async data loading",
		"frame_render": "Use virtualized lists and memoization",
		"keystroke":    "Use direct dispatch and avoid async round-trips",
		"file_op":      "Implement caching and batch operations",
		"lsp_startup":  "Use lazy initialization and server pooling",
		"doctor_run":   "Implement incremental analysis",
		"dashboard_load": "Use background data loading and pagination",
	}
	
	if rec, exists := recommendations[operation]; exists {
		return rec
	}
	return "Review and optimize the operation"
}

func determineStatus(score float64, issueCount int) string {
	if score >= 90 {
		return "excellent"
	} else if score >= 70 {
		return "good"
	} else if score >= 50 {
		return "fair"
	} else if score >= 30 {
		return "poor"
	} else {
		return "critical"
	}
}

func writeReport(report PerformanceReport, outputPath string) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal report: %v", err)
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

func runQualityGate(enforcer *ui.BudgetEnforcer, regression *ui.RegressionDetector, config *PerformanceTestConfig) {
	log.Println("Running performance quality gate...")
	
	// Check enforcement statistics
	enforcementStats := enforcer.GetEnforcementStats()
	enforcementStats.CalculateFailureRate()
	
	if enforcementStats.FailureRate > config.MaxRegression {
		log.Printf("Quality gate failed: Failure rate %.2f%% exceeds threshold %.2f%%", 
			enforcementStats.FailureRate*100, config.MaxRegression*100)
		os.Exit(1)
	}
	
	// Check regression count
	regressionAlerts := regression.DetectRegressions()
	if len(regressionAlerts) > 0 {
		log.Printf("Quality gate failed: %d regressions detected", len(regressionAlerts))
		os.Exit(1)
	}
	
	// Check budget compliance
	operations := []string{"first_paint", "frame_render", "keystroke"}
	for _, op := range operations {
		if compliance, exists := enforcer.GetBudgetCompliance(op); exists {
			if compliance < config.MinBudgetCompliance {
				log.Printf("Quality gate failed: %s compliance %.2f%% below threshold %.2f%%", 
					op, compliance*100, config.MinBudgetCompliance*100)
				os.Exit(1)
			}
		}
	}
	
	log.Println("Quality gate passed")
}

func getCommitHash() string {
	// In a real implementation, this would get the current commit hash
	return "unknown"
}

func getBranchName() string {
	// In a real implementation, this would get the current branch name
	return "unknown"
}