package lsp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// FallbackManager handles the universal LSP fallback
type FallbackManager struct {
	repoRoot     string
	adapter      *adapters.UniversalLSPAdapter
	symbols      []storage.Symbol
	lastSync     time.Time
	syncInterval time.Duration
}

// NewFallbackManager creates a new LSP fallback manager
func NewFallbackManager(repoRoot string) *FallbackManager {
	return &FallbackManager{
		repoRoot:     repoRoot,
		adapter:      &adapters.UniversalLSPAdapter{},
		syncInterval: 30 * time.Minute, // Sync every 30 minutes
	}
}

// Start starts the LSP server and loads the symbol graph
func (m *FallbackManager) Start() error {
	log.Printf("🔧 Starting universal LSP fallback for %s", m.repoRoot)

	// Check if gopls is available
	if !m.adapter.Detect() {
		return fmt.Errorf("gopls not available, LSP fallback disabled")
	}

	// Start the LSP server
	if err := m.adapter.StartLSP(m.repoRoot); err != nil {
		return fmt.Errorf("failed to start LSP server: %w", err)
	}

	// Load initial symbol graph
	if err := m.SyncSymbols(); err != nil {
		log.Printf("⚠ Failed to load initial symbol graph: %v", err)
		// Continue anyway
	}

	return nil
}

// Stop stops the LSP server
func (m *FallbackManager) Stop() error {
	log.Printf("🛑 Stopping universal LSP fallback")
	return m.adapter.Stop()
}

// SyncSymbols syncs symbols from the LSP server
func (m *FallbackManager) SyncSymbols() error {
	log.Printf("🔄 Syncing symbols from LSP server")

	// Get symbols from LSP
	symbols, err := m.adapter.GetSymbols()
	if err != nil {
		return fmt.Errorf("failed to get symbols: %w", err)
	}

	// Convert to our format
	m.symbols = make([]storage.Symbol, len(symbols))
	for i, symbol := range symbols {
		m.symbols[i] = storage.Symbol{
			Name:       symbol.Name,
			Kind:       symbol.Kind,
			Location:   symbol.Location,
		}
	}

	m.lastSync = time.Now()

	// Save symbol graph to cache
	if err := m.SaveSymbolGraph(); err != nil {
		log.Printf("⚠ Failed to save symbol graph: %v", err)
		// Continue anyway
	}

	log.Printf("✓ Synced %d symbols", len(m.symbols))
	return nil
}

// GetSymbols returns the cached symbols
func (m *FallbackManager) GetSymbols() []storage.Symbol {
	// Check if we need to sync
	if time.Since(m.lastSync) > m.syncInterval {
		if err := m.SyncSymbols(); err != nil {
			log.Printf("⚠ Failed to sync symbols: %v", err)
		}
	}

	return m.symbols
}

// SaveSymbolGraph saves the symbol graph to cache
func (m *FallbackManager) SaveSymbolGraph() error {
	graphFile := filepath.Join(m.repoRoot, ".caf", "mental-model", "graph.yml")
	
	graph := storage.SymbolGraph{
		Symbols: m.symbols,
	}

	return storage.WriteYAML(graphFile, graph)
}

// LoadSymbolGraph loads the symbol graph from cache
func (m *FallbackManager) LoadSymbolGraph() error {
	graphFile := filepath.Join(m.repoRoot, ".caf", "mental-model", "graph.yml")
	
	var graph storage.SymbolGraph
	if err := storage.ReadYAML(graphFile, &graph); err != nil {
		return fmt.Errorf("failed to load symbol graph: %w", err)
	}

	m.symbols = graph.Symbols
	m.lastSync = time.Now()

	return nil
}

// IsAvailable returns true if the LSP fallback is available
func (m *FallbackManager) IsAvailable() bool {
	return m.adapter.Detect()
}

// GetAdapter returns the underlying adapter
func (m *FallbackManager) GetAdapter() *adapters.UniversalLSPAdapter {
	return m.adapter
}

// LazyLoad lazily loads the LSP server when needed
func (m *FallbackManager) LazyLoad() error {
	if m.adapter == nil {
		m.adapter = &adapters.UniversalLSPAdapter{}
	}

	if !m.adapter.Detect() {
		return fmt.Errorf("gopls not available")
	}

	// Check if we already have symbols loaded
	if len(m.symbols) == 0 {
		return m.SyncSymbols()
	}

	return nil
}

// GetSymbolGraph returns the symbol graph for the repository
func (m *FallbackManager) GetSymbolGraph() storage.SymbolGraph {
	return storage.SymbolGraph{
		Symbols: m.symbols,
	}
}

// Initialize initializes the LSP fallback for a repository
func Initialize(repoRoot string) (*FallbackManager, error) {
	manager := NewFallbackManager(repoRoot)
	
	// Start the LSP server
	if err := manager.Start(); err != nil {
		return nil, err
	}

	return manager, nil
}

// GetManagerForRepo returns an LSP fallback manager for the given repository
func GetManagerForRepo(repoRoot string) (*FallbackManager, error) {
	// Check if .caf directory exists
	cafDir := storage.GetCAFDir(repoRoot)
	if _, err := os.Stat(cafDir); err != nil {
		return nil, fmt.Errorf(".caf directory not found: %w", err)
	}

	// Create or load manager
	manager := NewFallbackManager(repoRoot)
	
	// Try to load existing symbol graph
	if err := manager.LoadSymbolGraph(); err != nil {
		// If loading fails, start fresh
		if err := manager.Start(); err != nil {
			return nil, err
		}
	}

	return manager, nil
}