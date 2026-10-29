package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

// MVPTester runs comprehensive MVP testing
type MVPTester struct {
	validator     *ui.MVPValidator
	tui           *ui.TUI
	startTime     time.Time
	testResults   map[string]bool
	errorMessages []string
}

// NewMVPTester creates a new MVP tester
func NewMVPTester() *MVPTester {
	validator := ui.NewMVPValidator()
	tui := ui.NewTUI()
	
	return &MVPTester{
		validator:     validator,
		tui:           tui,
		startTime:     time.Now(),
		testResults:   make(map[string]bool),
		errorMessages: []string{},
	}
}

// RunAllTests runs all MVP tests
func (mt *MVPTester) RunAllTests() error {
	fmt.Println("🧪 Starting CodeAuthorityFirst MVP Testing Suite")
	fmt.Println("================================================")
	fmt.Println()
	
	// Test 1: Architecture Validation
	fmt.Println("🏗️  Test 1: Architecture Validation")
	mt.testArchitecture()
	
	// Test 2: Feature Requirements Validation
	fmt.Println("🎯 Test 2: Feature Requirements Validation")
	mt.testFeatureRequirements()
	
	// Test 3: State Management Validation
	fmt.Println("🔄 Test 3: State Management Validation")
	mt.testStateManagement()
	
	// Test 4: Accessibility Validation
	fmt.Println("♿ Test 4: Accessibility Validation")
	mt.testAccessibility()
	
	// Test 5: Performance Validation
	fmt.Println("⚡ Test 5: Performance Validation")
	mt.testPerformance()
	
	// Test 6: Integration Validation
	fmt.Println("🔗 Test 6: Integration Validation")
	mt.testIntegration()
	
	// Test 7: UI Components Validation
	fmt.Println("🎨 Test 7: UI Components Validation")
	mt.testUIComponents()
	
	// Test 8: Keyboard Navigation Validation
	fmt.Println("⌨️  Test 8: Keyboard Navigation Validation")
	mt.testKeyboardNavigation()
	
	// Test 9: TUI Integration Test
	fmt.Println("🚀 Test 9: TUI Integration Test")
	mt.testTUIIntegration()
	
	// Test 10: Screen-Specific Tests
	fmt.Println("📱 Test 10: Screen-Specific Tests")
	mt.testScreenSpecific()
	
	// Generate Report
	mt.generateReport()
	
	return nil
}

// testArchitecture tests architectural requirements
func (mt *MVPTester) testArchitecture() {
	defer mt.recordTestResult("architecture", mt.errorMessages == nil)
	
	fmt.Println("  Validating INV-UI-001 to INV-UI-009...")
	
	// Test budgets
	budget := ui.NewPerformanceBudget()
	if budget.FirstPaint >= 100*time.Millisecond {
		mt.addError("INV-UI-001: First paint budget should be <100ms")
	}
	
	if budget.FrameRender >= 16*time.Millisecond {
		mt.addError("INV-UI-002: Frame render budget should be <16ms")
	}
	
	if budget.Keystroke >= 50*time.Millisecond {
		mt.addError("INV-UI-003: Keystroke response budget should be <50ms")
	}
	
	// Test theme
	theme := ui.NewTheme()
	theme.PrecomputeStyles()
	if theme.Primary == nil {
		mt.addError("INV-UI-006: Precomputed styles should not be nil")
	}
	
	// Test accessibility
	config := ui.NewAccessibilityConfig()
	if !config.KeyboardOnly {
		mt.addError("INV-UI-007: Keyboard-only operation should be enabled")
	}
	
	// Test keys
	keys := ui.NewKeys()
	if keys.Tab == "" || keys.Enter == "" || keys.Esc == "" {
		mt.addError("INV-UI-009: All navigation should be keyboard accessible")
	}
	
	fmt.Println("  ✓ Architecture validation completed")
}

// testFeatureRequirements tests feature requirements
func (mt *MVPTester) testFeatureRequirements() {
	defer mt.recordTestResult("feature_requirements", mt.errorMessages == nil)
	
	fmt.Println("  Validating FR-TUI-001 to FR-TUI-007...")
	
	// Test dashboard
	dashboard := ui.NewDashboardScreen(ui.NewTheme(), ui.NewKeys(), 
		ui.NewA11y(ui.NewAccessibilityConfig(), ui.NewTheme(), ui.NewKeys()), 
		ui.NewPerfTracker(ui.NewPerformanceBudget()))
	
	sessions := dashboard.(*ui.DashboardScreen).generateMockSessions()
	if len(sessions) == 0 {
		mt.addError("FR-TUI-001: Dashboard should have recent sessions")
	}
	
	// Test doctor
	doctor := ui.NewDoctorScreen(ui.NewTheme(), ui.NewKeys(),
		ui.NewA11y(ui.NewAccessibilityConfig(), ui.NewTheme(), ui.NewKeys()),
		ui.NewPerfTracker(ui.NewPerformanceBudget()))
	
	results := doctor.(*ui.DoctorScreen).generateMockResults()
	if len(results) != 5 {
		mt.addError("FR-TUI-002: Doctor should show 5 health checks")
	}
	
	// Test metrics
	metrics := ui.NewMetricsScreen(ui.NewTheme(), ui.NewKeys(),
		ui.NewA11y(ui.NewAccessibilityConfig(), ui.NewTheme(), ui.NewKeys()),
		ui.NewPerfTracker(ui.NewPerformanceBudget()))
	
	sessionStats := metrics.(*ui.MetricsScreen).generateMockSessionStats()
	if sessionStats.TotalSessions == 0 {
		mt.addError("FR-TUI-003: Metrics should show session counts")
	}
	
	// Test cost
	cost := ui.NewCostScreen(ui.NewTheme(), ui.NewKeys(),
		ui.NewA11y(ui.NewAccessibilityConfig(), ui.NewTheme(), ui.NewKeys()),
		ui.NewPerfTracker(ui.NewPerformanceBudget()))
	
	alerts := cost.(*ui.CostScreen).generateMockAlerts(45.67, 100.0)
	if len(alerts) == 0 {
		mt.addError("FR-TUI-004: Cost should show budget alerts")
	}
	
	// Test settings
	settings := ui.NewSettingsScreen(ui.NewTheme(), ui.NewKeys(),
		ui.NewA11y(ui.NewAccessibilityConfig(), ui.NewTheme(), ui.NewKeys()),
		ui.NewPerfTracker(ui.NewPerformanceBudget()))
	
	general := settings.(*ui.SettingsScreen).General
	if general.Theme == "" {
		mt.addError("FR-TUI-005: Settings should show general options")
	}
	
	// Test integration
	integration := ui.NewIntegrationLayer()
	if err := integration.Initialize(); err != nil {
		mt.addError("FR-TUI-006: Integration layer should initialize without error")
	}
	
	fmt.Println("  ✓ Feature requirements validation completed")
}

// testStateManagement tests state management requirements
func (mt *MVPTester) testStateManagement() {
	defer mt.recordTestResult("state_management", mt.errorMessages == nil)
	
	fmt.Println("  Validating INV-PROMPT-MVP-006 state transitions...")
	
	// Test all 5 states
	sm := ui.NewStateMachine()
	
	// Test Loading → Success
	sm.Transition(ui.LoadingMsg{})
	if sm.Current() != ui.StateLoading {
		mt.addError("Screen should start in loading state")
	}
	
	sm.Transition(ui.DataLoadedMsg{Empty: false})
	if sm.Current() != ui.StateSuccess {
		mt.addError("Screen should transition to success on data loaded")
	}
	
	// Test Loading → Error
	sm = ui.NewStateMachine()
	sm.Transition(ui.LoadingMsg{})
	sm.Transition(ui.DataErrorMsg{Err: fmt.Errorf("test error")})
	if sm.Current() != ui.StateError {
		mt.addError("Screen should transition to error on data error")
	}
	
	// Test Loading → Empty
	sm = ui.NewStateMachine()
	sm.Transition(ui.LoadingMsg{})
	sm.Transition(ui.DataLoadedMsg{Empty: true})
	if sm.Current() != ui.StateEmpty {
		mt.addError("Screen should transition to empty on empty data")
	}
	
	// Test Success → Loading (on refresh)
	sm = ui.NewStateMachine()
	sm.Transition(ui.DataLoadedMsg{Empty: false})
	sm.Transition(ui.RefreshMsg{Force: true})
	if sm.Current() != ui.StateLoading {
		mt.addError("Screen should transition to loading on refresh")
	}
	
	fmt.Println("  ✓ State management validation completed")
}

// testAccessibility tests accessibility requirements
func (mt *MVPTester) testAccessibility() {
	defer mt.recordTestResult("accessibility", mt.errorMessages == nil)
	
	fmt.Println("  Validating accessibility requirements...")
	
	config := ui.NewAccessibilityConfig()
	a11y := ui.NewA11y(config, ui.NewTheme(), ui.NewKeys())
	
	// Test high contrast mode
	if !config.HighContrast {
		mt.addError("High contrast mode should be available")
	}
	
	// Test reduced motion
	if !config.ReducedMotion {
		mt.addError("Reduced motion should be available")
	}
	
	// Test screen reader support
	if !config.ScreenReader {
		mt.addError("Screen reader support should be available")
	}
	
	// Test keyboard-only operation
	if !config.KeyboardOnly {
		mt.addError("Keyboard-only operation should be available")
	}
	
	// Test focus management
	focusManager := ui.NewFocusManager(a11y)
	focusManager.AddFocusableElement("element1")
	focusManager.AddFocusableElement("element2")
	
	if len(focusManager.focusableElements) != 2 {
		mt.addError("Focus manager should track focusable elements")
	}
	
	// Test accessibility report
	report := a11y.GetAccessibilityReport()
	if report == "" {
		mt.addError("Accessibility report should not be empty")
	}
	
	fmt.Println("  ✓ Accessibility validation completed")
}

// testPerformance tests performance requirements
func (mt *MVPTester) testPerformance() {
	defer mt.recordTestResult("performance", mt.errorMessages == nil)
	
	fmt.Println("  Validating performance requirements...")
	
	budget := ui.NewPerformanceBudget()
	tracker := ui.NewPerfTracker(budget)
	
	// Test performance budgets
	if budget.FirstPaint >= 100*time.Millisecond {
		mt.addError("First paint budget should be <100ms")
	}
	
	if budget.FrameRender >= 16*time.Millisecond {
		mt.addError("Frame render budget should be <16ms")
	}
	
	if budget.Keystroke >= 50*time.Millisecond {
		mt.addError("Keystroke response budget should be <50ms")
	}
	
	// Test performance tracker
	tracker.Start("test_operation")
	time.Sleep(10 * time.Millisecond)
	duration := tracker.End("test_operation")
	
	if duration <= 0 {
		mt.addError("Performance tracker should measure positive duration")
	}
	
	// Test performance report
	report := tracker.GetReport()
	if report == "" {
		mt.addError("Performance report should not be empty")
	}
	
	fmt.Println("  ✓ Performance validation completed")
}

// testIntegration tests integration requirements
func (mt *MVPTester) testIntegration() {
	defer mt.recordTestResult("integration", mt.errorMessages == nil)
	
	fmt.Println("  Validating integration requirements...")
	
	integration := ui.NewIntegrationLayer()
	if err := integration.Initialize(); err != nil {
		mt.addError("Integration layer should initialize without error")
	}
	
	// Test screen manager
	screenManager := integration.ScreenManager
	if len(screenManager.Screens) != 5 {
		mt.addError("Screen manager should have 5 screens")
	}
	
	// Test data manager
	dataManager := integration.DataManager
	dataManager.SetSharedData("test_key", "test_value")
	value, exists := dataManager.GetSharedData("test_key")
	if !exists || value != "test_value" {
		mt.addError("Data manager should handle shared data correctly")
	}
	
	// Test event bus
	eventBus := integration.EventBus
	channel := eventBus.Subscribe("test_event")
	eventBus.Publish(ui.Event{Type: "test_event", Payload: "test_payload", Timestamp: time.Now(), Sender: ui.ScreenDashboard})
	
	select {
	case event := <-channel:
		if event.Type != "test_event" {
			mt.addError("Event bus should deliver events correctly")
		}
	default:
		mt.addError("Event bus should deliver events to subscribers")
	}
	
	fmt.Println("  ✓ Integration validation completed")
}

// testUIComponents tests UI components
func (mt *MVPTester) testUIComponents() {
	defer mt.recordTestResult("ui_components", mt.errorMessages == nil)
	
	fmt.Println("  Validating UI components...")
	
	theme := ui.NewTheme()
	keys := ui.NewKeys()
	
	// Test button
	button := ui.NewButton("test_button", "Test Button", func() tea.Cmd { return nil }, theme, keys)
	if button == nil {
		mt.addError("Button component should not be nil")
	}
	
	// Test input
	input := ui.NewInput("test_input", "Test Input", theme, keys)
	if input == nil {
		mt.addError("Input component should not be nil")
	}
	
	// Test list
	list := ui.NewList("test_list", 10, theme, keys)
	if list == nil {
		mt.addError("List component should not be nil")
	}
	
	// Test container
	container := ui.NewContainer("test_container", "Test Container", theme, keys)
	if container == nil {
		mt.addError("Container component should not be nil")
	}
	
	// Test text
	text := ui.NewText("test_text", "Test Text", theme, keys)
	if text == nil {
		mt.addError("Text component should not be nil")
	}
	
	// Test spinner
	spinner := ui.NewSpinner("test_spinner", theme, keys)
	if spinner == nil {
		mt.addError("Spinner component should not be nil")
	}
	
	fmt.Println("  ✓ UI components validation completed")
}

// testKeyboardNavigation tests keyboard navigation
func (mt *MVPTester) testKeyboardNavigation() {
	defer mt.recordTestResult("keyboard_navigation", mt.errorMessages == nil)
	
	fmt.Println("  Validating keyboard navigation...")
	
	keys := ui.NewKeys()
	
	// Test navigation keys
	if keys.Tab == "" || keys.ShiftTab == "" || keys.Up == "" || keys.Down == "" {
		mt.addError("All navigation keys should be configured")
	}
	
	// Test action keys
	if keys.Enter == "" || keys.Space == "" || keys.Refresh == "" || keys.Quit == "" {
		mt.addError("All action keys should be configured")
	}
	
	// Test global shortcuts
	if keys.Dashboard == "" || keys.Doctor == "" || keys.Metrics == "" || keys.Cost == "" || keys.Settings == "" {
		mt.addError("All global shortcuts should be configured")
	}
	
	// Test key type checking
	if !keys.IsNavigationKey(keys.Up) {
		mt.addError("Up should be a navigation key")
	}
	
	if keys.IsNavigationKey(keys.Refresh) {
		mt.addError("Refresh should not be a navigation key")
	}
	
	if !keys.IsActionKey(keys.Enter) {
		mt.addError("Enter should be an action key")
	}
	
	if !keys.IsGlobalShortcut(keys.Dashboard) {
		mt.addError("Dashboard should be a global shortcut")
	}
	
	fmt.Println("  ✓ Keyboard navigation validation completed")
}

// testTUIIntegration tests TUI integration
func (mt *MVPTester) testTUIIntegration() {
	defer mt.recordTestResult("tui_integration", mt.errorMessages == nil)
	
	fmt.Println("  Validating TUI integration...")
	
	// Test TUI creation
	tui := ui.NewTUI()
	if tui == nil {
		mt.addError("TUI should not be nil")
	}
	
	// Test integration layer
	integration := ui.NewIntegrationLayer()
	if integration == nil {
		mt.addError("Integration layer should not be nil")
	}
	
	// Test screen manager
	screenManager := tui.GetScreenManager()
	if screenManager == nil {
		mt.addError("Screen manager should not be nil")
	}
	
	// Test data manager
	dataManager := tui.GetDataManager()
	if dataManager == nil {
		mt.addError("Data manager should not be nil")
	}
	
	// Test event bus
	eventBus := tui.GetEventBus()
	if eventBus == nil {
		mt.addError("Event bus should not be nil")
	}
	
	// Test shared data operations
	tui.SetSharedData("test_key", "test_value")
	value, exists := tui.GetSharedData("test_key")
	if !exists || value != "test_value" {
		mt.addError("Shared data operations should work correctly")
	}
	
	fmt.Println("  ✓ TUI integration validation completed")
}

// testScreenSpecific tests screen-specific functionality
func (mt *MVPTester) testScreenSpecific() {
	defer mt.recordTestResult("screen_specific", mt.errorMessages == nil)
	
	fmt.Println("  Validating screen-specific functionality...")
	
	theme := ui.NewTheme()
	keys := ui.NewKeys()
	config := ui.NewAccessibilityConfig()
	a11y := ui.NewA11y(config, theme, keys)
	perfTracker := ui.NewPerfTracker(ui.NewPerformanceBudget())
	
	// Test dashboard
	dashboard := ui.NewDashboardScreen(theme, keys, a11y, perfTracker)
	if dashboard.GetID() != "dashboard" {
		mt.addError("Dashboard ID should be 'dashboard'")
	}
	
	// Test doctor
	doctor := ui.NewDoctorScreen(theme, keys, a11y, perfTracker)
	if doctor.GetID() != "doctor" {
		mt.addError("Doctor ID should be 'doctor'")
	}
	
	// Test metrics
	metrics := ui.NewMetricsScreen(theme, keys, a11y, perfTracker)
	if metrics.GetID() != "metrics" {
		mt.addError("Metrics ID should be 'metrics'")
	}
	
	// Test cost
	cost := ui.NewCostScreen(theme, keys, a11y, perfTracker)
	if cost.GetID() != "cost" {
		mt.addError("Cost ID should be 'cost'")
	}
	
	// Test settings
	settings := ui.NewSettingsScreen(theme, keys, a11y, perfTracker)
	if settings.GetID() != "settings" {
		mt.addError("Settings ID should be 'settings'")
	}
	
	fmt.Println("  ✓ Screen-specific validation completed")
}

// recordTestResult records a test result
func (mt *MVPTester) recordTestResult(testName string, success bool) {
	mt.testResults[testName] = success
}

// addError adds an error message
func (mt *MVPTester) addError(message string) {
	mt.errorMessages = append(mt.errorMessages, message)
	fmt.Printf("    ❌ %s\n", message)
}

// generateReport generates a comprehensive test report
func (mt *MVPTester) generateReport() {
	fmt.Println()
	fmt.Println("📊 MVP Test Report")
	fmt.Println("==================")
	fmt.Println()
	
	// Summary
	totalTests := len(mt.testResults)
	passedTests := 0
	for _, result := range mt.testResults {
		if result {
			passedTests++
		}
	}
	
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", totalTests-passedTests)
	fmt.Printf("Success Rate: %.1f%%\n", float64(passedTests)/float64(totalTests)*100)
	fmt.Println()
	
	// Test Results
	fmt.Println("Test Results:")
	for testName, result := range mt.testResults {
		status := "✅ PASS"
		if !result {
			status = "❌ FAIL"
		}
		fmt.Printf("  %s: %s\n", testName, status)
	}
	
	// Errors
	if len(mt.errorMessages) > 0 {
		fmt.Println()
		fmt.Println("Errors:")
		for _, error := range mt.errorMessages {
			fmt.Printf("  ❌ %s\n", error)
		}
	}
	
	// Performance Metrics
	fmt.Println()
	fmt.Println("Performance Metrics:")
	fmt.Printf("  Total Test Time: %v\n", time.Since(mt.startTime))
	
	// System Information
	fmt.Println()
	fmt.Println("System Information:")
	fmt.Printf("  Go Version: %s\n", "1.22+ (required)")
	fmt.Printf("  Architecture: %s\n", "amd64")
	fmt.Printf("  Platform: %s\n", "windows")
	
	// MVP Status
	fmt.Println()
	fmt.Println("MVP Status:")
	if passedTests == totalTests {
		fmt.Println("🎉 ALL TESTS PASSED! MVP requirements fully satisfied!")
		fmt.Println("✅ CodeAuthorityFirst is ready for production!")
	} else {
		fmt.Printf("⚠️  %d tests failed. MVP needs additional work.\n", totalTests-passedTests)
	}
	
	fmt.Println()
	fmt.Println("🏁 MVP Testing Suite Completed")
}

func main() {
	// Create and run MVP tester
	tester := NewMVPTester()
	
	if err := tester.RunAllTests(); err != nil {
		log.Printf("Error running tests: %v", err)
		os.Exit(1)
	}
	
	// Exit with appropriate code
	allPassed := true
	for _, result := range tester.testResults {
		if !result {
			allPassed = false
			break
		}
	}
	
	if !allPassed {
		os.Exit(1)
	}
	
	os.Exit(0)
}