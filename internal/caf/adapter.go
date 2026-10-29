package caf

import (
	"fmt"
)

// EventTier represents the capability tier for an event
type EventTier string

const (
	TierNative    EventTier = "native"
	TierWrapper   EventTier = "wrapper"
	TierWatcher   EventTier = "watcher"
	TierUniversal EventTier = "universal"
	TierNone      EventTier = "none"
)

// CapabilityMatrix represents the capabilities of an adapter
type CapabilityMatrix struct {
	PreSession  EventTier `yaml:"pre_session"`
	PostEdit    EventTier `yaml:"post_edit"`
	PreCommit   EventTier `yaml:"pre_commit"`
	MCP         bool      `yaml:"mcp"`
	Tier        int       `yaml:"tier"` // 1, 2, 3, 4
}

// Adapter is the interface that all adapters must implement
// INV-ARCH-002: Adapter interface (5 methods) is stable within a major version
type Adapter interface {
	// ID returns the unique identifier for this adapter
	ID() string
	
	// Capabilities returns the capability matrix for this adapter
	Capabilities() CapabilityMatrix
	
	// Install installs the adapter and its hooks/config
	Install() error
	
	// Uninstall removes the adapter and its hooks/config
	Uninstall() error
	
	// Detect returns true if this adapter is installed on the system
	Detect() bool
}

// LanguageAdapter is the interface for language-specific adapters
type LanguageAdapter interface {
	// ID returns the language identifier
	ID() string
	
	// Detect returns true if this language is present in the repo
	Detect() bool
	
	// ScanImports returns the import graph
	ScanImports() []Import
	
	// BuildSymbolGraph builds the symbol graph for the language
	BuildSymbolGraph() SymbolGraph
	
	// ResolveExternalRepo resolves an external package reference
	ResolveExternalRepo(pkg string) RepoRef
}

// Import represents an import statement
type Import struct {
	Path     string `yaml:"path"`
	Alias    string `yaml:"alias,omitempty"`
	Line     int    `yaml:"line"`
	Column   int    `yaml:"column"`
}

// Symbol represents a symbol in the code
type Symbol struct {
	Name       string `yaml:"name"`
	Kind       string `yaml:"kind"` // function, struct, interface, etc.
	Parent     string `yaml:"parent,omitempty"`
	Definition string `yaml:"definition,omitempty"`
	Location   string `yaml:"location"` // file:line:column
}

// SymbolGraph represents the symbol graph for a language
type SymbolGraph struct {
	Symbols   []Symbol `yaml:"symbols"`
	Imports   []Import `yaml:"imports"`
	Relations []Relation `yaml:"relations,omitempty"`
}

// Relation represents a relation between symbols
type Relation struct {
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Kind     string `yaml:"kind"` // calls, inherits, implements, etc.
	Location string `yaml:"location"` // where the relation is defined
}

// RepoRef represents a reference to an external repository
type RepoRef struct {
	URL      string `yaml:"url"`
	Commit   string `yaml:"commit"`
	Path     string `yaml:"path"`
	Resolved bool   `yaml:"resolved"`
}

// Registry holds all registered adapters
var Registry = make(map[string]Adapter)

// RegisterAdapter registers an adapter
func RegisterAdapter(adapter Adapter) {
	Registry[adapter.ID()] = adapter
}

// GetAdapter returns an adapter by ID
func GetAdapter(id string) (Adapter, bool) {
	adapter, exists := Registry[id]
	return adapter, exists
}

// ListAdapters returns all registered adapter IDs
func ListAdapters() []string {
	var ids []string
	for id := range Registry {
		ids = append(ids, id)
	}
	return ids
}

// GetLanguageAdapter returns a language adapter by ID
func GetLanguageAdapter(id string) (LanguageAdapter, bool) {
	adapter, exists := Registry[id]
	if !exists {
		return nil, false
	}
	
	langAdapter, ok := adapter.(LanguageAdapter)
	return langAdapter, ok
}

// ValidateCapabilityMatrix validates a capability matrix
func ValidateCapabilityMatrix(matrix CapabilityMatrix) error {
	// Validate event tiers
	validTiers := map[EventTier]bool{
		TierNative:    true,
		TierWrapper:   true,
		TierWatcher:   true,
		TierUniversal: true,
		TierNone:      true,
	}
	
	if !validTiers[matrix.PreSession] {
		return fmt.Errorf("invalid pre-session tier: %s", matrix.PreSession)
	}
	if !validTiers[matrix.PostEdit] {
		return fmt.Errorf("invalid post-edit tier: %s", matrix.PostEdit)
	}
	if !validTiers[matrix.PreCommit] {
		return fmt.Errorf("invalid pre-commit tier: %s", matrix.PreCommit)
	}
	
	// Validate tier number
	if matrix.Tier < 1 || matrix.Tier > 4 {
		return fmt.Errorf("invalid tier number: %d", matrix.Tier)
	}
	
	return nil
}