// Package adapters объединяет Tier-1 (opencode, Claude Code) и Tier-4
// (universal LSP) адаптеры под единым импортом для cmd/caf и lsp.
package adapters

import (
	claudecode "github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters/claude_code"
	opencode "github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters/opencode"
	universallsp "github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters/universal_lsp"
)

// UniversalLSPAdapter — алиас типадля Tier-4 fallback адаптера (ADR-013).
type UniversalLSPAdapter = universallsp.UniversalLSPAdapter

// Initialize регистрирует все встроенные адаптеры в реестре caf.
func Initialize() {
	opencode.Initialize()
	claudecode.Initialize()
	universallsp.Initialize()
}

// OpencodeInitialize регистрирует Tier-1 opencode адаптер (FR-ADP-001).
func OpencodeInitialize() {
	opencode.Initialize()
}

// ClaudeCodeInitialize регистрирует Tier-1 Claude Code адаптер (FR-ADP-001).
func ClaudeCodeInitialize() {
	claudecode.Initialize()
}

// UniversalLSPInitialize регистрирует Tier-4 universal LSP fallback (FR-ADP-004).
func UniversalLSPInitialize() {
	universallsp.Initialize()
}
