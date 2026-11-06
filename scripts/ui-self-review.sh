#!/bin/bash

# ui-self-review.sh - Quick UI self-review automation script
# This script helps developers quickly verify their UI implementations
# against the CodeAuthorityFirst UI self-review checklist

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🔍 CodeAuthorityFirst UI Self-Review"
echo "===================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS") echo -e "${GREEN}✓ $message${NC}" ;;
        "FAIL") echo -e "${RED}✗ $message${NC}" ;;
        "WARN") echo -e "${YELLOW}⚠ $message${NC}" ;;
        "INFO") echo -e "${BLUE}ℹ $message${NC}" ;;
    esac
}

# Function to check if file exists
check_file_exists() {
    local file=$1
    local description=$2
    if [ -f "$file" ]; then
        print_status "PASS" "$description: $file"
        return 0
    else
        print_status "FAIL" "$description: $file"
        return 1
    fi
}

# Function to check if directory exists
check_dir_exists() {
    local dir=$1
    local description=$2
    if [ -d "$dir" ]; then
        print_status "PASS" "$description: $dir"
        return 0
    else
        print_status "FAIL" "$description: $dir"
        return 1
    fi
}

# Function to check if command exists
check_command_exists() {
    local cmd=$1
    local description=$2
    if command -v "$cmd" &> /dev/null; then
        print_status "PASS" "$description: $cmd"
        return 0
    else
        print_status "FAIL" "$description: $cmd"
        return 1
    fi
}

# Function to check Go build
check_go_build() {
    local path=$1
    local name=$2
    if go build -o /tmp/test_ui "$path" &> /dev/null; then
        print_status "PASS" "Go build successful: $name"
        rm -f /tmp/test_ui
        return 0
    else
        print_status "FAIL" "Go build failed: $name"
        return 1
    fi
}

# Function to check performance
check_performance() {
    local path=$1
    local name=$2
    local max_time=$3
    
    if go build -o /tmp/test_ui "$path" &> /dev/null; then
        local start_time=$(date +%s%N)
        /tmp/test_ui &> /dev/null
        local end_time=$(date +%s%N)
        local duration=$((($end_time - $start_time) / 1000000))
        
        if [ $duration -lt $max_time ]; then
            print_status "PASS" "Performance $name: ${duration}ms < ${max_time}ms"
            rm -f /tmp/test_ui
            return 0
        else
            print_status "FAIL" "Performance $name: ${duration}ms >= ${max_time}ms"
            rm -f /tmp/test_ui
            return 1
        fi
    else
        print_status "FAIL" "Performance test failed: build error"
        return 1
    fi
}

# Function to check test coverage
check_test_coverage() {
    local path=$1
    local name=$2
    
    if go test -v "$path" 2>&1 | grep -q "PASS"; then
        print_status "PASS" "Tests passing: $name"
        return 0
    else
        print_status "FAIL" "Tests failing: $name"
        return 1
    fi
}

# Function to check UI files
check_ui_files() {
    echo "📁 Checking UI files..."
    
    local ui_files=(
        "docs/UI_SELF_REVIEW_CHECKLIST.md"
        "docs/UI_IMPLEMENTATION_EXAMPLES.md"
        "docs/UI_QUICK_REFERENCE.md"
        "pkg/ui/screens/ui_self_review_test.go"
    )
    
    local missing_files=0
    for file in "${ui_files[@]}"; do
        if ! check_file_exists "$PROJECT_ROOT/$file" "UI documentation"; then
            missing_files=$((missing_files + 1))
        fi
    done
    
    if [ $missing_files -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to check screen implementations
check_screen_implementations() {
    echo "🎬 Checking screen implementations..."
    
    local screens=(
        "dashboard.go"
        "onboarding.go"
        "adapters.go"
    )
    
    local missing_screens=0
    for screen in "${screens[@]}"; do
        local screen_file="$PROJECT_ROOT/pkg/ui/screens/$screen"
        if [ -f "$screen_file" ]; then
            # Check if screen has required methods
            if grep -q "func.*Init()" "$screen_file" && \
               grep -q "func.*Update()" "$screen_file" && \
               grep -q "func.*View()" "$screen_file"; then
                print_status "PASS" "Screen implementation: $screen"
            else
                print_status "FAIL" "Screen missing required methods: $screen"
                missing_screens=$((missing_screens + 1))
            fi
        else
            print_status "WARN" "Screen file not found: $screen"
        fi
    done
    
    if [ $missing_screens -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to check performance requirements
check_performance_requirements() {
    echo "⚡ Checking performance requirements..."
    
    local passed=0
    local total=0
    
    # Dashboard first paint <100ms
    total=$((total + 1))
    if check_performance "$PROJECT_ROOT/cmd/caf" "dashboard" 100; then
        passed=$((passed + 1))
    fi
    
    # Widget render <50ms
    total=$((total + 1))
    if check_performance "$PROJECT_ROOT/pkg/ui/screens/dashboard.go" "widget" 50; then
        passed=$((passed + 1))
    fi
    
    echo "Performance: $passed/$total checks passed"
    [ $passed -eq $total ]
}

# Function to check accessibility
check_accessibility() {
    echo "♿ Checking accessibility requirements..."
    
    local ui_files=(
        "pkg/ui/screens/dashboard.go"
        "pkg/ui/screens/onboarding.go"
        "pkg/ui/screens/adapters.go"
    )
    
    local passed=0
    local total=0
    
    for file in "${ui_files[@]}"; do
        total=$((total + 1))
        if [ -f "$file" ]; then
            # Check for focus visibility
            if grep -q "▸" "$file"; then
                passed=$((passed + 1))
                print_status "PASS" "Focus visibility: $file"
            else
                print_status "WARN" "Focus visibility missing: $file"
            fi
            
            # Check for symbol-text pairing
            if grep -q "✓.*[A-Za-z]" "$file" || grep -q "⚠.*[A-Za-z]" "$file"; then
                passed=$((passed + 1))
                print_status "PASS" "Symbol-text pairing: $file"
            else
                print_status "WARN" "Symbol-text pairing missing: $file"
            fi
        fi
    done
    
    echo "Accessibility: $passed/$total checks passed"
    [ $passed -gt 0 ]
}

# Function to check keyboard shortcuts
check_keyboard_shortcuts() {
    echo "⌨️ Checking keyboard shortcuts..."
    
    local ui_files=(
        "pkg/ui/screens/dashboard.go"
        "pkg/ui/screens/onboarding.go"
        "pkg/ui/screens/adapters.go"
    )
    
    local passed=0
    local total=0
    
    for file in "${ui_files[@]}"; do
        total=$((total + 1))
        if [ -f "$file" ]; then
            # Check for help overlay
            if grep -q "ShowHelp" "$file"; then
                passed=$((passed + 1))
                print_status "PASS" "Help overlay: $file"
            else
                print_status "WARN" "Help overlay missing: $file"
            fi
            
            # Check for keyboard navigation
            if grep -q "CurrentWidget" "$file" || grep -q "CurrentItem" "$file"; then
                passed=$((passed + 1))
                print_status "PASS" "Keyboard navigation: $file"
            else
                print_status "WARN" "Keyboard navigation missing: $file"
            fi
        fi
    done
    
    echo "Keyboard shortcuts: $passed/$total checks passed"
    [ $passed -gt 0 ]
}

# Function to check documentation
check_documentation() {
    echo "📚 Checking documentation..."
    
    local doc_files=(
        "docs/UI_SELF_REVIEW_CHECKLIST.md"
        "docs/UI_IMPLEMENTATION_EXAMPLES.md"
        "docs/UI_QUICK_REFERENCE.md"
    )
    
    local missing_sections=0
    
    for doc in "${doc_files[@]}"; do
        if [ -f "$doc" ]; then
            # Check for required sections
            if ! grep -q "Performance Requirements" "$doc"; then
                print_status "WARN" "Missing performance section: $doc"
                missing_sections=$((missing_sections + 1))
            fi
            
            if ! grep -q "Accessibility Compliance" "$doc"; then
                print_status "WARN" "Missing accessibility section: $doc"
                missing_sections=$((missing_sections + 1))
            fi
        fi
    done
    
    [ $missing_sections -eq 0 ]
}

# Function to check test coverage
check_test_coverage() {
    echo "🧪 Checking test coverage..."
    
    local test_files=(
        "pkg/ui/screens/dashboard_test.go"
        "pkg/ui/screens/onboarding_test.go"
        "pkg/ui/screens/adapters_test.go"
        "pkg/ui/screens/ui_self_review_test.go"
    )
    
    local missing_tests=0
    
    for test_file in "${test_files[@]}"; do
        if [ -f "$test_file" ]; then
            # Check for performance tests
            if ! grep -q "TestPerformance" "$test_file"; then
                print_status "WARN" "Missing performance test: $test_file"
                missing_tests=$((missing_tests + 1))
            fi
            
            # Check for accessibility tests
            if ! grep -q "TestAccessibility" "$test_file"; then
                print_status "WARN" "Missing accessibility test: $test_file"
                missing_tests=$((missing_tests + 1))
            fi
        else
            print_status "WARN" "Test file missing: $test_file"
            missing_tests=$((missing_tests + 1))
        fi
    done
    
    [ $missing_tests -eq 0 ]
}

# Main function
main() {
    echo "Starting UI self-review..."
    echo
    
    # Initialize counters
    local passed_checks=0
    local total_checks=0
    
    # Run all checks
    echo "📋 Running UI self-review checks..."
    echo
    
    # 1. Check UI files
    total_checks=$((total_checks + 1))
    if check_ui_files; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 2. Check screen implementations
    total_checks=$((total_checks + 1))
    if check_screen_implementations; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 3. Check performance requirements
    total_checks=$((total_checks + 1))
    if check_performance_requirements; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 4. Check accessibility
    total_checks=$((total_checks + 1))
    if check_accessibility; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 5. Check keyboard shortcuts
    total_checks=$((total_checks + 1))
    if check_keyboard_shortcuts; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 6. Check documentation
    total_checks=$((total_checks + 1))
    if check_documentation; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # 7. Check test coverage
    total_checks=$((total_checks + 1))
    if check_test_coverage; then
        passed_checks=$((passed_checks + 1))
    fi
    
    # Summary
    echo
    echo "📊 UI Self-Review Summary"
    echo "========================"
    echo "Passed: $passed_checks/$total_checks checks"
    
    if [ $passed_checks -eq $total_checks ]; then
        echo -e "${GREEN}🎉 All checks passed! UI implementation is ready for review.${NC}"
        exit 0
    else
        echo -e "${YELLOW}⚠ Some checks failed. Please address the issues above.${NC}"
        echo -e "${BLUE}💡 Use the detailed checklist for manual verification.${NC}"
        exit 1
    fi
}

# Run main function
main