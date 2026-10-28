---
doc: ARCHITECT
layer: architecture
version: 0.3.1-triz
language: bilingual
status: stable
audience: [agents, maintainers, contributors, anyone reasoning about system boundaries]
related: [ROADMAP.md, UI.md, BACKLOG.md, ADR.md]
depends_on: [ADR-001, ADR-002, ADR-003, ADR-004, ADR-005, ADR-009, ADR-010, ADR-011, ADR-012, ADR-013, ADR-014, ADR-015]
invariants:
  - {id: INV-ARCH-001, statement: "Storage layer (.caf/) is the single source of truth; index files are regenerable caches", criticality: critical, enforced_by: "caf doctor"}
  - {id: INV-ARCH-002, statement: "Adapter interface (5 methods) is stable within a major version; canonical events are exactly {pre-session, post-edit, pre-commit}", criticality: critical, enforced_by: "interface contract tests, 25×3=75 cases in CI"}
  - {id: INV-ARCH-003, statement: "Universal LSP fallback (Tier 4) is always present; lazy-loaded on caf sync / caf tour", criticality: high, enforced_by: "caf sync benchmark"}
  - {id: INV-ARCH-004, statement: "Hook execution is fail-safe: caf binary missing → warn-and-continue, never block session/commit", criticality: critical, enforced_by: "adapter contract tests"}
  - {id: INV-ARCH-005, statement: "TUI first paint <100ms on pilot repos; >2x CI regression triggers alert", criticality: high, enforced_by: "CI benchmark — see UI.md §0.3"}
  - {id: INV-ARCH-006, statement: "ease_factor ∈ [1.3, 3.0] for every active SRS card", criticality: high, enforced_by: "property-based tests"}
  - {id: INV-ARCH-007, statement: ".caf/ folder is committed in-repo and additive across major versions; .caf/ size per repo stays under 50 MB", criticality: high, enforced_by: "caf doctor size-check"}
owned_contracts:
  - "Adapter interface (Go): ID() string; Capabilities() CapabilityMatrix; Install() error; Uninstall() error; Detect() bool"
  - "LanguageAdapter interface (Go): ID() string; Detect() bool; ScanImports() []Import; BuildSymbolGraph() SymbolGraph; ResolveExternalRepo() RepoRef"
  - "Canonical event contract: {pre-session, post-edit, pre-commit} × {native, wrapper, watcher, universal}"
  - ".caf/ schema: sessions/*.yml, decisions/*.md, mental-model/graph.yml, manifest.yml, config.yml, cards.yml"
  - "MCP tool contract: caf.context, caf.provenance, caf.decisions, caf.mental_model, caf.first_person_explain"
agent_routing:
  - {task: "Add new AI agent adapter", see: [§1 Capability matrix, ADR-009, ADR-012], backlog: [E-007, H-002]}
  - {task: "Add new language adapter", see: [§6 FUTURE, ADR-011], backlog: [E-001, E-002, E-003, E-004, E-005, E-006, E-007, E-008]}
  - {task: "Modify .caf/ storage schema", see: [§7.4 Backwards compatibility, ADR-001, ADR-005], backlog: [A-005, J-008]}
  - {task: "Change TUI screen behavior", see: [UI.md, ADR-014], ui_budgets: UI.md §0.3}
  - {task: "Add new SRS card type", see: [§4 BETA, ADR-002], backlog: [H-001, D-007]}
  - {task: "Add MCP tool", see: [§5 PROD, ADR-006], backlog: [D-001, D-008]}
  - {task: "Investigate invariant violation", see: [front-matter invariants, ADR.md @invariant tags]}
  - {task: "Understand fail-safe behavior", see: [§2.5 Risks (fail-safe contract within)]}
trace_to_code:
  - "internal/caf/storage/ — .caf/ folder read/write (ADR-001, ADR-005)"
  - "internal/caf/adapter.go — Adapter, LanguageAdapter interfaces (ADR-009, ADR-011)"
  - "internal/caf/hooks/*.go — pre-session/post-edit/pre-commit hook contracts (ADR-003)"
  - "internal/caf/lsp/fallback.go — Tier 4 universal LSP fallback (ADR-013)"
  - "internal/caf/srs/sm2.go — SM-2 implementation (ADR-002)"
  - "internal/caf/mcp/server.go — MCP server with 5 tools + 3 prompts (ADR-006)"
  - "internal/caf/triviality/threshold.go — 50 LOC threshold (ADR-004)"
  - "pkg/ui/ — TUI screens (ADR-014)"
ui_budgets_enforced:
  - first_paint: "<100ms"
  - frame_render: "<16ms (60fps)"
  - keystroke_response: "<50ms"
  - llm_op: "<3s"
  - file_op_10k_repo: "<1s"
  - dashboard_full_load: "<500ms"
  - dashboard_auto_refresh: "<50ms (incremental, only changed widgets)"
  - card_render: "<50ms (no LLM)"
  - session_save: "<100ms (incremental write)"
  - list_render: "<50ms (virtualized)"
  - filter_response: "<16ms (in-memory index)"
  - autocomplete: "<50ms (pre-built index)"
  - heatmap_render: "<200ms (247 files)"
  - heatmap_recompute: "<500ms (incremental)"
  - adapter_health_check: "<5s/adapter (<30s/25)"
  - adapter_install: "<2s/adapter (parallel ≤5)"
  - file_tour_first_paint: "<100ms"
  - file_load: "<50ms (cached)"
  - timeline_build: "<50ms (from .caf/sessions/ index)"
  - lsp_startup: "~500ms (gopls, small repo)"
  - animated_transitions: "≤100ms"
---

# CodeAuthorityFirst — ARCHITECT.md

> **Версия / Version:** 0.3.1-triz (post-TRIZ pass: this file is no longer a byte-identical duplicate of ROADMAP.md; content rewritten to focus on architectural evolution per phase. Capability matrix, fail-safe contract, and config schema now live in §1 / §2.5 / §7.4 respectively, fixing the four dangling cross-references from ADR.md.)
> **Назначение / Purpose:** describes the **architectural components, layer boundaries, and code contracts** that ship per phase. Where ROADMAP answers "what features ship when", ARCHITECT answers "what architectural layers and contracts ship when, and what must remain invariant across them". Together with [ADR.md](./ADR.md), this file is the canonical reference for any agent reasoning about system boundaries before touching code.
> **Язык / Language:** bilingual — Russian body, English summary per phase
> **Аудитория / Audience:** maintainers, contributors proposing architectural changes, LLM agents traversing the codebase before writing code
> **Связанные документы:** [ROADMAP.md](./ROADMAP.md) · [UI.md](./UI.md) · [BACKLOG.md](./BACKLOG.md) · [ADR.md](./ADR.md)

> **How to read this file as an agent.** The YAML front-matter above is the machine-readable layer: it lists invariants, owned contracts, agent_routing (task → which sections to read), and trace_to_code (which Go paths implement which ADR). Read the front-matter first. Then read the §1 Capability matrix to know which adapter is in which tier. Only then drill into the per-phase sections.

---

## 0. Roadmap Philosophy (v0.3)

**EN summary.** The architecture is organized around four invariants: (1) `.caf/` is the single source of truth, all indices are regenerable caches; (2) the Adapter interface is stable within a major version, canonical events are exactly three; (3) hooks are fail-safe — `caf` binary missing means warn-and-continue, never block; (4) TUI-first, with web UI deferred to FUTURE and only for visualizations impossible in TUI. The 80/20 principle (ADR-015) determines which architectural components ship in which phase.

Архитектурная философия v0.3 построена на четырёх инвариантах, которые остаются неизменными через все фазы:

1. **🔥 INV-ARCH-001: `.caf/` — единственный source of truth.** Все `index.yml`, `graph.yml` и SQLite-кэши являются regenerable артефактами; `caf doctor` восстанавливает их из YAML source. Это означает, что любое предложение «давайте хранить graph в SQLite» встречает вопрос: «а является ли SQLite regenerable из YAML?» — если нет, нарушает инвариант.
2. **🔥 INV-ARCH-002: Adapter interface стабилен внутри major version.** Пять методов (`ID()`, `Capabilities()`, `Install()`, `Uninstall()`, `Detect()`); три canonical events (`pre-session`, `post-edit`, `pre-commit`). Новые события требуют нового ADR (см. ADR-009 trigger).
3. **🔥 INV-ARCH-004: Hooks fail-safe.** Если `caf` binary недоступен, агент продолжает работу с warning, но не блокируется. Это критично для принятия: один сбой pre-commit hook не должен ломать рабочий процесс автора.
4. **INV-ARCH-005: TUI-first.** Бюджеты (first paint <100ms, frame render <16ms, keystroke <50ms) — часть архитектурного контракта, не оптимизация. Web UI отложен до FUTURE и только для визуализаций, невозможных в TUI.

Дополнительно: v0.3 меняет v0.2 подход с "specific adapters first" на "universal fallback first" (см. ADR-013). Это значит, что 25+ агентов покрываются с первого дня MVP через universal LSP fallback, а явные адаптеры добавляются для enhanced UX, а не для base coverage.

### 0.1 High-leverage features list (🔥 80/20)

| # | Architectural component | Effort | Value | Why 80/20 | Phase |
|---|-------------------------|--------|-------|-----------|-------|
| 1 | Universal `.git/hooks/pre-commit` shell | 2d | 90% coverage of agents | One hook fires for all 25+ agents | MVP |
| 2 | Auto-detect (`caf init --auto-detect`) | 1d | 5min→30s onboarding | Removes friction; user has nothing to configure | MVP |
| 3 | Universal LSP fallback (ADR-013, Tier 4) | 5d | Long-tail coverage | Zero-config support for all LSP languages | MVP |
| 4 | TUI dashboard | 3d | Hub for all screens | Replaces 6 separate CLI commands | MVP |
| 5 | `caf.first_person_explain` MCP tool | 5d | Killer feature for agents | LLM gets narrative-grade context, grounded in `.caf/` | PROD |
| 6 | Single Go binary | 1d | Zero-dep distribution | One file, no runtime; SSH/Docker/CI universal | MVP |
| 7 | ASCII heatmap | 1d | One-screen overview | Replaces 100-line text report | BETA |
| 8 | SM-2 algorithm | 2d | Battle-tested SRS | Anki users recognize it immediately; ease ∈ [1.3, 3.0] | BETA |

**Итого:** 8 компонентов, ~20 дней effort, ~80% ценности. Это и есть MVP+ALPHA+BETA архитектурное ядро. См. также [ADR-015](./ADR.md#adr-015-8020-prioritization-principle) и [BACKLOG §0](./BACKLOG.md#0-backlog-philosophy) для процесса приоритезации.

---

## 1. Phase Overview (v0.3)

**EN summary.** Five phases ship architectural components in this order: MVP ships storage, hooks, universal LSP fallback, TUI dashboard — the universal-coverage layer; ALPHA adds Decision Log and pre-commit validation; BETA adds Mental Model Reconstructor with SM-2 and three card types; PROD adds MCP server, cross-repo sync automation, and `caf.first_person_explain`; FUTURE adds multi-user, web UI (read-only), CRDT, plugins. The per-adapter Capability matrix is the table below; it is referenced from ADR-009 as a first-class artifact.

| Phase | Duration (est.) | Architectural layer shipped | Modules live | Agents supported | UI screens live | 🔥 80/20 components |
|-------|-----------------|-----------------------------|--------------|-------------------|------------------|----------------------|
| MVP | 2-3 wk | Storage + Hooks + Universal fallback + TUI hub | Provenance Tracker, `context.md` gen, Universal pre-commit hook, Universal LSP fallback, Auto-detect, TUI dashboard, TUI onboarding, TUI adapters manager | 25+ via universal LSP + Tier 1 (2 explicit) | Dashboard, Onboarding, Adapters manager | Universal pre-commit hook, Auto-detect, LSP fallback, TUI dashboard, Single Go binary |
| ALPHA | 4-6 wk | Decision Log layer + pre-commit validation | Decision Log, TUI decisions list/editor/detail, pre-commit validation, Tier 1 +3 explicit adapters (Cursor, Cline, KiloCode) | 25+ via universal + 5 enhanced (Tier 1) | Decisions list, Decision detail, Decision editor | 80/20 autocomplete (90s→25s) |
| BETA | 6-8 wk | Mental Model Reconstructor + SRS | MMR, SM-2 SRS, 3 card types, dependency graph, blindspots detection, TUI review session, TUI heatmap, TUI file tour | 25+ via universal + 13 enhanced (+8 Tier 2) | Review session, Heatmap, File tour | SM-2 algorithm, ASCII heatmap, TUI review session |
| PROD | 8-10 wk | MCP server + cross-repo sync automation | MCP server (5 tools + 3 prompts), cross-repo sync, TUI reports view, success metrics instrumentation, `caf.first_person_explain` 🔥, +12 Tier 3 adapters | 25+ via universal + 25 enhanced (all explicit) | Reports view, Metrics report | `caf.first_person_explain` MCP tool, cross-repo sync |
| FUTURE | open-ended | Extensibility + multi-user | Multi-user (per-author attribution), Web UI (read-only), CRDT, Plugin API, IDE plugins, CI/CD integrations | (no change to base; additive) | Web dashboard (read-only), VSCode/JetBrains/Neovim plugins | (per-item 4 promotion criteria — see ROADMAP §6.9) |

> **🔥 Key shift from v0.2:** MVP теперь покрывает 25+ агентов через universal fallback, а не 2 явных адаптера. Это радикально расширяет initial coverage и меняет архитектурную стратегию с "specific-first" на "universal-first, enhanced-second".

### Capability matrix (first-class artifact, referenced from ADR-009 and ADR-013)

| Adapter | Tier | pre-session | post-edit | pre-commit | MCP | Phase |
|---------|------|-------------|-----------|------------|-----|-------|
| opencode | 1 | ✓ native | ✓ native | ✓ native | ✓ | MVP |
| Claude Code | 1 | ✓ native | ✓ native | ✓ native | ✓ | MVP |
| Cursor | 1 | ✓ wrapper (pre-session via shell alias) | ✓ native | ✓ native | ✓ | ALPHA |
| Cline | 1 | ✓ native | ✓ native | ✓ native | ✓ | ALPHA |
| KiloCode | 1 | ✓ native | ✓ native | ✓ native | ✓ | ALPHA |
| Aider | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Codex CLI | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Continue.dev | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| GitHub Copilot CLI | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Gemini CLI | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Amazon Q Developer CLI | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Goose | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Windsurf | 2 | ✗ | ✓ wrapper | ✓ wrapper | — | BETA |
| Tabnine | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Sourcegraph Cody | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Trae | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Bolt | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Zed AI | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Roo Code | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Augment | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Plandex | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| OpenHands | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Devika | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| GPT Engineer | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Mentat | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| Sweep | 3 | ✗ | ✗ | ✓ wrapper | — | PROD |
| **Universal LSP fallback** | **4** | ✓ `caf-wrap` wrapper | ✓ `caf-watch` watcher (5s debounce) | ✓ universal `.git/hooks/pre-commit` | ✓ `caf-wrap` prepends `context.md` | **MVP** |

**Notes:**
- Tier 1 = full native hook support; Tier 2 = partial native + wrapper; Tier 3 = wrapper-only (pre-commit); Tier 4 = universal fallback (zero-config, LSP-based).
- Total: 25 explicit + 1 universal = 26 rows in matrix (25 + universal).
- "Augment" (not "Augment Code") — vendor's preferred branding; consistent across ADR-012 and ROADMAP §5.1 in v0.3.1-triz.
- OryaAgent is the in-house reference implementation, not in the matrix; see §6.5 below for its role.
- New event types beyond {pre-session, post-edit, pre-commit} require a new ADR (see ADR-009 trigger).

---

## 2. Phase: MVP (2-3 weeks)

**EN summary.** MVP ships the universal-coverage architectural layer: `.caf/` storage (ADR-001), fail-safe hooks (ADR-003 + INV-ARCH-004), universal LSP fallback (ADR-013), TUI dashboard, single Go binary. After MVP, the tool is already useful for any agent, even if development stops.

### 2.1 Scope

Архитектурные компоненты, которые должны быть готовы к концу MVP:

- **Storage layer (ADR-001, ADR-005):** `.caf/` folder schema — `sessions/*.yml`, `decisions/*.md`, `mental-model/graph.yml`, `manifest.yml`, `config.yml`, `cards.yml` (empty in MVP). YAML 1.2 + `gopkg.in/yaml.v3`.
- **Hook contract (ADR-003, INV-ARCH-004):** three event types — `pre-session` (recorded by adapter when agent starts), `post-edit` (recorded after non-trivial file changes), `pre-commit` (universal `.git/hooks/pre-commit` shell, fires for ALL agents).
- **Adapter pattern (ADR-009, ADR-012):** `Adapter` interface with 5 methods; 4 adapters shipped (Universal LSP Tier 4 + opencode Tier 1 + Claude Code Tier 1 + Universal pre-commit hook).
- **Universal LSP fallback (ADR-013, INV-ARCH-003):** `textDocument/documentSymbol` + `workspace/symbol` for symbol graph; `caf-watch` daemon with 5-second debounce; lazy-load on `caf sync`/`caf tour`.
- **Auto-detect (`caf init --auto-detect`):** scans for installed agents, recommends with recommended marker.
- **TUI dashboard (ADR-014, INV-ARCH-005):** hub screen showing Today widget, Authorship Score widget, Active Agents widget, Recent Decisions widget, Cross-repo widget.
- **Single Go binary:** `caf` (3 chars, ADR-008), no runtime dependencies, ~10 MB.

CLI commands shipped (MVP): `caf init`, `caf init --auto-detect`, `caf init --agent opencode`, `caf dashboard`, `caf adapters list`, `caf who <file:line>`, `caf when --pattern <regex> <dir>`, `caf sessions --since 7d`, `caf context`, `caf hook pre-session`, `caf hook post-edit`, `caf hook pre-commit` 🔥.

### 2.2 Deliverables

- `cmd/caf/` — single-binary entry point.
- `internal/caf/storage/` — `.caf/` read/write, atomic writes (temp + rename) for crash safety (BACKLOG J-002).
- `internal/caf/hooks/` — three hook contracts + universal `.git/hooks/pre-commit` installer.
- `internal/caf/lsp/fallback.go` — Tier 4 universal LSP fallback (ADR-013).
- `internal/adapters/{opencode,claude-code}/` — two Tier 1 explicit adapters.
- `pkg/ui/screens/{dashboard,onboarding,adapters}.go` — three TUI screens (see UI.md §2, §3, §7).
- `pkg/ui/keys.go` — universal keyboard map (see UI.md §0.6).
- `pkg/ui/perf.go` — performance budget enforcer (see UI.md §0.3).
- `internal/caf/triviality/threshold.go` — 50 LOC threshold (ADR-004; this is shipped in MVP for `post-edit` to know whether to record).
- Contract tests: `internal/adapters/*/contract_test.go` — one test per cell of the 25-adapter × 3-event matrix (75 cases total by PROD; in MVP only 4 adapters × 3 events = 12 cases).

### 2.3 Exit criteria

1. Pilot repo uses `caf` 7 days straight, no manual intervention.
2. Universal pre-commit hook fires on 3+ different agents (opencode, Claude Code, +1 unknown via fallback).
3. LSP fallback builds symbol graph for Go repo up to 1000 files in <1s (INV-ARCH-003 enforced).
4. **TUI dashboard first paint <100ms** on pilot repo (INV-ARCH-005 enforced; CI benchmark in MVP).
5. Auto-detect finds opencode + Claude Code on a fresh machine in <5s.
6. Adapter contract test suite passes (12 cases for 4 adapters).

### 2.4 What's NOT in MVP

- Decision Log (architectural layer: ALPHA).
- SRS / Mental Model Reconstructor (architectural layer: BETA).
- MCP server (architectural layer: PROD).
- Tier 1 beyond opencode + Claude Code (3 more in ALPHA).
- Tier 2/3 adapters (BETA/PROD).
- Cross-repo sync (PROD).
- LLM features (PROD, ADR-006).

### 2.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| LSP fallback noisy (false-positive symbol graphs) | high | medium | `caf doctor --verify-graph` rebuilds from filesystem-only mode; `caf config set lsp.strict false` relaxes resolution |
| `caf-watch` filesystem watcher CPU on large repos | medium | medium | Throttle: poll every 1s, debounce 5s, ignore `.git/` and other noisy paths; daemon exits after 30 minutes of inactivity (INV-ARCH-003) |
| TUI rendering issues on weird terminals | medium | low | UI.md §12.1 compatibility matrix (iTerm2/Terminal.app/Alacritty/Kitty/Windows Terminal/VSCode integrated/Git Bash/Linux TTY/SSH); 16-color fallback via termenv |
| `caf init` confusing for new users | medium | high | Onboarding wizard (UI.md §2); auto-detect as default; 3-step flow <30s |
| **Universal pre-commit hook blocks legit commits** | **medium** | **high** | **Fail-safe contract** (see below): if `caf` binary missing, hook warns and lets commit through. If `caf` exits non-zero for reasons other than missing decision (ALPHA), hook also lets commit through with warning. **Hook never blocks a commit due to `caf` itself being broken** — only blocks if `caf hook pre-commit` explicitly returns "decision required" exit code (ALPHA). |

#### Fail-safe contract (referenced from ADR-003 INV-ADR-003 and ADR-010)

The fail-safe contract is the architectural guarantee that hooks never block agent operation due to `caf` infrastructure issues. Concretely:

1. **Binary missing.** If `.git/hooks/pre-commit` cannot find `caf` in `$PATH`, the hook writes a single-line warning to stderr and exits 0. The commit proceeds.
2. **Binary crashes.** If `caf hook pre-commit` exits with a non-zero status code other than the reserved "decision required" code (exit 2), the hook treats this as `caf` infrastructure failure, writes a warning, and exits 0. The commit proceeds.
3. **Reserved exit codes.** Exit 2 = "decision required" (ALPHA, blocks commit until decision added or `CAF-Trivial: true` trailer set). Exit 3 = "triviality threshold exceeded but no decision" (ALPHA, same as 2). Exit 0 = "all good, commit proceeds". No other exit code blocks the commit.
4. **`caf-watch` daemon.** If daemon is unavailable, post-edit hook is skipped (not recorded). Sessions continue. Graph stays stale until next `caf sync`.

This contract is **non-negotiable**: any PR that makes `caf` block a commit for reasons other than "decision required" violates INV-ARCH-004 and must be rejected in code review.

### 2.6 Phase: MVP — additional architectural notes

(MVP-specific architectural notes moved here from the historical scope. Section preserved from original structure.)

---

## 3. Phase: ALPHA (4-6 weeks)

**EN summary.** ALPHA adds the Decision Log layer (architectural): `.caf/decisions/*.md` schema, pre-commit validation requiring decisions for non-trivial commits, three more Tier 1 adapters (Cursor, Cline, KiloCode). The 80/20 win: autocomplete for files/tags/agent reduces decision fill time from 90s to 25s.

### 3.1 Scope

- **Decision Log layer (ADR-001, ADR-005):** `.caf/decisions/*.md` schema — title, date, agent, files[], context, chosen, alternatives, implications. Markdown-rendered in TUI via Glamour.
- **Pre-commit validation (ADR-003, ADR-004, INV-ADR-004):** `caf hook pre-commit` exit 2 if any non-trivial edit (≥50 LOC) lacks a decision. `CAF-Trivial: true` trailer bypasses for mechanical refactors.
- **Three more Tier 1 adapters:** Cursor (pre-session via wrapper shell alias), Cline (native all events), KiloCode (native all events).
- **TUI decisions list / detail / editor (ADR-014):** three new TUI screens (see UI.md §5.1, §5.2, §5.3).
- **🔥 80/20 autocomplete (UI.md §5.3):** autocomplete for `files` (from paths modified in current/recent sessions), `tags` (from previous decisions), `agent` (from active agents). Pre-built index, <50ms response (INV-ARCH-005).
- **Triviality threshold enforcement (ADR-004, INV-ADR-004):** `internal/caf/triviality/threshold.go` ships in MVP for `post-edit`; in ALPHA it gates pre-commit. 50 LOC excluding imports/formatting/test-duplicates.

New CLI commands (ALPHA): `caf why <file>`, `caf decisions` 🔥, `caf decisions new [--files ...]` 🔥, plus `--today`/`--file`/`--search`/`--gaps`/`--agent` flags on `caf decisions`.

### 3.2 Deliverables

- `internal/caf/decisions/` — Decision Log CRUD, markdown rendering.
- `internal/adapters/{cursor,cline,kilocode}/` — three more Tier 1 adapters.
- `pkg/ui/screens/{decisions_list,decisions_detail,decisions_editor}.go` — three new TUI screens.
- `internal/caf/autocomplete/` — pre-built indexes for files, tags, agents.
- 3 new adapter contract tests (3 adapters × 3 events = 9 cases; total 12+9=21 cases).

### 3.3 Exit criteria

1. Decision Coverage ≥90% over last 7 days on pilot repo (number of non-trivial commits with decisions / total non-trivial commits).
2. Pre-commit blocks non-trivial commits without decision for ALL 25+ agents (universal fallback + 5 explicit).
3. **TUI decision editor avg fill time ≤30s (vs 90s without autocomplete)** — measured via instrumentation.
4. All Tier 1 adapters work on macOS + Linux with no manual config.

### 3.4 What's NOT in ALPHA

- Mental Model Reconstructor (architectural layer: BETA).
- Tier 2/3 adapters (BETA/PROD).
- Cross-repo sync (PROD).
- LLM features (PROD).

### 3.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Autocomplete too noisy (suggests irrelevant files) | medium | medium | `caf config set decisions.autocomplete aggressive\|conservative\|off` |
| Pre-commit blocks CI commits (no human to write decision) | high | high | Detect `CI=true` env, skip pre-commit validation; CI runs `caf report --check` separately (BACKLOG C-001) |
| Decision template too verbose | high | high | Default compact template; `--verbose` for full template |
| Cursor wrapper breaks on Cursor version update | medium | medium | Pin wrapper version; auto-detect breaks, surface in `caf adapters list` |
| Inconsistent UX between adapters | medium | low | Capability matrix in §1 makes limitations explicit; UI.md §7 surfaces them |

---

## 4. Phase: BETA (6-8 weeks)

**EN summary.** BETA adds the Mental Model Reconstructor layer: SM-2 spaced-repetition algorithm (ADR-002, ease ∈ [1.3, 3.0]), three card types (`why`, `what_if_remove`, `where_handled`), single-repo dependency graph builder, blindspots detection, three new TUI screens (review session, heatmap, file tour). Eight Tier 2 adapters join the matrix.

### 4.1 Scope

- **Mental Model Reconstructor (ADR-002):** SM-2 algorithm; ease ∈ [1.3, 3.0] invariant (INV-ARCH-006 enforced by property-based tests).
- **Three card types:** `why` (Type 1, fact recall — baseline ease 2.5), `what_if_remove` (Type 2, abstract reasoning — baseline ease 2.0 per BACKLOG D-007 future), `where_handled` (Type 3, trace through code — baseline ease 2.5).
- **Single-repo dependency graph builder:** `go list -json` (Go-only in BETA; BACKLOG E-001..E-008 cover other languages in FUTURE). Graph cached in `.caf/mental-model/graph.yml`.
- **Blindspots detection:** files with low ACS or low coverage. Risk score = frequency × centrality × last_review (BACKLOG D-004 refines).
- **TUI Review Session (UI.md §4, ADR-014):** full-screen split-view, ~10 minutes, 10-15 cards per session, grade scale 0-5.
- **TUI Heatmap (UI.md §6, ADR-014):** each file = one cell, color = ACS or coverage. Pre-computed, cached.
- **TUI File Tour (UI.md §8, ADR-014):** split-view timeline (left) + numbered source (right). Phase-conflict fix in v0.3.1-triz: original UI.md §1.1 listed File tour as PROD; ROADMAP §4.1 BETA scope listed it as BETA deliverable. Canonical resolution: File tour = BETA (per ROADMAP §4.1); UI.md §1.1 will be updated in TRIZ-pass to BETA.

8 Tier 2 adapters (cumulative): Aider, Codex CLI, Continue.dev, GitHub Copilot CLI, Gemini CLI, Amazon Q Developer CLI, Goose, Windsurf. Coverage: 25+ via universal + 13 enhanced.

New CLI commands (BETA): `caf review` 🔥, `caf review --hard`/`--blindspots`/`--regenerate --since 7d`/`--stats`, `caf heatmap`, `caf tour <file>`, `caf sync`.

### 4.2 Deliverables

- `internal/caf/srs/sm2.go` — SM-2 implementation with ease ∈ [1.3, 3.0] invariant.
- `internal/caf/mentalmodel/` — graph builder, blindspots detection.
- `pkg/ui/screens/{review,heatmap,file_tour}.go` — three new TUI screens.
- 8 new adapter directories under `internal/adapters/`.
- Property-based tests for SM-2 (invariant: ease ∈ [1.3, 3.0]).
- Performance benchmark: `caf heatmap` renders 247 files in <200ms (INV-ARCH-005 sub-budget).

### 4.3 Exit criteria

1. Mental Model Reach ≥70% (cards due / cards generated).
2. Recall Rate ≥75% sustained over 4 weeks.
3. ACS (Authorship Confidence Score) ≥3.0 (target ≥3.5 in PROD; ≥3.5 is the §12 success threshold).
4. Review session completion ≥80% (cards answered / cards started).
5. All Tier 2 adapters work on macOS + Linux.
6. **Heatmap render <200ms on repo up to 1000 files** (INV-ARCH-005 sub-budget enforced).

### 4.4 What's NOT in BETA

- MCP server (architectural layer: PROD).
- Cross-repo sync (PROD).
- `caf.first_person_explain` LLM feature (PROD, ADR-006).
- Tier 3 adapters (PROD).
- Web UI (FUTURE, only for visualizations impossible in TUI).

### 4.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| SM-2 burnout (too many cards per session) | high | high | Default cap 15 cards/session; `caf config set review.max_cards 10`; BACKLOG D-007 adaptive per-card-type ease is the long-term fix |
| Cramped TUI on small terminals | medium | medium | UI.md §12.2 responsive: ≥140 cols full grid, 100-139 2-col, 60-99 single stacked, <60 compact mode with hint "Widen terminal for full dashboard." |
| Low-quality cards (LLM-generated, BACKLOG D-002 future) | high | medium | (Not in BETA — D-002 is idea); BETA ships hand-curated cards only |
| Tier 2 wrappers fail on agent version updates | medium | medium | Pin wrapper versions; auto-detect breaks, surface in `caf adapters list` |
| Graph builder fails on >10k files | medium | high | Streaming + lazy loading (BACKLOG A-001); benchmark on `kubernetes/kubernetes` before release |

---

## 5. Phase: PROD (8-10 weeks)

**EN summary.** PROD adds the MCP server layer (architectural): 5 MCP tools (`caf.context`, `caf.provenance`, `caf.decisions`, `caf.mental_model`, `caf.first_person_explain` 🔥) + 3 prompts (`caf.pre_session_brief`, `caf.code_review_with_provenance`, `caf.archaeology`); cross-repo sync automation; Tier 3 adapters (12 agents, wrapper-only); success metrics instrumentation.

### 5.1 Scope

- **MCP server (`caf-mcp`):** 5 tools + 3 prompts. Transport: stdio (`caf mcp serve`) and HTTP SSE (`caf mcp serve --transport sse`).
- **`caf.first_person_explain` 🔥 (ADR-006):** LLM-generated first-person narrative grounded in `.caf/` facts. INV-ADR-006: LLM only paraphrases `.caf/`-grounded facts, never invents.
- **Cross-repo sync automation (ADR-001):** `caf sync` detects locally-cloned repos via `~/.caf/repos.yml` registry; pulls `manifest.yml` from each dependee; creates `review-because-interface-changed` cards.
- **TUI Reports view (UI.md §10, ADR-014):** daily report, week report, metrics report.
- **Success metrics instrumentation:** ACS, Decision Coverage, Mental Model Reach, Recall Rate, Agent Diversity (Shannon entropy), Adapter Coverage.
- **12 Tier 3 adapters (cumulative):** Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Augment, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep. Note: in v0.3.1-triz, the list is reconciled to 12 (OryaAgent removed — it is the in-house consumer, listed separately; Augment renamed consistently). Total explicit: 5 (Tier 1) + 8 (Tier 2) + 12 (Tier 3) = 25.

MCP tools:
- `caf.context` — returns `context.md` for current repo.
- `caf.provenance` — returns provenance for `file:line` range.
- `caf.decisions` — returns decisions matching filter.
- `caf.mental_model` — returns graph + card inventory.
- `caf.first_person_explain` 🔥 — returns LLM-generated first-person narrative, grounded in `.caf/` facts (ADR-006).

MCP prompts:
- `caf.pre_session_brief` — what an agent should know before editing.
- `caf.code_review_with_provenance` — review a PR with authorship context.
- `caf.archaeology` — explain why a file looks the way it does.

New CLI commands (PROD): `caf report --day`/`--week`/`--ownership`/`--agent`/`--metrics`, `caf sync`, `caf hook post-commit --repo <path>`, `caf mcp serve [--transport sse]`.

### 5.2 Deliverables

- `internal/caf/mcp/server.go` — MCP server with 5 tools + 3 prompts.
- `internal/caf/sync/` — cross-repo sync automation, manifest cache (1h TTL in `~/.caf/manifest-cache/`).
- `pkg/ui/screens/reports.go` — TUI Reports view.
- 12 new adapter directories under `internal/adapters/` (Tier 3, wrapper-only).
- Success metrics instrumentation in `internal/caf/metrics/`.
- Contract tests for all 25 explicit + 1 universal adapter × 3 events = 78 cases in CI.

### 5.3 Exit criteria

1. All 5 success metrics hit targets on pilot repos over 8 weeks: **ACS ≥3.5, Decision Coverage ≥90%, Mental Model Reach ≥70%, Recall Rate ≥75%, Agent Diversity H ≥1.0, Adapter Coverage ≥70% native (NEW v0.3)**.
2. MCP works with opencode, Claude Code, Cline, KiloCode, Goose.
3. Cross-repo sync detects changed interfaces between dependent repos.
4. All 25 explicit adapters pass contract tests.
5. `caf.first_person_explain` produces factually accurate narratives (no hallucination on `.caf/`-grounded facts — INV-ADR-006 enforced).
6. README, ARCHITECT.md, UI.md, ROADMAP.md, BACKLOG.md, ADR.md all consistent with implementation.

### 5.4 What's NOT in PROD

- Multi-user mode (FUTURE).
- Web UI (FUTURE).
- CRDT for `cards.yml` (FUTURE).
- Plugin API / marketplace (FUTURE).
- IDE plugins — VSCode/JetBrains/Neovim (FUTURE).
- CI/CD integrations beyond `caf report --check` (FUTURE).

### 5.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| LLM hallucinates in `caf.first_person_explain` | high | high | Strict grounding in `.caf/` facts (ADR-006, INV-ADR-006); LLM only paraphrases, never invents; test suite verifies no hallucination on `.caf/`-grounded facts |
| MCP protocol changes (Anthropic updates spec) | medium | medium | Pin MCP spec version; contract tests against pinned version; ADR-028 future addresses adapter contribution governance |
| Cross-repo sync slow on large fleets | medium | medium | Parallel fetch (10 concurrent, rate-limited); cache `manifest.yml` for 1 hour in `~/.caf/manifest-cache/`; incremental sync (BACKLOG A-003) |
| Tier 3 brittle across 12 wrapper-only adapters | high | medium | Wrapper contracts minimal; failure in one adapter doesn't block others; `caf adapters health` surfaces breakages (BACKLOG J-003 — already done in v0.3) |
| TUI performance degrades as data grows | medium | high | Performance budgets enforced (UI.md §0.3, INV-ARCH-005); CI benchmarks; alert on >2x regression; BACKLOG A-001 lazy graph loading as long-term fix |

---

## 6. Phase: FUTURE (open-ended)

**EN summary.** FUTURE adds extensibility and multi-user layers. Each FUTURE item requires the 4 promotion criteria (ROADMAP §6.9): clear problem statement; design sketch / minimal ADR; pilot user; success metric. Without these four, item stays in BACKLOG.

### 6.1 Multi-language support

Per-language graph adapters beyond Go: Python (`ast` module + `pip show`), TypeScript (`tsc --listFiles` + tsconfig `paths`), Rust (`cargo metadata` + `cargo tree`), Java (Maven/Gradle), C# (.NET). Universal LSP fallback (ADR-013) covers all of them at lower quality meanwhile. See BACKLOG E-001..E-008.

### 6.2 Multi-user mode

Per-author attribution: `author: <id>` field in every session/decision. Conflict resolution rules for `ease_factor` (ADR-016 future). `caf report --team`. SSO via GitHub/GitLab OAuth. Triggers ADR-006 voice change ("As Anna, I chose...").

### 6.3 Web dashboard (visualizations only)

Read-only view on `.caf/`. **Interactive graph visualization (D3.js) — единственное, что нельзя в TUI.** Real-time heatmap. Export PDF/Markdown. **NOT replacing TUI** (ADR-014 invariant: TUI-first, web UI additive). See ADR-025 (future — web UI scope).

### 6.4 LLM-assisted features

`--ai-prefill` for `caf decisions new` (BACKLOG D-001). Smart card generation (D-002). Decision auto-summarization (D-003). All opt-in, LLM-free default for privacy (ADR-017 future — LLM backend selection).

### 6.5 CRDT for real-time collaboration

`yjs` or `automerge` for `cards.yml`. Auto-merge on `git pull`. Triggers ADR-001 mitigation for merge conflicts. See BACKLOG G-002.

> **OryaAgent architectural role.** OryaAgent is the in-house reference consumer of `caf` (built by the same team). In v0.3.1-triz it is **not** in the Tier 3 adapter list (was incorrectly listed there in v0.3.0 — that was the cause of the "12 vs 15" miscount). Architecturally, OryaAgent validates the universal LSP fallback path on a real production agent before third-party adapters ship; bugs found by OryaAgent feed back into ADR-013 and INV-ARCH-003.

### 6.6 Plugin API

Custom card types (BACKLOG H-001). `caf adapters register <path-to-binary>` (BACKLOG H-002). Custom report formats (BACKLOG H-003). Plugin marketplace like Homebrew taps (BACKLOG H-004). See ADR-018 (future), ADR-021 (future), ADR-028 (future).

### 6.7 IDE plugins

VSCode (inline provenance tooltips, "Why is this line here?" hover), JetBrains, Neovim. See BACKLOG B-001, B-002, B-003 (and K-003/4/5 cross-theme duplicates — consolidated in v0.3.1-triz).

### 6.8 CI/CD integrations

GitHub Action `caf report --check` as PR status check (BACKLOG C-001). GitLab CI (BACKLOG C-002). Slack daily `caf report --day` (BACKLOG C-003).

### 6.9 Exit criteria for FUTURE items

Each FUTURE item must satisfy these four promotion criteria before it can leave the BACKLOG:

1. **Clear problem statement** — what currently doesn't work or is inefficient.
2. **Design sketch / minimal ADR** — at least 3-4 lines of how it could work; an ADR if architectural.
3. **Pilot user** — someone who will use the feature within 2 weeks of shipping.
4. **Success metric** — concrete number, not "feels better".

> **TRIZ-pass note (v0.3.1-triz):** BACKLOG.md originally cited `ROADMAP.md §6.10` three times for these criteria. The actual criteria live here in §6.9 (there is no §6.10). All citations in BACKLOG.md have been corrected to point to this section.

"Без этих четырёх пунктов item остаётся в BACKLOG."

---

## 7. Cross-phase Concerns

### 7.1 Testing strategy

| Test layer | Coverage target | Phase | What it tests |
|------------|-----------------|-------|---------------|
| Unit | ≥80% | MVP | All non-trivial functions; pure logic |
| TUI snapshot tests | Per screen | MVP | Visual regression — see UI.md §14.4 |
| Adapter contract tests | All 25 × 3 events = 78 cases | MVP (12 cases), ALPHA (+9), BETA (+24), PROD (+33) | Per-adapter capability matrix |
| Integration `.caf/` round-trip | All YAML files | ALPHA | Write → read → diff = 0 |
| E2E 1 per adapter smoke | 25 adapters | BETA | Real adapter invocation, real `.caf/` |
| Performance benchmarks | First paint, frame render, keystroke, file op | MVP | UI.md §0.3 budgets; alert on >2x regression |
| **Property-based tests for SM-2** | ease ∈ [1.3, 3.0] invariant | BETA | INV-ARCH-006 enforcement |
| Terminal compat tests | 8 terminals | BETA | UI.md §12.1 compat matrix |

### 7.2 Documentation strategy

| Doc | Phase | Owner | What it covers |
|-----|-------|-------|----------------|
| README | MVP | maintainers | Quick start, install, 5-minute tour |
| ARCHITECT (this file) | MVP | maintainers | Architectural layers, capability matrix, invariants, fail-safe contract |
| UI | MVP | maintainers | All TUI screens, ASCII mockups, keyboard maps, performance budgets |
| ROADMAP | MVP | maintainers | Phases, deliverables, exit criteria, success metrics |
| BACKLOG | MVP | community | Catalog of ideas not yet committed; promotion criteria |
| ADR | MVP | maintainers | Architectural Decision Records — see ADR.md |
| Adapter dev guide | ALPHA | adapter contributors | How to write a new adapter (5-method interface, capability matrix) |
| User guide | ALPHA | maintainers | End-user docs; onboarding flow; common workflows |
| MCP integration guide | PROD | maintainers | How to plug `caf` into MCP-aware agents |
| Migration guides | per breaking change | maintainers | `caf migrate v0.2-to-v0.3`, etc. |

### 7.3 Release strategy

- semver `0.x.y` pre-PROD; `1.0.0` at PROD.
- `CHANGELOG.md` from conventional commits.
- GitHub Releases with Linux/macOS/Windows binaries.
- Pre-release artifacts tagged `0.x.y-rc.N`.

### 7.4 Backwards compatibility

- `.caf/` schema is additive across major versions (INV-ARCH-007, INV-ADR-001): new fields are optional; old fields are never removed within a major version.
- Adapter interface (5 methods) stable within a major version (INV-ARCH-002): new methods can be added as default-implemented; existing methods' signatures never change.
- `caf migrate` for breaking changes between major versions.
- **Config schema (referenced from ADR-004).** `.caf/config.yml` schema — keys: `triviality_threshold_loc` (default 50, ADR-004), `ui.heatmap_refresh_seconds` (default 30, INV-ARCH-005 sub-budget), `ui.theme` (dark/light/auto/high-contrast, ADR-014), `ui.language` (RU/EN, ADR-006 default), `decisions.autocomplete` (aggressive/conservative/off, ALPHA), `review.max_cards` (default 15, BETA), `lsp.strict` (true/false, MVP), `llm.backend` (PROD, future ADR-017), `telemetry.enabled` (default false, BACKLOG J-004). All keys are optional with sensible defaults; missing keys = use default. Schema version field: `version` (current `0.3.1-triz`). Migration path: `caf migrate v0.2-to-v0.3` adds `agent_tier`, `auto_detected`, `ui` sections to existing `.caf/config.yml`.

---

## 8. Phase Transitions

### 8.1 MVP → ALPHA

**Trigger:** MVP exit criteria met; universal fallback works for ≥3 different agents; ≥1 external early adopter using `caf` ≥2 weeks; adapter contract suite stable (12 cases pass on 4 adapters).

**Migration:** `caf init` adds `decisions/` directory to existing `.caf/` folder. `caf migrate v0.2-to-v0.3` for users coming from v0.2.

### 8.2 ALPHA → BETA

**Trigger:** ALPHA exit criteria met; Decision Coverage ≥90% on 2+ pilot repos; TUI decision editor avg ≤30s; ≥5 external users.

**Migration:** `caf init` adds `mental-model/` directory. First `caf review --regenerate --since 7d` populates initial cards.

### 8.3 BETA → PROD

**Trigger:** BETA exit criteria met; Mental Model Reach ≥70% on 3+ pilot repos; Recall Rate ≥75% sustained 4 weeks; all 25 explicit adapters shipped; `caf.first_person_explain` produces factually accurate narratives.

**Migration:** `caf init` adds MCP config (`.caf/mcp.yml`). `caf sync` enables cross-repo tracking.

### 8.4 PROD → FUTURE

**Trigger:** PROD 1.0 released; 6 months of stable operation; ≥50 active users; FUTURE items prioritized via community process (4 promotion criteria, §6.9).

---

## 9. Resourcing (rough estimates)

| Phase | Maintainer hours/week | Community PRs expected | Architectural components needing most attention |
|-------|-----------------------|------------------------|-------------------------------------------------|
| MVP | 15-20 | 0-2 | Storage layer (ADR-001, ADR-005); universal LSP fallback (ADR-013); TUI dashboard (ADR-014) |
| ALPHA | 20-25 | 2-5 | Decision Log schema (ADR-001); autocomplete pre-built indexes (UI.md §5.3); pre-commit validation (ADR-003, ADR-004) |
| BETA | 25-30 | 5-10 | SM-2 implementation (ADR-002); property-based tests (INV-ARCH-006); TUI review session (UI.md §4); graph builder on >10k files |
| PROD | 20-25 (stabilization) | 10-15 | MCP server (ADR-006); LLM grounding for `caf.first_person_explain` (INV-ADR-006); cross-repo sync; success metrics |
| FUTURE | variable | 15+ | Per FUTURE item |

Full-time engagement (40h/week) compresses timelines ~2x.

**🔥 80/20 implication:** high-leverage components (universal hook, LSP fallback, dashboard) — small effort, big value. Most of MVP budget should go to these, not to custom adapters for Tier 2/3 agents (those are ALPHA/BETA/PROD work).

---

## 10. Open Questions

Six numbered open architectural questions, to be resolved via ADRs:

1. **(MVP)** `caf init --auto-detect` — offer all detected adapters with recommended marker, or auto-install recommended only?
2. **(MVP)** Default LSP server for Go: `gopls` — bundle in `caf` binary or require install? See ADR-026 (future).
3. **(ALPHA)** Default `triviality_threshold_loc`: 50 (see ADR-004). Confirmed.
4. **(BETA)** LLM for `first_person_explain`: configurable, `z-ai-web-dev-sdk` as default. See ADR-017 (future).
5. **(PROD)** `caf sync` with private repos: yes, via git credential helper. Confirmed.
6. **(FUTURE)** Multi-user backward-compat: yes, optional `author` field. See BACKLOG G-001, ADR-016 (future).

"Эти вопросы будут решены через ADRs." Each accepted ADR retires one open question.

---

## 11. v0.2 → v0.3 Migration Notes

### 11.1 What changed

| Aspect | v0.2 | v0.3 | Architectural reason |
|--------|------|------|----------------------|
| Agent count | 9 explicit | 25 explicit + universal LSP | Coverage; ADR-012, ADR-013 |
| Default install (no-args) | nothing | universal hook + LSP 🔥 | "Provenance by default" principle; ADR-003, ADR-013 |
| Adapter tiers | 3 | 4 (added Tier 4 universal LSP) | Long-tail coverage; ADR-012, ADR-013 |
| Prioritization | feature-parity | 🔥 80/20 | Concentrated investment; ADR-015 |
| UI | CLI text | TUI-first | Universal SSH/Docker/CI; ADR-014 |
| Pre-commit | per-adapter | 🔥 universal | One hook for all agents; ADR-003 |
| Architecture file (this file) | (implicit, scattered across ADRs) | explicit, with capability matrix + invariants | Agent readability — see YAML front-matter |
| ROADMAP file | (separate, with same structure as this file but different content) | (separate, with same structure but roadmap-focused content) | TRIZ-pass: ARCHITECT and ROADMAP are no longer byte-identical duplicates |

### 11.2 Migration path

- `.caf/` schema is additive (`agent_tier`, `auto_detected`, `ui` config section added; existing fields preserved).
- v0.3 `caf init` no-args installs universal hook.
- All v0.2 commands still work (aliases preserved).
- New commands: `caf dashboard`, `caf heatmap`, `caf tour`, `caf adapters list` (TUI versions).

### 11.3 What's deprecated in v0.3

Nothing. v0.3 is purely additive. (v0.2 commands work as aliases to v0.3 commands; output format may differ — TUI vs text.)

---

## 12. Success Measurement

### 12.1 Per-phase metrics

| Phase | Primary metric | Target | Architectural invariant enforced |
|-------|----------------|--------|----------------------------------|
| MVP | Universal hook fires on 3+ different agents | 3+ agents in 7-day pilot | INV-ARCH-004 (fail-safe), INV-ARCH-003 (LSP fallback) |
| ALPHA | Decision Coverage ≥90% (last 7 days) | 90% | INV-ADR-004 (50 LOC threshold) |
| BETA | Recall Rate sustained 4 weeks ≥75% | 75% | INV-ARCH-006 (ease ∈ [1.3, 3.0]) |
| PROD | All 6 success metrics met (see §5.3) | all 6 | INV-ADR-006 (LLM no hallucination) |
| FUTURE | Per item (4 promotion criteria, §6.9) | per item | (per item) |

### 12.2 Adoption metrics (post-PROD)

| Metric | Target (12 months post-1.0) |
|--------|------------------------------|
| GitHub stars | ≥500 |
| Weekly active users | ≥100 |
| Community-contributed adapters | ≥5 new |
| Languages supported | ≥3 (Go + 2) |
| Issues closed | ≥80% within 30 days |

---

## Trace Matrix (architectural — TRIZ-pass addition)

| Architectural layer | ADR | Code path | UI screen | Invariant |
|---------------------|-----|-----------|-----------|-----------|
| Storage (`.caf/`) | ADR-001, ADR-005 | `internal/caf/storage/` | UI.md §3 Dashboard (Cross-repo widget) | INV-ARCH-001, INV-ARCH-007 |
| Hooks (universal + adapter) | ADR-003, ADR-010 | `internal/caf/hooks/`, `internal/caf/watch/` | UI.md §7 Adapters manager | INV-ARCH-004 |
| Adapter pattern | ADR-009, ADR-012 | `internal/caf/adapter.go`, `internal/adapters/` | UI.md §7 Adapters manager | INV-ARCH-002 |
| Universal LSP fallback (Tier 4) | ADR-013 | `internal/caf/lsp/fallback.go` | UI.md §7.4 Health check view | INV-ARCH-003 |
| Triviality threshold | ADR-004 | `internal/caf/triviality/threshold.go` | (none — internal) | INV-ADR-004 |
| SRS (SM-2) | ADR-002 | `internal/caf/srs/sm2.go` | UI.md §4 Review Session | INV-ARCH-006 |
| MCP server | ADR-006 (first-person) | `internal/caf/mcp/server.go` | (none — server-side) | INV-ADR-006 |
| TUI | ADR-014 | `pkg/ui/` | UI.md §0-§15 | INV-ARCH-005 |
| 80/20 prioritization | ADR-015 | (cross-cutting) | UI.md §0.3 Performance budget | (cross-cutting) |
| Multi-language adapters | ADR-011 | `internal/caf/lang/` | (none in MVP) | INV-ARCH-002 (LanguageAdapter interface) |
| Fail-safe contract | ADR-003 (INV-ADR-003), ADR-010 | `internal/caf/hooks/fail_safe.go` | (none — internal) | INV-ARCH-004 |
| Config schema | (within this §7.4) | `internal/caf/config/config.go` | UI.md §11 Settings Editor | (within §7.4) |
