// Package lsp — Tier-4 universal LSP fallback менеджер (ADR-013, INV-ARCH-003).
//
// FallbackManager владеет жизненным циклом gopls-процесса через
// universal_lsp.UniversalLSPAdapter, кэширует символный граф в
// .caf/mental-model/graph.yml (regenerable, C-09) и деградирует мягко:
// 3 краша за сессию → отказ от LSP до перезапуска процесса caf, без падения
// вызывающих команд (ADR-013 mitigation).
package lsp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/adapters"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

const (
	// startupBudget — soft-бюджет INV-ADR-013 (~500ms на малом репо).
	// Превышение логируется, но не отменяет старт: gopls на большом репо
	// стартует дольше, и это ожидаемо.
	startupBudget = 500 * time.Millisecond

	// startupTimeout — hard-потолок handshake. Превышение = ошибка старта:
	// без потолка `caf doctor` мог бы зависнуть на минуту на сломанном gopls.
	startupTimeout = 2 * time.Second

	// syncTimeout — потолок на workspace/symbol запрос.
	syncTimeout = 10 * time.Second

	// maxCrashes — после N фатальных сбоев LSP за сессию менеджер отказывается
	// перезапускать сервер (ADR-013: "fall back to filesystem-only mode").
	maxCrashes = 3
)

// FallbackManager — точка входа Tier-4 fallback для одной репозитории.
type FallbackManager struct {
	repoRoot     string
	adapter      *adapters.UniversalLSPAdapter
	symbols      []caf.Symbol
	lastSync     time.Time
	syncInterval time.Duration
	crashes      int
	startupOnce  bool
}

// NewFallbackManager создаёт менеджер; процесс LSP не запускается до Start
// (lazy-load, INV-ARCH-003).
func NewFallbackManager(repoRoot string) *FallbackManager {
	return &FallbackManager{
		repoRoot:     repoRoot,
		adapter:      &adapters.UniversalLSPAdapter{},
		syncInterval: 30 * time.Minute,
	}
}

// Start запускает LSP-сервер и загружает начальный символный граф.
// Идемпотентен: повторный вызов при живом процессе — no-op.
// При 3+ крашах за сессию возвращает ошибку без попытки старта (мягкая
// деградация ADR-013): вызывающий (doctor, dashboard) показывает warning,
// но не падает.
func (m *FallbackManager) Start() error {
	if m.crashes >= maxCrashes {
		return fmt.Errorf("lsp fallback disabled after %d crashes this session (ADR-013 degradation)", m.crashes)
	}
	if m.adapter.Started() {
		return nil
	}

	if !m.adapter.Detect() {
		return errors.New("gopls not available — install Go and gopls (ADR-026: MVP requires manual install)")
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	if err := m.adapter.Start(ctx, m.repoRoot); err != nil {
		m.crashes++
		return fmt.Errorf("failed to start LSP server: %w", err)
	}
	m.startupOnce = true

	if elapsed := time.Since(start); elapsed > startupBudget {
		// Soft budget (INV-ADR-013): warn, но не fail.
		log.Printf("⚠ LSP startup %v exceeds %v budget (INV-ADR-013 soft budget)", elapsed, startupBudget)
	}

	// Начальная загрузка графа — best-effort: крэш синка не должен откатить
	// успешный старт, doctor оценит содержимое отдельно.
	if err := m.SyncSymbols(); err != nil {
		log.Printf("⚠ initial symbol sync failed: %v", err)
	}
	return nil
}

// Stop останавливает LSP-процесс. Идемпотентен.
func (m *FallbackManager) Stop() error {
	return m.adapter.Stop()
}

// SyncSymbols запрашивает workspace/symbol и обновляет кэш графа.
// Контекст со syncTimeout: зависший gopls не должен блокировать вызывающего.
func (m *FallbackManager) SyncSymbols() error {
	if !m.adapter.Started() {
		return errors.New("lsp not started")
	}
	ctx, cancel := context.WithTimeout(context.Background(), syncTimeout)
	defer cancel()

	symbols, err := m.adapter.Symbols(ctx)
	if err != nil {
		// Фатальные ошибки канала (процесс умер) — считаем крашем.
		if errors.Is(err, context.DeadlineExceeded) || isConnError(err) {
			m.crashes++
		}
		return fmt.Errorf("failed to get symbols: %w", err)
	}

	m.symbols = symbols
	m.lastSync = time.Now()

	if err := m.SaveSymbolGraph(); err != nil {
		// Кэш regenerable (C-09): несохранённый граф — не повод терять синк.
		log.Printf("⚠ failed to save symbol graph cache: %v", err)
	}
	log.Printf("✓ synced %d symbols from LSP", len(symbols))
	return nil
}

// isConnError отличает «сервер умер/отвалился» от прочих ошибок запроса.
func isConnError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return errors.Is(err, os.ErrClosed) ||
		strings.Contains(s, "connection closed") ||
		strings.Contains(s, "broken pipe") ||
		strings.Contains(s, "use of closed")
}

// GetSymbols возвращает кэшированные символы; при устаревшем кэше (> syncInterval)
// пытается пересинхронизироваться (best-effort: при ошибке отдаём старые данные —
// stale лучше, чем пустота).
func (m *FallbackManager) GetSymbols() []caf.Symbol {
	if !m.lastSync.IsZero() && time.Since(m.lastSync) > m.syncInterval && m.adapter.Started() {
		if err := m.SyncSymbols(); err != nil {
			log.Printf("⚠ background symbol resync failed: %v", err)
		}
	}
	return m.symbols
}

// SaveSymbolGraph пишет граф в .caf/mental-model/graph.yml (ADR-013 invariant).
// Родительская директория создаётся: свежий репозиторий без .caf/mental-model
// не должен ронять запись кэша.
func (m *FallbackManager) SaveSymbolGraph() error {
	graphFile := filepath.Join(m.repoRoot, ".caf", "mental-model", "graph.yml")
	if err := storage.EnsureDir(filepath.Dir(graphFile)); err != nil {
		return fmt.Errorf("create mental-model dir: %w", err)
	}
	return storage.WriteYAML(graphFile, caf.SymbolGraph{Symbols: m.symbols})
}

// LoadSymbolGraph загружает граф из кэша (например, до доступного LSP —
// doctor может читать кэш без запуска gopls).
func (m *FallbackManager) LoadSymbolGraph() error {
	graphFile := filepath.Join(m.repoRoot, ".caf", "mental-model", "graph.yml")
	var graph caf.SymbolGraph
	if err := storage.ReadYAML(graphFile, &graph); err != nil {
		return fmt.Errorf("failed to load symbol graph: %w", err)
	}
	m.symbols = graph.Symbols
	m.lastSync = time.Now()
	return nil
}

// IsAvailable сообщает, может ли fallback стартовать (gopls в $PATH и нет
// деградации после 3 крашей).
func (m *FallbackManager) IsAvailable() bool {
	return m.adapter.Detect() && m.crashes < maxCrashes
}

// Degraded сообщает, отключён ли LSP после крашей (для doctor UI).
func (m *FallbackManager) Degraded() bool {
	return m.crashes >= maxCrashes
}

// Crashes возвращает число фатальных сбоев за сессию.
func (m *FallbackManager) Crashes() int {
	return m.crashes
}

// GetAdapter возвращает underlying адаптер (для adapters-тестов и UI).
func (m *FallbackManager) GetAdapter() *adapters.UniversalLSPAdapter {
	return m.adapter
}

// LazyLoad запускает LSP только при первом реальном запросе (INV-ARCH-003:
// lazy-load на caf sync / caf tour; в MVP эти команды BETA/PROD, точка
// вызова — dashboard/doctor).
func (m *FallbackManager) LazyLoad() error {
	if m.adapter.Started() {
		return nil
	}
	return m.Start()
}

// GetSymbolGraph возвращает текущий граф (пустой, если ещё не синкали).
func (m *FallbackManager) GetSymbolGraph() caf.SymbolGraph {
	return caf.SymbolGraph{Symbols: m.symbols}
}

// StartupMeasured сообщает, выполнялся ли старт LSP в этой сессии
// (врач использует для различения "кэш есть, сервер не поднимали" и
// "сервер отвечает").
func (m *FallbackManager) StartupMeasured() bool {
	return m.startupOnce
}

// Initialize создаёт менеджер и запускает LSP (ранний convenience-путь;
// предпочитайте GetManagerForRepo — он уважает кэш).
func Initialize(repoRoot string) (*FallbackManager, error) {
	manager := NewFallbackManager(repoRoot)
	if err := manager.Start(); err != nil {
		return nil, err
	}
	return manager, nil
}

// GetManagerForRepo возвращает менеджер для репозитория: читает кэш графа,
// при отсутствии/битом кэше пытается запустить LSP. .caf/ отсутствует →
// ошибка: fallback живёт внутри инициализированной репозитории.
func GetManagerForRepo(repoRoot string) (*FallbackManager, error) {
	cafDir := storage.GetCAFDir(repoRoot)
	if _, err := os.Stat(cafDir); err != nil {
		return nil, fmt.Errorf(".caf directory not found: %w", err)
	}

	manager := NewFallbackManager(repoRoot)

	if err := manager.LoadSymbolGraph(); err == nil {
		// Кэш есть — возвращаем без старта LSP (lazy-load, INV-ARCH-003).
		// Пересинк произойдёт по stale-интервалу или по явному SyncSymbols.
		return manager, nil
	}

	if err := manager.Start(); err != nil {
		// Без gopls manager остаётся в degraded-состоянии с пустым графом —
		// вызывающий (doctor) показывает warning, а не падает.
		log.Printf("⚠ LSP fallback unavailable: %v", err)
	}
	return manager, nil
}
