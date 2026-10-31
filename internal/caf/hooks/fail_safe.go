package hooks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// HookContext содержит контекст выполнения хука.
type HookContext struct {
	RepoRoot  string
	SessionID string
	EventType string
}

// universalHookScript — дословный текст universal pre-commit hook из
// PROMPT-MVP.md §2.2 (FR-ADP-003). Строки LF-only: sh не переваривает
// CRLF в начале файла (Git for Windows запускает хук через sh).
const universalHookScript = `#!/bin/sh
# .git/hooks/pre-commit — installed by caf init
# Universal: fires for all AI agents writing to this repo.
# Fail-safe: if caf binary is missing, warn and continue (never block commit).
CAF_PATH=$(command -v caf 2>/dev/null)
if [ -z "$CAF_PATH" ]; then
  echo "⚠ caf binary not in $PATH; pre-commit hook skipped (sessions continue)."
  exit 0
fi
exec "$CAF_PATH" hook pre-commit --repo "$(pwd)" "$@"
`

// cafHookMarker помечает хук, установленный caf init. По маркеру мы
// отличаем свой хук от чужого и отказываемся перезаписывать чужой.
const cafHookMarker = "installed by caf init"

// hookPathFor — путь к universal pre-commit hook в репозитории.
func hookPathFor(repoRoot string) string {
	return filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
}

// ExecuteHook выполняет событие хука через бинарник caf в $PATH —
// INV-ARCH-004 / FR-REL-001: warn-and-continue, никогда не блокирует
// сессию или коммит.
//
// Я исправил два бага на шаге 6: (1) аргументы передавались одним
// элементом "hook <event>" — кobra получала неизвестное событие и
// падала; (2) сбой выполнения возвращал непустую ошибку наружу, чем
// нарушал fail-safe-контракт §7 — теперь это WARN и nil. Сигнатура
// error сохранена под ALPHA-валидацию (FR-DEC-001).
func ExecuteHook(event string, ctx HookContext) error {
	cafPath, err := exec.LookPath("caf")
	if err != nil {
		// caf binary missing — warn and continue, never block.
		log.Printf("⚠ caf binary not in $PATH; hook %s skipped", event)
		return nil
	}

	args := []string{"hook", event}
	if ctx.RepoRoot != "" {
		args = append(args, "--repo", ctx.RepoRoot)
	}
	if ctx.SessionID != "" {
		args = append(args, "--session", ctx.SessionID)
	}

	output, err := exec.Command(cafPath, args...).CombinedOutput()
	if err != nil {
		// caf упал — warn-and-continue: коммит/сессия идут дальше.
		log.Printf("⚠ hook %s failed (non-blocking, INV-ARCH-004): %v; output: %s",
			event, err, strings.TrimSpace(string(output)))
		return nil
	}

	log.Printf("✓ Hook %s executed successfully: %s", event, strings.TrimSpace(string(output)))
	return nil
}

// IsHookInstalled проверяет, установлен ли universal pre-commit hook.
// На unix дополнительно требуется exec-бит; на Windows Git for Windows
// сам запускает .git/hooks/* через sh, поэтому важен только факт файла
// (иначе caf doctor давал бы false-negative на Windows).
func IsHookInstalled(repoRoot string) (bool, error) {
	info, err := os.Stat(hookPathFor(repoRoot))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0111 == 0 {
		return false, nil
	}
	return true, nil
}

// InstallUniversalHook ставит universal pre-commit hook (PROMPT-MVP §2.2).
//
// Я отказываюсь перезаписывать чужой хук (husky, golangci-lint и т.п.):
// молчаливая перезапись чужой инфраструктуры означала бы потерянные
// проверки пользователя. Повторный caf init поверх caf-хука —
// идемпотентная перезапись своего. Отказ — не фатален: caf init
// печатает WARN и продолжает (fail-safe дух §7).
func InstallUniversalHook(repoRoot string) error {
	path := hookPathFor(repoRoot)

	if existing, err := os.ReadFile(path); err == nil &&
		!strings.Contains(string(existing), cafHookMarker) {
		return fmt.Errorf("refusing to overwrite existing non-caf pre-commit hook at %s "+
			"— back it up, merge caf's hook manually, then rerun `caf init`", path)
	}

	// Без .git не создаём ложный каталог .git — иначе git перестал бы
	// считать репозиторий валидным.
	if _, err := os.Stat(filepath.Join(repoRoot, ".git")); err != nil {
		return fmt.Errorf("not a git repository (no .git under %s): "+
			"universal pre-commit hook not installed", repoRoot)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(universalHookScript), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}
	return nil
}

// RemoveUniversalHook удаляет только наш хук (идемпотентно); чужой —
// не трогаем, как и при установке.
func RemoveUniversalHook(repoRoot string) error {
	path := hookPathFor(repoRoot)
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // уже не установлен — не ошибка
		}
		return err
	}
	if !strings.Contains(string(existing), cafHookMarker) {
		return fmt.Errorf("refusing to remove non-caf pre-commit hook at %s", path)
	}
	return os.Remove(path)
}
