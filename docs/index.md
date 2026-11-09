# CodeAuthorityFirst Documentation

Welcome to CodeAuthorityFirst (caf) - the provenance tracker for AI coding sessions.

## Overview

CodeAuthorityFirst records who wrote what, when, and why in AI coding sessions. It tracks which agent touched which file, from which prompt, through which iterations, and shows it in a TUI dashboard.

### Key Features

- **Single binary** - Pure Go, no CGO, no runtime dependencies
- **Offline & LLM-free** - No cloud calls, no API keys required
- **Fail-safe hooks** - Never blocks commits, warns and continues
- **Universal pre-commit hook** - Works with all AI agents
- **TUI dashboard** - Interactive interface for session tracking

## Quick Start

### Installation

```bash
git clone https://github.com/Reider85/OryaCodeAuthorityFirst.git
cd OryaCodeAuthorityFirst
go build -o caf ./cmd/caf
mv caf /usr/local/bin/
```

### Initialize

```bash
caf init --auto-detect    # Create .caf/, install hook, detect agents
caf doctor                # Health checks
caf dashboard             # Open TUI
```

## Usage

### Basic Commands

| Command | Description |
|---------|-------------|
| `caf init` | Initialize repository and install hooks |
| `caf who <file>` | Find session that modified a file |
| `caf dashboard` | Launch TUI dashboard |
| `caf adapters list` | List available adapters |
| `caf doctor` | Run health checks |

### Adapters

CodeAuthorityFirst supports multiple AI agents through adapters:

- **Tier 1**: Native support for opencode and Claude Code
- **Tier 4**: LSP-based fallback for any Go repository
- **Universal**: Works with any AI agent via pre-commit hooks

## Configuration

The `.caf/` directory contains configuration and session data:

```
.caf/
├── sessions/       # Session records
├── decisions/     # Decision logs
├── config.yml     # Configuration file
└── context.md     # Generated context
```

## Development

### Building

```bash
go build -o caf ./cmd/caf
```

### Testing

```bash
go test ./...        # Run all tests
go vet ./...         # Lint code
go build ./...       # Verify builds
```

### Contributing

Please follow conventional commits:
- `feat:` New features
- `fix:` Bug fixes  
- `docs:` Documentation
- `chore:` Maintenance

## Support

- 📖 **Documentation**: Full docs available at GitHub Pages
- 💬 **Discussions**: Join discussions for questions and ideas
- 🐛 **Issues**: Report bugs and request features
- 🌟 **Star**: Star the project if it helps your workflow

## License

This project is open source and available under the MIT License.