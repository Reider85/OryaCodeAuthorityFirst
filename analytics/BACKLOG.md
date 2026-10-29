---
doc: BACKLOG
layer: backlog
version: 0.4.0-triz2
language: bilingual
status: stable
audience: [agents, maintainers, contributors, anyone proposing a feature]
related: [ARD.md, ARCHITECT.md, UI.md, ROADMAP.md, ADR.md]
depends_on: [ADR-001, ADR-009, ADR-011, ADR-013, ADR-015, ADR-017 (future), ADR-022 (future), ADR-028 (future)]
supercedes: [BACKLOG.md v0.3.1-triz (additive update; no items removed)]
invariants:
  - {id: INV-BL-001, statement: "Every backlog item carries an agent_hint field pointing to the ADR(s) or code contracts it would affect", criticality: high, enforced_by: "caf doc-check (proposed, ADR-029)"}
  - {id: INV-BL-002, statement: "Items already shipped in v0.3 are marked done; duplicate themes are merged with see-also pointers", criticality: medium, enforced_by: "quarterly backlog review"}
  - {id: INV-BL-003, statement: "Cross-references to ROADMAP §6.10 are corrected to §6.9 (the 4 promotion criteria)", criticality: high, enforced_by: "caf doc-check"}
  - {id: INV-BL-004, statement: "Every new item proposed after v0.4.0-triz2 must map to ≥1 ARD requirement (FR-*) or be explicitly marked 'out of current ARD scope'", criticality: high, enforced_by: "PR review checklist"}
agent_routing:
  - {task: "Add a new feature idea", see: [§0 Backlog Philosophy, §16 How to Add a New Backlog Item], template: "§0 Backlog item template"}
  - {task: "Promote backlog item to roadmap", see: [§15 How to Promote], criteria: "4 criteria in ROADMAP.md §6.9"}
  - {task: "Find what's been done", see: [§14 Rejected Items — done items tagged inline]}
  - {task: "Audit backlog health", see: [§17 Statistics]}
  - {task: "Find security/observability/cost items (NEW in v0.4-triz2)", see: [§13 Theme SEC, §14 Theme OBS, §15 Theme COST]}
trace_to_code:
  - "internal/caf/ — every backlog item that proposes implementation work links to the file it would touch"
ui_budgets_affected:
  - "Theme K items affect UI budgets (animated transitions K-007 ≤100ms, themes K-008, mouse K-011)"
  - "Theme A items affect caf sync / caf dashboard budgets"
  - "Theme SEC items add ≤10ms sandbox overhead (INV-ARCH-010)"
  - "Theme OBS items add ≤50ms tracing overhead per command (INV-ARCH-011)"
  - "Theme COST items affect LLM-dependent operations budget (<3s, INV-UI-004)"
---

# CodeAuthorityFirst — BACKLOG.md
# Бэклог / Backlog

> **Версия / Version:** 0.4.0-triz2 (post-TRIZ-2 pass: 3 new themes added — SEC, OBS, COST — totaling 23 new backlog items; TRIZ 80/20 verdicts recomputed across all themes; statistics recomputed; agent_hint fields updated to trace to ARD requirements where applicable)
> **Назначение / Purpose:** catalog of ideas, features, and improvements considered but **not committed to any roadmap phase**. Items here are "parked" — they may be promoted to ROADMAP via the 4 promotion criteria in [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items). See [ADR-015](./ADR.md#adr-015-8020-prioritization-principle) for the 🔥 80/20 prioritization principle.
> **Принцип приоритизации / Prioritization:** 🔥 80/20 — high-leverage items marked; these get disproportionate resources. Resources allocated 60% 🔥 / 30% ◐ / 10% ○ (see ARD §0).
> **Язык / Language:** bilingual — Russian body, English summary per theme
> **Аудитория / Audience:** maintainers, contributors, anyone proposing a feature, LLM agents searching for the next work item
> **Связанные документы:** [ARD.md](./ARD.md) · [ARCHITECT.md](./ARCHITECT.md) · [UI.md](./UI.md) · [ROADMAP.md](./ROADMAP.md) · [ADR.md](./ADR.md)

---

## 0. Backlog Philosophy

Бэклог — это не todo-list, а парковка идей. Каждая запись существует, потому что кто-то её предложил, но она не получила приоритета в текущем roadmap. Это нормально: лучше задокументировать идею и оставить её ждать своего часа, чем потерять её или впихнуть в неподходящий момент.

Чтобы идея жила, она должна иметь:

1. **Ясную проблему.** Не "хорошо бы иметь X", а "сейчас Y происходит, и это блокирует Z".
2. **Эскиз решения.** Минимум 3-4 строки, достаточно для нового контрибьютора чтобы начать.
3. **Связь с ADR или инвариантом.** Каждая идея должна указывать, какой ADR она затрагивает или какой инвариант она вводит/изменяет. Это `agent_hint` поле.
4. **Зависимость от фазы.** Какая фаза (MVP/ALPHA/BETA/PROD/FUTURE) является минимальной для реализации.
5. **Связь с ARD requirement (NEW v0.4-triz2).** Каждая новая идея должна указывать, какой FR-* requirement она реализует или расширяет. Если не maps ни к одному FR — идея "out of current ARD scope" и должна быть либо отклонена, либо добавлена в ARD сначала.

### Backlog item template

```markdown
### [ID] Title 🔥 80/20 (optional)

**Theme:** A | B | C | D | E | F | G | H | I | J | K | SEC | OBS | COST | Rejected
**Status:** idea | draft | proposed | done | rejected
**Origin:** <GitHub issue #, discussion link, or "internal proposal">
**Phase dependency:** MVP | ALPHA | BETA | PROD | FUTURE
**ARD requirement:** FR-<CLUSTER>-<NNN> (NEW v0.4-triz2; or "out of current ARD scope")

**Problem.** (что сейчас не работает или неэффективно)

**Sketch.** (минимум 3-4 строки, как это могло бы работать)

**Open questions.** (что нужно решить до реализации)

**Dependencies.** (другие backlog items, ADRs, фазы)

**agent_hint:** (какой ADR или инвариант затрагивает; что агенту прочитать в первую очередь; какой код-контракт меняется)
```

Критерии промоушена в roadmap — см. [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items) (четыре критерия: ясная проблема, эскиз решения, пилотный пользователь, метрика успеха). Без этих четырёх пунктов item остаётся в BACKLOG.

---

## 1. Themes

| # | Theme | Description | 🔥 80/20 | done (v0.4) | Total items |
|---|-------|-------------|----------|-------------|-------------|
| A | Performance & Scalability | Скорость `caf sync`, парсинг `.caf/`, графы | 2 | 0 | 6 |
| B | UX & DX | IDE-плагины, TUI-сretorno, onboarding | 5 | 5 (B-004/5/6/8/9) | 10 |
| C | Integrations | CI/CD, мессенджеры, таск-трекеры | 2 | 0 | 7 |
| D | AI & LLM Features | LLM-ассистенты, голосовой review | 1 | 0 | 9 |
| E | Languages (beyond Go) | Python, TS, Rust, Java, C#, Ruby, universal | 1 | 0 (E-007 subsumed by ADR-013) | 8 |
| F | Enterprise | SSO, аудит, compliance, on-prem | 0 | 0 | 7 |
| G | Multi-user & Collaboration | Атрибуция, CRDT, тим-аналитика | 0 | 0 | 6 |
| H | Extensibility | Плагины, маркетплейс, hooks | 1 | 0 | 6 |
| I | Analytics & Reporting | Тренды, модули, Prometheus | 2 | 1 (I-007) | 7 |
| J | Quality & Reliability | Self-heal, фаззинг, бэкапы | 2 | 1 (J-003) | 8 |
| K | UI Enhancements | TUI animations, themes, command palette | 3 | 0 | 12 |
| **SEC** | **Security (NEW v0.4-triz2)** | **Threat mitigations, sandboxing, signing** | **3** | **0** | **8** |
| **OBS** | **Observability (NEW v0.4-triz2)** | **Tracing, metrics, health, alerting** | **2** | **0** | **7** |
| **COST** | **Cost management (NEW v0.4-triz2)** | **LLM cost caps, budget, optimization** | **2** | **0** | **8** |

**14 тем** (was 11 in v0.3.1-triz). Three new themes added in v0.4.0-triz2: SEC, OBS, COST.

---

## 2. Theme A: Performance & Scalability

(unchanged from v0.3.1-triz — items A-001 through A-006 preserved)

### A-001: Lazy graph loading 🔥 80/20

**Theme:** A · **Status:** draft · **Phase:** BETA · **ARD requirement:** FR-PERF-003, FR-SRS-004

**Problem.** Граф `.caf/mental-model/graph.yml` на репо >10k файлов грузится в память целиком, что потребляет >100 MB и замедляет старт TUI dashboard.

**Sketch.** `git ls-files` для списка файлов; content-addressed cache (по hash от содержимого); ленивая загрузка подграфа при навигации. Аналогично `bubbles/list` lazy pagination.

**Open questions.** Хеш-функция? Хранить ли кэш в `.caf/cache/` (git-ignored) или `~/.caf/cache/`?

**Dependencies.** ADR-005 (YAML streaming), BACKLOG A-002 (streaming parser).
**agent_hint:** Затрагивает ADR-005 (формат) и инвариант "TUI first paint <100ms". Код-контракт: `internal/caf/mentalmodel/graph.go`.

---

### A-002: Streaming `.caf/` parser

**Theme:** A · **Status:** draft · **Phase:** BETA · **ARD requirement:** FR-PERF-005

**Problem.** `cards.yml` на 10k+ карточек парсится целиком в память, что потребляет >100 MB.

**Sketch.** `gopkg.in/yaml.v3` поддерживает streaming для arrays; `cards.idx.yml` — индекс с позиционными offset для random access. Парсер читает индекс при старте, карточки — по требованию.

**Open questions.** API парсера: channel-based или iterator? Бенчмарк на 100k карточек.

**Dependencies.** ADR-005 (YAML).
**agent_hint:** Затрагивает ADR-005 (формат); инвариант "file ops on 10k-file repo <1s". Код: `internal/caf/storage/yaml.go`.

---

### A-003 through A-006

(Preserved unchanged from v0.3.1-triz — see v0.3.1-triz BACKLOG.md for full text. Items: A-003 Incremental caf sync, A-004 Parallel adapter installation 🔥 80/20, A-005 Memory-mapped index 🔥 80/20, A-006 Distributed caf sync for monorepo fleets.)

---

## 3. Theme B: UX & DX

(Preserved from v0.3.1-triz — items B-001 through B-010. Items B-004, B-005, B-006, B-008, B-009 marked done in v0.3.1-triz remain done in v0.4.0-triz2.)

---

## 4. Theme C: Integrations

(Preserved from v0.3.1-triz — items C-001 through C-007.)

---

## 5. Theme D: AI & LLM Features

(Preserved from v0.3.1-triz — items D-001 through D-009.)

---

## 6. Theme E: Languages (beyond Go)

(Preserved from v0.3.1-triz — items E-001 through E-008. E-007 superseded by ADR-013.)

---

## 7. Theme F: Enterprise

(Preserved from v0.3.1-triz — items F-001 through F-007.)

---

## 8. Theme G: Multi-user & Collaboration

(Preserved from v0.3.1-triz — items G-001 through G-006.)

---

## 9. Theme H: Extensibility

(Preserved from v0.3.1-triz — items H-001 through H-006.)

---

## 10. Theme I: Analytics & Reporting

(Preserved from v0.3.1-triz — items I-001 through I-007. I-007 marked done in v0.3.1-triz remains done.)

---

## 11. Theme J: Quality & Reliability

(Preserved from v0.3.1-triz — items J-001 through J-008. J-003 marked done in v0.3.1-triz remains done.)

---

## 12. Theme K: UI Enhancements

(Preserved from v0.3.1-triz — items K-001 through K-012.)

---

## 13. Theme SEC: Security (NEW in v0.4-triz2)

> **TRIZ-pass changelog (v0.3.1-triz → v0.4.0-triz2):** This theme is entirely new. Security was implicit in v0.3.1-triz (scattered across ADR-001, ADR-005, BACKLOG F-002); v0.4.0-triz2 consolidates it as a first-class theme with 8 items. Each item maps to a FR-SEC-* requirement from ARD §4.8.

### SEC-001: Sandboxed YAML parser 🔥 80/20

**Theme:** SEC · **Status:** proposed · **Phase:** ALPHA · **ARD requirement:** FR-SEC-001

**Problem.** `.caf/` content from cloned untrusted repos may contain malicious YAML (e.g., custom tags that execute code, `!!python/object` attacks). The current parser (`gopkg.in/yaml.v3`) is safe by default but not formally sandboxed.

**Sketch.** Wrap `yaml.Unmarshal` in a sandbox layer: (1) reject any YAML tag starting with `!` (custom tags); (2) validate output struct (no `interface{}` at top level); (3) static analysis in CI to verify no `os.Execute`, no file writes outside `.caf/`, no network calls from parser code path. The sandbox is enforced at compile time via Go's package visibility and at runtime via the parser wrapper.

**Open questions.** Performance overhead? (Expected: <10ms per YAML file, negligible.) How to handle legitimate use cases for custom tags? (Answer: none — `.caf/` schema is fixed; custom tags not allowed.)

**Dependencies.** ADR-005 (YAML format), INV-ARCH-010 (new invariant).
**agent_hint:** Затрагивает ADR-005; вводит INV-ARCH-010 (".caf/ parser is sandboxed"). Код: `internal/caf/security/sandbox.go` (new). TRIZ Principle 22 (Blessing in disguise): the sandbox also makes the parser faster (no tag evaluation) and more predictable.

---

### SEC-002: Adapter signature verification 🔥 80/20

**Theme:** SEC · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-SEC-002

**Problem.** Community-contributed adapters (BACKLOG H-002) could be compromised in supply-chain attacks. `caf adapters install <user>/<repo>` currently installs without verification.

**Sketch.** Adapters are signed with GPG or cosign. `caf adapters install` verifies signature before install; rejects if signature missing or invalid. Signature registry: GitHub releases with `.sig` files alongside binaries. Maintainer key pinned in `caf` binary; community keys via trust-on-first-use.

**Open questions.** Key rotation? Maintainer key compromise recovery? Offline install of unsigned adapters (`--insecure` flag)?

**Dependencies.** ADR-022 (future — cross-repo sync security), ADR-028 (future — adapter contribution governance), BACKLOG H-002 (custom adapter registration).
**agent_hint:** Затрагивает ADR-022, ADR-028 (future). Код: `internal/caf/security/signature.go` (new). TRIZ Principle 11 (Beforehand cushioning): signature verification happens before install, not after.

---

### SEC-003: LLM data minimization 🔥 80/20

**Theme:** SEC · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-SEC-003

**Problem.** `caf.first_person_explain` sends data to LLM. If full file contents are sent, source code IP leaks to LLM provider. Even with `z-ai-web-dev-sdk` (default), enterprise users worry about data egress.

**Sketch.** LLM calls never include full file contents. Only: (1) decision records (already in `.caf/`, low risk); (2) diff summaries (≤500 chars, sanitized); (3) session metadata (timestamps, agent ID, prompt summary). Full file contents stay local. LLM paraphrases the facts, never sees the code.

**Open questions.** What counts as "diff summary"? (Answer: changed line count + function names touched, not actual code.) How to handle LLM prompts that "need" code context? (Answer: they don't — `caf.first_person_explain` is about decisions, not code review.)

**Dependencies.** ADR-006 (first-person voice), INV-ADR-006 (LLM only paraphrases), ADR-017 (future — LLM backend).
**agent_hint:** Затрагивает ADR-006, INV-ADR-006. Код: `internal/caf/mcp/first_person.go`. TRIZ Principle 22 (Blessing in disguise): data minimization also reduces hallucination (less context = less temptation to invent).

---

### SEC-004: Audit log with hash chain

**Theme:** SEC · **Status:** draft · **Phase:** FUTURE · **ARD requirement:** FR-SEC-005

**Problem.** Enterprise compliance (SOC2, ISO 27001) requires tamper-evident audit log. `.caf/audit.log` currently doesn't exist.

**Sketch.** `.caf/audit.log` is append-only. Each entry: `{timestamp, operation, user, hash_of_previous_entry, hash_of_this_entry}`. Hash chain (Merkle-style) makes tampering detectable. `caf audit --since 7d` queries; `caf audit --verify` checks chain integrity.

**Open questions.** Hash algorithm (SHA-256)? Rotation policy? What operations are audited (all `caf` commands, or only state-changing ones)?

**Dependencies.** ADR-001 (storage — `.caf/` size budget), BACKLOG F-002 (audit logs — this item subsumes F-002 for v0.4-triz2).
**agent_hint:** Затрагивает ADR-001, INV-ADR-001 (.caf/ size budget — audit log может расти). Mitigation: log rotation. Код: `internal/caf/security/audit.go` (new). TRIZ Principle 11 (Beforehand cushioning): hash chain is built as log is written, not retroactively.

---

### SEC-005: Telemetry opt-in with content redaction

**Theme:** SEC · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-SEC-006

**Problem.** Telemetry is opt-in (BACKLOG J-004), but the boundary of "what is sent" is unclear. Users may fear that prompts, paths, or decisions leak.

**Sketch.** Telemetry sends ONLY: (1) command name + duration; (2) exit code; (3) error class (not message); (4) aggregate counts (sessions, decisions). NEVER: prompts, file paths, decision content, user identifiers. All telemetry is reviewed by `caf config set telemetry.enabled true` (default false). Public dashboard: `caf telemetry stats`.

**Open questions.** How to verify no content leaks? (Answer: open-source the telemetry code; allow users to inspect network traffic via `caf config set telemetry.debug true`.)

**Dependencies.** ADR-019 (future), ADR-027 (future), BACKLOG J-004.
**agent_hint:** Затрагивает ADR-019, ADR-027 (future). Код: `internal/caf/observability/telemetry.go` (new). TRIZ Principle 22 (Blessing in disguise): strict content redaction also simplifies the telemetry pipeline (no PII to worry about).

---

### SEC-006: RBAC for `.caf/` operations

**Theme:** SEC · **Status:** idea · **Phase:** FUTURE · **ARD requirement:** (out of current ARD scope; maps to multi-user FUTURE)

**Problem.** In multi-user mode (FUTURE), not all users should be able to modify `.caf/decisions/` or `.caf/config.yml`.

**Sketch.** Roles: `reader` (read-only), `author` (can create sessions/decisions), `maintainer` (can modify config, install adapters). `.caf/config.yml` rbac section maps user IDs to roles.

**Open questions.** Role inheritance? Per-repo or global? Integration with SSO (BACKLOG F-001)?

**Dependencies.** BACKLOG F-001 (SSO), BACKLOG G-001 (multi-user mode).
**agent_hint:** Затрагивает ADR-001 storage schema. Код: `internal/caf/security/rbac.go` (future).

---

### SEC-007: Secret scanning in `.caf/`

**Theme:** SEC · **Status:** idea · **Phase:** PROD · **ARD requirement:** FR-SEC-001

**Problem.** Users may accidentally commit secrets (API keys, tokens) to `.caf/decisions/` or `.caf/sessions/` (e.g., a prompt that contains a hardcoded token).

**Sketch.** Pre-commit hook scans `.caf/` for common secret patterns (AWS keys, GitHub tokens, generic high-entropy strings). If found, blocks commit with "Potential secret detected in `.caf/<file>`. Remove or mark `CAF-Safe: true` if false positive." Uses `gitleaks` or similar library.

**Open questions.** False positive rate? Performance impact on pre-commit (expected <100ms)?

**Dependencies.** ADR-003 (hooks), BACKLOG J-005 (fuzzing — related but different).
**agent_hint:** Затрагивает ADR-003; extends INV-ARCH-004 (fail-safe — secret scan blocks but with clear bypass). Код: `internal/caf/security/secret_scan.go` (new).

---

### SEC-008: Supply-chain SBOM generation

**Theme:** SEC · **Status:** idea · **Phase:** PROD · **ARD requirement:** (out of current ARD scope; enterprise compliance)

**Problem.** Enterprise customers require SBOM (Software Bill of Materials) for compliance. `caf` binary should ship with SBOM.

**Sketch.** `caf sbom --format cyclonedx` generates SBOM from Go module dependencies. SBOM embedded in release artifacts (`.sig` + `.sbom.json` alongside binary).

**Open questions.** SBOM format (CycloneDX vs SPDX)? Automated generation in CI?

**Dependencies.** (none — standalone).
**agent_hint:** Код: `internal/caf/security/sbom.go` (new). TRIZ Principle 25 (Self-service): SBOM generated automatically in CI, not manually.

---

## 14. Theme OBS: Observability (NEW in v0.4-triz2)

> **TRIZ-pass changelog:** This theme is entirely new. Observability was implicit in v0.3.1-triz (scattered across BACKLOG I-005 Prometheus, J-003 adapter health); v0.4.0-triz2 consolidates it with 7 items, each mapping to FR-OBS-* from ARD §4.9.

### OBS-001: OpenTelemetry tracing integration 🔥 80/20

**Theme:** OBS · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-OBS-001

**Problem.** No way to trace `caf` command execution across components. When `caf sync` takes 30s, we don't know if it's LSP startup, graph build, or cross-repo fetch.

**Sketch.** Integrate OpenTelemetry Go SDK. Every `caf` command emits a root span; sub-operations (LSP call, YAML parse, graph build) emit child spans. Export via OTLP to any OTel collector (Jaeger, Tempo, Honeycomb). Opt-in via `caf config set observability.otel_enabled true`.

**Open questions.** Sampling rate (default 100% for CLI, 10% for daemon)? Performance overhead (expected <50ms per command, INV-ARCH-011)?

**Dependencies.** (new ADR needed for OTel integration), INV-ARCH-011 (new).
**agent_hint:** Затрагивает INV-ARCH-011 (new). Код: `internal/caf/observability/otel.go` (new). TRIZ Principle 25 (Self-service): tracing is automatic; no manual instrumentation per command needed (wrapper handles it).

---

### OBS-002: Metrics rolling log

**Theme:** OBS · **Status:** draft · **Phase:** PROD · **ARD requirement:** FR-OBS-002

**Problem.** Success metrics (ACS, Decision Coverage, Recall Rate) are computed on-demand but not stored historically. Cannot answer "what was our ACS 30 days ago?"

**Sketch.** `.caf/metrics/<date>.jsonl` — daily snapshot of all success metrics. Each line: `{date, metric, value, repo, agent}`. Rolling window: 365 days retained, older archived. `caf report --trend --metric acs --since 90d` reads from this log.

**Open questions.** Compression for old logs? Aggregation (hourly vs daily)?

**Dependencies.** ADR-001 (storage — `.caf/` size budget), BACKLOG I-001 (trend analysis).
**agent_hint:** Затрагивает ADR-001, INV-ADR-001 (.caf/ size — metrics log is small, ~10KB/day). Код: `internal/caf/metrics/rolling_log.go` (new).

---

### OBS-003: Prometheus metrics endpoint

**Theme:** OBS · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-OBS-003

**Problem.** Teams with existing Prometheus/Grafana stack cannot ingest `caf` metrics without custom integration.

**Sketch.** `caf metrics --export prometheus` starts HTTP server on port 9464, exposes `/metrics` endpoint with standard Prometheus format. Metrics: `caf_acs_score`, `caf_decision_coverage`, `caf_recall_rate`, `caf_llm_tokens_total`, `caf_adapter_health`. Grafana dashboard template shipped in `docs/grafana/caf-dashboard.json`.

**Open questions.** Port conflict (configurable via `observability.metrics_port`)? Authentication for remote Prometheus scrapers?

**Dependencies.** BACKLOG I-005 (Prometheus export — this item subsumes I-005 for v0.4-triz2), OBS-002 (rolling log — source of metrics).
**agent_hint:** Subsumes BACKLOG I-005. Код: `internal/caf/observability/prometheus.go` (new). TRIZ Principle 6 (Universality): standard Prometheus format works with any monitoring stack.

---

### OBS-004: caf doctor comprehensive health checks 🔥 80/20

**Theme:** OBS · **Status:** proposed · **Phase:** ALPHA · **ARD requirement:** FR-OBS-004

**Problem.** When something breaks (`.caf/` corrupt, adapter not firing, LSP down), user has no single command to diagnose.

**Sketch.** `caf doctor` runs all health checks:
1. `.caf/` folder structure valid
2. `config.yml` schema valid
3. `sessions/` index consistent with `git log`
4. `decisions/` no gaps (every non-trivial commit has decision)
5. `mental-model/graph.yml` fresh (regenerated in last 24h)
6. Universal pre-commit hook installed and executable
7. LSP server responds in <500ms
8. All installed adapters pass contract tests
9. `.caf/` size under 50 MB
10. (PROD) LLM endpoint reachable
11. (PROD) OTel collector reachable

Output: human-readable report with ✓/✗/⚠ per check. Exit code: 0 if all pass, 1 if any fail. CI integration: `caf doctor --ci` exits non-zero on any failure.

**Open questions.** Should `caf doctor` auto-fix issues (e.g., rebuild graph)? (Answer: no — suggest fixes, don't auto-apply; user controls.)

**Dependencies.** ADR-001, ADR-005, ADR-009, ADR-013, BACKLOG J-001 (self-heal — related but different; doctor diagnoses, self-heal fixes).
**agent_hint:** Затрагивает INV-ARCH-001, INV-ARCH-004, INV-ARCH-010. Код: `internal/caf/doctor/doctor.go`. TRIZ Principle 25 (Self-service): doctor runs on `caf init` and surfaces issues before they block work.

---

### OBS-005: Grafana dashboard template

**Theme:** OBS · **Status:** idea · **Phase:** FUTURE · **ARD requirement:** FR-OBS-005

**Problem.** Even with Prometheus endpoint, users must build their own Grafana dashboard.

**Sketch.** `docs/grafana/caf-dashboard.json` — pre-built dashboard with panels: ACS trend, Decision Coverage, Recall Rate, Agent Diversity, LLM token consumption, Adapter health. Importable via "Import dashboard" in Grafana.

**Open questions.** Dashboard versioning (track with `caf` version)? Multi-repo aggregation?

**Dependencies.** OBS-003 (Prometheus endpoint).
**agent_hint:** Код: `docs/grafana/caf-dashboard.json` (new). TRIZ Principle 26 (Copying): mirror the dashboard patterns users already know from other Go tools.

---

### OBS-006: Anomaly detection alerting

**Theme:** OBS · **Status:** idea · **Phase:** FUTURE · **ARD requirement:** (out of current ARD scope; extends FR-OBS-002)

**Problem.** ACS drops from 3.8 to 2.5 over 3 weeks — nobody notices until quarterly review.

**Sketch.** `caf report --anomalies` uses simple statistical methods (3-sigma, or Prophet for trends) on `.caf/metrics/` rolling log. Alerts when metric deviates >2σ from 30-day rolling mean. Integration: Slack webhook (BACKLOG C-003), email (BACKLOG I-006).

**Open questions.** Threshold tuning? False positive rate? ML model (simple regression vs Prophet)?

**Dependencies.** OBS-002 (rolling log), BACKLOG D-009 (anomaly detection — this item subsumes D-009 for v0.4-triz2), BACKLOG I-001 (trend analysis).
**agent_hint:** Subsumes BACKLOG D-009. Код: `internal/caf/observability/anomaly.go` (future).

---

### OBS-007: Distributed tracing across cross-repo sync

**Theme:** OBS · **Status:** idea · **Phase:** PROD · **ARD requirement:** FR-OBS-001 (extended)

**Problem.** `caf sync` across 10 repos — which repo was slow? Current tracing is per-process, not distributed.

**Sketch.** When `caf sync` fetches manifests from dependent repos, propagate trace context (W3C TraceContext headers if HTTP, or `CAF-Trace-Id` git trailer if git). Each repo's `caf sync` contributes a span to the same trace. Visualize in Jaeger/Tempo.

**Open questions.** Trace context propagation via git (unusual)? Privacy (trace IDs may leak repo relationships)?

**Dependencies.** OBS-001 (OTel integration), ADR-022 (future — cross-repo sync security).
**agent_hint:** Код: `internal/caf/observability/distributed.go` (future). TRIZ Principle 5 (Combining): trace context combines with git trailer (one mechanism, two purposes).

---

## 15. Theme COST: Cost Management (NEW in v0.4-triz2)

> **TRIZ-pass changelog:** This theme is entirely new. LLM cost was unaddressed in v0.3.1-triz (assumed free); v0.4.0-triz2 adds 8 items, recognizing that LLM calls cost money and need controls.

### COST-001: LLM token tracking 🔥 80/20

**Theme:** COST · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-AI-003 (extended)

**Problem.** No visibility into LLM token consumption. A heavy user could burn $100s/day on `caf.first_person_explain` without knowing.

**Sketch.** Every LLM call records: operation, tokens_in, tokens_out, cost (computed from model pricing table), timestamp. Stored in `.caf/metrics/llm-consumption.jsonl`. `caf cost --day` / `--week` / `--month` aggregates. Dashboard (TUI) shows breakdown by operation, cache hit rate, projected monthly cost.

**Open questions.** Pricing table source (hardcoded vs. configurable)? Cache invalidation when model pricing changes?

**Dependencies.** ADR-017 (future — LLM backend), OBS-002 (rolling log — same storage pattern).
**agent_hint:** Затрагивает ADR-017. Код: `internal/caf/cost/tracker.go` (new). TRIZ Principle 35 (Parameter change): token count is the parameter that adapts user behavior (visibility → conservation).

---

### COST-002: Per-user cost caps 🔥 80/20

**Theme:** COST · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-AI-003 (extended)

**Problem.** Without caps, a runaway script (e.g., `caf first_person_explain` in a loop) could exhaust a monthly LLM budget in hours.

**Sketch.** Config: `cost.llm_token_cap_daily` (default 100K), `cost.llm_token_cap_monthly` (default 2M), `cost.alert_threshold_pct` (default 80). When daily cap hit, LLM operations return error "Daily token cap exceeded; try tomorrow or `caf config set cost.llm_token_cap_daily <N>`". When alert threshold hit, TUI dashboard shows warning banner.

**Open questions.** Cap reset time (midnight local time vs. UTC)? Grace period? Per-repo vs. global caps?

**Dependencies.** COST-001 (token tracking), ADR-017 (future).
**agent_hint:** Затрагивает ADR-017. Код: `internal/caf/cost/caps.go` (new). TRIZ Principle 11 (Beforehand cushioning): caps prevent runaway costs before they happen, not after.

---

### COST-003: Narrative caching

**Theme:** COST · **Status:** draft · **Phase:** PROD · **ARD requirement:** FR-AI-001 (extended)

**Problem.** `caf.first_person_explain` on the same file (unchanged) returns identical narrative. Currently re-calls LLM each time.

**Sketch.** `.caf/cache/narratives/<file-hash>.md` — cache keyed by file content hash + decision IDs. Cache hit → return cached narrative (<10ms vs. <3s LLM call). Cache invalidated when: (a) file content changes; (b) linked decision records change; (c) `caf review --regenerate` invoked.

**Open questions.** Cache size limit? LRU eviction? Should cache be committed to `.caf/` (git) or `.caf/cache/` (git-ignored)?

**Dependencies.** ADR-006 (first-person voice), COST-001 (token tracking — measures cache effectiveness).
**agent_hint:** Затрагивает ADR-006. Код: `internal/caf/cost/narrative_cache.go` (new). TRIZ Principle 35 (Parameter change): cache hit rate is the parameter that reduces cost; 80% hit rate → 80% cost reduction.

---

### COST-004: Model tier selection

**Theme:** COST · **Status:** draft · **Phase:** PROD · **ARD requirement:** FR-AI-003

**Problem.** `caf.first_person_explain` uses one model (e.g., `glm-4.6`). For simple paraphrase, a smaller model (e.g., `glm-4-flash`) suffices and costs 10x less.

**Sketch.** Auto-select model based on task complexity: (1) simple paraphrase (decision has <500 chars context) → `glm-4-flash`; (2) multi-session synthesis (D-008) → `glm-4.6`; (3) decision summarization → `glm-4-air`. Configurable via `llm.model_tiers` in config.

**Open questions.** How to measure "task complexity" before calling LLM? (Answer: heuristic based on input length + decision count.) Fallback if smaller model produces lower quality?

**Dependencies.** ADR-017 (future — LLM backend), COST-001 (tracking — measures savings).
**agent_hint:** Затрагивает ADR-017. Код: `internal/caf/cost/model_selector.go` (new). TRIZ Principle 35 (Parameter change): model size is the parameter that scales cost vs. quality.

---

### COST-005: Batch decision summarization

**Theme:** COST · **Status:** idea · **Phase:** PROD · **ARD requirement:** FR-AI-004 (extended)

**Problem.** `caf decisions --summary` calls LLM once per decision. For 12 decisions, that's 12 LLM calls.

**Sketch.** Batch: one LLM call processes 5 decisions, returns 5 summaries. Reduces call count by 5x. Cache individual summaries after batch.

**Open questions.** Batch size (5? 10?)? Quality degradation in batch mode?

**Dependencies.** BACKLOG D-003 (decision auto-summarization — this item extends D-003), COST-001 (tracking).
**agent_hint:** Extends BACKLOG D-003. Код: `internal/caf/cost/batch.go` (new). TRIZ Principle 5 (Combining): one call combines multiple decisions.

---

### COST-006: Local Ollama for cost-zero mode

**Theme:** COST · **Status:** proposed · **Phase:** PROD · **ARD requirement:** FR-SEC-004

**Problem.** Cloud LLM costs add up. For users with a local GPU, Ollama provides cost-zero inference.

**Sketch.** `caf config set llm.backend local-ollama` + `llm.endpoint http://localhost:11434`. `caf` calls Ollama API (OpenAI-compatible). No per-token cost after GPU capex. Trade-off: slower (local GPU vs. cloud TPU), lower quality (7B vs. 175B).

**Open questions.** Model recommendation (`llama3.3-70b` for quality, `qwen2.5-coder-7b` for speed)? VRAM requirements?

**Dependencies.** ADR-017 (future — LLM backend), BACKLOG F-004 (on-prem LLM — this item overlaps; COST-006 is the user-facing version, F-004 is the enterprise version).
**agent_hint:** Overlaps BACKLOG F-004. Код: `internal/caf/llm/ollama.go` (new). TRIZ Principle 28 (Mechanics substitution): local inference substitutes cloud inference.

---

### COST-007: Cost projection and budgeting

**Theme:** COST · **Status:** idea · **Phase:** PROD · **ARD requirement:** (out of current ARD scope; extends COST-001)

**Problem.** User wants to know "if I continue at this rate, what will my monthly LLM cost be?"

**Sketch.** `caf cost --project` extrapolates from last 7 days of token consumption. Shows: projected monthly cost, days until monthly cap hit, recommended cap adjustment. Integrates with `caf config set cost.llm_token_cap_monthly` for "what-if" scenarios.

**Open questions.** Extrapolation method (linear vs. exponential smoothing)? Confidence interval?

**Dependencies.** COST-001 (tracking), COST-002 (caps).
**agent_hint:** Код: `internal/caf/cost/projection.go` (future). TRIZ Principle 11 (Beforehand cushioning): projection alerts before cap hit, not after.

---

### COST-008: Free tier for OSS contributors

**Theme:** COST · **Status:** idea · **Phase:** FUTURE · **ARD requirement:** (out of current ARD scope; business model)

**Problem.** OSS contributors to `caf` itself may not have LLM API budget. Want to offer free tier for approved contributors.

**Sketch.** `caf auth --oss-token <token>` enables free LLM calls (z-ai-web-dev-sdk sponsored) for contributors with ≥1 merged PR to `codeauthorityfirst` repo. Token issued via GitHub Action on PR merge.

**Open questions.** Abuse prevention? Monthly cap for free tier (10x lower than paid)? Sponsorship sustainability?

**Dependencies.** ADR-017 (future), BACKLOG F-007 (enterprise support tier — business model).
**agent_hint:** Business decision, not technical. TRIZ Principle 22 (Blessing in disguise): sponsoring OSS contributors increases adoption, which attracts enterprise customers.

---

## 16. Rejected Items

(Preserved from v0.3.1-triz — items R-001 through R-007: blockchain provenance, real-time collaborative editing, automatic decision generation, mobile app, gamification, public leaderboard, automatic decision generation from git log.)

### R-008 (NEW v0.4-triz2): Cryptocurrency-based LLM cost settlement

**Theme:** Rejected · **Status:** rejected

**Rejection reason.** Paying for LLM calls via cryptocurrency adds complexity (wallet management, price volatility) for no user value. Standard payment via LLM provider (z-ai-web-dev-sdk, OpenAI, Anthropic) suffices. TRIZ Principle 1 (Segmentation): don't add layers when one suffices.

---

## 17. Statistics (v0.4.0-triz2 recomputed)

| Theme | Items | 🔥 80/20 | proposed | draft | idea | done | rejected |
|-------|-------|----------|----------|-------|------|------|----------|
| A | 6 | 2 | 2 | 2 | 2 | 0 | 0 |
| B | 10 | 5 | 1 | 0 | 1 | 5 (B-004/5/6/8/9) | 0 |
| C | 7 | 2 | 2 | 1 | 4 | 0 | 0 |
| D | 9 | 1 | 1 | 0 | 8 | 0 | 0 |
| E | 8 | 1 | 0 | 2 | 5 | 0 | 1 (E-007 superseded) |
| F | 7 | 0 | 2 | 0 | 5 | 0 | 0 |
| G | 6 | 0 | 0 | 1 | 5 | 0 | 0 |
| H | 6 | 1 | 1 | 1 | 3 | 0 | 0 |
| I | 7 | 2 | 2 | 1 | 2 | 1 (I-007) | 0 |
| J | 8 | 2 | 2 | 2 | 2 | 1 (J-003) | 0 |
| K | 12 | 3 | 2 | 2 | 5 | 0 | 0 |
| **SEC (NEW)** | **8** | **3** | **4** | **1** | **3** | **0** | **0** |
| **OBS (NEW)** | **7** | **2** | **3** | **1** | **3** | **0** | **0** |
| **COST (NEW)** | **8** | **2** | **3** | **2** | **3** | **0** | **0** |
| **Subtotal active** | **109** | **26** | **25** | **15** | **48** | **7** | — |
| Rejected | 8 | — | — | — | — | — | 8 |
| **Total** | **117** | **26 (22%)** | — | — | — | — | — |

> **TRIZ-pass reconciliation (v0.4.0-triz2):**
> - 3 new themes added (SEC, OBS, COST) with 23 new items (8+7+8).
> - 7 🔥 80/20 items added across new themes (SEC-001/002/003, OBS-001/004, COST-001/002).
> - Total active items: 109 (was 86 in v0.3.1-triz; +23 new).
> - 🔥 80/20 percentage: 26/109 = 24% (within target 20-25%).
> - All new items map to ARD requirements (FR-SEC-*, FR-OBS-*, FR-AI-* extended for COST).
> - 1 new rejected item (R-008: cryptocurrency-based LLM cost settlement).

---

## 18. How to Promote a Backlog Item to Roadmap

Promotion criteria — see [ROADMAP.md §6.9](./ROADMAP.md#69-exit-criteria-for-future-items) (4 promotion criteria: clear problem; design sketch; pilot user; success metric).

Process:

1. Fill the §0 template (problem, sketch, open questions, dependencies, agent_hint, ARD requirement).
2. Find a pilot user (someone who will actually use the feature within 2 weeks of shipping).
3. Define a success metric — concrete number, not "feels better". E.g., "Decision Coverage ≥90%" not "more decisions get recorded".
4. Create an ADR if the change is architectural.
5. Map to ARD requirement (FR-*) — if no FR exists, add one to ARD first.
6. Open a GitHub Discussion with label `backlog-promotion`.
7. Maintainer review within 2 weeks.

---

## 19. How to Add a New Backlog Item

1. Open a GitHub Issue with label `backlog`.
2. Use the §0 template.
3. Specify the theme (A-K, SEC, OBS, COST, or new theme).
4. Specify the ARD requirement (FR-*) it implements — NEW in v0.4-triz2.
5. Optionally propose a phase (MVP/ALPHA/BETA/PROD/FUTURE).
6. Quarterly review: promote / reject / merge duplicates / clarify open questions.

---

## 20. Trace Matrix (TRIZ-2-pass addition)

| Goal | ADR | Backlog Item | ARD Requirement | Code Contract | UI Section |
|------|-----|-------------|------------------|---------------|------------|
| Provenance by default | ADR-001, ADR-003 | (none — implemented) | FR-PROV-001..006 | `internal/caf/storage/`, `internal/caf/hooks/` | UI.md §3.5 refresh |
| SRS-based mental model | ADR-002 | D-007 (adaptive difficulty) | FR-SRS-001..007 | `internal/caf/srs/sm2.go` | UI.md §4 Review Session |
| Triviality threshold | ADR-004 | H-006 (Lua rules) | FR-DEC-001 | `internal/caf/triviality/` | — |
| Adapter pattern | ADR-009, ADR-012 | H-002 (custom adapters) | FR-ADP-001..007 | `internal/caf/adapter.go`, `internal/adapters/` | UI.md §7 Adapters Manager |
| Universal LSP fallback | ADR-013 | E-007 (superseded) | FR-ADP-004 | `internal/caf/lsp/fallback.go` | UI.md §7.4 Health check |
| TUI-first UI | ADR-014 | K-001..K-012 | FR-UI-001..010 | `pkg/ui/` | UI.md §0-§15 |
| 80/20 prioritization | ADR-015 | (8 high-leverage features in ARD §0) | (all FRs tagged) | (cross-cutting) | UI.md §0.3 Performance budget |
| **Security (NEW)** | ADR-005 (sandbox), ADR-022/028 (future) | SEC-001..SEC-008 | FR-SEC-001..006 | `internal/caf/security/` | (internal) |
| **Observability (NEW)** | (new ADR for OTel) | OBS-001..OBS-007 | FR-OBS-001..005 | `internal/caf/observability/` | UI.md §10.3 Metrics |
| **Cost management (NEW)** | ADR-017 (future LLM backend) | COST-001..COST-008 | FR-AI-001..006 (extended) | `internal/caf/cost/` | UI.md §10.4 (proposed) |
| Invariant enforcement | ADR-029 (proposed) | (this file — agent_hint fields) | FR-DEC-007, FR-REL-001 | `internal/caf/invariants/` (future) | — |
