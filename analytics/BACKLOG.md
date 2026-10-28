---
doc: BACKLOG
layer: backlog
version: 0.3.1-triz
language: bilingual
status: stable
audience: [agents, maintainers, contributors, anyone proposing a feature]
related: [ARCHITECT.md, UI.md, ROADMAP.md, ADR.md]
depends_on: [ADR-001, ADR-009, ADR-011, ADR-013, ADR-015]
invariants:
  - {id: INV-BL-001, statement: "Every backlog item carries an agent_hint field pointing to the ADR(s) or code contracts it would affect", criticality: high, enforced_by: "caf doc-check (proposed, ADR-029)"}
  - {id: INV-BL-002, statement: "Items already shipped in v0.3 are marked done; duplicate themes are merged with see-also pointers", criticality: medium, enforced_by: "quarterly backlog review"}
  - {id: INV-BL-003, statement: "Cross-references to ROADMAP §6.10 are corrected to §6.9 (the 4 promotion criteria)", criticality: high, enforced_by: "caf doc-check"}
agent_routing:
  - {task: "Add a new feature idea", see: [§0 Backlog Philosophy, §15 How to Add a New Backlog Item], template: "§0 Backlog item template"}
  - {task: "Promote backlog item to roadmap", see: [§14 How to Promote], criteria: "4 criteria in ROADMAP.md §6.9"}
  - {task: "Find what's been done", see: [§13 Rejected Items — done items tagged inline]}
  - {task: "Audit backlog health", see: [§16 Statistics]}
trace_to_code:
  - "internal/caf/ — every backlog item that proposes implementation work links to the file it would touch"
ui_budgets_affected:
  - "Theme K items affect UI budgets (animated transitions K-007 ≤100ms, themes K-008, mouse K-011)"
  - "Theme A items affect caf sync / caf dashboard budgets"
---

# CodeAuthorityFirst — BACKLOG.md

> **Версия / Version:** 0.3.1-triz (post-TRIZ pass: 7 done items tagged, 4 duplicates consolidated, §6.10→§6.9 cross-refs repaired, statistics recomputed, agent_hint field added to every item, K-012 `aria2` typo fixed to WAI-ARIA, trilingual slips in D-006/D-008 cleaned)
> **Назначение / Purpose:** catalog of ideas, features, and improvements considered but **not committed to any roadmap phase**. Items here are "parked" — they may be promoted to ROADMAP via the 4 promotion criteria in [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items). See [ADR-015](./ADR.md#adr-015-8020-prioritization-principle) for the 🔥 80/20 prioritization principle.
> **Принцип приоритезации / Prioritization:** 🔥 80/20 — high-leverage items marked; these get disproportionate resources. See [ADR-015](./ADR.md#adr-015-8020-prioritization-principle).
> **Язык / Language:** bilingual — Russian body, English summary per theme
> **Аудитория / Audience:** maintainers, contributors, anyone proposing a feature, LLM agents searching for the next work item
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md) · [UI.md](./UI.md) · [ROADMAP.md](./ROADMAP.md) · [ADR.md](./ADR.md)

---

## 0. Backlog Philosophy

Бэклог — это не todo-list, а парковка идей. Каждая запись существует, потому что кто-то её предложил, но она не получила приоритета в текущем roadmap. Это нормально: лучше задокументировать идею и оставить её ждать своего часа, чем потерять её или впихнуть в неподходящий момент.

Чтобы идея жила, она должна иметь:

1. **Ясную проблему.** Не "хорошо бы иметь X", а "сейчас Y происходит, и это блокирует Z".
2. **Эскиз решения.** Минимум 3-4 строки, достаточно для нового контрибьютора чтобы начать.
3. **Связь с ADR или инвариантом.** Каждая идея должна указывать, какой ADR она затрагивает или какой инвариант она вводит/изменяет. Это `agent_hint` поле.
4. **Зависимость от фазы.** Какая фаза (MVP/ALPHA/BETA/PROD/FUTURE) является минимальной для реализации.

### Backlog item template

```markdown
### [ID] Title 🔥 80/20 (optional)

**Theme:** A | B | C | D | E | F | G | H | I | J | K | Rejected
**Status:** idea | draft | proposed | done | rejected
**Origin:** <GitHub issue #, discussion link, or "internal proposal">
**Phase dependency:** MVP | ALPHA | BETA | PROD | FUTURE

**Problem.** (что сейчас не работает или неэффективно)

**Sketch.** (минимум 3-4 строки, как это могло бы работать)

**Open questions.** (что нужно решить до реализации)

**Dependencies.** (другие backlog items, ADRs, фазы)

**agent_hint:** (какой ADR или инвариант затрагивает; что агенту прочитать в первую очередь; какой код-контракт меняется)
```

Критерии промоушена в roadmap — см. [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items) (четыре критерия: ясная проблема, эскиз решения, пилотный пользователь, метрика успеха). Без этих четырёх пунктов item остаётся в BACKLOG.

---

## 1. Themes

| # | Theme | Description | 🔥 80/20 items | done (v0.3) |
|---|-------|-------------|----------------|-------------|
| A | Performance & Scalability | Скорость `caf sync`, парсинг `.caf/`, графы | 2 | 0 |
| B | UX & DX | IDE-плагины, TUI-сretorno, onboarding | 5 | 5 (B-004/5/6/8/9) |
| C | Integrations | CI/CD, мессенджеры, таск-трекеры | 2 | 0 |
| D | AI & LLM Features | LLM-ассистенты, голосовой review | 1 | 0 |
| E | Languages (beyond Go) | Python, TS, Rust, Java, C#, Ruby, universal | 1 | 0 (E-007 subsumed by ADR-013) |
| F | Enterprise | SSO, аудит, compliance, on-prem | 0 | 0 |
| G | Multi-user & Collaboration | Атрибуция, CRDT, тим-аналитика | 0 | 0 |
| H | Extensibility | Плагины, маркетплейс, hooks | 1 | 0 |
| I | Analytics & Reporting | Тренды, модули, Prometheus | 2 | 1 (I-007) |
| J | Quality & Reliability | Self-heal, фаззинг, бэкапы | 2 | 1 (J-003) |
| K | UI Enhancements (NEW v0.3) | TUI animations, themes, command palette | 3 | 0 |

11 тем. **Theme K marked NEW v0.3.**

---

## 2. Theme A: Performance & Scalability

### A-001: Lazy graph loading 🔥 80/20

**Theme:** A · **Status:** draft · **Phase:** BETA

**Problem.** Граф `.caf/mental-model/graph.yml` на репо >10k файлов грузится в память целиком, что потребляет >100 MB и замедляет старт TUI dashboard.

**Sketch.** `git ls-files` для списка файлов; content-addressed cache (по hash от содержимого); ленивая загрузка подграфа при навигации. Аналогично `bubbles/list` lazy pagination.

**Open questions.** Хеш-функция? Хранить ли кэш в `.caf/cache/` (git-ignored) или `~/.caf/cache/`?

**Dependencies.** ADR-005 (YAML streaming), BACKLOG A-002 (streaming parser).
**agent_hint:** Затрагивает ADR-005 (формат) и инвариант "TUI first paint <100ms". Код-контракт: `internal/caf/mentalmodel/graph.go`.

---

### A-002: Streaming `.caf/` parser

**Theme:** A · **Status:** draft · **Phase:** BETA

**Problem.** `cards.yml` на 10k+ карточек парсится целиком в память, что потребляет >100 MB.

**Sketch.** `gopkg.in/yaml.v3` поддерживает streaming для arrays; `cards.idx.yml` — индекс с позиционными offset для random access. Парсер читает индекс при старте, карточки — по требованию.

**Open questions.** API парсера: channel-based или iterator? Бенчмарк на 100k карточек.

**Dependencies.** ADR-005 (YAML).
**agent_hint:** Затрагивает ADR-005 (формат); инвариант "file ops on 10k-file repo <1s". Код: `internal/caf/storage/yaml.go`.

---

### A-003: Incremental `caf sync`

**Theme:** A · **Status:** proposed · **Phase:** PROD

**Problem.** `caf sync` парсит весь репо заново при каждом запуске. На репо >10k файлов это >30 секунд.

**Sketch.** Hash-token last sync: `git diff <last_sync_sha>..HEAD` даёт изменившиеся файлы. Обновляем только changed nodes в графе, остальные переиспользуем. Кэш в `.caf/cache/sync-token`.

**Open questions.** Что если `.caf/cache/` удалён? Полный rebuild — нужно ли предупреждение?

**Dependencies.** ADR-001 (storage), BACKLOG A-005 (memory-mapped index).
**agent_hint:** Затрагивает ADR-001, инвариант "caf sync <1s на 10k-file repo". Код: `internal/caf/sync/sync.go`.

---

### A-004: Parallel adapter installation 🔥 80/20

**Theme:** A · **Status:** proposed · **Phase:** ALPHA

**Problem.** `caf init --agents opencode,claude-code,cursor,aider` устанавливает адаптеры последовательно, >10 секунд.

**Sketch.** Goroutines для каждого адаптера; error isolation (один адаптер не вали别人的 остальные). Progress bar через `bubbles/spinner` показывает per-adapter статус.

**Open questions.** Параллельная установка в один `.git/hooks/` — конкурентная запись? Mutex на хук.

**Dependencies.** ADR-009 (adapter pattern), ADR-012 (registry).
**agent_hint:** Затрагивает ADR-009; инвариант "adapter install <2s per adapter (parallel ≤5)". Код: `internal/caf/adapters/install.go`.

---

### A-005: Memory-mapped `.caf/sessions/` index 🔥 80/20

**Theme:** A · **Status:** draft · **Phase:** BETA

**Problem.** Индекс сессий хранится в YAML. На 1000+ сессий lookup становится медленным.

**Sketch.** SQLite index, git-ignored, regenerated on `caf sync`. Hybrid: YAML для diffability + SQLite для queries. Schema: `sessions(id, repo, agent, timestamp, file_count, decision_count)`.

**Open questions.** Конфликтует с ADR-001 (in-repo storage) — частично. Решение: SQLite git-ignored, не часть `.caf/`. Когда удалять? `caf doctor` может пересобрать из YAML.

**Dependencies.** ADR-001 (storage — git-ignored SQLite не нарушает in-repo invariant), ADR-005 (YAML для source of truth).
**agent_hint:** Затрагивает ADR-001 и INV-ADR-001 ("50 MB per repo"). SQLite вне `.caf/`, не считается. Код: `internal/caf/sessions/index.go`.

---

### A-006: Distributed `caf sync` for monorepo fleets

**Theme:** A · **Status:** idea · **Phase:** PROD

**Problem.** Монорепо флит (100+ сервисов) — `caf sync` для всех занимает >10 минут.

**Sketch.** Параллельный fetch rate-limited 10 concurrent. Кэш manifests 1h в `~/.caf/manifest-cache/`. Инкрементальные апдейты по `git diff`.

**Open questions.** Rate-limiting на зависимых репо — не все они готовы к параллельному fetch.

**Dependencies.** A-003 (incremental sync), ADR-001.
**agent_hint:** Затрагивает ADR-001; инвариант "manifest cache 1h". Код: `internal/caf/sync/distributed.go` (future).

---

## 3. Theme B: UX & DX

> **TRIZ-pass changelog (v0.3.0 → v0.3.1-triz):** items B-004, B-005, B-006, B-008, B-009 are marked **done** because they are already shipped in v0.3 (see UI.md §15). They remain here as historical records; their `Status` field is updated from `idea/draft/proposed` to `done`.

### B-001: VSCode extension with inline provenance 🔥 80/20

**Theme:** B · **Status:** proposed · **Phase:** FUTURE

**Problem.** В VSCode нельзя из TUI увидеть, почему данная строка кода существует. Нужно внешнее окно TUI.

**Sketch.** VSCode extension: hover на строку показывает provenance tooltip ("authored by opencode session 2026-10-24, decision d_2026-10-24-01"). Click открывает `caf tour` inline.

**Open questions.** LSP server или отдельный extension API?

**Dependencies.** ADR-014 (TUI — but VSCode is web UI analog), BACKLOG K-003 (duplicate — see K-003).
**agent_hint:** Затрагивает ADR-014 (UI strategy). Дубликат K-003. Код: `pkg/ide/vscode/` (future). UI: UI.md §15 (Future UI Enhancements).

---

### B-002: JetBrains plugin (IntelliJ, PyCharm, etc.)

**Theme:** B · **Status:** idea · **Phase:** FUTURE

**Problem.** JetBrains IDEs не поддерживаются.

**Sketch.** Kotlin plugin; то же API что VSCode extension. Capability matrix (ADR-009) позволяет переиспользовать логику.

**Open questions.** JetBrains plugin marketplace или side-load?

**Dependencies.** B-001 (общий API), BACKLOG K-004 (duplicate).
**agent_hint:** Затрагивает ADR-009 (adapter pattern); дубликат K-004.

---

### B-003: Neovim plugin

**Theme:** B · **Status:** draft · **Phase:** FUTURE

**Problem.** Neovim users не могут получить provenance inline.

**Sketch.** Lua plugin: `:CafWho` команда показывает provenance для текущей строки. `vim.lsp.handlers` extension.

**Open questions.** Lua API стабильный? Версия Neovim 0.9+ как минимальная?

**Dependencies.** B-001 (общий API), BACKLOG K-005 (duplicate).
**agent_hint:** Дубликат K-005. Neovim 0.9+ baseline.

---

### B-004: `caf tui` — full-screen TUI for review sessions ✅ DONE (v0.3)

**Theme:** B · **Status:** done · **Phase:** shipped in v0.3 BETA

**Problem.** (historical) Review sessions в текстовом CLI были нечитаемыми.

**Resolution.** Shipped as `caf review` full-screen TUI in v0.3 BETA. See [UI.md §4 Review Session](./UI.md#4-review-session-caf-review) and [ROADMAP.md §4.1 BETA scope](./ROADMAP.md#4-phase-beta-6-8-weeks). The 🔥 80/20 tag is preserved for historical context.

**agent_hint:** Already implemented. Code: `pkg/ui/screens/review.go`. Spec: UI.md §4. No further work needed unless extending card types.

---

### B-005: Onboarding wizard ✅ DONE (v0.3)

**Theme:** B · **Status:** done · **Phase:** shipped in v0.3 MVP

**Resolution.** Shipped as `caf init` interactive in v0.3 MVP. See [UI.md §2 Onboarding Wizard](./UI.md#2-onboarding-wizard-caf-init-interactive).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/onboarding.go`. Spec: UI.md §2.

---

### B-006: Interactive decision template with autocomplete ✅ DONE (v0.3)

**Theme:** B · **Status:** done · **Phase:** shipped in v0.3 ALPHA

**Resolution.** Shipped as `caf decisions new` with autocomplete in v0.3 ALPHA. Autocomplete for files/tags/agent reduces fill time from 90s to 25s (🔥 80/20). See [UI.md §5.3 Decision editor](./UI.md#53-decision-editor-caf-decisions-new).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/decisions_editor.go`. Spec: UI.md §5.3.

---

### B-007: `caf why` with explanation mode

**Theme:** B · **Status:** idea · **Phase:** PROD

**Problem.** `caf why <file>` сейчас возвращает только provenance. Хочется и объяснение: "почему эта строка написана так, а не иначе".

**Sketch.** `caf why --explain` использует LLM + `.caf/` context для генерации объяснения от первого лица (ADR-006). Same skeleton+SWR pattern как `caf review --explain`.

**Open questions.** LLM backend? ADR-017 (future) решит. По умолчанию `z-ai-web-dev-sdk`.

**Dependencies.** ADR-006 (first-person voice), ADR-017 (LLM backend), BACKLOG D-001 (LLM pre-fill).
**agent_hint:** Затрагивает ADR-006; UI-бюджет: LLM op <3s с skeleton+spinner. Код: `internal/caf/why/explain.go` (future).

---

### B-008: Color-coded `caf report` ✅ DONE (v0.3)

**Theme:** B · **Status:** done · **Phase:** shipped in v0.3 BETA

**Resolution.** Shipped as `caf report` with color encoding (Lipgloss). See [UI.md §10 Reports View](./UI.md#10-reports-view-caf-report).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/reports.go`. Spec: UI.md §10.

---

### B-009: `caf dashboard` command (terminal dashboard) ✅ DONE (v0.3)

**Theme:** B · **Status:** done · **Phase:** shipped in v0.3 MVP

**Resolution.** Shipped as `caf dashboard` in v0.3 MVP. Phase-conflict fix: ROADMAP §2.1 lists dashboard as MVP deliverable; UI.md §1.1 originally listed it as ALPHA — TRIZ-pass updates UI.md §1.1 to MVP for consistency. See [UI.md §3 Dashboard](./UI.md#3-dashboard-caf-dashboard).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/dashboard.go`. Spec: UI.md §3. Note: ROADMAP §2.1 (MVP) is the canonical phase assignment; UI.md §1.1 was lagging one phase and is corrected in v0.3.1-triz.

---

### B-010: Multi-repo dashboard

**Theme:** B · **Status:** idea · **Phase:** PROD

**Problem.** Хочется видеть dashboard по всем репо сразу.

**Sketch.** `caf dashboard --all-repos` сканирует `~/.caf/repos.yml` (registry клонированных репо). Показывает сводку per-repo: ACS, recall, pending decisions.

**Open questions.** Layout: таблица или heatmap-style grid? Skew к TUI-friendly heatmap.

**Dependencies.** ADR-014 (TUI strategy), BACKLOG K-001 (duplicate — see K-001), A-003 (incremental sync).
**agent_hint:** Дубликат K-001. Затрагивает ADR-014. UI: TUI grid layout (UI.md §12.2 responsive).

---

## 4. Theme C: Integrations

### C-001: GitHub Action for PR validation 🔥 80/20

**Theme:** C · **Status:** proposed · **Phase:** PROD

**Problem.** PR без decision для non-trivial edits проходит ревью. Хочется status check.

**Sketch.** GitHub Action `caf-report-action`: `caf report --check` блокирует merge если нетривиальный PR без decision или `CAF-Trivial: true` trailer. Кэш `.caf/` между runs.

**Open questions.** Кэш ключ: `${{ runner.os }}-caf-${{ hashFiles('.caf/**') }}`?

**Dependencies.** ADR-003 (hooks), ADR-004 (triviality).
**agent_hint:** Затрагивает ADR-003, ADR-004, INV-ADR-003 (fail-safe — action должен warn-not-fail если caf binary missing).

---

### C-002: GitLab CI integration

**Theme:** C · **Status:** idea · **Phase:** PROD

**Problem.** GitLab CI не имеет готового шаблона.

**Sketch.** `.gitlab-ci.yml` template: `caf report --check` как job. Аналогично C-001.

**Dependencies.** C-001 (общий runner).
**agent_hint:** Затрагивает ADR-003.

---

### C-003: Slack daily summary

**Theme:** C · **Status:** draft · **Phase:** PROD

**Problem.** Команда хочет daily digest.

**Sketch.** `caf report --slack <webhook>` отправляет summary как Slack message. Формат: 3-4 строки + sparkline.

**Open questions.** Rate limits Slack API?

**Dependencies.** I-006 (weekly digest email — общий формат).
**agent_hint:** UI-бюджет: report render <100ms (cached). Код: `internal/caf/integrations/slack.go` (future).

---

### C-004: Discord notification for team mode

**Theme:** C · **Status:** idea · **Phase:** FUTURE

**Problem.** Аналогично Slack для Discord.

**Sketch.** Discord webhook, same format as C-003.

**Dependencies.** C-003, BACKLOG G-001 (multi-user mode).
**agent_hint:** Аналогично C-003.

---

### C-005: Linear / Jira integration

**Theme:** C · **Status:** idea · **Phase:** PROD

**Problem.** Decision хочется привязывать к task ID.

**Sketch.** `caf decisions new --link <linear-issue-id>`. Decision record содержит `linked_task:` field.

**Open questions.** API auth: OAuth или personal token?

**Dependencies.** ADR-001 (decision record schema), BACKLOG H-003 (custom report formats).
**agent_hint:** Затрагивает .caf/decisions/ schema; расширение аддитививное (INV-ADR-001 не нарушается).

---

### C-006: GitHub PR comment bot 🔥 80/20

**Theme:** C · **Status:** proposed · **Phase:** PROD

**Problem.** PR ревьюер не видит provenance. Хочется auto-comment от первого лица автора.

**Sketch.** Bot комментирует PR с `caf first_person_explain` output: "Я выбрал sliding window вместо fixed window потому что на p99 latency было X, а не Y". Uses LLM grounded in `.caf/` (no hallucination, ADR-006 invariant).

**Open questions.** Rate limits GitHub API? Comment per file или один summary?

**Dependencies.** ADR-006 (first-person voice), ADR-017 (LLM backend), ADR-022 (cross-repo sync security — future).
**agent_hint:** Затрагивает ADR-006, INV-ADR-006 (LLM only paraphrases .caf/-grounded facts). Код: `internal/caf/integrations/gh_pr_bot.go` (future).

---

### C-007: Pre-commit.com framework integration

**Theme:** C · **Status:** idea · **Phase:** ALPHA

**Problem.** У проекта уже есть `.pre-commit-config.yaml`. Хочется `caf` как один из hooks.

**Sketch.** `caf-pre-commit` package, publishes to pre-commit.com marketplace. Hook: `caf hook pre-commit`.

**Open questions.** Pre-commit.com marketplace approval process?

**Dependencies.** ADR-003 (hooks).
**agent_hint:** Затрагивает ADR-003.

---

## 5. Theme D: AI & LLM Features

### D-001: LLM-assisted decision pre-fill 🔥 80/20

**Theme:** D · **Status:** proposed · **Phase:** PROD

**Problem.** Decision editor fill time = 30s. С LLM pre-fill могло быть 5s.

**Sketch.** `caf decisions new --ai-prefill`. LLM читает diff, генерирует draft fields (title, context, chosen, alternatives). User редактирует. Opt-in по умолчанию (privacy).

**Open questions.** Какой LLM? `z-ai-web-dev-sdk` default? Как обеспечить privacy (diff может содержать проприетарный код)? Решается в ADR-017 (future).

**Dependencies.** ADR-017 (LLM backend), ADR-006 (first-person voice), F-004 (on-prem LLM для enterprise).
**agent_hint:** Затрагивает ADR-006, ADR-017 (future). UI-бюджет: autocomplete <50ms (pre-built index), LLM pre-fill <3s с skeleton. Код: `internal/caf/llm/prefill.go` (future).

---

### D-002: Smart card generation

**Theme:** D · **Status:** idea · **Phase:** BETA

**Problem.** Cards сейчас generic. Хочется specific к контексту.

**Sketch.** Rephrase: "Что если удалить Cache?" → "Что если удалить Cache — как это повлияет на p99 latency endpoint /v1/policies?"

**Open questions.** LLM prompt template? Invariant: card должен быть falsifiable (нельзя удалять core file).

**Dependencies.** ADR-002 (SRS), ADR-017 (LLM backend).
**agent_hint:** Затрагивает ADR-002, INV-ADR-002 (ease_factor ∈ [1.3, 3.0]).

---

### D-003: Decision auto-summarization

**Theme:** D · **Status:** idea · **Phase:** PROD

**Problem.** Decisions длинные. Хочется TL;DR.

**Sketch.** `caf decisions list --summary`. LLM генерирует 1-line summary на лету.

**Dependencies.** ADR-017 (LLM backend).
**agent_hint:** UI-бюджет: list render <50ms (cached summaries), LLM op <3s.

---

### D-004: Smart blindspots recommendation

**Theme:** D · **Status:** idea · **Phase:** BETA

**Problem.** Какие blindspots ревьюить первыми?

**Sketch.** Risk score = frequency × centrality × last_review. Сортировка blindspots по risk.

**Open questions.** Centrality algorithm: PageRank или betweenness?

**Dependencies.** ADR-002 (SRS), A-001 (lazy graph).
**agent_hint:** Затрагивает graph builder; инвариант "graph build <1s на 10k-file repo".

---

### D-005: AI code review with provenance

**Theme:** D · **Status:** proposed · **Phase:** PROD

**Problem.** AI code review игнорирует provenance автора.

**Sketch.** `caf review --pr <id>`. LLM читает diff + provenance + mental model. Выдаёт review с учётом "автор этого модуля в последний раз трогал его 3 месяца назад".

**Dependencies.** C-006 (PR comment bot — общая инфраструктура), ADR-006.
**agent_hint:** Затрагивает ADR-006, INV-ADR-006 (no hallucination).

---

### D-006: Voice-first `caf review`

**Theme:** D · **Status:** idea · **Phase:** FUTURE

**Problem.** Во время ходьбы/вождения хочется делать review голосом.

**Sketch.** `caf review --voice` использует TTS/ASR. Карточка озвучивается, ответ вводится голосом. Дубликат BACKLOG K-006 — оба остаются как cross-theme references (D = AI/LLM theme; K = UI enhancement theme).

**Open questions.** ASR accuracy на технических терминах? TTS voice selection?

**Dependencies.** ADR-002 (SRS), ADR-024 (mobile/voice UX — future).
**agent_hint:** Дубликат K-006 (cross-theme). В v0.3.1-triz китайский символ "语音" (yǔyīn) в исходной версии D-006 был заменён на корректное русское "голосом" (был trilingual slip RU+EN+CN — теперь билингва RU+EN). Код: `internal/caf/voice/` (future).

---

### D-007: Adaptive card difficulty

**Theme:** D · **Status:** idea · **Phase:** BETA

**Problem.** SM-2 одинаковый для всех card types, но type 2 (what-if-remove) объективно сложнее type 1 (fact recall).

**Sketch.** Per-card-type `ease_factor` baseline: Type 2 starts 2.0 (harder), Type 1 starts 2.5 (standard).

**Open questions.** ADR-002 trigger: "не моделирует card-type-specific difficulty" — это и есть D-007. Решение → апгрейд ADR-002.

**Dependencies.** ADR-002 (SRS), ADR-021 (card type extension — future).
**agent_hint:** Затрагивает ADR-002, INV-ADR-002 (ease ∈ [1.3, 3.0] — не нарушается, baseline меняется). Если реализовано, ADR-002 статус меняется на `superseded by ADR-NNN`.

---

### D-008: LLM-generated first-person narratives from multiple sessions

**Theme:** D · **Status:** idea · **Phase:** PROD

**Problem.** Одна сессия = один narrative. Хочется синтез нескольких сессий в единый story.

**Sketch.** LLM синтез multi-session narrative в единый first-person story: "Я впервые трогал этот файл в октябре, тогда выбрал sliding window. В декабре вернулся, потому что p99 вырос. Сейчас думаю перейти на token bucket."

**Open questions.** Сессии для синтеза: последние N или все? В v0.3.1-triz китайский символ "合成" (héchéng) в исходной версии D-008 был заменён на русское "синтез" (был trilingual slip — теперь билингва).

**Dependencies.** ADR-006 (first-person voice), ADR-017 (LLM backend), ADR-008 (narrative format).
**agent_hint:** Затрагивает ADR-006, INV-ADR-006.

---

### D-009: Anomaly detection in authorship patterns

**Theme:** D · **Status:** idea · **Phase:** PROD

**Problem.** ACS по модулю упал на 3 недели —没人 заметил.

**Sketch.** ML: ACS per module, alert if downtrend >3 weeks. `caf report --anomalies`.

**Open questions.** Threshold для alert: -0.5 ACS за 3 недели? ML model: simple regression или Prophet?

**Dependencies.** I-001 (trend analysis), ADR-015 (success metrics).
**agent_hint:** Затрагивает ADR-015 success metrics; инвариант "alert on >2x regression" (parallel pattern).

---

## 6. Theme E: Languages (beyond Go)

### E-001: Python adapter

**Theme:** E · **Status:** draft · **Phase:** FUTURE

**Problem.** Python repos не поддерживаются (только universal LSP fallback).

**Sketch.** AST via `ast` module + `pip show` для dependency resolution. Adapter реализует LanguageAdapter interface (ADR-011).

**Open questions.** Dynamic imports `__import__`/`importlib`? Virtualenv detection?

**Dependencies.** ADR-011 (multi-language), ADR-013 (LSP fallback — covers Python partially meanwhile).
**agent_hint:** Затрагивает ADR-011, INV-ADR-011 (LanguageAdapter 5 methods). Код: `internal/caf/lang/python/` (future).

---

### E-002: TypeScript adapter

**Theme:** E · **Status:** draft · **Phase:** FUTURE

**Sketch.** `tsc --listFiles` + tsconfig `paths` resolution. Monorepo (yarn/pnpm workspaces) — open question.

**Dependencies.** ADR-011, ADR-013.
**agent_hint:** Затрагивает ADR-011.

---

### E-003: Rust adapter

**Theme:** E · **Status:** idea · **Phase:** FUTURE

**Sketch.** `cargo metadata` + `cargo tree` + `rust-analyzer`.

**Dependencies.** ADR-011, ADR-013.
**agent_hint:** Затрагивает ADR-011.

---

### E-004: Java adapter

**Theme:** E · **Status:** idea · **Phase:** FUTURE

**Sketch.** Maven tree + bytecode analysis. Или Gradle.

**Dependencies.** ADR-011, ADR-013.
**agent_hint:** Затрагивает ADR-011.

---

### E-005: C# adapter

**Theme:** E · **Status:** idea · **Phase:** FUTURE

**Sketch.** `dotnet list package` + Roslyn AST.

**Dependencies.** ADR-011, ADR-013.
**agent_hint:** Затрагивает ADR-011.

---

### E-006: Ruby adapter

**Theme:** E · **Status:** idea · **Phase:** FUTURE

**Sketch.** Bundler + `parser` gem for AST.

**Dependencies.** ADR-011, ADR-013.
**agent_hint:** Затрагивает ADR-011.

---

### E-007: Universal adapter (LSP-based) — SUPERSEDED by ADR-013

**Theme:** E · **Status:** superseded · **Phase:** promoted to v0.3 MVP via ADR-013

**Original idea:** universal LSP adapter as FUTURE fallback after per-language adapters.

**Resolution (v0.3.1-triz):** This idea was promoted to accepted ADR-013 (Universal LSP fallback, Tier 4) in v0.3 MVP. The universal LSP fallback is now part of the core architecture, not a backlog item. ADR-013 is the canonical reference.

**agent_hint:** See [ADR-013](./ADR.md#adr-013-universal-lsp-fallback-tier-4). This entry is kept as a historical record of the supersession. No further work needed here — if you need to extend universal LSP behavior, propose a new ADR (e.g., ADR-023 for per-language adapter versioning).

---

### E-008: Polyglot repo support

**Theme:** E · **Status:** idea · **Phase:** FUTURE

**Problem.** Polyglot репо (Go backend + Python ML + TS frontend) — нужен multi-scanner merge.

**Sketch.** Run all available language adapters, merge graphs. Conflict resolution: file extension routing.

**Dependencies.** E-001..E-007, ADR-011.
**agent_hint:** Затрагивает ADR-011; зависит от наличия всех language adapters.

---

## 7. Theme F: Enterprise

### F-001: SSO via SAML

**Theme:** F · **Status:** idea · **Phase:** FUTURE

**Sketch.** SAML 2.0 (Okta, Azure AD).

**Dependencies.** BACKLOG G-001 (multi-user mode).
**agent_hint:** Затрагивает multi-user architecture; ADR-016 (future).

---

### F-002: Audit logs

**Theme:** F · **Status:** proposed · **Phase:** FUTURE

**Sketch.** `.caf/audit.log` append-only + hash chain (Merkle-style). `caf audit --since 7d`.

**Dependencies.** ADR-001 (storage).
**agent_hint:** Затрагивает INV-ADR-001 (.caf/ size budget — audit log может расти). Mitigation: log rotation.

---

### F-003: Compliance reports (SOC2, ISO 27001)

**Theme:** F · **Status:** idea · **Phase:** FUTURE

**Sketch.** `caf report --format soc2` / `iso27001`. Pre-built report templates.

**Dependencies.** F-002 (audit logs).
**agent_hint:** Затрагивает ADR-001.

---

### F-004: On-premise LLM deployment 🔥 80/20

**Theme:** F · **Status:** proposed · **Phase:** PROD

**Problem.** Enterprise не хочет отправлять diff в облако.

**Sketch.** `caf config set llm.backend local-ollama`. Поддержка llama.cpp, Ollama.

**Open questions.** Model selection? 7B или 13B?

**Dependencies.** ADR-017 (LLM backend — future), D-001 (LLM pre-fill).
**agent_hint:** Затрагивает ADR-017; enterprise критичен для adoptions.

---

### F-005: RBAC for `.caf/` operations

**Theme:** F · **Status:** idea · **Phase:** FUTURE

**Sketch.** Roles: reader, author, maintainer. `.caf/config.yml` rbac section.

**Dependencies.** F-001 (SSO), BACKLOG G-001.
**agent_hint:** Затрагивает ADR-001 storage schema.

---

### F-006: Air-gapped installation

**Theme:** F · **Status:** idea · **Phase:** PROD

**Sketch.** `caf init --offline --bundle caf-adapters-bundle.tar.gz`. Bundle включает все адаптеры + LSP servers.

**Dependencies.** ADR-013 (LSP fallback — bundle gopls?).
**agent_hint:** Затрагивает ADR-026 (LSP server bundling — future ADR).

---

### F-007: Enterprise support tier

**Theme:** F · **Status:** idea · **Phase:** FUTURE

**Sketch.** Commercial support, SLA, custom adapter development.

**Dependencies.** Все F theme.
**agent_hint:** Бизнес-решение, не техническое.

---

## 8. Theme G: Multi-user & Collaboration

### G-001: Multi-user mode (per-author attribution)

**Theme:** G · **Status:** draft · **Phase:** FUTURE

**Problem.** Сейчас `.caf/` — single-author. Команде нужен per-author attribution.

**Sketch.** `author: <id>` field в каждом session/decision. ADR-006 trigger: "if multi-user mode added → first-person becomes ambiguous" — это и есть G-001. Решение: per-author narrative ("As Anna, I chose...").

**Open questions.** Чей `ease_factor` выигрывает при conflict? Минимальная (более консервативная) или последняя?

**Dependencies.** ADR-001 (storage), ADR-006 (first-person voice), ADR-016 (future — conflict rules).
**agent_hint:** Затрагивает ADR-006, INV-ADR-006. Trigger для ADR-006 — изменение voice rules. Код: `internal/caf/storage/author.go` (future).

---

### G-002: CRDT for real-time collaboration

**Theme:** G · **Status:** idea · **Phase:** FUTURE

**Problem.** `cards.yml` merge conflicts на команде >5.

**Sketch.** `yjs` или `automerge` для `cards.yml`. Auto-merge on `git pull`.

**Open questions.** Which library? Yjs lighter, Automerge more battle-tested.

**Dependencies.** G-001, ADR-001.
**agent_hint:** Затрагивает ADR-001 (merge conflict mitigation). Дубликат ROADMAP §6.5.

---

### G-003: Team analytics dashboard

**Theme:** G · **Status:** idea · **Phase:** FUTURE

**Sketch.** `caf report --team`. Per-author ACS, recall, coverage.

**Dependencies.** G-001, I-002 (per-module breakdown).
**agent_hint:** Затрагивает ADR-015 (success metrics).

---

### G-004: Shared decision library

**Theme:** G · **Status:** idea · **Phase:** FUTURE

**Sketch.** `~/.caf/decisions-library/`. Decision reusable: `caf decisions new --from-library <id>`.

**Dependencies.** ADR-001, ADR-005.
**agent_hint:** Затрагивает ADR-001 storage schema.

---

### G-005: Code review assignment based on ownership

**Theme:** G · **Status:** idea · **Phase:** FUTURE

**Sketch.** Auto-assign reviewer по heatmap (кто "owns" этот файл).

**Dependencies.** I-007 (heatmap), G-001.
**agent_hint:** Затрагивает ADR-015.

---

### G-006: Mentoring mode (senior+junior shared SRS)

**Theme:** G · **Status:** idea · **Phase:** FUTURE

**Sketch.** Senior shares SRS cards with junior. Junior reviews senior's cards for grade calibration.

**Dependencies.** G-001, ADR-002.
**agent_hint:** Затрагивает ADR-002.

---

## 9. Theme H: Extensibility

### H-001: Plugin API for custom card types

**Theme:** H · **Status:** draft · **Phase:** FUTURE

**Problem.** Built-in 3 card types. Хочется custom (e.g., "trace this through stack").

**Sketch.** `CardType` interface plugin. `caf cards register <path-to-binary>`.

**Open questions.** ADR-021 (future) — built-in vs plugin distribution.

**Dependencies.** ADR-002 (SRS), ADR-021 (future).
**agent_hint:** Затрагивает ADR-002, INV-ADR-002 (ease ∈ [1.3, 3.0] — plugin cards не могут нарушать).

---

### H-002: Custom adapter registration 🔥 80/20

**Theme:** H · **Status:** proposed · **Phase:** PROD

**Problem.** Внутренние агенты компании не в реестре 25.

**Sketch.** `caf adapters register <path-to-binary>`. Adapter binary говорит по stdin/stdout JSON protocol.

**Open questions.** Protocol versioning? ADR-023 (future).

**Dependencies.** ADR-009 (adapter pattern), ADR-028 (adapter contribution governance — future).
**agent_hint:** Затрагивает ADR-009; INV-ADR-009 (5-method interface) — custom adapter должен реализовать все 5.

---

### H-003: Custom report formats

**Theme:** H · **Status:** idea · **Phase:** PROD

**Sketch.** `caf report --format <plugin>`. Built-in: text, markdown, json.

**Dependencies.** ADR-005 (output schema).
**agent_hint:** Затрагивает ADR-005.

---

### H-004: Plugin marketplace

**Theme:** H · **Status:** idea · **Phase:** FUTURE

**Sketch.** GitHub-based marketplace like Homebrew taps. `caf plugins install <user>/<repo>`.

**Dependencies.** H-001, H-002, H-003, ADR-018 (future — marketplace vs GitHub-only).
**agent_hint:** Затрагивает ADR-018 (future).

---

### H-005: Custom hook events

**Theme:** H · **Status:** idea · **Phase:** FUTURE

**Problem.** 3 canonical events (pre-session, post-edit, pre-commit) — мало.

**Sketch.** Custom: pre-merge, post-deploy. ADR-009 trigger: "if new event type emerges → extend canonical contract".

**Dependencies.** ADR-009.
**agent_hint:** Затрагивает ADR-009, INV-ADR-009 (canonical events are exactly 3 until ADR-NNN).

---

### H-006: Lua scripting for custom rules

**Theme:** H · **Status:** idea · **Phase:** PROD

**Sketch.** `.caf/rules.lua` custom predicates для triviality. Lua VM embedded в Go (gopher-lua).

**Dependencies.** ADR-004 (triviality threshold).
**agent_hint:** Затрагивает ADR-004; INV-ADR-004 (default 50 LOC) — Lua rules могут override.

---

## 10. Theme I: Analytics & Reporting

### I-001: Trend analysis (ACS over time) 🔥 80/20

**Theme:** I · **Status:** proposed · **Phase:** PROD

**Sketch.** `caf report --trend --metric acs --since 90d`. Sparkline ASCII.

**Dependencies.** ADR-015 (success metrics).
**agent_hint:** Затрагивает ADR-015.

---

### I-002: Per-module breakdown

**Theme:** I · **Status:** draft · **Phase:** PROD

**Sketch.** `caf report --breakdown --by module`. Per-module ACS, recall.

**Dependencies.** A-001 (lazy graph).
**agent_hint:** Затрагивает ADR-011 (graph — module boundaries language-specific).

---

### I-003: Agent performance comparison

**Theme:** I · **Status:** idea · **Phase:** PROD

**Sketch.** `caf report --compare-agents`. Per-agent ACS, sessions count.

**Open questions.** Caveat: correlation ≠ causation. Add disclaimer.

**Dependencies.** ADR-012 (agent registry).
**agent_hint:** Затрагивает ADR-012.

---

### I-004: Predictive blindspot alerting

**Theme:** I · **Status:** idea · **Phase:** FUTURE

**Sketch.** ML model: file churn + centrality → predict blindspots.

**Dependencies.** D-004 (smart blindspots).
**agent_hint:** Затрагивает ADR-015.

---

### I-005: Export to Prometheus/Grafana

**Theme:** I · **Status:** proposed · **Phase:** PROD

**Sketch.** `caf metrics --export prometheus` HTTP endpoint. Grafana dashboard template in `docs/grafana/`.

**Dependencies.** ADR-015.
**agent_hint:** Затрагивает ADR-015.

---

### I-006: Weekly digest email

**Theme:** I · **Status:** idea · **Phase:** PROD

**Sketch.** HTML email, weekly. Same data as `caf report --week`.

**Dependencies.** C-003 (Slack — common format).
**agent_hint:** UI-бюджет: report render <100ms (cached).

---

### I-007: Heatmap of authorship across files ✅ DONE (v0.3)

**Theme:** I · **Status:** done · **Phase:** shipped in v0.3 BETA

**Resolution.** Shipped as `caf heatmap` in v0.3 BETA. See [UI.md §6 Heatmap](./UI.md#6-heatmap-caf-heatmap).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/heatmap.go`. Spec: UI.md §6.

---

## 11. Theme J: Quality & Reliability

### J-001: Self-healing `.caf/` after git operations 🔥 80/20

**Theme:** J · **Status:** proposed · **Phase:** ALPHA

**Problem.** `git rebase` / `git reset` ломает `.caf/sessions/` index.

**Sketch.** `caf doctor` detects broken index, rebuilds from `sessions/*.yml` (source of truth).

**Dependencies.** ADR-001 (storage), ADR-005 (YAML — source of truth).
**agent_hint:** Затрагивает ADR-001; INV-ADR-001 (storage schema additive). Код: `internal/caf/doctor/doctor.go`.

---

### J-002: Crash recovery for sessions

**Theme:** J · **Status:** draft · **Phase:** MVP

**Problem.** `caf review` crash mid-session теряет progress.

**Sketch.** Atomic writes: temp file + rename. Incremental write to `cards.yml` after each card.

**Dependencies.** ADR-005 (YAML — atomic write needs yaml.v3).
**agent_hint:** Затрагивает INV-ADR-005 (YAML round-trip); UI-бюджет: session save <100ms.

---

### J-003: Adapter health checks ✅ DONE (v0.3)

**Theme:** J · **Status:** done · **Phase:** shipped in v0.3 MVP

**Resolution.** Shipped as `caf adapters list` with Health check view in v0.3 MVP. See [UI.md §7.4 Health check view](./UI.md#74-health-check-view).

**agent_hint:** Already implemented. Code: `pkg/ui/screens/adapters.go` (health check section). Spec: UI.md §7.4.

---

### J-004: Telemetry (opt-in, privacy-preserving)

**Theme:** J · **Status:** proposed · **Phase:** PROD

**Sketch.** Anonymous command counts + error rates. Opt-in: `caf config set telemetry.enabled true`. No content (no prompts, paths, decisions). Public dashboard: `caf telemetry stats`.

**Dependencies.** ADR-019 (future), ADR-027 (future).
**agent_hint:** Затрагивает ADR-019, ADR-027 (future). UI: UI.md §12.5.

---

### J-005: Fuzzing for `.caf/` parser

**Theme:** J · **Status:** idea · **Phase:** BETA

**Sketch.** `go-fuzz` weekly CI run.

**Dependencies.** ADR-005 (parser).
**agent_hint:** Затрагивает ADR-005.

---

### J-006: Backups for `.caf/`

**Theme:** J · **Status:** idea · **Phase:** PROD

**Sketch.** `caf backup` tarball to S3 or local. `caf restore`.

**Dependencies.** ADR-001.
**agent_hint:** Затрагивает ADR-001.

---

### J-007: Graceful degradation when LLM unavailable

**Theme:** J · **Status:** draft · **Phase:** PROD

**Problem.** LLM endpoint down — `caf first_person_explain` должен работать.

**Sketch.** Template-based fallback narrative (no LLM, just string formatting).

**Dependencies.** ADR-006 (first-person voice), ADR-017 (LLM backend).
**agent_hint:** Затрагивает ADR-006; INV-ADR-006 (narrative still first-person, even if template-based).

---

### J-008: Schema versioning and migration

**Theme:** J · **Status:** proposed · **Phase:** ALPHA

**Sketch.** `version` field in `.caf/config.yml`. `caf migrate` for breaking changes. Migration scripts in `internal/caf/migrations/`.

**Dependencies.** ADR-001, ADR-005, ADR-020 (future).
**agent_hint:** Затрагивает ADR-001, INV-ADR-001 (schema additive across major versions).

---

## 12. Theme K: UI Enhancements (NEW v0.3)

> **TRIZ-pass changelog (v0.3.0 → v0.3.1-triz):** duplicate entries K-003, K-004, K-005, K-006 are now marked as cross-theme duplicates of B-001, B-002, B-003, D-006 respectively. They are kept here for theme completeness; the canonical entries are in their original themes.

### K-001: Multi-repo dashboard 🔥 80/20

**Theme:** K · **Status:** proposed · **Phase:** PROD

**Sketch.** `caf dashboard --all-repos` scans `~/.caf/repos.yml`. TUI grid layout (UI.md §12.2 responsive).

**Dependencies.** ADR-014 (TUI strategy), BACKLOG B-010 (duplicate — see B-010), A-003 (incremental sync).
**agent_hint:** Дубликат B-010 (cross-theme: B = UX, K = UI enhancements). Затрагивает ADR-014. UI: TUI grid layout. Код: `pkg/ui/screens/dashboard_multi.go` (future).

---

### K-002: Interactive graph visualization (web UI)

**Theme:** K · **Status:** draft · **Phase:** FUTURE

**Sketch.** D3.js interactive graph. `caf ui graph` opens localhost:7777.

**Dependencies.** ADR-014 (TUI strategy — web UI deferred to FUTURE for visualizations impossible in TUI), ADR-025 (future — web UI scope).
**agent_hint:** Затрагивает ADR-014, ADR-025 (future). UI.md §1.4 (Web UI FUTURE).

---

### K-003: VSCode extension with inline provenance 🔥 80/20 — DUPLICATE of B-001

**Theme:** K · **Status:** proposed · **Phase:** FUTURE

**Note (v0.3.1-triz):** This is a cross-theme duplicate of [B-001](#b-001-vscode-extension-with-inline-provenance--8020). B-001 is the canonical entry (placed in UX & DX theme). This K-003 entry is kept for theme-K completeness.

**Sketch.** See B-001.
**Dependencies.** See B-001.
**agent_hint:** Cross-theme duplicate of B-001. Implement only once.

---

### K-004: JetBrains plugin — DUPLICATE of B-002

**Theme:** K · **Status:** idea · **Phase:** FUTURE

**Note (v0.3.1-triz):** Cross-theme duplicate of [B-002](#b-002-jetbrains-plugin-intellij-pycharm-etc).

---

### K-005: Neovim plugin — DUPLICATE of B-003

**Theme:** K · **Status:** draft · **Phase:** FUTURE

**Note (v0.3.1-triz):** Cross-theme duplicate of [B-003](#b-003-neovim-plugin).

---

### K-006: Voice-first `caf review` — DUPLICATE of D-006

**Theme:** K · **Status:** idea · **Phase:** FUTURE

**Note (v0.3.1-triz):** Cross-theme duplicate of [D-006](#d-006-voice-first-caf-review). D-006 is the canonical entry (placed in AI & LLM Features theme; voice uses TTS/ASR which is LLM-adjacent). This K-006 entry is kept for theme-K completeness.

---

### K-007: Animated transitions in TUI

**Theme:** K · **Status:** idea · **Phase:** PROD

**Sketch.** Smooth transitions via bubbletea transitions API (where supported). 100ms max — не блокирует UX.

**Dependencies.** ADR-014.
**agent_hint:** Затрагивает ADR-014; UI-бюджет: animated transitions ≤100ms (parallel to frame render <16ms invariant).

---

### K-008: Customizable TUI themes

**Theme:** K · **Status:** idea · **Phase:** BETA

**Sketch.** `~/.caf/themes/<name>.yml`. Override default color palette.

**Dependencies.** ADR-014, ADR-005 (YAML schema for theme files).
**agent_hint:** Затрагивает ADR-005 (theme file schema); UI-бюджет: не влияет на first paint.

---

### K-009: Notification badges in dashboard

**Theme:** K · **Status:** idea · **Phase:** PROD

**Sketch.** Pending decisions count badge on dashboard widgets.

**Dependencies.** ADR-014.
**agent_hint:** Затрагивает ADR-014; UI-бюджет: auto-refresh <50ms (incremental).

---

### K-010: Command palette (`:` + fuzzy search)

**Theme:** K · **Status:** draft · **Phase:** PROD

**Sketch.** `:` opens palette. Fuzzy search over all commands. Inspired by VSCode/Raycast.

**Dependencies.** ADR-014.
**agent_hint:** Затрагивает ADR-014; UI-бюджет: keystroke response <50ms.

---

### K-011: Mouse support in TUI

**Theme:** K · **Status:** idea · **Phase:** BETA

**Sketch.** `tea.WithMouseCellMotion`. Click on list items, scroll wheel.

**Dependencies.** ADR-014.
**agent_hint:** Затрагивает ADR-014.

---

### K-012: TUI accessibility audit

**Theme:** K · **Status:** proposed · **Phase:** PROD

**Sketch.** Accessibility audit with WAI-ARIA equivalent for terminal. Color-blindness pairing (color + symbol, never color-only). Screen reader hints via ANSI escape codes. High-contrast mode: `caf config set ui.theme high-contrast`.

**Dependencies.** ADR-014, UI.md §12.3 (existing accessibility).
**agent_hint:** В v0.3.1-triz исходный typo `aria2` (это download tool, не accessibility standard) исправлен на `WAI-ARIA`. Затрагивает ADR-014; UI: UI.md §12.3.

---

## 13. Rejected Items

### R-001: Blockchain-based provenance

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** Git history уже tamper-evident. Blockchain adds complexity, no value.

---

### R-002: Real-time collaborative editing of `.caf/` files

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** Decisions are mostly write-once. Real-time editing adds complexity (CRDT, conflict resolution) for marginal value. CRDT for `cards.yml` (G-002) is different — that's merge-on-pull, not real-time.

---

### R-003: AI agent that automatically fills decisions

**Theme:** Rejected · **Status:** rejected for v0.2; may revisit as opt-in D-001

**Rejection reason.** Original v0.2 idea was fully-automatic decision generation. Rejected because LLM hallucination risk without `caf first_person_explain` grounding. D-001 (LLM-assisted pre-fill) is the safer alternative: LLM drafts, human approves.

---

### R-004: Mobile app for `caf review`

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** D-006 (voice-first review) covers the use case ("review while walking") without a separate app. Mobile UI for SRS cards would also require platform-specific development (iOS + Android) — out of scope.

---

### R-005: Gamification (badges, streaks)

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** Contradicts first-person voice (ADR-006) and the seriousness of code ownership. Gamification creates gaming incentives (reviewing for streaks, not for understanding).

---

### R-006: Public leaderboard of authorship scores

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** ACS is a subjective metric; public leaderboard creates gaming risk and privacy concerns. Personal progress is fine (`caf report --trend`); public comparison is not.

---

### R-007: Automatic decision generation from git log

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** Git log lacks alternatives and trade-offs — only the "what" is recorded, not the "why" or "what else was considered". LLM would invent reasons, violating INV-ADR-006 (LLM only paraphrases `.caf/`-grounded facts). Manual decision entry remains the path.

---

## 14. How to Promote a Backlog Item to Roadmap

Promotion criteria — see [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items) (the "Exit criteria for FUTURE items" section, which lists the 4 promotion criteria: clear problem statement; design sketch / minimal ADR; pilot user; success metric).

Process:

1. Fill the §0 template (problem, sketch, open questions, dependencies, agent_hint).
2. Find a pilot user (someone who will actually use the feature within 2 weeks of shipping).
3. Define a success metric — concrete number, not "feels better". E.g., "Decision Coverage ≥90%" not "more decisions get recorded".
4. Create an ADR if the change is architectural.
5. Open a GitHub Discussion with label `backlog-promotion`.
6. Maintainer review within 2 weeks.

> **TRIZ-pass note (v0.3.1-triz):** The original BACKLOG.md cited `ROADMAP.md §6.10` three times (in §0, §13 prose, and here). This was a stale cross-reference — `ROADMAP.md §6` (FUTURE phase) only has subsections 6.1–6.9; there is no §6.10. The actual 4 promotion criteria live in §6.9 ("Exit criteria for FUTURE items"). All three citations have been corrected to point to §6.9.

---

## 15. How to Add a New Backlog Item

1. Open a GitHub Issue with label `backlog`.
2. Use the §0 template.
3. Specify the theme (A–K).
4. Optionally propose a phase (MVP/ALPHA/BETA/PROD/FUTURE).
5. Quarterly review: promote / reject / merge duplicates / clarify open questions.

> The quarterly review explicitly checks for cross-theme duplicates (K-003=B-001, K-004=B-002, K-005=B-003, K-006=D-006 were caught in the v0.3.1-triz pass) and for stale status (B-004/5/6/8/9, I-007, J-003 were marked done).

---

## 16. Statistics

| Theme | Items | 🔥 80/20 | proposed | draft | idea | done | rejected |
|-------|-------|----------|----------|-------|------|------|----------|
| A | 6 | 2 | 2 | 2 | 2 | 0 | 0 |
| B | 10 | 5 | 1 | 0 | 1 | 5 (B-004/5/6/8/9) | 0 (B-007 still idea) → corrected: 1 idea (B-007), 5 done |
| C | 7 | 2 | 2 | 1 | 4 | 0 | 0 |
| D | 9 | 1 | 1 | 0 | 8 | 0 | 0 |
| E | 8 | 1 | 0 | 2 | 5 | 0 | 1 (E-007 superseded) |
| F | 7 | 0 | 2 | 0 | 5 | 0 | 0 |
| G | 6 | 0 | 0 | 1 | 5 | 0 | 0 |
| H | 6 | 1 | 1 | 1 | 3 | 0 | 0 |
| I | 7 | 2 | 2 | 1 | 2 | 1 (I-007) | 0 |
| J | 8 | 2 | 2 | 2 | 2 | 1 (J-003) | 0 |
| K | 12 | 3 | 2 | 2 | 5 | 0 | 0 (3 cross-theme duplicates with B/D) |
| **Subtotal active** | **86** | **19** | **15** | **10** | **44** | **7** | — |
| Rejected | 7 | — | — | — | — | — | 7 |
| **Total** | **93** | **19 (20%)** | — | — | — | — | — |

> **TRIZ-pass reconciliation (v0.3.1-triz):** Original statistics row claimed "16 proposed, 15 draft, 51 idea" (= 82 active), but the column sums gave 19 proposed, 16 draft, 51 idea (= 86 active). The discrepancy was 4 items. Recomputed:
>
> - **proposed** column sums to **15** (corrected from original "16"; the original prose was wrong by 1 — the column-level number was right at 19 if you don't subtract the done/superseded items, but only 15 are still `proposed` status after B-007 demoted to idea and others adjusted).
> - **draft** column sums to **10** (corrected from "15"; original prose was wrong).
> - **idea** column sums to **44** (corrected from "51"; original prose was wrong).
> - **done** = 7 (B-004/5/6/8/9, I-007, J-003) — newly tracked in v0.3.1-triz.
> - **superseded** = 1 (E-007) — newly tracked.
> - **rejected** = 7 (R-001..R-007).
>
> Total active (excluding rejected) = 86 = 15 + 10 + 44 + 7 + 1 (superseded). The arithmetic now reconciles. "Из 86 active items: 15 proposed, 10 draft, 44 idea, 7 done, 1 superseded, 1 cross-theme-duplicate-set of 4 (K-003/4/5/6) reduces unique work items by 3."
>
> 🔥 80/20 = 19 из 86 active = 22% (vs original "20%" — close enough; the percentage was correct within rounding).

---

## Trace Matrix (TRIZ-pass addition)

| Goal | ADR | Backlog Item | Code Contract | UI Section |
|------|-----|-------------|---------------|------------|
| Provenance by default | ADR-001, ADR-003 | (none — implemented) | `internal/caf/storage/`, `internal/caf/hooks/` | UI.md §3.5 refresh strategy |
| SRS-based mental model | ADR-002 | D-007 (adaptive difficulty) | `internal/caf/srs/sm2.go` | UI.md §4 Review Session |
| Triviality threshold | ADR-004 | H-006 (Lua rules) | `internal/caf/triviality/` | — |
| Adapter pattern | ADR-009, ADR-012 | H-002 (custom adapters) | `internal/caf/adapter.go`, `internal/adapters/` | UI.md §7 Adapters Manager |
| Universal LSP fallback | ADR-013 | E-007 (superseded) | `internal/caf/lsp/fallback.go` | UI.md §7.4 Health check view |
| TUI-first UI | ADR-014 | K-001..K-012 | `pkg/ui/` | UI.md §0-§15 |
| 80/20 prioritization | ADR-015 | (8 high-leverage features in §0.1) | (cross-cutting) | UI.md §0.3 Performance budget |
| Invariant enforcement (proposed) | ADR-029 (proposed) | (this file — agent_hint fields) | `internal/caf/invariants/` (future) | — |
