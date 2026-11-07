// ui_components_high_contrast_test.go - Comprehensive UI component testing in high-contrast mode
package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestUIComponentsHighContrastMode tests all UI components in high-contrast mode
func TestUIComponentsHighContrastMode(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	components := []struct {
		name     string
		style    lipgloss.Style
		testText string
		expectedColors []string
	}{
		{
			name:     "Header",
			style:    hcTheme.Header,
			testText: "Dashboard",
			expectedColors: []string{"#FFFFFF", "#000000"},
		},
		{
			name:     "StatusBar",
			style:    hcTheme.StatusBar,
			testText: "Ready",
			expectedColors: []string{"#CCCCCC", "#000000"},
		},
		{
			name:     "Title",
			style:    hcTheme.Title,
			testText: "Welcome",
			expectedColors: []string{"#FFFFFF"},
		},
		{
			name:     "Border",
			style:    hcTheme.Border,
			testText: "Content",
			expectedColors: []string{"#CCCCCC"},
		},
		{
			name:     "FocusBorder",
			style:    hcTheme.FocusBorder,
			testText: "Focused Widget",
			expectedColors: []string{"#FFFFFF"},
		},
		{
			name:     "Button",
			style:    hcTheme.Button,
			testText: "Click Me",
			expectedColors: []string{"#FFFFFF", "#000000"},
		},
		{
			name:     "Input",
			style:    hcTheme.Input,
			testText: "text input",
			expectedColors: []string{"#000000", "#FFFFFF"},
		},
		{
			name:     "Help",
			style:    hcTheme.Help,
			testText: "Press ? for help",
			expectedColors: []string{"#CCCCCC", "#000000"},
		},
		{
			name:     "ErrorBox",
			style:    hcTheme.ErrorBox,
			testText: "Error: Connection failed",
			expectedColors: []string{"#FFFF00", "#000000"},
		},
		{
			name:     "EmptyBox",
			style:    hcTheme.EmptyBox,
			testText: "No data available",
			expectedColors: []string{"#CCCCCC", "#000000"},
		},
	}

	for _, component := range components {
		t.Run(component.name, func(t *testing.T) {
			rendered := component.style.Render(component.testText)
			
			if rendered == "" {
				t.Errorf("%s should not be empty in high-contrast mode", component.name)
			}

			// Test high-contrast colors are present
			hasExpectedColors := false
			for _, color := range component.expectedColors {
				if strings.Contains(rendered, color) {
					hasExpectedColors = true
					break
				}
			}
			
			if !hasExpectedColors {
				t.Errorf("%s should use expected high-contrast colors: %s", component.name, rendered)
			}

			// Test text is visible
			if !strings.Contains(rendered, component.testText) {
				t.Errorf("%s should contain test text: %s", component.name, component.testText)
			}
		})
	}
}

// TestUIComponentsHighContrastFocusIndicators tests focus indicators for all components
func TestUIComponentsHighContrastFocusIndicators(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	components := []struct {
		name     string
		style    lipgloss.Style
		testText string
	}{
		{"Header", hcTheme.Header, "Dashboard"},
		{"StatusBar", hcTheme.StatusBar, "Ready"},
		{"Title", hcTheme.Title, "Welcome"},
		{"Border", hcTheme.Border, "Content"},
		{"FocusBorder", hcTheme.FocusBorder, "Focused Widget"},
		{"Button", hcTheme.Button, "Click Me"},
		{"Input", hcTheme.Input, "text input"},
		{"Help", hcTheme.Help, "Press ? for help"},
		{"ErrorBox", hcTheme.ErrorBox, "Error: Connection failed"},
		{"EmptyBox", hcTheme.EmptyBox, "No data available"},
	}

	for _, component := range components {
		t.Run(component.name, func(t *testing.T) {
			// Test normal rendering (no focus)
			normalRender := component.style.Render(component.testText)
			if normalRender == "" {
				t.Errorf("%s normal render should not be empty", component.name)
			}

			// Test focus rendering (simulate focus by adding focus marker)
			focusedRender := FocusTitle(component.testText, true)
			if !HasFocusMarker(focusedRender) {
				t.Errorf("%s focused render should have focus marker", component.name)
			}

			// Test focus border rendering
			focusBox := RenderFocusBox(component.name, true, hcTheme, termenv.TrueColor)
			if !HasFocusMarker(focusBox) {
				t.Errorf("%s focus box should have focus marker", component.name)
			}
		})
	}
}

// TestUIComponentsHighContrastStatusIndicators tests status indicators for components
func TestUIComponentsHighContrastStatusIndicators(t *testing.T) {
	a := newTestHighContrastA11y()

	statusTypes := []string{"healthy", "warning", "error", "unknown"}

	for _, status := range statusTypes {
		t.Run(status, func(t *testing.T) {
			// Test status prefix
			prefix := a.StatusPrefix(status)
			var expectedPrefix string
			switch status {
			case "healthy":
				expectedPrefix = "+"
			case "warning":
				expectedPrefix = "!"
			case "error":
				expectedPrefix = "X"
			default:
				expectedPrefix = "?"
			}
			
			if prefix != expectedPrefix {
				t.Errorf("Status prefix for %s should be %s, got %s", status, expectedPrefix, prefix)
			}

			// Test status pair
			statusPair := a.StatusPair(status, "component status")
			if !strings.Contains(statusPair, prefix) {
				t.Errorf("Status pair for %s should contain prefix %s", status, prefix)
			}

			// Test status symbol
			symbol := StatusSymbol(status)
			if symbol == "" {
				t.Errorf("Symbol should not be empty for status %s", status)
			}
		})
	}
}

// TestUIComponentsHighContrastTerminalProfiles tests component rendering across terminal profiles
func TestUIComponentsHighContrastTerminalProfiles(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	profiles := []termenv.Profile{termenv.TrueColor, termenv.ANSI256, termenv.ANSI, termenv.Ascii}

	components := []struct {
		name  string
		style lipgloss.Style
		text  string
	}{
		{"Button", hcTheme.Button, "Click"},
		{"Input", hcTheme.Input, "text"},
		{"Header", hcTheme.Header, "Header"},
		{"ErrorBox", hcTheme.ErrorBox, "Error"},
	}

	for _, profile := range profiles {
		for _, component := range components {
			t.Run(string(profile)+"-"+component.name, func(t *testing.T) {
				a.SetOutputProfile(profile)
				
				// Test component rendering
				rendered := component.style.Render(component.text)
				if rendered == "" {
					t.Errorf("%s should not be empty in %s profile", component.name, profile)
				}

				// Test focus border adaptation
				focusBorder := FocusBorderStyle(profile)
				if focusBorder == nil {
					t.Errorf("Focus border should not be nil for %s profile", profile)
				}

				// Test focus box rendering
				focusBox := RenderFocusBox(component.name, true, hcTheme, profile)
				if profile != termenv.Ascii && !HasFocusMarker(focusBox) {
					t.Errorf("%s focus box should have focus marker in %s profile", component.name, profile)
				}
			})
		}
	}
}

// TestUIComponentsHighContrastErrorStates tests error states for components
func TestUIComponentsHighContrastErrorStates(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	errorComponents := []struct {
		name     string
		style    lipgloss.Style
		testText string
		expectedPrefix string
	}{
		{"ErrorBox", hcTheme.ErrorBox, "Critical error: System failure", "X"},
		{"Button", hcTheme.Button, "Error Action", "X"},
		{"Input", hcTheme.Input, "Error Input", "X"},
	}

	for _, component := range errorComponents {
		t.Run(component.name, func(t *testing.T) {
			// Test error rendering
			rendered := component.style.Render(component.testText)
			if rendered == "" {
				t.Errorf("%s should not be empty in error state", component.name)
			}

			// Test error status with prefix
			errorStatus := a.StatusPair("error", component.testText)
			if !strings.Contains(errorStatus, component.expectedPrefix) {
				t.Errorf("Error status should contain %s prefix", component.expectedPrefix)
			}
			if !strings.Contains(errorStatus, "✗") {
				t.Errorf("Error status should contain ✗ symbol")
			}

			// Test error announcement
			announcement := a.GetAnnouncementText("CRITICAL", component.testText)
			if announcement != "CRITICAL: "+component.testText {
				t.Errorf("Error announcement should be status-first: %s", announcement)
			}
		})
	}
}

// TestUIComponentsHighContrastScreenReaderSupport tests screen reader support for components
func TestUIComponentsHighContrastScreenReaderSupport(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.ScreenReader = true

	components := []struct {
		name     string
		visual   string
		altText  string
	}{
		{"Button", "🔘 Click", "Button: Click"},
		{"Input", "█ Text", "Input field: Text"},
		{"Status", "● Active", "Status: Active"},
		{"Progress", "▁▂▃▄", "Progress: 40% complete"},
	}

	for _, component := range components {
		t.Run(component.name, func(t *testing.T) {
			// Test visual mode
			a.SetOutputProfile(termenv.TrueColor)
			visualResult := a.RenderAlt(component.visual, component.altText)
			if visualResult != component.visual {
				t.Errorf("Visual mode should return visual content: %s", visualResult)
			}

			// Test ASCII/screen reader mode
			a.SetOutputProfile(termenv.Ascii)
			asciiResult := a.RenderAlt(component.visual, component.altText)
			if asciiResult != component.altText {
				t.Errorf("ASCII mode should return alt-text: %s", asciiResult)
			}

			// Test announcement
			announcement := a.GetAnnouncementText("INFO", component.altText)
			if announcement != "INFO: "+component.altText {
				t.Errorf("Announcement should be status-first: %s", announcement)
			}
		})
	}
}

// TestUIComponentsHighContrastKeyboardNavigation tests keyboard navigation for components
func TestUIComponentsHighContrastKeyboardNavigation(t *testing.T) {
	a := newTestHighContrastA11y()
	a.Config.KeyboardOnly = true
	a.Config.TabNavigation = true

	components := []string{"dashboard", "adapters", "onboarding", "settings"}

	// Test focus manager
	fm := NewFocusManager(a)
	for _, component := range components {
		fm.AddFocusableElement(component)
	}

	// Test tab order
	tabOrder := a.GetTabOrder()
	if len(tabOrder) != len(components) {
		t.Errorf("Tab order should have %d elements, got %d", len(components), len(tabOrder))
	}

	// Test focus movement
	for i := 0; i < len(components); i++ {
		fm.MoveForward()
		expectedFocus := components[(i+1)%len(components)]
		if fm.GetFocusedElement() != expectedFocus {
			t.Errorf("Focus should be on %s, got %s", expectedFocus, fm.GetFocusedElement())
		}
	}

	// Test backward movement
	for i := len(components) - 1; i >= 0; i-- {
		fm.MoveBackward()
		expectedFocus := components[i%len(components)]
		if fm.GetFocusedElement() != expectedFocus {
			t.Errorf("Focus should be on %s, got %s", expectedFocus, fm.GetFocusedElement())
		}
	}

	// Test focus path
	focusPath := fm.GetFocusPath()
	if len(focusPath) != len(components) {
		t.Errorf("Focus path should have %d elements, got %d", len(components), len(focusPath))
	}
}

// TestUIComponentsHighContrastValidation tests validation for components in high-contrast mode
func TestUIComponentsHighContrastValidation(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	// Test valid high-contrast configuration
	errs := a.ValidateAccessibility(hcTheme)
	if len(errs) != 0 {
		t.Errorf("High-contrast mode should pass validation, got errors: %v", errs)
	}

	// Test validation with broken focus indicators
	brokenTheme := hcTheme.Clone()
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

	// Test validation with broken status colors
	brokenTheme2 := hcTheme.Clone()
	brokenTheme2.Primary = lipgloss.NewStyle() // Remove primary color
	errs = a.ValidateAccessibility(brokenTheme2)
	found = false
	for _, e := range errs {
		if strings.Contains(e, "color contrast") {
			found = true
		}
	}
	if !found {
		t.Errorf("Validation should detect color contrast issues in high-contrast mode")
	}
}

// TestUIComponentsHighContrastReport tests accessibility report for components in high-contrast mode
func TestUIComponentsHighContrastReport(t *testing.T) {
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

// TestUIComponentsHighContrastContrastRatios tests color contrast ratios for components
func TestUIComponentsHighContrastContrastRatios(t *testing.T) {
	a := newTestHighContrastA11y()
	hcTheme := a.GetTheme()

	// Test primary color contrast (white on black)
	contrast := ContrastRatio("#000000", "#FFFFFF")
	if contrast < 7.0 || contrast > 21.0 {
		t.Errorf("Primary contrast ratio: %f (should be ~21:1 for WCAG AAA)", contrast)
	}

	// Test warning/error color contrast (yellow on black)
	contrast = ContrastRatio("#000000", "#FFFF00")
	if contrast < 7.0 || contrast > 21.0 {
		t.Errorf("Warning/Error contrast ratio: %f (should be ~21:1 for WCAG AAA)", contrast)
	}

	// Test muted color contrast (light gray on black)
	contrast = ContrastRatio("#000000", "#CCCCCC")
	if contrast < 7.0 || contrast > 21.0 {
		t.Errorf("Muted contrast ratio: %f (should be ~21:1 for WCAG AAA)", contrast)
	}

	// Test component specific colors
	buttonColors := []struct {
		bg    string
		fg    string
		name  string
	}{
		{"#FFFFFF", "#000000", "Button"},
		{"#000000", "#FFFFFF", "Input"},
		{"#FFFF00", "#000000", "ErrorBox"},
	}

	for _, color := range buttonColors {
		contrast := ContrastRatio(color.bg, color.fg)
		if contrast < 7.0 {
			t.Errorf("%s contrast ratio: %f (should be >=7:1 for WCAG AAA)", color.name, contrast)
		}
	}
}

// TestUIComponentsHighContrastAnimationTests tests reduced motion for components
func TestUIComponentsHighContrastAnimationTests(t *testing.T) {
	a := newTestHighContrastA11y()

	// Test normal animation
	if !a.ShouldAnimate() {
		t.Errorf("Normal mode should allow animation")
	}
	if a.TransitionFrames() != 8 {
		t.Errorf("Normal mode should have 8 transition frames")
	}

	frames := a.SpinnerFrames()
	if len(frames) <= 1 {
		t.Errorf("Normal mode should have animated spinner frames")
	}

	// Test reduced motion
	a.Config.ReducedMotion = true
	if a.ShouldAnimate() {
		t.Errorf("Reduced motion should not allow animation")
	}
	if a.TransitionFrames() != 0 {
		t.Errorf("Reduced motion should have 0 transition frames")
	}

	reducedFrames := a.SpinnerFrames()
	if len(reducedFrames) != 1 || reducedFrames[0] != "…" {
		t.Errorf("Reduced motion should have static spinner …")
	}
}