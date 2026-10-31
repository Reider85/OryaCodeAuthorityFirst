package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// TestPostEditHookTrivialChangeSkips — тривиальное изменение не пишет
// сессию (ACV: порог triviality, ADR-004).
func TestPostEditHookTrivialChangeSkips(t *testing.T) {
	repo := t.TempDir()
	ctx := HookContext{RepoRoot: repo, SessionID: "2026-10-24-a1b2c3", EventType: "post-edit"}

	err := PostEditHook(ctx, []string{"main.go"}, "trivial comment tweak")
	if err != nil {
		t.Errorf("trivial change must be fail-safe, got %v", err)
	}

	m := storage.NewManager(filepath.Join(repo, ".caf", "sessions"))
	ids, listErr := m.ListSessions()
	if listErr != nil {
		t.Fatal(listErr)
	}
	if len(ids) != 0 {
		t.Errorf("trivial change must not create a session, got %v", ids)
	}
}

// TestPostEditHookCreatesAndAppendsSession — нетривиальное изменение
// создаёт сессию (i1), повторный вызов добавляет итерацию (i2).
func TestPostEditHookCreatesAndAppendsSession(t *testing.T) {
	repo := t.TempDir()
	ctx := HookContext{RepoRoot: repo, SessionID: "2026-10-24-a1b2c3", EventType: "post-edit"}

	if err := PostEditHook(ctx, []string{"internal/api/policies.go"}, "added middleware RateLimit(token, 10/min)"); err != nil {
		t.Fatalf("first PostEditHook failed: %v", err)
	}

	m := storage.NewManager(filepath.Join(repo, ".caf", "sessions"))
	s, err := m.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatalf("session not written: %v", err)
	}
	if len(s.Iterations) != 1 || s.Iterations[0].ID != "i1" || s.Iterations[0].Action != "edit" {
		t.Errorf("expected single i1/edit iteration, got %+v", s.Iterations)
	}
	if len(s.FilesTouched) != 1 || s.FilesTouched[0] != "internal/api/policies.go" {
		t.Errorf("files_touched wrong: %v", s.FilesTouched)
	}

	if err := PostEditHook(ctx, []string{"internal/api/middleware.go"}, "added rate limit bucket"); err != nil {
		t.Fatalf("second PostEditHook failed: %v", err)
	}
	s, err = m.LoadSession("2026-10-24-a1b2c3")
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Iterations) != 2 || s.Iterations[1].ID != "i2" {
		t.Errorf("expected i2 appended, got %+v", s.Iterations)
	}
}

// TestPostEditHookBrokenConfigNoPanic — регрессия fail-safe: битый
// config.yml раньше ронял хук nil-указателем (INV-ARCH-004).
func TestPostEditHookBrokenConfigNoPanic(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".caf"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".caf", "config.yml"), []byte("!custom-tag\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := HookContext{RepoRoot: repo, SessionID: "2026-10-24-a1b2c3", EventType: "post-edit"}
	err := PostEditHook(ctx, []string{"main.go"}, "refactored auth flow")
	if err != nil {
		t.Errorf("broken config must not fail the hook, got %v", err)
	}
}

func TestPostEditCommand(t *testing.T) {
	repo := t.TempDir()
	err := PostEditCommand(repo, "2026-10-24-a1b2c3", []string{"main.go"}, "trivial")
	if err != nil {
		t.Errorf("PostEditCommand must be fail-safe, got %v", err)
	}
}

func TestIsTrivialChange(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		summary string
		want    bool
	}{
		{"summary mentions trivial", []string{"a.go"}, "trivial formatting fix", true},
		{"summary mentions trivial case-insensitive", []string{"a.go"}, "TRIVIAL", true},
		{"many files, no marker", []string{"a.go", "b.go", "c.go"}, "renamed", false},
		{"few files, no marker", []string{"a.go"}, "real change", false},
		{"empty everything", nil, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTrivialChange(tt.files, tt.summary, 50); got != tt.want {
				t.Errorf("isTrivialChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPostEditSessionHelpers — детерминированные helper-функции (шаблон
// YAML-схемы §2.1) возвращают стабильные значения, а не паникуют.
func TestPostEditSessionHelpers(t *testing.T) {
	if got := getRepoName("/tmp/repo"); got != "unknown" {
		t.Errorf("getRepoName() = %q", got)
	}
	if getStartTime("2026-10-24-a1b2c3").IsZero() {
		t.Error("getStartTime() must return a time")
	}
	if getEndTime().IsZero() {
		t.Error("getEndTime() must return a time")
	}
	if got := getAgentFromSessionID("2026-10-24-opencode"); got != "opencode" {
		t.Errorf("getAgentFromSessionID() = %q", got)
	}
	if got := getAgentFromSessionID("2026-10-24-claude"); got != "claude-code" {
		t.Errorf("getAgentFromSessionID() = %q", got)
	}
	if got := getAgentFromSessionID("2026-10-24-other"); got != "unknown" {
		t.Errorf("getAgentFromSessionID() = %q", got)
	}
	if !strings.Contains(getToolCallFromSessionID("x"), "rate limiter") {
		t.Errorf("getToolCallFromSessionID() = %q", getToolCallFromSessionID("x"))
	}
	if getInitialPromptFromSessionID("x") == "" {
		t.Error("getInitialPromptFromSessionID() must not be empty")
	}
	if getVoiceSignatureFromSessionID("x") == "" {
		t.Error("getVoiceSignatureFromSessionID() must not be empty")
	}
}
