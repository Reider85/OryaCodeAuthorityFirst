#!/bin/bash

# build-ui-15-check.sh - Build the 15-item UI checker for distribution
# This script compiles the Go program and creates distribution packages

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

# Function to create distribution package
create_package() {
    local platform=$1
    local output_dir="$PROJECT_ROOT/dist/ui-15-check-$platform"
    
    print_status "INFO" "Creating distribution package for $platform..."
    
    # Create output directory
    mkdir -p "$output_dir"
    
    # Copy files
    cp "$SCRIPT_DIR/ui_15_check.go" "$output_dir/"
    cp "$SCRIPT_DIR/ui-15-check.sh" "$output_dir/"
    
    # Make shell script executable
    chmod +x "$output_dir/ui-15-check.sh"
    
    # Create README
    cat > "$output_dir/README.md" << EOF
# 15-Item Critical UI Checker

This tool provides automated verification of the 15 most critical UI requirements for CodeAuthorityFirst.

## Usage

### Quick Check (All Screens)
\`\`\`bash
./ui-15-check.sh --all
\`\`\`

### Check Specific Screen
\`\`\`bash
./ui-15-check.sh --screen /path/to/screen.go
\`\`\`

### Options
- \`--all\`: Check all UI screens
- \`--screen <path>\`: Check specific screen file
- \`--output <dir>\`: Output directory for results
- \`--verbose\`: Verbose output

## Requirements
- Go 1.19 or later
- Bash shell (Linux/macOS)
- Windows PowerShell (Windows)

## Output
- JSON results for each screen
- Summary report in JSON and Markdown formats
- Colored console output

## Exit Codes
- 0: All checks passed
- 1: Critical checks failed
- 2: Warnings only
EOF
    
    print_status "PASS" "Package created: $output_dir"
}

# Function to build Go binary
build_binary() {
    local platform=$1
    local output_dir="$PROJECT_ROOT/dist/ui-15-check-$platform"
    local binary_name="ui-15-check"
    
    print_status "INFO" "Building Go binary for $platform..."
    
    # Create output directory
    mkdir -p "$output_dir"
    
    # Build for the target platform
    if [ "$platform" = "windows" ]; then
        GOOS=windows GOARCH=amd64 go build -o "$output_dir/$binary_name.exe" "$SCRIPT_DIR/ui_15_check.go"
    else
        GOOS=linux GOARCH=amd64 go build -o "$output_dir/$binary_name" "$SCRIPT_DIR/ui_15_check.go"
        chmod +x "$output_dir/$binary_name"
    fi
    
    print_status "PASS" "Binary built: $output_dir/$binary_name"
}

# Function to create CI/CD integration files
create_ci_files() {
    local ci_dir="$PROJECT_ROOT/.github/workflows"
    
    print_status "INFO" "Creating CI/CD integration files..."
    
    # Create UI check workflow
    mkdir -p "$ci_dir"
    
    cat > "$ci_dir/ui-15-check.yml" << EOF
name: UI 15-Item Critical Checklist

on:
  pull_request:
    paths:
      - 'pkg/ui/screens/**/*.go'
      - 'scripts/ui-15-check*'
  push:
    branches: [ main, master ]
    paths:
      - 'pkg/ui/screens/**/*.go'
      - 'scripts/ui-15-check*'

jobs:
  ui-check:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Build UI checker
      run: |
        chmod +x scripts/build-ui-15-check.sh
        ./scripts/build-ui-15-check.sh
    
    - name: Run UI 15-item checklist
      run: |
        ./dist/ui-15-check-linux/ui-15-check.sh --all
    
    - name: Upload results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: ui-check-results
        path: |
          ui_15_check_summary_*.json
          ui_15_check_summary_*.md
          dist/ui-15-check-*/**/*_15_check_results.json
EOF

    # Create pre-commit hook
    mkdir -p "$PROJECT_ROOT/.git/hooks"
    
    cat > "$PROJECT_ROOT/.git/hooks/pre-commit" << EOF
#!/bin/bash

# Pre-commit hook for UI 15-item checklist
# This runs the UI checker before committing changes

echo "Running UI 15-item checklist..."

# Check if UI checker exists
if [ -f "scripts/ui-15-check.sh" ]; then
    chmod +x scripts/ui-15-check.sh
    
    # Run UI checker on changed screens
    changed_files=\$(git diff --cached --name-only --diff-filter=ACM | grep "pkg/ui/screens.*\.go$" || true)
    
    if [ -n "\$changed_files" ]; then
        echo "Checking changed UI screens..."
        for file in \$changed_files; do
            if [ -f "\$file" ]; then
                echo "Checking \$file..."
                ./scripts/ui-15-check.sh --screen "\$file"
                if [ \$? -ne 0 ]; then
                    echo "❌ UI checks failed for \$file"
                    echo "Please fix the issues before committing."
                    exit 1
                fi
            fi
        done
    fi
fi

echo "✅ UI 15-item checklist passed"
EOF

    chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"
    
    print_status "PASS" "CI/CD integration files created"
}

# Function to create Makefile targets
create_makefile() {
    local makefile="$PROJECT_ROOT/Makefile"
    
    print_status "INFO" "Creating Makefile targets..."
    
    # Check if Makefile exists, create if not
    if [ ! -f "$makefile" ]; then
        cat > "$makefile" << EOF
# Makefile for UI 15-item checklist

.PHONY: ui-15-check build-ui-15-check clean-ui-15-check

# Run UI 15-item checklist on all screens
ui-15-check:
	@echo "Running UI 15-item checklist..."
	@chmod +x scripts/ui-15-check.sh
	@./scripts/ui-15-check.sh --all

# Run UI 15-item checklist on specific screen
ui-15-check-screen:
	@echo "Running UI 15-item checklist on screen: \$(SCREEN)"
	@chmod +x scripts/ui-15-check.sh
	@./scripts/ui-15-check.sh --screen \$(SCREEN)

# Build UI 15-item checker
build-ui-15-check:
	@echo "Building UI 15-item checker..."
	@chmod +x scripts/build-ui-15-check.sh
	@./scripts/build-ui-15-check.sh

# Clean UI 15-item checker artifacts
clean-ui-15-check:
	@echo "Cleaning UI 15-item checker artifacts..."
	@rm -rf dist/
	@rm -f ui_15_check_summary_*.json
	@rm -f ui_15_check_summary_*.md

# Install UI 15-item checker system-wide
install-ui-15-check:
	@echo "Installing UI 15-item checker..."
	@mkdir -p \$(HOME)/.local/bin
	@cp dist/ui-15-check-linux/ui-15-check.sh \$(HOME)/.local/bin/ui-15-check
	@chmod +x \$(HOME)/.local/bin/ui-15-check
	@echo "✅ Installed to \$(HOME)/.local/bin/ui-15-check"
	@echo "📝 Add \$(HOME)/.local/bin to your PATH if not already there"
EOF
    else
        # Append to existing Makefile
        cat >> "$makefile" << EOF

# UI 15-item checklist targets
.PHONY: ui-15-check build-ui-15-check clean-ui-15-check

ui-15-check:
	@echo "Running UI 15-item checklist..."
	@chmod +x scripts/ui-15-check.sh
	@./scripts/ui-15-check.sh --all

ui-15-check-screen:
	@echo "Running UI 15-item checklist on screen: \$(SCREEN)"
	@chmod +x scripts/ui-15-check.sh
	@./scripts/ui-15-check.sh --screen \$(SCREEN)

build-ui-15-check:
	@echo "Building UI 15-item checker..."
	@chmod +x scripts/build-ui-15-check.sh
	@./scripts/build-ui-15-check.sh

clean-ui-15-check:
	@echo "Cleaning UI 15-item checker artifacts..."
	@rm -rf dist/
	@rm -f ui_15_check_summary_*.json
	@rm -f ui_15_check_summary_*.md

install-ui-15-check:
	@echo "Installing UI 15-item checker..."
	@mkdir -p \$(HOME)/.local/bin
	@cp dist/ui-15-check-linux/ui-15-check.sh \$(HOME)/.local/bin/ui-15-check
	@chmod +x \$(HOME)/.local/bin/ui-15-check
	@echo "✅ Installed to \$(HOME)/.local/bin/ui-15-check"
	@echo "📝 Add \$(HOME)/.local/bin to your PATH if not already there"
EOF
    fi
    
    print_status "PASS" "Makefile targets created"
}

# Function to create documentation
create_documentation() {
    local doc_dir="$PROJECT_ROOT/docs/ui-15-check"
    
    print_status "INFO" "Creating documentation..."
    
    # Create documentation directory
    mkdir -p "$doc_dir"
    
    # Create user guide
    cat > "$doc_dir/USER_GUIDE.md" << EOF
# UI 15-Item Critical Checklist User Guide

## Overview

The UI 15-Item Critical Checklist is an automated verification system that ensures all UI implementations meet the most critical requirements for performance, accessibility, and code quality.

## Quick Start

### 1. Run on All Screens
\`\`\`bash
make ui-15-check
\`\`\`

### 2. Run on Specific Screen
\`\`\`bash
make ui-15-check-screen SCREEN=/path/to/screen.go
\`\`\`

### 3. Build Checker
\`\`\`bash
make build-ui-15-check
\`\`\`

## 15 Critical Items

### Performance Requirements
1. **INV-UI-001**: First paint <100ms
2. **INV-UI-002**: Frame render <16ms
3. **INV-UI-003**: Keystroke response <50ms

### Accessibility Requirements
4. **INV-UI-009**: 100% keyboard-only operation
5. **INV-UI-010**: Help overlay functionality
6. **Accessibility**: Focus visibility
7. **Accessibility**: Symbol-text pairing

### Architecture Requirements
8. **Screen Structure**: ELM pattern compliance
9. **State Management**: ui.StateMachine usage

### Code Quality Requirements
10. **Performance**: Lipgloss pre-computed styles
11. **File Structure**: Required files present
12. **Testing**: Unit tests exist
13. **Documentation**: Header comment completeness
14. **Error Handling**: Consistent error states
15. **Keyboard Map**: Universal shortcuts implementation

## Integration

### CI/CD Integration
The system is integrated into the CI/CD pipeline and will run automatically on:
- Pull requests affecting UI screens
- Push to main/master branches

### Pre-commit Hooks
A pre-commit hook is installed that runs the UI checker on changed UI files before commit.

## Troubleshooting

### Common Issues
1. **Missing Dependencies**: Ensure Go 1.19+ is installed
2. **Permission Issues**: Run \`chmod +x scripts/ui-15-check.sh\`
3. **Path Issues**: Use absolute paths for screen files

### Debug Mode
Run with verbose output:
\`\`\`bash
./scripts/ui-15-check.sh --all --verbose
\`\`\`

## Output

### Console Output
Colored console output showing:
- ✅ Passed checks
- ❌ Failed checks
- ⚠ Warnings

### JSON Results
Detailed results in JSON format for each screen:
\`\`\`json
{
  "timestamp": "2026-11-07T10:00:00Z",
  "item": "INV-UI-001_First_Paint_LT100ms",
  "status": "PASS",
  "message": "First paint time acceptable"
}
\`\`\`

### Summary Reports
- JSON summary: \`ui_15_check_summary_*.json\`
- Markdown summary: \`ui_15_check_summary_*.md\`

## Contributing

### Adding New Checks
1. Update \`ui_15_check.go\` with new check function
2. Add check to \`critical15Items\` array
3. Update documentation
4. Test with existing screens

### Maintaining the System
- Regular review of critical items
- Update as UI requirements evolve
- Monitor performance metrics
- Collect feedback from developers
EOF

    # Create developer guide
    cat > "$doc_dir/DEVELOPER_GUIDE.md" << EOF
# UI 15-Item Critical Checklist Developer Guide

## Architecture

### Core Components
1. **ui_15_check.go**: Main checker implementation
2. **ui-15-check.sh**: Linux/macOS runner script
3. **ui-15-check.ps1**: Windows runner script
4. **build-ui-15-check.sh**: Build script for distribution

### Check Categories
The 15 critical items are organized into 5 categories:

#### Performance (3 items)
- INV-UI-001: First paint <100ms
- INV-UI-002: Frame render <16ms
- INV-UI-003: Keystroke response <50ms

#### Accessibility (3 items)
- INV-UI-009: 100% keyboard-only operation
- INV-UI-010: Help overlay functionality
- Focus visibility checks
- Symbol-text pairing checks

#### Architecture (2 items)
- ELM pattern compliance
- ui.StateMachine usage

#### Code Quality (4 items)
- Lipgloss pre-computed styles
- Required file structure
- Unit test coverage
- Header comment completeness

#### Error Handling & UX (3 items)
- Consistent error states
- Universal keyboard shortcuts
- Help overlay functionality

## Implementation Details

### Check Functions
Each check is implemented as a separate function that returns a \`UI15CheckResult\` struct:

\`\`\`go
func checkItem(item, screenPath string) UI15CheckResult {
    switch item {
    case "INV-UI-001_First_Paint_LT100ms":
        return checkFirstPaint(screenPath)
    case "INV-UI-002_Frame_Render_LT16ms":
        return checkFrameRender(screenPath)
    // ... other checks
    }
}
\`\`\`

### Performance Measurement
The system uses various techniques to measure performance:
- File size analysis for first paint estimation
- Code pattern analysis for frame render optimization
- AST parsing for structural checks
- Regular expressions for pattern matching

### Error Handling
Each check returns one of three statuses:
- **PASS**: Check passed successfully
- **FAIL**: Critical requirement not met
- **WARN**: Non-critical issue detected

### Exit Codes
- 0: All checks passed
- 1: Critical checks failed
- 2: Warnings only

## Extending the System

### Adding New Critical Items
1. Add item name to \`critical15Items\` array
2. Implement check function
3. Add case in \`checkItem\` function
4. Update documentation

### Adding New Platforms
1. Create platform-specific runner script
2. Update build script for platform
3. Add platform to distribution packages

### Improving Performance
1. Use AST parsing for more accurate analysis
2. Implement actual performance measurement
3. Add caching for repeated checks
4. Optimize regular expressions

## Testing

### Unit Testing
The system includes unit tests for:
- Check function correctness
- Pattern matching accuracy
- Error handling
- Performance estimation

### Integration Testing
Integration tests verify:
- End-to-end workflow
- CI/CD integration
- Pre-commit hooks
- Output formatting

### Performance Testing
Performance tests ensure:
- Fast execution time
- Low memory usage
- Scalability for large codebases

## Best Practices

### Code Organization
- Keep check functions simple and focused
- Use descriptive function names
- Provide detailed error messages
- Maintain consistent error handling

### Documentation
- Update documentation when adding new checks
- Provide clear examples for each check
- Document edge cases and limitations
- Keep user guides up to date

### Maintenance
- Regular review of critical items
- Monitor false positives/negatives
- Collect feedback from developers
- Adapt to changing requirements

## Future Enhancements

### Planned Features
1. **Actual Performance Measurement**: Replace estimations with real measurements
2. **Visual Regression Testing**: Add screenshot-based testing
3. **Accessibility Testing**: Integration with screen readers
4. **Performance Profiling**: Detailed performance analysis

### Integration Goals
1. **IDE Integration**: Plugin for VSCode, IntelliJ, etc.
2. **Language Server**: LSP integration for real-time feedback
3. **Web Interface**: Web-based dashboard for results
4. **API Access**: REST API for programmatic access

### Advanced Features
1. **Machine Learning**: ML-based anomaly detection
2. **Trend Analysis**: Performance trend tracking
3. **Predictive Analysis**: Predict potential issues
4. **Automated Fixes**: Suggest and apply fixes
EOF

    print_status "PASS" "Documentation created in $doc_dir"
}

# Main function
function main {
    print_status "INFO" "Building UI 15-Item Critical Checklist System"
    print_status "INFO" "Project Root: $PROJECT_ROOT"
    echo "================================================================"
    
    # Create dist directory
    mkdir -p "$PROJECT_ROOT/dist"
    
    # Build for different platforms
    platforms=("linux" "windows" "darwin")
    
    for platform in "${platforms[@]}"; do
        build_binary $platform
        create_package $platform
    done
    
    # Create integration files
    create_ci_files
    create_makefile
    create_documentation
    
    # Create final summary
    print_status "INFO" "Build complete!"
    echo "================================================================"
    print_status "PASS" "Distribution packages created in: $PROJECT_ROOT/dist/"
    print_status "PASS" "CI/CD integration files created"
    print_status "PASS" "Makefile targets added"
    print_status "PASS" "Documentation created in: docs/ui-15-check/"
    
    # Show usage examples
    echo ""
    print_status "INFO" "Usage Examples:"
    echo "  1. Run UI 15-item checklist: make ui-15-check"
    echo "  2. Build checker: make build-ui-15-check"
    echo "  3. Clean artifacts: make clean-ui-15-check"
    echo "  4. Install system-wide: make install-ui-15-check"
}

# Run main function
main "$@"