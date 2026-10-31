package universallsp

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// startAdapter запускает адаптер с тестовым контекстом (5s на handshake).
func startAdapter(t *testing.T, a *UniversalLSPAdapter) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.Start(ctx, t.TempDir())
}

// testCtx — контекст с таймаутом для запросов в тестах.
func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// stubGopls кладёт исполняемый стаб gopls в dir; Detect()/Install() ищут
// его через LookPath (на Windows LookPath берёт .bat из PATHEXT).
func stubGopls(t *testing.T, dir string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		if err := os.WriteFile(filepath.Join(dir, DefaultBinary+".bat"), []byte("@echo off\r\n"), 0644); err != nil {
			t.Fatalf("failed to write stub: %v", err)
		}
		return
	}
	if err := os.WriteFile(filepath.Join(dir, DefaultBinary), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("failed to write stub: %v", err)
	}
}

// TestDetectRequiresGopls — пустой $PATH → false; стаб в $PATH → true.
func TestDetectRequiresGopls(t *testing.T) {
	a := &UniversalLSPAdapter{}

	t.Setenv("PATH", "")
	if a.Detect() {
		t.Error("Detect() must be false with an empty $PATH")
	}

	dir := t.TempDir()
	stubGopls(t, dir)
	t.Setenv("PATH", dir)
	if !a.Detect() {
		t.Error("Detect() must be true when gopls is in $PATH")
	}
}

// TestInstallIdempotent — Install() дважды подряд без ошибки (ARCHITECT.md §3.1);
// без gopls — ошибка с hint про ADR-026, но не паника.
func TestInstallIdempotent(t *testing.T) {
	a := &UniversalLSPAdapter{}

	t.Setenv("PATH", "")
	err := a.Install()
	if err == nil {
		t.Fatal("Install() without gopls must fail")
	}
	if !strings.Contains(err.Error(), "gopls") || !strings.Contains(err.Error(), "go install") {
		t.Errorf("error must name gopls and the install command, got: %v", err)
	}

	dir := t.TempDir()
	stubGopls(t, dir)
	t.Setenv("PATH", dir)
	if err := a.Install(); err != nil {
		t.Fatalf("first Install: %v", err)
	}
	if err := a.Install(); err != nil {
		t.Fatalf("second Install must be idempotent: %v", err)
	}
}

// TestUninstallIdempotent — двойной Stop-путь без запущенного процесса: nil, nil.
func TestUninstallIdempotent(t *testing.T) {
	a := &UniversalLSPAdapter{}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("first Uninstall (no process): %v", err)
	}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("second Uninstall must be nil: %v", err)
	}
}

// TestStartSymbolsStopLifecycle — полный жизненный цикл против fake-gopls:
// Start → Started → Symbols → Stop → Started=false → Symbols-ошибка.
func TestStartSymbolsStopLifecycle(t *testing.T) {
	t.Setenv("CAF_FAKE_GOPLS", "1")
	a := &UniversalLSPAdapter{}
	// Подменяем бинарник: Start использует DefaultBinary, поэтому кладём
	// стаб test-бинарника под именем gopls во временный PATH.
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		// .bat-обёртка над test-бинарником с env-переменной.
		script := "@echo off\r\nset CAF_FAKE_GOPLS=1\r\n\"" + os.Args[0] + "\"\r\n"
		if err := os.WriteFile(filepath.Join(dir, DefaultBinary+".bat"), []byte(script), 0644); err != nil {
			t.Fatal(err)
		}
	} else {
		script := "#!/bin/sh\nexport CAF_FAKE_GOPLS=1\nexec \"" + os.Args[0] + "\"\n"
		if err := os.WriteFile(filepath.Join(dir, DefaultBinary), []byte(script), 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir)

	// Start требует контекста — берём фоновый с разумным лимитом через Start.
	// (адаптер сам оборачивает handshake в ctx вызывающего; здесь — тестовый)
	if err := startAdapter(t, a); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if !a.Started() {
		t.Fatal("Started() = false after Start")
	}

	symbols, err := a.Symbols(testCtx(t))
	if err != nil {
		t.Fatalf("Symbols: %v", err)
	}
	if len(symbols) != 3 {
		t.Errorf("got %d symbols, want 3", len(symbols))
	}

	if err := a.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if a.Started() {
		t.Error("Started() = true after Stop")
	}
	if _, err := a.Symbols(testCtx(t)); err == nil {
		t.Error("Symbols after Stop must fail")
	}
	if err := a.Stop(); err != nil {
		t.Errorf("second Stop must be nil: %v", err)
	}
}

// TestCapabilitiesMatchContract — матрица ARCHITECT.md §1 «Universal (Tier 4)».
func TestCapabilitiesMatchContract(t *testing.T) {
	a := &UniversalLSPAdapter{}
	want := caf.CapabilityMatrix{
		PreSession: caf.TierWrapper,
		PostEdit:   caf.TierWatcher,
		PreCommit:  caf.TierUniversal,
		MCP:        false,
		Tier:       4,
	}
	if got := a.Capabilities(); got != want {
		t.Errorf("capability matrix drifted:\n got %+v\nwant %+v", got, want)
	}
	if a.ID() != "universal-lsp" {
		t.Errorf("ID() = %q, want universal-lsp", a.ID())
	}
}

// TestSymbolKindName — граничные значения kind.
func TestSymbolKindName(t *testing.T) {
	tests := []struct {
		kind int
		want string
	}{
		{0, "symbol-0"},        // невалидный (LSP нумерует с 1)
		{1, "file"},            // нижняя граница валидного
		{12, "function"},       // типичная функция
		{26, "type-parameter"}, // верхняя граница LSP 3.17
		{27, "symbol-27"},      // вне диапазона
		{-1, "symbol--1"},      // отрицательный → formatted, не паника
	}
	for _, tt := range tests {
		if got := symbolKindName(tt.kind); got != tt.want {
			t.Errorf("symbolKindName(%d) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}
