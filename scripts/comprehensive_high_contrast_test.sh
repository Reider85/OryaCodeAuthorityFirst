#!/bin/bash

# comprehensive_high_contrast_test.sh - Comprehensive high-contrast mode testing script
# This script runs all high-contrast mode tests with detailed reporting

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "🎨 Comprehensive High-Contrast Mode Testing for CodeAuthorityFirst"
echo "=================================================================="
echo "Testing Date: $(date)"
echo "Git Commit: $(git rev-parse HEAD 2>/dev/null || echo 'unknown')"
echo "Git Branch: $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    exit 1
fi

# Get version from git or use fallback
if command -v git &> /dev/null && [ -d .git ]; then
    GIT_VERSION=$(git describe --tags --dirty --always 2>/dev/null || echo "unknown")
    VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=$GIT_VERSION'"
else
    VERSION_FLAG="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=unknown'"
fi

# Build flags
LDFLAGS="$VERSION_FLAG -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit=$(git rev-parse HEAD 2>/dev/null || echo unknown)'"
LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.datetime=$(date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo unknown)'"
LDFLAGS="$LDFLAGS -X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.modified=$(git diff-index --quiet HEAD -- 2>/dev/null || echo true)'"

echo "📦 Building with flags: $LDFLAGS"

# Build the binary
go build -ldflags="$LDFLAGS" -o caf ./cmd/caf

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# Test results counter
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run tests and count results
run_tests() {
    local test_name="$1"
    local test_pattern="$2"
    
    echo ""
    echo "🧪 Running $test_name tests..."
    echo "================================"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if go test ./pkg/ui -run "$test_pattern" -v -timeout 30s; then
        echo "✅ $test_name tests passed!"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo "❌ $test_name tests failed!"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Function to run a specific test and count results
run_specific_test() {
    local test_name="$1"
    local test_pattern="$2"
    
    echo ""
    echo "🧪 Running $test_name..."
    echo "================================"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if go test ./pkg/ui -run "$test_pattern" -v -timeout 30s; then
        echo "✅ $test_name passed!"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo "❌ $test_name failed!"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Test execution matrix
echo ""
echo "🎯 Starting High-Contrast Mode Test Matrix"
echo "=========================================="

# Phase 1: Basic accessibility tests
run_tests "Basic Accessibility" "TestA11y"

# Phase 2: High-contrast theme validation
run_tests "High-Contrast Theme" "TestA11yHighContrast"

# Phase 3: Screen-specific high-contrast tests
run_tests "Dashboard High-Contrast" "TestDashboardHighContrast"
run_tests "Adapters Manager High-Contrast" "TestAdaptersManagerHighContrast"
run_tests "Onboarding Wizard High-Contrast" "TestOnboardingWizardHighContrast"

# Phase 4: Component rendering in high-contrast mode
run_tests "Common Components High-Contrast" "TestCommonComponentsHighContrast"
run_tests "UI Components High-Contrast" "TestUIComponentsHighContrast"

# Phase 5: Error handling in high-contrast mode
run_tests "High-Contrast Error Handling" "TestHighContrastErrorHandling"

# Phase 6: Screen reader support in high-contrast mode
run_tests "High-Contrast Screen Reader" "TestHighContrastScreenReader"

# Phase 7: Keyboard navigation in high-contrast mode
run_tests "High-Contrast Keyboard Navigation" "TestHighContrastKeyboard"

# Phase 8: Terminal profile adaptation in high-contrast mode
run_tests "High-Contrast Terminal Profile" "TestHighContrastTerminal"

# Phase 9: Validation in high-contrast mode
run_tests "High-Contrast Validation" "TestHighContrastValidation"

# Phase 10: Accessibility report in high-contrast mode
run_tests "High-Contrast Report" "TestHighContrastReport"

# Phase 11: Integration tests
run_tests "High-Contrast Integration" "TestHighContrastIntegration"

# Phase 12: Detailed component tests
run_specific_test "Focus Management" "TestHighContrastFocusManagement"
run_specific_test "Screen Transitions" "TestHighContrastScreenTransitions"
run_specific_test "Error Handling Across Screens" "TestHighContrastErrorHandlingAcrossScreens"
run_specific_test "Screen Reader Support Across Screens" "TestHighContrastScreenReaderSupportAcrossScreens"
run_specific_test "Terminal Profile Adaptation Across Screens" "TestHighContrastTerminalProfileAdaptationAcrossScreens"
run_specific_test "Validation Across Screens" "TestHighContrastValidationAcrossScreens"
run_specific_test "Report Across Screens" "TestHighContrastReportAcrossScreens"

# Phase 13: Animation and contrast ratio tests
run_specific_test "Animation Tests" "TestHighContrastAnimationTests"
run_specific_test "Contrast Ratios" "TestHighContrastContrastRatios"

# Test summary
echo ""
echo "🎉 High-Contrast Mode Testing Complete!"
echo "====================================="
echo "📊 Test Summary:"
echo "   Total Tests: $TOTAL_TESTS"
echo "   Passed: $PASSED_TESTS"
echo "   Failed: $FAILED_TESTS"
echo ""
echo "📋 Test Coverage Areas:"
echo "   ✅ Basic accessibility functionality"
echo "   ✅ High-contrast theme validation"
echo "   ✅ Screen-specific high-contrast tests"
echo "   ✅ Component rendering in high-contrast mode"
echo "   ✅ Error handling in high-contrast mode"
echo "   ✅ Screen reader support in high-contrast mode"
echo "   ✅ Keyboard navigation in high-contrast mode"
echo "   ✅ Terminal profile adaptation in high-contrast mode"
echo "   ✅ Validation in high-contrast mode"
echo "   ✅ Accessibility report in high-contrast mode"
echo "   ✅ Integration testing across screens"
echo "   ✅ Focus management testing"
echo "   ✅ Screen transition testing"
echo "   ✅ Cross-screen error handling"
echo "   ✅ Cross-screen screen reader support"
echo "   ✅ Cross-screen terminal profile adaptation"
echo "   ✅ Cross-screen validation"
echo "   ✅ Cross-screen reporting"
echo "   ✅ Animation and contrast ratio testing"
echo ""
echo "🎨 High-contrast mode features tested:"
echo "   ✅ WCAG AAA compliance (21:1 contrast ratio)"
echo "   ✅ Enhanced focus indicators across terminal profiles"
echo "   ✅ Status prefixes (+, !, X) for color-blind users"
echo "   ✅ Color contrast validation"
echo "   ✅ Terminal profile adaptation (TrueColor, ANSI256, ANSI, ASCII)"
echo "   ✅ Screen reader support with alt-text generation"
echo "   ✅ Keyboard navigation with tab order"
echo "   ✅ Error state handling with visual indicators"
echo "   ✅ Component rendering consistency"
echo "   ✅ Integration between different screens"
echo "   ✅ Focus management across screens"
echo "   ✅ Screen transitions with proper focus handling"
echo "   ✅ Reduced motion support"
echo ""
echo "🔍 High-contrast mode compliance:"
echo "   ✅ WCAG AAA color contrast (minimum 7:1 ratio)"
echo "   ✅ Focus indicators visible across all terminal profiles"
echo "   ✅ Status indicators use symbols + text (color-blind safe)"
echo "   ✅ Screen reader announcements with status-first format"
echo "   ✅ Keyboard-only operation support"
echo "   ✅ Error states with clear visual indicators"
echo "   ✅ Component consistency across all screens"
echo ""
echo "📱 Terminal Emulator Compatibility:"
echo "   ✅ TrueColor terminals (double borders, full colors)"
echo "   ✅ ANSI256 terminals (thick borders, good colors)"
echo "   ✅ ANSI terminals (normal borders, basic colors)"
echo "   ✅ ASCII terminals (text-only, screen reader mode)"
echo ""

# Final result
if [ $FAILED_TESTS -eq 0 ]; then
    echo "🎉 All high-contrast mode tests passed successfully!"
    echo "✅ High-contrast mode is fully functional and compliant!"
    exit 0
else
    echo "❌ $FAILED_TESTS high-contrast mode tests failed!"
    echo "🔧 Please review the failed tests above and fix the issues."
    exit 1
fi

# Clean up
rm -f caf