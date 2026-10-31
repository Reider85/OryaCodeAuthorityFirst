package doctor

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// writeGoplsStub кладёт фейковый gopls в temp-каталог и ставит его первым
// в $PATH — чтобы exec.LookPath("gopls") находил его детерминированно
// (паттерн writeCafStub из hooks-тестов).
func writeGoplsStub(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		stub := filepath.Join(dir, "gopls.bat")
		if err := os.WriteFile(stub, []byte("@echo off\r\nexit /b 0\r\n"), 0755); err != nil {
			t.Fatal(err)
		}
	} else {
		stub := filepath.Join(dir, "gopls")
		if err := os.WriteFile(stub, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// installTestHook создаёт .git/hooks/pre-commit в temp-репо. На unix
// IsHookInstalled требует exec-бит, на Windows — только факт файла.
func installTestHook(t *testing.T, repoRoot string) {
	t.Helper()
	hookPath := filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
	if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\n# installed by caf init\n"), 0755); err != nil {
		t.Fatal(err)
	}
}

// writeConfig — валидный или (при version="") неполный config.yml.
func writeConfig(t *testing.T, repoRoot, version string) {
	t.Helper()
	content := "version: \"" + version + "\"\n"
	if err := os.WriteFile(filepath.Join(repoRoot, ".caf", "config.yml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeSessionFile кладёт .yml в .caf/sessions — ListSessions считает
// файлы по расширению, содержимое не парсится.
func writeSessionFile(t *testing.T, repoRoot, name string) {
	t.Helper()
	sessionsDir := filepath.Join(storage.GetCAFDir(repoRoot), "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionsDir, name+".yml"), []byte("session_id: "+name+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeSymbolGraph — кэш символом в .caf/mental-model/graph.yml, чтобы
// GetManagerForRepo вернул менеджер без старта LSP (lazy-load, INV-ARCH-003).
func writeSymbolGraph(t *testing.T, repoRoot string) {
	t.Helper()
	graphFile := filepath.Join(storage.GetCAFDir(repoRoot), "mental-model", "graph.yml")
	if err := os.MkdirAll(filepath.Dir(graphFile), 0755); err != nil {
		t.Fatal(err)
	}
	if err := storage.WriteYAML(graphFile, caf.SymbolGraph{
		Symbols: []caf.Symbol{{Name: "Cached", Kind: "function", Location: "x.go:1:1"}},
	}); err != nil {
		t.Fatal(err)
	}
}

// INV-PROMPT-MVP-006: каждая проверка имеет 4 состояния; здесь фиксируем
// ветки checkCAFValidity: error(нет .caf) / warning(нет subdir, нет manifest) / healthy.
func TestCheckCAFValidity(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, repo string)
		wantStatus  string
		wantContain string
	}{
		{
			name:        "no .caf directory",
			setup:       func(t *testing.T, repo string) {},
			wantStatus:  "error",
			wantContain: ".caf/ directory not found",
		},
		{
			name: "missing subdirectory",
			setup: func(t *testing.T, repo string) {
				if err := os.MkdirAll(filepath.Join(repo, ".caf"), 0755); err != nil {
					t.Fatal(err)
				}
			},
			wantStatus:  "warning",
			wantContain: "Missing subdirectory",
		},
		{
			name: "missing manifest",
			setup: func(t *testing.T, repo string) {
				cafDir := storage.GetCAFDir(repo)
				for _, sub := range []string{"sessions", "decisions", "mental-model", "cache"} {
					if err := os.MkdirAll(filepath.Join(cafDir, sub), 0755); err != nil {
						t.Fatal(err)
					}
				}
			},
			wantStatus:  "warning",
			wantContain: "manifest.yml not found",
		},
		{
			name: "valid via InitializeCAF",
			setup: func(t *testing.T, repo string) {
				if err := storage.InitializeCAF(repo); err != nil {
					t.Fatal(err)
				}
			},
			wantStatus:  "healthy",
			wantContain: ".caf/ directory is valid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			tt.setup(t, repo)

			check := NewDoctor(repo).checkCAFValidity()
			if check.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q (message: %s)", check.Status, tt.wantStatus, check.Message)
			}
			if !strings.Contains(check.Message, tt.wantContain) {
				t.Errorf("Message = %q, want contain %q", check.Message, tt.wantContain)
			}
			if check.Name != "CAF directory validity" {
				t.Errorf("Name = %q", check.Name)
			}
		})
	}
}

// FR-OBS-004: config schema — warning на отсутствии/пустом version,
// error на битом YAML, healthy на валидном.
func TestCheckConfigSchema(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, repo string)
		wantStatus  string
		wantContain string
	}{
		{
			name:        "config missing",
			setup:       func(t *testing.T, repo string) {},
			wantStatus:  "warning",
			wantContain: "config.yml not found",
		},
		{
			name: "invalid yaml",
			setup: func(t *testing.T, repo string) {
				cafDir := storage.GetCAFDir(repo)
				if err := os.MkdirAll(cafDir, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(cafDir, "config.yml"), []byte("invalid: yaml: content:"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			wantStatus:  "error",
			wantContain: "Invalid config.yml schema",
		},
		{
			name: "missing version field",
			setup: func(t *testing.T, repo string) {
				if err := os.MkdirAll(storage.GetCAFDir(repo), 0755); err != nil {
					t.Fatal(err)
				}
				writeConfig(t, repo, "")
			},
			wantStatus:  "warning",
			wantContain: "Missing version field",
		},
		{
			name: "valid config",
			setup: func(t *testing.T, repo string) {
				if err := os.MkdirAll(storage.GetCAFDir(repo), 0755); err != nil {
					t.Fatal(err)
				}
				writeConfig(t, repo, "0.4.1")
			},
			wantStatus:  "healthy",
			wantContain: "config.yml is valid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			tt.setup(t, repo)

			check := NewDoctor(repo).checkConfigSchema()
			if check.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q (message: %s)", check.Status, tt.wantStatus, check.Message)
			}
			if !strings.Contains(check.Message, tt.wantContain) {
				t.Errorf("Message = %q, want contain %q", check.Message, tt.wantContain)
			}
		})
	}
}

// Ноль сессий — unknown, не error: новый репозиторий ≠ сломанный.
func TestCheckSessionsPresent(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, repo string)
		wantStatus  string
		wantContain string
	}{
		{
			name:        "no sessions dir",
			setup:       func(t *testing.T, repo string) {},
			wantStatus:  "unknown",
			wantContain: "No sessions yet",
		},
		{
			name: "empty sessions dir",
			setup: func(t *testing.T, repo string) {
				if err := os.MkdirAll(filepath.Join(storage.GetCAFDir(repo), "sessions"), 0755); err != nil {
					t.Fatal(err)
				}
			},
			wantStatus:  "unknown",
			wantContain: "No sessions yet",
		},
		{
			name: "one session",
			setup: func(t *testing.T, repo string) {
				if err := os.MkdirAll(storage.GetCAFDir(repo), 0755); err != nil {
					t.Fatal(err)
				}
				writeSessionFile(t, repo, "2026-10-24-a1b2c3")
			},
			wantStatus:  "healthy",
			wantContain: "Found 1 sessions",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			tt.setup(t, repo)

			check := NewDoctor(repo).checkSessionsPresent()
			if check.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q (message: %s)", check.Status, tt.wantStatus, check.Message)
			}
			if !strings.Contains(check.Message, tt.wantContain) {
				t.Errorf("Message = %q, want contain %q", check.Message, tt.wantContain)
			}
		})
	}
}

// FR-ADP-003: hook отсутствует — warning (fail-safe, не error);
// установлен — healthy.
func TestCheckHookInstalled(t *testing.T) {
	repo := t.TempDir()

	check := NewDoctor(repo).checkHookInstalled()
	if check.Status != "warning" {
		t.Errorf("no hook: Status = %q, want warning (message: %s)", check.Status, check.Message)
	}
	if check.Suggestion == "" {
		t.Error("no hook: Suggestion must point at 'caf init'")
	}

	installTestHook(t, repo)
	check = NewDoctor(repo).checkHookInstalled()
	if check.Status != "healthy" {
		t.Errorf("installed hook: Status = %q, want healthy (message: %s)", check.Status, check.Message)
	}
}

// ADR-013: gopls нет — unknown; gopls есть, но .caf нет — warning
// (GetManagerForRepo отказывает); кэш графа есть — healthy.
func TestCheckLSPResponds(t *testing.T) {
	t.Run("gopls not available", func(t *testing.T) {
		t.Setenv("PATH", "")
		repo := t.TempDir()

		check := NewDoctor(repo).checkLSPResponds()
		if check.Status != "unknown" {
			t.Errorf("Status = %q, want unknown (message: %s)", check.Status, check.Message)
		}
		if !strings.Contains(check.Message, "gopls not available") {
			t.Errorf("Message = %q", check.Message)
		}
	})

	t.Run("manager unavailable without .caf", func(t *testing.T) {
		writeGoplsStub(t)
		repo := t.TempDir()

		check := NewDoctor(repo).checkLSPResponds()
		if check.Status != "warning" {
			t.Errorf("Status = %q, want warning (message: %s, details: %s)",
				check.Status, check.Message, check.Details)
		}
	})

	t.Run("healthy from cached graph", func(t *testing.T) {
		writeGoplsStub(t)
		repo := t.TempDir()
		if err := storage.InitializeCAF(repo); err != nil {
			t.Fatal(err)
		}
		writeSymbolGraph(t, repo)

		check := NewDoctor(repo).checkLSPResponds()
		if check.Status != "healthy" {
			t.Errorf("Status = %q, want healthy (message: %s)", check.Status, check.Message)
		}
		if !strings.Contains(check.Message, "1 symbols loaded") {
			t.Errorf("Message = %q, want symbol count", check.Message)
		}
	})
}

// Run() агрегирует ровно 5 проверок и считает summary; timestamp — RFC3339.
func TestRunAggregation(t *testing.T) {
	t.Run("broken repo overall error", func(t *testing.T) {
		t.Setenv("PATH", "")
		repo := t.TempDir()

		report, err := NewDoctor(repo).Run()
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if len(report.Checks) != 5 {
			t.Fatalf("len(Checks) = %d, want 5", len(report.Checks))
		}
		if report.Overall != "error" {
			t.Errorf("Overall = %q, want error (нет .caf — check 1 error)", report.Overall)
		}
		if report.Summary != "0 passed, 1 failed, 2 warnings, 2 unknown" {
			t.Errorf("Summary = %q", report.Summary)
		}
		if _, err := time.Parse(time.RFC3339, report.Timestamp); err != nil {
			t.Errorf("Timestamp not RFC3339: %v", err)
		}
	})

	t.Run("healthy repo overall healthy", func(t *testing.T) {
		writeGoplsStub(t)
		repo := t.TempDir()
		if err := storage.InitializeCAF(repo); err != nil {
			t.Fatal(err)
		}
		writeConfig(t, repo, "0.4.1")
		writeSessionFile(t, repo, "2026-10-24-a1b2c3")
		installTestHook(t, repo)
		writeSymbolGraph(t, repo)

		report, err := RunCommand(repo)
		if err != nil {
			t.Fatalf("RunCommand: %v", err)
		}
		if report.Overall != "healthy" {
			t.Errorf("Overall = %q, want healthy (summary: %s)", report.Overall, report.Summary)
		}
		if report.Summary != "5 passed, 0 failed, 0 warnings, 0 unknown" {
			t.Errorf("Summary = %q", report.Summary)
		}
		for _, c := range report.Checks {
			if c.Status != "healthy" {
				t.Errorf("check %q = %q, want healthy (%s)", c.Name, c.Status, c.Message)
			}
		}
	})

	t.Run("warning-only repo overall warning", func(t *testing.T) {
		t.Setenv("PATH", "")
		repo := t.TempDir()
		if err := storage.InitializeCAF(repo); err != nil {
			t.Fatal(err)
		}

		report, err := NewDoctor(repo).Run()
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		// config missing → warning, sessions empty → unknown, hook missing → warning
		if report.Overall != "warning" {
			t.Errorf("Overall = %q, want warning", report.Overall)
		}
	})
}

// ARCHITECT.md §9.4: 4 иконки + Suggestion/Details; цвет не единственный
// носитель смысла (INV-PROMPT-MVP-008) — иконка дублирует статус текстом.
func TestFormatReport(t *testing.T) {
	report := &HealthReport{
		Overall:   "warning",
		Summary:   "1 passed, 1 failed, 1 warnings, 1 unknown",
		Timestamp: "2026-10-31T12:00:00+03:00",
		Checks: []HealthCheck{
			{Name: "OK check", Status: "healthy", Message: "all good"},
			{Name: "Warn check", Status: "warning", Message: "degraded", Suggestion: "run caf init"},
			{Name: "Error check", Status: "error", Message: "broken", Details: "stat failed"},
			{Name: "Unknown check", Status: "unknown", Message: "not probed"},
		},
	}

	out := FormatReport(report)

	for _, want := range []string{
		"CodeAuthorityFirst Health Report",
		"Timestamp: 2026-10-31T12:00:00+03:00",
		"Overall: warning",
		"Summary: 1 passed, 1 failed, 1 warnings, 1 unknown",
		"  ✓ OK check: all good",
		"  ⚠ Warn check: degraded",
		"    Suggestion: run caf init",
		"  ✗ Error check: broken",
		"    Details: stat failed",
		"  ? Unknown check: not probed",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}
