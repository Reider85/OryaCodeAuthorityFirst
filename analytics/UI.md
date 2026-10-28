---
doc: UI
layer: ui
version: 0.3.1-triz
language: bilingual
status: stable
audience: [agents, UI/UX implementers, contributors adding TUI screens, design reviewers]
related: [ARCHITECT.md, ROADMAP.md, BACKLOG.md, ADR.md]
depends_on: [ADR-014, ADR-002, ADR-005, ADR-009, ADR-013, ADR-015]
invariants:
  - {id: INV-UI-001, statement: "TUI first paint <100ms on pilot repos; >2x CI regression triggers alert", criticality: high, enforced_by: "CI benchmark — pkg/ui/perf.go"}
  - {id: INV-UI-002, statement: "Frame render <16ms (60fps) for scroll/redraw; uses virtualized lists + cached layouts", criticality: high, enforced_by: "bubbles/list virtualization"}
  - {id: INV-UI-003, statement: "Keystroke response <50ms for simple actions; no async round-trip for typing", criticality: high, enforced_by: "tea.Cmd direct dispatch"}
  - {id: INV-UI-004, statement: "Every async data load returns Skeleton{} immediately, then DataLoadedMsg swaps to real data (SWR pattern)", criticality: high, enforced_by: "pkg/ui/components/skeleton.go pattern"}
  - {id: INV-UI-005, statement: "Screens mount lazily on first navigation (tea.Cmd returns LoadScreenMsg); closed screens cost zero render time (Islands/partial hydration analogue)", criticality: medium, enforced_by: "pkg/ui/screens/registry.go"}
  - {id: INV-UI-006, statement: "Lipgloss pre-computes style strings at app startup; only final render string emitted per frame (Critical CSS analogue)", criticality: medium, enforced_by: "Lipgloss API contract"}
  - {id: INV-UI-007, statement: "Heatmap render <200ms on 247-file repo; re-compute on file change <500ms incremental", criticality: high, enforced_by: "pre-computed + cached"}
owned_contracts:
  - "pkg/ui/tui.go — root model + Update(msg) + View() string"
  - "pkg/ui/screens/*.go — one file per screen, all follow Purpose/Layout/Components/Keyboard map/Performance budget/Edge cases schema"
  - "pkg/ui/components/ — header, statusbar, spinner, skeleton, heatmap"
  - "pkg/ui/keys.go — universal keyboard map"
  - "pkg/ui/perf.go — performance budget enforcer + CI benchmark hooks"
agent_routing:
  - {task: "Add new TUI screen", see: [§0 UI Philosophy, §1.3 Common UI patterns, §13.1 File organization], template: "Purpose/Layout/Components/Keyboard map/Performance budget/Edge cases"}
  - {task: "Optimize TUI performance", see: [§0.3 Performance budget, §13.2 State management, §13.3 Virtualization, §13.4 Caching], invariant: INV-UI-001 through INV-UI-007}
  - {task: "Add new keyboard shortcut", see: [§0.6 Universal keyboard map], conflict_check: "must not collide with universal map"}
  - {task: "Fix terminal compat issue", see: [§12.1 Terminal compatibility]}
  - {task: "Implement skeleton/SWR for new screen", see: [§3.5 Refresh strategy, §13.2 State management]}
trace_to_code:
  - "pkg/ui/tui.go — root model"
  - "pkg/ui/theme.go — color palette, typography"
  - "pkg/ui/components/{header,statusbar,spinner,skeleton,heatmap}.go"
  - "pkg/ui/screens/{onboarding,dashboard,review,decisions_list,decisions_detail,decisions_editor,heatmap,adapters,file_tour,reports,settings}.go"
  - "pkg/ui/keys.go"
  - "pkg/ui/perf.go"
ui_budgets:
  first_paint: "<100ms (any screen)"
  frame_render: "<16ms (60fps, scroll/redraw)"
  keystroke: "<50ms"
  llm_op: "<3s (with skeleton + spinner)"
  file_op_10k_repo: "<1s"
  dashboard_full_load: "<500ms"
  dashboard_auto_refresh: "<50ms (incremental, only changed widgets)"
  dashboard_memory: "<50MB"
  card_render: "<50ms (no LLM)"
  session_save: "<100ms (incremental write)"
  list_render: "<50ms (virtualized, only visible rows)"
  filter_response: "<16ms (in-memory index)"
  detail_view: "<100ms (markdown render via Glamour, cached)"
  autocomplete: "<50ms (pre-built index)"
  heatmap_render: "<200ms (247 files pre-computed)"
  heatmap_recompute: "<500ms (incremental)"
  heatmap_memory: "<20MB"
  adapter_health_check: "<5s/adapter (<30s/25 total)"
  adapter_install: "<2s/adapter (parallel ≤5)"
  file_tour_first_paint: "<100ms"
  file_load: "<50ms (cached)"
  timeline_build: "<50ms (from .caf/sessions/ index)"
  file_tour_memory: "<10MB"
  lsp_startup: "~500ms (gopls, small repo)"
  animated_transitions: "≤100ms"
  dashboard_auto_refresh_interval: "30s (configurable via ui.heatmap_refresh_seconds)"
  manifest_cache: "1h"
  ci_regression_alert_threshold: ">2x regression triggers alert"
---

# CodeAuthorityFirst — UI.md

> **Версия / Version:** 0.3.1-triz (post-TRIZ pass: phase conflicts resolved — Dashboard moved to MVP per ROADMAP §2.1, File tour moved to BETA per ROADMAP §4.1; UI-velocity patterns documented explicitly — Critical CSS analogue, Islands/partial hydration analogue, Skeleton/SWR analogue; @invariant tags added to every performance budget subsection; front-matter cross-refs to BACKLOG.md and ADR.md repaired.)
> **Назначение / Purpose:** detailed specification of all user-facing interfaces — primarily TUI (terminal UI), with progressive enhancements for IDE plugins and (future) web UI. Every screen has ASCII mockup, keyboard map, performance budget, and edge cases.
> **Философия / Philosophy:** TUI-first (🔥 80/20), sub-100ms response, keyboard-only operation possible, lazy rendering for large datasets. Three UI-velocity patterns (originally from web) are implemented as native TUI patterns: Critical CSS analogue (Lipgloss pre-computed styles), Islands/partial hydration analogue (lazy screen mount), Skeleton/SWR analogue (placeholder + DataLoadedMsg swap).
> **Язык / Language:** bilingual — Russian body, English summary per section
> **Аудитория / Audience:** UI/UX implementers, contributors adding TUI screens, design reviewers, LLM agents reasoning about UI velocity budgets
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md) · [ROADMAP.md](./ROADMAP.md) · [BACKLOG.md](./BACKLOG.md) · [ADR.md](./ADR.md)

---

## 0. UI Philosophy & Stack

### 0.1 Why TUI-first (🔥 80/20)

**EN summary.** TUI-first covers 80% of use cases with 20% of effort. Terminal apps start in <50ms; web apps need server + bundle. SSH/Docker/CI all have terminals, none have browsers. Keyboard-native aligns with vim/Emacs users. No Node.js, no browser, no Electron — single Go binary.

TUI-first — архитектурное решение (ADR-014), не fallback. Пять причин:

1. **🔥 80/20.** TUI покрывает 80% UX use cases 20% effort. Web UI — отложен до FUTURE.
2. **Скорость.** Terminal apps start in <50ms. Web apps 500-2000ms (server + bundle + render).
3. **Универсальность.** SSH, Docker, CI — везде есть terminal, нигде нет browser. TUI работает в tmux, в IDE-integrated terminal, в VSCode terminal.
4. **Keyboard-native.** Developers know vim. `?` for help is universal (less/vim/man).
5. **No deps.** Single Go binary. No Node.js, no browser, no Electron.

| Критерий | TUI | Web UI |
|----------|-----|--------|
| Time-to-first-paint | <100ms | 500-2000ms |
| Distribution | Single binary (~10 MB) | Server + bundle (~5 MB JS) |
| SSH/Docker/CI | ✓ universal | ✗ requires port-forward |
| Keyboard-first | ✓ native | requires tab-trapping JS |
| Maintenance cost | 1 stack (Go + Bubble Tea) | 2 stacks (Go + Node + browser quirks) |
| Coverage | 80% use cases | 95% use cases |

**Вывод:** TUI покрывает 80% use cases с 20% effort. Web UI отложен до FUTURE (см. §1.4).

### 0.2 TUI stack

**EN summary.** Stack: Bubble Tea (ELM architecture framework), Lipgloss (styling), Bubbles list (virtualization), Glamour (markdown), termenv (color detection). Rejected tview (older), gocui (lower-level), raw ANSI (error-prone). The stack implements three UI-velocity patterns natively (Critical CSS, Islands/partial hydration, Skeleton/SWR analogues — see §13.1 for the architectural mapping).

| Component | Choice | Reason |
|-----------|--------|--------|
| Framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | ELM architecture (Model-Update-View), `tea.Cmd` for async, mature |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Style strings pre-computed at startup; **Critical CSS analogue** — see §13.1 |
| List virtualization | [Bubbles list](https://github.com/charmbracelet/bubbles) | Handles >1000 items with viewport virtualization |
| Markdown | [Glamour](https://github.com/charmbracelet/glamour) | Renders decision records, README, ADRs in TUI |
| Progress bars | Bubbles spinner, progress | Standard component |
| Input fields | Bubbles textinput, textarea | Standard components |
| Tables | Bubbles table | Sortable, paginated |
| Color detection | [termenv](https://github.com/muesli/termenv) | Auto-detects 256/truecolor, falls back to 16-color |

**Rejected alternatives:** `tview` (older API, less maintained), `gocui` (lower-level, more boilerplate), raw ANSI (error-prone, no abstraction).

> **UI-velocity pattern alignment (TRIZ-pass addition, v0.3.1-triz).** The TUI stack already implements the equivalent of three web-velocity patterns. Documenting them explicitly helps agents apply web knowledge to TUI:
>
> 1. **Critical CSS analogue → Lipgloss pre-computed styles.** Lipgloss computes the final ANSI escape string for a style at definition time (e.g. `var title = lipgloss.NewStyle().Bold(true).Foreground(blue)`). On each frame, only the pre-computed string is emitted, never the raw style rules. This is structurally the same as inlining critical CSS and deferring non-critical: pre-compute once, emit minimal per frame. @invariant INV-UI-006.
> 2. **Islands / partial hydration analogue → Lazy screen mount.** Bubble Tea screens mount lazily on first navigation (`tea.Cmd` returns `LoadScreenMsg`). A screen that is never opened costs zero render time — its model is not initialized, its data is not loaded. This is structurally the same as Astro/RSC partial hydration — only the islands that need interactivity are hydrated. @invariant INV-UI-005.
> 3. **Skeleton / stale-while-revalidate analogue → Skeleton{} placeholder + DataLoadedMsg swap.** Every async data load returns a `Skeleton{}` placeholder immediately, then a `DataLoadedMsg` replaces it. The skeleton is the cached "stale" view; the new data is the "revalidate". This is structurally the same as SWR in React: show stale immediately, swap to fresh in background. @invariant INV-UI-004.
>
> These analogues are not metaphors — they are the same pattern, implemented for a TUI instead of a browser.

### 0.3 Performance budget (universal)

**EN summary.** Universal performance budget enforced for all screens. Sub-100ms first paint, sub-16ms frame render (60fps), sub-50ms keystroke, sub-3s LLM ops, sub-1s file ops on 10k-file repo. CI benchmarks; alert on >2x regression. @invariant INV-UI-001 through INV-UI-003.

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| First paint (any screen) | <100ms | Render skeleton first, lazy-load data | INV-UI-001 |
| Frame render (scroll, redraw) | <16ms (60fps) | Virtualized lists, cached layouts | INV-UI-002 |
| Keystroke response | <50ms | Direct dispatch, no async for simple actions | INV-UI-003 |
| LLM-dependent operations | <3s (acceptable) | Skeleton + spinner, cache results | (INV-UI-004 SWR pattern) |
| File operations on 10k-file repo | <1s | Streaming + index, no full load | (INV-ARCH-003 LSP fallback) |

Если операция превышает budget — показываем skeleton + spinner, не блокируем UI. Если превышение систематическое (CI benchmark regression >2x), alert срабатывает и PR не мержится без обоснования.

### 0.4 Color system

**EN summary.** Hex palette: Primary #61AFEF (blue), Secondary #98C379 (green), Warning #E5C07B (yellow), Error #E06C75 (red), Muted #5C6370 (gray), Accent #C678DD (purple), Background2 #282C34 (dark). termenv auto-detects 256/truecolor, falls back to 16-color.

```go
// pkg/ui/theme.go
var (
    Primary    = lipgloss.Color("#61AFEF") // blue — primary actions, headings
    Secondary  = lipgloss.Color("#98C379") // green — success, owned files
    Warning    = lipgloss.Color("#E5C07B") // yellow — warnings, pending decisions
    Error      = lipgloss.Color("#E06C75") // red — errors, blindspots
    Muted      = lipgloss.Color("#5C6370") // gray — muted text, meta
    Accent     = lipgloss.Color("#C678DD") // purple — accent, LLM-related
    Background = lipgloss.Color("#282C34") // dark — app background
)
```

Fallback (16-color): Primary→blue, Secondary→green, Warning→yellow, Error→red, Muted→white(dim), Accent→magenta, Background→black. termenv handles the detection at startup.

### 0.5 Typography

| Element | Style | Example |
|---------|-------|---------|
| H1 (screen title) | Bold + Accent color | `# CodeAuthorityFirst — Dashboard` |
| H2 (section) | Bold + Primary color | `## Today` |
| H3 (subsection) | Bold + Primary color | `### Review queue` |
| Body | Regular + Text color | `5 cards due, 2 pending decisions` |
| Muted meta | Muted color, smaller | `Last refresh: 30s ago` |
| Code | SarasaMonoSC | `caf who internal/middleware/ratelimit.go:42` |
| Error | Error color + Bold | `Error: caf binary missing. Sessions continue.` |

Word wrap via `github.com/muesli/reflow`.

### 0.6 Universal keyboard map

**EN summary.** Universal keyboard map enforced across all screens. `?` help, `q`/`Ctrl+C` quit, `Esc` back, `Tab`/`Shift+Tab` cycle widgets, `j`/`↓`/`k`/`↑` vim-style, `g`/`G` top/bottom, `/` search, `Enter` activate, `Space` toggle, `r` refresh, `1`-`9` screen switch, `:` command palette, `h`/`l` horizontal nav.

| Key | Action | Works in |
|-----|--------|----------|
| `?` | Help overlay | All screens |
| `q` / `Ctrl+C` | Quit | All screens |
| `Esc` | Back / Cancel | All screens |
| `Tab` / `Shift+Tab` | Cycle widgets / focus | Dashboard, multi-panel screens |
| `j` / `↓` | Down | Lists, tables, single-panel |
| `k` / `↑` | Up | Lists, tables, single-panel |
| `g` | Top | Lists, tables |
| `G` | Bottom | Lists, tables |
| `/` | Search | Lists, tables, decisions |
| `Enter` | Activate | Lists, tables, cards |
| `Space` | Toggle | Checkboxes, card selection |
| `r` | Refresh | Dashboard, decisions, adapters |
| `1`-`9` | Switch screen | Dashboard, navigation |
| `:` | Command palette | All screens (FUTURE K-010) |
| `h` / `l` | Left / Right nav | Multi-panel: review, file tour |

**Vim-style rationale:** developers know vim, no learning curve. `?` for help is universal (less/vim/man).

---

## 1. Screens Overview

### 1.1 Screen inventory

**EN summary.** 11 screens inventoried. Phase assignments reconciled with ROADMAP in v0.3.1-triz: Dashboard = MVP (was ALPHA in original UI.md §1.1, conflicting with ROADMAP §2.1); File tour = BETA (was PROD in original UI.md §1.1, conflicting with ROADMAP §4.1).

| # | Screen | Command | Phase | Spec section |
|---|--------|---------|-------|--------------|
| 1 | Onboarding Wizard | `caf init` | MVP | §2 |
| 2 | Dashboard | `caf dashboard` | **MVP** (v0.3.1-triz reconciliation with ROADMAP §2.1) | §3 |
| 3 | Review Session | `caf review` | BETA | §4 |
| 4 | Decisions list | `caf decisions` | ALPHA | §5.1 |
| 5 | Decision detail | (from list) | ALPHA | §5.2 |
| 6 | Decision editor | `caf decisions new` | ALPHA | §5.3 |
| 7 | Heatmap | `caf heatmap` | BETA | §6 |
| 8 | Adapters Manager | `caf adapters list` | MVP | §7 |
| 9 | File Tour | `caf tour <file>` | **BETA** (v0.3.1-triz reconciliation with ROADMAP §4.1) | §8 |
| 10 | Sessions browser | `caf sessions` | ALPHA | (not in this file — future §) |
| 11 | Reports | `caf report` | PROD | §10 |
| 12 | Settings | `caf config edit` | ALPHA | §11 |

### 1.2 Navigation flow

```
                        ┌──────────────────┐
                        │   Onboarding     │
                        │   (caf init)     │
                        └────────┬─────────┘
                                 │
                                 ▼
                        ┌──────────────────┐    1     ┌─────────────┐
                        │    Dashboard     │────────▶│  Review     │
                        │  (caf dashboard) │         │  Session    │
                        └────────┬─────────┘         └─────────────┘
                                 │ 2
                       ┌─────────┼─────────┐
                       │3        │4        │5
                       ▼         ▼         ▼
                  ┌────────┐ ┌────────┐ ┌────────────┐
                  │Decisions│ │Heatmap │ │Adapters    │
                  │(list)   │ │        │ │Manager     │
                  └────┬────┘ └────────┘ └────────────┘
                       │
                       ▼
                  ┌────────┐
                  │Decision│
                  │detail  │
                  └────────┘
                       │ 6 (file ref)
                       ▼
                  ┌─────────────────┐
                  │ File Tour       │
                  │ (caf tour <file>)│
                  └─────────────────┘
```

Dashboard is the hub; review/decisions/sessions branch out; file tour is leaf, collecting from all.

### 1.3 Common UI patterns

**EN summary.** Four common patterns: status bar (bottom, screen tabs + repo + agent + help), header bar (top, version + repo + counts + ACS), empty states (CTA box), loading states (skeleton + spinner + 1s warning), error states (`[r] Retry [d] Doctor [i] Ignore [q] Quit`).

**Status bar (bottom):**
```
[1] Dashboard [2] Review [3] Decisions │ repo: my-backend-service │ agent: opencode │ ? help  q quit
```

**Header bar (top):**
```
CodeAuthorityFirst v0.3.1-triz │ my-backend-service │ 47 sessions · 12 decisions · ACS 3.8 ↑
```

**Empty states:** CTA box centered on screen.
```
┌────────────────────────────────────────────────┐
│                                                │
│   No decisions yet.                            │
│                                                │
│   Run `caf decisions new` to record your first │
│   architectural decision.                      │
│                                                │
│   [n] New decision    [?] Help                 │
│                                                │
└────────────────────────────────────────────────┘
```

**Loading states:** skeleton + spinner.
```
CodeAuthorityFirst v0.3.1-triz │ my-backend-service │ Loading...
⠋ Loading dashboard... (this should take <1s. If it takes longer, run `caf doctor`)
```

**Error states:**
```
┌────────────────────────────────────────────────┐
│ Error: caf binary missing. Sessions continue.  │
│                                                │
│ Cause: $PATH does not contain `caf`            │
│ Recovery: run `make install` or `brew install  │
│ codeauthorityfirst/tap/caf`                    │
│                                                │
│ [r] Retry  [d] Doctor  [i] Ignore  [q] Quit    │
└────────────────────────────────────────────────┘
```

### 1.4 Web UI (FUTURE, not in v0.3)

**EN summary.** Web UI is deferred to FUTURE. Read-only view on `.caf/`, interactive graph visualization (D3.js — only thing impossible in TUI), team-mode real-time dashboard, export to PDF/Markdown. **Web UI НЕ заменяет TUI** (ADR-014 invariant). See BACKLOG K-002, ADR-025 (future).

Web UI НЕ заменяет TUI. Web UI добавляет только визуализации, невозможные в terminal:

- **Interactive graph visualization (D3.js).** Single thing TUI cannot do well — interactive node-link graph with zoom/pan.
- **Real-time heatmap** for team mode (BACKLOG G-003).
- **Export to PDF/Markdown** for compliance reports (BACKLOG F-003).
- **NOT replacing TUI** (ADR-014 invariant).

---

## 2. Onboarding Wizard (`caf init` interactive)

### 2.1 Purpose

**EN summary.** First `caf` run, minimum friction. ~30 seconds, 4 steps. Auto-detect installed agents. Install universal pre-commit hook. Show first-time setup summary.

### 2.2 Flow

**Step 1: Welcome**
```
┌────────────────────────────────────────────────┐
│                                                │
│   Welcome to CodeAuthorityFirst v0.3.1-triz     │
│                                                │
│   I'll set up provenance tracking for your     │
│   repo in ~30 seconds.                         │
│                                                │
│   Plan:                                        │
│   1. Detect installed agents                   │
│   2. Install universal pre-commit hook         │
│   3. Configure adapters                        │
│   4. Show summary                              │
│                                                │
│   [Enter] Continue    [q] Quit                 │
│                                                │
└────────────────────────────────────────────────┘
```

**Step 2: Detection**
```
┌────────────────────────────────────────────────┐
│                                                │
│   Detecting agents...                          │
│                                                │
│   ✓ opencode     (found in $PATH)              │
│   ✓ Claude Code  (found in ~/.claude/)         │
│   ✗ Cursor       (not found — wrapper avail.)  │
│   ✗ Aider        (not found — wrapper avail.)  │
│                                                │
│   Universal fallback (LSP-based) will track    │
│   all other agents automatically.              │
│                                                │
│   [Enter] Continue                             │
│                                                │
└────────────────────────────────────────────────┘
```

**Step 3: Installation progress**
```
┌────────────────────────────────────────────────┐
│                                                │
│   Installing adapters...                       │
│                                                │
│   opencode      ████████████████████ 100%      │
│   Claude Code   ████████████████████ 100%      │
│   Universal hook ████████████████████ 100%      │
│   LSP fallback  ████████████████████ 100%      │
│                                                │
│   [Enter] Continue                             │
│                                                │
└────────────────────────────────────────────────┘
```

**Step 4: Summary**
```
┌────────────────────────────────────────────────┐
│                                                │
│   Setup complete.                              │
│                                                │
│   Installed:                                   │
│   • .caf/ folder created                       │
│   • .git/hooks/pre-commit installed (universal)│
│   • Adapters: opencode, Claude Code            │
│   • Universal LSP fallback active               │
│                                                │
│   Next steps:                                  │
│   • Run `caf dashboard` to see overview        │
│   • Run `caf who <file:line>` for provenance  │
│   • Run `caf context` to generate context.md   │
│                                                │
│   [Enter] Finish    [d] Open dashboard        │
│                                                │
└────────────────────────────────────────────────┘
```

### 2.3 Edge cases

- Existing `.caf/` — `caf init` detects, offers `caf migrate` if version mismatch.
- No agent detected — universal fallback will still track; show message "universal LSP fallback active; no explicit adapter needed."
- Permission denied on `.git/hooks/` — fallback to user-level hooks in `~/.config/git/hooks/`.
- Non-git repo — `caf init` warns, offers to initialize git first.

---

## 3. Dashboard (`caf dashboard`)

### 3.1 Purpose

**EN summary.** Dashboard is the hub screen. Shows Today widget (review queue, pending decisions, new blindspots), Authorship Score widget (ACS weekly avg, target ≥3.5, Recall %, Coverage %), Active Agents widget, Recent Decisions widget, Cross-repo widget (dependencies with stale-decision indicators).

### 3.2 Layout

```
CodeAuthorityFirst v0.3.1-triz │ my-backend-service │ 47 sessions · 12 decisions · ACS 3.8 ↑
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
┌─ Today ────────────────────────────┐  ┌─ Authorship Score ──────────────────────────────────────────┐
│ Review queue:        5 cards due    │  │                                                              │
│ Pending decisions:  2 unrecorded    │  │   ACS (weekly avg): 3.8 ↑        Target: ≥3.5  ✓           │
│ New blindspots:      3 files flagged│  │                                                              │
│                                     │  │   Recall:    78%   Target: ≥75%  ✓                         │
│ [Enter] Start review session        │  │   Coverage:  92%   Target: ≥90%  ✓                         │
│ [n]     New decision                │  │                                                              │
└─────────────────────────────────────┘  │   Last 7 days: ▁▂▄▅▆▇█ 3.8                                  │
                                         │                                                              │
┌─ Active Agents ─────────────────────┐  │   [m] View metrics report                                   │
│ opencode        ●  (12 sessions)    │  └──────────────────────────────────────────────────────────────┘
│ Claude Code     ●  (8 sessions)     │
│ Universal LSP   ●  (3 sessions)     │  ┌─ Recent Decisions ───────────────────────────────────────────┐
│                                     │  │ 2026-10-24  Sliding window for rate limiter   [opencode]   │
└─────────────────────────────────────┘  │ 2026-10-22  Token bucket as fallback          [Claude Code]│
                                         │ 2026-10-20  Postgres connection pool size=20  [opencode]   │
┌─ Cross-repo dependencies ──────────┐  │                                                              │
│ my-backend-service ──▶ my-proto     │  │ [Enter] Open detail   [n] New decision                     │
│  └ stale decision (3d ago)         │  └──────────────────────────────────────────────────────────────┘
│ my-backend-service ──▶ my-utils    │
│  └ up to date                       │  ┌─ Cross-repo ────────────────────────────────────────────────┐
│                                     │  │ my-backend-service depends on:                              │
│ [x] Run `caf sync`                 │  │   • my-proto     (3d stale decision — ⚠)                   │
└─────────────────────────────────────┘  │   • my-utils     (up to date — ✓)                          │
                                         │                                                              │
                                         │   [x] Run `caf sync`                                        │
                                         └──────────────────────────────────────────────────────────────┘
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[1] Dashboard [2] Review [3] Decisions │ repo: my-backend-service │ agent: opencode │ ? help  q quit
```

### 3.3 Components

Five widgets:
1. **Today** — review queue, pending decisions, new blindspots.
2. **Authorship Score** — ACS weekly avg with sparkline, Recall %, Coverage %, targets comparison.
3. **Active Agents** — agents that have run sessions in last 7 days.
4. **Recent Decisions** — last 3-5 decisions, sorted by date desc.
5. **Cross-repo** — dependencies with stale-decision indicators (⚠) or up-to-date (✓).

### 3.4 Keyboard map

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Cycle widgets |
| `1`-`5` | Jump to widget N |
| `Enter` | Activate widget action (e.g. start review from Today widget) |
| `r` | Manual refresh |
| `R` | Force full refresh (skip cache) |
| `x` | Run cross-repo sync (from Cross-repo widget) |
| `n` | New decision |
| `m` | View metrics report |
| `?` | Help |

### 3.5 Refresh strategy

**EN summary.** Auto-refresh every 30 seconds (configurable via `.caf/config.yml` `ui.heatmap_refresh_seconds`). Manual `r` for incremental update (only changed widgets). `R` for full refresh (skip cache). File watcher on `.caf/` triggers incremental update.

- Auto-refresh: every 30 seconds (configurable via `.caf/config.yml` `ui.heatmap_refresh_seconds`). Note: the config key name includes "heatmap" for historical reasons; in v0.3.1-triz it controls dashboard refresh as well. Future ADR may rename to `ui.dashboard_refresh_seconds`.
- Manual `r`: incremental update — only changed widgets. Uses file watcher on `.caf/` to know which widgets need refresh.
- Manual `R`: full refresh, skip cache. Use when something looks wrong.
- File watcher on `.caf/sessions/` and `.caf/decisions/` triggers incremental update immediately on change.

### 3.6 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| First paint | <100ms | Render skeleton widgets immediately (INV-UI-004 SWR pattern), load real data in background | INV-UI-001 |
| Full data load | <500ms | Parallel widget loads; caching | (INV-UI-005 lazy mount) |
| Auto-refresh | <50ms | Only changed widgets; incremental diff | INV-UI-002 |
| Memory | <50MB | Heatmap pre-computed (not full file content); only visible widgets in memory | (INV-UI-005) |

### 3.7 Edge cases

- No sessions yet — empty state with CTA "Run `caf init` to start tracking."
- No decisions yet — empty state with CTA "Run `caf decisions new` to record your first decision."
- Cross-repo not synced — Cross-repo widget shows "Run `caf sync`" CTA.
- Terminal <120 cols — widgets stack vertically (UI.md §12.2 responsive).

---

## 4. Review Session (`caf review`)

### 4.1 Purpose

**EN summary.** Daily SRS session. ~10 minutes, 10-15 cards. Main contact point with Mental Model Reconstructor. Split-view layout: question (left) + context (file path, decision ref, sessions) + grade panel.

### 4.2 Layout — Card view

```
CodeAuthorityFirst v0.3.1-triz │ Review Session │ Card 3/12 │ Recall rate (7d): 78% │ ACS target: ≥3.5
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
┌─ Question ─────────────────────────────────────────────┐  ┌─ Context ───────────────────────────────────────┐
│                                                       │  │                                                  │
│  Type 1 (why)                                         │  │  File: internal/middleware/ratelimit.go         │
│                                                       │  │  Decision: d_2026-10-24-01                       │
│  Why does this code use a sliding window instead of   │  │  Sessions: 2026-10-24 (opencode), 2026-10-22     │
│  a fixed window for rate limiting?                   │  │            (Claude Code)                         │
│                                                       │  │                                                  │
│                                                       │  │  Related decisions:                              │
│                                                       │  │  • d_2026-10-22-01  Token bucket as fallback    │
│                                                       │  │                                                  │
│                                                       │  │  Related cards:                                  │
│                                                       │  │  • c_001  Type 1  (this card)                    │
│                                                       │  │  • c_005  Type 2  (what if remove)               │
│                                                       │  │                                                  │
└───────────────────────────────────────────────────────┘  └──────────────────────────────────────────────────┘
┌─ Grade ───────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                            │
│  [0] Blackout (no memory)                                                                                 │
│  [1] Wrong, but familiar                                                                                  │
│  [2] Wrong, but easy to find                                                                              │
│  [3] Hard recall, took effort                                                                             │
│  [4] Good recall, minor hesitation                                                                        │
│  [5] Perfect, instant recall                                                                              │
│                                                                                                            │
│  [n] Skip (= grade 0)    [s] Show expected answer                                                          │
│                                                                                                            │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[?] Help  [q] Quit  [s] Show answer  [0-5] Grade  [n] Skip  [←/→] Card history
```

### 4.3 Card types rendered differently

- **Type 1 (why):** blue accent. "Why does this code use X?" — tests fact recall.
- **Type 2 (what_if_remove):** yellow accent. "If we removed X, what would break?" — tests abstract reasoning.
- **Type 3 (where_handled):** green accent. "Where in the codebase is X handled from entry to retry?" — tests trace through code.

### 4.4 Expected answer reveal

```
┌─ Expected answer ──────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                            │
│  Decision d_2026-10-24-01 (recorded by opencode session 2026-10-24):                                       │
│                                                                                                            │
│  "Sliding window имеет смысл при 100+ req/sec, у нас на 3 порядка меньше. Fixed window достаточно и       │
│  проще в reasoning."                                                                                        │
│                                                                                                            │
│  ── from .caf/decisions/d_2026-10-24-01.md (grounded in .caf/, ADR-006 invariant: LLM only paraphrases)   │
│                                                                                                            │
│  [0-5] Grade  [n] Skip  [→] Next card                                                                      │
│                                                                                                            │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.5 Session summary (after last card)

```
CodeAuthorityFirst v0.3.1-triz │ Review Session │ Summary │ 12/12 cards
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
┌─ Grade distribution ──────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                            │
│  [5] ████████████ 4 cards (33%)                                                                            │
│  [4] ████████████ 4 cards (33%)                                                                            │
│  [3] ██████ 2 cards (17%)                                                                                  │
│  [2] ███ 1 card (8%)                                                                                       │
│  [1] █ 1 card (8%)                                                                                         │
│  [0] 0 cards                                                                                               │
│                                                                                                            │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

┌─ Recall rate ─────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                            │
│  This session: 75%  (9/12 cards ≥ grade 3)                                                                 │
│  7-day avg:    78%  Target: ≥75%  ✓                                                                        │
│                                                                                                            │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

┌─ Final reflection ────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                            │
│  "Мог бы я написать этот код с нуля?" (1-5)                                                                 │
│                                                                                                            │
│  [1] No  [2] With help  [3] Maybe  [4] Probably  [5] Yes                                                   │
│                                                                                                            │
└────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.6 Keyboard map

| Key | Action |
|-----|--------|
| `0`-`5` | Grade current card |
| `n` | Skip (= grade 0) |
| `s` | Show expected answer |
| `→` | Next card |
| `←` | Previous card (only in review mode, not after grade) |
| `?` | Help |
| `q` | Quit (saves partial progress) |
| `Ctrl+C` | Force quit (saves partial progress) |

### 4.7 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| Card render | <50ms | Pre-computed display, no LLM call | INV-UI-002 |
| LLM call (if `--explain` flag) | <3s | Skeleton + spinner (INV-UI-004 SWR pattern), cache results | (INV-UI-004) |
| Session save | <100ms | Incremental write to `cards.yml` after each card grade | (INV-ADR-005 atomic writes) |

### 4.8 Edge cases

- No cards due — "No cards due today. Come back tomorrow." with CTA to view heatmap.
- Missing decision for card — "Decision d_xxx missing. Run `caf sync`?" `[y/n]`
- Ctrl+C — saves partial progress (`cards.yml` is incrementally written).
- 100+ cards backlog — warns, offers `--max-cards 15` cap (mitigates SM-2 burnout risk, BACKLOG D-007).

---

## 5. Decisions (`caf decisions`)

### 5.1 Decisions list

**EN summary.** Filterable table of decisions. Filters: today/week/month/all, by agent, by file, by status. Actions: `Enter` open, `n` new, `/` search, `g` gaps.

```
CodeAuthorityFirst v0.3.1-triz │ Decisions │ 12 total │ 2 unrecorded non-trivial edits (gaps)
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Filter: [t] today  [w] week  [m] month  [a] all   │   [g] gaps   │   [/] search
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Date        │ Title                                    │ Agent       │ Files
────────────┼──────────────────────────────────────────┼─────────────┼──────────────────────
2026-10-24  │ Sliding window for rate limiter           │ opencode    │ internal/middleware/rl.go
2026-10-22  │ Token bucket as fallback                  │ Claude Code │ internal/middleware/tb.go
2026-10-20  │ Postgres connection pool size=20          │ opencode    │ internal/db/pool.go
2026-10-18  │ Context cancellation in HTTP handler      │ opencode    │ internal/handler/ctx.go
2026-10-15  │ Circuit breaker for payments service      │ Claude Code │ internal/payments/cb.go
...
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[Enter] Open detail  [n] New decision  [/] Search  [g] Gaps  [t/w/m/a] Filter  [?] Help  [q] Quit
```

### 5.2 Decision detail view

```
CodeAuthorityFirst v0.3.1-triz │ Decision d_2026-10-24-01
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
# Sliding window for rate limiter

**Date:** 2026-10-24
**Agent:** opencode
**Files:** internal/middleware/rl.go (+45 -12)

## Контекст

We need a rate limiter for the /v1/policies endpoint. Options considered: fixed window, sliding window, token bucket.

## Выбрано

Sliding window with 1-minute window size. Reasoning:

- Sliding window имеет смысл при 100+ req/sec, у нас на 3 порядка меньше.
- Fixed window проще в reasoning, но creates burst-at-boundary problem.
- Token bucket is over-engineered for our traffic.

## Альтернативы

- Fixed window — simpler but creates burst-at-boundary problem.
- Token bucket — over-engineered for our traffic; harder to debug.

## Импликации

- p99 latency target: <100ms.
- Memory: O(window_size × unique_clients) — bounded.
- TTL: 5 minutes for inactive clients.

---
Related sessions: 2026-10-24 (opencode) — 12 file edits, 3 iterations
Related cards: c_001 (Type 1), c_005 (Type 2)

[e] Edit  [d] Delete  [s] Show file  [Esc] Back
```

### 5.3 Decision editor (`caf decisions new`)

**🔥 80/20 TUI:** автозаполнение для `files`, `tags`, `agent` сокращает время с 90s до 25s.

```
CodeAuthorityFirst v0.3.1-triz │ New Decision
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Title:       [Sliding window for rate limiter                                    ]
Agent:       [opencode                                                           ]  (autocomplete: ↓ opencode, Claude Code, Cursor...)
Files:       [internal/middleware/rl.go                                          ]  (autocomplete: ↓ from current/recent sessions)
Tags:        [rate-limiting, middleware                                          ]  (autocomplete: ↓ from previous decisions)
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
## Контекст

We need a rate limiter for the /v1/policies endpoint. Options considered: fixed window,
sliding window, token bucket.

## Выбрано

[Sliding window with 1-minute window size. Reasoning: ...]

## Альтернативы

[- Fixed window — simpler but creates burst-at-boundary problem.]

## Импликации

[- p99 latency target: <100ms.]

─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[Tab] Next field  [Enter] Save  [Esc] Cancel  [?] Help
```

### 5.4 Keyboard map

| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Save |
| `Esc` | Cancel |
| `↓` (in autocomplete fields) | Show autocomplete suggestions |
| `?` | Help |

### 5.5 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| List render | <50ms | Virtualized (only visible rows rendered) | INV-UI-002 |
| Filter response | <16ms | In-memory index, no disk I/O | INV-UI-002 |
| Detail view | <100ms | Markdown render via Glamour, cached | (INV-UI-005 lazy mount) |
| Autocomplete suggestions | <50ms | Pre-built index, refreshed on session write | INV-UI-002 |

---

## 6. Heatmap (`caf heatmap`)

### 6.1 Purpose

**EN summary.** Heatmap shows each file as a cell, color = ACS or coverage. Single-screen overview of authorship health. Pre-computed and cached.

### 6.2 Layout

```
CodeAuthorityFirst v0.3.1-triz │ Heatmap │ 247 files │ Sort: [n] name  [o] ownership  [r] recency
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
internal/
  middleware/    ██ rl.go        ██ tb.go        ░░ ctx.go        ▓▓ auth.go       ▓▓ logging.go
  handler/       ██ policies.go  ██ users.go     ░░ admin.go      ▒▒ health.go      ██ metrics.go
  db/            ██ pool.go      ▓▓ conn.go      ░░ tx.go          ▒▒ migrate.go
  caf/           ██ adapter.go   ██ storage.go   ▒▒ hooks.go       ▒▒ sync.go         ░░ doctor.go

pkg/ui/
  screens/       ██ dashboard.go ██ review.go    ██ decisions.go   ▓▓ heatmap.go      ▒▒ adapters.go
  components/    ██ header.go     ██ statusbar.go ▓▓ spinner.go     ▓▓ skeleton.go

Legend: ░░ 0-20% (blindspot)  ▒▒ 20-50% (weak)  ▓▓ 50-75% (decent)  ██ 75-100% (owned)
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[b] Blindspots only  [/] Filter by path  [n/o/r] Sort  [Enter] File tour  [?] Help  [q] Quit
```

### 6.3 Color encoding

| Range | Color | Meaning |
|-------|-------|---------|
| 0-20% | dark gray (░░) | blindspot — needs SRS card |
| 20-50% | red (▒▒) | review overdue |
| 50-75% | yellow (▓▓) | interval low |
| 75-100% | green (██) | strong recall |

### 6.4 Keyboard map

| Key | Action |
|-----|--------|
| `j`/`↓` `k`/`↑` | Navigate files |
| `Enter` | Open file tour for selected file |
| `b` | Toggle blindspots-only filter |
| `/` | Filter by path regex |
| `n` | Sort by name |
| `o` | Sort by ownership (ASC) |
| `r` | Sort by recency (most recent first) |
| `?` | Help |
| `q` | Quit |

### 6.5 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| Render 247 files | <200ms | Pre-computed heatmap, cached | INV-UI-007 |
| Re-compute on file change | <500ms | Incremental — only changed files | INV-UI-007 |
| Memory | <20MB | Heatmap pre-computed (not full file content); only cell metadata in memory | INV-UI-005 |

### 6.6 Edge cases

- Empty repo — "No files tracked yet. Run `caf sync`."
- Repo with 10k+ files — render first 1000 by default, paginated. (Note: original UI.md §6.6 said "render first 1000" but did not specify a budget for the 10k-file case; the 1000-file budget <200ms applies. For >1000 files, render time scales linearly.)
- New file <7 days old — shows "new" indicator (cyan dot).
- Deleted file — shows `✗` fading over 30 days (gray with strikethrough).

---

## 7. Adapters Manager (`caf adapters list`)

### 7.1 Purpose

**EN summary.** Lists installed + available adapters, with capability matrix and health check. Surfaces which canonical events each adapter supports natively vs via fallback.

### 7.2 Layout

```
CodeAuthorityFirst v0.3.1-triz │ Adapters Manager │ 4 installed / 21 available
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Installed (4):
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Adapter           │ Tier │ pre-session │ post-edit │ pre-commit │ MCP │ Status
──────────────────┼──────┼─────────────┼───────────┼─────────────┼─────┼───────────
opencode          │ 1    │ ✓ native    │ ✓ native  │ ✓ native    │ ✓   │ ✓ healthy
Claude Code       │ 1    │ ✓ native    │ ✓ native  │ ✓ native    │ ✓   │ ✓ healthy
Universal LSP     │ 4    │ ✓ wrapper   │ ✓ watcher │ ✓ universal │ ✓   │ ✓ healthy
Universal hook    │ —    │ —           │ —         │ ✓ universal │ —   │ ✓ healthy
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Available (21):
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Cursor            │ 1    │ ✓ wrapper   │ ✓ native  │ ✓ native    │ ✓   │ [i] Install
Cline             │ 1    │ ✓ native    │ ✓ native  │ ✓ native    │ ✓   │ [i] Install
KiloCode          │ 1    │ ✓ native    │ ✓ native  │ ✓ native    │ ✓   │ [i] Install
Aider             │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Codex CLI         │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Continue.dev      │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
GitHub Copilot    │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Gemini CLI        │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Amazon Q          │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Goose             │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Windsurf          │ 2    │ ✗           │ ✓ wrapper │ ✓ wrapper   │ —   │ [i] Install
Tabnine           │ 3    │ ✗           │ ✗         │ ✓ wrapper   │ —   │ [i] Install
Sourcegraph Cody  │ 3    │ ✗           │ ✗         │ ✓ wrapper   │ —   │ [i] Install
...
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[↑↓] Navigate  [Enter] Install / Open detail  [h] Health check  [m] Migrate config  [?] Help  [q] Quit
```

### 7.3 Keyboard map

| Key | Action |
|-----|--------|
| `j`/`↓` `k`/`↑` | Navigate list |
| `Enter` | Install selected adapter (from Available) or open detail (from Installed) |
| `h` | Run health check on all installed adapters |
| `m` | Migrate config (e.g. `caf adapters install aider --migrate-config`) |
| `u` | Uninstall selected adapter |
| `?` | Help |
| `q` | Quit |

### 7.4 Health check view

```
Health check — 4 adapters
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
opencode
  ✓ pre-session hook fires (config: ~/.config/opencode/hooks.toml)
  ✓ post-edit writes to .caf/sessions/ (verified by test session)
  ✓ pre-commit validation passes
  ✓ config file present
  ✓ MCP server responds

Claude Code
  ✓ pre-session hook fires (config: .claude/settings.json)
  ✓ post-edit writes to .caf/sessions/
  ✓ pre-commit validation passes
  ✓ config file present
  ✓ MCP server responds

Universal LSP
  ✓ LSP server started (gopls, ~500ms startup)
  ✓ symbol graph built (247 files)
  ✓ saved to .caf/mental-model/graph.yml
  ✓ pre-commit hook installed in .git/hooks/pre-commit

Universal hook
  ✓ .git/hooks/pre-commit exists
  ✓ hook is executable
  ✓ caf binary in $PATH

─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Summary: 4 passed, 0 failed, 0 warnings                          [r] Re-run  [Esc] Back
```

### 7.5 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| List render | <50ms | Virtualized, only visible rows | INV-UI-002 |
| Health check | <5s per adapter (parallel) | Parallel goroutines; <30s total for 25 adapters | (INV-ARCH-003 lazy LSP) |
| Install | <2s per adapter (parallel up to 5) | Goroutines; error isolation | (INV-ARCH-002 stable interface) |

---

## 8. File Tour (`caf tour <file>`)

### 8.1 Purpose

**EN summary.** File tour shows the history of a single file: timeline (left, sessions/decisions/cards) + numbered source (right). Useful for "why does this file look the way it does?" archaeology.

### 8.2 Layout

```
CodeAuthorityFirst v0.3.1-triz │ File Tour │ internal/middleware/ratelimit.go (142 lines, 3 sessions)
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
┌─ Timeline ─────────────────────────────────────┐  ┌─ Source: internal/middleware/ratelimit.go ──────────────────────────────────────┐
│                                                │  │                                                                                  │
│  ● 2026-10-24 (post-edit, opencode)            │  │   1  package middleware                                                          │
│  "feat: sliding window for /v1/policies"      │  │   2                                                                                │
│  +45 -12 lines                                 │  │   3  import (                                                                      │
│                                                │  │   4      "context"                                                                 │
│  ● 2026-10-22 (post-edit, Claude Code)         │  │   5      "time"                                                                    │
│  "fix: race condition in token counter"        │  │   6      "sync"                                                                    │
│  +12 -3 lines                                  │  │   7  )                                                                             │
│                                                │  │   8                                                                                │
│  ● 2026-10-20 (pre-session, opencode)          │  │   9  // TokenBucket implements a sliding-window rate limiter.                     │
│  "refactor: extract rate limiter to its own    │  │  10  // See decision d_2026-10-24-01 for reasoning.                                │
│  file"                                         │  │  11  type TokenBucket struct {                                                    │
│  +142 -0 lines (file created)                  │  │  12      tokens chan struct{}                                                      │
│                                                │  │  13      rate   time.Duration                                                       │
│  Related decisions:                            │  │  14      mu     sync.Mutex                                                         │
│  • d_2026-10-24-01  Sliding window             │  │  15  }                                                                             │
│  • d_2026-10-22-01  Race condition fix         │  │  16                                                                                │
│                                                │  │  17  func NewTokenBucket(rate time.Duration, capacity int) *TokenBucket {          │
│  Related cards:                                │  │  18      tb := &TokenBucket{                                                       │
│  • c_001  Type 1  (why)                       │  │  19          tokens: make(chan struct{}, capacity),                                 │
│  • c_005  Type 2  (what if remove)            │  │  20          rate:   rate,                                                          │
│  • c_012  Type 3  (where handled)             │  │  21      }                                                                          │
│                                                │  │  22      go tb.refill()                                                             │
│                                                │  │  23      return tb                                                                  │
│                                                │  │  24  }                                                                              │
│                                                │  │  ... (118 more lines, paginated)                                                   │
│                                                │  │                                                                                   │
│  [j/k] Navigate timeline                       │  │  [j/k] Scroll source  [n/p] Next/prev page  [Enter] Jump to line                  │
│  [Enter] Jump to file:line                     │  │                                                                                   │
│  [d] Show decision detail                      │  │                                                                                   │
│  [c] Show card                                 │  │                                                                                   │
│  [Esc] Back                                    │  │                                                                                   │
└────────────────────────────────────────────────┘  └───────────────────────────────────────────────────────────────────────────────────┘
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[?] Help  [q] Quit
```

### 8.3 Keyboard map

| Key | Action |
|-----|--------|
| `j`/`↓` `k`/`↑` | Navigate timeline (left panel) or scroll source (right panel) |
| `Tab` | Switch between timeline and source panel |
| `Enter` | Jump to file:line from selected timeline event |
| `d` | Show decision detail for selected event |
| `c` | Show card for selected event |
| `n`/`p` | Next/prev page of source (for files >1000 lines) |
| `?` | Help |
| `q` | Quit |

### 8.4 Performance budget

| Operation | Budget | Strategy | Invariant |
|-----------|--------|----------|-----------|
| First paint | <100ms | Render skeleton (INV-UI-004 SWR); load timeline + source in background | INV-UI-001 |
| File load | <50ms | Cached after first open | (INV-UI-005 lazy mount) |
| Timeline build | <50ms | From `.caf/sessions/` index (pre-built) | (INV-ARCH-003) |
| Memory | <10MB | Only current file in memory; timeline from index | INV-UI-005 |

### 8.5 Edge cases

- File never touched by AI — "No AI sessions recorded for this file. Run `caf sync` to re-scan."
- File deleted — show last known content with `✗` indicator and "file deleted on YYYY-MM-DD" header.
- Binary file — show "Binary file, no source view." with timeline still available.
- File >5000 lines — paginate source 1000 lines per page. `[n]` next page. Timeline unaffected.
- File <50 lines — full file shown; no pagination needed.

---

## 9. Onboarding Wizard (extended)

### 9.1 Advanced mode (`caf init --advanced`)

```
┌────────────────────────────────────────────────┐
│ Advanced settings                              │
│                                                │
│ Triviality threshold (LOC): [50]              │
│ Pre-commit validation:        [✓] enabled      │
│ Auto-detect agents:           [✓] enabled      │
│ Dashboard auto-refresh (sec): [30]             │
│ Theme:                        ( ) dark         │
│                               (•) light        │
│                               ( ) auto         │
│                                                │
│ [Tab] Next field  [Enter] Continue             │
└────────────────────────────────────────────────┘
```

### 9.2 Migration from v0.1/v0.2

`caf init` detects existing `.orya/` (v0.1) or `.caf/` (v0.2) and offers migration:

```
┌────────────────────────────────────────────────┐
│ Existing v0.2 .caf/ folder detected.           │
│                                                │
│ Migration will:                                │
│   • Add agent_tier, auto_detected fields      │
│   • Add ui config section (theme, refresh)     │
│   • Preserve all existing data                 │
│                                                │
│ Estimated time: ~5 seconds (247 sessions)      │
│                                                │
│ [Enter] Migrate  [s] Skip  [?] Help            │
└────────────────────────────────────────────────┘
```

---

## 10. Reports View (`caf report`)

### 10.1 Daily report

```
CodeAuthorityFirst v0.3.1-triz │ Daily Report │ 2026-10-28
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
Sessions today: 4
  • opencode      2 sessions (45 + 28 minutes)
  • Claude Code   2 sessions (32 + 19 minutes)

Commits with `CAF-Session` trailer: 3/3 ✓

Decisions new: 2
  • d_2026-10-28-01  "Why we switched from sync.Map to RWMutex"
  • d_2026-10-28-02  "Context cancellation in retry logic"

Cards due today: 5  Cards reviewed: 5  Recall rate: 80%

What you did well:
  • 100% decision coverage on non-trivial commits
  • Recall rate above 7-day avg (78%)

What to improve:
  • 1 blindspot: internal/middleware/circuit_breaker.go (last reviewed 14 days ago)
  • Agent diversity below target (H=0.8, target ≥1.0)

Tomorrow's queue:
  • 3 cards due (Type 1: 2, Type 2: 1)
  • 1 pending decision from today's session

[Esc] Back  [w] Week report  [m] Metrics report  [?] Help
```

### 10.2 Week report

```
CodeAuthorityFirst v0.3.1-triz │ Week Report │ 2026-10-22 to 2026-10-28
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
ACS (7-day avg):        3.5 → 3.8  ↑+0.3   Target: ≥3.5  ✓
Recall (7-day avg):    72% → 80%   ↑+8pp   Target: ≥75%  ✓
Decision coverage:    88% → 92%   ↑+4pp   Target: ≥90%  ✓
Sessions/day:           4.2 → 4.5   ↑+0.3
Agent diversity (H):   0.9 → 1.2   ↑+0.3   Target: ≥1.0  ✓

7-day ACS sparkline: ▁▂▃▄▅▆▇█

Top contributors (this week):
  • opencode      12 sessions
  • Claude Code    8 sessions
  • Universal LSP 3 sessions
```

### 10.3 Metrics report

```
CodeAuthorityFirst v0.3.1-triz │ Metrics Report (30-day rolling)
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
ACS (Authorship Confidence Score):     3.8   Target: ≥3.5   ✓
Recall rate:                            78%   Target: ≥75%   ✓
Decision coverage:                      92%   Target: ≥90%  ✓
Mental model reach:                     68%   Target: ≥70%  ⚠ below target
  Action: `caf review --regenerate --since 30d`
Agent diversity (Shannon entropy H):     1.2   Target: ≥1.0  ✓

Per-agent distribution:
  opencode      ████████████████ 48%
  Claude Code   ███████████       32%
  Universal LSP ████              12%
  Cursor        ███                8%

Adapter coverage (native):              70%   Target: ≥70%  ✓ (4 native / 25 explicit + universal)
```

---

## 11. Settings Editor (`caf config edit`)

### 11.1 Layout

```
CodeAuthorityFirst v0.3.1-triz │ Settings
─────────────────────────────────────────────────────────────────────────────────────────────────────────────
General
  Triviality threshold (LOC):    [50]
  Pre-commit validation:        [✓] enabled
  Fail-safe mode:               [✓] enabled (recommended — see ARCHITECT.md §2.5 fail-safe contract)

Adapters
  [✓] opencode      (Tier 1, native)
  [✓] Claude Code   (Tier 1, native)
  [ ] Cursor        (Tier 1, wrapper)
  [ ] Aider         (Tier 2, wrapper)
  [✓] Universal LSP (Tier 4, fallback)

UI
  Theme:                         (•) dark   ( ) light   ( ) auto   ( ) high-contrast
  Review session max cards:     [15]
  Dashboard refresh (sec):      [30]

LLM (PROD)
  Backend:                       [z-ai-web-dev-sdk ▼]
  First-person voice:           [✓] enabled (recommended — ADR-006)

─────────────────────────────────────────────────────────────────────────────────────────────────────────────
[Ctrl+S] Save  [Esc] Cancel  [?] Help
```

---

## 12. Cross-cutting Concerns

### 12.1 Terminal compatibility

| Terminal | Tested | Notes |
|----------|--------|-------|
| iTerm2 (macOS) | ✓ | Full truecolor, mouse, all features |
| Terminal.app (macOS) | ✓ | 256-color only (no truecolor) |
| Alacritty | ✓ | Full truecolor |
| Kitty | ✓ | Full truecolor, graphics protocol |
| Windows Terminal | ✓ | Full truecolor |
| VSCode integrated | ✓ | 256-color, mouse works |
| Git Bash (Windows) | partial | Limited color, some rendering issues. Mitigation: fall back to 16-color via termenv; surface as `caf doctor` warning recommending WSL or Windows Terminal. |
| Linux TTY (Ctrl+Alt+F1) | partial | 16-color only; no mouse |
| SSH | ✓ | Full features (depends on client terminal) |

### 12.2 Responsive design

| Width (cols) | Layout |
|--------------|--------|
| ≥140 | Full grid (3-column dashboard) |
| 100-139 | 2-column dashboard |
| 60-99 | Single column, widgets stacked |
| <60 | Compact mode: "Widen terminal for full dashboard." + key info only |

### 12.3 Accessibility

- Color-blind friendly: pair color with text/symbol (✓ ⚠ ✗), never color-only.
- Screen reader hints via ANSI escape codes (`\033[?25h` to show cursor, etc.).
- Keyboard-only operation possible (no mouse required, FUTURE K-011 adds mouse support).
- No time-based-only cues (no flashing).
- High-contrast mode: `caf config set ui.theme high-contrast`.

### 12.4 Internationalization

- Body text: RU/EN depending on user locale (`ui.language` config).
- CLI names always English (universal).
- YAML field names always English (machine-readable).
- Decision content: user's language (recorded as-is, not translated).

### 12.5 Telemetry (opt-in, FUTURE, BACKLOG J-004)

Anonymous command counts + error rates. Opt-in via `caf config set telemetry.enabled true`. No content (no prompts, paths, decisions). Public dashboard: `caf telemetry stats`. See BACKLOG J-004.

---

## 13. Implementation Notes

### 13.1 File organization

```
pkg/ui/
├── tui.go              # root model + Update + View
├── theme.go            # color palette, typography
├── keys.go             # universal keyboard map (§0.6)
├── perf.go             # performance budget enforcer + CI benchmark hooks (§0.3)
├── components/
│   ├── header.go       # header bar
│   ├── statusbar.go    # status bar
│   ├── spinner.go      # loading spinner
│   ├── skeleton.go     # skeleton placeholder (INV-UI-004 SWR pattern)
│   └── heatmap.go      # heatmap cell renderer (§6)
└── screens/
    ├── onboarding.go     # §2 / §9
    ├── dashboard.go      # §3
    ├── review.go         # §4
    ├── decisions_list.go      # §5.1
    ├── decisions_detail.go    # §5.2
    ├── decisions_editor.go    # §5.3
    ├── heatmap.go             # §6
    ├── adapters.go            # §7
    ├── file_tour.go           # §8
    ├── reports.go             # §10
    └── settings.go            # §11
```

Each screen file maps to its spec section by file comment: `// §3 Dashboard (caf dashboard)`.

### 13.2 State management

**EN summary.** Bubble Tea ELM architecture. `model` struct holds activeScreen / repo / agent / data / theme / width / height. `Update(msg)` handles `tea.WindowSizeMsg`, `tea.KeyMsg`, `DataLoadedMsg`. Async loading via `tea.Cmd` returning `DataLoadedMsg{Sessions}`. The pattern is structurally equivalent to React's SWR (stale-while-revalidate) — skeleton shown immediately, fresh data swaps in.

```go
type model struct {
    activeScreen Screen
    repo         string
    agent        string
    data         map[string]interface{} // generic; typed per screen
    theme        Theme
    width        int
    height       int
    loading      map[string]bool         // INV-UI-004: which keys have Skeleton shown
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        return m, nil
    case tea.KeyMsg:
        // handle keystroke (<50ms, INV-UI-003 direct dispatch)
        return m, m.handleKey(msg)
    case DataLoadedMsg:
        // SWR pattern: replace Skeleton with real data (INV-UI-004)
        m.data[msg.Key] = msg.Data
        delete(m.loading, msg.Key)
        return m, nil
    }
    return m, nil
}
```

Async loading via `tea.Cmd` returning `DataLoadedMsg`:

```go
func loadSessionsCmd(repo string) tea.Cmd {
    return func() tea.Msg {
        sessions, _ := storage.Sessions(repo)
        return DataLoadedMsg{Key: "sessions", Data: sessions}
    }
}
```

**Skeleton/SWR pattern:** when a screen mounts, it immediately shows a `Skeleton{}` for each data key in `m.loading`. When `DataLoadedMsg` arrives, the skeleton is replaced with real data. This is the same shape as React's SWR hook: show stale immediately, swap to fresh in background.

### 13.3 Virtualization strategy

- `bubbles/list` for any list with >100 items. `SetFilteringEnabled(true)` for search.
- For file tour source view >5000 lines: paginate 1000 lines per page (`const pageSize = 1000`).
- Heatmap: pre-computed cells, never full file content.

### 13.4 Caching strategy

| Data | Cache location | TTL | Invalidation |
|------|----------------|-----|--------------|
| Dashboard widgets | in-memory | 30s | `.caf/` watcher fires |
| Decision detail markdown | `~/.caf/cache/decisions/` | 7d | `.caf/decisions/` watcher fires |
| Sessions index | in-memory | session-lifetime | (rebuilt on `caf sync`) |
| Heatmap | in-memory | 5m | `git diff` since last compute |
| Adapter capabilities | in-memory | process-lifetime | (never changes within process) |

---

## 14. Testing Strategy

### 14.1 Unit tests

Every screen has `*_test.go` with snapshot tests using `termenv` profile. Mock `DataLoader` interface.

### 14.2 Integration tests

Real `.caf/` fixtures; full TUI flow via [teatest](https://github.com/charmbracelet/teatest). CI with terminal sizes 60/100/140 cols.

### 14.3 Performance tests

```bash
# CI benchmarks (run on every PR)
benchmark: `caf dashboard` first paint <100ms on 1000-session repo
benchmark: `caf review` card render <50ms
benchmark: `caf heatmap` render <200ms on 1000-file repo
# Regressions >2x trigger alert, PR blocked
```

### 14.4 Visual regression tests

Snapshots in `pkg/ui/__snapshots__/`. Snapshot diff in CI; manual review for intentional changes.

---

## 15. Future UI Enhancements (not in v0.3)

| Enhancement | Where (BACKLOG ref) | Phase | Status |
|-------------|---------------------|-------|--------|
| VSCode extension with inline provenance | B-001 / K-003 (duplicate) | FUTURE | proposed |
| JetBrains plugin | B-002 / K-004 (duplicate) | FUTURE | idea |
| Neovim plugin | B-003 / K-005 (duplicate) | FUTURE | draft |
| Full-screen review TUI | B-004 | BETA | ✅ done (shipped in v0.3) |
| Onboarding wizard | B-005 | MVP | ✅ done (shipped in v0.3) |
| Color-coded `caf report` | B-008 | BETA | ✅ done (shipped in v0.3) |
| Terminal dashboard | B-009 | MVP | ✅ done (shipped in v0.3) |
| Multi-repo dashboard | B-010 / K-001 (duplicate) | PROD | proposed |
| Web dashboard | (FUTURE) | FUTURE | — |
| Voice-first review | D-006 / K-006 (duplicate) | FUTURE | idea |
| Mobile app | R-004 | rejected | rejected |

> **v0.3.1-triz reconciliation:** In the original UI.md §15, the "Status" column was missing; items B-004, B-005, B-008, B-009 were tagged inline as "already in v0.3! done" but the status field was inconsistent. v0.3.1-triz adds an explicit Status column. Cross-theme duplicates (K-003/4/5/6 = B-001/2/3 and D-006) are now explicitly tagged as duplicates per BACKLOG §12 cross-reference.
