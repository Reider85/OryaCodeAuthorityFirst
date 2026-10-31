package hooks

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf/storage"
)

// captureLog перенаправляет стандартный logger в буфер на время fn и
// возвращает записанный вывод — логи и есть контракт advisory-хука.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)
	fn()
	return buf.String()
}

// seedSession записывает тестовую сессию в <repo>/.caf/sessions.
func seedSession(t *testing.T, repoRoot, sessionID string) {
	t.Helper()
	m := storage.NewManager(filepath.Join(repoRoot, ".caf", "sessions"))
	s := &storage.Session{
		SessionID: sessionID,
		Repo:      "test/repo",
		StartedAt: time.Now(),
		Agent:     "opencode",
		Model:     "glm-4.6",
	}
	if err := m.SaveSession(s); err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}
}

func TestExtractCAFSessionTrailer(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{"absent", "feat: add rate limiter", ""},
		{"trivial only", "fix: typo\n\nCAF-Trivial: true", ""},
		{"present", "feat: x\n\nCAF-Session: 2026-10-24-a1b2c3", "2026-10-24-a1b2c3"},
		{"present with whitespace", "feat: x\n\n  CAF-Session:   2026-10-24-zzz  \nbody", "2026-10-24-zzz"},
		{"colon without space", "feat: x\nCAF-Session:abc", "abc"},
		{"not a trailer prefix", "CAF-Sessionless: 1", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractCAFSessionTrailer(tt.msg); got != tt.want {
				t.Errorf("extractCAFSessionTrailer() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"trivial bypass", "chore: bump\n\nCAF-Trivial: true", false},
		{"session trailer", "feat: x\n\nCAF-Session: 2026-10-24-a1b2c3", false},
		{"no trailers", "feat: x", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommitMessage(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommitMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddCAFSessionTrailer(t *testing.T) {
	got := AddCAFSessionTrailer("feat: x", "2026-10-24-a1b2c3")
	if !strings.Contains(got, "CAF-Session: 2026-10-24-a1b2c3") {
		t.Errorf("trailer not appended: %q", got)
	}
	again := AddCAFSessionTrailer(got, "other-session")
	if again != got {
		t.Errorf("AddCAFSessionTrailer must be idempotent, got %q", again)
	}
}

func TestAddCAFTrivialTrailer(t *testing.T) {
	got := AddCAFTrivialTrailer("chore: bump")
	if !strings.Contains(got, "CAF-Trivial: true") {
		t.Errorf("trailer not appended: %q", got)
	}
	again := AddCAFTrivialTrailer(got)
	if again != got {
		t.Errorf("AddCAFTrivialTrailer must be idempotent, got %q", again)
	}
	// Любой существующий CAF-трейлер означает «не добавляй».
	withSession := "feat: x\n\nCAF-Session: 2026-10-24-a1b2c3"
	if out := AddCAFTrivialTrailer(withSession); out != withSession {
		t.Errorf("must not add CAF-Trivial when a CAF trailer exists, got %q", out)
	}
}

// TestPreCommitHookNeverBlocks — ядро MVP-контракта: никакой набор
// входных данных не возвращает непустую ошибку (INV-ARCH-004, SC-8).
func TestPreCommitHookNeverBlocks(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{"empty message (git-driven)", ""},
		{"whitespace only", "   \n\t"},
		{"trivial trailer", "chore: bump\n\nCAF-Trivial: true"},
		{"no trailers at all", "feat: add rate limiter"},
		{"session trailer, session missing", "feat: x\n\nCAF-Session: 2026-10-24-nope"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			err := PreCommitHook(HookContext{RepoRoot: repo, EventType: "pre-commit"}, tt.msg)
			if err != nil {
				t.Errorf("PreCommitHook must never block (MVP advisory), got error: %v", err)
			}
		})
	}
}

func TestPreCommitHookAdvisoryLogMessages(t *testing.T) {
	repo := t.TempDir()

	t.Run("message unavailable is explained, not silent", func(t *testing.T) {
		out := captureLog(t, func() {
			_ = PreCommitHook(HookContext{RepoRoot: repo}, "")
		})
		if !strings.Contains(out, "commit message unavailable") {
			t.Errorf("expected advisory log about missing message, got: %s", out)
		}
	})

	t.Run("missing trailer is advisory, not blocking", func(t *testing.T) {
		out := captureLog(t, func() {
			_ = PreCommitHook(HookContext{RepoRoot: repo}, "feat: x")
		})
		if !strings.Contains(out, "commit NOT blocked") {
			t.Errorf("expected advisory log about missing trailer, got: %s", out)
		}
	})

	t.Run("existing session is verified", func(t *testing.T) {
		seedSession(t, repo, "2026-10-24-a1b2c3")
		out := captureLog(t, func() {
			_ = PreCommitHook(HookContext{RepoRoot: repo}, "feat: x\n\nCAF-Session: 2026-10-24-a1b2c3")
		})
		if !strings.Contains(out, "verified") {
			t.Errorf("expected verification success log, got: %s", out)
		}
	})
}

// TestPreCommitCommandNoRecursion — регрессия бага шага 6: команда
// вызывает PreCommitHook напрямую, а не ExecuteHook (перезапуск caf).
func TestPreCommitCommandNoRecursion(t *testing.T) {
	repo := t.TempDir()
	err := captureLogErr(t, func() error {
		return PreCommitCommand(repo, "", "feat: x")
	})
	if err != nil {
		t.Errorf("PreCommitCommand returned error: %v", err)
	}
}

// captureLogErr — обёртка captureLog для вызовов, возвращающих error.
func captureLogErr(t *testing.T, fn func() error) error {
	t.Helper()
	var err error
	_ = captureLog(t, func() { err = fn() })
	return err
}
