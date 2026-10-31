package claudecode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// settings — упрощённое зеркало settings.json для проверок в тестах.
type settings struct {
	Model string `json:"model"`
	Hooks map[string][]struct {
		Matcher string `json:"matcher"`
		Hooks   []struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		} `json:"hooks"`
	} `json:"hooks"`
}

// withTempHome подменяет домашний каталог, чтобы тесты никогда не трогали
// реальный ~/.claude/settings.json.
func withTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	return home
}

// stubBinary кладёт исполняемый стаб бинарника агента в dir (Detect ищет
// его в $PATH; на Windows LookPath берёт .bat из PATHEXT).
func stubBinary(t *testing.T, dir, name string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		if err := os.WriteFile(filepath.Join(dir, name+".bat"), []byte("@echo off\r\n"), 0644); err != nil {
			t.Fatalf("failed to write stub: %v", err)
		}
		return
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("failed to write stub: %v", err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(data)
}

// writeSettingsRaw кладёт произвольный контент в ~/.claude/settings.json.
func writeSettingsRaw(t *testing.T, body string) string {
	t.Helper()
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

// loadSettings декодирует settings.json в тестовую структуру.
func loadSettings(t *testing.T, path string) settings {
	t.Helper()
	var doc settings
	if err := json.Unmarshal([]byte(mustRead(t, path)), &doc); err != nil {
		t.Fatalf("settings.json is not valid JSON: %v", err)
	}
	return doc
}

// cafEntryCount считает hook-записи, принадлежащие caf, по всем событиям.
func cafEntryCount(doc settings) int {
	count := 0
	for _, entries := range doc.Hooks {
		for _, e := range entries {
			for _, h := range e.Hooks {
				if strings.Contains(h.Command, cafCommandPrefix) {
					count++
					break
				}
			}
		}
	}
	return count
}

func TestIDAndCapabilities(t *testing.T) {
	a := &ClaudeCodeAdapter{}
	if a.ID() != "claude-code" {
		t.Errorf("ID() = %q, want claude-code", a.ID())
	}
	// ARCHITECT.md §1: Claude Code — Tier 1, все события native, MCP да (FR-ADP-002).
	want := caf.CapabilityMatrix{
		PreSession: caf.TierNative,
		PostEdit:   caf.TierNative,
		PreCommit:  caf.TierNative,
		MCP:        true,
		Tier:       1,
	}
	if got := a.Capabilities(); got != want {
		t.Errorf("Capabilities() = %+v, want %+v", got, want)
	}
	if err := caf.ValidateCapabilityMatrix(a.Capabilities()); err != nil {
		t.Errorf("ValidateCapabilityMatrix: %v", err)
	}
}

func TestConfigPathUsesHome(t *testing.T) {
	home := withTempHome(t)
	got, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: %v", err)
	}
	want := filepath.Join(home, ".claude", "settings.json")
	if got != want {
		t.Errorf("ConfigPath() = %q, want %q", got, want)
	}
}

func TestDetectFindsBinaryInPath(t *testing.T) {
	dir := t.TempDir()
	stubBinary(t, dir, "claude")
	t.Setenv("PATH", dir)
	if !(&ClaudeCodeAdapter{}).Detect() {
		t.Error("Detect() = false, want true when claude is in $PATH")
	}
}

func TestDetectFalseWithoutBinary(t *testing.T) {
	t.Setenv("PATH", "")
	if (&ClaudeCodeAdapter{}).Detect() {
		t.Error("Detect() = true, want false when $PATH is empty")
	}
}

func TestInstallWritesCanonicalHooks(t *testing.T) {
	withTempHome(t)
	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	path, _ := ConfigPath()
	doc := loadSettings(t, path)

	if len(doc.Hooks["SessionStart"]) != 1 || doc.Hooks["SessionStart"][0].Hooks[0].Command != "caf hook pre-session" {
		t.Errorf("SessionStart hook missing/incorrect: %+v", doc.Hooks["SessionStart"])
	}
	if len(doc.Hooks["PostToolUse"]) != 1 ||
		doc.Hooks["PostToolUse"][0].Matcher != "Write|Edit|MultiEdit" ||
		doc.Hooks["PostToolUse"][0].Hooks[0].Command != "caf hook post-edit" {
		t.Errorf("PostToolUse hook missing/incorrect: %+v", doc.Hooks["PostToolUse"])
	}
	// pre-commit подключён через PreToolUse на Bash с stdin-фильтром —
	// Claude Code не имеет отдельного pre-commit события (шаг 7).
	if len(doc.Hooks["PreToolUse"]) != 1 ||
		doc.Hooks["PreToolUse"][0].Matcher != "Bash" ||
		doc.Hooks["PreToolUse"][0].Hooks[0].Command != "caf hook pre-commit --from-agent" {
		t.Errorf("PreToolUse hook missing/incorrect: %+v", doc.Hooks["PreToolUse"])
	}
	for _, e := range doc.Hooks["SessionStart"] {
		if e.Hooks[0].Type != "command" {
			t.Errorf("hook type must be %q, got %q", "command", e.Hooks[0].Type)
		}
	}
	if !a.Configured() {
		t.Error("Configured() = false after Install, want true")
	}
}

func TestInstallIsIdempotent(t *testing.T) {
	withTempHome(t)
	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("first Install: %v", err)
	}
	path, _ := ConfigPath()
	first := mustRead(t, path)

	if err := a.Install(); err != nil {
		t.Fatalf("second Install must be idempotent: %v", err)
	}
	if second := mustRead(t, path); first != second {
		t.Errorf("Install is not idempotent:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if got := cafEntryCount(loadSettings(t, path)); got != 3 {
		t.Errorf("expected exactly 3 caf hook entries, got %d", got)
	}
}

func TestInstallPreservesForeignKeysAndHooks(t *testing.T) {
	withTempHome(t)
	foreign := `{
  "model": "opus",
  "hooks": {
    "Stop": [
      { "hooks": [ { "type": "command", "command": "echo bye" } ] }
    ]
  }
}`
	path := writeSettingsRaw(t, foreign)

	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}

	doc := loadSettings(t, path)
	if doc.Model != "opus" {
		t.Errorf("foreign key \"model\" lost: %q", doc.Model)
	}
	if len(doc.Hooks["Stop"]) != 1 || doc.Hooks["Stop"][0].Hooks[0].Command != "echo bye" {
		t.Errorf("foreign Stop hook lost: %+v", doc.Hooks["Stop"])
	}
	if got := cafEntryCount(doc); got != 3 {
		t.Errorf("expected 3 caf entries, got %d", got)
	}
	if !strings.Contains(mustRead(t, path), "http") && strings.Contains(foreign, "http") {
		t.Error("HTML escaping must be disabled for foreign content")
	}
}

func TestInstallReplacesOwnEntriesWithoutDuplicates(t *testing.T) {
	withTempHome(t)
	// Старая установка (например, с другим matcher) + чужая запись.
	stale := `{
  "hooks": {
    "PostToolUse": [
      { "matcher": "Write", "hooks": [ { "type": "command", "command": "caf hook post-edit --legacy" } ] },
      { "matcher": "NotebookEdit", "hooks": [ { "type": "command", "command": "echo foreign" } ] }
    ]
  }
}`
	path := writeSettingsRaw(t, stale)

	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}

	doc := loadSettings(t, path)
	if got := cafEntryCount(doc); got != 3 {
		t.Fatalf("expected 3 caf entries after reinstall, got %d", got)
	}
	for _, e := range doc.Hooks["PostToolUse"] {
		if strings.Contains(e.Hooks[0].Command, "--legacy") {
			t.Errorf("stale caf entry survived: %+v", e)
		}
	}
	foreignKept := false
	for _, e := range doc.Hooks["PostToolUse"] {
		if e.Hooks[0].Command == "echo foreign" && e.Matcher == "NotebookEdit" {
			foreignKept = true
		}
	}
	if !foreignKept {
		t.Errorf("foreign hook entry was dropped: %+v", doc.Hooks["PostToolUse"])
	}
}

func TestInstallRejectsBrokenJSON(t *testing.T) {
	withTempHome(t)
	broken := "{ this is not json"
	path := writeSettingsRaw(t, broken)

	if err := (&ClaudeCodeAdapter{}).Install(); err == nil {
		t.Fatal("Install must refuse to touch invalid JSON")
	}
	if got := mustRead(t, path); got != broken {
		t.Errorf("file must stay untouched on parse error:\ngot:  %q\nwant: %q", got, broken)
	}
}

func TestInstallRejectsNonObjectHooksKey(t *testing.T) {
	withTempHome(t)
	body := `{"hooks": []}`
	path := writeSettingsRaw(t, body)

	if err := (&ClaudeCodeAdapter{}).Install(); err == nil {
		t.Fatal("Install must refuse a \"hooks\" key that is not an object")
	}
	if got := mustRead(t, path); got != body {
		t.Errorf("file must stay untouched:\ngot:  %q\nwant: %q", got, body)
	}
}

// TestInstallWithoutHomeFails — регрессия бага шага 7: раньше пустой
// Getenv("HOME") приводил к записи в относительный ".claude" от текущей
// директории (на Windows USERPROFILE и HOME пусты).
func TestInstallWithoutHomeFails(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := (&ClaudeCodeAdapter{}).Install(); err == nil {
		t.Fatal("Install must fail when no home directory can be resolved")
	}
	if _, err := os.Stat(filepath.Join(wd, ".claude")); !os.IsNotExist(err) {
		t.Errorf("Install must not create .claude under the current directory (wd=%s)", wd)
	}
}

func TestUninstallRemovesOwnFileOnly(t *testing.T) {
	withTempHome(t)
	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	path, _ := ConfigPath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Uninstall must remove the caf-owned file at %s", path)
	}
	if err := a.Uninstall(); err != nil {
		t.Errorf("second Uninstall must be nil, got %v", err)
	}
	if a.Configured() {
		t.Error("Configured() = true after Uninstall, want false")
	}
}

func TestUninstallKeepsForeignContent(t *testing.T) {
	withTempHome(t)
	foreign := `{
  "model": "opus",
  "hooks": {
    "Stop": [ { "hooks": [ { "type": "command", "command": "echo bye" } ] } ]
  }
}`
	path := writeSettingsRaw(t, foreign)

	a := &ClaudeCodeAdapter{}
	if err := a.Install(); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if err := a.Uninstall(); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	doc := loadSettings(t, path)
	if doc.Model != "opus" {
		t.Errorf("foreign key \"model\" lost: %q", doc.Model)
	}
	if len(doc.Hooks["Stop"]) != 1 {
		t.Errorf("foreign Stop hook lost: %+v", doc.Hooks)
	}
	if got := cafEntryCount(doc); got != 0 {
		t.Errorf("caf entries survived Uninstall: %d", got)
	}
	if a.Configured() {
		t.Error("Configured() = true after Uninstall, want false")
	}
}

func TestUninstallRefusesForeignHooks(t *testing.T) {
	withTempHome(t)
	foreign := `{
  "hooks": {
    "Stop": [ { "hooks": [ { "type": "command", "command": "echo bye" } ] } ]
  }
}`
	path := writeSettingsRaw(t, foreign)

	if err := (&ClaudeCodeAdapter{}).Uninstall(); err == nil {
		t.Fatal("Uninstall must refuse hooks caf did not write")
	}
	if got := mustRead(t, path); got != foreign {
		t.Errorf("refused file must stay untouched:\ngot:  %q\nwant: %q", got, foreign)
	}
}

func TestUninstallIgnoresFileWithoutCafHooks(t *testing.T) {
	withTempHome(t)
	// Ни caf-хуков, ни чужих хуков — файл не наш, но и трогать его нечего.
	path := writeSettingsRaw(t, `{"model": "opus"}`)
	before := mustRead(t, path)

	if err := (&ClaudeCodeAdapter{}).Uninstall(); err != nil {
		t.Fatalf("Uninstall of a file without caf hooks must be a no-op, got %v", err)
	}
	if got := mustRead(t, path); got != before {
		t.Errorf("file was rewritten without reason:\ngot:  %q\nwant: %q", got, before)
	}
}

func TestUninstallRejectsBrokenJSON(t *testing.T) {
	withTempHome(t)
	broken := "{ this is not json"
	path := writeSettingsRaw(t, broken)

	if err := (&ClaudeCodeAdapter{}).Uninstall(); err == nil {
		t.Fatal("Uninstall must refuse to touch invalid JSON")
	}
	if got := mustRead(t, path); got != broken {
		t.Errorf("file must stay untouched:\ngot:  %q\nwant: %q", got, broken)
	}
}

func TestInitializeRegistersEventAdapterOnly(t *testing.T) {
	Initialize()
	if _, ok := caf.GetAdapter("claude-code"); !ok {
		t.Fatal("Initialize() must register the claude-code event adapter")
	}
	// Языковые адаптеры с выдуманными symbol graph удалены на шаге 7 —
	// в MVP граф строит gopls через universal LSP fallback (ADR-013).
	for _, id := range []string{"go", "python"} {
		if _, ok := caf.GetAdapter(id); ok {
			t.Errorf("language adapter %q must not be registered (removed in step 7)", id)
		}
	}
}

func TestOwnedEntry(t *testing.T) {
	tests := []struct {
		name  string
		entry string
		want  bool
	}{
		{"caf entry", `{"matcher":"Bash","hooks":[{"type":"command","command":"caf hook pre-commit --from-agent"}]}`, true},
		{"foreign entry", `{"hooks":[{"type":"command","command":"echo hi"}]}`, false},
		{"mixed entry is ours", `{"hooks":[{"type":"command","command":"echo hi"},{"type":"command","command":"caf hook post-edit"}]}`, true},
		{"garbage entry", `not json`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ownedEntry(json.RawMessage(tt.entry)); got != tt.want {
				t.Errorf("ownedEntry(%s) = %v, want %v", tt.entry, got, tt.want)
			}
		})
	}
}
