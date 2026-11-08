# GITHUB.md — стратегия релизов, исправление CI и рост звёзд

> Файл подготовлен для репозитория **[Reider85/OryaCodeAuthorityFirst](https://github.com/Reider85/OryaCodeAuthorityFirst)** (CodeAuthorityFirst / `caf`, Go 1.27, single binary, статус MVP `0.4.1-mvp.1`).
> Главная цель: проект должен **чаще мелькать** в ленте GitHub и **набирать звёзды** на пушах, при этом максимум работы должен выполняться **на ресурсах GitHub**, а не на ваших.

---

## 1. Краткое резюме (TL;DR)

| Вопрос | Ответ |
|---|---|
| Почему CI всегда failed? | В `ci.yml` и `performance-enhanced-ci.yml` сломан триггер: `branches: ain]` вместо `branches: [ main ]`. GitHub Actions не может распарсить YAML и сразу валит запуск. Это **единственная глобальная причина** — она блокирует вообще все push-события. |
| Дополнительные причины | Устаревшие `actions/upload-artifact@v3`, `setup-go@v4` с захардкоженным Go 1.21 (проект на Go 1.27), обращения к несуществующим таргетам `./scripts/analyze-performance-trend`, `./scripts/generate-dashboard`, флагам `caf --performance-test`, `caf traceability check`, которые скорее всего не реализованы. |
| Что уже реализовано? | 4 workflow: `ci.yml`, `performance-enhanced-ci.yml`, `traceability-validation.yml`, `manual-traceability.yml`. В `ci.yml` есть `vet`, `test` + покрытие ≥80%, `benchmark`, `build` со smoke-тестом. |
| Релизы на каждом пуше? | Сейчас **не создаются вообще** (1 тег, 1 release за всю историю). Ниже — промпт на auto-prerelease на каждый push в `main` + полноценный release на semver-тег через goreleaser. |
| Trending | Прямой способ «попасть в тренды» — это **звёзды за последние 7–30 дней** + стабильная видимая активность (коммиты, PR, релизы, issue). Спам-релизы на каждый коммит приведут к обратному эффекту — GitHub фильтрует такие репо. Правильная схема: prerelease на push + «настоящий» release на тег + Dependabot для постоянной активности. |
| Что даст оптимизация? | Зелёные галочки на коммитах, бейджи доверия, кросс-платформенные бинарники в релизах, авто-ченджлог, постоянный поток Dependabot-PR — всё это напрямую влияет на доверие и звёзды. |

---

## 2. Что я нашёл в репозитории (исходное состояние)

**Стек:** Go 1.27, без CGO, single binary `cmd/caf`, есть директория `cmd/perf-test`.
**Активность:** 38 запусков workflow, **все failed**. Релизов: 1, тегов: 1, звёзд: 1, topics: 0 (пусто).
**Workflow-файлы (4 шт.):**

| Файл | Триггер | Состояние |
|---|---|---|
| `ci.yml` | `push.branches: ain]` (сломано), `tags: ["v*"]`, `pull_request` | **парсинг YAML не проходит** → все запуски failed |
| `performance-enhanced-ci.yml` | та же ошибка `branches: ain]` | failed + ссылается на несуществующие `scripts/` и флаги `caf` |
| `traceability-validation.yml` | `push.branches: [ main, develop ]` | парсится, но `setup-go@v4` + `go-version: '1.21'` (нужно 1.27), `upload-artifact@v3` (retired), зависит от подкоманды `caf traceability check` |
| `manual-traceability.yml` | `workflow_dispatch` | `setup-go@v4` + `go-version: '1.21'`, `context.issue.number` в не-PR контексте → выброс |

### 2.1 Точная причина «всегда failed»

Вот как сейчас выглядит триггер в `ci.yml` и `performance-enhanced-ci.yml`:

```yaml
on:
  push:
    branches: ain]      # ← сломано: должно быть [ main ]
    tags: ["v*"]
  pull_request:
```

`branches: ain]` — это невалидный YAML-скаляр (что-то съело `[ m` при прошлой правке). GitHub Actions падает ещё до запуска джобов, поэтому **все 38 запусков красные**. Это первая и обязательная вещь к починке.

### 2.2 Прочие технические долги

- `actions/upload-artifact@v3` и `actions/download-artifact@v3` — **retired**, в v4 изменился API (артефакты не шарятся между джобами автоматически, нужен явный upload/download с одинаковым именем).
- `actions/setup-go@v4` с `go-version: '1.21'` в двух workflow — проект на `go 1.27`, надо `go-version-file: go.mod`.
- `performance-enhanced-ci.yml` ссылается на:
  - `./scripts/analyze-performance-trend` — пакет/бинарник с таким путём нужно создавать;
  - `./scripts/generate-dashboard` — то же самое;
  - `./caf --performance-test --threshold=0.1` — флаг, скорее всего не реализован в `cmd/caf`;
  - `./perf-test --comprehensive`, `--report`, `--check-regression`, `--baseline`, `--quality-gate`, `--max-regression`, `--min-budget-compliance` — нужно проверить, что `cmd/perf-test` реально парсит эти флаги.
- Slack webhook `secrets.SLACK_WEBHOOK_URL` не задан — шаг молча skipped, не критично, но шаг в `performance-alert` лучше сделать `if: failure() && env.SLACK_WEBHOOK_URL != ''`.
- `manual-traceability.yml` использует `context.issue.number` — в `workflow_dispatch` этого контекста нет, GitHub-script упадёт. Нужно оборачивать в `if: github.event_name == 'pull_request'`.

---

## 3. Как оптимизировать и что это даст

| Что сделать | Что даст | Приоритет |
|---|---|---|
| Починить `branches: [ main ]` в двух workflow | CI из красного → зелёный за 1 коммит. Зелёный бейдж = +доверие. | P0 |
| Обновить actions до v4/v5 | Совместимость с актуальным API GitHub, артефакты снова шарятся между джобами. | P0 |
| `setup-go` через `go-version-file: go.mod` | Перестанет падать при апгрейде Go. | P0 |
| Удалить или «обезвредить» шаги со ссылками на несуществующие `scripts/`/флаги | Перестанут падать джобы `performance-gate`, `performance-trend`, `performance-dashboard`. | P0 |
| **goreleaser + auto-release на semver-тег** | Профессиональные релизы: кросс-платформенные бинарники, checksums, SBOM, авто-ченджлог из conventional commits. Каждый релиз = событие в ленте подписчиков. | P1 |
| **Auto-prerelease на каждый push в main** | Лента « releases » обновляется на каждом пуше → кажется живым. Но **только prerelease**, чтобы не спамить подписчиков «настоящими» релизами. | P1 |
| **Build matrix (linux/mac/windows × amd64/arm64)** | В каждом релизе появляются 6+ бинарников → больше платформ → больше скачиваний → GitHub считает репо активным. | P1 |
| `golangci-lint` + `govulncheck` + `gosec` | Три бейджа доверия в README. Особенно `govulncheck` — его смотрят. | P1 |
| Codecov upload | Бейдж покрытия, который обновляется на каждом пуше. | P2 |
| Dependabot (gomod + github-actions) | Постоянный поток авто-PR → постоянная активность → репо «дышит». | P2 |
| Topics + description + social preview | Без topics репо **не найдут** через topic-страницы — это бесплатный канал трафика. | P1 |
| Issue/PR templates + Discussions | Снижает порог входа, первые issue/PR = дополнительные события в ленте. | P2 |
| GitHub Pages для docs/dashboard | Бесплатный хостинг документации, использует ресурсы GitHub. | P3 |

> Все шаги, кроме Codecov (есть free-tier для OSS), используют **бесплатные** GitHub-ресурсы: Actions (безлимит для public), Releases, Pages, Dependabot, Container Registry. Ваши ресурсы вообще не задействуются.

---

## 4. Какие тесты можно выполнять через GitHub CI и улучшит ли это

**Уже есть:** `go vet`, unit + contract + TUI teatest snapshots, coverage gate ≥80%, UI-бенчмарки (first_paint ≤100ms, frame ≤16ms, keystroke ≤50ms).

**Что стоит добавить** (каждое — это и качество, и бейдж доверия):

| Тест | Что проверяет | Улучшит ли звёзды? |
|---|---|---|
| **`golangci-lint`** | 30+ статических анализаторов вместо одного `go vet` | Косвенно — бейдж в README |
| **`govulncheck`** | Уязвимости в зависимостях (golang.org/x/vuln) | Да, особенно для devtool |
| **`gosec`** | Security-анализ Go-кода | Да, отдельный бейдж |
| **Cross-OS build matrix** | Что `caf` собирается и запускается на linux/macOS/windows | Да — больше平台 в релизах |
| **Integration tests с реальным git-репо** | Что хуки реально срабатывают на `git commit` в temp-репо | Да — доказательство, что MVP не «на бумаге» |
| **`go test -race`** | Data races в concurrency (хуки, watcher) | Косвенно |
| **Reproducible build check** | Два билда дают одинаковый хеш | Косвенно, trust-сигнал |
| **Releasing smoke на artifacts** | После `goreleaser` — скачать бинарник из release и запустить `caf --version` | Гарантирует, что релиз рабочий |
| **CodeQL** | GitHub-овский security-сканер кода, бесплатный для OSS | Да, отдельный бейдж security |

> **Вывод:** «просто добавить тесты» само по себе звёзд не приносит. Но **зелёный CI + бейджи + кросс-платформенные релизы** — это ровно те сигналы, по которым разработчик решает «поставить звезду или пройти мимо». Для devtool-проекта (а `caf` — devtool) это критично.

---

## 5. Стратегия для GitHub Trending и звёзд (честно, без gaming)

GitHub Trending учитывает: **звёзды за последние 7–30 дней**, форки, активность (коммиты/PR/релизы), качество README, topics. Прямого «попасть в тренды по релизам» нет — но есть косвенные рычаги:

1. **Каждый push в main = prerelease.** Это создаёт «release event» в ленте тех, кто уже поставил звезду или watch. Не спамить настоящими релизами.
2. **Каждый semver-tag = настоящий release** через goreleaser, с бинарниками для 6 платформ и авто-ченджлогом. Это большие события.
3. **Dependabot** — 2–8 авто-PR в неделю на обновление зависимостей. Постоянная «дышащая» активность в Network-графе.
4. **Topics** — без них репо не попадает в страницы `github.com/topics/go`, `github.com/topics/developer-tools`, `github.com/topics/cli`. Это бесплатный SEO.
5. **Социальный preview** (карточка 1280×640) — то, что видно в Twitter/Telegram при шеринге. Без неё шеринг даёт унылую серую карточку.
6. **README с бейджами и GIF/скриншотом TUI** — первый кадр, по которому решают «звездовать или нет».
7. **Issue/PR templates + Discussions** — снижают порог первого контакта; каждый новый issue = событие в ленте подписчиков.

⚠ **Анти-паттерн**, которого стоит избегать: создавать «настоящий» release на каждый коммит. GitHub помечает такие репо как спам в алгоритмах trending, а подписчики отписываются. Поэтому ниже — **prerelease на push + real release на тег**, а не «release на каждый коммит».

---

## 6. Готовые промпты

Ниже — 12 самостоятельных промптов. Каждый можно скопировать и вставить в Copilot Chat / Claude / ChatGPT/ Cursor. Они не зависят друг от друга, но **применять строго по порядку P0 → P1 → P2**, иначе P1 будет падать на сломанном YAML.

### Промпт 1 (P0) — починить сломанный триггер `branches: ain]`

```
В репозитории Reider85/OryaCodeAuthorityFirst в файлах
.github/workflows/ci.yml и .github/workflows/performance-enhanced-ci.yml
поломан триггер push: сейчас написано `branches: ain]`, а должно быть
`branches: [ main ]`. Из-за этого YAML не парсится и все 38 запусков CI
падают ещё до старта джобов.

Задача:
1. В обоих файлах замени строку `branches: ain]` на `branches: [ main ]`.
2. В performance-enhanced-ci.yml в блоке `pull_request:` тоже замени
   `branches: ain]` на `branches: [ main ]`.
3. Больше ничего не меняй. Сделай отдельный коммит с сообщением
   `fix(ci): repair broken push trigger (branches: ain] -> [ main ])`.
4. В PR-описании укажи: «Восстанавливает парсинг workflow, после мержа
   CI перестанет быть красным для всех push-событий».

Критерий приёмки: после мержа следующий push в main запускает CI без
ошибок парсинга YAML, джобы переходят в статус in_progress.
```

### Промпт 2 (P0) — обновить устаревшие actions и Go-версию

```
В .github/workflows у нас 4 файла. Приведи их к актуальному состоянию:

1. Во всех файлах замени:
   - actions/upload-artifact@v3 -> actions/upload-artifact@v4
   - actions/download-artifact@v3 -> actions/download-artifact@v4
   - actions/setup-go@v4 -> actions/setup-go@v5
   - actions/checkout@v4 оставить как есть
   - actions/github-script@v6 -> actions/github-script@v7

2. В traceability-validation.yml и manual-traceability.yml убери
   жёстко заданный `go-version: '1.21'` и замени на:
     with:
       go-version-file: go.mod
   (проект сейчас на go 1.27, и должен следовать за go.mod автоматически)

3. В performance-enhanced-ci.yml шаг `performance-alert` -> `Create issue`
   оберни в условие if: failure(), чтобы issue создавался только при
   реальном падении. Slack-шаг оберни в
   `if: failure() && env.SLACK_WEBHOOK_URL != ''`, чтобы он не падал,
   когда секрет не задан.

4. В manual-traceability.yml шаг `Comment on workflow` (github-script)
   оберни в `if: github.event_name == 'pull_request'`, потому что
   context.issue.number отсутствует в workflow_dispatch.

5. Все изменения — отдельным коммитом
   `chore(ci): bump actions to v4/v5, use go-version-file`.

Критерий приёмки: после мержа нет ни одного deprecated warning в
вкладке Actions, джобы setup-go/upload-artifact проходят успешно.
```

### Промпт 3 (P0) — убрать шаги, которые ссылаются на несуществующие скрипты

```
В .github/workflows/performance-enhanced-ci.yml есть несколько шагов,
которые ссылаются на артефакты кода, которых нет в репо:

- ./scripts/analyze-performance-trend
- ./scripts/generate-dashboard
- флаги caf: --performance-test, --threshold
- флаги perf-test: --comprehensive, --report, --output,
  --check-regression, --baseline, --quality-gate, --max-regression,
  --min-budget-compliance

Задача:
1. Создай в репо временный аудит: проверь факт существования
   ./scripts/analyze-performance-trend и ./scripts/generate-dashboard
   (как Go-пакетов main), а также парсинг флагов в cmd/caf/main.go
   и cmd/perf-test/main.go.
2. Если их нет — временно закомментируй соответствующие шаги
   `performance-trend`, `performance-dashboard`, `Performance
   validation` в build-джобе, и шаги в performance-gate, которые
   используют несуществующие флаги. Оставь только то, что реально
   работает: `go vet`, `go test`, benchmark, `build` + smoke
   `./caf --version`.
3. Над каждым закомментированным шагом добавь комментарий:
   # TODO: implement <script/flag> and re-enable (see GITHUB.md)
4. Открой отдельный issue «Implement missing perf-test flags and
   scripts/* tools referenced in performance-enhanced-ci.yml» с
   чек-листом того, что нужно реализовать, чтобы шаги можно было
   включить обратно.
5. Коммит: `chore(ci): disable steps referencing unimplemented targets`.

Критерий приёмки: после мержа performance-enhanced-ci.yml проходит
целиком до зелёного статуса. В issue есть явный список того, что
нужно реализовать, чтобы вернуть отключённые шаги.
```

### Промпт 4 (P1) — auto-prerelease на каждый push в main

```
Цель: на каждый push в main автоматически создавать GitHub prerelease
с тегом вида v<last-tag>-pre.<short-sha>, чтобы в ленте релизов
проект обновлялся постоянно, но не спамил подписчиков настоящими
релизами.

Создай новый файл .github/workflows/auto-prerelease.yml:

name: Auto Prerelease

on:
  push:
    branches: [ main ]

permissions:
  contents: write

jobs:
  prerelease:
    name: Create prerelease for push
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Compute version
        id: ver
        run: |
          LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
          SHA=$(git rev-parse --short HEAD)
          DATE=$(date -u +%Y%m%d)
          PRERELEASE_TAG="${LAST_TAG}-pre.${DATE}.${SHA}"
          echo "tag=$PRERELEASE_TAG" >> "$GITHUB_OUTPUT"
          echo "last_tag=$LAST_TAG" >> "$GITHUB_OUTPUT"

      - name: Build binary (linux amd64)
        run: |
          go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=${{ steps.ver.outputs.tag }}'" -o caf-linux-amd64 ./cmd/caf
          ./caf-linux-amd64 --version

      - name: Delete previous prerelease
        continue-on-error: true
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release delete "${{ steps.ver.outputs.last_tag }}-pre" --yes --cleanup-tag || true

      - name: Create prerelease
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release create "${{ steps.ver.outputs.tag }}" \
            --title "Prerelease ${{ steps.ver.outputs.tag }}" \
            --notes "Автоматический prerelease для коммита ${{ github.sha }}.
          Список изменений со времени последнего тега:
          $(git log --oneline ${{ steps.ver.outputs.last_tag }}..HEAD)" \
            --prerelease \
            caf-linux-amd64

Требования:
- Используй только встроенный GITHUB_TOKEN, без внешних секретов.
- Prerelease создаётся с флагом --prerelease, чтобы он не считался
  настоящим релизом и не spam-ил подписчиков.
- Перед созданием нового prerelease удаляй предыдущий prerelease того
  же последнего «настоящего» тега, чтобы лента не плодила десятки
  prereleases. Для этого в шаге Compute version также экспортируй
  last_tag.
- Коммит: `feat(ci): auto prerelease on every push to main`.

Критерий приёмки: после push в main в Releases появляется свежий
prerelease с одним бинарником caf-linux-amd64, предыдущий prerelease
того же тега удалён, настоящие релизы (с semver-тегом) не затронуты.
```

### Промпт 5 (P1) — настоящий release на semver-тег через goreleaser

```
Цель: когда я ставлю semver-тег вида v0.5.0, автоматически создавать
полноценный релиз с бинарниками под 6 платформ, checksums, SBOM и
авто-ченджлогом. Это «большое событие» в ленте, в отличие от prerelease.

Шаги:

1. Добавь в корень .goreleaser.yml:

project_name: caf
before:
  hooks:
    - go mod tidy
builds:
  - id: caf
    main: ./cmd/caf
    binary: caf
    env:
      - CGO_ENABLED=0
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version={{.Version}}
      - -X github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.commit={{.Commit}}
      - -X github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.date={{.Date}}
archives:
  - id: default
    builds: [caf]
    name_template: >-
      {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
    format_overrides:
      - goos: windows
        format: zip
checksum:
  name_template: 'checksums.txt'
sbom:
  artifacts: archive
snapshot:
  name_template: "{{ incpatch .Version }}-next"
changelog:
  use: github
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^chore:'
  groups:
    - title: 'New features'
      regexp: '^feat:'
      order: 0
    - title: 'Bug fixes'
      regexp: '^fix:'
      order: 1
    - title: 'Other'
      order: 999
release:
  draft: false
  prerelease: auto
  name_template: "{{ .Tag }} — CodeAuthorityFirst"

2. Создай .github/workflows/release.yml:

name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    name: goreleaser
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

3. Коммит: `feat(release): goreleaser on semver tags, cross-OS binaries + SBOM`.

Критерий приёмки: после `git tag v0.5.0 && git push origin v0.5.0`
в Releases появляется полноценный релиз с архивами
caf_<ver>_linux_amd64, _linux_arm64, _darwin_amd64, _darwin_arm64,
_windows_amd64.zip, checksums.txt, SBOM-файлом и сгруппированным
changelog из conventional commits.
```

### Промпт 6 (P1) — заполнить topics, description и social preview

```
Цель: репо сейчас имеет 0 topics и пустое description, поэтому его
невозможно найти через github.com/topics/go, /topics/cli, /topics/developer-tools.
Это бесплатный канал трафика и звёзд.

Задача — выпиши мне готовый набор значений, я вставлю их руками в
Settings → General / Tags:

About → Description (≤350 символов):
"CodeAuthorityFirst (caf) — single-binary Go tool that records
provenance for AI coding sessions: which agent touched which file,
from which prompt, through which iteration. Offline, no LLM calls,
fail-safe pre-commit hooks."

About → Website: https://github.com/Reider85/OryaCodeAuthorityFirst#readme

About → Topics (ровно 20 шт., максимум что позволяет GitHub):
go, golang, cli, devtool, developer-tools, ai-coding, ai-agents,
provenance, traceability, audit, developer-experience, dx,
single-binary, pre-commit, tui, source-code, code-review,
observability, llm-tools, agent-ops

Дополнительно: сгенерируй SVG-карточку 1280×640 для Social Preview
(на тёмном фоне, с названием CodeAuthorityFirst, одной строкой
слогана «who wrote what, when, and why in AI coding sessions», и
маленьким значком caf в углу). Назови файл social-preview.svg,
положи в docs/. Затем дай мне пошаговую инструкцию, как загрузить
её в Settings → Social preview, чтобы при шеринге в Twitter/Telegram
появлялась красивая карточка.

Критерий приёмки: на странице репозитория видно 20 topics, заполненный
description, и в Settings → Social preview загружено изображение.
```

### Промпт 7 (P1) — `golangci-lint` + `govulncheck` + `gosec` в CI

```
Цель: три бейджа доверия в README — lint, vulnerability, security.

1. Создай .github/workflows/quality.yml:

name: Quality

on:
  push:
    branches: [ main ]
    tags: ["v*"]
  pull_request:
    branches: [ main ]

jobs:
  golangci-lint:
    name: golangci-lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

  govulncheck:
    name: govulncheck
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - run: govulncheck ./...

  gosec:
    name: gosec
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: securego/gosec-action@master
        with:
          args: -exclude=G104 ./...

2. Создай в корне .golangci.yml с минимально разумным набором линтеров:
   enabled: govet, errcheck, ineffassign, staticcheck, unused, gosec,
   misspell, revive, gocritic, gocyclo (max-complexity 15).

3. В README.md в шапку (сразу под заголовком) добавь три бейджа:

   [![golangci-lint](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml/badge.svg?branch=main&label=golangci-lint)](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml)
   [![govulncheck](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml/badge.svg?branch=main&label=govulncheck)](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml)
   [![gosec](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml/badge.svg?branch=main&label=gosec)](https://github.com/Reider85/OryaCodeAuthorityFirst/actions/workflows/quality.yml)

4. Коммит: `feat(quality): add golangci-lint, govulncheck, gosec + badges`.

Критерий приёмки: три зелёных бейджа в README, все три джобы в
quality.yml проходят на main.
```

### Промпт 8 (P1) — build matrix для smoke-теста на 6 платформах

```
Цель: доказать, что caf собирается и запускается на 6 платформах.
Это и сигнал качества, и подготовка к кросс-платформенным релизам.

1. В .github/workflows/ci.yml добавь новую джобу build-matrix после
   существующей build:

  build-matrix:
    name: build (${{ matrix.os }}/${{ matrix.arch }})
    runs-on: ${{ matrix.runner }}
    needs: [vet, test]
    strategy:
      fail-fast: false
      matrix:
        include:
          - { os: linux,   arch: amd64, runner: ubuntu-latest }
          - { os: linux,   arch: arm64, runner: ubuntu-24.04-arm }
          - { os: darwin,  arch: amd64, runner: macos-13 }
          - { os: darwin,  arch: arm64, runner: macos-14 }
          - { os: windows, arch: amd64, runner: windows-latest }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Build
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
          CGO_ENABLED: 0
        run: |
          go build -ldflags="-X 'github.com/Reider85/OryaCodeAuthorityFirst/cmd/caf.version=${{ github.sha }}'" -o caf${{ matrix.os == 'windows' && '.exe' || '' }} ./cmd/caf
      - name: Smoke test
        run: |
          ./caf${{ matrix.os == 'windows' && '.exe' || '' }} --version

2. Коммит: `feat(ci): build matrix across linux/macOS/windows × amd64/arm64`.

Критерий приёмки: в списке checks каждого коммита видно 5 зелёных
джоб (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64,
windows-amd64), все проходят --version.
```

### Промпт 9 (P2) — Dependabot для постоянной активности

```
Цель: Dependabot будет 1–3 раза в неделю автоматически открывать
PR на обновление go-модулей и GitHub Actions. Каждый такой PR — это
событие в Network-графе и в ленте подписчиков, репо «дышит».

Создай .github/dependabot.yml:

version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
      day: monday
      time: "06:00"
      timezone: Europe/Moscow
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - go
    commit-message:
      prefix: "chore(deps)"
      include: scope

  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    labels:
      - dependencies
      - github-actions
    commit-message:
      prefix: "chore(ci-deps)"
      include: scope

  - package-ecosystem: docker
    directory: /
    schedule:
      interval: monthly
    open-pull-requests-limit: 2
    labels:
      - dependencies
      - docker

Коммит: `chore(deps): enable Dependabot for gomod + actions`.

Критерий приёмки: в течение 7 дней в Issues/PRs появляется хотя бы
один авто-PR с меткой dependencies. Настроено auto-merge для PR с
проходящим CI — отдельным шагом через GitHub Actions или через
настройку Branch protection.
```

### Промпт 10 (P2) — Codecov +覆盖率 бейдж

```
Цель: публичный бейдж покрытия, обновляется на каждом пуше.
Бесплатно для OSS.

1. В .github/workflows/ci.yml в джобе test после шага
   `Enforce NFR-MAINT-01 coverage floor` добавь:

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.out
          fail_ci_if_error: false
        env:
          CODECOV_TOKEN: ${{ secrets.CODECOV_TOKEN }}

2. Зарегистрируй репо на codecov.io через GitHub login, получи
   CODECOV_TOKEN, добавь его в Settings → Secrets and variables →
   Actions как CODECOV_TOKEN.

3. В README.md рядом с другими бейджами добавь:

   [![coverage](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst/graph/badge.svg?token=ВАРЬИРУЙСЯ-ТОКЕН)](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst)

   где ВАРЬИРУЙСЯ-ТОКEN замени на публичный бейдж-токен из Codecov.

4. Коммит: `feat(ci): upload coverage to Codecov + badge`.

Критерий приёмки: бейдж покрытия показывает процент, после пуша в
main он обновляется. На codecov.io виден тренд покрытия по файлам.
```

### Промпт 11 (P2) — Issue/PR templates + Discussions

```
Цель: снизить порог первого контакта, каждое новое issue/PR — событие
в ленте подписчиков. Использует только ресурсы GitHub.

1. Создай .github/ISSUE_TEMPLATE/bug_report.yml (GitHub Form):
   поля: что произошло, шаги воспроизведения, ожидание/реальность,
   окружение (OS, Go version, caf --version), логи.
2. Создай .github/ISSUE_TEMPLATE/feature_request.yml:
   поля: проблема, предлагаемое решение, альтернативы.
3. Создай .github/ISSUE_TEMPLATE/config.yml:
   blank_issues_enabled: false
   contact_links:
     - name: 💬 Discussions
       url: https://github.com/Reider85/OryaCodeAuthorityFirst/discussions
       about: Questions, ideas, RFC — welcome in Discussions.
4. Создай .github/PULL_REQUEST_TEMPLATE.md:
   разделы: Summary, Motivation, Changes, Type (feat/fix/docs/refactor/
   perf/test/chore), Checklist (tests added, changelog updated,
   no breaking change).
5. Создай .github/FUNDING.yml:
   github: [Reider85]
   (появится кнопка "Sponsor", +trust signal для devtool)
6. Через GitHub UI включи Discussions: Settings → General → Features →
   Set up discussions. Создай 3 категории: Q&A, Ideas, Show & tell.

7. Коммит: `feat(meta): issue/PR templates, funding, discussions`.

Критерий приёмки: при создании нового issue доступен выбор между
Bug report / Feature request. Discussions включён, видны категории.
Кнопка Sponsor появилась на странице репо.
```

### Промпт 12 (P2) — Conventional Commits + auto changelog

```
Цель: все коммиты в едином формате feat:/fix:/docs:/chore:/refactor:/
perf:/test:. Это позволяет goreleaser автоматически собирать
группированный changelog для каждого релиза (см. Промпт 5) и
включает авто-определение «breaking change» через `feat!: ...`.

1. Создай .github/workflows/conventional-commits.yml:

name: Conventional Commits

on:
  pull_request:
    types: [opened, edited, reopened, synchronize]

jobs:
  validate:
    name: Validate commit messages
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: amannn/action-semantic-pull-request@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          types: |
            feat
            fix
            docs
            style
            refactor
            perf
            test
            build
            ci
            chore
            revert
          requireScope: false

2. Создай COMMIT_TEMPLATE.md в корне:

   # Conventional Commits
   <type>(<optional scope>): <subject>

   <optional body>

   <optional footer with BREAKING CHANGE:>

   Допустимые type: feat, fix, docs, style, refactor, perf, test,
   build, ci, chore, revert.

3. Коммит: `chore(meta): enforce conventional commits via PR check`.

Критерий приёмки: PR с заголовком «update stuff» блокируется
semantic-pull-request экшеном, PR с заголовком «feat: add X»
проходит. В будущих релизах goreleaser формирует сгруппированный
changelog автоматически.
```

### Промпт 13 (P3) — GitHub Pages для docs и perf-dashboard

```
Цель: бесплатный хостинг документации и performance-dashboard на
ресурсах GitHub Pages. После каждого push в main обновляется
страница — это ещё один видимый «activity».

1. Создай ветку docs/ или папку docs/ с простым mkdocs или
   static-html сайтом (минимум: index.md с кратким обзором caf,
   install.md, usage.md, architecture.md).
2. Создай .github/workflows/pages.yml:

name: Pages

on:
  push:
    branches: [ main ]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Generate docs from README + markdown files
        run: |
          mkdir -p _site
          cp docs/*.md _site/ 2>/dev/null || true
          cp README.md _site/index.md
      - uses: actions/upload-pages-artifact@v3
        with:
          path: _site

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4

3. В Settings → Pages выбери Source = GitHub Actions.
4. В README добавь ссылку на готовую страницу:
   📖 **Docs**: https://reider85.github.io/OryaCodeAuthorityFirst/

5. Коммит: `docs(pages): publish README + docs to GitHub Pages`.

Критерий приёмки: после push в main в Actions отрабатывает pages.yml,
по URL https://reider85.github.io/OryaCodeAuthorityFirstFirst/
открывается опубликованная страница.
```

---

## 7. Приоритет внедрения

| # | Промпт | Эффект | Сложность |
|---|---|---|---|
| 1 | P1 — починить `branches: ain]` | CI из красного в зелёное за 1 мерж | 5 минут |
| 2 | P2 — обновить actions до v4/v5 | Нет deprecated-предупреждений | 10 минут |
| 3 | P3 — отключить шаги с несуществующими скриптами | performance-enhanced-ci.yml зелёный | 15 минут |
| 4 | P6 — topics + description + social preview | +бесплатный трафик с topic-страниц | 20 минут |
| 5 | P4 — auto-prerelease на push в main | Лента релизов обновляется каждый push | 30 минут |
| 6 | P5 — goreleaser на semver-тег | «Большие» релизы с бинарниками и changelog | 45 минут |
| 7 | P7 — golangci-lint + govulncheck + gosec | 3 бейджа доверия | 30 минут |
| 8 | P8 — build matrix на 5 платформах | Доказательство кросс-платформенности | 30 минут |
| 9 | P9 — Dependabot | Постоянная активность 1–3 PR/неделю | 10 минут |
| 10 | P10 — Codecov | Бейдж покрытия, обновляется каждый push | 20 минут |
| 11 | P11 — Issue/PR templates + Discussions | Снижение порога входа | 30 минут |
| 12 | P12 — Conventional Commits | Авто-changelog в релизах | 15 минут |
| 13 | P13 — GitHub Pages | Бесплатная документация | 30 минут |

**Если внедрять по одному промпту в день**, за 2 недели репо превращается
из «1 звезда, 38 failed CI» в «зелёный CI, 5 бейджей, авто-релизы,
Dependabot-активность, 20 topics, кросс-платформенные бинарники».
Это ровно те сигналы, по которым разработчики ставят звёзды devtool-проектам.

---

## 8. Чек-лист «после применения»

- [ ] В `ci.yml` и `performance-enhanced-ci.yml` строка `branches: [ main ]`.
- [ ] Везде `actions/upload-artifact@v4`, `download-artifact@v4`, `setup-go@v5`.
- [ ] `setup-go` использует `go-version-file: go.mod`, нигде не захардкожен Go 1.21.
- [ ] В `performance-enhanced-ci.yml` нет шагов, ссылающихся на `./scripts/analyze-performance-trend` или `./scripts/generate-dashboard`, пока они не реализованы.
- [ ] Создан `auto-prerelease.yml`, после каждого push в main в Releases появляется prerelease.
- [ ] Создан `release.yml` + `.goreleaser.yml`, на `git tag v0.5.0` появляется полноценный релиз.
- [ ] В README 5+ бейджей: CI, golangci-lint, govulncheck, gosec, coverage, Go Report Card.
- [ ] Заполнены 20 topics, описание, загружен social preview.
- [ ] Включён Dependabot (gomod + github-actions).
- [ ] Включены Discussions, есть issue/PR-шаблоны.
- [ ] Включён GitHub Pages, по URL открывается опубликованная страница.

---

## 9. Что важно не делать

- ❌ **Не создавать настоящий release на каждый коммит.** Это спам, GitHub
  trending помечает такие репо как gaming, подписчики отписываются.
  Правильно — prerelease на push (Промпт 4) + real release на semver-тег
  (Промпт 5).
- ❌ **Не оставлять `actions/upload-artifact@v3`** — он retired, в v4 другой
  API, на больших объёмах артефактов будет терять данные.
- ❌ **Не захардкоживать Go-версию в workflow** — всегда `go-version-file:
  go.mod`, иначе при апгрейде Go CI падает.
- ❌ **Не включать `performance-gate` и зависимые джобы**, пока не
  реализованы `./scripts/analyze-performance-trend`, `./scripts/generate-dashboard`
  и флаги `caf --performance-test` / `perf-test --comprehensive`. Сначала
  имплементация, потом включатель.
- ❌ **Не рассчитывать «попасть в Trending» за счёт одного только CI.**
  Звёзды приносят: topics + README-бейджи + кросс-платформенные бинарники +
  активность Dependabot + живые issue/Discussions. CI — это фундамент
  доверия, а не источник звёзд сам по себе.

---

_Файл `GITHUB.md` создан как инженерный playbook для репозитория
Reider85/OryaCodeAuthorityFirst. Применяй промпты по порядку P0 → P1 →
P2 → P3, каждый — отдельным PR с conventional-сообщением._
