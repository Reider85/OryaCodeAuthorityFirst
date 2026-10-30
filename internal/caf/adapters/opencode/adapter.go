package adapters

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// OpencodeAdapter implements the Adapter interface for opencode
type OpencodeAdapter struct{}

// ID returns the unique identifier for this adapter
func (a *OpencodeAdapter) ID() string {
	return "opencode"
}

// Capabilities returns the capability matrix for this adapter
func (a *OpencodeAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession:  caf.TierNative,
		PostEdit:    caf.TierNative,
		PreCommit:   caf.TierNative,
		MCP:         true,
		Tier:        1,
	}
}

// Install installs the opencode adapter and its hooks/config
func (a *OpencodeAdapter) Install() error {
	// Check if opencode is already installed
	if a.Detect() {
		return fmt.Errorf("opencode is already installed")
	}

	// Install opencode (placeholder - actual installation would depend on opencode's installer)
	// For MVP, we assume opencode is already installed and we just need to configure it
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "config.yml")
	
	// Create opencode config directory
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create opencode config directory: %w", err)
	}

	// Create basic opencode config
	configContent := `# Opencode configuration for CodeAuthorityFirst
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
		return fmt.Errorf("failed to write opencode config: %w", err)
	}

	return nil
}

// Uninstall removes the opencode adapter and its hooks/config
func (a *OpencodeAdapter) Uninstall() error {
	// Remove opencode config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "config.yml")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove opencode config: %w", err)
	}

	// Remove opencode hooks directory
	hooksPath := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "hooks")
	if err := os.RemoveAll(hooksPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove opencode hooks: %w", err)
	}

	return nil
}

// Detect returns true if opencode is installed on the system
func (a *OpencodeAdapter) Detect() bool {
	// Check if opencode command exists
	if _, err := exec.LookPath("opencode"); err != nil {
		return false
	}

	// Check if opencode config directory exists
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "opencode")
	if _, err := os.Stat(configPath); err != nil {
		return false
	}

	return true
}

// GetOpencodeConfigPath returns the path to opencode config
func GetOpencodeConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "opencode", "config.yml")
}

// GetOpencodeHooksPath returns the path to opencode hooks directory
func GetOpencodeHooksPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "opencode", "hooks")
}

// OpencodeLanguageAdapter implements LanguageAdapter for Go
type OpencodeLanguageAdapter struct{}

// ID returns the language identifier
func (a *OpencodeLanguageAdapter) ID() string {
	return "go"
}

// Detect returns true if Go is present in the repo
func (a *OpencodeLanguageAdapter) Detect() bool {
	// Check if go.mod exists
	goModPath := filepath.Join(".", "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return true
	}

	// Check for Go source files
	goFiles, err := filepath.Glob("*.go")
	if err != nil {
		return false
	}

	return len(goFiles) > 0
}

// ScanImports returns the import graph
func (a *OpencodeLanguageAdapter) ScanImports() []caf.Import {
	// Placeholder - in real implementation, parse Go files
	return []caf.Import{
		{
			Path:   "fmt",
			Alias:  "",
			Line:   3,
			Column: 8,
		},
		{
			Path:   "github.com/Reider85/OryaCodeAuthorityFirst/internal/caf",
			Alias:  "",
			Line:   4,
			Column: 8,
		},
	}
}

// BuildSymbolGraph builds the symbol graph for Go
func (a *OpencodeLanguageAdapter) BuildSymbolGraph() caf.SymbolGraph {
	// Placeholder - in real implementation, use gopls or go list
	return caf.SymbolGraph{
		Symbols: []caf.Symbol{
			{
				Name:       "main",
				Kind:       "function",
				Parent:     "",
				Definition: "func main()",
				Location:   "main.go:6",
			},
			{
				Name:       "Hello",
				Kind:       "function",
				Parent:     "",
				Definition: "func Hello(name string) string",
				Location:   "hello.go:3",
			},
		},
		Imports: a.ScanImports(),
	}
}

// ResolveExternalRepo resolves an external package reference
func (a *OpencodeLanguageAdapter) ResolveExternalRepo(pkg string) caf.RepoRef {
	// Placeholder - in real implementation, use go list
	return caf.RepoRef{
		URL:      "https://github.com/golang/go",
		Commit:   "latest",
		Path:     pkg,
		Resolved: true,
	}
}

// Capabilities возвращает матрицу возможностей языкового адаптера.
// Языковые адаптеры не обрабатывают события хуков — все тиры TierNone.
func (a *OpencodeLanguageAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession: caf.TierNone,
		PostEdit:   caf.TierNone,
		PreCommit:  caf.TierNone,
		MCP:        false,
		Tier:       1,
	}
}

// Install языковому адаптеру не нужен — он детектирует язык, а не ставит хуки.
func (a *OpencodeLanguageAdapter) Install() error { return nil }

// Uninstall языковому адаптеру не нужен — см. Install.
func (a *OpencodeLanguageAdapter) Uninstall() error { return nil }

// Initialize initializes the opencode adapter
func Initialize() {
	adapter := &OpencodeAdapter{}
	caf.RegisterAdapter(adapter)
	
	// Register language adapter
	langAdapter := &OpencodeLanguageAdapter{}
	caf.RegisterAdapter(langAdapter)
}