# CodeAuthorityFirst — BACKLOG.md

> **Версия / Version:** 0.3 (80/20 prioritized, UI backlog added)
> **Назначение / Purpose:** catalog of ideas, features, and improvements considered but **not committed to any roadmap phase**. Items here are directional, not promises. Promotion to a roadmap phase requires: (1) clear problem statement, (2) design sketch, (3) pilot user, (4) success metric — see ROADMAP.md section 6.10.
> **Принцип приоритезации / Prioritization:** 🔥 80/20 — high-leverage items marked; these get disproportionate resources. See ADR-015.
> **Язык / Language:** bilingual — Russian body, English summary per theme
> **Аудитория / Audience:** maintainers, contributors, anyone proposing a feature

---

## 0. Backlog Philosophy

**EN summary.** This file is a parking lot — not a todo list. Items here exist because someone thought of them, but they have not been prioritized. An item can sit in backlog for years; that is fine. An item can be removed if it becomes obsolete or if a better idea supersedes it. Promotion to ROADMAP requires meeting four criteria (see ROADMAP.md §6.10).

Бэклог — это не todo-list, а парковка идей. Каждая запись существует, потому что кто-то её предложил, но она не получила приоритета в текущем roadmap. Это нормально — идея может лежать в бэклоге годами. Удаление из бэклога происходит, если идея устарела или если появилась лучшая альтернатива.

**Структура каждой записи:**

```markdown
### [ID] Title
- **Theme:** категория (см. секцию 1)
- **Status:** idea / draft / proposed / rejected
- **Origin:** кто и когда предложил
- **Problem:** что пользователь не может сделать сейчас
- **Sketch:** high-level idea (1-3 предложения)
- **Open questions:** что нужно решить
- **Dependencies:** на что опирается (другие backlog items, фазы roadmap)
```

---

## 1. Themes

| # | Theme | Description | 🔥 80/20 items |
|---|---|---|---|
| A | Performance & Scalability | оптимизация для больших репо, lazy loading, streaming | 2 |
| B | UX & DX | TUI improvements, IDE plugins, onboarding | 5 |
| C | Integrations | CI/CD, chat platforms, project management | 2 |
| D | AI & LLM Features | smart card generation, decision summarization, AI review | 1 |
| E | Languages | support beyond Go | 1 |
| F | Enterprise | SSO, audit, compliance, on-prem | 0 |
| G | Multi-user & Collaboration | team mode, real-time sync, conflict resolution | 0 |
| H | Extensibility | plugin API, custom card types, marketplace | 1 |
| I | Analytics & Reporting | advanced metrics, trend analysis, dashboards | 2 |
| J | Quality & Reliability | error handling, recovery, observability of CAF itself | 2 |
| **K** | **UI Enhancements (NEW v0.3)** | **TUI improvements beyond v0.3 baseline, web UI, IDE plugins** | **3** |

---

## 2. Theme A: Performance & Scalability

### A-001: Lazy graph loading

- **Status:** idea
- **Problem:** на репо >10k файлов `caf sync` занимает >30 секунд, что блокирует daily workflow.
- **Sketch:** загружать только узлы, reachable из current working set; остальное lazy-load по запросу. Использовать `git ls-files` + content-addressed cache.
- **Open questions:** какой размер working set по умолчанию? Как invalidate cache при `git checkout`?
- **Dependencies:** BETA (когда single-repo graph появится).

### A-002: Streaming `.caf/` parser

- **Status:** idea
- **Problem:** `cards.yml` на 10k+ карточек парсится целиком в память, что потребляет >100 MB.
- **Sketch:** streaming YAML parser (`gopkg.in/yaml.v3` уже поддерживает streaming для массивов). Index file `cards.idx.yml` для random access по `card_id`.
- **Dependencies:** BETA.

### A-003: Incremental `caf sync`

- **Status:** draft
- **Problem:** каждый `caf sync` пересканирует весь репо с нуля.
- **Sketch:** hash-токен последнего sync; сравнение с текущим `git diff <last_sync_sha>..HEAD`; обновление только изменённых узлов графа.
- **Open questions:** как обрабатывать force-push (история переписана)?
- **Dependencies:** PROD (когда cross-repo sync появится).

### A-004: Parallel adapter installation

- **Status:** idea
- **Problem:** `caf init --agents opencode,claude-code,cursor,aider` устанавливает адаптеры последовательно, занимает >10 секунд.
- **Sketch:** parallel installation с error isolation. Каждый адаптер в горутине, результаты агрегируются.
- **Dependencies:** ALPHA.

### A-005: Memory-mapped `.caf/sessions/` index

- **Status:** idea
- **Problem:** `caf who <file:line>` на репо с >5k sessions парсит все YAML файлы.
- **Sketch:** SQLite index (git-ignored, regenerated on demand) с mmaps для session files. Hybrid approach: text YAML для git diff, SQLite для queries.
- **Open questions:** конфликтует с ADR-001 (in-repo storage). Сделать `.caf/index.sqlite` git-ignored?
- **Dependencies:** BETA.

### A-006: Distributed `caf sync` for monorepo fleets

- **Status:** idea
- **Problem:** при наличии 50+ репо одного автора `caf sync` последовательно pullит каждый `manifest.yml`, что медленно.
- **Sketch:** parallel fetching с rate limiting (10 concurrent); cache manifest'ов на 1 час в `~/.caf/manifest-cache/`.
- **Dependencies:** PROD.

---

## 3. Theme B: UX & DX

### B-001: VSCode extension with inline provenance

- **Status:** proposed
- **Problem:** пользователь должен переключаться в терминал для `caf who`, что прерывает flow.
- **Sketch:** VSCode extension показывает provenance tooltip при hover над строкой. Inline gutter icons для файлов с active decisions.
- **Open questions:** какой minimum VSCode API version? Как синхронизировать с `.caf/sessions/` updates в real-time?
- **Dependencies:** FUTURE (IDE plugins).

### B-002: JetBrains plugin (IntelliJ, PyCharm, etc.)

- **Status:** idea
- **Problem:** тот же, что B-001, но для JetBrains IDE.
- **Sketch:** Kotlin-based plugin, тот же функционал.
- **Dependencies:** FUTURE.

### B-003: Neovim plugin

- **Status:** draft (community contribution welcome)
- **Problem:** тот же, что B-001, но для Neovim.
- **Sketch:** Lua plugin, integrates with existing LSP ecosystem. `:CafWho` command, virtual text для provenance hints.
- **Dependencies:** FUTURE.

### B-004: `caf tui` — full-screen TUI for review sessions

- **Status:** draft
- **Problem:** текущий `caf review` показывает карточки по одной в terminal; нет overview'а сессии.
- **Sketch:** full-screen TUI (на `bubbletea` или `tview`) с split-view: карточка слева, related decisions/sessions справа. Heatmap прогресса внизу.
- **Open questions:** keyboard shortcuts — vim-style или custom?
- **Dependencies:** BETA.

### B-005: Onboarding wizard

- **Status:** proposed
- **Problem:** new user не понимает, с чего начать после `caf init`.
- **Sketch:** `caf onboarding` запускает интерактивный wizard: detect agents → install adapters → run first `caf context --refresh` → suggest first `caf review`.
- **Dependencies:** ALPHA.

### B-006: Interactive decision template with autocomplete

- **Status:** draft
- **Problem:** TUI prompt для `caf decisions new` требует ручного ввода всех полей.
- **Sketch:** autocomplete для `files`, `tags`, `agent` (на основе истории решений в репо). Suggestions для `alternatives` section на основе code context.
- **Dependencies:** ALPHA.

### B-007: `caf why` with explanation mode

- **Status:** idea
- **Problem:** `caf why <file>` сейчас возвращает только decision ID; пользователь хочет полное объяснение.
- **Sketch:** `caf why <file> --explain` вызывает `caf.first_person_explain` (если LLM настроен) или формирует текстовое объяснение из session+decision данных.
- **Dependencies:** PROD (когда MCP/LLM integration появится).

### B-008: Color-coded `caf report`

- **Status:** idea
- **Problem:** отчёты в plain text; нет визуального выделения hotspots.
- **Sketch:** ANSI colors для heatmap (red = low ownership, green = high). Configurable thresholds.
- **Dependencies:** BETA.

### B-009: `caf dashboard` command (terminal dashboard)

- **Status:** idea
- **Problem:** пользователь хочет live overview метрик без запуска `caf report` каждый раз.
- **Sketch:** `caf dashboard` открывает TUI dashboard с real-time метриками (refresh каждые 30 секунд). Показывает ACS trend, blindspots count, today's review queue.
- **Dependencies:** PROD.

### B-010: Multi-repo dashboard

- **Status:** idea
- **Problem:** у автора 5+ репо, нужно видеть все метрики в одном месте.
- **Sketch:** `caf dashboard --all-repos` сканирует `~/.caf/repos.yml` и показывает aggregate metrics.
- **Dependencies:** PROD.

---

## 4. Theme C: Integrations

### C-001: GitHub Action for PR validation

- **Status:** proposed
- **Problem:** PR без decision reference проходит review без проверки на authorship.
- **Sketch:** `caf-report-action` GitHub Action. В PR status check: `caf report --check` блокирует merge если нетривиальный PR без decision или `CAF-Trivial: true`.
- **Open questions:** как определить "нетривиальный" PR в GitHub context? По LOC delta?
- **Dependencies:** PROD.

### C-002: GitLab CI integration

- **Status:** idea
- **Problem:** то же, что C-001, для GitLab.
- **Sketch:** `.gitlab-ci.yml` template с `caf-report-check` job.
- **Dependencies:** PROD.

### C-003: Slack daily summary

- **Status:** draft
- **Problem:** автор забывает запускать `caf report --day` утром.
- **Sketch:** Slack bot отправляет daily summary в личку или канал. Включает: today's review queue, yesterday's decisions, blindspots warning.
- **Dependencies:** PROD.

### C-004: Discord notification for team mode

- **Status:** idea
- **Problem:** то же, что C-003, для команд использующих Discord.
- **Dependencies:** FUTURE (multi-user mode).

### C-005: Linear / Jira integration

- **Status:** idea
- **Problem:** decisions не связаны с task tracker.
- **Sketch:** `caf decisions new --link <linear-issue-id>` добавляет ссылку. `caf why` показывает связанные задачи.
- **Dependencies:** PROD.

### C-006: GitHub PR comment bot

- **Status:** proposed
- **Problem:** в PR diff не видно provenance каждой изменённой строки.
- **Sketch:** bot комментирует PR с `caf first_person_explain` для каждого нетривиального изменённого файла.
- **Open questions:** не слишком ли шумно? Может быть opt-in через label `caf-explain`?
- **Dependencies:** PROD.

### C-007: Pre-commit.com framework integration

- **Status:** idea
- **Problem:** пользователи `pre-commit` framework не хотят отдельно настраивать `caf` hook.
- **Sketch:** опубликовать `caf-pre-commit` как pre-commit hook в marketplace. Один раз в `.pre-commit-config.yaml`, и `caf init` не нужен.
- **Dependencies:** ALPHA.

---

## 5. Theme D: AI & LLM Features

### D-001: LLM-assisted decision pre-fill

- **Status:** proposed (opt-in)
- **Problem:** заполнение decision template вручную — главное трение в workflow.
- **Sketch:** `caf decisions new --ai-prefill` использует LLM для генерации черновика на основе diff + session. Автор редактирует и подтверждает.
- **Open questions:** какой LLM? z-ai-web-dev-sdk default? Как обеспечить privacy (diff может содержать проприетарный код)?
- **Dependencies:** PROD (LLM integration в `first_person_explain` уже есть).

### D-002: Smart card generation

- **Status:** idea
- **Problem:** auto-generated карточки типов 2 и 3 иногда слишком общие.
- **Sketch:** LLM анализирует graph и формирует более конкретные вопросы. Например, вместо "Что если удалить Cache?" → "Что если удалить Cache — как это повлияет на p99 latency endpoint /v1/policies?"
- **Dependencies:** BETA.

### D-003: Decision auto-summarization

- **Status:** idea
- **Problem:** длинные decisions (>500 слов) трудно сканировать.
- **Sketch:** `caf decisions --summarize` добавляет TL;DR (1-2 предложения) в начало каждого decision. Generated lazily, cached.
- **Dependencies:** PROD.

### D-004: Smart blindspots recommendation

- **Status:** idea
- **Problem:** `caf review --blindspots` возвращает список, но не приоритезирует.
- **Sketch:** LLM приоритизирует blindspots по risk score (frequency of changes × centrality in graph × last review date).
- **Dependencies:** BETA.

### D-005: AI code review with provenance

- **Status:** proposed
- **Problem:** AI code review (через LLM) не знает о прошлых decisions и предлагает отвергнутые альтернативы.
- **Sketch:** `caf review --pr <id>` анализирует PR diff + provenance + decisions, даёт фидбек с учётом контекста.
- **Dependencies:** PROD.

### D-006: Voice-first `caf review`

- **Status:** idea
- **Problem:** на ходу (walking, commuting) неудобно отвечать на карточки текстом.
- **Sketch:** `caf review --voice` использует TTS/ASR. Карточка озвучивается, ответ语音 вводится.
- **Open questions:** ASR quality for technical terms? On-device или cloud?
- **Dependencies:** FUTURE.

### D-007: Adaptive card difficulty

- **Status:** idea
- **Problem:** SM-2 использует один `ease_factor` для всех карточек; абстрактные карточки (тип 2) труднее фактологических (тип 1).
- **Sketch:** per-card-type `ease_factor` baseline. Type 2 starts at ease=2.0 (harder), type 1 at ease=2.5.
- **Open questions:** не сломает ли это SM-2 invariants?
- **Dependencies:** BETA.

### D-008: LLM-generated first-person narratives from multiple sessions

- **Status:** idea
- **Problem:** `caf first_person_explain` сейчас фокусируется на одной session. Если файл редактировался 3 разными агентами в разное время, narrative фрагментирован.
- **Sketch:** LLM合成 multi-session narrative в единый first-person story: "Я впервые трогал этот файл в октябре, добавил базовую логику. В ноябре вернулся через Aider, добавил метрики..."
- **Dependencies:** PROD.

### D-009: Anomaly detection in authorship patterns

- **Status:** idea
- **Problem:** нет способа обнаружить, что ACS падает по конкретному модулю.
- **Sketch:** time-series anomaly detection на ACS per module. Alert если тренд вниз >3 недели.
- **Dependencies:** PROD.

---

## 6. Theme E: Languages (beyond Go)

### E-001: Python adapter

- **Status:** draft
- **Problem:** Python-репозитории не поддерживаются в graph builder.
- **Sketch:** AST-парсинг через `ast` module + `pip show` для resolved package paths. Wrapper `caf-python-adapter`.
- **Open questions:** как обрабатывать dynamic imports (`__import__`, `importlib`)? Virtualenv detection?
- **Dependencies:** FUTURE.

### E-002: TypeScript adapter

- **Status:** draft
- **Problem:** TypeScript/JavaScript репозитории не поддерживаются.
- **Sketch:** `tsc --listFiles` + tsconfig `paths` resolution. AST через `typescript` package (Node.js).
- **Open questions:** monorepo support (yarn workspaces, pnpm workspaces)?
- **Dependencies:** FUTURE.

### E-003: Rust adapter

- **Status:** idea
- **Problem:** Rust репозитории не поддерживаются.
- **Sketch:** `cargo metadata` + `cargo tree` для dependency graph. Source analysis через `rust-analyzer`.
- **Dependencies:** FUTURE.

### E-004: Java adapter

- **Status:** idea
- **Problem:** Java/Kotlin репозитории не поддерживаются.
- **Sketch:** Maven dependency tree + bytecode analysis. Или Gradle `dependencies` task.
- **Dependencies:** FUTURE.

### E-005: C# adapter

- **Status:** idea
- **Problem:** .NET репозитории не поддерживаются.
- **Sketch:** `dotnet list package` + Roslyn analyzer for call graph.
- **Dependencies:** FUTURE.

### E-006: Ruby adapter

- **Status:** idea
- **Problem:** Ruby репозитории не поддерживаются.
- **Sketch:** Bundler dependency tree + Ruby AST via `parser` gem.
- **Dependencies:** FUTURE.

### E-007: Universal adapter (LSP-based)

- **Status:** idea
- **Problem:** каждый язык требует отдельного adapter; поддержка long-tail языков не масштабируется.
- **Sketch:** universal adapter на базе LSP (Language Server Protocol). Любой язык с LSP-сервером автоматически поддерживается. Импорты — через `textDocument/documentSymbol` + custom resolution.
- **Open questions:** LSP не стандартизирует "external repo" resolution; нужно per-language configuration.
- **Dependencies:** FUTURE (после per-language adapters, как fallback).

### E-008: Polyglot repo support

- **Status:** idea
- **Problem:** репо с Go backend + TypeScript frontend не полностью покрывается одним language adapter.
- **Sketch:** multi-language adapter, который запускает несколько language-specific scanners и merge'ит результаты в единый graph.
- **Dependencies:** FUTURE (после минимум 2 language adapters).

---

## 7. Theme F: Enterprise

### F-001: SSO via SAML

- **Status:** idea
- **Problem:** enterprise пользователи требуют SSO для multi-user mode.
- **Sketch:** SAML 2.0 integration (Okta, Azure AD, etc.) для team identification.
- **Dependencies:** FUTURE (multi-user mode).

### F-002: Audit logs

- **Status:** proposed
- **Problem:** enterprise требует audit trail: кто и когда прочитал/изменил decision.
- **Sketch:** `.caf/audit.log` append-only log всех read/write operations на `.caf/`. Hash chain для tamper-evidence.
- **Open questions:** как ротейтить? Не раздувает ли репо?
- **Dependencies:** FUTURE.

### F-003: Compliance reports (SOC2, ISO 27001)

- **Status:** idea
- **Problem:** compliance officer требует отчёты в определённом формате.
- **Sketch:** `caf report --format soc2` / `--format iso27001` генерирует отчёты в нужном формате.
- **Dependencies:** FUTURE.

### F-004: On-premise LLM deployment

- **Status:** proposed
- **Problem:** enterprise не может отправлять code diff в cloud LLM для `first_person_explain`.
- **Sketch:** support for local LLM (llama.cpp, Ollama) как alternative backend. `caf config set llm.backend local-ollama`.
- **Dependencies:** PROD (LLM integration в `first_person_explain`).

### F-005: RBAC for `.caf/` operations

- **Status:** idea
- **Problem:** в multi-user mode не все пользователи должны писать decisions.
- **Sketch:** role-based access: `reader` (только read), `author` (read + write own decisions), `maintainer` (full access).
- **Dependencies:** FUTURE (multi-user mode).

### F-006: Air-gapped installation

- **Status:** idea
- **Problem:** air-gapped environments не могут скачать adapters из internet.
- **Sketch:** `caf init --offline --bundle caf-adapters-bundle.tar.gz` устанавливает все adapters из локального bundle.
- **Dependencies:** PROD.

### F-007: Enterprise support tier

- **Status:** idea (commercial)
- **Problem:** enterprise хочет SLA, support channel.
- **Sketch:** separate enterprise build with priority bug fixes, dedicated support channel. Open-source version remains free.
- **Dependencies:** FUTURE (commercial decision).

---

## 8. Theme G: Multi-user & Collaboration

### G-001: Multi-user mode (per-author attribution)

- **Status:** draft (in FUTURE roadmap)
- **Problem:** сейчас tool for one author; sync через git push/pull работает, но без слияния конфликтующих decisions.
- **Sketch:** per-author attribution в `cards.yml`: `author: <id>`. Conflict resolution rules для `ease_factor`.
- **Open questions:** чья `ease_factor` выигрывает при conflict? Минимальная (более консервативная) или последняя?
- **Dependencies:** FUTURE.

### G-002: CRDT for real-time collaboration

- **Status:** idea
- **Problem:** при concurrent edits в `cards.yml` git merge часто конфликтует на `ease_factor` полях.
- **Sketch:** CRDT-семантика (Yjs или Automerge) для `cards.yml`. Автоматическое слияние.
- **Dependencies:** FUTURE.

### G-003: Team analytics dashboard

- **Status:** idea
- **Problem:** тимлид хочет видеть aggregate ownership metrics по команде.
- **Sketch:** `caf report --team` aggregate metrics: total ACS, coverage per member, blindspots per member.
- **Open questions:** privacy — должны ли teammates видеть друг друга's recall rate?
- **Dependencies:** FUTURE (multi-user mode).

### G-004: Shared decision library

- **Status:** idea
- **Problem:** decisions часто повторяются между репозиториями одной команды.
- **Sketch:** central decision library (`~/.caf/decisions-library/`) с reusable decisions. `caf decisions new --from-library <id>` импортирует.
- **Dependencies:** FUTURE.

### G-005: Code review assignment based on ownership

- **Status:** idea
- **Problem:** PR review assignment не учитывает, кто лучше знает затронутый код.
- **Sketch:** `caf report --ownership --file <path>` возвращает heatmap владельцев. GitHub bot автоматически назначает reviewer'ов.
- **Dependencies:** FUTURE (multi-user mode).

### G-006: Mentoring mode (senior+junior shared SRS)

- **Status:** idea
- **Problem:** junior не владеет кодом, но SR-карточки генерируются из его сессий и слишком простые.
- **Sketch:** senior может "share" карточки с junior; junior видит карточки senior'а на тех же файлах.
- **Dependencies:** FUTURE.

---

## 9. Theme H: Extensibility

### H-001: Plugin API for custom card types

- **Status:** draft
- **Problem:** три встроенных типа карточек (why, what_if_remove, where_handled) могут не покрывать все сценарии.
- **Sketch:** `CardType` interface; пользователи могут регистрировать custom типы через plugins.
- **Open questions:** как distribution? Local plugins vs marketplace?
- **Dependencies:** FUTURE.

### H-002: Custom adapter registration

- **Status:** proposed
- **Problem:** пользователи proprietary AI agents (внутренние тулзы компаний) не могут интегрироваться с CAF без форка.
- **Sketch:** `caf adapters register <path-to-binary>` — external adapter, реализующий adapter contract через stdin/stdout JSON protocol.
- **Dependencies:** PROD.

### H-003: Custom report formats

- **Status:** idea
- **Problem:** `caf report` выводит в text/markdown; teams могут хотеть JSON/HTML/PDF.
- **Sketch:** `caf report --format <plugin>` plugin interface. Built-in: text, markdown, json. External plugins для других форматов.
- **Dependencies:** PROD.

### H-004: Plugin marketplace

- **Status:** idea (community-driven)
- **Problem:** нет способа распространять community plugins.
- **Sketch:** GitHub-based marketplace (как Homebrew taps). `caf plugins install <user>/<repo>`.
- **Dependencies:** FUTURE.

### H-005: Custom hook events

- **Status:** idea
- **Problem:** три встроенных события (pre-session, post-edit, pre-commit) могут не покрывать все scenarios (например, pre-merge, post-deploy).
- **Sketch:** plugin registration для custom events. Adapters могут emit'ить custom events.
- **Dependencies:** FUTURE.

### H-006: Lua scripting for custom rules

- **Status:** idea
- **Problem:** triviality rules (50 LOC) фиксированы; разные teams могут хотеть разные rules.
- **Sketch:** `.caf/rules.lua` с custom predicates. `caf hook post-edit` вызывает Lua function для определения triviality.
- **Dependencies:** PROD.

---

## 10. Theme I: Analytics & Reporting

### I-001: Trend analysis (ACS over time)

- **Status:** proposed
- **Problem:** ACS — это точечная метрика; нет понимания тренда.
- **Sketch:** `caf report --trend --metric acs --since 90d` показывает time-series. Визуализация в TUI (sparkline) или ASCII chart.
- **Dependencies:** PROD.

### I-002: Per-module breakdown

- **Status:** draft
- **Problem:** aggregate metrics скрывают problem modules.
- **Sketch:** `caf report --breakdown --by module` показывает ACS/Coverage/Recall per top-level package/directory.
- **Dependencies:** PROD.

### I-003: Agent performance comparison

- **Status:** idea
- **Problem:** нет данных, какой агент производит более "объяснимый" код.
- **Sketch:** `caf report --compare-agents` сравнивает ACS sessions, сделанных разными агентами.
- **Open questions:** корреляция ≠ causation; результаты могут вводить в заблуждение.
- **Dependencies:** PROD.

### I-004: Predictive blindspot alerting

- **Status:** idea
- **Problem:** blindspots обнаруживаются post-factum, после 30 дней без карточки.
- **Sketch:** ML model предсказывает, какие файлы скорее всего станут blindspots, на основе churn rate + centrality.
- **Dependencies:** FUTURE.

### I-005: Export to Prometheus/Grafana

- **Status:** proposed
- **Problem:** teams хотят интегрировать CAF metrics в existing observability stack.
- **Sketch:** `caf metrics --export prometheus` запускает HTTP endpoint в формате Prometheus exposition. Grafana dashboard template в `docs/grafana/`.
- **Dependencies:** PROD.

### I-006: Weekly digest email

- **Status:** idea
- **Problem:** автор не проверяет `caf report --week` регулярно.
- **Sketch:** `caf digest --weekly --email <addr>` отправляет HTML email с summary.
- **Dependencies:** PROD.

### I-007: Heatmap of authorship across files

- **Status:** draft
- **Problem:** `caf report --ownership` в текстовом виде; нет визуального heatmap.
- **Sketch:** TUI heatmap (ASCII art grid), где каждая ячейка — файл, цвет — ACS или coverage.
- **Dependencies:** PROD.

---

## 11. Theme J: Quality & Reliability

### J-001: Self-healing `.caf/` after git operations

- **Status:** proposed
- **Problem:** после `git rebase` или `git reset --hard` trailers могут потеряться, sessions могут ссылаться на несуществующие commit SHA.
- **Sketch:** `caf doctor` проверяет целостность `.caf/`, исправляет broken references, переписывает trailers при необходимости.
- **Dependencies:** ALPHA.

### J-002: Crash recovery for sessions

- **Status:** draft
- **Problem:** если `caf hook post-edit` падает посреди записи session, файл может остаться в некорректном YAML.
- **Sketch:** atomic writes через temp file + rename. Recovery on next `caf` invocation: detect partial YAML, prompt user to recover or discard.
- **Dependencies:** MVP.

### J-003: Adapter health checks

- **Status:** idea
- **Problem:** adapter может сломаться после обновления agent (например, Claude Code обновил hook format).
- **Sketch:** `caf adapters health` запускает smoke test для каждого установленного adapter'а.
- **Dependencies:** BETA.

### J-004: Telemetry (opt-in, privacy-preserving)

- **Status:** proposed
- **Problem:** maintainers не знают, какие фичи используются, какие падают.
- **Sketch:** opt-in anonymous telemetry: adapter usage, error rates, command frequency. NO code content, NO prompts, NO file paths — только aggregate counters.
- **Open questions:** как обеспечить trust? Открыть telemetry dashboard?
- **Dependencies:** PROD.

### J-005: Fuzzing for `.caf/` parser

- **Status:** idea
- **Problem:** parser может упасть на malformed YAML.
- **Sketch:** fuzzing с `go-fuzz` на YAML inputs. CI job раз в неделю.
- **Dependencies:** BETA.

### J-006: Backups for `.caf/`

- **Status:** idea
- **Problem:** accidental `git rm -r .caf/` или `git push --force` могут удалить всю ментальную карту.
- **Sketch:** `caf backup` создаёт tarball `.caf/` и загружает в configurable location (S3, local path). `caf restore` восстанавливает.
- **Dependencies:** PROD.

### J-007: Graceful degradation when LLM unavailable

- **Status:** draft
- **Problem:** `caf first_person_explain` падает, если LLM endpoint недоступен.
- **Sketch:** fallback на template-based narrative (без LLM). Quality ниже, но functional.
- **Dependencies:** PROD.

### J-008: Schema versioning and migration

- **Status:** proposed
- **Problem:** при эволюции `.caf/` schema старые репозитории могут ломаться.
- **Sketch:** `version` field в `.caf/config.yml`. `caf migrate` команда для upgrade. Backward compatibility на 1 major version.
- **Dependencies:** ALPHA (когда schema стабилизируется).

---

## 12. Theme K: UI Enhancements (NEW v0.3)

UI фичи, выходящие за рамки v0.3 TUI baseline (см. [UI.md](./UI.md)). Большинство — FUTURE.

### K-001: Multi-repo dashboard 🔥 80/20

- **Status:** proposed
- **Problem:** у автора 5+ репо, нужно видеть все метрики в одном `caf dashboard` вызове.
- **Sketch:** `caf dashboard --all-repos` сканирует `~/.caf/repos.yml` и показывает aggregate metrics в TUI. Multi-tab interface: один tab на репо + aggregate tab.
- **Open questions:** layout — вкладки или split-view? Memory budget при 50+ репо?
- **Dependencies:** PROD (после `caf sync` для multi-repo registry).

### K-002: Interactive graph visualization (web UI)

- **Status:** draft
- **Problem:** кросс-репо dependency graph неудобно смотреть в ASCII; нужны interactive zoom/pan/collapse.
- **Sketch:** web UI (read-only) с D3.js force-directed graph. Запускается `caf ui graph`, открывает localhost:7777. ONLY для visualization; editing остаётся в TUI.
- **Open questions:** security — только localhost binding? Auth?
- **Dependencies:** FUTURE (после cross-repo sync в PROD).

### K-003: VSCode extension with inline provenance 🔥 80/20

- **Status:** proposed
- **Problem:** пользователь должен переключаться в терминал для `caf who`, что прерывает flow.
- **Sketch:** VSCode extension показывает provenance tooltip при hover над строкой. Inline gutter icons для файлов с active decisions. Command palette: `CAF: Show file tour`, `CAF: Why this line`.
- **Open questions:** communication с `caf` binary — spawn process или LSP server mode?
- **Dependencies:** FUTURE.

### K-004: JetBrains plugin

- **Status:** idea
- **Problem:** тот же, что K-003, но для JetBrains IDE.
- **Sketch:** Kotlin-based plugin, тот же функционал.
- **Dependencies:** FUTURE.

### K-005: Neovim plugin

- **Status:** draft (community contribution welcome)
- **Problem:** тот же, что K-003, но для Neovim.
- **Sketch:** Lua plugin, integrates with existing LSP ecosystem. `:CafWho` command, virtual text для provenance hints.
- **Dependencies:** FUTURE.

### K-006: Voice-first `caf review`

- **Status:** idea
- **Problem:** на ходу (walking, commuting) неудобно отвечать на карточки текстом.
- **Sketch:** `caf review --voice` использует TTS/ASR. Карточка озвучивается, ответ голосом.
- **Open questions:** ASR quality for technical terms? On-device или cloud?
- **Dependencies:** FUTURE.

### K-007: Animated transitions in TUI

- **Status:** idea
- **Problem:** transitions между экранами резкие; тяжело понять, что изменилось.
- **Sketch:** плавные transitions через bubbletea transitions API (где поддерживается). 100ms max — не блокирует UX.
- **Dependencies:** PROD (после стабилизации v0.3 TUI).

### K-008: Customizable TUI themes

- **Status:** idea
- **Problem:** только 3 встроенных темы (dark/light/auto); пользователи могут хотеть custom.
- **Sketch:** `~/.caf/themes/<name>.yml` с color overrides. `caf config set ui.theme <name>`.
- **Dependencies:** BETA.

### K-009: Notification badges in dashboard

- **Status:** idea
- **Problem:** пользователь не видит, что появились новые blindspots, без явного refresh.
- **Sketch:** real-time notification badges на widgets (red dot с count).
- **Dependencies:** PROD.

### K-010: Command palette (`:` + fuzzy search)

- **Status:** draft
- **Problem:** 25+ команд CLI — пользователь не помнит все.
- **Sketch:** `:` в любом TUI screen открывает fuzzy-search command palette. Аналог VSCode `Ctrl+Shift+P`.
- **Open questions:** discovery — как пользователи узнают о `:`? Hint в status bar?
- **Dependencies:** PROD.

### K-011: Mouse support in TUI

- **Status:** idea
- **Problem:** keyboard-only — неудобно для новых пользователей, привыкших к mouse.
- **Sketch:** опциональная mouse support через bubbletea `tea.WithMouseCellMotion`. Click на widget → фокус. Scroll в списках. Не заменяет keyboard.
- **Dependencies:** BETA.

### K-012: TUI accessibility audit

- **Status:** proposed
- **Problem:** TUI может иметь accessibility issues (screen reader compat, color blind).
- **Sketch:** audit с `aria2` эквивалентом для terminal. Document в UI.md §12.3.
- **Dependencies:** PROD.

---

## 13. Rejected Items

Идеи, которые были рассмотрены и осознанно отклонены. Оставлены здесь для исторического контекста.

### R-001: Blockchain-based provenance

- **Status:** rejected
- **Problem:** (hypothetical) tamper-proof provenance.
- **Sketch:** store sessions in a blockchain.
- **Rejection reason:** git history уже даёт tamper-evidence для provenance. Blockchain добавляет complexity без реальной ценности. Слишком много overhead для single-author tool.

### R-002: Real-time collaborative editing of `.caf/` files

- **Status:** rejected
- **Problem:** (hypothetical) multiple authors editing same decision simultaneously.
- **Rejection reason:** decisions — это mostly-write-once records. Concurrent editing — анти-pattern. Git merge sufficient для редких конфликтов.

### R-003: AI agent that automatically fills decisions

- **Status:** rejected (for v0.2; may revisit as opt-in D-001)
- **Problem:** reduce friction of manual decision entry.
- **Rejection reason:** автоматически сгенерированные decisions теряют первую ценность — осознанный выбор автора. Если LLM заполняет, автор не ассимилирует. Допустимо только как opt-in pre-fill с обязательным review.

### R-004: Mobile app for `caf review`

- **Status:** rejected
- **Problem:** review on the go.
- **Rejection reason:** D-006 (voice-first review) покрывает этот use case без отдельного mobile app. Поддержка mobile app — существенный overhead без proportional value.

### R-005: Gamification (badges, streaks)

- **Status:** rejected
- **Problem:** мотивация ежедневного `caf review`.
- **Rejection reason:** противоречит first-person voice и серьёзности цели (authorship, не entertainment). Streaks создают pressure на количество, а не качество ответов.

### R-006: Public leaderboard of authorship scores

- **Status:** rejected
- **Problem:** (hypothetical) social motivation.
- **Rejection reason:** ACS — субъективная метрика. Public comparison стимулирует gaming the metric, не реальное владение кодом. Противоречит privacy-first approach.

### R-007: Automatic decision generation from git log

- **Status:** rejected
- **Problem:** backfill decisions для существующего кода.
- **Rejection reason:** git log не содержит альтернатив; LLM их выдумает. Лучше начать с пустого Decision Log и заполнять по мере новых changes. Для legacy кода — `caf tour` руководство, а не auto-generated decisions.

---

## 14. How to Promote a Backlog Item to Roadmap

Чтобы перевести item из BACKLOG в ROADMAP (MVP/ALPHA/BETA/PROD/FUTURE), необходимо:

1. **Заполнить шаблон** (см. секцию 0) полностью.
2. **Найти pilot user** — реальный человек или команда, готовые тестировать фичу.
3. **Определить success metric** — как поймём, что фича готова (см. ROADMAP.md §6.10).
4. **Создать ADR** — если фича затрагивает архитектурные решения (новая dependency, изменение data model, новый adapter pattern).
5. **Открыть GitHub Discussion** — для community feedback перед commitment.

Процесс:

- Issue в GitHub с label `backlog-promotion`.
- Maintainer review в течение 2 недель.
- При approval — item перемещается в ROADMAP.md с assigned phase.

---

## 15. How to Add a New Backlog Item

1. Создать GitHub Issue с label `backlog`.
2. Использовать шаблон из секции 0.
3. Указать theme (см. секцию 1) в поле `Theme:`.
4. Optionally предложить phase, в которую item мог бы попасть.

Maintainers периодически (раз в квартал) ревьюят backlog и:

- Promote items в ROADMAP (см. секцию 13).
- Reject items (см. секцию 12) с обоснованием.
- Merge duplicates.
- Clarify open questions через discussion.

---

## 16. Statistics

На момент v0.3:

| Theme | Items | 🔥 80/20 | Proposed | Draft | Idea |
|---|---|---|---|---|---|
| A: Performance | 6 | 2 | 0 | 2 | 4 |
| B: UX & DX | 10 | 5 | 3 | 2 | 5 |
| C: Integrations | 7 | 2 | 3 | 1 | 3 |
| D: AI & LLM | 9 | 1 | 2 | 0 | 7 |
| E: Languages | 8 | 1 | 0 | 2 | 6 |
| F: Enterprise | 7 | 0 | 2 | 0 | 5 |
| G: Multi-user | 6 | 0 | 0 | 1 | 5 |
| H: Extensibility | 6 | 1 | 1 | 1 | 4 |
| I: Analytics | 7 | 2 | 2 | 2 | 3 |
| J: Quality | 8 | 2 | 3 | 2 | 3 |
| **K: UI (NEW v0.3)** | **12** | **3** | **3** | **3** | **6** |
| **Rejected** | 7 | — | — | — | — |
| **Total** | **93** | **19** | **16** | **15** | **51** |

**🔥 80/20 items:** 19 из 93 (20%) — высокореверсивные фичи, prioritized в roadmap.

Из 86 active items: 16 proposed (готовы к promotion discussion), 15 draft (есть rough design), 51 idea (нужен design). 7 rejected (контекст сохранён).

Это означает, что backlog живой — 36% items имеют draft или proposed status, и любой contributor может взять любой item и развить до proposed.
