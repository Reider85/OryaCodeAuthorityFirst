# CodeAuthorityFirst — UI.md

> **Версия / Version:** 0.3
> **Назначение / Purpose:** detailed specification of all user-facing interfaces — primarily TUI (terminal UI), with progressive enhancements for IDE plugins and (future) web UI. Every screen has ASCII mockup, keyboard map, performance budget, and edge cases.
> **Философия / Philosophy:** TUI-first (🔥 80/20), sub-100ms response, keyboard-only operation possible, lazy rendering for large datasets.
> **Язык / Language:** bilingual — Russian body, English summary per section
> **Аудитория / Audience:** UI/UX implementers, contributors adding TUI screens, design reviewers
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md) · [ROADMAP.md](./ROADMAP.md)

---

## 0. UI Philosophy & Stack

### 0.1 Why TUI-first (🔥 80/20)

**EN summary.** TUI (terminal UI) is the primary interface for v0.3. Reasons: (1) 80/20 — TUI delivers 80% of UX value with 20% of web UI effort; (2) Speed — terminal apps start in <50ms, web apps need server + bundle; (3) Universal — works on any terminal, including SSH, Docker, CI; (4) Keyboard-native — matches developer workflow; (5) No dependencies — no Node.js, no browser, no Electron.

TUI-first стратегия для v0.3:

| Критерий | TUI | Web UI |
|---|---|---|
| Time-to-first-paint | <100ms | 500-2000ms |
| Distribution | single binary | needs server or static host |
| SSH/Docker compatibility | ✓ | ✗ |
| Keyboard-first | native | needs extra work |
| Maintenance cost | low | high (frontend deps, browser compat) |
| Coverage of use cases | 80% | 95% |

**Вывод:** TUI покрывает 80% use cases с 20% effort. Web UI отложен до FUTURE (см. §1.4).

### 0.2 TUI stack

| Component | Choice | Reason |
|---|---|---|
| Framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Idiomatic Go, ELM-like architecture, mature, MIT |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Pairs with Bubble Tea, declarative styles |
| List virtualization | [Bubbles list](https://github.com/charmbracelet/bubbles) | Built-in virtualization for >1000 items |
| Markdown rendering | [Glamour](https://github.com/charmbracelet/glamour) | Renders decision records with syntax highlighting |
| Progress bars | Bubbles spinner + progress | Standard patterns |
| Input fields | Bubbles textinput | With autocomplete via custom matcher |
| Tables | Lipglose table | Sortable, paginated |
| Color detection | [termenv](https://github.com/muesli/termenv) | Auto-detects terminal color support |

**Alternatives rejected:**
- `tview` — older, less idiomatic, less modular
- `gocui` — lower-level, more boilerplate
- Raw ANSI escapes — too error-prone across terminals

### 0.3 Performance budget (universal)

Любая команда `caf` должна соблюдать:

| Operation | Budget | Strategy |
|---|---|---|
| First paint (any screen) | <100ms | Render skeleton first, lazy-load data |
| Frame render (scroll, redraw) | <16ms (60fps) | Virtualized lists, cached layouts |
| Keystroke response | <50ms | Direct dispatch, no async for simple actions |
| LLM-dependent operations | <3s (acceptable) | Skeleton + spinner, cache results |
| File operations on 10k-file repo | <1s | Streaming + index, no full load |

Если операция превышает budget — показываем skeleton + spinner, не блокируем UI.

### 0.4 Color system

Terminal-safe ANSI colors. Поддержка: 256-color, truecolor (если терминал поддерживает). Авто-fallback на 16-color.

**Semantic palette:**

```
Background:   terminal default (no override)
Primary:      #61AFEF  (blue)        — primary actions, links
Secondary:    #98C379  (green)       — success, "good" metrics
Warning:      #E5C07B  (yellow)       — partial support, attention
Error:        #E06C75  (red)          — errors, low metrics
Muted:        #5C6370  (gray)         — secondary text, hints
Accent:       #C678DD  (purple)       — highlights, focus
Background2:  #282C34  (dark gray)    — selected items, hot rows
```

В коде:

```go
var (
    Primary   = lipgloss.Color("#61AFEF")
    Success   = lipgloss.Color("#98C379")
    Warning   = lipgloss.Color("#E5C07B")
    Error     = lipgloss.Color("#E06C75")
    Muted     = lipgloss.Color("#5C6370")
    Accent    = lipgloss.Color("#C678DD")
    Selected  = lipgloss.Color("#282C34")
)
```

Если `termenv` detect'ит 16-color terminal → fallback mapping:

```
Primary   → blue
Success   → green
Warning   → yellow
Error     → red
Muted     → white (dim)
Accent    → magenta
Selected  → reverse video
```

### 0.5 Typography

Terminal monospace. Шрифты зависят от terminal, но CAF задаёт:

| Element | Style | Example |
|---|---|---|
| H1 | bold + accent color | **Dashboard** |
| H2 | bold + primary | **Today's Review Queue** |
| H3 | bold | **Pending decisions** |
| Body | regular | Text content |
| Code | regular + background | `internal/api/policies.go` |
| Metadata | muted | 2026-10-24 · 11:02 |
| Highlight | accent underlined | token bucket |
| Truncated | muted `…` suffix | Very long title… |

Word wrap: `wordwrap` из `github.com/muesli/reflow` с сохранением code spans.

### 0.6 Universal keyboard map

| Key | Action | Works in |
|---|---|---|
| `?` | Show help / keyboard overlay | everywhere |
| `q` / `Ctrl+C` | Quit / cancel | everywhere |
| `Esc` | Back / close dialog | everywhere |
| `Tab` / `Shift+Tab` | Next / prev focusable | forms, dashboards |
| `j` / `↓` | Next item | lists, tables |
| `k` / `↑` | Prev item | lists, tables |
| `g` / `G` | Top / bottom | lists, tables |
| `/` | Search / filter | lists, tables |
| `Enter` | Confirm / open | everywhere |
| `Space` | Toggle / mark | lists, decisions |
| `r` | Refresh | dashboards |
| `1`-`9` | Jump to tab N | tabbed screens |
| `:` | Command palette | everywhere (optional) |
| `h` / `l` | Left / right (in panes) | split-view screens |

**Vim-style rationale:** developers know vim, no learning curve. `?` for help is universal (less/vim/man).

---

## 1. Screens Overview

### 1.1 Screen inventory

| # | Screen | Command | Phase | Spec section |
|---|---|---|---|---|
| 1 | Onboarding wizard | `caf init` (interactive) | MVP | §9 |
| 2 | Dashboard | `caf dashboard` | ALPHA | §3 |
| 3 | Review session | `caf review` | BETA | §4 |
| 4 | Decisions list | `caf decisions` | ALPHA | §5.1 |
| 5 | Decision editor | `caf decisions new` | ALPHA | §5.3 |
| 6 | Heatmap | `caf heatmap` | BETA | §6 |
| 7 | Adapters manager | `caf adapters list` | MVP | §7 |
| 8 | File tour | `caf tour <file>` | PROD | §8 |
| 9 | Sessions browser | `caf sessions` | ALPHA | §5.2 |
| 10 | Reports view | `caf report` | PROD | §10 |
| 11 | Settings editor | `caf config edit` | ALPHA | §11 |

### 1.2 Navigation flow

```text
                ┌──────────────┐
                │  caf init    │ (onboarding)
                └──────┬───────┘
                       │
                ┌──────▼───────┐
   ┌───────────►│ caf dashboard │◄───────────┐
   │            └──────┬───────┘            │
   │                   │                    │
   │       ┌───────────┼───────────┐        │
   │       │           │           │        │
   │  ┌────▼────┐ ┌────▼────┐ ┌────▼─────┐  │
   │  │ review  │ │decisions│ │ sessions │  │
   │  └────┬────┘ └────┬────┘ └────┬─────┘  │
   │       │          │            │        │
   │       │     ┌────▼────┐       │        │
   │       │     │ decision │       │        │
   │       │     │  editor  │       │        │
   │       │     └────┬────┘       │        │
   │       │          │            │        │
   │  ┌────▼──────────▼────────────▼─────┐  │
   └──┤         file tour (caf tour)      ├──┘
      └───────────────────────────────────┘
```

### 1.3 Common UI patterns

**Status bar (bottom of every screen):**

```text
[1] Dashboard  [2] Review  [3] Decisions  │  repo: my-backend  │  agent: claude-code  │  ? help  q quit
```

- Tabs по `1`-`9`
- Repo + активный агент в центре
- Help/quit справа

**Header bar (top of every screen):**

```text
CodeAuthorityFirst v0.3.0 │ my-backend-service │ 47 sessions · 12 decisions · ACS 3.8 ↑
```

- Версия CAF
- Repo name
- Summary stats (sessions / decisions / ACS с трендом)

**Empty states:**

Когда данных нет, показываем не пустой экран, а CTA:

```text
┌─────────────────────────────────────────────────┐
│             No decisions yet                    │
│                                                 │
│   Run your first AI agent session, then:       │
│                                                 │
│   $ caf decisions new --files <path>           │
│                                                 │
│   Or wait for pre-commit hook to prompt you    │
│   after the next non-trivial commit.           │
│                                                 │
│   [Enter] Dismiss    [?] Help                   │
└─────────────────────────────────────────────────┘
```

**Loading states:**

Skeleton + spinner. Никогда не блокируем UI:

```text
Loading sessions...

  ████████░░░░░░░░░░░░  42%
  Fetching 247 session files...

  This should take <1s. If it takes longer, run `caf doctor`.
```

**Error states:**

```text
┌─────────────────────────────────────────────────┐
│ ⚠ Failed to load sessions                       │
│                                                 │
│ Error: .caf/sessions/2026-10-24-a1b2c3.yml:    │
│        line 12: invalid YAML                    │
│                                                 │
│ [r] Retry  [d] Doctor  [i] Ignore  [q] Quit    │
└─────────────────────────────────────────────────┘
```

### 1.4 Web UI (FUTURE, not in v0.3)

Web UI отложен. Когда наступит FUTURE, web UI будет:

- Read-only view на `.caf/` данные (без редактирования — TUI для этого)
- Interactive graph visualization (D3.js) — единственное, что нельзя в TUI
- Real-time dashboard для team mode (multi-user)

Web UI НЕ заменяет TUI. TUI остаётся primary interface.

---

## 2. Onboarding Wizard (`caf init` interactive)

### 2.1 Purpose

Первый запуск `caf` в новом репо. Минимум friction, максимум прозрачности.

### 2.2 Flow

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   Welcome to CodeAuthorityFirst v0.3.0                          │
│                                                                 │
│   This wizard will:                                             │
│     1. Detect AI agents active in your repo                    │
│     2. Install universal pre-commit hook                      │
│     3. Install agent-specific adapters (optional)             │
│     4. Run first `caf context --refresh`                      │
│                                                                 │
│   Time: ~30 seconds                                            │
│                                                                 │
│   [Enter] Continue    [q] Cancel                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Step 1: Detection**

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   Scanning for AI agents...                                    │
│                                                                 │
│   Detected:                                                    │
│     ✓ opencode       (.opencode/ found)                       │
│     ✓ Claude Code    (.claude/ found)                          │
│     ✗ Cursor         (not detected)                            │
│     ✗ Aider          (not detected)                            │
│                                                                 │
│   Universal fallback (LSP): will be installed                  │
│                                                                 │
│   [Enter] Continue with detected agents                       │
│   [a] Add another agent manually                              │
│   [s] Skip adapters, universal only                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Step 2: Installation progress**

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   Installing adapters...                                       │
│                                                                 │
│   [✓] Universal pre-commit hook    ████████████████████ 100%   │
│   [✓] Universal LSP fallback       ████████████████████ 100%   │
│   [✓] opencode adapter             ████████████████████ 100%   │
│   [✓] Claude Code adapter          ████████████████████ 100%   │
│   [•] Generating context.md        ████████░░░░░░░░░░░░  42%   │
│                                                                 │
│   [Running] caf context --refresh                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Step 3: First-time setup**

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   ✓ Setup complete!                                            │
│                                                                 │
│   Summary:                                                    │
│     Repo:           my-backend-service                         │
│     Adapters:        opencode, Claude Code, Universal LSP      │
│     Pre-commit hook: installed (validates future commits)    │
│     Sessions logged: 0 (will accumulate as you use agents)   │
│                                                                 │
│   Next steps:                                                 │
│     1. Run an AI agent session normally                       │
│     2. After non-trivial commit, caf will prompt for decision│
│     3. Run `caf dashboard` to see your authorship metrics     │
│                                                                 │
│   [d] Open dashboard    [q] Done                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.3 Edge cases

- **Repo already has `.caf/`:** show "Setup already complete. [r] Re-run  [q] Cancel"
- **No AI agent detected:** show "No agents detected. Universal fallback will track any agent. [c] Continue  [q] Cancel"
- **Permission denied on `.git/hooks/`:** show error with `sudo` suggestion
- **Non-git repo:** show error "CAF requires git. Run `git init` first."

---

## 3. Dashboard (`caf dashboard`)

### 3.1 Purpose

Real-time overview ментального владения кодом. Главная точка входа после onboarding.

### 3.2 Layout

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ CodeAuthorityFirst v0.3.0 │ my-backend-service │ 47 sessions · 12 dec · ACS 3.8 ↑│
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─ Today ──────────────────────────┐  ┌─ Authorship Score ─────────────┐  │
│  │                                  │  │                                 │  │
│  │  Review queue: 7 cards due       │  │   ACS (weekly avg)              │  │
│  │  Pending decisions: 2            │  │                                 │  │
│  │  New blindspots: 3 files        │  │   3.8 ▲                         │  │
│  │                                  │  │   ▲ +0.3 vs last week           │  │
│  │  [r] Review now    [d] Decide   │  │                                 │  │
│  │                                  │  │   Target: ≥3.5  ✓              │  │
│  └──────────────────────────────────┘  │                                 │  │
│                                        │   Recall: 78%  Coverage: 64%  │  │
│  ┌─ Active Agents ─────────────────┐  │                                 │  │
│  │                                 │  │   [s] See detailed stats       │  │
│  │  claude-code    23 sessions     │  └─────────────────────────────────┘  │
│  │  opencode        12 sessions     │                                        │
│  │  aider            8 sessions     │  ┌─ Cross-repo ───────────────────┐  │
│  │  cursor           4 sessions     │  │                                 │  │
│  │                                 │  │  Dependencies:                  │  │
│  │  [a] Adapters manager           │  │   → my-auth-lib (1 stale dec)   │  │
│  └─────────────────────────────────┘  │   → my-telemetry-lib (clean)    │  │
│                                        │                                 │  │
│  ┌─ Recent Decisions ─────────────┐  │  Exposes:                       │  │
│  │                                 │  │   PolicyAPI (consumed by 2)     │  │
│  │  d_2026-10-24-01  Rate limiter │  │                                 │  │
│  │    agent: claude-code           │  │  [x] Run caf sync               │  │
│  │    files: 1                     │  └─────────────────────────────────┘  │
│  │                                 │                                        │
│  │  d_2026-10-23-02  Token format  │                                        │
│  │    agent: aider                 │                                        │
│  │    files: 2                     │                                        │
│  │                                 │                                        │
│  │  [Enter] Open    [/] Search     │                                        │
│  └─────────────────────────────────┘                                        │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ [1] Dashboard  [2] Review  [3] Decisions  [4] Sessions  │ repo: my-backend │ ? help  q quit │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Components

**Top bar:** версия, repo, summary stats с трендом.

**Today widget:**
- Review queue (карточки с `due ≤ сегодня`)
- Pending decisions (нетривиальные коммиты без decision)
- New blindspots (файлы, изменённые за 30 дней, без карточек)
- Quick actions: `r` review, `d` decide

**Authorship Score widget:**
- ACS (weekly avg) с trend arrow
- Target met indicator
- Recall rate, Coverage rate
- `s` открывает detailed stats

**Active Agents widget:**
- Список агентов, активных за 30 дней
- Session count per agent
- `a` открывает adapters manager

**Recent Decisions widget:**
- Последние 5 decisions
- `Enter` открыть, `/` search

**Cross-repo widget:**
- Dependencies (со stale-decision indicators)
- Exposes
- `x` запускает `caf sync`

**Bottom bar:** tabs, repo, help/quit hints.

### 3.4 Keyboard map

| Key | Action |
|---|---|
| `1`-`5` | Switch to widget 1-5 |
| `Tab` | Next widget |
| `r` | Open review screen |
| `d` | Open decisions list |
| `a` | Open adapters manager |
| `s` | Open detailed stats |
| `x` | Run `caf sync` |
| `R` | Refresh dashboard |
| `/` | Search across all data |
| `?` | Help |
| `q` | Quit |

### 3.5 Refresh strategy

- **Auto-refresh:** every 30 seconds (configurable via `.caf/config.yml` `ui.heatmap_refresh_seconds`)
- **Manual refresh:** `R`
- **Refresh only changed widgets:** incremental update, no full re-render
- **File watcher:** `.caf/` changes trigger automatic refresh of affected widgets only

### 3.6 Performance budget

- First paint: <100ms (skeleton + spinner for data)
- Full data load: <500ms
- Auto-refresh: <50ms (only changed widgets)
- Memory: <50MB for dashboard process

### 3.7 Edge cases

- **No sessions yet:** show onboarding CTA in Today widget
- **No decisions yet:** show "Run your first agent session" CTA in Recent Decisions
- **Cross-repo deps not synced:** show "Run `caf sync`" CTA in Cross-repo widget
- **Terminal too narrow (<120 cols):** stack widgets vertically instead of grid

---

## 4. Review Session (`caf review`)

### 4.1 Purpose

Ежедневная SRS-сессия повторения. ~10 минут, 10-15 карточек. Главная точка контакта с MMR.

### 4.2 Layout — Card view

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf review │ Card 3 of 12 │ 4 min elapsed │ ease: 2.5 │ interval: 1d → 3d    │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─ Question ─────────────────────────────────┐  ┌─ Context ──────────────┐ │
│  │                                            │  │                         │ │
│  │  Type 1 · Why                              │  │  File:                  │ │
│  │                                            │  │  internal/middleware/   │ │
│  │  Почему в                                   │  │    ratelimit.go:18     │ │
│  │  internal/middleware/ratelimit.go:18       │  │                         │ │
│  │  выбран token bucket?                     │  │  Decision:              │ │
│  │                                            │  │  d_2026-10-24-01       │ │
│  │  ─────────────────────────────────────    │  │                         │ │
│  │                                            │  │  Sessions:             │ │
│  │  Your answer (type below, then grade):     │  │  • 2026-10-24-a1b2c3    │ │
│  │                                            │  │    (claude-code)       │ │
│  │  ▌                                         │  │  • 2026-11-02-e5f6a7    │ │
│  │                                            │  │    (aider)             │ │
│  │                                            │  │                         │ │
│  └────────────────────────────────────────────┘  └─────────────────────────┘ │
│                                                                              │
│  ┌─ Grade your recall ──────────────────────────────────────────────────────┐ │
│  │                                                                          │ │
│  │  [0] Blackout — no memory           [3] Hard recall, much difficulty    │ │
│  │  [1] Wrong, but felt familiar       [4] Good recall, some hesitation   │ │
│  │  [2] Wrong, but easy to find         [5] Perfect, instant recall      │ │
│  │                                                                          │ │
│  │  [s] Show expected answer    [n] Skip    [q] Quit session              │ │
│  │                                                                          │ │
│  └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Card 3/12 │ ease 2.5 │ next interval: 1d → 3d (if grade 4) │ 4:12 elapsed  │ ? help  q quit │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Card types rendered differently

**Type 1 — Why (above):** blue accent, shows decision ref in context.

**Type 2 — What if remove:**

```text
┌─ Question ────────────────────────────────────┐
│                                                │
│  Type 2 · What if remove                       │
│                                                │
│  Что произойдёт, если удалить                  │
│  policy.Cache из                               │
│  pkg/policy/evaluator.go?                      │
│                                                │
│  ─────────────────────────────────────────    │
│                                                │
│  Your answer (predict the impact):             │
│                                                │
│  ▌                                             │
│                                                │
└────────────────────────────────────────────────┘
```

**Type 3 — Where handled:**

```text
┌─ Question ────────────────────────────────────┐
│                                                │
│  Type 3 · Where handled                       │
│                                                │
│  Где в коде my-backend-service                │
│  обрабатывается событие                        │
│  tool_call_failed?                             │
│                                                │
│  (List the file paths involved, from           │
│   ingestion to retry logic.)                    │
│                                                │
│  ▌                                             │
│                                                │
└────────────────────────────────────────────────┘
```

### 4.4 Expected answer reveal

After `s` (show expected answer):

```text
┌─ Expected answer ─────────────────────────────────────────────────────────┐
│                                                                          │
│  10 req/min, накладные расходы sliding window не оправданы.              │
│  Token bucket: 10 токенов в бакете, пополнение 1 токен/6 секунд.         │
│                                                                          │
│  From decision d_2026-10-24-01:                                          │
│  > Sliding window имеет смысл при 100+ req/sec, у нас на 3 порядка       │
│  > меньше.                                                               │
│                                                                          │
│  Related sessions:                                                       │
│  • 2026-10-24-a1b2c3 (claude-code) — initial implementation             │
│  • 2026-11-02-e5f6a7 (aider) — added metrics                            │
│                                                                          │
│  Now grade your recall:                                                  │
│  [0-5] Grade    [r] Reveal decision full text    [n] Skip               │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 4.5 Session summary (after last card)

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│   Review session complete!                                                   │
│                                                                              │
│   ┌─ Summary ────────────────────────────────────────────────────────────┐  │
│   │                                                                       │  │
│   │   Cards reviewed: 12/12                                               │  │
│   │   Time elapsed: 8m 47s                                               │  │
│   │   Avg per card: 44s                                                  │  │
│   │                                                                       │  │
│   │   Grade distribution:                                                │  │
│   │     [5] ████████████  6 cards  (perfect recall)                      │  │
│   │     [4] ████████      4 cards  (good recall)                         │  │
│   │     [3] ████          2 cards  (hard recall — flagged for review)     │  │
│   │     [2]                                               │  │
│   │     [1]                                               │  │
│   │     [0]                                               │  │
│   │                                                                       │  │
│   │   Recall rate this session: 83% (10/12 grade ≥4)                    │  │
│   │   Recall rate 7-day avg: 78% (target ≥75%)                          │  │
│   │                                                                       │  │
│   └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│   Authorship Confidence Score:                                              │
                                                                              │
│   "Мог бы я написать этот код с нуля?" (1-5): _                             │
│                                                                              │
│   [1-5] Grade    [s] Skip    [q] Done                                       │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 4.6 Keyboard map

| Key | Action |
|---|---|
| `0`-`5` | Grade current card |
| `s` | Show expected answer |
| `n` | Skip card (counts as grade 0) |
| `r` | Reveal decision full text |
| `Tab` | Switch between answer field and context pane |
| `Enter` | Submit answer text (then grades activate) |
| `Backspace` | Edit answer |
| `?` | Help overlay |
| `q` | Quit session (with confirmation if mid-session) |

### 4.7 Performance budget

- Card render: <50ms (pre-computed display, no LLM call)
- LLM call (if `--explain` flag): <3s with skeleton + spinner
- Session save: <100ms (incremental write to `cards.yml`)

### 4.8 Edge cases

- **No cards due:** show "Nothing to review today! Come back tomorrow. [d] Open dashboard"
- **Card references missing decision:** show "Decision d_xxx not found. Run `caf sync`? [y/n]"
- **Session interrupted (Ctrl+C):** save partial progress, prompt to resume on next `caf review`
- **100+ cards due (backlog):** warn user, offer `--max-cards 15` cap

---

## 5. Decisions (`caf decisions`)

### 5.1 Decisions list

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf decisions │ 12 decisions │ 90% coverage │ Filter: all agents · 30 days │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  / Filter: type here to search...                                           │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │ Date       Title                              Agent         Files  │  │
│  ├──────────────────────────────────────────────────────────────────────┤  │
│  │ 2026-10-24 Rate limiter: token bucket      claude-code    1      │  │
│  │            vs sliding window                                       │  │
│  │ 2026-10-23 Token format: JWT vs opaque     aider           2      │  │
│  │ 2026-10-22 Database: Postgres vs MySQL      claude-code    5      │  │
│  │ 2026-10-20 Auth: RBAC vs ABAC               opencode       8      │  │
│  │ 2026-10-18 Cache: in-memory vs Redis        claude-code    3      │  │
│  │ 2026-10-15 Logging: structured vs plain     aider           2      │  │
│  │ 2026-10-14 API: REST vs GraphQL             claude-code    12     │  │
│  │ ...                                                                 │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Filter:                                                                    │
│    [t] Today    [w] Week    [m] Month    [a] All time                      │
│    [g] By agent [f] By file  [s] By status                                  │
│                                                                              │
│  [Enter] Open    [n] New    [/] Search    [g] Gaps    [q] Quit            │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ [3] Decisions │ repo: my-backend │ ? help  q quit                            │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Decision detail view

After `Enter` on a row:

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf decisions │ d_2026-10-24-01                                            │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  # Rate limiter: token bucket vs sliding window                             │
│                                                                              │
│  Date: 2026-10-24  ·  Agent: claude-code  ·  Session: 2026-10-24-a1b2c3     │
│  Status: accepted  ·  Files: 1                                              │
│                                                                              │
│  ────────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  ## Контекст                                                                 │
│                                                                              │
│  POST /v1/policies принимает запросы создания политик. Без ограничения       │
│  любой валидный токен может выполнить DoS, забив очередь обработки.           │
│  Нужно ограничение 10 запросов в минуту на токен.                          │
│                                                                              │
│  ## Выбрано                                                                  │
│                                                                              │
│  Token bucket: 10 токенов в бакете, пополнение 1 токен/6 секунд.            │
│                                                                              │
│  ## Альтернативы                                                             │
│                                                                              │
│  ### Sliding window                                                          │
│  - Что: считать запросы в скользящем окне 60 секунд                          │
│  - Почему отверг: при 10 req/min накладные расходы не оправданы              │
│                                                                              │
│  ### Fixed window                                                            │
│  - Что: считать запросы в фиксированных минутных окнах                       │
│  - Почему отверг: на границе окна возможен всплеск 2x                       │
│                                                                              │
│  ### Redis-based distributed limiter                                        │
│  - Что: вынести счётчик в Redis                                              │
│  - Почему отверг: на данном этапе нет требования multi-instance             │
│                                                                              │
│  ## Импликации                                                               │
│  - Если нагрузка вырастет до 100+ req/sec — пересмотреть в сторону Redis   │
│  - Интерфейс Limiter абстрактен, миграция без боли                          │
│                                                                              │
│  ────────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  Related sessions:                                                           │
│    • 2026-10-24-a1b2c3 (claude-code) — created this decision                │
│    • 2026-11-02-e5f6a7 (aider) — touched same file                         │
│                                                                              │
│  Related cards:                                                              │
│    • c_001 (type 1, why): "Почему выбран token bucket?"                     │
│    • c_005 (type 2, what_if_remove): "Что если удалить rate limiter?"       │
│                                                                              │
│  [e] Edit    [d] Delete (confirm)    [s] Show file    [Esc] Back           │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Decision editor (`caf decisions new`)

**🔥 80/20 TUI:** автозаполнение для `files`, `tags`, `agent` сокращает время с 90s до 25s.

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf decisions new │ d_2026-10-28-01 (auto-generated)                        │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Title:                                                                       │
│  Rate limiter: token bucket vs sliding window                                │
│                                                                              │
│  Files (autocomplete with Tab):                                              │
│  ─────────────────────────────────────────────                              │
│  ✓ internal/middleware/ratelimit.go                                          │
│  ▌ internal/                                                                  │
│  ↓ internal/api/policies.go         (suggested, modified in session)        │
│  ↓ internal/middleware/ratelimit.go (suggested, modified in session)        │
│  ↓ internal/storage/diesel.go        (no changes)                            │
│                                                                              │
│  Tags (Tab to add):                                                          │
│  rate-limiting, middleware, token-bucket                                    │
│                                                                              │
│  ┌─ Контекст ────────────────────────────────────────────────────────────┐  │
│  │ POST /v1/policies принимает запросы создания политик. Без             │  │
│  │ ограничения любой валидный токен может выполнить DoS, забив            │  │
│  │ очередь обработки. Нужно ограничение 10 запросов в минуту             │  │
│  │ на токен.                                                            │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌─ Выбрано ─────────────────────────────────────────────────────────────┐  │
│  │ Token bucket: 10 токенов, пополнение 1/6 сек.                       │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  Альтернативы (press [a] to add):                                           │
│    1. Sliding window                                                        │
│       Что: ...                                                              │
│       Почему отверг: ...                                                    │
│    2. [Add alternative...]                                                  │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────       │
│                                                                              │
│  [Tab] Next field    [a] Add alternative    [s] Save                         │
│  [x] Save as trivial    [Esc] Cancel                                       │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│ Pre-filled from session 2026-10-24-a1b2c3 (claude-code)                     │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Autocomplete sources:**
- `files`: paths modified in current/recent sessions
- `tags`: tags used in previous decisions in this repo
- `agent`: agents active in this repo
- `title`: LLM suggestion if `--ai-prefill` (FUTURE)

### 5.4 Keyboard map

| Key | Action |
|---|---|
| `Tab` | Next field |
| `Shift+Tab` | Prev field |
| `Enter` | New line in text field |
| `Ctrl+S` | Save |
| `Ctrl+X` | Save as trivial |
| `Esc` | Cancel (with confirm) |
| `a` | Add alternative (in alternatives section) |
| `d` | Delete alternative (cursor on alternative) |
| `?` | Help |

### 5.5 Performance budget

- List render: <50ms (virtualized, only visible rows rendered)
- Filter response: <16ms (in-memory index)
- Detail view: <100ms (markdown render via Glamour, cached)
- Autocomplete suggestions: <50ms (pre-built index)

---

## 6. Heatmap (`caf heatmap`)

### 6.1 Purpose

Визуализация владения кодом. Каждый файл — ячейка, цвет — ACS или coverage.

### 6.2 Layout

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf heatmap │ my-backend-service │ 247 files │ sorted by: ownership (worst)│
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Legend:  ░░ 0-20% (blindspot)  ▒▒ 20-50%  ▓▓ 50-75%  ██ 75-100% (owned)   │
│                                                                              │
│  internal/                                                                   │
│    api/          ██ ██ ██ ░░ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██            │
│                  policies.go                                                  │
│                              ↑ worst (0% — never reviewed)                   │
│                                                                              │
│    middleware/   ▓▓ ██ ██ ▒▒ ██ ██ ██ ██                                    │
│                  ratelimit.go                                                │
│                              ↑ 25% — reviewed once, due in 14 days           │
│                                                                              │
│    storage/      ██ ██ ██ ██ ██ ░░ ██ ██                                    │
│                              ↑ 0% — modified but no SRS card                 │
│                                                                              │
│    handlers/     ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██ ██                          │
│                                                                              │
│  pkg/                                                                        │
│    policy/       ██ ██ ▓▓ ██ ██                                              │
│    api/          ██ ██ ██ ██                                                 │
│                                                                              │
│  ────────────────────────────────────────────────────────────────────        │
│                                                                              │
│  Total: 247 files                                                            │
│    ██ Owned (≥75%):       189 files (76%)                                    │
│    ▓▓ Decent (50-75%):     32 files (13%)                                    │
│    ▒▒ Weak (20-50%):       18 files ( 7%)                                    │
│    ░░ Blindspots (0-20%):   8 files ( 3%)                                    │
│                                                                              │
│  [s] Sort by name/ownership/recency    [/] Filter    [r] Refresh            │
│  [Enter] Open file tour    [b] Show blindspots only    [q] Quit            │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Color encoding

```
0-20%    ░░  dark gray      blindspot — needs SRS card
20-50%   ▒▒  red           weak — review overdue or never reviewed
50-75%   ▓▓  yellow        decent — reviewed but interval low
75-100%  ██  green         owned — strong recall, long interval
```

### 6.4 Keyboard map

| Key | Action |
|---|---|
| `j/k` | Navigate files |
| `Enter` | Open file tour |
| `s` | Cycle sort |
| `/` | Filter by path |
| `b` | Show only blindspots |
| `r` | Refresh |
| `q` | Quit |

### 6.5 Performance budget

- Render 247 files: <200ms (pre-computed heatmap, cached)
- Re-compute on file change: <500ms (incremental, only changed files)
- Memory: <20MB (heatmap pre-computed, not full file content)

### 6.6 Edge cases

- **Empty repo:** show "No files yet."
- **Repo with 10k+ files:** render only first 1000 by default, paginated
- **New file (<7 days):** show "new" indicator instead of color
- **Deleted file:** show `✗` marker, fade out over 30 days

---

## 7. Adapters Manager (`caf adapters list`)

### 7.1 Purpose

Показать все доступные адаптеры, их tier, status установки, capabilities.

### 7.2 Layout

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf adapters │ 25 available │ 4 installed in my-backend                    │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Installed (4):                                                              │
│  ───────────────                                                              │
│  ┌────────────┬──────┬─────────────┬─────────────┬─────────────┬──────────┐ │
│  │ Agent      │ Tier │ pre-session │ post-edit   │ pre-commit  │ Status   │ │
│  ├────────────┼──────┼─────────────┼─────────────┼─────────────┼──────────┤ │
│  │ opencode   │  1   │ ✓ native    │ ✓ native    │ ✓ universal │ healthy  │ │
│  │ claude-code│  1   │ ✓ native    │ ✓ native    │ ✓ universal │ healthy  │ │
│  │ aider      │  2   │ ✗ wrapper   │ ✓ native    │ ✓ universal │ healthy  │ │
│  │ LSP fallback│  4  │ ✗ wrapper  │ ✗ watcher   │ ✓ universal │ healthy  │ │
│  └────────────┴──────┴─────────────┴─────────────┴─────────────┴──────────┘ │
│                                                                              │
│  Available (21):                                                             │
│  ───────────────                                                              │
│  ┌────────────┬──────┬─────────────┬──────────────────────────────────────┐ │
│  │ Agent      │ Tier │ Detected?   │ Action                              │ │
│  ├────────────┼──────┼─────────────┼──────────────────────────────────────┤ │
│  │ Cursor     │  1   │ ✗          │ [i] Install                          │ │
│  │ Cline      │  1   │ ✗          │ [i] Install                          │ │
│  │ KiloCode   │  1   │ ✗          │ [i] Install                          │ │
│  │ Codex CLI  │  2   │ ✗          │ [i] Install                          │ │
│  │ Continue   │  2   │ ✗          │ [i] Install                          │ │
│  │ GH Copilot │  2   │ ✗          │ [i] Install                          │ │
│  │ Gemini CLI │  2   │ ✗          │ [i] Install                          │ │
│  │ Amazon Q   │  2   │ ✗          │ [i] Install                          │ │
│  │ Goose      │  2   │ ✗          │ [i] Install                          │ │
│  │ Windsurf   │  2   │ ✗          │ [i] Install                          │ │
│  │ Tabnine    │  3   │ ✗          │ [i] Install                          │ │
│  │ Cody       │  3   │ ✗          │ [i] Install                          │ │
│  │ ...                                                                       │ │
│  └────────────┴──────┴─────────────┴──────────────────────────────────────┘ │
│                                                                              │
│  [d] Detect agents    [r] Refresh    [a] Install all detected              │
│  [h] Health check     [q] Quit                                              │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Keyboard map

| Key | Action |
|---|---|
| `j/k` | Navigate rows |
| `Enter` | Open adapter detail |
| `i` | Install selected adapter |
| `u` | Uninstall selected adapter |
| `d` | Run auto-detect |
| `a` | Install all detected |
| `h` | Run health check |
| `r` | Refresh |
| `q` | Quit |

### 7.4 Health check view

After `h`:

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf adapters health │ Running smoke tests...                                │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  [✓] opencode         pre-session hook fires correctly                     │
│  [✓] opencode         post-edit writes to .caf/sessions/                   │
│  [✓] opencode         config file present at ~/.config/opencode/hooks.toml │
│                                                                              │
│  [✓] claude-code       pre-session hook fires correctly                    │
│  [✓] claude-code       post-edit writes to .caf/sessions/                  │
│  [✗] claude-code       config file missing at .claude/settings.json        │
│      → Reinstall: `caf adapters install claude-code`                       │
│                                                                              │
│  [✓] aider             wrapper script present                               │
│  [✓] aider             wrapper invoked correctly (test run)               │
│  [⚠] aider             .aider.conf.yml uses old format                     │
│      → Migrate: `caf adapters install aider --migrate-config`              │
│                                                                              │
│  [✓] LSP fallback      LSP server started                                  │
│  [✓] LSP fallback      symbol graph built                                  │
│  [✓] LSP fallback      graph saved to .caf/mental-model/graph.yml          │
│                                                                              │
│  ─────────────────────────────────────────────────────────────────────       │
│  Summary: 11 passed, 1 failed, 1 warning                                   │
│                                                                              │
│  [r] Re-run    [f] Fix all    [q] Close                                    │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 7.5 Performance budget

- List render: <50ms
- Health check: <5s per adapter (parallel), <30s total for 25 adapters
- Install: <2s per adapter (parallel up to 5)

---

## 8. File Tour (`caf tour <file>`)

### 8.1 Purpose

Guided экскурсия по конкретному файлу: когда, почему, какими агентами редактировался.

### 8.2 Layout

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf tour │ internal/middleware/ratelimit.go │ 142 lines │ 3 sessions         │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─ Timeline ──────────────────────┐  ┌─ File view ──────────────────────┐ │
│  │                                  │  │                                   │ │
│  │  ● 2026-10-24 (created)          │  │  1  package middleware             │ │
│  │    claude-code                   │  │  2                                 │ │
│  │    "add rate limiter"            │  │  3  import (                       │ │
│  │    +88 lines                     │  │  4    "context"                    │ │
│  │                                  │  │  5    "time"                       │ │
│  │  ● 2026-10-28 (metrics)          │  │  6    "sync"                       │ │
│  │    claude-code                   │  │  7  )                              │ │
│  │    "add prometheus metrics"      │  │  8                                 │ │
│  │    +12 lines                     │  │  9  type Limiter interface {       │ │
│  │                                  │  │ 10    Allow(token string) bool     │ │
│  │  ● 2026-11-02 (refactor)         │  │ 11  }                              │ │
│  │    aider                         │  │ 12                                 │ │
│  │    "extract TokenBucket type"    │  │ 13  type TokenBucket struct {      │ │
│  │    -5 +8 lines                   │  │ 14    tokens chan struct{}         │ │
│  │                                  │  │ 15    rate time.Duration          │ │
│  │  ─────────────────────────       │  │ 16  }                              │ │
│  │                                  │  │ 17                                 │ │
│  │  Decisions:                      │  │ 18  func New(...) *TokenBucket {   │ │
│  │  • d_2026-10-24-01 (token bucket)│  │ 19    // ...                       │ │
│  │                                  │  │ 20    return &TokenBucket{        │ │
│  │  Related cards:                  │  │ 21      tokens: make(chan struct{} │ │
│  │  • c_001 (why: token bucket)     │  │ 22      // ...                      │ │
│  │  • c_005 (what_if_remove)        │  │ 23    }                             │ │
│  │                                  │  │ 24  }                              │ │
│  │  [Enter] Open decision           │  │ 25                                 │ │
│  │  [s] Show full session           │  │ ...                                │ │
│  │  [r] Review related cards         │  │                                   │ │
│  └──────────────────────────────────┘  └───────────────────────────────────┘ │
│                                                                              │
│  [j/k] Scroll file    [Tab] Switch panes    [/] Search in file              │
│  [Enter] Open decision    [s] Open session    [r] Review cards    [q] Quit│
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Keyboard map

| Key | Action |
|---|---|
| `j/k` | Scroll file view |
| `Tab` | Switch between timeline and file panes |
| `Enter` | Open selected decision (when in timeline) |
| `s` | Open selected session details |
| `r` | Start review of related cards |
| `/` | Search in file |
| `?` | Help |
| `q` | Quit |

### 8.4 Performance budget

- First paint: <100ms
- File load: <50ms (cached after first open)
- Timeline build: <50ms (from `.caf/sessions/` index)
- Memory: <10MB (only current file in memory)

### 8.5 Edge cases

- **File never touched by AI agents:** show "This file has no AI session history. [d] Open dashboard"
- **File deleted:** show "File no longer exists in working tree. Showing last known version."
- **Binary file:** show "Cannot tour binary file."
- **File >5000 lines:** paginate, show only first 1000 with "press [n] for next page"

---

## 9. Onboarding Wizard (extended)

See §2 for basic flow. This section adds advanced options.

### 9.1 Advanced mode (`caf init --advanced`)

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   Advanced setup                                                │
│                                                                 │
│   Triviality threshold (LOC): [50]                              │
│   Pre-commit validation: [x] enable                             │
│   Auto-detect agents: [x] yes                                  │
│   Dashboard auto-refresh (sec): [30]                            │
│   Theme: ( ) dark  ( ) light  (•) auto                          │
│                                                                 │
│   [s] Save    [Esc] Cancel                                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 9.2 Migration from v0.1/v0.2

If `.orya/` exists (v0.1) or `.caf/` exists (v0.2):

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   Existing CodeAuthorityFirst setup detected                    │
│                                                                 │
│   Version: 0.2.0  →  0.3.0                                     │
│                                                                 │
│   Migration will:                                               │
│     • Add agent_tier field to existing sessions                │
│     • Add auto_detected field to config.yml                     │
│     • Add ui section to config.yml                              │
│                                                                 │
│   Time: ~5 seconds (247 sessions)                              │
│                                                                 │
│   [m] Migrate    [s] Skip (keep v0.2 format)    [q] Cancel   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 10. Reports View (`caf report`)

### 10.1 Daily report

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf report --day │ 2026-10-28                                                │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Today's activity:                                                           │
│                                                                              │
│  Sessions:     3  (claude-code ×2, opencode ×1)                              │
│  Commits:      2  (both with CAF-Session trailer ✓)                          │
│  Decisions:    1  new (d_2026-10-28-01: "Database indexing strategy")       │
│  Cards due:    7  (5 reviewed, recall rate 80%)                              │
│                                                                              │
│  ────────────────────────────────────────────────────────────────────        │
│                                                                              │
│  What you did well:                                                          │
│    ✓ All non-trivial commits have decision records                          │
│    ✓ Recall rate above target (75%)                                         │
│                                                                              │
│  What to improve:                                                            │
│    ⚠ 1 commit marked trivial (LOC: 87) — consider decision                 │
│    ⚠ 2 blindspots detected in internal/storage/                            │
│                                                                              │
│  Tomorrow's queue:                                                          │
│    • 8 cards due for review                                                 │
│    • 1 cross-repo sync needed (my-auth-lib changed TokenValidator)         │
│                                                                              │
│  [w] Week report    [m] Metrics    [o] Ownership heatmap    [q] Quit      │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Week report

Similar layout, but with 7-day trends visualized as ASCII sparklines:

```text
ACS trend (7 days):   3.5 3.6 3.4 3.7 3.8 3.7 3.8
                     ▁▂▁▃▄▃▄  ↑ +0.3

Recall rate:         72 75 70 78 80 75 80
                     ▃▄▂▆▇▅▇  ↑ +8pp

Sessions per day:    2  3  1  4  3  2  3
                     ▂▃▁▅▄▃▃  → stable
```

### 10.3 Metrics report

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf report --metrics │ 30-day rolling                                        │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─ Authorship ───────────────────────────────────────────────────────────┐ │
│  │                                                                        │ │
│  │  ACS (Authorship Confidence Score)                                    │ │
│  │    Current: 3.8  Target: ≥3.5  ✓                                       │ │
│  │    30-day trend: 3.4 → 3.8  (+0.4, ↑ improving)                      │ │
│  │                                                                        │ │
│  │  Recall Rate                                                          │ │
│  │    Current: 78%  Target: ≥75%  ✓                                       │ │
│  │    30-day trend: 65% → 78%  (+13pp, ↑ improving)                    │ │
│  │                                                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌─ Coverage ────────────────────────────────────────────────────────────┐  │
│  │                                                                        │ │
│  │  Decision Coverage                                                    │ │
│  │    Current: 92%  Target: ≥90%  ✓                                       │ │
│  │    (11/12 non-trivial commits have decisions)                         │ │
│  │                                                                        │ │
│  │  Mental Model Reach                                                   │ │
│  │    Current: 68%  Target: ≥70%  ⚠ below target                        │ │
│  │    Action: run `caf review --regenerate --since 30d`                  │ │
│  │                                                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌─ Agents ──────────────────────────────────────────────────────────────┐  │
│  │                                                                        │ │
│  │  Agent Diversity (Shannon entropy)                                    │ │
│  │    Current: 1.2  Target: ≥1.0  ✓                                       │ │
│  │    Distribution:                                                       │ │
│  │      claude-code   47 sessions  (53%)  ████████████                  │ │
│  │      opencode      18 sessions  (20%)  ████                          │ │
│  │      aider         14 sessions  (16%)  ███                            │ │
│  │      cursor         8 sessions  ( 9%)  ██                             │ │
│  │      LSP fallback   2 sessions  ( 2%)  ▌                              │ │
│  │                                                                        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  [d] Day report    [w] Week report    [o] Ownership heatmap    [q] Quit   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 11. Settings Editor (`caf config edit`)

### 11.1 Layout

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ caf config edit │ .caf/config.yml                                             │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  General:                                                                    │
│    Triviality threshold (LOC): [50]                                          │
│    Pre-commit validation:    [x] enabled                                     │
│    Fail-safe mode:            [x] enabled                                    │
│                                                                              │
│  Adapters:                                                                   │
│    Enabled agents:                                                           │
│      [x] opencode                                                            │
│      [x] claude-code                                                         │
│      [x] aider                                                               │
│      [ ] cursor                                                              │
│      [ ] cline                                                               │
│      ...                                                                     │
│    [a] Add agent    [d] Detect agents                                       │
│                                                                              │
│  UI:                                                                         │
│    Theme:        ( ) dark  ( ) light  (•) auto                               │
│    Review session max cards: [15]                                            │
│    Dashboard refresh (sec):  [30]                                            │
│                                                                              │
│  [Ctrl+S] Save    [Esc] Cancel (with confirm if changes)                    │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 12. Cross-cutting Concerns

### 12.1 Terminal compatibility

| Terminal | Tested | Notes |
|---|---|---|
| iTerm2 (macOS) | ✓ | full color, truecolor, mouse |
| Terminal.app (macOS) | ✓ | 256-color only, no truecolor |
| Alacritty | ✓ | full color, fast |
| Kitty | ✓ | full color, graphics |
| Windows Terminal | ✓ | full color |
| VSCode integrated terminal | ✓ | full color |
| Git Bash (Windows) | partial | limited color, some rendering issues |
| Linux TTY (no X) | partial | 16-color fallback |
| SSH from any terminal | ✓ | color detected remotely |

**Detection:** `termenv` auto-detects. Fallback to no-color if unsure.

### 12.2 Responsive design

| Width (cols) | Layout |
|---|---|
| ≥140 | Full grid (3-column widgets) |
| 100-139 | 2-column grid |
| 60-99 | Single column, stacked |
| <60 | Compact mode (essential info only, hints to widen) |

```text
< 60 cols:
┌──────────────────────────────────┐
│ CAF v0.3  my-backend  ACS 3.8  │
├──────────────────────────────────┤
│ Review queue: 7                  │
│ Pending decisions: 2             │
│ Blindspots: 3                    │
│                                  │
│ [r] Review  [d] Decide  [q] Q  │
└──────────────────────────────────┘
Widen terminal for full dashboard.
```

### 12.3 Accessibility

- **Color-blind friendly:** never use color as the sole encoding. Always pair with text/symbol (`✓`, `⚠`, `✗`).
- **Screen reader hints:** use ANSI escape `\033[?25h` to indicate focus changes.
- **Keyboard-only operation:** every action accessible via keyboard.
- **No time-based-only cues:** flashing/blinking forbidden; transitions are state-based.
- **High contrast mode:** `caf config set ui.theme high-contrast` uses brighter colors with thicker borders.

### 12.4 Internationalization

- Body text: RU/EN depending on user locale (configurable via `.caf/config.yml` `ui.language`)
- CLI command names: always English (universal)
- YAML field names: always English
- Decision record content: whatever language the user writes in

### 12.5 Telemetry (opt-in, FUTURE)

Not in v0.3. When added (BACKLOG J-004):

- Anonymous: only command counts and error rates
- Opt-in: explicit `caf config set telemetry.enabled true`
- No content: no prompts, no file paths, no decisions
- Public dashboard: aggregate stats visible at `caf telemetry stats`

---

## 13. Implementation Notes

### 13.1 File organization

```text
pkg/ui/
├── tui.go                    # TUI root, Bubble Tea program
├── theme.go                  # color palette, typography
├── components/
│   ├── header.go             # top bar
│   ├── statusbar.go          # bottom bar
│   ├── spinner.go            # loading spinner
│   ├── skeleton.go           # skeleton loaders
│   └── heatmap.go            # ASCII heatmap renderer
├── screens/
│   ├── onboarding.go         # §2 / §9
│   ├── dashboard.go          # §3
│   ├── review.go             # §4
│   ├── decisions_list.go     # §5.1
│   ├── decisions_detail.go   # §5.2
│   ├── decisions_editor.go   # §5.3
│   ├── heatmap.go            # §6
│   ├── adapters.go           # §7
│   ├── file_tour.go          # §8
│   ├── reports.go            # §10
│   └── settings.go           # §11
├── keys.go                   # keyboard bindings
└── perf.go                   # performance budgets enforcement
```

### 13.2 State management

Bubble Tea ELM-like architecture:

```go
type model struct {
    activeScreen Screen
    repo         string
    agent        string
    data         DataLoader  // async data fetching
    theme        Theme
    width        int
    height       int
    // ... per-screen state in screens/
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
    case tea.KeyMsg:
        return m.handleKey(msg)
    case DataLoadedMsg:
        return m.handleData(msg)
    }
    return m, nil
}
```

Async data loading via `tea.Cmd`:

```go
func loadData() tea.Msg {
    sessions, _ := caf.LoadSessions(repoPath)
    return DataLoadedMsg{Sessions: sessions}
}
```

### 13.3 Virtualization strategy

For lists with >100 items, use `bubbles/list` with virtualization:

```go
list := list.New(items, delegate, width, height)
list.SetShowHelp(false)
list.SetShowTitle(false)
list.SetFilteringEnabled(true)  // built-in /filter
// Only visible items rendered
```

For file tour with >5000 lines, paginate:

```go
const pageSize = 1000
currentPage := 0
// On 'n': currentPage++; render lines [currentPage*pageSize : (currentPage+1)*pageSize]
```

### 13.4 Caching strategy

| Data | Cache location | TTL | Invalidation |
|---|---|---|---|
| Dashboard widgets | in-memory | 30s | file watcher on `.caf/` |
| Decision detail (rendered markdown) | `~/.caf/cache/decisions/` | 7d | file change |
| Sessions index | in-memory | session-lifetime | new session writes |
| Heatmap | in-memory | 5m | `git diff` since last build |
| Adapter capabilities | in-memory | process-lifetime | never changes during run |

---

## 14. Testing Strategy

### 14.1 Unit tests

- Every screen has `*_test.go` with snapshot tests
- Snapshot tests use [termenv] profile to ensure consistent rendering
- Mock `DataLoader` for predictable test data

### 14.2 Integration tests

- Real `.caf/` directories with fixture data
- Full TUI flow tests via [teatest](https://github.com/charmbracelet/teatest)
- Tests run in CI with various terminal sizes (60, 100, 140 cols)

### 14.3 Performance tests

- Benchmark: `caf dashboard` first paint <100ms on 1000-session repo
- Benchmark: `caf review` card render <50ms
- Benchmark: `caf heatmap` render <200ms on 1000-file repo
- Benchmarks run in CI; regressions >2x trigger alert

### 14.4 Visual regression tests

- Snapshot tests capture rendered output
- Diffs reviewed in PR
- Snapshots stored in `pkg/ui/__snapshots__/`

---

## 15. Future UI Enhancements (not in v0.3)

| Enhancement | Where (BACKLOG ref) | Phase |
|---|---|---|
| VSCode extension | B-001 | FUTURE |
| JetBrains plugin | B-002 | FUTURE |
| Neovim plugin | B-003 | FUTURE |
| Full-screen review TUI | B-004 (already in v0.3!) | done |
| Onboarding wizard | B-005 (already in v0.3!) | done |
| Color-coded reports | B-008 (already in v0.3!) | done |
| Terminal dashboard | B-009 (already in v0.3!) | done |
| Multi-repo dashboard | B-010 | FUTURE |
| Web dashboard | FUTURE section | FUTURE |
| Voice-first review | D-006 | FUTURE |
| Mobile app | rejected (R-004) | — |
