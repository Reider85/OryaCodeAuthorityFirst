package hooks

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeCafStub кладёт в dir исполняемую заглушку `caf` с телом body.
// На Windows это .bat (Go запускает их через exec.Command), на unix —
// shell-скрипт с exec-битом.
func writeCafStub(t *testing.T, dir, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		path := filepath.Join(dir, "caf.bat")
		if err := os.WriteFile(path, []byte("@echo off\r\n"+body), 0644); err != nil {
			t.Fatalf("failed to write caf.bat stub: %v", err)
		}
		return path
	}
	path := filepath.Join(dir, "caf")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0755); err != nil {
		t.Fatalf("failed to write caf stub: %v", err)
	}
	return path
}

func TestHookPathFor(t *testing.T) {
	got := hookPathFor(`C:\repo`)
	want := filepath.Join(`C:\repo`, ".git", "hooks", "pre-commit")
	if got != want {
		t.Errorf("hookPathFor() = %q, want %q", got, want)
	}
}

// TestExecuteHookWithoutCafBinary — INV-ARCH-004: caf отсутствует в
// $PATH → warn-and-continue, никогда не ошибка.
func TestExecuteHookWithoutCafBinary(t *testing.T) {
	t.Setenv("PATH", "")
	out := captureLog(t, func() {
		if err := ExecuteHook("pre-commit", HookContext{RepoRoot: t.TempDir()}); err != nil {
			t.Errorf("ExecuteHook must return nil when caf is missing, got %v", err)
		}
	})
	if !strings.Contains(out, "not in $PATH") {
		t.Errorf("expected skip warning, got: %s", out)
	}
}

// TestExecuteHookSuccess — ветка успешного выполнения: argv собирается
// как ["hook", <event>, "--repo", ...] (баг шага 6 исправлен) и caf
// запускается без ошибки.
func TestExecuteHookSuccess(t *testing.T) {
	bin := t.TempDir()
	writeCafStub(t, bin, "echo caf hook ok")
	t.Setenv("PATH", bin)

	out := captureLog(t, func() {
		if err := ExecuteHook("pre-commit", HookContext{RepoRoot: t.TempDir(), SessionID: "2026-10-24-a1b2c3"}); err != nil {
			t.Errorf("ExecuteHook must return nil on success, got %v", err)
		}
	})
	if !strings.Contains(out, "executed successfully") {
		t.Errorf("expected success log, got: %s", out)
	}
}

// TestExecuteHookFailureIsNonBlocking — caf упал (ненулевой exit):
// warn-and-continue, nil наружу (fail-safe §7).
func TestExecuteHookFailureIsNonBlocking(t *testing.T) {
	bin := t.TempDir()
	writeCafStub(t, bin, "exit 1")
	t.Setenv("PATH", bin)

	out := captureLog(t, func() {
		if err := ExecuteHook("pre-commit", HookContext{RepoRoot: t.TempDir()}); err != nil {
			t.Errorf("ExecuteHook must return nil even when caf fails, got %v", err)
		}
	})
	if !strings.Contains(out, "non-blocking") {
		t.Errorf("expected non-blocking failure log, got: %s", out)
	}
}

// newGitRepo — temp-репозиторий с каталогом .git (без реального git init:
// для хука достаточно наличия .git).
func newGitRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git: %v", err)
	}
	return repo
}

func TestInstallUniversalHookRequiresGitRepo(t *testing.T) {
	err := InstallUniversalHook(t.TempDir())
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInstallUniversalHookCreatesHooksDir(t *testing.T) {
	repo := newGitRepo(t) // без каталога .git/hooks

	if err := InstallUniversalHook(repo); err != nil {
		t.Fatalf("InstallUniversalHook failed: %v", err)
	}

	path := hookPathFor(repo)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("hook not written: %v", err)
	}
	content := string(data)
	if content != universalHookScript {
		t.Errorf("hook content must match §2.2 script exactly")
	}
	if !strings.Contains(content, cafHookMarker) {
		t.Errorf("hook must contain caf marker")
	}
	if strings.Contains(content, "\r\n") {
		t.Errorf("hook must be LF-only (sh chokes on CRLF)")
	}
	if !strings.HasPrefix(content, "#!/bin/sh") {
		t.Errorf("hook must start with #!/bin/sh, got: %.20q", content)
	}
}

func TestInstallUniversalHookRefusesForeignHook(t *testing.T) {
	repo := newGitRepo(t)
	foreign := "#!/bin/sh\n# husky pre-commit\nexit 1\n"
	path := hookPathFor(repo)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(foreign), 0755); err != nil {
		t.Fatal(err)
	}

	err := InstallUniversalHook(repo)
	if err == nil {
		t.Fatal("expected refusal to overwrite foreign hook, got nil")
	}
	if !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Errorf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("foreign hook disappeared: %v", readErr)
	}
	if string(data) != foreign {
		t.Errorf("foreign hook must remain untouched, got: %q", data)
	}
}

func TestInstallUniversalHookIsIdempotent(t *testing.T) {
	repo := newGitRepo(t)
	if err := InstallUniversalHook(repo); err != nil {
		t.Fatalf("first install failed: %v", err)
	}
	if err := InstallUniversalHook(repo); err != nil {
		t.Fatalf("re-install over own hook must succeed, got: %v", err)
	}
	data, err := os.ReadFile(hookPathFor(repo))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != universalHookScript {
		t.Errorf("re-install must restore the canonical script")
	}
}

func TestIsHookInstalled(t *testing.T) {
	repo := newGitRepo(t)

	installed, err := IsHookInstalled(repo)
	if err != nil {
		t.Fatalf("IsHookInstalled errored: %v", err)
	}
	if installed {
		t.Error("hook must be reported absent before install")
	}

	if err := InstallUniversalHook(repo); err != nil {
		t.Fatal(err)
	}
	installed, err = IsHookInstalled(repo)
	if err != nil {
		t.Fatalf("IsHookInstalled errored: %v", err)
	}
	if !installed {
		t.Error("hook must be reported installed after install")
	}

	// exec-бит проверяется только не на Windows (см. IsHookInstalled).
	if runtime.GOOS != "windows" {
		path := hookPathFor(repo)
		if err := os.Chmod(path, 0644); err != nil {
			t.Fatal(err)
		}
		installed, err = IsHookInstalled(repo)
		if err != nil {
			t.Fatal(err)
		}
		if installed {
			t.Error("hook without exec bit must not count as installed")
		}
	}
}

func TestRemoveUniversalHook(t *testing.T) {
	repo := newGitRepo(t)

	// Идемпотентно: отсутствие — не ошибка.
	if err := RemoveUniversalHook(repo); err != nil {
		t.Errorf("removing absent hook must be a no-op, got %v", err)
	}

	if err := InstallUniversalHook(repo); err != nil {
		t.Fatal(err)
	}
	if err := RemoveUniversalHook(repo); err != nil {
		t.Fatalf("RemoveUniversalHook failed: %v", err)
	}
	if installed, _ := IsHookInstalled(repo); installed {
		t.Error("hook must be gone after remove")
	}
}

func TestRemoveUniversalHookRefusesForeignHook(t *testing.T) {
	repo := newGitRepo(t)
	path := hookPathFor(repo)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\n# someone else's hook\n"), 0755); err != nil {
		t.Fatal(err)
	}

	err := RemoveUniversalHook(repo)
	if err == nil {
		t.Fatal("expected refusal to remove foreign hook, got nil")
	}
	if !strings.Contains(err.Error(), "refusing to remove") {
		t.Errorf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(path); statErr != nil {
		t.Errorf("foreign hook must survive: %v", statErr)
	}
}
