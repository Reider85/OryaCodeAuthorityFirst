package adapters

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// ClaudeCodeAdapter implements the Adapter interface for Claude Code
type ClaudeCodeAdapter struct{}

// ID returns the unique identifier for this adapter
func (a *ClaudeCodeAdapter) ID() string {
	return "claude-code"
}

// Capabilities returns the capability matrix for this adapter
func (a *ClaudeCodeAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession:  caf.TierNative,
		PostEdit:    caf.TierNative,
		PreCommit:   caf.TierNative,
		MCP:         true,
		Tier:        1,
	}
}

// Install installs the Claude Code adapter and its hooks/config
func (a *ClaudeCodeAdapter) Install() error {
	// Check if claude code is already installed
	if a.Detect() {
		return fmt.Errorf("claude code is already installed")
	}

	// Install Claude Code (placeholder - actual installation would depend on Claude Code's installer)
	// For MVP, we assume Claude Code is already installed and we just need to configure it
	configPath := filepath.Join(os.Getenv("HOME"), ".claude", "config.yml")
	
	// Create Claude Code config directory
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create claude code config directory: %w", err)
	}

	// Create basic Claude Code config
	configContent := `# Claude Code configuration for CodeAuthorityFirst
hooks:
  pre_session:
    enabled: true
    command: caf hook pre-session --session {{session_id}}
  post_edit:
    enabled: true
    command: caf hook post-edit --session {{session_id}} --files "{{files_changed}}"
  pre_commit:
    enabled: true
    command: caf hook pre-commit --session {{session_id}} --message "{{commit_message}}"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write claude code config: %w", err)
	}

	return nil
}

// Uninstall removes the Claude Code adapter and its hooks/config
func (a *ClaudeCodeAdapter) Uninstall() error {
	// Remove Claude Code config
	configPath := filepath.Join(os.Getenv("HOME"), ".claude", "config.yml")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove claude code config: %w", err)
	}

	// Remove Claude Code hooks directory
	hooksPath := filepath.Join(os.Getenv("HOME"), ".claude", "hooks")
	if err := os.RemoveAll(hooksPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove claude code hooks: %w", err)
	}

	return nil
}

// Detect returns true if Claude Code is installed on the system
func (a *ClaudeCodeAdapter) Detect() bool {
	// Check if claude command exists
	if _, err := exec.LookPath("claude"); err != nil {
		return false
	}

	// Check if Claude Code config directory exists
	configPath := filepath.Join(os.Getenv("HOME"), ".claude")
	if _, err := os.Stat(configPath); err != nil {
		return false
	}

	return true
}

// GetClaudeCodeConfigPath returns the path to Claude Code config
func GetClaudeCodeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".claude", "config.yml")
}

// GetClaudeCodeHooksPath returns the path to Claude Code hooks directory
func GetClaudeCodeHooksPath() string {
	return filepath.Join(os.Getenv("HOME"), ".claude", "hooks")
}

// ClaudeCodeLanguageAdapter implements LanguageAdapter for Python
type ClaudeCodeLanguageAdapter struct{}

// ID returns the language identifier
func (a *ClaudeCodeLanguageAdapter) ID() string {
	return "python"
}

// Detect returns true if Python is present in the repo
func (a *ClaudeCodeLanguageAdapter) Detect() bool {
	// Check if pyproject.toml or setup.py exists
	pyProjectPath := filepath.Join(".", "pyproject.toml")
	if _, err := os.Stat(pyProjectPath); err == nil {
		return true
	}

	setupPath := filepath.Join(".", "setup.py")
	if _, err := os.Stat(setupPath); err == nil {
		return true
	}

	// Check for Python source files
	pyFiles, err := filepath.Glob("*.py")
	if err != nil {
		return false
	}

	return len(pyFiles) > 0
}

// ScanImports returns the import graph
func (a *ClaudeCodeLanguageAdapter) ScanImports() []caf.Import {
	// Placeholder - in real implementation, parse Python files
	return []caf.Import{
		{
			Path:   "os",
			Alias:  "",
			Line:   1,
			Column: 0,
		},
		{
			Path:   "typing",
			Alias:  "",
			Line:   2,
			Column: 0,
		},
	}
}

// BuildSymbolGraph builds the symbol graph for Python
func (a *ClaudeCodeLanguageAdapter) BuildSymbolGraph() caf.SymbolGraph {
	// Placeholder - in real implementation, use pyright or jedi
	return caf.SymbolGraph{
		Symbols: []caf.Symbol{
			{
				Name:       "main",
				Kind:       "function",
				Parent:     "",
				Definition: "def main():",
				Location:   "main.py:4",
			},
			{
				Name:       "calculate",
				Kind:       "function",
				Parent:     "",
				Definition: "def calculate(x, y):",
				Location:   "utils.py:3",
			},
		},
		Imports: a.ScanImports(),
	}
}

// ResolveExternalRepo resolves an external package reference
func (a *ClaudeCodeLanguageAdapter) ResolveExternalRepo(pkg string) caf.RepoRef {
	// Placeholder - in real implementation, use pip index
	return caf.RepoRef{
		URL:      "https://pypi.org/project/" + pkg,
		Commit:   "latest",
		Path:     pkg,
		Resolved: true,
	}
}

// Capabilities возвращает матрицу возможностей языкового адаптера.
// Языковые адаптеры не обрабатывают события хуков — все тиры TierNone.
func (a *ClaudeCodeLanguageAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession: caf.TierNone,
		PostEdit:   caf.TierNone,
		PreCommit:  caf.TierNone,
		MCP:        false,
		Tier:       1,
	}
}

// Install языковому адаптеру не нужен — он детектирует язык, а не ставит хуки.
func (a *ClaudeCodeLanguageAdapter) Install() error { return nil }

// Uninstall языковому адаптеру не нужен — см. Install.
func (a *ClaudeCodeLanguageAdapter) Uninstall() error { return nil }

// Initialize initializes the Claude Code adapter
func Initialize() {
	adapter := &ClaudeCodeAdapter{}
	caf.RegisterAdapter(adapter)
	
	// Register language adapter
	langAdapter := &ClaudeCodeLanguageAdapter{}
	caf.RegisterAdapter(langAdapter)
}