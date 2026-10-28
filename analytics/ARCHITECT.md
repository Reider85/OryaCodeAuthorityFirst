# CodeAuthorityFirst — ARCHITECT.md

> **Версия спецификации / Spec version:** 0.2 (universal)
> **Форм-фактор / Form factor:** CLI + MCP server + Agent Adapter Layer
> **Охват / Scope:** any codebase authored or co-authored by AI coding agents; agent-agnostic via pluggable adapters
> **Поддерживаемые агенты / Supported agents (v0.2 target):** opencode, Claude Code, Cursor, Cline, Aider, KiloCode, Codex CLI, Continue.dev, OryaAgent
> **Голос реконструкции / Reconstruction voice:** first-person («Я выбрал X потому что…»)
> **Предыдущая версия / Previous version:** 0.1 (Orya-specific, see git history)
> **Язык / Language:** bilingual — Russian body, English glossary and per-section summaries
> **Аудитория / Audience:** tool maintainers + committers using AI coding agents in any stack

---

## 0. Document Conventions

### 0.1 Bilingual policy

**Русский** — основной язык документа. Технические термины сохраняются на языке оригинала (CLI-команды, YAML-ключи, ASCII-диаграммы, идентификаторы агентов). **English** — для глоссария, заголовков таблиц и кратких `EN summary.` блоков в начале каждой крупной секции. Это даёт международным контрибьюторам точку входа без перевода, а русскоязычным мейнтейнерам — точные формулировки без лоскутного перевода.

### 0.2 Glossary

| Term (EN) | Термин (RU) | Definition |
|---|---|---|
| Authorship | Авторство | субъективная способность объяснить код так, будто он написан тобой |
| Provenance | Происхождение | трассировка от промпта AI-сессии до git-объекта |
| Decision Record | Запись решения | ADR-подобный документ с акцентом на отвергнутые альтернативы |
| Mental Model | Ментальная карта | внутреннее представление связей и причинно-следственных решений |
| Spaced Repetition | Интервальное повторение | алгоритм SM-2 (как в Anki) |
| Blindspot | Слепая зона | файл, изменённый за 30 дней, но не покрытый SRS-карточкой |
| Agent Adapter | Адаптер агента | модуль, транслирующий CAF-события в hook-систему конкретного AI-агента |
| Capability Matrix | Матрица возможностей | таблица "агент × поддерживаемое hook-событие" |
| Trivial commit | Тривиальный коммит | коммит с пометкой `CAF-Trivial: true`, не требующий decision-записи |
| Card | Карточка | SRS-единица в Mental Model Reconstructor |

### 0.3 Identifier conventions

- CLI-команды: `caf <subcommand>` (snake_case для флагов)
- YAML-поля: `snake_case`
- Файлы provenance: `.caf/sessions/<YYYY-MM-DD>-<short_id>.yml`, где `short_id` — первые 6 hex-символов SHA-256 от `started_at + initial_prompt + agent_id`
- Файлы решений: `.caf/decisions/d_<YYYY-MM-DD>-<seq>.md`
- SRS-карточки: единый файл `.caf/mental-model/cards.yml`
- Git trailer: `CAF-Session: <id>` (заменяет Orya-Session из v0.1)

### 0.4 Naming rationale

**Project name:** `CodeAuthorityFirst` (вместо `OryaSkill`/`OryaCodeAuthorityFirst` из v0.1) — отражает универсальность: инструмент не привязан к конкретному стеку.

**CLI binary:** `caf` (вместо `orya`) — короткое, легко печатается, не конфликтует с существующими инструментами в `$PATH`.

**Folder:** `.caf/` (вместо `.orya/`) — то же имя, что и CLI, не содержит информации о каком-либо фреймворке.

Подробное обоснование — см. ADR-007, ADR-008, ADR-009.

---

## 1. Problem Statement

**EN summary.** When code in any repository is generated or co-authored by an AI coding agent (opencode, Claude Code, Cursor, Aider, etc.), the human author loses the mental map of decisions and rejected alternatives. Existing tools capture git history but not the *agent session context*: the prompt, the alternatives the agent considered and rolled back, the reasoning behind choosing approach A over B. `CodeAuthorityFirst` restores authorship by recording provenance across **any** agent, archiving decisions, and applying spaced repetition to the author's own knowledge of the codebase.

Когда код в любом репозитории пишется или редактируется AI-агентом (opencode, Claude Code, Cursor, Aider и др.), возникает специфическое отчуждение. Файл существует, компилируется, проходит тесты — но в голове автора нет истории его появления. Не зафиксированы рассмотренные альтернативы, не сохранены тупиковые ветки, которые агент попробовал и откатил, не сохранено обоснование выбора подхода. Git log показывает «refactor auth middleware» — но это память агента, а не память человека.

Это критично, потому что владение кодом — это не знание синтаксиса, а **ментальная карта принятых решений и отвергнутых альтернатив**. Без неё невозможно безопасно эволюционировать систему: каждое следующее изменение принимается в вакууме, риск регрессий растёт, а обсуждение архитектуры превращается в «спрошу у агента».

`CodeAuthorityFirst` — инструмент, который восстанавливает ощущение авторства **для любого проекта, разрабатываемого любым AI-агентом**. Он не пишет код за автора и не заменяет агент. Он делает так, чтобы через неделю, месяц или год после AI-сессии автор мог открыть любой файл в любом репозитории и сказать: «Я понимаю, почему здесь именно так, и могу объяснить это другому разработчику».

Универсальность достигается через **Agent Adapter Layer** (см. секцию 3.3) — каждый AI-агент получает свой адаптер, транслирующий внутренние события агента в единый CAF-контракт.

---

## 2. Design Principles

**EN summary.** Six principles govern all architectural decisions. The same six as v0.1, with one reformulation: "Cross-repo mental graph" is now "Cross-project mental graph" (any repo, not just Orya*). The agent-agnostic principle is implicit in the Adapter Layer — any feature that requires a specific agent's API and cannot be expressed via the adapter contract is out of scope for v0.2.

Инструмент строится на шести принципах, которые определяют все архитектурные решения ниже.

**Authorship-first.** Главная метрика — не покрытие кода тестами и не скорость разработки, а субъективная уверенность автора: «я могу объяснить этот код так, как будто написал его сам». Любая фича оценивается через этот критерий; фичи, не двигающие Authorship Confidence Score (см. секцию 12.1), не попадают в v0.2.

**Provenance by default.** Происхождение каждого нетривиального коммита (какой промпт, какая сессия, какой агент, какая итерация) записывается автоматически, без участия человека. Хранение provenance — это базовый слой, на котором строятся остальные модули; отсутствие provenance делает два других модуля бесполезными.

**Decision > Implementation.** Понимание «как» вторично по отношению к пониманию «почему». Инструмент тратит больше ресурсов на фиксацию отвергнутых альтернатив, чем на объяснение финальной реализации — финальную реализацию можно прочитать в коде, альтернативы же исчезают без следа.

**Spaced repetition of code.** Ментальная карта затухает без подкрепления — эмпирически этот порог составляет 4-6 недель. Применяется алгоритм SM-2 (как в Anki) к ключевым файлам и решениям, чтобы через 30/60/90 дней после коммита они автоматически возвращались в зону внимания автора.

**First-person voice.** Все генерируемые объяснения пишутся от лица автора: «Я выбрал RBAC вместо ABAC, потому что на тот момент у нас было всего три роли». Это когнитивный трюк — первое лицо заставляет ассимилировать решение как собственное, а не как внешне навязанное.

**Cross-project mental graph.** Любой набор репозиториев, разработанных автором, — один продукт с точки зрения ментальной карты. Инструмент хранит граф зависимостей между ними: «Repo A зависит от Repo B через интерфейс `Foo`, определённый в `RepoB/pkg/foo.go:42`». При изменении интерфейса в одном репо инструмент подсказывает, какие решения в зависимых репо могут стать неактуальными.

---

## 3. System Architecture

**EN summary.** Four layers compose the universal system: Interface Layer (CLI + MCP server), Agent Adapter Layer (pluggable per-agent modules translating CAF events to agent-specific hook systems), Storage Layer (`.caf/` folder committed to each repo's git history), Hooks Layer (canonical pre-session / post-edit / pre-commit events that adapters implement). The Adapter Layer is the v0.2 addition that makes the system universal — adapters declare a capability matrix and gracefully degrade when an agent lacks a feature.

`CodeAuthorityFirst` состоит из четырёх слоёв. Каждый слой можно использовать независимо, но вместе они дают синергетический эффект. **Новинка v0.2:** добавлен Agent Adapter Layer, изолирующий специфику каждого AI-агента за единым контрактом событий.

### 3.1 Layered diagram

```text
┌──────────────────────────────────────────────────────────────────────────┐
│  Interface Layer                                                          │
│  ┌──────────────────┐    ┌──────────────────┐                             │
│  │  caf CLI         │    │  MCP server      │  ← any MCP-compatible      │
│  │  (terminal/IDE)  │    │  (stdio / SSE)   │    client                  │
│  └────────┬─────────┘    └────────┬─────────┘                             │
│           │                       │                                       │
└───────────┼───────────────────────┼───────────────────────────────────────┘
            │                       │
┌───────────┼───────────────────────┼───────────────────────────────────────┐
│  Agent Adapter Layer               │                                       │
│  ┌────────┴───────────────────────┴────────┐                              │
│  │  Adapter Registry                        │                              │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐     │                              │
│  │  │opencode │ │ClaudeCC │ │ Cursor  │ ... │                              │
│  │  └─────────┘ └─────────┘ └─────────┘     │                              │
│  └────────┬────────────────────────────────┘                              │
│           │ canonical CAF events:                                         │
│           │  pre-session / post-edit / pre-commit                          │
└───────────┼───────────────────────────────────────────────────────────────┘
            │
┌───────────┼───────────────────────────────────────────────────────────────┐
│  Storage Layer                       │                                     │
│  ┌────────▼─────────────────────────┐                                      │
│  │  .caf/  (per-repo, in git)       │                                      │
│  │  ├─ context.md                    │                                      │
│  │  ├─ manifest.yml                  │                                      │
│  │  ├─ sessions/*.yml                │                                      │
│  │  ├─ decisions/*.md                │                                      │
│  │  └─ mental-model/                 │                                      │
│  │     ├─ cards.yml                  │                                      │
│  │     └─ graph.yml                  │                                      │
│  └──────────────────────────────────┘                                       │
└───────────────────────────────────────────────────────────────────────────┘
            ▲
            │ canonical events fire through adapters
┌───────────┼───────────────────────────────────────────────────────────────┐
│  Hooks Layer (canonical CAF contract)                                     │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐                   │
│  │ pre-session  │   │  post-edit   │   │ pre-commit   │                   │
│  │ reads        │   │  writes      │   │  validates   │                   │
│  │ context.md   │   │  provenance  │   │  decisions   │                   │
│  └──────────────┘   └──────────────┘   └──────────────┘                   │
└───────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Why this composition

Хранение `.caf/` в самом репозитории (а не в отдельной БД) гарантирует, что provenance и decisions эволюционируют синхронно с кодом: `git checkout` на старый коммит даёт доступ к ментальной карте того момента. Agent Adapter Layer минимизирует трение — автор не должен явно вызывать `caf log`, журнал сессий пишется автоматически независимо от того, какой агент редактирует код. CLI и MCP-сервер дают два режима потребления: CLI для дневной работы в терминале, MCP — для автоматического инжектирования контекста в новые AI-сессии.

Ядро системы — три модуля: **Provenance Tracker**, **Decision Log**, **Mental Model Reconstructor**. Они описаны в секции 4.

---

## 3.3 Agent Adapter Layer

**EN summary.** The Adapter Layer is the v0.2 addition that makes the tool universal. Each AI coding agent has its own hook system, configuration format, and event vocabulary. Adapters translate between agent-specific concepts and the canonical CAF event contract (pre-session / post-edit / pre-commit). Adapters declare a capability matrix; the core never assumes an adapter supports all events — it gracefully degrades when an agent lacks support.

### 3.3.1 Canonical event contract

CAF определяет три канонических события, которые любой адаптер должен по возможности реализовать:

| Event | Trigger | Input | Action |
|---|---|---|---|
| `pre-session` | запуск AI-сессии агентом | `repo_path`, `session_id`, `agent_id` | читать `.caf/context.md`, инжектировать в system prompt агента |
| `post-edit` | после каждой итерации агента, изменившей файлы | `repo_path`, `session_id`, `diff` | писать `.caf/sessions/<id>.yml`, при нетривиальном изменении — TUI-промпт для decision |
| `pre-commit` | перед `git commit` | `repo_path`, staged diff | валидировать наличие trailer `CAF-Session` и ссылки на decision |

### 3.3.2 Adapter capability matrix

Не все агенты поддерживают все три события одинаково. Адаптер объявляет capability matrix; ядро CAF использует fallback-стратегии для отсутствующих capabilities.

| Agent | `pre-session` | `post-edit` | `pre-commit` | MCP native | Adapter status (v0.2) |
|---|---|---|---|---|---|
| opencode | ✓ (hooks.toml) | ✓ | ✓ | ✓ | Tier 1 — MVP |
| Claude Code | ✓ (settings.json hooks) | ✓ | ✓ | ✓ | Tier 1 — MVP |
| Cursor | partial (rules file) | ✓ (extension API) | ✓ (git hook) | ✗ | Tier 1 — ALPHA |
| Cline | ✓ (extension hooks) | ✓ | ✓ | ✗ | Tier 1 — ALPHA |
| Aider | ✗ (no pre-session) | ✓ (via .aider.conf) | ✓ (git hook) | ✗ | Tier 2 — BETA |
| KiloCode | ✓ | ✓ | ✓ | ✓ | Tier 2 — BETA |
| Codex CLI | partial | ✓ | ✓ | ✗ | Tier 2 — BETA |
| Continue.dev | ✓ (config.yaml hooks) | ✓ | ✓ | ✗ | Tier 2 — BETA |
| OryaAgent | ✓ (native) | ✓ (native) | ✓ (native) | ✓ | Tier 3 — PROD |

### 3.3.3 Fallback strategies

Для агентов с неполной поддержкой событий:

- **Нет `pre-session`** (Aider, частично Codex CLI): wrapper-скрипт `caf-wrap <agent> <prompt>` запускает агент, предварительно инжектировав `context.md` через `--system-prompt-file` или эквивалент. Если агент не поддерживает system prompt override — `context.md` печатается в stdout как преамбула, агент читает её как первый turn.
- **Нет `post-edit` (real-time)**: filesystem watcher daemon `caf-watch` отслеживает изменения в рабочем каталоге и пишет session-лог при каждом flush. Задержка до 5 секунд; итерации агрегируются по временным окнам.
- **Нет `pre-commit` (git hook)**: невозможно — `pre-commit` всегда реализуем через `.git/hooks/pre-commit`. Это universal fallback, работающий для любого агента.

### 3.3.4 Adapter interface (Go)

```go
// pkg/adapter/adapter.go
package adapter

type Event string

const (
    EventPreSession  Event = "pre-session"
    EventPostEdit     Event = "post-edit"
    EventPreCommit    Event = "pre-commit"
)

type Adapter interface {
    // ID returns the agent identifier (e.g. "opencode", "claude-code").
    ID() string

    // Capabilities returns the events this adapter supports natively.
    // Core uses fallbacks for missing capabilities.
    Capabilities() map[Event]Capability

    // Install registers the adapter's hook configuration in the agent.
    // Called once by `caf init --agent <id>`.
    Install(repoPath string) error

    // Uninstall removes the adapter's hook configuration.
    Uninstall(repoPath string) error

    // Detect checks whether this adapter's agent is active in the repo.
    // Returns true if e.g. .opencode/ or .claude/ exists.
    Detect(repoPath string) bool
}

type Capability struct {
    Native      bool   // true = real hook, false = fallback
    FallbackDesc string // human-readable explanation when Native=false
}
```

### 3.3.5 Adapter registry

Адаптеры регистрируются в `internal/registry/registry.go` через `init()`:

```go
package registry

import (
    "codeauthorityfirst/pkg/adapter"
    "codeauthorityfirst/pkg/adapter/opencode"
    "codeauthorityfirst/pkg/adapter/claudecode"
    "codeauthorityfirst/pkg/adapter/cursor"
    "codeauthorityfirst/pkg/adapter/cline"
    "codeauthorityfirst/pkg/adapter/aider"
    "codeauthorityfirst/pkg/adapter/kilocode"
    "codeauthorityfirst/pkg/adapter/codexcli"
    "codeauthorityfirst/pkg/adapter/continue_dev"
    "codeauthorityfirst/pkg/adapter/oryaagent"
)

func init() {
    for _, a := range []adapter.Adapter{
        opencode.New(),
        claudecode.New(),
        cursor.New(),
        cline.New(),
        aider.New(),
        kilocode.New(),
        codexcli.New(),
        continue_dev.New(),
        oryaagent.New(),
    } {
        Register(a)
    }
}
```

Команда `caf adapters list` показывает зарегистрированные адаптеры и их capabilities.

---

## 4. Module Specifications

### 4.1 Provenance Tracker

**EN summary.** Each AI agent session receives a unique `session_id` and auto-logs to `.caf/sessions/<id>.yml` with the agent identifier, original prompt, iterations, files touched, final commit SHA, and references to decisions made. `git blame` is patched via a `post-commit` trailer `CAF-Session: <id>` so plain `git log --grep=CAF-Session` finds provenance regardless of which agent produced the commit.

Provenance Tracker отвечает на вопрос «откуда взялся этот код». Каждая сессия любого поддерживаемого агента получает уникальный `session_id` и автоматически логирует себя в `.caf/sessions/<id>.yml`. Запись содержит идентификатор агента, исходный промпт, итерации, затронутые файлы, финальный коммит и ссылку на git-объект.

**Схема записи сессии (v0.2 — добавлено поле `agent`):**

```yaml
# .caf/sessions/2026-10-24-a1b2c3.yml
session_id: 2026-10-24-a1b2c3
repo: my-backend-service
agent: claude-code        # NEW in v0.2: which agent produced this session
agent_version: 1.0.21
model: claude-sonnet-4-5
started_at: 2026-10-24T11:02:14+03:00
ended_at: 2026-10-24T11:47:33+03:00
tool_call: claude-code run "add rate limiter to /v1/policies"
initial_prompt: |
  Нужно добавить rate limiter на endpoint POST /v1/policies.
  Сейчас любой клиент может спамить запросами. Ограничение:
  10 запросов в минуту на токен.
iterations:
- id: i1
  action: read
  target: internal/api/policies.go
- id: i2
  action: edit
  target: internal/api/policies.go
  diff_summary: "добавлен middleware RateLimit(token, 10/min)"
- id: i3
  action: edit
  target: internal/middleware/ratelimit.go  # новый файл
files_touched:
- internal/api/policies.go
- internal/middleware/ratelimit.go
commit_sha: 9f3e2a1  # null если не закоммичено
decisions_made:
- d_2026-10-24-01  # ссылка на .caf/decisions/d_2026-10-24-01.md
voice_signature: "rate limiter нужен был срочно, выбрал token bucket"
```

**Связь с `git blame`.** Provenance Tracker патчит `git blame` через `post-commit` hook: каждый коммит, сделанный через любого AI-агента, помечается trailer `CAF-Session: 2026-10-24-a1b2c3` и `CAF-Agent: claude-code`. Это позволяет `git log --grep=CAF-Session` находить происхождение кода без отдельного инструмента, а `git log --grep=CAF-Agent=claude-code` фильтровать по конкретному агенту.

**CLI-команды модуля:**

```bash
# Кто написал эту строку, в какой сессии и каким агентом?
caf who internal/api/policies.go:42
# → session 2026-10-24-a1b2c3, agent: claude-code,
#   prompt: "add rate limiter", iteration i2

# Когда в этом файле появился паттерн X?
caf when --pattern "token_bucket" internal/middleware/
# → 2026-10-24, session a1b2c3, agent: claude-code, commit 9f3e2a1

# Почему (вытягивает decision log для этого файла)
caf why internal/middleware/ratelimit.go
# → d_2026-10-24-01: выбран token bucket вместо sliding window

# Все сессии определённого агента
caf sessions --agent claude-code --since 7d
```

Provenance — это базовый слой. Без него два других модуля не работают.

---

### 4.2 Decision Log

**EN summary.** ADR-like records stored in `.caf/decisions/d_<date>-<seq>.md`. Triggered by `post-edit` adapter event after non-trivial edits (>50 LOC, excluding imports / formatting / test-duplicates) and enforced by `pre-commit`. The decision template is agent-agnostic; the same format applies whether Claude Code, Aider, or Cursor produced the change.

Decision Log — это архив отвергнутых альтернатив. Каждое нетривиальное решение получает файл `.caf/decisions/<id>.md` в формате, похожем на ADR, но с одним принципиальным отличием: **главная ценность — не в описании выбранного варианта, а в фиксации отвергнутых**.

**Схема decision-записи:**

```markdown
---
id: d_2026-10-24-01
title: Rate limiter: token bucket vs sliding window
date: 2026-10-24
repo: my-backend-service
agent: claude-code          # NEW in v0.2
session: 2026-10-24-a1b2c3
files:
- internal/middleware/ratelimit.go
status: accepted
---

## Контекст
POST /v1/policies принимает запросы создания политик. Без ограничения
любой валидный токен может выполнить DoS, забив очередь обработки.
Нужно ограничение 10 запросов в минуту на токен.

## Выбрано
Token bucket: 10 токенов в бакете, пополнение 1 токен/6 секунд.

## Альтернативы

### Sliding window
- Что: считать запросы в скользящем окне 60 секунд
- Почему отверг: при 10 req/min накладные расходы на хранение
  временных меток не оправданы.

### Fixed window
- Что: считать запросы в фиксированных минутных окнах
- Почему отверг: на границе окна возможен всплеск 2x (10+10),
  что для политик критично.

### Redis-based distributed limiter
- Что: вынести счётчик в Redis
- Почему отверг: на данном этапе нет требования multi-instance.

## Импликации
- Если нагрузка вырастет до 100+ req/sec — пересмотреть в сторону
  sliding window + Redis.
- Интерфейс Limiter абстрактен, миграция без боли.
- Метрики: размер бакета вынесен в конфиг, можно тюнить без релиза.
```

**Триггеры записи.** Decision log не пишется вручную. Записи инициируются событием `post-edit` от адаптера: после каждой AI-сессии, изменившей >50 строк нетривиального кода, ядро CAF предлагает автору за 30 секунд заполнить шаблон через TUI-промпт. Если автор пропускает — событие `pre-commit` напоминает ещё раз и блокирует коммит, пока не появится ссылка на decision (или явная пометка `CAF-Trivial: true`).

**CLI-команды модуля:**

```bash
# Все решения за сегодня
caf decisions --today

# Решения, влияющие на конкретный файл
caf decisions --file internal/middleware/ratelimit.go

# Решения, сделанные конкретным агентом
caf decisions --agent aider --since 30d

# Решения без связанных коммитов (прорехи в логе)
caf decisions --gaps

# Поиск по ключевому слову среди отвергнутых альтернатив
caf decisions --search "Redis"

# Интерактивное создание новой записи
caf decisions new --files internal/middleware/ratelimit.go
```

---

### 4.3 Mental Model Reconstructor (MMR)

**EN summary.** Converts passive code-reading ability into active ownership via SM-2 spaced repetition over three card types: "Why X here?", "What if Y is removed?", "Where is Z handled?". Cards are generated automatically from `.caf/sessions/` and `.caf/decisions/`. A cross-repo dependency graph in `.caf/mental-model/graph.yml` enables cards of types 2 and 3, and surfaces "blindspots". Cards are agent-agnostic — a card may reference sessions from multiple agents if the file was touched by different agents over time.

MMR применяет алгоритм SM-2 к трём типам карточек, генерируемых из `.caf/sessions/` и `.caf/decisions/`. Карточки не привязаны к конкретному агенту — если файл редактировался разными агентами в разное время, карточка ссылается на все релевантные сессии.

**Тип 1 — «Почему здесь X?»** Карточка привязана к конкретному файлу и строке. Вопрос: «Почему в `internal/middleware/ratelimit.go:18` выбран token bucket?». Ответ сверяется с `d_2026-10-24-01`. Генерируется автоматически при создании decision-записи.

**Тип 2 — «Что будет если удалить Y?»** Карточка симулирует удаление ключевого компонента. Пример: «Что произойдёт, если удалить `policy.Cache`?». Ответ: «Латентность eval вырастет с 2ms до 40ms». Генерируется из графа зависимостей.

**Тип 3 — «Где обрабатывается Z?»** Карточка проверяет знание потока данных. Пример: «Где обрабатывается событие `tool_call_failed`?». Выявляет слепые зоны.

**Алгоритм SM-2.** Каждая карточка имеет `ease_factor` (начально 2.5) и `interval` (начально 1 день). После ответа автор оценивает качество 0-5. Оценка <3 сбрасывает интервал, ≥4 увеличивает на `ease_factor`. Стандартный SM-2 без модификаций.

**Формат хранения:**

```yaml
# .caf/mental-model/cards.yml
cards:
- id: c_001
  type: why
  repo: my-backend-service
  file: internal/middleware/ratelimit.go
  line: 18
  question: "Почему выбран token bucket?"
  expected_answer: "10 req/min, накладные расходы sliding window не оправданы"
  decision_ref: d_2026-10-24-01
  sessions:              # NEW in v0.2: multiple sessions may contribute
  - 2026-10-24-a1b2c3    # claude-code session
  - 2026-11-02-e5f6a7    # aider session that touched same line
  ease: 2.5
  interval_days: 1
  due: 2026-10-25
  last_reviewed: null
  reviews: []
```

**Граф зависимостей.** MMR строит и поддерживает кросс-репо граф: импорты, вызовы интерфейсов, события. Граф хранится в `.caf/mental-model/graph.yml` и обновляется при каждом `caf sync`.

**CLI-команды модуля:**

```bash
# Ежедневная сессия повторения (10 минут, ~15 карточек)
caf review

# Жёсткая сессия: только карточки с низким ease
caf review --hard

# Просмотр слепых зон
caf review --blindspots

# Принудительно сгенерировать карточки из последних решений
caf review --regenerate --since 7d

# Статистика владения
caf review --stats
```

---

## 5. Data Model — `.caf/` Contracts

**EN summary.** `.caf/` lives at the root of any repo and is committed alongside code. Same structure as v0.1's `.orya/` but with `agent` field added to all session/decision/card records. Five subdirectories + three top-level files. All schemas are YAML/Markdown to maximize diffability.

### 5.1 File tree

```text
.caf/
├── context.md              # автоматически генерируемый контекст для pre-session
├── manifest.yml            # мета-информация о репо и связях с другими репо
├── config.yml              # NEW in v0.2: per-repo CAF configuration
│                           #   (enabled agents, capability overrides)
├── sessions/               # Provenance Tracker
│   ├── 2026-10-24-a1b2c3.yml
│   └── 2026-10-25-d4e5f6.yml
├── decisions/              # Decision Log
│   ├── d_2026-10-24-01.md
│   ├── d_2026-10-24-02.md
│   └── index.yml           # быстрый индекс для поиска
├── mental-model/           # Mental Model Reconstructor
│   ├── cards.yml           # SRS-карточки
│   └── graph.yml           # кросс-репо граф зависимостей
└── reviews/                # история caf review сессий
    └── 2026-10-24.log
```

### 5.2 `config.yml` (NEW in v0.2)

Per-repo конфигурация CAF, в частности — какие агенты включены и какие capability overrides применяются:

```yaml
# .caf/config.yml
version: 0.2
enabled_agents:
- opencode
- claude-code
- cursor              # fallback pre-session via wrapper
- aider               # fallback pre-session via wrapper
capability_overrides:
  cursor:
    pre-session:
      native: false
      fallback: wrapper
  aider:
    pre-session:
      native: false
      fallback: wrapper
triviality_threshold_loc: 50
hooks_fail_safe: true
```

### 5.3 `manifest.yml`

Связующий файл между репо. Пример для произвольного бэкенд-сервиса:

```yaml
# .caf/manifest.yml в репозитории my-backend-service
repo: my-backend-service
role: service  # service | library | cli | orchestrator | etc.
description: |
  Backend API for policy management. Go + PostgreSQL.
  Handles POST /v1/policies, GET /v1/policies/{id}, etc.
depends_on:
- repo: my-auth-lib            # внутренняя библиотека
  interface: github.com/myorg/my-auth-lib.TokenValidator
  used_in: internal/api/policies.go:42
  decision: d_2026-09-15-03    # почему именно этот интерфейс
- repo: my-telemetry-lib
  interface: github.com/myorg/my-telemetry-lib.Tracer
  used_in: internal/middleware/tracing.go:18
- repo: postgres               # внешняя зависимость
  interface: database/sql
  used_in: internal/storage/pg.go
exposes:
- interface: PolicyAPI
  package: github.com/myorg/my-backend-service/pkg/api
  consumers: [my-frontend, my-cli-tool]
principles:
- "Сервис не содержит бизнес-логики в handlers, только маршрутизацию"
- "Все долгоживящие состояния — в PostgreSQL, не в памяти"
- "Каждая ошибка логируется с trace_id"
primary_agents:                # NEW in v0.2: which agents typically work here
- claude-code
- cursor
```

`depends_on[].decision` — обратная ссылка на `.caf/decisions/`, объясняющая выбор интерфейса. При `caf sync` инструмент проверяет, что все referenced decisions существуют.

### 5.4 `context.md`

Генерируется командой `caf context` и перегенерируется автоматически при каждом коммите. Содержит выжимку для AI-агента: принципы, последние 10 decisions, активные слепые зоны, интерфейсы, которые репо экспортирует/потребляет. **Длина — строго до 2000 токенов**.

Алгоритм генерации:

1. Прочитать `manifest.yml` → секция `principles` + `exposes` + `depends_on` + `primary_agents`.
2. Прочитать `decisions/index.yml`, взять последние 10 по полю `date`.
3. Прочитать `mental-model/graph.yml`, вычислить blindspots.
4. Скомпоновать markdown ≤2000 токенов; приоритезация: `principles > blindspots > decisions > exposes`.

### 5.5 `sessions/*.yml`

См. секцию 4.1. Именование: `<YYYY-MM-DD>-<short_id>.yml`, где `short_id` — первые 6 hex-символов SHA-256 от `started_at + initial_prompt + agent_id`. Добавление `agent_id` в хэш гарантирует уникальность даже если два разных агента запускаются в одну секунду с одинаковым промптом.

### 5.6 `decisions/*.md` и `decisions/index.yml`

См. секцию 4.2 для формата отдельной записи. `index.yml`:

```yaml
# .caf/decisions/index.yml
- id: d_2026-10-24-01
  title: Rate limiter: token bucket vs sliding window
  date: 2026-10-24
  repo: my-backend-service
  agent: claude-code            # NEW in v0.2
  files:
  - internal/middleware/ratelimit.go
  tags: [rate-limiting, middleware, token-bucket]
  status: accepted
```

### 5.7 `mental-model/cards.yml` и `graph.yml`

См. секцию 4.3 для формата карточек. `graph.yml` — кросс-репо граф зависимостей:

```yaml
# .caf/mental-model/graph.yml
nodes:
- repo: my-backend-service
  file: pkg/policy/evaluator.go
  symbol: Cache
  kind: struct
- repo: my-frontend
  file: src/api/policies.ts
  symbol: fetchPolicies
  kind: usage
edges:
- from: {repo: my-frontend, file: src/api/policies.ts, line: 12}
  to:   {repo: my-backend-service, file: pkg/api/policy.go, symbol: PolicyAPI}
  kind: http_call
  decision: d_2026-09-15-03
last_synced: 2026-10-24T11:00:00+03:00
```

Граф строится language-agnostically. Для Go: `go list -json ./...`. Для Python: `pip show` + AST-парсинг импортов. Для TypeScript: `tsc --listFiles` + `tsconfig` paths. Language adapters добавляются по фазам roadmap (см. ROADMAP.md).

---

## 6. Hooks Contracts

**EN summary.** Three canonical CAF events: `pre-session` (injects `context.md` into agent's system prompt), `post-edit` (writes provenance draft after each iteration; for >50 LOC non-trivial edits, prompts author for a decision record), `pre-commit` (blocks commits lacking `CAF-Session` trailer or decision reference, unless `CAF-Trivial: true`). All events fire through adapters — the core never calls agent-specific APIs directly.

### 6.1 Pre-session event

**Trigger:** запуск AI-сессии агентом.
**Input:** `repo_path`, `session_id`, `agent_id`.
**Action:**

1. Читать `.caf/context.md`.
2. Если файл отсутствует или старше последнего коммита — запустить `caf context --refresh`.
3. Через адаптер инжектировать содержимое в system prompt агента (способ зависит от агента — см. capability matrix).

**Effect:** Агент с первой итерации знает об архитектурных решениях репо и не предложит отвергнутые альтернативы заново.

**Adapter implementations:**

| Agent | Способ инжекции |
|---|---|
| opencode | `hooks.toml` pre-session hook → `caf hook pre-session` выводит context в stdout, opencode читает как system prompt |
| Claude Code | `settings.json` PreToolUse hook → `caf hook pre-session` writes to `.claude/context.txt`, Claude Code reads via `--system-prompt-file` |
| Cursor | wrapper `caf-wrap cursor` prependует context в первый промпт |
| Cline | extension API `onActivate` → injects `context.md` content into Cline's system prompt slot |
| Aider | wrapper `caf-wrap aider` prependует context в `--message` |
| KiloCode | hooks config → `caf hook pre-session` |
| Codex CLI | wrapper (pre-session fallback) |
| Continue.dev | `config.yaml` hooks → `caf hook pre-session` |
| OryaAgent | native first-class integration |

### 6.2 Post-edit event

**Trigger:** после каждой итерации агента, изменившей файлы.
**Input:** `repo_path`, `session_id`, `diff`.
**Action:**

1. Автоматически создать/обновить `.caf/sessions/<session_id>.yml` с затронутыми файлами и diff-сводкой.
2. Подсчитать «нетривиальные» изменения: >50 строк, исключая imports, форматирование, тесты-дубликаты.
3. Если порог пройден — прервать сессию (если адаптер поддерживает) и предложить автору за 30 секунд заполнить decision-шаблон через TUI-промпт.
4. Автор может пропустить с пометкой `trivial: true`; в этом случае `pre-commit` спросит ещё раз.

### 6.3 Pre-commit event

**Trigger:** перед `git commit`.
**Input:** staged diff.
**Action:**

1. Проверить наличие trailer `CAF-Session: <id>` в message коммита.
2. Если коммит нетривиален и не ссылается ни на один decision — **заблокировать коммит**.
3. Вывести подсказку: «У этого изменения нет записи в Decision Log. Создать? (y/n/t — trivial)».
4. При `y` — запустить `caf decisions new` интерактивно.
5. При `t` — добавить trailer `CAF-Trivial: true` и разрешить коммит.
6. При `n` — отказ от коммита.

**Universal implementation:** `pre-commit` всегда реализуется через `.git/hooks/pre-commit` shell-скрипт, устанавливаемый `caf init`. Этот скрипт agent-agnostic — git hook не знает, какой агент редактировал файлы, он проверяет только trailer и decision-ссылку.

### 6.4 Adapter installation

Команда `caf init --agent <id>` устанавливает адаптер для конкретного агента в текущем репо:

```bash
# Установить один адаптер
caf init --agent claude-code

# Установить несколько адаптеров
caf init --agents opencode,claude-code,cursor

# Автоопределение: CAF сканирует репо на наличие .opencode/, .claude/, .cursor/
# и предлагает установить соответствующие адаптеры
caf init --auto-detect
```

Для каждого адаптера генерируется конфигурация в соответствующем месте:

| Agent | Конфигурационный файл |
|---|---|
| opencode | `~/.config/opencode/hooks.toml` |
| Claude Code | `.claude/settings.json` (committed to repo) |
| Cursor | `.cursor/rules/caf.mdc` (committed to repo) |
| Cline | `.cline/config.json` (committed to repo) |
| Aider | `.aider.conf.yml` (committed to repo) |
| KiloCode | `.kilocode/settings.json` |
| Codex CLI | `~/.codex/config.toml` |
| Continue.dev | `.continue/config.yaml` (committed to repo) |
| OryaAgent | `.orya/settings.yml` |

### 6.5 Fail-safe semantics

Хуки спроектированы так, чтобы не ломать workflow при сбое:

- Если `caf` binary недоступен → агент продолжает работать, в логе warning, сессия не блокируется.
- Если `.caf/` отсутствует → hook молча пропускается (репо не инициализирован через `caf init`).
- Если адаптер не поддерживает `pre-session` → fallback через wrapper; если wrapper не используется, `pre-session` пропускается с warning.
- При конфликте `CAF-Trivial: true` и нетривиального диффа → `pre-commit` всё равно требует decision или явного trailer.

См. ADR-003 для обоснования.

---

## 7. Sequence Flows

**EN summary.** Three canonical flows: (1) prompt → git commit through all three modules in a daily workflow; (2) SRS card lifecycle from generation to review and ease update; (3) cross-repo change propagation when an interface in one repo is modified, surfacing potentially stale decisions in dependents. All flows are agent-agnostic — they work identically whether Claude Code, Cursor, or Aider produced the change.

### 7.1 Prompt → git commit (дневной поток данных)

```text
[author]            [any agent]      [adapter]           [.caf/]
   │                    │                │                  │
   │  agent run         │                │                  │
   ├───────────────────►│ pre-session    │                  │
   │                    ├───────────────►│ read context.md  │
   │                    │                ├─────────────────►│
   │                    │                │◄─────────────────┤ markdown
   │                    │ inject to      │                  │
   │                    │ system prompt  │                  │
   │                    │                │                  │
   │  (agent works)     │                │                  │
   │                    │ iteration i1   │                  │
   │                    ├───────────────►│ post-edit        │
   │                    │                ├─────────────────►│ write sessions/
   │                    │                │                  │ <id>.yml
   │                    │ iteration i2   │                  │
   │                    ├───────────────►│ post-edit        │
   │                    │                │ >50 LOC? → TUI   │
   │                    │                │ prompt author    │
   │  (fill decision)   │                │                  │
   │◄───────────────────┤                │                  │
   │                    │                ├─────────────────►│ write decisions/
   │                    │                │                  │ d_*.md
   │  git commit        │                │                  │
   ├───────────────────►│                │ pre-commit       │
   │                    │ (git hook, universal)              │
   │                    │ validate trailer + decision ref   │
   │                    │ pass                              │
   │  commit created    │                │                  │
   │◄───────────────────┤                │ post-commit      │
   │                    │                ├─────────────────►│ add trailers
   │                    │                │                  │ CAF-Session,
   │                    │                │                  │ CAF-Agent
```

### 7.2 SRS card lifecycle

```text
[decision created]                 [MMR generator]              [cards.yml]
       │                                  │                          │
       │  d_2026-10-24-01.md              │                          │
       ├─────────────────────────────────►│ generate type-1 card     │
       │                                  ├─────────────────────────►│ c_001
       │                                  │ ease=2.5, interval=1d   │
       │                                  │ due=2026-10-25           │
       │                                  │                          │
       │  [next day, caf review]           │                          │
       │  due ≤ today? → show c_001        │                          │
       │  author answers, rates 4/5        │                          │
       │  interval = 1 * 2.5 = 2.5 → 3d   │                          │
       │  due = 2026-10-28                │                          │
       │                                  │                          │
       │  [graph update via caf sync]      │                          │
       │  graph.yml rebuilt                │                          │
       │  → generate type-2/3 cards        │                          │
       │  from new graph nodes            │                          │
       │                                  ├─────────────────────────►│ c_002, c_003
       │                                  │                          │
       │  [30 days later]                  │                          │
       │  file X changed but no card →    │                          │
       │  blindspot                       │                          │
       │  caf review --blindspots         │                          │
       │  → regenerate cards for X         │                          │
```

### 7.3 Cross-repo change propagation

```text
[Repo A: my-auth-lib]               [Repo B: my-backend-service]
   │                                            │
   │  interface TokenValidator changed          │
   │  git commit (with CAF-Session trailer)    │
   │                                            │
   │  caf sync (in Repo B)                      │
   │◄───────────────────────────────────────────┤
   │  reads my-auth-lib/.caf/manifest.yml       │
   │  detects: used_in internal/api/policies.go:42│
   │           references outdated interface    │
   │                                            │
   │  flags decision d_2026-09-15-03            │
   │  as "potentially stale"                    │
   │                                            │
   │  surfaces in next caf review as            │
   │  "review-because-interface-changed"        │
   │  card type                                 │
```

---

## 8. MCP Interface

**EN summary.** MCP server `caf-mcp` exposes 5 tools and 3 prompts to any MCP-compatible client. The flagship tool `caf.first_person_explain` returns a first-person narrative explaining a file's history, decisions, and rejected alternatives. MCP tools are agent-agnostic — they read from `.caf/` regardless of which agent produced the data.

MCP-сервер `caf-mcp` выставляет инструменты и промпты, которые может вызывать любой MCP-совместимый клиент: opencode, Claude Code, Cursor (через MCP bridge), Cline, KiloCode, OryaAgent.

### 8.1 MCP-tools

| Tool | Описание | Возвращает |
|---|---|---|
| `caf.context` | Текущий контекст репо (то же, что `context.md`) | Markdown, ~2000 токенов |
| `caf.provenance` | История сессий, затронувших файл/строку | Список сессий с промптами и идентификаторами агентов |
| `caf.decisions` | Поиск decisions по тегам/файлам/репо/агенту | Список decision summary |
| `caf.mental_model` | SRS-карточки для текущего пользователя | JSON с due-карточками |
| `caf.first_person_explain` | **Главное:** объясняет файл от первого лица | Markdown-нарратив |

### 8.2 `caf.first_person_explain` — flagship tool

Принимает путь к файлу и возвращает объяснение, написанное от лица автора:

```text
Tool: caf.first_person_explain
Input: { "file": "internal/middleware/ratelimit.go", "depth": "medium" }

Output (markdown):

## internal/middleware/ratelimit.go

Я добавил этот файл 24 октября, в сессии с Claude Code про rate
limiting на POST /v1/policies (session 2026-10-24-a1b2c3).
На тот момент у нас был открытый endpoint без ограничений, и
любой валидный токен мог заDoSить очередь обработки.

Я выбрал token bucket на 10 токенов с пополнением 1/6 сек.
Смотрел sliding window, но при 10 req/min накладные расходы
на хранение меток не оправданы. Redis-based вариант отверг,
потому что multi-instance пока не нужен, а точка отказа
не нужна. Интерфейс Limiter абстрактен, перенести в Redis
можно без боли.

Позже (session 2026-11-02-e5f6a7, Aider) я вернулся сюда,
чтобы добавить метрики — размер бакета вынесен в конфиг,
можно тюнить без релиза. Метрики экспортируются в
my-telemetry-lib как ratelimit_deny_total.
```

Глубина `depth` принимает значения `short` (1 параграф, только решение), `medium` (решение + главные альтернативы), `long` (полная археология: все итерации сессии, все отвергнутые альтернативы, метрики). Narратив генерируется LLM, но фактические данные берутся строго из `.caf/`, что устраняет риск галлюцинаций на фактологической части.

### 8.3 MCP-prompts

- `caf.pre_session_brief` — что показать агенту перед началом работы. Включает `context.md` + список активных blindspots + последние 3 decisions по теме предстоящей задачи.
- `caf.code_review_with_provenance` — ревью коммита с учётом исходного промпта и отвергнутых альтернатив.
- `caf.archaeology` — раскопки конкретного файла: когда, почему, какие альтернативы, какими агентами редактировался.

---

## 9. CLI Reference

```bash
# === Инициализация ===
caf init                          # создать .caf/ в текущем репо
caf init --agent <id>             # установить адаптер для конкретного агента
caf init --agents a,b,c           # установить несколько адаптеров
caf init --auto-detect            # автоопределение агентов по наличию конфиг-файлов

# === Adapter management (NEW in v0.2) ===
caf adapters list                 # показать зарегистрированные адаптеры + capabilities
caf adapters install <id>         # установить адаптер (alias для caf init --agent)
caf adapters uninstall <id>       # удалить адаптер
caf adapters detect               # определить активные агенты в репо

# === Provenance ===
caf who <file:line>               # кто, в какой сессии и каким агентом написал строку
caf when --pattern <regex> <dir>  # когда появился паттерн
caf why <file>                    # связать файл с decision log
caf sessions --since 7d           # список сессий за период
caf sessions --agent <id>         # сессии конкретного агента
caf sessions --gaps               # коммиты без session trailer

# === Decision Log ===
caf decisions                     # все решения (по умолчанию за 30 дней)
caf decisions --today             # только сегодняшние
caf decisions --file <path>       # влияющие на файл
caf decisions --agent <id>        # NEW: решения конкретного агента
caf decisions --search <query>    # полнотекстовый поиск
caf decisions --gaps              # нетривиальные коммиты без decision
caf decisions new [--files ...]   # интерактивное создание записи

# === Mental Model ===
caf review                        # ежедневная SRS-сессия (~10 мин)
caf review --hard                 # только карточки с низким ease
caf review --blindspots           # файлы без SRS-карточек
caf review --regenerate --since 7d
caf review --stats                # статистика владения

# === Контекст и синхронизация ===
caf context                       # вывести context.md в stdout
caf context --refresh             # перегенерировать context.md
caf sync                          # обновить кросс-репо граф
caf tour <file>                   # guided-экскурсия по файлу

# === Hook-интерфейс (вызывается адаптерами) ===
caf hook pre-session --repo <path> --agent <id>
caf hook post-edit --repo <path> --session <id> --agent <id>
caf hook pre-commit --repo <path>

# === Отчёты ===
caf report --day                  # дневной отчёт
caf report --week                 # недельная ретроспектива
caf report --ownership            # heatmap владения по файлам
caf report --agent <id>           # NEW: отчёт по конкретному агенту
caf report --metrics              # сводные метрики успеха
```

Все команды принимают `--repo <path>` для кросс-репо операций. По умолчанию репо определяется по текущему каталогу.

---

## 10. Cross-repo Topology

**EN summary.** Any set of repos developed by the same author forms one product from the mental-map perspective. Topology is encoded in each repo's `.caf/manifest.yml` and aggregated by `caf sync` into `.caf/mental-model/graph.yml`. The topology is not restricted to a single org — a developer's personal projects, work projects, and open-source contributions can all reference each other via local clones.

### 10.1 Topology example (universal)

```text
                    ┌──────────────────────────┐
                    │  my-backend-service      │  role: service
                    │  primary_agents:         │  exposes: PolicyAPI
                    │  - claude-code           │
                    │  - cursor                 │
                    └────────────┬─────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
              ▼                  ▼                  ▼
   ┌──────────────────┐ ┌────────────────┐ ┌────────────┐
   │  my-auth-lib     │ │ my-telemetry-  │ │ postgres   │
   │  (library)       │ │ lib (library)  │ │ (external) │
   │  primary_agents: │ │ primary_agents:│ └────────────┘
   │  - opencode      │ │ - aider        │
   │  - claude-code   │ │                │
   └──────────────────┘ └────────────────┘

                                 ▲
              ┌──────────────────┴──────────────────┐
              │                                     │
   ┌──────────────────┐                  ┌──────────────────┐
   │  my-frontend     │                  │  my-cli-tool     │
   │  (webapp)         │                  │  (cli)           │
   │  primary_agents:  │                  │  primary_agents:  │
   │  - cursor         │                  │  - claude-code   │
   │  - cline          │                  │  - aider         │
   └──────────────────┘                  └──────────────────┘
```

### 10.2 Roles

Roles не фиксированы — пользователь определяет их в `manifest.yml`:

| Repo | Role | Exposes | Consumes |
|---|---|---|---|
| `my-backend-service` | service | `PolicyAPI` | `my-auth-lib.TokenValidator`, `my-telemetry-lib.Tracer`, `postgres` |
| `my-auth-lib` | library | `TokenValidator`, `Limiter` | `postgres` |
| `my-telemetry-lib` | library | `Tracer`, `Metrics` | — |
| `my-frontend` | webapp | — | `my-backend-service.PolicyAPI` |
| `my-cli-tool` | cli | — | `my-backend-service.PolicyAPI` |

### 10.3 Sync algorithm

`caf sync` работает в три шага:

1. **Local scan.** Для текущего репо определяет language и запускает соответствующий adapter:
   - Go: `go list -json ./...`
   - Python: AST-парсинг импортов + `pip show`
   - TypeScript: `tsc --listFiles` + tsconfig paths
   - Rust: `cargo metadata` + `cargo tree`
2. **Manifest reconciliation.** Для каждого найденного external import читает соответствующий `<repo>/.caf/manifest.yml` (если репо клонирован локально) или pullит через git remote + raw URL.
3. **Graph update.** Записывает узлы и рёбра в `.caf/mental-model/graph.yml`, помечает stale decisions (те, где referenced interface изменился со времени последнего sync).

При обнаружении, что интерфейс в `depends_on` изменился, sync создаёт карточку типа `review-because-interface-changed` в `cards.yml`, которая появится в следующем `caf review`.

---

## 11. Trade-offs / ADRs

Подробное обоснование всех архитектурных решений вынесено в отдельный документ **`ADR.md`**. Краткий список ADR:

| ADR | Тема | Краткое решение |
|---|---|---|
| ADR-001 | Storage location | In-repo `.caf/` committed to git |
| ADR-002 | SRS algorithm | SM-2 (Anki-compatible) |
| ADR-003 | Hooks vs explicit CLI | Hooks with fail-safe fallbacks |
| ADR-004 | Triviality threshold | 50 LOC excluding imports/formatting/tests |
| ADR-005 | Schema format | YAML for structured, Markdown for human-readable |
| ADR-006 | First-person voice | First person in all generated narratives |
| ADR-007 | Project naming | `CodeAuthorityFirst` (universal, not Orya-specific) |
| ADR-008 | CLI binary name | `caf` (short, unique in $PATH) |
| ADR-009 | Agent adapter pattern | Pluggable adapters with capability matrix |
| ADR-010 | Fallback strategy | Wrapper scripts + filesystem watcher + universal git hook |
| ADR-011 | Multi-language support | Per-language graph adapters, Go-first |

Полные тексты ADR — см. `ADR.md`.

---

## 12. Success Metrics

**EN summary.** Same four metrics as v0.1, plus a new metric `Agent Diversity` measuring the distribution of sessions across agents — a healthy project uses multiple agents, not just one. All metrics must converge to target values within 8 weeks of usage.

Инструмент считается работающим, если пять метрик сходятся к целевым значениям за 8 недель использования.

### 12.1 Authorship Confidence Score (ACS)

Микро-опрос: после каждого `caf review` сессии автор отвечает на вопрос «Мог бы я написать этот код с нуля, если бы мне дали задачу?» от 1 до 5. Среднее за неделю — ACS.

**Target:** ≥3.5.

**Interpretation.** Ниже 3.0 → MMR генерирует слишком простые карточки, либо в проекте есть целые модули, которые автор не контролирует.

### 12.2 Decision Coverage

Доля нетривиальных коммитов со связанной decision-записью.

**Target:** ≥90%.

### 12.3 Mental Model Reach

Доля файлов в репо (взвешенная по LOC), покрытых хотя бы одной SRS-карточкой.

**Target:** ≥70%.

### 12.4 Recall Rate

Доля карточек в `caf review`, на которые автор ответил с оценкой ≥3.

**Target:** ≥75%.

### 12.5 Agent Diversity (NEW in v0.2)

Шенноновская энтропия распределения сессий по агентам за последние 30 дней. Высокая энтропия = автор использует несколько агентов осознанно (а не один агент для всего).

**Target:** H ≥ 1.0 (примерно соответствует 3 активным агентам с сопоставимым объёмом).

**Interpretation.** Если H близко к 0 — автор полагается на один агент, что повышает риск слепых зон (один агент имеет одинаковые biases). Если H очень высока (>2.0) — возможно, автор тратит слишком много времени на переключение контекста между агентами без ясной причины.

Метрики хранятся в `.caf/reviews/*.log` и агрегируются командой `caf report --metrics`.

---

## 13. Out of Scope (v0.2)

Сознательно вынесено за рамки текущей версии спецификации:

- **Web-дашборд.** CLI + MCP покрывают дневные задачи; визуализация графов отложена до FUTURE (см. ROADMAP.md).
- **Мульти-пользовательский режим.** Сейчас инструмент для одного автора; sync через `git push`/`pull` работает, но без слияния конфликтующих decisions. Multi-user требует per-author attribution в карточках.
- **Поддержка языков кроме Go в первой итерации.** Анализ импортов и графа зависимостей в v0.2 завязан на `go list`. Python/TypeScript/Rust добавляются по фазам roadmap.
- **Автоматическая генерация decision из LLM.** v0.2 требует ручного заполнения шаблона автором. LLM-assisted pre-fill — опциональная фича FUTURE с явным opt-in флагом `--ai-prefill`.
- **Real-time collaboration через CRDT.** Если несколько авторов одновременно работают в одном репо, конфликтующие `.caf/` изменения разрешаются через стандартный git merge. CRDT-семантика отложена.

Полный список отложенных идей — см. `BACKLOG.md`.
