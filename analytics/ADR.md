# CodeAuthorityFirst — ADR.md

> **Версия / Version:** 0.2
> **Назначение / Purpose:** Architecture Decision Records — formal log of architectural choices, rejected alternatives, and trigger conditions for reconsideration. ADRs follow the same format as Decision Log entries but cover system-wide choices rather than code-level ones.
> **Язык / Language:** bilingual — Russian body, English summary per ADR
> **Аудитория / Audience:** maintainers, contributors proposing architectural changes
> **Связанные документы:** [ARCHITECT.md](./ARCHITECT.md), [ROADMAP.md](./ROADMAP.md), [BACKLOG.md](./BACKLOG.md)

---

## 0. ADR Format

Каждая ADR следует единому шаблону:

```markdown
### ADR-NNN: Title

**Status:** proposed | accepted | deprecated | superseded by ADR-MMM
**Date:** YYYY-MM-DD
**Phase introduced:** MVP | ALPHA | BETA | PROD | FUTURE

**Context.** (почему нужно решение, какие варианты рассмотрены)

**Chosen.** (что выбрано)

**Rejected alternatives.** (что не выбрано и почему)

**Trade-offs.** (что мы отдали за выбранное решение)

**Trigger to reconsider.** (при каких условиях нужно вернуться к этому ADR)
```

ADRs **не удаляются** — при смене решения предыдущая ADR помечается `superseded by ADR-MMM`, и создаётся новая. Это сохраняет исторический контекст для будущих мейнтейнеров.

---

## ADR-001: Storage location — in-repo `.caf/` vs separate database

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP

**EN summary.** Storage for provenance, decisions, and SRS cards is committed alongside code as YAML/Markdown in `.caf/` folder. Rejected SQLite (breaks checkout time-travel) and centralized Postgres (SPOF, offline-unfriendly). Reconsider only if `.caf/` exceeds 50 MB in a single repo.

**Context.** Provenance, decisions, and SRS cards need persistent storage. Options: (a) YAML/Markdown files in `.caf/` committed to git, (b) SQLite database per repo (git-ignored), (c) centralized Postgres for all repos using CAF.

**Chosen.** (a) In-repo `.caf/` folder, committed.

**Rejected alternatives.**

- **SQLite per repo.** Gives query performance, but breaks the "checkout time travel" property — `git checkout` to an old commit would not restore the mental map state at that commit. Requires separate backup strategy and complicates cross-repo queries.
- **Centralized Postgres.** Allows cross-repo queries natively, but introduces a SPOF and forces auth/network model. Excludes offline work and contradicts the "Provenance by default" principle (any barrier between code and provenance reduces capture rate).

**Trade-offs.**

- (+) `git checkout` = mental map time travel; perfect sync between code and authorship data.
- (+) No infrastructure dependency; works offline.
- (+) `.caf/` diffable in PRs, reviewable in code review.
- (-) Slower queries on large `.caf/` (mitigated by `index.yml` files and lazy loading — see BACKLOG A-005).
- (-) Merge conflicts possible on frequently-updated files like `cards.yml` (mitigated by CRDT in FUTURE — see BACKLOG G-002).

**Trigger to reconsider.** If `.caf/` exceeds 50 MB in a single repo (unlikely for textual data, but possible for high-churn repos), evaluate hybrid: text files for current state, archived snapshots in compressed bundles committed as `git-lfs` assets.

---

## ADR-002: SRS algorithm — SM-2 vs FSRS vs Leitner

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

**Trigger to reconsider.** If recall rate consistently <60% across the fleet after 90 days of usage, evaluate FSRS — it might better model difficulty of "abstract reasoning" cards (type 2) vs "fact recall" cards (type 1).

---

## ADR-003: Hooks vs explicit CLI calls

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP (pre-session, post-edit), ALPHA (pre-commit)

**EN summary.** Provenance and decisions are recorded via agent hooks (automatic), not explicit CLI commands. Rejected explicit CLI due to capture rate collapse under time pressure. Hooks are fail-safe — if `caf` binary is unavailable, agents continue with a warning.

**Context.** Provenance and decisions could be recorded via hooks (automatic) or via explicit `caf log` commands (manual).

**Chosen.** Hooks.

**Rejected alternatives.**

- **Explicit CLI.** Author calls `caf log --session <id>` after each agent session. Most flexible, but fails the "Provenance by default" principle — humans skip manual steps under time pressure, and capture rate collapses to <30% within 2 weeks of real use.

**Trade-offs.**

- (+) Automatic capture; authorship data complete without user effort.
- (-) Hooks are slightly fragile (require agents to support hook events `pre-session` / `post-edit` / `pre-commit`).
- (-) Hook failures may block agent operation if not handled carefully.

**Mitigation:** fail-safe behavior — if `caf` binary is unavailable, agent continues with a warning. Sessions are never blocked by `caf` issues (see ARCHITECT.md §6.5).

**Trigger to reconsider.** If a major supported agent removes hook support entirely, fall back to a daemon polling git history with a 30-second interval. Less real-time, but preserves the "Provenance by default" property. See also BACKLOG A-003 (incremental sync) as related work.

---

## ADR-004: Triviality threshold — 50 LOC

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** ALPHA

**EN summary.** 50 lines (excluding imports, formatting, test-duplicates) chosen as threshold for "non-trivial" edits requiring decision records. 30 too aggressive; 100 too permissive; AST-based complexity too brittle. Reconsider if coverage <70% or false-positive rate >40%.

**Context.** `post-edit` and `pre-commit` hooks need a threshold above which an edit is considered "non-trivial" and requires a decision record.

**Chosen.** 50 lines, excluding imports, formatting, and test-duplicates.

**Rejected alternatives.**

- **30 lines.** Too aggressive; would block too many legitimate refactors and increase false-positive rate beyond 40%.
- **100 lines.** Too permissive; many architectural decisions happen in 50-90 line changes (e.g., adding a new middleware, replacing an interface implementation).
- **Per-file AST complexity delta.** Theoretically correct, but expensive to compute and brittle across languages — Go and Python have different AST models, and any future language support would require new heuristics.

**Trade-offs.**

- (+) Simple to compute; predictable.
- (+) Aligns with empirical observation that most architectural decisions happen in 40-90 line changes.
- (-) Misses "small but architectural" changes (e.g., single-line interface signature change).
- (-) May over-trigger on large but mechanical refactors (e.g., rename across 200 files).

**Mitigation for "small but architectural":** `caf decisions new --files <path>` allows manual creation of decision records for any change, regardless of LOC threshold.

**Mitigation for "large but mechanical":** `CAF-Trivial: true` trailer allows bypass for legitimate cases.

**Trigger to reconsider.** If decision coverage consistently <70% with 50-line threshold, lower to 30. If false-positive rate >40% (authors frequently mark `trivial: true` on substantive changes), raise to 75. Threshold is configurable per-repo via `.caf/config.yml` (see ARCHITECT.md §5.2).

---

## ADR-005: Schema format — YAML/Markdown vs JSON vs protobuf

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP

**EN summary.** YAML for structured data (`sessions/*.yml`, `cards.yml`, `graph.yml`, `manifest.yml`, `config.yml`); Markdown for human-readable documents (`decisions/*.md`, `context.md`). JSON rejected (no comments, poor diff); Protobuf rejected (binary, kills git diff). YAML whitespace-sensitivity mitigated by pinning spec and using `yaml.v3`.

**Context.** `.caf/` files need a serialization format.

**Chosen.** YAML for structured data; Markdown for human-readable documents.

**Rejected alternatives.**

- **JSON.** Strict schema, easy to parse, but no comments and poor diff readability for humans reviewing history.
- **Protobuf.** Compact, schema-enforced, but binary — kills `git diff`, requires compilation step, blocks ad-hoc inspection via `cat`.

**Trade-offs.**

- (+) YAML supports comments, multi-line strings (essential for `initial_prompt` and decision context blocks).
- (+) Markdown for decisions maximizes readability in GitHub PR view.
- (+) Both formats have stable, well-tested Go parsers.
- (-) YAML is whitespace-sensitive and has multiple parsers with subtle incompatibilities.

**Mitigation:** pin to YAML 1.2 spec, use `gopkg.in/yaml.v3` everywhere, lint with `yamllint` in CI.

**Trigger to reconsider.** None — the format is correct for the use case. Only revisit if a fundamentally different access pattern (e.g., streaming queries over millions of cards) emerges, which is out of scope for v0.2. If it ever happens, hybrid approach (YAML for source of truth, SQLite index for queries — see BACKLOG A-005) is preferable to switching the source format.

---

## ADR-006: First-person voice in generated narratives

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** PROD (with `caf.first_person_explain`)

**EN summary.** All generated explanations and decision records use first person ("I chose X because Y"). Third person rejected (psychological distance); passive voice rejected (zero agency). First-person forces cognitive assimilation — the brain accepts the decision as its own. Reconsider only if multi-user mode is added.

**Context.** `caf.first_person_explain` and decision records can be written in third person ("the author chose X") or first person ("I chose X").

**Chosen.** First person.

**Rejected alternatives.**

- **Third person.** More "objective" but creates psychological distance. The author reads "the author chose token bucket" as if reading about someone else, which defeats the assimilation goal.
- **Passive voice.** "Token bucket was chosen because..." Worst of both worlds — distance + zero agency.

**Trade-offs.**

- (+) Cognitive trick: first-person voice forces assimilation. When you say "I chose X because Y", your brain accepts the decision as your own.
- (+) Aligns with narrative psychology research on self-talk and identity consolidation.
- (-) May feel unnatural for decisions authored by AI agent on user's behalf (but that's the point — to reclaim ownership).

**Trigger to reconsider.** If multi-user mode is added (FUTURE phase), first-person voice becomes ambiguous — would need per-author attribution ("As Anna, I chose..."). This is the only trigger; first-person is the correct default for single-author use case.

---

## ADR-007: Project naming — `CodeAuthorityFirst` (universal) vs `OryaSkill` (Orya-specific)

**Status:** accepted (supersedes implicit v0.1 naming)
**Date:** 2026-10-28
**Phase introduced:** v0.2 universalization

**EN summary.** Project renamed from `OryaSkill`/`OryaCodeAuthorityFirst` (v0.1, Orya-stack-specific) to `CodeAuthorityFirst` (v0.2, universal). The tool now supports any codebase authored by any AI agent. Reconsider only if the tool becomes strongly associated with a specific ecosystem.

**Context.** v0.1 of the architecture was scoped to the Orya-stack (`OryaAgent`, `OryaSecurity`, etc.). v0.2 universalizes the tool to work with any codebase developed by any AI agent. The project name needs to reflect this universality.

**Chosen.** `CodeAuthorityFirst` — matches the original concept name from the analytics document; explicitly universal; not tied to any ecosystem.

**Rejected alternatives.**

- **`OryaSkill` (v0.1).** Tied to Orya ecosystem; misleading for non-Orya users.
- **`OryaCodeAuthorityFirst` (intermediate).** Hybrid name; the "Orya" prefix still implies Orya-specific scope.
- **`CodeAuthor` / `CodeOwnership` / `AuthorshipTool`.** Generic; doesn't capture the "first-person authorship" principle.
- **`Prompt2Code` / `CodeProvenance`.** Focuses on provenance module only, missing Decision Log and MMR.
- **`Anki4Code`.** Cute but reduces tool to one module (MMR); also potentially trademark issue.

**Trade-offs.**

- (+) Name reflects the core principle (authorship-first).
- (+) Universal — not tied to any ecosystem.
- (-) Longer than ideal for CLI (`CodeAuthorityFirst` → CLI is `caf`, see ADR-008).

**Trigger to reconsider.** None foreseeable. The name is correct for the universal scope. If the tool ever becomes strongly associated with a specific ecosystem (e.g., adopted as official tool of a major framework), the name may need re-evaluation, but this is hypothetical.

---

## ADR-008: CLI binary name — `caf` vs `orya` vs `caf-first`

**Status:** accepted (supersedes implicit v0.1 `orya`)
**Date:** 2026-10-28
**Phase introduced:** v0.2 universalization

**EN summary.** CLI binary is `caf` (Code Authority First). Short, unique in `$PATH`, easy to type. Rejected `orya` (v0.1, Orya-specific), `caf-first` (redundant), `codeauth` (longer). Reconsider only if a name conflict emerges in major package managers.

**Context.** v0.1 used `orya` as CLI binary name. v0.2 universalization requires a name not tied to Orya ecosystem.

**Chosen.** `caf` — Code Authority First.

**Rejected alternatives.**

- **`orya` (v0.1).** Tied to Orya ecosystem; misleading for universal tool.
- **`caf-first`.** Redundant; "first" is part of project name, not a useful CLI suffix.
- **`codeauth`.** Longer (8 chars vs 3); no clear advantage.
- **`caf-cli`.** Redundant; CLI is implied.
- **`auth`.** Conflicts with too many existing tools; ambiguous.
- **`first`.** Ambiguous out of context.

**Trade-offs.**

- (+) 3 characters — fast to type, hard to typo.
- (+) Unique in `$PATH` based on check against Homebrew, apt, and PyPI.
- (+) Matches `.caf/` folder name (consistent mental model).
- (-) Less self-descriptive than `codeauth`; new users may not immediately understand what it stands for.

**Mitigation:** README and `caf --help` explain the acronym clearly on first invocation.

**Pre-adoption check:** As of 2026-10-28, `caf` is not used by any package in:
- Homebrew core
- Debian apt repositories
- npm registry
- PyPI
- crates.io

There is an unrelated Python library `caf.toolkit` (climate science), but no CLI binary named `caf`. We monitor for conflicts.

**Trigger to reconsider.** If a major package manager (Homebrew, apt) adopts `caf` as binary name for another tool, rename to `caf-cli` or `codeauth` to avoid conflict.

---

## ADR-009: Agent adapter pattern — pluggable adapters with capability matrix

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** v0.2 (universalization)

**EN summary.** Each AI agent gets a pluggable adapter implementing the `Adapter` interface. Adapters declare a capability matrix (which canonical events they support natively). Core CAF uses fallback strategies for missing capabilities. Rejected alternatives: monolithic integration (hard to extend), per-agent forks (maintenance nightmare), runtime API detection (unreliable).

**Context.** v0.2 needs to support 9 AI agents: opencode, Claude Code, Cursor, Cline, Aider, KiloCode, Codex CLI, Continue.dev, OryaAgent. Each has its own hook system, configuration format, and event vocabulary. The architecture needs a way to handle this diversity without duplicating core logic.

**Chosen.** Pluggable adapter pattern with capability matrix.

Each adapter:

1. Implements the `Adapter` interface (`ID()`, `Capabilities()`, `Install()`, `Uninstall()`, `Detect()`).
2. Declares which of 3 canonical events it supports natively vs via fallback.
3. Translates between agent-specific concepts and CAF's canonical event contract.

Core CAF never calls agent-specific APIs directly — always through the adapter.

**Rejected alternatives.**

- **Monolithic integration (single codebase with `if agent == "opencode" { ... } else if agent == "claude-code" { ... }`).** Quickly becomes unmaintainable as agents are added. Mixing concerns; testing requires all agents installed.
- **Per-agent forks (separate `caf-opencode`, `caf-claude-code`, etc.).** Duplicated core logic; divergence risk; users with multiple agents must install multiple binaries.
- **Runtime API detection (probe agent API at runtime).** Unreliable; agents don't expose stable introspection APIs; brittle to version changes.
- **MCP-only integration (no adapter layer; rely on MCP for everything).** Some agents (Cursor, Aider, Codex CLI) don't support MCP natively or have limited support.

**Trade-offs.**

- (+) Adding a new agent = implementing one interface; no core changes.
- (+) Adapters isolated; failures in one adapter don't crash others.
- (+) Capability matrix makes limitations explicit; users know what to expect.
- (-) Adapter interface may need evolution as new agents introduce new event types.
- (-) Fallback strategies add complexity (wrapper scripts, filesystem watcher).

**Mitigation for interface evolution:** ADR-009 itself can be superseded by a new ADR if the interface needs breaking changes. Minor additions (new optional methods) are non-breaking.

**Capability matrix is a first-class artifact:** documented in ARCHITECT.md §3.3.2, tested in CI per adapter, surfaced to users via `caf adapters list`.

**Trigger to reconsider.** If a future agent introduces a fundamentally new event type (e.g., "pre-deploy" for agents that trigger deployments), the canonical event contract may need extension. This would be a non-breaking change (new event type, adapters may opt-in) but should be tracked via a new ADR.

---

## ADR-010: Fallback strategy — wrapper scripts + filesystem watcher + universal git hook

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** ALPHA (wrapper scripts), BETA (filesystem watcher)

**EN summary.** Three fallback strategies for agents with incomplete hook support: (1) wrapper scripts `caf-wrap <agent>` for missing pre-session, (2) filesystem watcher daemon `caf-watch` for missing real-time post-edit, (3) universal `.git/hooks/pre-commit` shell script (works for any agent). Rejected alternatives: polling git history, inotify-only, requiring all agents to support all events.

**Context.** Not all agents support all three canonical CAF events natively (see ARCHITECT.md §3.3.2 capability matrix). The tool needs fallback strategies that preserve the "Provenance by default" principle.

**Chosen.** Three-tier fallback:

1. **Wrapper scripts** for missing `pre-session`: `caf-wrap <agent> <prompt>` runs the agent with `context.md` prepended to the system prompt. Works for agents that support `--system-prompt-file` or accept prompt from stdin.
2. **Filesystem watcher daemon** (`caf-watch`) for missing real-time `post-edit`: watches working directory for changes, writes session logs with up to 5-second delay; aggregates iterations by time windows.
3. **Universal `.git/hooks/pre-commit`** shell script: works for any agent because git hooks are agent-agnostic. Always installed by `caf init`.

**Rejected alternatives.**

- **Polling git history every 30 seconds.** Latency too high; misses in-flight iterations.
- **inotify-only (no wrapper, no daemon).** Doesn't help with `pre-session` (which is before any file change).
- **Require all agents to support all events (specification by fiat).** Not feasible; agents are external dependencies we don't control.
- **Skip functionality for agents with incomplete support.** Violates "Provenance by default" principle; some agents would have no provenance at all.

**Trade-offs.**

- (+) Maximizes coverage; every agent gets at least partial CAF functionality.
- (+) Universal git hook ensures `pre-commit` validation regardless of agent.
- (-) Wrapper scripts add friction (`caf-wrap cursor` instead of `cursor`).
- (-) Filesystem watcher daemon is a long-running process; resource cost.
- (-) Quality of provenance varies by adapter capability (clearly communicated via capability matrix).

**Mitigation for wrapper friction:** `caf init --agent <id>` can install shell aliases (`alias cursor='caf-wrap cursor'`) so users don't change muscle memory.

**Mitigation for daemon resource cost:** `caf-watch` is opt-in; only runs if a configured adapter lacks native `post-edit` support. Exits when no agent session is active.

**Trigger to reconsider.** If a future version of an agent (e.g., Aider adds pre-session hooks natively), the corresponding adapter's capability matrix is updated, and fallback is no longer used for that agent. The fallback infrastructure remains for other agents.

---

## ADR-011: Multi-language support — per-language graph adapters, Go-first

**Status:** accepted
**Date:** 2026-10-28
**Phase introduced:** MVP (Go), FUTURE (Python, TypeScript, Rust, etc.)

**EN summary.** Language support via per-language graph adapters implementing a `LanguageAdapter` interface. Go is first-class in MVP; Python/TypeScript/Rust added in FUTURE. Rejected alternatives: LSP-only (insufficient for cross-repo resolution), universal AST (too complex), language-agnostic (loses precision).

**Context.** `caf sync` needs to build a dependency graph from source code. Different languages have different import systems, package managers, and AST structures.

**Chosen.** Per-language graph adapters implementing `LanguageAdapter` interface:

```go
type LanguageAdapter interface {
    ID() string  // "go", "python", "typescript", etc.
    Detect(repoPath string) bool  // is this a repo of this language?
    ScanImports(repoPath string) ([]Import, error)
    BuildSymbolGraph(repoPath string) (*SymbolGraph, error)
    ResolveExternalRepo(importPath string) (string, error)  // returns repo URL or local path
}
```

Go is the first supported language (MVP through PROD). Python, TypeScript, Rust, Java are FUTURE (see BACKLOG E-001 through E-006).

**Rejected alternatives.**

- **LSP-only (use Language Server Protocol for everything).** LSP standardizes editor-server communication but doesn't standardize "external repo resolution" (mapping an import to a repository URL or local clone). Per-language configuration is still needed; LSP alone is insufficient.
- **Universal AST (parse all languages to a common AST).** Tree-sitter can parse many languages, but symbol resolution (which function is called, which interface is implemented) is language-specific. Universal AST loses precision.
- **Language-agnostic (treat all imports as strings).** Loses semantic information (e.g., struct vs interface vs function); can't generate "what_if_remove" cards.
- **Single-language only (Go forever).** Excludes majority of real-world projects. Universalization (v0.2 goal) requires multi-language support eventually.

**Trade-offs.**

- (+) Each language gets a tailored adapter; maximum precision for that language.
- (+) Adding a language is self-contained; doesn't affect other adapters.
- (+) Polyglot repos can use multiple adapters simultaneously (see BACKLOG E-008).
- (-) Initial development cost per language (est. 2-4 weeks each).
- (-) Long-tail of languages (Ruby, C#, Haskell, etc.) may never get adapters.

**Mitigation for long-tail:** BACKLOG E-007 (universal LSP-based adapter) is a fallback for languages without dedicated adapters. Lower precision, but functional.

**Sequencing rationale:**

- **Go first** because the tool itself is in Go; dogfooding from day one. Also, Go has the most predictable import system.
- **Python second** because of AI/ML ecosystem prevalence; many AI agent users write Python.
- **TypeScript third** because of web ecosystem prevalence.
- **Rust, Java, etc.** follow based on community demand.

**Trigger to reconsider.** If community contributions for a specific language accumulate (e.g., 3+ community PRs for Python adapter), prioritize that language's official support regardless of original sequencing. If LSP-based universal adapter (E-007) proves sufficient for ≥80% of languages, consider deprecating per-language adapters in favor of universal one.

---

## ADR Status Summary

| ADR | Title | Status | Phase |
|---|---|---|---|
| ADR-001 | Storage location | accepted | MVP |
| ADR-002 | SRS algorithm | accepted | BETA |
| ADR-003 | Hooks vs explicit CLI | accepted | MVP/ALPHA |
| ADR-004 | Triviality threshold | accepted | ALPHA |
| ADR-005 | Schema format | accepted | MVP |
| ADR-006 | First-person voice | accepted | PROD |
| ADR-007 | Project naming | accepted | v0.2 |
| ADR-008 | CLI binary name | accepted | v0.2 |
| ADR-009 | Agent adapter pattern | accepted | v0.2 |
| ADR-010 | Fallback strategy | accepted | ALPHA/BETA |
| ADR-011 | Multi-language support | accepted | MVP/FUTURE |

**Total:** 11 ADRs, all `accepted`. None `deprecated` or `superseded` yet (v0.2 is the first universalized version).

---

## How to Add a New ADR

1. **Identify the decision.** If you're making an architectural choice that affects multiple modules, adapters, or future development, it likely warrants an ADR.
2. **Check existing ADRs.** The decision may extend or supersede an existing one. If superseding, mark the old one as `superseded by ADR-NNN` and create a new ADR.
3. **Draft the ADR** using the template in section 0.
4. **Open a PR** with the new ADR appended to this file. PR description should explain why the decision is needed now.
5. **Review period:** 1 week for minor ADRs (e.g., configuration defaults), 2 weeks for major ADRs (e.g., new adapter pattern, schema changes).
6. **Merge with `accepted` status** if approved. Mark as `proposed` if further discussion needed.

ADRs are never deleted — supersession is the only path. This preserves historical context.

---

## ADR Backlog (Future Decisions to Document)

Decisions that will likely need ADRs as the project evolves:

| Future ADR | Trigger | Phase |
|---|---|---|
| ADR-012 | Multi-user conflict resolution rules for `ease_factor` | FUTURE (multi-user mode) |
| ADR-013 | LLM backend selection (z-ai-web-dev-sdk vs local vs configurable) | PROD (with `first_person_explain`) |
| ADR-014 | Plugin distribution mechanism (marketplace vs GitHub-only) | FUTURE (plugin API) |
| ADR-015 | Telemetry opt-in mechanism and privacy guarantees | PROD (telemetry) |
| ADR-016 | Schema migration strategy for breaking changes | ALPHA (schema stabilization) |
| ADR-017 | Card type extension mechanism (built-in vs plugin) | FUTURE (plugin API) |
| ADR-018 | Cross-repo sync security (signed manifests?) | PROD (cross-repo sync) |
| ADR-019 | Per-language adapter versioning | FUTURE (multiple language adapters) |
| ADR-020 | Mobile/voice-first review UX | FUTURE (D-006) |

These are anticipated decisions; actual ADRs are created when the corresponding feature enters the roadmap.
