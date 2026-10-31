// Package opencode — Tier-1 нативный адаптер opencode (FR-ADP-001, ADR-009).
//
// Контракт ARCHITECT.md §3.1: Install() идемпотентен, Detect() отвечает
// «есть ли агент на машине», Uninstall() удаляет только своё. Языковых
// (LanguageAdapter) реализаций здесь больше нет: их выдуманные symbol graph
// не входили в скоуп MVP §2 — граф даёт gopls через universal LSP fallback
// (step 8, ADR-013).
package opencode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

const (
	// binaryName — имя бинарника opencode в $PATH (детекция, FR-ADP-005).
	binaryName = "opencode"

	// cafBlockStart/cafBlockEnd помечают managed-блок в hooks.toml.
	// По маркерам я отличаю свой блок от чужого — как cafHookMarker в
	// hooks/fail_safe.go: чужие хуки не затираются никогда.
	cafBlockStart = "# >>> caf managed block (installed by caf init) >>>"
	cafBlockEnd   = "# <<< caf managed block <<<"
)

// managedHooks — канонический тело managed-блока. Формат — дословно из
// OryaCodeAuthorityFirst.md §7 ([[hooks]] event + command), но команды
// указывают на caf hook, а не на orya (проект переименован, ADR-008: имя
// бинарника `caf`).
//
// Я сознательно НЕ подставляю {{session_id}}/{{files_changed}}-шаблоны из
// примера §7: opencode не документирует такой синтаксис шаблонизации в наших
// требованиях, а выдуманный шаблон молча сломал бы хук. MVP-вариант — команды
// без флагов: caf hook определяет репозиторий по cwd процесса агента.
const managedHooks = `[[hooks]]
event = "pre-session"
command = "caf hook pre-session"

[[hooks]]
event = "post-edit"
command = "caf hook post-edit"

[[hooks]]
event = "pre-commit"
command = "caf hook pre-commit"
`

// OpencodeAdapter implements the Adapter interface for opencode.
type OpencodeAdapter struct{}

// ID returns the unique identifier for this adapter.
func (a *OpencodeAdapter) ID() string {
	return "opencode"
}

// Capabilities returns the capability matrix: opencode поддерживает все три
// канонических события нативно (ARCHITECT.md §1, FR-ADP-002).
func (a *OpencodeAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession: caf.TierNative,
		PostEdit:   caf.TierNative,
		PreCommit:  caf.TierNative,
		MCP:        true,
		Tier:       1,
	}
}

// ConfigPath — путь к hooks.toml opencode.
//
// На non-Windows уважаю XDG_CONFIG_HOME (opencode следует XDG-спеке),
// иначе ~/.config/opencode/hooks.toml (OryaCodeAuthorityFirst.md §7).
// Домашний каталог беру через os.UserHomeDir(): на Windows это USERPROFILE,
// а Getenv("HOME") там пуст — старая версия писала в относительный
// ".config/opencode" от текущей директории (баг шага 7).
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("cannot resolve user home directory: %w", err)
	}
	base := filepath.Join(home, ".config")
	if runtime.GOOS != "windows" {
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			base = xdg
		}
	}
	return filepath.Join(base, "opencode", "hooks.toml"), nil
}

// Detect returns true, если opencode установлен на машине (бинарник в $PATH).
//
// Я исправил семантику на шаге 7: раньше Detect() требовал ещё и наличие
// конфиг-директории, из-за чего свежая машина с только что установленным
// opencode не детектилась — `caf init --auto-detect` и exit-criterion §12.5
// проваливались. Детекция отвечает на вопрос «есть ли агент», а не
// «настроен ли caf» — для второго есть Configured() (ARCHITECT.md §3.1).
func (a *OpencodeAdapter) Detect() bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

// Configured возвращает true, если managed-блок caf прописан в hooks.toml
// (статус hooks wired в `caf adapters list`).
func (a *OpencodeAdapter) Configured() bool {
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), cafBlockStart)
}

// Install прописывает caf-хуки в hooks.toml opencode.
//
// Идемпотентность (ARCHITECT.md §3.1): повторный вызов переписывает свой
// блок каноническим содержимым и не трогает ничего больше. Чужой hooks.toml
// (другие хуки пользователя) не затирается — я ДОПИСЫВАЮ managed-блок после
// его содержимого; отказ от append здесь был бы избыточен, в отличие от
// .git/hooks/pre-commit, где перезапись уничтожила бы чужой хук целиком
// (TOML-таблицы безопасно дописываются в конец файла).
func (a *OpencodeAdapter) Install() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return writeManagedBlock(path)
}

// Uninstall удаляет только managed-блок caf. Если файл после этого пуст —
// удаляет файл; если остались чужие хуки — сохраняет их. Файл без нашего
// маркера не трогаю (отказ, как RemoveUniversalHook).
func (a *OpencodeAdapter) Uninstall() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // уже не установлен — не ошибка (идемпотентность)
		}
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	body := string(existing)
	prefix, suffix, found, damaged := splitBlock(body)
	switch {
	case !found:
		if strings.TrimSpace(body) == "" {
			return os.Remove(path) // пустой файл — нечего сохранять
		}
		return fmt.Errorf("refusing to remove non-caf hooks config at %s "+
			"— caf did not write this file", path)
	case damaged:
		// Начало блока есть, но содержимое после него не наше: безопасно
		// вырезать нечего. Молча портить чужой файл я не стану — отказ.
		return fmt.Errorf("caf block at %s is damaged (missing end marker and "+
			"content is not caf-managed); repair it manually, then rerun uninstall", path)
	}

	remainder := prefix + suffix
	if strings.TrimSpace(remainder) == "" {
		return os.Remove(path)
	}
	if !strings.HasSuffix(remainder, "\n") {
		remainder += "\n"
	}
	if err := os.WriteFile(path, []byte(remainder), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

// splitBlock делит тело файла на (до managed-блока, после него).
//
// found=false — блока нет. damaged=true — встречается cafBlockStart, но
// содержимое после него не похоже на наш managed-блок и конец блока не
// найден: безопасно разделить невозможно, вызывающий принимает решение
// (Install — self-heal, Uninstall — отказ, файл не трогаем).
func splitBlock(body string) (prefix, suffix string, found, damaged bool) {
	start := strings.Index(body, cafBlockStart)
	if start < 0 {
		return body, "", false, false
	}
	prefix = body[:start]

	rest := body[start+len(cafBlockStart):]
	rest = strings.TrimPrefix(rest, "\r\n")
	rest = strings.TrimPrefix(rest, "\n")

	if i := strings.Index(rest, cafBlockEnd); i >= 0 {
		suffix = rest[i+len(cafBlockEnd):]
		// Один перевод строки после маркера принадлежит нашему блоку
		// (renderBlock добавляет его); лишние переводы — чужие, сохраняем.
		suffix = strings.TrimPrefix(suffix, "\r\n")
		suffix = strings.TrimPrefix(suffix, "\n")
		return prefix, suffix, true, false
	}
	// end-маркера нет: обрыв нашего блока либо чужой контент после маркера.
	// Перед сравнением нормализую ведущие пробелы/переводы строки — иначе
	// лишняя пустая строка внутри нашего блока выглядела бы «чужим» контентом
	// и Install продублировал бы [[hooks]]-таблицы.
	trimmed := strings.TrimLeft(rest, " \t\r\n")
	if strings.HasPrefix(trimmed, managedHooks) {
		return prefix, trimmed[len(managedHooks):], true, false
	}
	if trimmed == "" || strings.HasPrefix(managedHooks, trimmed) {
		return prefix, "", true, false // усечённый (оборванный) наш блок
	}
	return prefix, rest, true, true
}

// renderBlock возвращает канонический managed-блок с маркерами.
func renderBlock() string {
	return cafBlockStart + "\n" + managedHooks + cafBlockEnd + "\n"
}

// writeManagedBlock создаёт/обновляет файл с managed-блоком.
func writeManagedBlock(path string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var content string
	if err != nil {
		content = renderBlock()
	} else {
		content = mergeBlock(string(existing))
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create opencode config directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

// mergeBlock заменяет существующий managed-блок (self-heal по маркерам)
// или дописывает его в конец чужого файла. Результат детерминирован:
// mergeBlock(renderBlock()) == renderBlock() — в этом идемпотентность.
func mergeBlock(existing string) string {
	prefix, suffix, found, _ := splitBlock(existing)
	if found {
		// damaged=true тоже попадает сюда: хвост после маркера (чужой или
		// неузнаваемый) сохраняется ПОСЛЕ канонического блока — не теряем.
		return prefix + renderBlock() + suffix
	}
	// Блока нет — дописываем в конец чужого TOML (новые таблицы безопасно
	// добавляются в конец файла), с одной пустой строкой-разделителем.
	if existing != "" {
		existing = strings.TrimRight(existing, "\r\n") + "\n\n"
	}
	return existing + renderBlock()
}

// Initialize регистрирует Tier-1 opencode адаптер в реестре caf (FR-ADP-001).
// Языковые адаптеры здесь не регистрируются — см. комментарий пакета.
func Initialize() {
	caf.RegisterAdapter(&OpencodeAdapter{})
}
