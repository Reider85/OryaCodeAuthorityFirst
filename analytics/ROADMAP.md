---
doc: ROADMAP
layer: roadmap
version: 0.3.1-triz
language: bilingual
status: stable
audience: [agents, maintainers, early-adopters, anyone planning a feature]
related: [ARCHITECT.md, UI.md, BACKLOG.md, ADR.md]
depends_on: [ADR-015]
invariants:
  - {id: INV-RM-001, statement: "Each phase has exit criteria expressed as a measurable number, not a feeling", criticality: high, enforced_by: "caf report --metrics"}
  - {id: INV-RM-002, statement: "Phase transitions require all exit criteria met + at least one external early adopter", criticality: medium, enforced_by: "manual review"}
  - {id: INV-RM-003, statement: "FUTURE items must satisfy the 4 promotion criteria in §6.9 before leaving BACKLOG", criticality: high, enforced_by: "quarterly review"}
agent_routing:
  - {task: "Plan a release", see: [§7.3 Release strategy]}
  - {task: "Decide phase transitions", see: [§8 Phase Transitions]}
  - {task: "Promote a backlog item", see: [§6.9 Exit criteria for FUTURE items]}
  - {task: "Estimate resourcing", see: [§9 Resourcing]}
  - {task: "Track success metrics", see: [§12 Success Measurement]}
trace_to_code:
  - "(ROADMAP is product-planning; no direct code paths. Code is mapped via ARCHITECT.md trace matrix.)"
ui_budgets_referenced: "UI.md §0.3 (universal performance budget)"
---

# CodeAuthorityFirst — ROADMAP.md

> **Версия / Version:** 0.3.1-triz (post-TRIZ pass: this file is no longer a byte-identical duplicate of ARCHITECT.md; content rewritten to focus on product roadmap — deliverables, milestones, exit criteria, adoption metrics. Architectural concerns moved to ARCHITECT.md; capability matrix now lives in ARCHITECT.md §1.)
> **Формат фаз / Phase format:** MVP → ALPHA → BETA → PROD → FUTURE
> **Принцип приоритезации / Prioritization principle:** 🔥 80/20 — 20% усилий даёт 80% ценности; high-leverage фичи получают приоритет. See [ADR-015](./ADR.md#adr-015-8020-prioritization-principle).
> **Язык / Language:** bilingual — Russian body, English summary per phase
> **Аудитория / Audience:** tool maintainers + early adopters + anyone planning when a feature will ship
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md) · [UI.md](./UI.md) · [BACKLOG.md](./BACKLOG.md) · [ADR.md](./ADR.md)

> **How this file relates to ARCHITECT.md.** ROADMAP answers "what product capabilities ship in each phase, by when, with what measurable exit criteria". ARCHITECT answers "what architectural layers and code contracts ship in each phase, what must remain invariant". The two files share the same 12-section structure (per user constraint v0.3.1-triz: preserve structure 1:1) but have different content flavor. ROADMAP focuses on user-facing deliverables, exit criteria, adoption metrics; ARCHITECT focuses on invariants, interfaces, code paths.

---

## 0. Roadmap Philosophy (v0.3)

**EN summary.** v0.3 roadmap applies 80/20 prioritization: high-leverage features (universal pre-commit hook, auto-detect, LSP fallback, TUI dashboard) get disproportionate resources. Perfectionism on low-value features is forbidden. Each phase is self-sufficient — after each, the tool is more useful than before.

Дорожная карта v0.3 построена на четырёх принципах:

1. **🔥 80/20 priority.** Фичи помеченные `🔥 80/20` получают приоритет в ресурсах. Это фичи, дающие 80% ценности за 20% усилий. Перфекционизм на low-value фичах запрещён (см. [ADR-015](./ADR.md#adr-015-8020-prioritization-principle)).

2. **Self-sufficiency.** Каждая фаза добавляет ценность автономно. После MVP инструмент уже полезен, даже если development остановится. После ALPHA — ещё полезнее. И так далее.

3. **Universal before specific.** Сначала universal fallback (LSP, git hook, filesystem watcher) — покрывает ВСЕХ агентов сразу. Потом явные адаптеры для enhanced UX. Это обратно v0.2 подхода, где specific adapters were primary.

4. **TUI before Web.** TUI покрывает 80% UX ценностью 20% effort. Web UI отложен до FUTURE.

### 0.1 High-leverage features list (🔥 80/20)

| # | Feature | Effort | Value | Why 80/20 | Phase |
|---|---------|--------|-------|-----------|-------|
| 1 | Universal pre-commit hook | 2d | 90% coverage of agents | One hook fires for all 25+ agents | MVP |
| 2 | Auto-detect (`caf init --auto-detect`) | 1d | 5min→30s onboarding | Removes friction; user has nothing to configure | MVP |
| 3 | Universal LSP fallback (ADR-013) | 5d | Long-tail coverage | Zero-config support for all LSP languages | MVP |
| 4 | TUI dashboard | 3d | Hub for all screens | Replaces 6 separate CLI commands | MVP |
| 5 | `caf.first_person_explain` MCP tool | 5d | Killer feature for agents | LLM gets narrative-grade context | PROD |
| 6 | Single Go binary | 1d | Zero-dep distribution | One file, no runtime | MVP |
| 7 | ASCII heatmap | 1d | One-screen overview | Replaces 100-line text report | BETA |
| 8 | SM-2 algorithm | 2d | Battle-tested SRS | Anki users recognize it immediately | BETA |

**Итого:** 8 фич, ~20 дней effort, ~80% ценности. Это и есть MVP+ALPHA+BETA ядро. См. [ADR-015](./ADR.md#adr-015-8020-prioritization-principle).

---

## 1. Phase Overview (v0.3)

**EN summary.** v0.3 ships in five phases: MVP (2-3 wk, universal coverage), ALPHA (4-6 wk, decisions + autocomplete), BETA (6-8 wk, mental model + SRS), PROD (8-10 wk, MCP + cross-repo + LLM), FUTURE (open-ended, multi-user + web UI + plugins). The phase list is the same as in ARCHITECT.md §1 (by design — same 12-section structure); the difference is that ROADMAP §1 lists user-facing features per phase, while ARCHITECT §1 lists architectural components per phase.

| Phase | Duration (est.) | Primary goal | Modules live | Agents supported | UI screens | 🔥 80/20 features |
|-------|-----------------|--------------|--------------|-------------------|------------|---------------------|
| MVP | 2-3 wk | Universal coverage | Provenance, context.md, dashboard | 25+ via universal LSP + Tier 1 (2 explicit) | Dashboard, Onboarding, Adapters manager | Universal pre-commit hook, Auto-detect, LSP fallback, TUI dashboard, Single Go binary |
| ALPHA | 4-6 wk | Decision Log + pre-commit validation | Decision Log, TUI decisions list/editor/detail, pre-commit validation, +3 Tier 1 adapters (Cursor, Cline, KiloCode) | 25+ via universal + 5 enhanced (Tier 1) | Decisions list, Decision detail, Decision editor | 80/20 autocomplete (90s→25s) |
| BETA | 6-8 wk | Mental Model Reconstructor + SRS | MMR, SM-2 🔥, 3 card types, dependency graph, blindspots detection, TUI review 🔥, TUI heatmap, TUI file tour | 25+ via universal + 13 enhanced (+8 Tier 2) | Review session, Heatmap, File tour | SM-2 algorithm, ASCII heatmap, TUI review session |
| PROD | 8-10 wk | MCP server + cross-repo sync + LLM | MCP server (5 tools + 3 prompts), cross-repo sync, TUI reports view, success metrics, `caf.first_person_explain` 🔥, +12 Tier 3 adapters | 25+ via universal + 25 enhanced (all explicit) | Reports view, Metrics report | `caf.first_person_explain` MCP tool, cross-repo sync |
| FUTURE | open-ended | Extensibility + multi-user | Multi-user, Web UI (read-only), CRDT, Plugin API, IDE plugins, CI/CD integrations | (additive) | Web dashboard (read-only), VSCode/JetBrains/Neovim plugins | (per-item 4 promotion criteria — see §6.9) |

> **🔥 Key shift from v0.2:** MVP теперь покрывает 25+ агентов через universal fallback, а не 2 явных адаптера. Это радикально расширяет initial coverage.

> **Capability matrix** lives in [ARCHITECT.md §1 Capability matrix](./ARCHITECT.md#capability-matrix-first-class-artifact-referenced-from-adr-009-and-adr-013) — the canonical table showing which adapter is in which tier and which canonical events each supports natively vs via fallback.

---

## 2. Phase: MVP (2-3 weeks)

**EN summary.** MVP delivers universal coverage: `.caf/` storage, fail-safe hooks (universal pre-commit fires for ALL 25+ agents), auto-detect, TUI dashboard, single Go binary. After MVP, any agent gets provenance + decisions infrastructure, even without an explicit adapter.

### 2.1 Scope

User-facing features in MVP:

- **Provenance Tracker (partial — without decision linking):** records who/when/what for each agent session, no manual `caf log` calls.
- **`context.md` generator:** produces a markdown summary of recent sessions for agents to consume as context.
- **🔥 Universal pre-commit hook:** one shell script in `.git/hooks/pre-commit`, fires for all 25+ agents.
- **🔥 Universal LSP fallback adapter (Tier 4):** zero-config support for any LSP-supporting language (Go first via gopls).
- **🔥 Auto-detect (`caf init --auto-detect`):** scans for installed agents, recommends with marker.
- **opencode adapter (Tier 1 native):** explicit enhanced adapter for opencode.
- **Claude Code adapter (Tier 1 native):** explicit enhanced adapter for Claude Code.
- **🔥 TUI dashboard:** hub screen — Today widget (review queue, pending decisions, new blindspots), Authorship Score widget (ACS weekly avg, target ≥3.5, Recall %, Coverage %), Active Agents widget, Recent Decisions widget, Cross-repo widget (dependencies with stale-decision indicators).
- **TUI onboarding wizard:** first-run experience, ~30 seconds, 4-step flow.
- **TUI adapters manager:** lists installed + available adapters, health check view.

Adapters shipped (MVP):
- Universal LSP fallback (Tier 4)
- Universal pre-commit hook
- opencode (Tier 1 native, all events)
- Claude Code (Tier 1 native, all events)

Coverage: 25+ agents covered (via universal fallback) + 2 enhanced (Tier 1).

CLI commands shipped (MVP):
- `caf init` — interactive onboarding
- `caf init --auto-detect` — auto-detect agents
- `caf init --agent opencode` — install specific adapter
- `caf dashboard` — TUI dashboard
- `caf adapters list` — TUI adapters manager
- `caf who <file:line>` — provenance for line
- `caf when --pattern <regex> <dir>` — sessions by file pattern
- `caf sessions --since 7d` — list sessions
- `caf context` — generate `context.md`
- `caf hook pre-session` — invoked by adapter on agent start
- `caf hook post-edit` — invoked after non-trivial file changes
- `caf hook pre-commit` 🔥 — universal, fires in `.git/hooks/pre-commit`

Three hook types: pre-session, post-edit, pre-commit. See [ARCHITECT.md §2.1](./ARCHITECT.md#21-scope) for the architectural contract.

### 2.2 Deliverables

User-visible deliverables in MVP:

1. `caf` single binary for Linux/macOS/Windows.
2. `.caf/` folder structure created by `caf init`.
3. `.git/hooks/pre-commit` universal hook installed.
4. TUI dashboard with 5 widgets.
5. TUI onboarding wizard.
6. TUI adapters manager with health check.
7. `context.md` auto-generated.
8. Two Tier 1 explicit adapters (opencode, Claude Code) + Universal LSP fallback.

### 2.3 Exit criteria

1. Pilot repo uses `caf` 7 days straight.
2. Universal pre-commit hook fires on 3+ different agents (opencode, Claude Code, +1 unknown via fallback).
3. LSP fallback builds symbol graph for Go repo up to 1000 files.
4. **TUI dashboard first paint <100ms** on pilot repo.
5. Auto-detect finds opencode + Claude Code on fresh machine in <5s.
6. Adapter contract test suite passes (12 cases for 4 adapters).

### 2.4 What's NOT in MVP

- Decision Log (ALPHA)
- MMR / SRS (BETA)
- MCP (PROD)
- Tier 1 beyond opencode+Claude Code (ALPHA)
- Tier 2/3 (BETA/PROD)
- Cross-repo sync (PROD)
- LLM features (PROD)

### 2.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| LSP fallback noisy (false-positive symbol graphs) | high | medium | `caf doctor --verify-graph` rebuilds; `caf config set lsp.strict false` relaxes |
| fs watcher CPU on large repos | medium | medium | Throttle: poll every 1s, debounce 5s, ignore `.git/`; daemon exits after 30 min inactivity |
| TUI on weird terminals | medium | low | UI.md §12.1 compat matrix; 16-color fallback via termenv |
| `caf init` confusing for new users | medium | high | Onboarding wizard (UI.md §2); auto-detect default; 3-step flow <30s |
| Universal pre-commit blocks legit commits | medium | high | Fail-safe contract: `caf` binary missing → warn + let through; see [ARCHITECT.md §2.5](./ARCHITECT.md#25-risks) for the full fail-safe contract |

---

## 3. Phase: ALPHA (4-6 weeks)

**EN summary.** ALPHA delivers the Decision Log layer: `.caf/decisions/*.md` schema, TUI decisions list/detail/editor, pre-commit validation requiring decisions for non-trivial commits, three more Tier 1 adapters (Cursor, Cline, KiloCode). The 🔥 80/20 win: autocomplete for files/tags/agent reduces fill time from 90s to 25s.

### 3.1 Scope

User-facing features in ALPHA:

- **Decision Log:** `.caf/decisions/*.md` records — title, date, agent, files[], context, chosen, alternatives, implications.
- **TUI decisions list:** filterable table (Date/Title/Agent/Files). Filters: today/week/month/all, by agent/file/status. Actions: `Enter` open, `n` new, `/` search, `g` gaps.
- **TUI decision detail view:** markdown-rendered record with sections `## Контекст / ## Выбрано / ## Альтернативы / ## Импликации`, related sessions, related cards. Actions: `e` edit, `d` delete, `s` show file, `Esc` back.
- **TUI decision editor 🔥:** autocomplete for files/tags/agent reduces time 90s→25s. Pre-fills from session.
- **Pre-commit validation:** `caf hook pre-commit` blocks non-trivial commits (≥50 LOC, see ADR-004) without a decision or `CAF-Trivial: true` trailer. Applies to ALL 25+ agents via universal fallback chain.
- **Cursor adapter (Tier 1, pre-session via wrapper):** explicit enhanced adapter.
- **Cline adapter (Tier 1 native, all events):** explicit enhanced adapter.
- **KiloCode adapter (Tier 1 native, all events):** explicit enhanced adapter.

Adapters shipped (cumulative, ALPHA): 4 (MVP) + Cursor + Cline + KiloCode = 7. Coverage: 25+ via universal + 5 enhanced (Tier 1).

New CLI commands (ALPHA):
- `caf why <file>` — provenance + (eventually) explanation
- `caf decisions` 🔥 — TUI decisions list
- `caf decisions new [--files ...]` 🔥 — TUI decision editor with autocomplete
- `caf decisions --today` / `--file <path>` / `--search <query>` / `--gaps` / `--agent <id>` — list filters

### 3.2 Deliverables

1. `.caf/decisions/` folder created by `caf init` (additive to MVP).
2. TUI decisions list, detail, editor.
3. Pre-commit validation for non-trivial edits.
4. Cursor, Cline, KiloCode adapters.
5. `caf migrate v0.2-to-v0.3` migration path for users coming from v0.2.

### 3.3 Exit criteria

1. Decision Coverage ≥90% over last 7 days on pilot repo.
2. Pre-commit blocks non-trivial commits without decision for ALL 25+ agents.
3. **TUI decision editor avg fill time ≤30s (vs 90s without autocomplete)**.
4. All Tier 1 adapters work on macOS + Linux with no manual config.

### 3.4 What's NOT in ALPHA

- Mental Model Reconstructor (BETA)
- Tier 2/3 adapters (BETA/PROD)
- Cross-repo sync (PROD)
- LLM features (PROD)

### 3.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Autocomplete too noisy | medium | medium | `caf config set decisions.autocomplete aggressive\|conservative\|off` |
| Pre-commit blocks CI commits | high | high | Detect `CI=true` env, skip pre-commit validation; CI runs `caf report --check` separately (BACKLOG C-001) |
| Decision template too verbose | high | high | Default compact template; `--verbose` for full template |
| Cursor wrapper breaks on Cursor version update | medium | medium | Pin wrapper version; auto-detect breaks, surface in `caf adapters list` |
| Inconsistent UX between adapters | medium | low | Capability matrix in ARCHITECT.md §1 makes limitations explicit |

---

## 4. Phase: BETA (6-8 weeks)

**EN summary.** BETA delivers the Mental Model Reconstructor: SM-2 spaced-repetition algorithm (🔥 80/20), three card types (`why`, `what_if_remove`, `where_handled`), single-repo dependency graph, blindspots detection, three new TUI screens (review session, heatmap, file tour). Eight Tier 2 adapters join the matrix.

### 4.1 Scope

User-facing features in BETA:

- **Mental Model Reconstructor (MMR):** generates SRS cards for non-trivial code regions.
- **🔥 SM-2 algorithm (ADR-002):** Anki-compatible spaced repetition. ease ∈ [1.3, 3.0] invariant.
- **Three card types:** `why` (Type 1, fact recall), `what_if_remove` (Type 2, abstract reasoning — predict impact if removed), `where_handled` (Type 3, trace through code from entry point to handler).
- **Single-repo dependency graph:** `go list -json` (Go-only in BETA). Stored in `.caf/mental-model/graph.yml`. Other languages via Tier 4 LSP fallback.
- **Blindspots detection:** files with low ACS or low coverage, sorted by risk score.
- **🔥 TUI Review Session (`caf review`):** daily SRS session, ~10 minutes, 10-15 cards. Split-view layout: question (left) + context (file path, decision ref, sessions) + grade panel.
- **TUI Heatmap (`caf heatmap`):** each file = one cell, color = ACS or coverage. Pre-computed, cached.
- **TUI File Tour (`caf tour <file>`):** split-view: timeline (left) + numbered source (right).

8 Tier 2 adapters (cumulative, BETA): Aider, Codex CLI, Continue.dev, GitHub Copilot CLI, Gemini CLI, Amazon Q Developer CLI, Goose, Windsurf. Coverage: 25+ via universal + 13 enhanced.

New CLI commands (BETA):
- `caf review` 🔥 — full-screen TUI review session
- `caf review --hard` / `--blindspots` / `--regenerate --since 7d` / `--stats`
- `caf heatmap` — TUI heatmap
- `caf tour <file>` — TUI file tour
- `caf sync` — single-repo graph builder

### 4.2 Deliverables

1. SM-2 SRS engine with three card types.
2. Mental Model Reconstructor.
3. Dependency graph builder (Go-only; FUTURE adds other languages).
4. Blindspots detection.
5. Three TUI screens: review session, heatmap, file tour.
6. Eight Tier 2 adapters.
7. Property-based tests for SM-2 (ease ∈ [1.3, 3.0]).

### 4.3 Exit criteria

1. Mental Model Reach ≥70% (cards due / cards generated).
2. Recall Rate ≥75% sustained over 4 weeks.
3. ACS ≥3.0 (target ≥3.5 in PROD).
4. Review session completion ≥80% (cards answered / cards started).
5. All Tier 2 adapters work on macOS + Linux.
6. **Heatmap render <200ms on repo up to 1000 files**.

### 4.4 What's NOT in BETA

- MCP server (PROD)
- Cross-repo sync (PROD)
- `caf.first_person_explain` LLM feature (PROD)
- Tier 3 adapters (PROD)
- Web UI (FUTURE)

### 4.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| SM-2 burnout (too many cards per session) | high | high | Default cap 15 cards/session; `caf config set review.max_cards 10`; BACKLOG D-007 adaptive per-card-type ease |
| Cramped TUI on small terminals | medium | medium | UI.md §12.2 responsive: ≥140 cols full grid, 100-139 2-col, 60-99 single stacked, <60 compact mode |
| Low-quality cards | high | medium | BETA ships hand-curated cards only; LLM-generated cards (BACKLOG D-002) are FUTURE |
| Tier 2 wrappers fail | medium | medium | Pin wrapper versions; auto-detect breaks, surface in `caf adapters list` |
| Graph builder fails on >10k files | medium | high | Streaming + lazy loading (BACKLOG A-001); benchmark on `kubernetes/kubernetes` before release |

---

## 5. Phase: PROD (8-10 weeks)

**EN summary.** PROD delivers the MCP server layer (5 tools + 3 prompts), cross-repo sync automation, Tier 3 adapters (12 agents, wrapper-only), TUI reports view, success metrics instrumentation. The 🔥 killer feature: `caf.first_person_explain` MCP tool produces LLM-generated first-person narratives grounded in `.caf/` facts (no hallucination, ADR-006 invariant).

### 5.1 Scope

User-facing features in PROD:

- **MCP server (`caf-mcp`):** 5 tools + 3 prompts.
  - Tools: `caf.context`, `caf.provenance`, `caf.decisions`, `caf.mental_model`, `caf.first_person_explain` 🔥
  - Prompts: `caf.pre_session_brief`, `caf.code_review_with_provenance`, `caf.archaeology`
  - Transport: stdio (`caf mcp serve`) and HTTP SSE (`caf mcp serve --transport sse`)
- **🔥 `caf.first_person_explain` (ADR-006):** LLM-generated first-person narrative, grounded in `.caf/` facts. LLM only paraphrases, never invents.
- **Cross-repo sync automation:** `caf sync` detects locally-cloned repos via `~/.caf/repos.yml` registry; pulls `manifest.yml` from each dependee; creates `review-because-interface-changed` cards when interfaces change.
- **TUI Reports view:** daily report (`caf report --day`), week report (`caf report --week`), metrics report (`caf report --metrics`).
- **Success metrics instrumentation:** ACS, Decision Coverage, Mental Model Reach, Recall Rate, Agent Diversity (Shannon entropy), Adapter Coverage.
- **12 Tier 3 adapters (cumulative):** Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Augment, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep. Note (v0.3.1-triz): list reconciled to 12 — OryaAgent is in-house consumer (see ARCHITECT.md §6.5); Augment (not "Augment Code") consistent with vendor branding.

Tier 3 list (12 agents, wrapper-only):
Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Augment, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep.

Cumulative total: 25 explicit + universal. Coverage: 25+ via universal + 25 enhanced.

New CLI commands (PROD):
- `caf report --day` / `--week` / `--ownership` / `--agent` / `--metrics`
- `caf sync` — cross-repo sync
- `caf hook post-commit --repo <path>`
- `caf mcp serve [--transport sse]`

### 5.2 Deliverables

1. MCP server with 5 tools + 3 prompts.
2. `caf.first_person_explain` with strict grounding in `.caf/` facts.
3. Cross-repo sync automation, manifest cache (1h TTL).
4. TUI Reports view (daily/week/metrics).
5. Success metrics instrumentation.
6. 12 Tier 3 adapters (wrapper-only).
7. README, ARCHITECT.md, UI.md, ROADMAP.md (this file), BACKLOG.md, ADR.md all consistent with implementation.

### 5.3 Exit criteria

1. All 5 success metrics hit targets on pilot repos over 8 weeks:
   - **ACS ≥3.5**
   - **Decision Coverage ≥90%**
   - **Mental Model Reach ≥70%**
   - **Recall Rate ≥75%**
   - **Agent Diversity H ≥1.0**
   - **Adapter Coverage ≥70% native (NEW v0.3)**
2. MCP works with opencode, Claude Code, Cline, KiloCode, Goose.
3. Cross-repo sync detects changed interfaces between dependent repos.
4. All 25 explicit adapters pass contract tests.
5. `caf.first_person_explain` produces factually accurate narratives (no hallucination).
6. README, ARCHITECT.md, UI.md, ROADMAP.md, BACKLOG.md, ADR.md all consistent with implementation.

### 5.4 What's NOT in PROD

- Multi-user mode (FUTURE)
- Web UI (FUTURE)
- CRDT for `cards.yml` (FUTURE)
- Plugin API / marketplace (FUTURE)
- IDE plugins (FUTURE)
- CI/CD integrations beyond `caf report --check` (FUTURE)

### 5.5 Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| LLM hallucinates in `caf.first_person_explain` | high | high | Strict grounding in `.caf/` facts (ADR-006, INV-ADR-006); LLM only paraphrases; test suite verifies |
| MCP protocol changes (Anthropic updates spec) | medium | medium | Pin MCP spec version; contract tests against pinned version |
| Cross-repo sync slow on large fleets | medium | medium | Parallel fetch (10 concurrent, rate-limited); cache `manifest.yml` 1h; incremental sync (BACKLOG A-003) |
| Tier 3 brittle across 12 wrapper-only adapters | high | medium | Wrapper contracts minimal; failure in one doesn't block others; `caf adapters health` surfaces breakages |
| TUI performance degrades as data grows | medium | high | Performance budgets enforced (UI.md §0.3); CI benchmarks; alert on >2x regression; BACKLOG A-001 lazy graph loading |

---

## 6. Phase: FUTURE (open-ended)

**EN summary.** FUTURE items are not yet committed to a phase. Each must satisfy the 4 promotion criteria in §6.9 (clear problem; design sketch / minimal ADR; pilot user; success metric). Without these four, item stays in BACKLOG.

### 6.1 Multi-language support

Python (`ast` module + `pip show`), TypeScript (`tsc --listFiles` + tsconfig `paths`), Rust (`cargo metadata` + `cargo tree`), Java (Maven/Gradle), C# (.NET). Universal LSP fallback (ADR-013) covers all of them at lower quality meanwhile. Open questions: Python dynamic imports `__import__`/`importlib`; TS monorepo (yarn/pnpm workspaces); Rust `rust-analyzer` integration. See BACKLOG E-001..E-008.

### 6.2 Multi-user mode

Per-author attribution; conflict resolution rules for `ease_factor` (which wins: minimum or latest?); `caf report --team`; SSO via GitHub/GitLab OAuth. Triggers ADR-006 voice change ("As Anna, I chose..."). See BACKLOG G-001, ADR-016 (future).

### 6.3 Web dashboard (visualizations only)

Read-only view on `.caf/`. **Interactive graph visualization (D3.js) — единственное, что нельзя в TUI.** Real-time heatmap. Export PDF/Markdown. **NOT replacing TUI.** See ADR-025 (future — web UI scope), BACKLOG K-002.

### 6.4 LLM-assisted features

`--ai-prefill` for `caf decisions new` (BACKLOG D-001). Smart card generation (D-002). Decision auto-summarization (D-003). All opt-in, LLM-free default for privacy. LLM backend configurable via `caf config set llm.backend`. See ADR-017 (future — LLM backend selection).

### 6.5 CRDT for real-time collaboration

`yjs` or `automerge` for `cards.yml`. Auto-merge on `git pull`. Mitigates ADR-001 merge conflict risk. See BACKLOG G-002.

> **OryaAgent** is the in-house reference consumer of `caf` (built by the same team). It is not in the Tier 3 adapter list — it's the in-house production agent that validates the universal LSP fallback path on real workloads. See [ARCHITECT.md §6.5](./ARCHITECT.md#65-crdt-for-real-time-collaboration) for its architectural role.

### 6.6 Plugin API

Custom card types (BACKLOG H-001). `caf adapters register <path-to-binary>` (BACKLOG H-002). Custom report formats (BACKLOG H-003). Plugin marketplace like Homebrew taps (BACKLOG H-004). See ADR-018 (future), ADR-021 (future).

### 6.7 IDE plugins

VSCode (inline provenance tooltips, "Why is this line here?" hover). JetBrains. Neovim (`:CafWho` Lua command). See BACKLOG B-001, B-002, B-003.

### 6.8 CI/CD integrations

GitHub Action `caf report --check` as PR status check (BACKLOG C-001). GitLab CI (BACKLOG C-002). Slack daily `caf report --day` (BACKLOG C-003). Discord team mode (BACKLOG C-004). Linear/Jira integration (BACKLOG C-005). GitHub PR comment bot with `caf first_person_explain` (BACKLOG C-006).

### 6.9 Exit criteria for FUTURE items

Each FUTURE item must satisfy these four promotion criteria before it can leave the BACKLOG:

1. **Clear problem statement** — what currently doesn't work or is inefficient.
2. **Design sketch / minimal ADR** — at least 3-4 lines of how it could work; an ADR if architectural.
3. **Pilot user** — someone who will use the feature within 2 weeks of shipping.
4. **Success metric** — concrete number, not "feels better".

"Без этих четырёх пунктов item остаётся в BACKLOG."

> **TRIZ-pass note (v0.3.1-triz):** BACKLOG.md originally cited `ROADMAP.md §6.10` three times for these criteria. The actual criteria live here in §6.9 (there is no §6.10 — §6 only has 6.1 through 6.9). All citations in BACKLOG.md have been corrected to point to this section.

---

## 7. Cross-phase Concerns

### 7.1 Testing strategy

| Test layer | Coverage target | Phase |
|------------|-----------------|-------|
| Unit | ≥80% | MVP |
| TUI snapshot tests (UI.md §14.4) | Per screen | MVP |
| Adapter contract tests (all 25 + universal) | All 25×3=78 cases | MVP (4 adapters), ALPHA (+3), BETA (+8), PROD (+12) |
| Integration `.caf/` round-trip | All YAML files | ALPHA |
| E2E 1 per adapter smoke | 25 adapters | BETA |
| Performance benchmarks (UI.md §0.3) | First paint, frame render, keystroke, file op | MVP |
| **Property-based tests for SM-2: invariant `ease ∈ [1.3, 3.0]`** | 1000 random card generations | BETA |
| Terminal compat tests (8 terminals, UI.md §12.1) | All 8 | BETA |

### 7.2 Documentation strategy

| Doc | Phase | Owner |
|-----|-------|-------|
| README | MVP | maintainers |
| ARCHITECT | MVP | maintainers |
| UI | MVP | maintainers |
| ROADMAP (this file) | MVP | maintainers |
| BACKLOG | MVP | community |
| ADR | MVP | maintainers |
| Adapter dev guide | ALPHA | adapter contributors |
| User guide | ALPHA | maintainers |
| MCP integration guide | PROD | maintainers |
| Migration guides | per breaking change | maintainers |

### 7.3 Release strategy

- semver `0.x.y` pre-PROD, `1.0.0` at PROD.
- `CHANGELOG.md` from conventional commits.
- GitHub Releases with Linux/macOS/Windows binaries.

### 7.4 Backwards compatibility

- `.caf/` schema additive (INV-ARCH-007).
- Adapter interface stable within major version (INV-ARCH-002).
- `caf migrate` for breaking changes between major versions.

---

## 8. Phase Transitions

### 8.1 MVP → ALPHA

**Trigger:** MVP exit criteria met; universal fallback works for ≥3 different agents; ≥1 external early adopter using `caf` ≥2 weeks; adapter contract suite stable.

**Migration:** `caf init` adds `decisions/` dir; `caf migrate v0.2-to-v0.3`.

### 8.2 ALPHA → BETA

**Trigger:** ALPHA exit criteria met; Decision Coverage ≥90% on 2+ repos; TUI editor avg ≤30s; ≥5 external users.

**Migration:** `caf init` adds `mental-model/`; first `caf review --regenerate`.

### 8.3 BETA → PROD

**Trigger:** BETA exit criteria met; Mental Model Reach ≥70% on 3+ repos; Recall Rate ≥75% sustained 4 weeks; all 25 adapters shipped; `caf.first_person_explain` accurate.

**Migration:** `caf init` adds MCP config; `caf sync`.

### 8.4 PROD → FUTURE

**Trigger:** PROD 1.0 released; 6 months stable; ≥50 active users; FUTURE items prioritized via community.

---

## 9. Resourcing (rough estimates)

| Phase | Maintainer hours/week | Community PRs expected |
|-------|-----------------------|------------------------|
| MVP | 15-20 | 0-2 |
| ALPHA | 20-25 | 2-5 |
| BETA | 25-30 | 5-10 |
| PROD | 20-25 (stabilization) | 10-15 |
| FUTURE | variable | 15+ |

Full-time engagement (40h/week) compresses timelines ~2x.

**🔥 80/20 implication:** high-leverage features (universal hook, LSP fallback, dashboard) — small effort, big value. Most of MVP budget should go to these, not to custom adapters for Tier 2/3 agents.

---

## 10. Open Questions

1. **(MVP)** `caf init --auto-detect` — offer all detected adapters with recommended marker.
2. **(MVP)** Default LSP server for Go: `gopls` — bundle in `caf` binary or require install.
3. **(ALPHA)** Default `triviality_threshold_loc`: 50 (see ADR-004).
4. **(BETA)** LLM for `first_person_explain`: configurable, `z-ai-web-dev-sdk` as default.
5. **(PROD)** `caf sync` with private repos: yes, via git credential helper.
6. **(FUTURE)** Multi-user backward-compat: yes, optional `author` field.

"Эти вопросы будут решены через ADRs."

---

## 11. v0.2 → v0.3 Migration Notes

### 11.1 What changed

| Aspect | v0.2 | v0.3 |
|--------|------|------|
| Agent count | 9 explicit | 25 explicit + universal LSP |
| Default install (no-args) | nothing | universal hook + LSP 🔥 |
| Adapter tiers | 3 | 4 (added Tier 4 universal LSP) |
| Prioritization | feature-parity | 🔥 80/20 |
| UI | CLI text | TUI-first |
| Pre-commit | per-adapter | 🔥 universal |
| Documentation files | (intermingled, ARCHITECT=ROADMAP duplicate) | (separated; ARCHITECT focuses on architectural invariants, ROADMAP focuses on product deliverables — same 12-section structure, different content) |

### 11.2 Migration path

- `.caf/` schema additive (`agent_tier`, `auto_detected`, `ui` config section).
- v0.3 `caf init` no-args installs universal hook.
- All v0.2 commands still work.
- New: `caf dashboard`, `caf heatmap`, `caf tour`, `caf adapters list` (TUI).

### 11.3 What's deprecated in v0.3

Nothing. v0.3 is purely additive.

---

## 12. Success Measurement

### 12.1 Per-phase metrics

| Phase | Primary metric | Target |
|-------|----------------|--------|
| MVP | universal hook fires ≥3 agents | 3+ agents in 7-day pilot |
| ALPHA | Decision Coverage ≥90% | 90% (last 7 days) |
| BETA | Recall Rate sustained 4 weeks ≥75% | 75% |
| PROD | all 6 success metrics met | all 6 (see §5.3) |
| FUTURE | per item (4 promotion criteria, §6.9) | per item |

### 12.2 Adoption metrics (post-PROD)

| Metric | Target (12 months post-1.0) |
|--------|------------------------------|
| GitHub stars | ≥500 |
| Weekly active users | ≥100 |
| Community-contributed adapters | ≥5 new |
| Languages supported | ≥3 (Go + 2) |
| Issues closed | ≥80% within 30 days |
