// Package universallsp — Tier-4 universal LSP fallback адаптер (FR-ADP-004,
// ADR-013). Символьный граф строит gopls через стандартный LSP-протокол
// поверх stdio (см. client.go); MCP=false, все три события — через
// fallback-цепочку ADR-010 (wrapper/watcher/universal), не нативно.
package universallsp

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// UniversalLSPAdapter реализует Adapter (5 методов, INV-ADR-009) и добавляет
// жизненный цикл LSP-клиента (Start/Symbols/Stop) для internal/caf/lsp.
type UniversalLSPAdapter struct {
	mu     sync.Mutex
	client *Client
}

// ID — идентификатор в реестре (совпадает с контракт-тестом ARCHITECT.md §1).
func (a *UniversalLSPAdapter) ID() string {
	return "universal-lsp"
}

// Capabilities — строка ARCHITECT.md §1 «Universal (Tier 4)»:
// pre-session=wrapper, post-edit=watcher, pre-commit=universal, MCP=false.
func (a *UniversalLSPAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession: caf.TierWrapper,
		PostEdit:   caf.TierWatcher,
		PreCommit:  caf.TierUniversal,
		MCP:        false,
		Tier:       4,
	}
}

// Detect — есть ли gopls в $PATH (ADR-026: MVP требует ручной установки,
// бандлинга нет). Пустой PATH → false, без паник.
func (a *UniversalLSPAdapter) Detect() bool {
	_, err := exec.LookPath(DefaultBinary)
	return err == nil
}

// Install идемпотентен (ARCHITECT.md §3.1): проверяет, что gopls доступен,
// и указывает, как его поставить. Адаптер не пишет никаких конфигов —
// Tier-4 zero-config по определению (ADR-013); единственный артефакт,
// который он порождает, — кэш graph.yml, и тот regenerable (C-09).
func (a *UniversalLSPAdapter) Install() error {
	if !a.Detect() {
		return fmt.Errorf("gopls not found in $PATH; install it with: " +
			"go install golang.org/x/tools/gopls@latest " +
			"(MVP requires manual install, ADR-026)")
	}
	return nil
}

// Uninstall останавливает живой LSP-процесс (если запущен) и идемпотентен.
// Кэш graph.yml не трогаю: он regenerable и общий для FallbackManager —
// удалять его из Uninstall адаптера значило бы ломать контекст вызывающего.
func (a *UniversalLSPAdapter) Uninstall() error {
	return a.Stop()
}

// Start запускает gopls для repoRoot, если ещё не запущен (идемпотентно).
// Таймаут handshake ограничивает ctx — бюджет INV-ADR-013 (~500ms soft,
// жёсткий потолок задаёт caller).
func (a *UniversalLSPAdapter) Start(ctx context.Context, repoRoot string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.client != nil && a.client.Started() {
		return nil
	}
	client := NewClient()
	if err := client.Start(ctx, repoRoot); err != nil {
		return err
	}
	a.client = client
	return nil
}

// Symbols запрашивает workspace/symbol у живого клиента.
func (a *UniversalLSPAdapter) Symbols(ctx context.Context) ([]caf.Symbol, error) {
	a.mu.Lock()
	client := a.client
	a.mu.Unlock()
	if client == nil || !client.Started() {
		return nil, errors.New("lsp: not started — call Start first")
	}
	return client.WorkspaceSymbols(ctx, "")
}

// Stop останавливает LSP-процесс. Идемпотентен, ошибка — только если
// Kill завершился неуспешно и по не-штатной причине.
func (a *UniversalLSPAdapter) Stop() error {
	a.mu.Lock()
	client := a.client
	a.client = nil
	a.mu.Unlock()
	if client == nil {
		return nil
	}
	return client.Stop()
}

// Started сообщает, запущен ли LSP-процесс (для fallback-менеджера и doctor).
func (a *UniversalLSPAdapter) Started() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client != nil && a.client.Started()
}

// DefaultStartupBudget — soft-бюджет старта LSP (INV-ADR-013: ~500ms на
// малом репозитории). Превышение — предупреждение, не отказ: gopls на
// крупном репо стартует дольше, и это нормально.
const DefaultStartupBudget = 500 * time.Millisecond

// Initialize регистрирует Tier-4 universal LSP fallback в реестре caf
// (FR-ADP-004). Вызывается adapters.Initialize().
func Initialize() {
	caf.RegisterAdapter(&UniversalLSPAdapter{})
}
