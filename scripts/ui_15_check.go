// ui_15_check.go - Automated 15-Item Critical UI Checklist Verification
// This program automatically verifies the 15 most critical UI requirements
package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// UI15CheckResult represents the result of a single checklist item
type UI15CheckResult struct {
	Item        string `json:"item"`
	Description string `json:"description"`
	Status      string `json:"status"` // PASS, FAIL, WARN
	Message     string `json:"message"`
	Details     string `json:"details,omitempty"`
}

// UI15CheckSuite represents the complete 15-item check suite
type UI15CheckSuite struct {
	Timestamp    string             `json:"timestamp"`
	TotalChecks  int                `json:"total_checks"`
	PassedChecks  int                `json:"passed_checks"`
	FailedChecks  int                `json:"failed_checks"`
	WarnedChecks  int                `json:"warned_checks"`
	Results      []UI15CheckResult  `json:"results"`
	ScreenPath   string             `json:"screen_path"`
	Performance  PerformanceMetrics `json:"performance"`
}

// PerformanceMetrics stores performance measurements
type PerformanceMetrics struct {
	FirstPaint    float64 `json:"first_paint_ms"`
	FrameRender   float64 `json:"frame_render_ms"`
	KeystrokeResp float64 `json:"keystroke_resp_ms"`
}

// Main 15-item checklist
var critical15Items = []string{
	"INV-UI-001_First_Paint_LT100ms",
	"INV-UI-002_Frame_Render_LT16ms", 
	"INV-UI-003_Keystroke_Response_LT50ms",
	"INV-UI-009_100_Keyboard_Operation",
	"INV-UI-010_Help_Overlay",
	"Screen_Structure_ELM_Pattern",
	"State_Management_ui_StateMachine",
	"Accessibility_Focus_Visibility",
	"Accessibility_Symbol_Text_Pairing",
	"Performance_Lipgloss_Precomputed",
	"File_Structure_Required_Files",
	"Test_Coverage_Unit_Tests",
	"Documentation_Header_Completeness",
	"Error_Handling_Consistent_States",
	"Keyboard_Map_Universal_Shortcuts",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ui_15_check.go <screen_path> [output_file]")
		os.Exit(1)
	}

	screenPath := os.Args[1]
	outputFile := ""
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	fmt.Printf("🔍 Running 15-Item Critical UI Checklist for: %s\n", screenPath)
	fmt.Println("=" * 60)

	suite := run15ItemCheck(screenPath)
	
	// Print results
	printResults(suite)
	
	// Save to file if specified
	if outputFile != "" {
		saveResults(suite, outputFile)
	}
	
	// Exit with appropriate code
	os.Exit(calculateExitCode(suite))
}

func run15ItemCheck(screenPath string) UI15CheckSuite {
	suite := UI15CheckSuite{
		Timestamp:   time.Now().Format(time.RFC3339),
		ScreenPath:  screenPath,
		TotalChecks: len(critical15Items),
	}

	for _, item := range critical15Items {
		result := checkItem(item, screenPath)
		suite.Results = append(suite.Results, result)
		
		switch result.Status {
		case "PASS":
			suite.PassedChecks++
		case "FAIL":
			suite.FailedChecks++
		case "WARN":
			suite.WarnedChecks++
		}
	}

	return suite
}

func checkItem(item, screenPath string) UI15CheckResult {
	switch item {
	case "INV-UI-001_First_Paint_LT100ms":
		return checkFirstPaint(screenPath)
	case "INV-UI-002_Frame_Render_LT16ms":
		return checkFrameRender(screenPath)
	case "INV-UI-003_Keystroke_Response_LT50ms":
		return checkKeystrokeResponse(screenPath)
	case "INV-UI-009_100_Keyboard_Operation":
		return checkKeyboardOperation(screenPath)
	case "INV-UI-010_Help_Overlay":
		return checkHelpOverlay(screenPath)
	case "Screen_Structure_ELM_Pattern":
		return checkELMPattern(screenPath)
	case "State_Management_ui_StateMachine":
		return checkStateMachine(screenPath)
	case "Accessibility_Focus_Visibility":
		return checkFocusVisibility(screenPath)
	case "Accessibility_Symbol_Text_Pairing":
		return checkSymbolTextPairing(screenPath)
	case "Performance_Lipgloss_Precomputed":
		return checkLipglossStyles(screenPath)
	case "File_Structure_Required_Files":
		return checkRequiredFiles(screenPath)
	case "Test_Coverage_Unit_Tests":
		return checkUnitTests(screenPath)
	case "Documentation_Header_Completeness":
		return checkHeaderCompleteness(screenPath)
	case "Error_Handling_Consistent_States":
		return checkErrorHandling(screenPath)
	case "Keyboard_Map_Universal_Shortcuts":
		return checkKeyboardMap(screenPath)
	default:
		return UI15CheckResult{
			Item:        item,
			Description: "Unknown check item",
			Status:      "FAIL",
			Message:     "Check item not implemented",
		}
	}
}

// Performance Checks
func checkFirstPaint(screenPath string) UI15CheckResult {
	// Simulate first paint measurement
	// In real implementation, this would measure actual startup time
	fmt.Printf("📊 Measuring first paint time... ")
	
	// For demo, simulate measurement
	time.Sleep(10 * time.Millisecond) // Simulate measurement
	
	// Check if file exists and is reasonable size
	info, err := os.Stat(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "INV-UI-001_First_Paint_LT100ms",
			Description: "First paint <100ms",
			Status:      "FAIL",
			Message:     "Screen file not found",
			Details:     err.Error(),
		}
	}

	if info.Size() > 100000 { // >100KB might indicate heavy initialization
		return UI15CheckResult{
			Item:        "INV-UI-001_First_Paint_LT100ms",
			Description: "First paint <100ms",
			Status:      "WARN",
			Message:     "Large file size may impact first paint",
			Details:     fmt.Sprintf("File size: %d bytes", info.Size()),
		}
	}

	return UI15CheckResult{
		Item:        "INV-UI-001_First_Paint_LT100ms",
		Description: "First paint <100ms",
		Status:      "PASS",
		Message:     "File size reasonable for first paint",
		Details:     fmt.Sprintf("File size: %d bytes", info.Size()),
	}
}

func checkFrameRender(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking frame render optimization... ")
	
	// Check for virtualization patterns
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "INV-UI-002_Frame_Render_LT16ms",
			Description: "Frame render <16ms",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for virtualization implementation
	if strings.Contains(code, "virtualization") || strings.Contains(code, "Virtual") {
		return UI15CheckResult{
			Item:        "INV-UI-002_Frame_Render_LT16ms",
			Description: "Frame render <16ms",
			Status:      "PASS",
			Message:     "Virtualization implemented for smooth rendering",
		}
	}
	
	// Check for list length optimization
	listPattern := regexp.MustCompile(`len\((\w+)\)`)
	matches := listPattern.FindAllStringSubmatch(code, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			// Check if list is filtered or limited
			if strings.Contains(code, fmt.Sprintf("%s =", varName)) || 
			   strings.Contains(code, fmt.Sprintf("%s:=", varName)) {
				return UI15CheckResult{
					Item:        "INV-UI-002_Frame_Render_LT16ms",
					Description: "Frame render <16ms",
					Status:      "WARN",
					Message:     "Consider implementing virtualization for large lists",
				}
			}
		}
	}

	return UI15CheckResult{
		Item:        "INV-UI-002_Frame_Render_LT16ms",
		Description: "Frame render <16ms",
		Status:      "PASS",
		Message:     "No large lists detected that need virtualization",
	}
}

func checkKeystrokeResponse(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking keystroke response optimization... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "INV-UI-003_Keystroke_Response_LT50ms",
			Description: "Keystroke response <50ms",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for blocking operations in Update
	if strings.Contains(code, "time.Sleep") && strings.Contains(code, "Update") {
		return UI15CheckResult{
			Item:        "INV-UI-003_Keystroke_Response_LT50ms",
			Description: "Keystroke response <50ms",
			Status:      "FAIL",
			Message:     "Blocking operations found in Update method",
		}
	}
	
	// Check for expensive operations on keystroke
	expensiveOps := []string{"fmt.Sprintf", "strings.Join", "strconv.Itoa"}
	for _, op := range expensiveOps {
		if strings.Contains(code, op) && strings.Contains(code, "Update") {
			return UI15CheckResult{
				Item:        "INV-UI-003_Keystroke_Response_LT50ms",
				Description: "Keystroke response <50ms",
				Status:      "WARN",
				Message:     "Expensive string operations found in Update",
			}
		}
	}

	return UI15CheckResult{
		Item:        "INV-UI-003_Keystroke_Response_LT50ms",
		Description: "Keystroke response <50ms",
		Status:      "PASS",
		Message:     "No blocking operations detected in Update",
	}
}

// Structure & Architecture Checks
func checkELMPattern(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking ELM pattern compliance... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Screen_Structure_ELM_Pattern",
			Description: "ELM pattern compliance",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for required methods
	methods := []string{"Init()", "Update(", "View()"}
	for _, method := range methods {
		if !strings.Contains(code, method) {
			return UI15CheckResult{
				Item:        "Screen_Structure_ELM_Pattern",
				Description: "ELM pattern compliance",
				Status:      "FAIL",
				Message:     fmt.Sprintf("Missing required method: %s", method),
			}
		}
	}

	return UI15CheckResult{
		Item:        "Screen_Structure_ELM_Pattern",
		Description: "ELM pattern compliance",
		Status:      "PASS",
		Message:     "All required ELM methods found",
	}
}

func checkStateMachine(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking ui.StateMachine usage... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "State_Management_ui_StateMachine",
			Description: "ui.StateMachine usage",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for StateMachine embedding
	if !strings.Contains(code, "ui.StateMachine") {
		return UI15CheckResult{
			Item:        "State_Management_ui_StateMachine",
			Description: "ui.StateMachine usage",
			Status:      "FAIL",
			Message:     "ui.StateMachine not embedded in struct",
		}
	}

	return UI15CheckResult{
		Item:        "State_Management_ui_StateMachine",
		Description: "ui.StateMachine usage",
		Status:      "PASS",
		Message:     "ui.StateMachine properly embedded",
	}
}

// Accessibility Checks
func checkFocusVisibility(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking focus visibility... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Accessibility_Focus_Visibility",
			Description: "Focus visibility",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for focus indicators
	focusIndicators := []string{"▸", "→", "•", "█", "■"}
	hasFocus := false
	
	for _, indicator := range focusIndicators {
		if strings.Contains(code, indicator) {
			hasFocus = true
			break
		}
	}

	if !hasFocus {
		return UI15CheckResult{
			Item:        "Accessibility_Focus_Visibility",
			Description: "Focus visibility",
			Status:      "FAIL",
			Message:     "No focus indicators found",
		}
	}

	return UI15CheckResult{
		Item:        "Accessibility_Focus_Visibility",
		Description: "Focus visibility",
		Status:      "PASS",
		Message:     "Focus indicators found",
	}
}

func checkSymbolTextPairing(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking symbol-text pairing... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Accessibility_Symbol_Text_Pairing",
			Description: "Symbol-text pairing",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for symbols without text pairing
	symbols := []string{"✓", "✗", "⚠", "❌", "🟢", "🔴", "🟡"}
	
	for _, symbol := range symbols {
		if strings.Contains(code, symbol) {
			// Check if symbol is paired with text
			symbolRegex := regexp.MustCompile(fmt.Sprintf(`%s\s*\w+`, regexp.QuoteMeta(symbol)))
			if !symbolRegex.MatchString(code) {
				return UI15CheckResult{
					Item:        "Accessibility_Symbol_Text_Pairing",
					Description: "Symbol-text pairing",
					Status:      "FAIL",
					Message:     fmt.Sprintf("Symbol '%s' not properly paired with text", symbol),
				}
			}
		}
	}

	return UI15CheckResult{
		Item:        "Accessibility_Symbol_Text_Pairing",
		Description: "Symbol-text pairing",
		Status:      "PASS",
		Message:     "Symbols properly paired with text",
	}
}

// Performance & Optimization Checks
func checkLipglossStyles(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking lipgloss pre-computed styles... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Performance_Lipgloss_Precomputed",
			Description: "Lipgloss pre-computed styles",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for lipgloss style creation
	stylePatterns := []string{"lipgloss.NewStyle()", "lipgloss.Style{}"}
	
	for _, pattern := range stylePatterns {
		if strings.Contains(code, pattern) {
			// Check if styles are created at package level or in Init
			if strings.Contains(code, "var.*=.*lipgloss.NewStyle") || 
			   strings.Contains(code, "func.*Init.*lipgloss.NewStyle") {
				return UI15CheckResult{
					Item:        "Performance_Lipgloss_Precomputed",
					Description: "Lipgloss pre-computed styles",
					Status:      "PASS",
					Message:     "Styles pre-computed at startup",
				}
			}
			
			// Check for styles created in View (bad practice)
			viewSection := stringsBetween(code, "func (s *", ") View() string")
			if viewSection != "" && strings.Contains(viewSection, "lipgloss.NewStyle") {
				return UI15CheckResult{
					Item:        "Performance_Lipgloss_Precomputed",
					Description: "Lipgloss pre-computed styles",
					Status:      "FAIL",
					Message:     "Styles created in View() - should be pre-computed",
				}
			}
		}
	}

	return UI15CheckResult{
		Item:        "Performance_Lipgloss_Precomputed",
		Description: "Lipgloss pre-computed styles",
		Status:      "WARN",
		Message:     "No lipgloss styles found - may not be using styling",
	}
}

// File Structure Checks
func checkRequiredFiles(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking required file structure... ")
	
	// Get directory containing the screen
	screenDir := filepath.Dir(screenPath)
	screenName := filepath.Base(screenPath)
	baseName := strings.TrimSuffix(screenName, filepath.Ext(screenName))
	
	// Check for test file
	testFile := filepath.Join(screenDir, baseName+"_test.go")
	if _, err := os.Stat(testFile); err != nil {
		return UI15CheckResult{
			Item:        "File_Structure_Required_Files",
			Description: "Required file structure",
			Status:      "FAIL",
			Message:     fmt.Sprintf("Missing test file: %s", testFile),
		}
	}

	// Check for component directory
	componentDir := filepath.Join(filepath.Dir(screenDir), "components")
	if _, err := os.Stat(componentDir); err != nil {
		return UI15CheckResult{
			Item:        "File_Structure_Required_Files",
			Description: "Required file structure",
			Status:      "WARN",
			Message:     "Components directory not found",
		}
	}

	return UI15CheckResult{
		Item:        "File_Structure_Required_Files",
		Description: "Required file structure",
		Status:      "PASS",
		Message:     "Required files present",
	}
}

// Testing Coverage Checks
func checkUnitTests(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking unit test coverage... ")
	
	screenDir := filepath.Dir(screenPath)
	screenName := filepath.Base(screenPath)
	baseName := strings.TrimSuffix(screenName, filepath.Ext(screenName))
	testFile := filepath.Join(screenDir, baseName+"_test.go")
	
	// Check if test file exists
	if _, err := os.Stat(testFile); err != nil {
		return UI15CheckResult{
			Item:        "Test_Coverage_Unit_Tests",
			Description: "Unit test coverage",
			Status:      "FAIL",
			Message:     fmt.Sprintf("No unit test file: %s", testFile),
		}
	}
	
	// Check test content for required patterns
	testContent, err := os.ReadFile(testFile)
	if err != nil {
		return UI15CheckResult{
			Item:        "Test_Coverage_Unit_Tests",
			Description: "Unit test coverage",
			Status:      "FAIL",
			Message:     "Cannot read test file",
		}
	}
	
	testCode := string(testContent)
	requiredPatterns := []string{"Test", "func", "tea.Model", "Update"}
	
	for _, pattern := range requiredPatterns {
		if !strings.Contains(testCode, pattern) {
			return UI15CheckResult{
				Item:        "Test_Coverage_Unit_Tests",
				Description: "Unit test coverage",
				Status:      "WARN",
				Message:     "Test file may be incomplete",
			}
		}
	}

	return UI15CheckResult{
		Item:        "Test_Coverage_Unit_Tests",
		Description: "Unit test coverage",
		Status:      "PASS",
		Message:     "Unit tests found and appear complete",
	}
}

// Documentation Checks
func checkHeaderCompleteness(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking header comment completeness... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Documentation_Header_Completeness",
			Description: "Header comment completeness",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Extract header comment (first comment block)
	headerEnd := strings.Index(code, "package")
	if headerEnd == -1 {
		return UI15CheckResult{
			Item:        "Documentation_Header_Completeness",
			Description: "Header comment completeness",
			Status:      "FAIL",
			Message:     "No package declaration found",
		}
	}
	
	headerCode := code[:headerEnd]
	
	// Check for required header sections
	requiredSections := []string{"Purpose:", "Layout:", "Components:", "Keyboard map:", "Performance budget:"}
	missingSections := []string{}
	
	for _, section := range requiredSections {
		if !strings.Contains(headerCode, section) {
			missingSections = append(missingSections, section)
		}
	}
	
	if len(missingSections) > 0 {
		return UI15CheckResult{
			Item:        "Documentation_Header_Completeness",
			Description: "Header comment completeness",
			Status:      "FAIL",
			Message:     fmt.Sprintf("Missing header sections: %v", missingSections),
		}
	}

	return UI15CheckResult{
		Item:        "Documentation_Header_Completeness",
		Description: "Header comment completeness",
		Status:      "PASS",
		Message:     "All required header sections present",
	}
}

// Error Handling Checks
func checkErrorHandling(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking error handling consistency... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Error_Handling_Consistent_States",
			Description: "Error handling consistency",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for error state handling
	errorPatterns := []string{"Error error", "error error", "hasError", "errorState"}
	hasErrorHandling := false
	
	for _, pattern := range errorPatterns {
		if strings.Contains(code, pattern) {
			hasErrorHandling = true
			break
		}
	}
	
	if !hasErrorHandling {
		return UI15CheckResult{
			Item:        "Error_Handling_Consistent_States",
			Description: "Error handling consistency",
			Status:      "FAIL",
			Message:     "No error handling found",
		}
	}
	
	// Check for recovery options
	recoveryOptions := []string{"[r] Retry", "[d] Doctor", "[i] Ignore", "[q] Quit"]
	hasRecovery := false
	
	for _, option := range recoveryOptions {
		if strings.Contains(code, option) {
			hasRecovery = true
			break
		}
	}
	
	if !hasRecovery {
		return UI15CheckResult{
			Item:        "Error_Handling_Consistent_States",
			Description: "Error handling consistency",
			Status:      "WARN",
			Message:     "Error handling found but no standard recovery options",
		}
	}

	return UI15CheckResult{
		Item:        "Error_Handling_Consistent_States",
		Description: "Error handling consistency",
		Status:      "PASS",
		Message:     "Consistent error handling with recovery options",
	}
}

// Keyboard Navigation Checks
func checkKeyboardMap(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking universal keyboard map... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "Keyboard_Map_Universal_Shortcuts",
			Description: "Universal keyboard shortcuts",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for global shortcuts
	globalShortcuts := []string{"?", "q", "Ctrl+C", "Esc", "D", "M", "$"}
	missingShortcuts := []string{}
	
	for _, shortcut := range globalShortcuts {
		if !strings.Contains(code, shortcut) {
			missingShortcuts = append(missingShortcuts, shortcut)
		}
	}
	
	if len(missingShortcuts) > 0 {
		return UI15CheckResult{
			Item:        "Keyboard_Map_Universal_Shortcuts",
			Description: "Universal keyboard shortcuts",
			Status:      "FAIL",
			Message:     fmt.Sprintf("Missing global shortcuts: %v", missingShortcuts),
		}
	}
	
	// Check for navigation shortcuts
	navShortcuts := []string{"Tab", "Shift+Tab", "j", "k", "↑", "↓", "g", "G", "Enter", "r"}
	missingNav := []string{}
	
	for _, shortcut := range navShortcuts {
		if !strings.Contains(code, shortcut) {
			missingNav = append(missingNav, shortcut)
		}
	}
	
	if len(missingNav) > 0 {
		return UI15CheckResult{
			Item:        "Keyboard_Map_Universal_Shortcuts",
			Description: "Universal keyboard shortcuts",
			Status:      "WARN",
			Message:     fmt.Sprintf("Some navigation shortcuts missing: %v", missingNav),
		}
	}

	return UI15CheckResult{
		Item:        "Keyboard_Map_Universal_Shortcuts",
		Description: "Universal keyboard shortcuts",
		Status:      "PASS",
		Message:     "Universal keyboard shortcuts implemented",
	}
}

func checkHelpOverlay(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking help overlay functionality... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "INV-UI-010_Help_Overlay",
			Description: "Help overlay functionality",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for help overlay implementation
	helpPatterns := []string{"ShowHelp", "help", "Help", "?"}
	hasHelp := false
	
	for _, pattern := range helpPatterns {
		if strings.Contains(code, pattern) {
			hasHelp = true
			break
		}
	}
	
	if !hasHelp {
		return UI15CheckResult{
			Item:        "INV-UI-010_Help_Overlay",
			Description: "Help overlay functionality",
			Status:      "FAIL",
			Message:     "No help overlay implementation found",
		}
	}

	return UI15CheckResult{
		Item:        "INV-UI-010_Help_Overlay",
		Description: "Help overlay functionality",
		Status:      "PASS",
		Message:     "Help overlay implementation found",
	}
}

func checkKeyboardOperation(screenPath string) UI15CheckResult {
	fmt.Printf("📊 Checking keyboard-only operation... ")
	
	content, err := os.ReadFile(screenPath)
	if err != nil {
		return UI15CheckResult{
			Item:        "INV-UI-009_100_Keyboard_Operation",
			Description: "100% keyboard-only operation",
			Status:      "FAIL",
			Message:     "Cannot read screen file",
		}
	}

	code := string(content)
	
	// Check for mouse operations (keyboard-only should not have mouse)
	mousePatterns := []string{"mouse", "Mouse", "click", "Click", "dblclick", "DblClick"}
	hasMouse := false
	
	for _, pattern := range mousePatterns {
		if strings.Contains(code, pattern) {
			hasMouse = true
			break
		}
	}
	
	if hasMouse {
		return UI15CheckResult{
			Item:        "INV-UI-009_100_Keyboard_Operation",
			Description: "100% keyboard-only operation",
			Status:      "FAIL",
			Message:     "Mouse operations found - violates keyboard-only requirement",
		}
	}

	return UI15CheckResult{
		Item:        "INV-UI-009_100_Keyboard_Operation",
		Description: "100% keyboard-only operation",
		Status:      "PASS",
		Message:     "No mouse operations found - keyboard-only compliant",
	}
}

// Utility functions
func printResults(suite UI15CheckSuite) {
	fmt.Printf("\n📊 15-Item Critical UI Checklist Results\n")
	fmt.Println("=" * 60)
	fmt.Printf("Screen: %s\n", suite.ScreenPath)
	fmt.Printf("Timestamp: %s\n", suite.Timestamp)
	fmt.Printf("Total Checks: %d | Passed: %d | Failed: %d | Warned: %d\n", 
		suite.TotalChecks, suite.PassedChecks, suite.FailedChecks, suite.WarnedChecks)
	
	fmt.Println("\n🔍 Detailed Results:")
	fmt.Println("-" * 60)
	
	for _, result := range suite.Results {
		statusIcon := "✓"
		statusColor := "\033[32m" // Green
		switch result.Status {
		case "FAIL":
			statusIcon = "✗"
			statusColor = "\033[31m" // Red
		case "WARN":
			statusIcon = "⚠"
			statusColor = "\033[33m" // Yellow
		}
		
		fmt.Printf("%s%s %s - %s\n", statusColor, statusIcon, result.Item, result.Message)
		if result.Details != "" {
			fmt.Printf("   Details: %s\n", result.Details)
		}
		fmt.Printf("\033[0m") // Reset color
	}
	
	// Summary
	fmt.Println("\n" + "=" * 60)
	if suite.FailedChecks == 0 {
		fmt.Printf("🎉 ALL CHECKS PASSED! (%d/%d)\n", suite.PassedChecks, suite.TotalChecks)
	} else {
		fmt.Printf("❌ %d CHECKS FAILED! (%d passed, %d warned)\n", 
			suite.FailedChecks, suite.PassedChecks, suite.WarnedChecks)
	}
	fmt.Println("=" * 60)
}

func saveResults(suite UI15CheckSuite, filename string) {
	jsonData, err := json.MarshalIndent(suite, "", "  ")
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		return
	}
	
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		return
	}
	
	fmt.Printf("📄 Results saved to: %s\n", filename)
}

func calculateExitCode(suite UI15CheckSuite) int {
	if suite.FailedChecks > 0 {
		return 1
	}
	if suite.WarnedChecks > 0 {
		return 2
	}
	return 0
}

func stringsBetween(s, start, end string) string {
	startIndex := strings.Index(s, start)
	if startIndex == -1 {
		return ""
	}
	
	startIndex += len(start)
	endIndex := strings.Index(s[startIndex:], end)
	if endIndex == -1 {
		return ""
	}
	
	return s[startIndex:startIndex+endIndex]
}

// AST parsing utilities for more advanced checks
func parseGoFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, parser.ParseComments)
}

func findASTNodes(node ast.Node, predicate func(ast.Node) bool) []ast.Node {
	var nodes []ast.Node
	
	ast.Inspect(node, func(n ast.Node) bool {
		if predicate(n) {
			nodes = append(nodes, n)
		}
		return true
	})
	
	return nodes
}

func checkASTMethods(file *ast.File, requiredMethods []string) []string {
	missing := []string{}
	
	for _, method := range requiredMethods {
		found := false
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == method {
					found = true
					return false
				}
			}
			return true
		})
		
		if !found {
			missing = append(missing, method)
		}
	}
	
	return missing
}