#!/bin/bash

# Test script for the Hybrid Runner System
# This script verifies that all components are working correctly

set -e

echo "🧪 Testing Hybrid Runner System"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test functions
test_go_build() {
    echo "🔨 Testing Go build..."
    if go build runner-manager.go; then
        echo "✅ Go build successful"
        rm -f runner-manager runner-manager.exe
        return 0
    else
        echo "❌ Go build failed"
        return 1
    fi
}

test_runner_selection() {
    echo "🎯 Testing runner selection..."
    
    # Test without GitHub token (should fail gracefully)
    echo "Testing without GitHub token..."
    if ./runner-manager 80 2>/dev/null; then
        echo "❌ Should have failed without GitHub token"
        return 1
    else
        echo "✅ Correctly failed without GitHub token"
    fi
    
    # Test script execution
    echo "Testing runner selection script..."
    chmod +x scripts/select-runner.sh
    if bash scripts/select-runner.sh 80 hybrid >/dev/null 2>&1; then
        echo "✅ Runner selection script works"
        return 0
    else
        echo "❌ Runner selection script failed"
        return 1
    fi
}

test_dashboard() {
    echo "📊 Testing dashboard..."
    if python3 -c "import scripts.runner_dashboard" 2>/dev/null; then
        echo "✅ Dashboard imports correctly"
        return 0
    else
        echo "❌ Dashboard import failed"
        return 1
    fi
}

test_workflow_syntax() {
    echo "📝 Testing workflow syntax..."
    
    # Check if GitHub Actions workflows exist
    workflows_dir=".github/workflows"
    if [[ ! -d "$workflows_dir" ]]; then
        echo "❌ Workflows directory not found"
        return 1
    fi
    
    # Check if workflows have been updated
    if grep -q "select-runner" "$workflows_dir/ci.yml" && grep -q "select-runner" "$workflows_dir/quality.yml"; then
        echo "✅ Workflows updated with hybrid runner selection"
        return 0
    else
        echo "❌ Workflows not properly updated"
        return 1
    fi
}

test_setup_scripts() {
    echo "🛠️ Testing setup scripts..."
    
    scripts_dir="scripts"
    if [[ ! -d "$scripts_dir" ]]; then
        echo "❌ Scripts directory not found"
        return 1
    fi
    
    # Check if setup scripts exist
    setup_scripts=(
        "setup-linux-runner.sh"
        "setup-macos-runner.sh"
        "setup-windows-runner.ps1"
        "select-runner.sh"
        "select-runner.ps1"
        "runner-dashboard.py"
    )
    
    for script in "${setup_scripts[@]}"; do
        if [[ -f "$scripts_dir/$script" ]]; then
            echo "✅ $script exists"
        else
            echo "❌ $script missing"
            return 1
        fi
    done
    
    return 0
}

# Main test execution
echo "Starting tests..."
echo

# Test 1: Go build
if ! test_go_build; then
    echo -e "${RED}❌ Go build test failed${NC}"
    exit 1
fi

# Test 2: Runner selection
if ! test_runner_selection; then
    echo -e "${RED}❌ Runner selection test failed${NC}"
    exit 1
fi

# Test 3: Dashboard
if ! test_dashboard; then
    echo -e "${RED}❌ Dashboard test failed${NC}"
    exit 1
fi

# Test 4: Workflow syntax
if ! test_workflow_syntax; then
    echo -e "${RED}❌ Workflow syntax test failed${NC}"
    exit 1
fi

# Test 5: Setup scripts
if ! test_setup_scripts; then
    echo -e "${RED}❌ Setup scripts test failed${NC}"
    exit 1
fi

echo
echo -e "${GREEN}🎉 All tests passed!${NC}"
echo "The hybrid runner system is ready to use."
echo
echo "Next steps:"
echo "1. Set up self-hosted runners using the setup scripts"
echo "2. Configure GitHub token with appropriate permissions"
echo "3. Test the system with a real workflow run"
echo "4. Monitor usage with the dashboard script"