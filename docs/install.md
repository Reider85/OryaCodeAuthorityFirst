# Installation Guide

## Requirements

- **Go 1.22+** - Build-time only
- **git** - Required for pre-commit hooks
- Optional: `gopls` (for LSP support)

## Installation Methods

### From Source

```bash
git clone https://github.com/Reider85/OryaCodeAuthorityFirst.git
cd OryaCodeAuthorityFirst
go build -o caf ./cmd/caf
sudo mv caf /usr/local/bin/
```

### Verify Installation

```bash
caf --version
caf doctor
```

## Post-Installation Setup

### Initialize Your Repository

```bash
cd your-project
caf init --auto-detect
```

This will:
- Create `.caf/` directory structure
- Install universal pre-commit hook
- Detect available AI agents

### Health Check

```bash
caf doctor
```

Check that all components are working correctly.

## Configuration

### Manual Configuration

Create `.caf/config.yml` for custom settings:

```yaml
# .caf/config.yml
session:
  max_age: 24h
  format: yaml
adapters:
  - opencode
  - claude-code
  - universal-lsp
```

### Environment Variables

Optional environment variables:

```bash
export CAF_CONFIG_PATH="/path/to/config.yml"
export CAF_SESSION_DIR="/path/to/sessions"
```

## Integration with AI Agents

### Supported Agents

- **opencode** - Native integration
- **Claude Code** - Native integration  
- **Cursor** - Via universal LSP adapter
- **Cline** - Via universal LSP adapter
- **Any AI agent** - Via pre-commit hooks

### Auto-Detection

```bash
caf init --auto-detect
```

This automatically detects and configures available agents in your environment.