// high_contrast_integration_test.go - Comprehensive high-contrast mode integration testing
package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui/screens"
)

// TestHighContrastIntegration tests integration between different screens in high-contrast mode
func TestHighContrastIntegration(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test screen models with high-contrast mode
	screens := []struct {
		name     string
		model    func() interface{}
		testFunc func(interface{}) string
	}{
		{
			name:  "Dashboard",
			model: func() interface{} { return screens.NewDashboardModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) string {
				if model, ok := m.(*screens.DashboardModel); ok {
					return model.GetTodayWidget().Render()
				}
				return ""
			},
		},
		{
			name:  "Adapters Manager",
			model: func() interface{} { return screens.NewAdaptersModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) string {
				if model, ok := m.(*screens.AdaptersModel); ok {
					return model.GetCapabilityMatrix().Render()
				}
				return ""
			},
		},
		{
			name:  "Onboarding Wizard",
			model: func() interface{} { return screens.NewOnboardingModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) string {
				if model, ok := m.(*screens.OnboardingModel); ok {
					return model.GetStepContent().Render()
				}
				return ""
			},
		},
	}

	for _, screen := range screens {
		t.Run(screen.name, func(t *testing.T) {
			model := screen.model()
			content := screen.testFunc(model)
			
			if content == "" {
				t.Errorf("%s should not be empty in high-contrast mode", screen.name)
			}

			// Test high-contrast colors are present
			if !strings.Contains(content, "#FFFFFF") && !strings.Contains(content, "#000000") {
				t.Errorf("%s should use high-contrast colors", screen.name)
			}

			// Test focus indicators are present
			if !HasFocusMarker(content) {
				t.Errorf("%s should have focus indicators in high-contrast mode", screen.name)
			}

			// Test status prefixes are present
			if strings.Contains(content, "healthy") && !strings.Contains(content, "+") {
				t.Errorf("Healthy status in %s should have + prefix", screen.name)
			}
			if strings.Contains(content, "error") && !strings.Contains(content, "X") {
				t.Errorf("Error status in %s should have X prefix", screen.name)
			}
			if strings.Contains(content, "warning") && !strings.Contains(content, "!") {
				t.Errorf("Warning status in %s should have ! prefix", screen.name)
			}
		})
	}
}

// TestHighContrastScreenTransitions tests screen transitions in high-contrast mode
func TestHighContrastScreenTransitions(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test transitions between different screens
	transitions := []struct {
		from    string
		to      string
		test    func() string
	}{
		{
			from: "dashboard",
			to:   "adapters",
			test: func() string {
				// Simulate dashboard to adapters transition
				return "Adapters Manager"
			},
		},
		{
			from: "adapters",
			to:   "onboarding",
			test: func() string {
				// Simulate adapters to onboarding transition
				return "Onboarding Wizard"
			},
		},
		{
			from: "onboarding",
			to:   "dashboard",
			test: func() string {
				// Simulate onboarding to dashboard transition
				return "Dashboard"
			},
		},
	}

	for _, transition := range transitions {
		t.Run(transition.from+"->"+transition.to, func(t *testing.T) {
			content := transition.test()
			
			// Test high-contrast colors are present
			if !strings.Contains(content, "#FFFFFF") && !strings.Contains(content, "#000000") {
				t.Errorf("Screen transition should use high-contrast colors")
			}

			// Test focus indicators are present
			if !HasFocusMarker(content) {
				t.Errorf("Screen transition should have focus indicators")
			}

			// Test status prefixes are present
			if strings.Contains(content, "healthy") && !strings.Contains(content, "+") {
				t.Errorf("Healthy status should have + prefix")
			}
		})
	}
}

// TestHighContrastFocusManagement tests focus management across screens in high-contrast mode
func TestHighContrastFocusManagement(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.KeyboardOnly = true
	a.Config.TabNavigation = true

	// Test focus management across different screens
	screens := []string{"dashboard", "adapters", "onboarding", "settings"}

	// Test global focus manager
	globalFM := NewFocusManager(a)
	for _, screen := range screens {
		globalFM.AddFocusableElement(screen)
	}

	// Test focus movement across screens
	for i := 0; i < len(screens); i++ {
		globalFM.MoveForward()
		expectedScreen := screens[(i+1)%len(screens)]
		if globalFM.GetFocusedElement() != expectedScreen {
			t.Errorf("Focus should be on %s, got %s", expectedScreen, globalFM.GetFocusedElement())
		}
	}

	// Test focus path
	focusPath := globalFM.GetFocusPath()
	if len(focusPath) != len(screens) {
		t.Errorf("Focus path should have %d elements, got %d", len(screens), len(focusPath))
	}

	// Test screen-specific focus managers
	screenFM := NewFocusManager(a)
	screenFM.AddFocusableElement("widget1")
	screenFM.AddFocusableElement("widget2")
	screenFM.AddFocusableElement("widget3")

	// Test screen focus movement
	screenFM.MoveForward()
	if screenFM.GetFocusedElement() != "widget2" {
		t.Errorf("Screen focus should move to widget2")
	}

	screenFM.MoveForward()
	if screenFM.GetFocusedElement() != "widget3" {
		t.Errorf("Screen focus should move to widget3")
	}
}

// TestHighContrastKeyboardNavigationAcrossScreens tests keyboard navigation across screens in high-contrast mode
func TestHighContrastKeyboardNavigationAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.KeyboardOnly = true
	a.Config.TabNavigation = true

	// Test keyboard navigation across different screens
	screens := []struct {
		name   string
		widgets []string
	}{
		{
			name:   "Dashboard",
			widgets: []string{"today", "score", "agents", "decisions", "cross-repo"},
		},
		{
			name:   "Adapters Manager",
			widgets: []string{"matrix", "health", "installation", "details"},
		},
		{
			name:   "Onboarding Wizard",
			widgets: []string{"welcome", "detection", "installation", "summary"},
		},
	}

	for _, screen := range screens {
		t.Run(screen.name, func(t *testing.T) {
			// Test focus manager for screen
			fm := NewFocusManager(a)
			for _, widget := range screen.widgets {
				fm.AddFocusableElement(widget)
			}

			// Test focus movement within screen
			for i := 0; i < len(screen.widgets); i++ {
				fm.MoveForward()
				expectedWidget := screen.widgets[(i+1)%len(screen.widgets)]
				if fm.GetFocusedElement() != expectedWidget {
					t.Errorf("Focus in %s should be on %s, got %s", screen.name, expectedWidget, fm.GetFocusedElement())
				}
			}

			// Test tab order
			tabOrder := a.GetTabOrder()
			if len(tabOrder) != len(screen.widgets) {
				t.Errorf("Tab order for %s should have %d elements, got %d", screen.name, len(screen.widgets), len(tabOrder))
			}

			// Test focus path
			focusPath := fm.GetFocusPath()
			if len(focusPath) != len(screen.widgets) {
				t.Errorf("Focus path for %s should have %d elements, got %d", screen.name, len(screen.widgets), len(focusPath))
			}
		})
	}
}

// TestHighContrastErrorHandlingAcrossScreens tests error handling across screens in high-contrast mode
func TestHighContrastErrorHandlingAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.ScreenReader = true

	// Test error handling across different screens
	screens := []struct {
		name   string
		errors []string
	}{
		{
			name:   "Dashboard",
			errors: []string{"data loading failed", "connection timeout", "invalid input"},
		},
		{
			name:   "Adapters Manager",
			errors: []string{"installation failed", "health check failed", "configuration error"},
		},
		{
			name:   "Onboarding Wizard",
			errors: []string{"detection failed", "setup incomplete", "validation error"},
		},
	}

	for _, screen := range screens {
		t.Run(screen.name, func(t *testing.T) {
			// Test error rendering for each error type
			for _, errorText := range screen.errors {
				// Test error status with prefix
				errorStatus := a.StatusPair("error", errorText)
				if !strings.Contains(errorStatus, "X") || !strings.Contains(errorStatus, "✗") {
					t.Errorf("Error status for %s should contain X prefix and ✗ symbol", screen.name)
				}

				// Test error announcement
				announcement := a.GetAnnouncementText("CRITICAL", errorText)
				if announcement != "CRITICAL: "+errorText {
					t.Errorf("Error announcement should be status-first: %s", announcement)
				}

				// Test error box rendering
				errorBox := a.Theme.ErrorBox.Render(errorText)
				if !strings.Contains(errorBox, "#FFFF00") || !strings.Contains(errorBox, "#000000") {
					t.Errorf("Error box should use yellow on black in high-contrast mode")
				}
			}
		})
	}
}

// TestHighContrastScreenReaderSupportAcrossScreens tests screen reader support across screens in high-contrast mode
func TestHighContrastScreenReaderSupportAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.ScreenReader = true

	// Test screen reader support across different screens
	screens := []struct {
		name     string
		elements []struct {
			visual   string
			altText  string
		}
	}{
		{
			name: "Dashboard",
			elements: []struct {
				visual   string
				altText  string
			}{
				{"📊 Today", "Today widget: summary of due items"},
				{"📈 Score", "Authorship score: metrics and analysis"},
				{"🤖 Agents", "Active agents: detected AI agents status"},
				{"📝 Decisions", "Recent decisions: coding decisions history"},
				{"🔗 Cross-repo", "Cross-repo: multi-repository sync status"},
			},
		},
		{
			name: "Adapters Manager",
			elements: []struct {
				visual   string
				altText  string
			}{
				{"🔌 Matrix", "Capability matrix: adapter capabilities overview"},
				{"💚 Health", "Health status: adapter health indicators"},
				{"📦 Installation", "Installation progress: adapter setup status"},
				{"🔍 Details", "Adapter details: configuration and settings"},
			},
		},
		{
			name: "Onboarding Wizard",
			elements: []struct {
				visual   string
				altText  string
			}{
				{"👋 Welcome", "Welcome screen: introduction to setup"},
				{"🔍 Detection", "Detection screen: AI agent detection"},
				{"📦 Installation", "Installation screen: adapter installation"},
				{"✅ Summary", "Summary screen: setup completion overview"},
			},
		},
	}

	for _, screen := range screens {
		t.Run(screen.name, func(t *testing.T) {
			// Test screen reader support for each element
			for _, element := range screen.elements {
				// Test visual mode
				a.SetOutputProfile(termenv.TrueColor)
				visualResult := a.RenderAlt(element.visual, element.altText)
				if visualResult != element.visual {
					t.Errorf("Visual mode should return visual content: %s", visualResult)
				}

				// Test ASCII/screen reader mode
				a.SetOutputProfile(termenv.Ascii)
				asciiResult := a.RenderAlt(element.visual, element.altText)
				if asciiResult != element.altText {
					t.Errorf("ASCII mode should return alt-text: %s", asciiResult)
				}

				// Test announcement
				announcement := a.GetAnnouncementText("INFO", element.altText)
				if announcement != "INFO: "+element.altText {
					t.Errorf("Announcement should be status-first: %s", announcement)
				}
			}
		})
	}
}

// TestHighContrastTerminalProfileAdaptationAcrossScreens tests terminal profile adaptation across screens in high-contrast mode
func TestHighContrastTerminalProfileAdaptationAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test terminal profile adaptation across different screens
	screens := []struct {
		name   string
		widgets []string
	}{
		{
			name:   "Dashboard",
			widgets: []string{"today", "score", "agents", "decisions", "cross-repo"},
		},
		{
			name:   "Adapters Manager",
			widgets: []string{"matrix", "health", "installation", "details"},
		},
		{
			name:   "Onboarding Wizard",
			widgets: []string{"welcome", "detection", "installation", "summary"},
		},
	}

	profiles := []termenv.Profile{termenv.TrueColor, termenv.ANSI256, termenv.ANSI, termenv.Ascii}

	for _, profile := range profiles {
		for _, screen := range screens {
			t.Run(string(profile)+"-"+screen.name, func(t *testing.T) {
				a.SetOutputProfile(profile)
				
				// Test focus border adaptation
				focusBorder := FocusBorderStyle(profile)
				if focusBorder == nil {
					t.Errorf("Focus border should not be nil for %s profile", profile)
				}

				// Test focus box rendering for each widget
				for _, widget := range screen.widgets {
					focusBox := RenderFocusBox(widget, true, a.Theme, profile)
					if profile != termenv.Ascii && !HasFocusMarker(focusBox) {
						t.Errorf("Focus box for %s should have focus marker in %s profile", widget, profile)
					}
				}

				// Test focus title rendering
				for _, widget := range screen.widgets {
					focusTitle := FocusTitle(widget, true)
					if !strings.Contains(focusTitle, "▸") {
						t.Errorf("Focus title for %s should contain marker: %s", widget, focusTitle)
					}
				}
			})
		}
	}
}

// TestHighContrastValidationAcrossScreens tests validation across screens in high-contrast mode
func TestHighContrastValidationAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test validation across different screens
	screens := []struct {
		name     string
		model    func() interface{}
		testFunc func(interface{}) *Theme
	}{
		{
			name:  "Dashboard",
			model: func() interface{} { return screens.NewDashboardModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) *Theme {
				if model, ok := m.(*screens.DashboardModel); ok {
					return a.GetTheme()
				}
				return nil
			},
		},
		{
			name:  "Adapters Manager",
			model: func() interface{} { return screens.NewAdaptersModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) *Theme {
				if model, ok := m.(*screens.AdaptersModel); ok {
					return a.GetTheme()
				}
				return nil
			},
		},
		{
			name:  "Onboarding Wizard",
			model: func() interface{} { return screens.NewOnboardingModel(t.TempDir(), "0.4.1-test") },
			testFunc: func(m interface{}) *Theme {
				if model, ok := m.(*screens.OnboardingModel); ok {
					return a.GetTheme()
				}
				return nil
			},
		},
	}

	for _, screen := range screens {
		t.Run(screen.name, func(t *testing.T) {
			model := screen.model()
			theme := screen.testFunc(model)
			
			if theme == nil {
				t.Errorf("%s theme should not be nil", screen.name)
			}

			// Test validation
			errs := a.ValidateAccessibility(theme)
			if len(errs) != 0 {
				t.Errorf("%s should pass validation in high-contrast mode, got errors: %v", screen.name, errs)
			}

			// Test validation with broken theme
			brokenTheme := theme.Clone()
			brokenTheme.FocusBorder = lipgloss.NewStyle() // Remove focus border
			errs = a.ValidateAccessibility(brokenTheme)
			found := false
			for _, e := range errs {
				if strings.Contains(e, "focus indicators") {
					found = true
				}
			}
			if !found {
				t.Errorf("Validation should detect missing focus indicators in %s", screen.name)
			}
		})
	}
}

// TestHighContrastReportAcrossScreens tests accessibility report across screens in high-contrast mode
func TestHighContrastReportAcrossScreens(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.HighContrast = true
	a.Config.ReducedMotion = true
	a.Config.ScreenReader = true
	a.Config.KeyboardOnly = true

	// Test accessibility report across different screens
	screens := []string{"Dashboard", "Adapters Manager", "Onboarding Wizard"}

	for _, screen := range screens {
		t.Run(screen, func(t *testing.T) {
			// Test screen-specific accessibility report
			report := a.GetAccessibilityReport()
			if !strings.Contains(report, "High contrast mode enabled") {
				t.Errorf("Report for %s should mention high-contrast mode", screen)
			}
			if !strings.Contains(report, "WCAG AAA compliance") {
				t.Errorf("Report for %s should mention WCAG AAA compliance", screen)
			}
			if !strings.Contains(report, "Enhanced focus indicators") {
				t.Errorf("Report for %s should mention enhanced focus indicators", screen)
			}
			if !strings.Contains(report, "Status prefixes (+, !, X)") {
				t.Errorf("Report for %s should mention status prefixes", screen)
			}
			if !strings.Contains(report, "Screen reader support enabled") {
				t.Errorf("Report for %s should mention screen reader support", screen)
			}
			if !strings.Contains(report, "Keyboard-only operation enabled") {
				t.Errorf("Report for %s should mention keyboard-only operation", screen)
			}
		})
	}
}