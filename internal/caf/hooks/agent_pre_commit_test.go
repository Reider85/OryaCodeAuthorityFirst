package hooks

import (
	"strings"
	"testing"
)

// agentPayload собирает stdin-JSON в формате PreToolUse Claude Code.
func agentPayload(tool, command string) string {
	return `{"tool_name":"` + tool + `","tool_input":{"command":` + jsonString(command) + `}}`
}

// jsonString экранирует строку как JSON-литерал (минимально: кавычки и \n).
func jsonString(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	return `"` + r.Replace(s) + `"`
}

// TestPreCommitFromAgentRunsForGitCommit — ядро контракта шага 7: хук
// выполняется только когда агент реально запускает git commit, и сообщение
// коммита извлекается из команды (-m), чтобы advisory-проверка трейлеров
// CAF-Session сработала уже на стороне агента.
func TestPreCommitFromAgentRunsForGitCommit(t *testing.T) {
	repo := t.TempDir()
	seedSession(t, repo, "2026-10-24-a1b2c3")

	payload := agentPayload("Bash", "git commit -m \"feat: x\nCAF-Session: 2026-10-24-a1b2c3\"")
	out := captureLog(t, func() {
		if err := PreCommitFromAgent(HookContext{RepoRoot: repo}, strings.NewReader(payload)); err != nil {
			t.Errorf("PreCommitFromAgent must stay fail-safe (nil), got %v", err)
		}
	})
	if !strings.Contains(out, "verified") {
		t.Errorf("expected CAF-Session verification in log, got:\n%s", out)
	}
}

// TestPreCommitFromAgentAdvisoryWithoutMessage — git commit без -m: advisory-
// режим («message unavailable»), nil, ничего не блокируется.
func TestPreCommitFromAgentAdvisoryWithoutMessage(t *testing.T) {
	payload := agentPayload("Bash", "git commit --amend --no-edit")
	out := captureLog(t, func() {
		if err := PreCommitFromAgent(HookContext{RepoRoot: t.TempDir()}, strings.NewReader(payload)); err != nil {
			t.Errorf("must stay fail-safe, got %v", err)
		}
	})
	if !strings.Contains(out, "unavailable") {
		t.Errorf("expected advisory 'message unavailable' path, got:\n%s", out)
	}
}

func TestPreCommitFromAgentSkipsNonGitCommands(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantLog string
	}{
		{
			"other bash command",
			agentPayload("Bash", "ls -la"),
			"is not a git commit",
		},
		{
			"git status not commit",
			agentPayload("Bash", "git status"),
			"is not a git commit",
		},
		{
			"non-bash tool",
			`{"tool_name":"Read","tool_input":{"path":"main.go"}}`,
			"is not Bash",
		},
		{
			"unparseable json",
			`{ this is not json`,
			"unparseable hook input",
		},
		{
			"empty input",
			`   `,
			"empty hook input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureLog(t, func() {
				if err := PreCommitFromAgent(
					HookContext{RepoRoot: t.TempDir()},
					strings.NewReader(tt.payload),
				); err != nil {
					t.Errorf("must stay fail-safe (nil), got %v", err)
				}
			})
			if !strings.Contains(out, tt.wantLog) {
				t.Errorf("expected log containing %q, got:\n%s", tt.wantLog, out)
			}
			if strings.Contains(out, "verified") {
				t.Errorf("hook must not execute for a skipped input:\n%s", out)
			}
		})
	}
}

func TestGitCommitMessage(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want string
	}{
		{"short flag plain", "git commit -m hello", "hello"},
		{"short flag double quotes", `git commit -m "multi word"`, "multi word"},
		{"short flag single quotes", "git commit -m 'single word'", "single word"},
		{"long flag with equals", `git commit --message="via equals"`, "via equals"},
		{"attached short flag", "git commit -mtight", "tight"},
		{"no message flag", "git commit --amend --no-edit", ""},
		{"bare commit", "git commit", ""},
		{"message after other flags", "git commit -a -m \"after flags\"", "after flags"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gitCommitMessage(tt.cmd); got != tt.want {
				t.Errorf("gitCommitMessage(%q) = %q, want %q", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestSplitShellArgs(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"plain", "git commit -m hi", []string{"git", "commit", "-m", "hi"}},
		{"quoted with spaces", `a "b c" d`, []string{"a", "b c", "d"}},
		{"single quoted", "a 'b c'", []string{"a", "b c"}},
		{"empty quoted arg preserved", `a "" b`, []string{"a", "", "b"}},
		{"extra whitespace collapsed", "  a   b  ", []string{"a", "b"}},
		{"newline inside quotes", "a \"b\nc\"", []string{"a", "b\nc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitShellArgs(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("splitShellArgs(%q) = %q, want %q", tt.in, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitShellArgs(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
				}
			}
		})
	}
}
