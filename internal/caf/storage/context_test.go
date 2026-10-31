package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewContextGenerator(t *testing.T) {
	g := NewContextGenerator("sessions", "decisions", 5, 3)
	if g == nil {
		t.Fatal("NewContextGenerator returned nil")
	}
	if g.sessionsDir != "sessions" || g.decisionsDir != "decisions" {
		t.Errorf("unexpected dirs: %q %q", g.sessionsDir, g.decisionsDir)
	}
	if g.maxSessions != 5 || g.maxDecisions != 3 {
		t.Errorf("unexpected limits: %d %d", g.maxSessions, g.maxDecisions)
	}
}

func TestGenerateContextHappyPath(t *testing.T) {
	tempDir := t.TempDir()
	repoRoot := filepath.Join(tempDir, "repo")
	if _, err := EnsureCAFDir(repoRoot); err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}
	sessionsDir := filepath.Join(repoRoot, ".caf", "sessions")
	decisionsDir := filepath.Join(repoRoot, ".caf", "decisions")

	older := &Session{
		SessionID:     "2026-10-24-aaaaaa",
		Agent:         "opencode",
		Model:         "glm-4.6",
		StartedAt:     time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC),
		InitialPrompt: "older prompt",
		Iterations:    []Iteration{{ID: "i1", Action: "edit", Target: "a.go", Diff: "added func"}},
		FilesTouched:  []string{"a.go"},
		DecisionsMade: []Decision{{ID: "d1", Title: "older decision", Timestamp: time.Date(2026, 10, 24, 9, 30, 0, 0, time.UTC)}},
	}
	newer := &Session{
		SessionID:     "2026-10-25-bbbbbb",
		Agent:         "claude-code",
		Model:         "sonnet",
		StartedAt:     time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC),
		InitialPrompt: "newer prompt",
	}
	manager := NewManager(sessionsDir)
	for _, s := range []*Session{older, newer} {
		if err := manager.SaveSession(s); err != nil {
			t.Fatalf("SaveSession %s: %v", s.SessionID, err)
		}
	}
	if err := os.WriteFile(filepath.Join(sessionsDir, "corrupt.yml"), []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}

	if err := WriteYAML(filepath.Join(decisionsDir, "index.yml"), map[string]interface{}{
		"decisions": []string{"d_2", "d_1", "d_missing"},
	}); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := WriteYAML(filepath.Join(decisionsDir, "d_1.yml"), Decision{
		ID: "d_1", Title: "first decision", Description: "desc one",
		Timestamp: time.Date(2026, 10, 24, 8, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("write d_1: %v", err)
	}
	if err := WriteYAML(filepath.Join(decisionsDir, "d_2.yml"), Decision{
		ID: "d_2", Title: "second decision", Description: "desc two",
		Timestamp: time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("write d_2: %v", err)
	}

	gen := NewContextGenerator(sessionsDir, decisionsDir, 10, 10)
	if err := gen.GenerateContext(repoRoot); err != nil {
		t.Fatalf("GenerateContext: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(repoRoot, ".caf", "context.md"))
	if err != nil {
		t.Fatalf("read context.md: %v", err)
	}
	content := string(data)

	for _, want := range []string{
		"# CodeAuthorityFirst Context",
		"## Recent Sessions",
		"### Session 2026-10-25-bbbbbb (claude-code)",
		"newer prompt",
		"### Session 2026-10-24-aaaaaa (opencode)",
		"older prompt",
		"edit: a.go",
		"a.go",
		"## Recent Decisions",
		"### d_2: second decision",
		"### d_1: first decision",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("context.md missing %q", want)
		}
	}

	newerIdx := strings.Index(content, "2026-10-25-bbbbbb")
	olderIdx := strings.Index(content, "2026-10-24-aaaaaa")
	if newerIdx > olderIdx {
		t.Error("sessions not sorted newest-first")
	}
	d2Idx := strings.Index(content, "d_2: second decision")
	d1Idx := strings.Index(content, "d_1: first decision")
	if d2Idx > d1Idx {
		t.Error("decisions not sorted newest-first")
	}
	if strings.Contains(content, "d_missing") {
		t.Error("missing decision file should be skipped")
	}
}

func TestGenerateContextEmpty(t *testing.T) {
	tempDir := t.TempDir()
	repoRoot := filepath.Join(tempDir, "repo")
	if _, err := EnsureCAFDir(repoRoot); err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}

	gen := NewContextGenerator(
		filepath.Join(repoRoot, ".caf", "sessions"),
		filepath.Join(repoRoot, ".caf", "decisions"),
		10, 10,
	)
	if err := gen.GenerateContext(repoRoot); err != nil {
		t.Fatalf("GenerateContext on empty dirs: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(repoRoot, ".caf", "context.md"))
	if err != nil {
		t.Fatalf("read context.md: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "## Recent Sessions") {
		t.Error("missing sessions section")
	}
	if !strings.Contains(content, "## Recent Decisions") {
		t.Error("missing decisions section")
	}
	if strings.Contains(content, "### Session") {
		t.Error("no sessions expected")
	}
}

func TestGenerateContextMaxSessionsCap(t *testing.T) {
	tempDir := t.TempDir()
	repoRoot := filepath.Join(tempDir, "repo")
	if _, err := EnsureCAFDir(repoRoot); err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}
	sessionsDir := filepath.Join(repoRoot, ".caf", "sessions")
	manager := NewManager(sessionsDir)
	for i, id := range []string{"2026-10-24-aaaaaa", "2026-10-25-bbbbbb", "2026-10-26-cccccc"} {
		err := manager.SaveSession(&Session{
			SessionID: id,
			Agent:     "opencode",
			StartedAt: time.Date(2026, 10, 24+i, 9, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("SaveSession %s: %v", id, err)
		}
	}

	gen := NewContextGenerator(sessionsDir, filepath.Join(repoRoot, ".caf", "decisions"), 2, 10)
	if err := gen.GenerateContext(repoRoot); err != nil {
		t.Fatalf("GenerateContext: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(repoRoot, ".caf", "context.md"))
	if err != nil {
		t.Fatalf("read context.md: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "2026-10-26-cccccc") || !strings.Contains(content, "2026-10-25-bbbbbb") {
		t.Error("expected two newest sessions")
	}
	if strings.Contains(content, "2026-10-24-aaaaaa") {
		t.Error("oldest session should be capped out")
	}
}

func TestGenerateContextSessionsDirUnreadable(t *testing.T) {
	tempDir := t.TempDir()
	repoRoot := filepath.Join(tempDir, "repo")
	if _, err := EnsureCAFDir(repoRoot); err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}

	// NUL в пути даёт не-ENOENT ошибку при чтении каталога
	blocked := filepath.Join(tempDir, "bad\x00name")

	gen := NewContextGenerator(blocked, filepath.Join(repoRoot, ".caf", "decisions"), 10, 10)
	err := gen.GenerateContext(repoRoot)
	if err == nil {
		t.Fatal("GenerateContext should fail when sessions dir cannot be read")
	}
	if !strings.Contains(err.Error(), "failed to list recent sessions") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGenerateContextBadDecisionsIndex(t *testing.T) {
	tempDir := t.TempDir()
	repoRoot := filepath.Join(tempDir, "repo")
	if _, err := EnsureCAFDir(repoRoot); err != nil {
		t.Fatalf("EnsureCAFDir: %v", err)
	}
	decisionsDir := filepath.Join(repoRoot, ".caf", "decisions")
	if err := os.WriteFile(filepath.Join(decisionsDir, "index.yml"), []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("write bad index: %v", err)
	}

	gen := NewContextGenerator(filepath.Join(repoRoot, ".caf", "sessions"), decisionsDir, 10, 10)
	err := gen.GenerateContext(repoRoot)
	if err == nil {
		t.Fatal("GenerateContext should fail on corrupted decisions index")
	}
	if !strings.Contains(err.Error(), "failed to list recent decisions") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGenerateContextWriteError(t *testing.T) {
	tempDir := t.TempDir()
	// repoRoot is an existing file — writing <repoRoot>/.caf/context.md must fail
	repoRoot := filepath.Join(tempDir, "repo-file")
	if err := os.WriteFile(repoRoot, []byte("x"), 0644); err != nil {
		t.Fatalf("create repo file: %v", err)
	}

	gen := NewContextGenerator(
		filepath.Join(tempDir, "sessions"),
		filepath.Join(tempDir, "decisions"),
		10, 10,
	)
	err := gen.GenerateContext(repoRoot)
	if err == nil {
		t.Fatal("GenerateContext should fail when context.md cannot be written")
	}
	if !strings.Contains(err.Error(), "failed to write context file") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestListRecentDecisionsMissingIndex(t *testing.T) {
	tempDir := t.TempDir()
	g := NewContextGenerator(tempDir, filepath.Join(tempDir, "decisions"), 10, 10)
	decisions, err := g.listRecentDecisions()
	if err != nil {
		t.Fatalf("listRecentDecisions: %v", err)
	}
	if len(decisions) != 0 {
		t.Errorf("expected no decisions, got %d", len(decisions))
	}
}

func TestListRecentSessionsSkipsCorrupt(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)
	if err := manager.SaveSession(&Session{
		SessionID: "2026-10-24-aaaaaa",
		Agent:     "opencode",
		StartedAt: time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("SaveSession: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "corrupt.yml"), []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}

	g := NewContextGenerator(tempDir, filepath.Join(tempDir, "decisions"), 10, 10)
	ids, err := g.listRecentSessions()
	if err != nil {
		t.Fatalf("listRecentSessions: %v", err)
	}
	if len(ids) != 1 || ids[0] != "2026-10-24-aaaaaa" {
		t.Errorf("expected only valid session, got %v", ids)
	}
}
