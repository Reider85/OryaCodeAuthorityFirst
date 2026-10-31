package hooks

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPreSessionCommandFreshRepo — свежий temp-репозиторий: команда
// вызывает PreSessionHook напрямую (регрессия рекурсии шага 6),
// попытка генерации context.md не блокирует вызов.
func TestPreSessionCommandFreshRepo(t *testing.T) {
	repo := t.TempDir()
	if err := PreSessionCommand(repo, "2026-10-24-a1b2c3"); err != nil {
		t.Errorf("PreSessionCommand must be fail-safe, got %v", err)
	}
}

// TestPreSessionHookExistingContext — context.md уже есть: генерация
// пропускается (ветка os.Stat == nil).
func TestPreSessionHookExistingContext(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".caf"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".caf", "context.md"), []byte("# context\n"), 0644); err != nil {
		t.Fatal(err)
	}

	err := PreSessionHook(HookContext{RepoRoot: repo, SessionID: "2026-10-24-a1b2c3", EventType: "pre-session"})
	if err != nil {
		t.Errorf("PreSessionHook must be fail-safe, got %v", err)
	}
}

// TestPreSessionHookBrokenConfig — битый .caf/config.yml: warn-and-continue.
func TestPreSessionHookBrokenConfig(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".caf"), 0755); err != nil {
		t.Fatal(err)
	}
	// Строка, начинающаяся с «!», отклоняется SafeParseYAML → ошибка загрузки.
	if err := os.WriteFile(filepath.Join(repo, ".caf", "config.yml"), []byte("!custom-tag\n"), 0644); err != nil {
		t.Fatal(err)
	}

	err := PreSessionHook(HookContext{RepoRoot: repo, EventType: "pre-session"})
	if err != nil {
		t.Errorf("broken config must not fail the hook, got %v", err)
	}
}
