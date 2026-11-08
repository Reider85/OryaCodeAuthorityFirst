# Real-World Testing Infrastructure
**Phase 1: Setup - Monitoring and Test Environment**

**Date**: 2026-11-08  
**Project**: CodeAuthorityFirst (caf)  
**Phase**: MVP → ALPHA Transition - Real-World Testing Setup

---

## Infrastructure Overview

This directory contains the testing infrastructure for real-world validation of exit criteria 1,2,3,9,10,11,12 §12.

### Directory Structure
```
scenarios/
├── monitoring/           # Monitoring and logging infrastructure
├── environments/         # Test environment configurations
├── repositories/        # Test repositories for different scenarios
├── agents/              # AI agent configurations and testing
├── validation/          # Validation scripts and criteria checks
└── reports/             # Report generation and analysis
```

---

## Monitoring Infrastructure

### 1. Session Tracking Monitor
**Purpose**: Track session creation and metadata collection
**Files**: 
- `monitoring/session-tracker.go` - Session monitoring logic
- `monitoring/session-logger.go` - Session logging utilities
- `monitoring/session-analyzer.go` - Session analysis tools

**Features**:
- Real-time session creation monitoring
- Metadata validation and completeness checking
- Performance tracking for session operations
- Error detection and recovery monitoring

### 2. Hook Performance Monitor
**Purpose**: Monitor universal hook execution across agents
**Files**:
- `monitoring/hook-monitor.go` - Hook execution monitoring
- `monitoring/performance-tracker.go` - Performance metrics collection
- `monitoring/error-collector.go` - Error log collection

**Features**:
- Hook execution timing and success rate tracking
- Cross-agent compatibility monitoring
- Performance benchmarking against requirements
- Error pattern analysis

### 3. LSP Performance Monitor
**Purpose**: Monitor LSP fallback performance with large repositories
**Files**:
- `monitoring/lsp-monitor.go` - LSP performance monitoring
- `monitoring/memory-tracker.go` - Memory usage monitoring
- `monitoring/symbol-accuracy.go` - Symbol accuracy validation

**Features**:
- Symbol graph generation performance tracking
- Memory usage monitoring and limits
- Symbol accuracy validation against expected results
- Large repository handling validation

---

## Test Environment Configurations

### Environment Types
1. **Development Environment**: VS Code, IntelliJ, Vim/Emacs
2. **Platform Environment**: Windows, macOS, Linux
3. **Repository Size**: Small (<100 files), Medium (100-500 files), Large (1000+ files)

### Configuration Files
- `environments/vs-code-config.json` - VS Code testing configuration
- `environments/intellij-config.json` - IntelliJ testing configuration
- `environments/vim-config.json` - Vim testing configuration
- `environments/windows-config.json` - Windows environment setup
- `environments/linux-config.json` - Linux environment setup

---

## Test Repository Setup

### Repository Types
1. **Small Repository**: Basic Go project (<100 files)
2. **Medium Repository**: Standard Go project (100-500 files)
3. **Large Repository**: Enterprise Go project (1000+ files)

### Repository Structure
```
repositories/
├── small-go-project/     # Basic Go application
├── medium-go-project/   # Standard Go application
└── large-go-project/    # Enterprise Go application
```

Each repository includes:
- Complete caf initialization
- Pre-configured for different agent types
- Ready for session tracking
- Performance monitoring hooks

---

## Agent Configuration

### Supported Agents
1. **opencode** - Tier 1 native adapter
2. **claude-code** - Tier 1 native adapter
3. **universal-lsp** - Tier 4 fallback adapter

### Agent Configuration Files
- `agents/opencode-config.toml` - opencode configuration
- `agents/claude-code-config.json` - claude-code configuration
- `agents/universal-lsp-config.yml` - Universal LSP configuration

### Agent Testing Scripts
- `agents/test-opencode-hooks.ps1` - opencode hook testing
- `agents/test-claude-code-hooks.ps1` - claude-code hook testing
- `agents/test-universal-lsp.ps1` - Universal LSP testing

---

## Validation Infrastructure

### Criteria Validation Scripts
- `validation/criterion1-7day-test.ps1` - 7-day pilot testing
- `validation/criterion2-multiagent-test.ps1` - Multi-agent hook testing
- `validation/criterion3-lsp-test.ps1` - LSP large repository testing
- `validation/criterion9-ui-test.ps1` - UI 5-state audit
- `validation/criterion10-checklist.ps1` - UI checklist validation
- `validation/criterion11-focus-test.ps1` - Focus visibility testing
- `validation/criterion12-contrast-test.ps1` - High-contrast testing

### Validation Features
- Automated criteria checking
- Performance requirement validation
- Error detection and reporting
- Success/failure metrics collection

---

## Report Generation

### Report Templates
- `templates/daily-report.md` - Daily testing progress report
- `templates/weekly-summary.md` - Weekly testing summary
- `templates/criterion-report.md` - Criterion-specific report
- `templates/final-report.md` - Final testing report

### Analysis Tools
- `reports/analyze-sessions.ps1` - Session analysis script
- `reports/analyze-performance.ps1` - Performance analysis script
- `reports/analyze-errors.ps1` - Error analysis script
- `reports/generate-summary.ps1` - Summary report generator

---

## Setup Instructions

### Prerequisites
- Go 1.22+ (build-time only)
- git (required for hooks)
- opencode (for Tier 1 testing)
- claude-code (for Tier 1 testing)
- gopls (for LSP fallback testing)

### Installation Steps
1. **Build caf binary**:
   ```bash
   go build -o caf ./cmd/caf
   ```

2. **Initialize monitoring infrastructure**:
   ```bash
   cd real-world-testing/scenarios
   .\setup-monitoring.ps1
   ```

3. **Prepare test repositories**:
   ```bash
   .\setup-repositories.ps1
   ```

4. **Configure agent testing**:
   ```bash
   .\configure-agents.ps1
   ```

5. **Start monitoring**:
   ```bash
   .\start-monitoring.ps1
   ```

### Monitoring Commands
- `.\monitor-sessions.ps1` - Monitor session tracking
- `.\monitor-hooks.ps1` - Monitor hook execution
- `.\monitor-lsp.ps1` - Monitor LSP performance
- `.\generate-reports.ps1` - Generate testing reports

---

## Monitoring Dashboard

### Real-time Metrics
- Session creation rate
- Hook success rate
- Performance metrics
- Error detection
- Agent compatibility

### Alert System
- Performance threshold alerts
- Error rate alerts
- Session completeness alerts
- Memory usage alerts

### Data Collection
- Session metadata logs
- Hook execution logs
- Performance data
- Error logs
- User feedback

---

## Next Steps

1. **Complete Infrastructure Setup**: Build monitoring and validation tools
2. **Prepare Test Environments**: Setup VS Code, IntelliJ, Vim/Emacs
3. **Configure Agent Testing**: Prepare opencode/claude-code testing
4. **Setup Test Repositories**: Create small, medium, and large repositories
5. **Initialize Monitoring**: Start data collection and monitoring
6. **Begin Testing**: Execute real-world testing scenarios

---

## Documentation

- **Testing Guide**: `docs/testing-guide.md`
- **Monitoring Guide**: `docs/monitoring-guide.md`
- **Agent Configuration**: `docs/agent-config.md`
- **Repository Setup**: `docs/repository-setup.md`
- **Validation Scripts**: `docs/validation-scripts.md`
- **Report Templates**: `docs/report-templates.md`

---

## Quality Assurance

### Testing Standards
- All scripts must be idempotent
- Comprehensive error handling
- Detailed logging
- Performance monitoring
- User feedback collection

### Validation Standards
- Automated criteria validation
- Performance requirement checking
- Error detection and reporting
- Success metric tracking
- Continuous improvement

---

This infrastructure provides comprehensive monitoring and testing capabilities for real-world validation of all exit criteria. The system is designed to be robust, scalable, and provide detailed insights into the performance and reliability of the CodeAuthorityFirst system.