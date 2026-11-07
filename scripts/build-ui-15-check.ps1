# build-ui-15-check.ps1 - Build the 15-item UI checker for distribution (Windows)
# This script compiles the Go program and creates distribution packages

param(
    [string]$Platform = "all",
    [switch]$Verbose = $false
)

# Set error preference
$ErrorActionPreference = "Stop"

# Function to print colored output
function Write-Status {
    param(
        [string]$Status,
        [string]$Message
    )
    
    if ($PSVersionTable.PSVersion.Major -ge 6) {
        switch ($Status) {
            "PASS" { Write-Host "$($PSStyle.Foreground.Green)✓ $Message$($PSStyle.Reset)" }
            "FAIL" { Write-Host "$($PSStyle.Foreground.Red)✗ $Message$($PSStyle.Reset)" }
            "WARN" { Write-Host "$($PSStyle.Foreground.Yellow)⚠ $Message$($PSStyle.Reset)" }
            "INFO" { Write-Host "$($PSStyle.Foreground.Blue)ℹ $Message$($PSStyle.Reset)" }
        }
    } else {
        switch ($Status) {
            "PASS" { Write-Host "✓ $Message" -ForegroundColor Green }
            "FAIL" { Write-Host "✗ $Message" -ForegroundColor Red }
            "WARN" { Write-Host "⚠ $Message" -ForegroundColor Yellow }
            "INFO" { Write-Host "ℹ $Message" -ForegroundColor Blue }
        }
    }
}

# Function to create distribution package
function Create-Package {
    param(
        [string]$Platform
    )
    
    $outputDir = Join-Path $ProjectRoot "dist\ui-15-check-$Platform"
    
    Write-Status "INFO" "Creating distribution package for $Platform..."
    
    # Create output directory
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
    
    # Copy files
    Copy-Item "$scriptRoot\ui_15_check.go" "$outputDir\" -Force
    Copy-Item "$scriptRoot\ui-15-check.ps1" "$outputDir\" -Force
    
    # Create README
    $readmeContent = @"
# 15-Item Critical UI Checker

This tool provides automated verification of the 15 most critical UI requirements for CodeAuthorityFirst.

## Usage

### Quick Check (All Screens)
```powershell
.\ui-15-check.ps1 -AllScreens
```

### Check Specific Screen
```powershell
.\ui-15-check.ps1 -ScreenPath C:\path\to\screen.go
```

### Options
- `-AllScreens`: Check all UI screens
- `-ScreenPath <path>`: Check specific screen file
- `-OutputDir <dir>`: Output directory for results
- `-Verbose`: Verbose output

## Requirements
- Go 1.19 or later
- PowerShell 5.1 or later
- Windows 10 or later

## Output
- JSON results for each screen
- Summary report in JSON and Markdown formats
- Colored console output

## Exit Codes
- 0: All checks passed
- 1: Critical checks failed
- 2: Warnings only
"@
    
    Set-Content -Path "$outputDir\README.md" -Value $readmeContent
    
    Write-Status "PASS" "Package created: $outputDir"
}

# Function to build Go binary
function Build-Binary {
    param(
        [string]$Platform
    )
    
    $outputDir = Join-Path $ProjectRoot "dist\ui-15-check-$Platform"
    $binaryName = "ui-15-check"
    
    Write-Status "INFO" "Building Go binary for $Platform..."
    
    # Create output directory
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
    
    # Build for the target platform
    if ($Platform -eq "windows") {
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        go build -o "$outputDir\$binaryName.exe" "$scriptRoot\ui_15_check.go"
    } elseif ($Platform -eq "linux") {
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        go build -o "$outputDir\$binaryName" "$scriptRoot\ui_15_check.go"
        chmod +x "$outputDir\$binaryName"
    } elseif ($Platform -eq "darwin") {
        $env:GOOS = "darwin"
        $env:GOARCH = "amd64"
        go build -o "$outputDir\$binaryName" "$scriptRoot\ui_15_check.go"
        chmod +x "$outputDir\$binaryName"
    }
    
    Write-Status "PASS" "Binary built: $outputDir\$binaryName"
}

# Function to create CI/CD integration files
function Create-CIFiles {
    $ciDir = Join-Path $ProjectRoot ".github\workflows"
    
    Write-Status "INFO" "Creating CI/CD integration files..."
    
    # Create UI check workflow
    New-Item -ItemType Directory -Path $ciDir -Force | Out-Null
    
    $workflowContent = @"
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
"@
    
    Set-Content -Path "$ciDir\ui-15-check.yml" -Value $workflowContent
    
    # Create pre-commit hook
    $gitHooksDir = Join-Path $ProjectRoot ".git\hooks"
    New-Item -ItemType Directory -Path $gitHooksDir -Force | Out-Null
    
    $precommitContent = @"
#!/bin/bash

# Pre-commit hook for UI 15-item checklist
# This runs the UI checker before committing changes

echo "Running UI 15-item checklist..."

# Check if UI checker exists
if [ -f "scripts/ui-15-check.sh" ]; then
    chmod +x scripts/ui-15-check.sh
    
    # Run UI checker on changed screens
    changed_files=$(git diff --cached --name-only --diff-filter=ACM | grep "pkg/ui/screens.*\.go$" || true)
    
    if [ -n "$changed_files" ]; then
        echo "Checking changed UI screens..."
        for file in $changed_files; do
            if [ -f "$file" ]; then
                echo "Checking $file..."
                ./scripts/ui-15-check.sh --screen "$file"
                if [ $? -ne 0 ]; then
                    echo "❌ UI checks failed for $file"
                    echo "Please fix the issues before committing."
                    exit 1
                fi
            fi
        done
    fi
fi

echo "✅ UI 15-item checklist passed"
"@
    
    Set-Content -Path "$gitHooksDir\pre-commit" -Value $precommitContent
    
    # Make pre-commit hook executable
    if (Get-Command chmod -ErrorAction SilentlyContinue) {
        chmod +x "$gitHooksDir\pre-commit"
    }
    
    Write-Status "PASS" "CI/CD integration files created"
}

# Function to create PowerShell build targets
function Create-PowerShellTargets {
    $targetsFile = Join-Path $ProjectRoot "build-ui-15-check.ps1"
    
    Write-Status "INFO" "Creating PowerShell build targets..."
    
    $targetsContent = @"
# PowerShell build targets for UI 15-item checklist

function Invoke-UI15Check {
    param(
        [string]$ScreenPath = "",
        [switch]$AllScreens = $false,
        [string]$OutputDir = "",
        [switch]$Verbose = $false
    )
    
    Write-Host "Running UI 15-item checklist..." -ForegroundColor Blue
    
    if ($AllScreens) {
        & "$PSScriptRoot\ui-15-check.ps1" -AllScreens -Verbose:$Verbose
    } elseif ($ScreenPath) {
        & "$PSScriptRoot\ui-15-check.ps1" -ScreenPath $ScreenPath -Verbose:$Verbose
    } else {
        Write-Host "Please specify either -AllScreens or -ScreenPath" -ForegroundColor Red
        exit 1
    }
}

function Invoke-BuildUI15Check {
    param(
        [string]$Platform = "all",
        [switch]$Verbose = $false
    )
    
    Write-Host "Building UI 15-item checker..." -ForegroundColor Blue
    
    & "$PSScriptRoot\build-ui-15-check.ps1" -Platform $Platform -Verbose:$Verbose
}

function Invoke-CleanUI15Check {
    Write-Host "Cleaning UI 15-item checker artifacts..." -ForegroundColor Blue
    
    $distDir = Join-Path $ProjectRoot "dist"
    if (Test-Path $distDir) {
        Remove-Item $distDir -Recurse -Force
    }
    
    $summaryFiles = Get-ChildItem $ProjectRoot "ui_15_check_summary_*.json" -ErrorAction SilentlyContinue
    if ($summaryFiles) {
        $summaryFiles | Remove-Item -Force
    }
    
    $summaryFiles = Get-ChildItem $ProjectRoot "ui_15_check_summary_*.md" -ErrorAction SilentlyContinue
    if ($summaryFiles) {
        $summaryFiles | Remove-Item -Force
    }
    
    Write-Host "✅ Cleaned UI 15-item checker artifacts" -ForegroundColor Green
}

function Install-UI15Check {
    Write-Host "Installing UI 15-item checker..." -ForegroundColor Blue
    
    $localBin = Join-Path $env:USERPROFILE ".local\bin"
    New-Item -ItemType Directory -Path $localBin -Force | Out-Null
    
    $binaryPath = Join-Path $localBin "ui-15-check.ps1"
    Copy-Item "$PSScriptRoot\ui-15-check.ps1" $binaryPath -Force
    
    Write-Host "✅ Installed to $binaryPath" -ForegroundColor Green
    Write-Host "📝 Add $localBin to your PATH if not already there" -ForegroundColor Yellow
}

# Export functions
Export-ModuleMember -Function @(
    "Invoke-UI15Check",
    "Invoke-BuildUI15Check", 
    "Invoke-CleanUI15Check",
    "Install-UI15Check"
)
"@
    
    Set-Content -Path $targetsFile -Value $targetsContent
    
    Write-Status "PASS" "PowerShell build targets created"
}

# Function to create documentation
function Create-Documentation {
    $docDir = Join-Path $ProjectRoot "docs\ui-15-check"
    
    Write-Status "INFO" "Creating documentation..."
    
    # Create documentation directory
    New-Item -ItemType Directory -Path $docDir -Force | Out-Null
    
    # Create user guide
    $userGuide = @"
# UI 15-Item Critical Checklist User Guide (Windows)

## Overview

The UI 15-Item Critical Checklist is an automated verification system that ensures all UI implementations meet the most critical requirements for performance, accessibility, and code quality.

## Quick Start

### 1. Run on All Screens
```powershell
.\scripts\ui-15-check.ps1 -AllScreens
```

### 2. Run on Specific Screen
```powershell
.\scripts\ui-15-check.ps1 -ScreenPath C:\path\to\screen.go
```

### 3. Build Checker
```powershell
.\scripts\build-ui-15-check.ps1 -Platform windows
```

### 4. Using PowerShell Functions
```powershell
# Import build targets
. .\build-ui-15-check.ps1

# Run UI 15-item checklist
Invoke-UI15Check -AllScreens

# Build checker
Invoke-BuildUI15Check -Platform windows

# Clean artifacts
Invoke-CleanUI15Check

# Install system-wide
Install-UI15Check
```

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
2. **Execution Policy**: Run \`Set-ExecutionPolicy RemoteSigned -Scope Process\`
3. **Path Issues**: Use absolute paths for screen files

### Debug Mode
Run with verbose output:
```powershell
.\scripts\ui-15-check.ps1 -AllScreens -Verbose
```

## Output

### Console Output
Colored console output showing:
- ✅ Passed checks
- ❌ Failed checks
- ⚠ Warnings

### JSON Results
Detailed results in JSON format for each screen:
```json
{
  "timestamp": "2026-11-07T10:00:00Z",
  "item": "INV-UI-001_First_Paint_LT100ms",
  "status": "PASS",
  "message": "First paint time acceptable"
}
```

### Summary Reports
- JSON summary: \`ui_15_check_summary_*.json\`
- Markdown summary: \`ui_15_check_summary_*.md\`
"@
    
    Set-Content -Path "$docDir\USER_GUIDE_WINDOWS.md" -Value $userGuide
    
    # Create developer guide
    $developerGuide = @"
# UI 15-Item Critical Checklist Developer Guide (Windows)

## Architecture

### Core Components
1. **ui_15_check.go**: Main checker implementation
2. **ui-15-check.ps1**: Windows runner script
3. **build-ui-15-check.ps1**: Build script for distribution
4. **PowerShell build targets**: Reusable functions

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
Each check is implemented as a separate function that returns a \`UI15CheckResult\` object:

```powershell
function Check-Item {
    param(
        [string]$Item,
        [string]$ScreenPath
    )
    
    switch ($Item) {
        "INV-UI-001_First_Paint_LT100ms" {
            return Check-FirstPaint $ScreenPath
        }
        "INV-UI-002_Frame_Render_LT16ms" {
            return Check-FrameRender $ScreenPath
        }
        # ... other checks
    }
}
```

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

## PowerShell Integration

### Build Targets
The system includes PowerShell build targets for easy integration:

```powershell
# Import build targets
. .\build-ui-15-check.ps1

# Use functions
Invoke-UI15Check -AllScreens
Invoke-BuildUI15Check -Platform windows
Invoke-CleanUI15Check
Install-UI15Check
```

### CI/CD Integration
The system integrates with GitHub Actions and other CI/CD systems:

```yaml
name: UI 15-Item Critical Checklist
on:
  pull_request:
    paths:
      - 'pkg/ui/screens/**/*.go'
      - 'scripts/ui-15-check*'
jobs:
  ui-check:
    runs-on: windows-latest
    steps:
    - uses: actions/checkout@v4
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - name: Run UI 15-item checklist
      run: |
        .\scripts\ui-15-check.ps1 -AllScreens
```

## Extending the System

### Adding New Critical Items
1. Add item name to \`critical15Items\` array
2. Implement check function
3. Add case in \`Check-Item\` function
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

## Windows-Specific Considerations

### PowerShell Version Compatibility
- PowerShell 5.1: Basic functionality
- PowerShell 6+: Enhanced color support
- PowerShell 7+: Best experience with full feature set

### Path Handling
- Use \`Join-Path\` for cross-platform path construction
- Handle both forward and backward slashes
- Normalize paths for consistent comparison

### Execution Policy
- May need to set execution policy for running scripts
- Use \`Set-ExecutionPolicy RemoteSigned -Scope Process\` for current session

### Environment Variables
- Use \`$env:GOOS\` and \`$env:GOARCH\` for cross-platform builds
- Handle PATH environment for Go installation detection
"@
    
    Set-Content -Path "$docDir\DEVELOPER_GUIDE_WINDOWS.md" -Value $developerGuide
    
    Write-Status "PASS" "Documentation created in $docDir"
}

# Main function
function Main {
    Write-Status "INFO" "Building UI 15-Item Critical Checklist System (Windows)"
    Write-Status "INFO" "Project Root: $ProjectRoot"
    Write-Host "================================================================"
    
    # Create dist directory
    New-Item -ItemType Directory -Path "$ProjectRoot\dist" -Force | Out-Null
    
    # Build for different platforms
    if ($Platform -eq "all") {
        $platforms = @("windows", "linux", "darwin")
    } else {
        $platforms = @($Platform)
    }
    
    foreach ($plat in $platforms) {
        Build-Binary $plat
        Create-Package $plat
    }
    
    # Create integration files
    Create-CIFiles
    Create-PowerShellTargets
    Create-Documentation
    
    # Create final summary
    Write-Status "INFO" "Build complete!"
    Write-Host "================================================================"
    Write-Status "PASS" "Distribution packages created in: $ProjectRoot\dist\"
    Write-Status "PASS" "CI/CD integration files created"
    Write-Status "PASS" "PowerShell build targets created"
    Write-Status "PASS" "Documentation created in: docs\ui-15-check\"
    
    # Show usage examples
    Write-Host ""
    Write-Status "INFO" "Usage Examples:"
    Write-Host "  1. Run UI 15-item checklist: .\scripts\ui-15-check.ps1 -AllScreens"
    Write-Host "  2. Build checker: .\scripts\build-ui-15-check.ps1 -Platform windows"
    Write-Host "  3. Clean artifacts: Invoke-CleanUI15Check"
    Write-Host "  4. Install system-wide: Install-UI15Check"
}

# Initialize variables
$scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = Split-Path -Parent $scriptRoot

# Run main function
Main