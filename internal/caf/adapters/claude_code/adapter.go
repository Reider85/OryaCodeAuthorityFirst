// Package claudecode — Tier-1 нативный адаптер Claude Code (FR-ADP-001, ADR-009).
//
// Контракт ARCHITECT.md §3.1: Install() идемпотентен, Detect() отвечает
// «есть ли агент на машине», Uninstall() удаляет только своё. Языковых
// (LanguageAdapter) реализаций здесь больше нет: их выдуманные symbol graph
// не входили в скоуп MVP §2 — граф даёт gopls через universal LSP fallback
// (step 8, ADR-013).
package claudecode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

const (
	// binaryName — имя бинарника Claude Code в $PATH (детекция, FR-ADP-005).
	binaryName = "claude"

	// settingsFile — штатный конфиг Claude Code: ~/.claude/settings.json.
	settingsFile = "settings.json"

	// cafCommandPrefix помечает наши hook-записи: любая команда, содержащая
	// "caf hook", считается принадлежащей caf. По этому признаку Uninstall
	// вырезает только своё и не трогает чужие хуки пользователя.
	cafCommandPrefix = "caf hook"
)

// hookCommand — запись {"type":"command","command":"..."} в схеме хуков
// Claude Code (https://docs.anthropic.com/en/docs/claude-code/hooks).
type hookCommand struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// hookEntry — один элемент массива-хука: matcher (опционально) + команды.
type hookEntry struct {
	Matcher string        `json:"matcher,omitempty"`
	Hooks   []hookCommand `json:"hooks"`
}

// eventEntry — (событие Claude Code, наша запись) в каноническом порядке.
type eventEntry struct {
	Event string
	Entry hookEntry
}

// cafEventEntries — хуки, которые caf прописывает в settings.json.
//
// Соответствие каноническим событиям (ARCHITECT.md §4.1):
//   - pre-session → SessionStart (читает .caf/context.md, INV-ARCH-004 —
//     сам хук fail-safe: caf отсутствует в $PATH → warn, агент работает);
//   - post-edit   → PostToolUse на Write/Edit/MultiEdit (черновик сессии);
//   - pre-commit  → PreToolUse на Bash, команда `caf hook pre-commit
//     --from-agent`: Claude Code не имеет отдельного pre-commit события,
//     поэтому я фильтрую stdin-JSON и срабатываю только на `git commit`
//     (матрица native из ARCHITECT.md §1 остаётся честной).
//
// Двойное срабатывание с universal .git/hooks/pre-commit безвредно: оба
// пути advisory в MVP (INV-ARCH-004).
func cafEventEntries() []eventEntry {
	return []eventEntry{
		{"SessionStart", hookEntry{
			Hooks: []hookCommand{{Type: "command", Command: "caf hook pre-session"}},
		}},
		{"PostToolUse", hookEntry{
			Matcher: "Write|Edit|MultiEdit",
			Hooks:   []hookCommand{{Type: "command", Command: "caf hook post-edit"}},
		}},
		{"PreToolUse", hookEntry{
			Matcher: "Bash",
			Hooks:   []hookCommand{{Type: "command", Command: "caf hook pre-commit --from-agent"}},
		}},
	}
}

// ClaudeCodeAdapter implements the Adapter interface for Claude Code.
type ClaudeCodeAdapter struct{}

// ID returns the unique identifier for this adapter.
func (a *ClaudeCodeAdapter) ID() string {
	return "claude-code"
}

// Capabilities returns the capability matrix: Claude Code поддерживает все
// три канонических события нативно (ARCHITECT.md §1, FR-ADP-002).
func (a *ClaudeCodeAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession: caf.TierNative,
		PostEdit:   caf.TierNative,
		PreCommit:  caf.TierNative,
		MCP:        true,
		Tier:       1,
	}
}

// ConfigPath — путь к settings.json Claude Code (~/.claude/settings.json).
//
// Домашний каталог беру через os.UserHomeDir(): на Windows это USERPROFILE,
// а Getenv("HOME") там пуст — старая версия писала в относительный
// ".claude" от текущей директории (баг шага 7).
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("cannot resolve user home directory: %w", err)
	}
	return filepath.Join(home, ".claude", settingsFile), nil
}

// Detect returns true, если Claude Code установлен на машине (бинарник в $PATH).
//
// Семантика исправлена на шаге 7: раньше Detect() требовал ещё и наличие
// конфиг-директории, из-за чего свежая машина с установленным `claude` не
// детектилась (exit-criterion §12.5). Детекция отвечает «есть ли агент»;
// «настроен ли caf» — Configured() (ARCHITECT.md §3.1).
func (a *ClaudeCodeAdapter) Detect() bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

// Configured возвращает true, если caf-хуки прописаны в settings.json
// (статус hooks wired в `caf adapters list`).
func (a *ClaudeCodeAdapter) Configured() bool {
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), cafCommandPrefix)
}

// Install прописывает caf-хуки в settings.json Claude Code.
//
// Идемпотентность (ARCHITECT.md §3.1): существующие caf-записи вырезаются
// и добавляются заново в каноническом виде — повторный вызов даёт байт-в-байт
// тот же файл. Чужие ключи и чужие хуки сохраняются: верхний уровень парсится
// как map[string]json.RawMessage, поэтому не тронутые ключи уходят в файл
// без перекодирования (комментарии в JSON отсутствуют, а вот HTML-экранирование
// я отключаю — чужие URL не должны превращаться в \u0026).
//
// Битый JSON (например, jsonc с комментариями) я НЕ чиню и НЕ переписываю:
// возвращаю ошибку, CLI печатает WARN и продолжает (fail-safe INV-ARCH-004).
func (a *ClaudeCodeAdapter) Install() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	root, err := readSettings(path)
	if err != nil {
		return err
	}
	if root == nil {
		root = map[string]json.RawMessage{}
	}

	hooks, err := hooksMap(root, path)
	if err != nil {
		return err
	}
	// Свои записи вырезаем (защита от дублей при повторной установке),
	// чужие — оставляем на месте.
	for event, entries := range hooks {
		kept := make([]json.RawMessage, 0, len(entries))
		for _, e := range entries {
			if !ownedEntry(e) {
				kept = append(kept, e)
			}
		}
		if len(kept) > 0 {
			hooks[event] = kept
		} else {
			delete(hooks, event)
		}
	}
	for _, ee := range cafEventEntries() {
		raw, err := json.Marshal(ee.Entry)
		if err != nil {
			return fmt.Errorf("failed to marshal caf hook entry: %w", err)
		}
		hooks[ee.Event] = append(hooks[ee.Event], raw)
	}

	if len(hooks) == 0 {
		delete(root, "hooks")
	} else {
		raw, err := json.Marshal(hooks)
		if err != nil {
			return fmt.Errorf("failed to marshal hooks: %w", err)
		}
		root["hooks"] = raw
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", filepath.Dir(path), err)
	}
	return writeSettings(path, root)
}

// Uninstall удаляет только caf-записи из settings.json. Если после этого не
// осталось ни одного ключа — удаляет файл. Файл без наших записей не трогаю
// (отказ, как RemoveUniversalHook); битый JSON не трогаю (ошибка → WARN).
func (a *ClaudeCodeAdapter) Uninstall() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	root, err := readSettings(path)
	if err != nil {
		return err
	}
	if root == nil {
		return nil // файла нет — уже не установлен (идемпотентность)
	}

	hooks, err := hooksMap(root, path)
	if err != nil {
		return err
	}
	owned := false
	for event, entries := range hooks {
		kept := make([]json.RawMessage, 0, len(entries))
		for _, e := range entries {
			if ownedEntry(e) {
				owned = true
				continue
			}
			kept = append(kept, e)
		}
		if len(kept) > 0 {
			hooks[event] = kept
		} else {
			delete(hooks, event)
		}
	}
	if !owned {
		if len(hooks) > 0 {
			// Есть хуки, но ни одного caf-записи: файл писал не caf — отказ.
			return fmt.Errorf("refusing to remove non-caf hooks config at %s "+
				"— caf did not write these hooks", path)
		}
		// Caf-хуков не было вовсе — файл не переписываем даже для нормализации.
		return nil
	}

	if len(hooks) == 0 {
		delete(root, "hooks")
	} else {
		raw, err := json.Marshal(hooks)
		if err != nil {
			return fmt.Errorf("failed to marshal hooks: %w", err)
		}
		root["hooks"] = raw
	}
	if len(root) == 0 {
		// Файл существовал только ради наших хуков — удаляем целиком.
		return os.Remove(path)
	}
	return writeSettings(path, root)
}

// readSettings читает settings.json. Файла нет → (nil, nil): вызывающий
// решает сам (Install создаёт, Uninstall — не ошибка). Битый JSON → ошибка
// с понятным сообщением; caf файл не трогает.
func readSettings(path string) (map[string]json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]json.RawMessage{}, nil
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("%s is not valid JSON (%w); caf will not touch it — "+
			"fix the file manually, then rerun `caf init --adapters claude-code`", path, err)
	}
	return root, nil
}

// hooksMap извлекает объект hooks из root (нет ключа → пустая карта).
func hooksMap(root map[string]json.RawMessage, path string) (map[string][]json.RawMessage, error) {
	raw, ok := root["hooks"]
	if !ok || len(bytes.TrimSpace(raw)) == 0 {
		return map[string][]json.RawMessage{}, nil
	}
	var hooks map[string][]json.RawMessage
	if err := json.Unmarshal(raw, &hooks); err != nil {
		return nil, fmt.Errorf("%s: \"hooks\" is not a valid object (%w); "+
			"caf will not touch it", path, err)
	}
	if hooks == nil {
		return map[string][]json.RawMessage{}, nil
	}
	return hooks, nil
}

// ownedEntry определяет, записана ли hook-запись caf (команда содержит
// "caf hook"). Не разбираемая запись не наша — её не трогаем.
func ownedEntry(entry json.RawMessage) bool {
	var probe struct {
		Hooks []hookCommand `json:"hooks"`
	}
	if err := json.Unmarshal(entry, &probe); err != nil {
		return false
	}
	for _, h := range probe.Hooks {
		if strings.Contains(h.Command, cafCommandPrefix) {
			return true
		}
	}
	return false
}

// writeSettings сериализует root с отступом в 2 пробела и пишет файл.
func writeSettings(path string, root map[string]json.RawMessage) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false) // чужие URL и пути не должны экранироваться
	enc.SetIndent("", "  ")
	if err := enc.Encode(root); err != nil {
		return fmt.Errorf("failed to encode %s: %w", path, err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

// Initialize регистрирует Tier-1 Claude Code адаптер в реестре caf (FR-ADP-001).
// Языковые адаптеры здесь не регистрируются — см. комментарий пакета.
func Initialize() {
	caf.RegisterAdapter(&ClaudeCodeAdapter{})
}
