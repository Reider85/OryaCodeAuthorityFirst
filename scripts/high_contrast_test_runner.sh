#!/bin/bash

# high_contrast_test_runner.sh - Comprehensive high-contrast mode test runner
# This script runs all high-contrast mode tests for the CodeAuthorityFirst TUI

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "🎨 Running High-Contrast Mode Tests for CodeAuthorityFirst"
echo "============================================================"

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

# Function to run specific test categories
run_tests() {
    local test_name="$1"
    local test_pattern="$2"
    
    echo ""
    echo "🧪 Running $test_name tests..."
    echo "================================"
    
    if go test ./pkg/ui -run "$test_pattern" -v; then
        echo "✅ $test_name tests passed!"
    else
        echo "❌ $test_name tests failed!"
        return 1
    fi
}

# Run all high-contrast mode tests
echo ""
echo "🎨 Running Comprehensive High-Contrast Mode Tests"
echo "==============================================="

# Test 1: Basic accessibility functionality
run_tests "Basic Accessibility" "TestA11y"

# Test 2: High-contrast theme validation
run_tests "High-Contrast Theme" "TestA11yHighContrast"

# Test 3: Screen-specific high-contrast tests
run_tests "Dashboard High-Contrast" "TestDashboardHighContrast"
run_tests "Adapters Manager High-Contrast" "TestAdaptersManagerHighContrast"
run_tests "Onboarding Wizard High-Contrast" "TestOnboardingWizardHighContrast"

# Test 4: Component rendering in high-contrast mode
run_tests "Common Components High-Contrast" "TestCommonComponentsHighContrast"

# Test 5: Error handling in high-contrast mode
run_tests "High-Contrast Error Handling" "TestHighContrastErrorHandling"

# Test 6: Screen reader support in high-contrast mode
run_tests "High-Contrast Screen Reader" "TestHighContrastScreenReader"

# Test 7: Keyboard navigation in high-contrast mode
run_tests "High-Contrast Keyboard Navigation" "TestHighContrastKeyboard"

# Test 8: Terminal profile adaptation in high-contrast mode
run_tests "High-Contrast Terminal Profile" "TestHighContrastTerminal"

# Test 9: Validation in high-contrast mode
run_tests "High-Contrast Validation" "TestHighContrastValidation"

# Test 10: Accessibility report in high-contrast mode
run_tests "High-Contrast Report" "TestHighContrastReport"

echo ""
echo "🎉 All high-contrast mode tests completed successfully!"
echo ""
echo "📋 Test Coverage Summary:"
echo "========================"
echo "✅ Basic accessibility functionality"
echo "✅ High-contrast theme validation"
echo "✅ Screen-specific high-contrast tests"
echo "✅ Component rendering in high-contrast mode"
echo "✅ Error handling in high-contrast mode"
echo "✅ Screen reader support in high-contrast mode"
echo "✅ Keyboard navigation in high-contrast mode"
echo "✅ Terminal profile adaptation in high-contrast mode"
echo "✅ Validation in high-contrast mode"
echo "✅ Accessibility report in high-contrast mode"
echo ""
echo "🎨 High-contrast mode testing completed successfully!"
echo "All screens and components have been tested for:"
echo "  - WCAG AAA compliance"
echo "  - Focus indicators visibility"
echo "  - Status prefixes (+, !, X)"
echo "  - Color contrast validation"
echo "  - Terminal profile adaptation"
echo "  - Screen reader support"
echo "  - Keyboard navigation"
echo "  - Error state handling"
echo "  - Component rendering consistency"

# Clean up
rm -f caf

echo ""
echo "🧹 Cleanup completed!"