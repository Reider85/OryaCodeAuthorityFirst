package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
)

// maxAgentInputBytes — верхняя граница stdin-ввода от агента. Claude Code
// шлёт JSON хука целиком; мегабайта хватает с большим запасом, а ограничение
// защищает от зависания на неожиданном потоке (fail-safe дух §7).
const maxAgentInputBytes = 1 << 20

// agentToolInput — JSON, который Claude Code передаёт хукам PreToolUse на
// stdin (см. claudecode.cafEventEntries: команда `caf hook pre-commit
// --from-agent`). Схема — публичный формат хуков Claude Code: tool_name +
// tool_input с полями конкретного инструмента (у Bash — command).
type agentToolInput struct {
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
}

// PreCommitFromAgent — точка входа `caf hook pre-commit --from-agent`.
//
// Зачем это нужно. У Claude Code нет отдельного события pre-commit, но
// ARCHITECT.md §1 (замороженная матрица, FR-ADP-002) декларирует для него
// native. Решение шага 7: хук вешается на PreToolUse инструмента Bash, а
// здесь я читаю stdin-JSON и выполняю pre-commit ТОЛЬКО когда агент реально
// запускает `git commit`. Остальные Bash-команды — быстрый skip.
//
// Почему stdin читается только по флагу --from-agent: путь `git commit`
// через .git/hooks/pre-commit никогда не читает stdin (иначе git мог бы
// заблокироваться на терминале). Два срабатывания на один коммит (Claude
// PreToolUse + universal git hook) безвредны: оба advisory в MVP
// (INV-ARCH-004).
//
// Любая ошибка парсинга → warn + skip + nil: битый вход агента не должен
// ломать вызов инструмента (fail-safe INV-ARCH-004, хук никогда не блокирует).
func PreCommitFromAgent(ctx HookContext, in io.Reader) error {
	data, err := io.ReadAll(io.LimitReader(in, maxAgentInputBytes))
	if err != nil {
		log.Printf("⚠ pre-commit --from-agent: cannot read hook input (%v); "+
			"skipped, never blocks", err)
		return nil
	}
	if len(bytes.TrimSpace(data)) == 0 {
		log.Printf("ℹ pre-commit --from-agent: empty hook input; skipped")
		return nil
	}

	var payload agentToolInput
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("⚠ pre-commit --from-agent: unparseable hook input (%v); "+
			"skipped, never blocks", err)
		return nil
	}
	if payload.ToolName != "" && payload.ToolName != "Bash" {
		log.Printf("ℹ pre-commit --from-agent: tool %q is not Bash; skipped",
			payload.ToolName)
		return nil
	}

	cmd := strings.TrimSpace(payload.ToolInput.Command)
	if cmd != "git commit" && !strings.HasPrefix(cmd, "git commit ") {
		log.Printf("ℹ pre-commit --from-agent: %q is not a git commit; skipped",
			clip(cmd, 80))
		return nil
	}

	// git commit обнаружен. Сообщение коммита извлекаю из самой команды
	// (-m/--message): это даёт advisory-проверку трейлеров CAF-Session уже на
	// стороне агента. Разобрать не удалось — пустая строка, PreCommitHook
	// честно уйдёт в режим «message unavailable» и ничего не заблокирует.
	return PreCommitHook(ctx, gitCommitMessage(cmd))
}

// gitCommitMessage извлекает аргумент первого -m/--message из команды вида
// `git commit -m "..."`. Парсер намеренно простой (MVP): кавычки ' и ",
// без поддержки обратного экранирования внутри строк. Не вышло разобрать —
// возвращаем "", и pre-commit работает advisory.
func gitCommitMessage(command string) string {
	args := splitShellArgs(command)
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "-m" || args[i] == "--message":
			if i+1 < len(args) {
				return args[i+1]
			}
		case strings.HasPrefix(args[i], "--message="):
			return strings.TrimPrefix(args[i], "--message=")
		case strings.HasPrefix(args[i], "-m") && len(args[i]) > 2:
			return args[i][len("-m"):]
		}
	}
	return ""
}

// splitShellArgs — минимальный разделитель командной строки с учётом
// одинарных и двойных кавычек (аргумент в кавычках сохраняется целиком,
// включая пробелы; пустой аргумент "" не теряется).
func splitShellArgs(s string) []string {
	var (
		args    []string
		cur     strings.Builder
		quote   rune // 0 — вне кавычек
		started bool // есть токен (в т.ч. пустой в кавычках)
	)
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
			started = true
		case r == ' ' || r == '\t' || r == '\n':
			if started {
				args = append(args, cur.String())
				cur.Reset()
				started = false
			}
		default:
			cur.WriteRune(r)
			started = true
		}
	}
	if started {
		args = append(args, cur.String())
	}
	return args
}

// clip урезает строку для лога (полные команды агента бывают огромными).
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
