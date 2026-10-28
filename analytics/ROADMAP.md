# CodeAuthorityFirst — ROADMAP.md

> **Версия / Version:** 0.2
> **Формат фаз / Phase format:** MVP → ALPHA → BETA → PROD → FUTURE
> **Принцип прогресса / Progress principle:** каждая фаза самодостаточна — после завершения инструмент уже полезен, можно останавливаться. Каждая следующая фаза расширяет покрытие агентов и/или модулей.
> **Язык / Language:** bilingual — Russian body, English summary per phase
> **Аудитория / Audience:** tool maintainers + early adopters

---

## 0. Roadmap Philosophy

**EN summary.** Five phases, each self-sufficient. Phases are sequenced by *value per effort* — MVP delivers the most basic promise ("I can see which prompt produced which file") with the smallest surface area. Each subsequent phase adds a major capability (multi-agent support, decisions, mental model, MCP, multi-language, multi-user). The FUTURE phase is intentionally open-ended — items there are directional, not commitments.

Дорожная карта построена на трёх принципах:

1. **Self-sufficiency.** Каждая фаза добавляет ценность автономно. После MVP инструмент уже полезен, даже если development остановится.
2. **Coverage breadth over depth.** Сначала покрываем больше агентов (tier 1 → tier 2 → tier 3), потом углубляем фичи (multi-language, multi-user). Это даёт early feedback от разнообразных workflow.
3. **Metric-driven graduation.** Фаза считается завершённой, когда её целевые метрики достигнуты на пилотном репо (не на synthetic test suite, а на реальном коде).

---

## 1. Phase Overview

| Phase | Duration (est.) | Primary goal | Modules live | Agents supported | Languages |
|---|---|---|---|---|---|
| MVP | 2-3 weeks | Provenance for one agent | Provenance Tracker, context.md | opencode + Claude Code | Go |
| ALPHA | 4-6 weeks | Decision Log + multi-agent tier 1 | + Decision Log, hooks validation | + Cursor, Cline | Go |
| BETA | 6-8 weeks | Mental Model Reconstructor | + MMR (SM-2), graph (single-repo) | + Aider, KiloCode, Codex CLI, Continue.dev | Go |
| PROD | 8-10 weeks | MCP server + cross-repo | + MCP, cross-repo graph, OryaAgent | all 9 agents | Go |
| FUTURE | open-ended | Multi-language, multi-user, web UI | + language adapters, multi-user, dashboard | extensible | + Python, TypeScript, Rust |

---

## 2. Phase: MVP (2-3 weeks)

**EN summary.** Deliver the smallest end-to-end slice: one agent (opencode), Provenance Tracker, auto-generated `context.md`, pre-session + post-edit hooks. After MVP, the user can see which prompt produced which file, and which session made which commit. Decision Log and Mental Model are absent.

### 2.1 Scope

**Modules implemented:**

- Provenance Tracker (частично — без decision linking)
- `context.md` generator
- Adapter for **opencode** (Tier 1)
- Adapter for **Claude Code** (Tier 1) — опционально в MVP, но желательно для early feedback

**Adapters shipped:**

| Adapter | Tier | Status |
|---|---|---|
| opencode | 1 | ✓ native |
| Claude Code | 1 | ✓ native (settings.json hooks) |

**CLI commands:**

- `caf init` — создать `.caf/` с базовой структурой
- `caf init --agent opencode` — установить opencode adapter
- `caf init --agent claude-code` — установить Claude Code adapter
- `caf init --auto-detect` — автоопределение активных агентов
- `caf who <file:line>` — кто и в какой сессии написал строку
- `caf when --pattern <regex> <dir>` — когда появился паттерн
- `caf sessions --since 7d` — список сессий за период
- `caf context` — вывести `context.md` в stdout
- `caf context --refresh` — перегенерировать `context.md`
- `caf hook pre-session --repo <path>` — вызывается адаптером
- `caf hook post-edit --repo <path> --session <id>` — вызывается адаптером

**Hooks:**

- `pre-session` — читает `.caf/context.md`, инжектирует в system prompt агента
- `post-edit` — пишет `.caf/sessions/<id>.yml` после каждой итерации
- `pre-commit` — отсутствует (Decision Log ещё нет, валидировать нечего)

### 2.2 Deliverables

- Working `caf` binary для Linux/macOS (Go build)
- `.caf/` folder structure, созданная `caf init`
- Документация: quick-start для opencode и Claude Code
- Один пилотный репо (private,作者的 pet project), на котором ведётся dogfooding
- CI pipeline: unit tests для Provenance Tracker + adapter contract tests

### 2.3 Exit criteria

Фаза считается завершённой, когда:

1. Пилотный репо использует `caf` 7 дней подряд без ручных fallback'ов.
2. `caf who` корректно отвечает для ≥95% нетривиальных строк в пилотном репо.
3. `context.md` помещается в ≤2000 токенов и инжектируется в system prompt opencode без ошибок.
4. Adapter contract test suite проходит для opencode и Claude Code.

### 2.4 What's NOT in MVP

- Decision Log — приходит в ALPHA.
- Mental Model Reconstructor — приходит в BETA.
- `pre-commit` hook validation — приходит в ALPHA (нечего валидировать без Decision Log).
- MCP server — приходит в PROD.
- Любые агенты кроме opencode и Claude Code — приходят по фазам.
- Cross-repo sync — приходит в PROD.

### 2.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| opencode API unstable | medium | high | Pin to specific opencode version; contract tests catch breaking changes |
| Claude Code hook format changes | low | medium | Adapter isolated behind interface; only adapter code needs update |
| `context.md` over 2000 tokens in real repos | medium | medium | Implement priority-based truncation in MVP (principles > blindspots > decisions > exposes) |
| User adoption friction | high | medium | Single-command `caf init --auto-detect` should get a repo running in <5 minutes |

---

## 3. Phase: ALPHA (4-6 weeks)

**EN summary.** Add Decision Log with full hook validation (pre-commit blocks trivial commits without decision records), expand agent coverage to Tier 1 complete (Cursor, Cline), and ship the `caf decisions` command family. After ALPHA, every non-trivial commit must have a decision record or explicit `CAF-Trivial: true` trailer.

### 3.1 Scope

**Modules implemented:**

- Decision Log (полная функциональность)
- `pre-commit` hook с validation
- Adapter for **Cursor** (Tier 1, partial — pre-session via wrapper)
- Adapter for **Cline** (Tier 1, native)

**Adapters shipped:**

| Adapter | Tier | Status |
|---|---|---|
| opencode | 1 | ✓ native |
| Claude Code | 1 | ✓ native |
| Cursor | 1 | ✓ pre-session via wrapper, post-edit + pre-commit native |
| Cline | 1 | ✓ native (extension hooks) |

**New CLI commands:**

- `caf why <file>` — связать файл с Decision Log
- `caf decisions` — все решения за 30 дней
- `caf decisions --today` / `--file <path>` / `--search <query>` / `--gaps` / `--agent <id>`
- `caf decisions new [--files ...]` — интерактивное создание записи через TUI
- `caf hook pre-commit --repo <path>` — вызывается git hook'ом
- `caf adapters list` / `caf adapters install <id>` / `caf adapters detect`

**Hooks added:**

- `pre-commit` — валидирует наличие trailer `CAF-Session` и ссылки на decision (или `CAF-Trivial: true`)

**TUI prompt:**

- Decision template prompt в терминале (30-second budget per decision)
- Поддержка `trivial: true` shortcut
- `caf decisions new` интерактивный режим с автозаполнением из sessions

### 3.2 Deliverables

- Working Decision Log с TUI
- `pre-commit` hook (universal shell script, устанавливаемый `caf init`)
- Cursor + Cline adapters с полными capability matrix
- Wrapper script `caf-wrap cursor` / `caf-wrap cline` для pre-session fallback
- Документация: migration guide из MVP (как ретроактивно создать decisions для существующих нетривиальных коммитов)
- Pilot: расширение пилотного репо до 2-3 репо, чтобы проверить cross-repo linking через manifest.yml (но без sync automation — приходит в PROD)

### 3.3 Exit criteria

1. Decision Coverage ≥90% на пилотных репо за последние 7 дней.
2. `pre-commit` hook корректно блокирует нетривиальные коммиты без decision.
3. Cursor adapter работает на macOS + Linux без ручной конфигурации (только `caf init --agent cursor`).
4. Cline adapter работает с VSCode + Cline extension версии ≤2.0.
5. Пилот: автор может ответить на вопрос «почему мы выбрали X, а не Y» для любого нетривиального решения за последние 30 дней, используя только `caf decisions`.

### 3.4 What's NOT in ALPHA

- Mental Model Reconstructor — приходит в BETA.
- MCP server — приходит в PROD.
- Кросс-репо graph sync — приходит в PROD (manifest.yml можно писать вручную, но автоматического sync нет).
- Aider, KiloCode, Codex CLI, Continue.dev, OryaAgent — приходят в BETA/PROD.
- Поддержка языков кроме Go — приходит по фазам.

### 3.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| TUI prompt annoys users, gets disabled | high | high | `caf config set triviality_threshold_loc 75` lets users raise threshold; fail-safe means disabling doesn't break workflow |
| Cursor wrapper script breaks on Cursor updates | medium | medium | Pin Cursor version in CI; auto-detect at adapter runtime |
| Decision template too verbose, users skip with `trivial` | high | high | Track `trivial:true` rate per repo; if >30%, raise threshold via ADR-004 reconsideration |
| Pre-commit hook blocks CI commits | medium | high | Detect CI environment via env vars (`CI=true`), skip validation in CI |

---

## 4. Phase: BETA (6-8 weeks)

**EN summary.** Add Mental Model Reconstructor with SM-2 algorithm and three card types. Expand agent coverage to Tier 2 (Aider, KiloCode, Codex CLI, Continue.dev). Single-repo dependency graph enabled; cross-repo sync still manual. After BETA, the tool actively helps the author retain knowledge via spaced repetition.

### 4.1 Scope

**Modules implemented:**

- Mental Model Reconstructor (полная функциональность)
- SM-2 algorithm
- Три типа карточек: `why`, `what_if_remove`, `where_handled`
- Single-repo dependency graph (`.caf/mental-model/graph.yml`)
- Blindspots detection
- Adapter for **Aider** (Tier 2, pre-session via wrapper)
- Adapter for **KiloCode** (Tier 2, native)
- Adapter for **Codex CLI** (Tier 2, partial)
- Adapter for **Continue.dev** (Tier 2, native)

**Adapters shipped (cumulative):**

| Adapter | Tier | Status |
|---|---|---|
| opencode | 1 | ✓ |
| Claude Code | 1 | ✓ |
| Cursor | 1 | ✓ |
| Cline | 1 | ✓ |
| Aider | 2 | ✓ pre-session via wrapper |
| KiloCode | 2 | ✓ native |
| Codex CLI | 2 | ✓ pre-session partial |
| Continue.dev | 2 | ✓ native |

**New CLI commands:**

- `caf review` — ежедневная SRS-сессия (~10 минут)
- `caf review --hard` — только карточки с низким ease
- `caf review --blindspots` — файлы без SRS-карточек
- `caf review --regenerate --since 7d` — принудительная регенерация
- `caf review --stats` — статистика владения
- `caf tour <file>` — guided-экскурсия по файлу
- `caf sync` — обновить single-repo graph (cross-repo приходит в PROD)

**MMR features:**

- Автоматическая генерация карточек типа 1 (почему) при создании decision-записи
- Автоматическая генерация карточек типов 2 и 3 из graph.yml
- SM-2 с настройками по умолчанию (ease=2.5, interval=1d)
- Ежедневная сессия ~10-15 карточек с истекающим `due`
- Blindspots detection: файлы, изменённые за 30 дней, не покрытые карточками

### 4.2 Deliverables

- Working MMR с SM-2
- Single-repo graph builder (Go-only via `go list -json`)
- 4 новых адаптера (Aider, KiloCode, Codex CLI, Continue.dev)
- Wrapper scripts: `caf-wrap aider`, `caf-wrap codex`
- TUI для `caf review` (показ карточек, ввод ответа, оценка качества)
- Документация: card type taxonomy, SRS configuration, blindspots methodology
- Pilot: расширение пилотных репо до 5+ для проверки MMR на разнообразном коде

### 4.3 Exit criteria

1. Mental Model Reach ≥70% на пилотных репо.
2. Recall Rate ≥75% после 4 недель использования MMR.
3. Authorship Confidence Score ≥3.0 (целевое ≥3.5 достигается в PROD).
4. Все Tier 2 адаптеры работают на macOS + Linux.
5. Blindspots detection корректно находит файлы, не покрытые карточками, при `caf review --blindspots`.

### 4.4 What's NOT in BETA

- MCP server — приходит в PROD.
- Кросс-репо sync automation — приходит в PROD.
- OryaAgent adapter — приходит в PROD.
- Multi-language support (Python, TypeScript, Rust) — приходит в FUTURE.
- LLM-assisted card generation — приходит в FUTURE.
- Web dashboard — приходит в FUTURE.

### 4.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| SM-2 too aggressive, users burn out on reviews | high | high | Default to 10-min session cap; expose `--max-cards` flag; track session duration telemetry |
| Card generation produces low-quality questions | high | medium | Pilot with 3+ repos before BETA release; iterate on card templates |
| Tier 2 adapters have unstable APIs | medium | medium | Each adapter isolated; CI tests cover happy path only; manual smoke test before release |
| Graph builder fails on large repos (>10k files) | medium | high | Streaming graph build; lazy loading; benchmark on kubernetes/kubernetes before release |

---

## 5. Phase: PROD (8-10 weeks)

**EN summary.** Add MCP server (5 tools, 3 prompts), cross-repo sync automation, OryaAgent adapter (first-class integration). All success metrics tracked and reported. After PROD, the tool is production-ready for any Go codebase using any of the 9 supported agents.

### 5.1 Scope

**Modules implemented:**

- MCP server (`caf-mcp`) с 5 tools и 3 prompts
- Cross-repo sync automation (`caf sync` обновляет graph из зависимых репо)
- `caf.first_person_explain` через LLM
- Adapter for **OryaAgent** (Tier 3, native first-class)
- Success metrics instrumentation
- `caf report --metrics` aggregation

**Adapters shipped (cumulative — all 9):**

| Adapter | Tier | Status |
|---|---|---|
| opencode | 1 | ✓ |
| Claude Code | 1 | ✓ |
| Cursor | 1 | ✓ |
| Cline | 1 | ✓ |
| Aider | 2 | ✓ |
| KiloCode | 2 | ✓ |
| Codex CLI | 2 | ✓ |
| Continue.dev | 2 | ✓ |
| OryaAgent | 3 | ✓ native first-class |

**New CLI commands:**

- `caf report --day` / `--week` / `--ownership` / `--agent <id>` / `--metrics`
- `caf sync` — обновляет кросс-репо graph (теперь с real cross-repo support)
- `caf hook post-commit --repo <path>` — добавляет trailers в git history

**MCP tools:**

- `caf.context` — текущий контекст репо
- `caf.provenance` — история сессий, затронувших файл/строку
- `caf.decisions` — поиск decisions
- `caf.mental_model` — SRS-карточки для текущего пользователя
- `caf.first_person_explain` — флагманский tool

**MCP prompts:**

- `caf.pre_session_brief`
- `caf.code_review_with_provenance`
- `caf.archaeology`

**Cross-repo sync:**

- `caf sync` автоматически определяет локально клонированные репо через `~/.caf/repos.yml` registry
- Pullит `manifest.yml` из каждого dependee репо (через git remote или raw GitHub URL)
- Создаёт карточки типа `review-because-interface-changed` при обнаружении изменившихся интерфейсов

### 5.2 Deliverables

- MCP server runnable via `caf mcp serve` (stdio) и `caf mcp serve --transport sse` (HTTP)
- Cross-repo sync с реальными пилотными репо (5+ связанных репо)
- OryaAgent adapter (если OryaAgent публичен — иначе stub для future integration)
- Metrics dashboard в CLI: `caf report --metrics` выводит ACS / Coverage / Reach / Recall / Agent Diversity
- Документация: MCP integration guide, cross-repo setup, metrics interpretation
- 1.0 release tag

### 5.3 Exit criteria

1. Все 4 метрики успеха достигают целевых значений на пилотных репо за 8 недель:
   - ACS ≥3.5
   - Decision Coverage ≥90%
   - Mental Model Reach ≥70%
   - Recall Rate ≥75%
   - Agent Diversity H ≥1.0
2. MCP server работает с opencode, Claude Code, Cline, KiloCode (тесты на каждом клиенте).
3. Cross-repo sync корректно обнаруживает изменившиеся интерфейсы в пилотном multi-repo setup.
4. Все 9 адаптеров проходят contract tests на CI.
5. README и ARCHITECT.md согласованы с реализацией.

### 5.4 What's NOT in PROD

- Multi-user mode — FUTURE.
- Multi-language support (Python, TypeScript, Rust) — FUTURE.
- Web dashboard — FUTURE.
- LLM-assisted decision pre-fill — FUTURE.
- CRDT для real-time collaboration — FUTURE.

### 5.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| LLM hallucinates in `first_person_explain` | high | high | Strict grounding in `.caf/` data; LLM only paraphrases, never invents facts; ground-truth tests |
| MCP protocol changes break clients | medium | medium | Pin to MCP spec version; contract tests per client |
| Cross-repo sync slow on large topologies | medium | medium | Parallel fetching; cache `manifest.yml` for 1 hour |
| OryaAgent API unstable (since it's internal) | high | medium | Stub adapter if OryaAgent not yet released; document integration contract |

---

## 6. Phase: FUTURE (open-ended)

**EN summary.** Open-ended phase for items that are valuable but not committed to a specific timeline. Includes: multi-language support, multi-user mode, web dashboard, LLM-assisted features, CRDT for real-time collaboration, plugin API for custom card types. Items here may be promoted to a future numbered release or deferred indefinitely.

### 6.1 Multi-language support

- **Python adapter**: AST-парсинг импортов через `ast` module + `pip show` для resolved package paths
- **TypeScript adapter**: `tsc --listFiles` + tsconfig `paths` resolution
- **Rust adapter**: `cargo metadata` + `cargo tree` для dependency graph
- **Java adapter**: Maven/Gradle dependency tree + bytecode analysis for call graph

Каждый language adapter реализует интерфейс `LanguageAdapter` с методами `ScanImports(repoPath)`, `BuildSymbolGraph(repoPath)`, `ResolveExternalRepo(importPath)`.

### 6.2 Multi-user mode

- Per-author attribution в карточках: `cards.yml` получает поле `author: <id>`
- Conflict resolution rules для `ease_factor` при merge (брать минимальное, обоснование см. в ADR-012 future)
- Team analytics: `caf report --team` показывает aggregate metrics по команде
- SSO integration (GitHub OAuth, GitLab OAuth) для team identification

### 6.3 Web dashboard

- Read-only view на `.caf/` данные
- Interactive graph visualization (D3.js)
- SRS review в браузере (как альтернатива CLI)
- Real-time heatmap владения
- Export в PDF/Markdown для offline review

### 6.4 LLM-assisted features

- `--ai-prefill` для `caf decisions new`: LLM предлагает черновик decision из diff + session
- Smart card generation: LLM генерирует более качественные вопросы для карточек типов 2 и 3
- Decision auto-summarization: длинные decisions получают TL;DR через LLM
- All LLM features opt-in; default is LLM-free для privacy

### 6.5 CRDT for real-time collaboration

- CRDT-семантика для `cards.yml` (где `ease_factor` мутируется часто)
- Автоматическое слияние при `git pull` без ручного conflict resolution
- Library candidate: `yjs` или `automerge`

### 6.6 Plugin API

- Custom card types через plugin interface
- Custom adapters через `caf adapters register <path-to-binary>`
- Custom report formats через `caf report --format <plugin>`
- Plugin marketplace (community-driven)

### 6.7 IDE plugins

- VSCode extension: inline provenance tooltips, "Why is this line here?" hover
- JetBrains plugin: same functionality for IntelliJ-based IDEs
- Neovim plugin: Lua-based, integrates with existing LSP ecosystem

### 6.8 CI/CD integrations

- GitHub Action: `caf report --check` as PR status check (e.g., requires decision for non-trivial PRs)
- GitLab CI integration
- Slack notifications: daily `caf report --day` summary to channel

### 6.9 Enterprise features

- Audit logs (who read/wrote which decisions, when)
- Compliance reports (SOC2, ISO 27001 friendly)
- On-premise LLM deployment option for `first_person_explain`
- SSO via SAML

### 6.10 Exit criteria for FUTURE items

FUTURE items не имеют timeline, но должны иметь:

1. **Clear problem statement.** Что пользователь не может сделать сейчас, но хочет.
2. **Design sketch.** Минимальный ADR или design doc.
3. **Pilot user.** Хотя бы один реальный пользователь, готовый тестировать.
4. **Success metric.** Как поймём, что фича готова.

Без этих четырёх пунктов item остаётся в BACKLOG.

---

## 7. Cross-phase Concerns

### 7.1 Testing strategy

| Test layer | Coverage target | Phase introduced |
|---|---|---|
| Unit tests (core modules) | ≥80% | MVP |
| Adapter contract tests | all adapters | MVP (Tier 1), ALPHA (Tier 1 complete), BETA (Tier 2), PROD (Tier 3) |
| Integration tests (`.caf/` round-trip) | critical paths | ALPHA |
| E2E tests (real agent runs) | 1 per adapter, smoke test | BETA |
| Performance benchmarks | graph build <10s on 10k-file repo | BETA |
| Property-based tests (SM-2 algorithm) | invariant: ease ∈ [1.3, 3.0] | BETA |

### 7.2 Documentation strategy

| Doc | Phase | Owner |
|---|---|---|
| README.md (quick start) | MVP | maintainer |
| ARCHITECT.md | MVP, updates per phase | maintainer |
| ROADMAP.md | MVP, updates per phase | maintainer |
| BACKLOG.md | ALPHA, ongoing | maintainer |
| ADR.md | MVP, ongoing | maintainer |
| Adapter development guide | ALPHA | adapter authors |
| User guide (CLI reference) | BETA | maintainer |
| MCP integration guide | PROD | maintainer |
| Migration guides (between phases) | per phase | maintainer |

### 7.3 Release strategy

- Semantic versioning: `0.x.y` до PROD, `1.0.0` на PROD release.
- Changelog в `CHANGELOG.md`, генерируется из conventional commits.
- GitHub Releases с бинарниками для Linux/macOS/Windows.
- Migration guide для каждой мажорной версии.

### 7.4 Backwards compatibility

- `.caf/` schema evolves additively. Новые поля optional.
- Adapter interface стабилен внутри мажорной версии.
- При breaking changes в `.caf/` schema — `caf migrate` команда с auto-upgrade.

---

## 8. Phase Transitions

### 8.1 MVP → ALPHA

**Trigger:**

- MVP exit criteria met on pilot repo
- At least 1 external early adopter using `caf` for ≥2 weeks
- Adapter contract test suite stable

**Migration:**

- `caf init` для существующих MVP-репо: добавляет `decisions/` directory
- `caf migrate v0.1-to-v0.2` (если schema изменилась)

### 8.2 ALPHA → BETA

**Trigger:**

- ALPHA exit criteria met
- Decision Coverage ≥90% on 2+ pilot repos
- At least 5 external users

**Migration:**

- `caf init` добавляет `mental-model/` directory
- 第一次 запуск `caf review --regenerate` для создания карточек из существующих decisions

### 8.3 BETA → PROD

**Trigger:**

- BETA exit criteria met
- Mental Model Reach ≥70% on 3+ pilot repos
- Recall Rate ≥75% sustained for 4 weeks
- All 9 adapters shipped (или stub'ы для отсутствующих)

**Migration:**

- `caf init` добавляет MCP configuration
- `caf sync` для построения кросс-репо графа

### 8.4 PROD → FUTURE

**Trigger:**

- PROD 1.0 released
- 6 months of stable usage
- At least 50 active users
- FUTURE items prioritized via community feedback

FUTURE items могут быть промотированы в numbered release (1.1, 2.0) по мере готовности.

---

## 9. Resourcing (rough estimates)

Предполагается команда из 1-2 maintainers (part-time) + community contributors.

| Phase | Maintainer hours/week | Community PRs expected |
|---|---|---|
| MVP | 15-20 | 0-2 |
| ALPHA | 20-25 | 2-5 |
| BETA | 25-30 | 5-10 |
| PROD | 20-25 (stabilization) | 10-15 |
| FUTURE | variable | 15+ |

При full-time engagement (1 maintainer, 40h/week) сроки сжимаются примерно в 2 раза.

---

## 10. Open Questions

Вопросы, требующие решения до или во время соответствующих фаз:

1. **MVP:** Должен ли `caf init --auto-detect` предлагать все обнаруженные адаптеры или только первый? — Решение: предлагать все, с пометкой рекомендуемого.
2. **ALPHA:** Какой порог `triviality_threshold_loc` по умолчанию? — Решение: 50 (см. ADR-004), пересматривается по данным.
3. **BETA:** Какой LLM использовать для `first_person_explain` — z-ai-web-dev-sdk, локальную модель, или configurable? — Решение: configurable, с z-ai-web-dev-sdk как default.
4. **PROD:** Должен ли `caf sync` работать с приватными репо? — Решение: да, через git credential helper.
5. **FUTURE:** Должен ли multi-user mode быть backward-compatible с single-user `.caf/`? — Решение: да, через optional `author` field.

Эти вопросы будут решены через ADRs по мере приближения к соответствующим фазам.
