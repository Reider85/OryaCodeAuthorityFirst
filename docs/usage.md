# Usage Guide

## Basic Workflow

### 1. Initialize

```bash
caf init --auto-detect
```

Creates `.caf/` directory and installs pre-commit hook.

### 2. Work with AI Agents

Use your favorite AI agent (opencode, Claude Code, etc.) in your repository. The pre-commit hook will automatically track sessions.

### 3. Check Provenance

```bash
caf who src/main.go        # Who modified this file?
caf dashboard              # Open TUI dashboard
```

## Command Reference

### Initialization Commands

| Command | Description |
|---------|-------------|
| `caf init` | Initialize repository |
| `caf init --auto-detect` | Auto-detect and configure agents |
| `caf init --tui` | Interactive setup wizard |
| `caf doctor` | Health check |

### Session Tracking

| Command | Description |
|---------|-------------|
| `caf who <file>` | Find session that modified a file |
| `caf who <file:line>` | Find session that modified specific line |
| `caf when --pattern "*.go"` | Find sessions matching pattern (ALPHA) |
| `caf context` | Show/regenerate context.md |
| `caf context --refresh` | Force refresh context |

### Adapters Management

| Command | Description |
|---------|-------------|
| `caf adapters list` | List available adapters |
| `caf adapters list --tui` | Interactive adapter matrix |
| `caf auto-detect` | Scan for installed agents |

### Dashboard

| Command | Description |
|---------|-------------|
| `caf dashboard` | Launch TUI dashboard |
| `caf dashboard --tui` | Force TUI mode |

## TUI Dashboard

The dashboard provides five main widgets:

### Today
- Cards due, pending, and new blindspots
- Daily activity overview

### Authorship Score
- ACS (Authorship Score) and target
- Recall and coverage percentages
- 30-day trends

### Active Agents
- Status of configured adapters
- Health indicators

### Recent Decisions
- Recent coding decisions
- Context and reasoning

### Cross-repo
- Cross-repository analysis
- Stale decision warnings

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Cycle widgets |
| `j` / `k`, `↑` / `↓` | Navigate lists |
| `Enter` | Activate/expand |
| `Esc` | Back/dismiss |
| `r` | Refresh |
| `q` | Quit |
| `?` | Help |
| `/` | Search |

## Advanced Usage

### Custom Configuration

```bash
# Create custom config
cat > .caf/config.yml << EOF
session:
  max_age: 48h
  format: json
adapters:
  - opencode
  - claude-code
EOF
```

### Manual Session Tracking

```bash
# Hook events (usually automatic)
caf hook pre-session --repo . --session "my-session"
caf hook post-edit --repo . --session "my-session" --files "src/main.go"
```

### Debugging

```bash
# Check hook installation
ls -la .git/hooks/pre-commit

# View session data
ls .caf/sessions/

# Check logs
caf doctor --verbose
```

## Troubleshooting

### Common Issues

**Hook not working:**
```bash
# Reinstall hook
caf init --adapters opencode
```

**No sessions detected:**
```bash
# Check agent installation
caf auto-detect
# Use an AI agent to create sessions
```

**TUI not opening:**
```bash
# Check terminal compatibility
export TERM=xterm-256color
caf dashboard
```

### Performance Tips

- Regular cleanup: `find .caf/sessions/ -mtime +7 -delete`
- Use `gopls` for better symbol resolution
- Keep `.caf/` directory size manageable