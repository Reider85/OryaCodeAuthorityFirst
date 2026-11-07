#!/bin/bash

# ui-15-check.sh - Automated 15-Item Critical UI Checklist Runner
# This script runs the automated 15-item critical UI checklist verification

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

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

# Function to find UI screens
find_ui_screens() {
    local screens_dir="$PROJECT_ROOT/pkg/ui/screens"
    
    if [ ! -d "$screens_dir" ]; then
        print_status "FAIL" "UI screens directory not found: $screens_dir"
        exit 1
    fi
    
    # Find all Go files in screens directory (excluding test files)
    find "$screens_dir" -name "*.go" -not -name "*_test.go" | sort
}

# Function to run 15-item check for a single screen
check_screen() {
    local screen_file=$1
    local screen_name=$(basename "$screen_file" .go)
    local screen_dir=$(dirname "$screen_file")
    local output_file="$screen_dir/${screen_name}_15_check_results.json"
    
    print_status "INFO" "Checking $screen_name..."
    
    # Run the Go program
    cd "$PROJECT_ROOT"
    if go run "$SCRIPT_DIR/ui_15_check.go" "$screen_file" "$output_file"; then
        print_status "PASS" "$screen_name check completed"
        return 0
    else
        local exit_code=$?
        if [ $exit_code -eq 1 ]; then
            print_status "FAIL" "$screen_name failed critical checks"
        elif [ $exit_code -eq 2 ]; then
            print_status "WARN" "$screen_name passed with warnings"
        else
            print_status "FAIL" "$screen_name check failed with error code $exit_code"
        fi
        return $exit_code
    fi
}

# Function to generate summary report
generate_summary() {
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local summary_file="$PROJECT_ROOT/ui_15_check_summary_$timestamp.json"
    local summary_report="$PROJECT_ROOT/ui_15_check_summary_$timestamp.md"
    
    print_status "INFO" "Generating summary report..."
    
    # Create JSON summary
    cat > "$summary_file" << EOF
{
  "timestamp": "$(date -Iseconds)",
  "total_screens": 0,
  "passed_screens": 0,
  "failed_screens": 0,
  "warned_screens": 0,
  "screens": []
}
EOF
    
    # Create Markdown summary
    cat > "$summary_report" << EOF
# 15-Item Critical UI Checklist Summary Report

**Generated:** $(date)
**Total Screens Checked:** 0

## Results Overview

| Status | Count | Percentage |
|--------|-------|------------|
| ✅ PASS | 0 | 0% |
| ❌ FAIL | 0 | 0% |
| ⚠ WARN | 0 | 0% |

## Detailed Results

EOF
    
    print_status "INFO" "Summary reports generated:"
    print_status "INFO" "  JSON: $summary_file"
    print_status "INFO" "  Markdown: $summary_report"
}

# Function to check if we're in the right directory
check_project_root() {
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        print_status "FAIL" "Project root not found (go.mod missing)"
        exit 1
    fi
    
    if [ ! -d "$PROJECT_ROOT/pkg/ui/screens" ]; then
        print_status "FAIL" "UI screens directory not found"
        exit 1
    fi
}

# Main function
main() {
    print_status "INFO" "Starting 15-Item Critical UI Checklist Verification"
    print_status "INFO" "Project Root: $PROJECT_ROOT"
    echo "================================================================"
    
    # Check if we're in the right project
    check_project_root
    
    # Find all UI screens
    local screens=($(find_ui_screens))
    local total_screens=${#screens[@]}
    
    print_status "INFO" "Found $total_screens screens to check"
    echo
    
    # Initialize counters
    local passed_screens=0
    local failed_screens=0
    local warned_screens=0
    
    # Check each screen
    for screen_file in "${screens[@]}"; do
        if check_screen "$screen_file"; then
            ((passed_screens++))
        else
            local exit_code=$?
            if [ $exit_code -eq 1 ]; then
                ((failed_screens++))
            elif [ $exit_code -eq 2 ]; then
                ((warned_screens++))
            fi
        fi
        echo
    done
    
    # Generate summary
    generate_summary
    
    # Print final summary
    echo "================================================================"
    print_status "INFO" "15-Item Critical UI Checklist Complete"
    print_status "INFO" "Screens Checked: $total_screens"
    print_status "INFO" "Passed: $passed_screens"
    print_status "INFO" "Failed: $failed_screens"
    print_status "INFO" "Warnings: $warned_screens"
    
    if [ $failed_screens -gt 0 ]; then
        print_status "FAIL" "$failed_screens screens failed critical checks"
        exit 1
    elif [ $warned_screens -gt 0 ]; then
        print_status "WARN" "$warned_screens screens have warnings"
        exit 2
    else
        print_status "PASS" "All screens passed the 15-item critical checklist"
        exit 0
    fi
}

# Run main function
main "$@"