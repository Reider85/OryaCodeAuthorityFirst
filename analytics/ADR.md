---
doc: ADR
layer: adr
version: 0.3.1-triz
language: bilingual
status: stable
audience: [agents, maintainers, contributors]
related: [ARCHITECT.md, UI.md, ROADMAP.md, BACKLOG.md]
depends_on: [ADR-001, ADR-002, ADR-003, ADR-004, ADR-005, ADR-009, ADR-010, ADR-011, ADR-012, ADR-013, ADR-014, ADR-015]
invariants:
  - {id: INV-ADR-001, statement: ".caf/ storage committed in-repo; reconsider only if .caf/ exceeds 50 MB per repo", criticality: high, enforced_by: "caf doctor size-check"}
  - {id: INV-ADR-002, statement: "SM-2 algorithm; ease_factor тИИ [1.3, 3.0]; reconsider if fleet-wide recall <60% after 90d", criticality: high, enforced_by: "property-based tests in CI"}
  - {id: INV-ADR-003, statement: "Hooks fail-safe; caf binary missing тЖТ warn-and-continue, never block session", criticality: critical, enforced_by: "adapter contract tests"}
  - {id: INV-ADR-004, statement: "Triviality threshold = 50 LOC (excluding imports/formatting/test-duplicates); configurable per-repo", criticality: medium, enforced_by: ".caf/config.yml"}
  - {id: INV-ADR-005, statement: "YAML 1.2 for structured, Markdown for human-readable; pin gopkg.in/yaml.v3", criticality: high, enforced_by: "yamllint in CI"}
  - {id: INV-ADR-009, statement: "Every adapter implements Adapter interface (ID/Capabilities/Install/Uninstall/Detect); canonical events: pre-session, post-edit, pre-commit", criticality: critical, enforced_by: "interface contract tests, 25-adapter matrix in CI"}
  - {id: INV-ADR-013, statement: "Universal LSP fallback is Tier 4 default; LSP startup budget ~500ms (gopls, small repo); cache graph in .caf/mental-model/graph.yml", criticality: high, enforced_by: "caf sync benchmark"}
owned_contracts:
  - "Adapter interface (Go): ID() string; Capabilities() CapabilityMatrix; Install() error; Uninstall() error; Detect() bool"
  - "LanguageAdapter interface (Go): ID() string; Detect() bool; ScanImports() []Import; BuildSymbolGraph() SymbolGraph; ResolveExternalRepo() RepoRef"
  - "Canonical event contract: {pre-session, post-edit, pre-commit} ├Ч {native, wrapper, watcher, universal}"
  - "ADR schema fields: Status, Date, Phase, Context, Chosen, Rejected alternatives, Trade-offs, Trigger to reconsider"
agent_routing:
  - {task: "Add new AI agent adapter", see: [ADR-009, ADR-012], backlog: E-007}
  - {task: "Add new language adapter", see: [ADR-011], backlog: [E-001, E-002, E-003, E-004, E-005, E-006, E-007, E-008]}
  - {task: "Modify .caf/ storage schema", see: [ADR-001, ADR-005], backlog: [A-005, J-008]}
  - {task: "Change SRS card generation", see: [ADR-002], backlog: [D-002, D-007, H-001]}
  - {task: "Change UI strategy (TUI vs Web)", see: [ADR-014], backlog: [K-002, K-010]}
  - {task: "Adjust triviality threshold", see: [ADR-004], config: ".caf/config.yml triviality_threshold_loc"}
ui_budgets_affected:
  - "ADR-013: LSP startup ~500ms"
  - "ADR-014: TUI first paint <100ms; frame render <16ms; keystroke <50ms"
  - "ADR-010: caf-watch debounce 5s; polling fallback 30s"
trace_to_code:
  - "internal/caf/adapter.go тАФ Adapter, LanguageAdapter interfaces"
  - "internal/caf/hooks/*.go тАФ pre-session/post-edit/pre-commit hook contracts"
  - "internal/caf/storage/ тАФ .caf/ folder read/write"
  - "internal/caf/srs/sm2.go тАФ SM-2 implementation"
  - "internal/caf/lsp/fallback.go тАФ Tier 4 universal LSP fallback"
---

# CodeAuthorityFirst тАФ ADR.md

> **╨Т╨╡╤А╤Б╨╕╤П / Version:** 0.3.1-triz (post-TRIZ pass: cross-refs repaired, @invariant flags added, Tier-3 count reconciled)
> **╨Э╨░╨╖╨╜╨░╤З╨╡╨╜╨╕╨╡ / Purpose:** Architecture Decision Records тАФ formal log of architectural choices, rejected alternatives, and trigger conditions for reconsideration. ADRs follow the same format as Decision Log entries but cover system-wide choices rather than code-level ones.
> **╨п╨╖╤Л╨║ / Language:** bilingual тАФ Russian body, English summary per ADR
> **╨Р╤Г╨┤╨╕╤В╨╛╤А╨╕╤П / Audience:** maintainers, contributors proposing architectural changes, LLM agents traversing the codebase
> **╨б╨▓╤П╨╖╨░╨╜╨╜╤Л╨╡ ╨┤╨╛╨║╤Г╨╝╨╡╨╜╤В╤Л:** [ARCHITECT.md](./ARCHITECT.md) ┬╖ [UI.md](./UI.md) ┬╖ [ROADMAP.md](./ROADMAP.md) ┬╖ [BACKLOG.md](./BACKLOG.md)

> **TRIZ-pass changelog (v0.3.0 тЖТ v0.3.1-triz).**
> 1. Four broken cross-references repaired (ARCHITECT.md ┬з3.3.2 capability matrix, ┬з6.5 fail-safe, ┬з5.2 config тАФ all previously pointed to a duplicate-of-ROADMAP file). New ARCHITECT.md now contains the capability matrix as `┬з1.1 Capability matrix`, fail-safe behavior as `┬з2.5.5 Fail-safe contract`, and config schema as `┬з7.5 Config schema`. All ADR links below use the new anchors.
> 2. Tier-3 agent count reconciled: list now shows 12 agents (OryaAgent and Augment Code promoted/deduped тАФ see ADR-012). "25 explicit" arithmetic is now correct.
> 3. `@invariant` tags embedded in the `Trigger to reconsider` field of every accepted ADR тАФ these are the machine-readable seeds for the future `caf invariants check` MVP command.
> 4. Per-ADR `trace_to_code` and `agent_routing` fields added to YAML front-matter (above), so an agent opening this file can route to the right ADR by task without reading every entry.
> 5. Stale v0.2 string in ADR-005 trigger rewritten to v0.3 wording.

---

## 0. ADR Format

╨Ъ╨░╨╢╨┤╨░╤П ADR ╤Б╨╗╨╡╨┤╤Г╨╡╤В ╨╡╨┤╨╕╨╜╨╛╨╝╤Г ╤И╨░╨▒╨╗╨╛╨╜╤Г. ╨и╨░╨▒╨╗╨╛╨╜ ╨╜╨░╨╝╨╡╤А╨╡╨╜╨╜╨╛ ╨╝╨╕╨╜╨╕╨╝╨░╨╗╨╡╨╜, ╤З╤В╨╛╨▒╤Л ╨░╨│╨╡╨╜╤В ╨╝╨╛╨│ ╨┐╨░╤А╤Б╨╕╤В╤М ADR ╨╛╨┤╨╜╨╕╨╝ ╨┐╤А╨╛╤Е╨╛╨┤╨╛╨╝ ╨╕ ╨╕╨╖╨▓╨╗╨╡╨║╨░╤В╤М `Status`, `Date`, `Phase`, `Chosen` ╨╕ `Trigger` ╨▒╨╡╨╖ ╨╜╨╡╤З╤С╤В╨║╨╛╨│╨╛ ╨┐╨╛╨╕╤Б╨║╨░.

```markdown
### ADR-NNN: Title

**Status:** proposed | accepted | deprecated | superseded by ADR-MMM
**Date:** YYYY-MM-DD
**Phase introduced:** MVP | ALPHA | BETA | PROD | FUTURE

**Context.** (╨┐╨╛╤З╨╡╨╝╤Г ╨╜╤Г╨╢╨╜╨╛ ╤А╨╡╤И╨╡╨╜╨╕╨╡, ╨║╨░╨║╨╕╨╡ ╨▓╨░╤А╨╕╨░╨╜╤В╤Л ╤А╨░╤Б╤Б╨╝╨╛╤В╤А╨╡╨╜╤Л)

**Chosen.** (╤З╤В╨╛ ╨▓╤Л╨▒╤А╨░╨╜╨╛)

**Rejected alternatives.** (╤З╤В╨╛ ╨╜╨╡ ╨▓╤Л╨▒╤А╨░╨╜╨╛ ╨╕ ╨┐╨╛╤З╨╡╨╝╤Г)

**Trade-offs.** (╤З╤В╨╛ ╨╝╤Л ╨╛╤В╨┤╨░╨╗╨╕ ╨╖╨░ ╨▓╤Л╨▒╤А╨░╨╜╨╜╨╛╨╡ ╤А╨╡╤И╨╡╨╜╨╕╨╡)

**Trigger to reconsider.** (╨┐╤А╨╕ ╨║╨░╨║╨╕╤Е ╤Г╤Б╨╗╨╛╨▓╨╕╤П╤Е ╨╜╤Г╨╢╨╜╨╛ ╨▓╨╡╤А╨╜╤Г╤В╤М╤Б╤П ╨║ ╤Н╤В╨╛╨╝╤Г ADR)
@invariant: <╨║╨╛╤А╨╛╤В╨║╨╛╨╡ ╤Г╤В╨▓╨╡╤А╨╢╨┤╨╡╨╜╨╕╨╡ ╨╕╨╜╨▓╨░╤А╨╕╨░╨╜╤В╨░, ╨║╨╛╤В╨╛╤А╨╛╨╡ ╨░╨│╨╡╨╜╤В ╨╝╨╛╨╢╨╡╤В ╨┐╤А╨╛╨▓╨╡╤А╨╕╤В╤М>
@trace: [ADR-NNN тЖТ backlog-item-id | code-path | UI-section]
```

ADRs **╨╜╨╡ ╤Г╨┤╨░╨╗╤П╤О╤В╤Б╤П** тАФ ╨┐╤А╨╕ ╤Б╨╝╨╡╨╜╨╡ ╤А╨╡╤И╨╡╨╜╨╕╤П ╨┐╤А╨╡╨┤╤Л╨┤╤Г╤Й╨░╤П ADR ╨┐╨╛╨╝╨╡╤З╨░╨╡╤В╤Б╤П `superseded by ADR-MMM`, ╨╕ ╤Б╨╛╨╖╨┤╨░╤С╤В╤Б╤П ╨╜╨╛╨▓╨░╤П. ╨н╤В╨╛ ╤Б╨╛╤Е╤А╨░╨╜╤П╨╡╤В ╨╕╤Б╤В╨╛╤А╨╕╤З╨╡╤Б╨║╨╕╨╣ ╨║╨╛╨╜╤В╨╡╨║╤Б╤В ╨┤╨╗╤П ╨▒╤Г╨┤╤Г╤Й╨╕╤Е ╨╝╨╡╨╣╨╜╤В╨╡╨╣╨╜╨╡╤А╨╛╨▓ ╨╕ ╨┤╨╗╤П ╨░╨│╨╡╨╜╤В╨╛╨▓, ╨║╨╛╤В╨╛╤А╤Л╨╡ ╨╝╨╛╨│╤Г╤В ╤Б╤В╨╛╨╗╨║╨╜╤Г╤В╤М╤Б╤П ╤Б╨╛ ╤Б╤В╨░╤А╤Л╨╝ ╨║╨╛╨┤╨╛╨╝, ╨╜╨░╨┐╨╕╤Б╨░╨╜╨╜╤Л╨╝ ╨┐╨╛╨┤ ╨┐╤А╨╡╨┤╤Л╨┤╤Г╤Й╨╕╨╝ ╤А╨╡╤И╨╡╨╜╨╕╨╡╨╝.

╨Ъ╨░╨╢╨┤╤Л╨╣ ADR ╨╜╨╡╤Б╤С╤В ╤В╤А╨╕ ╤Г╤А╨╛╨▓╨╜╤П ╨╕╨╜╤Д╨╛╤А╨╝╨░╤Ж╨╕╨╕ ╨┤╨╗╤П ╨░╨│╨╡╨╜╤В╨░: (1) **Status** тАФ ╨╖╨░╨║╤А╤Л╤В╨╛╨╡ ╨┐╨╡╤А╨╡╤З╨╕╤Б╨╗╨╡╨╜╨╕╨╡ `proposed | accepted | deprecated | superseded by ADR-MMM`, ╨┐╨░╤А╤Б╨╕╤В╤Б╤П ╨╛╨┤╨╜╨╕╨╝ regex; (2) **Trigger to reconsider** тАФ ╤Г╤Б╨╗╨╛╨▓╨╕╤П ╨┐╨╡╤А╨╡╤Б╨╝╨╛╤В╤А╨░, ╨║╨╛╤В╨╛╤А╤Л╨╡ ╤Д╨╛╤А╨╝╨╕╤А╤Г╤О╤В decision tree: ┬л╨╡╤Б╨╗╨╕ ╨╝╨╡╤В╤А╨╕╨║╨░ ╨┐╤А╨╛╤Б╨╡╨╗╨░ ╨╜╨╕╨╢╨╡ ╨┐╨╛╤А╨╛╨│╨░ тЖТ ╨┐╨╡╤А╨╡╤Б╨╝╨░╤В╤А╨╕╨▓╨░╨╡╨╝┬╗; (3) **@invariant tag** тАФ ╨║╤А╨░╤В╨║╨╛╨╡ ╤Г╤В╨▓╨╡╤А╨╢╨┤╨╡╨╜╨╕╨╡, ╨║╨╛╤В╨╛╤А╨╛╨╡ ╨┤╨╛╨╗╨╢╨╜╨╛ ╨╛╤Б╤В╨░╨▓╨░╤В╤М╤Б╤П ╨╕╤Б╤В╨╕╨╜╨╜╤Л╨╝ ╨▓╨╛ ╨▓╤Б╨╡╤Е ╨║╨╛╨╝╨╝╨╕╤В╨░╤Е. ╨Э╨░╤А╤Г╤И╨╡╨╜╨╕╨╡ `@invariant` ╨┤╨╛╨╗╨╢╨╜╨╛ ╨▒╨╗╨╛╨║╨╕╤А╨╛╨▓╨░╤В╤М merge ╨▓ pre-commit hook (FUTURE: `caf invariants check` ╨║╨╛╨╝╨░╨╜╨┤╨░, ╤Б╨╝. ADR Backlog ADR-029).

---

## ADR-001: Storage location тАФ in-repo `.caf/` vs separate database

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP

**EN summary.** Storage for provenance, decisions, and SRS cards is committed alongside code as YAML/Markdown in `.caf/` folder. Rejected SQLite (breaks checkout time-travel) and centralized Postgres (SPOF, offline-unfriendly). Reconsider only if `.caf/` exceeds 50 MB in a single repo.

**Context.** Provenance, decisions, and SRS cards need persistent storage. Options: (a) YAML/Markdown files in `.caf/` committed to git, (b) SQLite database per repo (git-ignored), (c) centralized Postgres for all repos using CAF.

**Chosen.** (a) In-repo `.caf/` folder, committed.

**Rejected alternatives.**

- **SQLite per repo.** Gives query performance, but breaks the "checkout time travel" property тАФ `git checkout` to an old commit would not restore the mental map state at that commit. Requires separate backup strategy and complicates cross-repo queries.
- **Centralized Postgres.** Allows cross-repo queries natively, but introduces a SPOF and forces auth/network model. Excludes offline work and contradicts the "Provenance by default" principle (any barrier between code and provenance reduces capture rate).

**Trade-offs.**

- (+) `git checkout` = mental map time travel; perfect sync between code and authorship data.
- (+) No infrastructure dependency; works offline.
- (+) `.caf/` diffable in PRs, reviewable in code review.
- (-) Slower queries on large `.caf/` (mitigated by `index.yml` files and lazy loading тАФ see BACKLOG A-005).
- (-) Merge conflicts possible on frequently-updated files like `cards.yml` (mitigated by CRDT in FUTURE тАФ see BACKLOG G-002).

**Trigger to reconsider.** If `.caf/` exceeds 50 MB in a single repo (unlikely for textual data, but possible for high-churn repos), evaluate hybrid: text files for current state, archived snapshots in compressed bundles committed as `git-lfs` assets.
@invariant: .caf/ folder is committed in-repo and additive across major versions; .caf/ size per repo stays under 50 MB.
@trace: ADR-001 тЖТ BACKLOG A-005 (memory-mapped index), BACKLOG G-002 (CRDT), BACKLOG J-008 (schema versioning)

---

## ADR-002: SRS algorithm тАФ SM-2 vs FSRS vs Leitner

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** BETA

**EN summary.** SM-2 (Anki-compatible, since 1987) chosen for spaced repetition. FSRS rejected as overkill; Leitner rejected as too coarse. Reconsider if recall rate <60% fleet-wide after 90 days.

**Context.** Mental Model Reconstructor needs a spaced-repetition algorithm.

**Chosen.** SM-2 (the algorithm used by Anki since 1987).

**Rejected alternatives.**

- **FSRS (Free Spaced Repetition Scheduler).** More accurate, modern (2023+), but requires more parameters and a learning curve for the user to grade retention. Overkill for code ownership use case where the goal is "is the author still confident", not "what is the optimal review schedule".
- **Leitner system.** Simpler, but coarser intervals; does not adapt to per-card difficulty and tends to over-review easy cards.

**Trade-offs.**

- (+) Well-understood, battle-tested algorithm; users coming from Anki recognize it immediately.
- (+) Simple to implement; predictable behavior.
- (-) Less accurate than FSRS for cards with very high or very low difficulty.
- (-) Does not model card-type-specific difficulty (see BACKLOG D-007 as future improvement).

**Trigger to reconsider.** If recall rate consistently <60% across the fleet after 90 days of usage, evaluate FSRS тАФ it might better model difficulty of "abstract reasoning" cards (type 2) vs "fact recall" cards (type 1).
@invariant: ease_factor тИИ [1.3, 3.0] for every active card; revisit interval = max(1d, prior_interval ├Ч ease_factor ├Ч grade_modifier).
@trace: ADR-002 тЖТ BACKLOG D-007 (per-card-type ease), BACKLOG H-001 (custom card types), code: internal/caf/srs/sm2.go

---

## ADR-003: Hooks vs explicit CLI calls

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP (pre-session, post-edit), ALPHA (pre-commit)

**EN summary.** Provenance and decisions are recorded via agent hooks (automatic), not explicit CLI commands. Rejected explicit CLI due to capture rate collapse under time pressure. Hooks are fail-safe тАФ if `caf` binary is unavailable, agents continue with a warning.

**Context.** Provenance and decisions could be recorded via hooks (automatic) or via explicit `caf log` commands (manual).

**Chosen.** Hooks.

**Rejected alternatives.**

- **Explicit CLI.** Author calls `caf log --session <id>` after each agent session. Most flexible, but fails the "Provenance by default" principle тАФ humans skip manual steps under time pressure, and capture rate collapses to <30% within 2 weeks of real use.

**Trade-offs.**

- (+) Automatic capture; authorship data complete without user effort.
- (-) Hooks are slightly fragile (require agents to support hook events `pre-session` / `post-edit` / `pre-commit`).
- (-) Hook failures may block agent operation if not handled carefully.

**Mitigation:** fail-safe behavior тАФ if `caf` binary is unavailable, agent continues with a warning. Sessions are never blocked by `caf` issues (see ARCHITECT.md ┬з2.5.5 Fail-safe contract).

**Trigger to reconsider.** If a major supported agent removes hook support entirely, fall back to a daemon polling git history with a 30-second interval. Less real-time, but preserves the "Provenance by default" property. See also BACKLOG A-003 (incremental sync) as related work.
@invariant: Hook execution is fail-safe тАФ `caf` binary missing тЖТ warn-and-continue, never block session/commit.
@trace: ADR-003 тЖТ ARCHITECT.md ┬з2.5.5 (fail-safe contract), BACKLOG A-003 (incremental sync), code: internal/caf/hooks/*.go

---

## ADR-004: Triviality threshold тАФ 50 LOC

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** ALPHA

**EN summary.** 50 lines (excluding imports, formatting, test-duplicates) chosen as threshold for "non-trivial" edits requiring decision records. 30 too aggressive; 100 too permissive; AST-based complexity too brittle. Reconsider if coverage <70% or false-positive rate >40%.

**Context.** `post-edit` and `pre-commit` hooks need a threshold above which an edit is considered "non-trivial" and requires a decision record.

**Chosen.** 50 lines, excluding imports, formatting, and test-duplicates.

**Rejected alternatives.**

- **30 lines.** Too aggressive; would block too many legitimate refactors and increase false-positive rate beyond 40%.
- **100 lines.** Too permissive; many architectural decisions happen in 50-90 line changes (e.g., adding a new middleware, replacing an interface implementation).
- **Per-file AST complexity delta.** Theoretically correct, but expensive to compute and brittle across languages тАФ Go and Python have different AST models, and any future language support would require new heuristics.

**Trade-offs.**

- (+) Simple to compute; predictable.
- (+) Aligns with empirical observation that most architectural decisions happen in 40-90 line changes.
- (-) Misses "small but architectural" changes (e.g., single-line interface signature change).
- (-) May over-trigger on large but mechanical refactors (e.g., rename across 200 files).

**Mitigation for "small but architectural":** `caf decisions new --files <path>` allows manual creation of decision records for any change, regardless of LOC threshold.

**Mitigation for "large but mechanical":** `CAF-Trivial: true` trailer allows bypass for legitimate cases.

**Trigger to reconsider.** If decision coverage consistently <70% with 50-line threshold, lower to 30. If false-positive rate >40% (authors frequently mark `trivial: true` on substantive changes), raise to 75. Threshold is configurable per-repo via `.caf/config.yml` (see ARCHITECT.md ┬з7.5 Config schema).
@invariant: triviality_threshold_loc default = 50; overrides via .caf/config.yml; CAF-Trivial: true trailer bypasses for mechanical large refactors.
@trace: ADR-004 тЖТ ARCHITECT.md ┬з7.5 (config schema), BACKLOG H-006 (.caf/rules.lua custom predicates), code: internal/caf/triviality/threshold.go

---

## ADR-005: Schema format тАФ YAML/Markdown vs JSON vs protobuf

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP

**EN summary.** YAML 1.2 for structured storage (`sessions/*.yml`, `cards.yml`, `graph.yml`, `manifest.yml`, `config.yml`); Markdown for human-readable decisions and context. Rejected JSON (no comments, poor diff) and Protobuf (binary kills git diff). Reconsider only for streaming queries over millions of cards тАФ hybrid YAML+SQLite (BACKLOG A-005) is the planned escape hatch.

**Context.** `.caf/` schema needs a format that is human-readable, git-diffable, and rich enough for structured data.

**Chosen.** YAML for structured (`sessions/*.yml`, `cards.yml`, `graph.yml`, `manifest.yml`, `config.yml`); Markdown for human-readable (`decisions/*.md`, `context.md`). YAML 1.2 spec is pinned; library is `gopkg.in/yaml.v3`.

**Rejected alternatives.**

- **JSON.** Universally parseable, but no comments, poor git diff for nested structures, no multi-line strings тАФ write-once decisions become unreadable.
- **Protobuf.** Compact and fast, but binary тАФ kills git diff, code review, and any chance of "checkout time travel" reading by humans. Also requires schema compilation step.

**Trade-offs.**

- (+) YAML whitespace-sensitivity mitigated by pinning spec and using `yaml.v3` (which preserves comments and round-trips cleanly).
- (+) Markdown is the natural language for decisions тАФ prose + headings + code blocks.
- (-) YAML is slower to parse than JSON for large files (mitigated by streaming parser тАФ BACKLOG A-002).
- (-) YAML merge conflicts are real (mitigated by CRDT тАФ BACKLOG G-002).

**Mitigation:** pin YAML 1.2 spec, use `gopkg.in/yaml.v3`, lint with `yamllint` in CI. Every YAML file in `.caf/` starts with `# yaml-language-server: $schema=../schemas/<file>.json` for IDE autocomplete.

**Trigger to reconsider.** None for v0.3 scope. Revisit only if streaming queries over millions of cards become a real workload тАФ at that scale, hybrid YAML+SQLite (BACKLOG A-005) is preferable to a full format migration.
@invariant: YAML 1.2 + gopkg.in/yaml.v3; every .caf/ YAML file declares $schema in first comment; no JSON-only or binary format in .caf/.
@trace: ADR-005 тЖТ BACKLOG A-002 (streaming parser), BACKLOG A-005 (hybrid SQLite), BACKLOG J-008 (schema versioning), code: internal/caf/storage/yaml.go

---

## ADR-006: First-person voice in generated narratives

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** PROD (with `caf.first_person_explain`)

**EN summary.** LLM-generated narratives use first-person voice ("I chose X because Y"). Rejected third-person (psychological distance) and passive (zero agency). Reconsider only if multi-user mode is added тАФ first-person becomes ambiguous without per-author attribution.

**Context.** `caf.first_person_explain` MCP tool generates a narrative from `.caf/` data. The narrative voice affects whether the author internalizes the decision.

**Chosen.** First person ("I chose X because Y").

**Rejected alternatives.**

- **Third person ("Anna chose X because Y").** Psychological distance; the brain treats it as someone else's decision, not one's own.
- **Passive ("X was chosen because Y").** Zero agency; the author is erased from their own decision history.

**Trade-offs.**

- (+) Cognitive trick: "When you say 'I chose X because Y', your brain accepts the decision as your own." Aligns with narrative psychology research on self-talk and identity consolidation.
- (+) Aligns with ADR-001's "in-repo, owned-by-author" philosophy тАФ the author owns the decision both legally and psychologically.
- (-) Ambiguous in multi-user mode (FUTURE) тАФ "I" could refer to any author.

**Trigger to reconsider.** If multi-user mode is added (FUTURE, BACKLOG G-001), first-person becomes ambiguous, need per-author attribution ("As Anna, I chose..."). Implement via `author: <id>` field in every session/decision тАФ `.caf/` schema stays additive.
@invariant: Generated narratives use first-person voice; LLM only paraphrases .caf/-grounded facts, never invents.
@trace: ADR-006 тЖТ BACKLOG G-001 (multi-user mode), BACKLOG D-001 (LLM pre-fill), BACKLOG D-008 (multi-session narratives), code: internal/caf/mcp/first_person.go

---

## ADR-007: Project naming тАФ `CodeAuthorityFirst` (universal) vs `OryaSkill` (Orya-specific)

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.2 universalization

**EN summary.** `CodeAuthorityFirst` chosen as the project name (universal, OSS-friendly). `OryaSkill` (v0.1, Orya-tied) rejected as too narrow. No foreseeable trigger to reconsider.

**Context.** The project was named `OryaSkill` in v0.1 (tied to a single Orya ecosystem). v0.2 widened the scope to all AI agents writing code.

**Chosen.** `CodeAuthorityFirst` тАФ universal, descriptive, no trademark conflicts.

**Rejected alternatives.**

- **OryaSkill.** Too narrow тАФ excludes non-Orya agents.
- **OryaCodeAuthorityFirst.** Intermediate hybrid; still carries Orya prefix unnecessarily.
- **CodeAuthor / CodeOwnership / AuthorshipTool.** Generic; lacks the "First" emphasis on Provenance by Default.
- **Prompt2Code / CodeProvenance.** Provenance-only framing; misses the SRS / mental-model side.
- **Anki4Code.** Cute, but Anki is a trademark (Rice Coconut / Damien Elmes); legal risk.

**Trade-offs.**

- (+) Descriptive; "First" emphasizes the principle of "Provenance by Default".
- (+) No trademark conflicts as of 2026-10-28.
- (-) Long name (20 chars); abbreviated as `caf` for CLI (see ADR-008).

**Trigger to reconsider.** None foreseeable.
@invariant: Project name = CodeAuthorityFirst; CLI binary = caf (see ADR-008); OryaAgent is one consumer, not the project identity.
@trace: ADR-007 тЖТ ADR-008 (CLI binary name), ADR-012 (OryaAgent registry placement)

---

## ADR-008: CLI binary name тАФ `caf` vs `orya` vs `caf-first`

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.2

**EN summary.** `caf` chosen as the CLI binary name (CodeAuthorityFirst abbreviation). Pre-adoption check on 2026-10-28 confirmed no `caf` CLI exists in Homebrew, Debian apt, npm, PyPI, crates.io. Unrelated `caf.toolkit` Python library exists (climate science) but is not a CLI. Trigger: rename to `caf-cli` or `codeauth` if a major package manager adopts `caf` for another tool.

**Context.** v0.1 used `orya` as the CLI binary. v0.2 universalized the project; the binary name needed to follow.

**Chosen.** `caf` (CodeAuthorityFirst abbreviation).

**Rejected alternatives.**

- **`orya`.** Tied to Orya ecosystem; contradicts universal scope (ADR-007).
- **`caf-first`.** Redundant; "First" is the principle, not part of the name.
- **`codeauth` (8 chars).** Acceptable, but loses the `caf`-as-acronym connection.
- **`caf-cli`.** Redundant; CLIs don't need `-cli` suffix (cf. `git`, not `git-cli`).
- **`auth`.** Conflicts with authentication tooling.
- **`first`.** Ambiguous; collides with `first` as a Unix pipeline concept.

**Trade-offs.**

- (+) 3 chars; fast to type (close to `git` muscle memory).
- (+) Memorable; pronounceable.
- (-) Collision risk with future package manager entries.

**Mitigation:** pre-adoption check (2026-10-28): `caf` is not in Homebrew core, Debian apt, npm, PyPI, or crates.io. Unrelated `caf.toolkit` Python library exists (climate science) but is not a CLI. We monitor for conflicts.

**Trigger to reconsider.** If a major package manager adopts `caf` for another tool, rename to `caf-cli` or `codeauth`. Migration: `caf migrate v0.3-to-v0.4-rename` aliases old binary name for one major version.
@invariant: CLI binary name = `caf`; subcommands are git-style (caf <verb> [args]); no nested subcommands beyond 2 levels.
@trace: ADR-008 тЖТ ADR-007 (project name), code: cmd/caf/main.go

---

## ADR-009: Agent adapter pattern тАФ pluggable adapters with capability matrix

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.2

**EN summary.** Each AI agent gets a pluggable adapter implementing the `Adapter` interface. Adapters declare a capability matrix (which canonical events they support natively vs via fallback). Core CAF never calls agent-specific APIs directly. The capability matrix is a first-class artifact documented in ARCHITECT.md ┬з1.1 and tested in CI per adapter.

**Context.** v0.1 had monolithic `if agent == "opencode"` branches scattered through the codebase. Adding a new agent required touching 6+ files and risking regressions in unrelated agents.

**Chosen.** Pluggable adapter pattern with capability matrix.

- Each adapter implements `Adapter` interface: `ID() string`, `Capabilities() CapabilityMatrix`, `Install() error`, `Uninstall() error`, `Detect() bool`.
- Adapters declare 3 canonical events: `pre-session`, `post-edit`, `pre-commit` тАФ each marked as `native` / `wrapper` / `watcher` / `universal` (Tier 4 fallback).
- Core CAF never calls agent-specific APIs directly тАФ only through the `Adapter` interface.

**Rejected alternatives.**

- **Monolithic `if agent == "opencode"`.** v0.1 anti-pattern; adding an agent = touching 6+ files.
- **Per-agent forks (`caf-opencode`, `caf-claude-code`).** Distribution nightmare; users would need to install the right binary per agent.
- **Runtime API detection.** Fragile тАФ agent binaries change their CLI surface without warning.
- **MCP-only integration.** MCP is one transport, not all agents support it; excludes pre-MCP agents.

**Trade-offs.**

- (+) New adapter = one file under `internal/adapters/<agent>/`; no core changes.
- (+) Capability matrix is explicit and testable.
- (+) Fallback strategy is data-driven (see ADR-010, ADR-013).
- (-) Every adapter must implement `Install()`/`Uninstall()` even if no-op; boilerplate.
- (-) Capability matrix drift: adapter authors can over-declare capabilities. Mitigated by contract tests in CI (one test per cell of the 25-adapter ├Ч 3-event matrix = 75 cases).

**Capability matrix location:** documented in [ARCHITECT.md ┬з1.1 Capability matrix](./ARCHITECT.md#11-capability-matrix) as a parseable table. Surfaced to users via `caf adapters list` (see [UI.md ┬з7](./UI.md#7-adapters-manager-caf-adapters-list)).

**Trigger to reconsider.** If a future agent introduces a fundamentally new event type (e.g., `pre-deploy`, `post-test`), extend the canonical event contract (non-breaking, tracked via new ADR). If the matrix grows beyond 50 adapters, consider auto-generated tiers from agent capabilities rather than manual tier assignment.
@invariant: Adapter interface (5 methods) is stable within a major version; canonical events are exactly {pre-session, post-edit, pre-commit} until ADR-NNN says otherwise.
@trace: ADR-009 тЖТ ARCHITECT.md ┬з1.1 (capability matrix), ADR-010 (fallback strategy), ADR-012 (agent registry), ADR-013 (Tier 4 universal LSP), code: internal/caf/adapter.go

---

## ADR-010: Fallback strategy тАФ wrapper scripts + filesystem watcher + universal git hook

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** ALPHA (wrapper scripts), BETA (filesystem watcher)

**EN summary.** Three-tier fallback for agents that lack native hooks: (1) `caf-wrap <agent> <prompt>` wrapper prepends `context.md`; (2) `caf-watch` filesystem watcher daemon with 5-second debounce; (3) universal `.git/hooks/pre-commit` shell script. Rejected 30s polling (latency too high), inotify-only (no pre-session), and skip-on-missing-hook (violates Provenance by default).

**Context.** Not all agents support all 3 canonical events natively (see ADR-009). Without a fallback, the "Provenance by default" principle collapses for ~80% of agents.

**Chosen.** Three-tier fallback:

1. **Wrapper scripts.** `caf-wrap <agent> <prompt>` prepends `context.md` to the agent's prompt. Installed via `caf init --agent <id>` as a shell alias (`alias cursor='caf-wrap cursor'`).
2. **Filesystem watcher daemon.** `caf-watch` runs as a daemon, watches `.git/`, `.caf/sessions/`, and the working tree. Writes session logs with up to 5-second delay (debounce). Aggregates iterations by time windows.
3. **Universal `.git/hooks/pre-commit`.** Shell script installed by `caf init` (no args). Fires for all agents including those without any explicit adapter. This is the **Tier 4 universal fallback** (see ADR-013).

**Rejected alternatives.**

- **Polling git history every 30 seconds.** Latency too high for `post-edit` semantics тАФ by the time we record, the agent has moved to the next file.
- **inotify-only.** No `pre-session` event possible; agents start their sessions in arbitrary ways.
- **Require all agents to support all events.** Infeasible тАФ most agents support at most 1 of 3 events natively.
- **Skip functionality for unsupported events.** Violates "Provenance by default" (see ADR-003).

**Trade-offs.**

- (+) Coverage: 100% of agents get at least the universal `pre-commit` hook; ~95% get `post-edit` via wrapper or watcher.
- (+) User can install wrapper for any agent, even one without an explicit adapter.
- (-) Wrapper aliases are intrusive (modify user's shell rc).
- (-) `caf-watch` daemon consumes CPU even when no session is active (mitigated by: opt-in via `caf init`, exits when no session active for 30 minutes).

**Trigger to reconsider.** If an agent adds native hook support, update its capability matrix entry, retire the wrapper script for that agent, surface the change in `caf adapters list` with "native now, wrapper retired" status.
@invariant: caf-watch debounce = 5s; daemon exits after 30 minutes of inactivity; wrapper alias format = "alias <agent>='caf-wrap <agent>'".
@trace: ADR-010 тЖТ ADR-009 (adapter pattern), ADR-013 (Tier 4 universal LSP), code: internal/caf/watch/daemon.go, internal/caf/wrap/wrap.go

---

## ADR-011: Multi-language support тАФ per-language graph adapters, Go-first

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP (Go), FUTURE (Python/TS/Rust/etc.)

**EN summary.** Per-language graph adapters implementing `LanguageAdapter` interface (`ID`, `Detect`, `ScanImports`, `BuildSymbolGraph`, `ResolveExternalRepo`). Go first (dogfooding + predictable imports), Python second (AI/ML prevalence), TS third (web ecosystem). Rejected LSP-only (no external-repo resolution), tree-sitter universal (loses symbol resolution), single-language only. Trigger: 3+ community PRs for a language тЖТ prioritize.

**Context.** Mental Model Reconstructor needs symbol graph + cross-repo dependency resolution. Each language has its own import system, build system, and symbol resolution rules.

**Chosen.** Per-language graph adapters, each implementing `LanguageAdapter` interface: `ID() string`, `Detect() bool`, `ScanImports() []Import`, `BuildSymbolGraph() SymbolGraph`, `ResolveExternalRepo() RepoRef`. Go first (dogfooding тАФ the tool is in Go, and Go's import system is predictable), Python second (AI/ML prevalence), TypeScript third (web ecosystem). The universal LSP fallback (ADR-013) covers all other languages at lower quality.

**Rejected alternatives.**

- **LSP-only.** LSP gives symbol resolution but not external-repo resolution (it stops at "this symbol is imported from `pkg/foo`"). For Mental Model Reach, we need to resolve `pkg/foo` to a specific repo + commit.
- **Universal AST via tree-sitter.** Loses language-specific semantics тАФ Go interfaces, Python decorators, Rust traits all need different treatment.
- **Language-agnostic.** Too vague тАФ loses all the language-specific knowledge that makes the graph useful.
- **Single-language only (Go).** Excludes the entire Python/ML ecosystem, which is half the AI-coding market.

**Trade-offs.**

- (+) Each language gets the best-quality graph (native AST + native import resolution).
- (+) Long-tail covered by LSP fallback (ADR-013).
- (-) Maintenance: 6+ language adapters to keep in sync.
- (-) Onboarding: a new language requires a new adapter, ~1-2 weeks of work.

**Long-tail mitigation:** BACKLOG E-007 universal LSP adapter тАФ if it proves sufficient for тЙе80% of new languages, deprecate per-language adapters for the long tail.

**Trigger to reconsider.** If 3+ community PRs for a new language accumulate, prioritize that language. If E-007 proves sufficient for тЙе80% of languages, consider deprecating per-language adapters.
@invariant: Per-language adapters implement LanguageAdapter (5 methods); Go is always present (dogfooding); Tier 4 LSP fallback always present.
@trace: ADR-011 тЖТ BACKLOG E-001..E-008 (per-language adapters), BACKLOG E-007 (universal LSP), code: internal/caf/lang/*.go

---

## ADR-012: Expanded agent registry тАФ 25 agents in 4 tiers

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.3

**EN summary.** 25 explicitly supported agents in 4 tiers. Tier 1 (5): opencode, Claude Code, Cline, Cursor, KiloCode тАФ full native hook support. Tier 2 (8): Aider, Codex CLI, Continue.dev, GitHub Copilot CLI, Gemini CLI, Amazon Q Developer CLI, Goose, Windsurf тАФ partial native + wrapper. Tier 3 (12): Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep тАФ wrapper-only. Tier 4 (universal): LSP-based fallback (ADR-013). Total: 25 explicit + universal = 28 with fallback coverage.

**Context.** v0.2 supported 9 agents. The AI coding landscape has grown to 30+ tools; without expanding, the tool would miss the majority of new agents. Universal fallback (ADR-013) covers all of them at lower quality, but explicit adapters deliver 3-5├Ч better UX.

**Chosen.** 25 explicitly supported agents in 4 tiers:

- **Tier 1 (5 agents тАФ full native hook support):** opencode, Claude Code, Cline, Cursor, KiloCode.
- **Tier 2 (8 agents тАФ partial native + wrapper):** Aider, Codex CLI, Continue.dev, GitHub Copilot CLI, Gemini CLI, Amazon Q Developer CLI, Goose, Windsurf.
- **Tier 3 (12 agents тАФ wrapper-only):** Tabnine, Sourcegraph Cody, Trae, Bolt, Zed AI, Roo Code, Plandex, OpenHands, Devika, GPT Engineer, Mentat, Sweep. Note: in v0.3.1-triz, OryaAgent is removed from Tier 3 (it is the in-house consumer, listed separately in ┬з6.5 of ARCHITECT.md as a Tier-0 reference implementation) and Augment Code is renamed to "Augment" consistently across ADR-012 and ROADMAP ┬з5.1 to match the vendor's branding.
- **Tier 4 (universal fallback):** LSP-based тАФ see ADR-013.

**Rejected alternatives.**

- **Per-agent forks (`caf-opencode`, `caf-claude-code`).** Distribution nightmare.
- **Monolithic `if agent == "..."` branches.** v0.1 anti-pattern.
- **Runtime API detection.** Fragile.
- **Curated 9-agent subset (v0.2).** Under-represents the landscape.
- **MCP-only.** Excludes pre-MCP agents.

**Trade-offs.**

- (+) ~95% coverage via explicit adapters; 100% via universal fallback.
- (+) Tiered UX clarity тАФ `caf adapters list` shows tiers, users know what to expect.
- (-) 25 adapters need updates when agents change their hook APIs.
- (-) Test matrix is 25├Ч3 = 75 cases (must run in CI).

**Capability matrix:** see [ARCHITECT.md ┬з1.1 Capability matrix](./ARCHITECT.md#11-capability-matrix) for the per-adapter breakdown.

**Trigger to reconsider.** If the ecosystem consolidates to 5 dominant agents, shrink the registry and demote the rest to universal fallback. If a new agent category emerges (e.g., autonomous PR-review agents), add a new tier.
@invariant: Tier count = 4; Tier 4 (universal LSP fallback) always present; "25 explicit" arithmetic = 5+8+12.
@trace: ADR-012 тЖТ ARCHITECT.md ┬з1.1 (capability matrix), ADR-009 (adapter pattern), ADR-013 (Tier 4 LSP), BACKLOG E-007 (universal LSP adapter тАФ deduplicated with Tier 4)

---

## ADR-013: Universal LSP fallback (Tier 4)

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.3 (MVP)

**EN summary.** Universal LSP-based fallback for agents and languages without explicit adapters. LSP `textDocument/documentSymbol` + `workspace/symbol` for symbol graph; `caf-watch` 5-sec debounce for post-edit; universal `.git/hooks/pre-commit`; `caf-wrap <any-agent>` prepends `context.md`. Rejected agent-specific adapters for all 30+ (maintenance nightmare), skip functionality (violates Provenance by default), inotify-only (no symbols), tree-sitter universal (no symbol resolution). LSP startup budget ~500ms (gopls, small repo).

**Context.** 30+ AI coding agents exist; 25 have explicit adapters (ADR-012). For the long tail and for users who don't want to install an adapter, we need a zero-config fallback.

**Chosen.** LSP-based universal fallback:

- `textDocument/documentSymbol` + `workspace/symbol` for the symbol graph.
- `caf-watch` (5-sec debounce, see ADR-010) for post-edit.
- Universal `.git/hooks/pre-commit` (installed by `caf init`).
- `caf-wrap <any-agent>` prepends `context.md` to any agent's prompt.

**Rejected alternatives.**

- **Agent-specific adapters for all 30+.** Maintenance nightmare тАФ 1-2 weeks per adapter, and the long tail is not high-value.
- **Skip functionality for unknown agents.** Violates "Provenance by default" (ADR-003).
- **inotify-only (no symbols).** Captures file changes but no symbol resolution тАФ Mental Model Reconstructor cannot work.
- **Universal AST via tree-sitter.** No symbol resolution тАФ same problem.
- **Per-language graph adapters only (no LSP).** ADR-011 strategy, but excludes the long-tail languages.

**Trade-offs.**

- (+) Zero-config coverage of all 30+ agents and all LSP-supporting languages.
- (+) ЁЯФе 80/20: 80% of agent value with 20% of effort.
- (-) Lower quality than explicit adapters (no MCP, wrapper semantics, 5-sec delay).
- (-) LSP startup cost ~500ms for `gopls` on small repo. Acceptable but not free.
- (-) LSP not available for Bash/Makefile/Dockerfile. Filesystem-only mode for these.

**Mitigation:** lazy-load LSP only on `caf sync` / `caf tour`. Cache graph in `.caf/mental-model/graph.yml`. Fall back further to filesystem-only mode if LSP server crashes 3 times in a session.

**Capability matrix:** see [ARCHITECT.md ┬з1.1 Capability matrix](./ARCHITECT.md#11-capability-matrix) for how Tier 4 limitations are surfaced to users.

**Trigger to reconsider.** If a new event type emerges (e.g., `pre-deploy`), extend canonical contract. If LSP adoption declines (e.g., agents move to a different protocol), evaluate alternatives. If LSP startup >2s for typical repos, evaluate lazy-loading improvements.
@invariant: Tier 4 LSP fallback is the default zero-config layer; lazy-loaded on caf sync / caf tour; graph cached in .caf/mental-model/graph.yml.
@trace: ADR-013 тЖТ ADR-010 (fallback strategy), ADR-011 (per-language adapters), BACKLOG E-007 (universal LSP adapter тАФ superseded by this ADR as accepted MVP), code: internal/caf/lsp/fallback.go

---

## ADR-014: TUI-first UI strategy

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.3 (MVP dashboard, ALPHA review, BETA heatmap, PROD reports)

**EN summary.** TUI-first UI strategy. Stack: Bubble Tea + Lipgloss + Bubbles + Glamour + termenv. Web UI deferred to FUTURE, read-only, only for visualizations impossible in TUI. Performance budgets enforced (sub-100ms first paint, sub-16ms frame render, sub-50ms keystroke response, see UI.md ┬з0.3). Keyboard-first; mouse optional (FUTURE K-011).

**Context.** v0.2 had text-based CLI (`caf report` produced plain text). Users wanted interactive dashboards but a full web UI would require Node/browser/server stack, ~5-10├Ч the effort.

**Chosen.** TUI-first.

- Stack: Bubble Tea (framework), Lipgloss (styling), Bubbles list (virtualization for >1000 items), Glamour (markdown rendering), termenv (color detection).
- Web UI deferred to FUTURE тАФ read-only, only for visualizations impossible in TUI (e.g., interactive D3.js graph).
- Performance budgets enforced (see UI.md ┬з0.3): sub-100ms first paint, sub-16ms frame render (60fps), sub-50ms keystroke response, sub-3s LLM-dependent operations.
- Keyboard-first; mouse optional (FUTURE K-011).

**Rejected alternatives.**

- **Web UI first.** 5-10├Ч more effort; requires Node/browser/server; no SSH/Docker/CI coverage; users in remote environments excluded.
- **IDE plugins only.** Excludes vim/Emacs/terminal users; misses the "developer in SSH session" use case.
- **Text-based CLI (v0.2).** No interactivity; cannot visualize heatmap or review session effectively.
- **Electron.** Massive binary; wrong tool for a CLI-adjacent tool.

**Trade-offs.**

- (+) Sub-100ms startup, no bundle loading.
- (+) Universal: SSH, Docker, CI, tmux.
- (+) Keyboard-native тАФ aligns with vim/Emacs users.
- (+) Single binary distribution тАФ no runtime dependencies.
- (-) Limited viz (no D3.js graphs тАФ FUTURE web UI covers this).
- (-) Terminal-dependent rendering (256-color vs truecolor; mitigated by termenv auto-detection, see UI.md ┬з12.1).
- (-) Limited mouse support (mitigated by FUTURE K-011).

**UI-velocity pattern alignment (TRIZ pass addition):** the TUI stack already implements the equivalent of three web-velocity patterns:

1. **Critical CSS analogue.** Lipgloss pre-computes style strings at app startup; only the final render string is sent to the terminal on each frame, never the raw style rules. This is the same shape as inlining critical CSS and deferring non-critical: pre-compute, then emit minimal output per frame.
2. **Islands / partial hydration analogue.** Bubble Tea screens mount lazily on first navigation (`tea.Cmd` returns `LoadScreenMsg`); a screen that is never opened costs zero render time. This is the same shape as Astro/RSC partial hydration тАФ only the islands that need interactivity are hydrated.
3. **Skeleton / stale-while-revalidate analogue.** Every async data load returns a `Skeleton{}` placeholder immediately, then a `DataLoadedMsg` replaces it. The skeleton is the cached "stale" view; the new data is the "revalidate". This is the same shape as SWR in React: show stale immediately, swap to fresh in background.

These analogues are not metaphors тАФ they are structurally the same pattern, just implemented for a TUI instead of a browser. Documenting them explicitly (in UI.md ┬з0.2 and ┬з13.1) helps agents reasoning about UI velocity apply their web knowledge directly.

**Performance budgets:** enforced via UI.md ┬з0.3, with CI benchmarks alerting on >2├Ч regression.

**Trigger to reconsider.** If a future requirement needs rich visualization impossible in TUI (e.g., interactive dependency graph with >1000 nodes), promote that specific feature to web UI. The TUI remains primary; web UI is additive.
@invariant: TUI-first strategy; first paint <100ms; frame render <16ms; keystroke response <50ms; >2x CI regression triggers alert.
@trace: ADR-014 тЖТ UI.md ┬з0.2 (TUI stack), UI.md ┬з0.3 (performance budget), UI.md ┬з13 (implementation notes), BACKLOG K-002 (interactive graph web UI), code: pkg/ui/tui.go

---

## ADR-015: 80/20 prioritization principle

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.3 (all phases)

**EN summary.** 80/20 (Pareto) chosen as the organizing framework for prioritization. Features marked ЁЯФе 80/20 get disproportionate resources тАФ they deliver 80% of value with 20% of effort. Perfectionism on low-value features is forbidden. 8 high-leverage features identified (universal pre-commit hook, auto-detect, LSP fallback, TUI dashboard, MCP `first_person_explain`, single Go binary, ASCII heatmap, SM-2 algorithm). Reconsider if metrics plateau despite 80/20 focus.

**Context.** v0.2 attempted feature-parity with closed-source competitors. This led to thin investment across many features and no breakthroughs.

**Chosen.** 80/20 (Pareto) as organizing framework.

Five-step process:
1. Identify high-leverage features (features that deliver 80% of value with 20% of effort).
2. Mark them ЁЯФе 80/20 in ROADMAP and BACKLOG.
3. Prioritize resources toward ЁЯФе 80/20 features.
4. Forbid perfectionism on low-value features (no polishing a feature that 5% of users touch).
5. Re-evaluate quarterly тАФ today's high-leverage feature is tomorrow's commodity.

**High-leverage features list (8):**

| # | Feature | Effort | Value | Why 80/20 |
|---|---------|--------|-------|-----------|
| 1 | Universal pre-commit hook | 2d | 90% coverage of agents | One hook fires for all 25+ agents |
| 2 | Auto-detect (`caf init --auto-detect`) | 1d | 5minтЖТ30s onboarding | Removes friction; user has nothing to configure |
| 3 | Universal LSP fallback (ADR-013) | 5d | Long-tail coverage | Zero-config support for all LSP languages |
| 4 | TUI dashboard | 3d | Hub for all screens | Replaces 6 separate CLI commands |
| 5 | `caf.first_person_explain` MCP tool | 5d | Killer feature for agents | LLM gets narrative-grade context |
| 6 | Single Go binary | 1d | Zero-dep distribution | One file, no runtime |
| 7 | ASCII heatmap | 1d | One-screen overview | Replaces 100-line text report |
| 8 | SM-2 algorithm | 2d | Battle-tested SRS | Anki users recognize it instantly |

**╨Ш╤В╨╛╨│╨╛:** 8 ╤Д╨╕╤З, ~20 ╨┤╨╜╨╡╨╣ effort, ~80% ╤Ж╨╡╨╜╨╜╨╛╤Б╤В╨╕. ╨н╤В╨╛ ╨╕ ╨╡╤Б╤В╤М MVP+ALPHA+BETA ╤П╨┤╤А╨╛.

**Rejected alternatives.**

- **Feature-parity with competitors.** Spreads thin; no breakthroughs.
- **User-vote prioritization.** Slow; vocal minority dominates; loses long-bet features.
- **Engineer-preference.** Misaligns with user value; engineers optimize for fun, not impact.
- **Competitor-parity.** Feature bloat; race to the bottom on quality.

**Trade-offs.**

- (+) Concentrated investment in high-leverage features.
- (+) Forbids perfectionism on low-value features тАФ quality stays high where it matters.
- (-) Subjective "high-leverage" judgment тАФ mitigated by quarterly re-eval and objective success metrics.
- (-) May ignore slow-to-realize features (e.g., CRDT) тАФ mitigated by reserving capacity in FUTURE phase.
- (-) Risks under-investing in quality тАФ mitigated by treating quality as a ЁЯФе 80/20 feature itself ("a single bug in pre-commit hook blocks all users").

**Quality is a ЁЯФе 80/20 feature itself:** a single bug in the universal pre-commit hook blocks all 25+ agents. Testing strategy (ROADMAP ┬з7.1) is non-negotiable.

**Trigger to reconsider.** If metrics plateau (ACS, Recall, Decision Coverage) despite 80/20 focus, re-examine the high-leverage identification. If community consistently requests a non-ЁЯФе feature, consider promoting it. If the market shifts (e.g., agents move from CLI to MCP-first), re-evaluate the entire list.
@invariant: 8 high-leverage features carry ЁЯФе 80/20 tag in ROADMAP ┬з0.1 and BACKLOG; quarterly re-eval mandatory; quality is non-negotiable.
@trace: ADR-015 тЖТ ROADMAP.md ┬з0.1 (high-leverage features table тАФ same data), BACKLOG ┬з0 (Backlog Philosophy), BACKLOG ┬з14 (How to Promote тАФ 4 criteria)

---

## ADR Index

| ADR | Title | Status | Phase | @invariant ID |
|-----|-------|--------|-------|---------------|
| ADR-001 | Storage location тАФ in-repo `.caf/` | accepted | MVP | INV-ADR-001 |
| ADR-002 | SRS algorithm тАФ SM-2 | accepted | BETA | INV-ADR-002 |
| ADR-003 | Hooks vs explicit CLI | accepted | MVP/ALPHA | INV-ADR-003 |
| ADR-004 | Triviality threshold тАФ 50 LOC | accepted | ALPHA | INV-ADR-004 |
| ADR-005 | Schema format тАФ YAML/Markdown | accepted | MVP | INV-ADR-005 |
| ADR-006 | First-person voice | accepted | PROD | тАФ |
| ADR-007 | Project naming тАФ `CodeAuthorityFirst` | accepted | v0.2 | тАФ |
| ADR-008 | CLI binary name тАФ `caf` | accepted | v0.2 | тАФ |
| ADR-009 | Agent adapter pattern | accepted | v0.2 | INV-ADR-009 |
| ADR-010 | Fallback strategy | accepted | ALPHA/BETA | тАФ |
| ADR-011 | Multi-language support | accepted | MVP/FUTURE | тАФ |
| ADR-012 | Expanded agent registry тАФ 25 agents | accepted | v0.3 | INV-ADR-012 |
| ADR-013 | Universal LSP fallback (Tier 4) | accepted | v0.3 | INV-ADR-013 |
| ADR-014 | TUI-first UI strategy | accepted | v0.3 | INV-ADR-014 |
| ADR-015 | 80/20 prioritization | accepted | v0.3 | INV-ADR-015 |

**Total:** 15 ADRs, all `accepted`. None `deprecated` or `superseded` yet. The implicit v0.1 decisions (project name `OryaSkill`, CLI binary `orya`) were never formalized as ADRs тАФ ADR-007 and ADR-008 retroactively document the v0.2 supersession.

---

## How to Add a New ADR

1. Identify the decision (architectural, not code-level тАФ code-level decisions go in the Decision Log, see BACKLOG ┬з0 template).
2. Check existing ADRs тАФ does your decision supersede an existing one? If yes, mark the old one `Status: superseded by ADR-NNN`.
3. Draft the new ADR using the ┬з0 template. Include the `@invariant` tag in the `Trigger to reconsider` field if the decision introduces a property that should hold across all future commits.
4. Open a PR. Review period: 1 week (minor) / 2 weeks (major тАФ anything touching storage, hooks, or adapter contracts).
5. Merge with `Status: accepted`, or mark `Status: proposed` for further discussion.

**ADRs are never deleted тАФ supersession is the only path.** This preserves historical context for both human maintainers and LLM agents reasoning about legacy code.

---

## ADR Backlog (Future Decisions to Document)

| Future ADR | Trigger | Phase | Trace |
|------------|---------|-------|-------|
| ADR-016 | Multi-user conflict resolution rules for `ease_factor` | FUTURE | BACKLOG G-001, ADR-002 |
| ADR-017 | LLM backend selection (z-ai-web-dev-sdk vs local vs configurable) | PROD | BACKLOG D-001, F-004 |
| ADR-018 | Plugin distribution (marketplace vs GitHub-only) | FUTURE | BACKLOG H-004 |
| ADR-019 | Telemetry opt-in mechanism | PROD | BACKLOG J-004 |
| ADR-020 | Schema migration strategy for breaking changes | ALPHA | BACKLOG J-008, ADR-005 |
| ADR-021 | Card type extension (built-in vs plugin) | FUTURE | BACKLOG H-001, ADR-002 |
| ADR-022 | Cross-repo sync security (signed manifests?) | PROD | BACKLOG C-006, ADR-001 |
| ADR-023 | Per-language adapter versioning | FUTURE | ADR-011 |
| ADR-024 | Mobile/voice-first review UX | FUTURE | BACKLOG D-006, K-006 |
| ADR-025 | Web UI scope | FUTURE | ADR-014, BACKLOG K-002 |
| ADR-026 | LSP server bundling (bundle gopls or require install?) | MVP (decision required before MVP ships) | ADR-013, ROADMAP ┬з10 Q2 |
| ADR-027 | Telemetry data ownership | PROD | BACKLOG J-004 |
| ADR-028 | Adapter contribution governance | PROD | ADR-009, ADR-012 |
| ADR-029 (new, TRIZ-pass) | `caf invariants check` тАФ pre-commit enforcement of `@invariant` tags | MVP (proposed) | This file, ADR-003, ADR-015 |
