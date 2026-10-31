package adapters

import (
	"testing"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// TestInitializeRegistersAllBuiltins — один вызов Initialize() ставит в
// реестр оба Tier-1 адаптера и Tier-4 fallback (то, что вызывает cmd/caf).
func TestInitializeRegistersAllBuiltins(t *testing.T) {
	Initialize()
	for _, id := range []string{"opencode", "claude-code", "universal-lsp"} {
		if _, ok := caf.GetAdapter(id); !ok {
			t.Errorf("Initialize() must register %q", id)
		}
	}
}

// TestPerAdapterInitializeWrappers — точечные регистрации для шага 8
// (universal_lsp) и юнит-тестов, которым не нужен весь реестр.
func TestPerAdapterInitializeWrappers(t *testing.T) {
	tests := []struct {
		wrapper string
		initFn  func()
		id      string
	}{
		{"OpencodeInitialize", OpencodeInitialize, "opencode"},
		{"ClaudeCodeInitialize", ClaudeCodeInitialize, "claude-code"},
		{"UniversalLSPInitialize", UniversalLSPInitialize, "universal-lsp"},
	}
	for _, tt := range tests {
		t.Run(tt.wrapper, func(t *testing.T) {
			tt.initFn()
			if _, ok := caf.GetAdapter(tt.id); !ok {
				t.Errorf("%s() must register %q", tt.wrapper, tt.id)
			}
		})
	}
}
