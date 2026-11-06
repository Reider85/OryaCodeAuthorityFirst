# Traceability System (INV-PROMPT-MVP-005)

## Overview

The Traceability System ensures that all files in the repository have proper traceability links in their commit history. This system validates that every file modification is properly traced to relevant requirements (FR-XXX, ADR-XXX, INV-XXX) in commit messages.

## Features

- **Automated pre-commit validation** - Validates staged files before commits
- **Configurable file patterns** - Include/exclude specific file types
- **Multiple output formats** - CLI, JSON, and Markdown reports
- **Grace period for new files** - Allows time for initial setup
- **Advisory-only in MVP** - Never blocks commits (logs warnings only)
- **CI/CD integration** - Automated validation in GitHub Actions

## Configuration

The traceability system is configured via `.caf/traceability.yaml`:

```yaml
# File patterns to include in validation
include_patterns:
  - "**/*.go"
  - "**/*.md"
  - "**/*.yaml"
  - "**/*.yml"

# File patterns to exclude from validation
exclude_patterns:
  - "**/vendor/**"
  - "**/node_modules/**"
  - "**/*.test.go"

# Requirement prefixes to validate
requirement_prefixes:
  - "FR-"
  - "ADR-"
  - "INV-"

# Number of commits to look back
lookback_commits: 10

# Grace commits for new files
new_file_grace_commits: 5

# Reporting level (advisory, warning, error)
reporting_level: "advisory"

# Output format (cli, json, markdown)
output_format: "cli"

# Enable/disable validation
enabled: true

# Log level
log_level: "info"
```

## Usage

### CLI Commands

#### Check Staged Files
```bash
caf traceability check
caf traceability check --repo /path/to/repo
caf traceability check --verbose
```

#### Validate All Tracked Files
```bash
caf traceability validate
caf traceability validate --repo /path/to/repo
```

#### Generate Reports
```bash
# CLI format (default)
caf traceability report

# JSON format
caf traceability report --format json

# Markdown format
caf traceability report --format markdown

# With verbose logging
caf traceability report --verbose
```

### Pre-commit Integration

The traceability system is automatically integrated into the pre-commit hook:

```bash
# Pre-commit hook will automatically validate traceability
git commit -m "feat: add new feature

Traces to: FR-001, ADR-002"
```

### CI/CD Integration

The system includes GitHub Actions for automated traceability validation:

#### Automatic Validation
- Runs on push to main/develop branches
- Runs on pull requests
- Generates compliance reports
- Comments on PRs with results

#### Manual Validation
- Triggered workflow for manual validation
- Configurable parameters (path, format, verbose)
- Generates detailed reports

## Commit Message Format

For proper traceability, commit messages should include a "Traces to:" footer:

```bash
git commit -m "feat: implement user authentication

Add JWT-based authentication system with refresh tokens.
Implemented with bcrypt password hashing and secure token storage.

Traces to: FR-AUTH-001, ADR-SEC-003"
```

### Valid Requirement Formats

- **FR-XXX** - Functional Requirements (e.g., FR-PROV-001)
- **ADR-XXX** - Architectural Decision Records (e.g., ADR-003)
- **INV-XXX** - Invariants and constraints (e.g., INV-PROMPT-MVP-005)

### Bypass Options

#### Trivial Commits
For minor changes, use the CAF-Trivial bypass:

```bash
git commit -m "fix: typo in README

CAF-Trivial: true"
```

#### Session Bypass
CAF session validation can be bypassed with proper session tracking.

## Reporting

### CLI Output
```
📋 Traceability Validation Report
==============================

✅ src/main.go
❌ src/auth.go (1 issues)

📊 Summary
---------
Total files checked: 2
Passed: 1
Failed: 1
Total issues: 1

⚠️  Files with traceability issues:

src/auth.go:
  - Missing traceability in commit abc1234
```

### JSON Output
```json
{
  "summary": {
    "total_files": 2,
    "passed": 1,
    "failed": 1,
    "total_issues": 1
  },
  "results": [
    {
      "file": "src/main.go",
      "has_traceability": true,
      "traceability_issues": []
    },
    {
      "file": "src/auth.go",
      "has_traceability": false,
      "traceability_issues": ["Missing traceability in commit abc1234"]
    }
  ]
}
```

### Markdown Output
```markdown
# Traceability Validation Report

## Summary
- **Total Files**: 2
- **Passed**: 1
- **Failed**: 1
- **Total Issues**: 1

## Files with Traceability Issues

### src/auth.go
- Missing traceability in commit abc1234
```

## Compliance Thresholds

- **MVP**: Advisory-only (no blocking)
- **ALPHA**: 80% compliance threshold
- **BETA**: 90% compliance threshold
- **PROD**: 95% compliance threshold

## Troubleshooting

### Common Issues

#### No Traceability Issues Found
- Check file patterns in configuration
- Verify requirement prefixes match your project's format
- Ensure files have commit history

#### False Positives for New Files
- New files get a grace period (default: 5 commits)
- Check if files are properly tracked by git
- Verify file patterns don't exclude the files

#### Performance Issues
- Reduce `lookback_commits` in configuration
- Exclude unnecessary file patterns
- Use `--verbose` flag for debugging

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
caf traceability check --verbose
```

### Configuration Validation

Validate your configuration file:

```bash
# Check if config is valid
caf traceability check

# Manually validate config
cat .caf/traceability.yaml | yamllint
```

## Integration with Existing Workflow

The traceability system integrates seamlessly with the existing CAF workflow:

1. **Pre-commit**: Automatically validates staged files
2. **Session Tracking**: Works with existing CAF session system
3. **Configuration**: Uses existing .caf directory structure
4. **Logging**: Integrates with existing logging infrastructure

## Migration Guide

### For Existing Projects

1. **Add traceability configuration**:
   ```bash
   cp .caf/traceability.yaml.example .caf/traceability.yaml
   ```

2. **Update commit messages** to include "Traces to:" footers

3. **Gradual adoption**:
   - Start with advisory-only mode
   - Use the grace period for new files
   - Gradually improve traceability coverage

### For New Projects

1. **Initialize with traceability**:
   ```bash
   caf init
   ```

2. **Configure file patterns** for your project type

3. **Set up CI/CD integration** for automated validation

## Future Enhancements

### ALPHA Phase
- **Blocking mode** - Option to block commits with poor traceability
- **Requirement validation** - Verify requirements exist in documentation
- **Automated suggestions** - Suggest relevant requirements for changes

### BETA Phase
- **Machine learning** - Auto-suggest requirement links
- **Integration with issue trackers** - Link to GitHub issues/JIRA
- **Advanced reporting** - Trend analysis and metrics

### PROD Phase
- **Real-time validation** - Validation during development
- **Team-wide enforcement** - Organization-wide policies
- **Advanced analytics** - Traceability health dashboard

## Contributing

To contribute to the traceability system:

1. **Fork the repository**
2. **Create a feature branch**
3. **Add tests for new functionality**
4. **Update documentation**
5. **Submit a pull request**

## License

This project is licensed under the Apache License 2.0.