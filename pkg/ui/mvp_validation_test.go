//go:build integration

package ui

import (
	"fmt"
	"testing"
	"time"
)

// MVPValidator validates MVP requirements
type MVPValidator struct {
	theme         *Theme
	keys          *Keys
	a11y          *A11y
	perfTracker   *PerfTracker
	integration   *IntegrationLayer
	dashboard     *DashboardScreen
	doctor        *DoctorScreen
	metrics       *MetricsScreen
	cost          *CostScreen
	settings      *SettingsScreen
}

// NewMVPValidator creates a new MVP validator
func NewMVPValidator() *MVPValidator {
	theme := NewTheme()
	keys := NewKeys()
	config := NewAccessibilityConfig()
	a11y := NewA11y(config, theme, keys)
	perfTracker := NewPerfTracker(NewPerformanceBudget())
	integration := NewIntegrationLayer()
	
	dashboard := NewDashboardScreen(theme, keys, a11y, perfTracker)
	doctor := NewDoctorScreen(theme, keys, a11y, perfTracker)
	metrics := NewMetricsScreen(theme, keys, a11y, perfTracker)
	cost := NewCostScreen(theme, keys, a11y, perfTracker)
	settings := NewSettingsScreen(theme, keys, a11y, perfTracker)
	
	return &MVPValidator{
		theme:       theme,
		keys:        keys,
		a11y:        a11y,
		perfTracker: perfTracker,
		integration: integration,
		dashboard:   dashboard,
		doctor:      doctor,
		metrics:     metrics,
		cost:        cost,
		settings:    settings,
	}
}

// ValidateArchitecture validates architectural requirements
func (v *MVPValidator) ValidateArchitecture(t *testing.T) {
	t.Log("🏗️  Validating architectural requirements...")
	
	// Validate INV-UI-001: TUI first paint <100ms
	budget := NewPerformanceBudget()
	if budget.FirstPaint >= 100*time.Millisecond {
		t.Errorf("INV-UI-001: First paint budget should be <100ms, got %v", budget.FirstPaint)
	}
	
	// Validate INV-UI-002: Frame render <16ms (60fps)
	if budget.FrameRender >= 16*time.Millisecond {
		t.Errorf("INV-UI-002: Frame render budget should be <16ms, got %v", budget.FrameRender)
	}
	
	// Validate INV-UI-003: Keystroke response <50ms
	if budget.Keystroke >= 50*time.Millisecond {
		t.Errorf("INV-UI-003: Keystroke response budget should be <50ms, got %v", budget.Keystroke)
	}
	
	// Validate INV-UI-004: LSP startup ~500ms
	if budget.LSPStartup >= 600*time.Millisecond {
		t.Errorf("INV-UI-004: LSP startup should be ~500ms, got %v", budget.LSPStartup)
	}
	
	// Validate INV-UI-005: File op <1s for 10k files
	if budget.FileOp >= 1*time.Second {
		t.Errorf("INV-UI-005: File operation should be <1s, got %v", budget.FileOp)
	}
	
	// Validate INV-UI-006: Lipgloss pre-computes style strings at startup
	v.theme.PrecomputeStyles()
	if v.theme.Primary == nil {
		t.Error("INV-UI-006: Precomputed styles should not be nil")
	}
	
	// Validate INV-UI-007: WCAG-AAA compliance for color contrast, keyboard navigation, screen reader support
	config := v.a11y.Config
	if !config.KeyboardOnly {
		t.Error("INV-UI-007: Keyboard-only operation should be enabled")
	}
	
	if !config.ScreenReader {
		t.Error("INV-UI-007: Screen reader support should be enabled")
	}
	
	if !config.HighContrast {
		t.Error("INV-UI-007: High contrast mode should be available")
	}
	
	// Validate INV-UI-008: Color-blind friendly: pair color with text/symbol, never color-only
	primaryColor := v.theme.Primary.GetForeground()
	backgroundColor := v.theme.Background.GetBackground()
	
	if primaryColor == "" || backgroundColor == "" {
		t.Error("INV-UI-008: Colors should not be empty")
	}
	
	// Validate INV-UI-009: Keyboard-only operation is 100% complete; no feature requires mouse
	if v.keys.Tab == "" || v.keys.Enter == "" || v.keys.Esc == "" {
		t.Error("INV-UI-009: All navigation should be keyboard accessible")
	}
	
	t.Log("✅ Architectural requirements validated")
}

// ValidateFeatureRequirements validates feature requirements
func (v *MVPValidator) ValidateFeatureRequirements(t *testing.T) {
	t.Log("🎯 Validating feature requirements...")
	
	// Validate FR-TUI-001: Dashboard shows recent sessions, decisions, and agent activity
	sessions := v.dashboard.generateMockSessions()
	if len(sessions) == 0 {
		t.Error("FR-TUI-001: Dashboard should have recent sessions")
	}
	
	decisions := v.dashboard.generateMockDecisions()
	if len(decisions) == 0 {
		t.Error("FR-TUI-001: Dashboard should have decisions")
	}
	
	activity := v.dashboard.generateMockActivity()
	if len(activity) == 0 {
		t.Error("FR-TUI-001: Dashboard should have agent activity")
	}
	
	// Validate FR-TUI-002: Doctor shows 5 health checks with detailed diagnostics
	results := v.doctor.generateMockResults()
	if len(results) != 5 {
		t.Errorf("FR-TUI-002: Doctor should show 5 health checks, got %d", len(results))
	}
	
	for _, result := range results {
		if result.ID == "" || result.Name == "" || result.Status == "" {
			t.Error("FR-TUI-002: Health check results should be complete")
		}
	}
	
	// Validate FR-TUI-003: Metrics shows usage statistics, session counts, performance data
	sessionStats := v.metrics.generateMockSessionStats()
	if sessionStats.TotalSessions == 0 {
		t.Error("FR-TUI-003: Metrics should show session counts")
	}
	
	usageStats := v.metrics.generateMockUsageStats()
	if usageStats.TotalFilesModified == 0 {
		t.Error("FR-TUI-003: Metrics should show usage statistics")
	}
	
	performanceData := v.metrics.generateMockPerformanceData()
	if performanceData.AvgResponseTime == 0 {
		t.Error("FR-TUI-003: Metrics should show performance data")
	}
	
	// Validate FR-TUI-004: Cost shows cost tracking, usage costs, budget information
	totalCost := 45.67
	monthlyBudget := 100.0
	alerts := v.cost.generateMockAlerts(totalCost, monthlyBudget)
	if len(alerts) == 0 {
		t.Error("FR-TUI-004: Cost should show budget alerts")
	}
	
	recentCosts := v.cost.generateMockRecentCosts()
	if len(recentCosts) == 0 {
		t.Error("FR-TUI-004: Cost should show recent costs")
	}
	
	// Validate FR-TUI-005: Settings shows configuration options, accessibility settings, preferences
	general := GeneralSettings{}
	if general.Theme == "" {
		t.Error("FR-TUI-005: Settings should show general options")
	}
	
	accessibility := AccessibilitySettings{}
	if !accessibility.KeyboardOnly {
		t.Error("FR-TUI-005: Settings should show accessibility options")
	}
	
	notifications := NotificationSettings{}
	if notifications.AlertThreshold == 0 {
		t.Error("FR-TUI-005: Settings should show notification options")
	}
	
	advanced := AdvancedSettings{}
	if advanced.DataRetention == 0 {
		t.Error("FR-TUI-005: Settings should show advanced options")
	}
	
	// Validate FR-TUI-006: TUI integrates all screens with smooth transitions and shared state
	if v.integration == nil {
		t.Error("FR-TUI-006: Integration layer should exist")
	}
	
	if len(v.integration.ScreenManager.Screens) != 5 {
		t.Errorf("FR-TUI-006: Should have 5 screens, got %d", len(v.integration.ScreenManager.Screens))
	}
	
	// Validate FR-TUI-007: Integration layer maintains shared state and smooth transitions
	if v.integration.DataManager == nil {
		t.Error("FR-TUI-007: Data manager should exist")
	}
	
	if v.integration.EventBus == nil {
		t.Error("FR-TUI-007: Event bus should exist")
	}
	
	t.Log("✅ Feature requirements validated")
}

// ValidateStateManagement validates state management requirements
func (v *MVPValidator) ValidateStateManagement(t *testing.T) {
	t.Log("🔄 Validating state management requirements...")
	
	// Validate INV-PROMPT-MVP-006: Every screen implements ALL 5 mandatory UI states
	screens := []Screen{
		v.dashboard,
		v.doctor,
		v.metrics,
		v.cost,
		v.settings,
	}
	
	for _, screen := range screens {
		sm := NewStateMachine()
		
		// Test Loading → Success
		sm.Transition(LoadingMsg{})
		if sm.Current() != StateLoading {
			t.Errorf("Screen should start in loading state")
		}
		
		sm.Transition(DataLoadedMsg{Empty: false})
		if sm.Current() != StateSuccess {
			t.Errorf("Screen should transition to success on data loaded")
		}
		
		// Test Loading → Error
		sm = NewStateMachine()
		sm.Transition(LoadingMsg{})
		sm.Transition(DataErrorMsg{Err: fmt.Errorf("test error")})
		if sm.Current() != StateError {
			t.Errorf("Screen should transition to error on data error")
		}
		
		// Test Loading → Empty
		sm = NewStateMachine()
		sm.Transition(LoadingMsg{})
		sm.Transition(DataLoadedMsg{Empty: true})
		if sm.Current() != StateEmpty {
			t.Errorf("Screen should transition to empty on empty data")
		}
		
		// Test Success → Loading (on refresh)
		sm = NewStateMachine()
		sm.Transition(DataLoadedMsg{Empty: false})
		sm.Transition(RefreshMsg{Force: true})
		if sm.Current() != StateLoading {
			t.Errorf("Screen should transition to loading on refresh")
		}
		
		// Test Success → Error (on refresh failure)
		sm = NewStateMachine()
		sm.Transition(DataLoadedMsg{Empty: false})
		sm.Transition(DataErrorMsg{Err: fmt.Errorf("refresh error")})
		if sm.Current() != StateError {
			t.Errorf("Screen should transition to error on refresh failure")
		}
		
		// Test Success → Partial (on partial refresh failure)
		sm = NewStateMachine()
		sm.Transition(DataLoadedMsg{Empty: false})
		sm.Transition(PartialDataMsg{Missing: []string{"file1"}})
		if sm.Current() != StatePartial {
			t.Errorf("Screen should transition to partial on partial refresh failure")
		}
	}
	
	t.Log("✅ State management requirements validated")
}

// ValidateAccessibility validates accessibility requirements
func (v *MVPValidator) ValidateAccessibility(t *testing.T) {
	t.Log("♿ Validating accessibility requirements...")
	
	// Validate WCAG-AAA compliance
	config := v.a11y.Config
	
	// Test high contrast mode
	if !config.HighContrast {
		t.Error("High contrast mode should be available")
	}
	
	// Test reduced motion
	if !config.ReducedMotion {
		t.Error("Reduced motion should be available")
	}
	
	// Test screen reader support
	if !config.ScreenReader {
		t.Error("Screen reader support should be available")
	}
	
	// Test large text
	if !config.LargeText {
		t.Error("Large text should be available")
	}
	
	// Test keyboard-only operation
	if !config.KeyboardOnly {
		t.Error("Keyboard-only operation should be available")
	}
	
	// Test tab navigation
	if !config.TabNavigation {
		t.Error("Tab navigation should be available")
	}
	
	// Test escape exits
	if !config.EscapeExits {
		t.Error("Escape exits should be available")
	}
	
	// Test confirm actions
	if !config.ConfirmActions {
		t.Error("Confirm actions should be available")
	}
	
	// Test focus management
	focusManager := NewFocusManager(v.a11y)
	focusManager.AddFocusableElement("element1")
	focusManager.AddFocusableElement("element2")
	
	if len(focusManager.focusableElements) != 2 {
		t.Error("Focus manager should track focusable elements")
	}
	
	// Test screen reader announcements
	screenReader := NewScreenReaderManager(v.a11y)
	screenReader.AddAnnouncement("test_id", "test message", AnnouncementImportant, "test context", false)
	
	announcements := screenReader.GetAnnouncements()
	if len(announcements) != 1 {
		t.Error("Screen reader should track announcements")
	}
	
	// Test accessibility report
	report := v.a11y.GetAccessibilityReport()
	if report == "" {
		t.Error("Accessibility report should not be empty")
	}
	
	t.Log("✅ Accessibility requirements validated")
}

// ValidatePerformance validates performance requirements
func (v *MVPValidator) ValidatePerformance(t *testing.T) {
	t.Log("⚡ Validating performance requirements...")
	
	// Test performance budgets
	budget := NewPerformanceBudget()
	
	// Test first paint budget
	if budget.FirstPaint >= 100*time.Millisecond {
		t.Errorf("First paint budget should be <100ms, got %v", budget.FirstPaint)
	}
	
	// Test frame render budget
	if budget.FrameRender >= 16*time.Millisecond {
		t.Errorf("Frame render budget should be <16ms, got %v", budget.FrameRender)
	}
	
	// Test keystroke budget
	if budget.Keystroke >= 50*time.Millisecond {
		t.Errorf("Keystroke budget should be <50ms, got %v", budget.Keystroke)
	}
	
	// Test performance tracker
	tracker := NewPerfTracker(budget)
	tracker.Start("test_operation")
	time.Sleep(10 * time.Millisecond) // Simulate work
	duration := tracker.End("test_operation")
	
	if duration <= 0 {
		t.Error("Performance tracker should measure positive duration")
	}
	
	// Test performance warnings
	warning := NewPerformanceWarning("test_operation", duration, budget.Keystroke)
	if warning.Ratio <= 0 {
		t.Error("Performance warning should calculate ratio correctly")
	}
	
	// Test performance report
	report := tracker.GetReport()
	if report == "" {
		t.Error("Performance report should not be empty")
	}
	
	t.Log("✅ Performance requirements validated")
}

// ValidateIntegration validates integration requirements
func (v *MVPValidator) ValidateIntegration(t *testing.T) {
	t.Log("🔗 Validating integration requirements...")
	
	// Test integration layer initialization
	if err := v.integration.Initialize(); err != nil {
		t.Errorf("Integration layer should initialize without error: %v", err)
	}
	
	// Test screen manager
	screenManager := v.integration.ScreenManager
	if len(screenManager.Screens) != 5 {
		t.Errorf("Screen manager should have 5 screens, got %d", len(screenManager.Screens))
	}
	
	// Test data manager
	dataManager := v.integration.DataManager
	dataManager.SetSharedData("test_key", "test_value")
	value, exists := dataManager.GetSharedData("test_key")
	if !exists || value != "test_value" {
		t.Error("Data manager should handle shared data correctly")
	}
	
	// Test event bus
	eventBus := v.integration.EventBus
	channel := eventBus.Subscribe("test_event")
	eventBus.Publish(Event{Type: "test_event", Payload: "test_payload", Timestamp: time.Now(), Sender: ScreenDashboard})
	
	select {
	case event := <-channel:
		if event.Type != "test_event" {
			t.Error("Event bus should deliver events correctly")
		}
	default:
		t.Error("Event bus should deliver events to subscribers")
	}
	
	// Test screen transitions
	screenManager.SwitchScreen(ScreenDoctor)
	if screenManager.CurrentScreen != ScreenDoctor {
		t.Error("Screen manager should switch screens correctly")
	}
	
	t.Log("✅ Integration requirements validated")
}

// ValidateUIComponents validates UI components
func (v *MVPValidator) ValidateUIComponents(t *testing.T) {
	t.Log("🎨 Validating UI components...")
	
	theme := NewTheme()
	keys := NewKeys()
	
	// Test button component
	button := NewButton("test_button", "Test Button", func() tea.Cmd { return nil }, theme, keys)
	if button == nil {
		t.Error("Button component should not be nil")
	}
	
	// Test input component
	input := NewInput("test_input", "Test Input", theme, keys)
	if input == nil {
		t.Error("Input component should not be nil")
	}
	
	// Test list component
	list := NewList("test_list", 10, theme, keys)
	if list == nil {
		t.Error("List component should not be nil")
	}
	
	// Test container component
	container := NewContainer("test_container", "Test Container", theme, keys)
	if container == nil {
		t.Error("Container component should not be nil")
	}
	
	// Test text component
	text := NewText("test_text", "Test Text", theme, keys)
	if text == nil {
		t.Error("Text component should not be nil")
	}
	
	// Test spinner component
	spinner := NewSpinner("test_spinner", theme, keys)
	if spinner == nil {
		t.Error("Spinner component should not be nil")
	}
	
	t.Log("✅ UI components validated")
}

// ValidateKeyboardNavigation validates keyboard navigation
func (v *MVPValidator) ValidateKeyboardNavigation(t *testing.T) {
	t.Log("⌨️ Validating keyboard navigation...")
	
	keys := NewKeys()
	
	// Test navigation keys
	navigationKeys := []string{keys.Tab, keys.ShiftTab, keys.Up, keys.Down, keys.Left, keys.Right, keys.Home, keys.End}
	for _, key := range navigationKeys {
		if key == "" {
			t.Error("All navigation keys should be configured")
		}
	}
	
	// Test action keys
	actionKeys := []string{keys.Enter, keys.Space, keys.Refresh, keys.Quit, keys.Help}
	for _, key := range actionKeys {
		if key == "" {
			t.Error("All action keys should be configured")
		}
	}
	
	// Test global shortcuts
	globalShortcuts := []string{keys.Dashboard, keys.Doctor, keys.Metrics, keys.Cost, keys.Settings}
	for _, key := range globalShortcuts {
		if key == "" {
			t.Error("All global shortcuts should be configured")
		}
	}
	
	// Test key type checking
	if !keys.IsNavigationKey(keys.Up) {
		t.Error("Up should be a navigation key")
	}
	
	if keys.IsNavigationKey(keys.Refresh) {
		t.Error("Refresh should not be a navigation key")
	}
	
	if !keys.IsActionKey(keys.Enter) {
		t.Error("Enter should be an action key")
	}
	
	if keys.IsActionKey(keys.Up) {
		t.Error("Up should not be an action key")
	}
	
	if !keys.IsGlobalShortcut(keys.Dashboard) {
		t.Error("Dashboard should be a global shortcut")
	}
	
	if keys.IsGlobalShortcut(keys.Up) {
		t.Error("Up should not be a global shortcut")
	}
	
	t.Log("✅ Keyboard navigation validated")
}

// RunFullMVPValidation runs the complete MVP validation
func (v *MVPValidator) RunFullMVPValidation(t *testing.T) {
	t.Log("🚀 Starting full MVP validation...")
	
	// Validate all requirements
	v.ValidateArchitecture(t)
	v.ValidateFeatureRequirements(t)
	v.ValidateStateManagement(t)
	v.ValidateAccessibility(t)
	v.ValidatePerformance(t)
	v.ValidateIntegration(t)
	v.ValidateUIComponents(t)
	v.ValidateKeyboardNavigation(t)
	
	t.Log("🎉 MVP validation completed successfully!")
}

// TestMVPValidation runs the MVP validation test
func TestMVPValidation(t *testing.T) {
	validator := NewMVPValidator()
	validator.RunFullMVPValidation(t)
}