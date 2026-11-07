// high_contrast_screen_test.go — comprehensive high-contrast mode testing for all screens
package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui/screens"
)

// newTestHighContrastA11y creates a new accessibility manager with high-contrast mode enabled
func newTestHighContrastA11y() *A11y {
	config := NewAccessibilityConfig()
	config.HighContrast = true
	config.ScreenReader = true
	config.KeyboardOnly = true
	config.TabNavigation = true
	config.AnnounceChanges = true

	theme := NewTheme()
	keys := NewKeys()

	return NewA11y(config, theme, keys)
}

// TestDashboardHighContrastMode tests the Dashboard screen in high-contrast mode
func TestDashboardHighContrastMode(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	
	// Test dashboard model creation with high-contrast
	model := screens.NewDashboardModel(repoRoot, "0.4.1-test")
	if model == nil {
		t.Fatal("Dashboard model should not be nil")
	}

	// Test high-contrast theme application
	hcTheme := a.GetTheme()
	if hcTheme == nil {
		t.Fatal("High-contrast theme should not be nil")
	}

	// Test widget rendering in high-contrast mode
	widgets := []struct {
		name string
		test func() string
	}{
		{"Today Widget", func() string { return model.GetTodayWidget().Render() }},
		{"Authorship Score Widget", func() string { return model.GetAuthorshipWidget().Render() }},
		{"Active Agents Widget", func() string { return model.GetAgentsWidget().Render() }},
		{"Recent Decisions Widget", func() string { return model.GetDecisionsWidget().Render() }},
		{"Cross-repo Widget", func() string { return model.GetCrossRepoWidget().Render() }},
	}

	for _, widget := range widgets {
		widgetContent := widget.test()
		if widgetContent == "" {
			t.Errorf("%s should not be empty in high-contrast mode", widget.name)
		}

		// Test high-contrast colors are present
		if !strings.Contains(widgetContent, "#FFFFFF") && !strings.Contains(widgetContent, "#000000") {
			t.Errorf("%s should use high-contrast colors", widget.name)
		}

		// Test focus indicators are present
		if !HasFocusMarker(widgetContent) {
			t.Errorf("%s should have focus indicators in high-contrast mode", widget.name)
		}

		// Test status prefixes are present
		if strings.Contains(widgetContent, "healthy") && !strings.Contains(widgetContent, "+") {
			t.Errorf("Healthy status in %s should have + prefix in high-contrast mode", widget.name)
		}
		if strings.Contains(widgetContent, "error") && !strings.Contains(widgetContent, "X") {
			t.Errorf("Error status in %s should have X prefix in high-contrast mode", widget.name)
		}
		if strings.Contains(widgetContent, "warning") && !strings.Contains(widgetContent, "!") {
			t.Errorf("Warning status in %s should have ! prefix in high-contrast mode", widget.name)
		}
	}
}

// TestDashboardHighContrastFocusNavigation tests focus navigation in dashboard high-contrast mode
func TestDashboardHighContrastFocusNavigation(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	model := screens.NewDashboardModel(repoRoot, "0.4.1-test")

	// Test focus manager setup
	fm := NewFocusManager(a)
	fm.AddFocusableElement("today")
	fm.AddFocusableElement("score")
	fm.AddFocusableElement("agents")
	fm.AddFocusableElement("decisions")
	fm.AddFocusableElement("cross-repo")

	// Test focus movement
	focusPath := fm.GetFocusPath()
	if len(focusPath) != 5 {
		t.Errorf("Focus path should have 5 elements, got %d", len(focusPath))
	}

	// Test focus cycling
	currentFocus := fm.GetFocusedElement()
	if currentFocus != "" {
		t.Errorf("Initial focus should be empty")
	}

	// Move focus and verify
	fm.MoveForward()
	if fm.GetFocusedElement() != "score" {
		t.Errorf("First forward should land on score")
	}

	fm.MoveForward()
	if fm.GetFocusedElement() != "agents" {
		t.Errorf("Second forward should land on agents")
	}

	// Test focus wrapping
	fm.MoveForward()
	fm.MoveForward()
	fm.MoveForward() // Should wrap back to today
	if fm.GetFocusedElement() != "today" {
		t.Errorf("Focus should wrap back to today")
	}
}

// TestDashboardHighContrastStatusIndicators tests status indicators in dashboard high-contrast mode
func TestDashboardHighContrastStatusIndicators(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	model := screens.NewDashboardModel(repoRoot, "0.4.1-test")

	// Test status indicators across different terminal profiles
	profiles := []termenv.Profile{termenv.TrueColor, termenv.ANSI256, termenv.ANSI, termenv.Ascii}
	
	for _, profile := range profiles {
		a.SetOutputProfile(profile)
		
		// Test healthy status with prefix
		healthyStatus := a.StatusPair("healthy", "service running")
		if !strings.Contains(healthyStatus, "+") || !strings.Contains(healthyStatus, "✓") {
			t.Errorf("Healthy status should contain + prefix and ✓ symbol in profile %v", profile)
		}

		// Test error status with prefix
		errorStatus := a.StatusPair("error", "service failed")
		if !strings.Contains(errorStatus, "X") || !strings.Contains(errorStatus, "✗") {
			t.Errorf("Error status should contain X prefix and ✗ symbol in profile %v", profile)
		}

		// Test warning status with prefix
		warningStatus := a.StatusPair("warning", "service degraded")
		if !strings.Contains(warningStatus, "!") || !strings.Contains(warningStatus, "⚠") {
			t.Errorf("Warning status should contain ! prefix and ⚠ symbol in profile %v", profile)
		}
	}
}

// TestAdaptersManagerHighContrastMode tests the Adapters Manager screen in high-contrast mode
func TestAdaptersManagerHighContrastMode(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	
	// Test adapters manager model creation with high-contrast
	model := screens.NewAdaptersModel(repoRoot, "0.4.1-test")
	if model == nil {
		t.Fatal("Adapters Manager model should not be nil")
	}

	// Test high-contrast theme application
	hcTheme := a.GetTheme()
	if hcTheme == nil {
		t.Fatal("High-contrast theme should not be nil")
	}

	// Test adapter matrix rendering in high-contrast mode
	matrixContent := model.GetCapabilityMatrix().Render()
	if matrixContent == "" {
		t.Fatal("Capability matrix should not be empty in high-contrast mode")
	}

	// Test high-contrast colors are present
	if !strings.Contains(matrixContent, "#FFFFFF") && !strings.Contains(matrixContent, "#000000") {
		t.Error("Capability matrix should use high-contrast colors")
	}

	// Test status indicators with prefixes
	if strings.Contains(matrixContent, "healthy") && !strings.Contains(matrixContent, "+") {
		t.Error("Healthy adapters should have + prefix in high-contrast mode")
	}
	if strings.Contains(matrixContent, "unhealthy") && !strings.Contains(matrixContent, "X") {
		t.Error("Unhealthy adapters should have X prefix in high-contrast mode")
	}
	if strings.Contains(matrixContent, "degraded") && !strings.Contains(matrixContent, "!") {
		t.Error("Degraded adapters should have ! prefix in high-contrast mode")
	}
}

// TestAdaptersManagerHighContrastHealthIndicators tests health indicators in adapters manager high-contrast mode
func TestAdaptersManagerHighContrastHealthIndicators(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	model := screens.NewAdaptersModel(repoRoot, "0.4.1-test")

	// Test health status rendering across different terminal profiles
	healthStatuses := []string{"healthy", "degraded", "unhealthy", "unknown"}
	
	for _, status := range healthStatuses {
		// Test status prefix
		prefix := a.StatusPrefix(status)
		var expectedPrefix string
		switch status {
		case "healthy":
			expectedPrefix = "+"
		case "degraded":
			expectedPrefix = "!"
		case "unhealthy":
			expectedPrefix = "X"
		default:
			expectedPrefix = "?"
		}
		
		if prefix != expectedPrefix {
			t.Errorf("Status prefix for %s should be %s, got %s", status, expectedPrefix, prefix)
		}

		// Test status pair with symbol
		statusPair := a.StatusPair(status, "adapter name")
		if !strings.Contains(statusPair, prefix) {
			t.Errorf("Status pair for %s should contain prefix %s", status, prefix)
		}

		// Test symbol is present
		symbol := StatusSymbol(status)
		if symbol == "" {
			t.Errorf("Symbol should not be empty for status %s", status)
		}
	}
}

// TestOnboardingWizardHighContrastMode tests the Onboarding Wizard screen in high-contrast mode
func TestOnboardingWizardHighContrastMode(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	
	// Test onboarding wizard model creation with high-contrast
	model := screens.NewOnboardingModel(repoRoot, "0.4.1-test")
	if model == nil {
		t.Fatal("Onboarding Wizard model should not be nil")
	}

	// Test high-contrast theme application
	hcTheme := a.GetTheme()
	if hcTheme == nil {
		t.Fatal("High-contrast theme should not be nil")
	}

	// Test wizard step rendering in high-contrast mode
	steps := []string{"welcome", "detection", "installation", "summary"}
	
	for _, step := range steps {
		model.SetCurrentStep(step)
		stepContent := model.GetStepContent().Render()
		
		if stepContent == "" {
			t.Errorf("Step %s should not be empty in high-contrast mode", step)
		}

		// Test high-contrast colors are present
		if !strings.Contains(stepContent, "#FFFFFF") && !strings.Contains(stepContent, "#000000") {
			t.Errorf("Step %s should use high-contrast colors", step)
		}

		// Test focus indicators are present
		if !HasFocusMarker(stepContent) {
			t.Errorf("Step %s should have focus indicators in high-contrast mode", step)
		}
	}
}

// TestOnboardingWizardHighContrastProgressIndicators tests progress indicators in onboarding wizard high-contrast mode
func TestOnboardingWizardHighContrastProgressIndicators(t *testing.T) {
	a := newTestHighContrastA11y()
	repoRoot := t.TempDir()
	model := screens.NewOnboardingModel(repoRoot, "0.4.1-test")

	// Test progress rendering in high-contrast mode
	progressContent := model.GetProgressIndicator().Render()
	if progressContent == "" {
		t.Fatal("Progress indicator should not be empty in high-contrast mode")
	}

	// Test high-contrast colors are present
	if !strings.Contains(progressContent, "#FFFFFF") && !strings.Contains(progressContent, "#000000") {
		t.Error("Progress indicator should use high-contrast colors")
	}

	// Test alt-text generation for progress
	altText := ProgressAltText(50, 500, 1000, "steps", "total steps")
	if altText != "50% of total steps (500 / 1,000 steps)" {
		t.Errorf("Progress alt-text should be formatted correctly: %s", altText)
	}

	// Test progress announcement
	announcement := a.GetAnnouncementText("INFO", altText)
	if announcement != "INFO: 50% of total steps (500 / 1,000 steps)" {
		t.Errorf("Progress announcement should be status-first: %s", announcement)
	}
}

// TestCommonComponentsHighContrastMode tests common UI components in high-contrast mode
func TestCommonComponentsHighContrastMode(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	components := []struct {
		name string
		test func() lipgloss.Style
	}{
		{"Header", func() lipgloss.Style { return hcTheme.Header }},
		{"StatusBar", func() lipgloss.Style { return hcTheme.StatusBar }},
		{"Button", func() lipgloss.Style { return hcTheme.Button }},
		{"Input", func() lipgloss.Style { return hcTheme.Input }},
		{"ErrorBox", func() lipgloss.Style { return hcTheme.ErrorBox }},
		{"EmptyBox", func() lipgloss.Style { return hcTheme.EmptyBox }},
		{"Help", func() lipgloss.Style { return hcTheme.Help }},
	}

	for _, component := range components {
		style := component.test()
		rendered := style.Render("Test")
		
		if rendered == "" {
			t.Errorf("%s should not be empty in high-contrast mode", component.name)
		}

		// Test high-contrast colors are present
		hasHighContrast := strings.Contains(rendered, "#FFFFFF") || strings.Contains(rendered, "#000000") ||
			strings.Contains(rendered, "#FFFF00") || strings.Contains(rendered, "#CCCCCC")
		
		if !hasHighContrast {
			t.Errorf("%s should use high-contrast colors: %s", component.name, rendered)
		}
	}
}

// TestHighContrastErrorHandling tests error handling in high-contrast mode
func TestHighContrastErrorHandling(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	// Test error box rendering
	errorBox := hcTheme.ErrorBox.Render("Error: Connection failed")
	if !strings.Contains(errorBox, "#FFFF00") || !strings.Contains(errorBox, "#000000") {
		t.Error("Error box should use yellow on black in high-contrast mode")
	}

	// Test error status with prefix
	errorStatus := a.StatusPair("error", "critical failure")
	if !strings.Contains(errorStatus, "X") || !strings.Contains(errorStatus, "✗") {
		t.Error("Error status should contain X prefix and ✗ symbol")
	}

	// Test error announcement
	announcement := a.GetAnnouncementText("CRITICAL", "system error")
	if announcement != "CRITICAL: system error" {
		t.Errorf("Error announcement should be status-first: %s", announcement)
	}
}

// TestHighContrastScreenReaderSupport tests screen reader support in high-contrast mode
func TestHighContrastScreenReaderSupport(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.ScreenReader = true

	// Test alt-text rendering across terminal profiles
	profiles := []termenv.Profile{termenv.TrueColor, termenv.Ascii}
	
	for _, profile := range profiles {
		a.SetOutputProfile(profile)
		
		visual := "▁▂▃"
		altText := "30-day trend: rising"
		
		result := a.RenderAlt(visual, altText)
		
		if profile == termenv.TrueColor {
			if result != visual {
				t.Errorf("Visual mode should return visual content: %s", result)
			}
		} else {
			if result != altText {
				t.Errorf("ASCII mode should return alt-text: %s", result)
			}
		}
	}

	// Test screen reader manager
	srm := NewScreenReaderManager(a)
	srm.AddAnnouncement("test1", "dashboard loaded", AnnouncementImportant, "Dashboard", false)
	srm.Announce(AnnouncementWarning, "high-contrast mode enabled", "Dashboard", true)

	announcements := srm.GetAnnouncements()
	if len(announcements) != 2 {
		t.Errorf("Should have 2 announcements, got %d", len(announcements))
	}

	latest := srm.GetLatestAnnouncement()
	if latest == nil || latest.Message != "high-contrast mode enabled" {
		t.Errorf("Latest announcement should be about high-contrast mode")
	}
}

// TestHighContrastKeyboardNavigation tests keyboard navigation in high-contrast mode
func TestHighContrastKeyboardNavigation(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.KeyboardOnly = true
	a.Config.TabNavigation = true

	// Test focus manager
	fm := NewFocusManager(a)
	fm.AddFocusableElement("dashboard")
	fm.AddFocusableElement("adapters")
	fm.AddFocusableElement("onboarding")

	// Test tab order
	tabOrder := a.GetTabOrder()
	if len(tabOrder) != 3 {
		t.Errorf("Tab order should have 3 elements, got %d", len(tabOrder))
	}

	// Test focus movement
	if fm.MoveForward() != "adapters" {
		t.Errorf("First forward should land on adapters")
	}
	if fm.MoveForward() != "onboarding" {
		t.Errorf("Second forward should land on onboarding")
	}
	if fm.MoveForward() != "dashboard" {
		t.Errorf("Forward should wrap to dashboard")
	}

	// Test backward movement
	if fm.MoveBackward() != "onboarding" {
		t.Errorf("First backward should land on onboarding")
	}
	if fm.MoveBackward() != "adapters" {
		t.Errorf("Second backward should land on adapters")
	}
	if fm.MoveBackward() != "dashboard" {
		t.Errorf("Backward should wrap to dashboard")
	}
}

// TestHighContrastTerminalProfileAdaptation tests terminal profile adaptation in high-contrast mode
func TestHighContrastTerminalProfileAdaptation(t *testing.T) {
	a := newTestHighContrastA11y()

	profiles := []struct {
		name   string
		profile termenv.Profile
		border lipgloss.Border
	}{
		{"TrueColor", termenv.TrueColor, lipgloss.DoubleBorder()},
		{"ANSI256", termenv.ANSI256, lipgloss.ThickBorder()},
		{"ANSI", termenv.ANSI, lipgloss.NormalBorder()},
		{"Ascii", termenv.Ascii, lipgloss.NormalBorder()},
	}

	for _, p := range profiles {
		a.SetOutputProfile(p.profile)
		
		// Test focus border adaptation
		focusBorder := FocusBorderStyle(p.profile)
		if focusBorder != p.border {
			t.Errorf("Profile %s focus border mismatch: got %v, want %v", p.name, focusBorder, p.border)
		}

		// Test focus box rendering
		focusBox := RenderFocusBox("widget", true, a.Theme, p.profile)
		if !HasFocusMarker(focusBox) {
			t.Errorf("Profile %s focus box should have focus marker", p.name)
		}

		// Test focus title
		focusTitle := FocusTitle("Test Widget", true)
		if !strings.Contains(focusTitle, "▸") {
			t.Errorf("Profile %s focus title should contain marker: %s", p.name, focusTitle)
		}
	}
}

// TestHighContrastValidation tests validation in high-contrast mode
func TestHighContrastValidation(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test valid high-contrast configuration
	errs := a.ValidateAccessibility(a.Theme)
	if len(errs) != 0 {
		t.Errorf("High-contrast mode should pass validation, got errors: %v", errs)
	}

	// Test validation with missing focus indicators
	brokenTheme := a.Theme.Clone()
	brokenTheme.FocusBorder = lipgloss.NewStyle() // Remove focus border
	errs = a.ValidateAccessibility(brokenTheme)
	found := false
	for _, e := range errs {
		if strings.Contains(e, "focus indicators") {
			found = true
		}
	}
	if !found {
		t.Errorf("Validation should detect missing focus indicators in high-contrast mode")
	}

	// Test validation with missing status prefixes
	brokenConfig := NewAccessibilityConfig()
	brokenConfig.HighContrast = true
	brokenConfig.AnnounceChanges = false // This should cause an error
	brokenA11y := NewA11y(brokenConfig, a.Theme, a.Keys)
	
	errs = brokenA11y.ValidateAccessibility(brokenA11y.Theme)
	found = false
	for _, e := range errs {
		if strings.Contains(e, "announcements disabled") {
			found = true
		}
	}
	if !found {
		t.Errorf("Validation should detect disabled announcements in high-contrast mode")
	}
}

// TestHighContrastReport tests accessibility report in high-contrast mode
func TestHighContrastReport(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.HighContrast = true
	a.Config.ReducedMotion = true
	a.Config.ScreenReader = true
	a.Config.KeyboardOnly = true

	report := a.GetAccessibilityReport()
	if !strings.Contains(report, "High contrast mode enabled") {
		t.Errorf("Report should mention high-contrast mode")
	}
	if !strings.Contains(report, "WCAG AAA compliance") {
		t.Errorf("Report should mention WCAG AAA compliance")
	}
	if !strings.Contains(report, "Enhanced focus indicators") {
		t.Errorf("Report should mention enhanced focus indicators")
	}
	if !strings.Contains(report, "Status prefixes (+, !, X)") {
		t.Errorf("Report should mention status prefixes")
	}
	if !strings.Contains(report, "Screen reader support enabled") {
		t.Errorf("Report should mention screen reader support")
	}
	if !strings.Contains(report, "Keyboard-only operation enabled") {
		t.Errorf("Report should mention keyboard-only operation")
	}
}