//go:build integration

package ui

import (
	"testing"
	"time"
)

// TestDashboardScreen tests the dashboard screen functionality
func TestDashboardScreen(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test dashboard screen creation
	dashboard := NewDashboardScreen(theme, keys, a11y, perfTracker)
	if dashboard == nil {
		t.Fatal("Dashboard screen should not be nil")
	}
	
	// Test dashboard screen properties
	if dashboard.GetID() != "dashboard" {
		t.Error("Dashboard screen ID should be 'dashboard'")
	}
	
	if !dashboard.IsActive() {
		t.Error("Dashboard screen should be active")
	}
	
	// Test dashboard data generation
	sessions := dashboard.generateMockSessions()
	if len(sessions) == 0 {
		t.Error("Should generate mock sessions")
	}
	
	for _, session := range sessions {
		if session.ID == "" {
			t.Error("Session should have an ID")
		}
		if session.Agent == "" {
			t.Error("Session should have an agent")
		}
		if session.Prompt == "" {
			t.Error("Session should have a prompt")
		}
		if session.Duration <= 0 {
			t.Error("Session should have a positive duration")
		}
		if session.Status == "" {
			t.Error("Session should have a status")
		}
		if session.Timestamp.IsZero() {
			t.Error("Session should have a timestamp")
		}
	}
	
	decisions := dashboard.generateMockDecisions()
	if len(decisions) == 0 {
		t.Error("Should generate mock decisions")
	}
	
	for _, decision := range decisions {
		if decision.ID == "" {
			t.Error("Decision should have an ID")
		}
		if decision.Title == "" {
			t.Error("Decision should have a title")
		}
		if decision.Description == "" {
			t.Error("Decision should have a description")
		}
		if decision.Reason == "" {
			t.Error("Decision should have a reason")
		}
		if decision.Outcome == "" {
			t.Error("Decision should have an outcome")
		}
		if decision.Timestamp.IsZero() {
			t.Error("Decision should have a timestamp")
		}
	}
	
	activity := dashboard.generateMockActivity()
	if len(activity) == 0 {
		t.Error("Should generate mock activity")
	}
	
	for _, item := range activity {
		if item.Agent == "" {
			t.Error("Activity item should have an agent")
		}
		if item.Action == "" {
			t.Error("Activity item should have an action")
		}
		if item.Target == "" {
			t.Error("Activity item should have a target")
		}
		if item.Duration <= 0 {
			t.Error("Activity item should have a positive duration")
		}
		if item.Timestamp.IsZero() {
			t.Error("Activity item should have a timestamp")
		}
	}
}

// TestDoctorScreen tests the doctor screen functionality
func TestDoctorScreen(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test doctor screen creation
	doctor := NewDoctorScreen(theme, keys, a11y, perfTracker)
	if doctor == nil {
		t.Fatal("Doctor screen should not be nil")
	}
	
	// Test doctor screen properties
	if doctor.GetID() != "doctor" {
		t.Error("Doctor screen ID should be 'doctor'")
	}
	
	if !doctor.IsActive() {
		t.Error("Doctor screen should be active")
	}
	
	// Test doctor data generation
	results := doctor.generateMockResults()
	if len(results) == 0 {
		t.Error("Should generate mock results")
	}
	
	for _, result := range results {
		if result.ID == "" {
			t.Error("Result should have an ID")
		}
		if result.Name == "" {
			t.Error("Result should have a name")
		}
		if result.Status == "" {
			t.Error("Result should have a status")
		}
		if result.Message == "" {
			t.Error("Result should have a message")
		}
		if result.Timestamp.IsZero() {
			t.Error("Result should have a timestamp")
		}
		if result.Duration <= 0 {
			t.Error("Result should have a positive duration")
		}
	}
	
	// Test overall status calculation
	status := doctor.getOverallStatus()
	if status == "" {
		t.Error("Overall status should not be empty")
	}
	
	// Test status color
	color := doctor.getStatusColor(status)
	if color == "" {
		t.Error("Status color should not be empty")
	}
}

// TestMetricsScreen tests the metrics screen functionality
func TestMetricsScreen(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test metrics screen creation
	metrics := NewMetricsScreen(theme, keys, a11y, perfTracker)
	if metrics == nil {
		t.Fatal("Metrics screen should not be nil")
	}
	
	// Test metrics screen properties
	if metrics.GetID() != "metrics" {
		t.Error("Metrics screen ID should be 'metrics'")
	}
	
	if !metrics.IsActive() {
		t.Error("Metrics screen should be active")
	}
	
	// Test metrics data generation
	sessionStats := metrics.generateMockSessionStats()
	if sessionStats.TotalSessions == 0 {
		t.Error("Should generate session statistics")
	}
	
	if len(sessionStats.AgentBreakdown) == 0 {
		t.Error("Should have agent breakdown")
	}
	
	usageStats := metrics.generateMockUsageStats()
	if usageStats.TotalFilesModified == 0 {
		t.Error("Should generate usage statistics")
	}
	
	if len(usageStats.TopFiles) == 0 {
		t.Error("Should have top files")
	}
	
	if len(usageStats.TopAgents) == 0 {
		t.Error("Should have top agents")
	}
	
	performanceData := metrics.generateMockPerformanceData()
	if performanceData.AvgResponseTime <= 0 {
		t.Error("Should have performance data")
	}
	
	if performanceData.SuccessRate <= 0 {
		t.Error("Should have success rate")
	}
	
	if performanceData.Throughput <= 0 {
		t.Error("Should have throughput")
	}
}

// TestCostScreen tests the cost screen functionality
func TestCostScreen(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test cost screen creation
	cost := NewCostScreen(theme, keys, a11y, perfTracker)
	if cost == nil {
		t.Fatal("Cost screen should not be nil")
	}
	
	// Test cost screen properties
	if cost.GetID() != "cost" {
		t.Error("Cost screen ID should be 'cost'")
	}
	
	if !cost.IsActive() {
		t.Error("Cost screen should be active")
	}
	
	// Test cost data generation
	totalCost := 45.67
	monthlyBudget := 100.0
	remainingBudget := monthlyBudget - totalCost
	
	usageByAgent := map[string]float64{
		"opencode":      25.34,
		"claude-code":   15.89,
		"universal-lsp": 4.44,
	}
	
	costByFile := map[string]float64{
		"internal/api/server.go":      12.45,
		"internal/database/connection.go": 8.92,
		"cmd/main.go":                 6.78,
		"internal/auth/middleware.go": 5.34,
		"internal/config/config.go":   4.23,
	}
	
	recentCosts := cost.generateMockRecentCosts()
	if len(recentCosts) == 0 {
		t.Error("Should generate recent costs")
	}
	
	for _, costEntry := range recentCosts {
		if costEntry.ID == "" {
			t.Error("Cost entry should have an ID")
		}
		if costEntry.Agent == "" {
			t.Error("Cost entry should have an agent")
		}
		if costEntry.File == "" {
			t.Error("Cost entry should have a file")
		}
		if costEntry.Action == "" {
			t.Error("Cost entry should have an action")
		}
		if costEntry.Cost <= 0 {
			t.Error("Cost entry should have a positive cost")
		}
		if costEntry.Timestamp.IsZero() {
			t.Error("Cost entry should have a timestamp")
		}
	}
	
	costTrend := cost.generateMockCostTrend()
	if len(costTrend) == 0 {
		t.Error("Should generate cost trend")
	}
	
	for _, point := range costTrend {
		if point.Date.IsZero() {
			t.Error("Cost point should have a date")
		}
		if point.Cost <= 0 {
			t.Error("Cost point should have a positive cost")
		}
	}
	
	alerts := cost.generateMockAlerts(totalCost, monthlyBudget)
	if len(alerts) == 0 {
		t.Error("Should generate alerts")
	}
	
	for _, alert := range alerts {
		if alert.Type == "" {
			t.Error("Alert should have a type")
		}
		if alert.Message == "" {
			t.Error("Alert should have a message")
		}
	}
}

// TestSettingsScreen tests the settings screen functionality
func TestSettingsScreen(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test settings screen creation
	settings := NewSettingsScreen(theme, keys, a11y, perfTracker)
	if settings == nil {
		t.Fatal("Settings screen should not be nil")
	}
	
	// Test settings screen properties
	if settings.GetID() != "settings" {
		t.Error("Settings screen ID should be 'settings'")
	}
	
	if !settings.IsActive() {
		t.Error("Settings screen should be active")
	}
	
	// Test settings data generation
	general := GeneralSettings{
		AutoRefresh:     true,
		RefreshInterval: 30 * time.Second,
		Theme:           "default",
		Language:        "en",
		Timezone:        "UTC",
	}
	
	if general.AutoRefresh != true {
		t.Error("General settings should have auto refresh enabled")
	}
	
	accessibility := AccessibilitySettings{
		HighContrast:    false,
		ReducedMotion:   false,
		ScreenReader:   false,
		LargeText:      false,
		KeyboardOnly:   true,
		TabNavigation:  true,
		EscapeExits:    true,
		ConfirmActions: true,
	}
	
	if accessibility.KeyboardOnly != true {
		t.Error("Accessibility settings should have keyboard only enabled")
	}
	
	notifications := NotificationSettings{
		EnableEmail:     false,
		EnableDesktop:  true,
		EnableSound:    true,
		EmailAddress:   "",
		AlertThreshold: 80,
	}
	
	if notifications.EnableDesktop != true {
		t.Error("Notification settings should have desktop notifications enabled")
	}
	
	advanced := AdvancedSettings{
		DebugMode:      false,
		LogLevel:      "info",
		DataRetention: 30,
		AutoBackup:    true,
		BackupPath:    "~/.config/caf/backups/",
	}
	
	if advanced.DataRetention != 30 {
		t.Error("Advanced settings should have 30 day retention")
	}
}

// TestDashboardDataFormatting tests the dashboard data formatting
func TestDashboardDataFormatting(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	dashboard := NewDashboardScreen(theme, keys, a11y, perfTracker)
	
	// Test session statistics formatting
	sessionStats := dashboard.generateMockSessionStats()
	sessionContent := dashboard.formatSessionStats()
	
	if sessionContent == "" {
		t.Error("Session stats content should not be empty")
	}
	
	// Test decision statistics formatting
	decisions := dashboard.generateMockDecisions()
	decisionContent := dashboard.formatDecisionStats([]DecisionItem{})
	
	if decisionContent == "" {
		t.Error("Decision stats content should not be empty")
	}
	
	// Test activity formatting
	activity := dashboard.generateMockActivity()
	activityContent := dashboard.formatActivityStats([]AgentActivityItem{})
	
	if activityContent == "" {
		t.Error("Activity stats content should not be empty")
	}
}

// TestDoctorDataFormatting tests the doctor data formatting
func TestDoctorDataFormatting(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	doctor := NewDoctorScreen(theme, keys, a11y, perfTracker)
	
	// Test health check results formatting
	results := doctor.generateMockResults()
	
	for _, result := range results {
		statusIcon, statusColor := doctor.getStatusIconAndColor(result.Status)
		
		if statusIcon == "" {
			t.Error("Status icon should not be empty")
		}
		
		if statusColor == "" {
			t.Error("Status color should not be empty")
		}
	}
	
	// Test duration formatting
	duration := doctor.formatDuration(500 * time.Millisecond)
	if duration == "" {
		t.Error("Duration formatting should not be empty")
	}
}

// TestMetricsDataFormatting tests the metrics data formatting
func TestMetricsDataFormatting(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	metrics := NewMetricsScreen(theme, keys, a11y, perfTracker)
	
	// Test session statistics formatting
	sessionStats := metrics.generateMockSessionStats()
	sessionContent := metrics.formatSessionStats()
	
	if sessionContent == "" {
		t.Error("Session stats content should not be empty")
	}
	
	// Test usage statistics formatting
	usageStats := metrics.generateMockUsageStats()
	usageContent := metrics.formatUsageStats()
	
	if usageContent == "" {
		t.Error("Usage stats content should not be empty")
	}
	
	// Test performance data formatting
	performanceData := metrics.generateMockPerformanceData()
	performanceContent := metrics.formatPerformanceData()
	
	if performanceContent == "" {
		t.Error("Performance data content should not be empty")
	}
}

// TestCostDataFormatting tests the cost data formatting
func TestCostDataFormatting(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	cost := NewCostScreen(theme, keys, a11y, perfTracker)
	
	// Test cost overview formatting
	overviewContent := cost.formatOverview()
	
	if overviewContent == "" {
		t.Error("Cost overview content should not be empty")
	}
	
	// Test budget formatting
	budgetContent := cost.formatBudget()
	
	if budgetContent == "" {
		t.Error("Budget content should not be empty")
	}
	
	// Test agent cost formatting
	agentContent := cost.formatAgentCosts()
	
	if agentContent == "" {
		t.Error("Agent cost content should not be empty")
	}
	
	// Test file cost formatting
	fileContent := cost.formatFileCosts()
	
	if fileContent == "" {
		t.Error("File cost content should not be empty")
	}
	
	// Test trend formatting
	trendContent := cost.formatTrend()
	
	if trendContent == "" {
		t.Error("Trend content should not be empty")
	}
}

// TestSettingsDataFormatting tests the settings data formatting
func TestSettingsDataFormatting(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	settings := NewSettingsScreen(theme, keys, a11y, perfTracker)
	
	// Test general settings formatting
	generalContent := settings.formatGeneral()
	
	if generalContent == "" {
		t.Error("General settings content should not be empty")
	}
	
	// Test accessibility settings formatting
	accessibilityContent := settings.formatAccessibility()
	
	if accessibilityContent == "" {
		t.Error("Accessibility settings content should not be empty")
	}
	
	// Test notification settings formatting
	notificationsContent := settings.formatNotifications()
	
	if notificationsContent == "" {
		t.Error("Notification settings content should not be empty")
	}
	
	// Test advanced settings formatting
	advancedContent := settings.formatAdvanced()
	
	if advancedContent == "" {
		t.Error("Advanced settings content should not be empty")
	}
}

// TestScreenHelpText tests the help text for all screens
func TestScreenHelpText(t *testing.T) {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	
	// Test dashboard help
	dashboard := NewDashboardScreen(theme, keys, a11y, perfTracker)
	dashboardHelp := dashboard.getDashboardHelp()
	
	if dashboardHelp == "" {
		t.Error("Dashboard help text should not be empty")
	}
	
	// Test doctor help
	doctor := NewDoctorScreen(theme, keys, a11y, perfTracker)
	doctorHelp := doctor.getDoctorHelp()
	
	if doctorHelp == "" {
		t.Error("Doctor help text should not be empty")
	}
	
	// Test metrics help
	metrics := NewMetricsScreen(theme, keys, a11y, perfTracker)
	metricsHelp := metrics.getMetricsHelp()
	
	if metricsHelp == "" {
		t.Error("Metrics help text should not be empty")
	}
	
	// Test cost help
	cost := NewCostScreen(theme, keys, a11y, perfTracker)
	costHelp := cost.getCostHelp()
	
	if costHelp == "" {
		t.Error("Cost help text should not be empty")
	}
	
	// Test settings help
	settings := NewSettingsScreen(theme, keys, a11y, perfTracker)
	settingsHelp := settings.getSettingsHelp()
	
	if settingsHelp == "" {
		t.Error("Settings help text should not be empty")
	}
}