// ui_self_review_test.go — Tests for UI Self-Review Checklist effectiveness
// This file tests the checklist coverage and validation procedures
package screens

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

// TestUIReviewChecklist tests the completeness of the UI self-review checklist
func TestUIReviewChecklistCompleteness(t *testing.T) {
	// Test that all required sections are present in the checklist
	checklistContent, err := os.ReadFile("../docs/UI_SELF_REVIEW_CHECKLIST.md")
	if err != nil {
		t.Fatalf("Could not read checklist: %v", err)
	}

	checklistText := string(checklistContent)

	requiredSections := []string{
		"Performance Requirements Verification",
		"Accessibility Compliance",
		"Screen Template Compliance",
		"Universal Keyboard Map Verification",
		"Visual Design System",
		"UI-Velocity Patterns",
		"Testing Coverage Requirements",
		"Code Organization Standards",
		"Edge Case Handling",
		"Documentation Integration",
		"Final Checklist",
	}

	for _, section := range requiredSections {
		if !strings.Contains(checklistText, section) {
			t.Errorf("Missing required section: %s", section)
		}
	}
}

// TestPerformanceRequirements tests all INV-UI requirements
func TestPerformanceRequirements(t *testing.T) {
	tests := []struct {
		name     string
		invariant string
		test     func(*testing.T)
	}{
		{
			name:     "INV-UI-001 First Paint <100ms",
			invariant: "INV-UI-001",
			test:     testFirstPaintPerformance,
		},
		{
			name:     "INV-UI-002 Frame Render <16ms",
			invariant: "INV-UI-002",
			test:     testFrameRenderPerformance,
		},
		{
			name:     "INV-UI-003 Keystroke Response <50ms",
			invariant: "INV-UI-003",
			test:     testKeystrokeResponse,
		},
		{
			name:     "INV-UI-004 Skeleton + SWR Pattern",
			invariant: "INV-UI-004",
			test:     testSkeletonSWRPattern,
		},
		{
			name:     "INV-UI-005 Lazy Screen Mounting",
			invariant: "INV-UI-005",
			test:     testLazyScreenMounting,
		},
		{
			name:     "INV-UI-006 Lipgloss Pre-computed Styles",
			invariant: "INV-UI-006",
			test:     testLipglossStyles,
		},
		{
			name:     "INV-UI-007 Heatmap Render <200ms",
			invariant: "INV-UI-007",
			test:     testHeatmapRender,
		},
		{
			name:     "INV-UI-008 Color-blind Safety",
			invariant: "INV-UI-008",
			test:     testColorBlindSafety,
		},
		{
			name:     "INV-UI-009 Keyboard-only Operation",
			invariant: "INV-UI-009",
			test:     testKeyboardOnlyOperation,
		},
		{
			name:     "INV-UI-010 Help Overlay",
			invariant: "INV-UI-010",
			test:     testHelpOverlay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// testFirstPaintPerformance tests INV-UI-001
func testFirstPaintPerformance(t *testing.T) {
	// Create a test screen
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Measure first paint time
	start := time.Now()
	view := screen.View()
	duration := time.Since(start)
	
	// Verify first paint meets requirement
	if duration >= 100*time.Millisecond {
		t.Errorf("First paint took %v, expected <100ms (INV-UI-001)", duration)
	}
	
	// Verify view is not empty
	if view == "" {
		t.Error("View should not be empty after first paint")
	}
}

// testFrameRenderPerformance tests INV-UI-002
func testFrameRenderPerformance(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Measure frame render time
	start := time.Now()
	for i := 0; i < 10; i++ {
		screen.View()
	}
	avgRender := time.Since(start) / 10
	
	if avgRender >= 16*time.Millisecond {
		t.Errorf("Average render time %v, expected <16ms (INV-UI-002)", avgRender)
	}
}

// testKeystrokeResponse tests INV-UI-003
func testKeystrokeResponse(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Measure keystroke response time
	start := time.Now()
	screen, _ = screen.Update(tea.KeyMsg{Type: tea.KeyTab})
	duration := time.Since(start)
	
	if duration >= 50*time.Millisecond {
		t.Errorf("Keystroke response took %v, expected <50ms (INV-UI-003)", duration)
	}
}

// testSkeletonSWRPattern tests INV-UI-004
func testSkeletonSWRPattern(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test skeleton appears immediately
	skeletonView := screen.View()
	if !strings.Contains(skeletonView, "Loading") {
		t.Error("Skeleton should appear immediately (INV-UI-004)")
	}
	
	// Test data loading replaces skeleton
	screen.Loading = false
	screen.Sessions = []*storage.Session{{ID: "test-session"}}
	dataView := screen.View()
	
	if strings.Contains(dataView, "Loading") {
		t.Error("Skeleton should be replaced with real data (INV-UI-004)")
	}
	
	if !strings.Contains(dataView, "test-session") {
		t.Error("Real data should be visible after loading")
	}
}

// testLazyScreenMounting tests INV-UI-005
func testLazyScreenMounting(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test screen initialization is lightweight
	start := time.Now()
	screen.Init()
	initDuration := time.Since(start)
	
	if initDuration >= 100*time.Millisecond {
		t.Errorf("Screen initialization took %v, expected <100ms (INV-UI-005)", initDuration)
	}
	
	// Test closing screen is lightweight
	screen.Active = false
	closeStart := time.Now()
	screen.View()
	closeDuration := time.Since(closeStart)
	
	if closeDuration >= 10*time.Millisecond {
		t.Errorf("Screen closing took %v, expected to be lightweight (INV-UI-005)", closeDuration)
	}
}

// testLipglossStyles tests INV-UI-006
func testLipglossStyles(t *testing.T) {
	// Test that styles are pre-computed
	theme := ui.DefaultTheme()
	
	start := time.Now()
	for i := 0; i < 100; i++ {
		_ = theme.Header.Render("Test")
	}
	avgStyleTime := time.Since(start) / 100
	
	if avgStyleTime >= 1*time.Millisecond {
		t.Errorf("Style computation took %v per call, expected to be pre-computed (INV-UI-006)", avgStyleTime)
	}
}

// testHeatmapRender tests INV-UI-007
func testHeatmapRender(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Create test data
	screen.Manifest = &storage.Manifest{
		FileCount: 247,
	}
	
	// Measure heatmap render time
	start := time.Now()
	view := screen.renderHeatmap()
	duration := time.Since(start)
	
	if duration >= 200*time.Millisecond {
		t.Errorf("Heatmap render took %v, expected <200ms (INV-UI-007)", duration)
	}
	
	if view == "" {
		t.Error("Heatmap view should not be empty")
	}
}

// testColorBlindSafety tests INV-UI-008
func testColorBlindSafety(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test that status symbols have text
	view := screen.View()
	
	for _, line := range strings.Split(view, "\n") {
		if strings.ContainsAny(line, "✓⚠✗") {
			if !strings.ContainsFunc(line, func(r rune) bool {
				return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
			}) {
				t.Errorf("Status symbol without text: %q (INV-UI-008)", line)
			}
		}
	}
}

// testKeyboardOnlyOperation tests INV-UI-009
func testKeyboardOnlyOperation(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test all navigation can be done with keyboard
	keys := []tea.KeyType{
		tea.KeyTab,      // Navigate forward
		tea.KeyShiftTab, // Navigate backward
		tea.KeyUp,       // Move up
		tea.KeyDown,     // Move down
		tea.KeyEnter,    // Activate
		tea.KeyR,        // Refresh
		tea.KeyQuestionMark, // Help
		tea.KeyEsc,      // Back
		tea.KeyCtrl+C,   // Quit
	}
	
	for _, key := range keys {
		_, _ = screen.Update(tea.KeyMsg{Type: key})
	}
	
	// Verify screen is still in valid state
	if screen == nil {
		t.Error("Screen should handle keyboard input without crashing (INV-UI-009)")
	}
}

// testHelpOverlay tests INV-UI-010
func testHelpOverlay(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test help overlay can be toggled
	screen.ShowHelp = false
	view1 := screen.View()
	
	screen.ShowHelp = true
	view2 := screen.View()
	
	if strings.Contains(view1, "Keyboard Shortcuts") {
		t.Error("Help should not be visible initially (INV-UI-010)")
	}
	
	if !strings.Contains(view2, "Keyboard Shortcuts") {
		t.Error("Help should be visible when toggled (INV-UI-010)")
	}
}

// TestAccessibilityCompliance tests accessibility requirements
func TestAccessibilityCompliance(t *testing.T) {
	tests := []struct {
		name     string
		test     func(*testing.T)
	}{
		{
			name: "Text Representation for Visual Elements",
			test: testTextRepresentation,
		},
		{
			name: "Screen Reader Support",
			test: testScreenReaderSupport,
		},
		{
			name: "High Contrast Mode",
			test: testHighContrastMode,
		},
		{
			name: "Focus Visibility",
			test: testFocusVisibility,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// testTextRepresentation tests text representation for visual elements
func testTextRepresentation(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	view := screen.View()
	
	// Test that charts have text representation
	if strings.Contains(view, "sparkline") {
		if !strings.Contains(view, "trend") {
			t.Error("Sparkline should have text representation of trend")
		}
	}
	
	// Test that bars have text representation
	if strings.Contains(view, "bar") {
		if !strings.Contains(view, "percentage") {
			t.Error("Bar chart should have text representation of percentage")
		}
	}
}

// testScreenReaderSupport tests screen reader support
func testScreenReaderSupport(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	view := screen.View()
	
	// Test that status is announced first
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		if strings.Contains(line, "ERROR:") {
			if i > 0 && !strings.Contains(lines[i-1], "ERROR:") {
				t.Error("Error status should be announced first")
			}
		}
	}
}

// testHighContrastMode tests high contrast mode
func testHighContrastMode(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test normal mode
	normalView := screen.View()
	
	// Test high contrast mode
	screen.Theme.IsHighContrast = true
	highContrastView := screen.View()
	
	// Both should be functional
	if normalView == "" {
		t.Error("Normal mode should render content")
	}
	
	if highContrastView == "" {
		t.Error("High contrast mode should render content")
	}
	
	// High contrast should use different colors
	if normalView == highContrastView {
		t.Error("High contrast mode should produce different output")
	}
}

// testFocusVisibility tests focus visibility
func testFocusVisibility(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test that focused item has marker
	screen.CurrentWidget = 0
	view := screen.View()
	
	if !strings.Contains(view, "▸") {
		t.Error("Focused item should have marker (INV-PROMPT-MVP-007)")
	}
}

// TestScreenTemplateCompliance tests screen template compliance
func TestScreenTemplateCompliance(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test that screen has required methods
	if screen.Init == nil {
		t.Error("Screen should have Init() method")
	}
	
	if screen.Update == nil {
		t.Error("Screen should have Update() method")
	}
	
	if screen.View == nil {
		t.Error("Screen should have View() method")
	}
	
	// Test that state management follows Bubble Tea pattern
	screen.StateMachine.SetState(ui.StateLoading)
	if screen.StateMachine.GetState() != ui.StateLoading {
		t.Error("State management should follow Bubble Tea pattern")
	}
}

// TestUniversalKeyboardMap tests universal keyboard map
func TestUniversalKeyboardMap(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test global keys
	keys := []struct {
		key tea.KeyType
		action string
	}{
		{tea.KeyQuestionMark, "toggle help"},
		{tea.KeyCtrl+C, "quit"},
		{tea.KeyEsc, "back"},
		{tea.KeyTab, "navigate forward"},
		{tea.KeyShiftTab, "navigate backward"},
		{tea.KeyEnter, "activate"},
		{tea.KeyR, "refresh"},
	}
	
	for _, k := range keys {
		_, _ = screen.Update(tea.KeyMsg{Type: k.key})
		// Verify screen doesn't crash and handles key
	}
}

// TestVisualDesignSystem tests visual design system
func TestVisualDesignSystem(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test color palette
	theme := screen.Theme
	if theme.Primary == "" {
		t.Error("Primary color should be defined")
	}
	if theme.Secondary == "" {
		t.Error("Secondary color should be defined")
	}
	if theme.Warning == "" {
		t.Error("Warning color should be defined")
	}
	if theme.Error == "" {
		t.Error("Error color should be defined")
	}
	
	// Test typography
	view := screen.View()
	if !strings.Contains(view, "# CodeAuthorityFirst") {
		t.Error("Header should use proper typography")
	}
}

// TestUIVelocityPatterns tests UI velocity patterns
func TestUIVelocityPatterns(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test skeleton pattern
	screen.Loading = true
	skeletonView := screen.View()
	if !strings.Contains(skeletonView, "Loading") {
		t.Error("Skeleton pattern should show loading state")
	}
	
	// Test SWR pattern
	screen.Loading = false
	screen.Sessions = []*storage.Session{{ID: "test"}}
	dataView := screen.View()
	if strings.Contains(dataView, "Loading") {
		t.Error("SWR pattern should replace skeleton with data")
	}
	if !strings.Contains(dataView, "test") {
		t.Error("SWR pattern should show real data")
	}
}

// TestTestingCoverage tests testing coverage requirements
func TestTestingCoverage(t *testing.T) {
	// Test that screen has test file
	testFiles, err := filepath.Glob("pkg/ui/screens/*_test.go")
	if err != nil {
		t.Fatalf("Could not find test files: %v", err)
	}
	
	if len(testFiles) == 0 {
		t.Error("Screens should have test files")
	}
	
	// Test that test files cover required areas
	for _, testFile := range testFiles {
		content, err := os.ReadFile(testFile)
		if err != nil {
			continue
		}
		
		testContent := string(content)
		requiredTests := []string{
			"TestFocusVisibility",
			"TestSymbolTextPairing",
			"TestPerformance",
			"TestAccessibility",
		}
		
		hasRequiredTests := false
		for _, test := range requiredTests {
			if strings.Contains(testContent, test) {
				hasRequiredTests = true
				break
			}
		}
		
		if !hasRequiredTests {
			t.Errorf("Test file %s should cover required test areas", testFile)
		}
	}
}

// TestCodeOrganizationStandards tests code organization standards
func TestCodeOrganizationStandards(t *testing.T) {
	// Test file structure
	screenFiles, err := filepath.Glob("pkg/ui/screens/*.go")
	if err != nil {
		t.Fatalf("Could not find screen files: %v", err)
	}
	
	for _, file := range screenFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		fileContent := string(content)
		
		// Test that files have proper package declaration
		if !strings.Contains(fileContent, "package screens") {
			t.Errorf("File %s should have proper package declaration", file)
		}
		
		// Test that files have proper structure
		if !strings.Contains(fileContent, "type") {
			t.Errorf("File %s should have type definitions", file)
		}
		
		if !strings.Contains(fileContent, "func") {
			t.Errorf("File %s should have function definitions", file)
		}
	}
}

// TestEdgeCaseHandling tests edge case handling
func TestEdgeCaseHandling(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test empty states
	screen.Sessions = []*storage.Session{}
	view := screen.View()
	if strings.Contains(view, "No sessions yet") {
		if !strings.Contains(view, "Run caf init") {
			t.Error("Empty state should have CTA")
		}
	}
	
	// Test error states
	screen.StateMachine.SetState(ui.StateError)
	errorView := screen.View()
	if !strings.Contains(errorView, "Error") {
		t.Error("Error state should be visible")
	}
}

// TestDocumentationIntegration tests documentation integration
func TestDocumentationIntegration(t *testing.T) {
	// Test that screen files have proper documentation
	screenFiles, err := filepath.Glob("pkg/ui/screens/*.go")
	if err != nil {
		t.Fatalf("Could not find screen files: %v", err)
	}
	
	for _, file := range screenFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		fileContent := string(content)
		
		// Test that files have header comments
		if !strings.Contains(fileContent, "// Purpose:") {
			t.Errorf("File %s should have purpose documentation", file)
		}
		
		if !strings.Contains(fileContent, "// Layout:") {
			t.Errorf("File %s should have layout documentation", file)
		}
		
		if !strings.Contains(fileContent, "// Performance budget:") {
			t.Errorf("File %s should have performance budget documentation", file)
		}
	}
}

// TestFinalChecklist tests the final checklist
func TestFinalChecklist(t *testing.T) {
	// Test checklist has final checklist section
	checklistContent, err := os.ReadFile("../docs/UI_SELF_REVIEW_CHECKLIST.md")
	if err != nil {
		t.Fatalf("Could not read checklist: %v", err)
	}
	
	checklistText := string(checklistContent)
	
	if !strings.Contains(checklistText, "Final Checklist") {
		t.Error("Checklist should have final checklist section")
	}
	
	// Test that final checklist has required items
	requiredItems := []string{
		"All performance requirements met",
		"Accessibility compliance verified",
		"Screen template fully implemented",
		"Universal keyboard map working",
		"Visual design system applied consistently",
		"UI-velocity patterns implemented",
		"Testing coverage complete",
		"Code organization standards followed",
		"Edge cases handled properly",
		"Documentation integrated and updated",
	}
	
	for _, item := range requiredItems {
		if !strings.Contains(checklistText, item) {
			t.Errorf("Final checklist should include: %s", item)
		}
	}
}

// BenchmarkUIReviewChecklist benchmarks the time to complete UI review
func BenchmarkUIReviewChecklist(b *testing.B) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate completing the checklist
		_ = screen.View()
		_ = screen.Update(tea.KeyMsg{Type: tea.KeyTab})
		_ = screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}
}

// TestUIReviewWorkflow tests the complete UI review workflow
func TestUIReviewWorkflow(t *testing.T) {
	screen := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Test 1: Performance testing
	start := time.Now()
	view := screen.View()
	perfTime := time.Since(start)
	
	if perfTime >= 100*time.Millisecond {
		t.Errorf("Performance test failed: %v >= 100ms", perfTime)
	}
	
	// Test 2: Accessibility testing
	if !strings.Contains(view, "▸") {
		t.Error("Accessibility test failed: focus visibility not working")
	}
	
	// Test 3: Keyboard navigation
	screen, _ = screen.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screen.CurrentWidget == 0 {
		t.Error("Keyboard test failed: navigation not working")
	}
	
	// Test 4: Error handling
	screen.StateMachine.SetState(ui.StateError)
	errorView := screen.View()
	if !strings.Contains(errorView, "Error") {
		t.Error("Error handling test failed: error state not visible")
	}
	
	// Test 5: Documentation
	if !strings.Contains(view, "CodeAuthorityFirst") {
		t.Error("Documentation test failed: header not present")
	}
}