# CodeAuthorityFirst — ROADMAP.md

> **Версия / Version:** 0.3 (80/20 prioritized, expanded agents, UI milestones)
> **Формат фаз / Phase format:** MVP → ALPHA → BETA → PROD → FUTURE
> **Принцип приоритезации / Prioritization principle:** 🔥 80/20 — 20% усилий даёт 80% ценности; high-leverage фичи получают приоритет
> **Язык / Language:** bilingual — Russian body, English summary per phase
> **Аудитория / Audience:** tool maintainers + early adopters
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md) · [UI.md](./UI.md) · [BACKLOG.md](./BACKLOG.md) · [ADR.md](./ADR.md)

---

## 0. Roadmap Philosophy (v0.3)

**EN summary.** v0.3 roadmap applies 80/20 prioritization: high-leverage features (universal pre-commit hook, auto-detect, LSP fallback, TUI dashboard) get disproportionate resources. Perfectionism on low-value features is forbidden. Each phase is self-sufficient — after each, the tool is more useful than before.

Дорожная карта v0.3 построена на четырёх принципах:

1. **🔥 80/20 priority.** Фичи помеченные `🔥 80/20` получают приоритет в ресурсах. Это фичи, дающие 80% ценности за 20% усилий. Перфекционизм на low-value фичах запрещён (см. ADR-015).

2. **Self-sufficiency.** Каждая фаза добавляет ценность автономно. После MVP инструмент уже полезен, даже если development остановится.

3. **Universal before specific.** Сначала universal fallback (LSP, git hook, filesystem watcher) — покрывает ВСЕХ агентов сразу. Потом явные адаптеры для enhanced UX. Это обратно v0.2 подхода, где specific adapters were primary.

4. **TUI before Web.** TUI покрывает 80% UX ценностью 20% effort. Web UI отложен до FUTURE.

### 0.1 High-leverage features list (🔥 80/20)

| # | Feature | Effort | Value | Why 80/20 |
|---|---|---|---|---|
| 1 | Universal pre-commit hook | ~1 день | enforcement для 25+ агентов | один shell script, agent-agnostic |
| 2 | Auto-detect agents by config files | ~1 день | removes need for explicit `caf init --agent` | один config-driven matcher |
| 3 | Universal LSP fallback | ~5 дней | мгновенная поддержка любого agent'а | покрывает long-tail агентов |
| 4 | TUI dashboard | ~5 дней | заменяет web UI для v1 | bubbletea + lipgloss, well-trodden path |
| 5 | `caf.first_person_explain` MCP | ~3 дня | killer feature, 80% perceived value | LLM paraphrase + ground truth |
| 6 | Single Go binary | 0 (бесплатно) | zero-dependency distribution | Go gives this for free |
| 7 | ASCII heatmap | ~2 дня | works in any terminal | no rendering overhead |
| 8 | SM-2 algorithm | ~2 дня | proven, simple, no innovation risk | copied from Anki |

**Итого:** 8 фич, ~20 дней effort, ~80% ценности. Это и есть MVP scope.

---

## 1. Phase Overview (v0.3)

| Phase | Duration (est.) | Primary goal | Modules live | Agents supported | UI screens | 🔥 80/20 features |
|---|---|---|---|---|---|---|
| MVP | 2-3 weeks | Universal provenance + dashboard | Provenance, context.md, dashboard | 25+ via universal LSP + Tier 1 (2 explicit) | dashboard, onboarding, adapters list | 1, 2, 3, 6 |
| ALPHA | 4-6 weeks | Decision Log + TUI review | + Decision Log, TUI editor | + Tier 1 complete (5 explicit) | + decisions list, decisions editor | 4, 7 |
| BETA | 6-8 weeks | Mental Model Reconstructor | + MMR (SM-2), single-repo graph | + Tier 2 (8 explicit) | + review session, heatmap, file tour | 8 |
| PROD | 8-10 weeks | MCP server + cross-repo | + MCP, cross-repo sync, OryaAgent | + Tier 3 (12 explicit) | + reports view, settings | 5 |
| FUTURE | open-ended | Multi-language, web UI, plugins | + language adapters, multi-user | extensibility via plugins | + web UI (visualizations only) | — |

**🔥 Key shift from v0.2:** MVP теперь покрывает 25+ агентов через universal fallback, а не 2 явных адаптера. Это радикально расширяет initial coverage.

---

## 2. Phase: MVP (2-3 weeks)

**EN summary.** Deliver universal provenance for ALL agents (25+) via LSP fallback + 2 explicit Tier 1 adapters (opencode, Claude Code). After MVP, user can run `caf init` (zero args), universal pre-commit hook is installed, any AI agent automatically tracked. Dashboard TUI shows real-time metrics. Decision Log and Mental Model absent.

### 2.1 Scope

**🔥 80/20 MVP scope:** maximum value, minimum effort.

**Modules implemented:**

- Provenance Tracker (partial — without decision linking)
- `context.md` generator
- Universal pre-commit hook (🔥 — works for any agent)
- Universal LSP fallback adapter (🔥 — covers any agent)
- Auto-detect (`caf init --auto-detect`, 🔥)
- Adapter for **opencode** (Tier 1, native)
- Adapter for **Claude Code** (Tier 1, native)
- TUI dashboard (🔥, see UI.md §3)
- TUI onboarding wizard (see UI.md §2)
- TUI adapters manager (see UI.md §7)

**Adapters shipped:**

| Adapter | Tier | Status | Coverage mechanism |
|---|---|---|---|
| Universal LSP fallback 🔥 | 4 | ✓ | covers ANY agent |
| Universal pre-commit hook 🔥 | n/a | ✓ | covers ANY agent |
| opencode | 1 | ✓ native | enhanced UX (native pre-session injection) |
| Claude Code | 1 | ✓ native | enhanced UX |

**Coverage:** 25+ agents covered (via universal fallback) + 2 enhanced (via Tier 1 adapters).

**CLI commands (MVP):**

- `caf init` 🔥 — universal hook + LSP fallback (no args needed!)
- `caf init --auto-detect` 🔥 — universal + auto-detected explicit adapters
- `caf init --agent opencode` — universal + explicit
- `caf dashboard` 🔥 — TUI dashboard
- `caf adapters list` — TUI adapter manager
- `caf who <file:line>` — provenance lookup
- `caf when --pattern <regex> <dir>`
- `caf sessions --since 7d`
- `caf context` / `caf context --refresh`
- `caf hook pre-session --repo <path> --agent <id>`
- `caf hook post-edit --repo <path> --session <id> --agent <id>`
- `caf hook pre-commit --repo <path>` (🔥 universal)

**Hooks:**

- `pre-session` — via Tier 1 adapter (native) or wrapper (Tier 2+) or LSP (Tier 4)
- `post-edit` — via Tier 1 adapter (native) or filesystem watcher (Tier 2+)
- `pre-commit` 🔥 — universal git hook, agent-agnostic

### 2.2 Deliverables

- Working `caf` binary для Linux/macOS (Go build, single binary)
- `.caf/` folder structure, созданная `caf init`
- TUI dashboard with real-time metrics
- TUI onboarding wizard
- Universal LSP fallback (covers 25+ agents out of the box)
- Universal pre-commit git hook
- Документация: quick-start (5-min onboarding via `caf init`)
- Pilot: 1-2 pet project репозиториев на dogfooding
- CI pipeline: unit tests, snapshot tests for TUI

### 2.3 Exit criteria

1. Пилотный репо использует `caf` 7 дней подряд без ручных fallback'ов
2. Universal pre-commit hook корректно срабатывает на нетривиальные коммиты от 3+ разных агентов
3. Universal LSP fallback успешно строит symbol graph для Go репо до 1000 файлов
4. TUI dashboard first paint <100ms на пилотном репо
5. `caf init --auto-detect` корректно детектит opencode + Claude Code в репо, где они оба активны
6. Adapter contract test suite проходит для opencode, Claude Code, и universal LSP

### 2.4 What's NOT in MVP

- Decision Log — приходит в ALPHA
- Mental Model Reconstructor — приходит в BETA
- MCP server — приходит в PROD
- Tier 1 adapters кроме opencode + Claude Code — приходят в ALPHA
- Tier 2/3 explicit adapters — приходят в BETA/PROD
- Cross-repo sync automation — приходит в PROD
- LLM features (`first_person_explain`) — приходят в PROD

### 2.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| Universal LSP fallback produces noisy/incomplete graph | high | medium | Acceptable — fallback is "best effort", explicit adapters improve quality |
| Filesystem watcher consumes too much CPU | medium | medium | Throttle: poll every 1s, debounce 5s, ignore .git/ and other noisy paths |
| TUI rendering issues on weird terminals | medium | low | Comprehensive terminal compat matrix (UI.md §12.1), graceful fallback to 16-color |
| `caf init` UX confusing for new users | medium | high | Onboarding wizard with clear steps (UI.md §2); test with 5 external users before release |
| Universal pre-commit hook blocks legit commits | medium | high | Fail-safe mode: if `caf` binary missing, hook warns and lets commit through |

---

## 3. Phase: ALPHA (4-6 weeks)

**EN summary.** Add Decision Log with TUI editor (🔥 80/20 autocomplete reduces fill time from 90s to 25s). Expand Tier 1 to all 5 native agents (Cursor, Cline, KiloCode). Pre-commit validation now requires decision records (or `CAF-Trivial: true`).

### 3.1 Scope

**Modules implemented:**

- Decision Log (полная функциональность)
- TUI decisions list (`caf decisions`, see UI.md §5.1)
- TUI decision editor with autocomplete (🔥, see UI.md §5.3)
- TUI decision detail view (see UI.md §5.2)
- Pre-commit validation (decision required for non-trivial commits)
- Adapter for **Cursor** (Tier 1, partial — pre-session via wrapper)
- Adapter for **Cline** (Tier 1, native)
- Adapter for **KiloCode** (Tier 1, native)

**Adapters shipped (cumulative):**

| Adapter | Tier | Status |
|---|---|---|
| Universal LSP fallback 🔥 | 4 | ✓ |
| Universal pre-commit hook 🔥 | n/a | ✓ |
| opencode | 1 | ✓ |
| Claude Code | 1 | ✓ |
| Cursor | 1 | ✓ pre-session wrapper, post-edit + pre-commit native |
| Cline | 1 | ✓ native |
| KiloCode | 1 | ✓ native |

**Coverage:** 25+ agents via universal + 5 enhanced via Tier 1.

**New CLI commands:**

- `caf why <file>` — связать файл с Decision Log
- `caf decisions` 🔥 — TUI dashboard для decisions
- `caf decisions new [--files ...]` 🔥 — TUI editor с autocomplete
- `caf decisions --today` / `--file <path>` / `--search <query>` / `--gaps` / `--agent <id>`

**TUI screens added:**

- Decisions list (UI.md §5.1)
- Decision detail (UI.md §5.2)
- Decision editor (UI.md §5.3)
- Settings editor (UI.md §11)

### 3.2 Deliverables

- Working Decision Log с TUI editor (🔥 autocomplete)
- Universal pre-commit validation (decision required for non-trivial commits)
- Cursor + Cline + KiloCode adapters с полными capability matrix
- Wrapper scripts `caf-wrap cursor` / `caf-wrap cline` / `caf-wrap kilocode` для pre-session fallback
- Документация: migration guide из MVP (как ретроактивно создать decisions)
- Pilot: расширение до 3-5 репо

### 3.3 Exit criteria

1. Decision Coverage ≥90% на пилотных репо за последние 7 дней
2. Universal pre-commit hook корректно блокирует нетривиальные коммиты без decision (works for ALL 25+ agents)
3. TUI decision editor average fill time ≤30 seconds (vs 90s without autocomplete)
4. Все Tier 1 адаптеры работают на macOS + Linux без ручной конфигурации
5. Пилот: автор может ответить на вопрос «почему мы выбрали X, а не Y» для любого нетривиального решения за последние 30 дней через `caf decisions`

### 3.4 What's NOT in ALPHA

- Mental Model Reconstructor — приходит в BETA
- MCP server — приходит в PROD
- Кросс-репо graph sync automation — приходит в PROD
- Tier 2 explicit adapters (Aider, Codex CLI, etc.) — приходят в BETA
- Поддержка языков кроме Go — приходит по фазам

### 3.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| TUI editor annoys users with autocomplete noise | medium | medium | Configurable: `caf config set decisions.autocomplete aggressive|conservative|off` |
| Pre-commit hook blocks CI commits | high | high | Detect `CI=true` env var, skip validation in CI environments |
| Decision template too verbose | high | high | Default to compact template; `caf decisions new --verbose` for full |
| Cursor wrapper script breaks on Cursor updates | medium | medium | Pin Cursor version in CI; auto-detect at runtime |
| Tier 1 adapters have inconsistent UX | medium | low | Standardize via TUI screens — all adapters use same `caf decisions new` editor |

---

## 4. Phase: BETA (6-8 weeks)

**EN summary.** Add Mental Model Reconstructor with SM-2 algorithm. TUI review session (🔥 full-screen, split-view). Expand to all Tier 2 adapters (8 agents: Aider, Codex CLI, Continue.dev, GitHub Copilot CLI, Gemini CLI, Amazon Q, Goose, Windsurf). Single-repo dependency graph enabled.

### 4.1 Scope

**Modules implemented:**

- Mental Model Reconstructor (полная функциональность)
- SM-2 algorithm (🔥 — proven, simple)
- Три типа карточек: `why`, `what_if_remove`, `where_handled`
- Single-repo dependency graph (`.caf/mental-model/graph.yml`)
- Blindspots detection
- TUI review session (🔥, see UI.md §4) — full-screen, split-view, keyboard shortcuts
- TUI heatmap (see UI.md §6)
- TUI file tour (see UI.md §8)
- Adapter for **Aider** (Tier 2, pre-session via wrapper)
- Adapter for **Codex CLI** (Tier 2, partial)
- Adapter for **Continue.dev** (Tier 2, native)
- Adapter for **GitHub Copilot CLI** (Tier 2, wrapper)
- Adapter for **Gemini CLI** (Tier 2, wrapper)
- Adapter for **Amazon Q Developer CLI** (Tier 2, wrapper)
- Adapter for **Goose** (Tier 2, native)
- Adapter for **Windsurf** (Tier 2, partial)

**Adapters shipped (cumulative):**

| Adapter | Tier | Status |
|---|---|---|
| Universal LSP fallback 🔥 | 4 | ✓ |
| Universal pre-commit hook 🔥 | n/a | ✓ |
| opencode, Claude Code, Cursor, Cline, KiloCode | 1 | ✓ |
| Aider, Codex CLI, Continue.dev, GH Copilot, Gemini, Amazon Q, Goose, Windsurf | 2 | ✓ |
| **Total explicit:** 13 + universal | | |

**Coverage:** 25+ agents via universal + 13 enhanced.

**New CLI commands:**

- `caf review` 🔥 — full-screen TUI SRS session
- `caf review --hard` / `--blindspots` / `--regenerate --since 7d` / `--stats`
- `caf heatmap` — ASCII heatmap владения
- `caf tour <file>` — TUI guided file tour
- `caf sync` — обновить single-repo graph (cross-repo в PROD)

**TUI screens added:**

- Review session (UI.md §4)
- Heatmap (UI.md §6)
- File tour (UI.md §8)

### 4.2 Deliverables

- Working MMR с SM-2
- TUI review session (🔥 full-screen, split-view)
- Single-repo graph builder (Go-only via `go list -json`)
- 8 новых Tier 2 адаптеров
- Wrapper scripts: `caf-wrap aider`, `caf-wrap codex`, etc.
- Документация: card type taxonomy, SRS configuration, blindspots methodology
- Pilot: расширение до 5-10 репо

### 4.3 Exit criteria

1. Mental Model Reach ≥70% на пилотных репо
2. Recall Rate ≥75% после 4 недель использования MMR
3. Authorship Confidence Score ≥3.0 (целевое ≥3.5 достигается в PROD)
4. TUI review session completion rate ≥80% (users finish sessions, not abandon)
5. Все Tier 2 адаптеры работают на macOS + Linux
6. Heatmap render <200ms на репо до 1000 файлов

### 4.4 What's NOT in BETA

- MCP server — приходит в PROD
- Кросс-репо sync automation — приходит в PROD
- Tier 3 adapters (Tabnine, Cody, Trae, etc.) — приходят в PROD
- Multi-language support — приходит в FUTURE
- LLM-assisted card generation — приходит в FUTURE
- Web dashboard — приходит в FUTURE

### 4.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| SM-2 too aggressive, users burn out | high | high | Default cap 15 cards/session; `caf config set review.max_cards 10`; track session duration telemetry |
| TUI review UX feels cramped on small terminals | medium | medium | Responsive design (UI.md §12.2); minimum 60 cols; compact mode |
| Card generation produces low-quality questions | high | medium | Pilot with 3+ repos before BETA release; iterate on card templates |
| Tier 2 wrappers fail on agent updates | medium | medium | Each adapter isolated; CI tests cover happy path; manual smoke test before release |
| Graph builder fails on large repos (>10k files) | medium | high | Streaming graph build; lazy loading; benchmark on kubernetes/kubernetes before release |

---

## 5. Phase: PROD (8-10 weeks)

**EN summary.** Add MCP server with 5 tools + 3 prompts (🔥 `caf.first_person_explain` is the killer feature). Cross-repo sync automation. All Tier 3 adapters (12 agents). All success metrics tracked. 1.0 release.

### 5.1 Scope

**Modules implemented:**

- MCP server (`caf-mcp`) с 5 tools и 3 prompts
- `caf.first_person_explain` 🔥 через LLM (killer feature)
- Cross-repo sync automation (`caf sync` обновляет graph из зависимых репо)
- TUI reports view (`caf report`, see UI.md §10)
- Success metrics instrumentation
- `caf report --metrics` aggregation
- Tier 3 adapters (12 agents): Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Augment, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep, Void, OryaAgent

**Adapters shipped (cumulative):**

| Adapter | Tier | Status |
|---|---|---|
| Universal LSP fallback 🔥 | 4 | ✓ |
| Universal pre-commit hook 🔥 | n/a | ✓ |
| Tier 1 (5 agents) | 1 | ✓ |
| Tier 2 (8 agents) | 2 | ✓ |
| Tier 3 (12 agents) | 3 | ✓ (wrapper-only) |
| **Total explicit:** 25 + universal | | |

**Coverage:** 25+ agents via universal + 25 enhanced (all explicit).

**New CLI commands:**

- `caf report --day` / `--week` / `--ownership` / `--agent <id>` / `--metrics`
- `caf sync` — обновляет кросс-репо graph (real cross-repo support)
- `caf hook post-commit --repo <path>` — добавляет trailers в git history

**MCP tools:**

- `caf.context` — текущий контекст репо
- `caf.provenance` — история сессий, затронувших файл/строку
- `caf.decisions` — поиск decisions
- `caf.mental_model` — SRS-карточки
- `caf.first_person_explain` 🔥 — флагманский tool (first-person narrative)

**MCP prompts:**

- `caf.pre_session_brief`
- `caf.code_review_with_provenance`
- `caf.archaeology`

**Cross-repo sync:**

- `caf sync` автоматически определяет локально клонированные репо через `~/.caf/repos.yml` registry
- Pullит `manifest.yml` из каждого dependee репо
- Создаёт карточки типа `review-because-interface-changed`

### 5.2 Deliverables

- MCP server runnable via `caf mcp serve` (stdio) и `caf mcp serve --transport sse` (HTTP)
- Cross-repo sync с реальными пилотными репо (5+ связанных репо)
- All 25 explicit adapters + universal LSP
- TUI reports view with metrics
- 1.0 release tag

### 5.3 Exit criteria

1. Все 5 метрик успеха достигают целевых значений на пилотных репо за 8 недель:
   - ACS ≥3.5
   - Decision Coverage ≥90%
   - Mental Model Reach ≥70%
   - Recall Rate ≥75%
   - Agent Diversity H ≥1.0
   - Adapter Coverage ≥70% native (NEW v0.3)
2. MCP server работает с opencode, Claude Code, Cline, KiloCode, Goose (тесты на каждом клиенте)
3. Cross-repo sync корректно обнаруживает изменившиеся интерфейсы в пилотном multi-repo setup
4. Все 25 explicit adapters проходят contract tests на CI
5. `caf.first_person_explain` produces factually accurate narratives (no LLM hallucination on `.caf/`-grounded facts)
6. README, ARCHITECT.md, UI.md, ROADMAP.md, BACKLOG.md, ADR.md — все согласованы с реализацией

### 5.4 What's NOT in PROD

- Multi-user mode — FUTURE
- Multi-language support (Python, TypeScript, Rust) — FUTURE
- Web dashboard — FUTURE
- LLM-assisted decision pre-fill — FUTURE
- CRDT для real-time collaboration — FUTURE

### 5.5 Risks

| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| LLM hallucinates in `first_person_explain` | high | high | Strict grounding in `.caf/` data; LLM only paraphrases, never invents facts; ground-truth tests |
| MCP protocol changes break clients | medium | medium | Pin to MCP spec version; contract tests per client |
| Cross-repo sync slow on large topologies | medium | medium | Parallel fetching; cache `manifest.yml` for 1 hour |
| Tier 3 wrapper scripts brittle across 12 agents | high | medium | Each adapter isolated; CI smoke tests; community contributions welcome |
| TUI performance degrades on huge repos | medium | high | Performance budgets enforced (UI.md §0.3); benchmarks in CI; alert on >2x regression |

---

## 6. Phase: FUTURE (open-ended)

**EN summary.** Open-ended phase. Items here are directional, not committed. Includes: multi-language support, multi-user mode, web UI (visualizations only), LLM-assisted features, CRDT, plugin API, IDE plugins. Items may be promoted to numbered releases (1.1, 2.0) as they mature.

### 6.1 Multi-language support

- **Python adapter**: AST-парсинг через `ast` module + `pip show`
- **TypeScript adapter**: `tsc --listFiles` + tsconfig `paths`
- **Rust adapter**: `cargo metadata` + `cargo tree`
- **Java adapter**: Maven/Gradle dependency tree
- **Universal LSP enhancement**: extend Tier 4 fallback to handle more languages automatically

### 6.2 Multi-user mode

- Per-author attribution в карточках
- Conflict resolution rules для `ease_factor`
- Team analytics: `caf report --team`
- SSO integration (GitHub OAuth, GitLab OAuth)

### 6.3 Web dashboard (visualizations only)

- Read-only view на `.caf/` данные
- **Interactive graph visualization (D3.js)** — единственное, что нельзя в TUI
- Real-time heatmap владения
- Export в PDF/Markdown для offline review
- **NOT replacing TUI** — TUI remains primary

### 6.4 LLM-assisted features

- `--ai-prefill` для `caf decisions new`: LLM предлагает черновик decision
- Smart card generation: LLM генерирует более качественные вопросы
- Decision auto-summarization
- All LLM features opt-in; default is LLM-free для privacy

### 6.5 CRDT for real-time collaboration

- CRDT-семантика для `cards.yml`
- Автоматическое слияние при `git pull` без ручного conflict resolution
- Library candidate: `yjs` или `automerge`

### 6.6 Plugin API

- Custom card types через plugin interface
- Custom adapters через `caf adapters register <path-to-binary>`
- Custom report formats
- Plugin marketplace (community-driven)

### 6.7 IDE plugins

- VSCode extension: inline provenance tooltips, "Why is this line here?" hover
- JetBrains plugin
- Neovim plugin

### 6.8 CI/CD integrations

- GitHub Action: `caf report --check` as PR status check
- GitLab CI integration
- Slack notifications: daily `caf report --day` summary

### 6.9 Exit criteria for FUTURE items

FUTURE items не имеют timeline, но должны иметь:

1. Clear problem statement
2. Design sketch (минимальный ADR)
3. Pilot user
4. Success metric

Без этих четырёх пунктов item остаётся в BACKLOG.

---

## 7. Cross-phase Concerns

### 7.1 Testing strategy

| Test layer | Coverage target | Phase |
|---|---|---|
| Unit tests (core modules) | ≥80% | MVP |
| TUI snapshot tests (UI.md §14.4) | all screens | MVP |
| Adapter contract tests | all 25 adapters | MVP (universal + 2), ALPHA (+3), BETA (+8), PROD (+12) |
| Integration tests (`.caf/` round-trip) | critical paths | ALPHA |
| E2E tests (real agent runs) | 1 per adapter, smoke test | BETA |
| Performance benchmarks (UI.md §0.3) | all budgets enforced | MVP |
| Property-based tests (SM-2 algorithm) | invariant: ease ∈ [1.3, 3.0] | BETA |
| Terminal compat tests | 8 terminals (UI.md §12.1) | BETA |

### 7.2 Documentation strategy

| Doc | Phase | Owner |
|---|---|---|
| README.md (quick start) | MVP | maintainer |
| ARCHITECT.md | MVP, updates per phase | maintainer |
| UI.md | MVP (TUI screens), updates per phase | maintainer |
| ROADMAP.md | MVP, updates per phase | maintainer |
| BACKLOG.md | ALPHA, ongoing | maintainer |
| ADR.md | MVP, ongoing | maintainer |
| Adapter development guide | ALPHA | adapter authors |
| User guide (CLI reference) | BETA | maintainer |
| MCP integration guide | PROD | maintainer |
| Migration guides (between phases) | per phase | maintainer |

### 7.3 Release strategy

- Semantic versioning: `0.x.y` до PROD, `1.0.0` на PROD release
- Changelog в `CHANGELOG.md`, генерируется из conventional commits
- GitHub Releases с бинарниками для Linux/macOS/Windows
- Migration guide для каждой мажорной версии

### 7.4 Backwards compatibility

- `.caf/` schema evolves additively. Новые поля optional.
- Adapter interface стабилен внутри мажорной версии.
- При breaking changes — `caf migrate` команда с auto-upgrade.

---

## 8. Phase Transitions

### 8.1 MVP → ALPHA

**Trigger:**

- MVP exit criteria met on pilot repo
- Universal fallback works for ≥3 different agents (verified empirically)
- At least 1 external early adopter using `caf` for ≥2 weeks
- Adapter contract test suite stable

**Migration:**

- `caf init` для существующих MVP-репо: добавляет `decisions/` directory
- `caf migrate v0.2-to-v0.3` (если schema изменилась)

### 8.2 ALPHA → BETA

**Trigger:**

- ALPHA exit criteria met
- Decision Coverage ≥90% on 2+ pilot repos
- TUI decision editor avg fill time ≤30s
- At least 5 external users

**Migration:**

- `caf init` добавляет `mental-model/` directory
- Первый запуск `caf review --regenerate` для создания карточек из существующих decisions

### 8.3 BETA → PROD

**Trigger:**

- BETA exit criteria met
- Mental Model Reach ≥70% on 3+ pilot repos
- Recall Rate ≥75% sustained for 4 weeks
- All 25 adapters shipped
- `caf.first_person_explain` produces accurate narratives (no hallucination)

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

| Phase | Maintainer hours/week | Community PRs expected |
|---|---|---|
| MVP | 15-20 | 0-2 |
| ALPHA | 20-25 | 2-5 |
| BETA | 25-30 | 5-10 |
| PROD | 20-25 (stabilization) | 10-15 |
| FUTURE | variable | 15+ |

При full-time engagement (1 maintainer, 40h/week) сроки сжимаются примерно в 2 раза.

**🔥 80/20 implication:** high-leverage features (universal hook, LSP fallback, dashboard) — small effort, big value. Most of MVP budget should go to these, not to custom adapters for Tier 2/3 agents.

---

## 10. Open Questions

Вопросы, требующие решения до или во время соответствующих фаз:

1. **MVP:** Должен ли `caf init --auto-detect` предлагать все обнаруженные адаптеры или только первый? — Решение: предлагать все, с пометкой рекомендуемого.
2. **MVP:** Universal LSP fallback — какой LSP server по умолчанию для Go? `gopls` — но он может быть не установлен. — Решение: bundle `gopls` в `caf` binary или требовать явную установку.
3. **ALPHA:** Какой порог `triviality_threshold_loc` по умолчанию? — Решение: 50 (см. ADR-004), пересматривается по данным.
4. **BETA:** Какой LLM использовать для `first_person_explain` — z-ai-web-dev-sdk, локальная модель, или configurable? — Решение: configurable, с z-ai-web-dev-sdk как default.
5. **PROD:** Должен ли `caf sync` работать с приватными репо? — Решение: да, через git credential helper.
6. **FUTURE:** Должен ли multi-user mode быть backward-compatible с single-user `.caf/`? — Решение: да, через optional `author` field.

Эти вопросы будут решены через ADRs по мере приближения к соответствующим фазам.

---

## 11. v0.2 → v0.3 Migration Notes

### 11.1 What changed

| Aspect | v0.2 | v0.3 |
|---|---|---|
| Agent count | 9 explicit | 25 explicit + universal LSP |
| Default install | `caf init` (no args = nothing) | `caf init` (no args = universal hook + LSP) 🔥 |
| Adapter tiers | 3 tiers | 4 tiers (added Tier 4 universal LSP) |
| Prioritization | feature-parity across agents | 🔥 80/20 prioritization |
| UI | CLI output (text-based) | TUI-first (dashboard, full-screen review) |
| Pre-commit hook | per-adapter | 🔥 universal (works for all agents) |

### 11.2 Migration path

1. **`.caf/` schema:** additive changes only — `agent_tier`, `auto_detected`, `ui` config section. Old `.caf/` continues to work.
2. **`caf init` behavior:** v0.3 `caf init` (no args) installs universal hook. v0.2 users who ran `caf init --agent <id>` get the same explicit adapter, plus universal hook as bonus.
3. **CLI commands:** all v0.2 commands still work. New: `caf dashboard`, `caf heatmap`, `caf tour`, `caf adapters list` (TUI).
4. **No breaking changes:** v0.2 configs work in v0.3 without modification.

### 11.3 What's deprecated in v0.3

Nothing. v0.3 is purely additive.

---

## 12. Success Measurement

### 12.1 Per-phase metrics

| Phase | Primary metric | Target |
|---|---|---|
| MVP | Universal hook fires for ≥3 different agents | ≥3 agents |
| ALPHA | Decision Coverage on pilot repos | ≥90% |
| BETA | Recall Rate sustained 4 weeks | ≥75% |
| PROD | All 6 success metrics at target | all met |
| FUTURE | (per item, see BACKLOG) | per item |

### 12.2 Adoption metrics (post-PROD)

| Metric | Target (12 months post-1.0) |
|---|---|
| GitHub stars | ≥500 |
| Active users (weekly) | ≥100 |
| Adapters contributed by community | ≥5 new |
| Languages supported | ≥3 (Go + 2) |
| Issues closed | ≥80% within 30 days |
