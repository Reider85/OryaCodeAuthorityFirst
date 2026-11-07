# UI 15-Item Critical Checklist System

An automated verification system that ensures all UI implementations meet the 15 most critical requirements for performance, accessibility, and code quality in CodeAuthorityFirst.

## 🎯 Overview

This system provides automated checks for the 15 most critical UI requirements that cannot be fully verified by automated tests alone. It focuses on:

- **Performance**: First paint, frame render, keystroke response
- **Accessibility**: Keyboard-only operation, focus visibility, symbol-text pairing
- **Architecture**: ELM pattern, state management
- **Code Quality**: File structure, testing, documentation
- **User Experience**: Error handling, keyboard shortcuts

## 📋 15 Critical Items

### Performance Requirements (3 items)
1. **INV-UI-001**: First paint <100ms
2. **INV-UI-002**: Frame render <16ms (60fps)
3. **INV-UI-003**: Keystroke response <50ms

### Accessibility Requirements (3 items)
4. **INV-UI-009**: 100% keyboard-only operation
5. **INV-UI-010**: Help overlay functionality
6. **Accessibility**: Focus visibility indicators
7. **Accessibility**: Symbol-text pairing

### Architecture Requirements (2 items)
8. **Screen Structure**: ELM pattern compliance (Init/Update/View)
9. **State Management**: ui.StateMachine usage

### Code Quality Requirements (4 items)
10. **Performance**: Lipgloss pre-computed styles
11. **File Structure**: Required files present (test files, components)
12. **Testing**: Unit tests exist and complete
13. **Documentation**: Header comment completeness

### Error Handling & UX (3 items)
14. **Error Handling**: Consistent error states with recovery options
15. **Keyboard Map**: Universal shortcuts implementation

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or later
- Bash shell (Linux/macOS) or PowerShell (Windows)
- CodeAuthorityFirst project

### Installation
```bash
# Clone the project
git clone <repository-url>
cd OryaCodeAuthorityFirst

# Build the system
chmod +x scripts/build-ui-15-check.sh
./scripts/build-ui-15-check.sh

# Or on Windows
.\scripts\build-ui-15-check.ps1 -Platform windows
```

### Basic Usage

#### Check All Screens
```bash
# Linux/macOS
./scripts/ui-15-check.sh --all

# Windows
.\scripts\ui-15-check.ps1 -AllScreens
```

#### Check Specific Screen
```bash
# Linux/macOS
./scripts/ui-15-check.sh --screen pkg/ui/screens/dashboard.go

# Windows
.\scripts\ui-15-check.ps1 -ScreenPath "pkg\ui\screens\dashboard.go"
```

#### Using Make Targets
```bash
# Run on all screens
make ui-15-check

# Run on specific screen
make ui-15-check-screen SCREEN=pkg/ui/screens/dashboard.go

# Build the checker
make build-ui-15-check

# Clean artifacts
make clean-ui-15-check

# Install system-wide
make install-ui-15-check
```

## 📊 Output

### Console Output
The system provides colored console output:
- ✅ **PASS**: Check passed successfully
- ❌ **FAIL**: Critical requirement not met
- ⚠ **WARN**: Non-critical issue detected

### JSON Results
Each screen generates detailed JSON results:
```json
{
  "timestamp": "2026-11-07T10:00:00Z",
  "item": "INV-UI-001_First_Paint_LT100ms",
  "description": "First paint <100ms",
  "status": "PASS",
  "message": "File size reasonable for first paint",
  "details": "File size: 4567 bytes"
}
```

### Summary Reports
- **JSON Summary**: `ui_15_check_summary_YYYYMMDD_HHMMSS.json`
- **Markdown Summary**: `ui_15_check_summary_YYYYMMDD_HHMMSS.md`

## 🔧 Integration

### CI/CD Integration
The system is integrated into the CI/CD pipeline and runs automatically on:
- Pull requests affecting UI screens
- Push to main/master branches

### Pre-commit Hooks
A pre-commit hook is installed that runs the UI checker on changed UI files before commit.

### IDE Integration
You can integrate the system into your IDE:

#### VS Code
Create a `.vscode/tasks.json`:
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "UI 15-Item Check",
      "type": "shell",
      "command": "./scripts/ui-15-check.sh",
      "args": ["--all"],
      "group": "build",
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    }
  ]
}
```

## 📁 Project Structure

```
scripts/
├── ui_15_check.go              # Main checker implementation
├── ui-15-check.sh              # Linux/macOS runner script
├── ui-15-check.ps1            # Windows runner script
├── build-ui-15-check.sh        # Linux/macOS build script
└── build-ui-15-check.ps1      # Windows build script

docs/ui-15-check/
├── USER_GUIDE.md             # Comprehensive user guide
├── DEVELOPER_GUIDE.md        # Developer guide
├── USER_GUIDE_WINDOWS.md    # Windows-specific guide
└── DEVELOPER_GUIDE_WINDOWS.md # Windows-specific dev guide

dist/
├── ui-15-check-linux/         # Linux distribution package
├── ui-15-check-windows/       # Windows distribution package
└── ui-15-check-darwin/       # macOS distribution package
```

## 🛠️ Advanced Usage

### Custom Output Directory
```bash
# Linux/macOS
./scripts/ui-15-check.sh --all --output ./results/

# Windows
.\scripts\ui-15-check.ps1 -AllScreens -OutputDir "./results/"
```

### Verbose Mode
```bash
# Linux/macOS
./scripts/ui-15-check.sh --all --verbose

# Windows
.\scripts\ui-15-check.ps1 -AllScreens -Verbose
```

### Programmatic Usage
```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/your-org/ui-15-check"
)

func main() {
    checker := ui15check.NewChecker()
    
    // Check a specific screen
    result, err := checker.CheckScreen("pkg/ui/screens/dashboard.go")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Screen: %s\n", result.ScreenName)
    fmt.Printf("Passed: %d/%d\n", result.PassedChecks, result.TotalChecks)
    
    if result.FailedChecks > 0 {
        fmt.Println("Failed checks:")
        for _, check := range result.FailedChecks {
            fmt.Printf("  - %s: %s\n", check.Item, check.Message)
        }
    }
}
```

## 🔍 Check Details

### INV-UI-001: First Paint <100ms
**What it checks**: Screen initialization performance
**Verification**: File size analysis, initialization patterns
**Common issues**: Heavy data loading, expensive calculations

### INV-UI-002: Frame Render <16ms
**What it checks**: Rendering performance for smooth animations
**Verification**: Virtualization patterns, list optimization
**Common issues**: Large lists without virtualization, expensive string operations

### INV-UI-003: Keystroke Response <50ms
**What it checks**: Input handling performance
**Verification**: Update method analysis, blocking operations
**Common issues**: Async operations blocking main thread

### INV-UI-009: 100% Keyboard-Only Operation
**What it checks**: Full keyboard accessibility
**Verification**: Mouse operation detection, keyboard shortcut coverage
**Common issues**: Mouse-only interactions, missing keyboard shortcuts

### INV-UI-010: Help Overlay
**What it checks**: Help system functionality
**Verification**: Help overlay implementation, documentation
**Common issues**: Missing help screens, incomplete documentation

## 🎨 Customization

### Adding New Critical Items
1. Update `critical15Items` array in `ui_15_check.go`
2. Implement check function
3. Add case in `checkItem` function
4. Update documentation

### Modifying Performance Thresholds
```go
// In ui_15_check.go
const (
    FirstPaintThresholdMs = 100
    FrameRenderThresholdMs = 16
    KeystrokeThresholdMs = 50
)
```

### Custom Error Messages
```go
// Customize error messages in check functions
return UI15CheckResult{
    Item:        "Your_Item",
    Description: "Your description",
    Status:      "FAIL",
    Message:     "Custom error message",
    Details:     "Additional details",
}
```

## 🐛 Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Fix executable permissions
chmod +x scripts/ui-15-check.sh
chmod +x scripts/build-ui-15-check.sh
```

#### Go Version Issues
```bash
# Check Go version
go version

# Install required version
# See https://golang.org/doc/install
```

#### Path Issues
```bash
# Use absolute paths
./scripts/ui-15-check.sh --screen /absolute/path/to/screen.go
```

### Debug Mode
```bash
# Enable verbose output
./scripts/ui-15-check.sh --all --verbose

# Debug build process
./scripts/build-ui-15-check.sh --verbose
```

### Performance Issues
If the checker is slow:
1. Use AST parsing instead of regex for complex checks
2. Implement caching for repeated checks
3. Optimize regular expressions

## 📈 Performance Monitoring

The system includes performance monitoring for the checker itself:
- Execution time measurement
- Memory usage tracking
- Scalability analysis

### Benchmarking
```bash
# Benchmark checker performance
time ./scripts/ui-15-check.sh --all

# Benchmark specific screen
time ./scripts/ui-15-check.sh --screen pkg/ui/screens/dashboard.go
```

## 🔗 Related Tools

### Existing UI Review System
This complements the existing comprehensive UI review system:
- **UI_SELF_REVIEW_CHECKLIST.md**: Full 619-line checklist
- **UI_SELF_REVIEW_SYSTEM.md**: System documentation
- **UI_IMPLEMENTATION_EXAMPLES.md**: Implementation examples

### Performance Monitoring
- **Performance-enhanced CI**: Advanced CI with performance monitoring
- **Performance budgets**: Strict performance requirements
- **Regression detection**: Automated performance regression detection

## 🤝 Contributing

### Development Setup
```bash
# Clone the repository
git clone <repository-url>
cd OryaCodeAuthorityFirst

# Install development dependencies
go mod tidy

# Run tests
go test ./scripts/

# Build the system
./scripts/build-ui-15-check.sh
```

### Guidelines
1. Follow the existing code style
2. Add comprehensive tests for new checks
3. Update documentation when adding new features
4. Consider backward compatibility when making changes

### Testing
```bash
# Run unit tests
go test ./scripts/

# Run integration tests
./scripts/ui-15-check.sh --all

# Run performance benchmarks
go test -bench=. ./scripts/
```

## 📚 Documentation

### Core Documentation
- [User Guide](docs/ui-15-check/USER_GUIDE.md)
- [Developer Guide](docs/ui-15-check/DEVELOPER_GUIDE.md)
- [Windows User Guide](docs/ui-15-check/USER_GUIDE_WINDOWS.md)
- [Windows Developer Guide](docs/ui-15-check/DEVELOPER_GUIDE_WINDOWS.md)

### Related Documentation
- [UI Self-Review Checklist](docs/UI_SELF_REVIEW_CHECKLIST.md)
- [UI Self-Review System](docs/UI_SELF_REVIEW_SYSTEM.md)
- [UI Implementation Examples](docs/UI_IMPLEMENTATION_EXAMPLES.md)

## 📞 Support

### Issues
Report bugs and request features:
- GitHub Issues: [Create an issue](https://github.com/your-org/repo/issues)
- Discussions: [Ask questions](https://github.com/your-org/repo/discussions)

### Community
- Join discussions about UI quality standards
- Share feedback and suggestions
- Contribute improvements to the system

## 🏆 Success Criteria

The system is successful when:
1. **All UI implementations** pass the 15 critical checks
2. **Developers** use the system regularly during development
3. **CI/CD pipeline** enforces the critical requirements
4. **Code quality** improves consistently
5. **User experience** meets accessibility and performance standards

---

*This system is part of CodeAuthorityFirst's commitment to exceptional UI quality and user experience.*